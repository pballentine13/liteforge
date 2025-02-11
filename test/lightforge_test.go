package lightforge

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pballentine13/liteforge"
)

func TestOpenDB_Unit(t *testing.T) {
	testCases := []struct {
		name        string
		config      liteforge.Config
		expectedErr bool
	}{
		{
			name: "In-Memory Database",
			config: liteforge.Config{
				DriverName:     "sqlite3",
				DataSourceName: ":memory:", // In-memory database for testing
			},
			expectedErr: false,
		},
		{
			name: "File-Based Database",
			config: liteforge.Config{
				DriverName:     "sqlite3",
				DataSourceName: "test.db", // File-based database
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up the file-based database after the test (if applicable)
			if tc.name == "File-Based Database" {
				defer os.Remove("test.db") // Remove test.db after the test
			}

			db, err := liteforge.OpenDB(tc.config)

			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				if db != nil {
					db.Close() // Close the connection after the test
				}
			}
		})
	}
}
