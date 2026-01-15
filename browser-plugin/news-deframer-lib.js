/******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./src/helloworld.ts"
/*!***************************!*\
  !*** ./src/helloworld.ts ***!
  \***************************/
(__unused_webpack_module, exports) {


Object.defineProperty(exports, "__esModule", ({ value: true }));
exports.HelloWorld = void 0;
class HelloWorld {
    message(name = "World") {
        return `Hello ${name}`;
    }
}
exports.HelloWorld = HelloWorld;


/***/ },

/***/ "./src/universalstorage.class.ts"
/*!***************************************!*\
  !*** ./src/universalstorage.class.ts ***!
  \***************************************/
(__unused_webpack_module, exports) {


var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", ({ value: true }));
exports.UniversalStorage = void 0;
class UniversalStorage {
    /**
     * Set a value.
     * Priority: Chrome Storage -> Tampermonkey 5+ (GM.setValue) -> LocalStorage
     */
    set(key, value) {
        return __awaiter(this, void 0, void 0, function* () {
            // 1. Browser Extension
            if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.local) {
                return new Promise((resolve) => chrome.storage.local.set({ [key]: value }, resolve));
            }
            // 2. Tampermonkey 5+
            if (typeof GM !== 'undefined' && GM.setValue) {
                return yield GM.setValue(key, value);
            }
            // 3. Web Fallback
            try {
                localStorage.setItem(key, JSON.stringify(value));
            }
            catch (e) {
                console.warn("UniversalStorage: LocalStorage unavailable, using memory fallback.", e);
                UniversalStorage.memoryStorage[key] = value;
            }
            return Promise.resolve();
        });
    }
    /**
     * Get a value.
     */
    get(key, defaultValue) {
        return __awaiter(this, void 0, void 0, function* () {
            // 1. Browser Extension
            if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.local) {
                return new Promise((resolve) => {
                    chrome.storage.local.get([key], (res) => {
                        resolve(res[key] !== undefined ? res[key] : defaultValue);
                    });
                });
            }
            // 2. Tampermonkey 5+
            if (typeof GM !== 'undefined' && GM.getValue) {
                return yield GM.getValue(key, defaultValue);
            }
            // 3. Web Fallback
            try {
                const item = localStorage.getItem(key);
                return Promise.resolve(item ? JSON.parse(item) : defaultValue);
            }
            catch (e) {
                console.warn("UniversalStorage: LocalStorage unavailable, checking memory fallback.", e);
                const memVal = UniversalStorage.memoryStorage[key];
                return Promise.resolve(memVal !== undefined ? memVal : defaultValue);
            }
        });
    }
    /**
     * Remove a value.
     */
    remove(key) {
        return __awaiter(this, void 0, void 0, function* () {
            if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.local) {
                return new Promise((resolve) => chrome.storage.local.remove(key, resolve));
            }
            if (typeof GM !== 'undefined' && GM.deleteValue) {
                return yield GM.deleteValue(key);
            }
            try {
                localStorage.removeItem(key);
            }
            catch (e) {
                console.warn("UniversalStorage: LocalStorage unavailable (remove), clearing memory fallback.", e);
            }
            delete UniversalStorage.memoryStorage[key];
            return Promise.resolve();
        });
    }
    /**
     * Sync: Subscribe to changes from other Tabs / Scripts.
     */
    subscribe(key, callback) {
        // 1. Browser Extension
        if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.onChanged) {
            chrome.storage.onChanged.addListener((changes, namespace) => {
                if (namespace === 'local' && changes[key]) {
                    callback(changes[key].newValue, changes[key].oldValue);
                }
            });
            return;
        }
        // 2. Tampermonkey 5+
        // (Supports GM.addValueChangeListener for MV3/Async compliance)
        if (typeof GM !== 'undefined' && GM.addValueChangeListener) {
            GM.addValueChangeListener(key, (name, oldVal, newVal, remote) => {
                if (remote)
                    callback(newVal, oldVal);
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
                }
                catch (_a) {
                    callback(event.newValue, event.oldValue);
                }
            }
        });
    }
}
exports.UniversalStorage = UniversalStorage;
/**
 * Internal global variable for memory fallback when LocalStorage is unavailable.
 */
UniversalStorage.memoryStorage = {};


