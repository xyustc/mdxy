package repository

import (
	"github.com/xyu/mdxy/internal/model"
	"gorm.io/gorm"
)

type GameScoreRepository struct {
	db *gorm.DB
}

func NewGameScoreRepository(db *gorm.DB) *GameScoreRepository {
	return &GameScoreRepository{db: db}
}

// GetOrCreate 获取或创建玩家游戏记录
func (r *GameScoreRepository) GetOrCreate(toolID uint, playerID string) (*model.GameScore, error) {
	var score model.GameScore
	err := r.db.Where("tool_id = ? AND player_id = ?", toolID, playerID).First(&score).Error
	if err != nil {
		// 不存在则创建
		score = model.GameScore{
			ToolID:    toolID,
			PlayerID:  playerID,
			Score:     0,
			BestScore: 0,
		}
		if err := r.db.Create(&score).Error; err != nil {
			return nil, err
		}
	}
	return &score, nil
}

// UpdateScore 更新分数
func (r *GameScoreRepository) UpdateScore(toolID uint, playerID string, score int) (*model.GameScore, error) {
	record, err := r.GetOrCreate(toolID, playerID)
	if err != nil {
		return nil, err
	}

	record.Score = score
	if score > record.BestScore {
		record.BestScore = score
	}

	if err := r.db.Save(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// GetBestScore 获取最高分
func (r *GameScoreRepository) GetBestScore(toolID uint, playerID string) (int, error) {
	record, err := r.GetOrCreate(toolID, playerID)
	if err != nil {
		return 0, err
	}
	return record.BestScore, nil
}

// GetLeaderboard 获取排行榜（前N名）
func (r *GameScoreRepository) GetLeaderboard(toolID uint, limit int) ([]model.GameScore, error) {
	var scores []model.GameScore
	err := r.db.Where("tool_id = ?", toolID).
		Order("best_score DESC").
		Limit(limit).
		Find(&scores).Error
	return scores, err
}
