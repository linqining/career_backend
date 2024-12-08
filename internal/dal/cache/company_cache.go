package cache

//go:generate mockgen -source=internal/dal/cache/company_cache.go -destination=internal/mock/company_cache_mock.go  -package mock

import (
	model2 "career_backend/internal/dal/db/model"
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/utils"
	"github.com/redis/go-redis/v9"
)

var (
	// PrefixCompanyCacheKey cache prefix
	PrefixCompanyCacheKey = utils.ConcatString(prefix, "company:%d")
)

// CompanyCache define cache interface
type CompanyCache interface {
	SetCompanyCache(ctx context.Context, id int64, data *model2.Company, duration time.Duration) error
	GetCompanyCache(ctx context.Context, id int64) (data *model2.Company, err error)
	MultiGetCompanyCache(ctx context.Context, ids []int64) (map[string]*model2.Company, error)
	MultiSetCompanyCache(ctx context.Context, data []*model2.Company, duration time.Duration) error
	DelCompanyCache(ctx context.Context, id int64) error
	SetCacheWithNotFound(ctx context.Context, id int64) error
}

// companyCache define cache struct
type companyCache struct {
	cache cache.Cache
}

// NewCompanyCache new a cache
func NewCompanyCache(rdb *redis.Client) CompanyCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &companyCache{
		cache: cache.NewRedisCache(rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model2.Company{}
		}),
	}
}

// GetCompanyCacheKey get cache key
func (c *companyCache) GetCompanyCacheKey(id int64) string {
	return fmt.Sprintf(PrefixCompanyCacheKey, id)
}

// SetCompanyCache write to cache
func (c *companyCache) SetCompanyCache(ctx context.Context, id int64, data *model2.Company, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetCompanyCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetCompanyCache get from cache
func (c *companyCache) GetCompanyCache(ctx context.Context, id int64) (data *model2.Company, err error) {
	cacheKey := c.GetCompanyCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("[cache] GetCompanyCache err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetCompanyCache batch get cache
func (c *companyCache) MultiGetCompanyCache(ctx context.Context, ids []int64) (map[string]*model2.Company, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetCompanyCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	retMap := make(map[string]*model2.Company)
	err := c.cache.MultiGet(ctx, keys, retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

// MultiSetCompanyCache batch set cache
func (c *companyCache) MultiSetCompanyCache(ctx context.Context, data []*model2.Company, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetCompanyCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}
	return nil
}

// DelCompanyCache delete cache
func (c *companyCache) DelCompanyCache(ctx context.Context, id int64) error {
	cacheKey := c.GetCompanyCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *companyCache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.GetCompanyCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
