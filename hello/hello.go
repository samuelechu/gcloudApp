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
    "html/template"
	_ "github.com/samuelechu/cloudSQL"
    _ "github.com/samuelechu/oauth"
    _ "github.com/samuelechu/templateTest"
    _ "github.com/samuelechu/cookieTest"
    _ "google.golang.org/api/gmail/v1"
)

func main() {

     //fs := http.FileServer(http.Dir("static"))
     http.HandleFunc("/", index)

     http.Handle("/scripts", http.FileServer(http.Dir("static")))
    // http.HandleFunc("/index.html", index)

     http.HandleFunc("/_ah/health", healthCheckHandler)

     log.Print("Listening on port 8080")
     http.ListenAndServe(":8080", nil)
     appengine.Main()
}

type Person struct {
    UserName string
}

func index(w http.ResponseWriter, r *http.Request) {
    t := template.New("index.html")


    p := Person{UserName: "Sam"}
    t, _ = t.ParseFiles("static/index.html")
 
    t.Execute(w, p)

}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
     fmt.Fprint(w, "ok")
}

