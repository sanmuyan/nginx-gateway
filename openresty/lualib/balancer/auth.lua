local ngx = ngx
local _M = {}

function _M.auth_token(token)
    local x_token = nil
    local ok = pcall(function()
        x_token = ngx.req.get_headers()["authorization"]
    end)
    if not ok or token ~= x_token then
        ngx.status = ngx.HTTP_FORBIDDEN
        ngx.say("token auth is invalid")
    end
end

function _M.auth_basic(username, password)
    local req_basic = nil
    local basic = nil
    local ok = pcall(function()
        local authorization = ngx.req.get_headers()["authorization"]
        req_basic = string.sub(authorization, 7, string.len(authorization))
        basic = ngx.encode_base64(username .. ":" .. password)
    end)
    if not ok or req_basic ~= basic then
        ngx.status = ngx.HTTP_FORBIDDEN
        ngx.say("basic auth is invalid")
    end
end

function _M.auth_request(auth_path)
    local res = ngx.location.capture(auth_path)
    if res.status ~= ngx.HTTP_OK then
        ngx.status = ngx.HTTP_FORBIDDEN
        ngx.say("auth request is invalid")
    end
end

return _M
