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

// IsEndingForbidden 检查该牌是否不能作为最后一张牌打出
// 功能牌（+2、+4、Skip、Reverse）不能作为最后一张牌
func (t CardType) IsEndingForbidden() bool {
	return t == CardTypeSkip || t == CardTypeReverse || t == CardTypeDrawTwo || t == CardTypeWildDraw
}
