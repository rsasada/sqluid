package parser

import (
	"github.com/rsasada/sqluid/srcs/lexer"
)

func (t lexer.Token) isEqual(compare lexer.Token) bool {

	return compare.Value == t.Value && compare.Kind == t.Kind
}

func GenerateToken(kind lexer.TokenKind, value string) lexer.Token {

	return lexer.Token{
		Value: value,
		Kind: kind,
	}
}

func GeneratePipe() *Ast {

	return &Ast{
		NodeType:	BinaryPipeType,
		Pipe:		&BinaryPipeNoder{},
	}
}
