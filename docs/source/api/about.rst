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

.. _to-api-about:

*********
``about``
*********

``GET``
=======
Retrieves information about the Traffic Ops server.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  ``undefined`` - this endpoint returns a custom JSON payload with top-level keys not inside of a ``response`` object

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:commitHash: The short hash ID of the Git commit at which this Traffic Ops instance was built
:commits:    The total number of Git commits up to the point at which this Traffic Ops instance was built
:goVersion:  The version of the ``go`` toolchain used to build this Traffic Ops instance
:release:    An Enterprise Linux release version - only has meaning for RPM-based installations
:name:       The name of the package as defined in the RPM used to install Traffic Ops
:RPMVersion: A full RPM version that identifies this Traffic Ops instance - this is a concatenation of "``name``-``Version``-``commits``.``commitHash``.``release``"
:Version:    The semantic version number of this Traffic Ops instance

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: AXASDxtVHtMilNhRsyMl4I5CBQhOSJWrN2i+TTgBiaCbrCdO2ejOVoWkuyvCYrv2wSKCDE9rsul6YJKnOZb+Lw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Dec 2018 18:26:43 GMT
	Content-Length: 171

	{
		"commitHash": "7685a12e",
		"commits": "9851",
		"goVersion": "go1.9.4",
		"release": "el7",
		"name": "traffic_ops",
		"RPMVersion": "traffic_ops-3.0.0-9851.7685a12e.el7",
		"Version": "3.0.0"
	}
