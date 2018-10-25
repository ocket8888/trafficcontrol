package main

//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/net/publicsuffix"
)

const __version__ = "1.0.0"

// Environment variables used:
//   TO_URL      -- URL for reference Traffic Ops
//   TEST_URL    -- URL for test Traffic Ops
//   TO_USER     -- Username for both instances
//   TO_PASSWORD -- Password for both instances
type Creds struct {
	// common user/password
	User     string `json:"u" required:"true"`
	Password string `json:"p" required:"true"`
}

type Connect struct {
	// URL of reference traffic_ops
	URL         string       `required:"true"`
	Client      *http.Client `ignore:"true"`
	ResultsPath string       `ignore:"true"`
	creds       Creds        `ignore:"true"`
}

func (to *Connect) login(creds Creds) error {
	body, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	to.Client = &http.Client{Transport: tr}
	url := to.URL + `/api/1.3/user/login`
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Create cookiejar so created cookie will be reused
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return err
	}
	to.Client.Jar = jar

	resp, err := to.Client.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Logged in to %s: %s", to.URL, string(data))
	return nil
}

func testRoute(tos []*Connect, route string) {
	// keeps result along with instance -- no guarantee on order collected
	type result struct {
		TO  *Connect
		Res string
	}
	var res []result
	ch := make(chan result, len(tos))

	var wg sync.WaitGroup
	var m sync.Mutex

	for _, to := range tos {
		wg.Add(1)
		go func(to *Connect) {
			s, err := to.get(route)
			if err != nil {
				s = err.Error()
			}
			ch <- result{to, s}
			wg.Done()
		}(to)

		wg.Add(1)
		go func() {
			m.Lock()
			defer m.Unlock()
			res = append(res, <-ch)
			wg.Done()
		}()
	}
	wg.Wait()
	close(ch)

	if res[0].Res == res[1].Res {
		log.Printf("Identical results (%d bytes) from %s", len(res[0].Res), route)
	} else {
		log.Print("Diffs from ", route, " written to")
		for _, r := range res {
			p, err := r.TO.writeResults(route, r.Res)
			if err != nil {
				log.Fatal("Error writing results for ", route)
			}
			log.Print("  ", p)
		}
	}
}

func (to *Connect) writeResults(route string, res string) (string, error) {
	var dst bytes.Buffer
	if err := json.Indent(&dst, []byte(res), "", "  "); err != nil {
		dst.WriteString(res)
	}

	m := func(r rune) rune {
		if unicode.IsPunct(r) && r != '.' || unicode.IsSymbol(r) {
			return '-'
		}
		return r
	}

	err := os.MkdirAll(to.ResultsPath, 0755)
	if err != nil {
		return "", err
	}

	p := to.ResultsPath + "/" + strings.Map(m, route)
	err = ioutil.WriteFile(p, dst.Bytes(), 0644)
	return p, err
}

func (to *Connect) get(route string) (string, error) {
	var url string
	if route[0] != '/' {
		url = to.URL + `/` + route
	} else {
		url = to.URL + route
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := to.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	return string(data), err
}

func (to *Connect) getCDNNames() ([]string, error) {
	res, err := to.get(`api/1.3/cdns`)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)

	var cdnResp tc.CDNsResponse

	err = json.Unmarshal([]byte(res), &cdnResp)
	if err != nil {
		return nil, err
	}
	var cdnNames []string
	for _, c := range cdnResp.Response {
		cdnNames = append(cdnNames, c.Name)
	}
	return cdnNames, nil
}

