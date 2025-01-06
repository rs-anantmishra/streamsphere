package requests

import (
	"database/sql"
	"fmt"

	queries "github.com/rs-anantmishra/streamsphere/utils/processor/database/queries"
	domain "github.com/rs-anantmishra/streamsphere/utils/processor/domain"
)

type IRequestRepository interface {
	GetRequestsByRequestType(string) ([]domain.RequestWithStatusId, error)
	UpdateRequestStatus()

	GetRequestQueue()
	CreateRequestQueue()
	UpdateRequestQueue()
}

type repository struct {
	db *sql.DB
}

func NewRepository(Database *sql.DB) IRequestRepository {
	return &repository{
		db: Database,
	}
}

func (r *repository) GetRequestsByRequestType(requestType string) ([]domain.RequestWithStatusId, error) {

	var result []domain.RequestWithStatusId

	rows, err := r.db.Query(queries.GetRequestsByRequestType, requestType)
	if err != nil {
		return nil, fmt.Errorf("error fetching requests: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		item := domain.RequestWithStatusId{}
		if err := rows.Scan(&item.RequestStatusId, &item.Id, &item.RequestUrl, &item.RequestType, &item.Metadata, &item.Thumbnail,
			&item.Content, &item.ContentFormat, &item.Subtitles, &item.SubtitlesLanguage, &item.IsProxied, &item.Proxy,
			&item.Scheduled, &item.CreatedDate, &item.ModifiedDate); err != nil {
			if err == sql.ErrNoRows {
				return result, fmt.Errorf("no requests found for requests status: %q", requestType)
			}
			return result, fmt.Errorf("unhandled error while executing get-requests-by-request-type for request type: %q: %v", requestType, err)
		}
		result = append(result, item)
	}

	return result, err
}

func (r *repository) UpdateRequestStatus() {

}

func (r *repository) GetRequestQueue() {

}

func (r *repository) CreateRequestQueue() {

}

func (r *repository) UpdateRequestQueue() {

	//update status
	//update cancelled
	//update retry-count

}
