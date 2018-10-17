Feature: test cache get
  In order to test the cache can return images
  As a developer
  I need to test the gRPC Get interface

  Scenario: get file
    Given the server is running
    And a file exists in the cache
    When I Get that file
    Then the file contents should be returned
