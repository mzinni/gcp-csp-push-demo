// Copyright 2025 Google LLC
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

// [START cloudrun_helloworld_service]

// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	support "cloud.google.com/go/support/apiv2"
	//supportpb "cloud.google.com/go/support/apiv2/supportpb"
)

func main() {
	log.Print("starting server...")
	http.HandleFunc("/", handler)
	http.HandleFunc("/pubsub", handlePubSub)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

// handlePubSub is an HTTP Cloud Function with a request parameter.
func handlePubSub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, err := support.NewCaseClient(ctx)
	if err != nil {
		fmt.Fprintf(w, "Failed to create client: %v", err)
		return
	}
	defer c.Close()

	var d struct {
		Subscription string `json:"subscription"`
		Message      struct {
			//Attributes map[string]string `json:"attributes"`
			ResourceName     string `json:"resourceName"`
			NotificationType string `json:"notificationType"`
		} `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "PubSub message failed to decode!!\r\n")
		return
	}

	fmt.Fprintf(w, "Subscription: %s\r\n", html.EscapeString(d.Subscription))
	fmt.Fprintf(w, "ResourceName: %s\r\n", html.EscapeString(d.Message.ResourceName))
	fmt.Fprintf(w, "NotificationType: %s\r\n", html.EscapeString(d.Message.NotificationType))
}
