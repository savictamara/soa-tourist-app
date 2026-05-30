package tourpurchase

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

type GetTourForPurchaseRequest struct {
	TourID string `json:"tourId"`
}

type GetTourForPurchaseResponse struct {
	TourID   string  `json:"tourId"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
	Archived bool    `json:"archived"`
}

type ValidateTourPurchaseRequest struct {
	TourID string `json:"tourId"`
}

type ValidateTourPurchaseResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

type TourPurchaseRPCServer interface {
	GetTourForPurchase(context.Context, *GetTourForPurchaseRequest) (*GetTourForPurchaseResponse, error)
	ValidateTourPurchase(context.Context, *ValidateTourPurchaseRequest) (*ValidateTourPurchaseResponse, error)
}

func RegisterTourPurchaseRPCServer(server *grpc.Server, service TourPurchaseRPCServer) {
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "tourpurchase.TourPurchaseRpc",
		HandlerType: (*TourPurchaseRPCServer)(nil),
		Methods: []grpc.MethodDesc{
			{MethodName: "GetTourForPurchase", Handler: getTourForPurchaseHandler},
			{MethodName: "ValidateTourPurchase", Handler: validateTourPurchaseHandler},
		},
	}, service)
}

func getTourForPurchaseHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(GetTourForPurchaseRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TourPurchaseRPCServer).GetTourForPurchase(ctx, req)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tourpurchase.TourPurchaseRpc/GetTourForPurchase"}
	return interceptor(ctx, req, info, func(ctx context.Context, request interface{}) (interface{}, error) {
		return srv.(TourPurchaseRPCServer).GetTourForPurchase(ctx, request.(*GetTourForPurchaseRequest))
	})
}

func validateTourPurchaseHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(ValidateTourPurchaseRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TourPurchaseRPCServer).ValidateTourPurchase(ctx, req)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tourpurchase.TourPurchaseRpc/ValidateTourPurchase"}
	return interceptor(ctx, req, info, func(ctx context.Context, request interface{}) (interface{}, error) {
		return srv.(TourPurchaseRPCServer).ValidateTourPurchase(ctx, request.(*ValidateTourPurchaseRequest))
	})
}
