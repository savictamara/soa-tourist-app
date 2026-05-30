package purchasegateway

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

const CodecName = "json"

type JSONCodec struct{}

func (JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (JSONCodec) Name() string {
	return CodecName
}

func RegisterCodec() {
	encoding.RegisterCodec(JSONCodec{})
}

type AddToCartRequest struct {
	TouristID string  `json:"touristId"`
	TourID    string  `json:"tourId"`
	TourName  string  `json:"tourName"`
	Price     float64 `json:"price"`
}

type CheckoutRequest struct {
	TouristID string `json:"touristId"`
}

type OrderItemMessage struct {
	TourID   string  `json:"tourId"`
	TourName string  `json:"tourName"`
	Price    float64 `json:"price"`
}

type ShoppingCartMessage struct {
	TouristID  string             `json:"touristId"`
	Items      []OrderItemMessage `json:"items"`
	TotalPrice float64            `json:"totalPrice"`
}

type PurchaseTokenMessage struct {
	TouristID   string `json:"touristId"`
	TourID      string `json:"tourId"`
	PurchasedAt string `json:"purchasedAt"`
	Token       string `json:"token"`
}

type ShoppingCartResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Cart    ShoppingCartMessage `json:"cart"`
}

type CheckoutResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Tokens  []PurchaseTokenMessage `json:"tokens"`
}

type PurchaseGatewayRpcClient interface {
	AddToCart(ctx context.Context, in *AddToCartRequest, opts ...grpc.CallOption) (*ShoppingCartResponse, error)
	Checkout(ctx context.Context, in *CheckoutRequest, opts ...grpc.CallOption) (*CheckoutResponse, error)
}

type purchaseGatewayRpcClient struct {
	conn grpc.ClientConnInterface
}

func NewPurchaseGatewayRpcClient(conn grpc.ClientConnInterface) PurchaseGatewayRpcClient {
	return &purchaseGatewayRpcClient{conn: conn}
}

func (c *purchaseGatewayRpcClient) AddToCart(ctx context.Context, in *AddToCartRequest, opts ...grpc.CallOption) (*ShoppingCartResponse, error) {
	out := new(ShoppingCartResponse)
	err := c.conn.Invoke(ctx, "/purchasegateway.PurchaseGatewayRpc/AddToCart", in, out, opts...)
	return out, err
}

func (c *purchaseGatewayRpcClient) Checkout(ctx context.Context, in *CheckoutRequest, opts ...grpc.CallOption) (*CheckoutResponse, error) {
	out := new(CheckoutResponse)
	err := c.conn.Invoke(ctx, "/purchasegateway.PurchaseGatewayRpc/Checkout", in, out, opts...)
	return out, err
}

type PurchaseGatewayRpcServer interface {
	AddToCart(context.Context, *AddToCartRequest) (*ShoppingCartResponse, error)
	Checkout(context.Context, *CheckoutRequest) (*CheckoutResponse, error)
}

func RegisterPurchaseGatewayRpcServer(server *grpc.Server, implementation PurchaseGatewayRpcServer) {
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "purchasegateway.PurchaseGatewayRpc",
		HandlerType: (*PurchaseGatewayRpcServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "AddToCart",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(AddToCartRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return srv.(PurchaseGatewayRpcServer).AddToCart(ctx, in)
					}
					info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/purchasegateway.PurchaseGatewayRpc/AddToCart"}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(PurchaseGatewayRpcServer).AddToCart(ctx, req.(*AddToCartRequest))
					}
					return interceptor(ctx, in, info, handler)
				},
			},
			{
				MethodName: "Checkout",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(CheckoutRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return srv.(PurchaseGatewayRpcServer).Checkout(ctx, in)
					}
					info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/purchasegateway.PurchaseGatewayRpc/Checkout"}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(PurchaseGatewayRpcServer).Checkout(ctx, req.(*CheckoutRequest))
					}
					return interceptor(ctx, in, info, handler)
				},
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "purchase_gateway.proto",
	}, implementation)
}
