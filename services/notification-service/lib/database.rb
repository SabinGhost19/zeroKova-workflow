require 'sequel'

module Database
  def self.connection
    @connection ||= begin
      host = ENV['DB_HOST'] || 'localhost'
      port = ENV['DB_PORT'] || '5432'
      database = ENV['DB_NAME'] || 'testworkflow'
      user = ENV['DB_USER'] || 'postgres'
      password = ENV['DB_PASSWORD'] || 'postgres'

      db = Sequel.connect(
        adapter: 'postgres',
        host: host,
        port: port.to_i,
        database: database,
        user: user,
        password: password,
        search_path: ['notifications']
      )

      # Ensure schema exists
      db.run('CREATE SCHEMA IF NOT EXISTS notifications')

      # Create notifications table if not exists
      db.create_table?(:notifications) do
        uuid :id, primary_key: true, default: Sequel.lit('gen_random_uuid()')
        String :type, null: false
        String :message, text: true, null: false
        Boolean :sent, default: false
        DateTime :created_at, default: Sequel.lit('CURRENT_TIMESTAMP')
      end

      db
    end
  end

  def self.notifications
    connection[:notifications]
  end
end
