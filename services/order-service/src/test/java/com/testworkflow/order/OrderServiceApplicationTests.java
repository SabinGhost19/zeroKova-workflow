package com.testworkflow.order;

import com.testworkflow.order.model.Order;
import com.testworkflow.order.model.OrderItem;
import com.testworkflow.order.model.OrderStatus;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;

import java.math.BigDecimal;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest(classes = OrderServiceApplication.class)
@ActiveProfiles("test")
class OrderServiceApplicationTests {

    @Test
    void contextLoads() {
        // Verify application context loads
    }

    @Test
    void testOrderCreation() {
        Order order = Order.builder()
                .customerName("Test Customer")
                .status(OrderStatus.PENDING)
                .totalAmount(BigDecimal.ZERO)
                .build();

        assertNotNull(order);
        assertEquals("Test Customer", order.getCustomerName());
        assertEquals(OrderStatus.PENDING, order.getStatus());
    }

    @Test
    void testOrderItemAddition() {
        Order order = Order.builder()
                .customerName("Test Customer")
                .status(OrderStatus.PENDING)
                .totalAmount(BigDecimal.ZERO)
                .build();

        OrderItem item = OrderItem.builder()
                .productId("prod-1")
                .productName("Test Product")
                .quantity(2)
                .price(BigDecimal.valueOf(25.00))
                .build();

        order.addItem(item);
        order.calculateTotal();

        assertEquals(1, order.getItems().size());
        assertEquals(BigDecimal.valueOf(50.00), order.getTotalAmount());
    }
}
