// array-comma-sep.pipe.ts
import { Pipe, PipeTransform } from '@angular/core';
@Pipe({
    name: 'commaSepStringFromArray',
    standalone: true,
})
export class CommaSepStringFromArray implements PipeTransform {
    transform(value: string[], isTags = false): string {
        if (value !== null && value.length > 0) {
            if (isTags) {
                value = value.map((line) => `#${line}`);
            }
            let result = value.join(', ')
            return result
        } else {
            return 'not available'
        }
    }
}