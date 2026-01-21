package memory

import (
	"errors"
	"sync"

	"github.com/user/dcminigames/internal/domain/uno/entity"
)

type GameRepository struct {
	games map[string]*entity.Game
	mu    sync.RWMutex
}

func NewGameRepository() *GameRepository {
	return &GameRepository{games: make(map[string]*entity.Game)}
}

func (r *GameRepository) Save(game *entity.Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.games[game.ChannelID] = game
	return nil
}

func (r *GameRepository) FindByChannelID(channelID string) (*entity.Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	game, ok := r.games[channelID]
	if !ok {
		return nil, errors.New("游戏不存在")
	}
	return game, nil
}

func (r *GameRepository) Delete(channelID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.games, channelID)
	return nil
}

func (r *GameRepository) Exists(channelID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.games[channelID]
	return ok
}
