package config

import (
	"log"

	"github.com/qbitty/snippetbox/pkg/models/mysql"
)

type Application struct {
	ErrLog   *log.Logger
	InfoLog  *log.Logger
	Snippets *mysql.SnippetModel
}
