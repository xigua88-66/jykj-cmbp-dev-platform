package global

import (
	"github.com/qiniu/qmgo"
	"sync"

	"github.com/songzhibin97/gkit/cache/local_cache"
	"jykj-cmbp-dev-platform/server/utils/timer"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"jykj-cmbp-dev-platform/server/config"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	CMBP_DB     *gorm.DB
	CMBP_DBList map[string]*gorm.DB
	CMBP_REDIS  *redis.Client
	CMBP_MONGO  *qmgo.QmgoClient
	CMBP_CONFIG config.Server
	CMBP_VP     *viper.Viper
	// CMBP_LOG    *oplogging.Logger
	CMBP_LOG                 *zap.Logger
	CMBP_Timer               timer.Timer = timer.NewTimerTask()
	CMBP_Concurrency_Control             = &singleflight.Group{}

	BlackCache local_cache.Cache
	lock       sync.RWMutex
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return CMBP_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := CMBP_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}
