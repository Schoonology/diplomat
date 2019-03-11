function log_on_failure() {
  echo Failed with status "$status" and output:
  echo "$output"
}

function setup() {
  FIXTURES_ROOT=$BATS_TEST_DIRNAME/fixtures
  TEST_HOST=http://localhost:7357
}
