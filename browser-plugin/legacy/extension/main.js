(function() {
    'use strict';
    
    const LOG_PREFIX = "🧩 EXTENSION:";

    console.log(`${LOG_PREFIX} Loading...`);

    // 1. Initialize the Library
    // Since this is an extension, the 'news-deframer-lib.js' content script
    // has already run and attached __lib_ndf to the window object.
    if (typeof window.__lib_ndf !== 'function') {
        console.error(`${LOG_PREFIX} ❌ Library not found! Check manifest order.`);
        return;
    }

    const lib = window.__lib_ndf();
    const { UniversalStorage, HelloWorld } = lib;

    console.log(`${LOG_PREFIX} ✅ Library Initialized.`);

    // 2. Hello World with Base URL
    const hello = new HelloWorld();
    const baseUrl = window.location.origin;
    console.log(`${LOG_PREFIX} ${hello.message(baseUrl)}`);

    // 3. Test Storage (using the UniversalStorage Singleton)
    // This will transparently use chrome.storage.local because 'chrome' is defined.
    const runStorageTest = async () => {
        const key = 'ext_test_key';
        
        // Read existing
        const oldVal = await UniversalStorage.get(key, 'No Data');
        console.log(`${LOG_PREFIX} 📖 Previous Storage Value: "${oldVal}"`);

        // Write new
        const newVal = `Visited ${baseUrl} at ${new Date().toLocaleTimeString()}`;
        await UniversalStorage.set(key, newVal);
        console.log(`${LOG_PREFIX} 📝 Wrote to Storage: "${newVal}"`);
    };

    runStorageTest().catch(err => console.error(`${LOG_PREFIX} ❌ Storage Error:`, err));

})();
