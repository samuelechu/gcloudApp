package transferMail

import (
	"log"
	"net/http"
    "golang.org/x/net/context"
    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
    "google.golang.org/appengine/runtime"
	"github.com/samuelechu/oauth"
    "github.com/samuelechu/cloudSQL"
    "github.com/buger/jsonparser"
    "github.com/samuelechu/jsonHelper"
)

type Values struct {
    m map[string]string
}

func (v Values) Get(key string) string {
    return v.m[key]
}

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





    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/messages/15e5d6ed5bb68a29?format=raw"
//retrieve threads




    
    res, _, _, _ := jsonparser.Get(respBody, "resultSizeEstimate")
    log.Printf("jsonparser returned %v", string(res))
    
    // err = runtime.RunInBackground(ctx, func(ctx context.Context) {
    //     startTransfer(ctx, curUserID, sourceToken, sourceID, destToken, destID)
    // })

    // if err != nil {
    //         log.Printf("Could not start background thread: %v", err)
    //         return
    // }

    redirectString := "https://gotesting-175718.appspot.com"
    if appengine.IsDevAppServer(){
        redirectString = "https://8080-dot-2979131-dot-devshell.appspot.com"
    }
    http.Redirect(w, r, redirectString, 302)
}

