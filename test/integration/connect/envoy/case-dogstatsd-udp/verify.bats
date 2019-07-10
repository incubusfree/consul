#!/usr/bin/env bats

load helpers

@test "s1 proxy admin is up on :19000" {
  retry_default curl -f -s localhost:19000/stats -o /dev/null
}

@test "s2 proxy admin is up on :19001" {
  retry_default curl -f -s localhost:19001/stats -o /dev/null
}

@test "s2 proxy should be healthy" {
  assert_service_has_healthy_instances s2 1
}

@test "s1 upstream should be able to connect to s2" {
  run retry_default curl -s -f -d hello localhost:5000

  echo "OUTPUT: $output"

  [ "$status" == 0 ]
  [ "$output" == "hello" ]
}

@test "s1 proxy should be sending metrics to statsd" {
  run retry_default cat /workdir/statsd/statsd.log

  echo "METRICS:"
  echo "$output"
  echo "COUNT: $(echo "$output" | grep -Ec '^envoy\.')"

  [ "$status" == 0 ]
  [ $(echo $output | grep -Ec '^envoy\.') -gt "0" ]
}

@test "s1 proxy should be sending dogstatsd tagged metrics" {
  run retry_default must_match_in_statsd_logs '[#,]local_cluster:s1(,|$)'

  echo "OUTPUT: $output"

  [ "$status" == 0 ]
}

@test "s1 proxy should be adding cluster name as a tag" {
  run retry_default must_match_in_statsd_logs '[#,]envoy.cluster_name:s2(,|$)'

  echo "OUTPUT: $output"

  [ "$status" == 0 ]
}

@test "s1 proxy should be sending additional configured tags" {
  run retry_default must_match_in_statsd_logs '[#,]foo:bar(,|$)'

  echo "OUTPUT: $output"

  [ "$status" == 0 ]
}

@test "s1 proxy should have custom stats flush interval" {
  INTERVAL=$(get_envoy_stats_flush_interval localhost:19000)

  echo "INTERVAL = $INTERVAL"

  [ "$INTERVAL" == "1s" ]
}
