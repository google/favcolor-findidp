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

type Crumbs struct {
	trail []SearchResult
}

func addCrumb(c *Crumbs, s SearchResult) {
	c.trail = append(c.trail, s)
}

func crumbTrail(c *Crumbs) []SearchResult {
	return c.trail
}

func (c Crumbs) Len() int {
	return len(c.trail)
}

func (c Crumbs) Swap(i, j int) {
	c.trail[i], c.trail[j] = c.trail[j], c.trail[i]
}

func (c Crumbs) Less(i, j int) bool {
	iEmpty := (len(c.trail[i].idps) == 0)
	jEmpty := (len(c.trail[j].idps) == 0)

	// a little tricky; only get past here if they’re the same
	if iEmpty != jEmpty {
		return jEmpty
	}
	iStrength := ResultStrengths[c.trail[i].rtype]
	jStrength := ResultStrengths[c.trail[j].rtype]
	if iStrength.verified != jStrength.verified {
		return iStrength.verified
	}
	return iStrength.strength > jStrength.strength
}
