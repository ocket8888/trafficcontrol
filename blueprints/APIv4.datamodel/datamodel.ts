/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/**
 * This is my Typescript Data Model. I did it in Typescript because I love
 * Typescript, but also it was designed for modeling strict typing on
 * traditionally more generic data - which makes it perfect.
 * If you're reading this on my GitHub Pull Request, don't worry about this
 * file. It's just to help me better organize my thoughts and make sure
 * structures are programmatically consistent. Just focus on the LaTeX
 * documents, everything in here that matters will be properly written up in
 * them.
 * @note A lot of things have lastUpdated fields which I have chosen to omit
 * from the data model. Just a lot of copy/pasting.
 * @packageDocumentation
*/

// Just an alias to make things simpler.
export type int = bigint;

// export type uint = bigint;
// export function isUint(a: int | uint): a is uint {
// 	return a >= 0;
// }

/**
 * A CDN represents a full stack of hardware, components and configuration
 * required to delivery content.
*/
interface CDN {
	/** The numeric IDs of the CDN's cache servers */
	cacheServers: Set<int>;
	/** The names of all of the Delivery Services within the CDN */
	deliveryServices: Set<string>;
	/** Whether or not DNSSEC is enabled on the CDN's domain */
	dnssecEnabled: boolean;
	/** The CDN's domain. Must be a valid DNS label. */
	domain: string;
	/** The numeric IDs of the CDN's infrastructure servers */
	infrastructureServers: Set<int>;
	/** The Name of a CDN uniquely identifies it. */
	name: string;
	/** The numeric IDs of the CDN's origins */
	origins: Set<int>;
	/** The numeric IDs of the Traffic Monitor instances within the CDN */
	trafficMonitors: Set<int>;
	/** The numeric IDs of the Traffic Router instances within the CDN */
	trafficRouters: Set<int>;
	/**
	 * The numeric IDs of the Traffic Stats instances within the CDN.
	 * @todo Currently TS is a singleton; is there any desire/reason to change?
	*/
	trafficStatsServers: Set<int>;
}

/**
 * A Tag is just sort of a miscellaneous label applied to something, used to
 * group them.
 *
 * There's no way to restrict a tag to be capable of being assigned only to
 * certain object types.
*/
interface Tag {
	/**
	 * A short description of what the tag means.
	*/
	description: string;
	/**
	 * The name of the Tag (case-sensitive) which uniquely identifies it.
	*/
	name: string;
}

/**
 * Physical Locations are exactly what they sound like, a real, physical place
 * where component hardware is stored.
 *
 * Mainly this contains metadata for support purposes, but it can also be used
 * to provide some routing fallback behavior.
*/
interface PhysicalLocation {
	/**
	 * This is a text field with no defined structure, but semantically it
	 * represents a Physical Location's real-world, physical address.
	 *
	 * If it is not an empty value, it is assumed to contain enough information
	 * to send a letter through normal postage to the Physical Location (though
	 * the site may not actually be capable of receiving mail).
	 *
	 * Addresses may consist of alphanumeric characters, hyphens, periods,
	 * spaces, and newlines, but may neither begin nor end with a space or
	 * newline.
	*/
	address: string;
	/**
	 * An email address, which is used for communicating with the Physical
	 * Location's "point-of-contact".
	*/
	email: string;
	/** A uniquely identifying string that names the location. */
	name: string;
	/** Arbitrary text for miscellaneous purposes */
	notes: string;
	/**
	 * A string which, if not empty, is presumed to be a telephone number at
	 * which the Physical Location's "point-of-contact" may be contacted.
	 *
	 * It may only contain numerics and hyphens, and may neither start nor end
	 * with a hyphen.
	*/
	phoneNumber: string;
	/**
	 * The name of the person designated as the "point-of-contact" for this
	 * Physical Location.
	*/
	pointOfContactName: string;
}

