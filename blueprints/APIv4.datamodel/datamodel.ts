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

/**
 * Traffic Routers are Traffic Router instances.
 *
 * They contain all of the server information as well as configuration for the
 * Traffic Router service.
*/
interface TrafficRouter {
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
		 * An integer that defines the initial size of the Guava cache, used by
		 * the zone manager. default is 10000.
		 * @todo Should this be configurable in TO, or only server-side?
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
