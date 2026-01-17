require_relative 'database'
require_relative '../protos/notification_services_pb'

class NotificationServiceImpl < Notification::NotificationService::Service
  def send_order_notification(request, _call)
    puts "[NotificationService] Sending order notification for order #{request.order_id}"

    message = "Order #{request.order_id} for #{request.customer_name} - Status: #{request.status}, Total: $#{request.total_amount}"

    notification_id = save_notification('ORDER', message)

    puts "[NotificationService] Notification saved with ID: #{notification_id}"

    Common::StatusResponse.new(
      success: true,
      message: "Order notification sent successfully"
    )
  rescue => e
    puts "[NotificationService] Error: #{e.message}"
    Common::StatusResponse.new(
      success: false,
      message: "Failed to send notification: #{e.message}"
    )
  end

  def send_stock_alert(request, _call)
    puts "[NotificationService] Sending stock alert for product #{request.product_id}"

    message = "Stock Alert [#{request.alert_type}]: #{request.product_name} - Current quantity: #{request.current_quantity}"

    notification_id = save_notification('STOCK_ALERT', message)

    puts "[NotificationService] Stock alert saved with ID: #{notification_id}"

    Common::StatusResponse.new(
      success: true,
      message: "Stock alert sent successfully"
    )
  rescue => e
    puts "[NotificationService] Error: #{e.message}"
    Common::StatusResponse.new(
      success: false,
      message: "Failed to send stock alert: #{e.message}"
    )
  end

  def get_notifications(request, _call)
    puts "[NotificationService] Fetching notifications"

    limit = request.limit > 0 ? request.limit : 10
    offset = request.offset > 0 ? request.offset : 0

    total = Database.notifications.count

    notifications = Database.notifications
      .order(Sequel.desc(:created_at))
      .limit(limit)
      .offset(offset)
      .all
      .map do |row|
        Notification::Notification.new(
          id: row[:id].to_s,
          type: row[:type],
          message: row[:message],
          created_at: row[:created_at].iso8601,
          sent: row[:sent]
        )
      end

    puts "[NotificationService] Found #{notifications.length} notifications"

    Notification::NotificationsResponse.new(
      notifications: notifications,
      total: total
    )
  rescue => e
    puts "[NotificationService] Error: #{e.message}"
    Notification::NotificationsResponse.new(
      notifications: [],
      total: 0
    )
  end

  private

  def save_notification(type, message)
    id = SecureRandom.uuid
    Database.notifications.insert(
      id: id,
      type: type,
      message: message,
      sent: true,
      created_at: Time.now
    )
    id
  end
end
