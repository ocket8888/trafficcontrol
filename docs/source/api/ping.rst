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

.. _to-api-ping:

********
``ping``
********

``GET``
=======
This endpoint is merely used as an HTTP equivalent of the :manpage:`ping(8)` UNIX command-line utility; i.e. it checks that the server is running and reachable via HTTP.

:Auth. Required: No
:Roles Required: None\ [1]_
:Response Type:  ``undefined``

Request Structure
-----------------
No parameters available

Response Structure
------------------
This endpoint includes the non-standard, top-level object ``ping`` rather than ``response`` or ``alerts``.

:ping: A string which for all successful responses will be ``"pong"``

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Whole-Content-Sha512: 0HqcLcYHCB4AFYGFzcAsP2h+PCMlYxk/TqMajcy3fWCzY730Tv32k5trUaJLeSBbRx2FUi7z/sTAuzikdg0E4g==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 18 Dec 2018 15:45:11 GMT
	Content-Length: 16
	Content-Type: text/plain; charset=utf-8

	{
		"ping": "pong"
	}

.. [1] Because authentication is not required, users that cannot log in or clients which may not be registered users at all will be able to use this endpoint. This is technically more permissive than "None" for required roles.
