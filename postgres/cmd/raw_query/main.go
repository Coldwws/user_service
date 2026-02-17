package main

import (
	"context"
	"log"
	"user_service/internal/config"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("postgres/local.migrate.env")

	_ = godotenv.Load("local.env")


	cfg := config.LoadConfig()
	ctx := context.Background()

	conn,err := pgx.Connect(ctx,cfg.PG.DSN())
	if err != nil{
		log.Fatalf("connect error: %v",err)
	}
	defer conn.Close(context.Background())
	log.Println("Connected to Postgres")


	res, err := conn.Exec(ctx,"INSERT INTO users(first_name,last_name,phone_number,email) VALUES($1,$2,$3,$4)", gofakeit.FirstName(),gofakeit.LastName(),gofakeit.Phone(),gofakeit.Email())
	if err != nil{
		log.Fatalf("failed to insert user: %v",err)
	}
	log.Printf("Inserted rows: %d", res.RowsAffected())


	rows, err := conn.Query(ctx, "SELECT id, first_name, last_name, phone_number, email FROM users")
	if err != nil{
		log.Fatalf("failed to query users : %v",err)
	}

	defer rows.Close()

	for rows.Next(){
		var id int
		var firstName, lastName, phoneNumber, email string
		err := rows.Scan(&id, &firstName, &lastName, &phoneNumber, &email)
		if err != nil{
			log.Printf("failed to scan user: %v",err)
			continue
		}
		log.Printf("User: ID=%d, Name=%s %s, Phone=%s, Email=%s", id, firstName, lastName, phoneNumber, email)
	}

	
	deletedId := 1
	deleleteRes, err:= conn.Exec(ctx,"DELETE FROM users WHERE id = $1",deletedId)
	if err != nil{
		log.Fatalf("failed to delete user: %v",err)
	}
	log.Printf("Deleted rows: %d", deleleteRes.RowsAffected())

}