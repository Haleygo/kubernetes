/*
Copyright 2021 The Kubernetes Authors.

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

package routes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPreCheckLogFileNameLength(t *testing.T) {
	oversizeFileName := fmt.Sprintf("%0256s", "a")
	correctFileName := fmt.Sprintf("%0255s", "a")

	// check oversize filename
	isOversize := preCheckLogFileNameLength(oversizeFileName)
	if !isOversize {
		t.Error("failed to check oversize filename")
	}

	// check correct filename when file is not exist
	NotOversize := preCheckLogFileNameLength(correctFileName)
	if NotOversize {
		t.Error("failed to check allowable length filename")
	}

	// create test file with correct name and check it
	_, err := os.Create(correctFileName)
	if err != nil {
		t.Error("failed to create test file")
	}
	defer os.Remove(correctFileName)
	NotOversize = preCheckLogFileNameLength(correctFileName)
	if NotOversize {
		t.Error("failed to check allowable length filename")
	}
}

func TestHttpServeFile(t *testing.T) {
	oversizeFileName := fmt.Sprintf("%0256s", "a")
	correctFileName := fmt.Sprintf("%0255s", "a")

	request, _ := http.NewRequest("", "", nil)

	// create test file with correct name and check it
	_, err := os.Create(correctFileName)
	if err != nil {
		t.Error("failed to create test file")
	}
	defer os.Remove(correctFileName)
	w2 := httptest.NewRecorder()
	http.ServeFile(w2, request, correctFileName)
	if w2.Code != http.StatusOK {
		t.Errorf("expected response code 200, get %d", w2.Code)
	}

	// http.ServeFile return 500 for ENAMETOOLONG instead of 404
	w1 := httptest.NewRecorder()
	http.ServeFile(w1, request, oversizeFileName)
	if w1.Code != http.StatusInternalServerError {
		t.Errorf("expected response code 500, get %d", w1.Code)
	}
}
