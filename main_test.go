// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloWorldHandler(t *testing.T) {
	tests := []struct {
		label string
		want  string
		name  string
	}{
		{
			label: "default",
			want:  "Hello World!\n",
			name:  "",
		},
		{
			label: "override",
			want:  "Hello Override!\n",
			name:  "Override",
		},
	}

	for _, test := range tests {
		a := &app{pubsubVerificationToken: "testTokenNotNeededYet"}

		reader := strings.NewReader(`{"name": "` + test.name + `"}`)
		req := httptest.NewRequest("GET", "/", reader)
		rr := httptest.NewRecorder()
		a.helloWorldHandler(rr, req)

		if got := rr.Body.String(); got != test.want {
			t.Errorf("%s: got %q, want %q", test.label, got, test.want)
		}
	}
}

func TestCreateFromPushMessageGETFails(t *testing.T) {
	a := &app{pubsubVerificationToken: "testTokenNotNeededYet"}

	reader := strings.NewReader(`{"name": "ThisContentShouldntMatter"}`)
	req := httptest.NewRequest("GET", "/pubsub", reader)
	resp := httptest.NewRecorder()
	a.createFromPushRequestHandler(resp, req)

	want := http.StatusMethodNotAllowed
	if got := resp.Code; got != want {
		t.Errorf("got code=%d, want %d", got, want)
	}

	want = 0
	if got := len(a.pubSubMessages); got != want {
		t.Errorf("got len=%d, want %d", got, want)
	}
}

func TestCreateFromPushMessageIncorrectBodyFails(t *testing.T) {
	a := &app{pubsubVerificationToken: "testTokenNotNeededYet"}

	reader := strings.NewReader(`ThisShouldntPassParsing`)
	req := httptest.NewRequest("POST", "/pubsub", reader)
	resp := httptest.NewRecorder()
	a.createFromPushRequestHandler(resp, req)

	want := http.StatusBadRequest
	if got := resp.Code; got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestClearMessages(t *testing.T) {
	a := &app{pubsubVerificationToken: "testTokenNotNeededYet"}

	a.pubSubMessages = append(a.pubSubMessages, pushRequest{}, pushRequest{})
	if got := len(a.pubSubMessages); got != 2 {
		t.Errorf("got len=%d, want %d", got, 2)
	}

	reader := strings.NewReader(`{"name": "ThisContentShouldntMatter"}`)
	req := httptest.NewRequest("POST", "/clear", reader)
	resp := httptest.NewRecorder()
	a.clearMessagesHandler(resp, req)

	want := http.StatusOK
	if got := resp.Code; got != want {
		t.Errorf("got code=%d, want %d", got, want)
	}

	want = 0
	if got := len(a.pubSubMessages); got != want {
		t.Errorf("got len=%d, , want %d", got, want)
	}
}

// func TestCreateFromPushMessage(t *testing.T) {
// 	tests := []struct {
// 		label   string
// 		want    string
// 		wantLen uint
// 		name    string
// 	}{
// 		{
// 			label: "default",
// 			want:  "Hello World!\n",
// 			name:  "",
// 		},
// 	}
//
// 	for _, test := range tests {
// 		a := &app{pubsubVerificationToken: "testTokenNotNeededYet"}
//
// 		reader := strings.NewReader(`{"name": "` + test.name + `"}`)
// 		req := httptest.NewRequest("POST", "/", reader)
// 		rr := httptest.NewRecorder()
// 		a.helloWorldHandler(rr, req)
//
// 		if got := rr.Body.String(); got != test.want {
// 			t.Errorf("%s: got %q, want %q", test.label, got, test.want)
// 		}
// 	}
// }
