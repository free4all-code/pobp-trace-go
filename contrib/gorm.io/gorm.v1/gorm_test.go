

package gorm

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	sqltrace "git.proto.group/protoobp/pobp-trace-go/contrib/database/sql"
	"git.proto.group/protoobp/pobp-trace-go/contrib/internal/sqltest"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/mocktracer"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
	"git.proto.group/protoobp/pobp-trace-go/internal/globalconfig"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	mysqlgorm "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// tableName holds the SQL table that these tests will be run against. It must be unique cross-repo.
const (
	tableName           = "testgorm"
	pgConnString        = "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"
	sqlServerConnString = "sqlserver://sa:myPassw0rd@127.0.0.1:1433?database=master"
	mysqlConnString     = "test:test@tcp(127.0.0.1:3306)/test"
)

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
	sqlDb, err := sqltrace.Open("mysql", mysqlConnString)
	if err != nil {
		log.Fatal(err)
	}

	db, err := Open(mysqlgorm.New(mysqlgorm.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	internalDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	testConfig := &sqltest.Config{
		DB:         internalDB,
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
	sqltrace.Register("pgx", &stdlib.Driver{})
	sqlDb, err := sqltrace.Open("pgx", pgConnString)
	if err != nil {
		log.Fatal(err)
	}

	db, err := Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	internalDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	testConfig := &sqltest.Config{
		DB:         internalDB,
		DriverName: "pgx",
		TableName:  tableName,
		ExpectName: "pgx.query",
		ExpectTags: map[string]interface{}{
			ext.ServiceName: "pgx.db",
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
	sqlDb, err := sqltrace.Open("sqlserver", sqlServerConnString)
	if err != nil {
		log.Fatal(err)
	}

	db, err := Open(sqlserver.New(sqlserver.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	internalDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	testConfig := &sqltest.Config{
		DB:         internalDB,
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

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestCallbacks(t *testing.T) {
	a := assert.New(t)
	mt := mocktracer.Start()
	defer mt.Stop()

	sqltrace.Register("pgx", &stdlib.Driver{})
	sqlDb, err := sqltrace.Open("pgx", pgConnString)
	if err != nil {
		log.Fatal(err)
	}

	db, err := Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal(err)
	}

	t.Run("create", func(t *testing.T) {
		parentSpan, ctx := tracer.StartSpanFromContext(context.Background(), "http.request",
			tracer.ServiceName("fake-http-server"),
			tracer.SpanType(ext.SpanTypeWeb),
		)

		db = db.WithContext(ctx)
		var queryText string
		db.Callback().Create().After("testing").Register("query text", func(d *gorm.DB) {
			queryText = d.Statement.SQL.String()
		})
		db.Create(&Product{Code: "L1212", Price: 1000})

		parentSpan.Finish()

		spans := mt.FinishedSpans()
		a.True(len(spans) >= 2)

		span := spans[len(spans)-2]
		a.Equal("gorm.create", span.OperationName())
		a.Equal(ext.SpanTypeSQL, span.Tag(ext.SpanType))
		a.Equal(queryText, span.Tag(ext.ResourceName))
	})

	t.Run("query", func(t *testing.T) {
		parentSpan, ctx := tracer.StartSpanFromContext(context.Background(), "http.request",
			tracer.ServiceName("fake-http-server"),
			tracer.SpanType(ext.SpanTypeWeb),
		)

		db = db.WithContext(ctx)
		var queryText string
		db.Callback().Query().After("testing").Register("query text", func(d *gorm.DB) {
			queryText = d.Statement.SQL.String()
		})
		var product Product
		db.First(&product, "code = ?", "L1212")

		parentSpan.Finish()

		spans := mt.FinishedSpans()
		a.True(len(spans) >= 2)

		span := spans[len(spans)-2]
		a.Equal("gorm.query", span.OperationName())
		a.Equal(ext.SpanTypeSQL, span.Tag(ext.SpanType))
		a.Equal(queryText, span.Tag(ext.ResourceName))
	})

	t.Run("update", func(t *testing.T) {
		parentSpan, ctx := tracer.StartSpanFromContext(context.Background(), "http.request",
			tracer.ServiceName("fake-http-server"),
			tracer.SpanType(ext.SpanTypeWeb),
		)

		db = db.WithContext(ctx)
		var queryText string
		db.Callback().Update().After("testing").Register("query text", func(d *gorm.DB) {
			queryText = d.Statement.SQL.String()
		})
		var product Product
		db.First(&product, "code = ?", "L1212")
		db.Model(&product).Update("Price", 2000)

		parentSpan.Finish()

		spans := mt.FinishedSpans()
		a.True(len(spans) >= 2)

		span := spans[len(spans)-2]
		a.Equal("gorm.update", span.OperationName())
		a.Equal(ext.SpanTypeSQL, span.Tag(ext.SpanType))
		a.Equal(queryText, span.Tag(ext.ResourceName))
	})

	t.Run("delete", func(t *testing.T) {
		parentSpan, ctx := tracer.StartSpanFromContext(context.Background(), "http.request",
			tracer.ServiceName("fake-http-server"),
			tracer.SpanType(ext.SpanTypeWeb),
		)

		db = db.WithContext(ctx)
		var queryText string
		db.Callback().Delete().After("testing").Register("query text", func(d *gorm.DB) {
			queryText = d.Statement.SQL.String()
		})
		var product Product
		db.First(&product, "code = ?", "L1212")
		db.Delete(&product)

		parentSpan.Finish()

		spans := mt.FinishedSpans()
		a.True(len(spans) >= 2)

		span := spans[len(spans)-2]
		a.Equal("gorm.delete", span.OperationName())
		a.Equal(ext.SpanTypeSQL, span.Tag(ext.SpanType))
		a.Equal(queryText, span.Tag(ext.ResourceName))
	})
}

func TestAnalyticsSettings(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	sqltrace.Register("pgx", &stdlib.Driver{})
	sqlDb, err := sqltrace.Open("pgx", pgConnString)
	if err != nil {
		log.Fatal(err)
	}

	db, err := Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal(err)
	}

	assertRate := func(t *testing.T, mt mocktracer.Tracer, rate interface{}, opts ...Option) {
		parentSpan, ctx := tracer.StartSpanFromContext(context.Background(), "http.request",
			tracer.ServiceName("fake-http-server"),
			tracer.SpanType(ext.SpanTypeWeb),
		)

		db, err := Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{}, opts...)
		if err != nil {
			log.Fatal(err)
		}

		db = db.WithContext(ctx)
		db.Create(&Product{Code: "L1212", Price: 1000})

		parentSpan.Finish()

		spans := mt.FinishedSpans()
		assert.True(t, len(spans) > 2)
		s := spans[len(spans)-2]
		assert.Equal(t, rate, s.Tag(ext.EventSampleRate))
	}

	t.Run("defaults", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()

		assertRate(t, mt, nil)
	})

	t.Run("global", func(t *testing.T) {
		t.Skip("global flag disabled")
		mt := mocktracer.Start()
		defer mt.Stop()

		rate := globalconfig.AnalyticsRate()
		defer globalconfig.SetAnalyticsRate(rate)
		globalconfig.SetAnalyticsRate(0.4)

		assertRate(t, mt, 0.4)
	})

	t.Run("enabled", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()

		assertRate(t, mt, 1.0, WithAnalytics(true))
	})

	t.Run("disabled", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()

		assertRate(t, mt, nil, WithAnalytics(false))
	})

	t.Run("override", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()

		rate := globalconfig.AnalyticsRate()
		defer globalconfig.SetAnalyticsRate(rate)
		globalconfig.SetAnalyticsRate(0.4)

		assertRate(t, mt, 0.23, WithAnalyticsRate(0.23))
	})
}

func TestContext(t *testing.T) {
	sqltrace.Register("pgx", &stdlib.Driver{})
	sqlDb, err := sqltrace.Open("pgx", pgConnString)
	if err != nil {
		log.Fatal(err)
	}

	db, err := Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	t.Run("with", func(t *testing.T) {
		const contextKey = "text context"

		type key string
		testCtx := context.WithValue(context.Background(), key(contextKey), true)
		db := db.WithContext(testCtx)
		ctx := db.Statement.Context
		assert.Equal(t, testCtx.Value(key(contextKey)), ctx.Value(key(contextKey)))
	})
}

func TestError(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	assertErrCheck := func(t *testing.T, mt mocktracer.Tracer, errExist bool, opts ...Option) {
		sqltrace.Register("pgx", &stdlib.Driver{})
		sqlDb, err := sqltrace.Open("pgx", pgConnString)
		assert.Nil(t, err)

		db, err := Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{}, opts...)
		assert.Nil(t, err)
		db.AutoMigrate(&Product{})
		db.First(&Product{}, Product{Code: "L1210", Price: 2000})

		spans := mt.FinishedSpans()
		assert.True(t, len(spans) > 1)

		// Get last span (gorm.db)
		s := spans[len(spans)-1]

		assert.Equal(t, errExist, s.Tag(ext.Error) != nil)
	}

	t.Run("defaults", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()
		assertErrCheck(t, mt, true)
	})

	t.Run("errcheck", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()
		errFn := func(err error) bool {
			return err != gorm.ErrRecordNotFound
		}
		assertErrCheck(t, mt, false, WithErrorCheck(errFn))
	})
}