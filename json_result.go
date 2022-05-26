package main

type JsonProjectInfos struct {
	NumberOfClasses int `json:"numberOfClasses"`
	NumberOfInterfaces int `json:"numberOfInterfaces"`
	NumberOfAnnotations int `json:"numberOfAnnotations"`
	NumberOfRecords     int  `json:"numberOfRecords"`
	NumberOfEnums	    int  `json:"numberOfEnums"`
	TotNumberOfEntities int   `json:"totNumberOfEntities"`

	NumberOfEntitiesWithSuperClass int `json:"numberOfEntitiesWithSuperClass"`
	NumberOfEntitiesWithChildren   int `json:"numberOfEntitesWithChildren"`

	// metrics
	ANYJ float32 `json:"anyj"`
	DIR float32 `json:"dir"`
	ANYC float32 `json:"anyc"`
	WJPD float32 `json:"wjpd"`

}



type JsonResult struct {
	ProjectInfos JsonProjectInfos `json:"projectInfos"`
	Entities []Scope     `json:"entities"`
}



func NewJsonResult(entities  []Scope) JsonResult {

	totNumberOfEntities := len(entities)

	numberOfClasses := 0
	numberOfInterfaces := 0
	numberOfEnums := 0
	numberOfAnnotations := 0;
	numberOfRecords := 0

	entitiesWithSuperClass := 0
	entitesWithChildren :=  0


	// metrics
	declarationsWithJavaDoc := 0
	documentedItems := 0
	documentableItems := 0
	javaDocWords := 0
	decWithAnyComment := 0

	declarations := 0

	for _, e := range entities {
		declarationsWithJavaDoc += e.DeclatationsWithJavaDoc()
		documentedItems += e.DocumentedItems()
		documentableItems += e.DocumentableItems()
		javaDocWords += e.WordsInJavaDoc()
		decWithAnyComment += e.DecWithAnyComment()

		declarations += len(e.Methods)

		if e.IsClass() {
			numberOfClasses++
		} else if e.IsInterface() {
			numberOfInterfaces++
		} else if e.IsEnum() {
			numberOfEnums++
		} else if e.IsAnnotation() {
			numberOfAnnotations++
		} else if e.IsRecord() {
			numberOfRecords++
		}


		if e.HasChildren() {
			entitesWithChildren++
		}

		if e.HasSuper() {
			entitiesWithSuperClass++
		}
	}

	var anyj float32

	var dir float32

	var anyc float32

	var wjpd float32


	if declarations == 0 {
		anyj = 0
		anyc = 0
		wjpd = 0
	} else {
		anyj = float32(declarationsWithJavaDoc) / float32(declarations)
		anyc = float32(decWithAnyComment) / float32(declarations)
		wjpd = float32(javaDocWords) / float32(declarations)
	}

	if documentableItems == 0 {
		dir = 0
	} else if documentableItems >= documentedItems {
		dir = float32(documentedItems) / float32(documentableItems)
	} else {
		dir = 1.0 
	}

	return JsonResult{
		ProjectInfos: JsonProjectInfos{
			NumberOfClasses: numberOfClasses,
			NumberOfInterfaces: numberOfInterfaces,
			NumberOfAnnotations: numberOfAnnotations,
			NumberOfEnums: numberOfEnums,
			NumberOfRecords: numberOfRecords,
			TotNumberOfEntities: totNumberOfEntities,
			NumberOfEntitiesWithSuperClass: entitiesWithSuperClass,
			NumberOfEntitiesWithChildren: entitesWithChildren,
			ANYJ: anyj,
			DIR: dir,
			ANYC: anyc,
			WJPD: wjpd,
		},
		Entities: entities,
	}
}


