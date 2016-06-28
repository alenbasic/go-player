package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"syscall"
)

// GLOBALS DECLARED HERE

var root = "/path/to/media/"
var movies = MovieList{}
var player = Player{}

// THE MODEL CODE IS HERE

type Movie struct {
	FullFilePath string
	FileName     string
}

type MovieList struct {
	FileInfo    []Movie
	CurrentFilm string
}

var ext = [][]byte{{'.', 'm', 'k', 'v'},
	{'.', 'm', 'p', 'g'},
	{'.', 'a', 'v', 'i'},
	{'.', 'A', 'V', 'I'},
	{'.', 'm', '4', 'v'},
	{'.', 'm', 'p', '4'}}

// PLAYER OBJECT STRUCT AND METHODS

type Player struct {
	Playing bool
	Paused  string
	Name    string
	Film    *exec.Cmd
}

func (p *Player) StartFilm(name string) {
	p.Name = name
	movies.CurrentFilm = p.Name
	p.Paused = "Pause"
	p.Playing = true
	p.Film = exec.Command("omxplayer", "-o", "hdmi", name)
	p.Film.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	p.Film.Start()
}

func (p *Player) PauseFilm() {
	if p.Paused == "Pause" {
		p.Paused = "Play"
	} else {
		p.Paused = "Pause"
	}
}

func (p *Player) EndFilm() {
	pgid, err := syscall.Getpgid(p.Film.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, 15) // note the minus sign
	}
	p.Name = ""
	movies.CurrentFilm = ""
	p.Playing = false
}

// go back here https://gobyexample.com/spawning-processes

func (p *Player) SendCommandToFilm(command string) {
	if command == "pause" {
		p.PauseFilm()
	} else {
		// to be implemented
		fmt.Println("implement me!")
	}
}

// LOOKS FOR FILES ON THE FILESYSTEM

func visit(path string, f os.FileInfo, err error) error {
	bpath := []byte(path)
	bpath = bpath[len(bpath)-4:]
	for i := 0; i < len(ext); i++ {
		if reflect.DeepEqual(bpath, ext[i]) {
			movie := Movie{path, f.Name()}
			movies.FileInfo = append(movies.FileInfo, movie)
		}
	}

	return nil
}

func getFiles() {
	filepath.Walk(root, visit)
	fmt.Printf("file import complete: %d files imported\n", len(movies.FileInfo))
}

// THE VIEW CODE IS HERE

var templates_dir = "./templates/"

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles(templates_dir+"base.html", templates_dir+tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tmpl == "index.html" {
		err = t.Execute(w, movies)
	} else if tmpl == "movie.html" {
		err = t.Execute(w, player)
	} else if tmpl == "alreadyplaying.html" {
		err = t.Execute(w, movies)
	} else {
		err = t.Execute(w, nil)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HANDLERS ARE HERE

func indexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	method := r.Method
	if method == "GET" {
		if len(player.Name) != 0 && player.Name == movies.CurrentFilm {
			renderTemplate(w, "alreadyplaying.html")
		} else {
			renderTemplate(w, "index.html")
		}
	} else if method == "POST" {
		movies = MovieList{}
		if len(player.Name) != 0 {
			movies.CurrentFilm = player.Name
		}
		getFiles()
		renderTemplate(w, "index.html")
	}

}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.html")
}

func movieHandler(w http.ResponseWriter, r *http.Request) {
	command := r.URL.Query().Get("command")
	film := r.URL.Query().Get("movie")
	if player.Playing == false {
		player.StartFilm(film)
	} else if player.Playing == true && len(command) == 0 {
		if player.Name != film {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		if len(command) != 0 {
			if command == "kill" {
				player.EndFilm()
				http.Redirect(w, r, "/", http.StatusFound)
				return
			} else {
				player.SendCommandToFilm(command)
			}
		}
	}

	renderTemplate(w, "movie.html")
}

// IT ALL STARTS HERE

func main() {
	getFiles()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/movie", movieHandler)
	http.ListenAndServe(":8080", nil)
}
