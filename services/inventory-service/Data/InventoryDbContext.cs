using InventoryService.Models;
using Microsoft.EntityFrameworkCore;

namespace InventoryService.Data;

public class InventoryDbContext : DbContext
{
    public InventoryDbContext(DbContextOptions<InventoryDbContext> options) : base(options)
    {
    }

    public DbSet<Product> Products { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("inventory");

        modelBuilder.Entity<Product>(entity =>
        {
            entity.HasIndex(e => e.Name);
            entity.Property(e => e.Price).HasPrecision(18, 2);
        });

        // Seed some test data
        modelBuilder.Entity<Product>().HasData(
            new Product
            {
                Id = Guid.Parse("11111111-1111-1111-1111-111111111111"),
                Name = "Test Product 1",
                Description = "A test product for development",
                Price = 29.99m,
                Quantity = 100,
                Reserved = 0
            },
            new Product
            {
                Id = Guid.Parse("22222222-2222-2222-2222-222222222222"),
                Name = "Test Product 2",
                Description = "Another test product",
                Price = 49.99m,
                Quantity = 50,
                Reserved = 0
            },
            new Product
            {
                Id = Guid.Parse("33333333-3333-3333-3333-333333333333"),
                Name = "Test Product 3",
                Description = "Third test product",
                Price = 99.99m,
                Quantity = 25,
                Reserved = 0
            }
        );
    }
}
