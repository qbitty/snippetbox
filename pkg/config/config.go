package config

import (
	"html/template"
	"log"

	"github.com/qbitty/snippetbox/pkg/models/mysql"
)

type Application struct {
	ErrLog        *log.Logger
	InfoLog       *log.Logger
	Snippets      *mysql.SnippetModel
	TemplateCache map[string]*template.Template
}
