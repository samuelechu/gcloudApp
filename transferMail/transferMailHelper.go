package transferMail

import (
    "log"
	//"net/http"
	"io"
    //"bytes"
    "github.com/samuelechu/cloudSQL"
   // "github.com/buger/jsonparser"
    "golang.org/x/net/context"
    "google.golang.org/appengine/urlfetch"
    //"github.com/samuelechu/jsonHelper"
)

type nopCloser struct { 
    io.Reader 
} 

func (nopCloser) Close() error { return nil } 

func startTransfer(ctx context.Context, curUserID, sourceToken, sourceID, destToken, destID string) {
    client := urlfetch.Client(ctx)

//get threads
	sourceThreads := cloudSQL.GetThreadsForUser(curUserID)
	log.Printf("GetThreads returned %v", sourceThreads)
	log.Printf("curUserID : %v, sourceToken : %v, sourceID : %v, destToken : %v, destID : %v", curUserID, sourceToken, sourceID, destToken, destID)
	
	//15d3d8e8de90ebcc

    // fields := [][]string{
    //     []string{"name"},
    //     []string{"messageListVisibility"},
    //     []string{"labelListVisibility"},
    //     []string{"type"},
    // }

    // var name, messageListVisibility, labelListVisibility string// labelType string
    // jsonparser.EachKey(value, func(idx int, value []byte, vt jsonparser.ValueType, err error){
    //     switch idx {
    //         case 0:
    //             name, _ = jsonparser.ParseString(value)
    //         case 1:
    //             messageListVisibility, _ = jsonparser.ParseString(value)
    //         case 2:
    //             labelListVisibility, _ = jsonparser.ParseString(value)
    //         // case 3:
    //      //     labelType, _ = jsonparser.ParseString(value)
    //     }
    // }, fields...)
//     HTTP PostForm/GET returned {
//  "id": "15e827062708e520",
//  "threadId": "15e827062708e520",
//  "labelIds": [
//   "UNREAD",
//   "INBOX"
//  ]
// }


    insertThreads(client,sourceThreads,sourceToken,destToken)
}