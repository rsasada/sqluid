package sqluid

import (
	"fmt"
	"strings"
)


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

func longestMatch(source string, ic cursor, options []string) string {
    var value []byte
    var skipList []int
    var match string

    cur := ic

    for cur.pointer < uint(len(source)) {

        value = append(value, strings.ToLower(string(source[cur.pointer]))...)
        cur.pointer++

    match:
        for i, option := range options {
            for _, skip := range skipList {
                if i == skip {
                    continue match
                }
            }

            // Deal with cases like INT vs INTO
            if option == string(value) {
                skipList = append(skipList, i)
                if len(option) > len(match) {
                    match = option
                }

                continue
            }

            sharesPrefix := string(value) == option[:cur.pointer-ic.pointer]
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