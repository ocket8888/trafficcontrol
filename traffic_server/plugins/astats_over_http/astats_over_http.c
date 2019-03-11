/**
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

	  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
 */


#include <stdio.h>
#include <stdlib.h>
#include <ctype.h>
#include <limits.h>
#include <ts/ts.h>
#include <string.h>
#include <stdbool.h>
#include <sys/stat.h>
#include <time.h>

#include <inttypes.h>
#include <sys/types.h>
#include <dirent.h>

#include <unistd.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <errno.h>

#define STR_BUFFER_SIZE      (0x1000)
#define FREE_TMOUT           (300000)
#define SYSTEM_RECORD_TYPE   (0x100)
#define DEFAULT_RECORD_TYPES (SYSTEM_RECORD_TYPE | TS_RECORDTYPE_PROCESS | TS_RECORDTYPE_PLUGIN)
#define DEFAULT_PATH         "_astats"
#define PLUGIN_TAG           "astats_over_http"
#define DEFAULT_CONFIG_NAME  "astats.config"
#define DEFAULT_IP           "127.0.0.1"
#define DEFAULT_IP6          "::1"
#define PATH_FIELD           "path="
#define RECORD_FIELD         "record_types="
#define IP_FIELD             "allow_ip="
#define IP6_FIELD            "allow_ip6="

#ifdef __cplusplus
extern "C" {
#endif

// Represents an IPv4 CIDR - e.g. 127.0.0.1/32
typedef struct {
	in_addr addr;
	uint8_t netmask;
} ipv4;

// Represents an IPv6 CIDR - e.g. ::1/128
typedef struct {
	in6_addr addr;
	uint16_t netmask;
} ipv6;

// Stores the information from a configuration file
typedef struct {
	uint8_t recordTypes;
	char* stats_path;
	size_t stats_path_len;
	ipv4* allowIps;
	size_t ipCount;
	ipv6* allowIps6;
	size_t ip6Count;
} config_t;

// Contains a configuration structure as well as metadata pertaining to it
typedef struct {
	char* config_path;
	volatile time_t last_load;
	config_t* config;
} config_holder_t;

typedef struct stats_state_t {
	TSVConn net_vc;
	TSVIO read_vio;
	TSVIO write_vio;

	TSIOBuffer req_buffer;
	TSIOBuffer resp_buffer;
	TSIOBufferReader resp_reader;

	int output_bytes;
	int body_written;

	int globals_cnt;
	char** globals;
	char* interfaceName;
	char* query;
	unsigned int recordTypes;
} stats_state;

static int free_handler(TSCont cont, TSEvent event, void *edata);
static int config_handler(TSCont cont, TSEvent event, void *edata);
static config_t* get_config(TSCont cont);
static config_holder_t* new_config_holder(const char* path);
static const char RESP_HEADER[] = "HTTP/1.1 200 OK\r\n"
                                  "Content-Type: application/json\r\n"
                                  "Cache-Control: no-cache\r\n\r\n";

static size_t PATH_FIELD_LEN = strlen(PATH_FIELD);
static size_t RECORD_FIELD_LEN = strlen(RECORD_FIELD);
static size_t IP_FIELD_LEN = strlen(IP_FIELD);
static size_t IP6_FIELD_LEN = strlen(IP6_FIELD);
static size_t DEFAULT_PATH_LEN = strlen(DEFAULT_PATH);



unsigned int configReloadRequests = 0;
unsigned int configReloads = 0;
time_t lastReloadRequest = 0;
time_t lastReload = 0;
time_t astatsLoad = 0;

static bool is_ip_allowed(const config_t* config, const struct sockaddr* addr);

static char* nstr(const char *s) {
	char* mys = (char*) TSmalloc(strlen(s)+1);
	strcpy(mys, s);
	return mys;
}

static char* nstrl(const char* s, int len) {
	char* mys = (char*) TSmalloc(len + 1);
	memcpy(mys, s, len);
	mys[len] = 0;
	return mys;
}

static char** parseGlobals(char* str, int* globals_cnt) {
	char* tok = 0;
	char** globals = 0;
	char** old = 0;
	unsigned int globals_size = 0, cnt = 0;

	while (true) {
		tok = strtok_r(str, ";", &str);
		if (!tok) {
			break;
		}
		if (cnt >= globals_size) {
			old = globals;
			globals = (char **) TSmalloc(sizeof(char*) * (globals_size + 20));
			if (old) {
				memcpy(globals, old, sizeof(char*) * (globals_size));
				TSfree(old);
				old = NULL;
			}
			globals_size += 20;
		}
		globals[cnt] = tok;
		cnt++;
	}
	*globals_cnt = cnt;

	unsigned int i;
	for (i = 0; i < cnt; i++) {
		TSDebug(PLUGIN_TAG, "globals[%d]: '%s'", i, globals[i]);
	}

	return globals;
}

static void stats_fillState(stats_state* my_state, char* query, int query_len) {
	char* arg = 0;

	while (true) {
		arg = strtok_r(query, "&", &query);
		if (!arg)
			break;
		if (strstr(arg, "application=")) {
			arg = arg + strlen("application=");
			my_state->globals = parseGlobals(arg, &my_state->globals_cnt);
		} else if (strstr(arg, "inf.name=")) {
			my_state->interfaceName = arg + strlen("inf.name=");
		} else if(strstr(arg, "record.types=")) {
			my_state->recordTypes = strtol(arg + strlen("record.types="), NULL, 16);
		}
	}
}

static void stats_cleanup(TSCont contp, stats_state* my_state) {
	if (my_state->req_buffer) {
		TSIOBufferDestroy(my_state->req_buffer);
		my_state->req_buffer = NULL;
	}

	if (my_state->resp_buffer) {
		TSIOBufferDestroy(my_state->resp_buffer);
		my_state->resp_buffer = NULL;
	}

	TSVConnClose(my_state->net_vc);
	TSfree(my_state);
	my_state = NULL;
	TSContDestroy(contp);
}

static void stats_process_accept(TSCont contp, stats_state* my_state) {
	my_state->req_buffer = TSIOBufferCreate();
	my_state->resp_buffer = TSIOBufferCreate();
	my_state->resp_reader = TSIOBufferReaderAlloc(my_state->resp_buffer);
	my_state->read_vio = TSVConnRead(my_state->net_vc, contp, my_state->req_buffer, INT64_MAX);
}

static int stats_add_data_to_resp_buffer(const char* s, stats_state* my_state) {
	int s_len = strlen(s);

	TSIOBufferWrite(my_state->resp_buffer, s, s_len);

	return s_len;
}


static void stats_process_read(TSCont contp, TSEvent event, stats_state* my_state) {
	TSDebug(PLUGIN_TAG, "stats_process_read(%d)", event);

	switch (event) {
		case TS_EVENT_VCONN_READ_READY:
			my_state->output_bytes = stats_add_data_to_resp_buffer(RESP_HEADER, my_state);
			TSVConnShutdown(my_state->net_vc, 1, 0);
			my_state->write_vio = TSVConnWrite(my_state->net_vc, contp, my_state->resp_reader, INT64_MAX);
			break;
		case TS_EVENT_ERROR:
			TSError("stats_process_read: Received TS_EVENT_ERROR\n");
			break;
		case TS_EVENT_VCONN_EOS:
			break;
		case TS_EVENT_NET_ACCEPT_FAILED:
			TSError("stats_process_read: Received TS_EVENT_NET_ACCEPT_FAILED\n");
			break;
		default:
			printf("Unexpected Event %d\n", event);
			TSReleaseAssert(!"Unexpected Event");
			break;
	}
}

#define APPEND(a) my_state->output_bytes += stats_add_data_to_resp_buffer(a, my_state)
#define APPEND_STAT(a, fmt, v) do { \
		char b[3048]; \
		if (snprintf(b, sizeof(b), "   \"%s\": " fmt ",\n", a, v) < sizeof(b)) \
		APPEND(b); \
} while(0);

