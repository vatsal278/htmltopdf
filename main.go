package main

import (
	"bytes"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/vatsal278/go-redis-cache"
	"io/ioutil"
	"log"
)

func main() {
	pdfg := wkhtmltopdf.NewPDFPreparer()
	htmlfile, err := ioutil.ReadFile("newfile.html")
	if err != nil {
		log.Fatal(err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(htmlfile)))

	// The contents of htmlsimple.html are saved as base64 string in the JSON file
	jb, err := pdfg.ToJSON()
	if err != nil {
		log.Fatal(err)
		return
	}
	cacher := redis.NewCacher(redis.Config{Addr: "172.18.0.2:6379"})
	err = cacher.Set("newfile.html", jb, 0)
	if err != nil {
		log.Print(err)
		return
	}
	y, err := cacher.Get("newfile.html")
	if err != nil {
		log.Print(err)
		return
	}

	// Create PDF document in internal buffer
	pdfgFromJSON, err := wkhtmltopdf.NewPDFGeneratorFromJSON(bytes.NewReader(y))
	if err != nil {
		log.Fatal(err)
	}

	err = pdfgFromJSON.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile("./simplesample.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	// Output: Done
}
