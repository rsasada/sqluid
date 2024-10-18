package parser

import (
	"github.com/rsasada/sqluid/srcs/lexer"
)

func Parser(source string, tokens []*lexer.Token) (*Ast, error) {

	topAst := GeneratePipe()
	cursor := 0
	curAst := topAst

	
	for cursor < len(tokens) {
		newAst, newCursor, err := ParsingTokens(source, tokens, cursor)
		if err != nil {
			return nil, err	
		}
		if cursor < len(tokens) && newAst != nil {
			curAst = GeneratePipe()
			curAst.BinaryPipeNode.Left = newAst
		}
		else {
			curAst = newAst
		}
		
	}

}


func ParsingTokens(source string, tokens []*lexer.Token, curAst *Ast, cursor uint) (*Ast, uint, error) {

	newCur := cursor
	semicolonToken := lexer.Token{
		Value: string(lexer.SemicolonSymbol),
		Kind: SymbolKind
	}

	if newCur > len(tokens)
		return nil,
	
	if tokens[newCur].isEqual(GenerateToken(lexer.keywordKind, lexer.SelectKeyword)) {
		parseSelect()
	}
	else if tokens[newCur].isEqual(GenerateToken(lexer.keywordKind, lexer.SelectKeyword)) {
		parseCreate()
	}
	else if tokens[newCur].isEqual(GenerateToken(lexer.keywordKind, lexer.SelectKeyword)) {
		parseInsert()
	}
	else {
		return nil, cursor, 
	}

}

func ParseSelect()

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

