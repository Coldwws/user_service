package main

import (
	"context"
	"log"
	"user_service/internal/config"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("postgres/local.migrate.env")
	_ = godotenv.Load("local.env")

	cfg := config.LoadConfig()

	ctx := context.Background()

	poll,err := pgxpool.Connect(ctx, cfg.PG.DSN())
	if err != nil{
		log.Fatalf("Failed to connect database: %v",err)
	}
	defer poll.Close()

	log.Println("Connected to Postgres")

	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("first_name","last_name","phone_number","email").
		Values(gofakeit.FirstName(),gofakeit.LastName(),gofakeit.Phone(),gofakeit.Email()).
		Suffix("RETURNING id")

	query,args,err := builderInsert.ToSql()
	if err != nil{
		log.Fatalf("Failed to build query: %v",err)
	}
	
	var userID int
	err = poll.QueryRow(ctx,query,args...).Scan(&userID)
	if err != nil{
		log.Fatalf("Failed to execute query: %v",err)
	}

	log.Printf("Inserted user with ID: %d", userID)



	builderSelect := sq.Select("id","email","phone_number").From("users").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)


	query,args,err = builderSelect.ToSql()
	if err != nil{
		log.Fatalf("Failed to build query: %v",err)
	}
	rows,err := poll.Query(ctx,query,args...)
	if err != nil{
		log.Fatalf("Failed to execute query: %v",err)
	}
	
	var id int
	var email string
	var phoneNumber string

	for rows.Next(){
		err := rows.Scan(&id,&email,&phoneNumber)
		if err != nil{
			log.Printf("failed to scan users: %v",err)
		}

		log.Printf("User: ID=%d, Email=%s, Phone=%s",id,email,phoneNumber)
	}

	// var updateID = 2
	// builderUpdate := sq.Update("users").PlaceholderFormat(sq.Dollar).
	// Set("email","test@mail.ru").
	// Set("phone_number","+77787200477").
	// Where(sq.Eq{"id":updateID})

 
	// query,args,err = builderUpdate.ToSql()
	// if err != nil{
	// 	log.Fatalf("Failed to build query: %v",err)

	// }
	// res,err := poll.Exec(ctx,query,args...)
	// if err != nil{
	// 	log.Fatalf("Failed to execute query: %v",err)
	// }
	// log.Printf("Updated rows: %d", res.RowsAffected())



	// deletedID := 10
	// builderDelete := sq.Delete("users").
	// 	PlaceholderFormat(sq.Dollar).
	// 	Where(sq.Eq{"id": deletedID})

	// 	query,args,err = builderDelete.ToSql()
	// 	if err != nil{
	// 		log.Fatalf("Failed to build query: %v",err)
	// 	}

	// 	res, err = poll.Exec(ctx,query,args...)
	// 	if err != nil{
	// 		log.Fatalf("Failed to execute query: %v",err)
	// 	}

	// 	log.Printf("Deleted row :%v",res.RowsAffected())
}