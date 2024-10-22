package lexer

import (
	"fmt"
	"strings"
)

type position struct {
	line uint
	col  uint
}

type Cursor struct {
	index uint
	pos   position
}

type Keyword string

const (
	SelectKeyword Keyword = "select"
	FromKeyword   Keyword = "from"
	AsKeyword     Keyword = "as"
	TableKeyword  Keyword = "table"
	CreateKeyword Keyword = "create"
	InsertKeyword Keyword = "insert"
	IntoKeyword   Keyword = "into"
	ValuesKeyword Keyword = "values"
	IntKeyword    Keyword = "int"
	TextKeyword   Keyword = "text"
)

type Symbol string

const (
	SemicolonSymbol  Symbol = ";"
	AsteriskSymbol   Symbol = "*"
	CommaSymbol      Symbol = ","
	LeftparenSymbol  Symbol = "("
	RightparenSymbol Symbol = ")"
)

type TokenKind uint

const (
	KeywordKind TokenKind = iota
	SymbolKind
	IdentifierKind
	StringKind
	NumericKind
	EndKind
)

type Token struct {
	Value string
	Kind  TokenKind
	pos   position
}

type lexer func(string, Cursor) (*Token, Cursor, bool)

func Lexing(source string) ([]*Token, error) {
	tokens := []*Token{}
	cur := Cursor{}

lex:
	for cur.index < uint(len(source)) {
		lexers := []lexer{lexKeyword, lexSymbol, lexString, lexNumeric, lexIdentifier}
		for _, l := range lexers {
			if token, newCursor, ok := l(source, cur); ok {
				cur = newCursor

				if token != nil {
					tokens = append(tokens, token)
				}

				continue lex
			}
		}

		hint := ""
		if len(tokens) > 0 {
			hint = " after " + tokens[len(tokens)-1].Value
		}
		return nil, fmt.Errorf("Unable to lex token%s, at %d:%d", hint, cur.pos.line, cur.pos.col)
	}
	tokens = append(tokens, &Token{Kind: EndKind,})
	return tokens, nil
}

func lexNumeric(source string, ic Cursor) (*Token, Cursor, bool) {
	cur := ic
	periodFlag := false

	for ; cur.index < uint(len(source)); cur.index++ {
		c := source[cur.index]
		cur.pos.col++

		isDigit := c >= '0' && c <= '9'
		isPeriod := c == '.'
		isSign := c == '+' || c == '-'

		if cur.index == ic.index && isSign {
			continue
		}
		if isPeriod {
			if periodFlag == true {
				return nil, ic, false
			}
			periodFlag = true
			continue
		}
		if !isDigit {
			break
		}
	}

	if cur.index == ic.index {
		return nil, ic, false
	}

	return &Token{
		Value: source[ic.index:cur.index],
		pos:   ic.pos,
		Kind:  NumericKind,
	}, cur, true
}

func lexString(source string, ic Cursor) (*Token, Cursor, bool) {
	return lexCharacterDelimited(source, ic, '\'')
}

func lexSymbol(source string, ic Cursor) (*Token, Cursor, bool) {
	c := source[ic.index]
	cur := ic

	cur.index++
	cur.pos.col++

	switch c {

	case '\n':
		cur.pos.line++
		cur.pos.col = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, cur, true
	}

	// Syntax that should be kept
	Symbols := []Symbol{
		CommaSymbol,
		LeftparenSymbol,
		RightparenSymbol,
		SemicolonSymbol,
		AsteriskSymbol,
	}

	var options []string
	for _, s := range Symbols {
		options = append(options, string(s))
	}

	match := longestMatch(source, ic, options)

	if match == "" {
		return nil, ic, false
	}

	cur.index = ic.index + uint(len(match))
	cur.pos.col = ic.pos.col + uint(len(match))

	return &Token{
		Value: match,
		pos:   ic.pos,
		Kind:  SymbolKind,
	}, cur, true
}
func lexKeyword(source string, ic Cursor) (*Token, Cursor, bool) {
	cur := ic
	keywords := []Keyword{
		SelectKeyword,
		InsertKeyword,
		ValuesKeyword,
		TableKeyword,
		CreateKeyword,
		FromKeyword,
		IntoKeyword,
		IntKeyword,
		TextKeyword,
	}

	var options []string
	for _, k := range keywords {
		options = append(options, string(k))
	}

	match := longestMatch(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.index = ic.index + uint(len(match))
	cur.pos.col = ic.pos.col + uint(len(match))

	return &Token{
		Value: match,
		Kind:  KeywordKind,
		pos:   ic.pos,
	}, cur, true
}

func lexIdentifier(source string, ic Cursor) (*Token, Cursor, bool) {

	value := []byte{}
	if token, newCursor, ok := lexCharacterDelimited(source, ic, '"'); ok {
		return token, newCursor, true
	}

	cur := ic

	for ; cur.index < uint(len(source)); cur.index++ {
		c := source[cur.index]

		isAlpabet := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
		isNum := (c >= '0' && c <= '9')
		if isAlpabet || isNum || c == '$' || c == '_' {
			value = append(value, c)
			cur.pos.col++
			continue
		}

		break
	}

	if len(value) == 0 {
		return nil, ic, false
	}

	return &Token{
		Value: strings.ToLower(string(value)),
		pos:   ic.pos,
		Kind:  IdentifierKind,
	}, cur, true

}

func (t *Token) IsEqual(compare Token) bool {
	return compare.Value == t.Value && compare.Kind == t.Kind
}
