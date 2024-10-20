package parser

import (
	"github.com/rsasada/sqluid/srcs/lexer"
)

func Parser(source string, tokens []*lexer.Token) (*Ast, error) {

	topAst := GeneratePipe()
	cursor := 0
	curAst := topAst

	
	for cursor < uint(len(tokens)) {
		newAst, newCursor := ParsingTokens(source, tokens, cursor)
		if err != nil {
			return nil, err	
		}
		if cursor < uint(len(tokens)) && newAst != nil {
			curAst = GeneratePipe()
			curAst.BinaryPipeNode.Left = newAst
			curAst = curAst.right
		}
		else {
			curAst = newAst
		}

	}
}


func parsingTokens(source string, tokens []*lexer.Token, curAst *Ast, cursor uint) (*Ast, uint) {

	semicolonToken := lexer.Token{
		Value: string(lexer.SemicolonSymbol),
		Kind: SymbolKind
	}
	
	if tokens[cursor].isEqual(GenerateToken(lexer.keywordKind, lexer.SelectKeyword)) {
		parseSelect(tokens, curosr, semicolonToken)
	}
	else if tokens[cursor].isEqual(GenerateToken(lexer.keywordKind, lexer.SelectKeyword)) {
		parseCreate()
	}
	else if tokens[cursor].isEqual(GenerateToken(lexer.keywordKind, lexer.SelectKeyword)) {
		parseInsert()
	}
	else {
		return nil, cursor
	}

}

func parseSelect(tokens []*lexer.Token, cursor uint, delimiter token) (*SelectNode, uint, bool) {

	newCur := cursor + 1
	selectNode := {}SelectNode
	
	exp, expCursor, ok := ParseExpressions(tokens, newCursor,
		lexer.Token{GenerateToken(lexer.KeywordKind, lexer.FromKeyword), delimiter})
	if !ok {
		return nil, cursor, false
	}

}

func parseExpressions(tokens []*lexer.Token, cursor uint, delimiter token) (*[]*Expression, uint, bool) {

	newCur := cursor
	exps := []*Expression{}

extract:
	for {
		
		if cursor >= uint(len(tokens)) {
			return nil, 
		}

		for i = 0; i < len(delimiter); i ++ {
			if tokens[newCur].isEqual(delimite[i]) == true {
				break parser
			}
		}

		if cursor != newCur {
			if !tokens[newCur].isEqual(GenerateToken(lexer.SymbolKind, lexer.CommaSymbol)) {
				return nil, cursor, false
			}
			cursor ++;
		}

		allowKind := []lexer.TokenKind{lexer.IdentifierKind, lexer.NumericKind, lexer.StringKind}
		exp, expCur, ok := extractExpression(tokens, newCur, allowKind)
		if !ok {
			return nil, cursor, false
		}
		

	}

	return &exps, newCur, true
}

func extractExpression(tokens []*lexer.Token, cursor uint, allowKinds []lexer.TokenKind) (*Expression, uint, bool) {

	if cursor >= uint(len(tokens)) {
		return nil, cursor, false 
	}

	for kind := range allowKinds {
		if tokens[cursor].Kind == kind {
			return Expression{
				Literal:	tokens[],
				Kind:		LiteralType
			}, cursor + 1, true
		}
	}

	return nil, cursor, false
}

