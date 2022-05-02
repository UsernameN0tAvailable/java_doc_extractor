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
	


	for _, e := range entities {

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
		},
		Entities: entities,
	}
}


