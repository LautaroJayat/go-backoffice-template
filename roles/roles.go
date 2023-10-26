package roles

type Role uint8

const (
	ReadProduct Role = 1 << iota
	CreateProduct
	ModifyProduct
	DeleteProduct
	ReadUser
	CreateUser
	ModifyUser
	DeleteUser

	ProductAdmin       = ReadProduct | CreateProduct | ModifyProduct | DeleteProduct
	UserAdmin          = ReadUser | CreateUser | ModifyUser | DeleteUser
	SuperAdmin         = ProductAdmin | UserAdmin
	DecodedPermsHeader = "X-Decoded-Perms"
)
