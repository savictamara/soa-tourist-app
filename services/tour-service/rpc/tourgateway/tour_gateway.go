package tourgateway

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

type CompletedKeyPointMessage struct {
	KeyPointID   string `json:"keyPointId"`
	KeyPointName string `json:"keyPointName"`
	ReachedAt    string `json:"reachedAt"`
}

type TourExecutionMessage struct {
	ID                 string                     `json:"id"`
	TourID             string                     `json:"tourId"`
	TourName           string                     `json:"tourName"`
	TouristID          string                     `json:"touristId"`
	Status             string                     `json:"status"`
	StartedAt          string                     `json:"startedAt"`
	CompletedAt        string                     `json:"completedAt,omitempty"`
	AbandonedAt        string                     `json:"abandonedAt,omitempty"`
	LastActivityAt     string                     `json:"lastActivityAt"`
	StartLatitude      float64                    `json:"startLatitude"`
	StartLongitude     float64                    `json:"startLongitude"`
	CurrentLatitude    float64                    `json:"currentLatitude"`
	CurrentLongitude   float64                    `json:"currentLongitude"`
	CompletedKeyPoints []CompletedKeyPointMessage `json:"completedKeyPoints"`
}

type StartTourExecutionRequest struct {
	TourID    string  `json:"tourId"`
	TouristID string  `json:"touristId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type StartTourExecutionResponse struct {
	Success   bool                 `json:"success"`
	Message   string               `json:"message"`
	Execution TourExecutionMessage `json:"execution"`
}

type CheckTourExecutionLocationRequest struct {
	ExecutionID string  `json:"executionId"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type CheckTourExecutionLocationResponse struct {
	Success            bool                      `json:"success"`
	Message            string                    `json:"message"`
	Execution          TourExecutionMessage      `json:"execution"`
	KeyPointReached    bool                      `json:"keyPointReached"`
	ReachedKeyPoint    *CompletedKeyPointMessage `json:"reachedKeyPoint,omitempty"`
	DistanceMeters     float64                   `json:"distanceMeters"`
	LastActivityAt     string                    `json:"lastActivityAt"`
	CompletedCount     int                       `json:"completedCount"`
	TotalKeyPointCount int                       `json:"totalKeyPointCount"`
}

type TourGatewayRpcClient interface {
	StartTourExecution(ctx context.Context, in *StartTourExecutionRequest, opts ...grpc.CallOption) (*StartTourExecutionResponse, error)
	CheckTourExecutionLocation(ctx context.Context, in *CheckTourExecutionLocationRequest, opts ...grpc.CallOption) (*CheckTourExecutionLocationResponse, error)
}

type tourGatewayRpcClient struct {
	conn grpc.ClientConnInterface
}

func NewTourGatewayRpcClient(conn grpc.ClientConnInterface) TourGatewayRpcClient {
	return &tourGatewayRpcClient{conn: conn}
}

func (c *tourGatewayRpcClient) StartTourExecution(ctx context.Context, in *StartTourExecutionRequest, opts ...grpc.CallOption) (*StartTourExecutionResponse, error) {
	out := new(StartTourExecutionResponse)
	err := c.conn.Invoke(ctx, "/tourgateway.TourGatewayRpc/StartTourExecution", in, out, opts...)
	return out, err
}

func (c *tourGatewayRpcClient) CheckTourExecutionLocation(ctx context.Context, in *CheckTourExecutionLocationRequest, opts ...grpc.CallOption) (*CheckTourExecutionLocationResponse, error) {
	out := new(CheckTourExecutionLocationResponse)
	err := c.conn.Invoke(ctx, "/tourgateway.TourGatewayRpc/CheckTourExecutionLocation", in, out, opts...)
	return out, err
}

type TourGatewayRpcServer interface {
	StartTourExecution(context.Context, *StartTourExecutionRequest) (*StartTourExecutionResponse, error)
	CheckTourExecutionLocation(context.Context, *CheckTourExecutionLocationRequest) (*CheckTourExecutionLocationResponse, error)
}

func RegisterTourGatewayRpcServer(server *grpc.Server, implementation TourGatewayRpcServer) {
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "tourgateway.TourGatewayRpc",
		HandlerType: (*TourGatewayRpcServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "StartTourExecution",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(StartTourExecutionRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return srv.(TourGatewayRpcServer).StartTourExecution(ctx, in)
					}
					info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tourgateway.TourGatewayRpc/StartTourExecution"}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(TourGatewayRpcServer).StartTourExecution(ctx, req.(*StartTourExecutionRequest))
					}
					return interceptor(ctx, in, info, handler)
				},
			},
			{
				MethodName: "CheckTourExecutionLocation",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(CheckTourExecutionLocationRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return srv.(TourGatewayRpcServer).CheckTourExecutionLocation(ctx, in)
					}
					info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tourgateway.TourGatewayRpc/CheckTourExecutionLocation"}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(TourGatewayRpcServer).CheckTourExecutionLocation(ctx, req.(*CheckTourExecutionLocationRequest))
					}
					return interceptor(ctx, in, info, handler)
				},
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "tour_gateway.proto",
	}, implementation)
}
