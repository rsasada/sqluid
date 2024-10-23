package parser

import (
	"github.com/rsasada/sqluid/srcs/lexer"
)

func GenerateToken(kind lexer.TokenKind, value string) lexer.Token {

	return lexer.Token{
		Value: value,
		Kind: kind,
	}
}

func GeneratePipe() *Ast {

	return &Ast{
		Kind:		BinaryPipeType,
		Pipe:		&BinaryPipeNode{},
	}
}
