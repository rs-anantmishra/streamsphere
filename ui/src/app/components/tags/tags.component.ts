import { CommonModule } from '@angular/common';
import { Component, OnInit, effect } from '@angular/core';
import { Router, RouterModule, RouterOutlet } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { SharedDataService } from '../../services/shared-data.service';

@Component({
    selector: 'app-tags',
    standalone: true,
    imports: [CommonModule, RouterModule, ButtonModule],
    providers: [Router, SharedDataService],
    templateUrl: './tags.component.html',
    styleUrl: './tags.component.scss'
})
export class TagsComponent implements OnInit {

    constructor(private svcSharedData: SharedDataService) {
    }
    
    ngOnInit(): void {

    }
}
