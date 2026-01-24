package entity

import (
	"errors"
	"math/rand"
	"time"

	"github.com/user/dcminigames/internal/domain/uno/valueobject"
)

type GameState string

const (
	GameStateWaiting           GameState = "waiting"
	GameStatePlaying           GameState = "playing"
	GameStateFinished          GameState = "finished"
	GameStateWaitingChallenge  GameState = "waiting_challenge"  // 等待+4质疑
	GameStateWaitingUnoButton  GameState = "waiting_uno_button" // 等待UNO按钮
)

type Game struct {
	ID            string
	ChannelID     string
	Players       []*Player
	Deck          []*Card
	DiscardPile   []*Card
	CurrentPlayer int
	Direction     int
	CurrentColor  valueobject.Color
	State         GameState
	Winner        *Player
	CreatedAt     time.Time

	// 新增字段
	HasDrawnThisTurn    bool      // 本回合是否已摸牌
	PendingWildDraw     bool      // 是否有待处理的+4
	WildDrawPlayer      string    // 打出+4的玩家ID
	WildDrawVictim      string    // 被+4的玩家ID
	UnoButtonActive     bool      // UNO按钮是否激活
	UnoPlayerID         string    // 需要喊UNO的玩家ID
	UnoButtonPressedBy  string    // 按下UNO按钮的玩家ID
	UnoButtonTime       time.Time // UNO按钮激活时间
}

func NewGame(id, channelID string) *Game {
	return &Game{
		ID:               id,
		ChannelID:        channelID,
		Players:          make([]*Player, 0),
		Deck:             make([]*Card, 0),
		DiscardPile:      make([]*Card, 0),
		CurrentPlayer:    0,
		Direction:        1,
		State:            GameStateWaiting,
		CreatedAt:        time.Now(),
		HasDrawnThisTurn: false,
	}
}

func (g *Game) AddPlayer(player *Player) error {
	if g.State != GameStateWaiting {
		return errors.New("游戏已开始")
	}
	if len(g.Players) >= 10 {
		return errors.New("玩家已满")
	}
	for _, p := range g.Players {
		if p.ID == player.ID {
			return errors.New("已在游戏中")
		}
	}
	g.Players = append(g.Players, player)
	return nil
}

func (g *Game) Start() error {
	if g.State != GameStateWaiting {
		return errors.New("游戏已开始")
	}
	if len(g.Players) < 2 {
		return errors.New("至少需要2名玩家")
	}
	g.initDeck()
	g.shuffleDeck()
	g.dealCards()
	g.flipFirstCard()
	g.State = GameStatePlaying
	return nil
}

func (g *Game) initDeck() {
	colors := []valueobject.Color{
		valueobject.ColorRed,
		valueobject.ColorBlue,
		valueobject.ColorGreen,
		valueobject.ColorYellow,
	}
	for _, color := range colors {
		g.Deck = append(g.Deck, NewNumberCard(color, 0))
		for num := 1; num <= 9; num++ {
			g.Deck = append(g.Deck, NewNumberCard(color, num))
			g.Deck = append(g.Deck, NewNumberCard(color, num))
		}
		for i := 0; i < 2; i++ {
			g.Deck = append(g.Deck, NewActionCard(color, valueobject.CardTypeSkip))
			g.Deck = append(g.Deck, NewActionCard(color, valueobject.CardTypeReverse))
			g.Deck = append(g.Deck, NewActionCard(color, valueobject.CardTypeDrawTwo))
		}
	}
	for i := 0; i < 4; i++ {
		g.Deck = append(g.Deck, NewWildCard(valueobject.CardTypeWild))
		g.Deck = append(g.Deck, NewWildCard(valueobject.CardTypeWildDraw))
	}
}

func (g *Game) shuffleDeck() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(g.Deck), func(i, j int) {
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	})
}

func (g *Game) dealCards() {
	for i := 0; i < 7; i++ {
		for _, player := range g.Players {
			if card := g.drawCard(); card != nil {
				player.AddCard(card)
			}
		}
	}
}