/**
 * Conveys server status information.
 *
 * These aren't objects presented through the API in their own right. Instead
 * this is just here to specify and describe the static, unchanging allowed
 * values.
*/
enum Status {
	/**
	 * The Cache Server is considered unhealthy and its thresholds and
	 * connectivity state are not monitored. Its existence is not disclosed to
	 * Traffic Router(s).
	*/
	ADMIN_DOWN = 'ADMIN_DOWN',
	/**
	 * The Cache Server is considered unhealthy regardless of any thresholds or
	 * connectivity state (but it is still monitored).
	*/
	OFFLINE = 'OFFLINE',
	/**
	 * The Cache Server will always be considered healthy regardless of any
	 * thresholds or connectivity state (but it is still monitored).
	*/
	ONLINE = 'ONLINE',
	/**
	 * The Cache Server’s health is presented to the Traffic Router(s) as it is
	 * reported by its various thresholds, as determined by the Traffic
	 * Monitor(s).
	*/
	REPORTED = 'REPORTED'
}

/**
 * This is just a collection of the field common to **all** server types: cache,
 * infrastructure, Traffic Stats, etc.
 *
 * This isn't a "real" type in that it isn't exposed through the API.
*/
interface BaseServer {
	/** The Name of the CDN to which the server belongs. */
	cdn: string;
	/**
	 * The "domain" part of the server's
	 * <abbr title="Fully Qualified Domain Name">FQDN</abbr>
	*/
	domain: string;
	/** The server's hostname - *not* necessarily unique */
	hostName: string;
	/** Still numeric IDs because we want a single, unique identifier */
	id: int;
	/** arbitrary text for miscellaneous purposes */
	notes: string;
	/** The Name of the Physical Location in which the server resides */
	physicalLocation: string;
	/** The Names of Tags given to this server */
	tags: Set<string>;
}

/**
 * IPAddress represents a single IP address used by a single Interface of a
 * Cache Server.
 */
interface IPAddress {
	/**
	 * The actual IP (v4 or v6) address which is used by an interface.
	 * If it is an IPv6 address
	 */
	address: string;
	/**
	 * The IP (v4 or v6) address of the gateway used by this IP address.
	 */
	gateway: string | null;
	/**
	 * Tells whether or not this address of this interface is the server's
	 * "service" address.
	 * At least one address of EXACTLY ONE interface MUST have this set to
	 * 'true' for a server.
	 */
	serviceAddress: boolean;
}

/**
 * Interface is a network interface used by a Cache Server.
 */
interface Interface {
	/**
	 * These will be all of the IPv4/IPv6 addresses assigned to the interface,
	 * including gateways and "masks".
	 * It is illegal for an interface to not have at least one associated IP
	 * address.
	 */
	ipAddresses: Array<IPAddress> & {0: IPAddress};
	/**
	 * The maximum allowed bandwidth for this interface to be considered "healthy"
	 * by Traffic Monitor.
	 *
	 * This has no effect if `monitor` is not true.
	 * Values are in kb/s.
	 * The value `0` means "no limit".
	 */
	maxBandwidth: bigint;
	/**
	 * Whether or not Traffic Monitor should monitor this particular interface.
	 */
	monitor: boolean;
	/**
	 * The interface's Maximum Transmission Unit.
	 * If this is 'null' it is assumed that the interface's MTU is not known/is
	 * irrelevant.
	 */
	mtu: number | null;
	/**
	 * The name of the interface device on the server e.g. eth0.
	 */
	name: string;
}

/**
 * CacheServer is meant to represent specifically servers for caching content -
 * MIDs and EDGEs.
 *
 * This is splitting them out from the broader concept of a "server" in the
 * current data model.
*/
interface CacheServer extends BaseServer {
	/**
	 * Sets the maximum 'healthy' bandwidth in kbps.
	 *
	 * '0' means "no limit".
	*/
	bandwidthThreshold: int;
	/** The Name of the Cache Group to which the server belongs. */
	cacheGroup: string;
	/** What we today call "Server Capabilities" */
	capabilities: Set<string>;
	/** A set of filepaths to HDD block devices to use for caching content */
	/**
	 * The port on which the server listens for incoming HTTP requests.
	 *
	 * If it's 'null', the assumption is that the server doesn't do HTTP.
	*/
	httpPort: bigint;
	/**
	 * The port on which the server listens for incoming HTTPS requests.
	 *
	 * If it's 'null', the assumption is that the server doesn't do HTTPS.
	*/
	httpsPort: bigint;
	/** The Name of the Profile in use by the cache server */
	profile: string;
	/** Whether or not the server has pending revalidations */
	revalidationPending: boolean;
	/** Server's status */
	status: Status;
	/**
	 * The server's type.
	 *
	 * Unlike currently, it is only permitted to be EDGE (meaning that it's an
	 * edge-tier cache) or MID (meaning that it's a mid-tier cache).
	*/
	type: 'EDGE' | 'MID';
	/** Whether or not the server has pending updates */
	updatePending: boolean;
}

