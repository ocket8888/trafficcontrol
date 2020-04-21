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
 * Delivery Services are complex, so I had to split them out into their own
 * module to avoid things getting too crazy to navigate.
 * @packageDocumentation
*/

import {int} from './datamodel';

/**
 * This defines what other components of ATC will consider a Delivery Service
 * "active".
 *
 * It's not an object exposed through the API in its own right, just a
 * specification of the allowed values.
*/
enum DeliveryServiceActiveState {
	/**
	 * A Delivery Service that is ”active” is one that is functionally
	 * in service, and fully capable of delivering content.
	 *
	 * This means that its configuration is deployed to Cache Servers and it is
	 * available for routing traffic.
	*/
	ACTIVE = 'ACTIVE',
	/**
	 * A Delivery Service that is ”inactive” is not available for
	 * routing and has not had its configuration distributed to its assigned
	 * Cache Servers.
	*/
	INACTIVE = 'INACTIVE',
	/**
	 * A Delivery Service that is ”primed” has had its configuration
	 * distributed to the various servers required to serve its content.
	 * However, the content itself is still inaccessible via routing.
	*/
	PRIMED = 'PRIMED'
}

/**
 * This defines how content may be cached for a Delivery Service.
 *
 * It's not an object exposed through the API in its own right, just a
 * specification of the allowed values.
*/
enum CachingType {
	/** Cache normally */
	CACHE = 'CACHE',
	/** Don't cache, only proxy */
	NO_CACHE = 'NO_CACHE',
	/** Cache in RAM block devices only */
	RAM_ONLY = 'RAM_ONLY'
}

