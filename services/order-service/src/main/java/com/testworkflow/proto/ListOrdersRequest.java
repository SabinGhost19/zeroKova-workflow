package com.testworkflow.proto;

public final class ListOrdersRequest extends com.google.protobuf.GeneratedMessageV3 {
    private int limit_;
    private int offset_;

    private ListOrdersRequest() {}

    public int getLimit() { return limit_; }
    public int getOffset() { return offset_; }

    public static Builder newBuilder() { return new Builder(); }

    public static final class Builder {
        private int limit_;
        private int offset_;

        public Builder setLimit(int value) { limit_ = value; return this; }
        public Builder setOffset(int value) { offset_ = value; return this; }
        public ListOrdersRequest build() {
            ListOrdersRequest req = new ListOrdersRequest();
            req.limit_ = limit_;
            req.offset_ = offset_;
            return req;
        }
    }
}
