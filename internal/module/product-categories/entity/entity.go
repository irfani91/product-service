package entity

import "codebase-app/pkg/types"

type CreateProductCategoriesRequest struct {
	Name string `json:"name" validate:"required" db:"name"`
}

type CreateProductCategoriesResponse struct {
	Id string `json:"id" db:"id"`
}

type GetProductCategoriesRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type GetProductCategoriesResponse struct {
	Name string `json:"name" db:"name"`
}

type DeleteProductCategoriesRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type UpdateProductCategoriesRequest struct {
	Id   string `params:"id" validate:"uuid" db:"id"`
	Name string `json:"name" validate:"required" db:"name"`
}

type UpdateProductCategoriesResponse struct {
	Id string `json:"id" db:"id"`
}

type ProductCategoriesRequest struct {
	Page     int `query:"page" validate:"required"`
	Paginate int `query:"paginate" validate:"required"`
}

func (r *ProductCategoriesRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type ProductCategoriesItem struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type ProductCategoriesResponse struct {
	Items []ProductCategoriesItem `json:"items"`
	Meta  types.Meta              `json:"meta"`
}
