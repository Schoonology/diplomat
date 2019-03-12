#!/usr/bin/env bats

load helpers/helpers

@test "Missing arguments" {
  run ./main

  log_on_failure

  [ $status -eq 1 ]
  [[ "${lines[0]}" =~ "error: required argument" ]]
}

@test "File does not exist" {
  run ./main missing $TEST_HOST

  log_on_failure

  [ $status -eq 1 ]
  [[ "${lines[0]}" =~ "does not exist" ]]
}

@test "Host does not exist" {
  run ./main $FIXTURES_ROOT/match-get-200.txt http://wrong

  log_on_failure

  [ $status -eq 3 ]
  [[ "${lines[0]}" =~ "Could not resolve host" ]]
}

@test "Host is unreachable" {
  run ./main $FIXTURES_ROOT/match-get-200.txt http://localhost:7538

  log_on_failure

  [ $status -eq 3 ]
  [[ "${lines[0]}" =~ "Failed to connect" ]]
}
