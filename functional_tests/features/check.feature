Feature: test health check
  In order to test the cache returns a correct health check
  As a developer
  I need to test the gRPC Check interface

  Scenario: health ok
    Given the server is running
    When I call check
    Then the response should be running
