package lexer

import (
	"strings"
)

func lexCharacterDelimited(source string, ic Cursor, delimiter byte) (*Token, Cursor, bool) {
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
				return &Token{
					Value: string(value),
					pos:   ic.pos,
					Kind:  StringKind,
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

func longestMatch(source string, ic Cursor, options []string) string {
    var value []byte
    var skipList []int
    var match string

    cur := ic

    for cur.index < uint(len(source)) {

        value = append(value, strings.ToLower(string(source[cur.index]))...)
        cur.index++

    match:
        for i, option := range options {
            for _, skip := range skipList {
                if i == skip {
                    continue match
                }
            }

            if option == string(value) {
                skipList = append(skipList, i)
                if len(option) > len(match) {
                    match = option
                }

                continue
            }

            sharesPrefix := string(value) == option[:cur.index-ic.index]
            tooLong := len(value) > len(option)
            if tooLong || !sharesPrefix {
                skipList = append(skipList, i)
            }
        }

        if len(skipList) == len(options) {
            break
        }
    }

    return match
}