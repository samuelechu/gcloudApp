package jsonHelper

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"log"
	"io/ioutil"
    "encoding/json"
	"net/http"
	"net/url"
)

func GetJSONRespBody(w http.ResponseWriter, r *http.Request, url string, data url.Values, rbType interface{}) interface{} {

    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)

    resp, err := client.PostForm(url, data)

    if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return nil
    }

    log.Print("am here")
    defer resp.Body.Close()
    respBody, _ := ioutil.ReadAll(resp.Body)
    log.Printf("HTTP PostForm returned %v", string(respBody))
    // fmt.Fprintf(w, "HTTP Post returned %v", string(respBody))

    if resp.Body == nil {
        http.Error(w, "Response body not found", 400)
        return nil
    }

    switch rb := rbType.(type) {
		case IdTokenRespBody:
			rb = rbType.(IdTokenRespBody)
			json.Unmarshal(respBody, &rb)
			return rb

		case AccessTokenRespBody:
			rb = rbType.(AccessTokenRespBody)
			json.Unmarshal(respBody, &rb)
			return rb

		case OauthRespBody:
			rb = rbType.(OauthRespBody)
			json.Unmarshal(respBody, &rb)
			return rb
		
		default:
			return rbType


	} 
}