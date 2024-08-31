package entity

import (
	"codebase-app/pkg/types"
	"strconv"
	"time"
)

type CreateProductRequest struct {
	UserId string `query:"user_id" validate:"required,uuid"`

	ShopId      string  `json:"shop_id" validate:"required,uuid" db:"shop_id"`
	CategoryId  string  `json:"category_id" validate:"required,uuid" db:"category_id"`
	Name        string  `json:"name" validate:"required,max=255,min=3" db:"name"`
	Description *string `json:"description" validate:"omitempty,max=255,min=3" db:"description"`
	ImageUrl    *string `json:"image_url" validate:"omitempty,url" db:"image_url"`
	Brand       *string `json:"brand" validate:"omitempty,max=255,min=3" db:"brand"`
	Price       float64 `json:"price" validate:"required,numeric" db:"price"`
	Stock       int64   `json:"stock" validate:"required,numeric" db:"stock"`
}

type CreateProductResponse struct {
	Id          string    `json:"id" db:"id"`
	ShopId      string    `json:"shop_id" db:"shop_id"`
	CategoryId  string    `json:"category_id" db:"category_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	ImageUrl    *string   `json:"image_url" db:"image_url"`
	Brand       *string   `json:"brand" validate:"omitempty,max=255,min=3" db:"brand"`
	Price       float64   `json:"price" db:"price"`
	Stock       int       `json:"stock" db:"stock"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type GetProductRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type GetProductResponse struct {
	Id          string    `json:"id" db:"id"`
	ShopId      string    `json:"shop_id" db:"shop_id"`
	CategoryId  string    `json:"category_id" db:"category_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	ImageUrl    *string   `json:"image_url" db:"image_url"`
	Price       float64   `json:"price" db:"price"`
	Stock       int       `json:"stock" db:"stock"`
	Brand       *string   `json:"brand" validate:"omitempty,max=255,min=3" db:"brand"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type DeleteProductRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type UpdateProductRequest struct {
	UserId string `query:"user_id" validate:"required,uuid"`

	Id          string  `params:"id" validate:"required,uuid" db:"id"`
	CategoryId  string  `json:"category_id" validate:"required,uuid" db:"category_id"`
	Name        string  `json:"name" validate:"required,max=255,min=3" db:"name"`
	Description *string `json:"description" validate:"omitempty,max=255,min=3" db:"description"`
	ImageUrl    *string `json:"image_url" validate:"omitempty,url" db:"image_url"`
	Brand       *string `json:"brand" validate:"omitempty,max=255,min=3" db:"brand"`
	Price       float64 `json:"price" validate:"required,numeric" db:"price"`
	Stock       int64   `json:"stock" validate:"required,numeric" db:"stock"`
}

type UpdateProductResponse struct {
	Id          string    `json:"id" db:"id"`
	UserId      string    `json:"user_id" db:"user_id"`
	ShopId      string    `json:"shop_id" db:"shop_id"`
	CategoryId  string    `json:"category_id" db:"category_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	ImageUrl    *string   `json:"image_url" db:"image_url"`
	Price       float64   `json:"price" db:"price"`
	Stock       int       `json:"stock" db:"stock"`
	Brand       *string   `json:"brand" validate:"omitempty,max=255,min=3" db:"brand"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ProductsRequest struct {
	UserId      string `prop:"user_id" validate:"uuid"`
	ShopId      string `query:"shop_id" validate:"omitempty,uuid"`
	CategoryId  string `query:"category_id" validate:"omitempty,uuid"`
	Name        string `query:"name" validate:"omitempty,max=255,min=3"`
	PriceMinStr string `query:"price_min" validate:"omitempty,numeric,gte=0"`
	PriceMaxStr string `query:"price_max" validate:"omitempty,numeric,gte=0"`
	IsAvailable bool   `query:"is_available"`

	Page     int `query:"page" validate:"required"`
	Paginate int `query:"paginate" validate:"required"`

	PriceMin float64
	PriceMax float64
}

func (r *ProductsRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type ProductItem struct {
	Id         string    `json:"id" db:"id"`
	CategoryId string    `json:"category_id" db:"category_id"`
	ShopId     string    `json:"shop_id" db:"shop_id"`
	Name       string    `json:"name" db:"name"`
	ImageUrl   *string   `json:"image_url" db:"image_url"`
	Brand      *string   `json:"brand" db:"brand"`
	Price      float64   `json:"price" db:"price"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ProductsResponse struct {
	Items []ProductItem `json:"items"`
	Meta  types.Meta    `json:"meta"`
}

func (r *ProductsRequest) CostumValidation() (int, map[string][]string) {
	var (
		errors   = make(map[string][]string)
		err      error
		priceMin float64
		priceMax float64
	)

	if r.PriceMinStr != "" {
		priceMin, err = strconv.ParseFloat(r.PriceMinStr, 64)
		if err != nil {
			errors["price_min"] = append(errors["price_min"], "price_min must be a number.")
		}
		r.PriceMin = priceMin
	}

	if r.PriceMaxStr != "" {
		priceMax, err = strconv.ParseFloat(r.PriceMaxStr, 64)
		if err != nil {
			errors["price_max"] = append(errors["price_max"], "price_max must be a number.")
		}
		r.PriceMax = priceMax
	}

	if len(errors) > 0 {
		return 400, errors
	}

	errors = nil
	return 0, errors
}
