import { Component, Input, OnInit } from '@angular/core';
import { PlaylistsDataResponse, PlaylistsInfo, SelectedPlaylist } from '../../classes/playlists'
import { CardModule } from 'primeng/card';
import { TagModule } from 'primeng/tag';
import { TooltipModule } from 'primeng/tooltip';
import { CommonModule } from '@angular/common';
import { SharedDataService } from '../../services/shared-data.service';
import { UrlEncode } from '../../utilities/url-encode';
import { Router } from '@angular/router';


@Component({
  selector: 'app-playlist-card',
  standalone: true,
  imports: [CardModule, TagModule, TooltipModule, CommonModule],
  providers: [UrlEncode],
  templateUrl: './playlist-card.component.html',
  styleUrl: './playlist-card.component.scss'
})
export class PlaylistCardComponent implements OnInit {
    @Input() playlist: PlaylistsInfo = new PlaylistsInfo();

    constructor(private sharedDataSvc: SharedDataService, private router: Router, private urlEncode: UrlEncode) {
    }
    
    ngOnInit(): void {
        
        if (this.playlist.thumbnail == '') {
            this.playlist.thumbnail = './asstes/noimage.png'
        } else {
            this.playlist.thumbnail = this.urlEncode.encodedUrl(this.playlist.thumbnail);
        }
    }
    
    selectedPlaylist(playlist: PlaylistsInfo) {
        let selected: SelectedPlaylist = new SelectedPlaylist();
        selected.info = playlist
        this.sharedDataSvc.setPlaylist(selected)
        //send playlistId Object
        this.router.navigate(['/playlist-details']);
    }
}