static void json_out_stat(TSRecordType rec_type,
						  void* edata,
						  int registered,
						  const char* name,
						  TSRecordDataType data_type,
						  TSRecordData* datum) {
	stats_state* my_state = edata;

	if (my_state->globals_cnt) {
		bool found = false;
		unsigned int i;
		for (i = 0; i < my_state->globals_cnt; i++) {
			if (strstr(name, my_state->globals[i])) {
				found = true;
				break;
			}
		}

		if (!found) {
			return; // skip
		}
	}

	switch(data_type) {
		case TS_RECORDDATATYPE_COUNTER:
			APPEND_STAT(name, "%" PRIu64, datum->rec_counter);
			break;
		case TS_RECORDDATATYPE_INT:
			APPEND_STAT(name, "%" PRIu64, datum->rec_int);
			break;
		case TS_RECORDDATATYPE_FLOAT:
			APPEND_STAT(name, "%f", datum->rec_float);
			break;
		case TS_RECORDDATATYPE_STRING:
			APPEND_STAT(name, "\"%s\"", datum->rec_string);
			break;
		default:
			TSDebug(PLUGIN_TAG, "unkown type for %s: %d", name, data_type);
			break;
	}
}

static char* getFile(char* filename, char* buffer, int bufferSize) {
	TSFile f = 0;
	size_t s = 0;

	f = TSfopen(filename, "r");
	if (!f)
	{
		buffer[0] = 0;
		return buffer;
	}

	s = TSfread(f, buffer, bufferSize);
	if (s > 0)
		buffer[s] = 0;
	else
		buffer[0] = 0;

	TSfclose(f);

	return buffer;
}

static int getSpeed(char* inf, char* buffer, int bufferSize) {
	char* str;
	char b[256];
	int speed = 0;

	snprintf(b, sizeof(b), "/sys/class/net/%s/operstate", inf);
	str = getFile(b, buffer, bufferSize);
	if (str && strstr(str, "up"))
	{
		snprintf(b, sizeof(b), "/sys/class/net/%s/speed", inf);
		str = getFile(b, buffer, bufferSize);
		speed = strtol(str, 0, 10);
	}

	return speed;
}

static void appendSystemState(stats_state* my_state) {
	char* interface = my_state->interfaceName;
	char buffer[16384];
	int bsize = 16384;
	char* str;
	char* end;
	int speed = 0;

	APPEND_STAT("inf.name", "\"%s\"", interface);

	speed = getSpeed(interface, buffer, sizeof(buffer));

	APPEND_STAT("inf.speed", "%d", speed);

	str = getFile("/proc/net/dev", buffer, sizeof(buffer));
	if (str && interface) {
		str = strstr(str, interface);
		if (str) {
			end = strstr(str, "\n");
			if (end) {
				*end = 0;
			}
			APPEND_STAT("proc.net.dev", "\"%s\"", str);
		}
	}

	str = getFile("/proc/loadavg", buffer, sizeof(buffer));
	if (str) {
		end = strstr(str, "\n");
		if (end) {
			*end = 0;
		}
		APPEND_STAT("proc.loadavg", "\"%s\"", str);
	}
}

