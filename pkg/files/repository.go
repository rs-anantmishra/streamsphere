package files

import (
	"database/sql"
	"fmt"

	q "github.com/rs-anantmishra/metubeplus/database/queries"
)

type IRepository interface {
	DBStorageStatusInfo() (int64, error)
}

type repository struct {
	db *sql.DB
}

func NewFilesRepo(database *sql.DB) IRepository {
	return &repository{
		db: database,
	}
}

func (r *repository) DBStorageStatusInfo() (int64, error) {

	var filesize int64
	// Query for a value based on a single row.
	if err := r.db.QueryRow(q.GetStorageUsedInfo).Scan(&filesize); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("error: %v", err)
		}
		return 0, fmt.Errorf("error: %v", err)
	}
	return filesize, nil
}
