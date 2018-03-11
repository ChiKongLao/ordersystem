package manager

func GetRedisMenuKey(key ...string) string {
	return generationRedisKey("menu", key)
}

func GetRedisOrderKey(key ...string) string {
	return generationRedisKey("order", key)
}

// 设置单个值
func SetValue(key string, value interface{}) {
	SetValueWithExpire(key, value, -1)
}

// 设置单个值, 带过期时间
func SetValueWithExpire(key string, value interface{}, time int64) {
	GetRedisConn().Do(RedisSet, key, value, RedisEx, time)
}

// 获取单个值
func GetValue(key string) {
	GetRedisConn().Do(RedisGet, key)
}

// 设置多个值
func SetValues(key string, value interface{}) {
	GetRedisConn().Do(RedisRPush, key, value)
}

// 设置多个值, 带过期时间
func SetValuesWithExpire(key string, value interface{}, time int64) {
	SetValues(key, value)
	GetRedisConn().Do(RedisExpire, time)
}

// 获取多个值
func GetValues(key string) {
	GetRedisConn().Do(RedisGet, key)
}

// 生成redis key
func GenerationRedisKey(pre string, key ...string) string {
	return generationRedisKey(pre, key)

}

// 生成redis key
func generationRedisKey(pre string, key []string) string {
	res := pre
	for _, subItem := range key {
		res = res + "_" + subItem
	}
	return res

}
