/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package lua

// validateTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func validateTemplate() string {
	var tmpl = "-- Much simplified regexp for parsing dates.\n" +
		"-- Not as strict as it could be, ON PURPOSE.\n" +
		"-- Double-escaped because this gets generated into Go, and re-escaped.\n" +
		"local __time_fragment = \"[\\\\d]{2}:[\\\\d]{2}:[\\\\d]{2}\"\n" +
		"local __rfc1123_date_fragment = \"[\\\\d]{2}[\\\\s][A-Z][a-z]{2}[\\\\s][\\\\d]{4}\"\n" +
		"local __rfc850_date_fragment = \"[\\\\d]{2}-[A-Z][a-z]{2}-[\\\\d]{2}\"\n" +
		"local __asctime_date_fragment = \"[A-Z][a-z]{2}[\\\\s][\\\\d\\\\s][\\\\d]\"\n" +
		"local __date_fragment = \"(\"..__rfc1123_date_fragment..\"|\"..__rfc850_date_fragment..\"|\"..__asctime_date_fragment..\")\"\n" +
		"local date_regexp = \"[MTWFS][a-z]+,?[\\\\s]\"..__date_fragment..\"[\\\\s]+\"..__time_fragment..\"[\\\\s](GMT|[\\\\d]{4})\"\n" +
		"function is_date(value)\n" +
		"  return re.match(value, date_regexp) ~= nil\n" +
		"end\n" +
		"\n" +
		"function is_test(value)\n" +
		"  return value == \"test\"\n" +
		"end\n" +
		"\n" +
		"function json_schema(schema)\n" +
		"  return function (value)\n" +
		"    return __validateJSON(schema, value)\n" +
		"  end\n" +
		"end\n" +
		"\n" +
		"function regexp(pattern)\n" +
		"  return function (value)\n" +
		"    return re.match(value, pattern) ~= nil\n" +
		"  end\n" +
		"end\n" +
		""
	return tmpl
}
