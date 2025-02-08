package database

import (
	"log"

	"github.com/TheAlpha16/typi/api/models"
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

func GetVideoCount() (int64, error) {
	var count int64

	if err := DB.Model(&models.Video{}).Count(&count).Error; err != nil {
		log.Println(err)
		return 0, err
	}

	return count, nil
}

func GetVideos(offset, limit int) ([]models.Video, error) {
	var videos []models.Video

	if err := DB.
		Model(&models.Video{}).
		Offset(offset).
		Limit(limit).
		Find(&videos).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	return videos, nil
}
