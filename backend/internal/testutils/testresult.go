package testutils

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

type TestResult struct {
	SRSRef          string `json:"srs_ref"`
	FSDRef          string `json:"fsd_ref"`
	TSDRef          string `json:"tsd_ref"`
	NoTestScript    string `json:"no_test_script"`
	Functional      string `json:"functional"`
	Endpoint        string `json:"endpoint"`
	ReqType         string `json:"req_type"`
	Parameter       string `json:"parameter"`
	ShouldBeSuccess string `json:"should_be_success"`
	Expectation     string `json:"expectation"`
	IsSuccess       string `json:"is_success"`
	Response        string `json:"response"`
	Date            string `json:"date"`
}

func (tr *TestResult) Log(t *testing.T, actualPass bool, resp *http.Response, respBody []byte) {
	t.Helper()
	tr.Date = time.Now().Format("2006-01-02")
	if actualPass {
		tr.IsSuccess = "true"
	} else {
		tr.IsSuccess = "false"
	}
	tr.Response = string(respBody)
	b, _ := json.MarshalIndent(tr, "", "  ")
	t.Log(string(b))
	if !actualPass {
		t.Errorf("FAIL: %s — expected %s but got status %d", tr.NoTestScript, tr.Expectation, resp.StatusCode)
	}
}
