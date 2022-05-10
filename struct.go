package main

import "pdf-server/pdf"

type request struct {
	Id       int          `json:"id"`
	Raw      *pdf.Content `json:"raw"`
	Template *Template    `json:"template"`
	FileName string       `json:"fileName"`
}

type Template struct {
	Path    string            `json:"path"`
	Replace map[string]string `json:"replace"`
	Insert  map[string][]int  `json:"insert"`
}
