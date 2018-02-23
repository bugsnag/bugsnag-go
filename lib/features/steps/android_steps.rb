
When("I start Android emulator {string}") do |emulator|
  steps %Q{
    When I set environment variable "ANDROID_EMULATOR" to "#{emulator}"
    And I run the script "launch-android-emulator.sh"
    And I run the script "await-android-emulator.sh" synchronously
  }
end

When("I clear the {string} Android app data") do |app|
  step('I run the script "clear-android-app-data.sh" synchronously')
end

When("I install the {string} Android app from {string}") do |bundle, filepath|
  steps %Q{
    When I set environment variable "APP_BUNDLE" to "#{bundle}"
    And I set environment variable "APK_PATH" to "#{filepath}"
    And I run the script "install-android-app.sh" synchronously
  }
end

When("I start the {string} Android app using the {string} activity") do |app, activity|
  steps %Q{
    When I set environment variable "APP_BUNDLE" to "#{app}"
    When I set environment variable "APP_ACTIVITY" to "#{activity}"
    And I run the script "launch-android-app.sh" synchronously
    And I wait for 4 seconds
  }
end
