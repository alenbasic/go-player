package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// THE MODEL CODE IS HERE

type Movie struct {
	FullFilePath string
	FileName     string
}

type PageData struct {
	MovieList   []Movie
	CurrentFilm string
	Player      Player
}

// LOOKS FOR FILES ON THE FILESYSTEM

var extensionList = [][]byte{
	{'.', 'm', 'k', 'v'},
	{'.', 'm', 'p', 'g'},
	{'.', 'a', 'v', 'i'},
	{'.', 'A', 'V', 'I'},
	{'.', 'm', '4', 'v'},
	{'.', 'm', 'p', '4'}}

func visit(path string, f os.FileInfo, err error) error {
	bpath := []byte(path)
	bpath = bpath[len(bpath)-4:]
	for i := 0; i < len(extensionList); i++ {
		if reflect.DeepEqual(bpath, extensionList[i]) {
			movie := Movie{path, f.Name()}
			pageData.MovieList = append(pageData.MovieList, movie)
		}
	}
	return nil
}

func generateMovies(path string) error {
	err := filepath.Walk(path, visit)
	if err != nil {
		return err
	} else if len(pageData.MovieList) <= 0 {
		return fmt.Errorf("No media files were found in the given path: %s", path)
	}
	fmt.Printf("file import complete: %d files imported\n", len(pageData.MovieList))
	return err
}
