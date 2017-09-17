package transferMail

import (
    "log"
	"net/http"
	"io"
    //"bytes"
    "github.com/samuelechu/cloudSQL"
    "github.com/buger/jsonparser"
    "golang.org/x/net/context"
    "google.golang.org/appengine/urlfetch"
    "github.com/samuelechu/jsonHelper"
)

type nopCloser struct { 
    io.Reader 
} 

func (nopCloser) Close() error { return nil } 

func startTransfer(ctx context.Context, curUserID, sourceToken, sourceID, destToken, destID string) {
    client := urlfetch.Client(ctx)

    urlStr := "https://www.googleapis.com/gmail/v1/users/me/threads"
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    respBody := jsonHelper.GetRespBody(req, client)
    if len(respBody) == 0 {
         log.Print("Error: empty respBody")
         return
    }
    //log.Printf("HTTP PostForm/GET returned %v", string(respBody))

    nextPage, _ := jsonparser.GetString(respBody, "nextPageToken")
    
    jsonparser.ArrayEach(respBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        thread_id, _, _, _ := jsonparser.Get(value, "id")
        if string(thread_id) != "" {
            log.Printf("Inserting into database: Thread %v", string(thread_id))
            cloudSQL.InsertThread(curUserID, string(thread_id))

        }
        
    }, "threads")

    for nextPage != "" {
        urlStr = "https://www.googleapis.com/gmail/v1/users/me/threads?pageToken=" + nextPage 
        req, _ = http.NewRequest("GET", urlStr, nil)
        req.Header.Set("Authorization", "Bearer " + sourceToken)

        respBody = jsonHelper.GetRespBody(req, client)
        if len(respBody) == 0 {
             log.Print("Error: empty respBody")
             return
        }

        nextPage, _ = jsonparser.GetString(respBody, "nextPageToken")

        jsonparser.ArrayEach(respBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
            thread_id, _, _, _ := jsonparser.Get(value, "id")
            if string(thread_id) != "" {
                //log.Printf("Inserting into database: Thread %v", string(thread_id))
                cloudSQL.InsertThread(curUserID, string(thread_id))

            }
            
        }, "threads")

    }

    
//get threads
	sourceThreads := cloudSQL.GetThreadsForUser(curUserID)
	//log.Printf("GetThreads returned %v", sourceThreads)
	log.Printf("curUserID : %v, sourceToken : %v, sourceID : %v, destToken : %v, destID : %v", curUserID, sourceToken, sourceID, destToken, destID)

    insertThreads(ctx, sourceThreads,sourceToken,destToken,curUserID)
}