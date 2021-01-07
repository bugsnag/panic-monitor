require 'webrick'

class Event
  attr_reader :body, :headers

  def initialize req
    @body = req.body
    @headers = req.each {|k,v| [k,v]}
  end
end

class Webserver
  attr_reader :events, :address

  def initialize
    @events = []
    dev_null = WEBrick::Log::new("/dev/null", 7)
    @server = WEBrick::HTTPServer.new Port: 8000, AccessLog: dev_null, :Logger => dev_null
    @server.mount_proc '/events' do |req, res|
      @events.append(Event.new(req))
      res.status = 201
      res.content_length = 0
    end
    @address = "http://localhost:8000/events"
  end

  def start
    @server_thread = Thread.new { @server.start }
  end

  def stop
    @server.shutdown
    @server_thread.join
  end

  private

  attr_reader :server, :server_thread
end
