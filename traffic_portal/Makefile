# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

# Directories
RESOURCESDIR := app/dist/public/resources
SRCDIR := app/src

# Static Resources
STATIC := $(patsubst $(SRCDIR)/%,$(RESOURCESDIR)/%,$(shell find $(SRCDIR)/assets/ -type f))
STATIC += app/dist/public/index.html
STATIC += app/dist/public/traffic_portal_release.json
STATIC += app/dist/public/traffic_portal_properties.json

# Compiled CSS files
STYLESHEETS := $(patsubst $(SRCDIR)/%.scss,$(RESOURCESDIR)/%.css,$(wildcard $(SRCDIR)/styles/*.scss))

# Template sources
TEMPLATES := $(shell find $(SRCDIR)/ -name "*.tpl.html" -type f)

# Current UNIX timestamp
NOW := $(shell date +%s)

# SASS compiler
SASS := $(shell which sass)
ifeq ('SASS',)
$(error Could not find a valid SASS compiler!)
endif

SASS_FLAGS ?= --scss -I app/bower_components/
ifdef DEV
SASS_FLAGS += --style=expanded --trace
else
SASS_FLAGS += --style=compressed
endif

.PHONY: clean all dist

all: $(STATIC) $(STYLESHEETS)

dist: all $(SRCDIR)/package.json
	@mkdir -p $(@D)
	cp -f $(SRCDIR)/package.json app/dist/package.json

## Rules to 'build' static assets ##
$(RESOURCESDIR)/assets/fonts/%: $(SRCDIR)/assets/fonts/%
	@mkdir -p $(@D)
	cp -f $< $@

$(RESOURCESDIR)/assets/css/%.css: $(SRCDIR)/assets/css/%.css
	@mkdir -p $(@D)
	cp -f $< $@

$(RESOURCESDIR)/assets/images/%.png: $(SRCDIR)/assets/images/%.png
	@mkdir -p $(@D)
	cp -f $< $@

$(RESOURCESDIR)/assets/js/%.js: $(SRCDIR)/assets/js/%.js
	@mkdir -p $(@D)
	cp -f $< $@

app/dist/public/index.html: $(SRCDIR)/index.html
	@mkdir -p $(@D)
	#sed -re 's/'
	cp -f $< $@

app/dist/public/%.json: $(SRCDIR)/%.json
	@mkdir -p $(@D)
	cp -f $< $@

####################################


## Rules to build CSS from SCSS files ##
$(RESOURCESDIR)/styles/%.css: $(SRCDIR)/styles/%.scss
	@mkdir -p $(@D)
	$(SASS) $(SASS_FLAGS) $< $@

#########################################

## Rules to build app templates ##
.tmp/app-templates.js: $(TEMPLATES)

clean:
	$(RM) -r app/dist .tmp/
