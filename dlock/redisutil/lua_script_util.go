package redisutil

func GetLockScript() string {
	return `
		if redis.call("EXISTS", KEYS[1]) == 0 or redis.call("HEXISTS", KEYS[1], ARGV[1]) == 1 then
			redis.call("HINCRBY", KEYS[1], ARGV[1], 1)
			redis.call("EXPIRE", KEYS[1], ARGV[2])
			return 1
		else
			return 0
		end
	`
}

func GetUnlockScript() string {
	return `
		if redis.call("HEXISTS", KEYS[1], ARGV[1]) == 1 then
			if redis.call("HINCRBY", KEYS[1], ARGV[1], -1) == 0 then
				redis.call("DEL", KEYS[1])
			else
				redis.call("EXPIRE", KEYS[1], ARGV[2])
			end
			return 1
		else
			return 0
		end
	`
}

func GetAutoRenewScript() string {
	return `
		if redis.call("HEXISTS", KEYS[1], ARGV[1]) == 1 then
			return redis.call("EXPIRE", KEYS[1], ARGV[2])
		else
			return 0
		end
`
}
