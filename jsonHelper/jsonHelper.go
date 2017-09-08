package jsonHelper

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"log"
	"io"
	"io/ioutil"
    "encoding/json"
	"net/http"
	"net/url"
)

func GetJSONRespBodyGET(w http.ResponseWriter, r *http.Request, url string, rbType interface{}) interface{} {

    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)

    resp, err := client.Get(url)

    if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return nil
    }

    return UnmarshalJSON(w, r, resp.Body, rbType)
}

func GetJSONRespBody(w http.ResponseWriter, r *http.Request, url string, data url.Values, rbType interface{}) interface{} {

    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)

    resp, err := client.PostForm(url, data)

    if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return nil
    }

    return UnmarshalJSON(w, r, resp.Body, rbType)
}

func UnmarshalJSON(w http.ResponseWriter, r *http.Request, body io.ReadCloser, struct_type interface{}) interface{} {

	defer body.Close()

	if body == nil {
        http.Error(w, "Response body not found", 400)
        return nil
    }

    respBody, _ := ioutil.ReadAll(body)
    log.Printf("HTTP PostForm/GET returned %v", string(respBody))
    // fmt.Fprintf(w, "HTTP Post returned %v", string(respBody))

    switch values := struct_type.(type) {
		case IdTokenRespBody:
			values = struct_type.(IdTokenRespBody)
			json.Unmarshal(respBody, &values)
			return values

		case AccessTokenRespBody:
			values = struct_type.(AccessTokenRespBody)
			json.Unmarshal(respBody, &values)
			return values

		case OauthRespBody:
			values = struct_type.(OauthRespBody)
			json.Unmarshal(respBody, &values)
			return values

		case UserInfoRespBody:
			values = struct_type.(UserInfoRespBody)
			json.Unmarshal(respBody, &values)
			return values

		case User:
			values = struct_type.(User)
			json.Unmarshal(respBody, &values)
			return values
		
		default:
			return struct_type
	} 
}

//writes a json object to the responseWriter
func MarshalJSON(w http.ResponseWriter, r *http.Request, json_object interface{}){
	data, err := json.Marshal(json_object)
	
	if err != nil {
        http.Error(w, "Error marshaling json object", 500)
        return
    }
	w.Write(data)

}