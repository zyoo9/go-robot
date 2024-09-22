package rag_robot

import (
	"github.com/patrickmn/go-cache"
	"time"
)

// 内存缓存
var Cache = cache.New(10*time.Second, time.Minute)
