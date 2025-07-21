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
	"time"

	support "cloud.google.com/go/support/apiv2"
	"google.golang.org/api/pubsub/v1"
	//supportpb "cloud.google.com/go/support/apiv2/supportpb"
)

func main() {
	log.Print("starting server...")
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/pubsub", createFromPushRequestHandler)
	//http.HandleFunc("/listPushMessages", listPushMessagesHandler)
	//http.HandleFunc("/listCustomMessages", listCustomMessagesHandler)

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

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var d struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "message failed to decode!!\r\n")
		return
	}

	name := d.Name
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

// pushRequest represents the payload of a Pub/Sub push message.
type pushRequest struct {
	Message      pubsub.PubsubMessage `json:"message"`
	Subscription string               `json:"subscription"`
}

// createFromPushRequestHandler is an HTTP Cloud Function with a request parameter.
func createFromPushRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, err := support.NewCaseClient(ctx)
	if err != nil {
		fmt.Fprintf(w, "Failed to create client: %v", err)
		return
	}
	defer c.Close()

	var pr pushRequest
	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		fmt.Fprint(w, "PubSub message failed to decode!!\r\n")
		return
	}

	notificationType := pr.Message.Attributes["notificationType"]
	resourceName := pr.Message.Attributes["resourceName"]

	fmt.Fprintf(w, "Received Msg ID: %s at timestamp Z: %s\r\n", pr.Message.MessageId, time.Now())
	fmt.Fprintf(w, "Subscription: %s\r\n", html.EscapeString(pr.Subscription))
	fmt.Fprintf(w, "ResourceName: %s\r\n", html.EscapeString(resourceName))
	fmt.Fprintf(w, "NotificationType: %s\r\n", html.EscapeString(notificationType))
}
