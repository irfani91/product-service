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

func (r *productRepository) GetProducts(ctx context.Context, req *entity.ProductsRequest) (entity.ProductsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.Product
	}
	var (
		res  entity.ProductsResponse
		data = make([]dao, 0)
		arg  = make(map[string]any)
	)
	res.Meta.Page = req.Page
	res.Meta.Paginate = req.Paginate

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			products.id as id,
			product_categories.id as category_id,
			product_categories.name as category,
			products.shop_id as shop_id,
			shops.name as shop_name,
			products.name as name,
			image_url,
			price,
			brand,
			products.created_at as created_at,
			products.updated_at as updated_at
		FROM
			products
		JOIN shops 
			ON products.shop_id = shops.id
		JOIN product_categories
			ON products.category_id = product_categories.id
		WHERE
			products.deleted_at IS NULL
	`

	if req.ShopId != "" {
		query += " AND products.shop_id = :shop_id"
		arg["shop_id"] = req.ShopId
	}

	if req.CategoryId != "" {
		query += " AND products.category_id = :category_id"
		arg["category_id"] = req.CategoryId
	}

	if req.Name != "" {
		query += " AND products.name ILIKE '%' || :name || '%'"
		arg["name"] = req.Name
	}

	if req.Brand != "" {
		query += " AND products.brand ILIKE '%' || :brand || '%'"
		arg["brand"] = req.Brand
	}

	if req.PriceMinStr != "" {
		query += " AND products.price >= :price_min"
		arg["price_min"] = req.PriceMin
	}

	if req.PriceMaxStr != "" {
		query += " AND products.price <= :price_max"
		arg["price_max"] = req.PriceMax
	}

	if req.IsAvailable {
		query += " AND products.stock > 0"
	}

	query += `
		ORDER BY products.created_at DESC
		LIMIT :paginate
		OFFSET :offset
	`
	arg["paginate"] = req.Paginate
	arg["offset"] = (req.Page - 1) * req.Paginate

	nstmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: GetProducts failed")
		return res, err
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &data, arg)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: GetProducts failed")
		return res, err
	}

	for _, d := range data {
		res.Items = append(res.Items, entity.Product{
			Id:         d.Id,
			CategoryId: d.CategoryId,
			Category:   d.Category,
			ShopId:     d.ShopId,
			ShopName:   d.ShopName,
			Name:       d.Name,
			ImageUrl:   d.ImageUrl,
			Price:      d.Price,
			CreatedAt:  d.CreatedAt,
			UpdatedAt:  d.UpdatedAt,
		})

		res.Meta.TotalData = d.TotalData
	}

	res.Meta.CountTotalPage()
	return res, nil

}
