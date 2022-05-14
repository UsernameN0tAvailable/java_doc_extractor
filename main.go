package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)


var tot int = 0

var basePath string 

type Extractor struct {
	mu sync.Mutex
	classes []Scope
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

	var wg sync.WaitGroup

	for ci, class := range e.classes {
		for ui, usedByClass := range e.classes {
			wg.Add(1)
			go e.ClassUsesClass(class, usedByClass, ci, ui, &wg)
		} 
	}

	wg.Wait()

	//panic("done")
}

func (e*Extractor) ClassUsesClass(class Scope, usedByClass Scope, ci int, ui int, wg *sync.WaitGroup) {
	defer wg.Done()
	if usedByClass.GetName() != class.GetName() && usedByClass.UsesClass(&class) {

		e.mu.Lock()	
		defer e.mu.Unlock()

		// match tests and benchmarks
		if !class.IsATest() && usedByClass.IsATest() {	

			e.classes[ci].AppendTestCase(usedByClass.GetName())
		} else {

			e.classes[ci].AppendUsedBy(usedByClass.GetName())
		}

		e.classes[ui].AppendUses(class.GetName())
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

	jsonResult := NewJsonResult(classes)

	jsonOut, err := json.MarshalIndent(jsonResult, "", "\t")

	if err == nil {
		fmt.Println(string(jsonOut))
	} else {
		fmt.Println("error", err)
	}

}

func NewExtractor() Extractor {
	return Extractor{classes: make([]Scope, 0, 20000)}
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
	//clean := removeComment(content)

	e.parseFile(content, filePath)

}

func (e* Extractor) parseFile(content []byte, path string) {

	activeScopes := make([]*Scope, 0, 100)
	var activeScope *Scope 


	start := 0
	lastElementEnd := 0

	scopeStarts := make([]int, 0, 100)

	doc := ""

	imports, err := NewImports(content)

	if err != nil {
		return
	}

	parser := NewParser()

	//	fmt.Println(path)


	for i, _ := range content {

		nextIndex := i + 1
		event := parser.Parse(content, i)

		switch event {
		case EnterDocumentation:
			start = i
		case EnterComment:
			start = i
		case LeaveDocumentation:
			doc = string(content[start:nextIndex])
			lastElementEnd = i
			case LeaveMultilineComment: 
			doc = ""
			lastElementEnd = i
			case SemiColumn: // catch interface methods
			if activeScope != nil && activeScope.IsInterface() {

				signatureStart, signature := findSignature(i, content, lastElementEnd + 1) 	
				signature = string(removeComment([]byte(signature)))


				sigArr := make([]string, 0, 10)

				for _,s := range strings.Split(signature, "\n") {
					if len(s) > 0 && s[0] != slash {
						sigArr = append(sigArr, s)
					}
				}
				if isValidSignature(signature) {	
					e.storeSignature(signature, doc, path, &imports, activeScope, signatureStart, getCurrentLine(content, i)) 		
				}
			}

		case EnterScope:

			signatureStart, signature := findSignature(i, content, lastElementEnd + 1) 	
			signature = string(removeComment([]byte(signature)))


			sigArr := make([]string, 0, 10)

			for _,s := range strings.Split(signature, "\n") {
				if len(s) > 0 && s[0] != slash {
					sigArr = append(sigArr, s)
				}
			}

			isContainerScope := false
			if isValidSignature(signature) {	
				isContainerScope = e.storeSignature(signature, doc, path, &imports, activeScope, signatureStart, getCurrentLine(content, i)) 		
			} 

			if isContainerScope {
				active := &e.classes[len(e.classes) - 1]

				if activeScope != nil {
					activeScope.AddInnerClass(active.GetName())
				}

				activeScopes = append(activeScopes, active)
				activeScope = active
			}  else {
				activeScopes = append(activeScopes, nil)
			} 

			scopeStarts = append(scopeStarts, i)
			lastElementEnd = i

		case LeaveScope:

			lastElementEnd = i
			doc = ""

			if activeScopes[len(activeScopes) - 1] == nil && activeScope != nil {
				err, m := activeScope.GetLastMethod()
				if err != nil {
					//fmt.Println(string(content), activeScopes)
					//panic(err)
				} else {
					m.AddBody(string(content), i)
				}

			} 

			body := content[scopeStarts[len(scopeStarts) - 1]:i]

			if len(activeScopes) > 0 && activeScopes[len(activeScopes) - 1] != nil {
				activeScopes[len(activeScopes) - 1].AddBody(RemoveTemplate(string(removeComment(body))), &imports)
			} 

			scopeStarts = scopeStarts[:len(scopeStarts) - 1]
			activeScopes = activeScopes[:(len(activeScopes) - 1)]
			if len(activeScopes) > 0 {
				activeScope = activeScopes[len(activeScopes) - 1]
				active := activeScopes[len(activeScopes) - 1]

				if active == nil {
					// find last used class
					// because inner class could be inside
					// of method
					for i := len(activeScopes) - 1; i >= 0; i -- {
						if activeScopes[i] != nil {
							active = activeScopes[i]
							break
						}
					}

					activeScope = active
				}
				//fmt.Println(active.GetName(), scopeCount, len(e.activeScopes), e.activeScopes[0] == nil)
			} else {
				activeScope = nil
			}
		default:

		}
	}
}

func getCurrentLine(content []byte, index int) int {
	return len(strings.Split(string(content[:index]), "\n"))
}

func printLine(content []byte, index int) {

	chunk := content[:index]

	splt := strings.Split(string(chunk),"\n")

	if len(splt) > 0 {
		fmt.Println(len(splt), splt[len(splt) - 1])
	}
}


func removeComment(v []byte) []byte {

	start := 0

	parser := NewParser()

	for i,_ := range v  {

		nextIndex := i + 1

		switch parser.Parse(v, i) {
		case EnterJson:
			start = i
		case LeaveJson:
			firstChunk := string(v[:start])
			endChunk := string(v[nextIndex+2:])
			return removeComment([]byte(strings.Join([]string{firstChunk, endChunk}, "")))
		case EnterMultilineComment:
			start = i
		case LeaveMultilineComment:
			firstChunk := string(v[:start])
			endChunk := string(v[nextIndex+1:])
			return removeComment([]byte(strings.Join([]string{firstChunk, endChunk}, "")))
		case EnterComment:
			start = i
		case LeaveInlineComment:
			firstChunk := string(v[:start])
			endChunk := string(v[i:])
			return removeComment([]byte(strings.Join([]string{firstChunk, endChunk}, "")))

		}
	}

	return v
}


func (e*Extractor) storeSignature(s string, doc string, path string, imports *Imports, activeScope *Scope, signatureStart int, signatureLineStart int) bool {

	isContainerScope := false
	fields := strings.Fields(s)

	paramsScope := 0

	for _, f := range fields {
		fT := strings.TrimSpace(f)

		for _, c := range f {
			if byte(c) == roundOpen {
				paramsScope++
			} else if byte(c) == roundClose {
				paramsScope--
			}

		}

		if paramsScope == 0 && ( fT == "class" || fT == "enum" || fT == "record" || fT == "interface" || fT == "@interface") {
			isContainerScope = true
			break 
		}
	}

	if isContainerScope {
		e.mu.Lock()
		defer e.mu.Unlock()
		e.classes = append(e.classes, NewScope(path, s, doc, imports, activeScope))	
	} else if activeScope != nil {
		activeScope.AppendMethod(NewMethod(s, doc, signatureStart, signatureLineStart))
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

func isArrayDeclaration(s string) bool {

	content := []byte(s)

	parser := Parser{}

	for i := 0; i < len(content); i ++ {
		result := parser.Parse(content, i)
		if result == CloseSquareScope && parser.ParamScopeCount == 0 {
			return true
		}

	}

	return false
}

func isValidSignature(s string) bool {

	s = string(removeComment([]byte(s)))


	if isArrayDeclaration(s) {
		return false
	}

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
	return predicate != "for" && predicate != "if" && predicate != "while" && predicate != "else" && predicate != "try" && predicate != "catch" && predicate != "finally" && predicate != "->" && predicate != "switch" && predicate != "new" && predicate != "&&" && predicate != "||" && predicate != "==" && predicate != "!=" && predicate != "synchronized"  
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


func findSignature(i int, content []byte, lastElementEnd int) (int, string) {

	iterEnd := i

	if content[i] == semiColumn {
		iterEnd = i - 1
	}

	end := i	
	start := end

	paramsScope := 0

	for t := iterEnd; t > lastElementEnd; t-- {

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
	return start, strings.TrimSpace(string(content[start: end]))
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


	return Imports{imports: imports, pack: pack, importUses: make([]string, 0, 10)}, nil
}

func trySplit(hay string, prefixes [3]string, needle string ) []string {

	var contentSplit []string

	for _, p := range prefixes {
		contentSplit = strings.Split(hay, p + needle + "(")

		if len(contentSplit) == 1 {
			contentSplit = strings.Split(hay, p + needle + " ")
		}

		if len(contentSplit) == 1 {
			contentSplit = strings.Split(hay, p + needle + "\n")
		}

		if len(contentSplit) == 1 {
			contentSplit = strings.Split(hay, p + needle + ".")
		}

		if len(contentSplit) > 1 {
			return contentSplit
		}
	}

	return contentSplit

}


func (i*Imports) IsClassUsed(class *Scope, body string) bool {

	importSplit := strings.Split(class.GetName(), ".")
	ending := importSplit[len(importSplit) - 1]

	if class.IsAnnotation() {
		ending = "@" + ending
	} 

	contentSplit := trySplit(body, [3]string{" ", "\n", "("}, ending)

	if len(contentSplit) > 1 {
		for i := 1; i < len(contentSplit); i += 2 {
			chunk := contentSplit[i]
			if len(chunk) > 1 {	
				roundSplit := strings.Split(chunk, "(")	

				for _, m := range class.GetStaticMethods() {
					smSplit := strings.Split(m, ".")
					if smSplit[len(smSplit)- 1] == roundSplit[0] {
						return true
					}
				}


				if len(roundSplit) < 2 {
					roundSplit = strings.Split(chunk, ")")
				}

				if len(roundSplit) > 1 && len(roundSplit[0]) > 0 {	

					token := ""

					splt := strings.Fields(roundSplit[0])

					if len(splt) > 0 {
						token = splt[0]
					}

					token = strings.Split(token, ",")[0]
					token = strings.Split(token, ";")[0]
					token = strings.Split(token, "\"")[0]
					token = strings.Split(token, "..")[0]
					token = strings.Split(token, "*/")[0]
					token = strings.Split(token, "*")[0]
					token = strings.TrimSpace(token)

					return len(token) > 0
				}	
			}
		}
	}
	return false
}


func (i*Imports) IsStaticEntityImported(entity string, body string) bool {

	// in case it is a static function
	contentSplit := strings.Split(body, " " + entity + "(")

	if len(contentSplit) > 1 {
		for i := 1; i < len(contentSplit); i += 2 {
			chunk := contentSplit[i]
			if len(chunk) > 1 {	
				roundSplit := strings.Split(chunk, "(")	

				if len(roundSplit) < 2 {
					roundSplit = strings.Split(chunk, ")")
				}

				return len(roundSplit) > 1 && len(roundSplit[0]) > 0 		
			}
		}
	}
	return false
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

func (i*Imports) IsImported(searchedClass *Scope) bool {


	if searchedClass.GetPackage() == i.pack {
		return true
	}

	if searchedClass.IsPrivate {return false}

	searchedValue := searchedClass.GetName()
	if i.IsInPackage(searchedValue) {
		return true
	}

	for _, imp := range i.imports {
		if searchedValue == imp {
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

