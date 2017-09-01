package templateTest

import (
	"net/http"
	"html/template"
)

func init() {
        http.HandleFunc("/testTemplate", handler)
}

type Person struct {
    Title string
}

func handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("header.html")
    p := Person{Title: "Astaxie"}
    t.Execute(w, p)

}


