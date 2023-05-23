package main

import "github.com/nexentra/snippetbox/internal/models"
type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
	}