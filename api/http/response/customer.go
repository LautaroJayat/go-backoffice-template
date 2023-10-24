package response

import (
	users "github.com/lautarojayat/e_shop/users"
)

type CustomerResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	DeletedAt int64  `json:"deletedAt"`
}

func ToCustomersResponse(c []users.User) []CustomerResponse {
	out := make([]CustomerResponse, len(c))
	for i := 0; i < len(c); i += 1 {
		out = append(out, ToCustomerResponse(c[i]))
	}
	return out
}

func ToCustomerResponse(c users.User) CustomerResponse {
	return CustomerResponse{
		Id:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Local().UnixMilli(),
		UpdatedAt: c.UpdatedAt.UnixMilli(),
		DeletedAt: c.DeletedAt.Time.UnixMilli(),
	}
}

type PaginatedCustomers struct {
	Pagination Pagination         `json:"pagination"`
	Customers  []CustomerResponse `json:"users"`
}
