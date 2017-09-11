package transferMail

import (
    "log"
	"net/http"
    "io/ioutil"
    "time"
    "github.com/samuelechu/cloudSQL"
    //"github.com/buger/jsonparser"
)

func startTransfer(curUserID, sourceToken, sourceID, destToken, destID string) {
	threads := cloudSQL.GetThreadsForUser(curUserID)
	log.Printf("GetThreads returned %v", threads)
	log.Printf("curUserID : %v, sourceToken : %v, sourceID : %v, destToken : %v, destID : %v", curUserID, sourceToken, sourceID, destToken, destID)
	
	urlStr := "https://www.googleapis.com/gmail/v1/users/me/messages/15d3d8e8de90ebcc" //testTransfer label
    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels"
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

	var client = &http.Client{
	  Timeout: time.Second * 10,
	}

    resp, err := client.Do(req)

    if err != nil {
    		log.Printf("Error: %v", http.StatusInternalServerError)
            return
    }
    
    body := resp.Body
    defer body.Close()

    if body == nil {
    	log.Print("Error: Response body not found")
        return
    }

    respBody, _ := ioutil.ReadAll(body)
    log.Printf("HTTP PostForm/GET returned %v", string(respBody))


	//15d3d8e8de90ebcc
	// for _, thread := range threads {

	// }
}