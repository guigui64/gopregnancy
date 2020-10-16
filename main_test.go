package main

import (
	"fmt"
	"os"
	"testing"
)

type Message struct {
	Coded   string
	Decoded string
	Offset  byte
	Day     byte
	Month   byte
}

var cases = map[string]Message{
	"Zahi": {
		Coded:   "LODE PCSOXN PYCOGOC",
		Decoded: "BEST FRIEND FOREVER",
		Offset:  15,
		Day:     27, Month: 2,
	},
	"Thomas": {
		Coded:   "GEMIR HIGLEMRI",
		Decoded: "CAIEN DECHAINE",
		Offset:  21,
		Day:     11, Month: 8,
	},
	"Jerome": {
		Coded:   "RK IUDHUF JGTY RG CORRK",
		Decoded: "LE COWBOY DANS LA VILLE",
		Offset:  19,
		Day:     22, Month: 3,
	},
	"Yankel": {
		Coded:   "DFE SFQSX JGFIKYV",
		Decoded: "MON COACH SPORTIF",
		Offset:  9,
		Day:     1, Month: 7,
	},
	"Nicolas": {
		Coded:   "UROOLQJ LQ WKH GHHS",
		Decoded: "ROLLING IN THE DEEP",
		Offset:  22,
		Day:     5, Month: 5,
	},
}

var order = []string{"Nicolas", "Zahi", "Thomas", "Yankel", "Jerome"}

func TestMain(m *testing.M) {
	parsePasswords()
	os.Exit(m.Run())
}

func TestShift(t *testing.T) {
	for _, c := range cases {
		shifted := shift(c.Coded, c.Offset)
		if shifted != c.Decoded {
			t.Errorf("Expected %s but got %s instead", c.Decoded, shifted)
		}
	}
}

func TestEVGOrder(t *testing.T) {
	password := ""
	for _, buddy := range order {
		password += fmt.Sprintf("%02d", cases[buddy].Offset)
	}
	if password != passwords[0] {
		t.Errorf("Expected %s but got %s instead", passwords[0], password)
	}
}

func TestDates(t *testing.T) {
	password := ""
	for _, buddy := range order {
		c := cases[buddy]
		password += fmt.Sprintf("%02d", c.Offset+24-c.Day+7-c.Month)
	}
	if password != passwords[1] {
		t.Errorf("Expected %s but got %s instead", passwords[1], password)
	}
}
