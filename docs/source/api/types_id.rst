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

.. _to-api-types-id:

****************
``types/{{ID}}``
****************

``GET``
=======
.. deprecated:: 1.1
	Use the ``id`` query parameter of a ``GET`` request to the :ref:`to-api-types` endpoint instead.

Retrieves a specific type of some things configured in Traffic Ops. Yes, that is as specific as a description of a 'type' can be.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------+
	| Name | Description                                                 |
	+======+=============================================================+
	|  ID  | The integral, unique identifier of the type being inspected |
	+------+-------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/types/48 HTTP/1.1
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
	Date: Wed, 12 Dec 2018 23:50:13 GMT
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

``PUT``
=======
Updates a specific type of thing configured in Traffic Ops.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------+
	| Name | Description                                                |
	+======+============================================================+
	|  ID  | The integral, unique identifier of the type being modified |
	+------+------------------------------------------------------------+

:description: A string of miscellaneous information describing the type
:name:        The new name of the type
:useInTable:  The name of the Traffic Ops database table that contains objects which are grouped, identified, or described by this type

	.. note:: This table need not actually exist - ``useInTable`` can be any string.

	.. note:: If ``useInTable`` is not present in the request payload, then an HTTP ``400 Bad Request`` response will be returned with an ``alerts`` object that erroneously reports: ``'use_in_table' cannot be blank``. The correct field name is ``useInTable``, do not be fooled by the error message. This bug is tracked by `GitHub Issue #3147 <https://github.com/apache/trafficcontrol/issues/3147>`_.

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/types/50 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 88
	Content-Type: application/json

	{
		"name": "quest",
		"description": "A test type for API examples",
		"useInTable": "quest"
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
	Whole-Content-Sha512: wSMLnSHTVfDY2O+CpYLFKGI1IpL+VlWx7RulbYDkU/jGbYXLumEaILYAhrXTgUr27IBFL+krmVAwYS/kiuoCqg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Dec 2018 18:15:56 GMT
	Content-Length: 200

	{ "alerts": [
		{
			"text": "type was updated.",
			"level": "success"
		}
	],
	"response": {
		"id": 50,
		"lastUpdated": "2018-12-19 18:15:56+00",
		"name": "quest",
		"description": "A test type for API examples",
		"useInTable": "quest"
	}}

``DELETE``
==========
Deletes a type of thing configured in Traffic Ops.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------+
	| Name | Description                                               |
	+======+===========================================================+
	|  ID  | The integral, unique identifier of the type being deleted |
	+------+-----------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/types/50 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

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
	Whole-Content-Sha512: uo2GrMVsT/dZaZYJqXA0pc0U+LvMGDlKhNWNAi2thA77eRE9Fzxnf0pb88cLSdoCWn5qYUwdIoGUXuU4Xd8uJQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Dec 2018 18:20:50 GMT
	Content-Length: 59

	{ "alerts": [
		{
			"text": "type was deleted.",
			"level": "success"
		}
	]}
