package presenter

type StorageStatus struct {
	StorageUsedDB int64 `json:"storage_used_db"` //this is sum of filesizes from db - in bytes
	StorageUsedFS int64 `json:"storage_used_fs"` //this size of media directory from fs - in bytes
}
