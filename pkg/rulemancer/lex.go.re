//go:generate re2go $INPUT -o $OUTPUT --api simple
package rulemancer

import (
	"fmt"
	"errors"
)

const (
	ScopeMain = 0 + iota
	ScopeComment
)

type scopeStack struct {
	*Engine
	debugLevel int
	stack []int
}

func newScopeStack(g *Engine) *scopeStack {
	return &scopeStack{Engine: g, stack: make([]int, 0)}
}

func (s *scopeStack) push(scope int) int {
	s.stack = append(s.stack, scope)
	if s.Debug {
		fmt.Println(purple("scope:"), s)
	}
	return scope
}

func (s *scopeStack) pop() int {
	if len(s.stack) > 0 {
		s.stack = s.stack[:len(s.stack)-1]
		scope := s.stack[len(s.stack)-1]
		if s.Debug {
			fmt.Println(purple("scope:"), s)
		}
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
		result += cyan(scopeName(s.stack[i]))
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

func (g * Engine) Compile(yyinput string) error {
	yycursor:= 0
	yytext:= 0
	yymarker:= 0
	prev:= 0
	scope:= ScopeMain
	sS := newScopeStack(g)
	sS.push(scope)

	for {
		if g.Debug && g.DebugLevel >= debugLevelMax {
			fmt.Println("----")
			fmt.Println("yycursor:", yycursor, "yytext:", yytext, "yymarker:", yymarker)
			fmt.Println("scopeStack:", sS.String())
			fmt.Println("scope:", scopeName(scope))
			fmt.Println(">>>>")
		}

		switch scope {
		case ScopeMain:
		 	/*!re2c
			re2c:yyfill:enable = 0;
			re2c:YYCTYPE = byte;
			re2c:YYPEEK = "peek(str, cur)";
			re2c:YYSKIP = "cur += 1";

			comment = "//";
			varname = [a-z][A-Za-z0-9]*;
			w = [ \t]+;

			*      { return errors.New("Unexpected input: "+string(yyinput[prev:yycursor])) }
			[\x00] { return nil }
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
