package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// t.Setenv is only available on go 1.17+ while the
	// GAE standard env only has support up til 1.16 so
	// I'm a minor version off with the test runs
	test_version := "v1.1.1"
	test_sha := "2c16273"
	t.Setenv("APP_VERSION", test_version)
	t.Setenv("COMMIT_SHA", test_sha)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}

	resp := make(map[string]string)
	resp["version"] = test_version
	resp["sha"] = test_sha
	resp["status"] = "OK"
	jsonResp, _ := json.Marshal(resp)

	if rr.Body.String() != string(jsonResp) {
		t.Errorf(
			"unexpected body: got (%v) want (%v)",
			rr.Body.String(),
			string(jsonResp),
		)
	}
}

func TestIndexHandlerNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusNotFound,
		)
	}
}
