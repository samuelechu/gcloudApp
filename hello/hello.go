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
    "github.com/samuelechu/oauth"
	_ "github.com/samuelechu/cloudSQL"
    _ "github.com/samuelechu/templateTest"
    _ "github.com/samuelechu/cookieTest"
    _ "google.golang.org/api/gmail/v1"
)

var indexTemplate *template.Template

func main() {


    http.Handle("/scripts/", http.FileServer(http.Dir("static")))
    http.Handle("/css/", http.FileServer(http.Dir("static")))
    http.Handle("/img/", http.FileServer(http.Dir("static")))

    http.HandleFunc("/", index)
    http.HandleFunc("/_ah/health", healthCheckHandler)
    
    indexTemplate = template.Must(template.ParseFiles("static/home.html"))

    log.Print("Listening on port 8080")
    http.ListenAndServe(":8080", nil)
    appengine.Main()
}

type AccountNames struct {
    Source string
    Destination string
}

func index(w http.ResponseWriter, r *http.Request) {

    log.Print("index was triggered!")

    oauth.GetCookies(w, r)

    sourceToken := ""
    destToken := ""
    
    sourceCookie, err := r.Cookie("source")
    if err == nil {
        sourceToken = sourceCookie.Value
    }

    destCookie, err := r.Cookie("destination")
    if err == nil {
        destToken = destCookie.Value
    }

    log.Printf("Source Cookie: %v\n", sourceCookie)
    log.Printf("Dest Cookie: %v\n", destCookie)

    _, sourceName := oauth.VerifyIDToken(w, r, sourceToken)
    _, destName := oauth.VerifyIDToken(w, r, destToken)


    log.Printf("Source Name: %v\n", sourceName)
    log.Printf("Dest Name: %v\n", destName)

    //names := AccountNames{Source: sourceName, Destination: destName,}
 
    //indexTemplate.Execute(w, names)

}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
     fmt.Fprint(w, "ok")
}

