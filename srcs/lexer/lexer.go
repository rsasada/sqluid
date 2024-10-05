package sqluid

import (
	"fmt"
	"strings"
)

type position struct {
    line uint
    col  uint
}

type cursor struct {
	index	uint
	pos 	position
}

type tokenKind uint

const (
    keywordKind tokenKind = iota
    symbolKind
    identifierKind
    stringKind
    numericKind
)

type token struct {
	value	string
	kind 	tokenKind
	pos		position
}

type lexer func(string, cursor) (*token, cursor, bool)

func lexing(source string) ([]*token, error) {
    tokens := []*token{}
    cur := cursor{}

lex:
    for cur.index < uint(len(source)) {
        lexers := []lexer{lexKeyword, lexSymbol, lexString, lexNumeric, lexIdentifier}
        for _, l := range lexers {
            if token, newCursor, ok := l(source, cur); ok {
                cur = newCursor

                // Omit nil tokens for valid, but empty syntax like newlines
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

func lexNumeric(source string, ic cursor) (*token, cursor, bool) {
    cur := ic

    periodFound := false
    expMarkerFound := false

    for ; cur.index < uint(len(source)); cur.index++ {
        c := source[cur.index]
        cur.pos.col++

        isDigit := c >= '0' && c <= '9'
        isPeriod := c == '.'
        isExpMarker := c == 'e'

        if cur.index == ic.index {
            if !isDigit && !isPeriod {
                return nil, ic, false
            }

            periodFound = isPeriod
            continue
        }

        if isPeriod {
            if periodFound {
                return nil, ic, false
            }

            periodFound = true
            continue
        }

        if isExpMarker {
            if expMarkerFound {
                return nil, ic, false
            }

            periodFound = true
            expMarkerFound = true

            // expMarker must be followed by digits
            if cur.index == uint(len(source)-1) {
                return nil, ic, false
            }

            cNext := source[cur.index+1]
            if cNext == '-' || cNext == '+' {
                cur.index++
                cur.pos.col++
            }

            continue
        }

        if !isDigit {
            break
        }
    }

    // No characters accumulated
    if cur.index == ic.index {
        return nil, ic, false
    }

    return &token{
        value: source[ic.index:cur.index],
        pos:   ic.pos,
        kind:  numericKind,
    }, cur, true
}

func lexCharacterDelimited(source string, ic cursor, delimiter byte) (*token, cursor, bool) {
    cur := ic

    if len(source[cur.index:]) == 0 {
        return nil, ic, false
    }

    if source[cur.index] != delimiter {
        return nil, ic, false
    }

    cur.pos.col++
    cur.index++

    var value []byte
    for ; cur.index < uint(len(source)); cur.index++ {
        c := source[cur.index]

        if c == delimiter {

            if cur.index+1 >= uint(len(source)) || source[cur.index+1] != delimiter {
                return &token{
                    value: string(value),
                    pos:   ic.pos,
                    kind:  stringKind,
                }, cur, true
            } else {
                value = append(value, delimiter)
                cur.index++
                cur.pos.col++
            }
        }

        value = append(value, c)
        cur.pos.col++
    }

    return nil, ic, false
}

func lexString(source string, ic cursor) (*token, cursor, bool) {
    return lexCharacterDelimited(source, ic, '\'')
}

