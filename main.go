// Enigma website to announce our pregnancy to my buddies
// Copyright 2020 Guillaume Comte

package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// parsed passwords
var passwords []string

// html templates
var templates = template.Must(template.ParseFiles("index.html", "good0.html", "good1.html", "wrong.html", "404.html", "end.html"))

// template struct
type Page struct {
	User           string
	CodedMessage   string
	DecodedMessage string
	Tip            string
	Offset         byte
	Guess          string
	Code           string
	Steps          []string
	CurrentStep    int
}

func parsePasswords() {
	content, err := ioutil.ReadFile("passwords.txt")
	if err != nil {
		log.Panic(err)
	}
	passwords = strings.Split(string(content), "\n")
}

// Shift letters
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

func loadPage(name string) (*Page, error) {
	filename := name + ".txt"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	offset, _ := strconv.Atoi(lines[4])
	return &Page{User: lines[0], DecodedMessage: lines[1], CodedMessage: lines[2], Tip: lines[3], Offset: byte(offset), Code: lines[5], Steps: lines[6:]}, nil
}

func renderTemplate(w http.ResponseWriter, p *Page, src string) {
	err := templates.ExecuteTemplate(w, src+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func notFound(w http.ResponseWriter, reason string) {
	err := templates.ExecuteTemplate(w, "404.html", reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name, ok := q["user"]
	if !ok {
		notFound(w, "No user specified")
		return
	}
	p, err := loadPage(name[0])
	if err != nil {
		notFound(w, err.Error())
		return
	}
	src := "index"
	step := q.Get("step")
	guess := q.Get("guess")
	p.Guess = guess
	switch step {
	case p.Steps[0]:
		guessN, _ := strconv.Atoi(guess)
		if guessN == int(p.Offset) {
			p.CurrentStep = 1
			src = "good0"
		} else {
			p.CurrentStep = 0
			p.Guess = shift(p.CodedMessage, byte(guessN))
			src = "wrong"
		}
	case p.Steps[1]:
		if guess == passwords[0] {
			p.CurrentStep = 2
			src = "good1"
		} else {
			p.CurrentStep = 1
			src = "wrong"
		}
	case p.Steps[2]:
		if guess == passwords[1] {
			p.CurrentStep = 2
			src = "end"
		} else {
			p.CurrentStep = 2
			src = "wrong"
		}
	}
	log.Printf("[%s] step=%d guess=%s res=%s", name, p.CurrentStep, guess, src)
	renderTemplate(w, p, src)
}

func main() {
	log.Println("Parsing passwords file...")
	parsePasswords()

	log.Println("Starting server...")
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