interface CacheGroup {
	name: string;
	latitude: number;
	longitude: number;
	cacheServers: Set<bigint>;
	type: 'EDGE' | 'MID';
}

/**
 * Infrastructures are arbitrary servers that aren't ATC components; think
 * INFLUX, GRAFANA, etc.
*/
interface InfrastructureServer {
	/** The Name of the CDN to which the server belongs, if there is one. */
	cdn: string | null;
	/**
	 * The "domain" part of the server's
	 * <abbr title="Fully Qualified Domain Name">FQDN</abbr>
	*/
	domain: string;
	/** The server's hostname - *not* necessarily unique */
	hostName: string;
	/** Still numeric IDs because we want a single, unique identifier */
	id: int;
	/** arbitrary text for miscellaneous purposes */
	notes: string;
	/** The Name of the Physical Location in which the server resides */
	physicalLocation: string | null;
	/** The Names of Tags given to this server */
	tags: Set<string>;
	/**
	 * The port on which the server's service listens for connections, e.g. 80.
	*/
	servicePort: int | null;
	/**
	 * The protocol the server uses, e.g. 'HTTP'.
	 *
	 * This should be coerced to uppercase and compared in a case-insensitive
	 * manner.
	*/
	serviceProtocol: string;
}

/**
 * This is just a collection of fields common to ATC component server instances.
 *
 * It's not a "real" type in that it isn't exposed through the API.
 * @ignore
*/
interface ATCServer extends BaseServer {
	/**
	 * The port on which the component service listens for connections, e.g. 80.
	 *
	 * Note that this is generally not enforced on the servers, and merely a way
	 * of book-keeping that information in Traffic Ops.
	*/
	port: int;
}

/**
 * Traffic Monitors are Traffic Monitor instances.
 *
 * They contain all of the server information as well as configuration for the
 * Traffic Monitor service.
*/
interface TrafficMonitor extends BaseServer {
	/**
	 * The Name of the CDN to which the Traffic Monitor belongs.
	 *
	 * It will monitor the health of Cache Servers within this same CDN.
	*/
	cdn: string;
	/**
	 * Defines the number of events to keep track of internally.
	 *
	 * '0' means "no limit".
	*/
	eventCount: int;
	/**
	 * The interval in milliseconds on which to poll Cache Server health.
	*/
	healthPollingInterval: int;
	/**
	 * The interval in milliseconds on which to poll Cache Servers 'heartbeats'.
	 *
	 * This is different than asking for their actual health, it's just a litmus
	 * test to see if they can be reached over the network.
	*/
	heartbeatPollingInterval: int;
	/**
	 * Thread Count determines ~~how soft the Monitor is~~ the number of threads
	 * used to concurrently poll Cache Server health.
	 *
	 * '0' will cause Traffic Monitor to select its own value.
	 *
	 * @note I'd suggest it use one for each available CPU core, in that case,
	 * which I think is what most libraries do by default.
	*/
	threadCount: int;
	/**
	 * A 'padding time' - in milliseconds - to add to requests to spread them
	 * out for Traffic Control systems that use a large number of Traffic
	 * Monitors.
	*/
	 timePad: int;
	/**
	 * The interval in milliseconds on which to poll Traffic Ops for
	 * configuration changes.
	*/
	 configPollingInterval: int;
}

