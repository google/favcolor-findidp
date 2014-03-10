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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
)

func fetchJSON(uri string, handles Handles, task string) (res interface{}, err error) {
	resp, err := handles.client.Get(uri)
	var x string
	if err != nil {
		x = fmt.Sprintf("%s HTTP fetch error for %s: %s", task, uri, err)
		handles.logger.logError(x)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		x = fmt.Sprintf("%s HTTP status %d on %s", task, resp.StatusCode, uri)
		err = errors.New(x)
		handles.logger.logError(x)
		return
	}
	ctype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil ||
		(ctype != "application/json" && ctype != "application/jrd+json") {
		x = fmt.Sprintf("Bogus %s media type for %s: %s", task, uri, ctype)
		handles.logger.logError(x)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		x = fmt.Sprintf("Can't read %s response for %s: %s", task, uri, ctype)
		handles.logger.logError(x)
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		x = fmt.Sprintf("Can't parse %s json for %s: %s", task, uri, err)
		handles.logger.logError(x)
		return
	}
	return
}
