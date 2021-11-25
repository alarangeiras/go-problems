package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type yamlContent struct {
	Path string
	Url  string
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := mapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandlerFunc := yamlHandler([]byte(yaml), mapHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandlerFunc)
}

func yamlHandler(yamlBytes []byte, defaultHandler http.Handler) http.HandlerFunc {
	content := extractYamlSlice(yamlBytes)
	contentMap := extractMap(content)
	return mapHandler(contentMap, defaultHandler)
	
}

func extractMap(content []yamlContent) map[string]string {
	contentMap := make(map[string]string, 0)
	for _, element := range content {
		contentMap[element.Path] = element.Url
	}
	return contentMap
}

func extractYamlSlice(yamlBytes []byte) []yamlContent {
	y := make([]yamlContent, 0)
	err := yaml.Unmarshal(yamlBytes, &y)
	if err != nil {
		log.Fatalln(err)
	}
	return y
}

func mapHandler(paths map[string]string, defaultMux http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		url, exists := paths[urlPath]
		if exists {
			http.Redirect(w, r, url, 301)
		}
		defaultMux.ServeHTTP(w, r)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
