package product

import (
	"backend/pkg/internal/repositories"
	"backend/pkg/internal/storageClient"
	"backend/pkg/internal/structs"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strconv"
)

type repository struct {
	postgresClient storageClient.PostgresClient // *pgxpool.Pool
	redisClient    *redis.Client
	log            *slog.Logger
}

func (r *repository) Create(ctx context.Context, model *structs.Product) error {
	q := `
		INSERT INTO products (name, description, image_url, price, weight)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	if err := r.postgresClient.QueryRow(
		ctx, q,
		model.Name, model.Description,
		model.ImageUrl, model.Price,
		model.Weight,
	).Scan(&model.Id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Error("Postgres error", *pgErr)
			return err
		}
		return err
	}
	r.log.Info("Created new model in Postgres")
	stringResp, err := json.Marshal(model)
	if err != nil {
		r.log.Error("Product marshal error")
	} else {
		r.redisClient.HSet(ctx, "ids", model.Id, string(stringResp))
	}
	return nil
}

func (r *repository) ReadAll(ctx context.Context) ([]structs.Product, error) {
	ids, err := r.redisClient.HGetAll(ctx, "ids").Result()
	if err != nil {
		r.log.Error("Redis get all error")
		return nil, err
	}
	if len(ids) == 0 {
		q := `
		SELECT * FROM products;
		`
		rows, err := r.postgresClient.Query(ctx, q)
		if err != nil {
			return nil, err
		}
		var res []structs.Product
		for rows.Next() {
			var respModel structs.Product
			err = rows.Scan(
				&respModel.Id, &respModel.Name,
				&respModel.Description, &respModel.ImageUrl,
				&respModel.Price, &respModel.Weight,
			)
			if err != nil {
				return nil, err
			}
			res = append(res, respModel)
			stringResp, err := json.Marshal(respModel)
			if err != nil {
				r.log.Error("Product marshal error")
			} else {
				r.redisClient.HSet(ctx, "ids", respModel.Id, string(stringResp))
			}
		}
		return res, nil
	} else {
		var res []structs.Product
		for _, v := range ids {
			var respModel structs.Product
			err = json.Unmarshal([]byte(v), &respModel)
			if err != nil {
				return nil, err
			}
			res = append(res, respModel)
		}
		return res, nil
	}
}

func (r *repository) ReadOne(ctx context.Context, id int) (structs.Product, error) {
	ans, err := r.redisClient.HGet(ctx, "ids", strconv.Itoa(id)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		r.log.Error("Redis get error")
		return structs.Product{}, err
	}
	if errors.Is(err, redis.Nil) {
		q := `
		SELECT * FROM products WHERE id = $1;
		`
		var res structs.Product
		rw := r.postgresClient.QueryRow(ctx, q, id)
		if err := rw.Scan(
			&res.Id, &res.Name,
			&res.Description, &res.ImageUrl,
			&res.Price, &res.Weight,
		); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				r.log.Error("Postgres error", *pgErr)
				return structs.Product{}, err
			}
			return structs.Product{}, err
		}
		stringRes, err := json.Marshal(res)
		if err != nil {
			r.log.Error("json marshal error")
		} else {
			r.redisClient.HSet(ctx, "ids", strconv.Itoa(id), string(stringRes))
		}
		return res, nil
	} else {
		var res structs.Product
		err = json.Unmarshal([]byte(ans), &res)
		if err != nil {
			r.log.Error("Product unmarshal error")
			return structs.Product{}, err
		}
		return res, nil
	}
}

func (r *repository) Update(ctx context.Context, product *structs.Product) error {
	q := `
	UPDATE products SET 
		name = $2, description = $3,
		image_url = $4, price = $5,
		weight = $6
	WHERE id = $1;
	`
	_ = r.postgresClient.QueryRow(ctx, q,
		product.Id, product.Name,
		product.Description, product.ImageUrl,
		product.Price, product.Weight)
	stringRes, err := json.Marshal(product)
	if err != nil {
		r.log.Error("json marshal error")
	} else {
		r.redisClient.HSet(ctx, "ids", product.Id, string(stringRes))
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) (structs.Product, error) {
	q := `
	DELETE FROM products WHERE id = $1 RETURNING name, description, image_url, price, weight;
	`
	var res structs.Product
	err := r.postgresClient.QueryRow(ctx, q, id).Scan(
		&res.Name, &res.Description, &res.ImageUrl,
		&res.Price, &res.Weight,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Error("Postgres error", *pgErr)
			return structs.Product{}, err
		}
		return structs.Product{}, err
	}
	_, err = r.redisClient.HDel(ctx, "ids", strconv.Itoa(id)).Result()
	if err != nil {
		r.log.Error("Redis delete error")
		return structs.Product{}, err
	}
	return res, nil
}

func NewRepository(pClient storageClient.PostgresClient, rClient *redis.Client, log *slog.Logger) repositories.ProductRepository {
	return &repository{
		postgresClient: pClient,
		redisClient:    rClient,
		log:            log,
	}
}
