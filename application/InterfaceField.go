package application

import (
	"fmt"
	"strings"

	gui "github.com/grupawp/warships-gui/v2"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

type descriptionField struct {
	firstLine  *gui.Text
	secondLine *gui.Text
	thirdLine  *gui.Text
}

func NewDescriptionFieldYour(desc string) descriptionField {
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

type LegendField struct {
	hitLegend  *gui.Text
	missLegend *gui.Text
	shipLegend *gui.Text
}

func NewLegendField() LegendField {
	hitLegend := gui.NewText(100, 10, "H-Hit the ship", nil)
	missLegend := gui.NewText(100, 12, "M-Miss the ship", nil)
	shipLegend := gui.NewText(100, 14, "S-Ship field", nil)
	return LegendField{hitLegend: hitLegend, missLegend: missLegend, shipLegend: shipLegend}
}

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

func NewShipLeftCountField() ShipLeftCountField {
	fourMastField := gui.NewText(100, 20, "Four mast ship left: 1", nil)
	threeMastField := gui.NewText(100, 22, "Three mast ship left: 2", nil)
	twoMastField := gui.NewText(100, 24, "Two mast ship left: 3", nil)
	oneMastField := gui.NewText(100, 26, "One mast ship left: 4", nil)

	return ShipLeftCountField{FourMastField: fourMastField, FourMastCount: 1, ThreeMastField: threeMastField, ThreeMastCount: 2, TwoMastField: twoMastField, TwoMastCount: 3, OneMastField: oneMastField, OneMastCount: 4}
}

func (s *ShipLeftCountField) UpdateFourMastCount() {
	countText := fmt.Sprintf("Four mast ship left: %d", s.FourMastCount)
	s.FourMastField.SetText(countText)
}

func (s *ShipLeftCountField) UpdateThreeMastCount() {
	countText := fmt.Sprintf("Three mast ship left: %d", s.ThreeMastCount)
	s.ThreeMastField.SetText(countText)
}

func (s *ShipLeftCountField) UpdateTwoMastCount() {
	countText := fmt.Sprintf("Two mast ship left: %d", s.TwoMastCount)
	s.TwoMastField.SetText(countText)
}

func (s *ShipLeftCountField) UpdateOneMastCount() {
	countText := fmt.Sprintf("One mast ship left: %d", s.OneMastCount)
	s.OneMastField.SetText(countText)
}
