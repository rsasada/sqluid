package main

import (
	"fmt"

	"github.com/rsasada/sqluid/srcs/lexer"
)

// position, Cursor, Keyword, Symbol, TokenKind, Token構造体の定義があると仮定します

func printTokens(tokens []*lexer.Token) {
	for _, token := range tokens {
		fmt.Printf("Value: %s, Kind: %d)\n", token.Value, token.Kind)
	}
}

func main() {
	source := ""
	tokens, err := lexer.Lexing(source)
	if err != nil {
		fmt.Print(err)
	}
	printTokens(tokens)
}
