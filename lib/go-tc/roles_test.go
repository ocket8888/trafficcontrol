package tc

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"fmt"
	"testing"
	"time"
)

func ExampleRoleV40_String() {
	r := RoleV40{
		Description: "description",
		LastUpdated: time.Now(),
		Name:        "name",
		Permissions: []string{"test", "quest"},
	}

	fmt.Println(r)
	// Output: name
}

func TestRoleV4_Downgrade(t *testing.T) {
	r := RoleV40{
		Description: "description",
		LastUpdated: time.Now(),
		Name:        "name",
		Permissions: []string{"test", "quest"},
	}

	downgraded := r.Downgrade()

	if downgraded.Description == nil {
		t.Error("Description became nil after downgrade")
	} else if *downgraded.Description != r.Description {
		t.Errorf("Downgrading a RoleV4 with description '%s' yielded a Role with description '%s'", r.Description, *downgraded.Description)
	}
	if downgraded.Name == nil {
		t.Error("Name became nil after downgrade")
	} else if *downgraded.Name != r.Name {
		t.Errorf("Downgrading a RoleV4 with name '%s' yielded a Role with name '%s'", r.Name, *downgraded.Name)
	}
	if downgraded.Capabilities == nil {
		t.Error("Capabilities became nil after downgrade")
	} else if *downgraded.Capabilities == nil {
		t.Error("Capabilities became a reference to nil after downgrade")
	} else if len(*downgraded.Capabilities) != len(r.Permissions) {
		t.Errorf("Downgrading a RoleV4 with %d Permissions yielded a Role with %d Capabilities", len(r.Permissions), len(*downgraded.Capabilities))
	} else {
		// Order doesn't matter for this
		capabilities := make(map[string]struct{}, len(r.Permissions))
		for _, cap := range *downgraded.Capabilities {
			capabilities[cap] = struct{}{}
		}

		// check for duplicated values
		if len(capabilities) != len(r.Permissions) {
			t.Error("Downgrade caused a duplicate Capability to appear, somehow")
		} else {
			for _, perm := range r.Permissions {
				if _, ok := capabilities[perm]; !ok {
					t.Errorf("Permission '%s' not found in downgraded Role's Capabilities", perm)
				}
			}
		}
	}

	r.Permissions = nil
	downgraded = r.Downgrade()
	if downgraded.Capabilities == nil {
		t.Error("Capabilities became nil after downgrade with nil Permissions")
	} else if *downgraded.Capabilities != nil {
		t.Errorf("Expected downgrading a RoleV4 with nil Permissions to yield a Role with Capabilities being a reference to nil, but found a slice with %d elements", len(*downgraded.Capabilities))
	}
}

func TestRole_Upgrade(t *testing.T) {
	r := Role{
		Capabilities: new([]string),
		RoleV11: RoleV11{
			Description: new(string),
			ID:          nil,
			Name:        new(string),
			PrivLevel:   nil,
		},
	}

	*r.Description = "description"
	*r.Name = "name"
	*r.Capabilities = []string{"test", "quest"}

	upgraded := r.Upgrade()

	if upgraded.Description != *r.Description {
		t.Errorf("Upgrading a Role with description '%s' yielded a RoleV4 with description '%s'", *r.Description, upgraded.Description)
	}
	if upgraded.Name != *r.Name {
		t.Errorf("Upgrading a Role with name '%s' yielded a RoleV4 with name '%s'", *r.Name, upgraded.Name)
	}
	if upgraded.Permissions == nil {
		t.Error("Permissions became nil after upgrade")
	} else if *r.Capabilities == nil {
		t.Error("Capabilities was muted by upgrade: became a reference to nil")
	} else if len(upgraded.Permissions) != len(*r.Capabilities) {
		t.Errorf("Upgrading a Role with %d Capabilities yielded a RoleV4 with %d Permissions", len(*r.Capabilities), len(upgraded.Permissions))
	} else {
		// Order doesn't matter for this
		permissions := make(map[string]struct{}, len(*r.Capabilities))
		for _, perm := range upgraded.Permissions {
			permissions[perm] = struct{}{}
		}

		// check for duplicated values
		if len(permissions) != len(*r.Capabilities) {
			t.Error("Upgrade caused a duplicate Permission to appear, somehow")
		} else {
			for _, cap := range *r.Capabilities {
				if _, ok := permissions[cap]; !ok {
					t.Errorf("Capability '%s' not found in upgraded RoleV4's Permissions", cap)
				}
			}
		}
	}

	*r.Capabilities = nil
	upgraded = r.Upgrade()
	if upgraded.Permissions == nil {
		t.Errorf("Expected upgrading a Role with Capabilities being a reference to nil to yield a RoleV4 with nil Permissions, but found a slice with %d elements", len(upgraded.Permissions))
	}

	r.Capabilities = nil
	upgraded = r.Upgrade()
	if upgraded.Permissions != nil {
		t.Errorf("Expected upgrading a Role with nil Capabilities to yield a RoleV4 with nil Permissions, but found a slice with %d elements", len(upgraded.Permissions))
	}
}
