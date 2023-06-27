

package gorm_test

import (
	"context"
	"log"

	sqltrace "git.proto.group/protoobp/pobp-trace-go/contrib/database/sql"
	gormtrace "git.proto.group/protoobp/pobp-trace-go/contrib/gorm.io/gorm.v1"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"

	"github.com/jackc/pgx/v4/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

func ExampleOpen() {
	// Register augments the provided driver with tracing, enabling it to be loaded by gormtrace.Open.
	sqltrace.Register("pgx", &stdlib.Driver{}, sqltrace.WithServiceName("my-service"))
	sqlDb, err := sqltrace.Open("pgx", "postgres://pqgotest:password@localhost/pqgotest?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db, err := gormtrace.Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	var user User

	// All calls through gorm.DB are now traced.
	db.Where("name = ?", "jinzhu").First(&user)
}

func ExampleContext() {
	// Register augments the provided driver with tracing, enabling it to be loaded by gormtrace.Open.
	sqltrace.Register("pgx", &stdlib.Driver{}, sqltrace.WithServiceName("my-service"))
	sqlDb, err := sqltrace.Open("pgx", "postgres://pqgotest:password@localhost/pqgotest?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db, err := gormtrace.Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	var user User

	// Create a root span, giving name, server and resource.
	span, ctx := tracer.StartSpanFromContext(context.Background(), "my-query",
		tracer.SpanType(ext.SpanTypeSQL),
		tracer.ServiceName("my-db"),
		tracer.ResourceName("initial-access"),
	)
	defer span.Finish()

	// Subsequent spans inherit their parent from context.
	db.WithContext(ctx).Where("name = ?", "jinzhu").First(&user)
}
