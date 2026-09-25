package postgres

import (
	"context"
	"fmt"
	"manga/internal/models"
	"manga/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Category() store.CategoryRepository {
	if db.Categorys == nil {
		db.Categorys = NewCategoryRepository(db.conn)
	}
	return db.Categorys
}

type CategoriesRepository struct {
	conn *sqlx.DB
}

func NewCategoryRepository(conn *sqlx.DB) store.CategoryRepository {
	return &CategoriesRepository{conn: conn}
}

func (c CategoriesRepository) Create(ctx context.Context, category *models.Category) error {
	_, err := c.conn.Exec("INSERT INTO categories(name) VALUES ($1)", category.Name)
	if err != nil {
		return fmt.Errorf("Unknow err:%S", err)
	}
	return nil
}
func (c CategoriesRepository) Update(ctx context.Context, category *models.Category) error {
	_, err := c.conn.Exec("INSERT INTO categories(name) VALUES ($1)", category.Name)
	if err != nil {
		return fmt.Errorf("Unknow err:%S", err)
	}
	return nil
}
func (c CategoriesRepository) Get(ctx context.Context, filter *models.Categoryesfilter) ([]*models.Category, error) {
	basicQuery := "SELECT * FROM categories"
	args := []interface{}{}
	if filter.Query != nil {
		basicQuery += " WHERE name ILIKE $1"
		args = append(args, "%"+*filter.Query+"%")
	}
	categories := make([]*models.Category, 0)
	if err := c.conn.Select(&categories, basicQuery, args...); err != nil {
		return nil, fmt.Errorf("%S", err)
	}
	return categories, nil
}
func (c CategoriesRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM Categories WHERE user_id=$1", id)
	if err != nil {
		panic(err)
	}
	return nil
}
