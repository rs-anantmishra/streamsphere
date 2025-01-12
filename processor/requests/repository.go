package requests

import (
	"database/sql"
	"fmt"
	"time"

	queries "github.com/rs-anantmishra/streamsphere/utils/processor/database/queries"
	domain "github.com/rs-anantmishra/streamsphere/utils/processor/domain"
)

type IRequestRepository interface {
	GetRequestsByRequestType(string) ([]domain.Request, error)

	InsertRequestStatus(requestId int, requestStatus string) int //move to api
	UpdateRequestStatus(requestId int, requestStatus string) int

	InsertRequestQueue(request domain.RequestQueue) int
	UpdateRequestQueue(requestQueueId int, field string, value string) int
}

type repository struct {
	db *sql.DB
}

func NewRepository(Database *sql.DB) IRequestRepository {
	return &repository{
		db: Database,
	}
}

func (r *repository) GetRequestsByRequestType(requestType string) ([]domain.Request, error) {

	var result []domain.Request

	rows, err := r.db.Query(queries.GetRequestsByRequestType, requestType)
	if err != nil {
		return nil, fmt.Errorf("error fetching requests: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		item := domain.Request{}
		if err := rows.Scan(&item.Id, &item.RequestUrl, &item.RequestType, &item.Metadata, &item.Thumbnail,
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

// to be moved to api
func (r *repository) InsertRequestStatus(requestId int, requestStatus string) int {

	var args []any
	args = append(args, requestId)
	args = append(args, requestStatus)
	args = append(args, time.Now().Unix())
	args = append(args, time.Now().Unix())

	requestStatusId := genericSave(*r, args, queries.InsertRequestStatus)
	return requestStatusId
}

func (r *repository) UpdateRequestStatus(requestId int, requestStatus string) int {

	var rowsAffected = -1
	requestStatusId := genericCheck(*r, requestId, "RequestStatus", queries.UpdateRequestStatusCheck)
	if requestStatusId > 0 {
		var args []any
		args = append(args, requestStatusId)
		args = append(args, requestStatus)
		args = append(args, time.Now().Unix())
		rowsAffected = genericUpdate(*r, args, queries.UpdateRequestStatus)
	}
	return rowsAffected
}

func (r *repository) InsertRequestQueue(request domain.RequestQueue) int {

	//multi-clause-check
	requestQueueId := tagsOrCategoriesCheck(*r, request.ContentId, "RequestQueue", queries.InsertRequestQueueCheck, request.RequestId)
	if requestQueueId <= 0 {
		var args []any
		args = append(args, request.RequestId)
		args = append(args, request.ContentId)
		args = append(args, request.ProcessStatus)
		args = append(args, request.RetryCount)
		args = append(args, request.Message)
		args = append(args, request.Cancelled)
		args = append(args, time.Now().Unix())
		args = append(args, time.Now().Unix())

		requestQueueId = genericSave(*r, args, queries.InsertRequestStatus)
	}
	return requestQueueId
}

func (r *repository) UpdateRequestQueue(requestQueueId int, field string, value string) int {

	var updateQuery string
	switch field {
	case "ProcessStatus":
		updateQuery = queries.UpdateProcessStatus
	case "RetryCount":
		updateQuery = queries.UpdateRetryCount
	case "Message":
		updateQuery = queries.UpdateMessage
	case "Cancelled":
		updateQuery = queries.UpdateCancelled
	}

	var rowsAffected = -1
	var args []any
	args = append(args, requestQueueId)
	args = append(args, value)
	args = append(args, time.Now().Unix())
	rowsAffected = genericUpdate(*r, args, updateQuery)

	return rowsAffected
}

// region Private Methods ////////////////////////////////////////////////////

func genericSave(r repository, args []any, genericQuery string) int {
	resultId := -1

	if resultId < 0 {

		//generic save
		result, err := r.db.Exec(genericQuery, args...)

		//check for errors
		if err != nil {
			fmt.Println("error:", err)
			return resultId
		}

		//get the inserted records Id
		var id int64
		if id, err = result.LastInsertId(); err != nil {
			return resultId
		}
		resultId = int(id)
	}

	return resultId
}

func genericUpdate(r repository, args []any, genericQuery string) int {
	rowsAffected := -1

	if rowsAffected < 0 {

		//generic save
		result, err := r.db.Exec(genericQuery, args...)

		//check for errors
		if err != nil {
			fmt.Println("error:", err)
			return rowsAffected
		}

		//get the inserted records Id
		var id int64
		if id, err = result.RowsAffected(); err != nil {
			return rowsAffected
		}
		rowsAffected = int(id)
	}

	return rowsAffected
}

func genericCheck(r repository, Id any, idContext string, genericQuery string) int {
	resultId := -1

	//Check if entry for Id exists?
	chk := r.db.QueryRow(genericQuery, Id)
	if err := chk.Scan(&resultId); err == sql.ErrNoRows {
		fmt.Println("No Rows found for", idContext, "Id", Id)
	}

	return resultId
}

func tagsOrCategoriesCheck(r repository, Id any, idContext string, genericQuery string, videoId int) int {
	resultId := -1

	//Check if entry for Id exists?
	chk := r.db.QueryRow(genericQuery, Id, videoId)
	if err := chk.Scan(&resultId); err == sql.ErrNoRows {
		fmt.Println("No Rows found for", idContext, "Id", Id)
	}

	return resultId
}

// endregion Private Methods ////////////////////////////////////////////////////
