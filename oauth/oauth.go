package oauth

import (
	"google.golang.org/appengine"

	"fmt"
	"log"
	"os"
	"net/http"
    "net/http/cookiejar"
	"net/url"
    "github.com/samuelechu/cloudSQL"
    "github.com/samuelechu/jsonHelper"
)

var cookieJar *cookiejar.Jar

func init() {
     http.HandleFunc("/askPermissions", askPermissions)
     http.HandleFunc("/oauthCallback", oauthCallback)
     http.HandleFunc("/testrefToken", getAccessToken)
     http.HandleFunc("/deleteCookies", deleteCookies)

    log.Print("Initializing curCookies")
    curCookies = &cookies{}
    cookieJar, _ = cookiejar.New(nil)

    log.Printf("Curcookies: %v", curCookies)
}

func getAccessToken(w http.ResponseWriter, r *http.Request) {
    
    urlStr := "https://www.googleapis.com/oauth2/v4/token"
 
    bodyVals := url.Values{
        "client_id": {os.Getenv("CLIENT_ID")},
        "client_secret": {os.Getenv("CLIENT_SECRET")},
        "refresh_token":{"1/pI4NYPkOnY_73TvjIPvZZ8jy9x7sqgmltw43cQDc-4g"},
        "grant_type": {"refresh_token"},
    }

    var respBody jsonHelper.AccessTokenRespBody 
    if rb, ok := jsonHelper.GetJSONRespBody(w, r, urlStr, bodyVals, respBody).(jsonHelper.AccessTokenRespBody); ok {
        fmt.Fprintf(w, "HTTP Post returned %v %v %v", rb.Access_token, rb.Expires_in, rb.Token_type)

    }
}

var accountType string

//askPermissions from user, response is auth code
func askPermissions(w http.ResponseWriter, r *http.Request) {
    //request will be format :   /askPermissions?(source||destination)
    accountType = r.URL.Query().Get("type")
    permissions := ""

    log.Printf("In askPermissions: account type was %v", accountType)
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
        //"prompt" : {"consent"},
        "state" : {"state_parameter_passthrough_value"},
        "redirect_uri" : {redirectUri},
        "response_type" : {"code"},
        "client_id" : {os.Getenv("CLIENT_ID")},
    }

    queryString := queryVals.Encode()

    redirectString := "https://accounts.google.com/o/oauth2/v2/auth?" + queryString

    //exchange auth code for access/refresh token in oauthCallback
    http.Redirect(w, r, redirectString, 302)

    log.Print("woohooooooo redirect")
}

//exchange auth code for access token
func oauthCallback(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL.Query())
    log.Print(r.URL.Query().Get("code"))

    authCode := r.URL.Query().Get("code")
    
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


    // if uid != "" {
    //     fmt.Fprintf(w, "\n Token verified! Name: %v, UserId: %v, Refresh_token: %v, Access_token: %v",
    //                     name, uid, respBody.Refresh_token, respBody.Access_token)
    // } else {
    //     fmt.Fprint(w, "\n Token verification failed!")
    // }

    //store the user and refresh token into database
    cloudSQL.InsertUser(uid, name, respBody.Refresh_token)
    
    redirectString := "https://gotesting-175718.appspot.com"
    if appengine.IsDevAppServer(){
        redirectString = "https://8080-dot-2979131-dot-devshell.appspot.com"
    }

    log.Printf("In oauth Callback. The type is %v, id token is \n%v", accountType, respBody.Id_token)

    var cookies []*http.Cookie

    cookie := &http.Cookie{
        Name: accountType,
        Value: respBody.Id_token,
        //Path: "/",
        //Domain: "8080-dot-2979131-dot-devshell.appspot.com",
    }

    cookies = append(cookies, cookie)
    
    //setCookies(accountType, respBody.Id_token)
    u, _ := url.Parse("https://8080-dot-2979131-dot-devshell.appspot.com")
    cookieJar.SetCookies(u, cookies)

    log.Println(cookieJar.Cookies(u))
    log.Print(r)
    log.Print(r.Host)
    log.Printf("%v", r.url)
    accountType = ""
    http.Redirect(w, r, redirectString, 302)
}