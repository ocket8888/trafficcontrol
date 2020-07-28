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

import "fmt"
import "strings"
import "testing"
import "time"

func ExampleLegacyInterfaceDetails_ToInterfaces() {
	lid := LegacyInterfaceDetails{
		InterfaceMtu:  new(int),
		InterfaceName: new(string),
		IP6Address:    new(string),
		IP6Gateway:    new(string),
		IPAddress:     new(string),
		IPGateway:     new(string),
		IPNetmask:     new(string),
	}
	*lid.InterfaceMtu = 9000
	*lid.InterfaceName = "test"
	*lid.IP6Address = "::14/64"
	*lid.IP6Gateway = "::15"
	*lid.IPAddress = "1.2.3.4"
	*lid.IPGateway = "4.3.2.1"
	*lid.IPNetmask = "255.255.255.252"

	ifaces, err := lid.ToInterfaces(true, false)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	for _, iface := range ifaces {
		fmt.Printf("name=%s, monitor=%t\n", iface.Name, iface.Monitor)
		for _, ip := range iface.IPAddresses {
			fmt.Printf("\taddr=%s, gateway=%s, service address=%t\n", ip.Address, *ip.Gateway, ip.ServiceAddress)
		}
	}
	// Output: name=test, monitor=false
	// 	addr=1.2.3.4/30, gateway=4.3.2.1, service address=true
	// 	addr=::14/64, gateway=::15, service address=false
	//
}

