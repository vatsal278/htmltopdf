### References:
https://blog.gopheracademy.com/advent-2017/using-go-templates/

docker command to create the volumes

```
docker run -v "//c/Users/Perennial/Downloads/files:/app/files" -it --rm ff5b0f38fb11 wkhtmltopdf http://google.com google.pdf
```

docker command to generate pdf from html file

```
docker run -v "/c/Program Files/work/htmltopdf:/app" -v "//app:/c/Program Files/work/htmltopdf" -it --rm ff5b0f38fb11 go run main.go
```