package com.testworkflow.proto;

import java.util.ArrayList;
import java.util.List;

public final class Order extends com.google.protobuf.GeneratedMessageV3 {
    private String id_ = "";
    private String customerName_ = "";
    private List<OrderItem> items_ = new ArrayList<>();
    private double totalAmount_;
    private String status_ = "";
    private String createdAt_ = "";
    private String updatedAt_ = "";

    private Order() {}

    private Order(Builder builder) {
        this.id_ = builder.id_;
        this.customerName_ = builder.customerName_;
        this.items_ = new ArrayList<>(builder.items_);
        this.totalAmount_ = builder.totalAmount_;
        this.status_ = builder.status_;
        this.createdAt_ = builder.createdAt_;
        this.updatedAt_ = builder.updatedAt_;
    }

    public String getId() { return id_; }
    public String getCustomerName() { return customerName_; }
    public List<OrderItem> getItemsList() { return items_; }
    public double getTotalAmount() { return totalAmount_; }
    public String getStatus() { return status_; }
    public String getCreatedAt() { return createdAt_; }
    public String getUpdatedAt() { return updatedAt_; }

    public static Builder newBuilder() { return new Builder(); }

    public static final class Builder {
        private String id_ = "";
        private String customerName_ = "";
        private List<OrderItem> items_ = new ArrayList<>();
        private double totalAmount_;
        private String status_ = "";
        private String createdAt_ = "";
        private String updatedAt_ = "";

        public Builder setId(String value) { id_ = value; return this; }
        public Builder setCustomerName(String value) { customerName_ = value; return this; }
        public Builder addItems(OrderItem value) { items_.add(value); return this; }
        public Builder setTotalAmount(double value) { totalAmount_ = value; return this; }
        public Builder setStatus(String value) { status_ = value; return this; }
        public Builder setCreatedAt(String value) { createdAt_ = value; return this; }
        public Builder setUpdatedAt(String value) { updatedAt_ = value; return this; }
        public Order build() { return new Order(this); }
    }
}
