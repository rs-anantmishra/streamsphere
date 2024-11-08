import { Injectable } from '@angular/core';
import { VideoData, VideoDataRequest, QueueDownloads, VideoDataResponse } from '../classes/video-data';
import { HttpClient } from '@angular/common/http';
import { SharedDataService } from './shared-data.service';
import { environment } from '../../environments/environment';

const apiUrl: string = environment.baseUrl + "/api"

@Injectable({
    providedIn: 'root'
})
export class DownloadService {

    constructor(private http: HttpClient, private sharedData: SharedDataService) { }

    //metadata
    async getMetadata(request: VideoDataRequest): Promise<VideoDataResponse> {

        let url = '/download/metadata'

        return fetch(apiUrl + url, {
            method: 'POST',
            body: JSON.stringify(request),
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => { return response.json(); });
    }

    //media
    async getMedia(request: QueueDownloads): Promise<string> {
        let url = '/download/media'

        return fetch(apiUrl + url, {
            method: 'POST',
            body: JSON.stringify(request),
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => { return response.json(); });
    }

    //queued-items
    async getQueuedItems(state: string): Promise<VideoData[]> {
        let url = '/download/queued-items'
        let queryParams = '?state=' + state
        return fetch(apiUrl + url + queryParams, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
        }).then(response => { return response.json(); });
    }

}