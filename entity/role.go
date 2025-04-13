package entity

type Role struct {
	Id   int    `db:"id"`
	Role string `db:"role"`
}
