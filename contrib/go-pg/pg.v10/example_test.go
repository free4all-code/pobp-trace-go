

package pg_test

import (
	"log"

	pgtrace "git.proto.group/protoobp/pobp-trace-go/contrib/go-pg/pg.v10"

	"github.com/go-pg/pg/v10"
)

func Example() {
	conn := pg.Connect(&pg.Options{
		User:     "go-pg-test",
		Database: "protoobp",
	})

	// Wrap the connection with the APM hook.
	pgtrace.Wrap(conn)
	var user struct {
		Name string
	}
	_, err := conn.QueryOne(&user, "SELECT name FROM users")
	if err != nil {
		log.Fatal(err)
	}
}
