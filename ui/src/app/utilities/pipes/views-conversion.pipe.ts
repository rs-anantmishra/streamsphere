// views-conversion.pipe.ts
import { Pipe, PipeTransform } from '@angular/core';
@Pipe({
  name: 'minifiedViewCount',
  standalone: true,
})
export class MinifiedViewCount implements PipeTransform {
  transform(views: number, decimals = 2): string {
    if (!+views) return '0 views'
    if (views === -1) return '0 views'

    const k = 1000
    const dm = decimals < 0 ? 0 : decimals
    const sizes = [' views', 'K views', 'M views', 'B views', 'T views']

    const i = Math.floor(Math.log(views) / Math.log(k))
    let result = `${parseFloat((views / Math.pow(k, i)).toFixed(dm))}${sizes[i]}`

    return result
  }
}