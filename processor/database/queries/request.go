package queries

const GetRequestsByRequestType string = `Select RS.Id 'RequestStatusId', R.Id, R.RequestUrl, R.RequestType, R.Metadata, R.Thumbnail, R.Content, R.ContentFormat, R.Subtitles, R.SubtitlesLanguage, R.IsProxied, R.Proxy,  R.Scheduled, R.CreatedDate, IFNULL(R.ModifiedDate, 0)
FROM tblRequestStatus RS
INNER JOIN tblRequests R ON R.Id = RS.RequestId AND RS.RequestStatus = ?`
