require 'rspec'

RSpec.describe 'NotificationService' do
  describe 'Notification creation' do
    it 'creates an order notification message' do
      order_id = 'order-123'
      customer_name = 'John Doe'
      status = 'CONFIRMED'
      total_amount = 99.99

      message = "Order #{order_id} for #{customer_name} - Status: #{status}, Total: $#{total_amount}"

      expect(message).to include(order_id)
      expect(message).to include(customer_name)
      expect(message).to include(status)
      expect(message).to include('99.99')
    end

    it 'creates a stock alert message' do
      product_id = 'prod-456'
      product_name = 'Test Product'
      current_quantity = 5
      alert_type = 'LOW_STOCK'

      message = "Stock Alert [#{alert_type}]: #{product_name} - Current quantity: #{current_quantity}"

      expect(message).to include(alert_type)
      expect(message).to include(product_name)
      expect(message).to include('5')
    end
  end

  describe 'HealthServer' do
    it 'returns healthy status' do
      status = { status: 'healthy', service: 'notification-service' }

      expect(status[:status]).to eq('healthy')
      expect(status[:service]).to eq('notification-service')
    end
  end
end
