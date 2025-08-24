
import { Component, OnInit, OnDestroy, effect } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { InputGroupModule } from 'primeng/inputgroup';
import { InputGroupAddonModule } from 'primeng/inputgroupaddon';
import { InputTextModule } from 'primeng/inputtext';
import { ButtonModule } from 'primeng/button';
import { CommonModule } from '@angular/common';
import { CheckboxModule } from 'primeng/checkbox';
import { PanelModule } from 'primeng/panel';
import { CardModule } from 'primeng/card';
import { SidebarModule } from 'primeng/sidebar';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { ProgressBarModule } from 'primeng/progressbar';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { FieldsetModule } from 'primeng/fieldset';
import { DividerModule } from 'primeng/divider';
import { OverlayPanelModule } from 'primeng/overlaypanel';
import { webSocket } from 'rxjs/webSocket'
import { environment } from '../../../environments/environment';
import { UrlEncode } from '../../utilities/url-encode';


//Services & Classes
import { QueueDownloads, DownloadMedia, VideoData, VideoDataRequest } from '../../classes/video-data';
import { DownloadService } from '../../services/download.service';
import { SimplecardComponent } from "../simplecard/simplecard.component";
import { SharedDataService, Operation } from '../../services/shared-data.service';
import { Messages, Severity } from '../../constants/messages'
import { RemovePrefixPipe } from '../../utilities/pipes/remove-prefix.pipe'
import { FilesService } from '../../services/files.service';
import { FilesizeConversionPipe } from '../../utilities/pipes/filesize-conversion.pipe'


const apiUrl: string = environment.baseUrl
interface ExtractionOptions {
    Identifier: string;
    GetAudioOnly: boolean;
    GetSubs: boolean;
}


@Component({
    selector: 'app-downloads',
    standalone: true,
    imports: [ToastModule, ProgressBarModule, FieldsetModule, ProgressSpinnerModule, SidebarModule, CardModule, FormsModule,
        InputGroupModule, InputGroupAddonModule, InputTextModule, ButtonModule, CommonModule, CheckboxModule, PanelModule,
        SimplecardComponent, RemovePrefixPipe, ScrollPanelModule, DividerModule, OverlayPanelModule, FilesizeConversionPipe],
    providers: [DownloadService, MessageService, Messages, FilesService, UrlEncode],
    templateUrl: './downloads.component.html',
    styleUrl: './downloads.component.scss'
})
export class DownloadsComponent implements OnInit {

    wsApiURL: string = this.getWebSocketApiUrl()
    wsMessage: string;
    serverLogs: string;
    queuedItems: VideoData[] = []
    nilMetadata = new VideoData()

    loading: boolean = false;
    urlInputDisabled: boolean = false;

    sock = webSocket(this.wsApiURL)

    activeDLImage = ''
    activeDLChannel = ''
    activeDLTitle = ''
    activeDownload: VideoData = new VideoData()

    //css-classes
    homeBoxActive = 'home-box'
    contentBoxActive = 'content-box'

    urlPlaceholder = 'Video or Playlist URL'
    sidebarVisible: boolean = false;
    options: ExtractionOptions = { Identifier: '', GetAudioOnly: false, GetSubs: false }
    isPlaylist: boolean = false

    constructor(private messageService: MessageService,
        private svcDownload: DownloadService,
        private svcFiles: FilesService,
        private sharedData: SharedDataService,
        readonly urlEncode: UrlEncode,
        private msg: Messages) {

        this.wsMessage = msg.wsMessage
        this.serverLogs = msg.serverLogs

        this.callFunction()
    }

    async ngOnInit() {

        this.sharedData.isPlaylist = this.sharedData.getIsPlaylist()
        this.sharedData.isDownloadActive = this.sharedData.getIsDownloadActive()
        this.sharedData.activeDownloadMetadata = this.sharedData.getActiveDownloadMetadata()

        //get queued-items on reload
        await this.getQueuedItems(false)

        //if there is an active download
        this.getDownloadStatus()
        this.populateVideoMetadata('init')
    }

    async getDownloadStatus() {
        this.sock.subscribe();
        this.sock.next(this.wsMessage);

        this.sock.subscribe({
            next: msg => { this.updateLogs(JSON.stringify(msg)); /* console.log(JSON.stringify(msg)) */ },
            error: err => { this.updateLogs('{"download": "web-socket connection is closed."}'); /* console.log('err:', err) */ },
            complete: () => { this.wsCloseWithDownloadComplete() }
        });
        //this.sock.complete();
    }

    wsCloseWithDownloadComplete() {
        console.log('ws-close message frame recieved from server.')
        this.sharedData.isDownloadActive = false
        this.sharedData.setIsDownloadActive(false)
    }

    async updateLogs(message: string) {
        let isActiveDownload = this.sharedData.getIsDownloadActive()

        if (!isActiveDownload) {
            await this.getAndSaveActiveDownload()
            this.populateVideoMetadata()
            if (this.activeDownload !== null && this.activeDownload.title !== '') {
                this.sharedData.setIsDownloadActive(true)
            }
        }

        const log = JSON.parse(message)
        if (log.download === this.msg.downloadComplete) {
            this.serverLogs = log.download
            this.sharedData.setIsDownloadActive(false)
            await setTimeout(() => {
                this.sharedData.setRefreshAutoCompleteValue(true);
            }, 500);
        } else if (this.serverLogs === this.msg.downloadComplete) {
            this.serverLogs = this.serverLogs + ' ' + log.download
        } else if (log.download.indexOf(this.msg.downloadInfoIdentifier) !== -1) {
            this.serverLogs = log.download
        }
    }

