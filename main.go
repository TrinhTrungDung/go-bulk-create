package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Name     string
	Password string
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=testuser dbname=testdb password=123456 sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	fmt.Println("Successfully connect to new database")

	db.AutoMigrate(&User{})
}
