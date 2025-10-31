package orm

import "database/sql"

// Datastore holds the database connection and the database adapter.
type Datastore struct {
	DB      *sql.DB
	Adapter DBAdapter
}
