package htmltopdf

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"log"
	"net/http"
)

type Student struct {
	Name  string
	Marks int
	Id    string
}
type Class []Student

func HtmlToPdf(w http.ResponseWriter, v string, class Class, b []byte) error {

	var z map[string]interface{}
	err := json.NewDecoder(bytes.NewBuffer(b)).Decode(&z)
	if err != nil {
		log.Fatal("error unmarshaling JSON: %s", err)
	}
	for i, p := range z["Pages"].([]interface{}) {
		page := p.(map[string]interface{})
		buf, err := base64.StdEncoding.DecodeString(page["Base64PageData"].(string))
		if err != nil {
			log.Printf("error decoding base 64 input on page %d: %s", i, err)
			return err
		}
		t, err := template.New(v).Parse(string(buf))
		if err != nil {
			log.Fatal(err)
			return err
		}
		buffer := bytes.NewBuffer(nil)
		err = t.Execute(buffer, class)
		if err != nil {
			log.Println(err)
			return err
		}
		page["Base64PageData"] = base64.StdEncoding.EncodeToString(buffer.Bytes())
	}
	buff := bytes.NewBuffer(nil)
	err = json.NewEncoder(buff).Encode(z)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println(string(b))
	pdfgFromJSON, err := wkhtmltopdf.NewPDFGeneratorFromJSON(bytes.NewBuffer(buff.Bytes()))
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = pdfgFromJSON.Create()
	if err != nil {
		log.Fatal(err)
		return err
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+v+".pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(pdfgFromJSON.Bytes())
	return nil

}
