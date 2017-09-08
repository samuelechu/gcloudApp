package transferMail

import (
	"log"
	"net/http"
	"github.com/samuelechu/oauth"
)

func init() {
     http.HandleFunc("/transferStart", transferEmail)
}

func transferEmail(w http.ResponseWriter, r *http.Request) {
	var sourceToken, sourceID, destToken, destID string
    
    sourceCookie, err := r.Cookie("source")
    if err == nil {
        sourceToken = sourceCookie.Value
    }

    destCookie, err := r.Cookie("destination")
    if err == nil {
        destToken = destCookie.Value
    }

    sourceID, _ = oauth.GetUserInfo(w, r, sourceToken)
    destID, _ = oauth.GetUserInfo(w, r, destToken)

    log.Printf("Source ID: %v\n", sourceID)
    log.Printf("Dest ID: %v\n", destID)

    


}