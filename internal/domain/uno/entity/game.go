package entity

import (
	"errors"
	"math/rand"
	"time"

	"github.com/user/dcminigames/internal/domain/uno/valueobject"
)

type GameState string

const (
	GameStateWaiting  GameState = "waiting"
	GameStatePlaying  GameState = "playing"
	GameStateFinished GameState = "finished"
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
}

func NewGame(id, channelID string) *Game {
	return &Game{
		ID:            id,
		ChannelID:     channelID,
		Players:       make([]*Player, 0),
		Deck:          make([]*Card, 0),
		DiscardPile:   make([]*Card, 0),
		CurrentPlayer: 0,
		Direction:     1,
		State:         GameStateWaiting,
		CreatedAt:     time.Now(),
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
	currentPlayer.RemoveCard(cardIndex)
	g.DiscardPile = append(g.DiscardPile, card)

	if card.Type.IsWildCard() {
		if !chosenColor.IsValid() || chosenColor == valueobject.ColorWild {
			return errors.New("请选择颜色")
		}
		g.CurrentColor = chosenColor
	} else {
		g.CurrentColor = card.Color
	}

	if currentPlayer.HasWon() {
		g.State = GameStateFinished
		g.Winner = currentPlayer
		return nil
	}
	g.handleCardEffect(card)
	return nil
}

func (g *Game) handleCardEffect(card *Card) {
	switch card.Type {
	case valueobject.CardTypeSkip:
		g.nextTurn()
		g.nextTurn()
	case valueobject.CardTypeReverse:
		g.Direction *= -1
		if len(g.Players) == 2 {
			g.nextTurn()
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
		g.nextTurn()
		if next := g.GetCurrentPlayer(); next != nil {
			next.AddCards(g.DrawCards(4))
		}
		g.nextTurn()
	default:
		g.nextTurn()
	}
}

func (g *Game) nextTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + g.Direction + len(g.Players)) % len(g.Players)
}

func (g *Game) DrawCardForPlayer(playerID string) (*Card, error) {
	if g.State != GameStatePlaying {
		return nil, errors.New("游戏未开始")
	}
	currentPlayer := g.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.ID != playerID {
		return nil, errors.New("不是你的回合")
	}
	card := g.drawCard()
	if card == nil {
		return nil, errors.New("牌组已空")
	}
	currentPlayer.AddCard(card)
	return card, nil
}

func (g *Game) PassTurn(playerID string) error {
	if g.State != GameStatePlaying {
		return errors.New("游戏未开始")
	}
	currentPlayer := g.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.ID != playerID {
		return errors.New("不是你的回合")
	}
	g.nextTurn()
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
