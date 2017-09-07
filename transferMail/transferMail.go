package transferMail

import (
	"log"
	"net/http"
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

    log.Printf("Source ID: %v\n", sourceName)
    log.Printf("Dest ID: %v\n", destName)
}