package oauth

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"fmt"
	"log"
	"bytes"
	"os"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/samuelechu/rstring"
)

func init() {
     http.HandleFunc("/askPermissions", askPermissions)
     http.HandleFunc("/getToken", getToken)
     http.HandleFunc("/testrefToken", getAccessToken)
}

func getAccessToken(w http.ResponseWriter, r *http.Request) {
        log.Print(r.URL.Query())
    log.Print("heyo")
    log.Print(r.URL.Query().Get("code"))

    authCode := r.URL.Query().Get("code")
    
    urlStr := "https://www.googleapis.com/oauth2/v4/token"
 
    bodyVals := url.Values{
        "client_id": {os.Getenv("CLIENT_ID")},
        "client_secret": {os.Getenv("CLIENT_SECRET")},
        "refresh_token":{"1/iDMKVLsBI8QC2KSjqwbdIvUkcdSFo8edj70unSDfjCM"}
        "grant_type": {"refresh_token"},
    }

    body := bytes.NewBufferString(bodyVals.Encode())

    log.Print(body)
    req, _ := http.NewRequest("POST", urlStr, body)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    log.Print("finished marshaling")

    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)

    resp, err := client.Do(req)
    if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
    }
    log.Print("am here")
    defer resp.Body.Close()
    respBody, _ := ioutil.ReadAll(resp.Body)
    fmt.Fprintf(w, "HTTP Post returned %v", string(respBody))

}

func askPermissions(w http.ResponseWriter, r *http.Request) {
	redirectUri := "https%3a%2f%2fgotesting-175718.appspot.com%2FgetToken"
	if appengine.IsDevAppServer(){
    	redirectUri = "https%3a%2f%2f8080-dot-2979131-dot-devshell.appspot.com%2FgetToken"
	}

    redirectString := `https://accounts.google.com/o/oauth2/v2/
		auth?scope=https%3a%2f%2fwww.googleapis.com%2fauth%2fgmail.readonly
		&access_type=offline
		&include_granted_scopes=true
		&prompt=consent
		&state=state_parameter_passthrough_value
		&redirect_uri=` + redirectUri + 
		`&response_type=code
		&client_id=` + os.Getenv("CLIENT_ID")

    redirectString = rstring.RemoveWhitespace(redirectString)
    //fmt.Fprint(w, redirectString)

    http.Redirect(w, r, redirectString, 301)
}

func getToken(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL.Query())
    log.Print("heyo")
    log.Print(r.URL.Query().Get("code"))

    authCode := r.URL.Query().Get("code")
    
    urlStr := "https://www.googleapis.com/oauth2/v4/token"

    redirectUri := "https://gotesting-175718.appspot.com/getToken"
    if appengine.IsDevAppServer(){
        redirectUri = "https://8080-dot-2979131-dot-devshell.appspot.com/getToken"
    }

    bodyVals := url.Values{
        "code": {authCode},
        "client_id": {os.Getenv("CLIENT_ID")},
        "client_secret": {os.Getenv("CLIENT_SECRET")},
        "redirect_uri": {redirectUri},
        "grant_type": {"authorization_code"},
    }

    body := bytes.NewBufferString(bodyVals.Encode())

    log.Print(body)
    req, _ := http.NewRequest("POST", urlStr, body)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    log.Print("finished marshaling")

    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)

    resp, err := client.Do(req)
    if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
    }
    log.Print("am here")
    defer resp.Body.Close()
    respBody, _ := ioutil.ReadAll(resp.Body)
    fmt.Fprintf(w, "HTTP Post returned %v", string(respBody))
}