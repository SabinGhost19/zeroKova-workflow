package com.testworkflow.order.kafka;

import com.testworkflow.order.model.Order;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

import java.util.HashMap;
import java.util.Map;

@Component
@RequiredArgsConstructor
@Slf4j
public class OrderEventProducer {

    private final KafkaTemplate<String, Object> kafkaTemplate;

    @Value("${kafka.topic.order-events}")
    private String orderEventsTopic;

    public void sendOrderCreatedEvent(Order order) {
        Map<String, Object> event = createEventPayload(order, "ORDER_CREATED");
        sendEvent(order.getId().toString(), event);
    }

    public void sendOrderStatusUpdatedEvent(Order order) {
        Map<String, Object> event = createEventPayload(order, "ORDER_STATUS_UPDATED");
        sendEvent(order.getId().toString(), event);
    }

    private Map<String, Object> createEventPayload(Order order, String eventType) {
        Map<String, Object> event = new HashMap<>();
        event.put("eventType", eventType);
        event.put("orderId", order.getId().toString());
        event.put("customerName", order.getCustomerName());
        event.put("status", order.getStatus().name());
        event.put("totalAmount", order.getTotalAmount());
        event.put("itemCount", order.getItems().size());

        // Include items for inventory processing
        event.put("items", order.getItems().stream()
                .map(item -> Map.of(
                        "productId", item.getProductId(),
                        "productName", item.getProductName(),
                        "quantity", item.getQuantity()
                ))
                .toList());

        return event;
    }

    private void sendEvent(String key, Object payload) {
        kafkaTemplate.send(orderEventsTopic, key, payload)
                .whenComplete((result, ex) -> {
                    if (ex == null) {
                        log.info("Event sent successfully to topic {}: {}", orderEventsTopic, key);
                    } else {
                        log.error("Failed to send event to topic {}: {}", orderEventsTopic, ex.getMessage());
                    }
                });
    }
}
