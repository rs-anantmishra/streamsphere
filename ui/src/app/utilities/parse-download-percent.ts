export function parseDownloadPercent(line: string): number | null {
    const match = /^\[download\]\s+([\d.]+)%/.exec(line);
    return match ? parseFloat(match[1]) : null;
}