func (g *Game) flipFirstCard() {
	for {
		card := g.drawCard()
		if card == nil {
			break
		}
		if card.Type.IsWildCard() {
			g.Deck = append(g.Deck, card)
			g.shuffleDeck()
			continue
		}
		g.DiscardPile = append(g.DiscardPile, card)
		g.CurrentColor = card.Color
		break
	}
}

func (g *Game) drawCard() *Card {
	if len(g.Deck) == 0 {
		g.reshuffleDiscardPile()
	}
	if len(g.Deck) == 0 {
		return nil
	}
	card := g.Deck[0]
	g.Deck = g.Deck[1:]
	return card
}

func (g *Game) DrawCards(count int) []*Card {
	cards := make([]*Card, 0, count)
	for i := 0; i < count; i++ {
		if card := g.drawCard(); card != nil {
			cards = append(cards, card)
		}
	}
	return cards
}

func (g *Game) reshuffleDiscardPile() {
	if len(g.DiscardPile) <= 1 {
		return
	}
	topCard := g.DiscardPile[len(g.DiscardPile)-1]
	g.Deck = g.DiscardPile[:len(g.DiscardPile)-1]
	g.DiscardPile = []*Card{topCard}
	g.shuffleDeck()
}

func (g *Game) GetCurrentPlayer() *Player {
	if len(g.Players) == 0 {
		return nil
	}
	return g.Players[g.CurrentPlayer]
}

func (g *Game) GetTopCard() *Card {
	if len(g.DiscardPile) == 0 {
		return nil
	}
	return g.DiscardPile[len(g.DiscardPile)-1]
}

func (g *Game) PlayCard(playerID string, cardIndex int, chosenColor valueobject.Color) error {
	if g.State != GameStatePlaying {
		return errors.New("游戏未开始")
	}
	currentPlayer := g.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.ID != playerID {
		return errors.New("不是你的回合")
	}
	card := currentPlayer.GetCard(cardIndex)
	if card == nil {
		return errors.New("无效的卡牌")
	}
	topCard := g.GetTopCard()
	if !card.CanPlayOn(topCard, g.CurrentColor) {
		return errors.New("不能打出这张牌")
	}

	// 先验证万能牌颜色选择，避免卡牌丢失
	if card.Type.IsWildCard() {
		if !chosenColor.IsValid() || chosenColor == valueobject.ColorWild {
			return errors.New("请选择颜色")
		}
	}

	// 检查最后一张牌不能是功能牌
	if len(currentPlayer.Hand) == 1 && card.Type.IsEndingForbidden() {
		return errors.New("最后一张牌不能出功能牌")
	}

	// 验证通过后才移除卡牌
	currentPlayer.RemoveCard(cardIndex)
	g.DiscardPile = append(g.DiscardPile, card)

	if card.Type.IsWildCard() {
		g.CurrentColor = chosenColor
	} else {
		g.CurrentColor = card.Color
	}

	// 检查是否需要触发 UNO 按钮
	if len(currentPlayer.Hand) == 1 {
		g.activateUnoButton(playerID)
	}

	if currentPlayer.HasWon() {
		g.State = GameStateFinished
		g.Winner = currentPlayer
		return nil
	}

	g.handleCardEffect(card, playerID)
	return nil
}

func (g *Game) handleCardEffect(card *Card, playerID string) {
	switch card.Type {
	case valueobject.CardTypeSkip:
		g.nextTurn()
		g.nextTurn()
	case valueobject.CardTypeReverse:
		g.Direction *= -1
		if len(g.Players) == 2 {
			// 两人游戏中，Reverse 等同于 Skip
			g.nextTurn()
		} else {
			g.nextTurn()
		}
	case valueobject.CardTypeDrawTwo:
		g.nextTurn()
		if next := g.GetCurrentPlayer(); next != nil {
			next.AddCards(g.DrawCards(2))
		}
		g.nextTurn()
	case valueobject.CardTypeWildDraw:
		// +4 需要等待质疑
		g.nextTurn()
		if next := g.GetCurrentPlayer(); next != nil {
			g.PendingWildDraw = true
			g.WildDrawPlayer = playerID
			g.WildDrawVictim = next.ID
			g.State = GameStateWaitingChallenge
		}
	default:
		g.nextTurn()
	}
}

