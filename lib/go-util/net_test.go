package util

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
   http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Package rfc contains functions implementing RFC 7234, 2616, and other RFCs.
// When changing functions, be sure they still conform to the corresponding RFC.
// When adding symbols, document the RFC and section they correspond to.

import (
	"net"
	"testing"
)

func TestCoalesceIPs(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.1.2"),
		net.ParseIP("192.168.1.3"),
		net.ParseIP("192.168.2.1"),
		net.ParseIP("192.168.2.2"),
		net.ParseIP("192.168.2.3"),
		net.ParseIP("192.168.2.4"),
	}

	nets := CoalesceIPs(ips, 2, 24)

	for _, ipnet := range nets {
		if ipnet.String() != "192.168.1.0/24" && ipnet.String() != "192.168.2.0/24" {
			t.Errorf("expected '192.168.1.0/24' and '192.168.2.0/24', actual: %+v\n", ipnet)
		}
	}
}

func TestCoalesceIPsSmallerThanNum(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.1.2"),
		net.ParseIP("192.168.1.3"),
		net.ParseIP("192.168.2.1"),
		net.ParseIP("192.168.2.2"),
		net.ParseIP("192.168.2.3"),
		net.ParseIP("192.168.2.4"),
	}

	nets := CoalesceIPs(ips, 4, 24)

	expecteds := map[string]struct{}{
		"192.168.1.1/32": {},
		"192.168.1.2/32": {},
		"192.168.1.3/32": {},
		"192.168.2.0/24": {},
	}

	for _, ipnet := range nets {
		if _, ok := expecteds[ipnet.String()]; !ok {
			t.Errorf("expected: %+v actual: %+v\n", expecteds, ipnet)
		}
		delete(expecteds, ipnet.String())
	}
}

func TestCoalesceIPsV6(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("2001:db8::1"),
		net.ParseIP("2001:db8::2"),
		net.ParseIP("2001:db8::3"),
		net.ParseIP("2001:db8::4:1"),
		net.ParseIP("2001:db8::4:2"),
		net.ParseIP("2001:db8::4:3"),
	}

	nets := CoalesceIPs(ips, 3, 112)

	expecteds := map[string]struct{}{
		"2001:db8::/112":    {},
		"2001:db8::4:0/112": {},
	}

	for _, ipnet := range nets {
		if _, ok := expecteds[ipnet.String()]; !ok {
			t.Errorf("expected: %+v actual: %+v\n", expecteds, ipnet)
		}
		delete(expecteds, ipnet.String())
	}
}

func TestCoalesceIPsV6SmallerThanNum(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("2001:db8::1"),
		net.ParseIP("2001:db8::2"),
		net.ParseIP("2001:db8::3"),
		net.ParseIP("2001:db8::4:1"),
		net.ParseIP("2001:db8::4:2"),
		net.ParseIP("2001:db8::4:3"),
		net.ParseIP("2001:db8::4:4"),
	}

	nets := CoalesceIPs(ips, 4, 112)

	expecteds := map[string]struct{}{
		"2001:db8::1/128":   {},
		"2001:db8::2/128":   {},
		"2001:db8::3/128":   {},
		"2001:db8::4:0/112": {},
	}

	for _, ipnet := range nets {
		if _, ok := expecteds[ipnet.String()]; !ok {
			t.Errorf("expected: %+v actual: %+v\n", expecteds, ipnet)
		}
		delete(expecteds, ipnet.String())
	}
}

func TestRangeStr(t *testing.T) {
	inputExpecteds := map[string]string{
		"192.168.1.0/24":     "192.168.1.0-192.168.1.255",
		"192.168.1.0/16":     "192.168.0.0-192.168.255.255",
		"192.168.1.42/32":    "192.168.1.42",
		"2001:db8::4:42/128": "2001:db8::4:42",
		"2001:db8::4:0/112":  "2001:db8::4:0-2001:db8::4:ffff",
	}
	for input, expected := range inputExpecteds {
		_, ipn, err := net.ParseCIDR(input)
		if err != nil {
			t.Fatal(err.Error())
		}

		//	t.Errorf("ipn: " + ipn.String())

		actual := RangeStr(ipn)
		if expected != actual {
			t.Errorf("expected: '" + expected + "' actual '" + actual + "'")
		}
	}
}

func TestFirstIP(t *testing.T) {
	inputExpecteds := map[string]string{
		"192.168.1.0/24":     "192.168.1.0",
		"192.168.1.0/16":     "192.168.0.0",
		"192.168.1.42/32":    "192.168.1.42",
		"2001:db8::4:42/128": "2001:db8::4:42",
		"2001:db8::4:0/112":  "2001:db8::4:0",
	}
	for input, expected := range inputExpecteds {
		_, ipn, err := net.ParseCIDR(input)
		if err != nil {
			t.Fatal(err.Error())
		}

		//	t.Errorf("ipn: " + ipn.String())

		actual := FirstIP(ipn).String()
		if expected != actual {
			t.Errorf("expected: '" + expected + "' actual '" + actual + "'")
		}
	}
}

func TestLastIP(t *testing.T) {
	inputExpecteds := map[string]string{
		"192.168.1.0/24":     "192.168.1.255",
		"192.168.1.0/16":     "192.168.255.255",
		"192.168.1.42/32":    "192.168.1.42",
		"2001:db8::4:42/128": "2001:db8::4:42",
		"2001:db8::4:0/112":  "2001:db8::4:ffff",
	}
	for input, expected := range inputExpecteds {
		_, ipn, err := net.ParseCIDR(input)
		if err != nil {
			t.Fatal(err.Error())
		}

		//	t.Errorf("ipn: " + ipn.String())

		actual := LastIP(ipn).String()
		if expected != actual {
			t.Errorf("expected: '" + expected + "' actual '" + actual + "'")
		}
	}
}
