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

.. _to-api-servers-checks:

******************
``servers/checks``
******************

``GET``
=======
Retrieves the values of the Traffic Ops :ref:`to-check-ext` for each cache server.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. code-block:: http
	:caption: Request Example

	GET /api/1.4/servers/checks HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:adminState: The server's health state

	.. seealso:: :ref:`health-proto`

:cacheGroup:   The name of the Cache Group which
:hostName:     The server's (short) hostname
:id:           an integral, unique identifier for the server
:profile:      The name of the profile used by the server
:revalPending: ``true`` if the server has a revalidation pending, or ``false`` if it does not
:type:         The name of the type of this server

	EDGE
		This is an Edge-tier cache server
	MID
		This is a Mid-tier cache server

:updPending: ``true`` if the server has updates pending, or ``false`` if it does not

.. note:: Only the :ref:`to-check-ext` values for Edge-tier and Mid-tier cache servers.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Tue, 18 Dec 2018 17:43:49 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Tue, 18 Dec 2018 21:43:49 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: bdPnUNnyELT6vn/E2Du2mVKNI2U1Cp+lu/RsrCyqnGLc4ahS0x370PLUvrhHPqcUe+E7lWBUl23zqp3F5ATMJA==
	Content-Length: 350

	{ "response": [
		{
			"profile": "ATS_EDGE_TIER_CACHE",
			"cacheGroup": "CDN_in_a_Box_Edge",
			"updPending": false,
			"hostName": "edge",
			"revalPending": false,
			"adminState": "REPORTED",
			"id": 9,
			"type": "EDGE"
		},
		{
			"profile": "ATS_MID_TIER_CACHE",
			"cacheGroup": "CDN_in_a_Box_Mid",
			"updPending": false,
			"hostName": "mid",
			"revalPending": false,
			"adminState": "REPORTED",
			"id": 8,
			"type": "MID"
		}
	]}
