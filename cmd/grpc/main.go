// gRPC API for Catalog bounded context (L06 HW5 task 02).
package main

import (
	"context"
	"log"
	"net"

	catalogv1 "github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog/gen/catalog/v1"
	"github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog/internal/catalog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type catalogServer struct {
	catalogv1.UnimplementedCatalogServiceServer
	store *catalog.Store
}

func (s *catalogServer) ListProducts(
	_ context.Context,
	_ *catalogv1.ListProductsRequest,
) (*catalogv1.ListProductsResponse, error) {
	items := s.store.List()
	out := make([]*catalogv1.Product, 0, len(items))
	for _, p := range items {
		out = append(out, toProto(p))
	}
	return &catalogv1.ListProductsResponse{Products: out}, nil
}

func (s *catalogServer) GetProduct(
	_ context.Context,
	req *catalogv1.GetProductRequest,
) (*catalogv1.GetProductResponse, error) {
	p, ok := s.store.Get(req.GetId())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "product %q not found", req.GetId())
	}
	return &catalogv1.GetProductResponse{Product: toProto(p)}, nil
}

func toProto(p catalog.Product) *catalogv1.Product {
	return &catalogv1.Product{
		Id:          p.ID,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	catalogv1.RegisterCatalogServiceServer(srv, &catalogServer{store: catalog.NewStore()})

	log.Printf("catalog gRPC listening on %s", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