/**
 * Traffic Routers are Traffic Router instances.
 *
 * They contain all of the server information as well as configuration for the
 * Traffic Router service.
*/
interface TrafficRouter extends BaseServer {
	/**
	 * Sets the value - in seconds - to be used for the 'max-age' parameter of
	 * the Cache-Control HTTP header in HTTP responses to clients.
	*/
	cacheControlMaxAge: int;
	/**
	 * The interval - in milliseconds - on which Traffic Router should poll for
	 * an updated <abbr title="Coverage Zone File">CZF</abbr>.
	*/
	coverageZonePollingInterval: int;
	/**
	 * The interval - in milliseconds - on which Traffic Router should poll for
	 * an updated <abbr title="Deep Coverage Zone File">DCZF</abbr>.
	*/
	deepCoverageZonePollingInterval: int;
	/**
	 * This contains settings related to "Dynamic Zone Cache Priming".
	*/
	dynamicZoneCachePriming: {
		/**
		 * If `true`, this will allow Traffic Router to attempt to prime the
		 * dynamic zone cache.
		*/
		prime: boolean;
		/**
		 * Limit the number of permutations to prime when "Dynamic Zone Cache
		 * Priming".
		 *
		 * The value `0` instructs Traffic Router to use its pre-configured
		 * default (500 by default at the time of this writing).
		 *
		 * Has no effect if [[TrafficRouter.dynamicZoneCachePriming.prime]] is
		 * `false`.
		*/
		primingLimit: int;
	}
	/**
	 * Whether or not the EDNS0 DNS extension mechanism described in
	 * [RFC2671](https://tools.ietf.org/html/rfc2671) should be made available
	 * to clients.
	*/
	edns0ClientSubnetEnabled: boolean;
	/**
	 * Whether or not to enable “Client Steering Forced Diversity”.
	 *
	 * When `true`, this will cause the Traffic Router to diversify the list of
	 * Cache Servers returned in STEERING responses by including more unique
	 * Edge-Tier Cache Servers in the response to the client’s request.
	*/
	forcedDiversity: boolean;
	/**
	 * The interval - in milliseconds - on which Traffic Router should poll for
	 * an updated geographic IP mapping database.
	*/
	geolocationPollingInterval: int;
	/**
	 * Contains options related to DNSSEC keys.
	*/
	dnssec: {
		/**
		 * If `true`, this will allow Traffic Router to use expired DNSSEC keys to
		 * sign zones; default is “true”. This helps prevent DNSSEC-related
		 * outages due to failed Traffic Control components or connectivity
		 * issues.
		*/
		allowExpiredKeys: boolean;
		/**
		 * Used when creating an effective date for a new key set.
		 *
		 * New keys are generated with an effective date of that is the
		 * effective multiplier multiplied by the
		 * <abbr title="Time to Live">TTL</abbr> less than the old key’s
		 * expiration date.
		 *
		 * The value `0` instructs Traffic Router to use its configured default,
		 * which at the time of this writing is `2`.
		*/
		effectiveMultiplier: int;
		/**
		 * The interval in seconds on which Traffic Router will check the
		 * Traffic Ops API for new DNSSEC keys.
		*/
		fetchInterval: int;
		/**
		 * The number of times Traffic Router will attempt to load DNSSEC keys
		 * before giving up.
		 *
		 * The value `0` instructs Traffic Router to use its configured default
		 * - which is `5` by default at the time of this writing.
		*/
		fetchRetries: int;
		/**
		 * The timeout in milliseconds for requests to the DNSSEC Key management
		 * endpoint of the Traffic Ops API.
		 *
		 * `0` means "no timeout".
		*/
		fetchTimeout: int;
		/**
		 * The number of milliseconds Traffic Router will wait in between
		 * attempts to load DNSSEC keys.
		 *
		 * The value `0` instructs Traffic Router to use its configured default.
		*/
		fetchWait: int;
		/**
		 * Used to determine when new DNSSEC keys need to be generated.
		 *
		 * Keys are re-generated if expiration is less than the generation
		 * multiplier multiplied by the <abbr title="Time to Live">TTL</abbr>.
		 *
		 * The value `0` instructs Traffic Router to use its configured default,
		 * which at the time of this writing is `10`.
		*/
		generationMultiplier: int;
		/**
		 * If DNSSEC is enabled on the CDN to which this Traffic Router belongs,
		 * enabling this parameter allows Traffic Router to compare existing
		 * zones with newly generated zones.
		 *
		 * If the newly generated zone is the same as the existing zone, Traffic
		 * Router will simply reuse the existing signed zone instead of signing
		 * the same new zone.
		*/
		zoneComparisons: boolean;
	}
	/**
	 * Contains configuration information for the
	 * <abbr title="Start of Authority">SOA</abbr> records served by the Traffic
	 * Router.
	*/
	soa: {
		/**
		 * The email address of the administrator of the DNS zones for which the
		 * Traffic Router is authoritative.
		 * @note I'd suggest the email of the creating user be used as a default
		 * value in UI forms.
		*/
		admin: string;
		/**
		 * The value for the "expire" field the Traffic Router DNS Server will
		 * respond with on <abbr title="Start of Authority">SOA</abbr> records.
		 *
		 * The value is in seconds.
		*/
		expire: int;
		/**
		 * The value for the "minimum" field the Traffic Router DNS Server will
		 * respond with on <abbr title="Start of Authority">SOA</abbr> records.
		 *
		 * The value is in seconds.
		*/
		minimum: int;
		/**
		 * The value for the "refresh" field the Traffic Router DNS Server will
		 * respond with on <abbr title="Start of Authority">SOA</abbr> records.
		 *
		 * The value is in seconds.
		*/
		refresh: int;
		/**
		 * The value for the "retry" field the Traffic Router DNS Server will
		 * respond with on <abbr title="Start of Authority">SOA</abbr> records.
		 *
		 * The value is in seconds.
		*/
		retry: int;
	}
	/**
	 * The interval, in seconds, on which Traffic Control components should
	 * check for new steering mappings.
	*/
	steeringPollingInterval: int;
	/**
	 * This determines the <abbr title="Time to Live">TTL</abbr>s of various DNS
	 * record types returned by Traffic Router.
	 *
	 * The values are in seconds.
	*/
	ttls: {
		/** <abbr title="Time to Live">TTL</abbr> for A records */
		a: int;
		/** <abbr title="Time to Live">TTL</abbr> for AAAA records */
		aaaa: int;
		/** <abbr title="Time to Live">TTL</abbr> for DNSKEY records */
		dnskey: int;
		/** <abbr title="Time to Live">TTL</abbr> for DS records */
		ds: int;
		/** <abbr title="Time to Live">TTL</abbr> for NS records */
		ns: int;
		/**
		 * <abbr title="Time to Live">TTL</abbr> for
		 * <abbr title="Start of Authority">SOA</abbr> records
		*/
		soa: int;
	}
	/**
	 * Contains configuration options for DNS Zone management.
	*/
	zones: {
		/**
		 * The interval in seconds on which Traffic Router will check for zones
		 * that need to be re-signed or if dynamic zones need to be expired from
		 * its cache.
		*/
		cacheMaintenanceInterval: int;
		/**
		 * An integer that defines the initial size of the Guava cache, default
		 * is 10000.
		 * @todo Should this be configurable in TO, or only server-side? Also,
		 * what are the units?
		*/
		dynamicInitialCapacity: int;
		/**
		 * A duration in seconds that defines how long a dynamic zone will
		 * remain valid before expiring.
		*/
		dynamicResponseExpiration: int;
		/**
		 * An integer that defines the number of minutes to allow for zone
		 * generation; this bounds the zone priming activity.
		*/
		initTimeout: int;
		/**
		 * An integer that defines the size of the concurrency level (threads)
		 * of the Guava cache used by ZoneManager to store zone material.
		*/
		threadCount: int;
		/**
		 * Multiplier used to determine the number of CPU cores to use for zone
		 * signing operations.
		 *
		 * The value `0` instructs Traffic Router to use its configured default,
		 * which at the time of this writing is `0.75`.
		*/
		threadMultiplier: number;
	}
}

