import { Component, OnInit, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { SharedDataService } from '../../services/shared-data.service';
import { VideoData } from '../../classes/video-data';
import { PanelModule } from 'primeng/panel';
import { TagModule } from 'primeng/tag';
import { ChipModule } from 'primeng/chip';
import { FieldsetModule } from 'primeng/fieldset';
import { ConfirmPopup, ConfirmPopupModule } from 'primeng/confirmpopup';
import { ConfirmationService, MessageService } from 'primeng/api';
import { ToastModule } from 'primeng/toast';
import { PlaylistsInfo, SelectedPlaylist } from '../../classes/playlists';

import { MinifiedViewCount } from '../../utilities/pipes/views-conversion.pipe';
import { MinifiedLikeCount } from '../../utilities/pipes/likes-conversion.pipe';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { CommaSepStringFromArray } from "../../utilities/pipes/array-comma-sep.pipe";
import { FormattedResolutionPipe } from "../../utilities/pipes/format-resolution.pipe";
import { PlaylistsService } from '../../services/playlists.service';
import { MinifiedDatePipe } from "../../utilities/pipes/formatted-date.pipe";
import { FilesizeConversionPipe } from "../../utilities/pipes/filesize-conversion.pipe";
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { LinkifyPipe } from '../../utilities/pipes/linkify.pipe';
import { cp1252_to_utf8 } from '../../utilities/cp1252-to-utf8';
import { downloadBlob } from '../../utilities/download-blob';
import { confirmDelete } from '../../utilities/confirm-delete';
import { UrlEncode } from '../../utilities/url-encode';
import Plyr from 'plyr';

@Component({
    selector: 'app-playlist-details',
    imports: [CommonModule, RouterModule, ButtonModule, PanelModule, ScrollPanelModule, TagModule, ConfirmPopupModule, ToastModule,
        ChipModule, MinifiedViewCount, MinifiedLikeCount, CommaSepStringFromArray, FormattedResolutionPipe,
        MinifiedDatePipe, FieldsetModule, FilesizeConversionPipe, ProgressSpinnerModule, LinkifyPipe],
    providers: [ConfirmationService, MessageService, UrlEncode],
    templateUrl: './playlist-details.component.html',
    styleUrl: './playlist-details.component.scss'
})
export class PlaylistDetailsComponent implements OnInit {

    fieldsetVideoClass = ''
    isDarkMode!: boolean;
    playlist!: SelectedPlaylist;
    data!: any;
    // subscription!: Subscription;
    public player: any;
    selectedVideo: VideoData = new VideoData()
    playlistInfo!: PlaylistsInfo;
    playlistVideos: VideoData[] = [new VideoData()]
    loaded = false

    @ViewChild(ConfirmPopup) confirmPopup!: ConfirmPopup;
    constructor(private confirmationService: ConfirmationService, private messageService: MessageService, private svcSharedData: SharedDataService, 
        private svcPlaylists: PlaylistsService, private router: Router, readonly urlEncode: UrlEncode) {
    }

    async ngOnInit(): Promise<void> {
        this.isDarkMode = this.svcSharedData.getIsDarkMode()
        this.playlist = this.svcSharedData.getPlaylist();

        //nothing is saved when this page is opened without picking a playlist first
        if (!this.playlist?.info) {
            this.router.navigate(['/playlists'])
            return
        }

        await this.loadPlaylistVideos();
        //a playlist can be empty - the page still shows its title and an empty player
        if (this.playlist.video_data.length > 0) {
            this.selectActiveVideo();
            this.prepareSelectedVideo();
        }
        this.createPlayer();
        this.loaded = true
    }

    //nothing is playable until a playlist video is selected
    get hasVideo(): boolean {
        return this.selectedVideo.media_url !== ''
    }

    async loadPlaylistVideos() {
        //get playlist videos
        this.playlist.video_data = await this.getPlaylistsVideos(this.playlist.info.playlist_id);

        //sort playlist
        this.playlist.video_data.sort((a, b) => a.playlist_video_index - b.playlist_video_index)

        this.playlistVideos = this.playlist.video_data
    }

    selectActiveVideo() {
        //assign selected video
        if (this.playlist.active_video === -1) {
            //set 1st video as active 
            this.playlist.active_video = this.playlist.video_data[0].video_id
            this.selectedVideo = this.playlist.video_data[0]
        } else {
            //find active-video and assign it
            this.playlist.video_data.forEach((item) => {
                if (item.video_id === this.playlist.active_video) {
                    this.selectedVideo = item
                }
            });
        }
    }

    prepareSelectedVideo() {
        this.selectedVideo.thumbnail = this.selectedVideo.thumbnail
        this.selectedVideo.media_url = this.urlEncode.encodedUrl(this.selectedVideo.media_url)
        this.selectedVideo.webpage_url = this.urlEncode.encodedUrl(this.selectedVideo.webpage_url)

        //linkify
        this.selectedVideo.description = cp1252_to_utf8(this.selectedVideo.description)
        // this.selectedVideo.description = this.linkify(this.selectedVideo.description)
    }

    createPlayer() {
        this.player = new Plyr('#plyrId', { captions: { active: true }, keyboard: { global: true }, autoplay: true });
        // use this for loading next item in playlist or repeating a single item
        // this.player.on('ended', function() { console.log('hi there')});
    }

    changeContent(video_id: number) {
        this.player.pause()
        this.loaded = false
        this.playlist.active_video = video_id
        this.svcSharedData.setPlaylist(this.playlist)

        this.playlist.video_data.forEach((item) => {
            if (item.video_id === this.playlist.active_video) {
                this.selectedVideo = item
            }
        });
        window.location.reload()
    }

    //always an array - an empty playlist, an empty response or a failed request all give no videos
    async getPlaylistsVideos(playlistId: number): Promise<VideoData[]> {
        //check-cached
        const playlist: SelectedPlaylist = this.svcSharedData.getPlaylist();
        if (playlist.video_data?.length > 0) {
            return playlist.video_data
        }

        //get from service
        try {
            const response = await this.svcPlaylists.getPlaylistVideos(playlistId);
            if (response === null) {
                console.error('error: recieved null instead of playlist videos.')
            }
            return response?.data ?? []
        } catch (err) {
            console.error('error: failed to get playlist videos.', err)
            return []
        }
    }

    async download(): Promise<void> {
        if (!this.hasVideo) { return }
        (await this.svcPlaylists
            .download(this.selectedVideo.media_url))
            .subscribe(blob => downloadBlob(blob, this.selectedVideo.title + '.' + this.selectedVideo.extension))
    }

    getClass(video: VideoData): string {
        let result = ''
        this.isDarkMode = this.svcSharedData.getIsDarkMode()
        if (this.isDarkMode) {
            result = 'fieldset-playlist-content-dark'
        } else if (!this.isDarkMode) {
            result = 'fieldset-playlist-content-light'
        }
        //video is selected
        if (this.selectedVideo.video_id === video.video_id && !this.isDarkMode) {
            result = 'fieldset-playlist-content-light selected-video-light'
        } else if (this.selectedVideo.video_id === video.video_id && this.isDarkMode) {
            result = 'fieldset-playlist-content-dark selected-video-dark'
        }

        return result
    }

    accept() {
        this.confirmPopup.onAccept();
    }

    reject() {
        this.confirmPopup.onReject();
    }

    confirm(event: Event, contentId: number) {
        confirmDelete(event, contentId, id => this.svcPlaylists.deleteVideoById(id), this.confirmationService, this.messageService, this.router);
    }
}
