package hash

import (
	"context"
	"errors"
	"io"
	"learning-assistant/dal"
	"mime/multipart"
	"strconv"

	"learning-assistant/dal/hashKey"
	"learning-assistant/util"
	"learning-assistant/util/decode"
	"learning-assistant/util/log"
)

func BindImageHash(ctx context.Context, header *multipart.FileHeader, repo string) error {
	hash, err := decode.GetHash(header)
	if err != nil {
		return err
	}

	file, err := header.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileStr, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Upload the image todo 422 Unprocessable Entity
	uploadFile, err := util.GetUploader().UploadFile(header.Filename+"2", string(fileStr))
	if err != nil {
		return err
	}
	log.Info(uploadFile)

	// Store hashImage
	id, err := dal.CreateImageHash(ctx, uploadFile, hash.ToString(), "", 3)
	if err != nil {
		return err
	}

	// Store hash
	db, err := hashKey.GetDB(repo)
	if err != nil {
		return err
	}
	store := hashKey.NewBadgerStore(db)
	err = store.Set(strconv.Itoa(int(id)), &hashKey.HashValue{
		ImageId: id,
		Hash:    *hash,
	})
	return err
}

func QuerySimilarImage(ctx context.Context, header *multipart.FileHeader, repo string) (uint, error) {
	// 获取hash
	hash, err := decode.GetHash(header)
	if err != nil {
		return 0, err
	}
	// 遍历
	db, err := hashKey.GetDB(repo)
	if err != nil {
		return 0, err
	}
	store := hashKey.NewBadgerStore(db)
	hashValues, err := store.FindAll()
	similar := decode.TopKSimilar(hash, hashValues, 1)
	if len(similar) == 0 {
		return 0, errors.New("无数据")
	}
	if similar[0].Distance < 10 {
		log.Info("相似图片距离为：%v", similar[0].Distance)
		return similar[0].Hash.ImageId, nil
	}
	return 0, errors.New("无相似图片")
}
