package parser

import (
	"github.com/rsasada/sqluid/srcs/lexer"
)

func Parser(source string, tokens []*lexer.Token) (*Ast, error) {

	topAst := GeneratePipe()
	cursor := uint(0)
	curAst := topAst

	for cursor < uint(len(tokens)) {
		newAst, newCursor := parsingTokens(source, tokens, cursor)
		if err != nil {
			return nil, err
		}
		if newCursor < uint(len(tokens)) && newAst != nil {
			curAst = GeneratePipe()
			curAst.Pipe.Left = newAst
			curAst = curAst.Pipe.Right
		} else {
			curAst = newAst
		}
		cursor = newCursor
	}
}

func parsingTokens(source string, tokens []*lexer.Token, curAst *Ast, cursor uint) (*Ast, uint) {

	semicolonToken := lexer.Token{
		Value: string(lexer.SemicolonSymbol),
		Kind:  lexer.SymbolKind,
	}

	if tokens[cursor].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.SelectKeyword))) {
		parseSelect(tokens, cursor, semicolonToken)
	} else if tokens[cursor].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.SelectKeyword))) {
		parseCreate()
	} else if tokens[cursor].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.SelectKeyword))) {
		parseInsert()
	} else {
		return nil, cursor
	}

}

func parseSelect(tokens []*lexer.Token, cursor uint, delimiter lexer.Token) (*SelectNode, uint, bool) {

	newCur := cursor + 1
	selectNode := SelectNode{}

	exp, expCursor, ok := parseExpressions(tokens, newCur,
		lexer.Token{GenerateToken(lexer.KeywordKind, lexer.FromKeyword), delimiter})
	if !ok {
		return nil, cursor, false
	}
	newCur = expCursor
	selectNode.Item = exp

}

func parseExpressions(tokens []*lexer.Token, cursor uint, delimiter []lexer.Token) (*[]*Expression, uint, bool) {

	newCur := cursor
	exps := []*Expression{}

extract:
	for {

		if cursor >= uint(len(tokens)) {
			return nil, cursor, false
		}

		for i := 0; i < len(delimiter); i++ {
			if tokens[newCur].IsEqual(delimiter[i]) == true {
				break extract
			}
		}

		if cursor != newCur {
			if !tokens[newCur].IsEqual(GenerateToken(lexer.SymbolKind, string(lexer.CommaSymbol))) {
				return nil, cursor, false
			}
			cursor++
		}

		allowKind := []lexer.TokenKind{lexer.IdentifierKind, lexer.NumericKind, lexer.StringKind}
		exp, expCur, ok := extractExpression(tokens, newCur, allowKind)
		if !ok {
			return nil, cursor, false
		}

		exps = append(exps, exp)
		newCur = expCur
	}

	return &exps, newCur, true
}

func extractExpression(tokens []*lexer.Token, cursor uint, allowKinds []lexer.TokenKind) (*Expression, uint, bool) {

	if cursor >= uint(len(tokens)) {
		return nil, cursor, false
	}

	for _, kind := range allowKinds {
		if tokens[cursor].Kind == kind {
			return &Expression{
				Literal: tokens[cursor],
				Kind:    LiteralType,
			}, cursor + 1, true
		}
	}

	return nil, cursor, false
}
