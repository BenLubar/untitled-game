package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

func main() {
	r := rand.New(rand.NewSource(0))
	s := NewSpecies(r)
	fmt.Println(s)
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

	b.Separatable = r.Intn(3) != 0
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

	headTypes := r.Intn(3) + 1
	for i := 0; i < headTypes; i++ {
		h := NewHead(r)

		headCount := r.Intn(7) + 1
		for j := 0; j < headCount; j++ {
			t.Heads = append(t.Heads, h)
		}
	}

	limbTypes := r.Intn(2)
	for i := 0; i < limbTypes; i++ {
		l := NewLimb(r)

		limbCount := (r.Intn(2) + 1) * 2
		if r.Intn(5) == 0 {
			limbCount--
		}
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

	limbTypes := r.Intn(2)
	for i := 0; i < limbTypes; i++ {
		l := NewLimb(r)

		limbCount := (r.Intn(3) + 1) * 2
		if r.Intn(5) == 0 {
			limbCount--
		}
		for j := 0; j < limbCount; j++ {
			a.Limbs = append(a.Limbs, l)
		}
	}

	tailTypes := r.Intn(4)
	for i := 0; i < tailTypes; i++ {
		t := NewTail(r)

		tailCount := r.Intn(9) + 1
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

	l.Joints = uint8(r.Intn(5))
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

	eyeTypes := r.Intn(2)
	for i := 0; i < eyeTypes; i++ {
		e := NewEye(r)

		eyeCount := (r.Intn(3) + 1) * 2
		if r.Intn(5) == 0 {
			eyeCount--
		}
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
