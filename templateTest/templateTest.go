package templateTest

import (
	"net/http"
	"html/template"
)

func init() {
        http.HandleFunc("/testTemplate", handler)
}


func handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("header.html", "footer.html")
    t.Execute(w, map[string] string {"Title": "My title", "Body": "Hi this is my body"})
}



