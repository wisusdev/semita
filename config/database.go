package config

import (
	"database/sql"
	"os"
	"path/filepath"
	"semita/app/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func DatabaseConnect() *sql.DB {
	var databaseDriver = utils.GetEnv("DB_DRIVER")

	switch databaseDriver {
	case "mysql":
		return MysqlDatabaseConnect()
	case "postgres":
		return PostgresDatabaseConnect()
	case "sqlite":
		return SqliteDatabaseConnect()
	default:
		panic("Unsupported database driver: " + databaseDriver)
	}
}

func MysqlDatabaseConnect() *sql.DB {
	var godot = godotenv.Load(".env")
	if godot != nil {
		panic("Error loading .env file")
	}

	var driver = os.Getenv("DB_DRIVER")
	var host = os.Getenv("DB_HOST")
	var port = os.Getenv("DB_PORT")
	var dbname = os.Getenv("DB_NAME")
	var user = os.Getenv("DB_USER")
	var password = os.Getenv("DB_PASSWORD")

	var db *sql.DB
	var err error

	db, err = sql.Open(driver, user+":"+password+"@tcp("+host+":"+port+")/"+dbname)

	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	return db
}

func PostgresDatabaseConnect() *sql.DB {
	var godot = godotenv.Load(".env")
	if godot != nil {
		panic("Error loading .env file")
	}

	var driver = os.Getenv("DB_DRIVER")
	var host = os.Getenv("DB_HOST")
	var port = os.Getenv("DB_PORT")
	var dbname = os.Getenv("DB_NAME")
	var user = os.Getenv("DB_USER")
	var password = os.Getenv("DB_PASSWORD")

	var db *sql.DB
	var err error

	db, err = sql.Open(driver, "host="+host+" port="+port+" user="+user+" password="+password+" dbname="+dbname+" sslmode=disable")

	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	return db
}

func SqliteDatabaseConnect() *sql.DB {
	var dbname = utils.GetEnv("DB_NAME")
	var dbPath = filepath.Join("storage", dbname+".db")

	if errorMkdir := os.MkdirAll("storage", os.ModePerm); errorMkdir != nil {
		utils.Logs("error", "Error creating storage directory: "+errorMkdir.Error())
		panic("Error creating storage directory: " + errorMkdir.Error())
	}

	var db *sql.DB
	var err error

	db, err = sql.Open("sqlite3", dbPath)

	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	return db
}
