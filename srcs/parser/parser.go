package parser

import (
	"github.com/rsasada/sqluid/srcs/lexer"
)

func Parser(source string, tokens []*lexer.Token) (*Ast, bool) {

	topAst := &Ast{}
	cursor := uint(0)
	curAst := topAst

	for cursor < uint(len(tokens)) {
		newAst, newCursor, ok := parsingTokens(source, tokens, cursor)
		if !ok {
			return nil, false
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

	return topAst, true
}

func parsingTokens(source string, tokens []*lexer.Token, cursor uint) (*Ast, uint, bool) {

	semicolonToken := lexer.Token{
		Value: string(lexer.SemicolonSymbol),
		Kind:  lexer.SymbolKind,
	}

	if tokens[cursor].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.SelectKeyword))) {
		slct, newCur, ok := parseSelect(tokens, cursor, semicolonToken)
		if !ok {
			return nil, cursor, false
		}
		return &Ast{
			Kind: SelectType,
			Select: slct,
		}, newCur, true

	} else if tokens[cursor].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.SelectKeyword))) {
		parseCreate()

	} else if tokens[cursor].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.SelectKeyword))) {
		insert, newCur, ok := parseInsert(tokens, cursor, semicolonToken)
		if !ok {
			return nil, cursor, false
		}
		return &Ast{
			Kind: InsertType,
			Insert: insert,
		}, newCur, true

	} else {
		return nil, cursor, false
	}

}

func parseSelect(tokens []*lexer.Token, cursor uint, delimiter lexer.Token) (*SelectNode, uint, bool) {

	newCur := cursor + 1
	selectNode := SelectNode{}

	exp, expCursor, ok := parseExpressions(tokens, newCur,
		[]lexer.Token{GenerateToken(lexer.KeywordKind, string(lexer.FromKeyword)), delimiter})
	if !ok {
		return nil, cursor, false
	}
	newCur = expCursor
	selectNode.Item = exp

	if tokens[newCur].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.FromKeyword))) {
		cursor ++;
		if tokens[newCur].Kind == lexer.IdentifierKind {
			selectNode.From = tokens[newCur]
		}
	}

	return &selectNode, newCur+1, true
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

func parseInsert(tokens []*lexer.Token, cursor uint, delimiter lexer.Token) (*InsertNode, uint, bool) {

	insertNode := &InsertNode{}
	newCur := cursor
	newCur ++;

	if !tokens[newCur].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.IntoKeyword))) {
		return nil, cursor, false
	}
	newCur ++;

	//parse tablename
	if tokens[newCur].Kind != lexer.IdentifierKind {
		return nil, cursor, false
	}
	insertNode.Table = tokens[newCur]
	newCur ++;

	if !tokens[newCur].IsEqual(GenerateToken(lexer.SymbolKind, string(lexer.LeftparenSymbol))) {
		return nil, cursor, false
	}
	newCur++

	exps, expCur, ok := parseExpressions(tokens, newCur, []lexer.Token{GenerateToken(lexer.SymbolKind, string(lexer.LeftparenSymbol))})
	if !ok {
		return nil, cursor, false
	}
	insertNode.Values = exps
	newCur = expCur + 2

	return insertNode, newCur, true
}

func parseCreate(tokens []*lexer.Token, cursor uint, delimiter lexer.Token) {

	
}