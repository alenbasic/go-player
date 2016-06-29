package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"syscall"
)

// GLOBALS DECLARED HERE

var root = "/path/to/media"
var command_list = map[string]string{"pause": "p", "up": "\x1b[A", "down": "\x1b[B", "left": "\x1b[D", "right": "\x1b[C"}
var extension_list = [][]byte{{'.', 'm', 'k', 'v'},
	{'.', 'm', 'p', 'g'},
	{'.', 'a', 'v', 'i'},
	{'.', 'A', 'V', 'I'},
	{'.', 'm', '4', 'v'},
	{'.', 'm', 'p', '4'}}

var pageData = PageData{}
var player = Player{}

// THE MODEL CODE IS HERE

type Movie struct {
	FullFilePath string
	FileName     string
}

type PageData struct {
	MovieList   []Movie
	CurrentFilm string
}

// PLAYER OBJECT STRUCT AND METHODS

type Player struct {
	Playing  bool
	Paused   string
	FilmName string
	Film     *exec.Cmd
	PipeIn   io.WriteCloser
}

func (p *Player) StartFilm(name string) {
	p.FilmName = name
	pageData.CurrentFilm = p.FilmName
	p.Paused = "Pause"
	p.Playing = true
	p.Film = exec.Command("omxplayer", "-o", "hdmi", name)
	p.Film.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	p.PipeIn, _ = p.Film.StdinPipe()
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
	p.FilmName = ""
	pageData.CurrentFilm = ""
	p.Playing = false
}

func (p *Player) SendCommandToFilm(command string) {
	if command == "pause" {
		p.PauseFilm()
	}
	p.PipeIn.Write([]byte(command_list[command]))
}

// LOOKS FOR FILES ON THE FILESYSTEM

func visit(path string, f os.FileInfo, err error) error {
	bpath := []byte(path)
	bpath = bpath[len(bpath)-4:]
	for i := 0; i < len(extension_list); i++ {
		if reflect.DeepEqual(bpath, extension_list[i]) {
			movie := Movie{path, f.Name()}
			pageData.MovieList = append(pageData.MovieList, movie)
		}
	}

	return nil
}

func GenerateMovieList() {
	filepath.Walk(root, visit)
	fmt.Printf("file import complete: %d files imported\n", len(pageData.MovieList))
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
		err = t.Execute(w, pageData)
	} else if tmpl == "movie.html" {
		err = t.Execute(w, player)
	} else if tmpl == "alreadyplaying.html" {
		err = t.Execute(w, pageData)
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
		if len(player.FilmName) != 0 && player.FilmName == pageData.CurrentFilm {
			renderTemplate(w, "alreadyplaying.html")
		} else {
			renderTemplate(w, "index.html")
		}
	} else if method == "POST" {
		pageData = PageData{}
		if len(player.FilmName) != 0 {
			pageData.CurrentFilm = player.FilmName
		}
		GenerateMovieList()
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
		if player.FilmName != film {
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
	GenerateMovieList()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/movie", movieHandler)
	http.ListenAndServe(":8080", nil)
}
