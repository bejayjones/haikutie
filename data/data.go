package data

import (
	"context"
	"database/sql"
	"embed"
	db "haikutie/data/db"
	"log"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

//go:embed schema.sql
var schemaFS embed.FS

type Helper struct {
	*sql.DB
	*db.Queries
}

type Haiku = db.Haiku
type User = db.User
type GetAllUsersRow = db.GetAllUsersRow
type GetReceivedHaikusRow = db.GetReceivedHaikusRow

func InitDB() *Helper {
	// Open database with pure Go driver
	database, err := sql.Open("sqlite", "./sqlite/haiku.db")
	if err != nil {
		log.Fatal(err)
	}

	// Read and execute schema
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	// Create queries instance
	queries := db.New(database)

	return &Helper{
		DB:      database,
		Queries: queries,
	}
}

// Helper methods that wrap sqlc-generated code

func (h *Helper) CreateUserHelper(username, password string) error {
	_, err := h.Queries.CreateUser(context.Background(), db.CreateUserParams{
		Username: username,
		Password: password,
	})
	return err
}

func (h *Helper) GetUserHelper(username string) (User, error) {
	return h.Queries.GetUser(context.Background(), username)
}

func (h *Helper) GetUserByIDHelper(id int64) (User, error) {
	return h.Queries.GetUserByID(context.Background(), id)
}

func (h *Helper) GetAllUsersHelper() ([]GetAllUsersRow, error) {
	return h.Queries.GetAllUsers(context.Background())
}

func (h *Helper) CreateHaikuHelper(fromUserID, toUserID int64, line1, line2, line3 string) error {
	_, err := h.Queries.CreateHaiku(context.Background(), db.CreateHaikuParams{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Line1:      line1,
		Line2:      line2,
		Line3:      line3,
	})
	return err
}

func (h *Helper) GetReceivedHaikusHelper(userID int64) ([]GetReceivedHaikusRow, error) {
	return h.Queries.GetReceivedHaikus(context.Background(), userID)
}
