package main

import "strings"

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
	templateStart = byte('<')
	templateEnd = byte('>')
	equal = byte('=')
) 

const (
	Nothing = 0
	EnterComment = 1
	LeaveComment = 2
	EnterDocumentation = 3
	LeaveDocumentation = 4
	EnterScope = 5
	LeaveScope = 6
	EnterString = 7
	LeaveString = 8
	StartEscape = 9
	LeaveEscape = 10
	EnterParamsScope = 11
	LeaveParamsScope = 12
	EnterChar = 13
	Leavechar = 14
	LeaveInlineComment = 15
	EnterTemplate = 16
	LeaveTemplate = 17
	EnterJson = 18
	LeaveJson = 19
	EnterMultilineComment = 20
	LeaveMultilineComment = 21
)


type Parser struct {
	InComment bool
	InInlineComment bool
	InDocumentation bool
	InJson bool

	InString bool
	InChar bool
	Escape bool
	ParamScopeCount int
	TemplateScopeCount int
	ScopeCount int

}

func NewParser() Parser {
	return Parser{InComment: false, InInlineComment: false, InDocumentation: false, InString: false, InChar: false, Escape: false, ParamScopeCount: 0, ScopeCount: 0, TemplateScopeCount: 0}
}


// returns event
func (p *Parser) Parse(content []byte, index int) int {

	c := content[index]

	nextIndex := index + 1
	nextNextIndex := nextIndex + 1
	prevIndex := index - 1

	if c == slash && p.IsNotInStringLike() {

		if !p.InComment && nextIndex < len(content) && star == content[nextIndex] {
			p.InComment = true

			if nextNextIndex < len(content) && star == content[nextNextIndex] {
				p.InDocumentation = true
				return EnterDocumentation
			}
			return EnterMultilineComment 
		} else if !p.InComment && !p.InInlineComment && nextIndex < len(content) && slash == content[nextIndex] {
			p.InComment = true
			p.InInlineComment = true
			return EnterComment
		} else if p.InComment && !p.InInlineComment && prevIndex >= 0 && content[prevIndex] == star {
			p.InComment = false
			if p.InDocumentation {
				p.InDocumentation = false
				return LeaveDocumentation
			}

			return LeaveMultilineComment
		}
	} else if c == scopeOn && p.IsInEmptyBody() {
		p.ScopeCount++
		return EnterScope
	} else if c == scopeOff && p.IsInEmptyBody() {
		p.ScopeCount--
		return LeaveScope
	} else if c == str && p.CanSwitchString() {

		if len(content) > nextIndex && len(content) > nextNextIndex && content[nextIndex] == str && content[nextNextIndex] == str {
			if p.InJson {
				p.InJson = false
				return LeaveJson
			} else {
				p.InJson = true
				return EnterJson
			}
		}

		if p.InString {
			p.InString = false 
			return LeaveString 
		} else {
			p.InString = true
			return EnterString
		} 
	} else if c == newLine && p.InInlineComment && !p.InString {
		p.TemplateScopeCount = 0
		p.InComment = false
		p.InInlineComment = false
		return LeaveInlineComment
	} else if c == char && p.CanSwitchChar() {

		if p.InChar {
			p.InChar = false
			return Leavechar
		} else {
			p.InChar = true
			return EnterChar
		}
	} else if c == backSlash && !p.Escape && (p.InChar || p.InString) {
		p.Escape = true
		return StartEscape
	} else if p.Escape && (p.InString || p.InChar) {
		p.Escape = false
		return LeaveEscape
	} else if c == roundOpen && p.IsInLogic() {
	        p.ParamScopeCount++	
		return EnterParamsScope
	} else if c == roundClose && p.IsInLogic() {
		p.ParamScopeCount--
		return LeaveParamsScope
	} else if c == templateStart && !p.InConditional(content, index) && p.IsInLogic() {
		p.TemplateScopeCount++
		return EnterTemplate
	} else if c == templateEnd && !p.InConditional(content, index) && p.IsInLogic() {

		if p.TemplateScopeCount > 0 {
		p.TemplateScopeCount--
	}

		return LeaveTemplate
	}

	return Nothing
}

func (p * Parser) InConditional(content []byte, index int) bool {
	paramScopeCount := p.ParamScopeCount
	nextIndex := index + 1

	if len(content) > nextIndex && (content[nextIndex] == equal) {
		return true
	}

	if paramScopeCount > 0 {
		line := getLine(content, index)
		fields := strings.Fields(line)

		for _, f := range fields {
			if f == "if" || f == "while" || f == "for" {
				return true
			}
		}
	}
	return false
}

func getLine(content []byte, index int) string {
	for i := index; i >=0; i-- {
		if content[i] == newLine {
			return string(content[i+1:index+1])
		}
	}
	return ""
}

func (p * Parser) IsInEmptyBody()  bool {
	return p.IsInLogic() && p.ParamScopeCount == 0
}

func (p *Parser) IsInLogic() bool {
	return p.IsNotInStringLike() && !p.InComment && !p.InDocumentation
}

func (p *Parser) CanSwitchString() bool {
	return !p.InChar && !p.InComment && !p.Escape && !p.InJson
}

func (p *Parser) CanSwitchChar() bool {
	return !p.InString && !p.InComment && !p.Escape
}

func (p *Parser) IsNotInStringLike() bool {
	return !p.InString && !p.InChar && !p.InJson
}
