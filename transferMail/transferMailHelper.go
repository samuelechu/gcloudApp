package transferMail

import (
    "log"
	"net/http"
	"strings"
    "io/ioutil"
    "time"
    "github.com/samuelechu/cloudSQL"
    "github.com/buger/jsonparser"
)

type nopCloser struct { 
    io.Reader 
} 

func (nopCloser) Close() os.Error { return nil } 

func startTransfer(curUserID, sourceToken, sourceID, destToken, destID string) {
	threads := cloudSQL.GetThreadsForUser(curUserID)
	log.Printf("GetThreads returned %v", threads)
	log.Printf("curUserID : %v, sourceToken : %v, sourceID : %v, destToken : %v, destID : %v", curUserID, sourceToken, sourceID, destToken, destID)
	
	urlStr := "https://www.googleapis.com/gmail/v1/users/me/messages/15d3d8e8de90ebcc?format=raw" //testTransfer label
    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels"
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    transport := http.Transport{}

    client := &http.Client{
        Transport: &transport,
        Timeout: time.Second * 10,
    }     

    resp, err := client.Do(req)

    if err != nil {
    		log.Printf("Error: %v", err)
            return
    }
    
    body := resp.Body
    defer body.Close()

    if body == nil {
    	log.Print("Error: Response body not found")
        return
    }

    respBody, _ := ioutil.ReadAll(body)
    //log.Printf("HTTP PostForm/GET returned %v", string(respBody))

    raw, _, _, _ := jsonparser.Get(respBody, "raw")
	//log.Printf("Got labels: %v", string(labels))


    // jsonparser.ArrayEach(respBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
    //     label, _, _, _ := jsonparser.Get(value, "")
    //     if string(label) != "" {
    //         log.Printf("Got label: %v", string(label))

    //     }
        
    // }, "labelIds")

	//15d3d8e8de90ebcc
	// for _, thread := range threads {

	// }

    urlStr = "https://www.googleapis.com/upload/gmail/v1/users/me/messages?uploadType=multipart"

    body = nopCloser{bytes.NewBufferString("--foobar\nContent-Type: application/json; charset=UTF-8\n{" +
"\n\"raw\":\"" + string(raw) + "\"\n\"labelIds\": [\"INBOX\", \"UNREAD\"]\n}" +
"--foo_bar\nContent-Type: message/rfc822\n\nstringd\n--foo_bar--")} 


    insertReq, _ := http.NewRequest("POST", urlStr, body)
    insertReq.Header.Set("Authorization", "Bearer " + destToken)
    insertReq.Header.Set("Content-Type", "multipart/related; boundary=foo_bar")

    resp, err = client.Do(insertReq)

    if err != nil {
    		log.Printf("Error: %v", err)
            return
    }
    
    body = resp.Body
    defer body.Close()

    if body == nil {
    	log.Print("Error: Response body not found")
        return
    }

    respBody, _ = ioutil.ReadAll(body)
    log.Printf("HTTP PostForm/GET returned %v", string(respBody))


}