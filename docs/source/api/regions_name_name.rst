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

.. _to-api-regions-name-name:

*************************
``regions/name/{{name}}``
*************************

``GET``
=======
.. deprecated:: 1.1
	Use the ``name`` query parameter of a ``GET`` request to the :ref:`to-api-regions` instead.

Retrieves information about a specific region.

:Auth. Required: Yes
:Roles Request:  None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------------------------------------+
	| Name | Description                            |
	+======+========================================+
	| name | The name of the region to be retrieved |
	+------+----------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/regions/name/Montreal HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:division: An object representing the division that contains this region

	:id:   The integral, unique identifier of the division which contains this region
	:name: The name of the division which contains this region

:id:           An integral, unique identifier for this region
:lastUpdated:  The date and time at which this region was last updated, in ISO format
:name:         The region name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: RUK7RKSMXsn4CkMPiSEtU0aus5yQurJAGVs/uz5vV5TAcRsrPPiOAHBVD2v5B5VSUCF/dZFQMg+WHHrwb+Akgw==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 18 Dec 2018 16:58:05 GMT
	Content-Length: 77

	{ "response": [
		{
			"id": 2,
			"name": "Montreal",
			"division": {
				"id": 1,
				"name": "Quebec"
			}
		}
	]}