func (g *Game) nextTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + g.Direction + len(g.Players)) % len(g.Players)
	g.HasDrawnThisTurn = false
}

func (g *Game) DrawCardForPlayer(playerID string) (*Card, error) {
	if g.State != GameStatePlaying {
		return nil, errors.New("游戏未开始")
	}
	currentPlayer := g.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.ID != playerID {
		return nil, errors.New("不是你的回合")
	}
	if g.HasDrawnThisTurn {
		return nil, errors.New("本回合已经摸过牌了")
	}
	card := g.drawCard()
	if card == nil {
		return nil, errors.New("牌组已空")
	}
	currentPlayer.AddCard(card)
	g.HasDrawnThisTurn = true
	return card, nil
}

// MustDrawCard 强制摸牌（没有能打的牌时调用）
func (g *Game) MustDrawCard(playerID string) (*Card, bool, error) {
	if g.State != GameStatePlaying {
		return nil, false, errors.New("游戏未开始")
	}
	currentPlayer := g.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.ID != playerID {
		return nil, false, errors.New("不是你的回合")
	}
	
	// 检查是否有能打的牌
	if g.HasPlayableCard(playerID) {
		return nil, false, errors.New("你有能打的牌，不能摸牌")
	}
	
	card := g.drawCard()
	if card == nil {
		return nil, false, errors.New("牌组已空")
	}
	currentPlayer.AddCard(card)
	g.HasDrawnThisTurn = true
	
	// 检查摸到的牌是否能打
	canPlay := card.CanPlayOn(g.GetTopCard(), g.CurrentColor)
	return card, canPlay, nil
}

// HasPlayableCard 检查玩家是否有能打的牌
func (g *Game) HasPlayableCard(playerID string) bool {
	player := g.GetPlayer(playerID)
	if player == nil {
		return false
	}
	topCard := g.GetTopCard()
	for _, card := range player.Hand {
		if card.CanPlayOn(topCard, g.CurrentColor) {
			return true
		}
	}
	return false
}

func (g *Game) PassTurn(playerID string) error {
	if g.State != GameStatePlaying {
		return errors.New("游戏未开始")
	}
	currentPlayer := g.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.ID != playerID {
		return errors.New("不是你的回合")
	}
	// 必须先摸牌才能跳过
	if !g.HasDrawnThisTurn {
		return errors.New("没有能打的牌时必须先摸一张牌")
	}
	g.nextTurn()
	g.HasDrawnThisTurn = false
	return nil
}

func (g *Game) GetPlayer(playerID string) *Player {
	for _, p := range g.Players {
		if p.ID == playerID {
			return p
		}
	}
	return nil
}

// ========== +4 质疑机制 ==========

