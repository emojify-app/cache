Feature: test cache put
  In order to test the cache can save images
  As a developer
  I need to test the gRPC Put interface

  Scenario: put file
    Given the server is running
    When I put a file
    Then the file should exist in the cache
