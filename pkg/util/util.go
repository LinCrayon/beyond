package util

import (
	"math/rand"
	"strconv"
	"time"
)

func RandomNumeric(size int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //创建随机数生成器
	if size <= 0 {
		panic("{ size : " + strconv.Itoa(size) + " } must be more than 0 ")
	}
	value := ""
	for index := 0; index < size; index++ {
		value += strconv.Itoa(r.Intn(10))
	}
	return value
}

// EndOfDay 返回给定时间当天的结束时间（午夜）
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location()) //地理位置（时区信息）
}
