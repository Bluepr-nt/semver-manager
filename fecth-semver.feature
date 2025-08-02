Feature: Fetch semver tag on a git repository

  Scenario: Fetching the highest semver tag from a git repository
    Given A git repository with multiple tags
    When I fetch highest semver tag
    Then I get the highest semver tag
