import { Router } from '@angular/router';
import { ConfirmationService, MessageService } from 'primeng/api';

export function confirmDelete(event: Event, contentId: number, deleteFn: (contentId: number) => Promise<any>,
    confirmationService: ConfirmationService, messageService: MessageService, router: Router): void {
    confirmationService.confirm({
        target: event.target as EventTarget,
        message: 'Are you sure you want to delete this content?',
        accept: () => {
            deleteFn(contentId).then((result: any) => {
                if (result.data) {
                    messageService.add({ severity: 'info', summary: 'Confirmed', detail: 'Content Deleted', life: 3000 });
                    setTimeout(() => {
                        router.navigate(['/home'])
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
