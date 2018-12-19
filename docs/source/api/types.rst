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

.. _to-api-types:

*********
``types``
*********

``GET``
=======
Retrieves all of the types of things configured in Traffic Ops. Yes, that is as specific as a description of a 'type' can be.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------------+----------+--------------------------------------------------------------------------------------------------------------------------------+
	|    Name    | Required |                Description                                                                                                     |
	+============+==========+================================================================================================================================+
	|     id     | no       | Return only the type that is identified by this integral, unique identifier                                                    |
	+------------+----------+--------------------------------------------------------------------------------------------------------------------------------+
	|    name    | no       | Return only types with this name                                                                                               |
	+------------+----------+--------------------------------------------------------------------------------------------------------------------------------+
	| useInTable | no       | Return only types that are used to identify the type of the object stored in the Traffic Ops database table that has this name |
	+------------+----------+--------------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Structure

	GET /api/1.4/types?name=TC_LOC HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:description: A short description of this type
:id:          An integral, unique identifier for this type
:lastUpdated: The date and time at which this type was last updated, in ISO format
:name:        The name of this type
:useInTable:  The name of the Traffic Ops database table that contains objects which are grouped, identified, or described by this type

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: EH8jo8OrCu79Tz9xpgT3YRyKJ/p2NcTmbS3huwtqRByHz9H6qZLQjA59RIPaVSq3ZxsU6QhTaox5nBkQ9LPSAA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 12 Dec 2018 22:59:22 GMT
	Content-Length: 168

	{ "response": [
		{
			"id": 48,
			"lastUpdated": "2018-12-12 16:26:41+00",
			"name": "TC_LOC",
			"description": "Location for Traffic Control Component Servers",
			"useInTable": "cachegroup"
		}
	]}

``POST``
========
Creates a new type

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
:description: A string of miscellaneous information describing the new type
:name:        The name of the new type
:useInTable:  The name of the Traffic Ops database table that contains objects which are grouped, identified, or described by this type

	.. note:: This table need not actually exist - ``useInTable`` can be any string.

	.. note:: If ``useInTable`` is not present in the request payload, then an HTTP ``400 Bad Request`` response will be returned with an ``alerts`` object that erroneously reports: ``'use_in_table' cannot be blank``. The correct field name is ``useInTable``, do not be fooled by the error message. This bug is tracked by `GitHub Issue #3147 <https://github.com/apache/trafficcontrol/issues/3147>`_.

.. code-block:: http
	:caption: Request Example

	POST /api/1.1/types HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 86
	Content-Type: application/json

	{
		"name": "test",
		"description": "A test type for API examples",
		"useInTable": "test"
	}

Response Structure
------------------
:description: A short description of this type
:id:          An integral, unique identifier for this type
:lastUpdated: The date and time at which this type was last updated, in ISO format
:name:        The name of this type
:useInTable:  The name of the Traffic Ops database table that contains objects which are grouped, identified, or described by this type

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: emF1D5RjeDhb8dylrs1oJh3MhWhIAF3t+eHzAXKQmSgMYXud5Cg9M8F7uaWZwczi0GQrbV7YUj4isQQg76VGVg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Dec 2018 17:50:45 GMT
	Content-Length: 198

	{ "alerts": [
		{
			"text": "type was created.",
			"level": "success"
		}
	],
	"response": {
		"id": 50,
		"lastUpdated": "2018-12-19 17:50:45+00",
		"name": "test",
		"description": "A test type for API examples",
		"useInTable": "test"
	}}