// ChallengeWildDraw 质疑 +4
// 返回: 质疑是否成功, 罚牌的玩家ID, 罚牌数量, 错误
func (g *Game) ChallengeWildDraw(challengerID string) (bool, string, int, error) {
	if g.State != GameStateWaitingChallenge {
		return false, "", 0, errors.New("当前没有可质疑的+4")
	}
	if challengerID != g.WildDrawVictim {
		return false, "", 0, errors.New("只有被+4的玩家可以质疑")
	}

	// 检查打出+4的玩家当时是否有其他能打的牌
	wildDrawPlayer := g.GetPlayer(g.WildDrawPlayer)
	if wildDrawPlayer == nil {
		return false, "", 0, errors.New("找不到打出+4的玩家")
	}

	// 获取+4之前的牌堆顶牌（现在是倒数第二张）
	var previousTopCard *Card
	if len(g.DiscardPile) >= 2 {
		previousTopCard = g.DiscardPile[len(g.DiscardPile)-2]
	}

	// 检查+4玩家是否有其他能打的牌（不包括万能牌）
	hadPlayableCard := false
	for _, card := range wildDrawPlayer.Hand {
		if !card.Type.IsWildCard() && card.CanPlayOn(previousTopCard, g.CurrentColor) {
			hadPlayableCard = true
			break
		}
	}

	g.State = GameStatePlaying
	g.PendingWildDraw = false

	if hadPlayableCard {
		// 质疑成功：+4玩家罚4张牌，质疑者不用罚牌
		wildDrawPlayer.AddCards(g.DrawCards(4))
		g.WildDrawPlayer = ""
		g.WildDrawVictim = ""
		g.nextTurn() // 轮到下一个玩家
		return true, wildDrawPlayer.ID, 4, nil
	} else {
		// 质疑失败：质疑者罚6张牌（原来的4张+额外2张）
		challenger := g.GetPlayer(challengerID)
		if challenger != nil {
			challenger.AddCards(g.DrawCards(6))
		}
		g.WildDrawPlayer = ""
		g.WildDrawVictim = ""
		g.nextTurn() // 质疑者跳过回合
		return false, challengerID, 6, nil
	}
}

// AcceptWildDraw 接受 +4 不质疑
func (g *Game) AcceptWildDraw(playerID string) error {
	if g.State != GameStateWaitingChallenge {
		return errors.New("当前没有可接受的+4")
	}
	if playerID != g.WildDrawVictim {
		return errors.New("只有被+4的玩家可以接受")
	}

	victim := g.GetPlayer(playerID)
	if victim != nil {
		victim.AddCards(g.DrawCards(4))
	}

	g.State = GameStatePlaying
	g.PendingWildDraw = false
	g.WildDrawPlayer = ""
	g.WildDrawVictim = ""
	g.nextTurn() // 被+4的玩家跳过回合
	return nil
}

// ========== UNO 喊话机制 ==========

// activateUnoButton 激活 UNO 按钮
func (g *Game) activateUnoButton(playerID string) {
	g.UnoButtonActive = true
	g.UnoPlayerID = playerID
	g.UnoButtonPressedBy = ""
	g.UnoButtonTime = time.Now()
}

// PressUnoButton 按下 UNO 按钮
// 返回: 是否成功, 被罚牌的玩家ID（如果有）, 错误
func (g *Game) PressUnoButton(playerID string) (bool, string, error) {
	if !g.UnoButtonActive {
		return false, "", errors.New("UNO按钮未激活")
	}
	if g.UnoButtonPressedBy != "" {
		return false, "", errors.New("UNO按钮已被按下")
	}

	g.UnoButtonPressedBy = playerID
	g.UnoButtonActive = false

	if playerID == g.UnoPlayerID {
		// 玩家自己按下，成功喊 UNO
		return true, "", nil
	} else {
		// 其他玩家先按下，UNO 玩家罚2张牌
		unoPlayer := g.GetPlayer(g.UnoPlayerID)
		if unoPlayer != nil {
			unoPlayer.AddCards(g.DrawCards(2))
		}
		return true, g.UnoPlayerID, nil
	}
}

// IsUnoButtonActive 检查 UNO 按钮是否激活
func (g *Game) IsUnoButtonActive() bool {
	return g.UnoButtonActive
}

// GetUnoButtonInfo 获取 UNO 按钮信息
func (g *Game) GetUnoButtonInfo() (active bool, playerID string, elapsed time.Duration) {
	return g.UnoButtonActive, g.UnoPlayerID, time.Since(g.UnoButtonTime)
}

// CancelUnoButton 取消 UNO 按钮（超时时调用）
func (g *Game) CancelUnoButton() {
	g.UnoButtonActive = false
	g.UnoPlayerID = ""
	g.UnoButtonPressedBy = ""
}
