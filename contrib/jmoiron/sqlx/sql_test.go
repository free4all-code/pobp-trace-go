

package sqlx

import (
	"fmt"
	"log"
	"os"
	"testing"

	sqltrace "git.proto.group/protoobp/pobp-trace-go/contrib/database/sql"
	"git.proto.group/protoobp/pobp-trace-go/contrib/internal/sqltest"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
)

// tableName holds the SQL table that these tests will be run against. It must be unique cross-repo.
const tableName = "testsqlx"

func TestMain(m *testing.M) {
	_, ok := os.LookupEnv("INTEGRATION")
	if !ok {
		fmt.Println("--- SKIP: to enable integration test, set the INTEGRATION environment variable")
		os.Exit(0)
	}
	defer sqltest.Prepare(tableName)()
	os.Exit(m.Run())
}

func TestMySQL(t *testing.T) {
	sqltrace.Register("mysql", &mysql.MySQLDriver{}, sqltrace.WithServiceName("mysql-test"))
	dbx, err := Open("mysql", "test:test@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	testConfig := &sqltest.Config{
		DB:         dbx.DB,
		DriverName: "mysql",
		TableName:  tableName,
		ExpectName: "mysql.query",
		ExpectTags: map[string]interface{}{
			ext.ServiceName: "mysql-test",
			ext.SpanType:    ext.SpanTypeSQL,
			ext.TargetHost:  "127.0.0.1",
			ext.TargetPort:  "3306",
			ext.DBUser:      "test",
			ext.DBName:      "test",
		},
	}
	sqltest.RunAll(t, testConfig)
}

func TestPostgres(t *testing.T) {
	sqltrace.Register("postgres", &pq.Driver{})
	dbx, err := Open("postgres", "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	testConfig := &sqltest.Config{
		DB:         dbx.DB,
		DriverName: "postgres",
		TableName:  tableName,
		ExpectName: "postgres.query",
		ExpectTags: map[string]interface{}{
			ext.ServiceName: "postgres.db",
			ext.SpanType:    ext.SpanTypeSQL,
			ext.TargetHost:  "127.0.0.1",
			ext.TargetPort:  "5432",
			ext.DBUser:      "postgres",
			ext.DBName:      "postgres",
		},
	}
	sqltest.RunAll(t, testConfig)
}

func TestSQLServer(t *testing.T) {
	sqltrace.Register("sqlserver", &mssql.Driver{})
	dbx, err := Open("sqlserver", "sqlserver://sa:myPassw0rd@127.0.0.1:1433?database=master")
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	testConfig := &sqltest.Config{
		DB:         dbx.DB,
		DriverName: "sqlserver",
		TableName:  tableName,
		ExpectName: "sqlserver.query",
		ExpectTags: map[string]interface{}{
			ext.ServiceName: "sqlserver.db",
			ext.SpanType:    ext.SpanTypeSQL,
			ext.TargetHost:  "127.0.0.1",
			ext.TargetPort:  "1433",
			ext.DBUser:      "sa",
			ext.DBName:      "master",
		},
	}
	sqltest.RunAll(t, testConfig)
}

func TestOpenWithOptions(t *testing.T) {
	sqltrace.Register("mysql", &mysql.MySQLDriver{}, sqltrace.WithServiceName("mysql-test"))
	dbx, err := Open("mysql", "test:test@tcp(127.0.0.1:3306)/test", sqltrace.WithServiceName("other-service"))
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	testConfig := &sqltest.Config{
		DB:         dbx.DB,
		DriverName: "mysql",
		TableName:  tableName,
		ExpectName: "mysql.query",
		ExpectTags: map[string]interface{}{
			ext.ServiceName: "other-service",
			ext.SpanType:    ext.SpanTypeSQL,
			ext.TargetHost:  "127.0.0.1",
			ext.TargetPort:  "3306",
			ext.DBUser:      "test",
			ext.DBName:      "test",
		},
	}
	sqltest.RunAll(t, testConfig)
}