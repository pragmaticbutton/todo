package dba

const (
	dsn = "root:root@tcp(localhost:13306)/todo"
)

var (
	da *DatabaseAccess
)

func setupTestCase() func() {
	d, err := NewDatabaseAccess(dsn)
	if err != nil {
		panic(err)
	}
	da = d

	cleanDatabase()

	return func() {
		if err := da.db.Close(); err != nil {
			panic(err)
		}
	}
}

func cleanDatabase() {
	da.db.Exec("DELETE FROM role_has_permission")
	da.db.Exec("DELETE FROM role")
	da.db.Exec("DELETE FROM permission")
	da.db.Exec("DELETE FROM user")
}
