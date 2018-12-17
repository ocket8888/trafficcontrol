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

.. _to-api-keys-ping:

*************
``keys/ping``
*************
.. deprecated:: 1.1
	Use :ref:`to-api-riak-ping` instead.

``GET``
=======
Retrieves the status of the connected Traffic Vault instance.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object (string)

Request Structure
-----------------
No parameters available.

Response Structure
------------------
.. important:: If the Traffic Vault server is unreachable, the TCP connection will be dropped without sending a "FIN" to the client, resulting in no HTTP response at all. Thus, ``response`` should always be ``"OK"`` because no response will be received in the event that the server is *not* "OK".

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: dvkoZxvL+RKFeZSSksXyq71HCvEH4DeIeGosqt16d8bzbtaeix2+/URJM/YDp5jYaH9Gdq9PriyzwDJOrZRHOw==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 18 Dec 2018 16:30:58 GMT
	Content-Length: 17

	{
		"response": "OK"
	}
