package oauth

import (
	"os"
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
	sourceCookie *Cookie
	destCookie *Cookie
}

func GetCookies(w http.ResponseWriter, r *http.Request) {
    deleteCookies(w, r)

    http.SetCookie(w, curCookies.sourceCookie)
    http.SetCookie(w, curCookies.destCookie)
}

func setCookies(cookieStruct *cookies, accountType string, id_token string){
	switch accountType {
		case "source":
			cookieStruct.sourceCookie = &http.Cookie{
		        Name: "source"
		        Value: id_token,
			}

		case "destination":
			cookieStruct.destCookie = &http.Cookie{
		        Name: "destination"
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
}