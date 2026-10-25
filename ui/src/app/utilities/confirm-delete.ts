import { Router } from '@angular/router';
import { ConfirmationService, MessageService } from 'primeng/api';

export interface DeleteOptions {
    //lets a page take over once the delete went through (e.g. show a neighbouring video):
    //when it returns true there is no redirect
    onDeleted?: (contentId: number) => boolean | Promise<boolean>
    //where the user is sent 3 seconds after the delete when no page took over
    redirectTo?: string
}

export function confirmDelete(event: Event, contentId: number, deleteFn: (contentId: number) => Promise<any>,
    confirmationService: ConfirmationService, messageService: MessageService, router: Router,
    options: DeleteOptions = {}): void {
    confirmationService.confirm({
        target: event.target as EventTarget,
        message: 'Are you sure you want to delete this content?',
        accept: () => {
            deleteFn(contentId).then(async (result: any) => {
                if (result.data) {
                    messageService.add({ severity: 'info', summary: 'Confirmed', detail: 'Content Deleted', life: 3000 });
                    let handled = false
                    try {
                        handled = !!(await options.onDeleted?.(contentId))
                    } catch (err) {
                        console.error('error: handling the deleted content failed.', err)
                    }
                    if (handled) { return }
                    setTimeout(() => {
                        router.navigate([options.redirectTo ?? '/home'])
                    }, 3000);
                } else {
                    messageService.add({ severity: 'error', summary: 'Failed', detail: 'Deleted Failed', life: 3000 });
                }
            });
        },
        reject: () => {
            messageService.add({ severity: 'error', summary: 'Rejected', detail: 'Delete Cancelled', life: 3000 });
        }
    });
}
