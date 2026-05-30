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

type TourPurchaseRPCClient interface {
	GetTourForPurchase(ctx context.Context, in *GetTourForPurchaseRequest, opts ...grpc.CallOption) (*GetTourForPurchaseResponse, error)
	ValidateTourPurchase(ctx context.Context, in *ValidateTourPurchaseRequest, opts ...grpc.CallOption) (*ValidateTourPurchaseResponse, error)
}

type tourPurchaseRPCClient struct {
	conn grpc.ClientConnInterface
}

func NewTourPurchaseRPCClient(conn grpc.ClientConnInterface) TourPurchaseRPCClient {
	return &tourPurchaseRPCClient{conn: conn}
}

func (c *tourPurchaseRPCClient) GetTourForPurchase(ctx context.Context, in *GetTourForPurchaseRequest, opts ...grpc.CallOption) (*GetTourForPurchaseResponse, error) {
	out := new(GetTourForPurchaseResponse)
	err := c.conn.Invoke(ctx, "/tourpurchase.TourPurchaseRpc/GetTourForPurchase", in, out, opts...)
	return out, err
}

func (c *tourPurchaseRPCClient) ValidateTourPurchase(ctx context.Context, in *ValidateTourPurchaseRequest, opts ...grpc.CallOption) (*ValidateTourPurchaseResponse, error) {
	out := new(ValidateTourPurchaseResponse)
	err := c.conn.Invoke(ctx, "/tourpurchase.TourPurchaseRpc/ValidateTourPurchase", in, out, opts...)
	return out, err
}
