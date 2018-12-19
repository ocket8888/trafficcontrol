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

.. _to-api-statuses-id:

*******************
``statuses/{{ID}}``
*******************

``GET``
=======
Retrieves information about a particular :term:`Status`

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------+
	| Name | Description                                                           |
	+======+=======================================================================+
	| ID   | The integral, unique identifier of the :term:`Status` being inspected |
	+------+-----------------------------------------------------------------------+

.. table:: Request Query Parameters

	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                   |
	+===========+==========+===============================================================================================================+
	| orderby   | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` |
	|           |          | array                                                                                                         |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| sortOrder | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                      |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| limit     | no       | Choose the maximum number of results to return                                                                |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| offset    | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit          |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| page      | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long   |
	|           |          | and the first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be     |
	|           |          | defined to make use of ``page``.                                                                              |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/statuses/3 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:description: A short description of the status
:id:          The integral, unique identifier of this status
:lastUpdated: The date and time at which this status was last modified, in ISO format
:name:        The name of the status

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: dHNip9kpTGGS1w39/fWcFehNktgmXZus8XaufnmDpv0PyG/3fK/KfoCO3ZOj9V74/CCffps7doEygWeL/xRtKA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 21:04:20 GMT
	Content-Length: 150

	{ "response": [
		{
			"description": "Server is online and reported in the health protocol.",
			"id": 3,
			"lastUpdated": "2018-12-10 19:11:17+00",
			"name": "REPORTED"
		}
	]}

``PUT``
=======
Updates a status.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------+
	| Name |                Description                                   |
	+======+==============================================================+
	|  ID  | The integral, unique identifier of the status being modified |
	+------+--------------------------------------------------------------+

:description: An optional string containing miscellaneous information describing the status

	.. danger:: The endpoint will technically accept requests without this field, but such requests **will** *break the :ref:`to-api-statuses` and :ref:`to-api-statuses-id` endpoints*. For this reason it is **strongly advised** that this field always be present, even if it will only be an empty string. This bug is tracked by `GitHub Issue #3146 <https://github.com/apache/trafficcontrol/issues/3146>`_. Note that if this occurs, the bug can be fixed by deleting the status, but only if the integral, unique identifier of the status causing the problem is known - as it obviously can no longer be retrieved. Because Traffic Portal uses the now-broken endpoints in this scenario, Traffic Portal cannot be used to delete the problem status - it **must** be done by using the API directly.

:name: The new name of the status

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/statuses/7 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 66
	Content-Type: application/json

	{
		"name": "quest",
		"description": "A test status for API examples"
	}

Response Structure
------------------
:description: A short description of the status
:id:          The integral, unique identifier of this status
:lastUpdated: The date and time at which this status was last modified, in ISO format
:name:        The name of the status

.. code-block:: http
	:caption: Response Structure

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: dhlo8qW6Tw7dvjXK6RU8OhAPb3Z4TcSZW7ccZxvbNZMUfGDB9Yh5d4iV0GsmOMTMWxG/JSejFu0mg1mrABjHoQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Dec 2018 17:31:44 GMT
	Content-Length: 182

	{ "alerts": [
		{
			"text": "status was updated.",
			"level": "success"
		}
	],
	"response": {
		"description": "A test status for API examples",
		"id": 7,
		"lastUpdated": "2018-12-19 17:31:44+00",
		"name": "quest"
	}}

``DELETE``
==========
Deletes a status

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------+
	| Name |                Description                                  |
	+======+=============================================================+
	|  ID  | The integral, unique identifier of the status being deleted |
	+------+-------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/statuses/7 HTTP/1.1
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
	Whole-Content-Sha512: jyxrjaiCgmzWO1TGhNj1wdxIfkMd7WjaqOdsfH1FC1SnsbbnHGfefGQSM+k63vVldYOGjalhbr+4Vs44AV/dTw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Dec 2018 17:41:23 GMT
	Content-Length: 61

	{ "alerts": [
		{
			"text": "status was deleted.",
			"level": "success"
		}
	]}
