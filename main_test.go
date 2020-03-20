package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/thoas/go-funk"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func genRandomString(length int) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func stubUsers(b *testing.B) (users []*User) {
	for i := 0; i < b.N; i++ {
		user := &User{
			Name:     genRandomString(100),
			Password: genRandomString(100),
		}
		users = append(users, user)
	}

	return users
}

func BenchmarkOrmCreate(b *testing.B) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=testuser dbname=testdb password=123456 sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	users := stubUsers(b)
	for _, user := range users {
		db.Create(user)
	}
}

func BenchmarkCreate(b *testing.B) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=testuser dbname=testdb password=123456 sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	users := stubUsers(b)
	tx := db.Begin()
	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, user := range users {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, user.Name)
		valueArgs = append(valueArgs, user.Password)
	}

	stmt := fmt.Sprintf("INSERT INTO users (name, password) VALUES %s", strings.Join(valueStrings, ","))
	err = tx.Exec(stmt, valueArgs...).Error
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
	}
	err = tx.Commit().Error
	if err != nil {
		fmt.Println(err)
	}
}

func BenchmarkBulkCreate(b *testing.B) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=testuser dbname=testdb password=123456 sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	users := stubUsers(b)
	size := 500
	tx := db.Begin()
	chunkList := funk.Chunk(users, size)
	for _, chunk := range chunkList.([][]*User) {
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for _, user := range chunk {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, user.Name)
			valueArgs = append(valueArgs, user.Password)
		}

		stmt := fmt.Sprintf("INSERT INTO users (name, password) VALUES %s", strings.Join(valueStrings, ","))
		err = tx.Exec(stmt, valueArgs...).Error
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
		}
	}
	err = tx.Commit().Error
	if err != nil {
		fmt.Println(err)
	}
}

func benchmarkBulkCreate(size int, b *testing.B) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=testuser dbname=testdb password=123456 sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	users := stubUsers(b)
	tx := db.Begin()
	chunkList := funk.Chunk(users, size)
	for _, chunk := range chunkList.([][]*User) {
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for _, user := range chunk {
			now := time.Now()
			valueStrings = append(valueStrings, "(?, ?, ?, ?)")
			valueArgs = append(valueArgs, now)
			valueArgs = append(valueArgs, now)
			valueArgs = append(valueArgs, user.Name)
			valueArgs = append(valueArgs, user.Password)
		}

		stmt := fmt.Sprintf("INSERT INTO users (created_at, updated_at, name, password) VALUES %s", strings.Join(valueStrings, ","))
		err = tx.Exec(stmt, valueArgs...).Error
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
		}
	}
	err = tx.Commit().Error
	if err != nil {
		fmt.Println(err)
	}
}

func BenchmarkBulkCreateSize1(b *testing.B) {
	benchmarkBulkCreate(1, b)
}

func BenchmarkBulkCreateSize100(b *testing.B) {
	benchmarkBulkCreate(100, b)
}

func BenchmarkBulkCreateSize500(b *testing.B) {
	benchmarkBulkCreate(500, b)
}

func BenchmarkBulkCreateSize1000(b *testing.B) {
	benchmarkBulkCreate(1000, b)
}
