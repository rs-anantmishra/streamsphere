// likes-conversion.pipe.ts
import { Pipe, PipeTransform } from '@angular/core';
@Pipe({
  name: 'minifiedLikeCount',
  standalone: true,
})
export class MinifiedLikeCount implements PipeTransform {
  transform(likes: number, decimals = 2): string {
    if (!+likes) return '0 likes'
    if (likes === -1) return '0 likes'

    const k = 1000
    const dm = decimals < 0 ? 0 : decimals
    const sizes = [' likes', 'K likes', 'M likes', 'B likes', 'T likes']

    const i = Math.floor(Math.log(likes) / Math.log(k))
    let result = `${parseFloat((likes / Math.pow(k, i)).toFixed(dm))}${sizes[i]}`

    return result
  }
}