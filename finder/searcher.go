/*
Copyright [2014] Google, Inc.

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
package findIDP

import (
	"fmt"
)

type ResultType uint8

const (
	TimeoutType ResultType = iota
	MXType
	WebfingerType
	WebFistType
	DomainMatchType
	DatabaseType
	GoogleProbeType
	TestingType
)

var labels = []string{
	"Timeout", "DNS MX record", "Webfinger", "WebFist",
	"Domain-name matcher", "Local database", "Google hosting probe",
	"Unit testing",
}

func resultLabel(r ResultType) string {
	return labels[r]
}

type ResultClass struct {
	strength int
	verified bool
}

var ResultStrengths = map[ResultType]ResultClass{
	MXType:          {100, false},
	WebfingerType:   {200, true},
	WebFistType:     {200, true},
	DomainMatchType: {80, false},
	DatabaseType:    {1000, true},
	GoogleProbeType: {70, false},
	TestingType:     {300, false},
}

type SearchResult struct {
	rtype ResultType
	idps  []IDP
}

func (s SearchResult) String() string {
	return fmt.Sprintf("%s: v=%v len=%d, s=%d", labels[s.rtype],
		ResultStrengths[s.rtype].verified,
		len(s.idps), ResultStrengths[s.rtype].strength)
}

type Searcher interface {
	Search(email *EMail, c chan SearchResult, handles Handles)
	Label() string
}

var Searchers []Searcher

func RegisterSearcher(s Searcher) {
	Searchers = append(Searchers, s)
}
