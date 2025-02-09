package database

import (
	"github.com/TheAlpha16/typi/cron/models"
	"gorm.io/gorm/clause"
)

/*
UploadVideos uploads a batch of videos to the database. If a video with the same
'vid' already exists, it updates the 'title', 'description', 'published_at', and
'thumbnail' fields of the existing record.

Parameters:
	- videos: A pointer to a slice of Video models to be uploaded.

Returns:
	- error: An error if the operation fails, otherwise nil.
*/

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
