package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var DB *pgx.Conn

func Connect() {
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	//defer db.Close(context.Background())
	DB = db
	fmt.Println("DATABASE::CONNECTED")
}
