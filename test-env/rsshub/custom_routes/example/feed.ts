import { Route } from '@/types';
import puppeteer from '@/utils/puppeteer';
import { load } from 'cheerio';

export const route: Route = {
    path: '/feed',
    categories: ['custom'],
    example: '/example/feed',
    parameters: {},
    features: {
        requireConfig: false,
        requirePuppeteer: true,
        antiCrawler: false,
        supportBT: false,
        supportPodcast: false,
        supportScihub: false,
    },
    name: 'Example Custom Route (Puppeteer)',
    maintainers: ['custom'],
    handler: async (ctx) => {
        const rootUrl = 'https://example.com';
        const browser = await puppeteer();
        const page = await browser.newPage();
        
        await page.setRequestInterception(true);
        page.on('request', (request) => {
            request.resourceType() === 'document' || request.resourceType() === 'script' 
                ? request.continue() 
                : request.abort();
        });

        await page.goto(rootUrl, {
            waitUntil: 'domcontentloaded',
        });

        const html = await page.content();
        await page.close();

        const $ = load(html);
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
            title: 'Example.com Feed (Puppeteer)',
            link: rootUrl,
            item: items,
        };
    },
};
