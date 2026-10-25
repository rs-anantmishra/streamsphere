import { parseDownloadPercent } from './parse-download-percent';

describe('parseDownloadPercent', () => {
    it('parses a progress line', () => {
        expect(parseDownloadPercent('[download]  45.3% of  229.20MiB at 19.19MiB/s ETA 00:10')).toBe(45.3);
    });

    it('parses 100%', () => {
        expect(parseDownloadPercent('[download] 100% of    3.27MiB in 00:00:00 at 20.45MiB/s')).toBe(100);
    });

    it('returns null for non-progress lines', () => {
        const lines = [
            'WARNING: [youtube] No supported JavaScript runtime',
            '',
            'Download completed successfully.',
            '[download] Saving file to library...',
        ];
        for (const line of lines) {
            expect(parseDownloadPercent(line)).toBeNull();
        }
    });
});
