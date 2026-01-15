declare global {
    // 1. Browser Extension API
    const chrome: any;

    // 2. Tampermonkey 5+ API (GM.*)
    const GM: {
        setValue: (key: string, value: any) => Promise<void>;
        getValue: <T>(key: string, defaultValue?: T) => Promise<T>;
        deleteValue: (key: string) => Promise<void>;
        listValues: () => Promise<string[]>;
        // Note: Supported by Tampermonkey 5+ / Violentmonkey
        addValueChangeListener?: (
            key: string,
            callback: (name: string, oldVal: any, newVal: any, remote: boolean) => void
        ) => Promise<number>;
    };
}

export class UniversalStorage {

    /**
     * Internal global variable for memory fallback when LocalStorage is unavailable.
     */
    private static memoryStorage: { [key: string]: any } = {};

    /**
     * Set a value.
     * Priority: Chrome Storage -> Tampermonkey 5+ (GM.setValue) -> LocalStorage
     */
    async set(key: string, value: any): Promise<void> {
        // 1. Browser Extension
        if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.local) {
            return new Promise((resolve) => chrome.storage.local.set({ [key]: value }, resolve));
        }

        // 2. Tampermonkey 5+
        if (typeof GM !== 'undefined' && GM.setValue) {
            return await GM.setValue(key, value);
        }

        // 3. Web Fallback
        try {
            localStorage.setItem(key, JSON.stringify(value));
        } catch (e) {
            console.warn("UniversalStorage: LocalStorage unavailable, using memory fallback.", e);
            UniversalStorage.memoryStorage[key] = value;
        }
        return Promise.resolve();
    }

    /**
     * Get a value.
     */
    async get<T>(key: string, defaultValue: T): Promise<T> {
        // 1. Browser Extension
        if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.local) {
            return new Promise((resolve) => {
                chrome.storage.local.get([key], (res: any) => {
                    resolve(res[key] !== undefined ? res[key] : defaultValue);
                });
            });
        }

        // 2. Tampermonkey 5+
        if (typeof GM !== 'undefined' && GM.getValue) {
            return await GM.getValue(key, defaultValue);
        }

        // 3. Web Fallback
        try {
            const item = localStorage.getItem(key);
            return Promise.resolve(item ? JSON.parse(item) : defaultValue);
        } catch (e) {
            console.warn("UniversalStorage: LocalStorage unavailable, checking memory fallback.", e);
            const memVal = UniversalStorage.memoryStorage[key];
            return Promise.resolve(memVal !== undefined ? memVal : defaultValue);
        }
    }

    /**
     * Remove a value.
     */
    async remove(key: string): Promise<void> {
        if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.local) {
            return new Promise((resolve) => chrome.storage.local.remove(key, resolve));
        }
        if (typeof GM !== 'undefined' && GM.deleteValue) {
            return await GM.deleteValue(key);
        }
        try {
            localStorage.removeItem(key);
        } catch (e) {
            console.warn("UniversalStorage: LocalStorage unavailable (remove), clearing memory fallback.", e);
        }
        delete UniversalStorage.memoryStorage[key];
        return Promise.resolve();
    }

    /**
     * Sync: Subscribe to changes from other Tabs / Scripts.
     */
    subscribe(key: string, callback: (newValue: any, oldValue: any) => void): void {
        // 1. Browser Extension
        if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.onChanged) {
            chrome.storage.onChanged.addListener((changes: any, namespace: string) => {
                if (namespace === 'local' && changes[key]) {
                    callback(changes[key].newValue, changes[key].oldValue);
                }
            });
            return;
        }

        // 2. Tampermonkey 5+
        // (Supports GM.addValueChangeListener for MV3/Async compliance)
        if (typeof GM !== 'undefined' && GM.addValueChangeListener) {
            GM.addValueChangeListener(key, (name: string, oldVal: any, newVal: any, remote: boolean) => {
                if (remote) callback(newVal, oldVal);
            });
            return;
        }

        // 3. Web Fallback
        window.addEventListener('storage', (event) => {
            if (event.key === key) {
                try {
                    const n = event.newValue ? JSON.parse(event.newValue) : null;
                    const o = event.oldValue ? JSON.parse(event.oldValue) : null;
                    callback(n, o);
                } catch { callback(event.newValue, event.oldValue); }
            }
        });
    }
}
