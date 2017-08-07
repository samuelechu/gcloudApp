// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample helloworld is a basic App Engine flexible app.
package main

import (
     "google.golang.org/appengine"
     "fmt"
	"log"
	"net/http"
	"github.com/samuelechu/string"
	_ "github.com/samuelechu/cloudSQL"
)

func main() {
     http.HandleFunc("/", handle)
     http.HandleFunc("/_ah/health", healthCheckHandler)
     log.Print("Listening on port 8080")
     http.ListenAndServe(":8080", nil)
     appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, string.Reverse("Hallos world!"))

}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
     fmt.Fprint(w, "ok")
     }