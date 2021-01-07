Given("I build the executable") do
  Dir.chdir(BUILD_DIR) do
    `go build ..`
    expect(File.exists? "panic-monitor").to be_truthy
  end
end

Given("I build a sample app") do
  Dir.chdir(BUILD_DIR) do
    `go build ../features/fixtures/app`
    expect(File.exists? "app").to be_truthy
  end
end

When(/^I crash the app using (.*)$/) do |testcase|
  add_to_environment("TESTCASE", testcase)
  Dir.chdir(BUILD_DIR) do
    start_process(["./panic-monitor", "./app"])
  end
end

When('I run the monitor with arguments {string}') do |args|
  Dir.chdir(BUILD_DIR) do
    start_process(["./panic-monitor"] + args.split(' '))
  end
end

Given('I set the API key to {string}') do |key|
  add_to_environment("BUGSNAG_API_KEY", key)
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

Then(/^I receive an error event matching (.*)$/) do |filename|
  event = @server.events.last
  expect(event).not_to be_nil
  actual = JSON.parse(event.body)
  expect(actual["events"].length).to eq(1)

  # Remove variable components of report
  expect(actual["events"][0]["device"]["hostname"]).not_to be_nil
  expect(actual["events"][0]["device"]["osName"]).not_to be_nil
  actual["events"][0]["device"].delete("hostname")
  actual["events"][0]["device"].delete("osName")

  # Separate stacktrace for more complex testing
  actual_stack = actual["events"][0]["exceptions"][0].delete("stacktrace")
  expect(actual_stack).not_to be_nil

  # Test against fixture file
  File.open(File.join('features/fixtures', filename), 'r') do |f|
    payload = f.read.strip.gsub("[[GO_VERSION]]", GO_VERSION)
    expected = JSON.parse(payload)
    expected_stack = expected["events"][0]["exceptions"][0].delete("stacktrace")
    expect(expected_stack).not_to be_nil
    expect(actual).to eq(expected)

    # Validate in-project components of the stacktrace
    found = 0 # counts matching frames and ensures ordering is correct
    expected_len = expected_stack.length
    actual_stack.each do |frame|
      if found < expected_len and frame["inProject"] and
          frame["file"] == expected_stack[found]["file"] and
          frame["method"] == expected_stack[found]["method"]
        found = found + 1
      elsif found >= expected_len and frame["inProject"]
        found = found + 1 # detect excess frames without false negatives
      end
    end
    expect(found).to eq(expected_len), "expected #{expected_len} matching frames but found #{found}. stacktrace:\n#{actual_stack}"
  end
end
