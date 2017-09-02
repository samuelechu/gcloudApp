package templateTest

import (
	"net/http"
	"html/template"
)

func init() {
        http.HandleFunc("/testTemplate", handler)
}

type Person struct {
    UserName string
}

func handler(w http.ResponseWriter, r *http.Request) {
    t := template.New("header.html")
    t, _ = t.ParseFiles("header.html",nil)
    p := Person{UserName: "Astaxie"}

    t.Execute(w, p)

}

