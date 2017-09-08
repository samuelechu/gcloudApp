package transferMail

import (
	"log"
	"net/http"
    "io/ioutil"
    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
	"github.com/samuelechu/oauth"
)

func init() {
     http.HandleFunc("/transferStart", transferEmail)
}

func transferEmail(w http.ResponseWriter, r *http.Request) {
	var sourceToken, sourceID, destToken, destID string
    
    sourceCookie, err := r.Cookie("source")
    if err == nil {
        sourceToken = sourceCookie.Value
    }

    destCookie, err := r.Cookie("destination")
    if err == nil {
        destToken = destCookie.Value
    }

    sourceID, _ = oauth.GetUserInfo(w, r, sourceToken)
    destID, _ = oauth.GetUserInfo(w, r, destToken)

    log.Printf("Source ID: %v\n", sourceID)
    log.Printf("Dest ID: %v\n", destID)

    urlStr := "https://www.googleapis.com/gmail/v1/users/me/messages/15e5d6ed5bb68a29?format=raw"

    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)

    resp, err := client.Do(req)

    body := resp.Body
    defer body.Close()

    if body == nil {
        http.Error(w, "Response body not found", 400)
        return
    }

    respBody, _ := ioutil.ReadAll(body)
    log.Printf("HTTP PostForm/GET returned %v", string(respBody))





}