package oauth

import (
	"os"
    "log"
	"net/http"
	"net/url"
	"github.com/samuelechu/jsonHelper"
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

    redirectString := "https://gotesting-175718.appspot.com"
    if appengine.IsDevAppServer(){
        redirectString = "https://8080-dot-2979131-dot-devshell.appspot.com"
    }
    http.Redirect(w, r, redirectString, 302)
}