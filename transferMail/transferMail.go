package transferMail

import (
	"log"
	"net/http"
    "io/ioutil"
    //"golang.org/x/net/context"
    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
    //"google.golang.org/appengine/runtime"
	"github.com/samuelechu/oauth"
    "github.com/samuelechu/cloudSQL"
    "github.com/buger/jsonparser"
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

    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/messages/15e5d6ed5bb68a29?format=raw"
//retrieve threads

    urlStr := "https://www.googleapis.com/gmail/v1/users/me/threads?labelIds=Label_8" //testTransfer label
    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels"
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)

    resp, err := client.Do(req)

    if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
    }


    body := resp.Body
    defer body.Close()

    if body == nil {
        http.Error(w, "Response body not found", 400)
        return
    }

    respBody, _ := ioutil.ReadAll(body)
    log.Printf("HTTP PostForm/GET returned %v", string(respBody))

    // if message_id, ok := jsonparser.GetString(respBody, "id"); ok == nil{
    //     log.Printf("ID of messsage was %v", message_id)
    // }

    s, _ := jsonparser.GetString(respBody, "nextPageToken")
    log.Printf("Token is %v", s)
    
    jsonparser.ArrayEach(respBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        thread_id, _, _, _ := jsonparser.Get(value, "id")
        if string(thread_id) != "" {
            log.Printf("Inserting into database: Thread %v", string(thread_id))
            cloudSQL.InsertThread(curUserID, string(thread_id))

        }
        
    }, "threads")


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

