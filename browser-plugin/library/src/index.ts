import { HelloWorld } from './helloworld';
import { UniversalStorage } from './universalstorage.class';

console.log("📦 Library: Initializing singleton...");

// 1. Instantiate the Singleton HERE.
// This runs once when the library loads, capturing the Host's permissions.
const storageInstance = new UniversalStorage();

// 2. Export the Function as Default
// Webpack will take this function and assign it to window.__lib_ndf
export default function () {
    return {
        // Return the Class (if you ever need to make new Hellos)
        HelloWorld,

        // Return the PRE-MADE Instance (The Singleton)
        UniversalStorage: storageInstance,
    };
}
