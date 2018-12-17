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

.. _to-api-cdns-health:

***************
``cdns/health``
***************
Extract health information from all Cache Groups across all CDNs

.. seealso:: :ref:`health-proto`

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available

Response Structure
------------------
:cachegroups:  An array of objects describing the health of each Cache Group

	:name:    The name of the Cache Group
	:offline: The number of OFFLINE caches in the Cache Group
	:online:  The number of ONLINE caches in the Cache Group

:totalOffline: Total number of OFFLINE caches across all Cache Groups which are assigned to any CDN
:totalOnline:  Total number of ONLINE caches across all Cache Groups which are assigned to any CDN

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Mon, 17 Dec 2018 20:13:49 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Tue, 18 Dec 2018 00:13:49 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: KpXViXeAgch58ueQqdyU8NuINBw1EUedE6Rv2ewcLUajJp6kowdbVynpwW7XiSvAyHdtClIOuT3OkhIimghzSA==
	Content-Length: 115

	{ "response": {
		"totalOffline": 0,
		"totalOnline": 1,
		"cachegroups": [
			{
				"offline": 0,
				"name": "CDN_in_a_Box_Edge",
				"online": 1
			}
		]
	}}
