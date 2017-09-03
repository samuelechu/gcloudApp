package cookieTest

import(
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/testCookie", handleCookie)
}


func handleCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("my-cookie")
	fmt.Println(cookie, err)

	http.SetCookie(w, &http.Cookie{
		Name: "my-cookie",
		Value: "some value",
	})
}
