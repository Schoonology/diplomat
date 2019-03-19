#!/usr/bin/env bats

load helpers/helpers

@test "Invalid JSON Schema" {
  run bin/diplomat $FIXTURES_ROOT/failing/invalid-json-schema.txt $TEST_HOST

  log_on_failure

  [ $status -eq 1 ]
  [[ "$output" =~ "Error while running Lua script" ]]
  [[ "$output" =~ "has a primitive type that is NOT VALID" ]]
}

@test "Fail to match JSON Schema" {
  run bin/diplomat $FIXTURES_ROOT/failing/fail-to-match-schema.txt $TEST_HOST

  log_on_failure

  [ $status -eq 1 ]
  [[ "$output" = "Error while running Lua script:
	(root): Additional property slideshow is not allowed" ]]
}

@test "Multiple JSON errors in a Markdown file" {
  run bin/diplomat $FIXTURES_ROOT/failing/multiple-json-failures.md $TEST_HOST

  log_on_failure

  [ $status -eq 2 ]
  [ ${lines[0]} = "Error while running Lua script:" ]
  [[ ${lines[1]} =~ "has a primitive type that is NOT VALID" ]]
  [ ${lines[2]} = "Error while running Lua script:" ]
  [[ ${lines[3]} =~ "Additional property slideshow is not allowed" ]]
}

@test "Error from template function" {
  run bin/diplomat --script $FIXTURES_ROOT/failing/error-from-template.lua \
    $FIXTURES_ROOT/failing/error-from-template.txt $TEST_HOST

  log_on_failure

  [ $status -eq 1 ]
  [[ "$output" =~ "Error while running Lua script:" ]]
  [[ "$output" =~ "test/fixtures/failing/error-from-template.lua:2: Template failed!" ]]
}

@test "Error from validator function" {
  run bin/diplomat --script $FIXTURES_ROOT/failing/error-from-validator.lua \
    $FIXTURES_ROOT/failing/error-from-validator.txt $TEST_HOST

  log_on_failure

  [ $status -eq 1 ]
  [[ "$output" =~ "Error while running Lua script:" ]]
  [[ "$output" =~ "test/fixtures/failing/error-from-validator.lua:2: Validator failed!" ]]
}
