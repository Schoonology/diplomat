#!/usr/bin/env bats

load helpers/helpers

@test "Missing arguments" {
  run bin/diplomat

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~ "error: required flag --address not provided" ]]
}

@test "File does not exist" {
  run bin/diplomat missing --address $TEST_HOST

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~ "does not exist" ]]
}

@test "Host does not exist" {
  run bin/diplomat $FIXTURES_ROOT/match-get-200.txt --address http://wrong

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~ "Could not resolve host" ]]
}

@test "Host is unreachable" {
  run bin/diplomat $FIXTURES_ROOT/match-get-200.txt --address http://localhost:7538

  log_on_failure

  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~ "Failed to connect" ]]
}
