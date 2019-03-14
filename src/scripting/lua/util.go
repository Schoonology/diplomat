/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package lua

// utilTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func utilTemplate() string {
	var tmpl = "function file(filename)\n" +
		"  local handle = assert(io.open(filename, \"rb\"))\n" +
		"  local contents = assert(handle:read(\"*a\"))\n" +
		"  handle:close()\n" +
		"  return contents\n" +
		"end\n" +
		""
	return tmpl
}
