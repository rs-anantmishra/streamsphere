import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { SharedDataService } from './shared-data.service';
import { PlaylistsInfo, PlaylistsDataResponse } from '../classes/playlists';
import { VideoData, VideoDataResponse } from '../classes/video-data';
import { Observable } from 'rxjs';
import { apiBaseUrl, withContentLinks } from '../utilities/url-transforms';

const apiUrl: string = apiBaseUrl() + "/api"

@Injectable({
    providedIn: 'root'
})
export class PlaylistsService {

    constructor(private http: HttpClient, private sharedData: SharedDataService) { }

    //getAllPlaylists
    async getAllPlaylists(): Promise<PlaylistsDataResponse> {
        let url = '/homepage/playlists'

        return fetch(apiUrl + url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => response.json()).then(withContentLinks)

    }

    //getPlaylistVideos
    async getPlaylistVideos(playlistId: number): Promise<VideoDataResponse> {
        let url = '/homepage/playlists/' + playlistId

        return fetch(apiUrl + url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => response.json()).then(withContentLinks)

    }

    //download video
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
