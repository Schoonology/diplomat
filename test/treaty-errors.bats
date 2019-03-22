#!/usr/bin/env bats

load helpers/helpers

@test "Bad request header" {
  run bin/diplomat $FIXTURES_ROOT/broken/bad-request-header.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Error building spec: line 1" ]
  [[ "${lines[1]}" =~ "Failed to parse header: Content-Type" ]]
}

@test "Bad request line" {
  run bin/diplomat $FIXTURES_ROOT/broken/bad-request-line.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Error building spec: line 1" ]
  [[ "${lines[1]}" =~ "Failed to parse request line: INVALID" ]]
}

@test "Multiple bad requests in a markdown file" {
  run bin/diplomat $FIXTURES_ROOT/broken/bad-request-multiple.md $TEST_HOST

  log_on_failure

  [ "$status" -eq 2 ]
  [ "${lines[0]}" = "Error building spec: line 6" ]
  [[ "${lines[1]}" =~ "Failed to parse header: Content-Type" ]]
  [ "${lines[2]}" = "Error building spec: line 17" ]
  [[ "${lines[3]}" =~ "Failed to parse request line: INVALID" ]]
}

@test "Bad response header" {
  run bin/diplomat $FIXTURES_ROOT/broken/bad-response-header.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Error building spec: line 1" ]
  [[ "${lines[1]}" =~ "Failed to parse header: Content-Type" ]]
}

@test "Bad response status" {
  run bin/diplomat $FIXTURES_ROOT/broken/bad-response-status.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Error building spec: line 1" ]
  [[ "${lines[1]}" =~ "Failed to parse response line: OOPS" ]]
}

@test "Markdown: No header" {
  run bin/diplomat $FIXTURES_ROOT/broken/no-header.md $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [[ "$output" = "GET /status/200 -> 200" ]]
}

@test "Markdown: Unclosed code fence" {
  run bin/diplomat $FIXTURES_ROOT/broken/unclosed-code-fence.md $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [[ "$output" = "Unclosed Code Fence" ]]
}

@test "Request only" {
  run bin/diplomat $FIXTURES_ROOT/broken/request-only.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Error building spec: line 1" ]
  [[ "${lines[1]}" =~ "Found a request without a corresponding response." ]]
}

@test "Response only" {
  run bin/diplomat $FIXTURES_ROOT/broken/response-only.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Error building spec: line 1" ]
  [[ "${lines[1]}" =~ "Found a response without a corresponding request." ]]
}

@test "Request only (second assertion)" {
  run bin/diplomat $FIXTURES_ROOT/broken/request-only-second.md $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Error building spec: line 14" ]
  [[ "${lines[1]}" =~ "Found a request without a corresponding response." ]]
  [ "${lines[2]}" = "First: Correct" ]
}

@test "Response only (second assertion)" {
  run bin/diplomat $FIXTURES_ROOT/broken/response-only-second.md $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Error building spec: line 14" ]
  [[ "${lines[1]}" =~ "Found a response without a corresponding request." ]]
  [ "${lines[2]}" = "First: Correct" ]
}

@test "Missing template function" {
  run bin/diplomat $FIXTURES_ROOT/broken/missing-template-function.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "$output" = 'Template `missing` could not be found.' ]]
}

@test "Missing validator function" {
  run bin/diplomat $FIXTURES_ROOT/broken/missing-validator-function.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "$output" =~ "Error while running Lua script" ]]
  [[ "$output" =~ "attempt to call a non-function object" ]]
}

@test "Invalid script syntax" {
  run bin/diplomat --script $FIXTURES_ROOT/broken/invalid-script-syntax.lua \
    $FIXTURES_ROOT/match-get-200.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "$output" =~ "Syntax error while parsing custom script:" ]]
  [[ "$output" =~ "test/fixtures/broken/invalid-script-syntax.lua:1:41: {" ]]
}
