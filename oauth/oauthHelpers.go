package oauth

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"log"
	"bytes"
	"io/ioutil"
    "encoding/json"
	"net/http"
	"net/url"
)


func getJSONRespBody(w http.ResponseWriter, r *http.Request, url string, data url.Values, rbType interface{}) interface{} {

	body := bytes.NewBufferString(data.Encode())

    log.Print(body)
    req, _ := http.NewRequest("POST", url, body)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    log.Print("finished marshaling")

    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)

    resp, err := client.Do(req)
    if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return nil
    }
    log.Print("am here")
    defer resp.Body.Close()
    respBody, _ := ioutil.ReadAll(resp.Body)
    log.Printf("HTTP Post returned %v", string(respBody))
    // fmt.Fprintf(w, "HTTP Post returned %v", string(respBody))

    if resp.Body == nil {
        http.Error(w, "Response body not found", 400)
        return nil
    }

    switch _ := rbType.(type) {
		case idTokenRespBody:
			rb := rbType.(idTokenRespBody)
			json.Unmarshal(respBody, &rb)
			return rb

		case accessTokenRespBody:
			rb := rbType.(accessTokenRespBody)
			json.Unmarshal(respBody, &rb)
			return rb
		
		default:
			return rbType


	}
	
 
    
    
}