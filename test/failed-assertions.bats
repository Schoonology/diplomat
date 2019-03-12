#!/usr/bin/env bats

load helpers/helpers

@test "Invalid JSON Schema" {
  run bin/diplomat $FIXTURES_ROOT/failing/invalid-json-schema.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" =~ "Error while running Lua script" ]]
  [[ "$output" =~ "has a primitive type that is NOT VALID" ]]
}

@test "Fail to match JSON Schema" {
  run bin/diplomat $FIXTURES_ROOT/failing/fail-to-match-schema.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" = "Error while running Lua script:
	(root): Additional property slideshow is not allowed" ]]
}

@test "Error from template function" {
  run bin/diplomat --script $FIXTURES_ROOT/failing/error-from-template.lua \
    $FIXTURES_ROOT/failing/error-from-template.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" =~ "Error while running Lua script:" ]]
  [[ "$output" =~ "test/fixtures/failing/error-from-template.lua:2: Template failed!" ]]
}

@test "Error from validator function" {
  run bin/diplomat --script $FIXTURES_ROOT/failing/error-from-validator.lua \
    $FIXTURES_ROOT/failing/error-from-validator.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" =~ "Error while running Lua script:" ]]
  [[ "$output" =~ "test/fixtures/failing/error-from-validator.lua:2: Validator failed!" ]]
}
