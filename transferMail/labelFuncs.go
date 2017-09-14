package transferMail

import (
    "log"
    "fmt"
	"net/http"
    "bytes"
    "github.com/buger/jsonparser"
    "github.com/samuelechu/jsonHelper"
)

//returns map of source label id to corresponding destination label id
func getLabelMap(client *http.Client, sourceToken, destToken string) map[string]string {

	var sourceEmail string
	labelIdMap := make(map[string]string)
	destLabels := make(map[string]string)

	//get destLabels
	urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels"
    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels"
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + destToken)

    //get Labels from destination account
    respBodyDest := jsonHelper.GetRespBody(req, client)
    if len(respBodyDest) == 0 {
         log.Print("Error: empty respBody")
         return labelIdMap
    }

    jsonparser.ArrayEach(respBodyDest, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	    labelName, _ := jsonparser.GetString(value, "name")
	    labelId, _ := jsonparser.GetString(value, "id")

	    destLabels[labelName] = labelId
	    
	}, "labels")


    log.Print(string(respBodyDest))

    //get sourceEmail
    urlStr = "https://www.googleapis.com/oauth2/v1/userinfo"

    req, _ = http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    respBodyUserInfo := jsonHelper.GetRespBody(req, client)
    if len(respBodyUserInfo) == 0 {
         log.Print("Error: empty respBody")
         return labelIdMap
    }

    sourceEmail, _ = jsonparser.GetString(respBodyUserInfo, "email")

    //create map:   map[sourceLabel id] = destLabels[sourceEmail + name] <--maps to an id
    req, _ = http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    respBodySource := jsonHelper.GetRespBody(req, client)
    if len(respBodySource) == 0 {
         log.Print("Error: empty respBody")
         return labelIdMap
    }

    jsonparser.ArrayEach(respBodyDest, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	    labelName, _ := jsonparser.GetString(value, "name")
	    labelId, _ := jsonparser.GetString(value, "id")

	    labelIdMap[labelId] = destLabels[sourceEmail + "/" + labelName]
	    
	}, "labels")

    log.Print(string(respBodySource))

	log.Print("\n\n\nPrinting labelIdMap")
    for key, value := range labelIdMap {
    	log.Print("Key:", key, " Value:", value)
	}

	return labelIdMap
}

func createNewLabel(client *http.Client, access_token, name, messageVis, labelVis string){
	urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels"
    bodyStr := fmt.Sprintf(`{"name": "%v", "messageListVisibility": "%v", "labelListVisibility": "%v"}`, name, messageVis, labelVis)
    jsonStr := []byte(bodyStr)

    req, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonStr))
    req.Header.Set("Authorization", "Bearer " + access_token)
    req.Header.Set("Content-Type", "application/json")

    respBody := jsonHelper.GetRespBody(req, client)
    if len(respBody) == 0 {
         log.Print("Error: empty respBody")
         return
    }

    log.Print(string(respBody))
}

func addMissingLabels(client *http.Client, sourceToken, destToken string){

    var sourceEmail string
    destLabels := make(map[string]bool)

    urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels" //testTransfer label
    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels"
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + destToken)

    //get Labels from destination account
    respBodyDest := jsonHelper.GetRespBody(req, client)
    if len(respBodyDest) == 0 {
         log.Print("Error: empty respBody")
         return
    }

    jsonparser.ArrayEach(respBodyDest, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        labelName, err := jsonparser.GetString(value, "name")
        if err != nil {
            log.Print("Error: invalid label")
            return
        }
        destLabels[labelName] = true
        
    }, "labels")

    //get labels from source account and add if not in dest
    req, _ = http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    respBodySource := jsonHelper.GetRespBody(req, client)
    if len(respBodySource) == 0 {
         log.Print("Error: empty respBody")
         return
    }

    fields := [][]string{
        []string{"name"},
        []string{"messageListVisibility"},
        []string{"labelListVisibility"},
        []string{"type"},
    }

    urlStr = "https://www.googleapis.com/oauth2/v1/userinfo"

    req, _ = http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    respBodyUserInfo := jsonHelper.GetRespBody(req, client)
    if len(respBodySource) == 0 {
         log.Print("Error: empty respBody")
         return
    }

    sourceEmail, _ = jsonparser.GetString(respBodyUserInfo, "email")

    //add main email label
    if !destLabels[sourceEmail] {
    	createNewLabel(client, destToken, sourceEmail, "show", "labelShow")
    }

    //add nested email labels
    jsonparser.ArrayEach(respBodySource, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        var name, messageListVisibility, labelListVisibility string// labelType string
        
        jsonparser.EachKey(value, func(idx int, value []byte, vt jsonparser.ValueType, err error){
            switch idx {
	            case 0:
	                name, _ = jsonparser.ParseString(value)
	            case 1:
	                messageListVisibility, _ = jsonparser.ParseString(value)
	            case 2:
	                labelListVisibility, _ = jsonparser.ParseString(value)
	            // case 3:
             //    	labelType, _ = jsonparser.ParseString(value)
            }
        }, fields...)

        if !destLabels[sourceEmail + "/" + name] {
        	// if labelType == "system" {
        	// 	labelListVisibility = "labelHide"
        	// }
        	if labelListVisibility == "" {
        		messageListVisibility= "show"
        		labelListVisibility = "labelShow"
        	}

        	log.Printf("Adding new Label: %v", sourceEmail + "/" + name)
            createNewLabel(client, destToken, sourceEmail + "/" + name, messageListVisibility, labelListVisibility)
        }
    }, "labels")
}