/**
 * An Origin is a source of content to be delivered through the CDN.
*/
interface Origin {
	/**
	 * The full URL of the origin, including schema and port e.g.
	 * `https://origin.test:443`.
	 *
	 * If the port is omitted, it is assumed to be the protocol's standard port,
	 * e.g. 80 for HTTP.
	*/
	url: URL;
	/** Still numeric IDs because we want a single, unique identifier */
	// id: int; //TODO: necessary?
	/**
	 * The Origin's IPv4 address.
	 *
	 * If it's 'null', the assumption is that the Origin doesn't do IPv4.
	*/
	ipv4Address: string | null;
	/**
	 * The Origin's IPv6 address.
	 *
	 * If it's 'null', the assumption is that the Origin doesn't do IPv6.
	*/
	ipv6Address: string | null;
	/** Arbitrary text for miscellaneous purposes */
	notes: string;
	/** The Names of Tags given to this Origin */
	tags: Set<string>;
	/**
	 * The Name of the Tenant to whom this Origin belongs
	*/
	tenant: string;
}

/**
 * A Traffic Portal is a Traffic Portal instance.
 *
 * They contain all of the server information as well as configuration for the
 * Traffic Portal instance.
*/
interface TrafficPortal {
	/**
	 * The full URL of the Traffic Portal, including schema and port e.g.
	 * `https://origin.test:443`.
	 *
	 * If the port is omitted, it is assumed to be the protocol's standard port,
	 * e.g. 80 for HTTP.
	*/
	url: URL;
	/**
	 * The Traffic Portal's IPv4 address.
	 *
	 * If it's 'null', the assumption is that the Traffic Portal doesn't do
	 * IPv4.
	*/
	ipv4Address: string | null;
	/**
	 * The Traffic Portal's IPv6 address.
	 *
	 * If it's 'null', the assumption is that the Traffic Portal doesn't do
	 * IPv6.
	*/
	ipv6Address: string | null;
	/** Arbitrary text for miscellaneous purposes */
	notes: string;
	/** The Names of Tags given to this Traffic Portal */
	tags: Set<string>;
}

