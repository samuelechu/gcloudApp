package oauth

import (
	"google.golang.org/appengine"
	"fmt"
	"log"
	"os"
	"net/http"
	"net/url"
    "github.com/samuelechu/cloudSQL"
    "github.com/samuelechu/jsonHelper"
)

func init() {
     http.HandleFunc("/askPermissions", askPermissions)
     http.HandleFunc("/oauthCallback", oauthCallback)
     http.HandleFunc("/deleteCookies", deleteCookies)
     http.HandleFunc("/getAccessToken", getAccessToken)
}

func getAccessToken(w http.ResponseWriter, r *http.Request) {
    
    uid := r.URL.Query().Get("uid")
    refreshToken, err := cloudSQL.GetRefreshToken(uid)


    if err != nil {
        http.Error(w, fmt.Sprintf("DB err: %v.", err), 500)
        return
    }
    
    urlStr := "https://www.googleapis.com/oauth2/v4/token"
 
    bodyVals := url.Values{
        "client_id": {os.Getenv("CLIENT_ID")},
        "client_secret": {os.Getenv("CLIENT_SECRET")},
        "refresh_token":{refresh_token},
        "grant_type": {"refresh_token"},
    }

    var respBody jsonHelper.AccessTokenRespBody 
    if rb, ok := jsonHelper.GetJSONRespBody(w, r, urlStr, bodyVals, respBody).(jsonHelper.AccessTokenRespBody); ok {
        fmt.Fprintf(w, "HTTP Post returned %v %v %v", rb.Access_token, rb.Expires_in, rb.Token_type)

    }
}

//askPermissions from user, response is auth code
func askPermissions(w http.ResponseWriter, r *http.Request) {
    //request will be format :   /askPermissions?(source||destination)
    accountType := r.URL.Query().Get("type")
    permissions := ""

    //pass on account type to redirect
    http.SetCookie(w, &http.Cookie{
        Name: "account_type",
        Value: accountType,
    })

    switch accountType {
        case "source":
            permissions = "https://www.googleapis.com/auth/gmail.readonly"
        case "destination":
            permissions = "https://www.googleapis.com/auth/gmail.insert"
        default:
            http.Error(w, "must specify in queryString type : source || destination", 400)
            return
    }

    redirectUri := "https://gotesting-175718.appspot.com/oauthCallback"
	if appengine.IsDevAppServer(){
    	redirectUri = "https://8080-dot-2979131-dot-devshell.appspot.com/oauthCallback"
	}

    queryVals := url.Values{
        "scope" : {"profile " + permissions},
        "access_type" : {"offline"},
        "include_granted_scopes" : {"true"},
        "prompt" : {"select_account"},
        "state" : {"state_parameter_passthrough_value"},
        "redirect_uri" : {redirectUri},
        "response_type" : {"code"},
        "client_id" : {os.Getenv("CLIENT_ID")},
    }

    queryString := queryVals.Encode()

    redirectString := "https://accounts.google.com/o/oauth2/v2/auth?" + queryString

    //exchange auth code for access/refresh token in oauthCallback
    http.Redirect(w, r, redirectString, 302)
}

//exchange auth code for access token
func oauthCallback(w http.ResponseWriter, r *http.Request) {
    authCode := r.URL.Query().Get("code")

    //retrieve account type and reset cookie
    accountType := ""
    typeCookie, err := r.Cookie("account_type")
    if err == nil {
        accountType = typeCookie.Value
        typeCookie.MaxAge = -1
        http.SetCookie(w, typeCookie)
    }
    
    urlStr := "https://www.googleapis.com/oauth2/v4/token"

    redirectUri := "https://gotesting-175718.appspot.com/oauthCallback"
    if appengine.IsDevAppServer(){
        redirectUri = "https://8080-dot-2979131-dot-devshell.appspot.com/oauthCallback"
    }

    bodyVals := url.Values{
        "code": {authCode},
        "client_id": {os.Getenv("CLIENT_ID")},
        "client_secret": {os.Getenv("CLIENT_SECRET")},
        "redirect_uri": {redirectUri},
        "grant_type": {"authorization_code"},
    }

    //  OauthRespBody struct: Access_token, Expires_in, Token_type, Refresh_token, Id_token string
    var respBody jsonHelper.OauthRespBody
    if rb, ok := jsonHelper.GetJSONRespBody(w, r, urlStr, bodyVals, respBody).(jsonHelper.OauthRespBody); ok {
        respBody = rb
        //fmt.Fprintf(w, "HTTP Post returned %+v", rb)
    }

    //verify the signed in user
    uid, name := VerifyIDToken(w, r, respBody.Id_token)

    if uid == "" {
        http.Error(w, "Error with token", 500)
    }

    if uid != "" {
        log.Printf("\n Token verified! Name: %v, UserId: %v, Refresh_token: %v, Access_token: %v",
                        name, uid, respBody.Refresh_token, respBody.Access_token)
    } else {
        log.Print("\n Token verification failed!")
    }

    //store the user and refresh token into database
    cloudSQL.InsertUser(uid, name, respBody.Refresh_token)

    //send id_token to browser to identify the signed in user 
    http.SetCookie(w, &http.Cookie{
        Name: accountType,
        Value: respBody.Id_token,
    })

    redirectString := "https://gotesting-175718.appspot.com"
    if appengine.IsDevAppServer(){
        redirectString = "https://8080-dot-2979131-dot-devshell.appspot.com"
    }
    http.Redirect(w, r, redirectString, 302)
}