func TestServer_ToNullable(t *testing.T) {
	fqdn := "testFQDN"
	srv := Server{
		Cachegroup:       "testCachegroup",
		CachegroupID:     42,
		CDNID:            43,
		CDNName:          "testCDNName",
		DeliveryServices: map[string][]string{"test": []string{"quest"}},
		DomainName:       "testDomainName",
		FQDN:             &fqdn,
		FqdnTime:         time.Now(),
		GUID:             "testGUID",
		HostName:         "testHostName",
		HTTPSPort:        -1,
		ID:               44,
		ILOIPAddress:     "testILOIPAddress",
		ILOIPGateway:     "testILOIPGateway",
		ILOIPNetmask:     "testILOIPNetmask",
		ILOPassword:      "testILOPassword",
		ILOUsername:      "testILOUsername",
		InterfaceMtu:     -2,
		InterfaceName:    "testInterfaceName",
		IP6Address:       "testIP6Address",
		IP6IsService:     true,
		IP6Gateway:       "testIP6Gateway",
		IPAddress:        "testIPAddress",
		IPIsService:      false,
		IPGateway:        "testIPGateway",
		IPNetmask:        "testIPNetmask",
		LastUpdated:      TimeNoMod(Time{Time: time.Now().Add(time.Minute), Valid: true}),
		MgmtIPAddress:    "testMgmtIPAddress",
		MgmtIPGateway:    "testMgmtIPGateway",
		MgmtIPNetmask:    "testMgmtIPNetmask",
		OfflineReason:    "testOfflineReason",
		PhysLocation:     "testPhysLocation",
		PhysLocationID:   45,
		Profile:          "testProfile",
		ProfileDesc:      "testProfileDesc",
		ProfileID:        46,
		Rack:             "testRack",
		RevalPending:     true,
		RouterHostName:   "testRouterHostName",
		RouterPortName:   "testRouterPortName",
		Status:           "testStatus",
		StatusID:         47,
		TCPPort:          -3,
		Type:             "testType",
		TypeID:           48,
		UpdPending:       false,
		XMPPID:           "testXMPPID",
		XMPPPasswd:       "testXMPPasswd",
	}

	nullable := srv.ToNullable()

	if nullable.Cachegroup == nil {
		t.Error("nullable conversion gave nil Cachegroup")
	} else if *nullable.Cachegroup != srv.Cachegroup {
		t.Errorf("Incorrect Cachegroup after nullable conversion; want: '%s', got: '%s'", srv.Cachegroup, *nullable.Cachegroup)
	}

	if nullable.CachegroupID == nil {
		t.Error("nullable conversion gave nil CachegroupID")
	} else if *nullable.CachegroupID != srv.CachegroupID {
		t.Errorf("Incorrect CachegroupID after nullable conversion; want: %d, got: %d", srv.CachegroupID, *nullable.CachegroupID)
	}

	if nullable.CDNID == nil {
		t.Error("nullable conversion gave nil CDNID")
	} else if *nullable.CDNID != srv.CDNID {
		t.Errorf("Incorrect CDNID after nullable conversion; want: %d, got: %d", srv.CDNID, *nullable.CDNID)
	}

	if nullable.CDNName == nil {
		t.Error("nullable conversion gave nil CDNName")
	} else if *nullable.CDNName != srv.CDNName {
		t.Errorf("Incorrect CDNName after nullable conversion; want: '%s', got: '%s'", srv.CDNName, *nullable.CDNName)
	}

	if nullable.DeliveryServices == nil {
		t.Error("nullable conversion gave nil DeliveryServices")
	} else if len(*nullable.DeliveryServices) != len(srv.DeliveryServices) {
		t.Errorf("Incorrect number of DeliveryServices after nullable conversion; want: %d, got: %d", len(srv.DeliveryServices), len(*nullable.DeliveryServices))
	} else {
		for k, v := range srv.DeliveryServices {
			nullableV, ok := (*nullable.DeliveryServices)[k]
			if !ok {
				t.Errorf("Missing Delivery Service '%s' after nullable conversion", k)
				continue
			}
			if len(nullableV) != len(v) {
				t.Errorf("Delivery Service '%s' has incorrect length after nullable conversion; want: %d, got: %d", k, len(v), len(nullableV))
			}
			for i, ds := range v {
				nullableDS := nullableV[i]
				if nullableDS != ds {
					t.Errorf("Incorrect value at position %d in Delivery Service '%s' after nullable conversion; want: '%s', got: '%s'", i, k, ds, nullableDS)
				}
			}
		}
	}

	if nullable.DomainName == nil {
		t.Error("nullable conversion gave nil DomainName")
	} else if *nullable.DomainName != srv.DomainName {
		t.Errorf("Incorrect DomainName after nullable conversion; want: '%s', got: '%s'", srv.DomainName, *nullable.DomainName)
	}

	if nullable.FQDN == nil {
		t.Error("nullable conversion gave nil FQDN")
	} else if *nullable.FQDN != fqdn {
		t.Errorf("Incorrect FQDN after nullable conversion; want: '%s', got: '%s'", fqdn, *nullable.FQDN)
	}

	if nullable.FqdnTime != srv.FqdnTime {
		t.Errorf("Incorrect FqdnTime after nullable conversion; want: '%s', got: '%s'", srv.FqdnTime, nullable.FqdnTime)
	}

	if nullable.GUID == nil {
		t.Error("nullable conversion gave nil GUID")
	} else if *nullable.GUID != srv.GUID {
		t.Errorf("Incorrect GUID after nullable conversion; want: '%s', got: '%s'", srv.GUID, *nullable.GUID)
	}

	if nullable.HostName == nil {
		t.Error("nullable conversion gave nil HostName")
	} else if *nullable.HostName != srv.HostName {
		t.Errorf("Incorrect HostName after nullable conversion; want: '%s', got: '%s'", srv.HostName, *nullable.HostName)
	}

	if nullable.HTTPSPort == nil {
		t.Error("nullable conversion gave nil HTTPSPort")
	} else if *nullable.HTTPSPort != srv.HTTPSPort {
		t.Errorf("Incorrect HTTPSPort after nullable conversion; want: %d, got: %d", srv.HTTPSPort, *nullable.HTTPSPort)
	}

	if nullable.ID == nil {
		t.Error("nullable conversion gave nil ID")
	} else if *nullable.ID != srv.ID {
		t.Errorf("Incorrect ID after nullable conversion; want: %d, got: %d", srv.ID, *nullable.ID)
	}

	if nullable.ILOIPAddress == nil {
		t.Error("nullable conversion gave nil ILOIPAddress")
	} else if *nullable.ILOIPAddress != srv.ILOIPAddress {
		t.Errorf("Incorrect ILOIPAddress after nullable conversion; want: '%s', got: '%s'", srv.ILOIPAddress, *nullable.ILOIPAddress)
	}

	if nullable.ILOIPGateway == nil {
		t.Error("nullable conversion gave nil ILOIPGateway")
	} else if *nullable.ILOIPGateway != srv.ILOIPGateway {
		t.Errorf("Incorrect ILOIPGateway after nullable conversion; want: '%s', got: '%s'", srv.ILOIPGateway, *nullable.ILOIPGateway)
	}

	if nullable.ILOIPNetmask == nil {
		t.Error("nullable conversion gave nil ILOIPNetmask")
	} else if *nullable.ILOIPNetmask != srv.ILOIPNetmask {
		t.Errorf("Incorrect ILOIPNetmask after nullable conversion; want: '%s', got: '%s'", srv.ILOIPNetmask, *nullable.ILOIPNetmask)
	}

	if nullable.ILOPassword == nil {
		t.Error("nullable conversion gave nil ILOPassword")
	} else if *nullable.ILOPassword != srv.ILOPassword {
		t.Errorf("Incorrect ILOPassword after nullable conversion; want: '%s', got: '%s'", srv.ILOPassword, *nullable.ILOPassword)
	}

	if nullable.ILOUsername == nil {
		t.Error("nullable conversion gave nil ILOUsername")
	} else if *nullable.ILOUsername != srv.ILOUsername {
		t.Errorf("Incorrect ILOUsername after nullable conversion; want: '%s', got: '%s'", srv.ILOUsername, *nullable.ILOUsername)
	}

	if nullable.InterfaceMtu == nil {
		t.Error("nullable conversion gave nil InterfaceMtu")
	} else if *nullable.InterfaceMtu != srv.InterfaceMtu {
		t.Errorf("Incorrect InterfaceMtu after nullable conversion; want: %d, got: %d", srv.InterfaceMtu, *nullable.InterfaceMtu)
	}

	if nullable.InterfaceName == nil {
		t.Error("nullable conversion gave nil InterfaceName")
	} else if *nullable.InterfaceName != srv.InterfaceName {
		t.Errorf("Incorrect InterfaceName after nullable conversion; want: '%s', got: '%s'", srv.InterfaceName, *nullable.InterfaceName)
	}

	if nullable.IP6Address == nil {
		t.Error("nullable conversion gave nil IP6Address")
	} else if *nullable.IP6Address != srv.IP6Address {
		t.Errorf("Incorrect IP6Address after nullable conversion; want: '%s', got: '%s'", srv.IP6Address, *nullable.IP6Address)
	}

	if nullable.IP6IsService == nil {
		t.Error("nullable conversion gave nil IP6IsService")
	} else if *nullable.IP6IsService != srv.IP6IsService {
		t.Errorf("Incorrect IP6IsService after nullable conversion; want: %t, got: %t", srv.IP6IsService, *nullable.IP6IsService)
	}

	if nullable.IP6Gateway == nil {
		t.Error("nullable conversion gave nil IP6Gateway")
	} else if *nullable.IP6Gateway != srv.IP6Gateway {
		t.Errorf("Incorrect IP6Gateway after nullable conversion; want: '%s', got: '%s'", srv.IP6Gateway, *nullable.IP6Gateway)
	}

	if nullable.IPAddress == nil {
		t.Error("nullable conversion gave nil IPAddress")
	} else if *nullable.IPAddress != srv.IPAddress {
		t.Errorf("Incorrect IPAddress after nullable conversion; want: '%s', got: '%s'", srv.IPAddress, *nullable.IPAddress)
	}

	if nullable.IPIsService == nil {
		t.Error("nullable conversion gave nil IPIsService")
	} else if *nullable.IPIsService != srv.IPIsService {
		t.Errorf("Incorrect IPIsService after nullable conversion; want: %t, got: %t", srv.IPIsService, *nullable.IPIsService)
	}

	if nullable.IPGateway == nil {
		t.Error("nullable conversion gave nil IPGateway")
	} else if *nullable.IPGateway != srv.IPGateway {
		t.Errorf("Incorrect IPGateway after nullable conversion; want: '%s', got: '%s'", srv.IPGateway, *nullable.IPGateway)
	}

	if nullable.IPNetmask == nil {
		t.Error("nullable conversion gave nil IPNetmask")
	} else if *nullable.IPNetmask != srv.IPNetmask {
		t.Errorf("Incorrect IPNetmask after nullable conversion; want: '%s', got: '%s'", srv.IPNetmask, *nullable.IPNetmask)
	}

	if nullable.LastUpdated == nil {
		t.Error("nullable conversion gave nil LastUpdated")
	} else if *nullable.LastUpdated != srv.LastUpdated {
		t.Errorf("Incorrect LastUpdated after nullable conversion; want: '%s', got: '%s'", srv.LastUpdated, *nullable.LastUpdated)
	}

	if nullable.MgmtIPAddress == nil {
		t.Error("nullable conversion gave nil MgmtIPAddress")
	} else if *nullable.MgmtIPAddress != srv.MgmtIPAddress {
		t.Errorf("Incorrect MgmtIPAddress after nullable conversion; want: '%s', got: '%s'", srv.MgmtIPAddress, *nullable.MgmtIPAddress)
	}

	if nullable.MgmtIPGateway == nil {
		t.Error("nullable conversion gave nil MgmtIPGateway")
	} else if *nullable.MgmtIPGateway != srv.MgmtIPGateway {
		t.Errorf("Incorrect MgmtIPGateway after nullable conversion; want: '%s', got: '%s'", srv.MgmtIPGateway, *nullable.MgmtIPGateway)
	}

	if nullable.MgmtIPNetmask == nil {
		t.Error("nullable conversion gave nil MgmtIPNetmask")
	} else if *nullable.MgmtIPNetmask != srv.MgmtIPNetmask {
		t.Errorf("Incorrect MgmtIPNetmask after nullable conversion; want: '%s', got: '%s'", srv.MgmtIPNetmask, *nullable.MgmtIPNetmask)
	}

	if nullable.OfflineReason == nil {
		t.Error("nullable conversion gave nil OfflineReason")
	} else if *nullable.OfflineReason != srv.OfflineReason {
		t.Errorf("Incorrect OfflineReason after nullable conversion; want: '%s', got: '%s'", srv.OfflineReason, *nullable.OfflineReason)
	}

	if nullable.PhysLocation == nil {
		t.Error("nullable conversion gave nil PhysLocation")
	} else if *nullable.PhysLocation != srv.PhysLocation {
		t.Errorf("Incorrect PhysLocation after nullable conversion; want: '%s', got: '%s'", srv.PhysLocation, *nullable.PhysLocation)
	}

	if nullable.PhysLocationID == nil {
		t.Error("nullable conversion gave nil PhysLocationID")
	} else if *nullable.PhysLocationID != srv.PhysLocationID {
		t.Errorf("Incorrect PhysLocationID after nullable conversion; want: %d, got: %d", srv.PhysLocationID, *nullable.PhysLocationID)
	}

	if nullable.Profile == nil {
		t.Error("nullable conversion gave nil Profile")
	} else if *nullable.Profile != srv.Profile {
		t.Errorf("Incorrect Profile after nullable conversion; want: '%s', got: '%s'", srv.Profile, *nullable.Profile)
	}

	if nullable.ProfileDesc == nil {
		t.Error("nullable conversion gave nil ProfileDesc")
	} else if *nullable.ProfileDesc != srv.ProfileDesc {
		t.Errorf("Incorrect ProfileDesc after nullable conversion; want: '%s', got: '%s'", srv.ProfileDesc, *nullable.ProfileDesc)
	}

	if nullable.ProfileID == nil {
		t.Error("nullable conversion gave nil ProfileID")
	} else if *nullable.ProfileID != srv.ProfileID {
		t.Errorf("Incorrect ProfileID after nullable conversion; want: %d, got: %d", srv.ProfileID, *nullable.ProfileID)
	}

	if nullable.Rack == nil {
		t.Error("nullable conversion gave nil Rack")
	} else if *nullable.Rack != srv.Rack {
		t.Errorf("Incorrect Rack after nullable conversion; want: '%s', got: '%s'", srv.Rack, *nullable.Rack)
	}

	if nullable.RevalPending == nil {
		t.Error("nullable conversion gave nil RevalPending")
	} else if *nullable.RevalPending != srv.RevalPending {
		t.Errorf("Incorrect RevalPending after nullable conversion; want: %t, got: %t", srv.RevalPending, *nullable.RevalPending)
	}

	if nullable.RouterHostName == nil {
		t.Error("nullable conversion gave nil RouterHostName")
	} else if *nullable.RouterHostName != srv.RouterHostName {
		t.Errorf("Incorrect RouterHostName after nullable conversion; want: '%s', got: '%s'", srv.RouterHostName, *nullable.RouterHostName)
	}

	if nullable.RouterPortName == nil {
		t.Error("nullable conversion gave nil RouterPortName")
	} else if *nullable.RouterPortName != srv.RouterPortName {
		t.Errorf("Incorrect RouterPortName after nullable conversion; want: '%s', got: '%s'", srv.RouterPortName, *nullable.RouterPortName)
	}

	if nullable.Status == nil {
		t.Error("nullable conversion gave nil Status")
	} else if *nullable.Status != srv.Status {
		t.Errorf("Incorrect Status after nullable conversion; want: '%s', got: '%s'", srv.Status, *nullable.Status)
	}

	if nullable.StatusID == nil {
		t.Error("nullable conversion gave nil StatusID")
	} else if *nullable.StatusID != srv.StatusID {
		t.Errorf("Incorrect StatusID after nullable conversion; want: %d, got: %d", srv.StatusID, *nullable.StatusID)
	}

	if nullable.TCPPort == nil {
		t.Error("nullable conversion gave nil TCPPort")
	} else if *nullable.TCPPort != srv.TCPPort {
		t.Errorf("Incorrect TCPPort after nullable conversion; want: %d, got: %d", srv.TCPPort, *nullable.TCPPort)
	}

	if nullable.Type != srv.Type {
		t.Errorf("Incorrect Type after nullable conversion; want: '%s', got: '%s'", srv.Type, nullable.Type)
	}

	if nullable.TypeID == nil {
		t.Error("nullable conversion gave nil TypeID")
	} else if *nullable.TypeID != srv.TypeID {
		t.Errorf("Incorrect TypeID after nullable conversion; want: %d, got: %d", srv.TypeID, *nullable.TypeID)
	}

	if nullable.UpdPending == nil {
		t.Error("nullable conversion gave nil UpdPending")
	} else if *nullable.UpdPending != srv.UpdPending {
		t.Errorf("Incorrect UpdPending after nullable conversion; want: %t, got: %t", srv.UpdPending, *nullable.UpdPending)
	}

	if nullable.XMPPID == nil {
		t.Error("nullable conversion gave nil XMPPID")
	} else if *nullable.XMPPID != srv.XMPPID {
		t.Errorf("Incorrect XMPPID after nullable conversion; want: '%s', got: '%s'", srv.XMPPID, *nullable.XMPPID)
	}

	if nullable.XMPPPasswd == nil {
		t.Error("nullable conversion gave nil XMPPPasswd")
	} else if *nullable.XMPPPasswd != srv.XMPPPasswd {
		t.Errorf("Incorrect XMPPPasswd after nullable conversion; want: '%s', got: '%s'", srv.XMPPPasswd, *nullable.XMPPPasswd)
	}
}

