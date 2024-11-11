package handler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tittuvarghese/core/logger"
	"github.com/tittuvarghese/product-service/core/database"
	"github.com/tittuvarghese/product-service/models"
	"github.com/tittuvarghese/product-service/proto"
	"github.com/tittuvarghese/product-service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	proto.UnimplementedProductServiceServer
	GrpcServer  *grpc.Server
	RdbInstance *database.RelationalDatabase
}

var log = logger.NewLogger("product-service")

func NewGrpcServer() *Server {
	return &Server{GrpcServer: grpc.NewServer()}
}

func (s *Server) Run(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error("Failed to listen", err)
	}

	proto.RegisterProductServiceServer(s.GrpcServer, s)

	// Register reflection service on gRPC server
	reflection.Register(s.GrpcServer)
	log.Info("GRPC server is listening on port " + port)
	if err := s.GrpcServer.Serve(lis); err != nil {
		log.Error("failed to serve", err)
	}
}

func (s *Server) mustEmbedUnimplementedAuthServiceServer() {
	log.Error("implement me", nil)
}

func (s *Server) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	var product models.Product

	product.Name = req.Product.Name
	product.Quantity = req.Product.Quantity
	product.Type = req.Product.Type
	product.Category = req.Product.Category
	//product.ImageURLs = req.Product.ImageUrls
	product.Price = req.Product.Price
	product.Width = req.Product.Size.Width
	product.Height = req.Product.Size.Height
	product.Weight = req.Product.Weight
	product.ShippingBasePrice = req.Product.ShippingBasePrice
	product.BaseDeliveryTimelines = req.Product.BaseDeliveryTimelines

	sellerId, err := uuid.Parse(req.Product.SellerId)
	if err != nil {
		return &proto.CreateProductResponse{
			Message: "Unable to parse seller id",
		}, err
	}
	product.SellerID = sellerId

	// Image Parsing
	imageUrlsJson, err := json.Marshal(req.Product.ImageUrls)
	if err != nil {
		log.Error("Error marshaling data: %v", err)
	}
	product.ImageURLs = string(imageUrlsJson)

	err = service.CreateProduct(product, s.RdbInstance)
	if err != nil {
		return &proto.CreateProductResponse{
			Message: "Failed to register the user. error: " + err.Error(),
		}, err
	}

	// Return the created product
	return &proto.CreateProductResponse{Message: "Successfully created the product listing"}, nil
}
