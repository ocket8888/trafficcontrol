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

import "encoding/json"
import "fmt"

func ExampleHealthThreshold_String() {
	ht := HealthThreshold{Comparator: ">=", Val: 500}
	fmt.Println(ht)
	// Output: >=500.000000
}

func ExampleTMParameters_UnmarshalJSON() {
	const data = `{
		"health.connection.timeout": 5,
		"health.polling.url": "https://example.com/",
		"health.polling.format": "stats_over_http",
		"history.count": 1,
		"health.threshold.bandwidth": ">50",
		"health.threshold.foo": "<=500"
	}`

	var params TMParameters
	if err := json.Unmarshal([]byte(data), &params); err != nil {
		fmt.Printf("Failed to unmarshal: %v\n", err)
		return
	}
	fmt.Printf("timeout: %d\n", params.HealthConnectionTimeout)
	fmt.Printf("url: %s\n", params.HealthPollingURL)
	fmt.Printf("format: %s\n", params.HealthPollingFormat)
	fmt.Printf("history: %d\n", params.HistoryCount)
	fmt.Printf("# of Thresholds: %d - foo: %s\n", len(params.Thresholds), params.Thresholds["foo"])
	fmt.Printf("# of Aggregate Thresholds: %d - bandwidth: %s\n", len(params.AggregateThresholds), params.AggregateThresholds["bandwidth"])

	// Output: timeout: 5
	// url: https://example.com/
	// format: stats_over_http
	// history: 1
	// # of Thresholds: 1 - foo: <=500.000000
	// # of Aggregate Thresholds: 1 - bandwidth: >50.000000
}
