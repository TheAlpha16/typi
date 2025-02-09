package database

import (
	"github.com/TheAlpha16/typi/cron/models"
	"gorm.io/gorm/clause"
)

func UploadVideos(videos *[]models.Video) error {
	if err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "vid"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "description", "published_at", "thumbnail"}),
	}).
		CreateInBatches(videos, 100).Error; err != nil {
		return err
	}

	return nil
}
