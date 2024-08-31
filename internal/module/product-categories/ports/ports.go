package ports

import (
	"codebase-app/internal/module/product-categories/entity"
	"context"
)

type ProductCategoriesRepository interface {
	CreateProductCategories(ctx context.Context, req *entity.CreateProductCategoriesRequest) (*entity.CreateProductCategoriesResponse, error)
	GetProductCategories(ctx context.Context, req *entity.GetProductCategoriesRequest) (*entity.GetProductCategoriesResponse, error)
	DeleteProductCategories(ctx context.Context, req *entity.DeleteProductCategoriesRequest) error
	UpdateProductCategories(ctx context.Context, req *entity.UpdateProductCategoriesRequest) (*entity.UpdateProductCategoriesResponse, error)
	GetProductCategoriess(ctx context.Context, req *entity.ProductCategoriesRequest) (*entity.ProductCategoriesResponse, error)
}

type ProductCategoriesService interface {
	CreateProductCategories(ctx context.Context, req *entity.CreateProductCategoriesRequest) (*entity.CreateProductCategoriesResponse, error)
	GetProductCategories(ctx context.Context, req *entity.GetProductCategoriesRequest) (*entity.GetProductCategoriesResponse, error)
	DeleteProductCategories(ctx context.Context, req *entity.DeleteProductCategoriesRequest) error
	UpdateProductCategories(ctx context.Context, req *entity.UpdateProductCategoriesRequest) (*entity.UpdateProductCategoriesResponse, error)
	GetProductCategoriess(ctx context.Context, req *entity.ProductCategoriesRequest) (*entity.ProductCategoriesResponse, error)
}
