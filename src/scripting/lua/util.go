/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package lua

// utilTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func utilTemplate() string {
	var tmpl = "local function compare(proto, value)\n" +
		"  if type(proto) == \"table\" then\n" +
		"    for k,v in pairs(proto) do\n" +
		"      if not compare(proto[k], value[k]) then return false end\n" +
		"    end\n" +
		"    return true\n" +
		"  end\n" +
		"\n" +
		"  return proto == value\n" +
		"end\n" +
		"\n" +
		"function chain(...)\n" +
		"  local pipeline = {...}\n" +
		"  return function (value)\n" +
		"    local result = value\n" +
		"\n" +
		"    for _,segment in ipairs(pipeline) do\n" +
		"      result = segment(result)\n" +
		"    end\n" +
		"\n" +
		"    return result\n" +
		"  end\n" +
		"end\n" +
		"\n" +
		"function env(key)\n" +
		"  return os.getenv(key)\n" +
		"end\n" +
		"\n" +
		"function equal(proto)\n" +
		"  return function (value)\n" +
		"    return compare(proto, value)\n" +
		"  end\n" +
		"end\n" +
		"\n" +
		"function file(filename)\n" +
		"  local handle = assert(io.open(filename, \"rb\"))\n" +
		"  local contents = assert(handle:read(\"*a\"))\n" +
		"  handle:close()\n" +
		"  return contents\n" +
		"end\n" +
		"\n" +
		"function get(key, validator)\n" +
		"  return function (value)\n" +
		"    return value[key]\n" +
		"  end\n" +
		"end\n" +
		""
	return tmpl
}
