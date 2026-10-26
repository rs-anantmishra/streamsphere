import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { SharedDataService } from './shared-data.service';
import { StorageStatus, StorageStatusResponse } from '../classes/storage';
import { environment } from '../../environments/environment';
import { appBaseUrl } from '../utilities/url-transforms';

const apiUrl: string = appBaseUrl(environment.baseUrl) + "/api"

@Injectable({
    providedIn: 'root'
})
export class FilesService {
    constructor(private http: HttpClient, private sharedData: SharedDataService) { }

    //getStorageStatus
    async getStorageStatus(): Promise<StorageStatusResponse> {
        let url = '/storage/status'
        return fetch(apiUrl + url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => { return response.json(); })
    }
}
