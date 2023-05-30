package configs

import (
	"html/template"
	"log"

	"github.com/go-playground/form/v4"
	"github.com/nexentra/snippetbox/internal/models"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Snippets *models.SnippetModel
	TemplateCache map[string]*template.Template
	FormDecoder *form.Decoder
}

// func SetLogger(logger *Logger) {
// 	var LoggerInstance = &Logger{}
// 	LoggerInstance.ErrorLog = logger.ErrorLog
// 	LoggerInstance.InfoLog = logger.InfoLog
// }