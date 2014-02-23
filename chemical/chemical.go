package chemical

type Chemical uint16

const (
	ChemAloe Chemical = iota
	ChemVitriol
	ChemHeparin
	ChemNepeta

	chemicalCount
)

type ChemicalInfo struct {
	// Toxicity limits the amount of a chemical that can be safely consumed.
	Toxicity int8

	// Healing increases the rate at which wounds repair themselves.
	Healing int8

	// Venom increases the chance of poisoning an enemy in combat.
	Venom int8

	// Bleed increases the chance of causing an enemy to bleed in combat.
	Bleed int8

	// Stun increases the chance of stunning an enemy in combat.
	Stun int8
}

var ChemInfo = [chemicalCount]ChemicalInfo{
	ChemAloe: {
		Toxicity: 10,
		Healing:  10,
	},
	ChemVitriol: {
		Toxicity: 25,
		Venom:    15,
	},
	ChemHeparin: {
		Toxicity: 15,
		Bleed:    10,
	},
	ChemNepeta: {
		Toxicity: 30,
		Stun:     15,
	},
}
