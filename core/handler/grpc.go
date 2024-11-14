package handler

import (
	"context"
	"encoding/json"
	"fmt"
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
	//product.ImageUrls = req.Product.ImageUrls
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
	product.SellerId = sellerId

	// Image Parsing
	imageUrlsJson, err := json.Marshal(req.Product.ImageUrls)
	if err != nil {
		log.Error("Error marshaling data: %v", err)
	}
	product.ImageUrls = string(imageUrlsJson)

	err = service.CreateProduct(product, s.RdbInstance)
	if err != nil {
		return &proto.CreateProductResponse{
			Message: "Failed to create the product. error: " + err.Error(),
		}, err
	}

	// Return the created product
	return &proto.CreateProductResponse{Message: "Successfully created the product listing"}, nil
}

func (s *Server) GetProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductResponse, error) {

	productResult, err := service.GetProduct(req.GetProductId(), s.RdbInstance)

	fmt.Println(err)

	if err != nil {
		return nil, err
	}

	if len(productResult) <= 0 {
		log.Error("no products found", nil)
		return &proto.GetProductResponse{
			Message: "No products found",
		}, fmt.Errorf("no products found")
	}

	product := productResult[0]

	response := &proto.Product{
		ProductId:             product.ID.String(),
		Name:                  product.Name,
		Quantity:              product.Quantity,
		Type:                  product.Type,
		Category:              product.Category,
		Price:                 product.Price,
		Size:                  &proto.Product_Size{Width: product.Width, Height: product.Height},
		Weight:                product.Weight,
		ShippingBasePrice:     product.ShippingBasePrice,
		BaseDeliveryTimelines: product.BaseDeliveryTimelines,
		SellerId:              product.SellerId.String(),
	}

	err = json.Unmarshal([]byte(product.ImageUrls), &response.ImageUrls)
	if err != nil {
		log.Error("Error unmarshalling JSON: %v", err)
		return &proto.GetProductResponse{
			Message: "No products found",
		}, err
	}

	return &proto.GetProductResponse{Message: "Successfully retrieved the product", Product: response}, nil
}

func (s *Server) GetProducts(ctx context.Context, req *proto.GetProductsRequest) (*proto.GetProductsResponse, error) {

	products, err := service.GetProducts(s.RdbInstance)
	if err != nil {
		return nil, err
	}
	var response []*proto.Product

	for _, product := range *products {
		res := &proto.Product{
			ProductId:             product.ID.String(),
			Name:                  product.Name,
			Quantity:              product.Quantity,
			Type:                  product.Type,
			Category:              product.Category,
			Price:                 product.Price,
			Size:                  &proto.Product_Size{Width: product.Width, Height: product.Height},
			Weight:                product.Weight,
			ShippingBasePrice:     product.ShippingBasePrice,
			BaseDeliveryTimelines: product.BaseDeliveryTimelines,
			SellerId:              product.SellerId.String(),
		}
		err = json.Unmarshal([]byte(product.ImageUrls), &res.ImageUrls)
		if err != nil {
			log.Error("Error unmarshalling JSON: %v", err)
		}
		response = append(response, res)
	}
	return &proto.GetProductsResponse{Message: "Successfully retrieved the product", Products: response}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {

	productResult, err := service.GetProduct(req.GetProductId(), s.RdbInstance)
	if err != nil {
		return nil, err
	}

	product := productResult[0]

	if product.SellerId.String() != req.Product.SellerId {
		return &proto.UpdateProductResponse{
			Message: "Unauthorized to perform this operation",
		}, fmt.Errorf("unauthorized to perform this operation")
	}

	if req.Product.Name != "" {
		product.Name = req.Product.Name
	}
	if req.Product.Quantity > 0 {
		product.Quantity = req.Product.Quantity
	}
	if req.Product.Type != "" {
		product.Type = req.Product.Type
	}
	if req.Product.Category != "" {
		product.Category = req.Product.Category
	}
	if req.Product.Price > 0 {
		product.Price = req.Product.Price
	}

	if req.Product.GetSize() != nil && req.Product.Size.Width > 0 {
		product.Width = req.Product.Size.Width
	}
	if req.Product.GetSize() != nil && req.Product.Size.Height > 0 {
		product.Height = req.Product.Size.Height
	}
	if req.Product.ShippingBasePrice > 0 {
		product.ShippingBasePrice = req.Product.ShippingBasePrice
	}
	if req.Product.BaseDeliveryTimelines > 0 {
		product.BaseDeliveryTimelines = req.Product.BaseDeliveryTimelines
	}

	if len(req.Product.ImageUrls) > 0 {
		// Image Parsing
		imageUrlsJson, err := json.Marshal(req.Product.ImageUrls)
		if err != nil {
			log.Error("Error marshaling data: %v", err)
		}
		product.ImageUrls = string(imageUrlsJson)
	}

	err = service.UpdateProduct(product, s.RdbInstance)
	if err != nil {
		return &proto.UpdateProductResponse{
			Message: "Failed to update the product. error: " + err.Error(),
		}, err
	}

	// Return the created product
	return &proto.UpdateProductResponse{Message: "Successfully updated the product listing"}, nil

}
