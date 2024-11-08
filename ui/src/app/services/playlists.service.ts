import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { SharedDataService } from './shared-data.service';
import { PlaylistsInfo, PlaylistsDataResponse } from '../classes/playlists';
import { VideoData, VideoDataResponse } from '../classes/video-data';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

const apiUrl: string = environment.baseUrl + "/api"

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
        }).then(response => { return response.json(); })

    }

    //getPlaylistVideos
    async getPlaylistVideos(playlistId: number): Promise<VideoDataResponse> {
        let url = '/homepage/playlists/' + playlistId

        return fetch(apiUrl + url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => { return response.json(); })

    }

    //download video
    async download(url: string): Promise<Observable<Blob>> {
        return this.http.get(url, {
            responseType: 'blob'
        })
    }
}
