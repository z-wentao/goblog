package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts/{slug}", PostHandler(FileReader{}))
	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal(err)
	}
}

type SlugReader interface {
	Read(slug string) (string, error)
}

type FileReader struct{}

func (fr FileReader) Read(slug string) (string, error) {
	f, err := os.Open(slug + ".md")
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func PostHandler(sl SlugReader) http.HandlerFunc {
	markRender := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
			),
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		postMarkdown, err := sl.Read(slug)
		if err != nil {
			http.Error(w, "Post Not found", http.StatusNotFound)
			return
		}
		var buf bytes.Buffer
		err = markRender.Convert([]byte(postMarkdown), &buf)
		if err != nil {
			panic(err)
		}
		io.Copy(w, &buf)
	}
}
