package services

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : services/CookieService.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-23 14:47:46
// Description :
//****************************************************/
type CookieService struct {
	BaseService
}

func (cs *CookieService) Get(key string) string {
	return cs.ctx.Input.Cookie(key)
}

func (cs *CookieService) Set(key string, value string) string {
	if key != "" && value != "" {
		cs.ctx.SetCookie(key, value, 10*60) // s
		return value
	}
	return "" // errors.New("Key And Value Not Null")
}
func (cs *CookieService) SetWithT(key string, value string, t interface{}) string {
	if key != "" && value != "" {
		cs.ctx.SetCookie(key, value, t) // s
		return value
	}
	return "" // errors.New("Key And Value Not Null")
}
