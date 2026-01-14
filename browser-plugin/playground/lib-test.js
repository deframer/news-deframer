// this library is hosted on a server
// we dynamicly download them from the
// broowser plugin

window._lib_test = function () {
    function helloWorld() {
        console.log("Hello World from Library");
    }

    return {
        helloWorld: helloWorld
    };
};
