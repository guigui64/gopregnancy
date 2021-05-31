// Enigma website to announce our pregnancy to my buddies
// Copyright 2020 Guillaume Comte

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// parsed passwords
var passwords []string

// parsed speech
var speech []string

// html templates
var templates = template.Must(template.ParseFiles("step0.html", "step1.html", "step2.html", "step3.html", "step4.html", "step5.html", "wrong.html", "404.html", "end.html"))

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

func parseFiles() {
	content, err := ioutil.ReadFile("passwords.txt")
	if err != nil {
		log.Panic(err)
	}
	passwords = strings.Split(string(content), "\n")
	passwords = passwords[:len(passwords)-1]
	content, err = ioutil.ReadFile("discours.txt")
	if err != nil {
		log.Panic(err)
	}
	speech = strings.Split(string(content), "\n")
	speech = speech[:len(speech)-1]
}

// Shift letters
func shift(s string, offset byte) string {
	b := []byte(s)
	for i, bi := range b {
		if 'A' <= bi && bi <= 'Z' {
			b[i] += offset
			b[i] = (b[i]-'A')%('Z'-'A'+1) + 'A'
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
	src := "index"
	q := r.URL.Query()
	name, ok := q["user"]
	if !ok {
		http.ServeFile(w, r, "index.html")
		return
	}
	p, err := loadPage(name[0])
	if err != nil {
		notFound(w, err.Error())
		return
	}
	src = "step0"
	step := q.Get("step")
	guess := q.Get("guess")
	p.Guess = guess
	wrong := false
	switch step {
	case p.Steps[0]:
		guessN, _ := strconv.Atoi(guess)
		if guessN == int(p.Offset) {
			p.CurrentStep = 1
		} else {
			p.CurrentStep = 0
			p.Guess = shift(p.CodedMessage, byte(guessN))
			wrong = true
		}
	case p.Steps[1]:
		if guess == passwords[0] {
			p.CurrentStep = 2
		} else {
			p.CurrentStep = 1
			wrong = true
		}
	case p.Steps[2]:
		if guess == passwords[1] {
			p.CurrentStep = 3
		} else {
			p.CurrentStep = 2
			wrong = true
		}
	case p.Steps[3]:
		p.CurrentStep = 4
		for i, word := range speech {
			g := strings.TrimSpace(q.Get(strconv.Itoa(i + 1)))
			guess += g + "/"
			if g != word {
				p.CurrentStep = 3
				wrong = true
			}
		}
		p.Guess = guess
	case p.Steps[4]:
		p.CurrentStep = 5
	case p.Steps[5]:
		guess = q.Get("boy") + "/" + q.Get("girl")
		src = "end"
	}
	if wrong {
		src = "wrong"
	} else if src != "end" {
		src = fmt.Sprintf("step%d", p.CurrentStep)
	}
	log.Printf("[%s] step=%d guess=%s res=%s", p.User, p.CurrentStep, guess, src)
	renderTemplate(w, p, src)
}

func main() {
	log.Println("Parsing files...")
	parseFiles()

	log.Println("Starting server...")
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
