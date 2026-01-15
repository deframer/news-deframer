import { UniversalStorage } from './universalstorage.class';

describe('UniversalStorage Fallback', () => {
    let storage: UniversalStorage;

    beforeEach(() => {
        storage = new UniversalStorage();
        jest.restoreAllMocks();
    });

    it('should warn and use internal global variable when localStorage fails', async () => {
        const key = 'fallback_test_key';
        const value = { test: 'data' };

        // Mock console.warn to suppress output but verify calls
        const warnSpy = jest.spyOn(console, 'warn').mockImplementation(() => {});

        // Mock localStorage to throw error.
        // We handle cases where Storage is not defined (e.g. Node environment) by defining global.localStorage directly.
        const originalLocalStorage = (global as any).localStorage;

        const mockThrow = () => { throw new Error('Access Denied'); };
        const mockLocalStorage = {
            setItem: jest.fn(mockThrow),
            getItem: jest.fn(mockThrow),
            removeItem: jest.fn(mockThrow),
        };

        Object.defineProperty(global, 'localStorage', {
            value: mockLocalStorage,
            writable: true,
            configurable: true,
        });

        try {
            // 1. Test Set Fallback
            await storage.set(key, value);
            expect(warnSpy).toHaveBeenCalledWith(
                expect.stringContaining('LocalStorage unavailable'),
                expect.any(Error)
            );

            // 2. Test Get Fallback
            const result = await storage.get(key, 'default');
            expect(result).toEqual(value);

            // 3. Test Remove Fallback
            await storage.remove(key);
            const resultAfterRemove = await storage.get(key, 'default');
            expect(resultAfterRemove).toBe('default');
        } finally {
            if (originalLocalStorage) {
                Object.defineProperty(global, 'localStorage', {
                    value: originalLocalStorage,
                    writable: true,
                    configurable: true,
                });
            } else {
                delete (global as any).localStorage;
            }
        }
    });
});