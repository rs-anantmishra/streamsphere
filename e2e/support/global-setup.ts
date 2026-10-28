import { setUp, tearDown } from './stack';

export default async function globalSetup() {
    await setUp();
    return async () => { await tearDown() };
}
