package service

import (
	"github.com/xyu/mdxy/internal/repository"
	"gorm.io/gorm"
)

type GameScoreService struct {
	repo *repository.GameScoreRepository
}

func NewGameScoreService(db *gorm.DB) *GameScoreService {
	return &GameScoreService{
		repo: repository.NewGameScoreRepository(db),
	}
}

// UpdateScore 更新游戏分数
func (s *GameScoreService) UpdateScore(toolID uint, playerID string, score int) (int, error) {
	record, err := s.repo.UpdateScore(toolID, playerID, score)
	if err != nil {
		return 0, err
	}
	return record.BestScore, nil
}

// GetBestScore 获取最高分
func (s *GameScoreService) GetBestScore(toolID uint, playerID string) (int, error) {
	return s.repo.GetBestScore(toolID, playerID)
}

// GetLeaderboard 获取排行榜
func (s *GameScoreService) GetLeaderboard(toolID uint, limit int) (interface{}, error) {
	scores, err := s.repo.GetLeaderboard(toolID, limit)
	if err != nil {
		return nil, err
	}

	// 格式化返回数据，隐藏完整 player_id
	result := make([]map[string]interface{}, len(scores))
	for i, score := range scores {
		pid := score.PlayerID
		if len(pid) > 8 {
			pid = pid[:8]
		}
		playerName := "玩家" + pid
		result[i] = map[string]interface{}{
			"rank":       i + 1,
			"player":     playerName,
			"best_score": score.BestScore,
		}
	}
	return result, nil
}
