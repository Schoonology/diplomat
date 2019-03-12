#!/usr/bin/env bats

load helpers/helpers

@test "Bad request header" {
  run ./main $FIXTURES_ROOT/broken/bad-request-header.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" = "Failed to parse header: Content-Type" ]]
}

@test "Bad request line" {
  run ./main $FIXTURES_ROOT/broken/bad-request-line.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" = "Failed to parse request line: INVALID" ]]
}

@test "Bad response header" {
  run ./main $FIXTURES_ROOT/broken/bad-response-header.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" = "Failed to parse header: Content-Type" ]]
}

@test "Bad response status" {
  run ./main $FIXTURES_ROOT/broken/bad-response-status.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" = "Failed to parse response line: OOPS" ]]
}

@test "Markdown: No header" {
  run ./main $FIXTURES_ROOT/broken/no-header.md $TEST_HOST

  log_on_failure

  [ $status -eq 0 ]
  [[ "$output" = "GET /status/200 -> 200" ]]
}

@test "Markdown: Unclosed code fence" {
  run ./main $FIXTURES_ROOT/broken/unclosed-code-fence.md $TEST_HOST

  log_on_failure

  [ $status -eq 0 ]
  [[ "$output" = "Unclosed Code Fence" ]]
}

@test "Request only" {
  run ./main $FIXTURES_ROOT/broken/request-only.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" = "Found a request without a corresponding response." ]]
}

@test "Response only" {
  run ./main $FIXTURES_ROOT/broken/response-only.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" = "Found a response without a corresponding request." ]]
}

@test "Request only (second assertion)" {
  run ./main $FIXTURES_ROOT/broken/request-only-second.md $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" = "Found a request without a corresponding response." ]]
}

@test "Response only (second assertion)" {
  run ./main $FIXTURES_ROOT/broken/response-only-second.md $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" = "Found a response without a corresponding request." ]]
}

@test "Missing template function" {
  run ./main $FIXTURES_ROOT/broken/missing-template-function.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" =~ "Error while running Lua script" ]]
  [[ "$output" =~ "attempt to call a non-function object" ]]
}

@test "Missing validator function" {
  run ./main $FIXTURES_ROOT/broken/missing-validator-function.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" =~ "Error while running Lua script" ]]
  [[ "$output" =~ "attempt to call a non-function object" ]]
}

@test "Invalid script syntax" {
  run ./main --script $FIXTURES_ROOT/broken/invalid-script-syntax.lua \
    $FIXTURES_ROOT/match-get-200.txt $TEST_HOST

  log_on_failure

  [ $status -eq 3 ]
  [[ "$output" =~ "Syntax error while parsing custom script:" ]]
  [[ "$output" =~ "test/fixtures/broken/invalid-script-syntax.lua:1:41: {" ]]
}