/**
 * This is a collection of all of the fields that **all** Delivery Services
 * have, regardless of routing type.
 *
 * It's not a "real" type, the only "real" types are the ones that make up
 * [[DeliveryService]].
*/
interface BaseDeliveryService {
	/**
	 * A Delivery Service that has ”Anonymous Blocking” tells Traffic Router to
	 * block requests from anonymized IP addresses.
	 *
	 * No guarantee is made that this will work very well for DNS-routed
	 * Delivery Services.
	*/
	anonymousBlocking: boolean;
	/**
	 * This is an amalgamation of the various bypasses we have now.
	 *
	 * It's always a string, but it's interpreted differently based on the
	 * Routing Type of the Delivery Service.
	 *
	 * - HTTP: This must be an FQDN and is interpreted/validated as such.
	 * - DNS or STATIC: If this can be parsed as an IPv4 address then it is
	 * 	interpreted as such and presented as an A record, otherwise if it
	 * 	can be parsed as an IPv6 address then it is interpreted as such and
	 * 	presented as an AAAA record, and finally if all else fails and it
	 * 	can be parsed as an FQDN then it is is treated as such and presented
	 * 	as a CNAME record.
	 * - STEERING: This MUST be the Name of an existing Delivery Service (but
	 * 	need not be one of the Delivery Service's Targets).
	*/
	bypassDestination: string;
	/** The Name of the CDN to which the Delivery Service belongs */
	cdn: string;
	/**
	 * The network location to which clients will be directed if they are denied
	 * access on the basis of Anonymous Blocking and/or Geographic Limiting
	 * settings.
	 *
	 * The rules regarding how it's interpreted are the same as those for
	 * Bypass Destination.
	*/
	deniedAccessRedirect: string;
	/**
	 * Whether or not the EDNS0 DNS extension mechanism described in
	 * [RFC2671](https://tools.ietf.org/html/rfc2671) should be made available
	 * to clients.
	 *
	 * Note that the ability of a Traffic Router to actually implement this
	 * setting depends on its own [[TrafficRouter.edns0ClientSubnetEnabled]]
	 * value.
	*/
	edns0ClientSubnetEnabled: boolean;
	/**
	 * This property describes limitations to the availability of this Delivery
	 * Service’s content on the basis of the requesting client’s geographic
	 * location.
	 *
	 * It is a set of strings, each of which is an ISO 3166-1 alpha-2 country
	 * code, optionally with ISO 3166-2 subdivisional alphabetic code. This is a
	 * ”white list” of countries/subdivisions wherein content is to be made
	 * available.
	 *
	 * Content is always available to clients whose IP addresses are found
	 * within the Traffic Routers’ Coverage Zone File(s). With that in mind,
	 * when this property is an empty set it means that no geographic regions
	 * are ”whitelisted” and thus only clients whose IP addresses are found
	 * within a Coverage Zone File will be granted access to content. When this
	 * property has a ”Null” type, there is no geographic restriction placed on
	 * the Delivery Service’s content access.
	*/
	geographicLimiting: Set<string> | null;
	/**
	 * A Delivery Service's Name uniquely identifies it.
	 *
	 * It's only allowed to contain alphanumerics, spaces, underscores and
	 * hyphens and can only begin and end with an alphanumeric. The thing we
	 * currently call an XML_ID can be generated from this by replacing
	 * everything that isn't alphanumeric or a hyphen with a hyphen.
	*/
	name: string;
	/** Arbitrary, structure-less text field */
	notes: string;
	/**
	 * The Delivery Service's origin.
	 *
	 * This has a different meaning for STATIC Delivery Services; it needs to be
	 * an external CDN's resolver.
	*/
	origin: string | int;
	/**
	 * Sets a little vanity name, same as the current Routing Name
	*/
	routingName: string;
	/**
	 * The Routing Type of a Delivery Service defines how content is served.
	 *
	 * The values are more completely explained on each sub-type.
	 * Note, though, that this is different from *protocol*.
	*/
	routingType: 'HTTP' | 'DNS' | 'STATIC' | 'STEERING';
	/**
	 * Determines whether or not the Delivery Service is currenly active or
	 * routed.
	*/
	status: DeliveryServiceActiveState;
	/**
	 * A set of the names of Tags that have been applied to the Delivery
	 * Service.
	*/
	tags: Set<string>;
	/**
	 * The Name of the Tenant that owns this Delivery Service.
	*/
	tenant: string;
	/**
	 * URLs that should be used as aliases of this Delivery Service's accessible
	 * routing FQDN(s) just because they look good.
	 *
	 * These must still be rooted in the Domain of the CDN to which this
	 * Delivery Service is assigned. They cannot be the same as the concatenation
	 * of any Delivery Service's (including this one) Routing Name, Name (with
	 * the appropriate replacements) and CDN Domain. No two Delivery Services
	 * may share a single Vanity URL.
	*/
	vanityURLs: Set<URL>;
}

/**
 * This is a collection of all of the fields that all *non-STATIC* Delivery
 * Services have.
 *
 * It's not a "real" type, the only "real" types are the ones that make up
 * [[DeliveryService]].
 * @todo I'm sure there's more than just type constraints to be done here.
*/
interface NonStaticDeliveryService extends BaseDeliveryService {
	/**
	 * The ID of an Origin for which this Delivery Service is responsible for
	 * serving content.
	*/
	origin: int;
	/**
	 * These are the routingTypes that are considered "Non-Static".
	 * (Which is to say, all of them except 'STATIC')
	*/
	routingType: 'HTTP' | 'DNS' | 'STEERING';
}

/**
 * This defines how Query Strings are handled by cache servers serving content
 * for a Delivery Service.
 *
 * It's not an object exposed through the API in its own right, just a
 * specification of the allowed values.
 * @todo See the "TODO" for [[HTTPDeliveryService.consistentHashingExpression]]
*/
enum QueryStringHandling {
	/**
	 * Caches strip query strings before processing requests.
	*/
	DROP = 'DROP',
	/**
	 * The query string is not stripped, but it is not considered for caching.
	*/
	IGNORE = 'IGNORE',
	/**
	 * Query strings are considered for caching (and are therefore not stripped)
	*/
	USE = 'USE'
}

