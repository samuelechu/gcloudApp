package main
import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/drive/v2"
    "google.golang.org/api/plus/v1"
    "google.golang.org/appengine"
    "google.golang.org/appengine/log"
    //"html/template"
    "os"
    "net/http"
)

const htmlIndex = `<html><body>
<a href="/GoogleLogin">Log in with Google</a>
</body></html>
`
var conf = &oauth2.Config{
    ClientID:     os.Getenv("CLIENT_ID"),       // Replace with correct ClientID
    ClientSecret: os.Getenv("CLIENT_SECRET"),   // Replace with correct ClientSecret
    RedirectURL:  "https://gotesting-175718.appspot.com/googleCallback",
    Scopes: []string{
        "https://www.googleapis.com/auth/drive",
        "profile",
    },
    Endpoint: google.Endpoint,
}

func init() {
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/authorize", handleAuthorize)
    http.HandleFunc("/oauth2callback", handleOAuth2Callback)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, htmlIndex)
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {

    if appengine.IsDevAppServer(){
        conf.RedirectURL = "https://8080-dot-2979131-dot-devshell.appspot.com/googleCallback"
    }

    c := appengine.NewContext(r)
    url := conf.AuthCodeURL("")
    http.Redirect(w, r, url, http.StatusFound)
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    code := r.FormValue("code")
    tok, err := conf.Exchange(c, code)
    if err != nil {
        log.Errorf(c, "%v", err)
    }
    client := conf.Client(c, tok)

    // PLUS SERVICE CLIENT
    pc, err := plus.New(client)
    if err != nil {
        log.Errorf(c, "An error occurred creating Plus client: %v", err)
    }
    person, err := pc.People.Get("me").Do()
    if err != nil {
        log.Errorf(c, "Person Error: %v", err)
    }
    log.Infof(c, "Name: %v", person.DisplayName)

    // DRIVE CLIENT
    dc, err := drive.New(client)
    if err != nil {
        log.Errorf(c, "An error occurred creating Drive client: %v", err)
    }
    files, err := dc.Files.List().Do()
    for _, value := range files.Items {
        log.Infof(c, "Files: %v", value.Title)
    }
}

/*
import (
    "fmt"
    "net/http"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/appengine"
    "os"
    "log"
    "google.golang.org/api/drive/v2"
    "golang.org/x/net/context"
)

var (
    googleOauthConfig = &oauth2.Config{
        RedirectURL:    "https://gotesting-175718.appspot.com/googleCallback",
        ClientID:     os.Getenv("CLIENT_ID"), // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
        ClientSecret: os.Getenv("CLIENT_SECRET"), // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
        Scopes:       []string{"https://www.googleapis.com/auth/drive", "https://www.googleapis.com/auth/drive.file", "https://www.googleapis.com/auth/gmail.readonly"},
        Endpoint:     google.Endpoint,
    }
// Some random string, random for each request
    oauthStateString = "random"
)

const htmlIndex = `<html><body>
<a href="/GoogleLogin">Log in with Google</a>
</body></html>
`

func main() {
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/GoogleLogin", handleGoogleLogin)
    http.HandleFunc("/googleCallback", handleGoogleCallback)

    log.Print("Listening on port 8080")
    http.ListenAndServe(":8080", nil)
    appengine.Main()
}

func handleMain(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, htmlIndex)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
    if appengine.IsDevAppServer(){
        googleOauthConfig.RedirectURL = "https://8080-dot-2979131-dot-devshell.appspot.com/googleCallback"
    }

    url := googleOauthConfig.AuthCodeURL(oauthStateString)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
    state := r.FormValue("state")
    if state != oauthStateString {
        fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    code := r.FormValue("code")

    token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
    if err != nil {
        log.Print("oauthConf.Exchange() failed with '%s'\n", err)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }
    client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

    //response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
    driveService, err := drive.New(client)

    myR, err := driveService.Files.List().MaxResults(10).Do()
    if err != nil {
        fmt.Fprintf(w, "Couldn't retrieve files ", err)
    }
    if len(myR.Items) > 0 {
        for _, i := range myR.Items {
            fmt.Fprintf(w, i.Title, " ", i.Id)
        }
    } else {
        fmt.Fprintf(w, "No files found.")
    }

    //defer response.Body.Close()
    //contents, err := ioutil.ReadAll(response.Body)
    //fmt.Fprintf(w, "Content: %s\n", contents)
}
*/