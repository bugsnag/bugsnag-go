#!/usr/bin/env ruby

require 'webrick'

root = File.expand_path File.dirname(__FILE__)
server = WEBrick::HTTPServer.new :Port => 8991, :DocumentRoot => root
server.start
