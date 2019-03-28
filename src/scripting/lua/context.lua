local function create_context()
  local data = {}

  local function debug()
    print("Context: ")
    for k,v in pairs(data) do
      print(k,type(v),v)
    end
    return ""
  end

  local function get(key)
    return data[key]
  end

  local function set(key, ...)
    local has_value = select('#', ...) > 0

    local function _set(value)
      data[key] = value
      return value
    end

    if has_value then
      return _set(select(1, ...))
    else
      return _set
    end
  end

  return {
    debug = debug,
    get = get,
    set = set,
  }
end

ctx = create_context()
