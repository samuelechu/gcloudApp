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
   	t := template.New("index.html")


    p := Person{UserName: "Sam"}
   	t, _ = t.ParseFiles("static/index.html")
 
   	t.Execute(w, p)

}


