#!/usr/bin/env bats

@test "Text: GET 422 Unprocessable Entity: Test Title" {
  run ./main test/fixtures/match-get-422.txt http://httpbin.org

  echo actual status: "$status"
  [ "$status" -eq 0 ]

  echo actual output: "$output" 
  [ "$output" = "GET /status/422 -> 422" ]
}

@test "Text: GET 422 Unprocessable Entity: Diff" {
  run ./main test/fixtures/match-get-422.txt http://httpbin.org

  echo actual status: "$status"
  [ "$status" -eq 0 ]

  echo actual lines[1]: "${lines[1]}"
  [ "${lines[1]}" = "" ]
}

@test "Text: GET 200 OK" {
  run ./main test/fixtures/match-get-200.txt http://httpbin.org

  echo actual status: "$status"
  [ "$status" -eq 0 ]

  echo actual lines[1]: "${lines[1]}"
  [ "${lines[1]}" = "" ]
}

@test "Text: GET 200 OK with Body" {
  run ./main test/fixtures/match-get-200-with-body.txt http://httpbin.org

  echo actual status: "$status"
  [ "$status" -eq 0 ]

  echo actual lines[1]: "${lines[1]}"
  [ "${lines[1]}" = "" ]
}

@test "Text: POST" {
  run ./main test/fixtures/match-post-200.txt http://httpbin.org

  echo actual status: "$status"
  [ "$status" -eq 0 ]

  echo actual lines[1]: "${lines[1]}"
  [ "${lines[1]}" = "" ]
}

# TODO: this test should pass, despite having brackets in the body
@test "Text: POST with HTTP response" {
  skip
  run ./main test/fixtures/match-post-200-http.txt http://httpbin.org

  echo actual status: "$status"
  [ "$status" -eq 0 ]

  echo actual output: "$output" 
  [ "$output" = "POST /post -> 200" ]
}

@test "Markdown: .markdown: GET 422 Unprocessable Entity" {
  run ./main test/fixtures/match-get-422.markdown http://httpbin.org

  echo actual status: "$status"
  [ "$status" -eq 0 ]

  echo actual output: "$output" 
  [ "$output" = "Markdown: GET /status/422" ]
}

@test "Markdown: .md: GET 422 Unprocessable Entity" {
  run ./main test/fixtures/match-get-422.md http://httpbin.org

  echo actual status: "$status"
  [ "$status" -eq 0 ]

  echo actual output: "$output" 
  [ "$output" = "Markdown: GET /status/422" ]
}

@test "Markdown: Multiple specs" {
  run ./main test/fixtures/multiple.markdown http://httpbin.org

  echo actual status: "$status"
  [ "$status" -eq 0 ]

  echo actual output: $output
  [ "${lines[0]}" = "Markdown: GET /status/422" ]
  [ "${lines[1]}" = "Markdown: GET /status/200" ]
  [ "${lines[2]}" = "" ]
}

@test "--tap Option" {
  run ./main --tap test/fixtures/match-get-422.txt http://httpbin.org

  echo actual status: "$status"
  [ "$status" -eq 0 ]

  echo actual output: $output
  [ "${lines[0]}" = "TAP version 13" ]
  [ "${lines[1]}" = "1..1" ]
  [ "${lines[2]}" = "ok 0 GET /status/422 -> 422" ]
  [ "${lines[3]}" = "" ]
}