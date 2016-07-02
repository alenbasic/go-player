package main

import (
	"html/template"
	"net/http"
)

// GLOBALS DECLARED HERE

var root = "/path/to/media"
var tmpl_dir = "./templates/"
var templates map[string]*template.Template
var pageData = PageData{}

// THE VIEW CODE IS HERE

func GenerateTemplates() {
	templates = make(map[string]*template.Template)
	modulus := template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }}
	templates_list := []string{"index.html", "about.html", "movie.html", "alreadyplaying.html"}
	for _, tmpl := range templates_list {
		t := template.New("base.html").Funcs(modulus)
		templates[tmpl] = template.Must(t.ParseFiles(tmpl_dir+"base.html", tmpl_dir+tmpl))

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
	r.ParseForm()
	tmpl := "index.html"
	if r.Method == "GET" {
		if pageData.Player.Playing {
			tmpl = "alreadyplaying.html"
		}
	} else if r.Method == "POST" {
		if pageData.Player.Playing {
			player := pageData.Player
			currentFilm := pageData.CurrentFilm
			pageData = PageData{}
			GenerateMovies()
			pageData.CurrentFilm = currentFilm
			pageData.Player = player
			tmpl = "alreadyplaying.html"
		} else {
			pageData = PageData{}
			GenerateMovies()
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
		pageData.Player.StartFilm(film)
		pageData.CurrentFilm = film
	} else if pageData.Player.Playing && (film == "" || pageData.Player.FilmName == film) {
		if len(command) != 0 {
			if command == "kill" {
				pageData.Player.EndFilm()
				http.Redirect(w, r, "/", http.StatusFound)
				return
			} else {
				pageData.Player.SendCommandToFilm(command)
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
	GenerateMovies()
	GenerateTemplates()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/movie", movieHandler)
	http.ListenAndServe(":8080", nil)
}
