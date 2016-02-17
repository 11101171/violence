package services

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/cache"
)

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : cacheService.go
// Author      : ningzhong.zeng
// Revision    : 2016-01-12 19:47:20
// Description :
//****************************************************/

const (
	timeout int64 = 1000
)

type CacheService struct {
	SeftCache cache.Cache
}

// basie
func (this *CacheService) Put(key string, value interface{}) {
	this.SeftCache.Put(key, value, timeout)
}
func (this *CacheService) Get(key string) interface{} {
	return this.SeftCache.Get(key)
}

// 设置 api 配置的内容
func (this *CacheService) PutServerContent(key string, value string) {
	value = strings.Replace(value, " ", "", -1)
	this.Put(key, value)
}

// 获取 api设置的内容
func (this *CacheService) GetServerContent(key string) string {
	return fmt.Sprintf("%s", this.Get(key))
}
