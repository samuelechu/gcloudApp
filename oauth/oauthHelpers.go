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

type cookies struct{
	sourceCookie *http.Cookie
	destCookie *http.Cookie
}

var curCookies *cookies

func GetCookies(w http.ResponseWriter, r *http.Request) {
    deleteCookies(w, r)

    if curCookies == nil {
        log.Print("Initializing curCookies")
        curCookies = &cookies{}
        return
    }

    if curCookies.sourceCookie != nil {
        http.SetCookie(w, curCookies.sourceCookie)
    }
    
    if curCookies.destCookie != nil {
        http.SetCookie(w, curCookies.destCookie)
    }

}

func setCookies(accType string, id_token string){
	switch accType {
		case "source":
			curCookies.sourceCookie = &http.Cookie{
		        Name: "source",
		        Value: id_token,
			}

		case "destination":
			curCookies.destCookie = &http.Cookie{
		        Name: "destination",
		        Value: id_token,
			}
	}
}

func deleteCookies(w http.ResponseWriter, r *http.Request) {

	sourceCookie, err := r.Cookie("source")
    if err == nil {
        sourceCookie.MaxAge = -1
        http.SetCookie(w, sourceCookie)
    }


    destinationCookie, err := r.Cookie("destination")
    if err == nil {
        destinationCookie.MaxAge = -1
        http.SetCookie(w, destinationCookie)
    }

    if r.URL.Query().Get("resetStruct") == "true" {
        curCookies = &cookies{}
    }
}