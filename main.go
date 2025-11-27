package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

type Shortcut struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

type Shortcuts map[string]Shortcut

//go:embed index.html
var index []byte

//go:embed styles.css
var styles []byte

//go:embed reset.css
var reset []byte

func loadShortcuts() (Shortcuts, error) {
	var s Shortcuts

	data, err := os.ReadFile("/data/shortcuts.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func main() {
	shortcuts, err := loadShortcuts()
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
			if shortcut, ok := shortcuts[split[len(split)-1][1:]]; ok {
				t, err := template.New("").Parse(shortcut.Query)
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
	})

	fmt.Println("serving", os.Args[1])
	if err = http.ListenAndServe(os.Args[1], mux); err != nil {
		panic(err)
	}
}
