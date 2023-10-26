package users

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UsersOp struct {
	Payload User   `json:"payload"`
	Op      string `json:"op"`
}

type User struct {
	gorm.Model
	Name string `json:"name"`
}

type Repo struct {
	db     *gorm.DB
	pubFun func(ctx context.Context, payload UsersOp)
}

func NewRepo(db *gorm.DB, propagationFunction func(ctx context.Context, payload UsersOp)) *Repo {

	return &Repo{db, propagationFunction}
}

func (r *Repo) List(offset int, limit int) ([]User, error) {
	var c []User
	result := r.db.Find(&c).Limit(limit).Offset(offset)
	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil

}

func (r *Repo) CreateOne(name string) (*User, error) {
	c := User{Name: name}

	result := r.db.Create(&c)

	if result.Error != nil {
		return nil, result.Error
	}

	go func(c User) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()
		r.pubFun(ctx, UsersOp{c, "create"})
	}(c)

	return &c, nil
}

func (r *Repo) FindById(id int) (*User, error) {
	c := &User{}
	result := r.db.First(c, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return c, nil
}

func (r *Repo) UpdateOne(id int, name string) (bool, error) {
	result := r.db.Model(User{}).Where("id = ?", id).Update("name", name)
	return result.RowsAffected > 0, result.Error
}
