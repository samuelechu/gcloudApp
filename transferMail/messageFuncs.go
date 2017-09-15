package transferMail

import (
    "log"
	"net/http"
    "bytes"
    //"github.com/samuelechu/cloudSQL"
    "github.com/buger/jsonparser"
    "github.com/samuelechu/jsonHelper"
)

func insertMessage(client *http.Client, labelMap map[string]string, messageId, sourceToken, destToken string){

	urlStr := "https://www.googleapis.com/gmail/v1/users/me/messages/" + messageId + "?format=raw" 
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    respBody := jsonHelper.GetRespBody(req, client)
    if len(respBody) == 0 {
         log.Print("Error: empty respBody")
         return
    }
    //log.Printf("HTTP PostForm/GET returned %v", string(respBody))

    raw, _, _, _ := jsonparser.Get(respBody, "raw")
    var messageLabels []string

    jsonparser.ArrayEach(respBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        labelId, _ := jsonparser.ParseString(value)
        messageLabels = append(messageLabels, "\"" + labelMap[labelId] + "\", ")
        
    }, "labelIds")

    //remove comma from last array element
    changedNdx := messageLabels[len(messageLabels)-1]
	changedNdx = changedNdx[0:len(changedNdx)-2]
	messageLabels[len(messageLabels)-1] = changedNdx


    log.Print("printing labels")
    log.Print(messageLabels)

	//post message
	    urlStr = "https://www.googleapis.com/upload/gmail/v1/users/me/messages?uploadType=multipart"
	    body := nopCloser{bytes.NewBufferString("--foo_bar\nContent-Type: application/json; charset=UTF-8\n\n{" +
	"\n\"raw\":\"" + raw + "\",\n\"labelIds\": " + messageLabels + "\n}" +
	"\n--foo_bar\nContent-Type: message/rfc822\n\nstringd\n--foo_bar--")} 

	    req, _ = http.NewRequest("POST", urlStr, body)
	    req.Header.Set("Authorization", "Bearer " + destToken)
	    req.Header.Set("Content-Type", "multipart/related; boundary=\"foo_bar\"")

	    respBody = jsonHelper.GetRespBody(req, client)
	    if len(respBody) == 0 {
	         log.Print("Error: empty respBody")
	         return
	    }
	    log.Printf("HTTP PostForm/GET returned %v", string(respBody))
}