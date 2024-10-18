package parser

import (
	"github.com/rsasada/sqluid/srcs/lexer"
)

func returnToken()

func Parser(source string, tokens []*lexer.Token) (*Ast, error) {

	topAst := Ast{}
	cursor := 0

}

func ParsingTokens(source string, tokens []*lexer.Token, cur uint) (*Ast, error) {

}

func (t lexer.Token) isEqual(compare lexer.Token) bool {
	return compare.Value == t.Value && compare.Kind == t.Kind
}

func GenerateToken(kind lexer.tokenKind, value string) token {

}
