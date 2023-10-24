package response

import (
	"github.com/lautarojayat/e_shop/products"
)

type ProductResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Price     uint   `json:"price"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	DeletedAt int64  `json:"deletedAt"`
}

func ToProductsResponse(p []products.Product) []ProductResponse {
	out := make([]ProductResponse, len(p))
	for i := 0; i < len(p); i += 1 {
		out = append(out, ToProductResponse(p[i]))
	}
	return out
}

func ToProductResponse(p products.Product) ProductResponse {
	return ProductResponse{
		Id:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		CreatedAt: p.CreatedAt.Local().UnixMilli(),
		UpdatedAt: p.UpdatedAt.UnixMilli(),
		DeletedAt: p.DeletedAt.Time.UnixMilli(),
	}
}

type PaginatedProducts struct {
	Pagination Pagination        `json:"pagination"`
	Products   []ProductResponse `json:"products"`
}
