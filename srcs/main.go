package main

import (
	"fmt"

	"github.com/rsasada/sqluid/srcs/lexer"
	"github.com/rsasada/sqluid/srcs/parser"
)


func printTokens(tokens []*lexer.Token) {
	for _, token := range tokens {
		fmt.Printf("Value: %s, Kind: %d)\n", token.Value, token.Kind)
	}
}

func main() {
	source := "INSERT INTO table1 (name);"
	tokens, err := lexer.Lexing(source)
	ast, ok := parser.Parser(source, tokens)
	if !ok {
		fmt.Printf("ast failed")
	}
	parser.PrintAst(ast, 3)
	if err != nil {
		fmt.Print(err)
	}
	//printTokens(tokens)
}