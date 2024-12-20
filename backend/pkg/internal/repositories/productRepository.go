package repositories

import (
	"backend/pkg/internal/structs"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, model *structs.Product) error
	ReadAll(ctx context.Context) ([]structs.Product, error)
	ReadOne(ctx context.Context, id int) (structs.Product, error)
	Update(ctx context.Context, model *structs.Product) error
	Delete(ctx context.Context, id int) (structs.Product, error)
}
