function is_test(value)
  return value == "test"
end

function json_schema(schema)
  return function (value)
    return __validateJSON(schema, value)
  end
end
