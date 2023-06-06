package dba

import "database/sql"

// User struct represents user table.
type User struct {
	Id       int
	Username string
	FkRole   int
}

// Role struct represents role table.
type Role struct {
	Id          int
	Name        string
	Description sql.NullString
}

// Permission struct represents permission table.
type Permission struct {
	Id          int
	Name        string
	Description sql.NullString
}

// RoleHasPermission struct represents role_has_permission table.
type RoleHasPermission struct {
	Id           int
	FkRole       int
	FkPermission int
}
