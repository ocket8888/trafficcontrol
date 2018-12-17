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

.. _to-api-cdns-dnsseckeys-generate:

****************************
``cdns/dnsseckeys/generate``
****************************

``POST``
========
Generates new :abbr:`DNSSEC (DNS Security Extensions)` keys for a CDN.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object (string)

Request Structure
-----------------
:effectiveDate:     An optional integer which, if present in the request payload, should be the UNIX timestamp when the :abbr:`DNSSEC (DNS Security Extensions)` keys will come into effect (default: immediately i.e. the current date and time)
:key:               The name of the CDN for which :abbr:`DNSSEC (DNS Security Extensions)` keys will be generated
:kskExpirationDays: The number of days after which the long-term :abbr:`KSK (Key-Signing Key)` will expire and need to be refreshed
:name:              The :abbr:`TLD (Top-Level Domain)` for which the :abbr:`DNSSEC (DNS Security Extensions)` keys shall be generated
:ttl:               The :abbr:`TTL (Time To Live)` in seconds for secure responses to DNS queries, after which downstream routers will need to revalidate them
:zskExpirationDays: The number of days after which the short-term :abbr:`ZSK (Zone-Signing Key)` will expire and need to be refreshed

.. code-block:: http
	:caption: Request Example

	POST /api/1.1/cdns/dnsseckeys/generate HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 110
	Content-Type: application/json

	{
		"key": "CDN-in-a-Box",
		"name": "mycdn.ciab.test",
		"kskExpirationDays": 1,
		"zskExpirationDays": 1,
		"ttl": 60
	}

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: LVQMc/XuQk1qHa7Uwa0ymNPs6KCZqag8QSguiAZr5jUEJOOa4PaxMu3n+Cnce8/o6lCDmTGB78BN3tY08SjytQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 17 Dec 2018 20:49:34 GMT
	Content-Length: 64

	{
		"response": "Successfully created dnssec keys for CDN-in-a-Box"
	}
