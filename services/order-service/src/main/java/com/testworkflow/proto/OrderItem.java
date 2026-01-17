package com.testworkflow.proto;

public final class OrderItem extends com.google.protobuf.GeneratedMessageV3 {
    private String productId_ = "";
    private String productName_ = "";
    private int quantity_;
    private double price_;

    private OrderItem() {}

    private OrderItem(Builder builder) {
        this.productId_ = builder.productId_;
        this.productName_ = builder.productName_;
        this.quantity_ = builder.quantity_;
        this.price_ = builder.price_;
    }

    public String getProductId() { return productId_; }
    public String getProductName() { return productName_; }
    public int getQuantity() { return quantity_; }
    public double getPrice() { return price_; }

    public static Builder newBuilder() { return new Builder(); }

    public static final class Builder {
        private String productId_ = "";
        private String productName_ = "";
        private int quantity_;
        private double price_;

        public Builder setProductId(String value) { productId_ = value; return this; }
        public Builder setProductName(String value) { productName_ = value; return this; }
        public Builder setQuantity(int value) { quantity_ = value; return this; }
        public Builder setPrice(double value) { price_ = value; return this; }
        public OrderItem build() { return new OrderItem(this); }
    }
}
