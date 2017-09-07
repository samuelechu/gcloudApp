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
	sourceToken := ""
    destToken := ""
    
    sourceCookie, err := r.Cookie("source")
    if err == nil {
        sourceToken = sourceCookie.Value
    }

    destCookie, err := r.Cookie("destination")
    if err == nil {
        destToken = destCookie.Value
    }

    log.Printf("Source Cookie: %v\n", sourceCookie)
    log.Printf("Dest Cookie: %v\n", destCookie)

    sourceID, _ := oauth.VerifyIDToken(w, r, sourceToken)
    destID, _ := oauth.VerifyIDToken(w, r, destToken)

    log.Printf("Source ID: %v\n", sourceID)
    log.Printf("Dest ID: %v\n", destID)


    sAccess := oauth.GetAccessToken(w, r, sourceID)
    dAccess := oauth.GetAccessToken(w, r, destID)

    log.Printf("Source ID: %v\n", sAccess)
    log.Printf("Dest ID: %v\n", sAccess)

}