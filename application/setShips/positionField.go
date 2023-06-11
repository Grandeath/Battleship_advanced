package setships

import gui "github.com/grupawp/warships-gui/v2"

type PositionField struct {
	currentShipFIeld *gui.Text
	shipCount        int
	currentMastCount int
}

func NewPositionField() PositionField {
	currentShipField := gui.NewText(50, 6, "Position 1 four mast ship", nil)
	return PositionField{currentShipFIeld: currentShipField, currentMastCount: 4, shipCount: 0}
}

func (p *PositionField) NextShip() {
	p.shipCount++
	switch p.shipCount {
	case 1:
		p.currentShipFIeld.SetText("Position 2 three mast ships")
		p.currentMastCount = 3
	case 2:
		p.currentShipFIeld.SetText("Position 1 three mast ship")
		p.currentMastCount = 3
	case 3:
		p.currentShipFIeld.SetText("Position 3 two mast ships")
		p.currentMastCount = 2
	case 4:
		p.currentShipFIeld.SetText("Position 2 two mast ships")
		p.currentMastCount = 2
	case 5:
		p.currentShipFIeld.SetText("Position 1 two mast ship")
		p.currentMastCount = 2
	case 6:
		p.currentShipFIeld.SetText("Position 4 one mast ships")
		p.currentMastCount = 1
	case 7:
		p.currentShipFIeld.SetText("Position 3 one mast ships")
		p.currentMastCount = 1
	case 8:
		p.currentShipFIeld.SetText("Position 2 one mast ships")
		p.currentMastCount = 1
	case 9:
		p.currentShipFIeld.SetText("Position 1 one mast ship")
		p.currentMastCount = 1
	case 10:
		p.currentShipFIeld.SetText("Done! You can leave click ctrl+f")
		p.currentMastCount = 0
	}
}
