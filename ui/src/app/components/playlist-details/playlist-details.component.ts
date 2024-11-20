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
import Plyr from 'plyr';

@Component({
    selector: 'app-playlist-details',
    standalone: true,
    imports: [CommonModule, RouterModule, ButtonModule, PanelModule, ScrollPanelModule, TagModule, ConfirmPopupModule, ToastModule,
        ChipModule, MinifiedViewCount, MinifiedLikeCount, CommaSepStringFromArray, FormattedResolutionPipe,
        MinifiedDatePipe, FieldsetModule, FilesizeConversionPipe, ProgressSpinnerModule, LinkifyPipe],
    providers: [Router,  ConfirmationService, MessageService],
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
    constructor(private confirmationService: ConfirmationService, private messageService: MessageService, private svcSharedData: SharedDataService, private svcPlaylists: PlaylistsService, private router: Router) {
    }

    async ngOnInit(): Promise<void> {
        this.isDarkMode = this.svcSharedData.getIsDarkMode()
        this.playlist = this.svcSharedData.getPlaylist();
        
        //get playlist videos
        this.playlist.video_data = await this.getPlaylistsVideos(this.playlist.info.playlist_id);

        //sort playlist
        this.playlist.video_data.sort((a, b) => a.playlist_video_index - b.playlist_video_index)

        //assign selected video
        this.playlistVideos = this.playlist.video_data
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
        
        this.selectedVideo.media_url = this.selectedVideo.media_url.replaceAll('#', '%23')
        this.selectedVideo.thumbnail = this.selectedVideo.thumbnail.replaceAll('#', '%23')
        this.selectedVideo.webpage_url = this.selectedVideo.webpage_url.replaceAll('#', '%23')

        //linkify
        this.selectedVideo.description = this.cp1252_to_utf8(this.selectedVideo.description)
        // this.selectedVideo.description = this.linkify(this.selectedVideo.description)
        this.player = new Plyr('#plyrId', { captions: { active: true }, keyboard: { global: true }, autoplay: true });
        // use this for loading next item in playlist or repeating a single item
        // this.player.on('ended', function() { console.log('hi there')});
        this.loaded = true
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

    async getPlaylistsVideos(playlistId: number): Promise<VideoData[]> {
        let result!: VideoData[];
        //check-cached
        let playlist: SelectedPlaylist = await this.svcSharedData.getPlaylist();
        if (playlist.video_data.length > 0) {
            result = playlist.video_data
        } else {
            //get from service
            let response = await this.svcPlaylists.getPlaylistVideos(playlistId);
            if (response !== null && response.data.length > 0) {
                result = response.data
            } else if (response === null) {
                console.error('error: recieved null instead of playlist videos.')
            }
        }
        return result
    }

    async download(): Promise<void> {
        (await this.svcPlaylists
            .download(this.selectedVideo.media_url))
            .subscribe(blob => {
                const a = document.createElement('a')
                const objectUrl = URL.createObjectURL(blob)
                a.href = objectUrl
                a.download = (this.selectedVideo.title + '.' + this.selectedVideo.extension);
                a.click();
                URL.revokeObjectURL(objectUrl);
            })
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
        this.confirmPopup.accept();
    }

    reject() {
        this.confirmPopup.reject();
    }

    confirm(event: Event, contentId: number) {
        this.confirmationService.confirm({
            target: event.target as EventTarget,
            message: 'Are you sure you want to delete this content?',
            accept: () => {
                this.svcPlaylists.deleteVideoById(contentId).then((result: any) => {
                    if (result.data) {
                        this.messageService.add({ severity: 'info', summary: 'Confirmed', detail: 'Content Deleted', life: 3000 });
                        setTimeout(() => {
                            this.router.navigate(['/home'])
                        }, 3000);
                    } else {
                        this.messageService.add({ severity: 'error', summary: 'Failed', detail: 'Deleted Failed', life: 3000 });
                    }
                });
            },
            reject: () => {
                this.messageService.add({ severity: 'error', summary: 'Rejected', detail: 'Delete Cancelled', life: 3000 });
            }
        });
    }

    linkify(text: string) {
        var urlRegex = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
        return text.replace(urlRegex, function (url) {
            return '<a href="' + url + '">' + url + '</a>';
        });
    }

    cp1252_to_utf8(text: string) {
        let chars_map: IObjectKeys = {
            "21": "!", "22": '"', "23": "#", "24": "$", "25": "%", "26": "&", "27": "'", "28": "(", "29": ")", "2a": "*", "2b": "+", "2c": ",",
            "2d": "-", "2e": ".", "2f": "/", "30": "0", "31": "1", "32": "2", "33": "3", "34": "4", "35": "5", "36": "6", "37": "7", "38": "8",
            "39": "9", "3a": ":", "3b": ";", "3c": "<", "3d": "=", "3e": ">", "3f": "?", "40": "@", "41": "A", "42": "B", "43": "C", "44": "D",
            "45": "E", "46": "F", "47": "G", "48": "H", "49": "I", "4a": "J", "4b": "K", "4c": "L", "4d": "M", "4e": "N", "4f": "O", "50": "P",
            "51": "Q", "52": "R", "53": "S", "54": "T", "55": "U", "56": "V", "57": "W", "58": "X", "59": "Y", "5a": "Z", "5b": "[", "5c": "\\",
            "5d": "]", "5e": "^", "5f": "_", "60": "`", "61": "a", "62": "b", "63": "c", "64": "d", "65": "e", "66": "f", "67": "g", "68": "h",
            "69": "i", "6a": "j", "6b": "k", "6c": "l", "6d": "m", "6e": "n", "6f": "o", "70": "p", "71": "q", "72": "r", "73": "s", "74": "t",
            "75": "u", "76": "v", "77": "w", "78": "x", "79": "y", "7a": "z", "7b": "{", "7c": "|", "7d": "}", "7e": "~", "a1": "¡", "a2": "¢",
            "a3": "£", "a4": "¤", "a5": "¥", "a6": "¦", "a7": "§", "a8": "¨", "a9": "©", "aa": "ª", "ab": "«", "ac": "¬", "a0": " ", "ae": "®",
            "af": "¯", "ad": " ", "b0": "°", "b1": "±", "b2": "²", "b3": "³", "b4": "´", "b5": "µ", "b6": "¶", "b7": "·", "b8": "¸", "b9": "¹",
            "ba": "º", "bb": "»", "bc": "¼", "bd": "½", "be": "¾", "bf": "¿", "c0": "À", "c1": "Á", "c2": "Â", "c3": "Ã", "c4": "Ä", "c5": "Å",
            "c6": "Æ", "c7": "Ç", "c8": "È", "c9": "É", "ca": "Ê", "cb": "Ë", "cc": "Ì", "cd": "Í", "ce": "Î", "cf": "Ï", "d0": "Ð", "d1": "Ñ",
            "d2": "Ò", "d3": "Ó", "d4": "Ô", "d5": "Õ", "d6": "Ö", "d7": "×", "d8": "Ø", "d9": "Ù", "da": "Ú", "db": "Û", "dc": "Ü", "dd": "Ý",
            "de": "Þ", "df": "ß", "e0": "à", "e1": "á", "e2": "â", "e3": "ã", "e4": "ä", "e5": "å", "e6": "æ", "e7": "ç", "e8": "è", "e9": "é",
            "ea": "ê", "eb": "ë", "ec": "ì", "ed": "í", "ee": "î", "ef": "ï", "f0": "ð", "f1": "ñ", "f2": "ò", "f3": "ó", "f4": "ô", "f5": "õ",
            "f6": "ö", "f7": "÷", "f8": "ø", "f9": "ù", "fa": "ú", "fb": "û", "fc": "ü", "fd": "ý", "fe": "þ", "ff": "ÿ", "91": "‘", "92": "’",
            "80": "€", "83": "ƒ", "85": "…", "86": "†", "87": "‡", "88": "ˆ", "89": "‰", "8a": "Š", "8b": "‹", "8c": "Œ", "8e": "Ž", "93": "“",
            "94": "”", "95": "•", "96": "–", "97": "—", "98": "˜", "99": "™", "9a": "š", "9b": "›", "9c": "œ", "9e": "ž", "9f": "Ÿ"
        };
        return text.replace(/\\x([0-9abcdef]{2})/ig, function (match, code) {
            return chars_map[code as keyof IObjectKeys];
        });
    }
}

interface IObjectKeys {
    [key: string]: string | string;
}
