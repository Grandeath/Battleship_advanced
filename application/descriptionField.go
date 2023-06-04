package application

import (
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
