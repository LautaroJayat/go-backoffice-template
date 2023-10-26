package response

import (
	users "github.com/lautarojayat/backoffice/users"
)

type UserResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	DeletedAt int64  `json:"deletedAt"`
}

func ToUsersResponse(c []users.User) []UserResponse {
	out := make([]UserResponse, len(c))
	for i := 0; i < len(c); i += 1 {
		out = append(out, ToUserResponse(c[i]))
	}
	return out
}

func ToUserResponse(c users.User) UserResponse {
	return UserResponse{
		Id:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Local().UnixMilli(),
		UpdatedAt: c.UpdatedAt.UnixMilli(),
		DeletedAt: c.DeletedAt.Time.UnixMilli(),
	}
}

type PaginatedUsers struct {
	Pagination Pagination     `json:"pagination"`
	Users      []UserResponse `json:"users"`
}
