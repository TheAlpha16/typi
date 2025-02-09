package database

import (
	"time"

	"github.com/TheAlpha16/typi/cron/config"
	"github.com/TheAlpha16/typi/cron/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func UpdateLastFetch() error {
	if err := DB.
		Model(&models.Config{}).
		Where("key = ?", "last_fetch").
		Update("value", config.LAST_FETCH.Format(time.RFC3339)).Error; err != nil {
		logrus.WithError(err).Error("failed to update last fetch")
		return err
	}

	return nil
}

func GetLastFetch() (time.Time, error) {
	last_fetch, err := GetConfig("last_fetch")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			_ = SetConfig("last_fetch", config.LAST_FETCH.Format(time.RFC3339))
			return config.LAST_FETCH, nil
		}
		logrus.WithError(err).Error("unable to get last fetch")
		return config.LAST_FETCH, err
	}

	t, err := time.Parse(time.RFC3339, last_fetch)
	if err != nil {
		logrus.WithError(err).Error("error parsing last fetch")
		return config.LAST_FETCH, err
	}

	return t, nil
}
