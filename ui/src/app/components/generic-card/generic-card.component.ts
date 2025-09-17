import { Component, Input, OnInit } from '@angular/core';
import { CardGroupDataResponse, CardGroup, SelectedCardGroup } from '../../classes/cards'
import { CardModule } from 'primeng/card';
import { TagModule } from 'primeng/tag';
import { TooltipModule } from 'primeng/tooltip';
import { CommonModule } from '@angular/common';
import { SharedDataService } from '../../services/shared-data.service';
import { UrlEncode } from '../../utilities/url-encode';
import { Router } from '@angular/router';

@Component({
    selector: 'app-generic-card',
    standalone: true,
    imports: [CardModule, TagModule, TooltipModule, CommonModule],
    providers: [UrlEncode],
    templateUrl: './generic-card.component.html',
    styleUrl: './generic-card.component.scss'
})
export class GenericCardComponent implements OnInit {

    @Input() group: CardGroup = new CardGroup();
    @Input() groupType: string = '';
    navigationUrl: string = `/${this.groupType}-details`;

    constructor(private sharedDataSvc: SharedDataService, private router: Router, private urlEncode: UrlEncode) {
    }

    ngOnInit(): void {

        if (this.group.thumbnail == '') {
            this.group.thumbnail = './asstes/noimage.png'
        } else {
            this.group.thumbnail = this.urlEncode.encodedUrl(this.group.thumbnail);
        }
    }

    selectedGroup(group: CardGroup) {
        let selected: SelectedCardGroup = new SelectedCardGroup();
        selected.groupInfo = group
        this.sharedDataSvc.setGroup(selected, this.groupType);
        this.router.navigate([this.navigationUrl]);
    }
}
