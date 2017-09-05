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

var curCookies *cookies //initialized in init

func GetCookies(w http.ResponseWriter, r *http.Request) {

    u, _ := url.Parse("https://8080-dot-2979131-dot-devshell.appspot.com")
    //cookieJar.SetCookies(u, cookies)

    log.Println(cookieJar.Cookies(u))
    // deleteCookies(w, r)

    // if curCookies.sourceCookie != nil {
    //     log.Printf("setting the sourceCookie: %v", curCookies.sourceCookie)
    //     http.SetCookie(w, curCookies.sourceCookie)
    // }
    
    // if curCookies.destCookie != nil {
    //     log.Printf("setting the destCookie: %v", curCookies.destCookie)
    //     http.SetCookie(w, curCookies.destCookie)
    // }

}

func setCookies(accType string, id_token string){
    log.Printf("In setCookies. The type is %v, id token is \n%v", accType, id_token)

    log.Printf("Curcookies: %v", curCookies)

	switch accType {
		case "source":
            log.Print("setting source cookie")
			curCookies.sourceCookie = &http.Cookie{
		        Name: "source",
		        Value: id_token,
			}

		case "destination":
            log.Print("setting dest cookie")
			curCookies.destCookie = &http.Cookie{
		        Name: "destination",
		        Value: id_token,
			}
	}

    log.Printf("Curcookies: %v", curCookies)
}

func deleteCookies(w http.ResponseWriter, r *http.Request) {


    log.Printf("In deleteCookies. Curcookies: %v", curCookies)

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

    if r.URL.Query().Get("resetStruct") == "true" {
        log.Print("resetting struct")
        curCookies = &cookies{}
    }

    log.Printf("end of deleteCookies. Curcookies: %v", curCookies)
}