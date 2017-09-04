package cookieTest

import(
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/testCookie", handleCookie)
}


func handleCookie(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name: "source",
		Value: "id_token source",
	})

	http.SetCookie(w, &http.Cookie{
		Name: "destination",
		Value: "id_token dest",
	})

	cookie, err := r.Cookie("source")
	cookie1, err1 := r.Cookie("destination")

	fmt.Fprintf(w, "Source Cookie: %v, Err: %v\n", cookie, err)
	fmt.Fprintf(w, "Dest Cookie: %v, Err: %v\n", cookie1, err1)
}
