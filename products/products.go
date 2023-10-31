package products

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

type ProductOp struct {
	Payload Product `json:"payload"`
	Op      string  `json:"op"`
}

type Repo struct {
	db     *gorm.DB
	pubFun func(context.Context, ProductOp)
}

func NewRepo(db *gorm.DB, pubFun func(context.Context, ProductOp)) *Repo {
	return &Repo{db, pubFun}
}

func (r *Repo) List(offset int, limit int) ([]Product, error) {
	var c []Product
	result := r.db.Find(&c).Limit(limit).Offset(offset)
	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil

}

func (r *Repo) CreateOne(name string, price uint) (*Product, error) {
	c := &Product{Name: name, Price: price}

	result := r.db.Create(c)

	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil
}

func (r *Repo) FindById(id int) (*Product, error) {
	c := &Product{}
	result := r.db.First(c, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, nil
	}
	return c, nil
}

func (r *Repo) UpdateOne(id int, p Product) (bool, error) {
	toUpdate := make(map[string]interface{})
	if p.Name != "" {
		toUpdate["name"] = p.Name
	}
	if p.Price != 0 {
		toUpdate["price"] = p.Price
	}

	result := r.db.Model(Product{}).Where("id = ?", id).Updates(toUpdate)
	return result.RowsAffected > 0, result.Error
}
