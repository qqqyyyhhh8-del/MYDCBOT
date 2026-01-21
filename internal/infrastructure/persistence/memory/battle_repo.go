package memory

import (
	"errors"
	"sync"

	"github.com/user/dcminigames/internal/domain/pokemon/entity"
)

// BattleRepository 对战内存仓储
type BattleRepository struct {
	battles map[string]*entity.Battle
	mu      sync.RWMutex
}

// NewBattleRepository 创建对战仓储
func NewBattleRepository() *BattleRepository {
	return &BattleRepository{battles: make(map[string]*entity.Battle)}
}

// Save 保存对战
func (r *BattleRepository) Save(battle *entity.Battle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.battles[battle.ChannelID] = battle
	return nil
}

// FindByChannelID 通过频道ID查找对战
func (r *BattleRepository) FindByChannelID(channelID string) (*entity.Battle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	battle, ok := r.battles[channelID]
	if !ok {
		return nil, errors.New("对战不存在")
	}
	return battle, nil
}

// Delete 删除对战
func (r *BattleRepository) Delete(channelID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.battles, channelID)
	return nil
}

// Exists 检查对战是否存在
func (r *BattleRepository) Exists(channelID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.battles[channelID]
	return ok
}

// FindByPlayerID 通过玩家ID查找对战
func (r *BattleRepository) FindByPlayerID(playerID string) (*entity.Battle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, battle := range r.battles {
		if battle.Player1 != nil && battle.Player1.ID == playerID {
			return battle, nil
		}
		if battle.Player2 != nil && battle.Player2.ID == playerID {
			return battle, nil
		}
	}
	return nil, errors.New("未找到对战")
}
