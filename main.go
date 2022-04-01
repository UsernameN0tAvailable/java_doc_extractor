package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)


const (
	scopeOn = byte('{')
	scopeOff = byte('}')

	slash = byte('/')
	star = byte('*')
	str = byte('"')
	newLine = byte('\n') // only works on unix systems
	tab = byte('\t')
	at = byte('@') 
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

	fmt.Println("tot files: " + string(tot))
	fmt.Println("tot java files: " + fmt.Sprint(tot))
}

func listDirs(root string) {

	files, err := ioutil.ReadDir(root)

	if err != nil {
		fmt.Println(err)
		parseJavaFile(root)

		return
	}
	for fileIndex := range files {
		file := files[fileIndex]

		if ext := filepath.Ext(file.Name()); !file.IsDir() && ext == ".java" {
			parseJavaFile(root + string(os.PathSeparator) + file.Name())
		} else if file.IsDir() {
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

	//	os.Exit(3)
}

func parseFile(content []byte) {

	inComment := false
	inInlineComment := false
	inDocumentation := false
	inString := false

	documentations := make([]string, 0)
	comments := make([]string, 0)

	start := 0

	scopeCount := 0

	for i, c := range content {

		if c == slash && !inString {
			nextIndex := i + 1
			prevIndex := i - 1
			if !inComment && nextIndex < len(content) && star == content[nextIndex] {
				inComment = true
				nextNextIndex := nextIndex + 1
				inDocumentation = nextNextIndex < len(content) && star == content[nextNextIndex]
				start = i
			} else if !inComment && !inInlineComment && nextIndex < len(content) && slash == content[nextIndex] {
				inComment = true
				inInlineComment = true
				start = i
			} else if inComment && !inInlineComment && prevIndex >= 0 && content[prevIndex] == star {

				if inDocumentation {
					documentations = append(documentations, string(content[start:nextIndex]))
					inDocumentation = false
				} else if inComment {
					comments = append(comments, string(content[start:nextIndex]))
				}
				inComment = false
			}
		} else if c == scopeOn && !inComment && !inString {

			var signature string

			// probably a class
			if scopeCount == 0 {
				signature = string(findFirstSignature(i, content))
			} else {
				signature = findSignature(i, content)
			}

			if isValidSignature(signature) {
				fmt.Println(signature)
			}

			scopeCount++

		} else if c == scopeOff && !inComment && !inString {
			scopeCount--
		} else if c == str {
			inString = !inString
		} else if c == newLine && inInlineComment && !inString  {
			inComment = false
			inInlineComment = false
		}
	}
}


func isValidSignature(s string) bool {

	trimmed := strings.TrimSpace(s)

	if len(trimmed) == 0 {
		return false
	} else if trimmed[0] == at {
		return false
	} else {
		fields := strings.Fields(trimmed)
		predicate := fields[0]
		if len(predicate) == 0 {
			return false
		}
		for _,field := range fields {
			subfields := strings.Fields(field)
			for _,f := range subfields {
				for _, f0 := range strings.Split(f, "(") {
					for _, f1 := range strings.Split(f0, ")") {	
						if !isValidSignatureKeyWord(f1) {
							return false
						}
					}

				}
			}
		}
		return true
	}
}

func isValidSignatureKeyWord(predicate string) bool {
	return predicate != "for" && predicate != "if" && predicate != "while" && predicate != "else" && predicate != "try" && predicate != "catch" && predicate != "finally" && predicate != "->" && predicate != "switch" && predicate != "new" && predicate != "&&" && predicate != "||" && predicate != "==" && predicate != "!="
}


func findFirstSignature(i int, content []byte) []byte {

	end := i

	for true {
		if content[i] == newLine {
			return content[i:end]
		} else if i >= 1 {
			i--
		}
	}

	return nil
}


func findSignature(i int, content []byte) string {

	end := i

	for true {
		if content[i] == scopeOff || content[i] == slash {
			s := strings.TrimSpace(string(content[i+1:end]))
			// remove newLines
			s = removeChar(s, string(newLine))
			// remove tab
			return removeChar(s, string(tab))
		} else if i >= 1 {
			i--
		}
	}

	return "" 
}



func removeChar(s string, c string) string {
	split := strings.Split(s, c)
	return strings.Join(split, "")
}


