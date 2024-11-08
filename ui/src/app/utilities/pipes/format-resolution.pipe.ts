// format-resolution.pipe.ts
import { Pipe, PipeTransform } from '@angular/core';
@Pipe({
  name: 'formattedResolution',
  standalone: true,
})
export class FormattedResolutionPipe implements PipeTransform {
  transform(value: string): string {
    if (value === '') return 'not available'
        //616 - 1920x1080 (Premium)+251 - audio only (medium)
        let videoSplit = value.split('+')[0]
        let videoFormatDetails = videoSplit.split('-')[1].trim()
        let resolution = videoFormatDetails.split(' ')[0].trim()

        return resolution
  }
}