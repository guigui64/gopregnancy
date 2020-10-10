// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Page struct {
	User    string
	Message string
	Tip     string
	Offset  int
}

var templates = template.Must(template.ParseFiles("index.html"))

func loadPage(name string) (*Page, error) {
	filename := name + ".txt"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	offset, _ := strconv.Atoi(lines[4])
	return &Page{User: lines[0], Message: lines[2], Tip: lines[3], Offset: offset}, nil
}

func renderTemplate(w http.ResponseWriter, p *Page) {
	err := templates.ExecuteTemplate(w, "index.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// TODO queries
	name := "z4z401b7e977"
	p, err := loadPage(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	renderTemplate(w, p)
}

func shift(s string, offset byte) string {
	b := []byte(s)
	for i, bi := range b {
		if 'A' <= bi && bi <= 'Z' {
			b[i] += offset
			b[i] = (b[i]-'A')%('Z'-'A') + 'A'
		}
	}
	return string(b)
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
