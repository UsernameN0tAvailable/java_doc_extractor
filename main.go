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
	semiColumn = byte(';')
	roundOpen = byte('(')
	roundClose = byte(')')
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
	//fmt.Println("tot java files: " + fmt.Sprint(tot))
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
	//fmt.Println("PARSE FILE: " + filePath)
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
	lastElementEnd := 0
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
				lastElementEnd = i
			}
		} else if c == scopeOn && !inComment && !inString {

			var signature string

			// probably a class
			if scopeCount == 0 {
				signature = string(findFirstSignature(i, content))
			} else {
				signature = findSignature(i - 1, content, lastElementEnd)
			}

			scopeCount++
			if isValidSignature(signature) {
				fmt.Println(blankSpace(scopeCount) + signature)
			} else {
			//	fmt.Println("not valid")
			}


		} else if c == scopeOff && !inComment && !inString {
			scopeCount--
			lastElementEnd = i
		} else if c == str {
			inString = !inString
		} else if c == newLine && inInlineComment && !inString  {
			inComment = false
			inInlineComment = false
			lastElementEnd = i
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
	return predicate != "for" && predicate != "if" && predicate != "while" && predicate != "else" && predicate != "try" && predicate != "catch" && predicate != "finally" && predicate != "->" && predicate != "switch" && predicate != "new" && predicate != "&&" && predicate != "||" && predicate != "==" && predicate != "!=" && predicate != "synchronized"
}


func findFirstSignature(i int, content []byte) []byte {

	end := i

	for true {
		if content[i] == newLine || content[i] == semiColumn {
			return content[i:end]
		} else if i >= 1 {
			i--
		}
	}

	return nil
}


func findSignature(i int, content []byte, lastElementEnd int) string {

	end := i

	inBrackets := false

	for true {
		if !inBrackets && (content[i] == scopeOff || content[i] == slash || i == lastElementEnd || content[i] == semiColumn || content[i] == scopeOn) {
			s := strings.TrimSpace(string(content[i+1:end]))

			splt := strings.Split(s, "\n")

			startIndex := 0

			for i, c := range splt {
				if strings.TrimSpace(string(c))[0] != at {
					startIndex = i
					break
				}
			}

			return strings.TrimSpace( strings.Join(splt[startIndex:], ""))
		} else if i >= 1 {
			i--
		}

		if content[i] == roundClose {
			inBrackets = true
		} else if content[i] == roundOpen {
			inBrackets = false
		}
	}

	return "" 
}



func removeChar(s string, c string) string {
	split := strings.Split(s, c)
	return strings.Join(split, "")
}


func blankSpace(count int) string {
	out := ""
	for i := 0; i< count; i++ {
		out += "\t"
	}

	return out
}


