package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type file struct {
	name          string
	isHtml, isDir bool
}

func (f file) String() string {
	return fmt.Sprintf("[%v]\nisHtml: %v\nisDir: %v\n", f.name, f.isHtml, f.isDir)
}

func main() {
	var root string

	if len(os.Args) == 1 {
		root = "."
	} else {
		root = os.Args[1]
	}
	addTrailingSlash(&root)
	buildIndices(root)
}

func buildIndices(dir string) (err error) {
	files, err := readDir(dir)
	if err != nil {
		return
	}
	if len(files) == 1 && files[0].isHtml {
		return
	}
	index := fmt.Sprintf("%sindex.html", dir)
	for _, f := range files {
		switch {
		case f.isHtml:
			appendHyperlinkToFile(index, f.name)
		case f.isDir:
			addTrailingSlash(&(f.name))
			var empty bool // getting a shadowed return error without this line
			if empty, err = dirIsEmpty(fmt.Sprintf("%s%s", dir, f.name)); err != nil {
				return
			} else if empty {
				continue
			}
			var files []file
			files, err = readDir(fmt.Sprintf("%s%s", dir, f.name))
			if err != nil {
				return
			}
			if len(files) == 1 && files[0].isHtml {
				appendHyperlinkToFile(index, fmt.Sprintf("%s%s", f.name, files[0].name))
				continue
			}
			appendHyperlinkToFile(index, fmt.Sprintf("%sindex.html", f.name))
			if err = buildIndices(fmt.Sprintf("%s%s", dir, f.name)); err != nil {
				return
			}
		}
	}
	return
}

func appendHyperlinkToFile(receivingFile string, fileToHyperlink string) {
	fmt.Printf("              index = %s\nhyperlink to append = %s\n\n", receivingFile, fileToHyperlink)
}

func readDir(dir string) (files []file, err error) {
	allFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, f := range allFiles {
		if f.IsDir() {
			files = append(files, file{name: f.Name(), isDir: true})
		} else {
			name := f.Name()
			if len(name) >= 4 { // the shortest possible name for a html file is ".htm"
				if strings.ToLower(name)[len(name)-4:] == ".htm" || (len(name) > 4 && strings.ToLower(name)[len(name)-5:] == ".html") {
					files = append(files, file{name: f.Name(), isHtml: true})
				}
			}
		}
	}
	return
}

func dirIsEmpty(dir string) (answer bool, err error) {
	files, err := readDir(dir)
	if err != nil {
		return
	}
	answer = len(files) == 0
	return
}

func addTrailingSlash(str *string) {
	if (*str)[len(*str)-1] != '/' {
		*str = fmt.Sprintf("%s/", *str)
	}
}
