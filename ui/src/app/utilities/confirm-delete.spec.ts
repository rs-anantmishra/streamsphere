import { fakeAsync, flushMicrotasks, tick } from '@angular/core/testing';
import { Router } from '@angular/router';
import { ConfirmationService, MessageService } from 'primeng/api';

import { confirmDelete } from './confirm-delete';

describe('confirmDelete', () => {
    let router: jasmine.SpyObj<Router>;
    let messages: jasmine.SpyObj<MessageService>;
    let confirmation: { confirm: (options: any) => void };

    beforeEach(() => {
        router = jasmine.createSpyObj('Router', ['navigate']);
        messages = jasmine.createSpyObj('MessageService', ['add']);
        // the user clicks Delete straight away
        confirmation = { confirm: options => options.accept() };
    });

    function run(deleteResult: any, onDeleted?: (id: number) => boolean | Promise<boolean>, redirectTo?: string) {
        confirmDelete({ target: null } as any, 7, () => Promise.resolve(deleteResult),
            confirmation as unknown as ConfirmationService, messages as unknown as MessageService, router as unknown as Router, { onDeleted, redirectTo });
        flushMicrotasks();
    }

    it('goes home after 3 seconds when nothing takes over', fakeAsync(() => {
        run({ data: true });

        expect(messages.add).toHaveBeenCalledWith(jasmine.objectContaining({ detail: 'Content Deleted' }));
        expect(router.navigate).not.toHaveBeenCalled();
        tick(3000);
        expect(router.navigate).toHaveBeenCalledWith(['/home']);
    }));

    it('goes to the page the caller names instead of home', fakeAsync(() => {
        run({ data: true }, undefined, '/videos');
        tick(3000);

        expect(router.navigate).toHaveBeenCalledOnceWith(['/videos']);
    }));

    it('goes to the named page when the page declines to take over', fakeAsync(() => {
        run({ data: true }, () => false, '/playlists');
        tick(3000);

        expect(router.navigate).toHaveBeenCalledOnceWith(['/playlists']);
    }));

    it('stays put when the page takes over', fakeAsync(() => {
        const onDeleted = jasmine.createSpy('onDeleted').and.returnValue(true);

        run({ data: true }, onDeleted);
        tick(3000);

        expect(onDeleted).toHaveBeenCalledWith(7);
        expect(router.navigate).not.toHaveBeenCalled();
    }));

    it('stays put when the page takes over asynchronously', fakeAsync(() => {
        run({ data: true }, () => Promise.resolve(true));
        tick(3000);

        expect(router.navigate).not.toHaveBeenCalled();
    }));

    it('goes home when the page declines to take over', fakeAsync(() => {
        run({ data: true }, () => false);
        tick(3000);

        expect(router.navigate).toHaveBeenCalledWith(['/home']);
    }));

    it('goes home when the page fails while taking over', fakeAsync(() => {
        spyOn(console, 'error');
        run({ data: true }, () => Promise.reject(new Error('boom')));
        tick(3000);

        expect(router.navigate).toHaveBeenCalledWith(['/home']);
    }));

    it('reports a failed delete and neither hands over nor leaves', fakeAsync(() => {
        const onDeleted = jasmine.createSpy('onDeleted').and.returnValue(true);

        run({ data: false }, onDeleted);
        tick(3000);

        expect(messages.add).toHaveBeenCalledWith(jasmine.objectContaining({ severity: 'error', summary: 'Failed' }));
        expect(onDeleted).not.toHaveBeenCalled();
        expect(router.navigate).not.toHaveBeenCalled();
    }));
});
