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

.. _to-api-deliveryservice_matches:

***************************
``deliveryservice_matches``
***************************

``GET``
=======
Retrieves a list of regular expressions that are used for routing :term:`Delivery Service`\ s.

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type:  Array

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:dsName: The 'xml_id' that uniquely identifies the :term:`Delivery Service` that uses ``patterns``

	.. warning:: This is **not** - as the name implies - the :term:`Delivery Service`'s display name (although the two are often - and should be - the same).

:patterns: An array of regular expressions used for routing the :term:`Delivery Service` identified by ``dsName``

.. note:: The response will only contain entries for :term:`Delivery Service`\ s that are active.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: TgzZrqpvxdm2Hz40nzvrfCZkKqRZTvRvvGacaEZYYKG4UgWZ+6qDDvnA6lGpu68Wz1tEqjC/p8HPa6oe6NiR8w==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 18 Dec 2018 17:23:54 GMT
	Content-Length: 56

	{ "response": [
		{
			"patterns": [
				".demo1."
			],
			"dsName": "demo1"
		}
	]}


.. [1] Users will only be able to query for the patterns associated with :term:`Delivery Service`\ s their tenant has permission to see.
