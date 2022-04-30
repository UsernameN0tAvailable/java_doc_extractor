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
	backSlash = byte('\\')
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

var basePath string 

type Extractor struct {
	classes []Scope
	activeScopes []*Scope
	activeScope *Scope 
}

func (e*Extractor) GetScopees() []Scope {
	return e.classes
}

func (e *Extractor) Extract(rootArg string) []Scope {

	splitRootPath := strings.Split(rootArg, "/")
	projectName := splitRootPath[len(splitRootPath) - 1]	
	basePath = strings.Split(rootArg, projectName)[0]

	root, err := filepath.Abs(rootArg)

	if err != nil {
		fmt.Println(err)
		panic("no file")	
	}

	e.listDirs(root)

	e.SecondaryPackageMatches()
	e.MatchUsages()

	return e.classes
}


/*
* Very rudimentary test matching
*/
func (e *Extractor) MatchUsages() {

	for ci, class := range e.classes {

		for ui, usedByClass := range e.classes {
			if usedByClass.Imports(class.GetName()) {
				// match tests and benchmarks
				if !class.IsATest() && usedByClass.IsATest() {	
					e.classes[ci].AppendTestCase(usedByClass.GetName())
				} else {
					e.classes[ci].AppendUsedBy(usedByClass.GetName())
				}
				e.classes[ui].AppendUses(class.GetName())
			}
		} 
	}
}


