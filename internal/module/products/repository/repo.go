package repository

import (
	"codebase-app/internal/module/products/entity"
	"codebase-app/internal/module/products/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductRepository = &productRepository{}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error) {
	var resp = new(entity.CreateProductResponse)
	// Your code here
	query := `
		INSERT INTO
			products (
				shop_id,
				category_id,
				name,
				description,
				image_url,
				price,
				brand,
				stock
			)
			VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9 )
			RETURNING
				id, shop_id,category_id, name, description, image_url, price, brand, stock, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, r.db.Rebind(query),
		req.ShopId,
		req.CategoryId,
		req.Name,
		req.Description,
		req.ImageUrl,
		req.Price,
		req.Brand,
		req.Stock).Scan(&resp.Id, &resp.ShopId, &resp.CategoryId, &resp.Name, &resp.Description, &resp.ImageUrl, &resp.Price, &resp.Brand, &resp.Stock, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::CreateProduct - Failed to create Product")
		return nil, err
	}
	return resp, nil
}

func (r *productRepository) GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error) {
	var resp = new(entity.GetProductResponse)
	// Your code here
	query := `
		SELECT 
		id,
			category_id,
			shop_id,
			name,
			image_url,
			price,
			stock,
			brand,
			created_at,
			updated_at
		FROM
			products
		WHERE
			deleted_at IS NULL
		AND id = ?
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(resp)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProduct - Failed to get Product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error {
	query := `
		UPDATE products
		SET deleted_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteProduct - Failed to delete Product")
		return err
	}

	return nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductResponse, error) {
	var resp = new(entity.UpdateProductResponse)

	query := `
		UPDATE
			products
		SET
			category_id = $1,
			name = $2,
			description = $3,
			image_url = $4,
			price = $5,
			stock = $6,
			brand = $7,
			updated_at = NOW()
		WHERE
			id = $8
			AND deleted_at IS NULL
		RETURNING
			id, shop_id, category_id, name, description, image_url, price, stock, brand, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.CategoryId,
		req.Name,
		req.Description,
		req.ImageUrl,
		req.Price,
		req.Stock,
		req.Brand,
		req.Id).Scan(&resp.Id, &resp.ShopId, &resp.CategoryId, &resp.Name, &resp.Description, &resp.ImageUrl, &resp.Price, &resp.Stock, &resp.Brand, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::UpdateProduct - Failed to update Product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) GetProducts(ctx context.Context, req *entity.ProductsRequest) (*entity.ProductsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.ProductItem
	}

	var (
		resp = new(entity.ProductsResponse)
		data = make([]dao, 0, req.Paginate)
		arg  = make(map[string]any)
	)
	resp.Items = make([]entity.ProductItem, 0, req.Paginate)

	query := `
		SELECT
			COUNT(id) OVER() as total_data,
			id,
			category_id,
			shop_id,
			name,
			image_url,
			price,
			brand,
			created_at,
			updated_at
		FROM
			products
		WHERE
			deleted_at IS NULL
	`
	if req.ShopId != "" {
		query += " AND shop_id = :shop_id"
		arg["shop_id"] = req.ShopId
	}

	if req.CategoryId != "" {
		query += " AND category_id = :category_id"
		arg["category_id"] = req.CategoryId
	}

	if req.Name != "" {
		query += " AND name ILIKE '%' || :name || '%'"
		arg["name"] = req.Name
	}

	if req.PriceMinStr != "" {
		query += " AND price >= :price_min"
		arg["price_min"] = req.PriceMin
	}

	if req.PriceMaxStr != "" {
		query += " AND price <= :price_max"
		arg["price_max"] = req.PriceMax
	}

	if req.IsAvailable {
		query += " AND stock > 0"
	}

	query += `
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query),
		req.UserId,
		req.Paginate,
		req.Paginate*(req.Page-1),
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProducts - Failed to get Products")
		return nil, err
	}

	if len(data) > 0 {
		resp.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		resp.Items = append(resp.Items, d.ProductItem)
	}

	resp.Meta.CountTotalPage(req.Page, req.Paginate, resp.Meta.TotalData)

	return resp, nil
}