/**
 * A Role is essentially just a grouping of permissions that can be given to a
 * user.
 * @note No privilege level, only permissions.
*/
interface Role {
	/**
	 * The Role's unique identifier.
	*/
	name: string;
	/**
	 * A short description of the Role.
	*/
	description: string;
	/**
	 * Defines what users with this Role can do.
	 * (Currently what we call 'capabilities').
	*/
	permissions: Set<Permission>;
}

/**
 * A Tenant is essentially a scope of a user's ability to affect Delivery
 * Services.
*/
interface Tenant {
	/**
	 * Defines whether or not a Tenant is able to actively modify/access its
	 * resources.
	*/
	active: boolean;
	/**
	 * The uniquely identifiying name of the Tenant.
	*/
	name: string;
	/**
	 * The Name of the Tenant's 'parent'.
	 *
	 * A Tenant has access to all of the resources scoped within itself, as well
	 * as all of the resources to which any Tenant having it as the 'parent' has
	 * access.
	*/
	parent: string;
}

/**
 * A user is... well it's a user.
*/
interface User {
	/**
	 * This is a text field with no defined structure, but semantically it
	 * represents a user’s real-world, physical address.
	 *
	 * If it is not an empty value, it is assumed to contain enough information
	 * to send a letter through normal postage to the user.
	 *
	 * Addresses may consist of alphanumeric characters, hyphens, periods,
	 * spaces, and newlines, but may neither begin nor end with a space or
	 * newline.
	*/
	address: string;
	/**
	 * A user’s email address, which is used for initial registration and
	 * password recovery.
	 *
	 * A User object is not guaranteed to have an Email because the initial,
	 * default User will not have one. Emails are unique among all Users, unless
	 * they are ”Null”-valued.
	*/
	email: string | null;
	/**
	 * A user’s ”full” name, as it would appear on a letter mailed to them.
	 *
	 * This field, if not empty, is presumed to be the user’s name as it would
	 * appear on normal postage. This field may not contain non-alphabetic
	 * characters.
	*/
	fullName: string;
	/**
	 * A string which, if not empty, is presumed to be a telephone number at
	 * which the user may be contacted.
	 *
	 * It may only contain numerics and hyphens, and may neither start nor end
	 * with a hyphen.
	*/
	phoneNumber: string;
	/**
	 * The Name of the user's Role
	*/
	role: string;
	/**
	 * The Name of the Tenant to which the user belongs.
	*/
	tenant: string;
	/**
	 * The user's unique username.
	*/
	username: string;
}

