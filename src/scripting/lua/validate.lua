-- Much simplified regexp for parsing dates.
-- Not as strict as it could be, ON PURPOSE.
-- Double-escaped because this gets generated into Go, and re-escaped.
local __time_fragment = "[\\\\d]{2}:[\\\\d]{2}:[\\\\d]{2}"
local __rfc1123_date_fragment = "[\\\\d]{2}[\\\\s][A-Z][a-z]{2}[\\\\s][\\\\d]{4}"
local __rfc850_date_fragment = "[\\\\d]{2}-[A-Z][a-z]{2}-[\\\\d]{2}"
local __asctime_date_fragment = "[A-Z][a-z]{2}[\\\\s][\\\\d\\\\s][\\\\d]"
local __date_fragment = "("..__rfc1123_date_fragment.."|"..__rfc850_date_fragment.."|"..__asctime_date_fragment..")"
local date_regexp = "[MTWFS][a-z]+,?[\\\\s]"..__date_fragment.."[\\\\s]+"..__time_fragment.."[\\\\s](GMT|[\\\\d]{4})"
function is_date(value)
  return re.match(value, date_regexp) ~= nil
end

function is_test(value)
  return value == "test"
end

function json_schema(schema)
  return function (value)
    return __validateJSON(schema, value)
  end
end

function regexp(pattern)
  return function (value)
    return re.match(value, pattern) ~= nil
  end
end
