import { Route } from '@/types';
import ofetch from '@/utils/ofetch';
import { load } from 'cheerio';

export const route: Route = {
    path: '/simple',
    categories: ['custom'],
    example: '/example/feed',
    parameters: {},
    features: {
        requireConfig: false,
        requirePuppeteer: false,
        antiCrawler: false,
        supportBT: false,
        supportPodcast: false,
        supportScihub: false,
    },
    name: 'Example Custom Route',
    maintainers: ['custom'],
    handler: async (ctx) => {
        const rootUrl = 'https://example.com';
        const response = await ofetch(rootUrl);
        const $ = load(response);

        const title = $('h1').text();
        const description = $('p').first().text();

        const items = [
            {
                title: title,
                link: rootUrl,
                description: description,
                pubDate: new Date().toUTCString(),
            },
        ];

        return {
            title: 'Example.com Feed',
            link: rootUrl,
            item: items,
        };
    },
};
