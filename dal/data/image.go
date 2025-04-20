package data

import (
	"context"
	"log"

	"gorm.io/gorm"

	"ar-app-api/dal"
	"ar-app-api/dal/schema"
)

// CreateImageHash 创建一个新的 schema.ImageHash 记录
func CreateImageHash(ctx context.Context, url, hash, desc string, mediaID uint) (uint, error) {
	imageHash := schema.ImageHash{
		URL:     url,
		Hash:    hash,
		Desc:    desc,
		MediaId: mediaID,
	}
	err := dal.DB.Create(&imageHash).Error
	if err != nil {
		log.Println("Error creating schema.ImageHash:", err)
	}
	return imageHash.ID, err
}

// ReadImageHashByID 根据 ID 查询 schema.ImageHash 记录
func ReadImageHashByID(ctx context.Context, id uint) (*schema.ImageHash, error) {
	var imageHash schema.ImageHash
	err := dal.DB.First(&imageHash, id).Error
	if err != nil {
		log.Println("Error reading schema.ImageHash by ID:", err)
	}
	return &imageHash, err
}

// UpdateImageHash 更新 schema.ImageHash 记录
func UpdateImageHash(ctx context.Context, id uint, url, hash, desc string, mediaID uint) error {
	imageHash := schema.ImageHash{
		Model:   gorm.Model{ID: id},
		URL:     url,
		Hash:    hash,
		Desc:    desc,
		MediaId: mediaID,
	}
	err := dal.DB.Save(&imageHash).Error
	if err != nil {
		log.Println("Error updating schema.ImageHash:", err)
	}
	return err
}

// DeleteImageHashByID 根据 ID 删除 schema.ImageHash 记录
func DeleteImageHashByID(ctx context.Context, id uint) error {
	err := dal.DB.Delete(&schema.ImageHash{}, id).Error
	if err != nil {
		log.Println("Error deleting schema.ImageHash by ID:", err)
	}
	return err
}
