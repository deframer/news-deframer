import { HelloWorld } from './helloworld';
import { UniversalStorage } from './universalstorage.class';

// 1. Instantiate the Singleton HERE (Top-level scope).
// Because this runs when the Host script loads the library,
// this specific instance captures the Host's @grant GM permissions.
const storageInstance = new UniversalStorage();

export default function () {
  return {
    HelloWorld,

    // Return the PRE-MADE Instance (so consumers share the Host's storage)
    UniversalStorage: storageInstance,
  };
}