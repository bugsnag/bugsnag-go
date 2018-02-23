# Possibly generic steps to upstream

When("I run {string} with the defaults") do |eventType|
  steps %Q{
    When I start Android emulator "newnexus"
    And I install the "com.bugsnag.android.mazerunner" Android app from "app/build/outputs/apk/release/app-release.apk"
    And I clear the "com.bugsnag.android.mazerunner" Android app data
    And I set environment variable "BUGSNAG_API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And I set environment variable "EVENT_TYPE" to "#{eventType}"
    And I start the "com.bugsnag.android.mazerunner" Android app using the "com.bugsnag.android.mazerunner.MainActivity" activity
  }
end
