package orm

import (
	"database/sql"
	"fmt"
)

// BeginTx starts a database transaction using the Datastore's adapter.
func BeginTx(ds *Datastore) (*sql.Tx, error) {
	if ds == nil || ds.DB == nil || ds.Adapter == nil {
		return nil, fmt.Errorf("datastore, database connection, or adapter was nil")
	}
	return ds.Adapter.BeginTx(ds.DB)
}
