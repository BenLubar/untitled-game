// +build english !lojban

package language

import (
	"github.com/BenLubar/untitled-game/chemical"
	"strings"
)

func Number(i int) string {
	var buf []byte
	var segment func(int64)
	segment = func(n int64) {
		if n >= 100 {
			segment(n / 100)
			buf = append(buf, "hundred "...)
			n %= 100
			if n != 0 {
				segment(n)
			}
		} else if n >= 20 {
			buf = append(buf, [10]string{"", "", "twenty ", "thirty ", "forty ", "fifty ", "sixty ", "seventy ", "eighty ", "ninety "}[n/10]...)
			n %= 10
			if n != 0 {
				segment(n)
			}
		} else {
			buf = append(buf, [20]string{"zero ", "one ", "two ", "three ", "four ", "five ", "six ", "seven ", "eight ", "nine ", "ten ", "eleven ", "twelve ", "thirteen ", "fourteen ", "fifteen ", "sixteen ", "seventeen ", "eighteen ", "nineteen "}[n]...)
		}
	}

	n := int64(i)

	if n < 0 {
		buf = append(buf, "negative "...)
		n = -n
	}

	if n >= 1000000000000000000 {
		segment(n / 1000000000000000000)
		buf = append(buf, "quintillion "...)
		n %= 1000000000000000000
	}
	if n >= 1000000000000000 {
		segment(n / 1000000000000000)
		buf = append(buf, "quadrillion "...)
		n %= 1000000000000000
	}
	if n >= 1000000000000 {
		segment(n / 1000000000000)
		buf = append(buf, "trillion "...)
		n %= 1000000000000
	}
	if n >= 1000000000 {
		segment(n / 1000000000)
		buf = append(buf, "billion "...)
		n %= 1000000000
	}
	if n >= 1000000 {
		segment(n / 1000000)
		buf = append(buf, "million "...)
		n %= 1000000
	}
	if n >= 1000 {
		segment(n / 1000)
		buf = append(buf, "thousand "...)
		n %= 1000
	}
	if n != 0 || len(buf) == 0 {
		segment(n)
	}
	return string(buf[:len(buf)-1])
}

func Ordinal(i int) string {
	s := Number(i)

	if strings.HasSuffix(s, "ty") {
		return s[:len(s)-len("ty")] + "tieth"
	}
	if strings.HasSuffix(s, "t") {
		return s[:len(s)-len("t")] + "th"
	}
	if strings.HasSuffix(s, "one") {
		return s[:len(s)-len("one")] + "first"
	}
	if strings.HasSuffix(s, "two") {
		return s[:len(s)-len("two")] + "second"
	}
	if strings.HasSuffix(s, "three") {
		return s[:len(s)-len("three")] + "third"
	}
	if strings.HasSuffix(s, "five") {
		return s[:len(s)-len("five")] + "fifth"
	}
	if strings.HasSuffix(s, "nine") {
		return s[:len(s)-len("nine")] + "ninth"
	}
	return s + "th"
}

func ChemName(c chemical.Chemical) string {
	switch c {
	case chemical.ChemAloe:
		return "aloe"
	case chemical.ChemVitriol:
		return "vitriol"
	case chemical.ChemHeparin:
		return "heparin"
	case chemical.ChemNepeta:
		return "nepeta"
	}
	panic(c)
}
