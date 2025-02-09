package database

import (
	"github.com/TheAlpha16/typi/api/models"

	"github.com/sirupsen/logrus"
)

func GetVideoCount() (int64, error) {
	var count int64

	if err := DB.Model(&models.Video{}).Count(&count).Error; err != nil {
		logrus.WithError(err).Error("unable to get video count")
		return 0, err
	}

	return count, nil
}

func GetVideos(offset, limit int) ([]models.Video, error) {
	var videos []models.Video

	if err := DB.
		Model(&models.Video{}).
		Order("published_at desc").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error; err != nil {
		logrus.WithError(err).Error("unable to get videos")
		return nil, err
	}

	return videos, nil
}
