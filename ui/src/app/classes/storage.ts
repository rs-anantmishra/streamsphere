export class StorageStatusResponse {
    data: StorageStatus = new StorageStatus();
    message: string = '';
    status: string = '';
}

export class StorageStatus {
    storage_used_db: number = 0;
    storage_used_fs: number = 0;
}
