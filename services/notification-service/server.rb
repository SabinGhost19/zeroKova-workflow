#!/usr/bin/env ruby

$LOAD_PATH.unshift(File.join(File.dirname(__FILE__), 'lib'))
$LOAD_PATH.unshift(File.join(File.dirname(__FILE__), 'protos'))

require 'grpc'
require 'securerandom'
require_relative 'lib/notification_service'
require_relative 'lib/database'

class GrpcServer
  def initialize(port = '50053')
    @port = port
  end

  def start
    # Initialize database connection
    puts "[NotificationService] Connecting to database..."
    Database.connection
    puts "[NotificationService] Database connected"

    # Start gRPC server
    @server = GRPC::RpcServer.new
    @server.add_http2_port("0.0.0.0:#{@port}", :this_port_is_insecure)
    @server.handle(NotificationServiceImpl.new)

    puts "[NotificationService] gRPC server starting on port #{@port}..."
    @server.run_till_terminated
  end
end

if __FILE__ == $0
  port = ENV['GRPC_PORT'] || '50053'

  # Start HTTP health check server in a thread
  Thread.new do
    require 'webrick'
    require 'json'

    http_port = ENV['HTTP_PORT'] || '8083'
    server = WEBrick::HTTPServer.new(Port: http_port.to_i)

    server.mount_proc '/health' do |req, res|
      res['Content-Type'] = 'application/json'
      res.body = { status: 'healthy', service: 'notification-service' }.to_json
    end

    puts "[NotificationService] HTTP health server starting on port #{http_port}..."
    server.start
  end

  # Give the HTTP server time to start
  sleep 1

  # Start gRPC server
  GrpcServer.new(port).start
end
