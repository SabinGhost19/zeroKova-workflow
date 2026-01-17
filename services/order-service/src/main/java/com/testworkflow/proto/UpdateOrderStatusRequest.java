package com.testworkflow.proto;

public final class UpdateOrderStatusRequest extends com.google.protobuf.GeneratedMessageV3 {
    private String orderId_ = "";
    private String status_ = "";

    private UpdateOrderStatusRequest() {}

    public String getOrderId() { return orderId_; }
    public String getStatus() { return status_; }

    public static Builder newBuilder() { return new Builder(); }

    public static final class Builder {
        private String orderId_ = "";
        private String status_ = "";

        public Builder setOrderId(String value) { orderId_ = value; return this; }
        public Builder setStatus(String value) { status_ = value; return this; }
        public UpdateOrderStatusRequest build() {
            UpdateOrderStatusRequest req = new UpdateOrderStatusRequest();
            req.orderId_ = orderId_;
            req.status_ = status_;
            return req;
        }
    }
}
