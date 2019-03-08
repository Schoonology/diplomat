/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package lua

// validateTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func validateTemplate() string {
	var tmpl = "function is_test(value)\n" +
		"  return value == \"test\"\n" +
		"end\n" +
		"\n" +
		"function json_schema(schema)\n" +
		"  return function (value)\n" +
		"    return __validateJSON(schema, value)\n" +
		"  end\n" +
		"end\n" +
		""
	return tmpl
}
