local function compare(proto, value)
  if type(proto) == "table" then
    for k,v in pairs(proto) do
      if not compare(proto[k], value[k]) then return false end
    end
    return true
  end

  return proto == value
end

function chain(...)
  local pipeline = {...}
  return function (value)
    local result = value

    for _,segment in ipairs(pipeline) do
      result = segment(result)
    end

    return result
  end
end

function env(key)
  return os.getenv(key)
end

function equal(proto)
  return function (value)
    return compare(proto, value)
  end
end

function file(filename)
  local handle = assert(io.open(filename, "rb"))
  local contents = assert(handle:read("*a"))
  handle:close()
  return contents
end

function get(key, validator)
  return function (value)
    return value[key]
  end
end
