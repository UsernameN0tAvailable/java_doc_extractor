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

var classes []Class = make([]Class, 0)
var interfaces []string = make([]string, 0)


func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("too few arguments!!")
		return
	}

	root, _ := filepath.Abs(os.Args[1])

	listDirs(root)


	for _, c := range classes {
		super := c.GetSuper()
		in := c.GetInterfaces()
		fmt.Println("")
		fmt.Println(c.GetName())
		fmt.Println("\tdoc:", c.GetDocLinesCount())
		if len(super) > 0 {
			fmt.Println("\tsuper:",super)
		}
		if len(in) > 0 {
			fmt.Println("\tinterfaces:", in)
		}	
	} 

	fmt.Println("\ntot classes", len(classes))
	fmt.Println("\ntot interfaces", len(interfaces))
	fmt.Println("\ntot files scanned: " + fmt.Sprint(tot))	

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
	tot += 1

	content, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Couldnt read file at: " + filePath)
		return
	}

	parseFile(content, filePath)

	//os.Exit(3)

}

func parseFile(content []byte, path string) {

	//fmt.Println(path)

	inComment := false
	inInlineComment := false
	inDocumentation := false
	inString := false

	start := 0
	lastElementEnd := 0
	scopeCount := 0

	doc := ""

	imports := NewImports(content)

	//fmt.Println(imports)

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
					doc = string(content[start:nextIndex])
					inDocumentation = false
				} else if inComment {
					doc = ""
				}
				inComment = false
				lastElementEnd = i
			}

		} else if c == scopeOn && !inComment && !inString {

			var signature string
			if scopeCount == 0 {
				signature = string(findFirstSignature(i, content))
			} else {
				signature = findSignature(i - 1, content, lastElementEnd)
			}
			scopeCount++
			if isValidSignature(signature) {	
				//fmt.Println(doc)
				//fmt.Println(signature)
				storeSignature(signature, doc, path, &imports)
			} 

			lastElementEnd = i

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


func storeSignature(s string, doc string, path string, imports *Imports) {

	isClass := false
	isInterface := false
	fields := strings.Fields(s)

	for _, f := range fields {
		fT := strings.TrimSpace(f)
		if fT == "class" {
			isClass = true
			break
		} else if fT == "interface" {
			isInterface = true
			break
		}
	}

	if isClass {
		classes = append(classes, NewClass(s, doc, strings.Split(path, "java/")[1], imports))
	} else if isInterface {
		interfaces = append(interfaces, s)
		//fmt.Print(s)
	}


	//fmt.Println(fields)


	//os.Exit(3)
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
		if i == 0 || ( content[i] == newLine || content[i] == semiColumn) {
			return content[i:end]
		} else if i >= 1 {
			i--
		}
	}

	return nil
}


func findSignature(i int, content []byte, lastElementEnd int) string {

	end := i

	bracketScopeCount := 0 

	for true {
		if i == lastElementEnd ||( bracketScopeCount == 0 && (content[i] == scopeOff || content[i] == slash || i == lastElementEnd || content[i] == semiColumn || content[i] == scopeOn)) {

			var s string

			if i < end {
				s = strings.TrimSpace(string(content[i+1:end]))
			} else {
				s = ""
			}


			splt := strings.Split(s, "\n")
			if len(splt) <= 1 {return ""}

			startIndex := 0

			for i, c := range splt {
				strBef := strings.TrimSpace(string(c))
				if len(strBef) > 0 && strBef[0] != at {
					startIndex = i
					break
				}
			}
			out := strings.TrimSpace( strings.Join(splt[startIndex:], ""))

			return  out
		} else if i >= 1 {
			i--
		}

		if content[i] == roundClose {
			bracketScopeCount++
		} else if content[i] == roundOpen {
			bracketScopeCount--
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

type Imports struct {
	imports []string
	packages []string
}


func NewImports(c []byte) Imports {
	content := string(c)

	imports := make([]string, 0)

	packages := make([]string, 0)


	lines := strings.Split(content, "\n")

	for _,line := range lines {
		splt := strings.Split(line, " ")

		if len(splt) > 0 && splt[0] == "import" {
			imports = append(imports, strings.Split(splt[len(splt) -1], ";")[0])
		}
		if len(splt) > 0 && splt[0] == "package" {
			packages = append(packages, strings.Split(splt[len(splt) -1], ";")[0])
		}
	}

	if len(packages) > 1 {
		fmt.Println("impossible")
	}

	//fmt.Println(len(packages))

	return Imports{imports: imports, packages: packages} 
}


func (i*Imports) GetPath(name string) string {

	for _,i := range i.imports {
		splt := strings.Split(i, ".")

		if splt[len(splt) - 1] == name {
			return i
		}
	}

	return i.packages[0] + "." + name
}

func (i*Imports) Print() {
	//fmt.Println(i.imports)
}








