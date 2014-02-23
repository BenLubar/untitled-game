// +build lojban

package language

import (
	"github.com/BenLubar/untitled-game/chemical"
)

func Number(i int) string {
	var buf []byte
	segment := func(n int64, comma bool) {
		n %= 1000
		digits := [10]string{"no", "pa", "re", "ci", "vo", "mu", "xa", "ze", "bi", "so"}
		if n >= 100 {
			buf = append(buf, digits[n/100]...)
		}
		if n >= 10 {
			buf = append(buf, digits[(n/10)%10]...)
		}
		if n > 0 {
			buf = append(buf, digits[n%10]...)
		}
		if comma {
			buf = append(buf, "ki'o"...)
		}
		buf = append(buf, ' ')
	}

	n := int64(i)

	if n < 0 {
		buf = append(buf, "ni'u "...)
		n = -n
	}

	if n >= 1000000000000000000 {
		segment(n/1000000000000000000, true)
	}
	if n >= 1000000000000000 {
		segment(n/1000000000000000, true)
	}
	if n >= 1000000000000 {
		segment(n/1000000000000, true)
	}
	if n >= 1000000000 {
		segment(n/1000000000, true)
	}
	if n >= 1000000 {
		segment(n/1000000, true)
	}
	if n >= 1000 {
		segment(n/1000, true)
	}
	if n != 0 {
		segment(n%1000, false)
	} else {
		buf = append(buf, "no "...)
	}
	return string(buf[:len(buf)-1])
}

func Ordinal(i int) string {
	return Number(i) + " moi"
}

func ChemName(c chemical.Chemical) string {
	switch c {
	case chemical.ChemAloe:
		return "sparalo'e"
	case chemical.ChemVitriol:
		return "slami"
	case chemical.ChemHeparin:
		return "flecu"
	case chemical.ChemNepeta:
		return "latfekspa"
	}
	panic(c)
}
