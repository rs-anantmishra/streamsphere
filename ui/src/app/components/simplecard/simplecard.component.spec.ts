import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { SimplecardComponent } from './simplecard.component';
import { SharedDataService } from '../../services/shared-data.service';
import { VideoData } from '../../classes/video-data';

describe('SimplecardComponent click', () => {
    let fixture: ComponentFixture<SimplecardComponent>;
    let component: SimplecardComponent;
    let navigate: jasmine.Spy;
    let setPlayVideo: jasmine.Spy;

    beforeEach(() => {
        TestBed.configureTestingModule({ imports: [SimplecardComponent] });
        navigate = spyOn(TestBed.inject(Router), 'navigate').and.resolveTo(true);
        setPlayVideo = spyOn(TestBed.inject(SharedDataService), 'setPlayVideo');
        fixture = TestBed.createComponent(SimplecardComponent);
        component = fixture.componentInstance;
        component.metadata = new VideoData();
    });

    afterEach(() => fixture.destroy());

    it('opens the video page with its own data by default', () => {
        component.selectedVideo(component.metadata);

        expect(setPlayVideo).toHaveBeenCalledWith(component.metadata);
        expect(navigate).toHaveBeenCalledWith(['/videos/play']);
    });

    it('only reports the click when the list opens the video itself', () => {
        const clicked: VideoData[] = [];
        component.navigate = false;
        component.cardClick.subscribe(v => clicked.push(v));

        component.selectedVideo(component.metadata);

        expect(clicked).toEqual([component.metadata]);
        expect(setPlayVideo).not.toHaveBeenCalled();
        expect(navigate).not.toHaveBeenCalled();
    });
});
