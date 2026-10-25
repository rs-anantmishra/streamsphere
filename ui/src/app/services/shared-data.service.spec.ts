import { TestBed } from '@angular/core/testing';

import { SharedDataService } from './shared-data.service';

describe('SharedDataService', () => {
    afterEach(() => localStorage.clear());

    it('clears the download keys left by earlier versions and keeps the rest', () => {
        localStorage.setItem('isDownloadActive', 'true');
        localStorage.setItem('activeDownloadMetadata', '[]');
        localStorage.setItem('isDarkMode', 'true');

        TestBed.inject(SharedDataService);

        expect(localStorage.getItem('isDownloadActive')).toBeNull();
        expect(localStorage.getItem('activeDownloadMetadata')).toBeNull();
        expect(localStorage.getItem('isDarkMode')).toBe('true');
    });
});
