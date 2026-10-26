import { Injectable } from '@angular/core';
import { VideoData, VideoDataRequest, QueueDownloads, VideoDataResponse } from '../classes/video-data';
import { HttpClient } from '@angular/common/http';
import { SharedDataService } from './shared-data.service';
import { environment } from '../../environments/environment';
import { Observable, repeat, retry } from 'rxjs';
import { webSocket } from 'rxjs/webSocket';
import { DownloadStatusEvent } from '../interfaces/download-status-event';
import { appBaseUrl, toWsUrl } from '../utilities/url-transforms';

const apiUrl: string = appBaseUrl(environment.baseUrl) + "/api"

@Injectable({
    providedIn: 'root'
})
export class DownloadService {

    constructor(private http: HttpClient, private sharedData: SharedDataService) { }

    // live download status; each subscription opens its own socket and reconnects after error (retry) and clean close (repeat)
    statusEvents(): Observable<DownloadStatusEvent> {
        const wsUrl = toWsUrl(appBaseUrl(environment.baseUrl)) + "/ws/status"
        return webSocket<DownloadStatusEvent>(wsUrl).pipe(retry({ delay: 2000 }), repeat({ delay: 2000 }))
    }

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