package server

import (
	"log"
	"net/http"
	Backend "github.com/rsasada/sqluid/srcs/backend"
	Lexer"github.com/rsasada/sqluid/srcs/lexer"
	Parser"github.com/rsasada/sqluid/srcs/parser"
)

type ApiServer struct {
	db	*SqluiDB
}

type SqluiDB struct {
	tokens	[]*Lexer.Token
	ast		*Parser.Ast
	backend	*Backend.MemoryBackend
}

func NewApiServer() *ApiServer {
	return &ApiServer {
		db: NewSqluiDB(),
	}
}

func NewSqluiDB() *SqluiDB {
	return &SqluiDB{}
}
