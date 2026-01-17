import { ApiClient, Config } from './api';

export default {
    start: async (config: Config) => {
        try {
            const client = new ApiClient(config);
            const domains = await client.getDomains();
            const hostname = window.location.hostname;
            const isMonitored = domains.some(d => hostname === d || hostname.endsWith('.' + d));

            if (isMonitored) {
                console.log("News Deframer Library: Hello from the Guest Library!");
                console.log("Configuration:", config);
                document.body.style.border = "5px solid green"; // Visual proof
            }
        } catch (e) {
            console.error("News Deframer Library Error:", e);
        }
    }
};
