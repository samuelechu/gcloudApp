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
	"github.com/samuelechu/rstring"
	_ "github.com/samuelechu/cloudSQL"
    _ "github.com/samuelechu/oauth"
    _ "google.golang.org/api/gmail/v1"
)

func main() {

     fs := http.FileServer(http.Dir("static"))
     http.Handle("/", fs)


     http.HandleFunc("/testTemplate", handler)
     http.HandleFunc("/rstring", handle)
     http.HandleFunc("/_ah/health", healthCheckHandler)
     http.HandleFunc("/authSuccess", authSuccessful)

     log.Print("Listening on port 8080")
     http.ListenAndServe(":8080", nil)
     appengine.Main()
}

func handler(w http.ResponseWriter, r *http.Request) {
    t := template.New("test") // Create a template.
    t, _ = t.ParseFiles("static/index.html")  // Parse template file.
    //user := GetUser() // Get current user infomration.
    t.Execute(w, nil)  // merge.
}

func authSuccessful(w http.ResponseWriter, r *http.Request){
     fmt.Fprintf(w, "Hallo!")
}

func handle(w http.ResponseWriter, r *http.Request) {

     if r.URL.Path != "/rstring" {
                http.NotFound(w, r)
                return
     }

	fmt.Fprint(w, rstring.Reverse("Hallos world!"))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
     fmt.Fprint(w, "ok")
}

