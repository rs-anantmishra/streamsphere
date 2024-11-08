// filesize-conversion.pipe.ts
import { Pipe, PipeTransform } from '@angular/core';
@Pipe({
  name: 'filesizeToMiB',
  standalone: true,
})
export class FilesizeConversionPipe implements PipeTransform {
  transform(bytes: number, decimals = 2): string {
    //let bytes = parseInt(value, 10); //base 10 for decimal value
    if (!+bytes) return '0 Bytes'
    if (bytes === -1) return '0 Bytes'

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

    const i = Math.floor(Math.log(bytes) / Math.log(k))
    let result = `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`

    return result
  }
}