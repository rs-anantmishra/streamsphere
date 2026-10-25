// linkify.pipe.ts
import { Pipe, PipeTransform } from '@angular/core';
import { linkify } from '../linkify';
@Pipe({
    name: 'linkify',
    standalone: true,
})
export class LinkifyPipe implements PipeTransform {
    transform(text: string): string {
        return linkify(text);
    }
}