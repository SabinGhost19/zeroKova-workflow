require 'sinatra/base'
require 'json'

class HealthServer < Sinatra::Base
  get '/health' do
    content_type :json
    { status: 'healthy', service: 'notification-service' }.to_json
  end
end
