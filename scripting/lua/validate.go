/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package lua

// validateTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func validateTemplate() string {
	var tmpl = "function is_true(value)\n" +
		"  return value == \"true\"\n" +
		"end\n" +
		""
	return tmpl
}
