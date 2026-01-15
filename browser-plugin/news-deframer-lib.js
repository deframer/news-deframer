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
console.log("📦 Library: Initializing singleton...");
// 1. Instantiate the Singleton HERE.
// This runs once when the library loads, capturing the Host's permissions.
const storageInstance = new universalstorage_class_1.UniversalStorage();
// 2. Export the Function as Default
// Webpack will take this function and assign it to window.__lib_ndf
function default_1() {
    return {
        // Return the Class (if you ever need to make new Hellos)
        HelloWorld: helloworld_1.HelloWorld,
        // Return the PRE-MADE Instance (The Singleton)
        UniversalStorage: storageInstance,
    };
}

})();

window.__lib_ndf = __webpack_exports__["default"];
/******/ })()
;
//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibmV3cy1kZWZyYW1lci1saWIuanMiLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7OztBQUFhO0FBQ2IsOENBQTZDLEVBQUUsYUFBYSxFQUFDO0FBQzdELGtCQUFrQjtBQUNsQjtBQUNBO0FBQ0Esd0JBQXdCLEtBQUs7QUFDN0I7QUFDQTtBQUNBLGtCQUFrQjs7Ozs7Ozs7Ozs7QUNSTDtBQUNiO0FBQ0EsNEJBQTRCLCtEQUErRCxpQkFBaUI7QUFDNUc7QUFDQSxvQ0FBb0MsTUFBTSwrQkFBK0IsWUFBWTtBQUNyRixtQ0FBbUMsTUFBTSxtQ0FBbUMsWUFBWTtBQUN4RixnQ0FBZ0M7QUFDaEM7QUFDQSxLQUFLO0FBQ0w7QUFDQSw4Q0FBNkMsRUFBRSxhQUFhLEVBQUM7QUFDN0Qsd0JBQXdCO0FBQ3hCO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLDJFQUEyRSxjQUFjO0FBQ3pGO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxTQUFTO0FBQ1Q7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLHFCQUFxQjtBQUNyQixpQkFBaUI7QUFDakI7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsU0FBUztBQUNUO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxTQUFTO0FBQ1Q7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLGFBQWE7QUFDYjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsYUFBYTtBQUNiO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxTQUFTO0FBQ1Q7QUFDQTtBQUNBLHdCQUF3QjtBQUN4QjtBQUNBO0FBQ0E7QUFDQTs7Ozs7OztVQ2pJQTtVQUNBOztVQUVBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBOztVQUVBO1VBQ0E7O1VBRUE7VUFDQTtVQUNBOzs7Ozs7Ozs7O0FDNUJhO0FBQ2IsOENBQTZDLEVBQUUsYUFBYSxFQUFDO0FBQzdELGtCQUFlO0FBQ2YscUJBQXFCLG1CQUFPLENBQUMseUNBQWM7QUFDM0MsaUNBQWlDLG1CQUFPLENBQUMsaUVBQTBCO0FBQ25FO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EiLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly9fX2xpYl9uZGYvLi9zcmMvaGVsbG93b3JsZC50cyIsIndlYnBhY2s6Ly9fX2xpYl9uZGYvLi9zcmMvdW5pdmVyc2Fsc3RvcmFnZS5jbGFzcy50cyIsIndlYnBhY2s6Ly9fX2xpYl9uZGYvd2VicGFjay9ib290c3RyYXAiLCJ3ZWJwYWNrOi8vX19saWJfbmRmLy4vc3JjL2luZGV4LnRzIl0sInNvdXJjZXNDb250ZW50IjpbIlwidXNlIHN0cmljdFwiO1xuT2JqZWN0LmRlZmluZVByb3BlcnR5KGV4cG9ydHMsIFwiX19lc01vZHVsZVwiLCB7IHZhbHVlOiB0cnVlIH0pO1xuZXhwb3J0cy5IZWxsb1dvcmxkID0gdm9pZCAwO1xuY2xhc3MgSGVsbG9Xb3JsZCB7XG4gICAgbWVzc2FnZShuYW1lID0gXCJXb3JsZFwiKSB7XG4gICAgICAgIHJldHVybiBgSGVsbG8gJHtuYW1lfWA7XG4gICAgfVxufVxuZXhwb3J0cy5IZWxsb1dvcmxkID0gSGVsbG9Xb3JsZDtcbiIsIlwidXNlIHN0cmljdFwiO1xudmFyIF9fYXdhaXRlciA9ICh0aGlzICYmIHRoaXMuX19hd2FpdGVyKSB8fCBmdW5jdGlvbiAodGhpc0FyZywgX2FyZ3VtZW50cywgUCwgZ2VuZXJhdG9yKSB7XG4gICAgZnVuY3Rpb24gYWRvcHQodmFsdWUpIHsgcmV0dXJuIHZhbHVlIGluc3RhbmNlb2YgUCA/IHZhbHVlIDogbmV3IFAoZnVuY3Rpb24gKHJlc29sdmUpIHsgcmVzb2x2ZSh2YWx1ZSk7IH0pOyB9XG4gICAgcmV0dXJuIG5ldyAoUCB8fCAoUCA9IFByb21pc2UpKShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XG4gICAgICAgIGZ1bmN0aW9uIGZ1bGZpbGxlZCh2YWx1ZSkgeyB0cnkgeyBzdGVwKGdlbmVyYXRvci5uZXh0KHZhbHVlKSk7IH0gY2F0Y2ggKGUpIHsgcmVqZWN0KGUpOyB9IH1cbiAgICAgICAgZnVuY3Rpb24gcmVqZWN0ZWQodmFsdWUpIHsgdHJ5IHsgc3RlcChnZW5lcmF0b3JbXCJ0aHJvd1wiXSh2YWx1ZSkpOyB9IGNhdGNoIChlKSB7IHJlamVjdChlKTsgfSB9XG4gICAgICAgIGZ1bmN0aW9uIHN0ZXAocmVzdWx0KSB7IHJlc3VsdC5kb25lID8gcmVzb2x2ZShyZXN1bHQudmFsdWUpIDogYWRvcHQocmVzdWx0LnZhbHVlKS50aGVuKGZ1bGZpbGxlZCwgcmVqZWN0ZWQpOyB9XG4gICAgICAgIHN0ZXAoKGdlbmVyYXRvciA9IGdlbmVyYXRvci5hcHBseSh0aGlzQXJnLCBfYXJndW1lbnRzIHx8IFtdKSkubmV4dCgpKTtcbiAgICB9KTtcbn07XG5PYmplY3QuZGVmaW5lUHJvcGVydHkoZXhwb3J0cywgXCJfX2VzTW9kdWxlXCIsIHsgdmFsdWU6IHRydWUgfSk7XG5leHBvcnRzLlVuaXZlcnNhbFN0b3JhZ2UgPSB2b2lkIDA7XG5jbGFzcyBVbml2ZXJzYWxTdG9yYWdlIHtcbiAgICAvKipcbiAgICAgKiBTZXQgYSB2YWx1ZS5cbiAgICAgKiBQcmlvcml0eTogQ2hyb21lIFN0b3JhZ2UgLT4gVGFtcGVybW9ua2V5IDUrIChHTS5zZXRWYWx1ZSkgLT4gTG9jYWxTdG9yYWdlXG4gICAgICovXG4gICAgc2V0KGtleSwgdmFsdWUpIHtcbiAgICAgICAgcmV0dXJuIF9fYXdhaXRlcih0aGlzLCB2b2lkIDAsIHZvaWQgMCwgZnVuY3Rpb24qICgpIHtcbiAgICAgICAgICAgIC8vIDEuIEJyb3dzZXIgRXh0ZW5zaW9uXG4gICAgICAgICAgICBpZiAodHlwZW9mIGNocm9tZSAhPT0gJ3VuZGVmaW5lZCcgJiYgY2hyb21lLnN0b3JhZ2UgJiYgY2hyb21lLnN0b3JhZ2UubG9jYWwpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUpID0+IGNocm9tZS5zdG9yYWdlLmxvY2FsLnNldCh7IFtrZXldOiB2YWx1ZSB9LCByZXNvbHZlKSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyAyLiBUYW1wZXJtb25rZXkgNStcbiAgICAgICAgICAgIGlmICh0eXBlb2YgR00gIT09ICd1bmRlZmluZWQnICYmIEdNLnNldFZhbHVlKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIHlpZWxkIEdNLnNldFZhbHVlKGtleSwgdmFsdWUpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gMy4gV2ViIEZhbGxiYWNrXG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIGxvY2FsU3RvcmFnZS5zZXRJdGVtKGtleSwgSlNPTi5zdHJpbmdpZnkodmFsdWUpKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGNhdGNoIChlKSB7XG4gICAgICAgICAgICAgICAgY29uc29sZS53YXJuKFwiVW5pdmVyc2FsU3RvcmFnZTogTG9jYWxTdG9yYWdlIHVuYXZhaWxhYmxlLCB1c2luZyBtZW1vcnkgZmFsbGJhY2suXCIsIGUpO1xuICAgICAgICAgICAgICAgIFVuaXZlcnNhbFN0b3JhZ2UubWVtb3J5U3RvcmFnZVtrZXldID0gdmFsdWU7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICByZXR1cm4gUHJvbWlzZS5yZXNvbHZlKCk7XG4gICAgICAgIH0pO1xuICAgIH1cbiAgICAvKipcbiAgICAgKiBHZXQgYSB2YWx1ZS5cbiAgICAgKi9cbiAgICBnZXQoa2V5LCBkZWZhdWx0VmFsdWUpIHtcbiAgICAgICAgcmV0dXJuIF9fYXdhaXRlcih0aGlzLCB2b2lkIDAsIHZvaWQgMCwgZnVuY3Rpb24qICgpIHtcbiAgICAgICAgICAgIC8vIDEuIEJyb3dzZXIgRXh0ZW5zaW9uXG4gICAgICAgICAgICBpZiAodHlwZW9mIGNocm9tZSAhPT0gJ3VuZGVmaW5lZCcgJiYgY2hyb21lLnN0b3JhZ2UgJiYgY2hyb21lLnN0b3JhZ2UubG9jYWwpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUpID0+IHtcbiAgICAgICAgICAgICAgICAgICAgY2hyb21lLnN0b3JhZ2UubG9jYWwuZ2V0KFtrZXldLCAocmVzKSA9PiB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXNvbHZlKHJlc1trZXldICE9PSB1bmRlZmluZWQgPyByZXNba2V5XSA6IGRlZmF1bHRWYWx1ZSk7XG4gICAgICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gMi4gVGFtcGVybW9ua2V5IDUrXG4gICAgICAgICAgICBpZiAodHlwZW9mIEdNICE9PSAndW5kZWZpbmVkJyAmJiBHTS5nZXRWYWx1ZSkge1xuICAgICAgICAgICAgICAgIHJldHVybiB5aWVsZCBHTS5nZXRWYWx1ZShrZXksIGRlZmF1bHRWYWx1ZSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyAzLiBXZWIgRmFsbGJhY2tcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgY29uc3QgaXRlbSA9IGxvY2FsU3RvcmFnZS5nZXRJdGVtKGtleSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuIFByb21pc2UucmVzb2x2ZShpdGVtID8gSlNPTi5wYXJzZShpdGVtKSA6IGRlZmF1bHRWYWx1ZSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBjYXRjaCAoZSkge1xuICAgICAgICAgICAgICAgIGNvbnNvbGUud2FybihcIlVuaXZlcnNhbFN0b3JhZ2U6IExvY2FsU3RvcmFnZSB1bmF2YWlsYWJsZSwgY2hlY2tpbmcgbWVtb3J5IGZhbGxiYWNrLlwiLCBlKTtcbiAgICAgICAgICAgICAgICBjb25zdCBtZW1WYWwgPSBVbml2ZXJzYWxTdG9yYWdlLm1lbW9yeVN0b3JhZ2Vba2V5XTtcbiAgICAgICAgICAgICAgICByZXR1cm4gUHJvbWlzZS5yZXNvbHZlKG1lbVZhbCAhPT0gdW5kZWZpbmVkID8gbWVtVmFsIDogZGVmYXVsdFZhbHVlKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgfVxuICAgIC8qKlxuICAgICAqIFJlbW92ZSBhIHZhbHVlLlxuICAgICAqL1xuICAgIHJlbW92ZShrZXkpIHtcbiAgICAgICAgcmV0dXJuIF9fYXdhaXRlcih0aGlzLCB2b2lkIDAsIHZvaWQgMCwgZnVuY3Rpb24qICgpIHtcbiAgICAgICAgICAgIGlmICh0eXBlb2YgY2hyb21lICE9PSAndW5kZWZpbmVkJyAmJiBjaHJvbWUuc3RvcmFnZSAmJiBjaHJvbWUuc3RvcmFnZS5sb2NhbCkge1xuICAgICAgICAgICAgICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSkgPT4gY2hyb21lLnN0b3JhZ2UubG9jYWwucmVtb3ZlKGtleSwgcmVzb2x2ZSkpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgaWYgKHR5cGVvZiBHTSAhPT0gJ3VuZGVmaW5lZCcgJiYgR00uZGVsZXRlVmFsdWUpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4geWllbGQgR00uZGVsZXRlVmFsdWUoa2V5KTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgbG9jYWxTdG9yYWdlLnJlbW92ZUl0ZW0oa2V5KTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGNhdGNoIChlKSB7XG4gICAgICAgICAgICAgICAgY29uc29sZS53YXJuKFwiVW5pdmVyc2FsU3RvcmFnZTogTG9jYWxTdG9yYWdlIHVuYXZhaWxhYmxlIChyZW1vdmUpLCBjbGVhcmluZyBtZW1vcnkgZmFsbGJhY2suXCIsIGUpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgZGVsZXRlIFVuaXZlcnNhbFN0b3JhZ2UubWVtb3J5U3RvcmFnZVtrZXldO1xuICAgICAgICAgICAgcmV0dXJuIFByb21pc2UucmVzb2x2ZSgpO1xuICAgICAgICB9KTtcbiAgICB9XG4gICAgLyoqXG4gICAgICogU3luYzogU3Vic2NyaWJlIHRvIGNoYW5nZXMgZnJvbSBvdGhlciBUYWJzIC8gU2NyaXB0cy5cbiAgICAgKi9cbiAgICBzdWJzY3JpYmUoa2V5LCBjYWxsYmFjaykge1xuICAgICAgICAvLyAxLiBCcm93c2VyIEV4dGVuc2lvblxuICAgICAgICBpZiAodHlwZW9mIGNocm9tZSAhPT0gJ3VuZGVmaW5lZCcgJiYgY2hyb21lLnN0b3JhZ2UgJiYgY2hyb21lLnN0b3JhZ2Uub25DaGFuZ2VkKSB7XG4gICAgICAgICAgICBjaHJvbWUuc3RvcmFnZS5vbkNoYW5nZWQuYWRkTGlzdGVuZXIoKGNoYW5nZXMsIG5hbWVzcGFjZSkgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChuYW1lc3BhY2UgPT09ICdsb2NhbCcgJiYgY2hhbmdlc1trZXldKSB7XG4gICAgICAgICAgICAgICAgICAgIGNhbGxiYWNrKGNoYW5nZXNba2V5XS5uZXdWYWx1ZSwgY2hhbmdlc1trZXldLm9sZFZhbHVlKTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9KTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICAvLyAyLiBUYW1wZXJtb25rZXkgNStcbiAgICAgICAgLy8gKFN1cHBvcnRzIEdNLmFkZFZhbHVlQ2hhbmdlTGlzdGVuZXIgZm9yIE1WMy9Bc3luYyBjb21wbGlhbmNlKVxuICAgICAgICBpZiAodHlwZW9mIEdNICE9PSAndW5kZWZpbmVkJyAmJiBHTS5hZGRWYWx1ZUNoYW5nZUxpc3RlbmVyKSB7XG4gICAgICAgICAgICBHTS5hZGRWYWx1ZUNoYW5nZUxpc3RlbmVyKGtleSwgKG5hbWUsIG9sZFZhbCwgbmV3VmFsLCByZW1vdGUpID0+IHtcbiAgICAgICAgICAgICAgICBpZiAocmVtb3RlKVxuICAgICAgICAgICAgICAgICAgICBjYWxsYmFjayhuZXdWYWwsIG9sZFZhbCk7XG4gICAgICAgICAgICB9KTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICAvLyAzLiBXZWIgRmFsbGJhY2tcbiAgICAgICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ3N0b3JhZ2UnLCAoZXZlbnQpID0+IHtcbiAgICAgICAgICAgIGlmIChldmVudC5rZXkgPT09IGtleSkge1xuICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IG4gPSBldmVudC5uZXdWYWx1ZSA/IEpTT04ucGFyc2UoZXZlbnQubmV3VmFsdWUpIDogbnVsbDtcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgbyA9IGV2ZW50Lm9sZFZhbHVlID8gSlNPTi5wYXJzZShldmVudC5vbGRWYWx1ZSkgOiBudWxsO1xuICAgICAgICAgICAgICAgICAgICBjYWxsYmFjayhuLCBvKTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgY2F0Y2ggKF9hKSB7XG4gICAgICAgICAgICAgICAgICAgIGNhbGxiYWNrKGV2ZW50Lm5ld1ZhbHVlLCBldmVudC5vbGRWYWx1ZSk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICB9XG59XG5leHBvcnRzLlVuaXZlcnNhbFN0b3JhZ2UgPSBVbml2ZXJzYWxTdG9yYWdlO1xuLyoqXG4gKiBJbnRlcm5hbCBnbG9iYWwgdmFyaWFibGUgZm9yIG1lbW9yeSBmYWxsYmFjayB3aGVuIExvY2FsU3RvcmFnZSBpcyB1bmF2YWlsYWJsZS5cbiAqL1xuVW5pdmVyc2FsU3RvcmFnZS5tZW1vcnlTdG9yYWdlID0ge307XG4iLCIvLyBUaGUgbW9kdWxlIGNhY2hlXG52YXIgX193ZWJwYWNrX21vZHVsZV9jYWNoZV9fID0ge307XG5cbi8vIFRoZSByZXF1aXJlIGZ1bmN0aW9uXG5mdW5jdGlvbiBfX3dlYnBhY2tfcmVxdWlyZV9fKG1vZHVsZUlkKSB7XG5cdC8vIENoZWNrIGlmIG1vZHVsZSBpcyBpbiBjYWNoZVxuXHR2YXIgY2FjaGVkTW9kdWxlID0gX193ZWJwYWNrX21vZHVsZV9jYWNoZV9fW21vZHVsZUlkXTtcblx0aWYgKGNhY2hlZE1vZHVsZSAhPT0gdW5kZWZpbmVkKSB7XG5cdFx0cmV0dXJuIGNhY2hlZE1vZHVsZS5leHBvcnRzO1xuXHR9XG5cdC8vIENoZWNrIGlmIG1vZHVsZSBleGlzdHMgKGRldmVsb3BtZW50IG9ubHkpXG5cdGlmIChfX3dlYnBhY2tfbW9kdWxlc19fW21vZHVsZUlkXSA9PT0gdW5kZWZpbmVkKSB7XG5cdFx0dmFyIGUgPSBuZXcgRXJyb3IoXCJDYW5ub3QgZmluZCBtb2R1bGUgJ1wiICsgbW9kdWxlSWQgKyBcIidcIik7XG5cdFx0ZS5jb2RlID0gJ01PRFVMRV9OT1RfRk9VTkQnO1xuXHRcdHRocm93IGU7XG5cdH1cblx0Ly8gQ3JlYXRlIGEgbmV3IG1vZHVsZSAoYW5kIHB1dCBpdCBpbnRvIHRoZSBjYWNoZSlcblx0dmFyIG1vZHVsZSA9IF9fd2VicGFja19tb2R1bGVfY2FjaGVfX1ttb2R1bGVJZF0gPSB7XG5cdFx0Ly8gbm8gbW9kdWxlLmlkIG5lZWRlZFxuXHRcdC8vIG5vIG1vZHVsZS5sb2FkZWQgbmVlZGVkXG5cdFx0ZXhwb3J0czoge31cblx0fTtcblxuXHQvLyBFeGVjdXRlIHRoZSBtb2R1bGUgZnVuY3Rpb25cblx0X193ZWJwYWNrX21vZHVsZXNfX1ttb2R1bGVJZF0uY2FsbChtb2R1bGUuZXhwb3J0cywgbW9kdWxlLCBtb2R1bGUuZXhwb3J0cywgX193ZWJwYWNrX3JlcXVpcmVfXyk7XG5cblx0Ly8gUmV0dXJuIHRoZSBleHBvcnRzIG9mIHRoZSBtb2R1bGVcblx0cmV0dXJuIG1vZHVsZS5leHBvcnRzO1xufVxuXG4iLCJcInVzZSBzdHJpY3RcIjtcbk9iamVjdC5kZWZpbmVQcm9wZXJ0eShleHBvcnRzLCBcIl9fZXNNb2R1bGVcIiwgeyB2YWx1ZTogdHJ1ZSB9KTtcbmV4cG9ydHMuZGVmYXVsdCA9IGRlZmF1bHRfMTtcbmNvbnN0IGhlbGxvd29ybGRfMSA9IHJlcXVpcmUoXCIuL2hlbGxvd29ybGRcIik7XG5jb25zdCB1bml2ZXJzYWxzdG9yYWdlX2NsYXNzXzEgPSByZXF1aXJlKFwiLi91bml2ZXJzYWxzdG9yYWdlLmNsYXNzXCIpO1xuY29uc29sZS5sb2coXCLwn5OmIExpYnJhcnk6IEluaXRpYWxpemluZyBzaW5nbGV0b24uLi5cIik7XG4vLyAxLiBJbnN0YW50aWF0ZSB0aGUgU2luZ2xldG9uIEhFUkUuXG4vLyBUaGlzIHJ1bnMgb25jZSB3aGVuIHRoZSBsaWJyYXJ5IGxvYWRzLCBjYXB0dXJpbmcgdGhlIEhvc3QncyBwZXJtaXNzaW9ucy5cbmNvbnN0IHN0b3JhZ2VJbnN0YW5jZSA9IG5ldyB1bml2ZXJzYWxzdG9yYWdlX2NsYXNzXzEuVW5pdmVyc2FsU3RvcmFnZSgpO1xuLy8gMi4gRXhwb3J0IHRoZSBGdW5jdGlvbiBhcyBEZWZhdWx0XG4vLyBXZWJwYWNrIHdpbGwgdGFrZSB0aGlzIGZ1bmN0aW9uIGFuZCBhc3NpZ24gaXQgdG8gd2luZG93Ll9fbGliX25kZlxuZnVuY3Rpb24gZGVmYXVsdF8xKCkge1xuICAgIHJldHVybiB7XG4gICAgICAgIC8vIFJldHVybiB0aGUgQ2xhc3MgKGlmIHlvdSBldmVyIG5lZWQgdG8gbWFrZSBuZXcgSGVsbG9zKVxuICAgICAgICBIZWxsb1dvcmxkOiBoZWxsb3dvcmxkXzEuSGVsbG9Xb3JsZCxcbiAgICAgICAgLy8gUmV0dXJuIHRoZSBQUkUtTUFERSBJbnN0YW5jZSAoVGhlIFNpbmdsZXRvbilcbiAgICAgICAgVW5pdmVyc2FsU3RvcmFnZTogc3RvcmFnZUluc3RhbmNlLFxuICAgIH07XG59XG4iXSwibmFtZXMiOltdLCJpZ25vcmVMaXN0IjpbXSwic291cmNlUm9vdCI6IiJ9