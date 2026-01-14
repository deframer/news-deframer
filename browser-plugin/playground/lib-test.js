// this library is hosted on a server
// we dynamically download them from the
// browser plugin

window._lib_test = function () {
    function helloWorld() {
        console.log("Hello World from Library");
    }

    return {
        helloWorld: helloWorld
    };
};
