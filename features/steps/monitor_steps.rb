When(/^I crash the app using (.*)$/) do |testcase|
  add_to_environment("TESTCASE", testcase)
  Dir.chdir(BUILD_DIR) do
    start_process([executable("panic-monitor"), executable("app")])
  end
end

When('I run the monitor with arguments {string}') do |args|
  Dir.chdir(BUILD_DIR) do
    start_process([executable("panic-monitor")] + args.split(' '))
  end
end

# More robust step to handle when more than split-every-space is needed
When("I run the monitor with:") do |table|
  Dir.chdir(BUILD_DIR) do
    args = table.raw.flatten
    start_process([executable("panic-monitor")] + args)
  end
end

Given('I set the API key to {string}') do |key|
  step("I set \"BUGSNAG_API_KEY\" to \"#{key}\" in the environment")
end

When('I set {string} to {string} in the environment') do |key, value|
  add_to_environment(key, value)
end

When('I set {string} to the sample app directory') do |key|
  add_to_environment(key, File.join(FIXTURE_DIR, 'app/'))
end

Then("payload field {string} equals {string}") do |keypath, expected_value|
  event = @server.events.last
  expect(event).not_to be_nil
  actual = JSON.parse(event.body)
  expect(actual["events"].length).to eq(1)
  expect(read_key_path(actual["events"][0], keypath)).to eq(expected_value)
end

Then("1 request was received") do
  step("1 requests were received")
end

Then("{int} requests were received") do |count|
  expect(@server.events.length).to equal(count)
end

Then("the monitor process exited with an error") do
  status = PROCESSES[-1][:thread].value
  expect(status.exited?).to be_truthy
end

Then("{string} was printed to stdout") do |contents|
  expect(PROCESSES[-1][:stdout].read).to include contents
end

Then("{string} was printed to stderr") do |contents|
  expect(PROCESSES[-1][:stderr].read).to include contents
end

Then('the following messages were printed to stderr:') do |table|
  buffer = PROCESSES[-1][:stderr].read
  table.raw.each do |message|
    expect(buffer).to include message[0]
  end
end

Then(/^I receive an error event matching (.*)$/) do |filename|
  event = @server.events.last
  expect(event).not_to be_nil
  actual = JSON.parse(event.body)
  expect(actual["events"].length).to eq(1)

  # Remove variable components of report
  expect(actual["events"][0]["device"]["hostname"]).not_to be_nil
  expect(actual["events"][0]["device"]["osName"]).not_to be_nil
  expect(actual["notifier"]["version"]).not_to be_nil
  actual["events"][0]["device"].delete("hostname")
  actual["events"][0]["device"].delete("osName")
  actual["notifier"].delete("version")

  # Separate stacktrace for more complex testing
  actual_stack = actual["events"][0]["exceptions"][0].delete("stacktrace")
  expect(actual_stack).not_to be_nil
  # Separate message - some have gotten more detailed in newer versions
  actual_message = actual["events"][0]["exceptions"][0].delete("message")
  expect(actual_message).not_to be_nil

  # Test against fixture file
  File.open(File.join('features/fixtures', filename), 'r') do |f|
    payload = f.read.strip.gsub("[[GO_VERSION]]", GO_VERSION)
    expected = JSON.parse(payload)
    expected_message = expected["events"][0]["exceptions"][0].delete("message")
    expect(actual_message).to start_with(expected_message)
    expected["notifier"].delete("version")
    expected_stack = expected["events"][0]["exceptions"][0].delete("stacktrace")
    expect(expected_stack).not_to be_nil
    expect(actual).to eq(expected)

    # Validate in-project components of the stacktrace
    validate_stacktrace(actual_stack, expected_stack)
  end
end

# Validate in-project components of the stacktrace
def validate_stacktrace actual_stack, expected_stack
  found = 0 # counts matching frames and ensures ordering is correct
  expected_len = expected_stack.length
  actual_stack.each do |frame|
    if found < expected_len and frame["inProject"] and
        frame["file"] == expected_stack[found]["file"] and
        frame["method"] == expected_stack[found]["method"] and
        frame["lineNumber"] == expected_stack[found]["lineNumber"].to_i
      found = found + 1
    elsif found >= expected_len and frame["inProject"]
      found = found + 1 # detect excess frames without false negatives
    end
  end
  expect(found).to eq(expected_len), "expected #{expected_len} matching frames but found #{found}. stacktrace:\n#{actual_stack}\nexpected these in-project frames:\n#{expected_stack}"
end

Then('the payload contains the following in-project stack frames:') do |table|
  event = @server.events.last
  expect(event).not_to be_nil
  actual = JSON.parse(event.body)
  expect(actual["events"].length).to eq(1)
  validate_stacktrace(actual["events"][0]["exceptions"][0]["stacktrace"], table.hashes)
end
