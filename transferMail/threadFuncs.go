package transferMail

import (
    "log"
	"net/http"
    //"bytes"
    //"github.com/samuelechu/cloudSQL"
    "github.com/buger/jsonparser"
    "github.com/samuelechu/jsonHelper"
)

func insertThreads(client *http.Client, sourceThreads []string, sourceToken, destToken string){

	labelMap := getLabelMap(client,sourceToken,destToken)
    log.Print("\n\n\nPrinting labelIdMap")
        for key, value := range labelMap {
        log.Print("Key:", key, " Value:", value)
    }

	threadId := sourceThreads[0]
	insertThread(client, labelMap, threadId, sourceToken, destToken)
}

func insertThread(client *http.Client, labelMap map[string]string, threadId, sourceToken, destToken string){

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

    messageID, _ := jsonparser.GetString(respBody, "messages", "[0]", "id")

    log.Printf("MessageId is :: %v", messageID)

    threadID := insertMessage(client, labelMap, "", messageID, sourceToken, destToken)

    messageID2, _ := jsonparser.GetString(respBody, "messages", "[1]", "id")

    threadID := insertMessage(client, labelMap, threadID, messageID2, sourceToken, destToken)

    //client *http.Client, messageId, sourceToken, destToken string

 //    jsonparser.ArrayEach(respBodyDest, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	//     labelName, _ := jsonparser.GetString(value, "name")
	//     labelId, _ := jsonparser.GetString(value, "id")

	//     destLabels[labelName] = labelId
	    
	// }, "labels")
}