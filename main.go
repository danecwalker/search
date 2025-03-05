package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

type Resource struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

type Resources map[string]Resource

//go:embed resources.json
var resource_file []byte

//go:embed index.html
var index []byte

//go:embed styles.css
var styles []byte

//go:embed reset.css
var reset []byte

func loadResources() (Resources, error) {
	var r Resources

	err := json.Unmarshal(resource_file, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func main() {
	resources, err := loadResources()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(styles)
	})

	mux.HandleFunc("/reset.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(reset)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := strings.TrimSpace(r.URL.Query().Get("q"))

		if q == "" {
			w.WriteHeader(200)
			w.Write(index)
			return
		}

		split := strings.Split(q, " ")
		if split[len(split)-1][0] == '@' {
			if res, ok := resources[split[len(split)-1][1:]]; ok {
				t, err := template.New("").Parse(res.Query)
				if err != nil {
					return
				}

				var loc strings.Builder
				t.Execute(&loc, map[string]string{
					"s": url.QueryEscape(strings.Join(split[0:len(split)-1], " ")),
				})
				w.Header().Set("Location", loc.String())
				w.WriteHeader(303)
				return
			}
		}

		t, err := template.New("").Parse("https://www.google.com/search?q={{.s}}")
		if err != nil {
			return
		}
		var loc strings.Builder
		t.Execute(&loc, map[string]string{
			"s": url.QueryEscape(strings.Join(split, " ")),
		})
		w.Header().Set("Location", loc.String())
		w.WriteHeader(303)
		return
	})

	fmt.Println("serving")
	if err = http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}
