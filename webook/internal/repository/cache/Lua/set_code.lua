-- set_code.lua
local key    = KEYS[1]
local cntkey = key..":cnt"
local val    = ARGV[1]

-- 1）先读现有的 TTL
local ttl = tonumber(redis.call("ttl", key))

-- 2）判断 ttl 情况
if ttl == -1 then
    -- 说明 key 存在但没有过期时间，不安全，返回 -2
    return -2
elseif ttl and ttl > 540 then
    -- 冷却中，不允许重复发送
    return -1
end

-- 3）设置验证码、10分钟有效、尝试次数
redis.call("set", key, val)
redis.call("expire", key, 600)
redis.call("set", cntkey, 3)
redis.call("expire", cntkey, 600)

return 0
