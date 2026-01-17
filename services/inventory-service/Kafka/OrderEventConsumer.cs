using System.Text.Json;
using Confluent.Kafka;
using InventoryService.Data;
using Microsoft.EntityFrameworkCore;

namespace InventoryService.Kafka;

public class OrderEventConsumer : BackgroundService
{
    private readonly ILogger<OrderEventConsumer> _logger;
    private readonly IServiceProvider _serviceProvider;
    private readonly IConfiguration _configuration;

    public OrderEventConsumer(
        ILogger<OrderEventConsumer> logger,
        IServiceProvider serviceProvider,
        IConfiguration configuration)
    {
        _logger = logger;
        _serviceProvider = serviceProvider;
        _configuration = configuration;
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        var bootstrapServers = Environment.GetEnvironmentVariable("KAFKA_BOOTSTRAP_SERVERS")
                               ?? _configuration["Kafka:BootstrapServers"]
                               ?? "kafka:9092";
        var groupId = _configuration["Kafka:GroupId"] ?? "inventory-service";
        var topic = _configuration["Kafka:Topic"] ?? "order-events";

        var config = new ConsumerConfig
        {
            BootstrapServers = bootstrapServers,
            GroupId = groupId,
            AutoOffsetReset = AutoOffsetReset.Earliest,
            EnableAutoCommit = true
        };

        _logger.LogInformation("Starting Kafka consumer. Bootstrap: {Bootstrap}, Topic: {Topic}",
            bootstrapServers, topic);

        // Wait for Kafka to be available
        await Task.Delay(10000, stoppingToken);

        try
        {
            using var consumer = new ConsumerBuilder<string, string>(config).Build();
            consumer.Subscribe(topic);

            _logger.LogInformation("Subscribed to topic: {Topic}", topic);

            while (!stoppingToken.IsCancellationRequested)
            {
                try
                {
                    var consumeResult = consumer.Consume(TimeSpan.FromSeconds(5));
                    if (consumeResult == null) continue;

                    _logger.LogInformation("Received message: {Key}", consumeResult.Message.Key);

                    await ProcessOrderEvent(consumeResult.Message.Value, stoppingToken);
                }
                catch (ConsumeException ex)
                {
                    _logger.LogError(ex, "Kafka consume error");
                }
            }

            consumer.Close();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Kafka consumer error");
        }
    }

    private async Task ProcessOrderEvent(string messageValue, CancellationToken cancellationToken)
    {
        try
        {
            var orderEvent = JsonSerializer.Deserialize<OrderEvent>(messageValue,
                new JsonSerializerOptions { PropertyNameCaseInsensitive = true });

            if (orderEvent == null)
            {
                _logger.LogWarning("Could not deserialize order event");
                return;
            }

            _logger.LogInformation("Processing order event: {EventType} for order {OrderId}",
                orderEvent.EventType, orderEvent.OrderId);

            if (orderEvent.EventType == "ORDER_CREATED" && orderEvent.Items != null)
            {
                using var scope = _serviceProvider.CreateScope();
                var dbContext = scope.ServiceProvider.GetRequiredService<InventoryDbContext>();

                foreach (var item in orderEvent.Items)
                {
                    if (Guid.TryParse(item.ProductId, out var productId))
                    {
                        var product = await dbContext.Products.FindAsync(new object[] { productId }, cancellationToken);
                        if (product != null)
                        {
                            product.Reserved += item.Quantity;
                            product.UpdatedAt = DateTime.UtcNow;
                            _logger.LogInformation("Reserved {Quantity} of product {ProductId}",
                                item.Quantity, productId);
                        }
                    }
                }

                await dbContext.SaveChangesAsync(cancellationToken);
            }
            else if (orderEvent.EventType == "ORDER_STATUS_UPDATED" && orderEvent.Status == "CONFIRMED")
            {
                // When order is confirmed, move from reserved to actual deduction
                using var scope = _serviceProvider.CreateScope();
                var dbContext = scope.ServiceProvider.GetRequiredService<InventoryDbContext>();

                if (orderEvent.Items != null)
                {
                    foreach (var item in orderEvent.Items)
                    {
                        if (Guid.TryParse(item.ProductId, out var productId))
                        {
                            var product = await dbContext.Products.FindAsync(new object[] { productId }, cancellationToken);
                            if (product != null)
                            {
                                product.Reserved -= item.Quantity;
                                product.Quantity -= item.Quantity;
                                product.UpdatedAt = DateTime.UtcNow;
                                _logger.LogInformation("Confirmed {Quantity} of product {ProductId}",
                                    item.Quantity, productId);
                            }
                        }
                    }

                    await dbContext.SaveChangesAsync(cancellationToken);
                }
            }
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error processing order event");
        }
    }

    private class OrderEvent
    {
        public string EventType { get; set; } = "";
        public string OrderId { get; set; } = "";
        public string CustomerName { get; set; } = "";
        public string Status { get; set; } = "";
        public decimal TotalAmount { get; set; }
        public List<OrderEventItem>? Items { get; set; }
    }

    private class OrderEventItem
    {
        public string ProductId { get; set; } = "";
        public string ProductName { get; set; } = "";
        public int Quantity { get; set; }
    }
}
