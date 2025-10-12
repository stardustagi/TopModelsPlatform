package constants

var (
	PlatformJwtExpire       string = "72h"  // 72小时
	PlatformUserTokenExpire string = "720h" // 节点用户Token过期时间
	PlatformUserJwtExpire   string = "720h" // 节点用户Jwt过期时间
	RedisPrefix             string
	ModelsKeyPrefix         string
)
