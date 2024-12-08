package repository

//go:generate mockgen -source=company_repo.go -destination=../../internal/mocks/company_repo_mock.go  -package mocks

import (
	"context"
	"time"

	localCache "github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"career_backend/internal/dal"
	"career_backend/internal/dal/cache"
	"career_backend/internal/dal/db/dao"
	"career_backend/internal/dal/db/model"
)

var _ CompanyRepo = (*companyRepo)(nil)

// CompanyRepo define a repo interface
type CompanyRepo interface {
	CreateCompany(ctx context.Context, data *model.Company) (id int64, err error)
	UpdateCompany(ctx context.Context, id int64, data *model.Company) error
	GetCompany(ctx context.Context, id int64) (ret *model.Company, err error)
	BatchGetCompany(ctx context.Context, ids []int64) (ret []*model.Company, err error)
	GetCompanyByAddress(ctx context.Context, address string) (ret *model.Company, err error)
}

type companyRepo struct {
	db         *dal.DBClient
	tracer     trace.Tracer
	cache      cache.CompanyCache
	localCache localCache.Cache
	sg         singleflight.Group
}

// NewCompanyRepo new a repository and return
func NewCompanyRepo(db *dal.DBClient, cache cache.CompanyCache) CompanyRepo {
	return &companyRepo{
		db:         db,
		tracer:     otel.Tracer("company"),
		cache:      cache,
		localCache: localCache.NewMemoryCache("local:company:", encoding.JSONEncoding{}),
		sg:         singleflight.Group{},
	}
}

// CreateCompany create a item
func (r *companyRepo) CreateCompany(ctx context.Context, data *model.Company) (id int64, err error) {
	if data == nil {
		return 0, errors.New("[repo] CreateCompany data cannot be nil")
	}

	err = dao.Company.WithContext(ctx).Create(data)
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create Company err")
	}

	return data.ID, nil
}

// UpdateCompany update item
func (r *companyRepo) UpdateCompany(ctx context.Context, id int64, data *model.Company) error {
	if id == 0 {
		return errors.New("[repo] UpdateCompany id cannot be equal to 0")
	}
	if data == nil {
		return errors.New("[repo] UpdateCompany data cannot be nil")
	}

	_, err := dao.Company.WithContext(ctx).Where(dao.Company.ID.Eq(id)).Updates(data)
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelCompanyCache(ctx, id)
	return nil
}

// GetCompany get a record
func (r *companyRepo) GetCompany(ctx context.Context, id int64) (ret *model.Company, err error) {
	if id == 0 {
		return nil, errors.New("[repo] GetCompany id cannot be equal to 0")
	}
	// get data from local cache
	err = r.localCache.Get(ctx, cast.ToString(id), &ret)
	if err != nil {
		return nil, err
	}
	if ret != nil && ret.ID > 0 {
		return ret, nil
	}

	// read redis cache
	ret, err = r.cache.GetCompanyCache(ctx, id)
	if err != nil {
		return nil, err
	}
	if ret != nil && ret.ID > 0 {
		return ret, nil
	}

	// get data from db
	// 避免缓存击穿(瞬间有大量请求过来)
	val, err, _ := r.sg.Do("sg:company:"+cast.ToString(id), func() (interface{}, error) {
		// read db or rpc
		data, err := dao.Company.WithContext(ctx).Where(dao.Company.ID.Eq(id)).Last()
		if err != nil {
			// cache not found and set empty cache to avoid 缓存穿透
			// Note: 如果缓存空数据太多，会大大降低缓存命中率，可以改为使用布隆过滤器
			if errors.Is(err, gorm.ErrRecordNotFound) {
				r.cache.SetCacheWithNotFound(ctx, id)
			}
			return nil, errors.Wrapf(err, "[repo] get Company from db error, id: %d", id)
		}

		// write cache
		if data != nil && data.ID > 0 {
			// write redis
			err = r.cache.SetCompanyCache(ctx, id, data, 5*time.Minute)
			if err != nil {
				return nil, errors.WithMessage(err, "[repo] GetCompany SetCompanyCache error")
			}

			// write local cache
			err = r.localCache.Set(ctx, cast.ToString(id), data, 2*time.Minute)
			if err != nil {
				return nil, errors.WithMessage(err, "[repo] GetCompany localCache set error")
			}
		}
		return data, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*model.Company), nil
}

// GetCompany get a record
func (r *companyRepo) GetCompanyByAddress(ctx context.Context, walletAddress string) (ret *model.Company, err error) {
	if walletAddress == "" {
		return nil, errors.New("[repo] GetCompany walletAddress cannot be equal to 0")
	}

	// get data from db
	// 避免缓存击穿(瞬间有大量请求过来)
	val, err, _ := r.sg.Do("sg:company:address:"+walletAddress, func() (interface{}, error) {
		// read db or rpc
		data, err := dao.Company.WithContext(ctx).Where(dao.Company.WalletAddress.Eq(walletAddress)).Last()
		if err != nil {
			return nil, errors.Wrapf(err, "[repo] get Company from db error, walletAddress: %s", walletAddress)
		}

		// write cache
		if data != nil && data.ID > 0 {
			// write redis
			err = r.cache.SetCompanyCache(ctx, data.ID, data, 5*time.Minute)
			if err != nil {
				return nil, errors.WithMessage(err, "[repo] GetCompanyByAddress SetCompanyCache error")
			}

			// write local cache
			err = r.localCache.Set(ctx, cast.ToString(data.ID), data, 2*time.Minute)
			if err != nil {
				return nil, errors.WithMessage(err, "[repo] GetCompanyByAddress localCache set error")
			}
		}
		return data, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*model.Company), nil
}

// BatchGetCompany batch get items
func (r *companyRepo) BatchGetCompany(ctx context.Context, ids []int64) (ret []*model.Company, err error) {
	// read cache
	itemMap, err := r.cache.MultiGetCompanyCache(ctx, ids)
	if err != nil {
		return nil, err
	}
	var missedID []int64
	for _, v := range ids {
		item, ok := itemMap[cast.ToString(v)]
		if !ok {
			missedID = append(missedID, v)
			continue
		}
		ret = append(ret, item)
	}
	// get missed data
	if len(missedID) > 0 {
		missedData, err := dao.Company.WithContext(ctx).Where(dao.Company.ID.In(ids...)).Find()
		if err != nil {
			// you can degrade to ignore error
			return nil, err
		}
		if len(missedData) > 0 {
			ret = append(ret, missedData...)
			err = r.cache.MultiSetCompanyCache(ctx, missedData, 5*time.Minute)
			if err != nil {
				// you can degrade to ignore error
				return nil, err
			}
		}
	}
	return ret, nil
}
