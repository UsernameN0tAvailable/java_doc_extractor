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

type Extractor struct {
	classes []Class
	interfaces []Interface
	activeClasses []*Class
	activeClass *Class
}

func (e*Extractor) GetClasses() []Class {
	return e.classes
}

func (e*Extractor) GetInterfaces() []Interface {
	return e.interfaces
}

func (e *Extractor) Extract(rootArg string) {

	root, err := filepath.Abs(rootArg)

	if err != nil {
		fmt.Println(err)
		return
	}

	e.listDirs(root)

	/*
	for _, c := range e.classes {
		super := c.GetSuper()
		in := c.GetInterfaces()
		methods := c.GetMethods()
		fmt.Println("")
		fmt.Println(c.GetName())
		fmt.Println("  doc:", c.GetDocLinesCount())
		if len(super) > 0 {
			fmt.Println("  super:",super)
		}
		if len(in) > 0 {
			fmt.Println("  interfaces:", in)
		}	

		if len(methods) > 0 {
			fmt.Println("  methods")
			for _,m := range methods {
				fmt.Println("    ",m.GetDoc())
				fmt.Println("    ",m.GetSignature())
			}
		}
	} 


	for _, c := range e.interfaces {
		super := c.GetSuper()
		fmt.Println("")
		fmt.Println(c.GetName())
		fmt.Println("\tdoc:", c.GetDocLinesCount())
		if len(super) > 0 {
			fmt.Println("\tsuper:",super)
		}	
	}  */


	//	fmt.Println("\ntot classes", len(e.classes))
	//	fmt.Println("tot interfaces", len(e.interfaces))
	//	fmt.Println("tot files scanned: " + fmt.Sprint(tot))	

	// search matching super inside of project

	notFound := make([]string, 0, 10000)

	for _,class := range e.classes {

		superClass := class.GetSuper()

		if len(superClass) > 0 {


			found := false

			for _,inClass := range e.classes {

				if inClass.GetName() == superClass {
					found = true
					break
				}
			}

			if !found {

				foundInNotFound := false

				for _,fn := range notFound {
					if fn == superClass {
						foundInNotFound = true
						break
					}

				}

				isInProject := false

				for _,s := range strings.Split(superClass, ".") {
					if s == "elasticsearch" {
						isInProject = true
						break
					}
				}

				//fmt.Println(ii)

				if !foundInNotFound && isInProject {
					notFound = append(notFound, superClass)
					fmt.Println("=================================")
					fmt.Println(class.GetFullPath())
					fmt.Println(class.GetName())
					fmt.Println(superClass)
				}
			}
		}
	}

	fmt.Println("not found", len(notFound))
}



func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("too few arguments!!")
		return
	}

	extractor := Extractor{classes: make([]Class, 0, 20000), interfaces: make([]Interface, 0, 10000), activeClasses: make([]*Class, 0, 200), activeClass: nil}

	extractor.Extract(os.Args[1])
}


func (e *Extractor) listDirs(root string) {

	files, err := ioutil.ReadDir(root)

	if err != nil {
		fmt.Println(err)
		e.parseJavaFile(root)

		return
	}
	for fileIndex := range files {
		file := files[fileIndex]

		if ext := filepath.Ext(file.Name()); !file.IsDir() && ext == ".java" {
			e.parseJavaFile(root + string(os.PathSeparator) + file.Name())
		} else if file.IsDir() {
			e.listDirs(root + string(os.PathSeparator) + file.Name())
		}
	}
}

func (e* Extractor) parseJavaFile(filePath string) {
	tot += 1

	content, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Couldnt read file at: " + filePath)
		return
	}

	e.parseFile(content, filePath)

	//os.Exit(3)

}

