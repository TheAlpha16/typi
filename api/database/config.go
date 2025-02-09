package database

import (
	"log"

	"github.com/TheAlpha16/typi/api/models"
	"gorm.io/gorm/clause"
)

func GetConfig(key string) (string, error) {
	var config models.Config

	if err := DB.
		Model(&models.Config{}).
		Where("key = ?", key).
		First(&config).Error; err != nil {
		log.Println(err)
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
		log.Println(err)
		return err
	}

	return nil
}
