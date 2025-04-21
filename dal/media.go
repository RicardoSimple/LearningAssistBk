package dal

import (
	"context"
	"log"

	"gorm.io/gorm"

	"learning-assistant/dal/schema"
)

// CreateMedia 创建一个新的 Media 记录
func CreateMedia(ctx context.Context, content, mediaType, url string) error {
	media := schema.Media{
		Content: content,
		Type:    mediaType,
		Url:     url,
	}
	err := DB.Create(&media).Error
	if err != nil {
		log.Println("Error creating Media:", err)
	}
	return err
}

// ReadMediaByID 根据 ID 查询 Media 记录
func ReadMediaByID(ctx context.Context, id uint) (*schema.Media, error) {
	var media schema.Media
	err := DB.First(&media, id).Error
	if err != nil {
		log.Println("Error reading Media by ID:", err)
	}
	return &media, err
}

// UpdateMedia 更新 Media 记录
func UpdateMedia(ctx context.Context, id uint, content, mediaType, url string) error {
	media := schema.Media{
		Model:   gorm.Model{ID: id},
		Content: content,
		Type:    mediaType,
		Url:     url,
	}
	err := DB.Save(&media).Error
	if err != nil {
		log.Println("Error updating Media:", err)
	}
	return err
}

// DeleteMediaByID 根据 ID 删除 Media 记录
func DeleteMediaByID(ctx context.Context, id uint) error {
	err := DB.Delete(&schema.Media{}, id).Error
	if err != nil {
		log.Println("Error deleting Media by ID:", err)
	}
	return err
}