func TestServerNullableV2_Upgrade(t *testing.T) {
	fqdn := "testFQDN"
	srv := Server{
		Cachegroup:       "testCachegroup",
		CachegroupID:     42,
		CDNID:            43,
		CDNName:          "testCDNName",
		DeliveryServices: map[string][]string{"test": []string{"quest"}},
		DomainName:       "testDomainName",
		FQDN:             &fqdn,
		FqdnTime:         time.Now(),
		GUID:             "testGUID",
		HostName:         "testHostName",
		HTTPSPort:        -1,
		ID:               44,
		ILOIPAddress:     "testILOIPAddress",
		ILOIPGateway:     "testILOIPGateway",
		ILOIPNetmask:     "testILOIPNetmask",
		ILOPassword:      "testILOPassword",
		ILOUsername:      "testILOUsername",
		InterfaceMtu:     2,
		InterfaceName:    "testInterfaceName",
		IP6Address:       "::1/64",
		IP6IsService:     true,
		IP6Gateway:       "::2",
		IPAddress:        "0.0.0.1",
		IPIsService:      false,
		IPGateway:        "0.0.0.2",
		IPNetmask:        "255.255.255.0",
		LastUpdated:      TimeNoMod(Time{Time: time.Now().Add(time.Minute), Valid: true}),
		MgmtIPAddress:    "testMgmtIPAddress",
		MgmtIPGateway:    "testMgmtIPGateway",
		MgmtIPNetmask:    "testMgmtIPNetmask",
		OfflineReason:    "testOfflineReason",
		PhysLocation:     "testPhysLocation",
		PhysLocationID:   45,
		Profile:          "testProfile",
		ProfileDesc:      "testProfileDesc",
		ProfileID:        46,
		Rack:             "testRack",
		RevalPending:     true,
		RouterHostName:   "testRouterHostName",
		RouterPortName:   "testRouterPortName",
		Status:           "testStatus",
		StatusID:         47,
		TCPPort:          -3,
		Type:             "testType",
		TypeID:           48,
		UpdPending:       false,
		XMPPID:           "testXMPPID",
		XMPPPasswd:       "testXMPPasswd",
	}

	// this is so much easier than double the lines to manually construct a
	// nullable v2 server
	nullable := srv.ToNullable()

	upgraded, err := nullable.Upgrade()
	if err != nil {
		t.Fatalf("Unexpected error upgrading server: %v", err)
	}

	if nullable.Cachegroup == nil {
		t.Error("Unexpectedly nil Cachegroup in nullable-converted server")
	} else if upgraded.Cachegroup == nil {
		t.Error("upgraded conversion gave nil Cachegroup")
	} else if *upgraded.Cachegroup != *nullable.Cachegroup {
		t.Errorf("Incorrect Cachegroup after upgraded conversion; want: '%s', got: '%s'", *nullable.Cachegroup, *upgraded.Cachegroup)
	}

	if nullable.CachegroupID == nil {
		t.Error("Unexpectedly nil CachegroupID in nullable-converted server")
	} else if upgraded.CachegroupID == nil {
		t.Error("upgraded conversion gave nil CachegroupID")
	} else if *upgraded.CachegroupID != *nullable.CachegroupID {
		t.Errorf("Incorrect CachegroupID after upgraded conversion; want: %d, got: %d", *nullable.CachegroupID, *upgraded.CachegroupID)
	}

	if nullable.CDNID == nil {
		t.Error("Unexpectedly nil CDNID in nullable-converted server")
	} else if upgraded.CDNID == nil {
		t.Error("upgraded conversion gave nil CDNID")
	} else if *upgraded.CDNID != *nullable.CDNID {
		t.Errorf("Incorrect CDNID after upgraded conversion; want: %d, got: %d", *nullable.CDNID, *upgraded.CDNID)
	}

	if nullable.CDNName == nil {
		t.Error("Unexpectedly nil CDNName in nullable-converted server")
	} else if upgraded.CDNName == nil {
		t.Error("upgraded conversion gave nil CDNName")
	} else if *upgraded.CDNName != *nullable.CDNName {
		t.Errorf("Incorrect CDNName after upgraded conversion; want: '%s', got: '%s'", *nullable.CDNName, *upgraded.CDNName)
	}

	if nullable.DeliveryServices == nil {
		t.Error("Unexpectedly nil DeliveryServices in nullable-converted server")
	} else if upgraded.DeliveryServices == nil {
		t.Error("upgraded conversion gave nil DeliveryServices")
	} else if len(*upgraded.DeliveryServices) != len(*nullable.DeliveryServices) {
		t.Errorf("Incorrect number of DeliveryServices after upgraded conversion; want: %d, got: %d", len(*nullable.DeliveryServices), len(*upgraded.DeliveryServices))
	} else {
		for k, v := range *nullable.DeliveryServices {
			upgradedV, ok := (*upgraded.DeliveryServices)[k]
			if !ok {
				t.Errorf("Missing Delivery Service '%s' after upgraded conversion", k)
				continue
			}
			if len(upgradedV) != len(v) {
				t.Errorf("Delivery Service '%s' has incorrect length after upgraded conversion; want: %d, got: %d", k, len(v), len(upgradedV))
			}
			for i, ds := range v {
				upgradedDS := upgradedV[i]
				if upgradedDS != ds {
					t.Errorf("Incorrect value at position %d in Delivery Service '%s' after upgraded conversion; want: '%s', got: '%s'", i, k, ds, upgradedDS)
				}
			}
		}
	}

	if nullable.DomainName == nil {
		t.Error("Unexpectedly nil DomainName in nullable-converted server")
	} else if upgraded.DomainName == nil {
		t.Error("upgraded conversion gave nil DomainName")
	} else if *upgraded.DomainName != *nullable.DomainName {
		t.Errorf("Incorrect DomainName after upgraded conversion; want: '%s', got: '%s'", *nullable.DomainName, *upgraded.DomainName)
	}

	if nullable.FQDN == nil {
		t.Error("Unexpectedly nil FQDN in nullable-converted server")
	} else if upgraded.FQDN == nil {
		t.Error("upgraded conversion gave nil FQDN")
	} else if *upgraded.FQDN != fqdn {
		t.Errorf("Incorrect FQDN after upgraded conversion; want: '%s', got: '%s'", fqdn, *upgraded.FQDN)
	}

	if upgraded.FqdnTime != nullable.FqdnTime {
		t.Errorf("Incorrect FqdnTime after upgraded conversion; want: '%s', got: '%s'", nullable.FqdnTime, upgraded.FqdnTime)
	}

	if nullable.GUID == nil {
		t.Error("Unexpectedly nil GUID in nullable-converted server")
	} else if upgraded.GUID == nil {
		t.Error("upgraded conversion gave nil GUID")
	} else if *upgraded.GUID != *nullable.GUID {
		t.Errorf("Incorrect GUID after upgraded conversion; want: '%s', got: '%s'", *nullable.GUID, *upgraded.GUID)
	}

	if nullable.HostName == nil {
		t.Error("Unexpectedly nil HostName in nullable-converted server")
	} else if upgraded.HostName == nil {
		t.Error("upgraded conversion gave nil HostName")
	} else if *upgraded.HostName != *nullable.HostName {
		t.Errorf("Incorrect HostName after upgraded conversion; want: '%s', got: '%s'", *nullable.HostName, *upgraded.HostName)
	}

	if nullable.HTTPSPort == nil {
		t.Error("Unexpectedly nil HTTPSPort in nullable-converted server")
	} else if upgraded.HTTPSPort == nil {
		t.Error("upgraded conversion gave nil HTTPSPort")
	} else if *upgraded.HTTPSPort != *nullable.HTTPSPort {
		t.Errorf("Incorrect HTTPSPort after upgraded conversion; want: %d, got: %d", *nullable.HTTPSPort, *upgraded.HTTPSPort)
	}

	if nullable.ID == nil {
		t.Error("Unexpectedly nil ID in nullable-converted server")
	} else if upgraded.ID == nil {
		t.Error("upgraded conversion gave nil ID")
	} else if *upgraded.ID != *nullable.ID {
		t.Errorf("Incorrect ID after upgraded conversion; want: %d, got: %d", *nullable.ID, *upgraded.ID)
	}

	if nullable.ILOIPAddress == nil {
		t.Error("Unexpectedly nil ILOIPAddress in nullable-converted server")
	} else if upgraded.ILOIPAddress == nil {
		t.Error("upgraded conversion gave nil ILOIPAddress")
	} else if *upgraded.ILOIPAddress != *nullable.ILOIPAddress {
		t.Errorf("Incorrect ILOIPAddress after upgraded conversion; want: '%s', got: '%s'", *nullable.ILOIPAddress, *upgraded.ILOIPAddress)
	}

	if nullable.ILOIPGateway == nil {
		t.Error("Unexpectedly nil ILOIPGateway in nullable-converted server")
	} else if upgraded.ILOIPGateway == nil {
		t.Error("upgraded conversion gave nil ILOIPGateway")
	} else if *upgraded.ILOIPGateway != *nullable.ILOIPGateway {
		t.Errorf("Incorrect ILOIPGateway after upgraded conversion; want: '%s', got: '%s'", *nullable.ILOIPGateway, *upgraded.ILOIPGateway)
	}

	if nullable.ILOIPNetmask == nil {
		t.Error("Unexpectedly nil ILOIPNetmask in nullable-converted server")
	} else if upgraded.ILOIPNetmask == nil {
		t.Error("upgraded conversion gave nil ILOIPNetmask")
	} else if *upgraded.ILOIPNetmask != *nullable.ILOIPNetmask {
		t.Errorf("Incorrect ILOIPNetmask after upgraded conversion; want: '%s', got: '%s'", *nullable.ILOIPNetmask, *upgraded.ILOIPNetmask)
	}

	if nullable.ILOPassword == nil {
		t.Error("Unexpectedly nil ILOPassword in nullable-converted server")
	} else if upgraded.ILOPassword == nil {
		t.Error("upgraded conversion gave nil ILOPassword")
	} else if *upgraded.ILOPassword != *nullable.ILOPassword {
		t.Errorf("Incorrect ILOPassword after upgraded conversion; want: '%s', got: '%s'", *nullable.ILOPassword, *upgraded.ILOPassword)
	}

	if nullable.ILOUsername == nil {
		t.Error("Unexpectedly nil ILOUsername in nullable-converted server")
	} else if upgraded.ILOUsername == nil {
		t.Error("upgraded conversion gave nil ILOUsername")
	} else if *upgraded.ILOUsername != *nullable.ILOUsername {
		t.Errorf("Incorrect ILOUsername after upgraded conversion; want: '%s', got: '%s'", *nullable.ILOUsername, *upgraded.ILOUsername)
	}

	checkInterfaces := true
	if nullable.InterfaceMtu == nil {
		t.Error("Unexpectedly nil InterfaceMtu in nullable-converted server")
		checkInterfaces = false
	}
	if nullable.InterfaceName == nil {
		t.Error("Unexpectedly nil InterfaceName in nullable-converted server")
		checkInterfaces = false
	}
	if nullable.IP6Address == nil {
		t.Error("Unexpectedly nil IP6Address in nullable-converted server")
		checkInterfaces = false
	}
	if nullable.IP6IsService == nil {
		t.Error("Unexpectedly nil IP6IsService in nullable-converted server")
		checkInterfaces = false
	}
	if nullable.IP6Gateway == nil {
		t.Error("Unexpectedly nil IP6Gateway in nullable-converted server")
		checkInterfaces = false
	}
	if nullable.IPAddress == nil {
		t.Error("Unexpectedly nil IPAddress in nullable-converted server")
		checkInterfaces = false
	}
	if nullable.IPIsService == nil {
		t.Error("Unexpectedly nil IPIsService in nullable-converted server")
		checkInterfaces = false
	}
	if nullable.IPGateway == nil {
		t.Error("Unexpectedly nil IPGateway in nullable-converted server")
		checkInterfaces = false
	}
	if nullable.IPNetmask == nil {
		t.Error("Unexpectedly nil IPNetmask in nullable-converted server")
		checkInterfaces = false
	}

	if checkInterfaces {
		infLen := len(upgraded.Interfaces)
		if infLen < 1 {
			t.Error("Expected exactly one interface after upgrade, got: 0")
		} else {
			if infLen > 1 {
				t.Errorf("Expected exactly one interface after upgrade, got: %d", infLen)
			}

			inf := upgraded.Interfaces[0]
			if inf.Name != *nullable.InterfaceName {
				t.Errorf("Incorrect interface name after upgrade; want: '%s', got: '%s'", *nullable.InterfaceName, inf.Name)
			}

			if inf.MTU == nil {
				t.Error("Unexpectedly nil Interface MTU after upgrade")
			} else if *inf.MTU != uint64(*nullable.InterfaceMtu) {
				t.Errorf("Incorrect Interface MTU after upgrade; want: %d, got: %d", *nullable.InterfaceMtu, *inf.MTU)
			}

			if inf.Monitor {
				t.Error("Incorrect Interface Monitor after upgrade; want: false, got: true")
			}

			if inf.MaxBandwidth != nil {
				t.Error("Unexpectedly non-nil Interface MaxBandwidth after upgrade")
			}

			if len(inf.IPAddresses) != 2 {
				t.Errorf("Incorrect number of IP Addresses after upgrade; want: 2, got: %d", len(inf.IPAddresses))
			} else {
				ip := inf.IPAddresses[0]
				cidrIndex := strings.Index(ip.Address, "/")
				addr := ip.Address
				if cidrIndex >= 0 {
					addr = addr[:cidrIndex]
				}

				// TODO: calculate and verify netmask
				if addr == *nullable.IPAddress {
					if ip.Gateway == nil {
						t.Error("Unexpectedly nil IPv4 Gateway after upgrade")
					} else if *ip.Gateway != *nullable.IPGateway {
						t.Errorf("Incorrect IPv4 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IPGateway, *ip.Gateway)
					}

					if ip.ServiceAddress != *nullable.IPIsService {
						t.Errorf("Incorrect IPv4 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IPIsService, ip.ServiceAddress)
					}

					secondIP := inf.IPAddresses[1]
					if secondIP.Address != *nullable.IP6Address {
						t.Errorf("Incorrect IPv6 Address after upgrade; want: '%s', got: '%s'", *nullable.IP6Address, secondIP.Address)
					} else {
						if secondIP.Gateway == nil {
							t.Error("Unexpectedly nil IPv6 Gateway after upgrade")
						} else if *secondIP.Gateway != *nullable.IP6Gateway {
							t.Errorf("Incorrect IPv6 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IP6Gateway, *secondIP.Gateway)
						}

						if secondIP.ServiceAddress != *nullable.IP6IsService {
							t.Errorf("Incorrect IPv6 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IP6IsService, secondIP.ServiceAddress)
						}
					}
				} else if ip.Address == *nullable.IP6Address {
					if ip.Gateway == nil {
						t.Error("Unexpectedly nil IPv6 Gateway after upgrade")
					} else if *ip.Gateway != *nullable.IP6Gateway {
						t.Errorf("Incorrect IPv6 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IP6Gateway, *ip.Gateway)
					}

					if ip.ServiceAddress != *nullable.IP6IsService {
						t.Errorf("Incorrect IPv6 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IP6IsService, ip.ServiceAddress)
					}

					secondIP := inf.IPAddresses[1]
					cidrIndex = strings.Index(secondIP.Address, "/")
					addr = secondIP.Address
					if cidrIndex >= 0 {
						addr = addr[:cidrIndex]
					}
					// TODO: calculate and verify netmask
					if addr != *nullable.IPAddress {
						t.Errorf("Incorrect IPv4 Address after upgrade; want: '%s', got: '%s'", *nullable.IPAddress, secondIP.Address)
					} else {
						if secondIP.Gateway == nil {
							t.Error("Unexpectedly nil IPv4 Gateway after upgrade")
						} else if *secondIP.Gateway != *nullable.IPGateway {
							t.Errorf("Incorrect IPv4 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IPGateway, *secondIP.Gateway)
						}

						if secondIP.ServiceAddress != *nullable.IPIsService {
							t.Errorf("Incorrect IPv4 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IPIsService, secondIP.ServiceAddress)
						}
					}

				} else {
					t.Errorf("Unknown IP address '%s' found in interface after upgrade", ip.Address)
					ip = inf.IPAddresses[1]
					cidrIndex = strings.Index(ip.Address, "/")
					addr = ip.Address
					if cidrIndex >= 0 {
						addr = addr[:cidrIndex]
					}

					if addr == *nullable.IPAddress {
						t.Error("Missing IPv6 address after upgrade")
						if ip.Gateway == nil {
							t.Error("Unexpectedly nil IPv4 Gateway after upgrade")
						} else if *ip.Gateway != *nullable.IPGateway {
							t.Errorf("Incorrect IPv4 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IPGateway, *ip.Gateway)
						}

						if ip.ServiceAddress != *nullable.IPIsService {
							t.Errorf("Incorrect IPv4 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IPIsService, ip.ServiceAddress)
						}
					} else if ip.Address == *nullable.IP6Address {
						t.Error("Missing IPv4 address after upgrade")
						if ip.Gateway == nil {
							t.Error("Unexpectedly nil IPv6 Gateway after upgrade")
						} else if *ip.Gateway != *nullable.IP6Gateway {
							t.Errorf("Incorrect IPv6 Gateway after upgrade; want: '%s', got: '%s'", *nullable.IP6Gateway, *ip.Gateway)
						}

						if ip.ServiceAddress != *nullable.IP6IsService {
							t.Errorf("Incorrect IPv6 ServiceAddress value after upgrade; want: %t, got: %t", *nullable.IP6IsService, ip.ServiceAddress)
						}
					} else {
						t.Errorf("Unknown IP address '%s' found in interface after upgrade", ip.Address)
						t.Error("Missing both IPv4 and IPv6 address after upgrade")
					}
				}
			}
		}
	}

	if nullable.LastUpdated == nil {
		t.Error("Unexpectedly nil LastUpdated in nullable-converted server")
	} else if upgraded.LastUpdated == nil {
		t.Error("upgraded conversion gave nil LastUpdated")
	} else if *upgraded.LastUpdated != *nullable.LastUpdated {
		t.Errorf("Incorrect LastUpdated after upgraded conversion; want: '%s', got: '%s'", *nullable.LastUpdated, *upgraded.LastUpdated)
	}

	if nullable.MgmtIPAddress == nil {
		t.Error("Unexpectedly nil MgmtIPAddress in nullable-converted server")
	} else if upgraded.MgmtIPAddress == nil {
		t.Error("upgraded conversion gave nil MgmtIPAddress")
	} else if *upgraded.MgmtIPAddress != *nullable.MgmtIPAddress {
		t.Errorf("Incorrect MgmtIPAddress after upgraded conversion; want: '%s', got: '%s'", *nullable.MgmtIPAddress, *upgraded.MgmtIPAddress)
	}

	if nullable.MgmtIPGateway == nil {
		t.Error("Unexpectedly nil MgmtIPGateway in nullable-converted server")
	} else if upgraded.MgmtIPGateway == nil {
		t.Error("upgraded conversion gave nil MgmtIPGateway")
	} else if *upgraded.MgmtIPGateway != *nullable.MgmtIPGateway {
		t.Errorf("Incorrect MgmtIPGateway after upgraded conversion; want: '%s', got: '%s'", *nullable.MgmtIPGateway, *upgraded.MgmtIPGateway)
	}

	if nullable.MgmtIPNetmask == nil {
		t.Error("Unexpectedly nil MgmtIPNetmask in nullable-converted server")
	} else if upgraded.MgmtIPNetmask == nil {
		t.Error("upgraded conversion gave nil MgmtIPNetmask")
	} else if *upgraded.MgmtIPNetmask != *nullable.MgmtIPNetmask {
		t.Errorf("Incorrect MgmtIPNetmask after upgraded conversion; want: '%s', got: '%s'", *nullable.MgmtIPNetmask, *upgraded.MgmtIPNetmask)
	}

	if nullable.OfflineReason == nil {
		t.Error("Unexpectedly nil OfflineReason in nullable-converted server")
	} else if upgraded.OfflineReason == nil {
		t.Error("upgraded conversion gave nil OfflineReason")
	} else if *upgraded.OfflineReason != *nullable.OfflineReason {
		t.Errorf("Incorrect OfflineReason after upgraded conversion; want: '%s', got: '%s'", *nullable.OfflineReason, *upgraded.OfflineReason)
	}

	if nullable.PhysLocation == nil {
		t.Error("Unexpectedly nil PhysLocation in nullable-converted server")
	} else if upgraded.PhysLocation == nil {
		t.Error("upgraded conversion gave nil PhysLocation")
	} else if *upgraded.PhysLocation != *nullable.PhysLocation {
		t.Errorf("Incorrect PhysLocation after upgraded conversion; want: '%s', got: '%s'", *nullable.PhysLocation, *upgraded.PhysLocation)
	}

	if nullable.PhysLocationID == nil {
		t.Error("Unexpectedly nil PhysLocationID in nullable-converted server")
	} else if upgraded.PhysLocationID == nil {
		t.Error("upgraded conversion gave nil PhysLocationID")
	} else if *upgraded.PhysLocationID != *nullable.PhysLocationID {
		t.Errorf("Incorrect PhysLocationID after upgraded conversion; want: %d, got: %d", *nullable.PhysLocationID, *upgraded.PhysLocationID)
	}

	if nullable.Profile == nil {
		t.Error("Unexpectedly nil Profile in nullable-converted server")
	} else if upgraded.Profile == nil {
		t.Error("upgraded conversion gave nil Profile")
	} else if *upgraded.Profile != *nullable.Profile {
		t.Errorf("Incorrect Profile after upgraded conversion; want: '%s', got: '%s'", *nullable.Profile, *upgraded.Profile)
	}

	if nullable.ProfileDesc == nil {
		t.Error("Unexpectedly nil ProfileDesc in nullable-converted server")
	} else if upgraded.ProfileDesc == nil {
		t.Error("upgraded conversion gave nil ProfileDesc")
	} else if *upgraded.ProfileDesc != *nullable.ProfileDesc {
		t.Errorf("Incorrect ProfileDesc after upgraded conversion; want: '%s', got: '%s'", *nullable.ProfileDesc, *upgraded.ProfileDesc)
	}

	if nullable.ProfileID == nil {
		t.Error("Unexpectedly nil ProfileID in nullable-converted server")
	} else if upgraded.ProfileID == nil {
		t.Error("upgraded conversion gave nil ProfileID")
	} else if *upgraded.ProfileID != *nullable.ProfileID {
		t.Errorf("Incorrect ProfileID after upgraded conversion; want: %d, got: %d", *nullable.ProfileID, *upgraded.ProfileID)
	}

	if nullable.Rack == nil {
		t.Error("Unexpectedly nil Rack in nullable-converted server")
	} else if upgraded.Rack == nil {
		t.Error("upgraded conversion gave nil Rack")
	} else if *upgraded.Rack != *nullable.Rack {
		t.Errorf("Incorrect Rack after upgraded conversion; want: '%s', got: '%s'", *nullable.Rack, *upgraded.Rack)
	}

	if nullable.RevalPending == nil {
		t.Error("Unexpectedly nil RevalPending in nullable-converted server")
	} else if upgraded.RevalPending == nil {
		t.Error("upgraded conversion gave nil RevalPending")
	} else if *upgraded.RevalPending != *nullable.RevalPending {
		t.Errorf("Incorrect RevalPending after upgraded conversion; want: %t, got: %t", *nullable.RevalPending, *upgraded.RevalPending)
	}

	if nullable.RouterHostName == nil {
		t.Error("Unexpectedly nil RouterHostName in nullable-converted server")
	} else if upgraded.RouterHostName == nil {
		t.Error("upgraded conversion gave nil RouterHostName")
	} else if *upgraded.RouterHostName != *nullable.RouterHostName {
		t.Errorf("Incorrect RouterHostName after upgraded conversion; want: '%s', got: '%s'", *nullable.RouterHostName, *upgraded.RouterHostName)
	}

	if nullable.RouterPortName == nil {
		t.Error("Unexpectedly nil RouterPortName in nullable-converted server")
	} else if upgraded.RouterPortName == nil {
		t.Error("upgraded conversion gave nil RouterPortName")
	} else if *upgraded.RouterPortName != *nullable.RouterPortName {
		t.Errorf("Incorrect RouterPortName after upgraded conversion; want: '%s', got: '%s'", *nullable.RouterPortName, *upgraded.RouterPortName)
	}

	if nullable.Status == nil {
		t.Error("Unexpectedly nil Status in nullable-converted server")
	} else if upgraded.Status == nil {
		t.Error("upgraded conversion gave nil Status")
	} else if *upgraded.Status != *nullable.Status {
		t.Errorf("Incorrect Status after upgraded conversion; want: '%s', got: '%s'", *nullable.Status, *upgraded.Status)
	}

	if nullable.StatusID == nil {
		t.Error("Unexpectedly nil StatusID in nullable-converted server")
	} else if upgraded.StatusID == nil {
		t.Error("upgraded conversion gave nil StatusID")
	} else if *upgraded.StatusID != *nullable.StatusID {
		t.Errorf("Incorrect StatusID after upgraded conversion; want: %d, got: %d", *nullable.StatusID, *upgraded.StatusID)
	}

	if nullable.TCPPort == nil {
		t.Error("Unexpectedly nil TCPPort in nullable-converted server")
	} else if upgraded.TCPPort == nil {
		t.Error("upgraded conversion gave nil TCPPort")
	} else if *upgraded.TCPPort != *nullable.TCPPort {
		t.Errorf("Incorrect TCPPort after upgraded conversion; want: %d, got: %d", *nullable.TCPPort, *upgraded.TCPPort)
	}

	if upgraded.Type != nullable.Type {
		t.Errorf("Incorrect Type after upgraded conversion; want: '%s', got: '%s'", nullable.Type, upgraded.Type)
	}

	if nullable.TypeID == nil {
		t.Error("Unexpectedly nil TypeID in nullable-converted server")
	} else if upgraded.TypeID == nil {
		t.Error("upgraded conversion gave nil TypeID")
	} else if *upgraded.TypeID != *nullable.TypeID {
		t.Errorf("Incorrect TypeID after upgraded conversion; want: %d, got: %d", *nullable.TypeID, *upgraded.TypeID)
	}

	if nullable.UpdPending == nil {
		t.Error("Unexpectedly nil UpdPending in nullable-converted server")
	} else if upgraded.UpdPending == nil {
		t.Error("upgraded conversion gave nil UpdPending")
	} else if *upgraded.UpdPending != *nullable.UpdPending {
		t.Errorf("Incorrect UpdPending after upgraded conversion; want: %t, got: %t", *nullable.UpdPending, *upgraded.UpdPending)
	}

	if nullable.XMPPID == nil {
		t.Error("Unexpectedly nil XMPPID in nullable-converted server")
	} else if upgraded.XMPPID == nil {
		t.Error("upgraded conversion gave nil XMPPID")
	} else if *upgraded.XMPPID != *nullable.XMPPID {
		t.Errorf("Incorrect XMPPID after upgraded conversion; want: '%s', got: '%s'", *nullable.XMPPID, *upgraded.XMPPID)
	}

	if nullable.XMPPPasswd == nil {
		t.Error("Unexpectedly nil XMPPPasswd in nullable-converted server")
	} else if upgraded.XMPPPasswd == nil {
		t.Error("upgraded conversion gave nil XMPPPasswd")
	} else if *upgraded.XMPPPasswd != *nullable.XMPPPasswd {
		t.Errorf("Incorrect XMPPPasswd after upgraded conversion; want: '%s', got: '%s'", *nullable.XMPPPasswd, *upgraded.XMPPPasswd)
	}
}
