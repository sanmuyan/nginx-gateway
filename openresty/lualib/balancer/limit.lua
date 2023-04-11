local limit_req = require "resty.limit.req"
local limit_conn = require "resty.limit.conn"
local ngx = ngx
local _M = {}

function _M.limit_req_handler(rate, burst)
    -- rate: 每秒允许的请求数 burst: 允许额外的请求但是要延迟处理
    local lim, err = limit_req.new("limit_req_store", rate, burst)
    if not lim then
        ngx.log(ngx.ERR, "failed to instantiate a resty.limit.req object: ", err)
        return
    end

    local key = ngx.var.Host
    local delay, err = lim:incoming(key, true)
    if not delay then
        if err == "rejected" then
            ngx.status = ngx.HTTP_TOO_MANY_REQUESTS
            ngx.say("request rejected")
            return
        end
        ngx.log(ngx.ERR, "failed to limit req: ", err)
        return ngx.exit(ngx.HTTP_INTERNAL_SERVER_ERROR)
    end

    if delay > 0.001 then
        ngx.log(ngx.WARN, "delaying request by ", delay, " seconds")
        ngx.sleep(delay)
    end

end

function _M.limit_conn_handler(max, burst)
    -- max: 允许的最大连接数 burst: 允许额外的连接但是要延迟处理 0.1: 默认延迟时间
    local lim, err = limit_conn.new("limit_conn_store", max, burst, 0.1)
    if not lim then
        ngx.log(ngx.ERR, "failed to instantiate a resty.limit.conn object: ", err)
        return
    end

    local key = ngx.var.Host
    local delay, err = lim:incoming(key, true)
    if not delay then
        if err == "rejected" then
            ngx.status = ngx.HTTP_TOO_MANY_REQUESTS
            ngx.say("connection rejected")
            return
        end
        ngx.log(ngx.ERR, "failed to limit conn: ", err)
        return ngx.exit(ngx.HTTP_INTERNAL_SERVER_ERROR)
    end

    if lim:is_committed() then
        local ctx = ngx.ctx
        ctx.limit_conn = lim
        ctx.limit_conn_key = key
        ctx.limit_conn_delay = delay
    end

    local conn = err

    if delay > 0.001 then
        ngx.log(ngx.WARN, "delaying connection by ", delay, " seconds")
        ngx.sleep(delay)
    end

end

return _M
