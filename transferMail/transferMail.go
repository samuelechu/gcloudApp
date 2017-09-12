package transferMail

import (
	"log"
	"net/http"
    "io/ioutil"
   // "time"

    "golang.org/x/net/context"
    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
    "google.golang.org/appengine/runtime"
	"github.com/samuelechu/oauth"
    "github.com/samuelechu/cloudSQL"
    "github.com/buger/jsonparser"
)

type Values struct {
    m map[string]string
}

func (v Values) Get(key string) string {
    return v.m[key]
}

func init() {
     http.HandleFunc("/transferStart", transferEmail)
}

func transferEmail(w http.ResponseWriter, r *http.Request) {
	var curUserID, sourceToken, sourceID, destToken, destID string

    curUserCookie, err := r.Cookie("current_user")
    if err == nil {
        curUserID = curUserCookie.Value
    }
    
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

    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/messages/15e5d6ed5bb68a29?format=raw"
//retrieve threads

    urlStr := "https://www.googleapis.com/gmail/v1/users/me/threads?labelIds=Label_8" //testTransfer label
    //urlStr := "https://www.googleapis.com/gmail/v1/users/me/labels"
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("Authorization", "Bearer " + sourceToken)

    cookieInfo := Values{map[string]string{
        "curUserID": curUserID,
        "sourceToken": sourceToken,
        "sourceID": sourceID,
        "destToken": destToken,
        "destID": destID,
    }}

    c := appengine.NewContext(r)
    ctx := context.WithValue(c, "cookieInfo", cookieInfo)

    client := urlfetch.Client(ctx)

    resp, err := client.Do(req)

    if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
    }


    body := resp.Body
    defer body.Close()

    if body == nil {
        http.Error(w, "Response body not found", 400)
        return
    }

    respBody, _ := ioutil.ReadAll(body)
    log.Printf("HTTP PostForm/GET returned %v", string(respBody))

    // if message_id, ok := jsonparser.GetString(respBody, "id"); ok == nil{
    //     log.Printf("ID of messsage was %v", message_id)
    // }
    
    jsonparser.ArrayEach(respBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        thread_id, _, _, _ := jsonparser.Get(value, "id")
        if string(thread_id) != "" {
            log.Printf("Inserting into database: Thread %v", string(thread_id))
            cloudSQL.InsertThread(curUserID, string(thread_id))

        }
        
    }, "threads")


    res, _, _, _ := jsonparser.Get(respBody, "resultSizeEstimate")
    log.Printf("jsonparser returned %v", string(res))

    //go startTransfer(curUserID, sourceToken, sourceID, destToken, destID)
    // urlStr := "https://www.googleapis.com/upload/gmail/v1/users/me/messages?uploadType=media"

    // bodyVals := url.Values{
    //     "raw": {"RGVsaXZlcmVkLVRvOiBtaWNoYWVsa2VyZEBnbWFpbC5jb20NClJlY2VpdmVkOiBieSAxMC4zMS4xNzMuNSB3aXRoIFNNVFAgaWQgdzVjc3A0MjIxMzl2a2U7DQogICAgICAgIFRodSwgNyBTZXAgMjAxNyAxMDo0MTozOCAtMDcwMCAoUERUKQ0KWC1SZWNlaXZlZDogYnkgMTAuOTkuOTYuMjMgd2l0aCBTTVRQIGlkIHUyM21yODYyNjJwZ2IuMjU3LjE1MDQ4MDYwOTc5OTk7DQogICAgICAgIFRodSwgMDcgU2VwIDIwMTcgMTA6NDE6MzcgLTA3MDAgKFBEVCkNCkFSQy1TZWFsOiBpPTE7IGE9cnNhLXNoYTI1NjsgdD0xNTA0ODA2MDk3OyBjdj1ub25lOw0KICAgICAgICBkPWdvb2dsZS5jb207IHM9YXJjLTIwMTYwODE2Ow0KICAgICAgICBiPW5PYW9oY08xU2xhbkg4K3pMd3ZTczF1bmtzc1llZGtEalVJZmtpcmRQejYvQkNrcXcwNkN2b0RQNzgyd3VITXorWg0KICAgICAgICAgUmVsZVlLM2VGbEtHL1lITHMrYUJGbTF2L3pWUWZ2WjBFcnVmZlhJVGl4ZGViY044NmtuUlVIRjMrMXBwUmUxbkNvZ1YNCiAgICAgICAgIHNPUUJUZUp5WFRvdEx2K0NxKy90cXZaV2l6Mk8xNmtLVGRPMElZdGlvMkcwamVFMk4rRDZ3dUcxWmhWMXJxZXU3RHY1DQogICAgICAgICAzQzB3dnZZYnVCbWt6MHVXRXlDUlA0cFdVWnRxUCtSSTJSeDlHekg0TzN4a1lJMDNFWFBScUZJaVgyWERtOEMvWDBGaw0KICAgICAgICAgOTFEbUJYZUN0eU1pbmV6Y1UxaFdOUXBxdmdaTHZGTXJUYmtJeFlCZXBOMERDdXZaUWRjR3I4NDdQMXdXZnluYWIzaWMNCiAgICAgICAgIHlPcWc9PQ0KQVJDLU1lc3NhZ2UtU2lnbmF0dXJlOiBpPTE7IGE9cnNhLXNoYTI1NjsgYz1yZWxheGVkL3JlbGF4ZWQ7IGQ9Z29vZ2xlLmNvbTsgcz1hcmMtMjAxNjA4MTY7DQogICAgICAgIGg9dG86ZnJvbTpzdWJqZWN0Om1lc3NhZ2UtaWQ6ZmVlZGJhY2staWQ6ZGF0ZTptaW1lLXZlcnNpb24NCiAgICAgICAgIDpka2ltLXNpZ25hdHVyZTphcmMtYXV0aGVudGljYXRpb24tcmVzdWx0czsNCiAgICAgICAgYmg9eHd6eTg1dmU3VjBKUVdJTDdHQ0tLdUVKWXZpTDhPUy9pZXNJa1M5VnlmYz07DQogICAgICAgIGI9Q3g3TTkxZ1Y2czNzajN3Y2Q2blFjOENyRzZtNktTN1IzWGV6VktjRlVOdzJIaXAzWUI4WEMyT3A1MmtxNElFQzJkDQogICAgICAgICBrT0xFbDdLNmZHbXYvdGttZmVGYTZ1d2Q0Tlgvcks2cHpBYml5Yzkrc1V0S0RaVkxJcUZhb1VOeGRObTJSSjFLaVlkTQ0KICAgICAgICAgcmJKTEFFU0pEcHliSEt0aXZJQ1ZjNnNGc2R4VjQwYStORDhWcGJocEtLZ3liczF3WWowTnRqQUptNkdXLzkvdVZySmENCiAgICAgICAgIG1KTzF5Ynl4T3BEU0NRY1lNTlkwUU13Y3doa0tQakdrakdzUEhRNlY4bzZJZEs3U3ZFVWUwQUp2Ti94VHRLcXNTTDZVDQogICAgICAgICBZQy9vYm5jTEpBZWdWZzcvcWtTQmNwaTRMdGVZY2tObHkxT0ZQY0o0UGZEUUVJdnhSWjVxeExTdGl1R09hdnExcExjTQ0KICAgICAgICAgWFF6dz09DQpBUkMtQXV0aGVudGljYXRpb24tUmVzdWx0czogaT0xOyBteC5nb29nbGUuY29tOw0KICAgICAgIGRraW09cGFzcyBoZWFkZXIuaT1AYWNjb3VudHMuZ29vZ2xlLmNvbSBoZWFkZXIucz0yMDE2MTAyNSBoZWFkZXIuYj1ISXFsUkFZUDsNCiAgICAgICBzcGY9cGFzcyAoZ29vZ2xlLmNvbTogZG9tYWluIG9mIDMweXN4d3FndGMzcWZnLWp3aGRxc3V1Z21mbGsueWdneWR3LnVnZWVhdXpzd2Rjd2p2eWVzYWQudWdlQGdhaWEuYm91bmNlcy5nb29nbGUuY29tIGRlc2lnbmF0ZXMgMjA5Ljg1LjIyMC42OSBhcyBwZXJtaXR0ZWQgc2VuZGVyKSBzbXRwLm1haWxmcm9tPTMwWVN4V1FnVEMzUWZnLWpXaGRxU1VVZ21mbGsuWWdnWWRXLlVnZWVhVVpTV2RjV2pWWWVTYWQuVWdlQGdhaWEuYm91bmNlcy5nb29nbGUuY29tOw0KICAgICAgIGRtYXJjPXBhc3MgKHA9UkVKRUNUIHNwPVJFSkVDVCBkaXM9Tk9ORSkgaGVhZGVyLmZyb209YWNjb3VudHMuZ29vZ2xlLmNvbQ0KUmV0dXJuLVBhdGg6IDwzMFlTeFdRZ1RDM1FmZy1qV2hkcVNVVWdtZmxrLllnZ1lkVy5VZ2VlYVVaU1dkY1dqVlllU2FkLlVnZUBnYWlhLmJvdW5jZXMuZ29vZ2xlLmNvbT4NClJlY2VpdmVkOiBmcm9tIG1haWwtc29yLWY2OS5nb29nbGUuY29tIChtYWlsLXNvci1mNjkuZ29vZ2xlLmNvbS4gWzIwOS44NS4yMjAuNjldKQ0KICAgICAgICBieSBteC5nb29nbGUuY29tIHdpdGggU01UUFMgaWQgbTFzb3IxMjcyMjVwbGkuMzIuMjAxNy4wOS4wNy4xMC40MS4zNw0KICAgICAgICBmb3IgPG1pY2hhZWxrZXJkQGdtYWlsLmNvbT4NCiAgICAgICAgKEdvb2dsZSBUcmFuc3BvcnQgU2VjdXJpdHkpOw0KICAgICAgICBUaHUsIDA3IFNlcCAyMDE3IDEwOjQxOjM3IC0wNzAwIChQRFQpDQpSZWNlaXZlZC1TUEY6IHBhc3MgKGdvb2dsZS5jb206IGRvbWFpbiBvZiAzMHlzeHdxZ3RjM3FmZy1qd2hkcXN1dWdtZmxrLnlnZ3lkdy51Z2VlYXV6c3dkY3dqdnllc2FkLnVnZUBnYWlhLmJvdW5jZXMuZ29vZ2xlLmNvbSBkZXNpZ25hdGVzIDIwOS44NS4yMjAuNjkgYXMgcGVybWl0dGVkIHNlbmRlcikgY2xpZW50LWlwPTIwOS44NS4yMjAuNjk7DQpBdXRoZW50aWNhdGlvbi1SZXN1bHRzOiBteC5nb29nbGUuY29tOw0KICAgICAgIGRraW09cGFzcyBoZWFkZXIuaT1AYWNjb3VudHMuZ29vZ2xlLmNvbSBoZWFkZXIucz0yMDE2MTAyNSBoZWFkZXIuYj1ISXFsUkFZUDsNCiAgICAgICBzcGY9cGFzcyAoZ29vZ2xlLmNvbTogZG9tYWluIG9mIDMweXN4d3FndGMzcWZnLWp3aGRxc3V1Z21mbGsueWdneWR3LnVnZWVhdXpzd2Rjd2p2eWVzYWQudWdlQGdhaWEuYm91bmNlcy5nb29nbGUuY29tIGRlc2lnbmF0ZXMgMjA5Ljg1LjIyMC42OSBhcyBwZXJtaXR0ZWQgc2VuZGVyKSBzbXRwLm1haWxmcm9tPTMwWVN4V1FnVEMzUWZnLWpXaGRxU1VVZ21mbGsuWWdnWWRXLlVnZWVhVVpTV2RjV2pWWWVTYWQuVWdlQGdhaWEuYm91bmNlcy5nb29nbGUuY29tOw0KICAgICAgIGRtYXJjPXBhc3MgKHA9UkVKRUNUIHNwPVJFSkVDVCBkaXM9Tk9ORSkgaGVhZGVyLmZyb209YWNjb3VudHMuZ29vZ2xlLmNvbQ0KREtJTS1TaWduYXR1cmU6IHY9MTsgYT1yc2Etc2hhMjU2OyBjPXJlbGF4ZWQvcmVsYXhlZDsNCiAgICAgICAgZD1hY2NvdW50cy5nb29nbGUuY29tOyBzPTIwMTYxMDI1Ow0KICAgICAgICBoPW1pbWUtdmVyc2lvbjpkYXRlOmZlZWRiYWNrLWlkOm1lc3NhZ2UtaWQ6c3ViamVjdDpmcm9tOnRvOw0KICAgICAgICBiaD14d3p5ODV2ZTdWMEpRV0lMN0dDS0t1RUpZdmlMOE9TL2llc0lrUzlWeWZjPTsNCiAgICAgICAgYj1ISXFsUkFZUEliMlpnc09vbGJaNUZjaVpCcTZQdHBMVlJVNE9JcG5KdFZIYzEwWm1ab24zVHFIQi9YaEdHRytybS8NCiAgICAgICAgIGJnMHdEYkFMYUV3NjNnMFR1RlZGc0Qyd2kwT3FPTUlBSEFQZWN2RnVTNE9hWHBLWWNDdEtWd3JtNkdlQ0twTUM5SFhiDQogICAgICAgICBlU1YxT0hib0JsWlU4UHBZV1gwaDVCUnd2VXh1VHh6d1NYOHI1VWNGRzBZSnZCYmFSa1pCWDFiWUVxbWJ6MjF3K1ZOWg0KICAgICAgICAgWDgzamNWLzV1bSt0eHRZUHMrVGhncXh3T215bVg2akpNRHh3RjRCY015aGd2NkZWRUVReEFBR1BhMHBQQjZ4dTg4QWYNCiAgICAgICAgIDczYVRvZURmQ1BMbnQzRnhPdWFnTjhKOVR1Zm5zSGVTeHg0UHNMOUVabVg5cGhtWVVWTGRRUHVmUWRUZVVwL1lYcjNrDQogICAgICAgICBaR3p3PT0NClgtR29vZ2xlLURLSU0tU2lnbmF0dXJlOiB2PTE7IGE9cnNhLXNoYTI1NjsgYz1yZWxheGVkL3JlbGF4ZWQ7DQogICAgICAgIGQ9MWUxMDAubmV0OyBzPTIwMTYxMDI1Ow0KICAgICAgICBoPXgtZ20tbWVzc2FnZS1zdGF0ZTptaW1lLXZlcnNpb246ZGF0ZTpmZWVkYmFjay1pZDptZXNzYWdlLWlkOnN1YmplY3QNCiAgICAgICAgIDpmcm9tOnRvOw0KICAgICAgICBiaD14d3p5ODV2ZTdWMEpRV0lMN0dDS0t1RUpZdmlMOE9TL2llc0lrUzlWeWZjPTsNCiAgICAgICAgYj1ZQjhYU3RMUjdwaXY5UHBkWnlrNW1UbzIwQ2tWd2xhRVh1WDd3aGlWYXZFV2xha004L242SmJQaGpCYjhleCtPV3cNCiAgICAgICAgIHdBcklVS1E1Q3NpUlAva3VCNjBQNnprRHd4N2tIcFZlVWZXNk83SnpTcVFGZmxOckJFS09BanNtandXaGRyMGoxdGVvDQogICAgICAgICB6RXQyc3kweTBML2JJSTl6eU1BbDhZWmJ4aGhzQXUxVkdBTmFyVFQ3dnd3TUgwS1hyUWxVdzJQQ1FTalNQR1Fuck92Mg0KICAgICAgICAgaTlPcVZnUXBzMWxvYS9tNHh5Z0RXb3RXSDhNT29EU1hkdWcvaUJyampXMUQ5dXRXM2xtK25EdytRNGZIL3plbi94cGMNCiAgICAgICAgIEV6WVgremlyTnRjeFQ1KzMzYjA3RndpeVdwTis1YWFTV0N1RUcrSDB5a0xQVFoxZ0dQOWRMKytaM1RidDhFZTFqNXY2DQogICAgICAgICBjdGtnPT0NClgtR20tTWVzc2FnZS1TdGF0ZTogQUhQampVZzZwS3pmWVVnVEUzaFBEczc0MzBBaWNoVXBUOVpWUTdHMy9qaWxCRmVia1lUMDBvTDYNCgl0eCtNa3FjZ3c3UGpnRW9sMlFub1hjbFlaSEIzcjJuWA0KWC1Hb29nbGUtU210cC1Tb3VyY2U6IEFES0NOYjVHSFEwSWxkUlg4N2drVVVPSXlZeW1IMUYxZHdzYk5kYnhKVndwRVpFdEYwUnRyM01BUmJyY1ppQzNkNUZvdVR0aWJxRVV0SVlEeVJNWnNKSUk4RVhUTGc9PQ0KTUlNRS1WZXJzaW9uOiAxLjANClgtUmVjZWl2ZWQ6IGJ5IDEwLjg0LjIzMi42NSB3aXRoIFNNVFAgaWQgZjFtcjc4NTYzcGxuLjQ0LjE1MDQ4MDYwOTc4MzA7IFRodSwNCiAwNyBTZXAgMjAxNyAxMDo0MTozNyAtMDcwMCAoUERUKQ0KRGF0ZTogVGh1LCA3IFNlcCAyMDE3IDE3OjQxOjMwICswMDAwIChVVEMpDQpYLU5vdGlmaWNhdGlvbnM6IFhFQUFBQUw4Umx4TWtpZjNfQmdLaV9ST0gzbnMNClgtQWNjb3VudC1Ob3RpZmljYXRpb24tVHlwZTogMTI3DQpGZWVkYmFjay1JRDogMTI3OmFjY291bnQtbm90aWZpZXINCk1lc3NhZ2UtSUQ6IDxEbkg5cXNBMnlWRWtuaWxLaVJSdkNRQG5vdGlmaWNhdGlvbnMuZ29vZ2xlLmNvbT4NClN1YmplY3Q6IE1haWxNaWdyYXRpb24gY29ubmVjdGVkIHRvIHlvdXIgR29vZ2xlIEFjY291bnQNCkZyb206IEdvb2dsZSA8bm8tcmVwbHlAYWNjb3VudHMuZ29vZ2xlLmNvbT4NClRvOiBtaWNoYWVsa2VyZEBnbWFpbC5jb20NCkNvbnRlbnQtVHlwZTogbXVsdGlwYXJ0L2FsdGVybmF0aXZlOyBib3VuZGFyeT0iOTRlYjJjMWEyNjQyZWFjNDcyMDU1ODljZjg1OSINCg0KLS05NGViMmMxYTI2NDJlYWM0NzIwNTU4OWNmODU5DQpDb250ZW50LVR5cGU6IHRleHQvcGxhaW47IGNoYXJzZXQ9IlVURi04IjsgZm9ybWF0PWZsb3dlZDsgZGVsc3A9eWVzDQpDb250ZW50LVRyYW5zZmVyLUVuY29kaW5nOiBiYXNlNjQNCg0KVFdGcGJFMXBaM0poZEdsdmJpQmpiMjV1WldOMFpXUWdkRzhnZVc5MWNpQkhiMjluYkdVZ1FXTmpiM1Z1ZEEwS0RRb05DZzBLU0drZw0KVFdsamFHRmxiQ3dOQ2cwS1RXRnBiRTFwWjNKaGRHbHZiaUJ1YjNjZ2FHRnpJR0ZqWTJWemN5QjBieUI1YjNWeUlFZHZiMmRzWlNCQg0KWTJOdmRXNTBJRzFwWTJoaFpXeHJaWEprUUdkdFlXbHNMbU52YlM0TkNnMEtUV0ZwYkUxcFozSmhkR2x2YmlCallXNDZEUW9OQ2lBZw0KSUNBdElFbHVjMlZ5ZENCdFlXbHNJR2x1ZEc4Z2VXOTFjaUJ0WVdsc1ltOTREUW9OQ2cwS1dXOTFJSE5vYjNWc1pDQnZibXg1SUdkcA0KZG1VZ2RHaHBjeUJoWTJObGMzTWdkRzhnWVhCd2N5QjViM1VnZEhKMWMzUXVJRkpsZG1sbGR5QnZjaUJ5WlcxdmRtVWdZWEJ3Y3cwSw0KWTI5dWJtVmpkR1ZrSUhSdklIbHZkWElnWVdOamIzVnVkQ0JoYm5rZ2RHbHRaU0JoZENCTmVTQkJZMk52ZFc1MERRbzhhSFIwY0hNNg0KTHk5dGVXRmpZMjkxYm5RdVoyOXZaMnhsTG1OdmJTOXdaWEp0YVhOemFXOXVjejR1RFFvTkNreGxZWEp1SUcxdmNtVWdQR2gwZEhCeg0KT2k4dmMzVndjRzl5ZEM1bmIyOW5iR1V1WTI5dEwyRmpZMjkxYm5SekwyRnVjM2RsY2k4ek5EWTJOVEl4UGlCaFltOTFkQ0IzYUdGMA0KRFFwcGRDQnRaV0Z1Y3lCMGJ5QmpiMjV1WldOMElHRnVJR0Z3Y0NCMGJ5QjViM1Z5SUdGalkyOTFiblF1RFFwVWFHVWdSMjl2WjJ4bA0KSUVGalkyOTFiblJ6SUhSbFlXME5DZzBLRFFvTkNsUm9hWE1nWlcxaGFXd2dZMkZ1SjNRZ2NtVmpaV2wyWlNCeVpYQnNhV1Z6TGlCRw0KYjNJZ2JXOXlaU0JwYm1admNtMWhkR2x2Yml3Z2RtbHphWFFnZEdobElFZHZiMmRzWlEwS1FXTmpiM1Z1ZEhNZ1NHVnNjQ0JEWlc1MA0KWlhJZ1BHaDBkSEJ6T2k4dmMzVndjRzl5ZEM1bmIyOW5iR1V1WTI5dEwyRmpZMjkxYm5SekwyRnVjM2RsY2k4ek5EWTJOVEl4UGk0Tg0KQ2cwS0RRb05DbGx2ZFNCeVpXTmxhWFpsWkNCMGFHbHpJRzFoYm1SaGRHOXllU0JsYldGcGJDQnpaWEoyYVdObElHRnVibTkxYm1ObA0KYldWdWRDQjBieUIxY0dSaGRHVWdlVzkxSUdGaWIzVjBEUXBwYlhCdmNuUmhiblFnWTJoaGJtZGxjeUIwYnlCNWIzVnlJRWR2YjJkcw0KWlNCd2NtOWtkV04wSUc5eUlHRmpZMjkxYm5RdURRb05Dc0twSURJd01UY2dSMjl2WjJ4bElFbHVZeTRzSURFMk1EQWdRVzF3YUdsMA0KYUdWaGRISmxJRkJoY210M1lYa3NJRTF2ZFc1MFlXbHVJRlpwWlhjc0lFTkJJRGswTURRekxDQlZVMEVOQ21WME9qRXlOdzBLDQotLTk0ZWIyYzFhMjY0MmVhYzQ3MjA1NTg5Y2Y4NTkNCkNvbnRlbnQtVHlwZTogdGV4dC9odG1sOyBjaGFyc2V0PSJVVEYtOCINCkNvbnRlbnQtVHJhbnNmZXItRW5jb2Rpbmc6IHF1b3RlZC1wcmludGFibGUNCg0KPGh0bWwgbGFuZz0zRCJlbiI-PGhlYWQ-PG1ldGEgbmFtZT0zRCJmb3JtYXQtZGV0ZWN0aW9uIiBjb250ZW50PTNEImRhdGU9M0RuPQ0KbyIvPjxtZXRhIG5hbWU9M0QiZm9ybWF0LWRldGVjdGlvbiIgY29udGVudD0zRCJlbWFpbD0zRG5vIi8-PHN0eWxlPkBtZWRpYSBzPQ0KY3JlZW4gYW5kIChtaW4td2lkdGg6IDYwMHB4KSB7LnYyc3Age3BhZGRpbmc6IDZweCAzNHB4IDBweDt9fTwvc3R5bGU-PC9oZWFkPQ0KPjxib2R5IHN0eWxlPTNEIm1hcmdpbjogMDsgcGFkZGluZzogMDsiIGJnY29sb3I9M0QiI0ZGRkZGRiI-PHRhYmxlIHdpZHRoPTNEPQ0KIjEwMCUiIGhlaWdodD0zRCIxMDAlIiBzdHlsZT0zRCJtaW4td2lkdGg6IDM0OHB4OyIgYm9yZGVyPTNEIjAiIGNlbGxzcGFjaW5nPQ0KPTNEIjAiIGNlbGxwYWRkaW5nPTNEIjAiPjx0ciBoZWlnaHQ9M0QiMzJweCI-PC90cj48dHIgYWxpZ249M0QiY2VudGVyIj48dGQgPQ0Kd2lkdGg9M0QiMzJweCI-PC90ZD48dGQ-PHRhYmxlIGJvcmRlcj0zRCIwIiBjZWxsc3BhY2luZz0zRCIwIiBjZWxscGFkZGluZz0NCj0zRCIwIiBzdHlsZT0zRCJtYXgtd2lkdGg6IDYwMHB4OyI-PHRyPjx0ZD48dGFibGUgd2lkdGg9M0QiMTAwJSIgYm9yZGVyPTNEIj0NCjAiIGNlbGxzcGFjaW5nPTNEIjAiIGNlbGxwYWRkaW5nPTNEIjAiPjx0cj48dGQgYWxpZ249M0QibGVmdCI-PGltZyB3aWR0aD0zRD0NCiI5MiIgaGVpZ2h0PTNEIjMyIiBzcmM9M0QiaHR0cHM6Ly93d3cuZ3N0YXRpYy5jb20vYWNjb3VudGFsZXJ0cy9lbWFpbC9nb29nbD0NCmVsb2dvX2NvbG9yXzE4OHg2NGRwLnBuZyIgc3R5bGU9M0QiZGlzcGxheTogYmxvY2s7IHdpZHRoOiA5MnB4OyBoZWlnaHQ6IDMycD0NCng7Ij48L3RkPjx0ZCBhbGlnbj0zRCJyaWdodCI-PGltZyB3aWR0aD0zRCIzMiIgaGVpZ2h0PTNEIjMyIiBzdHlsZT0zRCJkaXNwbD0NCmF5OiBibG9jazsgd2lkdGg6IDMycHg7IGhlaWdodDogMzJweDsiIHNyYz0zRCJodHRwczovL3d3dy5nc3RhdGljLmNvbS9hY2NvdT0NCm50YWxlcnRzL2VtYWlsL2tleWhvbGUucG5nIj48L3RkPjwvdHI-PC90YWJsZT48L3RkPjwvdHI-PHRyIGhlaWdodD0zRCIxNiI-PD0NCi90cj48dHI-PHRkPjx0YWJsZSBiZ2NvbG9yPTNEIiM0MTg0RjMiIHdpZHRoPTNEIjEwMCUiIGJvcmRlcj0zRCIwIiBjZWxsc3BhYz0NCmluZz0zRCIwIiBjZWxscGFkZGluZz0zRCIwIiBzdHlsZT0zRCJtaW4td2lkdGg6IDMzMnB4OyBtYXgtd2lkdGg6IDYwMHB4OyBibz0NCnJkZXI6IDFweCBzb2xpZCAjRjBGMEYwOyBib3JkZXItYm90dG9tOiAwOyBib3JkZXItdG9wLWxlZnQtcmFkaXVzOiAzcHg7IGJvcj0NCmRlci10b3AtcmlnaHQtcmFkaXVzOiAzcHg7Ij48dHI-PHRkIGhlaWdodD0zRCI3MnB4IiBjb2xzcGFuPTNEIjMiPjwvdGQ-PC90cj0NCj48dHI-PHRkIHdpZHRoPTNEIjMycHgiPjwvdGQ-PHRkIHN0eWxlPTNEImZvbnQtZmFtaWx5OiBSb2JvdG8tUmVndWxhcixIZWx2ZT0NCnRpY2EsQXJpYWwsc2Fucy1zZXJpZjsgZm9udC1zaXplOiAyNHB4OyBjb2xvcjogI0ZGRkZGRjsgbGluZS1oZWlnaHQ6IDEuMjU7ID0NCm1pbi13aWR0aDogMzAwcHg7Ij5NYWlsTWlncmF0aW9uIGNvbm5lY3RlZCB0byB5b3VyIEdvb2dsZSBBY2NvdW50PC90ZD48dGQgdz0NCmlkdGg9M0QiMzJweCI-PC90ZD48L3RyPjx0cj48dGQgaGVpZ2h0PTNEIjE4cHgiIGNvbHNwYW49M0QiMyI-PC90ZD48L3RyPjwvdD0NCmFibGU-PC90ZD48L3RyPjx0cj48dGQ-PHRhYmxlIGJnY29sb3I9M0QiI0ZBRkFGQSIgd2lkdGg9M0QiMTAwJSIgYm9yZGVyPTNEIj0NCjAiIGNlbGxzcGFjaW5nPTNEIjAiIGNlbGxwYWRkaW5nPTNEIjAiIHN0eWxlPTNEIm1pbi13aWR0aDogMzMycHg7IG1heC13aWR0aD0NCjogNjAwcHg7IGJvcmRlcjogMXB4IHNvbGlkICNGMEYwRjA7IGJvcmRlci1ib3R0b206IDFweCBzb2xpZCAjQzBDMEMwOyBib3JkZT0NCnItdG9wOiAwOyBib3JkZXItYm90dG9tLWxlZnQtcmFkaXVzOiAzcHg7IGJvcmRlci1ib3R0b20tcmlnaHQtcmFkaXVzOiAzcHg7Ij0NCj48dHIgaGVpZ2h0PTNEIjE2cHgiPjx0ZCB3aWR0aD0zRCIzMnB4IiByb3dzcGFuPTNEIjMiPjwvdGQ-PHRkPjwvdGQ-PHRkIHdpZD0NCnRoPTNEIjMycHgiIHJvd3NwYW49M0QiMyI-PC90ZD48L3RyPjx0cj48dGQ-PHRhYmxlIHN0eWxlPTNEIm1pbi13aWR0aDogMzAwcD0NCng7IiBib3JkZXI9M0QiMCIgY2VsbHNwYWNpbmc9M0QiMCIgY2VsbHBhZGRpbmc9M0QiMCI-PHRyPjx0ZCBzdHlsZT0zRCJmb250LT0NCmZhbWlseTogUm9ib3RvLVJlZ3VsYXIsSGVsdmV0aWNhLEFyaWFsLHNhbnMtc2VyaWY7IGZvbnQtc2l6ZTogMTNweDsgY29sb3I6ID0NCiMyMDIwMjA7IGxpbmUtaGVpZ2h0OiAxLjU7cGFkZGluZy1ib3R0b206IDRweDsiPkhpIE1pY2hhZWwsPC90ZD48L3RyPjx0cj48dD0NCmQgc3R5bGU9M0QiZm9udC1mYW1pbHk6IFJvYm90by1SZWd1bGFyLEhlbHZldGljYSxBcmlhbCxzYW5zLXNlcmlmOyBmb250LXNpej0NCmU6IDEzcHg7IGNvbG9yOiAjMjAyMDIwOyBsaW5lLWhlaWdodDogMS41O3BhZGRpbmc6IDRweCAwOyI-PGJyPk1haWxNaWdyYXRpbz0NCm4gbm93IGhhcyBhY2Nlc3MgdG8geW91ciBHb29nbGUgQWNjb3VudCA8YT5taWNoYWVsa2VyZEBnbWFpbC5jb208L2E-Ljxicj48Yj0NCnI-TWFpbE1pZ3JhdGlvbiBjYW46PHVsIHN0eWxlPTNEIm1hcmdpbjogMDsiPjxsaT5JbnNlcnQgbWFpbCBpbnRvIHlvdXIgbWFpbD0NCmJveDwvbGk-PC91bD48YnI-WW91IHNob3VsZCBvbmx5IGdpdmUgdGhpcyBhY2Nlc3MgdG8gYXBwcyB5b3UgdHJ1c3QuIFJldmlldz0NCiBvciByZW1vdmUgYXBwcyBjb25uZWN0ZWQgdG8geW91ciBhY2NvdW50IGFueSB0aW1lIGF0IDxhIGhyZWY9M0QiaHR0cHM6Ly9teT0NCmFjY291bnQuZ29vZ2xlLmNvbS9wZXJtaXNzaW9ucyIgc3R5bGU9M0QidGV4dC1kZWNvcmF0aW9uOiBub25lOyBjb2xvcjogIzQyOD0NCjVGNDsiIHRhcmdldD0zRCJfYmxhbmsiPk15IEFjY291bnQ8L2E-Ljxicj48YnI-PGEgaHJlZj0zRCJodHRwczovL3N1cHBvcnQuZz0NCm9vZ2xlLmNvbS9hY2NvdW50cy9hbnN3ZXIvMzQ2NjUyMSIgc3R5bGU9M0QidGV4dC1kZWNvcmF0aW9uOiBub25lOyBjb2xvcjogIz0NCjQyODVGNDsiIHRhcmdldD0zRCJfYmxhbmsiPkxlYXJuIG1vcmU8L2E-IGFib3V0IHdoYXQgaXQgbWVhbnMgdG8gY29ubmVjdCBhbj0NCiBhcHAgdG8geW91ciBhY2NvdW50LjwvdGQ-PC90cj48dHI-PHRkIHN0eWxlPTNEImZvbnQtZmFtaWx5OiBSb2JvdG8tUmVndWxhcj0NCixIZWx2ZXRpY2EsQXJpYWwsc2Fucy1zZXJpZjsgZm9udC1zaXplOiAxM3B4OyBjb2xvcjogIzIwMjAyMDsgbGluZS1oZWlnaHQ6ID0NCjEuNTsgcGFkZGluZy10b3A6IDI4cHg7Ij5UaGUgR29vZ2xlIEFjY291bnRzIHRlYW08L3RkPjwvdHI-PHRyIGhlaWdodD0zRCIxNj0NCnB4Ij48L3RyPjx0cj48dGQ-PHRhYmxlIHN0eWxlPTNEImZvbnQtZmFtaWx5OiBSb2JvdG8tUmVndWxhcixIZWx2ZXRpY2EsQXJpYT0NCmwsc2Fucy1zZXJpZjsgZm9udC1zaXplOiAxMnB4OyBjb2xvcjogI0I5QjlCOTsgbGluZS1oZWlnaHQ6IDEuNTsiPjx0cj48dGQ-VD0NCmhpcyBlbWFpbCBjYW4ndCByZWNlaXZlIHJlcGxpZXMuIEZvciBtb3JlIGluZm9ybWF0aW9uLCB2aXNpdCB0aGUgPGEgaHJlZj0zRD0NCiJodHRwczovL3N1cHBvcnQuZ29vZ2xlLmNvbS9hY2NvdW50cy9hbnN3ZXIvMzQ2NjUyMSIgZGF0YS1tZXRhLWtleT0zRCJoZWxwIj0NCiBzdHlsZT0zRCJ0ZXh0LWRlY29yYXRpb246IG5vbmU7IGNvbG9yOiAjNDI4NUY0OyIgdGFyZ2V0PTNEIl9ibGFuayI-R29vZ2xlID0NCkFjY291bnRzIEhlbHAgQ2VudGVyPC9hPi48L3RkPjwvdHI-PC90YWJsZT48L3RkPjwvdHI-PC90YWJsZT48L3RkPjwvdHI-PHRyID0NCmhlaWdodD0zRCIzMnB4Ij48L3RyPjwvdGFibGU-PC90ZD48L3RyPjx0ciBoZWlnaHQ9M0QiMTYiPjwvdHI-PHRyPjx0ZCBzdHlsZT0NCj0zRCJtYXgtd2lkdGg6IDYwMHB4OyBmb250LWZhbWlseTogUm9ib3RvLVJlZ3VsYXIsSGVsdmV0aWNhLEFyaWFsLHNhbnMtc2VyaT0NCmY7IGZvbnQtc2l6ZTogMTBweDsgY29sb3I6ICNCQ0JDQkM7IGxpbmUtaGVpZ2h0OiAxLjU7Ij48dHI-PHRkPjx0YWJsZSBzdHlsZT0NCj0zRCJmb250LWZhbWlseTogUm9ib3RvLVJlZ3VsYXIsSGVsdmV0aWNhLEFyaWFsLHNhbnMtc2VyaWY7IGZvbnQtc2l6ZTogMTBweD0NCjsgY29sb3I6ICM2NjY2NjY7IGxpbmUtaGVpZ2h0OiAxOHB4OyBwYWRkaW5nLWJvdHRvbTogMTBweCI-PHRyPjx0ZD5Zb3UgcmVjZT0NCml2ZWQgdGhpcyBtYW5kYXRvcnkgZW1haWwgc2VydmljZSBhbm5vdW5jZW1lbnQgdG8gdXBkYXRlIHlvdSBhYm91dCBpbXBvcnRhbj0NCnQgY2hhbmdlcyB0byB5b3VyIEdvb2dsZSBwcm9kdWN0IG9yIGFjY291bnQuPC90ZD48L3RyPjx0ciBoZWlnaHQ9M0QiNnB4Ij48Lz0NCnRyPjx0cj48dGQ-PGRpdiBzdHlsZT0zRCJkaXJlY3Rpb246IGx0cjsgdGV4dC1hbGlnbjogbGVmdCI-JmNvcHk7IDIwMTcgR29vZz0NCmxlIEluYy4sIDE2MDAgQW1waGl0aGVhdHJlIFBhcmt3YXksIE1vdW50YWluIFZpZXcsIENBIDk0MDQzLCBVU0E8L2Rpdj48ZGl2ID0NCnN0eWxlPTNEImRpc3BsYXk6IG5vbmUgIWltcG9ydGFudDsgbXNvLWhpZGU6YWxsOyBtYXgtaGVpZ2h0OjBweDsgbWF4LXdpZHRoOj0NCjBweDsiPmV0OjEyNzwvZGl2PjwvdGQ-PC90cj48L3RhYmxlPjwvdGQ-PC90cj48L3RkPjwvdHI-PC90YWJsZT48L3RkPjx0ZCB3aT0NCmR0aD0zRCIzMnB4Ij48L3RkPjwvdHI-PHRyIGhlaWdodD0zRCIzMnB4Ij48L3RyPjwvdGFibGU-PC9ib2R5PjwvaHRtbD4NCi0tOTRlYjJjMWEyNjQyZWFjNDcyMDU1ODljZjg1OS0tDQo="},
    //     "labelIds": {["INBOX", "UNREAD"]},
    // }
        
    // body := bytes.NewBufferString(bodyVals.Encode())

    // req, _ := http.NewRequest("POST", urlStr, body)
    // req.Header.Set("Authorization", "Bearer " + sourceToken)

    // req, _ := http.NewRequest("POST", urlStr, nil)
    

    // ctx := appengine.NewContext(r)
    // client := urlfetch.Client(ctx)

    // resp, err := client.Do(req)

    // body := resp.Body
    // defer body.Close()

    // if body == nil {
    //     http.Error(w, "Response body not found", 400)
    //     return
    // }

    // respBody, _ := ioutil.ReadAll(body)
    // log.Printf("HTTP PostForm/GET returned %v", string(respBody))

    log.Print("Printing Source Token:::!!!!!")
    log.Print(ctx.Value("cookieInfo").(Values).Get("sourceToken"))
    
    err = runtime.RunInBackground(ctx, startTransfer)
    if err != nil {
            log.Printf("Could not start background thread: %v", err)
            return
    }






}
