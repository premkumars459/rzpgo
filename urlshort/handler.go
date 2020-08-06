package urlshort

import (
	yaml "gopkg.in/yaml.v2"
	"net/http"
)


// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//fmt.Printf("good going")
	return func (w http.ResponseWriter , r *http.Request){
		toUrl, present := pathsToUrls[r.URL.Path]
		if present {
			http.Redirect (w, r ,toUrl,301)
		}else {
			fallback.ServeHTTP(w,r) }
	}
	//	TODO: Implement this...
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	var yamlrecords []yamlparse
	err := yaml.Unmarshal(yml, &yamlrecords)
	if err != nil {
		return nil, err 
	}
	//fmt.Println(yamlrecords)    // to check if the yaml was successfully parsed.

	yamlpathsToUrls := map[string]string{}
	for _, record:= range yamlrecords {
		yamlpathsToUrls[record.Path] = record.Url
	}



	return MapHandler(yamlpathsToUrls, fallback), nil 
}


type yamlparse struct {
	Path  string `yaml:"path"`
	Url string  `yaml:"url"`
}
