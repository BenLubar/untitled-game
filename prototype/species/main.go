package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
)

var (
	seed  = flag.Int64("seed", 0, "random seed")
	count = flag.Int("count", 10, "number of species to generate")
)

func main() {
	flag.Parse()

	r := rand.New(rand.NewSource(*seed))

	for i := 0; i < *count; i++ {
		if i != 0 {
			fmt.Println()
		}
		s := NewSpecies(r)
		fmt.Println(s)
	}
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
	str := s.Body.String()
	return str[:len(str)-1]
}

type Size uint32

const (
	SizeMiniscule Size = 1 << (iota << 2)
	SizeTiny
	SizeSmall
	SizeMedium
	SizeLarge
	SizeHuge
	SizeGigantic
	SizeEnormous
)

func (s Size) String() string {
	switch {
	case s >= SizeEnormous:
		return "enormous"
	case s >= SizeGigantic:
		return "gigantic"
	case s >= SizeHuge:
		return "huge"
	case s >= SizeLarge:
		return "large"
	case s >= SizeMedium:
		return "medium-sized"
	case s >= SizeSmall:
		return "small"
	case s >= SizeTiny:
		return "tiny"
	case s >= SizeMiniscule:
		return "miniscule"
	default:
		return "nonexistent"
	}
}

func (s Size) Adverb(big, little string) string {
	switch {
	case s >= SizeEnormous:
		return "extremely " + big
	case s >= SizeGigantic:
		return "very " + big
	case s >= SizeHuge:
		return "quite " + big
	case s >= SizeLarge:
		return big
	case s >= SizeMedium:
		return little
	case s >= SizeSmall:
		return "quite " + little
	case s >= SizeTiny:
		return "very " + little
	case s >= SizeMiniscule:
		return "extremely " + little
	default:
		return "not"
	}
}

func randomSize(r *rand.Rand) Size {
	s := Size(1 << uint(intn(r, 10, 30, 60, 80, 80, 40, 30, 10)<<2))
	s += Size(r.Int63n(int64(s)<<(1<<2) - int64(s)))
	return s
}

type Body struct {
	Upper *Thorax
	Lower *Abdomen

	Size Size

	Separatable bool
}

func NewBody(r *rand.Rand) *Body {
	var b Body

	b.Separatable = intn(r, 7, 1) > 0
	b.Upper = NewThorax(r)
	b.Lower = NewAbdomen(r)
	b.Size = randomSize(r)

	return &b
}

func (b *Body) String() string {
	var buf []byte

	buf = append(buf, "It is "...)
	buf = append(buf, b.Size.String()...)
	buf = append(buf, ". "...)

	buf = append(buf, b.Upper.String()...)
	buf = append(buf, b.Lower.String()...)
	if b.Separatable {
		buf = append(buf, "Its body separates in the middle like an insect. "...)
	}
	return string(buf)
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

	Width  Size
	Length Size
}

func NewThorax(r *rand.Rand) *Thorax {
	var t Thorax

	headTypes := intn(r, 0, 30, 10, 5, 3, 2, 1, 1, 1, 1, 1)
	for i := 0; i < headTypes; i++ {
		t.Heads = append(t.Heads, NewHead(r))
	}

	limbTypes := intn(r, 2, 4, 1)
	for i := 0; i < limbTypes; i++ {
		t.Limbs = append(t.Limbs, NewLimb(r))
	}

	t.Width = randomSize(r)
	t.Length = randomSize(r)

	return &t
}

func (t *Thorax) String() string {
	var buf []byte

	buf = append(buf, "Its upper body is "...)
	buf = append(buf, t.Width.Adverb("wide", "narrow")...)
	buf = append(buf, " and "...)
	buf = append(buf, t.Length.Adverb("long", "short")...)
	buf = append(buf, ". "...)

	for _, h := range t.Heads {
		buf = append(buf, h.String()...)
	}

	for _, l := range t.Limbs {
		buf = append(buf, l.String()...)
	}

	return string(buf)
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

	Width  Size
	Length Size
}

func NewAbdomen(r *rand.Rand) *Abdomen {
	var a Abdomen

	limbTypes := intn(r, 2, 4, 1)
	for i := 0; i < limbTypes; i++ {
		a.Limbs = append(a.Limbs, NewLimb(r))
	}

	tailTypes := intn(r, 15, 10, 2, 1, 1)
	for i := 0; i < tailTypes; i++ {
		a.Tails = append(a.Tails, NewTail(r))
	}

	a.Width = randomSize(r)
	a.Length = randomSize(r)

	return &a
}

