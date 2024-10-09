package sqluid

import (
	"fmt"
	"strings"
)

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