package config

type dbConfig struct {
	DB_URI        string
	DB_NAME       string
}

func GetDBConfig() dbConfig {
	return dbConfig{
		DB_URI:        "mongodb://localhost:27017",
		DB_NAME:       "db_kita",
	}
}
