package repository

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ak-repo/wim/internal/db"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/jackc/pgx/v5"
)

type RefCodeGenerator interface {
	GenerateUserRefCode(ctx context.Context) (string, error)
	GenerateCustomerRefCode(ctx context.Context) (string, error)
	GenerateProductRefCode(ctx context.Context) (string, error)
	GenerateWarehouseRefCode(ctx context.Context) (string, error)
	GenerateLocationRefCode(ctx context.Context) (string, error)
}

type refCodeGenerator struct {
	db    *db.DB
	redis *db.Redis
}

func NewRefCodeGenerator(database *db.DB, redis *db.Redis) RefCodeGenerator {
	return &refCodeGenerator{db: database, redis: redis}
}

func (r *refCodeGenerator) GenerateUserRefCode(ctx context.Context) (string, error) {
	return r.GenerateRefCode(ctx, "USR", "users", "ref_code")
}

func (r *refCodeGenerator) GenerateCustomerRefCode(ctx context.Context) (string, error) {
	return r.GenerateRefCode(ctx, "CUS", "customers", "ref_code")
}

func (r *refCodeGenerator) GenerateProductRefCode(ctx context.Context) (string, error) {
	return r.GenerateRefCode(ctx, "PRD", "products", "ref_code")
}

func (r *refCodeGenerator) GenerateWarehouseRefCode(ctx context.Context) (string, error) {
	return r.GenerateRefCode(ctx, "WH", "warehouses", "ref_code")
}

func (r *refCodeGenerator) GenerateLocationRefCode(ctx context.Context) (string, error) {
	return r.GenerateRefCode(ctx, "LOC", "locations", "ref_code")
}

func (r *refCodeGenerator) GenerateRefCode(ctx context.Context, prefix string, tableName string, columnName string) (string, error) {

	// Redis Cache
	cacheKey := fmt.Sprintf("refCode:%s", prefix)
	if r.redis != nil {
		cached, err := r.redis.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			return r.nextRefcode(ctx, prefix, cached)
		}
	}

	lastCode, err := r.getLastRefCode(ctx, tableName, columnName, prefix)
	if err != nil {
		return "", err
	}
	if lastCode == "" {
		return fmt.Sprintf("%s-%03d", prefix, 1), nil
	}
	return r.nextRefcode(ctx, prefix, lastCode)

}
func (r *refCodeGenerator) getLastRefCode(ctx context.Context, tableName string, columnName string, prefix string) (string, error) {

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIKE $1 ORDER BY %s DESC LIMIT 1", columnName, tableName, columnName, columnName)

	var refCode string
	err := r.db.Pool.QueryRow(ctx, query, prefix+"-%").Scan(&refCode)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return refCode, nil
}

func (r *refCodeGenerator) nextRefcode(ctx context.Context, prefix, lastRefCode string) (string, error) {
	if lastRefCode == "" {
		return fmt.Sprintf("%s-%03d", prefix, 1), nil
	}

	re := regexp.MustCompile("^" + regexp.QuoteMeta(prefix) + `-(\d+)$`)
	matches := re.FindStringSubmatch(lastRefCode)
	if len(matches) < 2 {
		return "", apperrors.New(apperrors.CodeRefCodeFailed, "invalid reference code format")
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", apperrors.Wrap(err, apperrors.CodeRefCodeFailed, "failed to parse reference code")
	}
	nextNum := num + 1
	nextRefCode := fmt.Sprintf("%s-%03d", prefix, nextNum)

	// Redis Cache
	cacheKey := fmt.Sprintf("refCode:%s", prefix)
	if r.redis != nil {
		err := r.redis.Set(ctx, cacheKey, nextRefCode, 5)
		if err != nil {
			log.Printf("failed to cache %s: %v", cacheKey, err)
		}
	}

	return nextRefCode, nil
}
