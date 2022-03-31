package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("too few arguments!!")
		return
	}

	root, _ := filepath.Abs(os.Args[1])

	listDirs(root)
}

func listDirs(root string) {

	fmt.Println("root", root)

	files, err := ioutil.ReadDir(root)

	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(files)

	for fileIndex := range files {
		file := files[fileIndex]

		if ext := filepath.Ext(file.Name()); !file.IsDir() && ext == ".java" {
			parseJava(root + string(os.PathSeparator) + file.Name())
		} else {
			fmt.Println("dir", file.Name())
		}
	}
}

func parseJava(filePath string) {

}