    flipCheckbox(event: any, option: string): void {
        if (option === 'GetAudioOnly') {
            this.options.GetAudioOnly = !event.checked;
        } else if (option === 'GetSubs') {
            this.options.GetSubs = !event.checked;
        }
    }

    async GetMedia() {
        this.loading = true;
        this.activeDLTitle = this.msg.getTitle
        this.activeDLChannel = this.msg.getChannel
        this.serverLogs = this.msg.getInfo

        this.urlInputDisabled = true;
        let metadataRequest = await this.GetMetadataRequest(this.options)
        if (metadataRequest.Indicator === '') {
            this.showMessage('No URL or Identifier provided', 'error', 'error')
            return
        }

        let metadataResponse = await this.svcDownload.getMetadata(metadataRequest)
        if (metadataResponse.status === 'failure') {
            this.showMessage(metadataResponse.message, 'error', 'error')
            //Completion Process
            this.GetMediaCompleteResult(false)
        } else {
            let metadata = metadataResponse.data;
            let isDownloadActive = this.sharedData.getIsDownloadActive()
            if (!isDownloadActive) {
                await this.getAndSaveActiveDownload()
                this.populateVideoMetadata()
                this.sharedData.setIsDownloadActive(true)
            }

            //refresh autocomplete cache
            this.sharedData.setRefreshAutoCompleteValue(true);

            //trigger stats checker if all goes well
            setTimeout(() => { this.getDownloadStatus(); }, 250);

            //Completion Process
            this.GetMediaCompleteResult()
        }
    }

    GetMediaCompleteResult(showMessage: boolean = true) {

        //failure case
        if (this.activeDLTitle == this.msg.getTitle) { this.activeDLTitle = "" }
        if (this.activeDLChannel == this.msg.getChannel) { this.activeDLChannel = "" }
        if (this.serverLogs == this.msg.getInfo) { this.serverLogs = "" }

        //Complete Result
        this.resetDownloadOptions();
        this.loading = false;
        this.urlInputDisabled = false;
        if (showMessage) {
            this.showMessage('Video/Playlist queued', 'info', 'Info')
        }
    }

    GetMetadataRequest(options: ExtractionOptions) {

        let request = new VideoDataRequest()

        request.Indicator = options.Identifier
        request.SubtitlesReq = options.GetSubs
        request.IsAudioOnly = options.GetAudioOnly

        if (request.Indicator === '') {
            request.Indicator = 'UMBEkWFMacc'
            this.showMessage('Using Test Video Indicator', 'info', 'Info')
            return request
        }
        return request
    }

    populateVideoMetadata(calledFrom: string = '') {
        //get ActiveDownloadMetadata
        if (this.sharedData.getActiveDownloadMetadata() !== null) {
            this.activeDownload = this.sharedData.getActiveDownloadMetadata()[0]
        }

        if (this.activeDownload === undefined) {
            this.activeDownload = new VideoData()
        }

        this.activeDLChannel = this.activeDownload.channel
        this.activeDLTitle = this.activeDownload.title
        this.activeDLImage = this.activeDownload.thumbnail

        //show download complete status of the last download if called from init
        if (calledFrom !== 'init') { this.serverLogs = ">>>waiting for server logs<<<" }
    }

    async resetDownloadOptions() {
        this.options.GetAudioOnly = false
        this.options.GetSubs = false
        this.options.Identifier = ''
    }

    async getQueuedItems(openSidebar: boolean) {
        if (openSidebar) {
            this.sidebarVisible = true
        }
        await this.svcDownload.getQueuedItems("queued").then(item => { this.queuedItems = item; })
    }

    async getAndSaveActiveDownload() {
        await this.svcDownload.getQueuedItems("downloading").then(item => { this.sharedData.setActiveDownloadMetadata(item); })
    }

    fsStorageStatus: number = 0
    dbStorageStatus: number = 0
    async getStorageStatus() {
        try {
            let filesResult = await this.svcFiles.getStorageStatus();
            this.fsStorageStatus = filesResult.data.storage_used_fs;
            this.dbStorageStatus = filesResult.data.storage_used_db;
        } catch (e) {
            console.log('error fetching storage status:', e);
        }
    }

    getWebSocketApiUrl(): string {
        let result = apiUrl.replace("http://", "ws://")
        result = result + "/ws/status"
        return result;
    }

    ngOnDestroy() {
        // prevent memory leak when component destroyed
        if (this.sock.closed) {
            this.sock.complete();
            this.sock.unsubscribe();
        } else {
            this.sock.unsubscribe();
        }
    }

    //placeholders = ['Video or Playlist URL', 'youtube.com/watch?v=cXdwGt_xPCI', 'cXdwGt_xPCI', 'youtu.be/cXdwGt_xPCI?si=r7GvF9rndAtAkb0n', 'youtube.com/watch?v=cABfEaP2IHo&list=PL3uDtbb3OvDMn16H-YSIHpdntVMZ6V4WU', 'youtube.com/playlist?list=PL3uDtbb3OvDMn16H-YSIHpdntVMZ6V4WU', 'PL3uDtbb3OvDMn16H-YSIHpdntVMZ6V4WU'];
    placeholders = ['Video or Playlist URL', 'youtube.com/watch?v=videoId', 'videoId', 'youtu.be/videoId', 'youtube.com/watch?v=videoId&list=playlistId', 'youtube.com/playlist?list=playlistId', 'playlistId'];
    index = 0;
    callFunction() {
        setInterval(() => {
            if (this.index === this.placeholders.length - 1) { this.index = 0; } //reset
            this.index = this.index + 1;
            this.urlPlaceholder = this.placeholders[this.index]
        }, 3000);
    }

    //Toast Messages
    showMessage(message: string, severity: string, summary: string) {
        this.messageService.add({ severity: severity, summary: summary, detail: message });
    }
}


