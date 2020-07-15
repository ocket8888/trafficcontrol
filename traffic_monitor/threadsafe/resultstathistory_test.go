package threadsafe

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
	"math/rand"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/traffic_monitor/srvhttp"
	jsoniter "github.com/json-iterator/go"
)

func randResultStatHistory() ResultStatHistory {
	hist := NewResultStatHistory()

	num := 5
	for i := 0; i < num; i++ {
		hist.Store(tc.CacheName(randStr()), randResultStatValHistory())
	}
	return hist
}

func randResultStatValHistory() ResultStatValHistory {
	a := NewResultStatValHistory()
	num := 5
	numSlice := 5
	for i := 0; i < num; i++ {
		cacheName := randStr()
		vals := []cache.ResultStatVal{}
		for j := 0; j < numSlice; j++ {
			vals = append(vals, randResultStatVal())
		}
		a.Store(cacheName, vals)
	}
	return a
}

func randResultStatVal() cache.ResultStatVal {
	return cache.ResultStatVal{
		Val:  uint64(rand.Int63()),
		Time: time.Now(),
		Span: uint64(rand.Int63()),
	}
}

func randResultInfoHistory() cache.ResultInfoHistory {
	hist := cache.ResultInfoHistory{}

	num := 5
	infNum := 5
	for i := 0; i < num; i++ {
		cacheName := tc.CacheName(randStr())
		for j := 0; j < infNum; j++ {
			hist[cacheName] = append(hist[cacheName], randResultInfo())
		}
	}
	return hist
}

func randResultInfo() cache.ResultInfo {
	return cache.ResultInfo{
		ID:          randStr(),
		Error:       fmt.Errorf(randStr()),
		Time:        time.Now(),
		RequestTime: time.Millisecond * time.Duration(rand.Int()),
		Vitals:      randVitals(),
		PollID:      uint64(rand.Int63()),
		Available:   randBool(),
	}
}

func randVitals() cache.Vitals {
	return cache.Vitals{
		LoadAvg:    rand.Float64(),
		BytesOut:   rand.Uint64(),
		BytesIn:    rand.Uint64(),
		KbpsOut:    rand.Int63(),
		MaxKbpsOut: rand.Int63(),
	}
}

func randStr() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
	num := 100
	s := ""
	for i := 0; i < num; i++ {
		s += string(chars[rand.Intn(len(chars))])
	}
	return s
}

func randBool() bool {
	return rand.Int()%2 == 0
}

type DummyFilterNever struct {
}

func (f DummyFilterNever) UseStat(name string) bool {
	return false
}

func (f DummyFilterNever) UseCache(name tc.CacheName) bool {
	return false
}

func (f DummyFilterNever) WithinStatHistoryMax(i int) bool {
	return false
}

func TestStatsMarshall(t *testing.T) {
	statHist := randResultStatHistory()
	infHist := randResultInfoHistory()
	filter := DummyFilterNever{}
	params := url.Values{}
	beforeStatsMarshall := time.Now()
	bytes, err := StatsMarshall(statHist, infHist, tc.CRStates{}, tc.LegacyTrafficMonitorConfigMap{}, cache.Kbpses{}, filter, params)
	afterStatsMarshall := time.Now()
	if err != nil {
		t.Fatalf("StatsMarshall return expected nil err, actual err: %v", err)
	}
	// if len(bytes) > 0 {
	// 	t.Errorf("expected empty bytes, actual: %v", string(bytes))
	// }

	stats := cache.Stats{}
	json := jsoniter.ConfigFastest // TODO make configurable
	if err := json.Unmarshal(bytes, &stats); err != nil {
		t.Fatalf("unmarshalling expected nil err, actual err: %v", err)
	}

	if stats.CommonAPIData.QueryParams != "" {
		t.Errorf(`unmarshalling stats.CommonAPIData.QueryParams expected "", actual %v`, stats.CommonAPIData.QueryParams)
	}

	statsDate, err := time.Parse(srvhttp.CommonAPIDataDateFormat, stats.CommonAPIData.DateStr)
	if err != nil {
		t.Errorf(`stats.CommonAPIData.DateStr expected format %v, actual %v`, srvhttp.CommonAPIDataDateFormat, stats.CommonAPIData.DateStr)
	}
	if beforeStatsMarshall.Truncate(time.Second).After(statsDate) || statsDate.Truncate(time.Second).After(afterStatsMarshall.Truncate(time.Second)) { // round to second, because CommonAPIDataDateFormat is second-precision
		t.Errorf(`unmarshalling stats.CommonAPIData.DateStr expected between %v and %v, actual %v`, beforeStatsMarshall, afterStatsMarshall, stats.CommonAPIData.DateStr)
	}
	if len(stats.Caches) > 0 {
		t.Errorf(`unmarshalling stats.Caches expected empty, actual %+v`, stats.Caches)
	}
}

