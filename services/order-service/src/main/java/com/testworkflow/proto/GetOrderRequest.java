package com.testworkflow.proto;

public final class GetOrderRequest extends com.google.protobuf.GeneratedMessageV3 {
    private String orderId_ = "";

    private GetOrderRequest() {}

    public String getOrderId() { return orderId_; }

    public static Builder newBuilder() { return new Builder(); }

    public static final class Builder {
        private String orderId_ = "";

        public Builder setOrderId(String value) { orderId_ = value; return this; }
        public GetOrderRequest build() {
            GetOrderRequest req = new GetOrderRequest();
            req.orderId_ = orderId_;
            return req;
        }
    }
}
