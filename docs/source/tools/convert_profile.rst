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

.. program:: convert_profile
.. _convert_profile:

***************
convert_profile
***************
:program:`convert_profile` is a tool for converting :abbr:`ATC (Apache Traffic Control)` :term:`Profiles` meant for use with :abbr:`ATS (Apache Traffic Server)` :term:`cache servers` using older versions of :abbr:`ATS (Apache Traffic Server)` to use newer options compatible with newer versions of :abbr:`ATS (Apache Traffic Server)`. Specifically, a ruleset is given in :atc-file:`tools/convert_profile/convert622to713.json` (and in YAML format in :atc-file:`tools/convert_profile/convert622to713.yaml) for converting :term:`Profiles` made for version 6.22 (or possibly lower) to use newer options compatible with :abbr:`ATS (Apache Traffic Server)` 7.13.

Building
========
:program:`convert_profile` is located in :atc-file:`tools/convert_profile` and built using Go, so to build just use :manpage:`go-build(1)`.

Usage
=====
``convert_profile -i IN -r RULES [-o OUT] [-f]``

Options
-------

.. option:: -i IN, --input_profile IN

	Reads the file ``IN`` as an :abbr:`ATC (Apache Traffic Control)` :term:`Profile` object encoded in JSON. This option is not optional; there must be input to process.

.. option:: -r RULES, --rules RULES

	Reads the file ``RULES`` as a set of rules to use when converting :term:`Profiles`. It may be encoded in JSON or YAML. This option is not optional; there must be rules to use when converting.

.. option:: -o OUT, --out OUT

	Writes the resulting :term:`Profile` object to the file given by ``OUT``. If not given, the result is printed to stdout.

.. option:: -f, --force

	If given, existing values in the input :term:`Profile` will be "clobbered", replacing them with the recommendations given by the rules even if the value doesn't match a pattern given by the rule set.

Ruleset Format
==============
Rulesets may be encoded in either YAML or JSON format. In any case, the important keys and their respective values are as follows.

.. tip:: Because of JSON parsing rules, backslashes in regular expressions will need to be escaped. For example, a string that encodes the regular expression :regexp:`\\..+` in JSON would look like: ``"\\..+"``. This does not apply to YAML encoding.

:description: This isn't mandated - or even parsed - by :program:`convert_profile`, but the default conversion ruleset comes with this key set to a description of the conversion being performed, and it is suggested that such descriptions be provided in all rulesets to help people understand what they do.
:conversion_actions: An array of rules for doing replacements on :term:`Parameters`.

	:action:          An optional string defining an alternative action to take instead of doing replacements. If this is specified, ``new_name``, ``new_config_file``, and ``new_action`` are ignored. Currently, the only allowed value is "delete" which will cause any matching :term:`Parameters` to be removed from the resulting :term:`Profile`.
	:match_parameter: An object containing regular expressions to match against the properties of :term:`Parameters` to determine if the rule's replacements should be done.

		:config_file: A regular expression that matches the :ref:`parameter-config-file` of a :term:`Parameter`.
		:name:        A regular expression that matches the :ref:`parameter-name` of a :term:`Parameter`.
		:value:       A regular expression that matches the :ref:`parameter-value` of a :term:`Parameter`.

	:new_config_file: An optional string that gives a new value to use for the resulting :term:`Parameter`'s :ref:`parameter-config-file`. If not given, the value is unchanged.
	:new_name:        An optional string that gives a new value to use for the resulting :term:`Parameter`'s :ref:`parameter-name`. If not given, the value is unchanged.
	:new_value:       An optional string that gives a new value to use for the resulting :term:`Parameter`'s :ref:`parameter-value`. If not given, the value is unchanged.

:replace_description: An optional definition of a word or phrase in the :term:`Profile`'s :ref:`profile-description` to be replaced with another word or phrase.

	:new: The new value to be used instead of ``old`` wherever ``old`` is found within the :ref:`profile-description`.
	:old: A string which will be removed from the :ref:`profile-description` and replaced with ``new``\ [#emptyold]_.

:replace_name: An optional definition of a word or phrase in the :term:`Profile`'s :ref:`profile-name` to be replaced with another word or phrase.

	:new: The new value to be used instead of ``old`` wherever ``old`` is found within the :ref:`profile-name`.
	:old: A string which will be removed from the :ref:`profile-name` and replaced with ``new``\ [#emptyold]_.

:validate_parameters: An array of rules that perform no replacements, but simply validate that the matching :term:`Parameters` exist and match the respective property regular expressions.

	:config_file: A regular expression that matches the :ref:`parameter-config-file` of a :term:`Parameter`.
	:name:        A regular expression that matches the :ref:`parameter-name` of a :term:`Parameter`.
	:value:       A regular expression that matches the :ref:`parameter-value` of a :term:`Parameter`.

Testing
=======
Unit testing is available via :manpage:`go-test(1)`.

.. [#emptyold]  If this is an empty string, then ``new`` will be inserted before every character and at the end. Unless ``new`` is also an empty string, in which case nothing happens. This is a result of how Go's :godoc:`strings.Replace` function works.
