package hashKey

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/corona10/goimagehash"
	"github.com/dgraph-io/badger/v4"
)

type HashValue struct {
	ImageId uint `json:"imageId"`
	Hash    goimagehash.ImageHash
}

type BadgerStore struct {
	db *badger.DB
}

func NewBadgerStore(db *badger.DB) *BadgerStore {
	return &BadgerStore{db: db}
}

func (store *BadgerStore) Get(imageId string) (*HashValue, error) {
	var result HashValue

	err := store.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(imageId))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &result)
		})
		return err
	})

	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (store *BadgerStore) Set(imageId string, value *HashValue) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = store.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(imageId), bytes)
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to set value in Badger: %v", err)
	}

	return nil
}

func (store *BadgerStore) FindAll() ([]*HashValue, error) {
	results := make([]*HashValue, 0)
	err := store.db.View(func(txn *badger.Txn) error {
		// 创建迭代器，参数设置为 nil 表示从数据库的起始位置开始遍历
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		// 迭代器移动到数据库的起始位置
		for it.Rewind(); it.Valid(); it.Next() {
			// 获取当前迭代器指向的键值对
			item := it.Item()
			// 获取键
			key := item.Key()
			// 获取值
			var value []byte
			err := item.Value(func(val []byte) error {
				value = append([]byte{}, val...)
				return nil
			})
			if err != nil {
				return err
			}
			v := &HashValue{}
			err = json.Unmarshal(value, v)
			if err != nil {
				return err
			}
			results = append(results, v)
			// 处理键值对
			fmt.Printf("Key: %s, Value: %s\n", key, value)
		}
		// 返回 nil 表示遍历完成
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return results, nil
}