static void json_out_stats(stats_state* my_state) {
	const char* version;
	TSDebug(PLUGIN_TAG, "recordTypes: '0x%x'", my_state->recordTypes);
	APPEND("{ \"ats\": {\n");
	TSRecordDump(my_state->recordTypes, json_out_stat, my_state);
	version = TSTrafficServerVersionGet();
	APPEND("   \"server\": \"");
	APPEND(version);
	APPEND("\"\n");
	APPEND("  }");

	if (my_state->recordTypes & SYSTEM_RECORD_TYPE) {
		APPEND(",\n \"system\": {\n");
		appendSystemState(my_state);
		APPEND_STAT("configReloadRequests", "%d", configReloadRequests);
		APPEND_STAT("lastReloadRequest", "%" PRIu64, lastReloadRequest);
		APPEND_STAT("configReloads", "%d", configReloads);
		APPEND_STAT("lastReload", "%" PRIu64, lastReload);
		APPEND_STAT("astatsLoad", "%" PRIu64, astatsLoad);
		APPEND("\"something\": \"here\"");
		APPEND("\n  }");
	}

	APPEND("\n}\n");
}

static void stats_process_write(TSCont contp, TSEvent event, stats_state* my_state) {
	if (event == TS_EVENT_VCONN_WRITE_READY) {
		if (my_state->body_written == 0) {
			TSDebug(PLUGIN_TAG, "plugin adding response body");
			my_state->body_written = 1;
			json_out_stats(my_state);
			TSVIONBytesSet(my_state->write_vio, my_state->output_bytes);
		}
		TSVIOReenable(my_state->write_vio);
		TSfree(my_state->globals);
		my_state->globals = NULL;
		TSfree(my_state->query);
		my_state->query = NULL;
	} else if (TS_EVENT_VCONN_WRITE_COMPLETE) {
		stats_cleanup(contp, my_state);
	}
	else if (event == TS_EVENT_ERROR) {
		TSError("stats_process_write: Received TS_EVENT_ERROR\n");
	}
	else {
		TSReleaseAssert(!"Unexpected Event");
	}
}

static int stats_dostuff(TSCont contp, TSEvent event, void* edata) {
	stats_state* my_state = TSContDataGet(contp);
	if (event == TS_EVENT_NET_ACCEPT) {
		my_state->net_vc = (TSVConn) edata;
		stats_process_accept(contp, my_state);
	} else if (edata == my_state->read_vio) {
		stats_process_read(contp, event, my_state);
	}
	else if (edata == my_state->write_vio) {
		stats_process_write(contp, event, my_state);
	}
	else {
		TSReleaseAssert(!"Unexpected Event");
	}

	return 0;
}

static int astats_origin(TSCont cont, TSEvent event, void* edata) {
	TSCont icontp;
	stats_state *my_state;
	TSMBuffer reqp;
	TSMLoc hdr_loc = NULL, url_loc = NULL;
	TSEvent reenable = TS_EVENT_HTTP_CONTINUE;
	config_t* config = get_config(cont);

	TSHttpTxn txnp = (TSHttpTxn) edata;
	TSDebug(PLUGIN_TAG, "in the read stuff");

	if (TSHttpTxnClientReqGet(txnp, &reqp, &hdr_loc) != TS_SUCCESS)
		goto cleanup;

	if (TSHttpHdrUrlGet(reqp, hdr_loc, &url_loc) != TS_SUCCESS)
		goto cleanup;

	int path_len = 0;
	const char* path = TSUrlPathGet(reqp,url_loc,&path_len);
	TSDebug(PLUGIN_TAG,"Path: %.*s",path_len,path);
	TSDebug(PLUGIN_TAG,"Path: %.*s",path_len,path);

	if (!(path_len == config->stats_path_len && !memcmp(path, config->stats_path, config->stats_path_len))) {
//		TSDebug(PLUGIN_TAG, "not right path: %.*s",path_len,path);
		goto notforme;
	}

	const struct sockaddr *addr = TSHttpTxnClientAddrGet(txnp);
	if(!is_ip_allowed(config, addr)) {
		TSDebug(PLUGIN_TAG, "not right ip");
		goto notforme;
	}
//	TSDebug(PLUGIN_TAG,"Path...: %.*s",path_len,path);

	int query_len;
	char* query = (char*)TSUrlHttpQueryGet(reqp,url_loc,&query_len);
	TSDebug(PLUGIN_TAG,"query: %.*s",query_len,query);

	TSSkipRemappingSet(txnp,1); //not strictly necessary, but speed is everything these days

	/* This is us -- register our intercept */
	TSDebug(PLUGIN_TAG, "Intercepting request");

	icontp = TSContCreate(stats_dostuff, TSMutexCreate());
	my_state = (stats_state *) TSmalloc(sizeof(*my_state));
	memset(my_state, 0, sizeof(*my_state));

	my_state->recordTypes = config->recordTypes;
	if (query_len) {
		my_state->query = nstrl(query, query_len);
		TSDebug(PLUGIN_TAG,"new query: %s", my_state->query);
		stats_fillState(my_state, my_state->query, query_len);
	}

	TSContDataSet(icontp, my_state);
	TSHttpTxnIntercept(icontp, txnp);

	goto cleanup;

	notforme:

	cleanup:
#if (TS_VERSION_NUMBER < 2001005)
	if (path) {
		TSHandleStringRelease(reqp, url_loc, path);
	}
#endif
	if (url_loc) {
		TSHandleMLocRelease(reqp, hdr_loc, url_loc);
	}
	if (hdr_loc) {
		TSHandleMLocRelease(reqp, TS_NULL_MLOC, hdr_loc);
	}

	TSHttpTxnReenable(txnp, reenable);

	return 0;
}


