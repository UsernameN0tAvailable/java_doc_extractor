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
	char = byte('\'')
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
	activeClasses []Scope
	activeClass Scope 
}

func (e*Extractor) GetClasses() []Class {
	return e.classes
}

func (e*Extractor) GetInterfaces() []Interface {
	return e.interfaces
}

func (e *Extractor) Extract(rootArg string) ([]Class, []Interface) {

	root, err := filepath.Abs(rootArg)

	if err != nil {
		fmt.Println(err)
		panic("no file")	
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

	return e.classes, e.interfaces
}



func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("too few arguments!!")
		return
	}

	extractor := Extractor{classes: make([]Class, 0, 20000), interfaces: make([]Interface, 0, 10000), activeClasses: make([]Scope, 0, 200), activeClass: nil}

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

//	actualContent := removeInlineComments(content)

	//fmt.Println(string(actualContent))

	e.parseFile(content, filePath)

}
/*
func removeInlineComments(data []byte) []byte {

	inString := false


	for i,d := range data {
		if d == byte('"')  {
			inString = !inString
		} 

		if !inString && len(data)  > i + 1 && d == slash && data[d+1] == slash {

			start := d
			toCut := i

			for in := i; in < len(data); in++ {
				if data[in] != newLine {
					toCut = in + 1
				} else {
					toCut = in
					break
				}
			}

			


		}


	}

	content := string(data)





	spl := strings.Split(content, "//")

	out := make([]string, 0, len(spl))

	for _, s := range spl {
		nlspl := strings.Split(s, "\n")

		out = append(out, "\n" + strings.Join(nlspl[1:], "\n"))
	}

	return []byte(strings.Join(out, ""))
} */

func (e* Extractor) parseFile(content []byte, path string) {

	//fmt.Println(path)

	inComment := false
	inInlineComment := false
	inDocumentation := false
	inString := false
	inChar := false

	start := 0
	lastElementEnd := 0
	scopeCount := 0

	doc := ""

	imports := NewImports(content)

	//fmt.Println(imports)

	for i, c := range content {

		if c == slash && !inString && !inChar {
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

		} else if c == scopeOn && !inComment && !inString && !inChar {

			var signature string
			if scopeCount == 0 {
				signature = string(findFirstSignature(i, content))
			} else {
				signature = findSignature(i - 1, content, lastElementEnd)
			}

			sigArr := make([]string, 0, 10)

			for _,s := range strings.Split(signature, "\n") {
				if len(s) > 0 && s[0] != slash {
					sigArr = append(sigArr, s)
				}
			}


			signature = strings.Join(sigArr, "\n")
			//fmt.Println("dude", signature)

			scopeCount++
			isClass := false
			isInterface := false
			if isValidSignature(signature) {	
				isClass, isInterface = e.storeSignature(signature, doc, path, &imports) 	
			}

			if isClass {
				active := &e.classes[len(e.classes) - 1]
				e.activeClasses = append(e.activeClasses, active)
				e.activeClass = active
			} else if isInterface {
				active := &e.interfaces[len(e.interfaces) - 1]
				e.activeClasses = append(e.activeClasses, active)
				e.activeClass = active
				//active = nil
			} else {
				e.activeClasses = append(e.activeClasses, nil)
			} 
			lastElementEnd = i

		} else if c == scopeOff && !inComment && !inString && !inChar {
			scopeCount--
			lastElementEnd = i
			doc = ""

			e.activeClasses = e.activeClasses[:(len(e.activeClasses) - 1)]
			if len(e.activeClasses) > 0 {
				e.activeClass = e.activeClasses[len(e.activeClasses) - 1]
				active := e.activeClasses[len(e.activeClasses) - 1]

				if active == nil {
					// find last used class
					// because inner class could be inside
					// of method
					for i := len(e.activeClasses) - 1; i >= 0; i -- {
						if e.activeClasses[i] != nil {
							active = e.activeClasses[i]
							break
						}
					}

					e.activeClass = active
				}
				//fmt.Println(active.GetName(), scopeCount, len(e.activeClasses), e.activeClasses[0] == nil)
			} else {
				e.activeClass = nil
			}	

		} else if c == str {
			inString = !inString
		} else if c == newLine && inInlineComment && !inString && !inChar  {
			inComment = false
			inInlineComment = false
			lastElementEnd = i
		} else if c == char {
			inChar = !inChar
		}

	}
}


func (e*Extractor) storeSignature(s string, doc string, path string, imports *Imports) (bool, bool) {

	isClass := false
	isInterface := false
	fields := strings.Fields(s)

	for _, f := range fields {
		fT := strings.TrimSpace(f)
		if fT == "class" || fT == "enum" || fT == "record" {
			isClass = true
			break
		} else if fT == "interface" {
			isInterface = true
			break
		} 
	}

	var pathIn string

	p := strings.Split(path, "/java/")

	if len(p) < 2 {
		pathIn = path
	} else {
		pathIn = p[1] 
	}

	if isClass {
		e.classes = append(e.classes, NewClass(path, s, doc, pathIn, imports, e.activeClass))
	} else if isInterface {
		e.interfaces = append(e.interfaces, NewInterface(s, doc, pathIn, imports))
	} else
	{
		e.activeClass.AppendMethod(NewMethod(s, doc))
	}

	return isClass, isInterface
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
		if i == 0 || ( content[i] == slash || content[i] == semiColumn) {

			sig := strings.Split(strings.TrimSpace(string(content[i:end])), "\n")

			startI := 0
			if len(sig) > 0 {
				for i := 0; i< len(sig); i++ {
					if len(sig[i]) == 0 || sig[i][0] == at {
						startI = i + 1
					}
				}
			}
			return []byte(strings.Join(sig[startI:], "\n"))
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

	return name 
}


func (i*Imports) Print() {
	//fmt.Println(i.imports)
}








