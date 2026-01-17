package com.testworkflow.proto;

import io.grpc.BindableService;
import io.grpc.ServerServiceDefinition;
import io.grpc.stub.StreamObserver;

import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;

public final class OrderServiceGrpc {

    private OrderServiceGrpc() {}

    public static final String SERVICE_NAME = "order.OrderService";

    public static abstract class OrderServiceImplBase implements BindableService {

        public void createOrder(CreateOrderRequest request, StreamObserver<OrderResponse> responseObserver) {
            io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCreateOrderMethod(), responseObserver);
        }

        public void getOrder(GetOrderRequest request, StreamObserver<OrderResponse> responseObserver) {
            io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetOrderMethod(), responseObserver);
        }

        public void listOrders(ListOrdersRequest request, StreamObserver<ListOrdersResponse> responseObserver) {
            io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListOrdersMethod(), responseObserver);
        }

        public void updateOrderStatus(UpdateOrderStatusRequest request, StreamObserver<OrderResponse> responseObserver) {
            io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getUpdateOrderStatusMethod(), responseObserver);
        }

        @Override
        public ServerServiceDefinition bindService() {
            return io.grpc.ServerServiceDefinition.builder(SERVICE_NAME)
                    .addMethod(
                            getCreateOrderMethod(),
                            asyncUnaryCall((request, responseObserver) -> createOrder((CreateOrderRequest) request, (StreamObserver<OrderResponse>) responseObserver)))
                    .addMethod(
                            getGetOrderMethod(),
                            asyncUnaryCall((request, responseObserver) -> getOrder((GetOrderRequest) request, (StreamObserver<OrderResponse>) responseObserver)))
                    .addMethod(
                            getListOrdersMethod(),
                            asyncUnaryCall((request, responseObserver) -> listOrders((ListOrdersRequest) request, (StreamObserver<ListOrdersResponse>) responseObserver)))
                    .addMethod(
                            getUpdateOrderStatusMethod(),
                            asyncUnaryCall((request, responseObserver) -> updateOrderStatus((UpdateOrderStatusRequest) request, (StreamObserver<OrderResponse>) responseObserver)))
                    .build();
        }
    }

    private static io.grpc.MethodDescriptor<CreateOrderRequest, OrderResponse> getCreateOrderMethod() {
        return io.grpc.MethodDescriptor.<CreateOrderRequest, OrderResponse>newBuilder()
                .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CreateOrder"))
                .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(CreateOrderRequest.getDefaultInstance()))
                .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(OrderResponse.getDefaultInstance()))
                .build();
    }

    private static io.grpc.MethodDescriptor<GetOrderRequest, OrderResponse> getGetOrderMethod() {
        return io.grpc.MethodDescriptor.<GetOrderRequest, OrderResponse>newBuilder()
                .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetOrder"))
                .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(GetOrderRequest.getDefaultInstance()))
                .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(OrderResponse.getDefaultInstance()))
                .build();
    }

    private static io.grpc.MethodDescriptor<ListOrdersRequest, ListOrdersResponse> getListOrdersMethod() {
        return io.grpc.MethodDescriptor.<ListOrdersRequest, ListOrdersResponse>newBuilder()
                .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListOrders"))
                .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(ListOrdersRequest.getDefaultInstance()))
                .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(ListOrdersResponse.getDefaultInstance()))
                .build();
    }

    private static io.grpc.MethodDescriptor<UpdateOrderStatusRequest, OrderResponse> getUpdateOrderStatusMethod() {
        return io.grpc.MethodDescriptor.<UpdateOrderStatusRequest, OrderResponse>newBuilder()
                .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                .setFullMethodName(generateFullMethodName(SERVICE_NAME, "UpdateOrderStatus"))
                .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(UpdateOrderStatusRequest.getDefaultInstance()))
                .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(OrderResponse.getDefaultInstance()))
                .build();
    }
}