/**
 * This defines how Range Requests are handled by cache servers when serving
 * content for a Delivery Service.
 *
 * It's not an object exposed through the API in its own right, just a
 * specification of the allowed values.
*/
enum RangeRequestHandling {
	/**
	 * Cache the range request object as an object in its own right.
	*/
	CACHE = 'CACHE',
	/**
	 * Don't cache range requests, just proxy.
	*/
	NO_CACHE = 'NO_CACHE',
	/**
	 * Transparently serve range requests while caching the whole object.
	*/
	WHOLE_OBJECT = 'WHOLE_OBJECT'
}

/**
 * This defines what protocols are supported by the Delivery Service at the
 * content delivery level (i.e. after the DNS step has completed).
 *
 * This value must be obeyed by both the routing and caching infrastructures.
*/
enum Protocol {
	/**
	 * This Delivery Service handles **only** unsecured HTTP traffic.
	*/
	'HTTP',
	/**
	 * This Delivery Service handles **only** secured HTTPS traffic.
	*/
	'HTTPS',
	/**
	 * This Delivery Service handles both HTTP **and** HTTPS traffic.
	*/
	'HTTP_AND_HTTPS',
	/**
	 * This Delivery Service handles HTTPS normally, and handles HTTP traffic by
	 * redirecting it to use HTTPS.
	*/
	'HTTP_TO_HTTPS'
}

/**
 * This is a collection of all of the fields that all *non-STATIC* and
 * *non-STEERING* Delivery Services have.
 *
 * It's not a "real" type, the only "real" types are the ones that make up
 * [[DeliveryService]].
*/
interface NonSteeringDeliveryService extends NonStaticDeliveryService {
	/** Defines how the Delivery Service's content may be cached */
	caching: CachingType;
	/**
	 * Sets the maximum allowed connections to the Delivery Service's Origin.
	 *
	 * Typically this means "from the Mid-tier" but also possibly Edge-tier
	 * if the used Topology is such that Mid-tier Cache Servers are not used.
	 * The value "null" has the special meaning "no limit".
	*/
	maxOriginConnections: int | null;
	/**
	 * These are the routingTypes that are considered "Non-Steering" and
	 * "Non-Static".
	 * (Which is to say, all of them except 'STATIC' and 'STEERING')
	*/
	routingType: 'HTTP' | 'DNS';
	/**
	 * This is the name of the Topology used by the Delivery Service.
	 */
	topology: string;
}

/**
 * A Delivery Service that uses HTTP-based routing.
*/
interface HTTPDeliveryService extends NonStaticDeliveryService {
	/**
	 * Key-value pairs where the key is an HTTP header to be returned by the
	 * Traffic Router in all (successful) redirect responses to the client.
	 *
	 * The value is (obviously) the value of the header.
	*/
	additionalResponseHeaders: Map<string, string>;
	/**
	 * A regular expression that is applied to requested URLs to determine what
	 * is hashed.
	 * @todo Jon thinks that we should make people "think about their cache
	 * key," which sounds great except that I need to know more about the
	 * implementation.
	*/
	consistentHashingExpression: RegExp;
	/** Whether or not "Deep Caching" can be used for this Delivery Service */
	deepCaching: boolean;
	/** The DSCP to use for the Delivery Service's traffic */
	dscp: int;
	/**
	 * A set of HTTP headers for Traffic Router to take note of in its logs.
	*/
	loggedRequestHEaders: Set<string>;
	/**
	 * Determines the protocols supported by the Delivery Service's content
	 * delivery system.
	*/
	protocol: Protocol;
	/**
	 * Defines how query strings are handled by this Delivery Service
	*/
	queryStringHandling: QueryStringHandling;
	/** Defines how Range Requests are cached */
	rangeRequestHandling: RangeRequestHandling;
	/** The Names of the Capabilities this Delivery Service requires to operate */
	requiredCapabilities: Set<string>;
	/**
	 * Traffic Router responds to DNS requests with its own address, and
	 * redirects clients to cache servers using HTTP 3XX responses.
	*/
	routingType: 'HTTP';
}

