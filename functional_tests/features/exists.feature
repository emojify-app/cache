Feature: test cache exists
  In order to test the cache returns if an image exists
  As a developer
  I need to test the gRPC Put interface

  Scenario: exists file
    Given the server is running
    And a file exists in the cache
    When I call exists 
    Then the response should be true
