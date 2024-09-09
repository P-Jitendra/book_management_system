package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// const (
// 	host     = "host.docker.internal"
// 	port     = 5432
// 	user     = "me"
// 	password = "password"
// 	dbname   = "api"
// 	sslmode  = "disable"
// )

var DB *gorm.DB

func ConnectDB() {

	loadErr := godotenv.Load(".env")

	if loadErr != nil {
		fmt.Println("Error loading .env file")
	}
	connStr := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		"disable",
	)

	// connInfo := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
	// 	host, user, password, dbname, port, sslmode)

	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})

	// if err != nil{
	// 	panic(err)
	// }
	// defer db.Close();

	if err != nil {
		fmt.Printf("Failed to connect to database, err:%v\n", err)
	} else {
		fmt.Println("Connected to Postgres database successfully!")
	}
	//It's important to call Ping() methoed as sql.Open method does not create a connection to the database
	//instead it simply validates the arguments provided
	//By calling db.Ping() we actually open up a connection to the database which will validate whether
	//or not our connection string was 100% correct.
	// err = db.Ping()
	// if err != nil{
	// 	panic(err);
	// }
	fmt.Println("Printing DB data: ", DB)
}
