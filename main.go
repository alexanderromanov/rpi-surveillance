package main

import (
	"flag"
	"github.com/alexanderromanov/rpi-surveillance/camera"
	"html/template"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

var port = flag.String("port", "8080", "HTTP listen port")
var pictureWidth = flag.Int("width", 1024, "picture width")
var pictureHeight = flag.Int("height", 800, "picture height")
var interval = flag.Int("interval", 5000, "interval between pictures")

var photoUrl = "/photo"

var homePageTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Live photo</title>
	</head>
	<body>
		<h1>Live photo</h1>
		<img src="{{.PhotoUrl}}"/>
	</body>
</html>
`

var homepageData = struct {
	PhotoUrl string
}{
	PhotoUrl: photoUrl,
}

func check(context string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", context, err)
	}
}

func main() {
	homepage, err := template.New("homepage").Parse(homePageTemplate)
	check("parse homepage template", err)

	go camera.StartTakingPictures(*pictureWidth, *pictureHeight, *interval)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			writer.WriteHeader(404)
			return
		}

		err := homepage.Execute(writer, homepageData)
		check("render main page", err) // todo: don't exit program when single write fails
	})

	http.HandleFunc(photoUrl, func(writer http.ResponseWriter, _ *http.Request) {
		writeStreamOutput(writer)
	})

	log.Println("start listening")
	err = http.ListenAndServe(":"+*port, nil)
	check("start HTTP server", err)
}

func writeStreamOutput(w http.ResponseWriter) {
	headers := w.Header()
	headers.Set("Content-Type", "image/jpg")
	headers.Set("Cache-Control", "no-cache")

	picture := camera.LatestPicture()

	_, err := w.Write(picture)
	check("write response", err)
}
