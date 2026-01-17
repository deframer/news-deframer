import { ApiClient, Config } from './api';

export default {
    start: async (config: Config) => {
        try {
            const client = new ApiClient(config);
            const domains = await client.getDomains();
            const hostname = window.location.hostname;
            const isMonitored = domains.some(d => hostname === d || hostname.endsWith('.' + d));

            if(!isMonitored) return;

            console.log("News Deframer Library: You are on an active domain");
            document.body.style.border = "15px dashed green"; // Visual proof

        } catch (e) {
            console.error("News Deframer Library Error:", e);
        }
    }
};
