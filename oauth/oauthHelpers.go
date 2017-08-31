package oauth

import (
	"os"
	"net/http"
	"net/url"
	"github.com/samuelechu/jsonHelper"
)

//verifies that the id_token that identifies user is genuine
func verifyIDToken(w http.ResponseWriter, r *http.Request, token string) (string, string) {

    urlStr := "https://www.googleapis.com/oauth2/v3/tokeninfo"

    bodyVals := url.Values{
        "id_token": {token},
    }

    var respBody IdTokenRespBody
    if rb, ok := jsonHelper.GetJSONRespBody(w, r, urlStr, bodyVals, respBody).(IdTokenRespBody); ok {

        if rb.Aud == os.Getenv("CLIENT_ID") {
            return rb.Sub, rb.Name
        }
    } else {
        http.Error(w, "Error: incorrect responsebody", 400)
    }

    return "",""
}