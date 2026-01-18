// Inventory Service - Manages product inventory
// Version: 1.0.0
using InventoryService.Data;
using InventoryService.Kafka;
using InventoryService.Services;
using Microsoft.EntityFrameworkCore;

var builder = WebApplication.CreateBuilder(args);

// Add services
builder.Services.AddGrpc();

// Database
var connectionString = $"Host={Environment.GetEnvironmentVariable("DB_HOST") ?? "localhost"};" +
                       $"Port={Environment.GetEnvironmentVariable("DB_PORT") ?? "5432"};" +
                       $"Database={Environment.GetEnvironmentVariable("DB_NAME") ?? "testworkflow"};" +
                       $"Username={Environment.GetEnvironmentVariable("DB_USER") ?? "postgres"};" +
                       $"Password={Environment.GetEnvironmentVariable("DB_PASSWORD") ?? "postgres"};" +
                       "Search Path=inventory";

builder.Services.AddDbContext<InventoryDbContext>(options =>
    options.UseNpgsql(connectionString));

// Kafka consumer
builder.Services.AddHostedService<OrderEventConsumer>();

var app = builder.Build();

// Ensure database is created
using (var scope = app.Services.CreateScope())
{
    var db = scope.ServiceProvider.GetRequiredService<InventoryDbContext>();
    db.Database.EnsureCreated();
}

// Map gRPC service
app.MapGrpcService<InventoryGrpcService>();

// Health check endpoint
app.MapGet("/health", () => Results.Ok(new { status = "healthy", service = "inventory-service" }));

app.Run();
