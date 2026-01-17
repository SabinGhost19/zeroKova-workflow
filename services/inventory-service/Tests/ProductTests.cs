using InventoryService.Models;
using Xunit;

namespace InventoryService.Tests;

public class ProductTests
{
    [Fact]
    public void Product_ShouldHaveCorrectDefaults()
    {
        var product = new Product
        {
            Name = "Test Product",
            Price = 29.99m,
            Quantity = 100
        };

        Assert.NotEqual(Guid.Empty, product.Id);
        Assert.Equal("Test Product", product.Name);
        Assert.Equal(29.99m, product.Price);
        Assert.Equal(100, product.Quantity);
        Assert.Equal(0, product.Reserved);
    }

    [Fact]
    public void Product_AvailableQuantity_ShouldBeCalculatedCorrectly()
    {
        var product = new Product
        {
            Name = "Test Product",
            Price = 29.99m,
            Quantity = 100,
            Reserved = 20
        };

        Assert.Equal(80, product.AvailableQuantity);
    }

    [Fact]
    public void Product_AvailableQuantity_WhenNoReservations()
    {
        var product = new Product
        {
            Name = "Test Product",
            Price = 29.99m,
            Quantity = 50,
            Reserved = 0
        };

        Assert.Equal(50, product.AvailableQuantity);
    }
}
