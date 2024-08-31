package service

import (
	"codebase-app/internal/module/product-categories/entity"
	"codebase-app/internal/module/product-categories/ports"
	"context"
)

var _ ports.ProductCategoriesService = &productCategoriesService{}

type productCategoriesService struct {
	repo ports.ProductCategoriesRepository
}

func NewProductCategoriesService(repo ports.ProductCategoriesRepository) *productCategoriesService {
	return &productCategoriesService{
		repo: repo,
	}
}

func (s *productCategoriesService) CreateProductCategories(ctx context.Context, req *entity.CreateProductCategoriesRequest) (*entity.CreateProductCategoriesResponse, error) {
	return s.repo.CreateProductCategories(ctx, req)
}

func (s *productCategoriesService) GetProductCategories(ctx context.Context, req *entity.GetProductCategoriesRequest) (*entity.GetProductCategoriesResponse, error) {
	return s.repo.GetProductCategories(ctx, req)
}

func (s *productCategoriesService) DeleteProductCategories(ctx context.Context, req *entity.DeleteProductCategoriesRequest) error {
	return s.repo.DeleteProductCategories(ctx, req)
}

func (s *productCategoriesService) UpdateProductCategories(ctx context.Context, req *entity.UpdateProductCategoriesRequest) (*entity.UpdateProductCategoriesResponse, error) {
	return s.repo.UpdateProductCategories(ctx, req)
}

func (s *productCategoriesService) GetProductCategoriess(ctx context.Context, req *entity.ProductCategoriesRequest) (*entity.ProductCategoriesResponse, error) {
	return s.repo.GetProductCategoriess(ctx, req)
}
