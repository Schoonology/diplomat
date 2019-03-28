/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package lua

// contextTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func contextTemplate() string {
	var tmpl = "local function create_context()\n" +
		"  local data = {}\n" +
		"\n" +
		"  local function debug()\n" +
		"    print(\"Context: \")\n" +
		"    for k,v in pairs(data) do\n" +
		"      print(k,type(v),v)\n" +
		"    end\n" +
		"    return \"\"\n" +
		"  end\n" +
		"\n" +
		"  local function get(key)\n" +
		"    return data[key]\n" +
		"  end\n" +
		"\n" +
		"  local function set(key, ...)\n" +
		"    local has_value = select('#', ...) > 0\n" +
		"\n" +
		"    local function _set(value)\n" +
		"      data[key] = value\n" +
		"      return value\n" +
		"    end\n" +
		"\n" +
		"    if has_value then\n" +
		"      return _set(select(1, ...))\n" +
		"    else\n" +
		"      return _set\n" +
		"    end\n" +
		"  end\n" +
		"\n" +
		"  return {\n" +
		"    debug = debug,\n" +
		"    get = get,\n" +
		"    set = set,\n" +
		"  }\n" +
		"end\n" +
		"\n" +
		"ctx = create_context()\n" +
		""
	return tmpl
}
