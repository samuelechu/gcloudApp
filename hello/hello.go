// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample helloworld is a basic App Engine flexible app.
package main

import (
    "google.golang.org/appengine"
    "fmt"
	"log"
    //"io/ioutil"
	"net/http"
	"github.com/samuelechu/rstring"
	_ "github.com/samuelechu/cloudSQL"
    _ "google.golang.org/api/gmail/v1"
)

func main() {

     fs := http.FileServer(http.Dir("static"))
     http.Handle("/", fs)

     http.HandleFunc("/rstring", handle)
     http.HandleFunc("/_ah/health", healthCheckHandler)
     http.HandleFunc("/drivePermissions", askPermissions)
     http.HandleFunc("/checkToken", checkToken)

     log.Print("Listening on port 8080")
     http.ListenAndServe(":8080", nil)
     appengine.Main()
}

func checkToken(w http.ResponseWriter, r *http.Request) {
    log.Print(r.URL.Query())
    log.Print("heyo")
    log.Print(r.URL.Query().Get("code"))

    redirectUri := "https://gotesting-175718.appspot.com"
    if appengine.IsDevAppServer(){
        redirectUri = "https://8080-dot-2979131-dot-devshell.appspot.com"
    }

    http.Redirect(w, r, redirectUri, 301)


}

func askPermissions(w http.ResponseWriter, r *http.Request) {

redirectUri := "https%3a%2f%2fgotesting-175718.appspot.com/checkToken"
if appengine.IsDevAppServer(){
    redirectUri = "https%3a%2f%2f8080-dot-2979131-dot-devshell.appspot.com/checkToken"
}

    redirectString := `https://accounts.google.com/o/oauth2/v2/
auth?scope=https%3a%2f%2fwww.googleapis.com%2fauth%2fgmail.readonly
&access_type=offline
&include_granted_scopes=true
&state=state_parameter_passthrough_value
&redirect_uri=` + redirectUri + 
`&response_type=code
&client_id=65587295914-kbl4e2chuddg9ml7d72f6opqhddl62fv.apps.googleusercontent.com`

    redirectString = rstring.RemoveWhitespace(redirectString)
    //fmt.Fprint(w, redirectString)

    http.Redirect(w, r, redirectString, 301)

    
    // resp, err := http.Get(`https://accounts.google.com/o/oauth2/v2/
    //     auth?scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive.metadata.readonly
    //     &state=state_parameter_passthrough_value
    //     &redirect_uri=http%3a%2f%2fwww.example.com%2foauth2callback
    //     &response_type=token
    //     &client_id=65587295914-kbl4e2chuddg9ml7d72f6opqhddl62fv.apps.googleusercontent.com`)

    // if err != nil {
    //     body, _ := ioutil.ReadAll(resp.Body)
    //     bodyString := string(body)
    //     fmt.Println(bodyString)
    // }
    
    


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

