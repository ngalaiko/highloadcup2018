package answers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

// Answer is a single test data entry.
type Answer struct {
	Method     string                 `json:",omitempty"`
	StatusCode int                    `json:",omitempty"`
	URL        *url.URL               `json:",omitempty"`
	Response   map[string]interface{} `json:",omitempty"`
}

// Report is a test report with results.
type Report struct {
	Expected     *Answer  `json:"expected"`
	Got          *Answer  `json:"got"`
	URL          *url.URL `json:"url"`
	ResponseTime string   `json:"response_time"`
}

// Test tests an answer.
func (a *Answer) Test() (bool, *Report, error) {
	req, err := http.NewRequest(a.Method, a.URL.String(), http.NoBody)
	if err != nil {
		return false, nil, err
	}

	report := &Report{
		Expected: &Answer{},
		Got:      &Answer{},
	}

	start := time.Now()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()

	report.ResponseTime = time.Since(start).String()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, nil, err
	}

	jsonBody := map[string]interface{}{}
	_ = json.Unmarshal(body, &jsonBody)

	success := true
	if resp.StatusCode != a.StatusCode {
		report.URL = a.URL
		report.Expected.StatusCode = a.StatusCode
		report.Got.StatusCode = resp.StatusCode
		success = false
	}

	if !reflect.DeepEqual(a.Response, jsonBody) {
		report.URL = a.URL
		report.Expected.Response = a.Response
		report.Got.Response = jsonBody
		success = false
	}

	return success, report, nil
}
