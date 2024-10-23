package lexer_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLexer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lexer Suite")
}
