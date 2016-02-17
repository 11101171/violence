package services

import (
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/context"
)

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : services/ServiceFactory.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-23 15:03:25
// Description :
//****************************************************/

type BaseServiceInterface interface {
	WithContext(ctx *context.Context, fn func())
	Factory() *ServiceFactory
}

type BaseService struct {
	ctx *context.Context
}

type ServiceFactory struct {
	CookieService *CookieService
	AuthService   *AuthService
	CacheService  *CacheService
	SSHService    *SSHService
}

var once sync.Once
var instance *ServiceFactory

func GetInstance() *ServiceFactory {
	once.Do(func() {
		if instance == nil {
			bm, err := cache.NewCache("memory", `{"interval":60}`)
			if err != nil {
				beego.Error("cache init fail=>", err)
			}
			cacheService := CacheService{
				SeftCache: bm,
			}
			instance = &ServiceFactory{
				CookieService: &CookieService{},
				AuthService:   &AuthService{},
				CacheService:  &cacheService,
				SSHService:    &SSHService{},
			}
		}
	})
	return instance
}

// Factory
func (bs *BaseService) Factory() *ServiceFactory {
	return GetInstance()
}

// Add Context
func (bs *BaseService) WithContext(ctx *context.Context, fn func()) {
	bs.ctx = ctx
	fn()
}

// Add Context And Return bool
func (bs *BaseService) WithContextBool(ctx *context.Context, fn func() bool) bool {
	bs.ctx = ctx
	return fn()
}

// Add Context And Return interface{}
func (bs *BaseService) WithContextReturn(ctx *context.Context, fn func() interface{}) interface{} {
	bs.ctx = ctx
	return fn()
}