func (e *Extractor) SecondaryPackageMatches() {

	// match classes
	for bi,class := range e.classes {

		super := class.GetSuper()

		if len(super) > 0 && len(strings.Split(super, ".")) == 1 {
			pack := class.GetPackage()

			for si,superScope := range e.classes {
				if bi != si && superScope.IsInPackage(pack) {
					sSplit := strings.Split(superScope.GetName(), ".")
					lastExt := sSplit[len(sSplit) - 1]

					newName := pack + "." + super
					if superScope.GetName() == newName {	
						e.classes[bi].SetSuper(newName)	
						e.classes[si].AppendSubClass(class.GetName())
					} else if lastExt == super {
						e.classes[bi].SetSuper(superScope.GetName())
						e.classes[si].AppendSubClass(class.GetName())
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
				pack := class.GetPackage()

				for si, superScope := range e.classes {

					if superScope.IsInterface() && superScope.IsInPackage(pack) {
						newName := pack + "." + inter 
						if superScope.GetName() == newName {	
							e.classes[bi].SetInterface(newName, ii)	
							e.classes[si].AppendImplementedBy(class.GetName())
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

	extractor := NewExtractor()

	classes := extractor.Extract(os.Args[1])

	jsonOut, err := json.MarshalIndent(classes, "", "\t")

	if err == nil {
		fmt.Println(string(jsonOut))
	} else {
		fmt.Println("error", err)
	}

}

func NewExtractor() Extractor {
	return Extractor{classes: make([]Scope, 0, 20000), activeScopes: make([]*Scope, 0, 200), activeScope: nil}
}


func (e*Extractor) evaluate() {

	notFoundCount := 0
	foundCount := 0
	withSuper := 0 

	notFoundInterfaces := 0
	foundInterfaces := 0
	importedInterfaces := 0

	classImports := make([]string, 0, 10000)
	interfaceImports := make([]string, 0, 10000)

	fmt.Println("Tot Scopes",len(e.classes))

	for _,class := range e.classes {

		super := class.GetSuper()

		found := false

		if len(super) > 0 {

			withSuper++


			for _,superScope := range e.classes {

				if superScope.GetName() == super && superScope.IsClass() {
					found = true
					break
				}
			}

			if found {
				foundCount++
			} else {
				notFoundCount++
				classImports = addUnique(classImports, super)
			}

		}


		interfaces := class.GetInterfaces()


		for _,inter := range interfaces {

			found := false

			for _,superScope := range e.classes {

				if superScope.GetName() == inter && superScope.IsInterface() {
					found = true
					break
				}


			}

			if found {
				foundInterfaces++
			} else {
				interfaceImports = addUnique(interfaceImports, inter)
				notFoundInterfaces++
			}

			importedInterfaces++
		}
	}

	fmt.Println("Classes: Not Found(imports):", notFoundCount,"Found:", foundCount,"extends:", withSuper)
	fmt.Println("Interfaces: Not Found(imports):", notFoundInterfaces,"Found:", foundInterfaces,"implements:", importedInterfaces)

	fmt.Println("unique interface imports:", len(interfaceImports))
	fmt.Println("unique class imports:", len(classImports))
	/*
	for _,i := range interfaceImports {
		//fmt.Println(i)
	} */

}


func addUnique(vals []string, v string) []string {

	found := false

	for _,e := range vals {
		if e == v {
			found = true
			break
		}
	}

	if !found {
		return append(vals, v)
	}

	return vals
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
		} else if file.IsDir()  {
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
	clean := removeComment(content)

	e.parseFile(clean, filePath)

}

func (e* Extractor) parseFile(content []byte, path string) {

	inComment := false
	inInlineComment := false
	inDocumentation := false
	inString := false
	inChar := false
	escape := false
	paramsScope := 0

	start := 0
	lastElementEnd := 0
	scopeCount := 0

	doc := ""

	imports, err := NewImports(content)

	if err != nil {
		return
	}

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

		} else if c == scopeOn && !inComment && !inString && !inChar && paramsScope == 0 {

			var signature string
			if scopeCount == 0 {

				signature = string(findFirstSignature(i, content, lastElementEnd))	
			} else {
				signature = findSignature(i, content, lastElementEnd + 1)
				//fmt.Println(signature)
			}

			//			fmt.Println(signature)

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

		} else if c == scopeOff && !inComment && !inString && !inChar && paramsScope == 0 {
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

		} else if c == str && !inChar && !inComment && !escape {
			inString = !inString
		} else if c == newLine && inInlineComment && !inString {
			inComment = false
			inInlineComment = false
			lastElementEnd = i
		} else if c == char && !inString && !inComment && !escape {
			inChar = !inChar
		} else if c == backSlash && !escape && (inString || inChar) {
			escape = true
		} else if escape && (inString || inChar) {
			escape = false
		} else if c == roundOpen {
			paramsScope ++
		} else if c == roundClose {
			paramsScope--
		}


	}
}


func printLine(content []byte, index int) {

	chunk := content[:index]

	splt := strings.Split(string(chunk),"\n")

	if len(splt) > 0 {
		fmt.Println(len(splt), splt[len(splt) - 1])
	}
}


func removeComment(v []byte) []byte {


	isMultiline := false
	isSingleLine := false
	inString := false
	inChar := false
	isEscape := false
	isJson := false

	l := len(v)

	start := 0

	for i,b := range v  {

		nextIndex := i + 1
		if b == byte('"') && len(v) > nextIndex + 1 && v[nextIndex] == byte('"') && v[nextIndex + 1] == byte('"') && !isEscape && !inChar && !inString && !isJson {
			isJson = true
			start = i
		} else if b == byte('"') && isJson && len(v) > nextIndex + 1 && v[nextIndex] == byte('"') && v[nextIndex + 1] == byte('"') && !isEscape && !inChar && !inString {
			firstChunk := string(v[:start])
			endChunk := string(v[nextIndex+2:])
			return removeComment([]byte(strings.Join([]string{firstChunk, endChunk}, "")))

		} else if b == slash && !inString && !inChar && !isJson {
			if !isMultiline && !isSingleLine && l > nextIndex && v[nextIndex] == star && (!(len(v) > nextIndex + 1 && v[nextIndex+1] == star) || (len(v) > nextIndex + 2 && v[nextIndex+1] == star && v[nextIndex + 2] == star)){
				isMultiline = true
				start = i
			} else if !isMultiline && !isSingleLine && l > nextIndex && v[nextIndex] == slash {
				isSingleLine = true
				start = i
			}
		} else if b == star && !inString && !inChar && !isJson {
			if isMultiline && l > nextIndex && v[nextIndex] == slash {
				firstChunk := string(v[:start])
				endChunk := string(v[nextIndex+1:])
				return removeComment([]byte(strings.Join([]string{firstChunk, endChunk}, "")))
			}
		} else if b == newLine && isSingleLine && !inChar && !inString && !isJson {
			firstChunk := string(v[:start])
			endChunk := string(v[i:])
			return removeComment([]byte(strings.Join([]string{firstChunk, endChunk}, "")))
		} else if b == byte('"') && !isMultiline && !isSingleLine && !inChar && !isEscape && !isJson {
			inString = !inString
		} else if b == byte('\'') && !isMultiline && !isSingleLine && !inString && !isEscape && !isJson {
			inChar = !inChar
		} else if b == backSlash && (inString || inChar) && !isEscape && !isJson {
			isEscape = true
		} else if (inString || inChar) && isEscape {
			isEscape = false
		}

	}

	return v
}


func (e*Extractor) storeSignature(s string, doc string, path string, imports *Imports) bool {

	isContainerScope := false
	fields := strings.Fields(s)

	//fmt.Println(fields)

	paramsOpen := false

	for _, f := range fields {
		fT := strings.TrimSpace(f)

		if !paramsOpen && contains(fT, roundOpen) {
			paramsOpen = true
		} else if paramsOpen && contains(fT, roundClose) {
			paramsOpen = false
		}

		if !paramsOpen && ( fT == "class" || fT == "enum" || fT == "record" || fT == "interface") {
			isContainerScope = true
			break 
		}
	}

	if isContainerScope {
		e.classes = append(e.classes, NewScope(path, s, doc, imports, e.activeScope))	
	} else if e.activeScope != nil {
		e.activeScope.AppendMethod(NewMethod(s, doc))
	}

	return isContainerScope 
}

func contains(stack string, hay byte) bool {

	for _,c := range stack {
		if byte(c) == hay {
			return true
		}
	}

	return false
}

func isValidSignature(s string) bool {

	trimmed := strings.TrimSpace(s)

	if len(trimmed) == 0 {
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

	if lastElementEnd == 0 {
		lastElementEnd = -1
	}

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
	start := end

	paramsScope := 0

	for t := end; t > lastElementEnd; t-- {

		char := content[t]

		if char == roundClose {
			paramsScope++
		} else if char == roundOpen {
			paramsScope--
		}

		if char == semiColumn && paramsScope == 0 {
			break
		}
		start = t
	}

	return strings.TrimSpace(string(content[start: end]))
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
	importUses [] string // InnerClasses used inside of code
	pack string
}


func NewImports(c []byte) (Imports, error) {
	content := string(c)

	imports := make([]string, 0)

	pack:= ""

	lines := strings.Split(content, "\n")

	for _,line := range lines {
		splt := strings.Split(strings.TrimSpace(line), " ")
		if len(splt) > 0 { 
			if splt[0] == "import" {
				imports = append(imports, strings.Split(splt[len(splt) -1], ";")[0])
			} else if splt[0] == "package" {
				if len(pack) > 0 {
					return *new(Imports), errors.New("Two packages found!!")
				}
				pack = strings.Split(splt[len(splt) -1], ";")[0]

			}
		}
	}	

	if len(pack) == 0 {
		return *new(Imports), errors.New("No Package Found!!")
	}

	// fetch inner uses
	importUses := make([]string, 0, len(imports) * 4)

	contentWOTemplates := RemoveTemplate(content)

	for _, imp := range imports {

		importSplit := strings.Split(imp, ".")
		ending :=  importSplit[len(importSplit) - 1]

		contentSplit := strings.Split(contentWOTemplates, ending + ".")

		if len(contentSplit) > 1 {
			for i := 1; i < len(contentSplit); i += 2 {
				chunk := contentSplit[i]
				if len(chunk) > 1 {
					firstChar := string(chunk[0])

					if firstChar == strings.ToUpper(firstChar) {
						roundSplit := strings.Split(chunk, ")")	

						if len(roundSplit) < 2 {
							roundSplit = strings.Split(chunk, "(")
						}


						if len(roundSplit) > 1 && len(roundSplit[0]) > 0 {	
							token := strings.Fields(roundSplit[0])[0]
							token = strings.Split(token, ",")[0]
							token = strings.Split(token, ";")[0]
							token = strings.Split(token, "\"")[0]
							token = strings.Split(token, "..")[0]
							token = strings.Split(token, "*/")[0]
							token = strings.Split(token, "*")[0]
							token = strings.TrimSpace(token)
							if len(token) > 0 {
								importUses = append(importUses, RemoveTemplate(imp + "." + token))
							}
						}
					}
				}
			}
		}
	}

	return Imports{imports: imports, pack: pack, importUses: importUses}, nil
}

func (i*Imports) GetPackage() string {
	return i.pack

}

func (i*Imports) IsInPackage(searchedValue string) bool {
	if i.pack == searchedValue  {
		return true
	}
	return false
}

func (i*Imports) IsImported(searchedValue string) bool {


	if i.IsInPackage(searchedValue) {
		return true
	}

	for _, imp := range i.imports {
		if searchedValue == imp {
			return true
		}
	}

	for _, use := range i.importUses {
		if searchedValue == use {
			return true
		}
	}

	return false
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

