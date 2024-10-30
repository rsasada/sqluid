package parser

import (
	"fmt"

	"github.com/rsasada/sqluid/srcs/lexer"
)

func Parser(source string, tokens []*lexer.Token) (*Ast, bool) {

	var topAst *Ast
	var curAst *Ast

	cursor := uint(0)
	numAsts := 0

	for cursor < uint(len(tokens)) && tokens[cursor].Kind != lexer.EndKind {
		newAst, newCursor, ok := parsingTokens(source, tokens, cursor)
		if !ok {
			return nil, false
		}
		if newCursor < uint(len(tokens)) && newAst != nil && tokens[newCursor].Kind != lexer.EndKind {
			curAst = GeneratePipe()
			curAst.Pipe.Left = newAst
			if numAsts == 0 {
				topAst = curAst
			}
			curAst = curAst.Pipe.Right
		} else {
			curAst = newAst
		}
		cursor = newCursor
		numAsts++
	}
	if numAsts == 1 {
		topAst = curAst
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
			Kind:   SelectType,
			Select: slct,
		}, newCur, true

	} else if tokens[cursor].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.CreateKeyword))) {
		createNode, newCur, ok := parseCreateTable(tokens, cursor, semicolonToken)
		if !ok {
			return nil, cursor, false
		}
		return &Ast{
			Kind:   CreateTableType,
			Create: createNode,
		}, newCur, true

	} else if tokens[cursor].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.InsertKeyword))) {
		insert, newCur, ok := parseInsert(tokens, cursor, semicolonToken)
		if !ok {
			return nil, cursor, false
		}
		return &Ast{
			Kind:   InsertType,
			Insert: insert,
		}, newCur, true

	} else {
		return nil, cursor, false
	}

}

func parseSelect(tokens []*lexer.Token, cursor uint, delimiter lexer.Token) (*SelectNode, uint, bool) {

	newCur := cursor + 1
	selectNode := SelectNode{}

	exps, expCursor, ok := parseExpressions(tokens, newCur,
		[]lexer.Token{GenerateToken(lexer.KeywordKind, string(lexer.FromKeyword)), delimiter})
	if !ok {
		return nil, cursor, false
	}
	newCur = expCursor
	selectNode.Item = exps

	if tokens[newCur].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.FromKeyword))) {
		newCur++
		if tokens[newCur].Kind == lexer.IdentifierKind {
			selectNode.From = tokens[newCur]
		}
		newCur++
	}

	if tokens[newCur].IsEqual(delimiter) {
		newCur++
		return &selectNode, newCur, true
	} else {
		return nil, cursor, false
	}
}

func parseExpressions(tokens []*lexer.Token, cursor uint, delimiter []lexer.Token) (*[]*Expression, uint, bool) {

	newCur := cursor
	exps := []*Expression{}

extract:
	for {

		if cursor >= uint(len(tokens)) && tokens[cursor].Kind == lexer.EndKind {
			return nil, cursor, false
		}

		for i := 0; i < len(delimiter); i++ {
			if tokens[newCur].IsEqual(delimiter[i]) == true {
				break extract
			}
		}

		if cursor != newCur {
			fmt.Print(tokens[newCur].Value)
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

	if len(exps) == 0 {
		return nil, cursor, false
	}

	return &exps, newCur, true
}

func extractExpression(tokens []*lexer.Token, cursor uint, allowKinds []lexer.TokenKind) (*Expression, uint, bool) {

	if cursor >= uint(len(tokens)) && tokens[cursor].Kind == lexer.EndKind {
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

	newCur := cursor
	newCur++

	if !tokens[newCur].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.IntoKeyword))) {
		return nil, cursor, false
	}
	newCur++

	//parse tablename
	if tokens[newCur].Kind != lexer.IdentifierKind {
		return nil, cursor, false
	}
	table := tokens[newCur]
	newCur++

	if !tokens[newCur].IsEqual(GenerateToken(lexer.SymbolKind, string(lexer.LeftparenSymbol))) {
		return nil, cursor, false
	}
	newCur++

	exps, expCur, ok := parseExpressions(tokens, newCur, []lexer.Token{GenerateToken(lexer.SymbolKind, string(lexer.RightparenSymbol))})
	if !ok {
		return nil, cursor, false
	}
	newCur = expCur + 1

	if tokens[newCur].IsEqual(delimiter) {
		newCur++
		return &InsertNode{
			Table:  table,
			Values: exps,
		}, newCur, true
	} else {
		return nil, cursor, false
	}
}

func parseCreateTable(tokens []*lexer.Token, cursor uint, delimiter lexer.Token) (*CreateTableNode, uint, bool) {

	newCur := cursor + 1
	createNode := &CreateTableNode{}
	delimiterColumn := GenerateToken(lexer.SymbolKind, string(lexer.RightparenSymbol))

	if !tokens[newCur].IsEqual(GenerateToken(lexer.KeywordKind, string(lexer.TableKeyword))) {
		return nil, cursor, false
	}
	newCur++

	if tokens[newCur].Kind != lexer.IdentifierKind {
		return nil, cursor, false
	}
	createNode.TableName = tokens[newCur]
	newCur++

	if !tokens[newCur].IsEqual(GenerateToken(lexer.SymbolKind, string(lexer.LeftparenSymbol))) {
		return nil, cursor, false
	}
	newCur++

	colmuns, colCur, ok := parseTableColumn(tokens, newCur, delimiterColumn)
	if !ok {
		return nil, cursor, false
	}
	createNode.Cols = colmuns
	newCur = colCur + 1

	if tokens[newCur].IsEqual(delimiter) {
		newCur++
		return createNode, newCur, true
	} else {
		return nil, cursor, false
	}
}

func parseTableColumn(tokens []*lexer.Token, cursor uint, delimiter lexer.Token) (*[]*TableColumn, uint, bool) {

	columns := []*TableColumn{}
	newCur := cursor

extract:
	for {
		if newCur >= uint(len(tokens)) && tokens[newCur].Kind == lexer.EndKind {
			return nil, cursor, false
		}

		if tokens[newCur].IsEqual(delimiter) {
			break extract
		}

		if len(columns) > 0 {
			if !tokens[newCur].IsEqual(GenerateToken(lexer.SymbolKind, string(lexer.CommaSymbol))) {
				return nil, cursor, false
			}
			newCur++
		}

		if tokens[newCur].Kind != lexer.IdentifierKind {
			return nil, cursor, false
		}
		columnName := tokens[newCur]
		newCur++

		if tokens[newCur].Kind != lexer.IdentifierKind {
			return nil, cursor, false
		}
		columnType := tokens[newCur]
		newCur++

		columns = append(columns, &TableColumn{
			Name:     columnName,
			DataType: columnType,
		})
	}

	return &columns, newCur, true
}
