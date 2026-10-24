import { Component, OnInit, ViewEncapsulation, ViewChild } from '@angular/core';
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
import { UrlEncode } from '../../utilities/url-encode';
import { linkify } from '../../utilities/linkify';
import { cp1252_to_utf8 } from '../../utilities/cp1252-to-utf8';
import { downloadBlob } from '../../utilities/download-blob';
import { confirmDelete } from '../../utilities/confirm-delete';
import Plyr from 'plyr';

import { MinifiedViewCount } from '../../utilities/pipes/views-conversion.pipe'
import { MinifiedLikeCount } from '../../utilities/pipes/likes-conversion.pipe';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { CommaSepStringFromArray } from "../../utilities/pipes/array-comma-sep.pipe";
import { FormattedResolutionPipe } from "../../utilities/pipes/format-resolution.pipe";
import { VideosService } from '../../services/videos.service';
import { MinifiedDatePipe } from "../../utilities/pipes/formatted-date.pipe";
import { FilesizeConversionPipe } from "../../utilities/pipes/filesize-conversion.pipe";

@Component({
    selector: 'app-video-details',
    imports: [RouterModule, ButtonModule, PanelModule, ScrollPanelModule, TagModule, ChipModule, ConfirmPopupModule, ToastModule,
        MinifiedViewCount, MinifiedLikeCount, CommaSepStringFromArray, FormattedResolutionPipe, MinifiedDatePipe, FieldsetModule, FilesizeConversionPipe],
    providers: [Router, SharedDataService, ConfirmationService, MessageService, UrlEncode],
    templateUrl: './video-details.component.html',
    styleUrl: './video-details.component.scss',
    encapsulation: ViewEncapsulation.None
})
export class VideoDetailsComponent implements OnInit {

    public player: any;
    selectedVideo: VideoData = new VideoData()
    @ViewChild(ConfirmPopup) confirmPopup!: ConfirmPopup;
    constructor(private confirmationService: ConfirmationService, private messageService: MessageService, private router: Router, private urlEncode: UrlEncode,
        private svcSharedData: SharedDataService, private svcVideos: VideosService) {
    }

    async ngOnInit(): Promise<void> {
        this.selectedVideo = await this.svcSharedData.getActivePlayerMetadata();

        this.selectedVideo.media_url = this.selectedVideo.media_url
        this.selectedVideo.thumbnail = this.selectedVideo.thumbnail
        this.selectedVideo.webpage_url = this.selectedVideo.webpage_url

        this.selectedVideo.description = cp1252_to_utf8(this.selectedVideo.description)
        this.selectedVideo.description = linkify(this.selectedVideo.description)

        this.player = new Plyr('#plyrId', { captions: { active: true }, keyboard: { global: true }, autoplay: true });
    }

    async download(): Promise<void> {
        (await this.svcVideos
            .download(this.selectedVideo.media_url))
            .subscribe(blob => downloadBlob(blob, this.selectedVideo.title + '.' + this.selectedVideo.extension))
    }

    accept() {
        this.confirmPopup.onAccept();
    }

    reject() {
        this.confirmPopup.onReject();
    }

    confirm(event: Event, contentId: number) {
        confirmDelete(event, contentId, id => this.svcVideos.deleteVideoById(id), this.confirmationService, this.messageService, this.router);
    }
}