func (a *Abdomen) String() string {
	var buf []byte

	buf = append(buf, "Its lower body is "...)
	buf = append(buf, a.Width.Adverb("wide", "narrow")...)
	buf = append(buf, " and "...)
	buf = append(buf, a.Length.Adverb("long", "short")...)
	buf = append(buf, ". "...)

	for _, l := range a.Limbs {
		buf = append(buf, l.String()...)
	}

	for _, t := range a.Tails {
		buf = append(buf, t.String()...)
	}

	return string(buf)
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
	Count  uint16
	Width  Size
	Length Size
}

func NewLimb(r *rand.Rand) *Limb {
	var l Limb

	l.Joints = uint8(intn(r, 20, 30, 15, 10, 2, 1))
	l.Type = LimbType(r.Intn(int(limbTypeCount)))
	l.Count = uint16(intn(r, 0, 20, 80, 15, 60, 10, 40, 5, 20, 1, 4))

	l.Width = randomSize(r)
	l.Length = randomSize(r)

	return &l
}

func (l *Limb) String() string {
	var buf []byte

	buf = append(buf, "It has "...)
	buf = strconv.AppendUint(buf, uint64(l.Count), 10)
	buf = append(buf, " "...)
	buf = append(buf, l.Width.Adverb("wide", "narrow")...)
	buf = append(buf, ", "...)
	buf = append(buf, l.Length.Adverb("long", "short")...)
	buf = append(buf, " "...)
	buf = append(buf, limbTypeName[l.Type]...)
	buf = append(buf, "-limbs, each with "...)
	buf = strconv.AppendUint(buf, uint64(l.Joints), 10)
	buf = append(buf, " joints. "...)

	return string(buf)
}

func (l *Limb) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "Count:\n\t"...)
	buf = append(buf, indent...)
	buf = strconv.AppendUint(buf, uint64(l.Count), 10)

	buf = append(buf, '\n')

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
	Eyes  []*Eye
	Mouth *Mouth

	Width  Size
	Length Size
}

func NewHead(r *rand.Rand) *Head {
	var h Head

	eyeTypes := intn(r, 20, 40, 15, 10, 5, 1)
	for i := 0; i < eyeTypes; i++ {
		h.Eyes = append(h.Eyes, NewEye(r))
	}

	h.Mouth = NewMouth(r)

	h.Width = randomSize(r)
	h.Length = randomSize(r)

	return &h
}

func (h *Head) String() string {
	var buf []byte

	buf = append(buf, "It has a head that is "...)
	buf = append(buf, h.Width.Adverb("wide", "narrow")...)
	buf = append(buf, " and "...)
	buf = append(buf, h.Length.Adverb("long", "short")...)
	buf = append(buf, ". "...)
	buf = append(buf, h.Mouth.String()...)

	for _, e := range h.Eyes {
		buf = append(buf, e.String()...)
	}

	return string(buf)
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
	if first {
		first = false
	} else {
		buf = append(buf, '\n')
	}
	buf = append(buf, indent...)
	buf = append(buf, "Mouth:\n"...)
	buf = h.Mouth.Indent(buf, indent)

	return buf
}

type Tail struct {
	Count  uint16
	Width  Size
	Length Size
}

func NewTail(r *rand.Rand) *Tail {
	var t Tail

	t.Count = uint16(intn(r, 0, 30, 10, 6, 4, 3, 2, 1, 1, 5))
	t.Width = randomSize(r)
	t.Length = randomSize(r)

	return &t
}

func (t *Tail) String() string {
	var buf []byte

	buf = append(buf, "It has "...)
	buf = strconv.AppendUint(buf, uint64(t.Count), 10)
	buf = append(buf, " "...)
	buf = append(buf, t.Width.Adverb("wide", "narrow")...)
	buf = append(buf, ", "...)
	buf = append(buf, t.Length.Adverb("long", "short")...)
	buf = append(buf, " tails. "...)

	return string(buf)
}

func (t *Tail) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "Count:\n\t"...)
	buf = append(buf, indent...)
	buf = strconv.AppendUint(buf, uint64(t.Count), 10)

	return buf
}

