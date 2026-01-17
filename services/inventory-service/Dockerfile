# Build stage
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS builder

WORKDIR /app

# Copy csproj and restore
COPY *.csproj ./
RUN dotnet restore

# Copy everything else and build
COPY . ./
RUN dotnet publish -c Release -o out

# Final stage
FROM mcr.microsoft.com/dotnet/aspnet:8.0

WORKDIR /app

# Copy the published app
COPY --from=builder /app/out .

# Expose ports (HTTP and gRPC)
EXPOSE 8082 50052

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:8082/health || exit 1

# Run the application
ENTRYPOINT ["dotnet", "InventoryService.dll"]
