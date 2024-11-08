// remove-prefix.pipe.ts
import { Pipe, PipeTransform } from '@angular/core';
@Pipe({
  name: 'removePrefix',
  standalone: true,
})
export class RemovePrefixPipe implements PipeTransform {
  transform(value: string): string {
    let transformed = value.replace("http://","").replace("https://","")
    return transformed;
  }
}