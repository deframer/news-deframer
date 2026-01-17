interface Config {
    backendUrl: string;
    username?: string;
    password?: string;
}

export default {
    start: (config: Config) => {
        // check if the current windows url is in
        // backendUrl +
        console.log("News Deframer Library: Hello from the Guest Library!");
        console.log("Configuration:", config);
        document.body.style.border = "5px solid green"; // Visual proof
    }
};
