package parser

import (
	"fmt"
	"strings"
	"github.com/rsasada/sqluid/srcs/lexer"
)

func GenerateToken(kind lexer.TokenKind, value string) lexer.Token {

	return lexer.Token{
		Value: value,
		Kind: kind,
	}
}

func GeneratePipe() *Ast {

	return &Ast{
		Kind:		BinaryPipeType,
		Pipe:		&BinaryPipeNode{},
	}
}

func PrintAst(ast *Ast, indent int) {
	if ast == nil {
		return
	}
	indentation := strings.Repeat("  ", indent)

	switch ast.Kind {
	case SelectType:
		fmt.Printf("%sSelect Node:\n", indentation)
		printSelectNode(ast.Select, indent+1)
	case CreateTableType:
		fmt.Printf("%sCreate Table Node:\n", indentation)
		printCreateTableNode(ast.Create, indent+1)
	case InsertType:
		fmt.Printf("%sInsert Node:\n", indentation)
		printInsertNode(ast.Insert, indent+1)
	case BinaryPipeType:
		fmt.Printf("%sBinary Pipe Node:\n", indentation)
		fmt.Printf("%sLeft:\n", indentation)
		PrintAst(ast.Pipe.Left, indent+1)
		fmt.Printf("%sRight:\n", indentation)
		PrintAst(ast.Pipe.Right, indent+1)
	}
}

func printSelectNode(selectNode *SelectNode, indent int) {
	if selectNode == nil {
		return
	}
	indentation := strings.Repeat("  ", indent)
	fmt.Printf("%sItems:\n", indentation)
	for _, expr := range *selectNode.Item {
		printExpression(expr, indent+1)
	}
	fmt.Printf("%sFrom: %s\n", indentation, selectNode.From.Value)
}

func printCreateTableNode(createNode *CreateTableNode, indent int) {
	if createNode == nil {
		return
	}
	indentation := strings.Repeat("  ", indent)
	fmt.Printf("%sTable Name: %s\n", indentation, createNode.TableName.Value)
	fmt.Printf("%sColumns:\n", indentation)
	for _, col := range *createNode.Cols {
		fmt.Printf("%s  Name: %s, Type: %s\n", indentation, col.Name.Value, col.DataType.Value)
	}
}

func printInsertNode(insertNode *InsertNode, indent int) {
	if insertNode == nil {
		return
	}
	indentation := strings.Repeat("  ", indent)
	fmt.Printf("%sTable: %s\n", indentation, insertNode.Table.Value)
	fmt.Printf("%sValues:\n", indentation)
	for _, expr := range *insertNode.Values {
		printExpression(expr, indent+1)
	}
}

func printExpression(expr *Expression, indent int) {
	indentation := strings.Repeat("  ", indent)
	switch expr.Kind {
	case LiteralType:
		fmt.Printf("%sLiteral: %s\n", indentation, expr.Literal.Value)
	}
}
