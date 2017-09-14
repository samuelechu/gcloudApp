package transferMail

import (
    "log"
	"net/http"
    //"bytes"
    //"github.com/samuelechu/cloudSQL"
    //"github.com/buger/jsonparser"
    "github.com/samuelechu/jsonHelper"
)

func insertThreads(client *http.Client, sourceThreads []string, sourceToken, destToken string){

	threadId := sourceThreads[0]
	insertThread(client,threadId,sourceToken, destToken)
}

func insertThread(client *http.Client, threadId, sourceToken, destToken string){

	urlStr := "https://www.googleapis.com/gmail/v1/users/me/threads/" + threadId + "?format=minimal"
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

 //    jsonparser.ArrayEach(respBodyDest, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	//     labelName, _ := jsonparser.GetString(value, "name")
	//     labelId, _ := jsonparser.GetString(value, "id")

	//     destLabels[labelName] = labelId
	    
	// }, "labels")
}