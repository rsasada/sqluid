package parser

import (
	"fmt"
	"strings"
	"github.com/rsasada/sqluid/srcs/lexer"
)

func Parser(source string, tokens []*lexer.Token) (*Ast, error) {

	topAst := Ast{}
	cursor := 0
	
	ParsingTokens(source, )

}

func ParsingTokens (source string, tokens []*lexer.Token, cur uint) (*Ast, error) {
	if (tokens[cur].isEqual())

	else if ()
	
}

func ParseSelect()

func (t lexer.Token) isEqual(compare *lexer.Token) bool {
	return compare.value == t.value && compare.kind == t.kind
}