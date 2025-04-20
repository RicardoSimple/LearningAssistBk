package hashKey

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/robfig/cron/v3"

	"ar-app-api/util/log"
)

const (
	Path        = "storage"
	SurviveTime = time.Minute * 30
)

var repo Repo

type Repo struct {
	badgers      map[string]*Database
	databaseLock *sync.RWMutex
	cron         *cron.Cron
}

type Database struct {
	db    *badger.DB
	timer *time.Timer
}

func init() {
	repo = Repo{
		badgers:      make(map[string]*Database),
		databaseLock: new(sync.RWMutex),
		cron:         cron.New(cron.WithSeconds()),
	}
	repo.cron.AddFunc("0 4 * * * *", repo.RunValueLogGC)
	repo.cron.Start()
}
func GetDB(name string) (*badger.DB, error) {
	if name == "" {
		return nil, fmt.Errorf("database name is empty")
	}
	repo.databaseLock.RLock()
	if database := repo.getDB(name); database != nil {
		repo.databaseLock.RUnlock()
		return database.db, nil
	}
	repo.databaseLock.RUnlock()
	return repo.open(name)
}

func (repo *Repo) getDB(name string) *Database {
	db, ok := repo.badgers[name]
	if !ok {
		return nil
	}
	db.timer.Reset(SurviveTime)
	return db
}
func (repo *Repo) RunValueLogGC() {
	for _, database := range repo.badgers {
		for {
			if err := database.db.RunValueLogGC(0.5); err != nil {
				break
			}
		}
	}
	return
}
func (repo *Repo) open(name string) (*badger.DB, error) {
	repo.databaseLock.Lock()
	defer repo.databaseLock.Unlock()
	// 再次获取数据库
	if database := repo.getDB(name); database != nil {

		return database.db, nil
	}
	dbPath := path.Join(Path, name)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		err = nil
		if err = os.Mkdir(dbPath, os.ModePerm); err != nil {
			log.Error(err.Error())
			return nil, err
		}
	}
	config := badger.DefaultOptions(dbPath).
		WithNumMemtables(3).WithNumLevelZeroTables(3).WithNumLevelZeroTablesStall(5).
		WithNumCompactors(2).WithBaseTableSize(1 << 20).WithValueLogFileSize(5 << 20).
		WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(config)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	// 定时关闭
	timer := time.NewTimer(SurviveTime)
	go func() {
		<-timer.C
		repo.Close(name)
	}()
	repo.badgers[name] = &Database{
		db:    db,
		timer: timer,
	}
	return db, nil
}

func (repo *Repo) Close(name string) {
	repo.databaseLock.Lock()
	defer repo.databaseLock.Unlock()
	database, ok := repo.badgers[name]
	if !ok {
		return
	}
	if err := database.db.Close(); err != nil {
		log.Error("hashKey close err:%v", err)
		return
	}
	database.timer.Stop()
	delete(repo.badgers, name)
	return
}
