package main

import (
	"html/template"
	"log"
	"net/http"
)

// GLOBALS DECLARED HERE

var mediaDir = "/path/to/media"
var tmplDir = "./templates/"
var templates map[string]*template.Template
var pageData = PageData{}

// THE VIEW CODE IS HERE

func generateTemplates() {
	templates = make(map[string]*template.Template)
	modulus := template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }}
	templatesList := []string{"index.html", "about.html", "movie.html", "alreadyplaying.html"}
	for _, tmpl := range templatesList {
		t := template.New("base.html").Funcs(modulus)
		templates[tmpl] = template.Must(t.ParseFiles(tmplDir+"base.html", tmplDir+tmpl))
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates[tmpl].Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HANDLERS ARE HERE

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	tmpl := "index.html"
	switch r.Method {
	case "GET":
		if pageData.Player.Playing {
			tmpl = "alreadyplaying.html"
		}
	case "POST":
		if pageData.Player.Playing {
			player := pageData.Player
			currentFilm := pageData.CurrentFilm
			pageData = PageData{}
			err = GenerateMovies(mediaDir)
			if err != nil {
				panic(err)
			}
			pageData.CurrentFilm = currentFilm
			pageData.Player = player
			tmpl = "alreadyplaying.html"
		} else {
			pageData = PageData{}
			err = GenerateMovies(mediaDir)
			if err != nil {
				panic(err)
			}
		}
	}
	renderTemplate(w, tmpl)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.html")
}

func movieHandler(w http.ResponseWriter, r *http.Request) {
	command := r.URL.Query().Get("command")
	film := r.URL.Query().Get("movie")

	if pageData.Player.Playing == false {
		if film == "" {
			log.Println("No film was selected")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		err := pageData.Player.StartFilm(film)
		if err != nil {
			log.Printf("Following error occurred: %v\n", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		pageData.CurrentFilm = film
	} else if pageData.Player.Playing && (film == "" || pageData.Player.FilmName == film) {
		if command == "kill" {
			err := pageData.Player.EndFilm()
			if err != nil {
				log.Printf("Following error occurred: %v\n", err)
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		} else if command != "" {
			err := pageData.Player.SendCommandToFilm(command)
			if err != nil {
				log.Printf("Following error occurred: %v\n", err)
			}
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	renderTemplate(w, "movie.html")
}

// IT ALL STARTS HERE

func main() {
	err := generateMovies(mediaDir)
	if err == nil {
		generateTemplates()
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
		http.HandleFunc("/", indexHandler)
		http.HandleFunc("/about", aboutHandler)
		http.HandleFunc("/movie", movieHandler)
		http.ListenAndServe(":8080", nil)
	} else {
		panic(err)
	}
}