/**
 * All of the pre-defined Permissions.
*/
enum Permission {
	/** Ability to authenticate */
	'auth',
	/** Ability to view Cache Groups */
	'cache-groups-read',
	/** Ability to edit Cache Groups */
	'cache-groups-write',
	/** Ability to view Cache Servers */
	'cache-servers-read',
	/** Ability to edit Cache Servers */
	'cache-servers-write',
	/** Ability to view Capabilities */
	'capabilities-read',
	/** Ability to edit Capabilities */
	'capabilities-write',
	/** Ability to view CDNS */
	'cdns-read',
	/** Ability to edit CDNS */
	'cdns-write',
	/** Ability to snapshot a CDN */
	'cdns-snapshot',
	/** Ability to queue updates on a CDN */
	'cdns-queue-updates',
	/** Ability to view CDN security keys */
	'cdn-security-keys-read',
	/** Ability to edit CDN security keys */
	'cdn-security-keys-write',
	/** Ability to view change logs */
	'change-logs-read',
	/** Ability to use Pattern-Based Consistent Hash Test Tool */
	'consistenthash-read',
	/** Ability to view Delivery Services */
	'delivery-services-read',
	/** Ability to view Delivery Services */
	'delivery-services-write',
	/** Ability to edit the Capabilities required by a Delivery Service */
	'delivery-service-capabilities-write',
	/** Ability to view Delivery Service security keys */
	'delivery-service-security-keys-read',
	/** Ability to edit Delivery Service security keys */
	'delivery-service-security-keys-write',
	/** Ability to view Delivery Service Requests */
	'delivery-service-requests-read',
	/** Ability to edit Delivery Service Requests */
	'delivery-service-requests-write',
	/** Ability to edit Delivery Service / server assignments */
	'delivery-service-servers-write',
	/** Ability to change the Status of a Delivery Service */
	'delivery-service-status-write',
	/** Ability to edit STEERING Delivery Service Targets */
	'delivery-service-targets-write',
	/** Ability to view infrastructure servers */
	'infrastructure-servers-read',
	/** Ability to edit infrastructure servers */
	'infrastructure-servers-write',
	/** Ability to view Content Invalidation Requests */
	'jobs-read',
	/** Ability to edit Content Invalidation Requests */
	'jobs-write',
	/** Ability to view Origins */
	'origins-read',
	/** Ability to edit Origins */
	'origins-write',
	/** Ability to view Parameters */
	'parameters-read',
	/** Ability to edit Parameters */
	'parameters-write',
	/** Ability to view Physical Locations */
	'physical-locations-read',
	/** Ability to edit Physical Locations */
	'physical-locations-write',
	/** Ability to view Profiles */
	'profiles-read',
	/** Ability to edit Profiles */
	'profiles-write',
	/** Ability to view Roles */
	'roles-read',
	/** Ability to edit Roles. */
	'roles-write',
	/** Ability to edit the Capabilities assigned to a Cache Server */
	'server-capabilities-write',
	/** Ability to view Tags */
	'tags-read',
	/** Ability to edit Tags */
	'tags-write',
	/** Ability to view system info */
	'system-info-read',
	/** Ability to view tenants */
	'tenants-read',
	/** Ability to edit tenants */
	'tenants-write',
	/** Ability to view Traffic Monitors */
	'traffic-monitors-read',
	/** Ability to edit Traffic Monitors */
	'traffic-monitors-write',
	/** Ability to view Traffic Portals */
	'traffic-portals-read',
	/** Ability to edit Traffic Portals */
	'traffic-portals-write',
	/** Ability to view Traffic Routers */
	'traffic-routers-read',
	/** Ability to edit Traffic Routers */
	'traffic-routers-write',
	/** Ability to view Traffic Stats instances */
	'traffic-stats-read',
	/** Ability to edit Traffic Stats instances */
	'traffic-stats-write',
	/** Ability to view Traffic Vaults */
	'traffic-vaults-read',
	/** Ability to edit Traffic Vaults */
	'traffic-vaults-write',
	/** Ability to register new users */
	'users-register',
	/** Ability to view users */
	'users-read',
	/** Ability to edit users */
	'users-write'
}
