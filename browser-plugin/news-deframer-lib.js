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
function default_1() {
    return {
        HelloWorld: helloworld_1.HelloWorld,
        UniversalStorage: universalstorage_class_1.UniversalStorage,
    };
}

})();

window.__lib_ndf = __webpack_exports__["default"];
/******/ })()
;
//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibmV3cy1kZWZyYW1lci1saWIuanMiLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7OztBQUFhO0FBQ2IsOENBQTZDLEVBQUUsYUFBYSxFQUFDO0FBQzdELGtCQUFrQjtBQUNsQjtBQUNBO0FBQ0Esd0JBQXdCLEtBQUs7QUFDN0I7QUFDQTtBQUNBLGtCQUFrQjs7Ozs7Ozs7Ozs7QUNSTDtBQUNiO0FBQ0EsNEJBQTRCLCtEQUErRCxpQkFBaUI7QUFDNUc7QUFDQSxvQ0FBb0MsTUFBTSwrQkFBK0IsWUFBWTtBQUNyRixtQ0FBbUMsTUFBTSxtQ0FBbUMsWUFBWTtBQUN4RixnQ0FBZ0M7QUFDaEM7QUFDQSxLQUFLO0FBQ0w7QUFDQSw4Q0FBNkMsRUFBRSxhQUFhLEVBQUM7QUFDN0Qsd0JBQXdCO0FBQ3hCO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLDJFQUEyRSxjQUFjO0FBQ3pGO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxTQUFTO0FBQ1Q7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLHFCQUFxQjtBQUNyQixpQkFBaUI7QUFDakI7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsU0FBUztBQUNUO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxTQUFTO0FBQ1Q7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBLGFBQWE7QUFDYjtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0EsYUFBYTtBQUNiO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSxTQUFTO0FBQ1Q7QUFDQTtBQUNBLHdCQUF3QjtBQUN4QjtBQUNBO0FBQ0E7QUFDQTs7Ozs7OztVQ2pJQTtVQUNBOztVQUVBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBO1VBQ0E7VUFDQTtVQUNBOztVQUVBO1VBQ0E7O1VBRUE7VUFDQTtVQUNBOzs7Ozs7Ozs7O0FDNUJhO0FBQ2IsOENBQTZDLEVBQUUsYUFBYSxFQUFDO0FBQzdELGtCQUFlO0FBQ2YscUJBQXFCLG1CQUFPLENBQUMseUNBQWM7QUFDM0MsaUNBQWlDLG1CQUFPLENBQUMsaUVBQTBCO0FBQ25FO0FBQ0E7QUFDQTtBQUNBO0FBQ0E7QUFDQSIsInNvdXJjZXMiOlsid2VicGFjazovL19fbGliX25kZi8uL3NyYy9oZWxsb3dvcmxkLnRzIiwid2VicGFjazovL19fbGliX25kZi8uL3NyYy91bml2ZXJzYWxzdG9yYWdlLmNsYXNzLnRzIiwid2VicGFjazovL19fbGliX25kZi93ZWJwYWNrL2Jvb3RzdHJhcCIsIndlYnBhY2s6Ly9fX2xpYl9uZGYvLi9zcmMvaW5kZXgudHMiXSwic291cmNlc0NvbnRlbnQiOlsiXCJ1c2Ugc3RyaWN0XCI7XG5PYmplY3QuZGVmaW5lUHJvcGVydHkoZXhwb3J0cywgXCJfX2VzTW9kdWxlXCIsIHsgdmFsdWU6IHRydWUgfSk7XG5leHBvcnRzLkhlbGxvV29ybGQgPSB2b2lkIDA7XG5jbGFzcyBIZWxsb1dvcmxkIHtcbiAgICBtZXNzYWdlKG5hbWUgPSBcIldvcmxkXCIpIHtcbiAgICAgICAgcmV0dXJuIGBIZWxsbyAke25hbWV9YDtcbiAgICB9XG59XG5leHBvcnRzLkhlbGxvV29ybGQgPSBIZWxsb1dvcmxkO1xuIiwiXCJ1c2Ugc3RyaWN0XCI7XG52YXIgX19hd2FpdGVyID0gKHRoaXMgJiYgdGhpcy5fX2F3YWl0ZXIpIHx8IGZ1bmN0aW9uICh0aGlzQXJnLCBfYXJndW1lbnRzLCBQLCBnZW5lcmF0b3IpIHtcbiAgICBmdW5jdGlvbiBhZG9wdCh2YWx1ZSkgeyByZXR1cm4gdmFsdWUgaW5zdGFuY2VvZiBQID8gdmFsdWUgOiBuZXcgUChmdW5jdGlvbiAocmVzb2x2ZSkgeyByZXNvbHZlKHZhbHVlKTsgfSk7IH1cbiAgICByZXR1cm4gbmV3IChQIHx8IChQID0gUHJvbWlzZSkpKGZ1bmN0aW9uIChyZXNvbHZlLCByZWplY3QpIHtcbiAgICAgICAgZnVuY3Rpb24gZnVsZmlsbGVkKHZhbHVlKSB7IHRyeSB7IHN0ZXAoZ2VuZXJhdG9yLm5leHQodmFsdWUpKTsgfSBjYXRjaCAoZSkgeyByZWplY3QoZSk7IH0gfVxuICAgICAgICBmdW5jdGlvbiByZWplY3RlZCh2YWx1ZSkgeyB0cnkgeyBzdGVwKGdlbmVyYXRvcltcInRocm93XCJdKHZhbHVlKSk7IH0gY2F0Y2ggKGUpIHsgcmVqZWN0KGUpOyB9IH1cbiAgICAgICAgZnVuY3Rpb24gc3RlcChyZXN1bHQpIHsgcmVzdWx0LmRvbmUgPyByZXNvbHZlKHJlc3VsdC52YWx1ZSkgOiBhZG9wdChyZXN1bHQudmFsdWUpLnRoZW4oZnVsZmlsbGVkLCByZWplY3RlZCk7IH1cbiAgICAgICAgc3RlcCgoZ2VuZXJhdG9yID0gZ2VuZXJhdG9yLmFwcGx5KHRoaXNBcmcsIF9hcmd1bWVudHMgfHwgW10pKS5uZXh0KCkpO1xuICAgIH0pO1xufTtcbk9iamVjdC5kZWZpbmVQcm9wZXJ0eShleHBvcnRzLCBcIl9fZXNNb2R1bGVcIiwgeyB2YWx1ZTogdHJ1ZSB9KTtcbmV4cG9ydHMuVW5pdmVyc2FsU3RvcmFnZSA9IHZvaWQgMDtcbmNsYXNzIFVuaXZlcnNhbFN0b3JhZ2Uge1xuICAgIC8qKlxuICAgICAqIFNldCBhIHZhbHVlLlxuICAgICAqIFByaW9yaXR5OiBDaHJvbWUgU3RvcmFnZSAtPiBUYW1wZXJtb25rZXkgNSsgKEdNLnNldFZhbHVlKSAtPiBMb2NhbFN0b3JhZ2VcbiAgICAgKi9cbiAgICBzZXQoa2V5LCB2YWx1ZSkge1xuICAgICAgICByZXR1cm4gX19hd2FpdGVyKHRoaXMsIHZvaWQgMCwgdm9pZCAwLCBmdW5jdGlvbiogKCkge1xuICAgICAgICAgICAgLy8gMS4gQnJvd3NlciBFeHRlbnNpb25cbiAgICAgICAgICAgIGlmICh0eXBlb2YgY2hyb21lICE9PSAndW5kZWZpbmVkJyAmJiBjaHJvbWUuc3RvcmFnZSAmJiBjaHJvbWUuc3RvcmFnZS5sb2NhbCkge1xuICAgICAgICAgICAgICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSkgPT4gY2hyb21lLnN0b3JhZ2UubG9jYWwuc2V0KHsgW2tleV06IHZhbHVlIH0sIHJlc29sdmUpKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIDIuIFRhbXBlcm1vbmtleSA1K1xuICAgICAgICAgICAgaWYgKHR5cGVvZiBHTSAhPT0gJ3VuZGVmaW5lZCcgJiYgR00uc2V0VmFsdWUpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4geWllbGQgR00uc2V0VmFsdWUoa2V5LCB2YWx1ZSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyAzLiBXZWIgRmFsbGJhY2tcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgbG9jYWxTdG9yYWdlLnNldEl0ZW0oa2V5LCBKU09OLnN0cmluZ2lmeSh2YWx1ZSkpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2F0Y2ggKGUpIHtcbiAgICAgICAgICAgICAgICBjb25zb2xlLndhcm4oXCJVbml2ZXJzYWxTdG9yYWdlOiBMb2NhbFN0b3JhZ2UgdW5hdmFpbGFibGUsIHVzaW5nIG1lbW9yeSBmYWxsYmFjay5cIiwgZSk7XG4gICAgICAgICAgICAgICAgVW5pdmVyc2FsU3RvcmFnZS5tZW1vcnlTdG9yYWdlW2tleV0gPSB2YWx1ZTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHJldHVybiBQcm9taXNlLnJlc29sdmUoKTtcbiAgICAgICAgfSk7XG4gICAgfVxuICAgIC8qKlxuICAgICAqIEdldCBhIHZhbHVlLlxuICAgICAqL1xuICAgIGdldChrZXksIGRlZmF1bHRWYWx1ZSkge1xuICAgICAgICByZXR1cm4gX19hd2FpdGVyKHRoaXMsIHZvaWQgMCwgdm9pZCAwLCBmdW5jdGlvbiogKCkge1xuICAgICAgICAgICAgLy8gMS4gQnJvd3NlciBFeHRlbnNpb25cbiAgICAgICAgICAgIGlmICh0eXBlb2YgY2hyb21lICE9PSAndW5kZWZpbmVkJyAmJiBjaHJvbWUuc3RvcmFnZSAmJiBjaHJvbWUuc3RvcmFnZS5sb2NhbCkge1xuICAgICAgICAgICAgICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSkgPT4ge1xuICAgICAgICAgICAgICAgICAgICBjaHJvbWUuc3RvcmFnZS5sb2NhbC5nZXQoW2tleV0sIChyZXMpID0+IHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlc29sdmUocmVzW2tleV0gIT09IHVuZGVmaW5lZCA/IHJlc1trZXldIDogZGVmYXVsdFZhbHVlKTtcbiAgICAgICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyAyLiBUYW1wZXJtb25rZXkgNStcbiAgICAgICAgICAgIGlmICh0eXBlb2YgR00gIT09ICd1bmRlZmluZWQnICYmIEdNLmdldFZhbHVlKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIHlpZWxkIEdNLmdldFZhbHVlKGtleSwgZGVmYXVsdFZhbHVlKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIDMuIFdlYiBGYWxsYmFja1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICBjb25zdCBpdGVtID0gbG9jYWxTdG9yYWdlLmdldEl0ZW0oa2V5KTtcbiAgICAgICAgICAgICAgICByZXR1cm4gUHJvbWlzZS5yZXNvbHZlKGl0ZW0gPyBKU09OLnBhcnNlKGl0ZW0pIDogZGVmYXVsdFZhbHVlKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGNhdGNoIChlKSB7XG4gICAgICAgICAgICAgICAgY29uc29sZS53YXJuKFwiVW5pdmVyc2FsU3RvcmFnZTogTG9jYWxTdG9yYWdlIHVuYXZhaWxhYmxlLCBjaGVja2luZyBtZW1vcnkgZmFsbGJhY2suXCIsIGUpO1xuICAgICAgICAgICAgICAgIGNvbnN0IG1lbVZhbCA9IFVuaXZlcnNhbFN0b3JhZ2UubWVtb3J5U3RvcmFnZVtrZXldO1xuICAgICAgICAgICAgICAgIHJldHVybiBQcm9taXNlLnJlc29sdmUobWVtVmFsICE9PSB1bmRlZmluZWQgPyBtZW1WYWwgOiBkZWZhdWx0VmFsdWUpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICB9XG4gICAgLyoqXG4gICAgICogUmVtb3ZlIGEgdmFsdWUuXG4gICAgICovXG4gICAgcmVtb3ZlKGtleSkge1xuICAgICAgICByZXR1cm4gX19hd2FpdGVyKHRoaXMsIHZvaWQgMCwgdm9pZCAwLCBmdW5jdGlvbiogKCkge1xuICAgICAgICAgICAgaWYgKHR5cGVvZiBjaHJvbWUgIT09ICd1bmRlZmluZWQnICYmIGNocm9tZS5zdG9yYWdlICYmIGNocm9tZS5zdG9yYWdlLmxvY2FsKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlKSA9PiBjaHJvbWUuc3RvcmFnZS5sb2NhbC5yZW1vdmUoa2V5LCByZXNvbHZlKSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBpZiAodHlwZW9mIEdNICE9PSAndW5kZWZpbmVkJyAmJiBHTS5kZWxldGVWYWx1ZSkge1xuICAgICAgICAgICAgICAgIHJldHVybiB5aWVsZCBHTS5kZWxldGVWYWx1ZShrZXkpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICBsb2NhbFN0b3JhZ2UucmVtb3ZlSXRlbShrZXkpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2F0Y2ggKGUpIHtcbiAgICAgICAgICAgICAgICBjb25zb2xlLndhcm4oXCJVbml2ZXJzYWxTdG9yYWdlOiBMb2NhbFN0b3JhZ2UgdW5hdmFpbGFibGUgKHJlbW92ZSksIGNsZWFyaW5nIG1lbW9yeSBmYWxsYmFjay5cIiwgZSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBkZWxldGUgVW5pdmVyc2FsU3RvcmFnZS5tZW1vcnlTdG9yYWdlW2tleV07XG4gICAgICAgICAgICByZXR1cm4gUHJvbWlzZS5yZXNvbHZlKCk7XG4gICAgICAgIH0pO1xuICAgIH1cbiAgICAvKipcbiAgICAgKiBTeW5jOiBTdWJzY3JpYmUgdG8gY2hhbmdlcyBmcm9tIG90aGVyIFRhYnMgLyBTY3JpcHRzLlxuICAgICAqL1xuICAgIHN1YnNjcmliZShrZXksIGNhbGxiYWNrKSB7XG4gICAgICAgIC8vIDEuIEJyb3dzZXIgRXh0ZW5zaW9uXG4gICAgICAgIGlmICh0eXBlb2YgY2hyb21lICE9PSAndW5kZWZpbmVkJyAmJiBjaHJvbWUuc3RvcmFnZSAmJiBjaHJvbWUuc3RvcmFnZS5vbkNoYW5nZWQpIHtcbiAgICAgICAgICAgIGNocm9tZS5zdG9yYWdlLm9uQ2hhbmdlZC5hZGRMaXN0ZW5lcigoY2hhbmdlcywgbmFtZXNwYWNlKSA9PiB7XG4gICAgICAgICAgICAgICAgaWYgKG5hbWVzcGFjZSA9PT0gJ2xvY2FsJyAmJiBjaGFuZ2VzW2tleV0pIHtcbiAgICAgICAgICAgICAgICAgICAgY2FsbGJhY2soY2hhbmdlc1trZXldLm5ld1ZhbHVlLCBjaGFuZ2VzW2tleV0ub2xkVmFsdWUpO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgICAgIC8vIDIuIFRhbXBlcm1vbmtleSA1K1xuICAgICAgICAvLyAoU3VwcG9ydHMgR00uYWRkVmFsdWVDaGFuZ2VMaXN0ZW5lciBmb3IgTVYzL0FzeW5jIGNvbXBsaWFuY2UpXG4gICAgICAgIGlmICh0eXBlb2YgR00gIT09ICd1bmRlZmluZWQnICYmIEdNLmFkZFZhbHVlQ2hhbmdlTGlzdGVuZXIpIHtcbiAgICAgICAgICAgIEdNLmFkZFZhbHVlQ2hhbmdlTGlzdGVuZXIoa2V5LCAobmFtZSwgb2xkVmFsLCBuZXdWYWwsIHJlbW90ZSkgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZW1vdGUpXG4gICAgICAgICAgICAgICAgICAgIGNhbGxiYWNrKG5ld1ZhbCwgb2xkVmFsKTtcbiAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgICAgIC8vIDMuIFdlYiBGYWxsYmFja1xuICAgICAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignc3RvcmFnZScsIChldmVudCkgPT4ge1xuICAgICAgICAgICAgaWYgKGV2ZW50LmtleSA9PT0ga2V5KSB7XG4gICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgbiA9IGV2ZW50Lm5ld1ZhbHVlID8gSlNPTi5wYXJzZShldmVudC5uZXdWYWx1ZSkgOiBudWxsO1xuICAgICAgICAgICAgICAgICAgICBjb25zdCBvID0gZXZlbnQub2xkVmFsdWUgPyBKU09OLnBhcnNlKGV2ZW50Lm9sZFZhbHVlKSA6IG51bGw7XG4gICAgICAgICAgICAgICAgICAgIGNhbGxiYWNrKG4sIG8pO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICBjYXRjaCAoX2EpIHtcbiAgICAgICAgICAgICAgICAgICAgY2FsbGJhY2soZXZlbnQubmV3VmFsdWUsIGV2ZW50Lm9sZFZhbHVlKTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuICAgIH1cbn1cbmV4cG9ydHMuVW5pdmVyc2FsU3RvcmFnZSA9IFVuaXZlcnNhbFN0b3JhZ2U7XG4vKipcbiAqIEludGVybmFsIGdsb2JhbCB2YXJpYWJsZSBmb3IgbWVtb3J5IGZhbGxiYWNrIHdoZW4gTG9jYWxTdG9yYWdlIGlzIHVuYXZhaWxhYmxlLlxuICovXG5Vbml2ZXJzYWxTdG9yYWdlLm1lbW9yeVN0b3JhZ2UgPSB7fTtcbiIsIi8vIFRoZSBtb2R1bGUgY2FjaGVcbnZhciBfX3dlYnBhY2tfbW9kdWxlX2NhY2hlX18gPSB7fTtcblxuLy8gVGhlIHJlcXVpcmUgZnVuY3Rpb25cbmZ1bmN0aW9uIF9fd2VicGFja19yZXF1aXJlX18obW9kdWxlSWQpIHtcblx0Ly8gQ2hlY2sgaWYgbW9kdWxlIGlzIGluIGNhY2hlXG5cdHZhciBjYWNoZWRNb2R1bGUgPSBfX3dlYnBhY2tfbW9kdWxlX2NhY2hlX19bbW9kdWxlSWRdO1xuXHRpZiAoY2FjaGVkTW9kdWxlICE9PSB1bmRlZmluZWQpIHtcblx0XHRyZXR1cm4gY2FjaGVkTW9kdWxlLmV4cG9ydHM7XG5cdH1cblx0Ly8gQ2hlY2sgaWYgbW9kdWxlIGV4aXN0cyAoZGV2ZWxvcG1lbnQgb25seSlcblx0aWYgKF9fd2VicGFja19tb2R1bGVzX19bbW9kdWxlSWRdID09PSB1bmRlZmluZWQpIHtcblx0XHR2YXIgZSA9IG5ldyBFcnJvcihcIkNhbm5vdCBmaW5kIG1vZHVsZSAnXCIgKyBtb2R1bGVJZCArIFwiJ1wiKTtcblx0XHRlLmNvZGUgPSAnTU9EVUxFX05PVF9GT1VORCc7XG5cdFx0dGhyb3cgZTtcblx0fVxuXHQvLyBDcmVhdGUgYSBuZXcgbW9kdWxlIChhbmQgcHV0IGl0IGludG8gdGhlIGNhY2hlKVxuXHR2YXIgbW9kdWxlID0gX193ZWJwYWNrX21vZHVsZV9jYWNoZV9fW21vZHVsZUlkXSA9IHtcblx0XHQvLyBubyBtb2R1bGUuaWQgbmVlZGVkXG5cdFx0Ly8gbm8gbW9kdWxlLmxvYWRlZCBuZWVkZWRcblx0XHRleHBvcnRzOiB7fVxuXHR9O1xuXG5cdC8vIEV4ZWN1dGUgdGhlIG1vZHVsZSBmdW5jdGlvblxuXHRfX3dlYnBhY2tfbW9kdWxlc19fW21vZHVsZUlkXS5jYWxsKG1vZHVsZS5leHBvcnRzLCBtb2R1bGUsIG1vZHVsZS5leHBvcnRzLCBfX3dlYnBhY2tfcmVxdWlyZV9fKTtcblxuXHQvLyBSZXR1cm4gdGhlIGV4cG9ydHMgb2YgdGhlIG1vZHVsZVxuXHRyZXR1cm4gbW9kdWxlLmV4cG9ydHM7XG59XG5cbiIsIlwidXNlIHN0cmljdFwiO1xuT2JqZWN0LmRlZmluZVByb3BlcnR5KGV4cG9ydHMsIFwiX19lc01vZHVsZVwiLCB7IHZhbHVlOiB0cnVlIH0pO1xuZXhwb3J0cy5kZWZhdWx0ID0gZGVmYXVsdF8xO1xuY29uc3QgaGVsbG93b3JsZF8xID0gcmVxdWlyZShcIi4vaGVsbG93b3JsZFwiKTtcbmNvbnN0IHVuaXZlcnNhbHN0b3JhZ2VfY2xhc3NfMSA9IHJlcXVpcmUoXCIuL3VuaXZlcnNhbHN0b3JhZ2UuY2xhc3NcIik7XG5mdW5jdGlvbiBkZWZhdWx0XzEoKSB7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgSGVsbG9Xb3JsZDogaGVsbG93b3JsZF8xLkhlbGxvV29ybGQsXG4gICAgICAgIFVuaXZlcnNhbFN0b3JhZ2U6IHVuaXZlcnNhbHN0b3JhZ2VfY2xhc3NfMS5Vbml2ZXJzYWxTdG9yYWdlLFxuICAgIH07XG59XG4iXSwibmFtZXMiOltdLCJpZ25vcmVMaXN0IjpbXSwic291cmNlUm9vdCI6IiJ9