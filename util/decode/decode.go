package decode

import (
	"bufio"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"path"
	"sort"
	"strings"

	"github.com/corona10/goimagehash"

	"learning-assistant/dal/hashKey"
)

func MFileToImage(header *multipart.FileHeader) (image image.Image, err error) {
	file, err := header.Open()
	ext := strings.ToLower(path.Ext(header.Filename))
	switch ext {
	case ".jpeg", ".jpg":
		image, err = jpeg.Decode(bufio.NewReader(file))
	case ".png":
		image, err = png.Decode(bufio.NewReader(file))
	}
	return image, err
}

func GetHash(header *multipart.FileHeader) (*goimagehash.ImageHash, error) {
	// 获取文件类型
	// 转换为图片
	img, err := MFileToImage(header)
	if err != nil {
		return nil, err
	}
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// HashDistance 包含哈希值和距离信息
type HashDistance struct {
	Hash     *hashKey.HashValue
	Distance int
}

// TopKSimilar 返回距离最小的 K 个哈希值
func TopKSimilar(h *goimagehash.ImageHash, values []*hashKey.HashValue, k int) []*HashDistance {
	// 定义存储哈希距离信息的数组
	distances := make([]*HashDistance, 0, len(values))

	// 计算每个哈希值与目标哈希值的距离，并保存到数组中
	for _, value := range values {
		dist, _ := h.Distance(&value.Hash)
		distances = append(distances, &HashDistance{
			Hash:     value,
			Distance: dist,
		})
	}

	// 根据距离升序排序
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].Distance < distances[j].Distance
	})

	// 只保留前 K 个距离最小的哈希值
	if len(distances) > k {
		distances = distances[:k]
	}

	// 返回距离最小的 K 个哈希值
	result := make([]*HashDistance, len(distances))
	copy(result, distances)
	return result
}
