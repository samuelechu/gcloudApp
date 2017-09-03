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
		Name: "my-cookie",
		Value: "some value",
	})

	cookie, err := r.Cookie("my-cookie")
	fmt.Fprintf(w, "Cookie: %v, Err: %v", cookie, err)
}
