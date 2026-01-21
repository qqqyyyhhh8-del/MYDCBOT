package uno

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/user/dcminigames/internal/domain/uno/entity"
	"github.com/user/dcminigames/internal/domain/uno/valueobject"
	"github.com/user/dcminigames/internal/infrastructure/imaging"
	"github.com/user/dcminigames/internal/infrastructure/persistence/memory"
)

type Handler struct {
	repo     *memory.GameRepository
	renderer *imaging.CardRenderer
}

func NewHandler(repo *memory.GameRepository, renderer *imaging.CardRenderer) *Handler {
	return &Handler{repo: repo, renderer: renderer}
}

func (h *Handler) CreateGame(channelID string) (*entity.Game, error) {
	if h.repo.Exists(channelID) {
		return nil, fmt.Errorf("该频道已有游戏进行中")
	}
	game := entity.NewGame(uuid.New().String(), channelID)
	if err := h.repo.Save(game); err != nil {
		return nil, err
	}
	return game, nil
}

func (h *Handler) JoinGame(channelID, playerID, username string) error {
	game, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return fmt.Errorf("没有进行中的游戏，请先使用 /uno create")
	}
	player := entity.NewPlayer(playerID, username)
	if err := game.AddPlayer(player); err != nil {
		return err
	}
	return h.repo.Save(game)
}

func (h *Handler) StartGame(channelID, playerID string) error {
	game, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return err
	}
	if len(game.Players) == 0 || game.Players[0].ID != playerID {
		return fmt.Errorf("只有房主可以开始游戏")
	}
	if err := game.Start(); err != nil {
		return err
	}
	return h.repo.Save(game)
}

func (h *Handler) GetGame(channelID string) (*entity.Game, error) {
	return h.repo.FindByChannelID(channelID)
}

func (h *Handler) GetPlayerHand(channelID, playerID string) ([]*entity.Card, error) {
	game, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return nil, err
	}
	player := game.GetPlayer(playerID)
	if player == nil {
		return nil, fmt.Errorf("你不在游戏中")
	}
	return player.Hand, nil
}

func (h *Handler) RenderPlayerHand(channelID, playerID string) ([]byte, error) {
	cards, err := h.GetPlayerHand(channelID, playerID)
	if err != nil {
		return nil, err
	}
	if len(cards) == 0 {
		return nil, fmt.Errorf("手牌为空")
	}
	return h.renderer.RenderHand(cards)
}

// RenderSingleCard 渲染单张卡牌
func (h *Handler) RenderSingleCard(card *entity.Card) ([]byte, error) {
	return h.renderer.RenderSingleCard(card)
}

// PlayCardAndGetCard 打出卡牌并返回打出的卡牌（用于显示图片）
func (h *Handler) PlayCardAndGetCard(channelID, playerID string, cardIndex int, chosenColor valueobject.Color) (*entity.Card, error) {
	game, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return nil, err
	}
	player := game.GetPlayer(playerID)
	if player == nil {
		return nil, fmt.Errorf("你不在游戏中")
	}
	if cardIndex < 0 || cardIndex >= len(player.Hand) {
		return nil, fmt.Errorf("无效的卡牌索引")
	}
	playedCard := player.Hand[cardIndex]
	if err := game.PlayCard(playerID, cardIndex, chosenColor); err != nil {
		return nil, err
	}
	if err := h.repo.Save(game); err != nil {
		return nil, err
	}
	return playedCard, nil
}

func (h *Handler) PlayCard(channelID, playerID string, cardIndex int, chosenColor valueobject.Color) error {
	game, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return err
	}
	if err := game.PlayCard(playerID, cardIndex, chosenColor); err != nil {
		return err
	}
	return h.repo.Save(game)
}

func (h *Handler) DrawCard(channelID, playerID string) (*entity.Card, error) {
	game, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return nil, err
	}
	card, err := game.DrawCardForPlayer(playerID)
	if err != nil {
		return nil, err
	}
	if err := h.repo.Save(game); err != nil {
		return nil, err
	}
	return card, nil
}

func (h *Handler) PassTurn(channelID, playerID string) error {
	game, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return err
	}
	if err := game.PassTurn(playerID); err != nil {
		return err
	}
	return h.repo.Save(game)
}

func (h *Handler) EndGame(channelID string) error {
	return h.repo.Delete(channelID)
}