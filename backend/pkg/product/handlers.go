package product

import (
	"backend/pkg/internal/repositories"
	"backend/pkg/internal/structs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

func NewFindAll(log *slog.Logger, repository repositories.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "product.handlers.NewFindAll"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		res, err := repository.ReadAll(r.Context())
		if err != nil {
			log.Error("Failed to read products", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, struct{ Products []structs.Product }{
			Products: res,
		})
	}
}

func NewFindOne(log *slog.Logger, repository repositories.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "product.handlers.NewFindOne"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		prodId, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			log.Error("Failed to get product Id", err)
			return
		}
		model, err := repository.ReadOne(r.Context(), prodId)
		if err != nil {
			log.Error("Failed to get product", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, model)
	}
}

func NewCreate(log *slog.Logger, repository repositories.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "product.handlers.NewAdd"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req structs.Product
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", err)
			return
		}
		err = repository.Create(r.Context(), &req)
		if err != nil {
			log.Error("Failed to insert product", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, req)
	}
}

func NewUpdate(log *slog.Logger, repository repositories.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "product.handlers.NewUpdate"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		prodId, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			log.Error("Failed to get product Id", err)
			return
		}
		var req structs.Product
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", err)
			return
		}
		req.Id = prodId
		err = repository.Update(r.Context(), &req)
		if err != nil {
			log.Error("Failed to update product", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, req)
	}
}

func NewDelete(log *slog.Logger, repository repositories.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "product.handlers.NewDelete"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		prodId, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			log.Error("Failed to get product Id", err)
			return
		}
		prod, err := repository.Delete(r.Context(), prodId)
		if err != nil {
			log.Error("Failed to delete product", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, prod)
	}
}
