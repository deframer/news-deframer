// server sided

window._lib_test = function () {
    function helloWorld() {
        console.log("Hello World from Library");
    }

    return {
        helloWorld: helloWorld
    };
};
