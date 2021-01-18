require 'open3'

BUILD_DIR = File.join(Dir.pwd, "build")
FIXTURE_DIR = File.expand_path(File.join(File.dirname(__FILE__), '..', 'fixtures'))
PROCESSES = []
VERBOSE = ENV['VERBOSE'] || ARGV.include?('--verbose')
GO_VERSION =`go version`.split[2]

FileUtils.mkdir_p BUILD_DIR
# Build executables for the tests
Dir.chdir(BUILD_DIR) do
  `go build ..`
  raise "Failed to build monitor" unless File.exists? "panic-monitor"
  `go build ../features/fixtures/app`
  raise "Failed to build sample app" unless File.exists? "app"
end

Before do
  PROCESSES.clear
  @server = Webserver.new
  @server.start
  @env = {"BUGSNAG_NOTIFY_ENDPOINT" => @server.address }
end

After do
  PROCESSES.each do |p|
    begin
      Kernel.puts p[:stderr].read if VERBOSE
      Process.kill 'KILL', p[:thread][:pid]
    rescue
    end
  end
  @server.stop
end

at_exit do
  FileUtils.rm_r BUILD_DIR
end

def start_process args
  stdin, stdout, stderr, thread = Open3.popen3(@env, *args)
  PROCESSES << {
    thread: thread,
    stdout: stdout,
    stderr: stderr,
    stdin: stdin
  }
end

def add_to_environment key, value
  @env[key] = value
end

def read_key_path hash, key_path
  value = hash
  key_path.split('.').each do |key|
    if key =~ /^(\d+)$/
      key = key.to_i
      if value.length > key
        value = value[key.to_i]
      else
        return nil
      end
    else
      if value.keys.include? key
        value = value[key]
      else
        return nil
      end
    end
  end
  value
end
