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
