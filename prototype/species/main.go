package main

import (
	"flag"
	"fmt"
	"github.com/BenLubar/untitled-game/language"
	"math/rand"
)

var (
	seed  = flag.Int64("seed", 0, "random seed")
	skip  = flag.Int("skip", 0, "number of species to generate but not print")
	count = flag.Int("count", 10, "number of species to generate")
)

func main() {
	flag.Parse()

	r := rand.New(rand.NewSource(*seed))

	for i := 0; i < *skip; i++ {
		_ = NewSpecies(r)
	}

	for i := 0; i < *count; i++ {
		if i != 0 {
			fmt.Println()
			fmt.Println("--------------------------------------------------------------------------------")
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
	return s.Body.String()
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
		return "infinitesimally " + little
	}
}

func (s Size) Length() string {
	return s.Adverb("long", "short")
}

func (s Size) Width() string {
	return s.Adverb("wide", "narrow")
}

func randomSize(r *rand.Rand) Size {
	s := Size(1 << uint(intn(r, 10, 30, 60, 80, 80, 40, 30, 10)<<2))
	s += Size(r.Int63n(int64(s)<<(1<<2) - int64(s)))
	return s
}

type Texture uint16

const (
	TexSmoothLeather Texture = iota
	TexRoughLeather
	TexWrinkledLeather
	TexSmoothFur
	TexStripedFur
	TexSpottedFur
	TexSmoothScale
	TexRoughScale
	TexSmoothChitin
	TexRoughChitin
	TexSmoothWood
	TexRoughWood
	TexMetal
	TexRock
	TexSlime
	TexVapor

	texCount
)

var texName = [texCount]string{
	TexSmoothLeather:   "smooth leathery skin",
	TexRoughLeather:    "rough leathery skin",
	TexWrinkledLeather: "wrinkled leathery skin",
	TexSmoothFur:       "smooth fur",
	TexStripedFur:      "striped fur",
	TexSpottedFur:      "spotted fur",
	TexSmoothScale:     "smooth scales",
	TexRoughScale:      "rough scales",
	TexSmoothChitin:    "smooth chitin",
	TexRoughChitin:     "rough chitin",
	TexSmoothWood:      "smooth wood-like skin",
	TexRoughWood:       "rough wood-like skin",
	TexMetal:           "metalic skin",
	TexRock:            "rock-hard skin",
	TexSlime:           "slimy skin",
	TexVapor:           "vapor-like skin",
}

type Color uint16

const (
	ColorWhite Color = iota
	ColorGray
	ColorBlack
	ColorBeige
	ColorTan
	ColorBrown
	ColorPink
	ColorRed
	ColorMaroon
	ColorYellow
	ColorOrange
	ColorGreen
	ColorTurquoise
	ColorBlue
	ColorNavy
	ColorViolet
	ColorPurple
	ColorIndigo

	colorCount
)

var colorName = [colorCount]string{
	ColorWhite:     "white",
	ColorGray:      "gray",
	ColorBlack:     "black",
	ColorBeige:     "beige",
	ColorTan:       "tan",
	ColorBrown:     "brown",
	ColorPink:      "pink",
	ColorRed:       "red",
	ColorMaroon:    "maroon",
	ColorYellow:    "yellow",
	ColorOrange:    "orange",
	ColorGreen:     "green",
	ColorTurquoise: "turquoise",
	ColorBlue:      "blue",
	ColorNavy:      "navy",
	ColorViolet:    "violet",
	ColorPurple:    "purple",
	ColorIndigo:    "indigo",
}

type Body struct {
	Upper *Thorax
	Lower *Abdomen

	Size  Size
	Skin  Texture
	Color Color

	Separatable bool
}

func NewBody(r *rand.Rand) *Body {
	var b Body

	b.Separatable = intn(r, 7, 1) > 0
	b.Upper = NewThorax(r)
	b.Lower = NewAbdomen(r)
	b.Size = randomSize(r)
	b.Skin = Texture(r.Intn(int(texCount)))
	b.Color = Color(r.Intn(int(colorCount)))

	return &b
}

func (b *Body) String() string {
	buf := b.Append(nil, 0, 0)
	return string(buf[:len(buf)-1])
}

func (b *Body) Append(buf []byte, n, total int) []byte {
	buf = append(buf, "It is "...)
	buf = append(buf, b.Size.String()...)
	buf = append(buf, " with "...)
	buf = append(buf, colorName[b.Color]...)
	buf = append(buf, " "...)
	buf = append(buf, texName[b.Skin]...)
	buf = append(buf, ". "...)

	buf = b.Upper.Append(buf, 0, 0)
	buf = b.Lower.Append(buf, 0, 0)

	if b.Separatable {
		buf = append(buf, "Its body separates in the middle like an insect. "...)
	}

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

func (t *Thorax) Append(buf []byte, n, total int) []byte {
	buf = append(buf, "Its upper body is "...)
	buf = append(buf, t.Length.Length()...)
	buf = append(buf, " and "...)
	buf = append(buf, t.Width.Width()...)
	buf = append(buf, ". "...)

	for i, h := range t.Heads {
		buf = h.Append(buf, i+1, len(t.Heads))
	}

	for i, l := range t.Limbs {
		buf = l.Append(buf, i+1, len(t.Limbs))
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

func (a *Abdomen) Append(buf []byte, n, total int) []byte {
	buf = append(buf, "Its lower body is "...)
	buf = append(buf, a.Length.Length()...)
	buf = append(buf, " and "...)
	buf = append(buf, a.Width.Width()...)
	buf = append(buf, ". "...)

	for i, l := range a.Limbs {
		buf = l.Append(buf, i+1, len(a.Limbs))
	}

	for i, t := range a.Tails {
		buf = t.Append(buf, i+1, len(a.Tails))
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

func (l *Limb) Append(buf []byte, n, total int) []byte {
	buf = append(buf, "It has "...)
	buf = append(buf, language.Number(int(l.Count))...)
	buf = append(buf, " "...)
	buf = append(buf, l.Length.Length()...)
	buf = append(buf, ", "...)
	buf = append(buf, l.Width.Width()...)
	buf = append(buf, " "...)
	buf = append(buf, limbTypeName[l.Type]...)
	if l.Count == 1 {
		buf = append(buf, "-limb with "...)
	} else {
		buf = append(buf, "-limbs, each with "...)
	}
	if l.Joints == 0 {
		buf = append(buf, "no joints. "...)
	} else if l.Joints == 1 {
		buf = append(buf, "one joint. "...)
	} else {
		buf = append(buf, language.Number(int(l.Joints))...)
		buf = append(buf, " joints. "...)
	}

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

func (h *Head) Append(buf []byte, n, total int) []byte {
	if total == 1 {
		buf = append(buf, "Its head is "...)
	} else {
		buf = append(buf, "Its "...)
		buf = append(buf, language.Ordinal(n)...)
		buf = append(buf, " head is "...)
	}
	buf = append(buf, h.Length.Length()...)
	buf = append(buf, " and "...)
	buf = append(buf, h.Width.Width()...)
	buf = append(buf, ". "...)
	buf = h.Mouth.Append(buf, 0, 0)

	for i, e := range h.Eyes {
		buf = e.Append(buf, i+1, len(h.Eyes))
	}

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

func (t *Tail) Append(buf []byte, n, total int) []byte {
	buf = append(buf, "It has "...)
	buf = append(buf, language.Number(int(t.Count))...)
	buf = append(buf, " "...)
	buf = append(buf, t.Length.Length()...)
	buf = append(buf, ", "...)
	buf = append(buf, t.Width.Width()...)
	if t.Count == 1 {
		buf = append(buf, " tail. "...)
	} else {
		buf = append(buf, " tails. "...)
	}

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

func (e *Eye) Append(buf []byte, n, total int) []byte {
	if e.Count == 1 {
		buf = append(buf, "On its head there is "...)
	} else {
		buf = append(buf, "On its head there are "...)
	}
	buf = append(buf, language.Number(int(e.Count))...)
	buf = append(buf, " "...)
	buf = append(buf, e.Size.String()...)
	if e.Count == 1 {
		buf = append(buf, " eye. "...)
	} else {
		buf = append(buf, " eyes. "...)
	}

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

func (m *Mouth) Append(buf []byte, n, total int) []byte {
	buf = append(buf, "It has a "...)
	buf = append(buf, m.Length.Length()...)
	buf = append(buf, ", "...)
	buf = append(buf, m.Width.Width()...)
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
		buf = t.Append(buf, i+1, len(m.Teeth))
	}
	if len(m.Teeth) == 0 {
		buf = append(buf, "no teeth"...)
	}
	buf = append(buf, ". "...)

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
	ToothMolar:   "molar",
	ToothCuspid:  "fang",
	ToothIncisor: "incisor",
	ToothTusk:    "tusk",
}

var teethTypeName = [toothTypeCount]string{
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
		t.Count = uint16(intn(r, 1, 10, 10, 15, 15, 15, 15, 15, 15, 15, 15, 10, 10, 9, 9, 8, 8, 7, 7, 5, 5, 3, 3, 2, 2, 1, 1, 1, 1, 1, 1))
		if t.Count == 0 {
			t.Count = uint16(r.Intn(5000) + 1)
		}
	}

	t.Size = randomSize(r)

	return &t
}

func (t *Tooth) Append(buf []byte, n, total int) []byte {
	buf = append(buf, language.Number(int(t.Count))...)
	buf = append(buf, " "...)
	buf = append(buf, t.Size.String()...)
	buf = append(buf, " "...)
	if t.Count == 1 {
		buf = append(buf, toothTypeName[t.Type]...)
	} else {
		buf = append(buf, teethTypeName[t.Type]...)
	}
	return buf
}
