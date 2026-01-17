package com.testworkflow.proto;

import java.util.ArrayList;
import java.util.List;

public final class ListOrdersResponse extends com.google.protobuf.GeneratedMessageV3 {
    private List<Order> orders_ = new ArrayList<>();
    private int total_;

    private ListOrdersResponse() {}

    public List<Order> getOrdersList() { return orders_; }
    public int getTotal() { return total_; }

    public static Builder newBuilder() { return new Builder(); }

    public static final class Builder {
        private List<Order> orders_ = new ArrayList<>();
        private int total_;

        public Builder addOrders(Order value) { orders_.add(value); return this; }
        public Builder setTotal(int value) { total_ = value; return this; }
        public ListOrdersResponse build() {
            ListOrdersResponse resp = new ListOrdersResponse();
            resp.orders_ = new ArrayList<>(orders_);
            resp.total_ = total_;
            return resp;
        }
    }
}
