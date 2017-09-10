package transferMail

import (
	"log"
	"net/http"
    "io/ioutil"
    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
	"github.com/samuelechu/oauth"
    "github.com/samuelechu/cloudSQL"
    "github.com/buger/jsonparser"
)

func startTransfer(curUserID, sourceToken, sourceID, destToken, destID string) {
	threads := cloudSQL.GetThreadsForUser(curUserID)
	log.Printf("GetThreads returned %v", threads)
}