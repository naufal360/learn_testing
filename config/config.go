package config

import (
	"fmt"
	"learn_testing/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Init() {
	InitDB()
	InitialMigrate()
}

func InitDB() {
	config := models.Config{
		DB_Username: ViperEnvVariable("DB_USERNAME"),
		DB_Password: ViperEnvVariable("DB_PASSWORD"),
		DB_Port:     ViperEnvVariable("DB_PORT"),
		DB_Host:     ViperEnvVariable("DB_HOST"),
		DB_Name:     ViperEnvVariable("DB_NAME"),
	}
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_Name,
	)

	var err error

	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

// auto migrate with db
func InitialMigrate() {
	DB.AutoMigrate(&models.Users{}, &models.Books{})
}

// Test Func
// func InitDBTest() {

// 	const (
// 		DB_USERNAME_TEST = "root"
// 		DB_PASSWORD_TEST = ""
// 		DB_PORT_TEST     = "3306"
// 		DB_HOST_TEST     = "localhost"
// 		DB_NAME_TEST     = "crud_go_mvc_test"
// 	)
// 	// DB_Username_Test := ViperEnvVariable("DB_USERNAME_TEST")
// 	// DB_Password_Test := ViperEnvVariable("DB_PASSWORD_TEST")
// 	// DB_Port_Test := ViperEnvVariable("DB_PORT_TEST")
// 	// DB_Host_Test := ViperEnvVariable("DB_HOST_TEST")
// 	// DB_Name_Test := ViperEnvVariable("DB_NAME_TEST")

// 	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		DB_USERNAME_TEST,
// 		DB_PASSWORD_TEST,
// 		DB_HOST_TEST,
// 		DB_PORT_TEST,
// 		DB_NAME_TEST,
// 	)

// 	var err error

// 	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	InitialMigrateTest()
// }

// func InitialMigrateTest() {
// 	DB.Migrator().DropTable(&models.User{})
// 	DB.Migrator().DropTable(&models.Book{})
// 	DB.AutoMigrate(&models.User{})
// 	DB.AutoMigrate(&models.Book{})
// }
