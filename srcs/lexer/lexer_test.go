package lexer_test

import (

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rsasada/sqluid/srcs/lexer"
)

var _ = Describe("Lexer", func() {
	Context("Lexing simple SQL queries", func() {

		It("should correctly lex a simple SELECT query", func() {
			source := "SELECT * FROM table1;"
			tokens, err := lexer.Lexing(source)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(tokens)).To(Equal(6))

			Expect(tokens[0].Value).To(Equal("select"))
			Expect(tokens[0].Kind).To(Equal(lexer.KeywordKind))

			Expect(tokens[1].Value).To(Equal("*"))
			Expect(tokens[1].Kind).To(Equal(lexer.SymbolKind))

			Expect(tokens[2].Value).To(Equal("from"))
			Expect(tokens[2].Kind).To(Equal(lexer.KeywordKind))

			Expect(tokens[3].Value).To(Equal("table1"))
			Expect(tokens[3].Kind).To(Equal(lexer.IdentifierKind))

			Expect(tokens[4].Value).To(Equal(";"))
			Expect(tokens[4].Kind).To(Equal(lexer.SymbolKind))

		})

		It("should return an error for invalid input", func() {
			source := "SELECT ^ FROM table;"
			tokens, err := lexer.Lexing(source)

			Expect(err).To(HaveOccurred())
			Expect(tokens).To(BeNil())
		})

		It("should correctly lex a CREATE TABLE query", func() {
			source := "CREATE TABLE users (id INT, name TEXT);"
			tokens, err := lexer.Lexing(source)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(tokens)).To(Equal(12))

			Expect(tokens[0].Value).To(Equal("create"))
			Expect(tokens[0].Kind).To(Equal(lexer.KeywordKind))

			Expect(tokens[1].Value).To(Equal("table"))
			Expect(tokens[1].Kind).To(Equal(lexer.KeywordKind))

			Expect(tokens[2].Value).To(Equal("users"))
			Expect(tokens[2].Kind).To(Equal(lexer.IdentifierKind))

			Expect(tokens[3].Value).To(Equal("("))
			Expect(tokens[3].Kind).To(Equal(lexer.SymbolKind))

			Expect(tokens[4].Value).To(Equal("id"))
			Expect(tokens[4].Kind).To(Equal(lexer.IdentifierKind))

			Expect(tokens[5].Value).To(Equal("int"))
			Expect(tokens[5].Kind).To(Equal(lexer.KeywordKind))

			Expect(tokens[6].Value).To(Equal(","))
			Expect(tokens[6].Kind).To(Equal(lexer.SymbolKind))

			Expect(tokens[7].Value).To(Equal("name"))
			Expect(tokens[7].Kind).To(Equal(lexer.IdentifierKind))

			Expect(tokens[8].Value).To(Equal("text"))
			Expect(tokens[8].Kind).To(Equal(lexer.KeywordKind))

			Expect(tokens[9].Value).To(Equal(")"))
			Expect(tokens[9].Kind).To(Equal(lexer.SymbolKind))

			Expect(tokens[10].Value).To(Equal(";"))
			Expect(tokens[10].Kind).To(Equal(lexer.SymbolKind))
		})

		It("should lex numeric values", func() {
			source := "INSERT INTO table VALUES (1, 2.5);"
			tokens, err := lexer.Lexing(source)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(tokens)).To(Equal(11))

			Expect(tokens[5].Value).To(Equal("1")) 
			Expect(tokens[5].Kind).To(Equal(lexer.NumericKind))

			Expect(tokens[7].Value).To(Equal("2.5"))
			Expect(tokens[7].Kind).To(Equal(lexer.NumericKind))
		})
	})
})
