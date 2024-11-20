import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { SharedDataService } from './shared-data.service';
import { VideoData, VideoDataRequest, QueueDownloads, VideoDataResponse } from '../classes/video-data';
import { BehaviorSubject, Observable, Subscription } from 'rxjs'
import { ContentSearch, ContentSearchResponse } from '../classes/search';
import { environment } from '../../environments/environment';

const apiUrl: string = environment.baseUrl + "/api"

@Injectable({
    providedIn: 'root'
})
export class VideosService {

    constructor(private http: HttpClient, private sharedData: SharedDataService) { }

    //getAllVideos
    async getAllVideos(): Promise<VideoData[]> {
        let url = '/homepage/videos'

        return fetch(apiUrl + url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => { return response.json(); })

    }

    //getContentSearchInfo
    async getContentSearchInfo(): Promise<ContentSearchResponse> {
        let url = '/search/info'

        return fetch(apiUrl + url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => { return response.json(); })
    }

    //getPlaylistVideos
    async getContentById(contentId: number): Promise<VideoDataResponse> {
        let url = '/homepage/video/' + contentId

        return fetch(apiUrl + url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Access-Control': 'Allow-Origin'
            }
        }).then(response => { return response.json(); })

    }

    async download(url: string): Promise<Observable<Blob>> {
        return this.http.get(url, {
            responseType: 'blob'
        })
    }

    //getAllVideos
    async deleteVideoById(contentId: number): Promise<boolean> {
        let url = '/homepage/video/' + contentId + '/delete'

        return fetch(apiUrl + url, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => { return response.json(); })
    }

}
