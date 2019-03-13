#!/usr/bin/env bats

load helpers/helpers

@test "Text: Fallback title" {
  run bin/diplomat $FIXTURES_ROOT/match-get-422.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "$output" = "GET /status/422 -> 422" ]
}

@test "Text: Empty diff on success" {
  run bin/diplomat $FIXTURES_ROOT/match-get-422.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "${#lines[@]}" = 1 ]
}

@test "Text: Status diff on incorrect status" {
  run bin/diplomat $FIXTURES_ROOT/fail-get-422.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "$output" = "GET /status/422 -> 400
Status:
	- 400 BAD REQUEST
	+ 422 UNPROCESSABLE ENTITY" ]
}

@test "Text: GET 200 OK" {
  run bin/diplomat $FIXTURES_ROOT/match-get-200.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "${#lines[@]}" = 1 ]
}

@test "Text: GET 200 OK with Body" {
  run bin/diplomat $FIXTURES_ROOT/match-get-200-with-body.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "${#lines[@]}" = 1 ]
}

@test "Text: POST" {
  run bin/diplomat $FIXTURES_ROOT/match-post-200.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "${#lines[@]}" = 1 ]
}

# TODO: this test should pass, despite having brackets in the body
@test "Text: POST with HTTP response" {
  skip
  run bin/diplomat $FIXTURES_ROOT/match-post-200-http.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "${#lines[@]}" = 1 ]
}

@test "Markdown: .markdown: Title from header" {
  run bin/diplomat $FIXTURES_ROOT/match-get-422.markdown $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "$output" = "Markdown: GET /status/422" ]
}

@test "Markdown: .md: Title from header" {
  run bin/diplomat $FIXTURES_ROOT/match-get-422.md $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "$output" = "Markdown: GET /status/422" ]
}

@test "Markdown: Multiple specs" {
  run bin/diplomat $FIXTURES_ROOT/multiple.markdown $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "Markdown: GET /status/422" ]
  [ "${lines[1]}" = "Markdown: GET /status/200" ]
  [ "${lines[2]}" = "" ]
}

@test "--tap Option" {
  run bin/diplomat --tap $FIXTURES_ROOT/match-get-422.txt $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "TAP version 13" ]
  [ "${lines[1]}" = "ok 0 GET /status/422 -> 422" ]
  [ "${lines[2]}" = "" ]
}

@test "JSON Schema" {
  run bin/diplomat $FIXTURES_ROOT/json-schema.md $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "$output" = "JSON Schema Test" ]
}

@test "Custom Script" {
  run bin/diplomat --script $FIXTURES_ROOT/custom.lua $FIXTURES_ROOT/custom-script.md $TEST_HOST

  log_on_failure

  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "Custom Script Test" ]
}

@test "Help" {
  run bin/diplomat --help

  log_on_failure

  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Usage: diplomat [<flags>] <filename> <address>" ]
  [[ "$output" =~ "Flags:" ]]
  [[ "$output" =~ "--debug" ]]
  [[ "$output" =~ "--help" ]]
  [[ "$output" =~ "--script" ]]
  [[ "$output" =~ "--tap" ]]
  [[ "$output" =~ "--version" ]]
  [[ "$output" =~ "Args:" ]]
}
