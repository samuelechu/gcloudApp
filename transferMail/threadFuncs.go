package transferMail

import (
	"time"
    "log"
	"net/http"
    //"bytes"
    "google.golang.org/appengine/urlfetch"
    "github.com/samuelechu/cloudSQL"
    "github.com/buger/jsonparser"
    "github.com/samuelechu/jsonHelper"
)


func accessTokenUpdater(client *http.Client, done chan int, curUserID string, sourceToken, destToken *string) {
	sourceID, destID := cloudSQL.GetJob(curUserID)
	log.Printf("sourceID: %v, destID: %v", sourceID, destID)
	*sourceToken = getAccessToken(client, sourceID)
	*destToken = getAccessToken(client, destID)

	for {
		select {
			case <-time.After(10 * time.Second):
				*sourceToken = getAccessToken(client, sourceID)
				*destToken = getAccessToken(client, destID)

			case <-done:
				return

		}
	}

}

func insertThreads(ctx context.Context, sourceThreads []string, sourceToken, destToken, curUserID string){

	client := urlfetch.Client(ctx)

	done := make(chan int)

	err = runtime.RunInBackground(ctx, func(ctx context.Context) {
    	accessTokenUpdater(client, done, curUserID, &sourceToken, &destToken)    
    })

    if err != nil {
        log.Printf("Could not start background thread: %v", err)
        return
    }

	labelMap := getLabelMap(client,sourceToken,destToken)
    log.Print("\n\n\nPrinting labelIdMap")
        for key, value := range labelMap {
        log.Print("Key:", key, " Value:", value)
    }

	for _, threadId := range sourceThreads {
		log.Print("The sourceToken is %v, destToken: %v", sourceToken, destToken)
		insertThread(client, labelMap, threadId, sourceToken, destToken, curUserID)
	}

	//stop background accessTokenUpdating thread
	done <- 1
	<-time.After(3 * time.Second)
	
}

func insertThread(client *http.Client, labelMap map[string]string, threadID, sourceToken, destToken, curUserID string){

	urlStr := "https://www.googleapis.com/gmail/v1/users/me/threads/" + threadID + "?format=minimal"
    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels"
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    //get Labels from destination account
    respBody := jsonHelper.GetRespBody(req, client)
    if len(respBody) == 0 {
         log.Print("Error: empty respBody")
         return
    }
    log.Print(string(respBody))

    threadId := ""

    jsonparser.ArrayEach(respBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        messageId, _ := jsonparser.GetString(value, "id")

        threadId = insertMessage(client, labelMap, threadId, messageId, sourceToken, destToken)

        if threadId == "" {
            log.Printf("Error: insertMessage failed for message %v", messageId)
            return
        }
        
    }, "messages")

    cloudSQL.MarkThreadDone(curUserID, threadID)
}