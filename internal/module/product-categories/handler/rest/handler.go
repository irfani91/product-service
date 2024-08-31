package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/product-categories/entity"
	"codebase-app/internal/module/product-categories/ports"
	"codebase-app/internal/module/product-categories/repository"
	"codebase-app/internal/module/product-categories/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type productCategoriesHandler struct {
	service ports.ProductCategoriesService
}

func NewProductCategoriesHandler() *productCategoriesHandler {
	var (
		handler = new(productCategoriesHandler)
		repo    = repository.NewProductCategoriesRepository(adapter.Adapters.ShopeefunPostgres)
		service = service.NewProductCategoriesService(repo)
	)
	handler.service = service

	return handler
}

func (h *productCategoriesHandler) Register(router fiber.Router) {
	router.Get("/categories", middleware.UserIdHeader, h.GetProductCategoriess)
	router.Post("/category", middleware.UserIdHeader, h.CreateProductCategories)
	router.Get("/category/:id", h.GetProductCategories)
	router.Delete("/category/:id", middleware.UserIdHeader, h.DeleteProductCategories)
	router.Patch("/category/:id", middleware.UserIdHeader, h.UpdateProductCategories)
}

func (h *productCategoriesHandler) CreateProductCategories(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateProductCategoriesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::ProductCategories - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::CreateProductCategories - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.CreateProductCategories(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))

}

func (h *productCategoriesHandler) GetProductCategories(c *fiber.Ctx) error {
	var (
		req = new(entity.GetProductCategoriesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetProductCategories - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetProductCategories(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *productCategoriesHandler) DeleteProductCategories(c *fiber.Ctx) error {
	var (
		req = new(entity.DeleteProductCategoriesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::DeleteProductCategories - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.DeleteProductCategories(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *productCategoriesHandler) UpdateProductCategories(c *fiber.Ctx) error {
	var (
		req = new(entity.UpdateProductCategoriesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateShop - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::UpdateProductCategories - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.UpdateProductCategories(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *productCategoriesHandler) GetProductCategoriess(c *fiber.Ctx) error {
	var (
		req = new(entity.ProductCategoriesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetProductCategories - Parse request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetProductCategories - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetProductCategoriess(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))

}
