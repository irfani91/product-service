package repository

import (
	"codebase-app/internal/module/product-categories/entity"
	"codebase-app/internal/module/product-categories/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductCategoriesRepository = &productCategoriesRepository{}

type productCategoriesRepository struct {
	db *sqlx.DB
}

func NewProductCategoriesRepository(db *sqlx.DB) *productCategoriesRepository {
	return &productCategoriesRepository{
		db: db,
	}
}

func (r *productCategoriesRepository) CreateProductCategories(ctx context.Context, req *entity.CreateProductCategoriesRequest) (*entity.CreateProductCategoriesResponse, error) {
	var resp = new(entity.CreateProductCategoriesResponse)
	// Your code here
	query := `
		INSERT INTO product_categories (name)
		VALUES (?) RETURNING id
	`

	err := r.db.QueryRowContext(ctx, r.db.Rebind(query),
		req.Name,
	).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::CreateProductCategories - Failed to create ProductCategories")
		return nil, err
	}

	return resp, nil
}

func (r *productCategoriesRepository) GetProductCategories(ctx context.Context, req *entity.GetProductCategoriesRequest) (*entity.GetProductCategoriesResponse, error) {
	var resp = new(entity.GetProductCategoriesResponse)
	// Your code here
	query := `
		SELECT name
		FROM product_categories
		WHERE id = ?
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(resp)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProductCategories - Failed to get ProductCategories")
		return nil, err
	}

	return resp, nil
}

func (r *productCategoriesRepository) DeleteProductCategories(ctx context.Context, req *entity.DeleteProductCategoriesRequest) error {
	query := `
		UPDATE product_categories
		SET deleted_at = NOW()
		WHERE id = ? 
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteProductCategories - Failed to delete ProductCategories")
		return err
	}

	return nil
}

func (r *productCategoriesRepository) UpdateProductCategories(ctx context.Context, req *entity.UpdateProductCategoriesRequest) (*entity.UpdateProductCategoriesResponse, error) {
	var resp = new(entity.UpdateProductCategoriesResponse)

	query := `
		UPDATE product_categories
		SET name = ?, updated_at = NOW()
		WHERE id = ?
		RETURNING id
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.Name,
		req.Id).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::UpdateProductCategories - Failed to update ProductCategories")
		return nil, err
	}

	return resp, nil
}

func (r *productCategoriesRepository) GetProductCategoriess(ctx context.Context, req *entity.ProductCategoriesRequest) (*entity.ProductCategoriesResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.ProductCategoriesItem
	}

	var (
		resp = new(entity.ProductCategoriesResponse)
		data = make([]dao, 0, req.Paginate)
	)
	resp.Items = make([]entity.ProductCategoriesItem, 0, req.Paginate)

	query := `
		SELECT
			COUNT(id) OVER() as total_data,
			id,
			name
		FROM product_categories
		WHERE
			deleted_at IS NULL
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query),
		req.Paginate,
		req.Paginate*(req.Page-1),
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProductCategories - Failed to get ProductCategories")
		return nil, err
	}

	if len(data) > 0 {
		resp.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		resp.Items = append(resp.Items, d.ProductCategoriesItem)
	}

	resp.Meta.CountTotalPage(req.Page, req.Paginate, resp.Meta.TotalData)

	return resp, nil
}
