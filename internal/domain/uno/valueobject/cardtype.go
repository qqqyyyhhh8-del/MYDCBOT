package valueobject

type CardType string

const (
	CardTypeNumber  CardType = "Number"
	CardTypeSkip    CardType = "Skip"
	CardTypeReverse CardType = "Reverse"
	CardTypeDrawTwo CardType = "Drawtwo"
	CardTypeWild    CardType = "Wild"
	CardTypeWildDraw CardType = "WildDraw"
)

func (t CardType) IsActionCard() bool {
	return t == CardTypeSkip || t == CardTypeReverse || t == CardTypeDrawTwo
}

func (t CardType) IsWildCard() bool {
	return t == CardTypeWild || t == CardTypeWildDraw
}
