package main

import (
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/payfazz/dbmap/pkg/dbmap"
	"github.com/payfazz/testpath/test"
	"golang.org/x/net/context"
)

func connectDBY() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "postgres://postgres:cashfazz@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Println(err)
		return nil
	}
	return db
}

type Flat struct {
	PersonId      int    `db:"person.id"`
	PersonName    string `db:"person.name"`
	PersonJobId   int    `db:"person.job_id"`
	PacarId       int    `db:"pacar.id"`
	PacarName     string `db:"pacar.name"`
	PacarPersonId int    `db:"pacar.person_id"`
	ShoesId       int    `db:"shoes.id"`
	ShoesName     string `db:"shoes.name"`
	ShoesPacarId  int    `db:"shoes.pacar_id"`
	JobId         int    `db:"job.id"`
	JobName       string `db:"job.name"`
}

func BenchmarkRaw(b *testing.B) {
	var person []dbmap.Person
	test.Test()
	db, _ := sqlx.Connect("postgres", "postgres://postgres:cashfazz@localhost/postgres?sslmode=disable")
	query :=
		`SELECT person.id as "person.id",
	person.name as "person.name",
	person.job_id as "person.job_id",
	pacar.id as "pacar.id",
	pacar.name as "pacar.name",
	pacar.person_id as "pacar.person_id" ,
	job.id as "job.id",
	job.name as "job.name",
	shoes.id as "shoes.id",
	shoes."name" as "shoes.name",
	shoes.pacar_id as "shoes.pacar_id"
	FROM person JOIN pacar ON person.id=pacar.person_id
		join job on job.id = person.job_id
		join shoes on shoes.pacar_id = pacar.id
	WHERE person.id=204 or person.id=203
	`
	for n := 0; n < b.N; n++ {
		stmt, _ := db.Preparex(query)
		rows, _ := stmt.Queryx()
		var a []Flat
		var tsh []dbmap.Shoes
		var tpcr []dbmap.Pacar
		var b Flat
		var jb dbmap.Job
		person = nil
		var p dbmap.Person
		for rows.Next() {
			err := rows.StructScan(&b)
			if err == nil {
				a = append(a, b)
			}
		}
		for i := 0; i < len(a); i++ {
			personId := a[i].PersonId
			for j := i; j < len(a); j++ {
				pacarId := a[j].PacarId
				for k := j; k < len(a); k++ {
					if a[k].PacarId != pacarId {
						j = k
						break
					}
					sh := dbmap.Shoes{
						Id:      a[k].ShoesId,
						Name:    a[k].ShoesName,
						PacarId: a[k].ShoesPacarId,
					}
					tsh = append(tsh, sh)
				}
				if personId != a[j].PersonId {
					i = j
					break
				}
				pr := dbmap.Pacar{
					Id:         a[j].PacarId,
					Name:       a[j].PacarName,
					PersonId:   a[j].PacarPersonId,
					PacarShoes: tsh,
				}
				tpcr = append(tpcr, pr)
				tsh = nil
				jb = dbmap.Job{
					Id:   a[j].JobId,
					Name: a[j].JobName,
				}

			}
			p = dbmap.Person{
				Name:   a[i].PersonName,
				Id:     a[i].PersonId,
				JobId:  a[i].PersonJobId,
				MyJob:  jb,
				Pacars: tpcr,
			}
			person = append(person, p)
			tpcr = nil
		}
	}
}

func BenchmarkSingleJoin(b *testing.B) {

	var person []dbmap.Person
	d := connectDBY()
	dbMap := dbmap.CreateDBMap(d)
	dbMap.AddTable("person", []string{"id"}, true, dbmap.Person{})
	dbMap.AddTable("job", []string{"id"}, true, dbmap.Job{})
	dbMap.AddTable("pacar", []string{"id"}, true, dbmap.Pacar{})
	dbMap.AddTable("shoes", []string{"id"}, true, dbmap.Shoes{})
	dbMap.AddRelation(dbmap.Person{}, dbmap.ONE_TO_ONE, dbmap.Job{}, "job", "job_id", "id")
	dbMap.AddRelation(dbmap.Person{}, dbmap.ONE_TO_MANY, dbmap.Pacar{}, "job", "id", "person_id")
	dbMap.AddRelation(dbmap.Pacar{}, dbmap.ONE_TO_MANY, dbmap.Shoes{}, "shoes", "id", "pacar_id")
	dbMap.Init()

	for n := 0; n < b.N; n++ {
		person = nil
		dbMap.Select(context.Background(), &person, "person.id=204 or person.id=203")
	}
}
