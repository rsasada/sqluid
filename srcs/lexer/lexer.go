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

type keyword string

const (
	selectKeyword keyword = "select"
	fromKeyword   keyword = "from"
	asKeyword     keyword = "as"
	tableKeyword  keyword = "table"
	createKeyword keyword = "create"
	insertKeyword keyword = "insert"
	intoKeyword   keyword = "into"
	valuesKeyword keyword = "values"
	intKeyword    keyword = "int"
	textKeyword   keyword = "text"
)

type symbol string

const (
	semicolonSymbol  symbol = ";"
	asteriskSymbol   symbol = "*"
	commaSymbol      symbol = ","
	leftparenSymbol  symbol = "("
	rightparenSymbol symbol = ")"
)

type tokenKind uint

const (
	keywordKind tokenKind = iota
	symbolKind
	identifierKind
	stringKind
	numericKind
)

type Token struct {
	value string
	kind  tokenKind
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
			hint = " after " + tokens[len(tokens)-1].value
		}
		return nil, fmt.Errorf("Unable to lex token%s, at %d:%d", hint, cur.pos.line, cur.pos.col)
	}

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
		if (!isDigit) {
			break
		}
	}

	if cur.index == ic.index {
		return nil, ic, false
	}

	return &Token{
		value: source[ic.index:cur.index],
		pos:   ic.pos,
		kind:  numericKind,
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
	symbols := []symbol{
		commaSymbol,
		leftparenSymbol,
		rightparenSymbol,
		semicolonSymbol,
		asteriskSymbol,
	}

	var options []string
	for _, s := range symbols {
		options = append(options, string(s))
	}

	match := longestMatch(source, ic, options)

	if match == "" {
		return nil, ic, false
	}

	cur.index = ic.index + uint(len(match))
	cur.pos.col = ic.pos.col + uint(len(match))

	return &Token{
		value: match,
		pos:   ic.pos,
		kind:  symbolKind,
	}, cur, true
}
func lexKeyword(source string, ic Cursor) (*Token, Cursor, bool) {
    cur := ic
    keywords := []keyword{
        selectKeyword,
        insertKeyword,
        valuesKeyword,
        tableKeyword,
        createKeyword,
        fromKeyword,
        intoKeyword,
        textKeyword,
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
        value: match,
        kind:  keywordKind,
        pos:   ic.pos,
    }, cur, true
}

func lexIdentifier(source string, ic Cursor) (*Token, Cursor, bool) {

	value := []byte{}
	if token, newCursor, ok := lexCharacterDelimited(source, ic, '"'); ok {
		return token, newCursor, true
	}

	cur := ic

	for ; cur.index < uint(len(source)); cur.index ++ {
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

	if (len(value) == 0) {
		return nil, ic, false
	}

	return  &Token {
		value: strings.ToLower(string(value)),
		pos: ic.pos,
		kind: identifierKind,
	}, cur, true

}
