package parser_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rsasada/sqluid/srcs/lexer"
	"github.com/rsasada/sqluid/srcs/parser"
)

var _ = Describe("Parser", func() {
	var (
		tokens []*lexer.Token
		source string
	)

	BeforeEach(func() {
		
		var err error
		source = "SELECT name FROM table1;"
		tokens, err = lexer.Lexing(source)
		if err != nil {
			return
		}
	})

	Describe("Parsing SELECT statements", func() {
		It("parses a simple SELECT statement successfully", func() {
			ast, ok := parser.Parser(source, tokens)
			Expect(ok).To(BeTrue())
			Expect(ast).ToNot(BeNil())
			Expect(ast.Kind).To(Equal(parser.SelectType))
			Expect(ast.Select).ToNot(BeNil())
			Expect(ast.Select.From.Value).To(Equal("table1"))
			Expect(len(*ast.Select.Item)).To(Equal(1))
			Expect((*ast.Select.Item)[0].Literal.Value).To(Equal("name"))
		})

		It("returns false for invalid SELECT syntax", func() {

			invalidTokens := []*lexer.Token{
				{Kind: lexer.KeywordKind, Value: "SELECT"},
				{Kind: lexer.IdentifierKind, Value: "name"},
				{Kind: lexer.KeywordKind, Value: "WHERE"}, // "FROM" の代わりに "WHERE"
				{Kind: lexer.IdentifierKind, Value: "table1"},
				{Kind: lexer.SymbolKind, Value: ";"},
			}
			ast, ok := parser.Parser(source, invalidTokens)
			Expect(ok).To(BeFalse())
			Expect(ast).To(BeNil())
		})
	})

	Describe("Parsing INSERT statements", func() {
		BeforeEach(func() {
			// INSERT用のトークンリスト
			var err error
			source = "INSERT INTO table1 (name);"
			tokens, err = lexer.Lexing(source)
			if err != nil {
				return
			}
		})

		It("parses a simple INSERT statement successfully", func() {
			ast, ok := parser.Parser(source, tokens)
			Expect(ok).To(BeTrue())
			Expect(ast).ToNot(BeNil())
			Expect(ast.Kind).To(Equal(parser.InsertType))
			Expect(ast.Insert).ToNot(BeNil())
			Expect(ast.Insert.Table.Value).To(Equal("table1"))
			Expect(len(*ast.Insert.Values)).To(Equal(1))
			Expect((*ast.Insert.Values)[0].Literal.Value).To(Equal("name"))
		})
	})

	Describe("Parsing CREATE TABLE statements", func() {
		BeforeEach(func() {
			// CREATE TABLE用のトークンリスト
			var err error
			source = "CREATE TABLE table1 (id INT);"
			tokens, err = lexer.Lexing(source)
			if err != nil {
				return
			}
		})

		It("parses a simple CREATE TABLE statement successfully", func() {
			ast, ok := parser.Parser(source, tokens)
			Expect(ok).To(BeTrue())
			Expect(ast).ToNot(BeNil())
			Expect(ast.Kind).To(Equal(parser.CreateTableType))
			Expect(ast.Create).ToNot(BeNil())
			Expect(ast.Create.TableName.Value).To(Equal("table1"))
			Expect(len(*ast.Create.Cols)).To(Equal(1))
			Expect((*ast.Create.Cols)[0].Name.Value).To(Equal("id"))
			Expect((*ast.Create.Cols)[0].DataType.Value).To(Equal("int"))
		})
	})

	Describe("Error handling for unsupported syntax", func() {
		It("returns false for unrecognized keywords", func() {
			unsupportedTokens := []*lexer.Token{
				{Kind: lexer.KeywordKind, Value: "DELETE"},
				{Kind: lexer.IdentifierKind, Value: "name"},
				{Kind: lexer.KeywordKind, Value: "FROM"},
				{Kind: lexer.IdentifierKind, Value: "table1"},
				{Kind: lexer.SymbolKind, Value: ";"},
			}
			ast, ok := parser.Parser(source, unsupportedTokens)
			Expect(ok).To(BeFalse())
			Expect(ast).To(BeNil())
		})
	})
})
