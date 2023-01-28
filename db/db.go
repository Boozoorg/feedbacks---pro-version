package db

func SetupDB() {
	setupPostgres()
	setupRedis()
}

func CloseDB() {
	closePostgres()
	closeRedis()
}