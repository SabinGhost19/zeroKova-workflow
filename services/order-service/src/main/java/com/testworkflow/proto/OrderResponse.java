package com.testworkflow.proto;

public final class OrderResponse extends com.google.protobuf.GeneratedMessageV3 {
    private boolean success_;
    private String message_ = "";
    private Order order_;

    private OrderResponse() {}

    public boolean getSuccess() { return success_; }
    public String getMessage() { return message_; }
    public Order getOrder() { return order_; }

    public static Builder newBuilder() { return new Builder(); }

    public static final class Builder {
        private boolean success_;
        private String message_ = "";
        private Order order_;

        public Builder setSuccess(boolean value) { success_ = value; return this; }
        public Builder setMessage(String value) { message_ = value; return this; }
        public Builder setOrder(Order value) { order_ = value; return this; }
        public OrderResponse build() {
            OrderResponse resp = new OrderResponse();
            resp.success_ = success_;
            resp.message_ = message_;
            resp.order_ = order_;
            return resp;
        }
    }
}
