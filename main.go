package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)


const (
	scopeOn = byte('{')
	scopeOff = byte('}')

	slash = byte('/')
	star = byte('*')
)

var tot int = 0

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("too few arguments!!")
		return
	}

	root, _ := filepath.Abs(os.Args[1])

	listDirs(root)

	//fmt.Println("tot files: " + string(tot))
	fmt.Println("tot java files: " + fmt.Sprint(tot))
}

func listDirs(root string) {

	//fmt.Println("root", root)

	files, err := ioutil.ReadDir(root)

	if err != nil {
		fmt.Println(err)
		return
	}
	for fileIndex := range files {
		file := files[fileIndex]

		if ext := filepath.Ext(file.Name()); !file.IsDir() && ext == ".java" {
			parseJavaFile(root + string(os.PathSeparator) + file.Name())
		} else if file.IsDir() {
			//fmt.Println("DIR: ", file.Name())
			listDirs(root + string(os.PathSeparator) + file.Name())
		}
	}
}

func parseJavaFile(filePath string) {
	fmt.Println("PARSE FILE: " + filePath)
	tot += 1

	content, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Couldnt read file at: " + filePath)
		return
	}

	parseFile(content)

	os.Exit(3)
}

func parseFile(content []byte) {

	inScope := false
	inComment := false
	inDocumentation := false

	inSlash := false
	inStar := false
	
	for _, c := range content {

		switch c {
		case slash:
			if inSlash && inComment && !inStar {
				inSlash = false
				inComment = false
			} else if !inSlash && !inStar {
				inSlash = true
			}
		case 
		}

		//fmt.Println(inScope, inComment, )
	}

}


