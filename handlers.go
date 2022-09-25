package urlshortener

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type yamlBody struct {
	YamlPath string `yaml:"path"`
	YamlURL  string `yaml:"url"`
}

type jsonBody struct {
	JsonPath string `json:"path"`
	JsonURL  string `json:"url"`
}

func urlShortnerHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
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
	var yData []yamlBody
	err := yaml.Unmarshal([]byte(yml), &yData)
	if err != nil {
		log.Printf("cannot unmarshal data: %v", err)
		return nil, err
	}
	// create a map to match the path entered
	pathurls := make(map[string]string)
	for _, v := range yData {
		pathurls[v.YamlPath] = v.YamlURL
	}
	return urlShortnerHandler(pathurls, fallback), nil
}

func JSONHandler(jsonFile []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var jData []jsonBody

	err := json.Unmarshal(jsonFile, &jData)
	if err != nil {
		log.Printf("cannot unmarshal data: %v", err)
		return nil, err
	}

	pathurls := make(map[string]string)
	for _, v := range jData {
		pathurls[v.JsonPath] = v.JsonURL
	}
	return urlShortnerHandler(pathurls, fallback), nil
}
