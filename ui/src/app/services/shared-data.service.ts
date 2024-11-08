import { Injectable, WritableSignal, signal } from '@angular/core';
import { VideoData } from '../classes/video-data';
import { PlaylistsDataResponse, PlaylistsInfo, SelectedPlaylist } from '../classes/playlists';
import { BehaviorSubject, Observable, Subject, Subscription } from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class SharedDataService {

    constructor() { }
    isDarkMode: boolean = false;
    lstVideos: VideoData[] = [];
    lstPlaylists: PlaylistsInfo[] = [];
    isDownloadActive: boolean = false;
    isPlaylist: boolean = false;
    playlist!: any;
    queuedItemsMetadata: VideoData[] = [];
    activeDownloadMetadata: VideoData[] = [];
    videosPageSizeCount: number = -1;
    playlistsPageSizeCount: number = -1;
    activePlayerMetadata: VideoData = new VideoData();

    //add or remove from queuedItems
    setQueuedItemsMetadata(metadata: VideoData[], ops: Operation) {
        if (Operation.Insert == ops) {
            let value = this.getQueuedItemsMetadata();
            value.push(...metadata)
            localStorage.setItem('queuedItemsMetadata', JSON.stringify(value));
            this.queuedItemsMetadata = value
        } else if (Operation.RemoveAtIndexZero == ops) {
            let deleted = metadata.splice(0, 1)
            localStorage.setItem('queuedItemsMetadata', JSON.stringify(metadata));
            this.queuedItemsMetadata = metadata
        } else if (Operation.Replace == ops) {
            localStorage.setItem('queuedItemsMetadata', JSON.stringify(metadata));
        }
    }

    //getQueuedItems
    getQueuedItemsMetadata() {
        let stringResult = localStorage.getItem('queuedItemsMetadata') !== null ? localStorage.getItem('queuedItemsMetadata') : JSON.stringify([])
        let queuedItemsMeta = stringResult === null ? [new VideoData()] : JSON.parse(stringResult);
        this.queuedItemsMetadata = queuedItemsMeta;

        return this.queuedItemsMetadata
    }

    setIsDownloadActive(value: boolean) {
        localStorage.setItem('isDownloadActive', JSON.stringify(value));
    }

    getIsDownloadActive(): boolean {
        let stringResult = localStorage.getItem('isDownloadActive') !== null ? localStorage.getItem('isDownloadActive') : 'false'
        let isActive = stringResult === null ? false : JSON.parse(stringResult);
        this.isDownloadActive = isActive;

        return this.isDownloadActive
    }

    setIsDarkMode(value: boolean) {
        localStorage.setItem('isDarkMode', JSON.stringify(value));
    }

    getIsDarkMode(): boolean {
        let stringResult = localStorage.getItem('isDarkMode') !== null ? localStorage.getItem('isDarkMode') : 'false'
        let isActive = stringResult === null ? false : JSON.parse(stringResult);
        this.isDarkMode = isActive;

        return this.isDarkMode
    }

    setIsPlaylist(value: boolean) {
        localStorage.setItem('isPlaylist', JSON.stringify(value));
    }

    getIsPlaylist(): boolean {
        let stringResult = localStorage.getItem('isPlaylist') !== null ? localStorage.getItem('isPlaylist') : 'false'
        let isPl = stringResult === null ? false : JSON.parse(stringResult);
        this.isPlaylist = isPl;

        return this.isPlaylist
    }

    setActiveDownloadMetadata(value: any) {
        localStorage.setItem('activeDownloadMetadata', JSON.stringify(value));
    }

    getActiveDownloadMetadata() {
        let stringResult = localStorage.getItem('activeDownloadMetadata') !== null ? localStorage.getItem('activeDownloadMetadata') : JSON.stringify([])
        let activeDownloadMeta = stringResult === null ? [new VideoData()] : JSON.parse(stringResult);
        this.activeDownloadMetadata = activeDownloadMeta;

        return this.activeDownloadMetadata
    }

    setlstVideos(value: any) {
        localStorage.setItem('lstVideos', JSON.stringify(value));
    }

    getlstVideos() {
        let stringResult = localStorage.getItem('lstVideos') !== null ? localStorage.getItem('lstVideos') : JSON.stringify([])
        let lstVideosData = stringResult === null ? [new VideoData()] : JSON.parse(stringResult);
        this.lstVideos = lstVideosData;

        return this.lstVideos
    }

    setVideosPageSizeCount(value: any) {
        localStorage.setItem('videosPageSizeCount', JSON.stringify(value));
    }

    getVideosPageSizeCount() {
        let stringResult = localStorage.getItem('videosPageSizeCount') !== null ? localStorage.getItem('videosPageSizeCount') : JSON.stringify('-1')
        let pageSizeCount = stringResult === null ? -1 : JSON.parse(stringResult);
        this.videosPageSizeCount = pageSizeCount;

        return this.videosPageSizeCount
    }

    setlstPlaylists(value: any) {
        localStorage.setItem('lstPlaylists', JSON.stringify(value));
    }

    getlstPlaylists() {
        let stringResult = localStorage.getItem('lstPlaylists') !== null ? localStorage.getItem('lstPlaylists') : JSON.stringify([])
        let lstPlaylistsData = stringResult === null ? [new VideoData()] : JSON.parse(stringResult);
        this.lstPlaylists = lstPlaylistsData;

        return this.lstPlaylists
    }

    setPlaylistsPageSizeCount(value: any) {
        localStorage.setItem('playlistsPageSizeCount', JSON.stringify(value));
    }

    getPlaylistsPageSizeCount() {
        let stringResult = localStorage.getItem('playlistsPageSizeCount') !== null ? localStorage.getItem('playlistsPageSizeCount') : JSON.stringify('-1')
        let pageSizeCount = stringResult === null ? -1 : JSON.parse(stringResult);
        this.playlistsPageSizeCount = pageSizeCount;

        return this.playlistsPageSizeCount
    }

    setActivePlayerMetadata(value: any) {
        localStorage.setItem('activePlayerMetadata', JSON.stringify(value));
    }

    getActivePlayerMetadata() {
        let stringResult = localStorage.getItem('activePlayerMetadata') !== null ? localStorage.getItem('activePlayerMetadata') : JSON.stringify(new VideoData())
        let activePlayerMeta = stringResult === null ? new VideoData() : JSON.parse(stringResult);
        this.activePlayerMetadata = activePlayerMeta;

        return this.activePlayerMetadata
    }

    setPlaylist(value: SelectedPlaylist) {
        localStorage.setItem('playlist', JSON.stringify(value));
    }

    getPlaylist(): SelectedPlaylist {
        let stringResult = localStorage.getItem('playlist') !== null ? localStorage.getItem('playlist') : '{}'
        let playlist = stringResult === null ? false : JSON.parse(stringResult);
        this.playlist = playlist;

        return this.playlist
    }

    private playVideo: BehaviorSubject<VideoData> = new BehaviorSubject(new VideoData());
    onPlayVideoChange(): Observable<VideoData> {
        //check localstorage
        let activeVideo = this.getActivePlayerMetadata()
        if (activeVideo.media_url != '') {
            this.setPlayVideo(activeVideo)
        }
        return this.playVideo.asObservable();
    }

    setPlayVideo(data: VideoData): void {
        this.setActivePlayerMetadata(data);
        this.playVideo.next(data);
    }

    resetPlayVideo(): void {
        this.playVideo.next(new VideoData());
    }

    private pageSizeCount: BehaviorSubject<number> = new BehaviorSubject(-1)
    setPageSizeCount(count: number): void {
        this.setVideosPageSizeCount(count);
        this.pageSizeCount.next(count);
    }

    getPageSizeCount(): Observable<number> {
        return this.pageSizeCount.asObservable()
    }

    private pageSizeCountPlaylist: BehaviorSubject<number> = new BehaviorSubject(-1)
    setPageSizeCountPlaylist(count: number): void {
        this.setPlaylistsPageSizeCount(count);
        this.pageSizeCountPlaylist.next(count);
    }

    getPageSizeCountPlaylist(): Observable<number> {
        return this.pageSizeCountPlaylist.asObservable()
    }

    private refreshAutoCompleteSubject = new BehaviorSubject(false);
    public _refreshAutoComplete$ = this.refreshAutoCompleteSubject.asObservable()
    setRefreshAutoCompleteValue(value: boolean) {
        this.refreshAutoCompleteSubject.next(value)
    }

}


export enum Operation {
    Insert = 1,
    RemoveAtIndexZero = 2,
    Replace = 3
}