package main

import (
	"bytes"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gorilla/mux"
	"github.com/vatsal278/go-redis-cache"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Student struct {
	Name  string
	Marks int
	Id    string
}

type Class []Student

func main() {
	/*pdfg := wkhtmltopdf.NewPDFPreparer()
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
	//pdfgFromJSON.SetOutput()
	err = pdfgFromJSON.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	err = pdfgFromJSON.WriteFile("./simplesample.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	// Output: Done*/
	cacher := redis.NewCacher(redis.Config{Addr: "127.0.0.1:9096"})
	r := mux.NewRouter()
	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("File Upload Endpoint Hit")
		err := r.ParseMultipartForm(10240)
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

		// The contents of htmlsimple.html are saved as base64 string in the JSON file
		jb, err := pdfg.ToJSON()
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Print(handler.Filename)
		err = cacher.Set(handler.Filename, jb, 0)
		if err != nil {
			log.Printf("failed writing to file: %s", err)
			return
		}
		log.Println("Successfully Uploaded File\n")
	}).Methods("POST")
	r.HandleFunc("/topdf", func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query().Get("file")
		//b, err := ioutil.ReadFile("uploads/" + v)
		b, err := cacher.Get(v + ".html")
		if err != nil {
			log.Print(err.Error())
			return
		}
		err = ioutil.WriteFile("uploads/"+v+".html", b, 0644)
		if err != nil {
			log.Printf("failed writing to file: %s", err)
			return
		}
		var class Class
		// defining struct instance
		std1 := Student{"A", 90, "1"}
		std2 := Student{"B", 100, "2"}
		std3 := Student{"C", 88, "3"}
		std4 := Student{"D", 25, "4"}
		std5 := Student{"E", 35, "5"}
		class = append(class, std4, std2, std3, std1, std5)
		//os.Setenv("list", "")
		//list := os.Getenv("list")

		t, err := template.ParseFiles("uploads/" + v + ".html")
		if err != nil {
			log.Print(err)
		}
		f, err := os.Create("tempfile.html")
		if err != nil {
			log.Println("create file: ", err)
			return
		}
		err = t.Execute(f, class)
		if err != nil {
			log.Print(err)
		}
		/* standard output to print merged data
		err = t.Execute(os.Stdout, class)
		if err != nil {
			log.Print(err)
		}*/
		w.Header().Set("Content-Disposition", "attachment; filename="+v)
		w.Header().Set("Content-Type", "application/octet-stream")
		// Create PDF document in internal buffer
		pdfgFromJSON, err := wkhtmltopdf.NewPDFGeneratorFromJSON(bytes.NewReader(b))
		if err != nil {
			log.Fatal(err)
		}
		//pdfgFromJSON.SetOutput()
		err = pdfgFromJSON.Create()
		if err != nil {
			log.Fatal(err)
		}

		// Write buffer contents to file on disk
		err = pdfgFromJSON.WriteFile("./simplesample.pdf")
		if err != nil {
			log.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			log.Println("unable to write image.")
			return
		}
	}).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9080", r))
	/*r := gin.Default()
	r.GET("/pong/:id", func(c *gin.Context) {
		var post Post
		x := c.Param("id")

		data, err := cacher.RedisGet(x)
		if err != nil && err != redis.Nil {
			log.Print(err)
			return
		}
		if err == nil {
			err = json.Unmarshal(data, &post)
			if err != nil {
				log.Print(err)
			}
			log.Print("cached data")
			json.NewEncoder(c.Writer).Encode(post)
			return
		}
		result, err := db.Query("SELECT id, title FROM posts WHERE id = ?", x)
		if err != nil {
			log.Print(err)
		}
		defer result.Close()
		for result.Next() {
			err := result.Scan(&post.ID, &post.Title)
			if err != nil {
				log.Print(err.Error())
			}
		}
		y, err := json.Marshal(post)
		err = sdk.RedisSet(post.ID, y, 5)
		if err != nil {
			log.Print(err)
		}
		json.NewEncoder(c.Writer).Encode(post)
		return
		log.Print()

	})
	r.POST("/ping", func(c *gin.Context) {

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err.Error())
		}

		json.Unmarshal(body, &posts)

		c.JSON(http.StatusCreated, posts)
	})

	r.Run(":8086")*/
}
