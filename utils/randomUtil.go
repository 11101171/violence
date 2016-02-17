package utils

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : util/RandomUtil.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-18 18:24:57
// Description :
//****************************************************/

import (
	"fmt"
	"math/rand"
	"time"
)

var randomUtil *RandomUtil = &RandomUtil{}

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type RandomUtil struct{}

func Random() *RandomUtil {
	return randomUtil
}

func (ru *RandomUtil) RandomString(pixff string, strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return time.Now().Format("20151212010203") + "-" + pixff + "-" + string(result)
}

func main() {
	fmt.Println("Start Main func()")
}
