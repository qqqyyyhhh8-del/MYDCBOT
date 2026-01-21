package entity

import (
	"fmt"

	"github.com/user/dcminigames/internal/domain/uno/valueobject"
)

type Card struct {
	ID       string
	Color    valueobject.Color
	Type     valueobject.CardType
	Number   int
	ImageKey string
}

func NewNumberCard(color valueobject.Color, number int) *Card {
	return &Card{
		ID:       fmt.Sprintf("%s%d", color, number),
		Color:    color,
		Type:     valueobject.CardTypeNumber,
		Number:   number,
		ImageKey: fmt.Sprintf("%s%d.jpg", color, number),
	}
}

func NewActionCard(color valueobject.Color, cardType valueobject.CardType) *Card {
	return &Card{
		ID:       fmt.Sprintf("%s%s", color, cardType),
		Color:    color,
		Type:     cardType,
		Number:   -1,
		ImageKey: fmt.Sprintf("%s%s.jpg", color, cardType),
	}
}

func NewWildCard(cardType valueobject.CardType) *Card {
	return &Card{
		ID:       string(cardType),
		Color:    valueobject.ColorWild,
		Type:     cardType,
		Number:   -1,
		ImageKey: fmt.Sprintf("%s.jpg", cardType),
	}
}

func (c *Card) CanPlayOn(target *Card, currentColor valueobject.Color) bool {
	if c.Type.IsWildCard() {
		return true
	}
	if c.Color == currentColor {
		return true
	}
	if c.Type == valueobject.CardTypeNumber && target.Type == valueobject.CardTypeNumber && c.Number == target.Number {
		return true
	}
	if c.Type == target.Type && c.Type.IsActionCard() {
		return true
	}
	return false
}

func (c *Card) String() string {
	if c.Type == valueobject.CardTypeNumber {
		return fmt.Sprintf("%s%d", c.Color, c.Number)
	}
	return fmt.Sprintf("%s%s", c.Color, c.Type)
}
