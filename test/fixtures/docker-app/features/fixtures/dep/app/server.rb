require 'rack'
require 'rack/request'
require 'rack/response'

class DepServer
  def call(env)
    req = Rack::Request.new(env)
    res = Rack::Response.new
    res.status = 200
    res.finish
  end
end

Rack::Server.start(app: DepServer.new)