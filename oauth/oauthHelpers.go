package oauth

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"log"
	"os"
	"io/ioutil"
    "encoding/json"
	"net/http"
	"net/url"
)

type idTokenRespBody struct{
    Aud     string
    Sub     string
    Name	string
}

type accessTokenRespBody struct{
    Access_token    string
    Expires_in      float64
    Token_type      string
}

//response after user grants permissions
type oauthRespBody struct{
	Access_token    string
    Expires_in      float64
    Token_type      string
    Refresh_token 	string
    Id_token 		string
}

func getJSONRespBody(w http.ResponseWriter, r *http.Request, url string, data url.Values, rbType interface{}) interface{} {

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
		case idTokenRespBody:
			rb = rbType.(idTokenRespBody)
			json.Unmarshal(respBody, &rb)
			return rb

		case accessTokenRespBody:
			rb = rbType.(accessTokenRespBody)
			json.Unmarshal(respBody, &rb)
			return rb

		case oauthRespBody:
			rb = rbType.(oauthRespBody)
			json.Unmarshal(respBody, &rb)
			return rb
		
		default:
			return rbType


	} 
}

//verifies that the id_token that identifies user is genuine
func verifyIDToken(w http.ResponseWriter, r *http.Request, token string) string, string {

    urlStr := "https://www.googleapis.com/oauth2/v3/tokeninfo"

    bodyVals := url.Values{
        "id_token": {token},
    }

    var respBody idTokenRespBody
    if rb, ok := getJSONRespBody(w, r, urlStr, bodyVals, respBody).(idTokenRespBody); ok {

        if rb.Aud == os.Getenv("CLIENT_ID") {
            return rb.Sub, rb.Name
        }
    } else {
        http.Error(w, "Error: incorrect responsebody", 400)
    }

    return "",""
}