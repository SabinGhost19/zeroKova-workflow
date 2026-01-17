package com.testworkflow.proto;

import java.util.ArrayList;
import java.util.List;

public final class CreateOrderRequest extends com.google.protobuf.GeneratedMessageV3 {
    private String customerName_ = "";
    private List<OrderItem> items_ = new ArrayList<>();

    private CreateOrderRequest() {}

    public String getCustomerName() { return customerName_; }
    public List<OrderItem> getItemsList() { return items_; }

    public static Builder newBuilder() { return new Builder(); }

    public static final class Builder {
        private String customerName_ = "";
        private List<OrderItem> items_ = new ArrayList<>();

        public Builder setCustomerName(String value) { customerName_ = value; return this; }
        public Builder addItems(OrderItem value) { items_.add(value); return this; }
        public CreateOrderRequest build() {
            CreateOrderRequest req = new CreateOrderRequest();
            req.customerName_ = customerName_;
            req.items_ = new ArrayList<>(items_);
            return req;
        }
    }
}
