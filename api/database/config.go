package database

import (
	"github.com/TheAlpha16/typi/api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

func GetConfig(key string) (string, error) {
	var config models.Config

	if err := DB.
		Model(&models.Config{}).
		Where("key = ?", key).
		First(&config).Error; err != nil {
		logrus.WithError(err).Error("unable to read config")
		return "", err
	}

	return config.Value, nil
}

func SetConfig(key, value string) error {
	config := models.Config{
		Key:   key,
		Value: value,
	}

	if err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).
		Create(&config).Error; err != nil {
		logrus.WithError(err).Error("unable to set config")
		return err
	}

	return nil
}
