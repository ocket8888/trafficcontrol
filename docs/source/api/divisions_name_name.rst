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

.. _to-api-divisions-name-name:

***************************
``divisions/name/{{name}}``
***************************

``GET``
=======
.. deprecated:: 1.1
	Use the ``name`` query parameter of a ``GET`` request to :ref:`to-api-divisions` instead.

Get information about a specific Division.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------+
	| Name | Description                              |
	+======+==========================================+
	| name | The name of the Division to be inspected |
	+------+------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.1/divisions/name/USA HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:id:          An integral, unique identifier for this Division
:lastUpdated: The date and time at which this Division was last modified, in ISO format
:name:        The Division name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: RmIoHRrkiUNiVew2/G5A5ZA+1x0lB+JGH3VPHD8nBFfjZOIsYtAImsU+JFggBIotVM1qeciStuLgV1DFKJiQDw==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 17 Dec 2018 21:33:03 GMT
	Content-Length: 75

	{ "response": [
		{
			"id": 2,
			"lastUpdated": "2018-12-17 19:26:03+00",
			"name": "USA"
		}
	]}
