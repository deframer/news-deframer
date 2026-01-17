interface Config {
    backendUrl: string;
    username?: string;
    password?: string;
}

export default {
    start: (config: Config) => {
        console.log("News Deframer Library: Hello from the Guest Library!");
        console.log("Configuration:", config);
        document.body.style.border = "5px solid green"; // Visual proof
    }
};