func TestSystemComputedStats(t *testing.T) {
	stats := cache.ComputedStats()

	for stat, function := range stats {
		if strings.HasPrefix(stat, "system.") {
			computedStat := function(cache.ResultInfo{}, tc.LegacyTrafficServer{}, tc.TMProfile{}, tc.IsAvailable{})
			_, err := newStatEqual([]cache.ResultStatVal{{Val: float64(0)}}, computedStat)
			if err != nil {
				t.Errorf("expected no errors from newStatEqual: %s", err)
			}
		}
	}
}

func TestcompareAndAppendStatForInterface(t *testing.T) {
	var errs strings.Builder
	var limit uint64 = 1
	stat := interfaceStat{
		InterfaceName: "test",
		Stat:          uint64(5),
		StatName:      "test",
		Time:          time.Now(),
	}

	history := compareAndAppendStatForInterface(nil, errs, limit, stat)
	if errs.Len() > 0 {
		t.Errorf("Unexpected errors comparing previously non-existent interface stat: %s", errs.String())
	}
	if len(history) == 0 {
		t.Fatal("Empty history after comparing previously non-existent interface stat")
	}
	if len(history) > 1 {
		t.Fatalf("Too many stats returned from comparing previously non-existent interface stat: %d", len(history))
	}

	result := history[0]
	if result.Span != 1 {
		t.Errorf("Incorrect span comparing previously non-existent interface stat; want: 1, got: %d", result.Span)
	}

	if result.Time != stat.Time {
		t.Errorf("Incorrect time comparing previously non-existent interface stat; want: %v, got: %v", stat.Time, result.Time)
	}

	if v, ok := result.Val.(uint64); !ok {
		t.Errorf("Incorrect value type from comparing previously non-existent interface stat; want: uint64, got: %v", stat.Val)
	} else if v != stat.Stat {
		t.Errorf("Incorrect value from comparing previously non-existent interface stat; want: %d, got: %d", stat.Stat, v)
	}

	errs.Reset()

	history = compareAndAppendStatForInterface(history, errs, limit, stat)
	if errs.Len() > 0 {
		t.Errorf("Unexpected errors comparing previously non-existent interface stat: %s", errs.String())
	}
	if len(history) == 0 {
		t.Fatal("Empty history after comparing previously non-existent interface stat")
	}
	if len(history) > 1 {
		t.Fatalf("Too many stats returned from comparing previously non-existent interface stat: %d", len(history))
	}

	result = history[0]
	if result.Span != 2 {
		t.Errorf("Incorrect span comparing previously non-existent interface stat; want: 2, got: %d", result.Span)
	}

	if result.Time != stat.Time {
		t.Errorf("Incorrect time comparing previously non-existent interface stat; want: %v, got: %v", stat.Time, result.Time)
	}

	if v, ok := result.Val.(uint64); !ok {
		t.Errorf("Incorrect value type from comparing previously non-existent interface stat; want: uint64, got: %v", stat.Val)
	} else if v != stat.Stat {
		t.Errorf("Incorrect value from comparing previously non-existent interface stat; want: %d, got: %d", stat.Stat, v)
	}

	errs.Reset()
	stat.Stat = uint64(6)

	history = compareAndAppendStatForInterface(history, errs, limit, stat)
	if errs.Len() > 0 {
		t.Errorf("Unexpected errors comparing previously non-existent interface stat: %s", errs.String())
	}
	if len(history) == 0 {
		t.Fatal("Empty history after comparing previously non-existent interface stat")
	}
	if len(history) > 1 {
		t.Fatalf("Too many stats returned from comparing previously non-existent interface stat: %d", len(history))
	}

	result = history[0]
	if result.Span != 1 {
		t.Errorf("Incorrect span comparing previously non-existent interface stat; want: 1, got: %d", result.Span)
	}

	if result.Time != stat.Time {
		t.Errorf("Incorrect time comparing previously non-existent interface stat; want: %v, got: %v", stat.Time, result.Time)
	}

	if v, ok := result.Val.(uint64); !ok {
		t.Errorf("Incorrect value type from comparing previously non-existent interface stat; want: uint64, got: %v", stat.Val)
	} else if v != stat.Stat {
		t.Errorf("Incorrect value from comparing previously non-existent interface stat; want: %d, got: %d", stat.Stat, v)
	}
}
