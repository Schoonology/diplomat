function file(filename)
  local handle = assert(io.open(filename, "rb"))
  local contents = assert(handle:read("*a"))
  handle:close()
  return contents
end
