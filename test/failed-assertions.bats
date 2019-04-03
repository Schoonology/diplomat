#!/usr/bin/env bats

load helpers/helpers

@test "Invalid JSON Schema" {
  run bin/diplomat $FIXTURES_ROOT/failing/invalid-json-schema.txt --address $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "$output" =~ "Error while running Lua script" ]]
  [[ "$output" =~ "has a primitive type that is NOT VALID" ]]
}

@test "Fail to match JSON Schema" {
  run bin/diplomat $FIXTURES_ROOT/failing/fail-to-match-schema.txt --address $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "$output" = "✗ GET /json -> 200
Error while running Lua script:
	(root): Additional property slideshow is not allowed" ]]
}

@test "Multiple JSON errors in a Markdown file" {
  run bin/diplomat $FIXTURES_ROOT/failing/multiple-json-failures.md --address $TEST_HOST

  log_on_failure

  [ "$status" -eq 2 ]
  [ "${lines[0]}" = "✗ Invalid JSON Schema" ]
  [ "${lines[1]}" = "Error while running Lua script:" ]
  [[ "${lines[2]}" =~ "has a primitive type that is NOT VALID" ]]
  [ "${lines[3]}" = "✗ Fail to Match JSON Schema" ]
  [ "${lines[4]}" = "Error while running Lua script:" ]
  [[ "${lines[5]}" =~ "Additional property slideshow is not allowed" ]]
}

@test "Multiple JSON errors in a Markdown file with TAP output" {
  run bin/diplomat $FIXTURES_ROOT/failing/multiple-json-failures.md --address $TEST_HOST --tap

  log_on_failure

  [ "$status" -eq 2 ]
  [ "${lines[0]}" = "TAP version 13" ]
  [ "${lines[1]}" = "not ok 1 Invalid JSON Schema" ]
  [ "${lines[2]}" = "not ok 2 Fail to Match JSON Schema" ]
  [ "${lines[3]}" = "Invalid JSON Schema:" ]
  [ "${lines[4]}" = "Error while running Lua script:" ]
  [[ "${lines[5]}" =~ "has a primitive type that is NOT VALID" ]]
  [ "${lines[6]}" = "Fail to Match JSON Schema:" ]
  [ "${lines[7]}" = "Error while running Lua script:" ]
  [[ "${lines[8]}" =~ "Additional property slideshow is not allowed" ]]
}

@test "Error from template function" {
  run bin/diplomat --script $FIXTURES_ROOT/failing/error-from-template.lua \
    $FIXTURES_ROOT/failing/error-from-template.txt \
    --address $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "$output" =~ "Error while running Lua script:" ]]
  [[ "$output" =~ "test/fixtures/failing/error-from-template.lua:2: Template failed!" ]]
}

@test "Error from validator function" {
  run bin/diplomat --script $FIXTURES_ROOT/failing/error-from-validator.lua \
    $FIXTURES_ROOT/failing/error-from-validator.txt \
    --address $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "$output" =~ "Error while running Lua script:" ]]
  [[ "$output" =~ "test/fixtures/failing/error-from-validator.lua:2: Validator failed!" ]]
}
