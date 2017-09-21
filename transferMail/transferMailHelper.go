package transferMail

import (
    "log"
	"io"
    //"bytes"
    "github.com/samuelechu/cloudSQL"
    "golang.org/x/net/context"
    "google.golang.org/appengine/urlfetch"
)

type nopCloser struct { 
    io.Reader 
} 

func (nopCloser) Close() error { return nil } 

func startTransfer(ctx context.Context, selectedLabels []string, curUserID, sourceToken, sourceID, destToken, destID string) {

    client := urlfetch.Client(ctx)
    labelMap := GetLabels(client, sourceToken)

    for _, val := range selectedLabels {
        labelId := labelMap[val]
        addThreadsWithLabel(client, curUserID, labelId, sourceToken)
    }

    cloudSQL.UpdateThreadInfoForJob(curUserID)

    //get threads
	sourceThreads := cloudSQL.GetThreadsForUser(curUserID)
	//log.Printf("GetThreads returned %v", sourceThreads)
	log.Printf("curUserID : %v, sourceToken : %v, sourceID : %v, destToken : %v, destID : %v", curUserID, sourceToken, sourceID, destToken, destID)

    insertThreads(ctx, sourceThreads,sourceToken,destToken,curUserID)
}



