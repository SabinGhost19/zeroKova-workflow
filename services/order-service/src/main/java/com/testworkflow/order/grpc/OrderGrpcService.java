package com.testworkflow.order.grpc;

import com.testworkflow.order.model.Order;
import com.testworkflow.order.model.OrderStatus;
import com.testworkflow.order.service.OrderService;
import com.testworkflow.proto.*;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

@GrpcService
@RequiredArgsConstructor
@Slf4j
public class OrderGrpcService extends OrderServiceGrpc.OrderServiceImplBase {

    private final OrderService orderService;

    @Override
    public void createOrder(CreateOrderRequest request, StreamObserver<OrderResponse> responseObserver) {
        try {
            List<OrderService.OrderItemDto> items = request.getItemsList().stream()
                    .map(item -> new OrderService.OrderItemDto(
                            item.getProductId(),
                            item.getProductName(),
                            item.getQuantity(),
                            BigDecimal.valueOf(item.getPrice())
                    ))
                    .toList();

            Order order = orderService.createOrder(request.getCustomerName(), items);

            OrderResponse response = OrderResponse.newBuilder()
                    .setSuccess(true)
                    .setMessage("Order created successfully")
                    .setOrder(toProtoOrder(order))
                    .build();

            responseObserver.onNext(response);
            responseObserver.onCompleted();
        } catch (Exception e) {
            log.error("Error creating order", e);
            OrderResponse response = OrderResponse.newBuilder()
                    .setSuccess(false)
                    .setMessage("Failed to create order: " + e.getMessage())
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
        }
    }

    @Override
    public void getOrder(GetOrderRequest request, StreamObserver<OrderResponse> responseObserver) {
        try {
            UUID orderId = UUID.fromString(request.getOrderId());
            orderService.getOrder(orderId)
                    .ifPresentOrElse(
                            order -> {
                                OrderResponse response = OrderResponse.newBuilder()
                                        .setSuccess(true)
                                        .setMessage("Order found")
                                        .setOrder(toProtoOrder(order))
                                        .build();
                                responseObserver.onNext(response);
                            },
                            () -> {
                                OrderResponse response = OrderResponse.newBuilder()
                                        .setSuccess(false)
                                        .setMessage("Order not found")
                                        .build();
                                responseObserver.onNext(response);
                            }
                    );
            responseObserver.onCompleted();
        } catch (Exception e) {
            log.error("Error getting order", e);
            OrderResponse response = OrderResponse.newBuilder()
                    .setSuccess(false)
                    .setMessage("Failed to get order: " + e.getMessage())
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
        }
    }

    @Override
    public void listOrders(ListOrdersRequest request, StreamObserver<ListOrdersResponse> responseObserver) {
        try {
            int page = request.getOffset() / Math.max(request.getLimit(), 1);
            int size = request.getLimit() > 0 ? request.getLimit() : 10;

            var ordersPage = orderService.listOrders(page, size);

            ListOrdersResponse.Builder responseBuilder = ListOrdersResponse.newBuilder()
                    .setTotal((int) ordersPage.getTotalElements());

            ordersPage.getContent().forEach(order ->
                    responseBuilder.addOrders(toProtoOrder(order)));

            responseObserver.onNext(responseBuilder.build());
            responseObserver.onCompleted();
        } catch (Exception e) {
            log.error("Error listing orders", e);
            responseObserver.onNext(ListOrdersResponse.newBuilder().build());
            responseObserver.onCompleted();
        }
    }

    @Override
    public void updateOrderStatus(UpdateOrderStatusRequest request, StreamObserver<OrderResponse> responseObserver) {
        try {
            UUID orderId = UUID.fromString(request.getOrderId());
            OrderStatus status = OrderStatus.valueOf(request.getStatus().toUpperCase());

            orderService.updateOrderStatus(orderId, status)
                    .ifPresentOrElse(
                            order -> {
                                OrderResponse response = OrderResponse.newBuilder()
                                        .setSuccess(true)
                                        .setMessage("Order status updated")
                                        .setOrder(toProtoOrder(order))
                                        .build();
                                responseObserver.onNext(response);
                            },
                            () -> {
                                OrderResponse response = OrderResponse.newBuilder()
                                        .setSuccess(false)
                                        .setMessage("Order not found")
                                        .build();
                                responseObserver.onNext(response);
                            }
                    );
            responseObserver.onCompleted();
        } catch (Exception e) {
            log.error("Error updating order status", e);
            OrderResponse response = OrderResponse.newBuilder()
                    .setSuccess(false)
                    .setMessage("Failed to update order status: " + e.getMessage())
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
        }
    }

    private com.testworkflow.proto.Order toProtoOrder(Order order) {
        com.testworkflow.proto.Order.Builder builder = com.testworkflow.proto.Order.newBuilder()
                .setId(order.getId().toString())
                .setCustomerName(order.getCustomerName())
                .setTotalAmount(order.getTotalAmount().doubleValue())
                .setStatus(order.getStatus().name())
                .setCreatedAt(order.getCreatedAt().toString())
                .setUpdatedAt(order.getUpdatedAt().toString());

        order.getItems().forEach(item ->
                builder.addItems(OrderItem.newBuilder()
                        .setProductId(item.getProductId())
                        .setProductName(item.getProductName())
                        .setQuantity(item.getQuantity())
                        .setPrice(item.getPrice().doubleValue())
                        .build()));

        return builder.build();
    }
}
