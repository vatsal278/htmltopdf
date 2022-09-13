package main

import (
	"bytes"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gorilla/mux"
	"github.com/varsal278/htmltopdf/htmltopdf"
	"github.com/vatsal278/go-redis-cache"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var class htmltopdf.Class
	// defining struct instance
	std1 := htmltopdf.Student{"A", 90, "1"}
	std2 := htmltopdf.Student{"B", 100, "2"}
	std3 := htmltopdf.Student{"C", 88, "3"}
	std4 := htmltopdf.Student{"D", 25, "4"}
	std5 := htmltopdf.Student{"E", 35, "5"}
	class = append(class, std4, std2, std3, std1, std5)
	cacher := redis.NewCacher(redis.Config{Addr: "172.19.0.2:6379"})
	r := mux.NewRouter()
	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("File Upload Endpoint Hit")
		err := r.ParseMultipartForm(102400)
		if err != nil {
			log.Print(err.Error())
			return
		}
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println("Error Retrieving the File")
			log.Println(err)
			return
		}
		defer file.Close()
		log.Printf("Uploaded File: %+v\n", handler.Filename)
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		pdfg := wkhtmltopdf.NewPDFPreparer()
		pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(fileBytes)))
		jb, err := pdfg.ToJSON()
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Print(handler.Filename)
		err = cacher.Set(handler.Filename, jb, 0)
		if err != nil {
			log.Printf("failed to save cache: %s", err)
			return
		}
		log.Println("Successfully Uploaded File\n")
	}).Methods("POST")
	r.HandleFunc("/topdf", func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query().Get("file")
		b, err := cacher.Get(v + ".html")
		if err != nil {
			log.Print(err.Error())
			return
		}
		ioutil.WriteFile("before.json", b, 0777)
		err = htmltopdf.HtmlToPdf(w, v, class, b)
		if err != nil {
			log.Print(err.Error())
			return
		}
	}).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9080", r))
}