/*
 * Plug-in entry point. This will handle plug-in registration and set up the necessary HTTP hooks,
 * initial data, and handlers to run the plug-in.
*/
void TSPluginInit(int argc, const char* argv[]) {
	TSPluginRegistrationInfo info;
	info.plugin_name = PLUGIN_TAG;
	info.vendor_name = "Apache";
	info.support_email = "dev@trafficcontrol.apache.org";

	astatsLoad = time(NULL);

	#if (TS_VERSION_NUMBER < 3000000)
	if (TSPluginRegister(TS_SDK_VERSION_2_0, &info) != TS_SUCCESS) {
	#elif (TS_VERSION_NUMBER < 6000000)
	if (TSPluginRegister(TS_SDK_VERSION_3_0, &info) != TS_SUCCESS) {
	#else
	if (TSPluginRegister(&info) != TS_SUCCESS) {
	#endif
		TSError("Plugin registration failed. \n");
		return;
	}

	config_holder_t* config_holder = new_config_holder(argc > 1 ? argv[1] : NULL);
	if (config_holder == NULL) {
		TSError("Plug-in initialization failed.\n");
		return;
	}

	TSCont main_cont = TSContCreate(astats_origin, NULL);
	TSContDataSet(main_cont, (void*) config_holder);
	TSHttpHookAdd(TS_HTTP_READ_REQUEST_HDR_HOOK, main_cont);

	TSCont config_cont = TSContCreate(config_handler, TSMutexCreate());
	TSContDataSet(config_cont, (void*) config_holder);
	TSMgmtUpdateRegister(config_cont, PLUGIN_TAG);
	/* Create a continuation with a mutex as there is a shared global structure
	   containing the headers to add */
	TSDebug(PLUGIN_TAG, "astats module registered, path: '%s'", config_holder->config->stats_path);
}

static bool is_ip_match(const char* ip, char* ipmask, char mask) {
	unsigned int i, j, k;
	// to be able to set mask to 128
	unsigned int umask = 0xff & mask;

	for(j=0, i=0; ((i+1)*8) <= umask; i++) {
		if(ip[i] != ipmask[i]) {
			return false;
		}
		j+=8;
	}
	char cm = 0;
	for(k=0; j<umask;j++,k++) {
		cm |= 1<<(7-k);
	}

	if((ip[i]&cm) != (ipmask[i]&cm)) {
		return false;
	}
	return true;
}

static bool is_ip_allowed(const config_t* config, const struct sockaddr* addr) {
	char ip_port_text_buffer[INET6_ADDRSTRLEN];
	int i;
	char* ipmask;
	if(!addr) {
		return true;
	}

	if (addr->sa_family == AF_INET && config->allowIps) {
		const struct sockaddr_in* addr_in = (struct sockaddr_in*) addr;
		const char *ip = (char*) &addr_in->sin_addr;

		for(i=0; i < config->ipCount; i++) {
			ipmask = config->allowIps + (i*(sizeof(struct in_addr) + 1));
			if(is_ip_match(ip, ipmask, ipmask[4])) {
				TSDebug(PLUGIN_TAG, "clientip is %s--> ALLOW", inet_ntop(AF_INET,ip,ip_port_text_buffer,INET6_ADDRSTRLEN));
				return true;
			}
		}
		TSDebug(PLUGIN_TAG, "clientip is %s--> DENY", inet_ntop(AF_INET,ip,ip_port_text_buffer,INET6_ADDRSTRLEN));
		return false;

	} else if (addr->sa_family == AF_INET6 && config->allowIps6) {
		const struct sockaddr_in6* addr_in6 = (struct sockaddr_in6*) addr;
		const char* ip = (char*) &addr_in6->sin6_addr;

		for(i=0; i < config->ip6Count; i++) {
			ipmask = config->allowIps6 + (i*(sizeof(struct in6_addr) + 1));
			if(is_ip_match(ip, ipmask, ipmask[sizeof(struct in6_addr)])) {
				TSDebug(PLUGIN_TAG, "clientip6 is %s--> ALLOW", inet_ntop( AF_INET6,ip,ip_port_text_buffer,INET6_ADDRSTRLEN));
				return true;
			}
		}
		TSDebug(PLUGIN_TAG, "clientip6 is %s--> DENY", inet_ntop( AF_INET6,ip,ip_port_text_buffer,INET6_ADDRSTRLEN));
		return false;
	}
	return true;
}

/*
 * Copies `n` ipv4 structures from `buf` into the 'allowIps' field of `config`.
 * This does NOT de-allocate `buf`'s memory, but WILL expand/allocate the memory pointed to by the
 * 'allowIps' field of `config`.
 * Will also update the 'ipCount' field of `config` to accurately represent the number of stored IPs
*/
static void copyIPv4InPlace(config_t* config, ipv4* buf, size_t n) {
	ipv4* dest = config->allowIps;
	const size_t totalSize = n * sizeof(ipv4);
	if (config->allowIps == NULL) {
		config->allowIps = TSmalloc(totalSize);
		dest = config->allowIps;
	} else {
		TSrealloc(config->allowIps, totalSize + config->ipCount*sizeof(ipv4));
		dest += config->ipCount * sizeof(ipv4);
	}
	memcpy(dest, buf, totalSize);
	config->ipCount += n;
}

/*
 * Copies `n` ipv6 structures from `buf` into the 'allowIps6' field of `config`.
 * This does NOT de-allocate `buf`'s memory, but WILL expand/allocate the memory pointed to by the
 * 'allowIps' field of `config`.
 * Will also update the 'ipCount' field of `config` to accurately represent the number of stored IPs
*/
static void copyIPv6InPlace(config_t* config, ipv4* buf, size_t n) {
	ipv6* dest = config->allowIps6;
	const size_t totalSize = n*sizeof(ipv6);
	if (config->allowIps6 == NULL) {
		config->allowIps6 = TSmalloc(totalSize);
		dest = config->allowIps6;
	} else {
		TSrealloc(config->allowIps6, totalSize + config->ip6Count*sizeof(ipv6));
		dest += config->ip6Count * sizeof(ipv6);
	}
	memcpy(dest, buf, totalSize);
	config->ip6Count += n;
}

/*
 * Parses the passed string as an IPv4 address, optionally with a CIDR mask. The constructed ipv4
 * structure is stored into `dest`.
 * Returns `true` on success, `false` on parse error. The state of the structure stored at `dest` is
 * undefined on error (but likely is left untouched) as it depends on the implementation of
 * inet_pton(3) as well as strtoul(3).
*/
static bool parseIPv4(const char* ipStr, ipv4* dest) {
	char* p = NULL;
	if (p=strstr(ipStr, "/")) {
		*p = '\0';
		p += 1;
	}

	if (inet_pton(AF_INET, ipStr, &(dest->addr)) != 1) {
		return false;
	}

	if (p != NULL) {
		*(p-1) = '/';
		char* error_position = p;
		dest->netmask = (uint8_t)strtoul(p, error_position, 10);
		if (*error_position != '\0') {
			TSError("[%s] '%s' could not be parsed as a valid CIDR netmask!\n", PLUGIN_TAG, p);
			TSDebug(PLUGIN_TAG, "error encountered starting here: '%s'", error_position);
			return false;
		} else if (dest->netmask > 32) {
			TSError("[%s] Netmask found to be %u, should be at most 32!\n", PLUGIN_TAG, dest->netmask);
			return false;
		}
	} else {
		dest->netmask = 32;
	}
	return true;
}

/*
 * Parses the passed string as an IPv6 address, optionally with a CIDR mask. The constructed ipv6
 * structure is stored into `dest`.
 * Returns `true` on success, `false` on parse error. The state of the structure stored at `dest` is
 * undefined on error (but likely is left untouched) as it depends on the implementation of
 * inet_pton(3) as well as strtoul(3).
*/
static bool parseIPv6(const char* ipStr, ipv6* dest) {
	char* p = NULL;
	if (p=strstr(ipStr, "/")) {
		*p = '\0';
		p += 1;
	}

	if (inet_pton(AF_INET6, ipStr, &(dest->addr)) != 1) {
		return false;
	}

	if (p != NULL) {
		*(p-1) = '/';
		char* error_position = p;
		dest->netmask = (uint16_t)strtoul(p, error_position, 10);
		if (*error_position != '\0') {
			TSError("[%s] '%s' could not be parsed as a valid CIDR netmask!\n", PLUGIN_TAG, p);
			TSDebug(PLUGIN_TAG, "error encountered starting here: '%s'", error_position);
			return false;
		} else if (dest->netmask > 128) {
			TSError("[%s] Netmask found to be %u, should be at most 32!\n", PLUGIN_TAG, dest->netmask);
			return false;
		}
	} else {
		dest->netmask = 128;
	}
	return true;
}

/*
 * Parses part of a configuration file for IPv4 addresses, appending them to the configuration
 * structures as they are parsed. `ipStr` is expected to point into the configuration file contents
 * to the point where a list of IP addresses begins. If `ipStr` is NULL (or points to NULL), this
 * will instead set it to use DEFAULT_IP. NOTE: this behavior overwrites existing stored IP
 * addresses!
 * Sets `ipStr` to point to to whatever was remaining when parsing finished (the idea being that the
 * rest of the data should be appended to a new string and then this function called again to
 * complete parsing all IP addresses). If `ipStr` was NULL, it will be unchanged.
 * Returns the number of succesfully parsed IP addresses, or `-1` on error.
*/
static int parseIps(config_t* config, char** ipStr) {

	char* pos = ipStr == NULL ? NULL : *ipStr;

	if (pos == NULL) {
		if (config->allowIps != NULL) {
			TSDebug(PLUGIN_TAG, "Overwriting existing IPv4 addresses");
			TSfree(config->allowIps);
		}
		config->ipCount = 1;
		config->allowIps = TSmalloc(sizeof(ipv4));
		inet_pton(AF_INET, DEFAULT_IP, config->allowIps);
		config->allowIps[0].netmask = 32;
		return 1;
	}

	// This is my best guess at a sensible buffer size; if every address in the string is minimum
	// possible length this is the maximum number of them that can fit in the passed string.
	const size_t buffsize = strlen(pos) / 8 + 1;
	ipv4* ipBuffer = (ipv4*)TSmalloc(sizeof(ipv4)*buffsize);

	char* anIP = strtok(pos, ", \n");
	size_t totalIPs = 0;
	do {
		size_t ipnum;
		for (ipnum=0; anIP != NULL && ipnum < buffsize; ++ipnum) {
			if (!parseIPv4(anIP, ipBuffer + ipnum*sizeof(ipv4))) {
				TSfree(ipBuffer);
				TSError("[%s] Couldn't parse '%s' as a valid IP address!\n", PLUGIN_TAG, anIP);
				return -1;
			}
			*ipStr = anIP;
			anIP = strtok(NULL, ", \n");
		}

		copyIPv4InPlace(config, ipBuffer, ipnum);
		totalIPs += ipnum;
	} while (anIP != NULL);

	TSfree(ipBuffer);
	return (int)totalIPs;

	// char buffer[STR_BUFFER_SIZE];
	// char *p, *tok1, *tok2, *ip;
	// int i, mask;
	// char ip_port_text_buffer[INET6_ADDRSTRLEN];

	// if(!ipStr) {
	// 	config->ip6Count = 1;
	// 	ip = config->allowIps6 = TSmalloc(sizeof(struct in6_addr) + 1);
	// 	inet_pton(AF_INET6, DEFAULT_IP6, ip);
	// 	ip[sizeof(struct in6_addr)] = 128;
	// 	return;
	// }

	// strcpy(buffer, ipStr);
	// p = buffer;
	// while(strtok_r(p, ", \n", &p)) {
	// 	config->ipCount++;
	// }
	// if(!config->ipCount) {
	// 	return;
	// }
	// config->allowIps = TSmalloc(5*config->ipCount); // 4 bytes for ip + 1 for bit mask
	// strcpy(buffer, ipStr);
	// p = buffer;
	// i = 0;
	// while((tok1 = strtok_r(p, ", \n", &p))) {
	// 	TSDebug(PLUGIN_TAG, "%d) parsing: %s", i+1,tok1);
	// 	tok2 = strtok_r(tok1, "/", &tok1);
	// 	ip = config->allowIps+((sizeof(struct in_addr) + 1)*i);
	// 	if(!inet_pton(AF_INET, tok2, ip)) {
	// 		TSDebug(PLUGIN_TAG, "%d) skipping: %s", i+1,tok1);
	// 		continue;
	// 	}
	// 	tok2 = strtok_r(tok1, "/", &tok1);
	// 	if(!tok2) {
	// 		mask = 32;
	// 	} else {
	// 		mask = atoi(tok2);
	// 	}
	// 	ip[4] = mask;
	// 	TSDebug(PLUGIN_TAG, "%d) adding netmask: %s/%d", i+1,
	// 			inet_ntop(AF_INET,ip,ip_port_text_buffer,INET_ADDRSTRLEN),ip[4]);
	// 	i++;
	// }
}

/*
 * Parses part of a configuration file for IPv6 addresses, appending them to the configuration
 * structures as they are parsed. `ipStr` is expected to point into the configuration file contents
 * to the point where a list of IP addresses begins. If `ipStr` is NULL (or points to NULL), this
 * will instead set it to use DEFAULT_IP. NOTE: this behavior overwrites existing stored IP
 * addresses!
 * Sets `ipStr` to point to to whatever was remaining when parsing finished (the idea being that the
 * rest of the data should be appended to a new string and then this function called again to
 * complete parsing all IP addresses). If `ipStr` was NULL, it will be unchanged.
 * Returns the number of succesfully parsed IP addresses, or `-1` on error.
*/
static void parseIps6(config_t* config, char** ipStr) {

	char* pos = ipStr == NULL ? NULL : *ipStr;

	if (pos == NULL) {
		if (config->allowIps6 != NULL) {
			TSDebug(PLUGIN_TAG, "Overwriting existing IPv6 addresses");
			TSfree(config->allowIps6);
		}
		config->ip6Count = 1;
		config->allowIps6 = TSmalloc(sizeof(ipv6));
		inet_pton(AF_INET, DEFAULT_IP, config->allowIps6);
		config->allowIps6[0].netmask = 128;
		return 1;
	}

	// This is my best guess at a sensible buffer size; IPv6 addresses support a lot of short-hand
	// notation, so I aimed for the average case by allocating enough memory for a number of IP
	// addresses that will fit in the input string assuming they're all half the maximum size.
	const size_t buffsize = strlen(pos) / 20 + 1;
	ipv6* ipBuffer = (ipv6*)TSmalloc(sizeof(ipv6)*buffsize);

	char* anIP = strtok(pos, ", \n");
	size_t totalIPs = 0;
	do {
		size_t ipnum;
		for (ipnum=0; anIP != NULL && ipnum < buffsize; ++ipnum) {
			if (!parseIPv6(anIP, ipBuffer + ipnum*sizeof(ipv6))) {
				TSfree(ipBuffer);
				TSError("[%s] Couldn't parse '%s' as a valid IP address!\n", PLUGIN_TAG, anIP);
				return -1;
			}
			*ipStr = anIP;
			anIP = strtok(NULL, ", \n");
		}

		copyIPv6InPlace(config, ipBuffer, ipnum);
		totalIPs += ipnum;
	} while (anIP != NULL);

	TSfree(ipBuffer);
	return (int)totalIPs;

	// char buffer[STR_BUFFER_SIZE];
	// char *p, *tok1, *tok2, *ip;
	// int i, mask;
	// char ip_port_text_buffer[INET6_ADDRSTRLEN];

	// if(!ipStr) {
	// 	config->ip6Count = 1;
	// 	ip = config->allowIps6 = TSmalloc(sizeof(struct in6_addr) + 1);
	// 	inet_pton(AF_INET6, DEFAULT_IP6, ip);
	// 	ip[sizeof(struct in6_addr)] = 128;
	// 	return;
	// }

	// strcpy(buffer, ipStr);
	// p = buffer;
	// while(strtok_r(p, ", \n", &p)) {
	// 	config->ip6Count++;
	// }
	// if(!config->ip6Count) {
	// 	return;
	// }

	// config->allowIps6 = TSmalloc((sizeof(struct in6_addr) + 1)*config->ip6Count); // 16 bytes for ip + 1 for bit mask
	// strcpy(buffer, ipStr);
	// p = buffer;
	// i = 0;
	// while((tok1 = strtok_r(p, ", \n", &p))) {
	// 	TSDebug(PLUGIN_TAG, "%d) parsing: %s", i+1,tok1);
	// 	tok2 = strtok_r(tok1, "/", &tok1);
	// 	ip = config->allowIps6+((sizeof(struct in6_addr)+1)*i);
	// 	if(!inet_pton(AF_INET6, tok2, ip)) {
	// 		TSDebug(PLUGIN_TAG, "%d) skipping: %s", i+1,tok1);
	// 		continue;
	// 	}
	// 	tok2 = strtok_r(tok1, "/", &tok1);
	// 	if(!tok2) {
	// 		mask = 128;
	// 	} else {
	// 		mask = atoi(tok2);
	// 	}
	// 	ip[sizeof(struct in6_addr)] = mask;
	// 	TSDebug(PLUGIN_TAG, "%d) adding netmask: %s/%d", i+1,
	// 			inet_ntop(AF_INET6,ip,ip_port_text_buffer,INET6_ADDRSTRLEN),ip[sizeof(struct in6_addr)]);
	// 	i++;
	// }
}

/*
 * Constructs a configuration structure from the contents of the file identified by `fh`.
 * If `fh` is NULL, constructs a configuration structure with the default options.
 * Returns NULL if the file could not be read/parsed or an error occured doing either.
*/
static config_t* new_config(TSFile fh) {
	config_t* config = (config_t*)TSmalloc(sizeof(config_t));
	config->stats_path = NULL;
	config->stats_path_len = 0;
	config->allowIps = NULL;
	config->ipCount = 0;
	config->allowIps6 = NULL;
	config->ip6Count = 0;
	config->recordTypes = DEFAULT_RECORD_TYPES;

	if (fh == NULL) {
		config->stats_path = (char*) TSmalloc(sizeof(char)*(DEFAULT_PATH_LEN+1));
		strcpy(config->stats_path, DEFAULT_PATH)
		config->stats_path_len = DEFAULT_PATH_LEN;

		parseIps(config, NULL);
		parseIps6(config, NULL);

		TSDebug(PLUGIN_TAG, "No config, using defaults");
		return config;
	}

	// Read in the configuration file. This will use the last definition of each field, discarding
	// previous ones as it goes. Note that comments MUST begin at the beginning of a line.
	char buffer[STR_BUFFER_SIZE];
	while (TSfgets(fh, buffer, STR_BUFFER_SIZE - 1)) {
		if (*buffer == '#') {
			continue;
		}

		char* tok;
		char* p = NULL;
		if (p = strstr(buffer, PATH_FIELD)) {
			p += PATH_FIELD_LEN;
			tok = strtok(p, " \n", &p);

			if (strlen(tok) < 1) {
				TSError("[%s] Invalid configuration file - path must have a value!\n", PLUGIN_TAG);
				delete_config(config);
				return NULL;
			}

			if (config->stats_path != NULL) {
				TSDebug(PLUGIN_TAG, "multiple 'path=' in config file");
				TSfree(config->stats_path);
			}

			config->stats_path = (char*) TSmalloc(sizeof(char) * (PATH_FIELD_LEN + 1));

			if (strcpy(config->stats_path, tok) == NULL) {
				TSDebug(PLUGIN_TAG, "new_config: strcpy failed reading path field");
				delete_config(config);
				return NULL;
			}

			config->stats_path_len = strlen(config->stats_path);


		} else if (p = strstr(buffer, RECORD_FIELD)) {
			p += RECORD_FIELD_LEN;
			tok = strtok(p, " \n", &p);

			if (strlen(tok) < 1) {
				TSError("[%s] Invalid configuration file - record_types (if present) must have a value!\n", PLUGIN_TAG);
				delete_config(config);
				return NULL;
			}

			char* error_position = NULL;
			config->recordTypes = (uint8_t)strtoul(tok, error_position, 16);

			if (*error_position != '\0') {
				delete_config(config);
				if (errno = ERANGE) {
					TSError("[%s] Invalid configuration file - record_types value '%s' out of range!\n", PLUGIN_TAG, tok);
				} else {
					TSError("[%s] Invalid configuration file - record_types value '%s' is not a hexidecimal integer!\n", PLUGIN_TAG);
					TSDebug(PLUGIN_TAG, "record_types value was '%s', error encountered starting here: '%s'", tok, error_position);
				}
				return NULL;
			}


		} else if (p = strstr(buffer, IP_FIELD)) {
			p += IP_FIELD_LEN;

			// To deal with extremely long lines we get more input in a loop until EOL/EOF
			do {
				if (parseIps(config, p) < 0) {
					TSError("[%s] Failed to parse '%s' value!\n", PLUGIN_TAG, IP_FIELD);
					delete_config(config);
					return NULL;
				}
				if (*p=='\n' || *p=='\0' || !TSfgets(fh, buffer, STR_BUFFER_SIZE - 1)) {
					break;
				}
				p = buffer;
			} while(true);

		} else if (p = strstr(buffer, IP6_FIELD)) {
			p += IP6_FIELD_LEN;

			// To deal with extremely long lines we get more input in a loop until EOL/EOF
			do {
				if (parseIps6(config, p) < 0) {
					TSError("[%s] Failed to parse '%s' value!\n", PLUGIN_TAG, IP6_FIELD);
					delete_config(config);
					return NULL;
				}
				if (*p=='\n' || *p=='\0' || !TSfgets(fh, buffer, STR_BUFFER_SIZE - 1)) {
					break;
				}
				p = buffer;
			} while(true);
		}
	}

	if (config->ipCount == 0 && parseIps(config, NULL) != 1) {
		delete_config(config);
		return NULL;
	}

	if (config->ip6Count == 0 && parseIps6(config, NULL) != 1) {
		delete_config(config);
		return NULL;
	}

	if (config->stats_path == NULL) {
		config->stats_path = (char*) TSmalloc(sizeof(char)*(DEFAULT_PATH_LEN+1));
		strcpy(config->stats_path, DEFAULT_PATH)
		config->stats_path_len = DEFAULT_PATH_LEN;
	}

	TSDebug(PLUGIN_TAG, "config path=%s", config->stats_path);

	return config;
}

/*
 * Safely de-allocates all dynamically allocated fields of `config`, and then de-allocates `config`
 * itself. ('Safely' meaning it won't attempt to de-allocate unallocated fields)
*/
static void delete_config(config_t* config) {
	TSDebug(PLUGIN_TAG, "Freeing config");
	if (config == NULL) {
		TSDebug(PLUGIN_TAG, "Config was null...");
		return;
	}
	if (config->allowIps != NULL) {
		TSfree(config->allowIps);
	}
	if (config->allowIPs6 != NULL) {
		TSfree(config->allowIPs6);
	}
	if (config->stats_path != NULL) {
		TSfree(config->stats_path);
	}
	TSfree(config);
}


// standard api below...

/*
 * Extracts a configuration structure from a Continuation. Returns NULL if no configuration holder
 * is found within the continuation data.
*/
static config_t* get_config(TSCont cont) {
	config_holder_t* configh = (config_holder_t *) TSContDataGet(cont);
	if (configh == NULL) {
		return NULL;
	}
	return configh->config;
}

// Uses the infromation contained in the passed `config_holder` to read in and parse a configuration
// file. Returns `true` if the load succeeded (or was unnecessary), `false` otherwise.
static bool load_config_file(config_holder_t* config_holder) {

	configReloadRequests++;
	lastReloadRequest = time(NULL);

	struct stat s;
	if (stat(config_holder->config_path, &s) < 0) {
		int err = errno;
		TSError("[%s] Could not stat configuration file\n", PLUGIN_TAG);
		TSDebug(PLUGIN_TAG, "Could not stat %s: %s", config_holder->config_path, strerror(err));

		// if there's an existing config, then we don't bail at this point.
		return config_holder->config != NULL;

	// don't read files that haven't changed since the last time we looked
	} else {
		TSDebug(PLUGIN_TAG, "s.st_mtime=%lu, last_load=%lu", s.st_mtime, config_holder->last_load);
		if (s.st_mtime < config_holder->last_load) {
			return true;
		}
	}

	TSDebug(PLUGIN_TAG, "Opening config file: %s", config_holder->config_path);
	TSFile fh = TSfopen(config_holder->config_path, "r");

	if (fh == NULL) {
		TSError("[%s] Unable to open config: %s.\n", PLUGIN_TAG, config_holder->config_path);

		// if there's an existing config, then we don't bail at this point
		return config_holder->config != NULL;
	}

	config_t* oldconfig;
	TSCont free_cont;

	config_t* newconfig = new_config(fh);
	if (newconfig == NULL) {
		return false;
	}

	configReloads++;
	lastReload = lastReloadRequest;
	config_holder->last_load = lastReloadRequest;
	config_t** confp = &(config_holder->config);
	oldconfig = __sync_lock_test_and_set(confp, newconfig);
	if (oldconfig) {
		TSDebug(PLUGIN_TAG, "scheduling free: %p (%p)", oldconfig, newconfig);
		free_cont = TSContCreate(free_handler, TSMutexCreate());
		TSContDataSet(free_cont, (void*) oldconfig);
		TSContSchedule(free_cont, FREE_TMOUT, TS_THREAD_POOL_TASK);
	}

	if (fh != NULL) {
		TSfclose(fh);
	}
	return true;
}

// Creates a new configuration structure for a configuration file located at `path`.
// If `path` is NULL, this will attempt to read in a file located at DEFAULT_CONFIG_NAME within the
// ATS configuration directory.
static config_holder_t* new_config_holder(const char* path) {
	config_holder_t* config_holder = TSmalloc(sizeof(config_holder_t));
	config_holder->config = NULL;
	config_holder->last_load = 0;

	if (path != NULL) {
		config_holder->config_path = (char*) TSmalloc(sizeof(char) * (strlen(path) + 1));
		if (strcpy(config_holder->config_path, path) == NULL) {
			TSError("[%s] Failed to initialize pathspec\n", PLUGIN_TAG);
			TSDebug(PLUGIN_TAG, "new_config_holder: pathspec was: %s", path);
			TSfree(config_holder->config_path);
			TSfree(config_holder);
			return NULL;
		}
	} else {
		char* cfgDir = TSConfigDirGet();
		char* fname = (char*) TSmalloc(sizeof(char) * (strlen(cfgDir) + 1));
		int n = sprintf(fname, "%s/"DEFAULT_CONFIG_NAME, cfgDir);
		if (n < 0) {
			TSError("[%s] Encoding error trying to set configuration file path\n", PLUGIN_TAG);
			TSDebug(PLUGIN_TAG, "new_config_holder: configuration dir was: %s", cfgDir);
			TSfree(fname);
			TSfree(config_holder);
			return NULL;
		}
		config_holder->config_path = fname;
	}

	if (load_config_file(config_holder)) {
		return config_holder;
	}
	TSfree(config_holder->config_path);
	TSfree(config_holder);
	return NULL;
}

/*
 * An asynchronous handler for deleting a configuration structure when it is no longer needed.
 * Always returns 0.
*/
static int free_handler(TSCont cont, TSEvent event, void* edata) {
	TSDebug(PLUGIN_TAG, "Freeing old config");
	config_t* config = (config_t*) TSContDataGet(cont);
	delete_config(config);
	TSContDestroy(cont);
	return 0;
}

static int config_handler(TSCont cont, TSEvent event, void* edata) {
	config_holder_t *config_holder;

	TSDebug(PLUGIN_TAG, "In config Handler");
	config_holder = (config_holder_t *) TSContDataGet(cont);
	load_config_file(config_holder);
	return 0;
}

#ifdef __cplusplus
};
#endif