/**
 * A Delivery Service that uses DNS-based routing.
*/
interface DNSDeliveryService extends NonStaticDeliveryService {
	/**
	 * Sets the TTL for DNS responses containing the bypass destination (in
	 * seconds).
	*/
	bypassTTL: int;
	/** All of the IDs of the Cache Servers assigned to this Delivery Service */
	cacheServers: Set<int>;
	/** Defines how the Delivery Service's content may be cached */
	caching: CachingType;
	/** Whether or not "Deep Caching" can be used for this Delivery Service */
	deepCaching: boolean;
	/** The DSCP to use for the Delivery Service's traffic */
	dscp: int;
	/**
	 * Sets the maximum number of returned DNS records.
	 *
	 * '0' means "no limit".
	*/
	maxRecords: int;
	/**
	 * Determines the protocols supported by the Delivery Service's content
	 * delivery system.
	*/
	protocol: Protocol;
	/**
	 * Defines how query strings are handled by this Delivery Service
	*/
	queryStringHandling: QueryStringHandling;
	/** Defines how Range Requests are cached */
	rangeRequestHandling: RangeRequestHandling;
	/**
	 * The Names of the Capabilities this Delivery Service requires to operate.
	*/
	requiredCapabilities: Set<string>;
	/**
	 * Traffic Router responds to DNS requests for this Delivery Service with
	 * the address of a cache server.
	*/
	routingType: 'DNS';
}

/**
 * This is the type definition for the objects that appear as entries in a
 * STEERING Delivery Service's [[SteeringDeliveryService.targets]] set.
 *
 * It's not a "real" type, the only "real" types are those that make up
 * [[DeliveryService]].
*/
interface Target {
	/** The type of target, same as always */
	type: 'STEERING_WEIGHT' | 'STEERING_ORDER' | 'STEERING_GEO_WEIGHT' | 'STEERING_GEO_ORDER',
	/** The value of the target, with a meaning that depends on its type */
	value: int;
	/**
	 * The Name of the Delivery Service that is this Target.
	 *
	 * This **must** refer to a DNS or HTTP Delivery Service.
	*/
	target: string;

}

/**
 * A STEERING Delivery Service serves content by redirecting clients to Cache
 * Servers assigned to its constituent Delivery Services (targets).
*/
interface SteeringDeliveryService extends NonStaticDeliveryService {
	/**
	 * The ID of an Origin for which this Delivery Service is responsible for
	 * serving content.
	*/
	origin: int;
	/**
	 * Traffic Router responds to DNS requests with its own address, and
	 * redirects clients to cache servers assigned to this Delivery Service's
	 * [[SteeringDeliveryService.targets]] using HTTP 3XX responses.
	*/
	routingType: 'STEERING';
	/** A collection of targets for Steering */
	targets: Set<Target>;
}

/**
 * A STATIC Delivery Service serves content by redirecting clients to external
 * resolvers.
 *
 * It's basically what we currently refer to as a "Federation".
 * @todo This possibly requires multiple origins.
*/
interface StaticDeliveryService extends BaseDeliveryService {
	/**
	 * The resolver CNAME to which clients are redirected.
	*/
	origin: string;
	/**
	 * Traffic Router responds to DNS requests for this Delivery Service with
	 * a resolver's address.
	*/
	routingType: 'STATIC';
}

/**
 * This is all of the different kinds of Delivery Service.
 *
 * That is, anything that is a "Delivery Service" is one of these types.
*/
type DeliveryService = DNSDeliveryService | HTTPDeliveryService | StaticDeliveryService | SteeringDeliveryService;
