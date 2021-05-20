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

import "time"

// RolesResponse is a list of Roles as a response.
// swagger:response RolesResponse
// in: body
type RolesResponse struct {
	// in: body
	Response []Role `json:"response"`
	Alerts
}

// RoleResponse is a single Role response for Update and Create to depict what
// changed.
// swagger:response RoleResponse
// in: body
type RoleResponse struct {
	// in: body
	Response Role `json:"response"`
	Alerts
}

// Role ...
type Role struct {
	RoleV11

	// Capabilities associated with the Role
	//
	// required: true
	Capabilities *[]string `json:"capabilities" db:"-"`
}

// RoleV11 ...
type RoleV11 struct {
	// ID of the Role
	//
	// required: true
	ID *int `json:"id" db:"id"`

	// Name of the Role
	//
	// required: true
	Name *string `json:"name" db:"name"`

	// Description of the Role
	//
	// required: true
	Description *string `json:"description" db:"description"`

	// Priv Level of the Role
	//
	// required: true
	PrivLevel *int `json:"privLevel" db:"priv_level"`
}

// RoleV40 is a Role as it appears in Traffic Ops API version 4.0.
type RoleV40 struct {
	// Description of the Role
	Description string `json:"description"`
	// Date/time at which the Role was last updated.
	LastUpdated time.Time `json:"lastUpdated"`
	// Name of the Role.
	Name string `json:"name"`
	// Permissions afforded to users with this Role.
	Permissions []string `json:"permissions"`
}

// String implements the fmt.Stringer interface to return a textual
// representation of the Role.
func (r RoleV40) String() string {
	return r.Name
}

// RoleV4 is a Role as it appears in the latest version of Traffic Ops API
// version 4.0.
type RoleV4 = RoleV40

// Downgrade converts a Role from API versions 4.0 and later to a Role as it
// appeared in API versions between 1.3 and 3.x (inclusive).
//
// This makes a "deep" copy of the Role, so the return value and the original
// Role can both be freely manipulated without affecting each other.
func (r RoleV4) Downgrade() Role {
	role := Role{
		Capabilities: new([]string),
		RoleV11: RoleV11{
			ID:          nil,
			Name:        new(string),
			Description: new(string),
			PrivLevel:   nil,
		},
	}

	*role.Description = r.Description
	*role.Name = r.Name

	if r.Permissions == nil {
		*role.Capabilities = nil
	} else {
		*role.Capabilities = make([]string, len(r.Permissions))
		copy(*role.Capabilities, r.Permissions)
	}

	return role
}

// Upgrade converts a Role from API versions between 1.3 and 3.x (inclusive) to
// a Role as it appears in API version 4.0 and later.
//
// This makes a "deep" copy of the Role, so the return value and the original
// Role can both be freely manipulated without affecting each other.
func (r Role) Upgrade() RoleV4 {
	var role RoleV40

	if r.Name != nil {
		role.Name = *r.Name
	}
	if r.Description != nil {
		role.Description = *r.Description
	}

	if r.Capabilities == nil {
		role.Permissions = nil
	} else {
		role.Permissions = make([]string, len(*r.Capabilities))
		copy(role.Permissions, *r.Capabilities)
	}

	return role
}

// RolesResponseV4 is the type of a response from Traffic Ops to a GET request
// made to its /roles API endpoint in the latest minor version of API version 4.
type RolesResponseV4 struct {
	Response []RoleV4 `json:"response"`
	Alerts
}

// RoleResponseV4 is the type of a response from Traffic Ops to a PUT, POST or
// DELETE request made to its /roles API endpoint in the latest minor version
// of API version 4.
type RoleResponseV4 struct {
	Response RoleV4 `json:"response"`
	Alerts
}
