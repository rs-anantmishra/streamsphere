import { Component, Input, OnChanges, OnInit, SimpleChanges } from '@angular/core';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { VideoData } from '../../classes/video-data'
import { CardModule } from 'primeng/card';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { SharedDataService } from '../../services/shared-data.service';
import { TooltipModule } from 'primeng/tooltip';
import { TagModule } from 'primeng/tag';
import { FilesizeConversionPipe } from '../../utilities/pipes/filesize-conversion.pipe'
import { CommaSepStringFromArray } from '../../utilities/pipes/array-comma-sep.pipe'
import { MinifiedViewCount } from '../../utilities/pipes/views-conversion.pipe';
import { MinifiedDatePipe } from '../../utilities/pipes/formatted-date.pipe';
import { FormattedResolutionPipe } from '../../utilities/pipes/format-resolution.pipe';

@Component({
    selector: 'app-simplecard',
    standalone: true,
    imports: [ToastModule, CardModule, CommonModule, TooltipModule, TagModule, FilesizeConversionPipe, CommaSepStringFromArray, MinifiedViewCount, MinifiedDatePipe, FormattedResolutionPipe],
    providers: [MessageService, Router],
    templateUrl: './simplecard.component.html',
    styleUrl: './simplecard.component.scss'
})
export class SimplecardComponent implements OnInit {

    @Input() metadata!: VideoData;
    constructor(private router: Router, private svcSharedData: SharedDataService) {
    }

    ngOnInit(): void {
        if (this.metadata.thumbnail == '') {
            this.metadata.thumbnail = './assets/noimage.png'
        } else {
            this.metadata.media_url = this.metadata.media_url.replaceAll('#', '%23')
            this.metadata.thumbnail = this.metadata.thumbnail.replaceAll('#', '%23')
            this.metadata.webpage_url = this.metadata.webpage_url.replaceAll('#', '%23')
        }
    }

    selectedVideo(playVideo: VideoData) {
        this.svcSharedData.setPlayVideo(playVideo);
        this.router.navigate(['/videos/play'])
    }

    getFormattedDuration(duration: number) {
        // Calculate minutes and seconds
        const minutes = Math.floor(duration / 60);
        const seconds = duration % 60;

        // Format the result as MM:SS
        return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
    }

}
