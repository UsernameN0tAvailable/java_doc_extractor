package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"encoding/json"
	"errors"
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
	classes []Scope
	activeScopes []*Scope
	activeScope *Scope 
}

func (e*Extractor) GetScopees() []Scope {
	return e.classes
}

func (e *Extractor) Extract(rootArg string) []Scope {

	root, err := filepath.Abs(rootArg)

	if err != nil {
		fmt.Println(err)
		panic("no file")	
	}

	e.listDirs(root)

	e.SecondaryPackageMatches()

	e.evaluate()
	os.Exit(3)

	return e.classes
}


func (e *Extractor) SecondaryPackageMatches() {

	// match classes
	for bi,class := range e.classes {

		super := class.GetSuper()

		if len(super) > 0 && len(strings.Split(super, ".")) == 1 {

			pack, err := class.GetPackage()

			if err == nil {
				for si,superScope := range e.classes {
					if bi != si && superScope.IsClass() && superScope.IsInPackage(pack) {
						newName := pack + "." + super
						if superScope.GetName() == newName {	
							e.classes[bi].SetSuper(newName)	
						}

					} 

				} 
			}
		}
	}

	// match interfaces 
	for bi,class := range e.classes {

		interfaces := class.GetInterfaces()

		for ii, inter := range interfaces {

			if len(inter) > 0 && len(strings.Split(inter, ".")) == 1 {
				pack, err := class.GetPackage()

				if err == nil {
					for _, superScope := range e.classes {
						if superScope.IsInterface() && superScope.IsInPackage(pack) {
							newName := pack + "." + inter 
							if superScope.IsClass() && superScope.GetName() == newName {	
								fmt.Println("this")
								e.classes[bi].SetInterface(newName, ii)	
							}

						} 

					} 
				}
			}
		}
	}
}



func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("too few arguments!!")
		return
	}

	extractor := Extractor{classes: make([]Scope, 0, 20000), activeScopes: make([]*Scope, 0, 200), activeScope: nil}

	classes := extractor.Extract(os.Args[1])

	jsonOut, err := json.MarshalIndent(classes, "", "\t")

	if err == nil {
		fmt.Println(string(jsonOut))
	} else {
		fmt.Println("error", err)
	}

}


func (e*Extractor) evaluate() {

	notFoundCount := 0
	foundCount := 0
	withSuper := 0 

	fmt.Println("tot",len(e.classes))


	for _,class := range e.classes {

		super := class.GetSuper()

		found := false

		if len(super) > 0 {

			withSuper++

			if true {

				for _,superScope := range e.classes {

					if superScope.GetName() == super {
						found = true
						break
					}
				}

				if found {
					foundCount++
				} else {
					notFoundCount++
				}
			} else {
				//				fmt.Println(super)
			}
		}
	}

	fmt.Println("not found:", notFoundCount,"found:", foundCount,"with super:", withSuper)

}

func inProject(path string,projectName string) bool {

	split := strings.Split(path, ".")


	isMatch := false

	for _,s := range split {
		if s == projectName {
			isMatch = true
			break
		}
	}


	return isMatch
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
		} else if file.IsDir() && file.Name() != "build" {
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

}

func (e* Extractor) parseFile(content []byte, path string) {


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
				signature = string(findFirstSignature(i, content, lastElementEnd))	
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

			scopeCount++
			isContainerScope := false
			if isValidSignature(signature) {	
				isContainerScope = e.storeSignature(signature, doc, path, &imports) 	
			}

			if isContainerScope {
				active := &e.classes[len(e.classes) - 1]
				e.activeScopes = append(e.activeScopes, active)
				e.activeScope = active
			}  else {
				e.activeScopes = append(e.activeScopes, nil)
			} 
			lastElementEnd = i

		} else if c == scopeOff && !inComment && !inString && !inChar {
			scopeCount--
			lastElementEnd = i
			doc = ""

			e.activeScopes = e.activeScopes[:(len(e.activeScopes) - 1)]
			if len(e.activeScopes) > 0 {
				e.activeScope = e.activeScopes[len(e.activeScopes) - 1]
				active := e.activeScopes[len(e.activeScopes) - 1]

				if active == nil {
					// find last used class
					// because inner class could be inside
					// of method
					for i := len(e.activeScopes) - 1; i >= 0; i -- {
						if e.activeScopes[i] != nil {
							active = e.activeScopes[i]
							break
						}
					}

					e.activeScope = active
				}
				//fmt.Println(active.GetName(), scopeCount, len(e.activeScopes), e.activeScopes[0] == nil)
			} else {
				e.activeScope = nil
			}	

		} else if c == str && !inChar && !inComment {
			inString = !inString
		} else if c == newLine && inInlineComment && !inString && !inChar  {
			inComment = false
			inInlineComment = false
			lastElementEnd = i
		} else if c == char && !inString && !inComment {
			inChar = !inChar
		}

	}
}


func (e*Extractor) storeSignature(s string, doc string, path string, imports *Imports) bool {

	isContainerScope := false
	fields := strings.Fields(s)

	for _, f := range fields {
		fT := strings.TrimSpace(f)
		if fT == "class" || fT == "enum" || fT == "record" || fT == "interface" {
			isContainerScope = true
			break 
		}
	}

	var pathIn string

	p := strings.Split(path, "/org/")

	if len(p) < 2 {
		pathIn = path
	} else {
		pathIn = "org/" + p[len(p) - 1] 
	}

	if isContainerScope {
		e.classes = append(e.classes, NewScope(path, s, doc, pathIn, imports, e.activeScope))	
	} else	{
		e.activeScope.AppendMethod(NewMethod(s, doc))
	}

	return isContainerScope 
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
	return predicate != "for" && predicate != "if" && predicate != "while" && predicate != "else" && predicate != "try" && predicate != "catch" && predicate != "finally" && predicate != "->" && predicate != "switch" && predicate != "new" && predicate != "&&" && predicate != "||" && predicate != "==" && predicate != "!=" && predicate != "synchronized" && predicate != "="
}


func findFirstSignature(i int, content []byte, lastElementEnd int) []byte {

	end := i

	for true {
		if i == lastElementEnd + 1 || ( content[i] == slash || content[i] == semiColumn) {

			sig := strings.Split(strings.TrimSpace(string(content[i:end])), "\n")

			startI := 0
			if len(sig) > 0 {
				for i := 0; i< len(sig); i++ {
					if len(sig[i]) == 0 || sig[i][0] == at {
						startI = i + 1 
					}
				}
			}

			tmp := strings.Join(sig[startI:], "\n")

			if len(tmp) > 0 && tmp[0] == slash {
				spltTmp := strings.Split(tmp, "\n")
				if len(spltTmp) > 0 {
					return []byte(spltTmp[len(spltTmp) - 1])
				} 
				return []byte(tmp[1:])
			} else {
				return []byte(tmp)
			}

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

	return Imports{imports: imports, packages: packages} 
}

func (i*Imports) GetPackage() (string, error) {
	if len(i.packages) > 0 {
		return i.packages[0], nil
	} else {
		return "", errors.New("No Package")
	}
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

