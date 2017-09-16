package transferMail

import (
    "log"
	"net/http"
    //"bytes"
    "github.com/samuelechu/cloudSQL"
    "github.com/buger/jsonparser"
    "github.com/samuelechu/jsonHelper"
)

func insertThreads(client *http.Client, sourceThreads []string, sourceToken, destToken, curUserID string){

	labelMap := getLabelMap(client,sourceToken,destToken)
    log.Print("\n\n\nPrinting labelIdMap")
        for key, value := range labelMap {
        log.Print("Key:", key, " Value:", value)
    }

	for _, threadId := range sourceThreads {
		insertThread(client, labelMap, threadId, sourceToken, destToken, curUserID)
	}
	
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