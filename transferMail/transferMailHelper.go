package transferMail

import (
    "log"
	"net/http"
    "io/ioutil"
    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
    "github.com/samuelechu/cloudSQL"
    //"github.com/buger/jsonparser"
)

func startTransfer(w http.ResponseWriter, client *http.Client, curUserID, sourceToken, sourceID, destToken, destID string) {
	threads := cloudSQL.GetThreadsForUser(curUserID)
	log.Printf("GetThreads returned %v", threads)

	
	urlStr := "https://www.googleapis.com/gmail/v1/users/me/messages/15d3d8e8de90ebcc" //testTransfer label
    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels"
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

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


	//15d3d8e8de90ebcc
	// for _, thread := range threads {

	// }
}