package com.testworkflow.order.service;

import com.testworkflow.order.kafka.OrderEventProducer;
import com.testworkflow.order.model.Order;
import com.testworkflow.order.model.OrderItem;
import com.testworkflow.order.model.OrderStatus;
import com.testworkflow.order.repository.OrderRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Slf4j
public class OrderService {

    private final OrderRepository orderRepository;
    private final OrderEventProducer orderEventProducer;

    @Transactional
    public Order createOrder(String customerName, List<OrderItemDto> items) {
        log.info("Creating order for customer: {}", customerName);

        Order order = Order.builder()
                .customerName(customerName)
                .status(OrderStatus.PENDING)
                .totalAmount(BigDecimal.ZERO)
                .build();

        for (OrderItemDto itemDto : items) {
            OrderItem item = OrderItem.builder()
                    .productId(itemDto.productId())
                    .productName(itemDto.productName())
                    .quantity(itemDto.quantity())
                    .price(itemDto.price())
                    .build();
            order.addItem(item);
        }

        order.calculateTotal();
        Order savedOrder = orderRepository.save(order);

        // Send event to Kafka for inventory update
        orderEventProducer.sendOrderCreatedEvent(savedOrder);

        log.info("Order created with ID: {}", savedOrder.getId());
        return savedOrder;
    }

    public Optional<Order> getOrder(UUID orderId) {
        return orderRepository.findById(orderId);
    }

    public Page<Order> listOrders(int page, int size) {
        return orderRepository.findAllByOrderByCreatedAtDesc(PageRequest.of(page, size));
    }

    @Transactional
    public Optional<Order> updateOrderStatus(UUID orderId, OrderStatus status) {
        return orderRepository.findById(orderId)
                .map(order -> {
                    order.setStatus(status);
                    Order updatedOrder = orderRepository.save(order);

                    // Send status update event
                    orderEventProducer.sendOrderStatusUpdatedEvent(updatedOrder);

                    log.info("Order {} status updated to {}", orderId, status);
                    return updatedOrder;
                });
    }

    public record OrderItemDto(
            String productId,
            String productName,
            Integer quantity,
            BigDecimal price
    ) {}
}
