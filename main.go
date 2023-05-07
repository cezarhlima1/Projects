package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "crudgo"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database!")

	// Create a user
	user := User{Name: "John Doe", Email: "johndoe@example.com"}
	err = createUser(db, &user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User created with ID: %d\n", user.ID)

	// Read a user
	readUser, err := readUser(db, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User read: %+v\n", readUser)

	// Update a user
	user.Name = "Jane Doe"
	err = updateUser(db, user.ID, &user)
	if err != nil {
		log.Fatal(err)
	}

}

func createUser(db *sql.DB, user *User) error {
	sqlStatement := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(sqlStatement, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func readUser(db *sql.DB, id int) (*User, error) {
	sqlStatement := `SELECT * FROM users WHERE id=$1;`
	row := db.QueryRow(sqlStatement, id)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func updateUser(db *sql.DB, id int, user *User) error {
	sqlStatement := `UPDATE users SET name=$2, email=$3 WHERE id=$1;`
	_, err := db.Exec(sqlStatement, id, user.Name, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func deleteUser(db *sql.DB, id int) error {
	sqlStatement := `DELETE FROM users WHERE id=$1;`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}
