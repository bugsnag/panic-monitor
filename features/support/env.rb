require 'open3'

BUILD_DIR = File.join(Dir.pwd, "build")
PROCESSES = []
VERBOSE = ENV['VERBOSE'] || ARGV.include?('--verbose')
GO_VERSION =`go version`.split[2]

Before do
  FileUtils.mkdir_p BUILD_DIR
  PROCESSES.clear
  @server = Webserver.new
  @server.start
  @env = {"BUGSNAG_ENDPOINT" => @server.address }
end

After do
  PROCESSES.each do |p|
    begin
      puts p[:stderr].read if VERBOSE
      Process.kill 'KILL', p[:thread][:pid]
    rescue
    end
  end
  @server.stop
end

at_exit do
  # FileUtils.rm_r BUILD_DIR
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
