local key = KEYS[1]
local expectcode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"
local cnt = tonumber(redis.call("get", cntKey))
if not cnt or cnt <= 0 then
    return -1  -- 尝试次数用完
end

if not code then
    return -3  -- 验证码不存在或过期
end

if expectcode == code then
    redis.call("del", key)
    redis.call("del", cntKey)
    return 0   -- 验证成功
else
    redis.call("decr", cntKey)
    return -2  -- 验证码错误，尝试次数减一
end

