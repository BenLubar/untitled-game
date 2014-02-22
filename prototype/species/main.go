package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
)

var seed = flag.Int64("seed", 0, "random seed")

func main() {
	flag.Parse()

	r := rand.New(rand.NewSource(*seed))
	s := NewSpecies(r)
	fmt.Println(s)
}

func intn(r *rand.Rand, chances ...int) int {
	n := 0
	for _, c := range chances {
		n += c
	}
	n = r.Intn(n)
	for i, c := range chances {
		n -= c
		if n < 0 {
			return i
		}
	}
	panic("unreachable")
}

type Species struct {
	Body *Body
}

func NewSpecies(r *rand.Rand) *Species {
	var s Species

	s.Body = NewBody(r)

	return &s
}

func (s *Species) String() string {
	return "Body:\n" + string(s.Body.Indent(nil, nil))
}

type Body struct {
	Upper *Thorax
	Lower *Abdomen

	Separatable bool
}

func NewBody(r *rand.Rand) *Body {
	var b Body

	b.Separatable = intn(r, 7, 1) > 0
	b.Upper = NewThorax(r)
	b.Lower = NewAbdomen(r)

	return &b
}

func (b *Body) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "Upper:\n"...)
	buf = b.Upper.Indent(buf, indent)

	buf = append(buf, '\n')

	buf = append(buf, indent...)
	buf = append(buf, "Lower:\n"...)
	buf = b.Lower.Indent(buf, indent)

	buf = append(buf, '\n')

	buf = append(buf, indent...)
	buf = append(buf, "Separatable:\n\t"...)
	buf = append(buf, indent...)
	buf = strconv.AppendBool(buf, b.Separatable)

	return buf
}

type Thorax struct {
	Heads []*Head
	Limbs []*Limb
}

func NewThorax(r *rand.Rand) *Thorax {
	var t Thorax

	headTypes := intn(r, 0, 30, 10, 5, 3, 2, 1, 1, 1, 1, 1)
	for i := 0; i < headTypes; i++ {
		t.Heads = append(t.Heads, NewHead(r))
	}

	limbTypes := intn(r, 2, 4, 1)
	for i := 0; i < limbTypes; i++ {
		l := NewLimb(r)

		limbCount := intn(r, 0, 20, 80, 15, 60, 10, 40, 5, 20, 1, 4)
		for j := 0; j < limbCount; j++ {
			t.Limbs = append(t.Limbs, l)
		}
	}

	return &t
}

func (t *Thorax) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	first := true

	for _, h := range t.Heads {
		if first {
			first = false
		} else {
			buf = append(buf, '\n')
		}
		buf = append(buf, indent...)
		buf = append(buf, "Head:\n"...)
		buf = h.Indent(buf, indent)
	}
	for _, l := range t.Limbs {
		if first {
			first = false
		} else {
			buf = append(buf, '\n')
		}
		buf = append(buf, indent...)
		buf = append(buf, "Limb:\n"...)
		buf = l.Indent(buf, indent)
	}

	return buf
}

type Abdomen struct {
	Tails []*Tail
	Limbs []*Limb
}

func NewAbdomen(r *rand.Rand) *Abdomen {
	var a Abdomen

	limbTypes := intn(r, 2, 4, 1)
	for i := 0; i < limbTypes; i++ {
		l := NewLimb(r)

		limbCount := intn(r, 0, 20, 80, 15, 60, 10, 40, 5, 20, 1, 4)
		for j := 0; j < limbCount; j++ {
			a.Limbs = append(a.Limbs, l)
		}
	}

	tailTypes := intn(r, 15, 10, 2, 1, 1)
	for i := 0; i < tailTypes; i++ {
		t := NewTail(r)

		tailCount := intn(r, 0, 30, 10, 6, 4, 3, 2, 1, 1, 5)
		for j := 0; j < tailCount; j++ {
			a.Tails = append(a.Tails, t)
		}
	}

	return &a
}

func (a *Abdomen) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	first := true

	for _, l := range a.Limbs {
		if first {
			first = false
		} else {
			buf = append(buf, '\n')
		}
		buf = append(buf, indent...)
		buf = append(buf, "Limb:\n"...)
		buf = l.Indent(buf, indent)
	}

	for _, t := range a.Tails {
		if first {
			first = false
		} else {
			buf = append(buf, '\n')
		}
		buf = append(buf, indent...)
		buf = append(buf, "Tail:\n"...)
		buf = t.Indent(buf, indent)
	}

	return buf
}

type LimbType uint16

const (
	LimbPaw LimbType = iota
	LimbHand
	LimbFoot
	LimbTalon
	LimbHoof
	LimbTentacle
	LimbWing

	limbTypeCount
)

var limbTypeName = [limbTypeCount]string{
	LimbPaw:      "paw",
	LimbHand:     "hand",
	LimbFoot:     "foot",
	LimbTalon:    "talon",
	LimbHoof:     "hoof",
	LimbTentacle: "tentacle",
	LimbWing:     "wing",
}

type Limb struct {
	Joints uint8
	Type   LimbType
}

func NewLimb(r *rand.Rand) *Limb {
	var l Limb

	l.Joints = uint8(intn(r, 20, 30, 7, 2, 1))
	l.Type = LimbType(r.Intn(int(limbTypeCount)))

	return &l
}

func (l *Limb) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "Joints:\n\t"...)
	buf = append(buf, indent...)
	buf = strconv.AppendUint(buf, uint64(l.Joints), 10)

	buf = append(buf, '\n')

	buf = append(buf, indent...)
	buf = append(buf, "Type:\n\t"...)
	buf = append(buf, indent...)
	buf = append(buf, limbTypeName[l.Type]...)

	return buf
}

type Head struct {
	Eyes []*Eye
}

func NewHead(r *rand.Rand) *Head {
	var h Head

	eyeTypes := intn(r, 10, 20, 15, 10, 5, 1)
	for i := 0; i < eyeTypes; i++ {
		e := NewEye(r)

		eyeCount := intn(r, 0, 10, 20, 2, 4, 1, 2)
		for j := 0; j < eyeCount; j++ {
			h.Eyes = append(h.Eyes, e)
		}
	}

	return &h
}

func (h *Head) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	first := true

	for _, e := range h.Eyes {
		if first {
			first = false
		} else {
			buf = append(buf, '\n')
		}
		buf = append(buf, indent...)
		buf = append(buf, "Eye:\n"...)
		buf = e.Indent(buf, indent)
	}

	return buf
}

type Tail struct {
}

func NewTail(r *rand.Rand) *Tail {
	var t Tail

	return &t
}

func (t *Tail) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	return buf
}

type Eye struct {
}

func NewEye(r *rand.Rand) *Eye {
	var e Eye

	return &e
}

func (e *Eye) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	return buf
}
