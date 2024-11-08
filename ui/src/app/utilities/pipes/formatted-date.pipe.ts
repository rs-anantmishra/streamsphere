// formatted-date.pipe.ts
import { Pipe, PipeTransform } from '@angular/core';
@Pipe({
  name: 'minifiedDate',
  standalone: true,
})
export class MinifiedDatePipe implements PipeTransform {
  transform(value: string): string {
    if (value === '') return 'not available'
        //recieved as: yyyymmdd
        let space = ' '
        let comma = ','
        const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
        let year = value.slice(0,4)
        let month = parseInt(value.slice(4,6), 10)

        let result = 'Uploaded' + space +  months[month] + comma + space + year
        return result
  }
}