type Eye struct {
	Count uint16

	Size Size
}

func NewEye(r *rand.Rand) *Eye {
	var e Eye

	e.Count = uint16(intn(r, 0, 10, 20, 2, 4, 1, 2))
	e.Size = randomSize(r)

	return &e
}

func (e *Eye) String() string {
	var buf []byte

	buf = append(buf, "On its head there are "...)
	buf = strconv.AppendUint(buf, uint64(e.Count), 10)
	buf = append(buf, " "...)
	buf = append(buf, e.Size.String()...)
	buf = append(buf, " eyes. "...)

	return string(buf)
}

func (e *Eye) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "Count:\n\t"...)
	buf = append(buf, indent...)
	buf = strconv.AppendUint(buf, uint64(e.Count), 10)

	return buf
}

type Mouth struct {
	Teeth []*Tooth

	Width  Size
	Length Size
}

func NewMouth(r *rand.Rand) *Mouth {
	var m Mouth

	toothTypes := intn(r, 5, 30, 20, 12, 5, 2, 1, 1, 1, 1)
	for i := 0; i < toothTypes; i++ {
		m.Teeth = append(m.Teeth, NewTooth(r))
	}

	m.Width = randomSize(r)
	m.Length = randomSize(r)

	return &m
}

func (m *Mouth) String() string {
	var buf []byte

	buf = append(buf, "It has a "...)
	buf = append(buf, m.Width.Adverb("wide", "narrow")...)
	buf = append(buf, ", "...)
	buf = append(buf, m.Length.Adverb("long", "short")...)
	buf = append(buf, " mouth containing "...)
	for i, t := range m.Teeth {
		if i > 0 && len(m.Teeth) > 2 {
			buf = append(buf, ", "...)
		}
		if i > 0 && i == len(m.Teeth)-1 {
			if len(m.Teeth) == 2 {
				buf = append(buf, " "...)
			}
			buf = append(buf, "and "...)
		}
		buf = append(buf, t.String()...)
	}
	if len(m.Teeth) == 0 {
		buf = append(buf, "no teeth"...)
	}
	buf = append(buf, ". "...)

	return string(buf)
}

func (m *Mouth) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	first := true

	for _, t := range m.Teeth {
		if first {
			first = false
		} else {
			buf = append(buf, '\n')
		}
		buf = append(buf, indent...)
		buf = append(buf, "Tooth:\n"...)
		buf = t.Indent(buf, indent)
	}

	return buf
}

type ToothType uint16

const (
	ToothMolar   ToothType = iota // Plant-eaters
	ToothCuspid                   // Fangs
	ToothIncisor                  // Meat-eaters
	ToothTusk                     // Elephants

	toothTypeCount
)

var toothTypeName = [toothTypeCount]string{
	ToothMolar:   "molars",
	ToothCuspid:  "fangs",
	ToothIncisor: "incisors",
	ToothTusk:    "tusks",
}

type Tooth struct {
	Type  ToothType
	Count uint16

	Size Size
}

func NewTooth(r *rand.Rand) *Tooth {
	var t Tooth

	t.Type = ToothType(r.Intn(int(toothTypeCount)))

	if t.Type == ToothTusk || t.Type == ToothCuspid {
		// fewer tusks and fangs
		t.Count = uint16(intn(r, 0, 10, 25, 10, 9, 7, 5, 3, 2, 1, 1, 1))
	} else {
		t.Count = uint16(intn(r, 0, 10, 10, 15, 15, 15, 15, 15, 15, 15, 15, 10, 10, 9, 8, 7, 5, 3, 2, 1, 1, 1))
	}

	t.Size = randomSize(r)

	return &t
}

func (t *Tooth) String() string {
	return strconv.FormatUint(uint64(t.Count), 10) + " " + t.Size.String() + " " + toothTypeName[t.Type]
}

func (t *Tooth) Indent(buf, indent []byte) []byte {
	indent = append(indent, '\t')

	buf = append(buf, indent...)
	buf = append(buf, "Count:\n\t"...)
	buf = append(buf, indent...)
	buf = strconv.AppendUint(buf, uint64(t.Count), 10)

	buf = append(buf, '\n')

	buf = append(buf, indent...)
	buf = append(buf, "Type:\n\t"...)
	buf = append(buf, indent...)
	buf = append(buf, toothTypeName[t.Type]...)

	return buf
}