/***/ }

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Check if module exists (development only)
/******/ 		if (__webpack_modules__[moduleId] === undefined) {
/******/ 			var e = new Error("Cannot find module '" + moduleId + "'");
/******/ 			e.code = 'MODULE_NOT_FOUND';
/******/ 			throw e;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/************************************************************************/
var __webpack_exports__ = {};
// This entry needs to be wrapped in an IIFE because it needs to be isolated against other modules in the chunk.
(() => {
var exports = __webpack_exports__;
/*!**********************!*\
  !*** ./src/index.ts ***!
  \**********************/

Object.defineProperty(exports, "__esModule", ({ value: true }));
exports["default"] = default_1;
const helloworld_1 = __webpack_require__(/*! ./helloworld */ "./src/helloworld.ts");
const universalstorage_class_1 = __webpack_require__(/*! ./universalstorage.class */ "./src/universalstorage.class.ts");
// 1. Instantiate the Singleton HERE (Top-level scope).
// Because this runs when the Host script loads the library,
// this specific instance captures the Host's @grant GM permissions.
const storageInstance = new universalstorage_class_1.UniversalStorage();
function default_1() {
    return {
        HelloWorld: helloworld_1.HelloWorld,
        // Return the PRE-MADE Instance (so consumers share the Host's storage)
        UniversalStorage: storageInstance,
    };
}
// Ensure global attachment for Webpack/Rollup IIFE builds
window.__lib_ndf = function () {
    return { HelloWorld: helloworld_1.HelloWorld, UniversalStorage: storageInstance };
};

})();

window.__lib_ndf = __webpack_exports__["default"];
/******/ })()
;
//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibmV3cy1kZWZyYW1lci1saWIuanMiLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7OztBQUFhO0FBQ2IsOENBQTZDLEVBQUUsYUFBYSxFQUFDO0FBQzdELGtCQUFrQjtBQUNsQjtBQUNBO0FBQ0Esd0JBQXdCLEtBQUs7QUFDN0I7QUFDQTtBQUNBLGtCQUFrQjs7Ozs7Ozs7Ozs7QUNSTDtBQUNiO0FBQ0EsNEJBQTRCLCtEQUErRCxpQkFBaUI7QUFDNUc7QUFDQSxvQ0FBb0MsTUFBTSwrQkFBK0IsWUFBWTtBQUNyRixtQ0FBbUMsTUFBTSxtQ0FBbUMsWUFBWTtBQUN4RixnQ0FBZ0M7QUFDaEM7QUFDQSxLQUFLO0FBQ0w7QUFDQSw4Q0FBNkMsRUFBRSxhQUFhLEVBQUM7QUFDN0Qsd0JBQXdCO0FBQ3hCO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLDJFQUEyRSxjQUFjO0FBQ3pGO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxTQUFTO0FBQ1Q7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLHFCQUFxQjtBQUNyQixpQkFBaUI7QUFDakI7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsU0FBUztBQUNUO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxTQUFTO0FBQ1Q7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLGFBQWE7QUFDYjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsYUFBYTtBQUNiO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxTQUFTO0FBQ1Q7QUFDQTtBQUNBLHdCQUF3QjtBQUN4QjtBQUNBO0FBQ0E7QUFDQTs7Ozs7OztVQ2pJQTtVQUNBOztVQUVBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBOztVQUVBO1VBQ0E7O1VBRUE7VUFDQTtVQUNBOzs7Ozs7Ozs7O0FDNUJhO0FBQ2IsOENBQTZDLEVBQUUsYUFBYSxFQUFDO0FBQzdELGtCQUFlO0FBQ2YscUJBQXFCLG1CQUFPLENBQUMseUNBQWM7QUFDM0MsaUNBQWlDLG1CQUFPLENBQUMsaUVBQTBCO0FBQ25FO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsYUFBYTtBQUNiIiwic291cmNlcyI6WyJ3ZWJwYWNrOi8vX19saWJfbmRmLy4vc3JjL2hlbGxvd29ybGQudHMiLCJ3ZWJwYWNrOi8vX19saWJfbmRmLy4vc3JjL3VuaXZlcnNhbHN0b3JhZ2UuY2xhc3MudHMiLCJ3ZWJwYWNrOi8vX19saWJfbmRmL3dlYnBhY2svYm9vdHN0cmFwIiwid2VicGFjazovL19fbGliX25kZi8uL3NyYy9pbmRleC50cyJdLCJzb3VyY2VzQ29udGVudCI6WyJcInVzZSBzdHJpY3RcIjtcbk9iamVjdC5kZWZpbmVQcm9wZXJ0eShleHBvcnRzLCBcIl9fZXNNb2R1bGVcIiwgeyB2YWx1ZTogdHJ1ZSB9KTtcbmV4cG9ydHMuSGVsbG9Xb3JsZCA9IHZvaWQgMDtcbmNsYXNzIEhlbGxvV29ybGQge1xuICAgIG1lc3NhZ2UobmFtZSA9IFwiV29ybGRcIikge1xuICAgICAgICByZXR1cm4gYEhlbGxvICR7bmFtZX1gO1xuICAgIH1cbn1cbmV4cG9ydHMuSGVsbG9Xb3JsZCA9IEhlbGxvV29ybGQ7XG4iLCJcInVzZSBzdHJpY3RcIjtcbnZhciBfX2F3YWl0ZXIgPSAodGhpcyAmJiB0aGlzLl9fYXdhaXRlcikgfHwgZnVuY3Rpb24gKHRoaXNBcmcsIF9hcmd1bWVudHMsIFAsIGdlbmVyYXRvcikge1xuICAgIGZ1bmN0aW9uIGFkb3B0KHZhbHVlKSB7IHJldHVybiB2YWx1ZSBpbnN0YW5jZW9mIFAgPyB2YWx1ZSA6IG5ldyBQKGZ1bmN0aW9uIChyZXNvbHZlKSB7IHJlc29sdmUodmFsdWUpOyB9KTsgfVxuICAgIHJldHVybiBuZXcgKFAgfHwgKFAgPSBQcm9taXNlKSkoZnVuY3Rpb24gKHJlc29sdmUsIHJlamVjdCkge1xuICAgICAgICBmdW5jdGlvbiBmdWxmaWxsZWQodmFsdWUpIHsgdHJ5IHsgc3RlcChnZW5lcmF0b3IubmV4dCh2YWx1ZSkpOyB9IGNhdGNoIChlKSB7IHJlamVjdChlKTsgfSB9XG4gICAgICAgIGZ1bmN0aW9uIHJlamVjdGVkKHZhbHVlKSB7IHRyeSB7IHN0ZXAoZ2VuZXJhdG9yW1widGhyb3dcIl0odmFsdWUpKTsgfSBjYXRjaCAoZSkgeyByZWplY3QoZSk7IH0gfVxuICAgICAgICBmdW5jdGlvbiBzdGVwKHJlc3VsdCkgeyByZXN1bHQuZG9uZSA/IHJlc29sdmUocmVzdWx0LnZhbHVlKSA6IGFkb3B0KHJlc3VsdC52YWx1ZSkudGhlbihmdWxmaWxsZWQsIHJlamVjdGVkKTsgfVxuICAgICAgICBzdGVwKChnZW5lcmF0b3IgPSBnZW5lcmF0b3IuYXBwbHkodGhpc0FyZywgX2FyZ3VtZW50cyB8fCBbXSkpLm5leHQoKSk7XG4gICAgfSk7XG59O1xuT2JqZWN0LmRlZmluZVByb3BlcnR5KGV4cG9ydHMsIFwiX19lc01vZHVsZVwiLCB7IHZhbHVlOiB0cnVlIH0pO1xuZXhwb3J0cy5Vbml2ZXJzYWxTdG9yYWdlID0gdm9pZCAwO1xuY2xhc3MgVW5pdmVyc2FsU3RvcmFnZSB7XG4gICAgLyoqXG4gICAgICogU2V0IGEgdmFsdWUuXG4gICAgICogUHJpb3JpdHk6IENocm9tZSBTdG9yYWdlIC0+IFRhbXBlcm1vbmtleSA1KyAoR00uc2V0VmFsdWUpIC0+IExvY2FsU3RvcmFnZVxuICAgICAqL1xuICAgIHNldChrZXksIHZhbHVlKSB7XG4gICAgICAgIHJldHVybiBfX2F3YWl0ZXIodGhpcywgdm9pZCAwLCB2b2lkIDAsIGZ1bmN0aW9uKiAoKSB7XG4gICAgICAgICAgICAvLyAxLiBCcm93c2VyIEV4dGVuc2lvblxuICAgICAgICAgICAgaWYgKHR5cGVvZiBjaHJvbWUgIT09ICd1bmRlZmluZWQnICYmIGNocm9tZS5zdG9yYWdlICYmIGNocm9tZS5zdG9yYWdlLmxvY2FsKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlKSA9PiBjaHJvbWUuc3RvcmFnZS5sb2NhbC5zZXQoeyBba2V5XTogdmFsdWUgfSwgcmVzb2x2ZSkpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gMi4gVGFtcGVybW9ua2V5IDUrXG4gICAgICAgICAgICBpZiAodHlwZW9mIEdNICE9PSAndW5kZWZpbmVkJyAmJiBHTS5zZXRWYWx1ZSkge1xuICAgICAgICAgICAgICAgIHJldHVybiB5aWVsZCBHTS5zZXRWYWx1ZShrZXksIHZhbHVlKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIDMuIFdlYiBGYWxsYmFja1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICBsb2NhbFN0b3JhZ2Uuc2V0SXRlbShrZXksIEpTT04uc3RyaW5naWZ5KHZhbHVlKSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBjYXRjaCAoZSkge1xuICAgICAgICAgICAgICAgIGNvbnNvbGUud2FybihcIlVuaXZlcnNhbFN0b3JhZ2U6IExvY2FsU3RvcmFnZSB1bmF2YWlsYWJsZSwgdXNpbmcgbWVtb3J5IGZhbGxiYWNrLlwiLCBlKTtcbiAgICAgICAgICAgICAgICBVbml2ZXJzYWxTdG9yYWdlLm1lbW9yeVN0b3JhZ2Vba2V5XSA9IHZhbHVlO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgcmV0dXJuIFByb21pc2UucmVzb2x2ZSgpO1xuICAgICAgICB9KTtcbiAgICB9XG4gICAgLyoqXG4gICAgICogR2V0IGEgdmFsdWUuXG4gICAgICovXG4gICAgZ2V0KGtleSwgZGVmYXVsdFZhbHVlKSB7XG4gICAgICAgIHJldHVybiBfX2F3YWl0ZXIodGhpcywgdm9pZCAwLCB2b2lkIDAsIGZ1bmN0aW9uKiAoKSB7XG4gICAgICAgICAgICAvLyAxLiBCcm93c2VyIEV4dGVuc2lvblxuICAgICAgICAgICAgaWYgKHR5cGVvZiBjaHJvbWUgIT09ICd1bmRlZmluZWQnICYmIGNocm9tZS5zdG9yYWdlICYmIGNocm9tZS5zdG9yYWdlLmxvY2FsKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlKSA9PiB7XG4gICAgICAgICAgICAgICAgICAgIGNocm9tZS5zdG9yYWdlLmxvY2FsLmdldChba2V5XSwgKHJlcykgPT4ge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmVzb2x2ZShyZXNba2V5XSAhPT0gdW5kZWZpbmVkID8gcmVzW2tleV0gOiBkZWZhdWx0VmFsdWUpO1xuICAgICAgICAgICAgICAgICAgICB9KTtcbiAgICAgICAgICAgICAgICB9KTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIDIuIFRhbXBlcm1vbmtleSA1K1xuICAgICAgICAgICAgaWYgKHR5cGVvZiBHTSAhPT0gJ3VuZGVmaW5lZCcgJiYgR00uZ2V0VmFsdWUpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4geWllbGQgR00uZ2V0VmFsdWUoa2V5LCBkZWZhdWx0VmFsdWUpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gMy4gV2ViIEZhbGxiYWNrXG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIGNvbnN0IGl0ZW0gPSBsb2NhbFN0b3JhZ2UuZ2V0SXRlbShrZXkpO1xuICAgICAgICAgICAgICAgIHJldHVybiBQcm9taXNlLnJlc29sdmUoaXRlbSA/IEpTT04ucGFyc2UoaXRlbSkgOiBkZWZhdWx0VmFsdWUpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2F0Y2ggKGUpIHtcbiAgICAgICAgICAgICAgICBjb25zb2xlLndhcm4oXCJVbml2ZXJzYWxTdG9yYWdlOiBMb2NhbFN0b3JhZ2UgdW5hdmFpbGFibGUsIGNoZWNraW5nIG1lbW9yeSBmYWxsYmFjay5cIiwgZSk7XG4gICAgICAgICAgICAgICAgY29uc3QgbWVtVmFsID0gVW5pdmVyc2FsU3RvcmFnZS5tZW1vcnlTdG9yYWdlW2tleV07XG4gICAgICAgICAgICAgICAgcmV0dXJuIFByb21pc2UucmVzb2x2ZShtZW1WYWwgIT09IHVuZGVmaW5lZCA/IG1lbVZhbCA6IGRlZmF1bHRWYWx1ZSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuICAgIH1cbiAgICAvKipcbiAgICAgKiBSZW1vdmUgYSB2YWx1ZS5cbiAgICAgKi9cbiAgICByZW1vdmUoa2V5KSB7XG4gICAgICAgIHJldHVybiBfX2F3YWl0ZXIodGhpcywgdm9pZCAwLCB2b2lkIDAsIGZ1bmN0aW9uKiAoKSB7XG4gICAgICAgICAgICBpZiAodHlwZW9mIGNocm9tZSAhPT0gJ3VuZGVmaW5lZCcgJiYgY2hyb21lLnN0b3JhZ2UgJiYgY2hyb21lLnN0b3JhZ2UubG9jYWwpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUpID0+IGNocm9tZS5zdG9yYWdlLmxvY2FsLnJlbW92ZShrZXksIHJlc29sdmUpKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGlmICh0eXBlb2YgR00gIT09ICd1bmRlZmluZWQnICYmIEdNLmRlbGV0ZVZhbHVlKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIHlpZWxkIEdNLmRlbGV0ZVZhbHVlKGtleSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIGxvY2FsU3RvcmFnZS5yZW1vdmVJdGVtKGtleSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBjYXRjaCAoZSkge1xuICAgICAgICAgICAgICAgIGNvbnNvbGUud2FybihcIlVuaXZlcnNhbFN0b3JhZ2U6IExvY2FsU3RvcmFnZSB1bmF2YWlsYWJsZSAocmVtb3ZlKSwgY2xlYXJpbmcgbWVtb3J5IGZhbGxiYWNrLlwiLCBlKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGRlbGV0ZSBVbml2ZXJzYWxTdG9yYWdlLm1lbW9yeVN0b3JhZ2Vba2V5XTtcbiAgICAgICAgICAgIHJldHVybiBQcm9taXNlLnJlc29sdmUoKTtcbiAgICAgICAgfSk7XG4gICAgfVxuICAgIC8qKlxuICAgICAqIFN5bmM6IFN1YnNjcmliZSB0byBjaGFuZ2VzIGZyb20gb3RoZXIgVGFicyAvIFNjcmlwdHMuXG4gICAgICovXG4gICAgc3Vic2NyaWJlKGtleSwgY2FsbGJhY2spIHtcbiAgICAgICAgLy8gMS4gQnJvd3NlciBFeHRlbnNpb25cbiAgICAgICAgaWYgKHR5cGVvZiBjaHJvbWUgIT09ICd1bmRlZmluZWQnICYmIGNocm9tZS5zdG9yYWdlICYmIGNocm9tZS5zdG9yYWdlLm9uQ2hhbmdlZCkge1xuICAgICAgICAgICAgY2hyb21lLnN0b3JhZ2Uub25DaGFuZ2VkLmFkZExpc3RlbmVyKChjaGFuZ2VzLCBuYW1lc3BhY2UpID0+IHtcbiAgICAgICAgICAgICAgICBpZiAobmFtZXNwYWNlID09PSAnbG9jYWwnICYmIGNoYW5nZXNba2V5XSkge1xuICAgICAgICAgICAgICAgICAgICBjYWxsYmFjayhjaGFuZ2VzW2tleV0ubmV3VmFsdWUsIGNoYW5nZXNba2V5XS5vbGRWYWx1ZSk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgLy8gMi4gVGFtcGVybW9ua2V5IDUrXG4gICAgICAgIC8vIChTdXBwb3J0cyBHTS5hZGRWYWx1ZUNoYW5nZUxpc3RlbmVyIGZvciBNVjMvQXN5bmMgY29tcGxpYW5jZSlcbiAgICAgICAgaWYgKHR5cGVvZiBHTSAhPT0gJ3VuZGVmaW5lZCcgJiYgR00uYWRkVmFsdWVDaGFuZ2VMaXN0ZW5lcikge1xuICAgICAgICAgICAgR00uYWRkVmFsdWVDaGFuZ2VMaXN0ZW5lcihrZXksIChuYW1lLCBvbGRWYWwsIG5ld1ZhbCwgcmVtb3RlKSA9PiB7XG4gICAgICAgICAgICAgICAgaWYgKHJlbW90ZSlcbiAgICAgICAgICAgICAgICAgICAgY2FsbGJhY2sobmV3VmFsLCBvbGRWYWwpO1xuICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgLy8gMy4gV2ViIEZhbGxiYWNrXG4gICAgICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdzdG9yYWdlJywgKGV2ZW50KSA9PiB7XG4gICAgICAgICAgICBpZiAoZXZlbnQua2V5ID09PSBrZXkpIHtcbiAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICBjb25zdCBuID0gZXZlbnQubmV3VmFsdWUgPyBKU09OLnBhcnNlKGV2ZW50Lm5ld1ZhbHVlKSA6IG51bGw7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IG8gPSBldmVudC5vbGRWYWx1ZSA/IEpTT04ucGFyc2UoZXZlbnQub2xkVmFsdWUpIDogbnVsbDtcbiAgICAgICAgICAgICAgICAgICAgY2FsbGJhY2sobiwgbyk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIGNhdGNoIChfYSkge1xuICAgICAgICAgICAgICAgICAgICBjYWxsYmFjayhldmVudC5uZXdWYWx1ZSwgZXZlbnQub2xkVmFsdWUpO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgfVxufVxuZXhwb3J0cy5Vbml2ZXJzYWxTdG9yYWdlID0gVW5pdmVyc2FsU3RvcmFnZTtcbi8qKlxuICogSW50ZXJuYWwgZ2xvYmFsIHZhcmlhYmxlIGZvciBtZW1vcnkgZmFsbGJhY2sgd2hlbiBMb2NhbFN0b3JhZ2UgaXMgdW5hdmFpbGFibGUuXG4gKi9cblVuaXZlcnNhbFN0b3JhZ2UubWVtb3J5U3RvcmFnZSA9IHt9O1xuIiwiLy8gVGhlIG1vZHVsZSBjYWNoZVxudmFyIF9fd2VicGFja19tb2R1bGVfY2FjaGVfXyA9IHt9O1xuXG4vLyBUaGUgcmVxdWlyZSBmdW5jdGlvblxuZnVuY3Rpb24gX193ZWJwYWNrX3JlcXVpcmVfXyhtb2R1bGVJZCkge1xuXHQvLyBDaGVjayBpZiBtb2R1bGUgaXMgaW4gY2FjaGVcblx0dmFyIGNhY2hlZE1vZHVsZSA9IF9fd2VicGFja19tb2R1bGVfY2FjaGVfX1ttb2R1bGVJZF07XG5cdGlmIChjYWNoZWRNb2R1bGUgIT09IHVuZGVmaW5lZCkge1xuXHRcdHJldHVybiBjYWNoZWRNb2R1bGUuZXhwb3J0cztcblx0fVxuXHQvLyBDaGVjayBpZiBtb2R1bGUgZXhpc3RzIChkZXZlbG9wbWVudCBvbmx5KVxuXHRpZiAoX193ZWJwYWNrX21vZHVsZXNfX1ttb2R1bGVJZF0gPT09IHVuZGVmaW5lZCkge1xuXHRcdHZhciBlID0gbmV3IEVycm9yKFwiQ2Fubm90IGZpbmQgbW9kdWxlICdcIiArIG1vZHVsZUlkICsgXCInXCIpO1xuXHRcdGUuY29kZSA9ICdNT0RVTEVfTk9UX0ZPVU5EJztcblx0XHR0aHJvdyBlO1xuXHR9XG5cdC8vIENyZWF0ZSBhIG5ldyBtb2R1bGUgKGFuZCBwdXQgaXQgaW50byB0aGUgY2FjaGUpXG5cdHZhciBtb2R1bGUgPSBfX3dlYnBhY2tfbW9kdWxlX2NhY2hlX19bbW9kdWxlSWRdID0ge1xuXHRcdC8vIG5vIG1vZHVsZS5pZCBuZWVkZWRcblx0XHQvLyBubyBtb2R1bGUubG9hZGVkIG5lZWRlZFxuXHRcdGV4cG9ydHM6IHt9XG5cdH07XG5cblx0Ly8gRXhlY3V0ZSB0aGUgbW9kdWxlIGZ1bmN0aW9uXG5cdF9fd2VicGFja19tb2R1bGVzX19bbW9kdWxlSWRdLmNhbGwobW9kdWxlLmV4cG9ydHMsIG1vZHVsZSwgbW9kdWxlLmV4cG9ydHMsIF9fd2VicGFja19yZXF1aXJlX18pO1xuXG5cdC8vIFJldHVybiB0aGUgZXhwb3J0cyBvZiB0aGUgbW9kdWxlXG5cdHJldHVybiBtb2R1bGUuZXhwb3J0cztcbn1cblxuIiwiXCJ1c2Ugc3RyaWN0XCI7XG5PYmplY3QuZGVmaW5lUHJvcGVydHkoZXhwb3J0cywgXCJfX2VzTW9kdWxlXCIsIHsgdmFsdWU6IHRydWUgfSk7XG5leHBvcnRzLmRlZmF1bHQgPSBkZWZhdWx0XzE7XG5jb25zdCBoZWxsb3dvcmxkXzEgPSByZXF1aXJlKFwiLi9oZWxsb3dvcmxkXCIpO1xuY29uc3QgdW5pdmVyc2Fsc3RvcmFnZV9jbGFzc18xID0gcmVxdWlyZShcIi4vdW5pdmVyc2Fsc3RvcmFnZS5jbGFzc1wiKTtcbi8vIDEuIEluc3RhbnRpYXRlIHRoZSBTaW5nbGV0b24gSEVSRSAoVG9wLWxldmVsIHNjb3BlKS5cbi8vIEJlY2F1c2UgdGhpcyBydW5zIHdoZW4gdGhlIEhvc3Qgc2NyaXB0IGxvYWRzIHRoZSBsaWJyYXJ5LFxuLy8gdGhpcyBzcGVjaWZpYyBpbnN0YW5jZSBjYXB0dXJlcyB0aGUgSG9zdCdzIEBncmFudCBHTSBwZXJtaXNzaW9ucy5cbmNvbnN0IHN0b3JhZ2VJbnN0YW5jZSA9IG5ldyB1bml2ZXJzYWxzdG9yYWdlX2NsYXNzXzEuVW5pdmVyc2FsU3RvcmFnZSgpO1xuZnVuY3Rpb24gZGVmYXVsdF8xKCkge1xuICAgIHJldHVybiB7XG4gICAgICAgIEhlbGxvV29ybGQ6IGhlbGxvd29ybGRfMS5IZWxsb1dvcmxkLFxuICAgICAgICAvLyBSZXR1cm4gdGhlIFBSRS1NQURFIEluc3RhbmNlIChzbyBjb25zdW1lcnMgc2hhcmUgdGhlIEhvc3QncyBzdG9yYWdlKVxuICAgICAgICBVbml2ZXJzYWxTdG9yYWdlOiBzdG9yYWdlSW5zdGFuY2UsXG4gICAgfTtcbn1cbi8vIEVuc3VyZSBnbG9iYWwgYXR0YWNobWVudCBmb3IgV2VicGFjay9Sb2xsdXAgSUlGRSBidWlsZHNcbndpbmRvdy5fX2xpYl9uZGYgPSBmdW5jdGlvbiAoKSB7XG4gICAgcmV0dXJuIHsgSGVsbG9Xb3JsZDogaGVsbG93b3JsZF8xLkhlbGxvV29ybGQsIFVuaXZlcnNhbFN0b3JhZ2U6IHN0b3JhZ2VJbnN0YW5jZSB9O1xufTtcbiJdLCJuYW1lcyI6W10sImlnbm9yZUxpc3QiOltdLCJzb3VyY2VSb290IjoiIn0=