package lexer_test

import (
	"github.com/rsasada/sqluid/srcs/lexer/lexer.go"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = Describe("Lexer", func() {
	Describe("Tokenization", func() {
		tests := []struct {
			input  string
			Tokens []Token
			err    error
		}{
			{
				input: "select a",
				Tokens: []Token{
					{
						Loc:   Location{Col: 0, Line: 0},
						Value: string(SelectKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 7, Line: 0},
						Value: "a",
						Kind:  IdentifierKind,
					},
				},
			},
			{
				input: "select true",
				Tokens: []Token{
					{
						Loc:   Location{Col: 0, Line: 0},
						Value: string(SelectKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 7, Line: 0},
						Value: "true",
						Kind:  BoolKind,
					},
				},
			},
			{
				input: "select 1",
				Tokens: []Token{
					{
						Loc:   Location{Col: 0, Line: 0},
						Value: string(SelectKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 7, Line: 0},
						Value: "1",
						Kind:  NumericKind,
					},
				},
				err: nil,
			},
			{
				input: "select 'foo' || 'bar';",
				Tokens: []Token{
					{
						Loc:   Location{Col: 0, Line: 0},
						Value: string(SelectKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 7, Line: 0},
						Value: "foo",
						Kind:  StringKind,
					},
					{
						Loc:   Location{Col: 13, Line: 0},
						Value: string(ConcatSymbol),
						Kind:  SymbolKind,
					},
					{
						Loc:   Location{Col: 16, Line: 0},
						Value: "bar",
						Kind:  StringKind,
					},
					{
						Loc:   Location{Col: 21, Line: 0},
						Value: string(SemicolonSymbol),
						Kind:  SymbolKind,
					},
				},
				err: nil,
			},
			{
				input: "CREATE TABLE u (id INT, name TEXT)",
				Tokens: []Token{
					{
						Loc:   Location{Col: 0, Line: 0},
						Value: string(CreateKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 7, Line: 0},
						Value: string(TableKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 13, Line: 0},
						Value: "u",
						Kind:  IdentifierKind,
					},
					{
						Loc:   Location{Col: 15, Line: 0},
						Value: "(",
						Kind:  SymbolKind,
					},
					{
						Loc:   Location{Col: 16, Line: 0},
						Value: "id",
						Kind:  IdentifierKind,
					},
					{
						Loc:   Location{Col: 19, Line: 0},
						Value: "int",
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 22, Line: 0},
						Value: ",",
						Kind:  SymbolKind,
					},
					{
						Loc:   Location{Col: 24, Line: 0},
						Value: "name",
						Kind:  IdentifierKind,
					},
					{
						Loc:   Location{Col: 29, Line: 0},
						Value: "text",
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 33, Line: 0},
						Value: ")",
						Kind:  SymbolKind,
					},
				},
			},
			{
				input: "insert into users Values (105, 233)",
				Tokens: []Token{
					{
						Loc:   Location{Col: 0, Line: 0},
						Value: string(InsertKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 7, Line: 0},
						Value: string(IntoKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 12, Line: 0},
						Value: "users",
						Kind:  IdentifierKind,
					},
					{
						Loc:   Location{Col: 18, Line: 0},
						Value: string(ValuesKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 25, Line: 0},
						Value: "(",
						Kind:  SymbolKind,
					},
					{
						Loc:   Location{Col: 26, Line: 0},
						Value: "105",
						Kind:  NumericKind,
					},
					{
						Loc:   Location{Col: 30, Line: 0},
						Value: ",",
						Kind:  SymbolKind,
					},
					{
						Loc:   Location{Col: 32, Line: 0},
						Value: "233",
						Kind:  NumericKind,
					},
					{
						Loc:   Location{Col: 36, Line: 0},
						Value: ")",
						Kind:  SymbolKind,
					},
				},
				err: nil,
			},
			{
				input: "SELECT id FROM users;",
				Tokens: []Token{
					{
						Loc:   Location{Col: 0, Line: 0},
						Value: string(SelectKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 7, Line: 0},
						Value: "id",
						Kind:  IdentifierKind,
					},
					{
						Loc:   Location{Col: 10, Line: 0},
						Value: string(FromKeyword),
						Kind:  KeywordKind,
					},
					{
						Loc:   Location{Col: 15, Line: 0},
						Value: "users",
						Kind:  IdentifierKind,
					},
					{
						Loc:   Location{Col: 20, Line: 0},
						Value: ";",
						Kind:  SymbolKind,
					},
				},
				err: nil,
			},
		}

		for _, test := range tests {
			test := test // ループ変数を保持するためのシャドウイング
			It(test.input, func() {
				tokens, err := lexing(test.input)
				Expect(err).To(Equal(test.err))

				Expect(len(tokens)).To(Equal(len(test.Tokens)))

				for i, tok := range tokens {
					Expect(&test.Tokens[i]).To(Equal(tok))
				}
			})
		}
	})
})