package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// DATA STRUCTURES

type ConfigDetails struct {
	firstStart        bool
	templateDirectory string
	FilePathList      []string
	templateFileList  []string
}

type Movie struct {
	FullFilePath string
	FileName     string
	Folder       string
}

type PageData struct {
	MovieList   []Movie
	CurrentFilm string
	Player      Player
}

// GLOBALS DECLARED HERE

var extensionList = [][]byte{
	{'.', 'm', 'k', 'v'},
	{'.', 'm', 'p', 'g'},
	{'.', 'a', 'v', 'i'},
	{'.', 'm', '4', 'v'},
	{'.', 'm', 'p', '4'},
	{'.', 'd', 'a', 't'}}

var config = ConfigDetails{}
var templates map[string]*template.Template
var pageData = PageData{}
var currentPath string
var currentBase string

// FUNCTIONS FROM HERE

func visit(path string, f os.FileInfo, err error) error {
	bpath := []byte(strings.ToLower(path))
	bpath = bpath[len(bpath)-4:]
	for i := 0; i < len(extensionList); i++ {
		if reflect.DeepEqual(bpath, extensionList[i]) {
			// get folder name
			movie := Movie{path, f.Name(), getFolderName(path)}
			pageData.MovieList = append(pageData.MovieList, movie)
		}
	}
	return nil
}

func getFolderName(path string) string {
	f := filepath.Dir(path)
	if f == currentPath {
		return currentBase
	}
	return strings.Replace(f, currentPath, currentBase, 1)
}

// func getFolderNameByDepth(path string, depth int) string {
// 	f := filepath.Dir(path)
// 	sep := string(os.PathSeparator)
// 	tmp := strings.Split(f, sep)
// 	d := len(tmp)
// 	if d <= depth {
// 		return f
// 	}
// 	f = ""
// 	for i := depth; i > 0; i-- {
// 		f += tmp[d-i]
// 		if i > 1 {
// 			f += sep
// 		}
// 	}
// 	return f
// }

func generateMovies() error {
	startingCounter := len(pageData.MovieList)
	if len(config.FilePathList) > 0 {
		for index, path := range config.FilePathList {
			currentPath = filepath.Dir(path)
			currentBase = filepath.Base(currentPath)
			fmt.Println("Traversing: ", path)
			err := filepath.Walk(path, visit)
			if err != nil || len(pageData.MovieList) == startingCounter {
				config.FilePathList = append(config.FilePathList[:index], config.FilePathList[index+1:]...)
				if err == nil {
					err = fmt.Errorf("No files found.")
				}
				return err
			}
			startingCounter = len(pageData.MovieList)
		}
	}

	fmt.Printf("file import complete: %d files imported\n", len(pageData.MovieList))
	return nil
}

// THE VIEW CODE IS HERE

func generateTemplates() {
	templates = make(map[string]*template.Template)
	modulus := template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }}
	for _, tmpl := range config.templateFileList {
		t := template.New("base.html").Funcs(modulus)
		templates[tmpl] = template.Must(t.ParseFiles(config.templateDirectory+"base.html", config.templateDirectory+tmpl))
	}
}

func renderTemplate(pageStruct interface{}, w http.ResponseWriter, tmpl string) {
	var err error
	if pageStruct == nil {
		err = templates[tmpl].Execute(w, pageData)
	} else {
		err = templates[tmpl].Execute(w, pageStruct)
	}
	if err != nil {
		log.Printf("The follwing error occurred: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func refreshList() error {
	if pageData.Player.Playing {
		player := pageData.Player
		currentFilm := pageData.CurrentFilm
		pageData = PageData{}
		err := generateMovies()
		pageData.CurrentFilm = currentFilm
		pageData.Player = player
		return err
	}
	pageData = PageData{}
	err := generateMovies()
	return err

}

// HANDLERS ARE HERE

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if config.firstStart {
		http.Redirect(w, r, "/setup", http.StatusFound)
		return
	}
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	group := r.URL.Query().Get("grp")
	var tmpl string
	if pageData.Player.Playing {
		if group == "f" {
			tmpl = "alreadyplayingf.html"
		} else {
			tmpl = "alreadyplaying.html"
		}
	} else {
		if group == "f" {
			tmpl = "indexf.html"
		} else {
			tmpl = "index.html"
		}
	}
	if r.Method == "POST" {
		err := refreshList()
		if err != nil {
			panic(err)
		}
	}
	renderTemplate(nil, w, tmpl)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(nil, w, "about.html")
}

func setupHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := "setup.html"
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	if r.Method == "POST" {
		config.firstStart = false
		if _, ok := r.Form["addFilePath"]; ok {
			alreadyExists := false
			newPath := r.Form["filepath"][0]
			for _, path := range config.FilePathList {
				if path == newPath {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				config.FilePathList = append(config.FilePathList, r.Form["filepath"][0])
			}
			refreshCheck(&tmpl)
		} else if _, ok := r.Form["deleteFilePath"]; ok {
			if i, err := strconv.Atoi(r.Form["deleteFilePath"][0]); err == nil {
				config.FilePathList = append(config.FilePathList[:i], config.FilePathList[i+1:]...)
				if len(config.FilePathList) == 0 {
					config.firstStart = true
				}
				refreshCheck(&tmpl)
			} else {
				panic(err)
			}
		} else if _, ok := r.Form["saveSetup"]; ok {
			// check if there is something to save
			if len(config.FilePathList) > 0 {
				// Open file using READ & WRITE permission.
				var file, err = os.OpenFile("config", os.O_RDWR, 0644)
				if err != nil {
					panic(err)
				}
				defer file.Close()

				for _, path := range config.FilePathList {
					_, err := file.WriteString(path + "\n")
					if err != nil {
						panic(err)
					}
				}
				err = file.Sync()
				if err != nil {
					panic(err)
				}
				tmpl = "index.html"
			}
		}
	}
	renderTemplate(config, w, tmpl)
}

func refreshCheck(tmpl *string) {
	if err := refreshList(); err != nil {
		config.firstStart = true
		*tmpl = "nothingfound.html"
	}
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
	renderTemplate(nil, w, "movie.html")
}

// IT ALL STARTS HERE

func initConfigDetails() {
	config.firstStart = true
	config.templateDirectory = "./templates/"
	config.templateFileList = append(config.templateFileList,
		"index.html", "indexf.html", "about.html", "movie.html", "alreadyplaying.html", "alreadyplayingf.html", "setup.html", "nothingfound.html")
}

func initPaths() {
	file, err := os.Open("config")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		config.FilePathList = append(config.FilePathList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if len(config.FilePathList) > 0 {
		config.firstStart = false
	}
}

func main() {
	initConfigDetails()
	generateTemplates()
	initPaths()
	refreshList()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/setup", setupHandler)
	http.HandleFunc("/movie", movieHandler)
	http.ListenAndServe(":8080", nil)
}
