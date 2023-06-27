

package gorm_test

import (
	"log"

	"github.com/lib/pq"
	"gopkg.in/jinzhu/gorm.v1"

	sqltrace "git.proto.group/protoobp/pobp-trace-go/contrib/database/sql"
	gormtrace "git.proto.group/protoobp/pobp-trace-go/contrib/gopkg.in/jinzhu/gorm.v1"
)

func ExampleOpen() {
	// Register augments the provided driver with tracing, enabling it to be loaded by gormtrace.Open.
	sqltrace.Register("postgres", &pq.Driver{}, sqltrace.WithServiceName("my-service"))

	// Open the registered driver, allowing all uses of the returned *gorm.DB to be traced, with the specified service name.
	db, err := gormtrace.Open("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=disable", gormtrace.WithServiceName("my-gorm-service"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	user := struct {
		gorm.Model
		Name string
	}{}

	// All calls through gorm.DB are now traced.
	db.Where("name = ?", "jinzhu").First(&user)
}