func main() {

	routesFileLong := flag.String("file", "", "File listing routes to test (will read from stdin if not given)")
	routesFileShort := flag.String("f", "", "File listing routes to test (will read from stdin if not given)")
	resultsPathLong := flag.String("results_path", "", "Directory where results will be written")
	resultsPathShort := flag.String("r", "", "Directory where results will be written")
	doSnapshotLong := flag.Bool("snapshot", false, "Perform comparison of all CDN's snapshotted CRConfigs")
	doSnapshotShort := flag.Bool("s", false, "Perform comparison of all CDN's snapshotted CRConfigs")
	refURL := flag.String("ref_url", "", "The URL for the reference Traffic Ops instance (overrides TO_URL environment variable)")
	testURL := flag.String("test_url", "", "The URL for the testing Traffic Ops instance (overrides TEST_URL environment variable)")
	refUser := flag.String("ref_user", "", "The username for logging into the reference Traffic Ops instance (overrides TO_USER environment variable)")
	refPasswd := flag.String("ref_passwd", "", "The password for logging into the reference Traffic Ops instance (overrides TO_PASSWORD environment variable)")
	testUser := flag.String("test_user", "", "The username for logging into the testing Traffic Ops instance (overrides TEST_USER environment variable)")
	testPasswd := flag.String("test_passwd", "", "The password for logging into the testing Traffic Ops instance (overrides TEST_PASSWORD environment variable)")
	versionLong := flag.Bool("version", false, "Print version information and exit")
	versionShort := flag.Bool("v", false, "Print version information and exit")
	helpLong := flag.Bool("help", false, "Print usage information and exit")
	helpShort := flag.Bool("h", false, "Print usage information and exit")
	flag.Parse()

	// Coalesce long/short form options
	doSnapshot := *doSnapshotShort || *doSnapshotLong

	version := *versionLong || *versionShort
	if version {
		fmt.Printf("Traffic Control 'compare' tool v%s\n", __version__)
		os.Exit(0)
	}

	help := *helpLong || *helpShort
	if help {
		flag.Usage()
		os.Exit(0)
	}

	var resultsPath string
	if *resultsPathLong == "" {
		if *resultsPathShort == "" {
			resultsPath = "results"
		} else {
			resultsPath = *resultsPathShort
		}
	} else if *resultsPathShort == "" || *resultsPathShort == *resultsPathLong {
		resultsPath = *resultsPathLong
	} else {
		log.Fatal("Duplicate specification of results path! (Hint: try '-h'/'--help')")
	}

	var routesFile *os.File
	var err error
	if *routesFileLong == "" {
		if *routesFileShort == "" {
			routesFile = os.Stdin
		} else {
			if routesFile, err = os.Open(*routesFileShort); err != nil {
				log.Fatal(err)
			}
			defer routesFile.Close()
		}
	} else if *routesFileShort == "" || *routesFileLong == *routesFileShort {
		if routesFile, err = os.Open(*routesFileLong); err != nil {
			log.Fatal(err)
		}
		defer routesFile.Close()
	} else {
		log.Fatal("Duplicate specification of input file! (Hint: try '-h'/'--help')")
	}

	// refTO, testTO are connections to the two Traffic Ops instances
	var refTO = &Connect{ResultsPath: resultsPath + `/ref`}
	var testTO = &Connect{ResultsPath: resultsPath + `/test`}

	if *refUser != "" && *refPasswd != "" {
		refTO.creds = Creds{*refUser, *refPasswd}
	} else {
		err := envconfig.Process("TO", &refTO.creds)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if *testUser != "" && *testPasswd != "" {
		testTO.creds = Creds{*testUser, *testPasswd}
	} else {
		err := envconfig.Process("TEST", &testTO.creds)
		if err != nil {
			// if not provided, re-use the same credentials
			testTO.creds = refTO.creds
		}
	}

	if *refURL != "" {
		refTO.URL = *refURL
	} else {
		err := envconfig.Process("TO", refTO)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if *testURL != "" {
		testTO.URL = *testURL
	} else {
		err := envconfig.Process("TEST", testTO)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	tos := []*Connect{refTO, testTO}

	// Login to the 2 Traffic Ops instances concurrently
	var wg sync.WaitGroup
	wg.Add(len(tos))
	for _, t := range tos {
		go func(to *Connect) {
			log.Print("Login to ", to.URL)
			err := to.login(to.creds)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(t)
	}
	wg.Wait()

	var testRoutes []string

	scanner := bufio.NewScanner(routesFile)
	for scanner.Scan() {
		testRoutes = append(testRoutes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Add(len(testRoutes))
	for _, route := range testRoutes {
		go func(r string) {
			testRoute(tos, r)
			wg.Done()
		}(route)
	}
	wg.Wait()

	if doSnapshot {
		cdnNames, err := refTO.getCDNNames()
		if err != nil {
			panic(err)
		}
		log.Printf("CDNNames are %+v", cdnNames)

		wg.Add(2 * len(cdnNames))
		for _, cdnName := range cdnNames {
			log.Print("CDN ", cdnName)
			go func(c string) {
				testRoute(tos, `api/1.3/cdns/`+c+`/snapshot`)
				wg.Done()
			}(cdnName)
			go func(c string) {
				testRoute(tos, `api/1.3/cdns/`+c+`/snapshot/new`)
				wg.Done()
			}(cdnName)
		}
		wg.Wait()
	}
}
