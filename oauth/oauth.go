package oauth

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"fmt"
	"log"
	"bytes"
	"os"
	"io/ioutil"
    "encoding/json"
	"net/http"
	"net/url"
)

func init() {
     http.HandleFunc("/askPermissions", askPermissions)
     http.HandleFunc("/getToken", getToken)
     http.HandleFunc("/testrefToken", getAccessToken)
}

type RespBody struct{
    Access_token    string
    Expires_in      float64
    Token_type      string
}

func getAccessToken(w http.ResponseWriter, r *http.Request) {

    
    urlStr := "https://www.googleapis.com/oauth2/v4/token"
 
    bodyVals := url.Values{
        "client_id": {os.Getenv("CLIENT_ID")},
        "client_secret": {os.Getenv("CLIENT_SECRET")},
        "refresh_token":{"1/08fGrbeZdKkEJmoNHhKqWxZuVvNWjSc_JjN1aMExhaU"},
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
    log.Printf("HTTP Post returned %v", string(respBody))
    // fmt.Fprintf(w, "HTTP Post returned %v", string(respBody))

    var rb RespBody 
    if resp.Body == nil {
        http.Error(w, "Please send a request body", 400)
        return
    }
    err = json.Unmarshal(respBody, &rb)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    fmt.Fprintf(w, "HTTP Post returned %v %v %v", rb.Access_token, rb.Expires_in, rb.Token_type)
    
    

}

func askPermissions(w http.ResponseWriter, r *http.Request) {
	
    //request will be format :   /askPermissions?(source||destination)
    callType := r.URL.Query().Get("type")
    permissions := ""

    switch callType {
        case "source":
            permissions = "https://www.googleapis.com/auth/gmail.readonly"
        case "destination":
            permissions = "https://www.googleapis.com/auth/gmail.insert"
        default:
            http.Error(w, "must specify in queryString type : source || destination", 400)
            return
    }

    redirectUri := "https://gotesting-175718.appspot.com/getToken"
	if appengine.IsDevAppServer(){
    	redirectUri = "https://8080-dot-2979131-dot-devshell.appspot.com/getToken"
	}

    queryVals := url.Values{
        "scope" : {"profile " + permissions},
        "access_type" : {"offline"},
        "include_granted_scopes" : {"true"},
        "prompt" : {"consent"},
        "state" : {"state_parameter_passthrough_value"},
        "redirect_uri" : {redirectUri},
        "response_type" : {"code"},
        "client_id" : {os.Getenv("CLIENT_ID")},
    }

    queryString := queryVals.Encode()

    redirectString := "https://accounts.google.com/o/oauth2/v2/auth?" + queryString

    fmt.Fprint(w, redirectString)

    //http.Redirect(w, r, redirectString, 301)
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