//go:generate re2go $INPUT -o $OUTPUT --api simple
package game

import (
	"fmt"
	"errors"
)

const (
	ScopeMain = 0 + iota
	ScopeComment
)

type scopeStack struct {
	stack []int
}

func newScopeStack() *scopeStack {
	return &scopeStack{stack: make([]int, 0)}
}

func (s *scopeStack) push(scope int) int {
	s.stack = append(s.stack, scope)
	return scope
}

func (s *scopeStack) pop() int {
	if len(s.stack) > 0 {
		s.stack = s.stack[:len(s.stack)-1]
		scope := s.stack[len(s.stack)-1]
		return scope
	}
	return -1
}

func (s *scopeStack) descend(scope int) int {
	for {
		if s.stack[len(s.stack)-1] == scope {
			return scope
		}

		curr := s.pop()
		if curr == -1 {
			return -1
		}
	}
	return -1
}

func (s *scopeStack) String() string {
	result := "["
	for i := 0 ; i < len(s.stack); i++ {
		result += scopeName(s.stack[i])
		if i < len(s.stack)-1 {
			result += ", "
		}
	}
	result += "]"
	return result
}

func scopeName(scope int) string {
	switch scope {
	case ScopeMain:
		return "main"
	case ScopeComment:
		return "comment"
	}
	return "unknown"
}

// Returns "fake" terminating null if cursor has reached limit.
func peek(str string, cur int) byte {
	if cur >= len(str) {
		return 0 // fake null
	} else {
		return str[cur]
	}
}

type GameParser struct {
	*Config
} 

func (h * GameParser) Lex(relation string, yyinput string) error {
	yycursor:= 0
	yytext:= 0
	yymarker:= 0
	prev:= 0
	scope:= ScopeMain
	sS := newScopeStack()
	sS.push(scope)

	for {
		if h.Debug {
			fmt.Println("----")
			fmt.Println("yycursor:", yycursor, "yytext:", yytext, "yymarker:", yymarker)
			fmt.Println("scopeStack:", sS.String())
			fmt.Println("scope:", scopeName(scope))
			fmt.Println(">>>>")
		}

		switch scope {
		case ScopeMain:
			if h.Debug {
				fmt.Println("Scope Main entered")
			}
		 	/*!re2c
			re2c:yyfill:enable = 0;
			re2c:YYCTYPE = byte;
			re2c:YYPEEK = "peek(str, cur)";
			re2c:YYSKIP = "cur += 1";

			comment = ";";
			capitalname = [A-Z][A-Za-z0-9]*;
			varname = [A-Za-z0-9][A-Za-z0-9-]*;

			w = [ \t]+;

			*      { return errors.New("Unexpected input: "+string(yyinput[prev:yycursor])) }
			[\x00] { return nil }
			varname {
					scope = sS.push(ScopeMain)
					prev = yycursor
					continue
				}
			comment {
					scope = sS.push(ScopeComment)
					prev = yycursor
					continue
				}
			w 	{
					//fmt.Println("space:", yyinput[prev:yycursor])
					prev = yycursor
					continue
				}
			[\n]+ {
					//fmt.Println("newline:", yyinput[prev:yycursor])
					prev = yycursor
					continue
				}
			*/
		case ScopeComment:
			if h.Debug {
				fmt.Println("Scope Comment entered")
			}
			/*!re2c
			re2c:yyfill:enable = 0;
			re2c:YYCTYPE = byte;
			re2c:YYPEEK = "peek(str, cur)";
			re2c:YYSKIP = "cur += 1";
			
			[\x00] { return errors.New("Unexpected end of input") }
			"\n"   {
					scope = sS.pop()
					prev = yycursor
					continue
				}
			*      {
					prev = yycursor
					continue
				}
			*/
		}
	}
}
