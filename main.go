package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
	//var root string
	//if len(os.Args) == 1 { // assume current dir if no arg given
	//root = "."
	//} else {
	//root = os.Args[1]
	//}

	if len(os.Args) == 1 {
		log.Fatalln("missing argument; specify directory")
	}
	root := os.Args[1]
	addTrailingSlash(&root)

	var err error
	if err = buildIndices(root); err != nil {
		log.Fatalln(err)
	}

	// get list of all html files in tree
	var htmls []file
	if htmls, err = htmlFiles(root); err != nil {
		log.Fatalln(err)
	}

	//for _, f := range htmls { // debugging
	//fmt.Println(f.name)
	//}
	//fmt.Println()

	if err = chainHtmlFiles(htmls); err != nil {
		log.Fatalln(err)
	}
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
			if err = appendHyperlinkToFile(index, removeExt(f.name), f.name); err != nil {
				return
			}
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
				if err = appendHyperlinkToFile(index, removeExt(files[0].name), fmt.Sprintf("%s%s", f.name, files[0].name)); err != nil {
					return
				}
				continue
			}
			if err = appendHyperlinkToFile(index, f.name, fmt.Sprintf("%sindex.html", f.name)); err != nil {
				return
			}
			if err = buildIndices(fmt.Sprintf("%s%s", dir, f.name)); err != nil {
				return
			}
		}
	}
	return
}

func htmlFiles(dir string) (files []file, err error) {
	files = make([]file, 0)
	htmlsAndDirs, err := readDir(dir)
	if err != nil {
		return
	}
	for _, f := range htmlsAndDirs {
		switch {
		case f.isHtml:
			files = append(files, file{name: fmt.Sprintf("%s%s", dir, f.name), isHtml: true})
		case f.isDir:
			addTrailingSlash(&(f.name))
			var subFiles []file
			if subFiles, err = htmlFiles(fmt.Sprintf("%s%s", dir, f.name)); err != nil {
				return
			}
			files = append(files, subFiles...)
		}
	}
	return
}

func chainHtmlFiles(files []file) (err error) {
	lasti := len(files) - 1
	for i := range files {
		if i > 0 {
			if err = appendHyperlinkToFile(files[i].name, "previous", relativeTo(files[i].name, files[i-1].name)); err != nil {
				return
			}
		}
		if i < lasti {
			if err = appendHyperlinkToFile(files[i].name, "next", relativeTo(files[i].name, files[i+1].name)); err != nil {
				return
			}
		}
	}
	return
}

func appendHyperlinkToFile(receivingFile string, text string, fileToHyperlink string) (err error) {
	//fmt.Printf("              index = %s\nhyperlink to append = %s\n\n", receivingFile, fileToHyperlink)
	f, err := os.OpenFile(receivingFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	_, err = f.Write([]byte(fmt.Sprintf("<a href=\"%s\">%s</a><br>\n", fileToHyperlink, text)))
	if err != nil {
		f.Close()
		return
	}
	err = f.Close()
	return
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

func removeExt(str string) string {
	return str[:strings.LastIndex(str, ".")]
}

func relativeTo(pathOfCurrent, path string) string {
	pOC := strings.Split(pathOfCurrent, "/")
	pOC = pOC[:len(pOC)-1]

	p := strings.Split(path, "/")
	file := p[len(p)-1]
	p = p[:len(p)-1]

	i := 0
	for ; i < len(pOC) && i < len(p) && pOC[i] == p[i]; i++ {
	}
	pOC = pOC[i:]
	p = p[i:]

	backUp := ""
	for i = 0; i < len(pOC); i++ {
		backUp += "../"
	}

	goForward := ""
	goForward += strings.Join(p, "/")
	if goForward != "" {
		goForward += "/"
	}

	return backUp + goForward + file
}

/*
relativeTo notes

This was a fun problem!

Given a list like this:

testFiles/Squirrelerella_Gets_Married/Ch1.html
testFiles/Squirrelerella_Gets_Married/Ch2.html
testFiles/Squirrelerella_Gets_Married/Ch3.html
testFiles/Writing_Performant_COBOL/Chapter_1/footNotes.html
testFiles/Writing_Performant_COBOL/Chapter_1/part1.html
testFiles/Writing_Performant_COBOL/Chapter_1/part2.html
testFiles/Writing_Performant_COBOL/Chapter_2/text.HTM
testFiles/Writing_Performant_COBOL/Chapter_3.html
testFiles/Writing_Performant_COBOL/Chapter_4/text.html
testFiles/Writing_Performant_COBOL/Introduction.html

how do you make relative links?
for instance, footNotes.html needs to point back to ../../Squirrelerella_Gets_Married/Ch3.html

after studying these examples:

a/x.html
a/y.html

a/b/x.html
a/b/y.html

/a/b/c/m/x.html
/a/b/d/n/y.html

a/b/c/d/x.html
a/y.html

a/x.html
a/b/c/d/y.html

i came up with this algorithm:

path to link to: /a/b/c/m/x.html
current file:    /a/b/d/n/y.html

split each string

a b c m x.html
a b d n y.html

set the files aside (keep the destination file)

a b c m
a b d n

starting from the left, remove matching folders

c m
d n

with what you have left, replace the folders from the current file (d and n) with ..'s
and then append those ..'s in front the folders from the path to link to

../../c/m/

reappend destination file

*/
