package oauth

import (
	"os"
    "log"
	"net/http"
	"net/url"
	"github.com/samuelechu/jsonHelper"
   // "github.com/samuelechu/cloudSQL"
)

//verifies that the id_token that identifies user is genuine
func VerifyIDToken(w http.ResponseWriter, r *http.Request, token string) (string, string) {

    urlStr := "https://www.googleapis.com/oauth2/v3/tokeninfo"

    bodyVals := url.Values{
        "id_token": {token},
    }

    var respBody jsonHelper.IdTokenRespBody
    if rb, ok := jsonHelper.GetJSONRespBody(w, r, urlStr, bodyVals, respBody).(jsonHelper.IdTokenRespBody); ok {

        if rb.Aud == os.Getenv("CLIENT_ID") {
            return rb.Sub, rb.Name
        } else {
        	return "",""
        }

    } else {
        http.Error(w, "Error: incorrect responsebody", 400)
    }

    return "",""   
}

// func GetAccessToken(uid string) string{

//     refreshToken, err := cloudSQL.GetRefreshToken(uid)

//     if err != nil {
//         log.Print("DB err: %v", err)
//         return ""
//     }
    
//     urlStr := "https://www.googleapis.com/oauth2/v4/token"
 
//     bodyVals := url.Values{
//         "client_id": {os.Getenv("CLIENT_ID")},
//         "client_secret": {os.Getenv("CLIENT_SECRET")},
//         "refresh_token":{refreshToken},
//         "grant_type": {"refresh_token"},
//     }

//     var respBody jsonHelper.AccessTokenRespBody 
//     if rb, ok := jsonHelper.GetJSONRespBody(w, r, urlStr, bodyVals, respBody).(jsonHelper.AccessTokenRespBody); ok {
//         fmt.Fprintf(w, "HTTP Post returned %v %v %v", rb.Access_token, rb.Expires_in, rb.Token_type)

//     }
// }

func deleteCookies(w http.ResponseWriter, r *http.Request) {

	sourceCookie, err := r.Cookie("source")
    if err == nil {
        log.Print("deleting source cookie")
        sourceCookie.MaxAge = -1
        http.SetCookie(w, sourceCookie)
    }

    destinationCookie, err := r.Cookie("destination")
    if err == nil {
        log.Print("deleting dest cookie")
        destinationCookie.MaxAge = -1
        http.SetCookie(w, destinationCookie)
    }
}