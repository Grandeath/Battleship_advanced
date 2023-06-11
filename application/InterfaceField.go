// application implements logic of the game
package application

import (
	"fmt"
	"math"
	"strings"

	gui "github.com/grupawp/warships-gui/v2"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

// descriptionField contain description of players and devide it to three lines
type descriptionField struct {
	firstLine  *gui.Text
	secondLine *gui.Text
	thirdLine  *gui.Text
}

// NewDescriptionFielYours take player description string and return it devided to up to three lines and create descriptionField
func NewDescriptionFieldYours(desc string) descriptionField {
	wrappedText := wordwrap.WrapString(desc, 45)

	lines := strings.Split(wrappedText, "\n")

	switch len(lines) {
	case 3:
		firstLane := gui.NewText(1, 28, lines[0], nil)
		secondLane := gui.NewText(1, 30, lines[1], nil)
		thirdLane := gui.NewText(1, 32, lines[2], nil)
		return descriptionField{firstLine: firstLane, secondLine: secondLane, thirdLine: thirdLane}
	case 2:
		firstLane := gui.NewText(1, 28, lines[0], nil)
		secondLane := gui.NewText(1, 30, lines[1], nil)
		thirdLane := gui.NewText(1, 32, "", nil)
		return descriptionField{firstLine: firstLane, secondLine: secondLane, thirdLine: thirdLane}
	case 1:
		firstLane := gui.NewText(1, 28, lines[0], nil)
		secondLane := gui.NewText(1, 30, "", nil)
		thirdLane := gui.NewText(1, 32, "", nil)
		return descriptionField{firstLine: firstLane, secondLine: secondLane, thirdLine: thirdLane}
	default:
		firstLane := gui.NewText(1, 28, "", nil)
		secondLane := gui.NewText(1, 30, "", nil)
		thirdLane := gui.NewText(1, 32, "", nil)
		return descriptionField{firstLine: firstLane, secondLine: secondLane, thirdLine: thirdLane}
	}
}

// NewDescriptionFielEnemy take enemy description string and return it devided to up to three lines and create descriptionField in appriopriate coordinates
func NewDescriptionFieldEnemy(desc string) descriptionField {
	wrappedText := wordwrap.WrapString(desc, 45)
	lines := strings.Split(wrappedText, "\n")

	switch len(lines) {
	case 3:
		firstLane := gui.NewText(50, 28, lines[0], nil)
		secondLane := gui.NewText(50, 30, lines[1], nil)
		thirdLane := gui.NewText(50, 32, lines[2], nil)
		return descriptionField{firstLine: firstLane, secondLine: secondLane, thirdLine: thirdLane}
	case 2:
		firstLane := gui.NewText(50, 28, lines[0], nil)
		secondLane := gui.NewText(50, 30, lines[1], nil)
		thirdLane := gui.NewText(50, 32, "", nil)
		return descriptionField{firstLine: firstLane, secondLine: secondLane, thirdLine: thirdLane}
	case 1:
		firstLane := gui.NewText(50, 28, lines[0], nil)
		secondLane := gui.NewText(50, 30, "", nil)
		thirdLane := gui.NewText(50, 32, "", nil)
		return descriptionField{firstLine: firstLane, secondLine: secondLane, thirdLine: thirdLane}
	default:
		firstLane := gui.NewText(50, 28, "", nil)
		secondLane := gui.NewText(50, 30, "", nil)
		thirdLane := gui.NewText(50, 32, "", nil)
		return descriptionField{firstLine: firstLane, secondLine: secondLane, thirdLine: thirdLane}
	}
}

// LegendField create legend explaining fields on the board
// hitLegend explain hit field
// missLegend explain miss field
// shipLegend explain ship field
type LegendField struct {
	hitLegend  *gui.Text
	missLegend *gui.Text
	shipLegend *gui.Text
}

// NewLegendField creates LegenField in appriopriate coordinates
func NewLegendField() LegendField {
	hitLegend := gui.NewText(100, 10, "H-Hit the ship", nil)
	missLegend := gui.NewText(100, 12, "M-Miss the ship", nil)
	shipLegend := gui.NewText(100, 14, "S-Ship field", nil)
	return LegendField{hitLegend: hitLegend, missLegend: missLegend, shipLegend: shipLegend}
}

// ShipLeftCountField creating fields teling ship left to strike
// FourMastField text field telling number of four mast ship left
// FourMastCount  number telling how many four mast ship left to strike
// ThreeMastField text field telling number of three mast ship left
// ThreeMastCount number telling how many three mast ship left to strike
// TwoMastField text field telling number of two mast ship left
// TwoMastCount number telling how many two mast ship left to strike
// OneMastField text field telling number of one mast ship left
// OneMastCount number telling how many one mast ship left to strike
type ShipLeftCountField struct {
	FourMastField  *gui.Text
	FourMastCount  int
	ThreeMastField *gui.Text
	ThreeMastCount int
	TwoMastField   *gui.Text
	TwoMastCount   int
	OneMastField   *gui.Text
	OneMastCount   int
}

// NewShipLeftCountField creates ShipLeftCountField on appropriate coordinates and starting values of ships number
func NewShipLeftCountField() ShipLeftCountField {
	fourMastField := gui.NewText(100, 20, "Four mast ship left: 1", nil)
	threeMastField := gui.NewText(100, 22, "Three mast ship left: 2", nil)
	twoMastField := gui.NewText(100, 24, "Two mast ship left: 3", nil)
	oneMastField := gui.NewText(100, 26, "One mast ship left: 4", nil)

	return ShipLeftCountField{FourMastField: fourMastField, FourMastCount: 1, ThreeMastField: threeMastField, ThreeMastCount: 2, TwoMastField: twoMastField, TwoMastCount: 3, OneMastField: oneMastField, OneMastCount: 4}
}

// UpdateFourMastField update FourMastField
func (s *ShipLeftCountField) UpdateFourMastField() {
	countText := fmt.Sprintf("Four mast ship left: %d", s.FourMastCount)
	s.FourMastField.SetText(countText)
}

// UpdateThreeMastField update ThreeMastField
func (s *ShipLeftCountField) UpdateThreeMastField() {
	countText := fmt.Sprintf("Three mast ship left: %d", s.ThreeMastCount)
	s.ThreeMastField.SetText(countText)
}

// UpdateTwoMastField update TwoMastField
func (s *ShipLeftCountField) UpdateTwoMastField() {
	countText := fmt.Sprintf("Two mast ship left: %d", s.TwoMastCount)
	s.TwoMastField.SetText(countText)
}

// UpdateOneMastField update OneMastField
func (s *ShipLeftCountField) UpdateOneMastField() {
	countText := fmt.Sprintf("One mast ship left: %d", s.OneMastCount)
	s.OneMastField.SetText(countText)
}

// Accuracy show and monitor you percentage of hit shot / total number shot
type Accuracy struct {
	accuracyField *gui.Text
	ShotNumber    uint16
	HitNumber     uint16
}

// NewAccuracyField creates Accuracy in appropriate coordinates
func NewAccuracyField() Accuracy {
	accuracyField := gui.NewText(100, 10, "Accuracy: 0 Number of shots: 0 hit: 0", nil)
	return Accuracy{accuracyField: accuracyField}
}

// updateField update accuracyField with current percentage of hit shots
func (a *Accuracy) updateField() {
	percentAccuracy := (float64(a.HitNumber) / float64(a.ShotNumber))
	percentAccuracy = math.Round(percentAccuracy*100) / 100.0
	countText := fmt.Sprintf("Accuracy: %.2f Number of shots: %d hit: %d", percentAccuracy, a.ShotNumber, a.HitNumber)
	a.accuracyField.SetText(countText)
}
