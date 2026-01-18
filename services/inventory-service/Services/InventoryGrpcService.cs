using Grpc.Core;
using InventoryService.Data;
using InventoryService.Models;
using Microsoft.EntityFrameworkCore;
using TestWorkflow.Proto;

namespace InventoryService.Services;

public class InventoryGrpcService : TestWorkflow.Proto.InventoryService.InventoryServiceBase
{
    private readonly InventoryDbContext _dbContext;
    private readonly ILogger<InventoryGrpcService> _logger;

    public InventoryGrpcService(InventoryDbContext dbContext, ILogger<InventoryGrpcService> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public override async Task<StockResponse> GetStock(GetStockRequest request, ServerCallContext context)
    {
        try
        {
            if (!Guid.TryParse(request.ProductId, out var productId))
            {
                return new StockResponse { Success = false, Message = "Invalid product ID format" };
            }

            var product = await _dbContext.Products.FindAsync(productId);
            if (product == null)
            {
                return new StockResponse { Success = false, Message = "Product not found" };
            }

            return new StockResponse
            {
                Success = true,
                Message = "Product found",
                Product = ToProtoProduct(product)
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting stock for product {ProductId}", request.ProductId);
            return new StockResponse { Success = false, Message = $"Error: {ex.Message}" };
        }
    }

    public override async Task<StockResponse> UpdateStock(UpdateStockRequest request, ServerCallContext context)
    {
        try
        {
            if (!Guid.TryParse(request.ProductId, out var productId))
            {
                return new StockResponse { Success = false, Message = "Invalid product ID format" };
            }

            var product = await _dbContext.Products.FindAsync(productId);
            if (product == null)
            {
                return new StockResponse { Success = false, Message = "Product not found" };
            }

            var newQuantity = product.Quantity + request.QuantityChange;
            if (newQuantity < 0)
            {
                return new StockResponse { Success = false, Message = "Insufficient stock" };
            }

            product.Quantity = newQuantity;
            product.UpdatedAt = DateTime.UtcNow;
            await _dbContext.SaveChangesAsync();

            _logger.LogInformation("Stock updated for product {ProductId}: {Change}", productId, request.QuantityChange);

            return new StockResponse
            {
                Success = true,
                Message = "Stock updated successfully",
                Product = ToProtoProduct(product)
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error updating stock for product {ProductId}", request.ProductId);
            return new StockResponse { Success = false, Message = $"Error: {ex.Message}" };
        }
    }

    public override async Task<StatusResponse> ReserveStock(ReserveStockRequest request, ServerCallContext context)
    {
        try
        {
            if (!Guid.TryParse(request.ProductId, out var productId))
            {
                return new StatusResponse { Success = false, Message = "Invalid product ID format" };
            }

            var product = await _dbContext.Products.FindAsync(productId);
            if (product == null)
            {
                return new StatusResponse { Success = false, Message = "Product not found" };
            }

            if (product.AvailableQuantity < request.Quantity)
            {
                return new StatusResponse { Success = false, Message = "Insufficient available stock" };
            }

            product.Reserved += request.Quantity;
            product.UpdatedAt = DateTime.UtcNow;
            await _dbContext.SaveChangesAsync();

            _logger.LogInformation("Reserved {Quantity} units of product {ProductId} for order {OrderId}",
                request.Quantity, productId, request.OrderId);

            return new StatusResponse { Success = true, Message = "Stock reserved successfully" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error reserving stock for product {ProductId}", request.ProductId);
            return new StatusResponse { Success = false, Message = $"Error: {ex.Message}" };
        }
    }

    public override async Task<StatusResponse> ReleaseStock(ReleaseStockRequest request, ServerCallContext context)
    {
        try
        {
            if (!Guid.TryParse(request.ProductId, out var productId))
            {
                return new StatusResponse { Success = false, Message = "Invalid product ID format" };
            }

            var product = await _dbContext.Products.FindAsync(productId);
            if (product == null)
            {
                return new StatusResponse { Success = false, Message = "Product not found" };
            }

            product.Reserved = Math.Max(0, product.Reserved - request.Quantity);
            product.UpdatedAt = DateTime.UtcNow;
            await _dbContext.SaveChangesAsync();

            _logger.LogInformation("Released {Quantity} units of product {ProductId} for order {OrderId}",
                request.Quantity, productId, request.OrderId);

            return new StatusResponse { Success = true, Message = "Stock released successfully" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error releasing stock for product {ProductId}", request.ProductId);
            return new StatusResponse { Success = false, Message = $"Error: {ex.Message}" };
        }
    }

    public override async Task<ListProductsResponse> ListProducts(ListProductsRequest request, ServerCallContext context)
    {
        try
        {
            var limit = request.Limit > 0 ? request.Limit : 10;
            var offset = request.Offset > 0 ? request.Offset : 0;

            var totalCount = await _dbContext.Products.CountAsync();
            var products = await _dbContext.Products
                .OrderBy(p => p.Name)
                .Skip(offset)
                .Take(limit)
                .ToListAsync();

            var response = new ListProductsResponse { Total = totalCount };
            response.Products.AddRange(products.Select(ToProtoProduct));

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error listing products");
            return new ListProductsResponse();
        }
    }

    public override async Task<ProductResponse> AddProduct(AddProductRequest request, ServerCallContext context)
    {
        try
        {
            var product = new Models.Product
            {
                Name = request.Name,
                Description = request.Description,
                Price = (decimal)request.Price,
                Quantity = request.Quantity,
                Reserved = 0
            };

            _dbContext.Products.Add(product);
            await _dbContext.SaveChangesAsync();

            _logger.LogInformation("Product added: {ProductId} - {ProductName}", product.Id, product.Name);

            return new ProductResponse
            {
                Success = true,
                Message = "Product added successfully",
                Product = ToProtoProduct(product)
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error adding product");
            return new ProductResponse { Success = false, Message = $"Error: {ex.Message}" };
        }
    }

    private static TestWorkflow.Proto.Product ToProtoProduct(Models.Product product)
    {
        return new TestWorkflow.Proto.Product
        {
            Id = product.Id.ToString(),
            Name = product.Name,
            Description = product.Description ?? "",
            Price = (double)product.Price,
            Quantity = product.Quantity,
            Reserved = product.Reserved,
            CreatedAt = product.CreatedAt.ToString("O"),
            UpdatedAt = product.UpdatedAt.ToString("O")
        };
    }
}
