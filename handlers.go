package urlshortener

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type yamlBody struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		// if the path is matched in the map, then redirect
		if dest, ok := pathsToUrls[path]; ok {
			log.Println("redirecting ....", pathsToUrls[path])
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		// else fallback to the fallback http.Handler
		fallback.ServeHTTP(w, r)

	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// create a slice of yaml body to get unmarshalled objects
	var yb []yamlBody
	err := yaml.Unmarshal([]byte(yml), &yb)
	if err != nil {
		log.Printf("cannot unmarshal data: %v", err)
		return nil, err
	}

	// create a map to match the path entered
	pathurls := make(map[string]string)

	for _, v := range yb {
		pathurls[v.Path] = v.URL
	}

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathurls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}, nil
}
