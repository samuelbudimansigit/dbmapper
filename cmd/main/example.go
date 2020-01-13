package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/payfazz/dbmap/pkg/dbmap"
	"github.com/payfazz/dbmap/types"
	"github.com/payfazz/testpath/test"
)

func connectDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "postgres://postgres:cashfazz@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Println(err)
		return nil
	}
	return db
}

type Person struct {
	Id     int     `db:"id"`
	Name   string  `db:"name"`
	JobId  int     `db:"job_id"`
	MyJob  Job     `join:"job"`
	Pacars []Pacar `join:"pacar"`
}

type Pacar struct {
	Id         int     `db:"id"`
	Name       string  `db:"name"`
	PersonId   int     `db:"person_id"`
	PacarShoes []Shoes `join:"shoes"`
}

//Shoes is a test struct
type Shoes struct {
	Id      int    `db:"id"`
	Name    string `db:"name"`
	PacarId int    `db:"pacar_id"`
}

//Job is a test struct
type Job struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

//TestStruct is a test struct
type TestStruct struct {
	Id     int            `db:"id"`
	Name   string         `db:"name"`
	Status types.Metadata `db:"status"`
}

func main() {
	d := connectDB()
	dbMap := dbmap.CreateDBMap(d)
	dbMap.AddTable("person", []string{"id"}, true, Person{})
	dbMap.AddTable("job", []string{"id"}, true, Job{})
	dbMap.AddTable("pacar", []string{"id"}, true, Pacar{})
	dbMap.AddTable("shoes", []string{"id"}, true, Shoes{})
	dbMap.AddRelation(Person{}, dbmap.ONE_TO_ONE, Job{}, "job", "job_id", "id")
	dbMap.AddRelation(Person{}, dbmap.ONE_TO_MANY, Pacar{}, "job", "id", "person_id")
	dbMap.AddRelation(Pacar{}, dbmap.ONE_TO_MANY, Shoes{}, "shoes", "id", "pacar_id")
	dbMap.Init()
	test.Test()

	var person []Person
	dbMap.Select(context.Background(), &person, "person.id=204")
	log.Println(person)
}
