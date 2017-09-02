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
   // t := template.New("header.tmpl")


    p := Person{UserName: "Astaxie"}
   // t, _ = t.ParseFiles("header.tmpl")
 


    template.Must(template.ParseFiles("header.tmpl")).Execute(w, p)


}


