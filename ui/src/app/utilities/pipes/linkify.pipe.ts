// linkify.pipe.ts
import { Pipe, PipeTransform } from '@angular/core';
@Pipe({
    name: 'linkify',
    standalone: true,
})
export class LinkifyPipe implements PipeTransform {
    transform(text: string): string {
        var urlRegex = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
        return text.replace(urlRegex, function (url: string) {
            return '<a href="' + url + '">' + url + '</a>';
        });
    }
}