..
..
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
..
..     http://www.apache.org/licenses/LICENSE-2.0
..
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
..

.. _to-api-servers-details:

*******************
``servers/details``
*******************
.. versionadded:: 1.2

``GET``
=======
Gets details about servers.

.. seealso:: The only reason to use this endpoint over :ref:`to-api-servers` is because this one supports pagination. If pagination is not required, :ref:`to-api-servers` should be used instead.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+----------------+-----------+----------------------------------------------------------------------------------------------------------+
	| Name           | Required  | Description                                                                                              |
	+================+===========+==========================================================================================================+
	| hostName       | yes\ [1]_ | Return only the server that has this (short) hostname                                                    |
	+----------------+-----------+----------------------------------------------------------------------------------------------------------+
	| limit          | no        | Limit the number of returned servers to this maximum (default: 1000)                                     |
	+----------------+-----------+----------------------------------------------------------------------------------------------------------+
	| orderBy        | no        | Order the ``response`` array by this key (default: ``hostName``)                                         |
	+----------------+-----------+----------------------------------------------------------------------------------------------------------+
	| physLocationID | yes\ [1]_ | Return only servers that reside in the physical location  identified by this integral, unique identifier |
	+----------------+-----------+----------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.2/servers/details?hostName=edge HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

.. [1] Exactly one of "hostName" or "physLocationID" is required. Both is also acceptable.

Response Structure
------------------
:cachegroup:       The name of the Cache Group to which this server belongs
:cdnName:          Name of the CDN to which the server belongs
:deliveryservices: An array of integral, unique identifiers to which this server has been assigned
:domainName:       The domain part of the server's Fully Qualified Domain Name (FQDN)
:guid:             An identifier used to uniquely identify the server

	.. deprecated:: 1.1
		This is a legacy key which only still exists for compatibility reasons - it should always be ``null``

:hardwareInfo: An object containing information about a server's hardware

	.. deprecated:: 1.1
		This is a legacy key which only still exists for compatibility reasons - it should always be ``null``

:hostName:       The (short) hostname of the server
:httpsPort:      The port on which the server listens for incoming HTTPS connections/requests
:id:             An integral, unique identifier for this server
:iloIpAddress:   The IPv4 address of the server's Integrated Lights-Out (ILO) service\ [1]_
:iloIpGateway:   The IPv4 gateway address of the server's ILO service\ [1]_
:iloIpNetmask:   The IPv4 subnet mask of the server's ILO service\ [1]_
:iloPassword:    The password of the of the server's ILO service user\ [1]_ - displays as simply ``******`` if the currently logged-in user does not have the 'admin' or 'operations' role(s)
:iloUsername:    The user name for the server's ILO service\ [1]_
:interfaceMtu:   The Maximum Transmission Unit (MTU) to configured on ``interfaceName``
:interfaceName:  The name of the primary network interface used by the server
:ip6Address:     The IPv6 address and subnet mask of ``interfaceName``
:ip6Gateway:     The IPv6 address of the gateway used by ``interfaceName``
:ipAddress:      The IPv4 address of ``interfaceName``
:ipGateway:      The IPv4 address of the gateway used by ``interfaceName``
:ipNetmask:      The IPv4 subnet mask used by ``interfaceName``
:mgmtIpAddress:  The IPv4 address of some network interface on the server used for 'management'
:mgmtIpGateway:  The IPv4 address of a gateway used by some network interface on the server used for 'management'
:mgmtIpNetmask:  The IPv4 subnet mask used by some network interface on the server used for 'management'
:offlineReason:  A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status
:physLocation:   The name of the physical location where the server resides
:profile:        The name of the profile this server uses
:profileDesc:    A description of the profile this server uses
:rack:           A string indicating "server rack" location
:routerHostName: The human-readable name of the router responsible for reaching this server
:routerPortName: The human-readable name of the port used by the router responsible for reaching this server
:status:         The status of the server

	.. seealso:: :ref:`health-proto`

:tcpPort: The port on which this server listens for incoming TCP connections

	.. note:: This is typically thought of as synonymous with "HTTP port", as the port specified by ``httpsPort`` may also be used for incoming TCP connections.

:type:       The name of the 'type' of this server
:xmppId:     An identifier to be used in XMPP communications with the server - in nearly all cases this will be the same as ``hostName``
:xmppPasswd: The password used in XMPP communications with the server

The response JSON payload also contains three non-standard, top-level keys (i.e. not "alerts" or "response")

:limit:   The maximum number of servers to which the ``response`` array was limited
:orderby: The key by which the elements of ``response`` were ordered
:size:    The actual number of elements in the ``response`` array - not to exceed ``limit``

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: SEGdGB4Ogpgp5GUJ0PHNaCf7vjBT8Zne8wnA3psJHg7gtAfI2bV/V0Gg/DG0lae0IWgghUoXsOVA6xHNIhHuEA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Dec 2018 15:57:08 GMT
	Content-Length: 865

	{
		"limit": 1000,
		"orderby": "hostName",
		"response": [{
			"cachegroup": "CDN_in_a_Box_Edge",
			"cdnName": "CDN-in-a-Box",
			"deliveryservices": [
				1
			],
			"domainName": "infra.ciab.test",
			"guid": null,
			"hardwareInfo": null,
			"hostName": "edge",
			"httpsPort": 443,
			"id": 9,
			"iloIpAddress": "",
			"iloIpGateway": "",
			"iloIpNetmask": "",
			"iloPassword": "",
			"iloUsername": "",
			"interfaceMtu": 1500,
			"interfaceName": "eth0",
			"ip6Address": "fc01:9400:1000:8::100",
			"ip6Gateway": "fc01:9400:1000:8::1",
			"ipAddress": "172.16.239.100",
			"ipGateway": "172.16.239.1",
			"ipNetmask": "255.255.255.0",
			"mgmtIpAddress": "",
			"mgmtIpGateway": "",
			"mgmtIpNetmask": "",
			"offlineReason": "",
			"physLocation": "Apachecon North America 2018",
			"profile": "ATS_EDGE_TIER_CACHE",
			"profileDesc": "Edge Cache - Apache Traffic Server",
			"rack": "",
			"routerHostName": "",
			"routerPortName": "",
			"status": "REPORTED",
			"tcpPort": 80,
			"type": "EDGE",
			"xmppId": "edge",
			"xmppPasswd": ""
		}],
		"size": 1
	}