func (e* Extractor) parseFile(content []byte, path string) {

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


			//fmt.Println(signature)

			scopeCount++
			isClass := false
			if isValidSignature(signature) {	
				isClass = e.storeSignature(signature, doc, path, &imports) 	
			} 

			if isClass {
				active := &e.classes[len(e.classes) - 1]
				e.activeClasses = append(e.activeClasses, active)
				e.activeClass = active
			} else {
				//active = nil
				e.activeClasses = append(e.activeClasses, nil)
			}

			lastElementEnd = i

		} else if c == scopeOff && !inComment && !inString {
			scopeCount--
			lastElementEnd = i
			doc = ""

			// pop class
			active := e.activeClasses[len(e.activeClasses) - 1]

			if active != nil {
				e.activeClasses = e.activeClasses[:(len(e.activeClasses) - 1)]

				// find last used class
				// because inner class could be inside
				// of method
				for i := len(e.activeClasses) - 1; i >= 0; i -- {
					if e.activeClasses[i] != nil {
						e.activeClass = e.activeClasses[i]
						e.activeClasses = e.activeClasses[:(i+1)]
						break
					}
				}
			}

		} else if c == str {
			inString = !inString
		} else if c == newLine && inInlineComment && !inString  {
			inComment = false
			inInlineComment = false
			lastElementEnd = i
		}
	}
}


func (e*Extractor) storeSignature(s string, doc string, path string, imports *Imports) bool {

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

	p := strings.Split(path, "java/")

	var pathIn string

	if len(p) < 2 {
		pathIn = path
	} else {
		pathIn = p[1] 
	}

	if isClass {
		e.classes = append(e.classes, NewClass(path, s, doc, pathIn, imports, e.activeClass != nil))
	} else if isInterface {
		e.interfaces = append(e.interfaces, NewInterface(s, doc, pathIn, imports))
	} else // method
	{
		e.activeClass.AppendMethod(NewMethod(s, doc))
	}

	return isClass
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

	//fmt.Println("first", string(content[i]))

	for true {
		if i == 0 || ( content[i] == slash || content[i] == semiColumn) {
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

	//fmt.Println("content")
	//o := string(content[lastElementEnd:i])

	for true {

		//	fmt.Println(o)

		if i == lastElementEnd ||( bracketScopeCount == 0 && (content[i] == scopeOff || content[i] == slash || i == lastElementEnd || content[i] == semiColumn || content[i] == scopeOn)) {

			var s string

			if i < end {
				s = strings.TrimSpace(string(content[i+1:end]))
			} else {
				s = ""
			}


			splt := strings.Split(s, "\n")
			//fmt.Println(splt)
			//if len(splt) <= 1 {return ""}

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
		splt := strings.Split(strings.TrimSpace(line), " ")
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

	spltName := strings.Split(name, ".")

	outPath := make([]string, 0, 100)


	for _,imp := range i.imports {
		spltImport := strings.Split(imp, ".")

		nameIndex := 0

		for ni, nameChunk := range spltName {
			matchCount := 0

			matching := false

			for importChunkIndex, importChunk := range spltImport {

				if !matching && importChunk == nameChunk {
					matching = true
					for a := 0; a < importChunkIndex; a++ {
						outPath = append(outPath, spltImport[a])
						//fmt.Println(spltImport[a])
					}

					outPath = append(outPath, importChunk)
					matchCount++
				} else if ni + matchCount >= len(spltName) || ( matching && importChunk !=  spltName[ni + matchCount] ){
					matching = false
					outPath = make([]string, 0, 100)
					break
				} else if matching {
					outPath = append(outPath, importChunk)
					nameIndex = ni
					matchCount++
				}
			}
			if matching {
				for i := nameIndex + matchCount; i < len(spltName); i++  {
					outPath = append(outPath, spltName[i])
				}

				return strings.Join(outPath, ".")
			} else {
				outPath = make([]string, 0, 100)
				break
			}
		}


	}

	return i.packages[0] + "." + name
}


func (i*Imports) Print() {
	//fmt.Println(i.imports)
}








