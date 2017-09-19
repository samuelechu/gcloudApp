package transferMail

import (
	"log"
	"net/http"
    "golang.org/x/net/context"
    "google.golang.org/appengine"
    //"google.golang.org/appengine/urlfetch"
    "google.golang.org/appengine/runtime"
	"github.com/samuelechu/oauth"
    "github.com/samuelechu/cloudSQL"
    //"github.com/buger/jsonparser"
    //"github.com/samuelechu/jsonHelper"
)

func init() {
     http.HandleFunc("/transferStart", transferEmail)
}

func transferEmail(w http.ResponseWriter, r *http.Request) {
	var curUserID, sourceToken, sourceID, destToken, destID string

    curUserCookie, err := r.Cookie("current_user")
    if err == nil {
        curUserID = curUserCookie.Value
    }
    
    sourceCookie, err := r.Cookie("source")
    if err == nil {
        sourceToken = sourceCookie.Value
    }

    destCookie, err := r.Cookie("destination")
    if err == nil {
        destToken = destCookie.Value
    }

    sourceID, _, _ = oauth.GetUserInfo(w, r, sourceToken)
    destID, _, _ = oauth.GetUserInfo(w, r, destToken)

    log.Printf("Source ID: %v\n", sourceID)
    log.Printf("Dest ID: %v\n", destID)

    //send job to database
    cloudSQL.InsertJob(curUserID, sourceID, destID)

    ctx := appengine.NewContext(r)

    err = runtime.RunInBackground(ctx, func(ctx context.Context) {
        startTransfer(ctx, curUserID, sourceToken, sourceID, destToken, destID)
    })

    if err != nil {
            log.Printf("Could not start background thread: %v", err)
            return
    }

    redirectString := "https://gotesting-175718.appspot.com"
    if appengine.IsDevAppServer(){
        redirectString = "https://8080-dot-2979131-dot-devshell.appspot.com"
    }
    http.Redirect(w, r, redirectString, 302)
}

