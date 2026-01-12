<?php
/*
(function() {
    const containers = document.querySelectorAll('.teaser, .teaser-xs, .teaser-nano');
    const results = [];

    containers.forEach(el => {
        // Robust Extraction
        const titleEl = el.querySelector('.teaser__headline, .teaser-nano__headline, h3');
        const linkEl = el.querySelector('a[href]');
        const summaryEl = el.querySelector('.teaser__shorttext');
        const categoryEl = el.querySelector('.teaser__label, .teaser__topline, .label');

        // High Quality Image Handling (Lazy loading + srcset)
        let imgUrl = '';
        const picture = el.querySelector('picture');
        if (picture) {
            // Get the highest resolution webp or fallback to jpg
            const source = picture.querySelector('source[srcset]');
            if (source) {
                const srcset = source.getAttribute('srcset');
                // Extract the last (usually largest) URL from the srcset
                imgUrl = srcset.split(',').pop().trim().split(' ')[0];
            } else {
                imgUrl = picture.querySelector('img')?.src;
            }
        }

        // Type Detection
        let contentType = 'Standard Article';
        if (el.querySelector('.icon--play') || el.querySelector('[data-v-type="MediaPlayer"]')) {
            contentType = 'Video/Multimedia';
        } else if (el.querySelector('.icon--audio')) {
            contentType = 'Audio/Podcast';
        } else if (el.innerText.toLowerCase().includes('exklusiv')) {
            contentType = 'Investigativ/Exklusiv';
        }

        if (titleEl && linkEl && linkEl.getAttribute('href').startsWith('/')) {
            results.push({
                Title: titleEl.innerText.replace(/\s+/g, ' ').trim(),
                URL: new URL(linkEl.getAttribute('href'), window.location.origin).href,
                Category: categoryEl ? categoryEl.innerText.trim() : 'Nachrichten',
                Summary: summaryEl ? summaryEl.innerText.trim() : 'Kein Teaser-Text vorhanden.',
                Type: contentType,
                Image: imgUrl
            });
        }
    });

    console.table(results);
    console.log(results);
})();
*/
class TagesschauBridge extends BridgeAbstract {
    const NAME = 'Tagesschau.de Bridge (unfinished)';
    const URI = 'https://www.tagesschau.de/';
    const DESCRIPTION = 'Liefert die aktuellen Nachrichten von der Tagesschau Startseite.';
    const MAINTAINER = 'Expert Scraper';
    const CACHE_TIMEOUT = 300; // 5 min

    public function collectData() {
        $html = getSimpleHTMLDOM(self::URI) or returnServerError('Could not request Tagesschau.');

        // Extracting all relevant teaser types
        $teasers = $html->find('.teaser, .teaser-xs, .teaser-nano');

        foreach($teasers as $element) {
            $item = [];

            // Title
            $title = $element->find('.teaser__headline, .teaser-nano__headline, h3', 0);
            if (!$title) continue;

            // Link
            $link = $element->find('a', 0);
            if (!$link) continue;
            $fullUrl = $this->getURI() . ltrim($link->href, '/');

            // Category/Label
            $label = $element->find('.teaser__label, .teaser__topline, .label', 0);
            $item['categories'] = $label ? [trim($label->plaintext)] : ['News'];

            // Content/Summary
            $summary = $element->find('.teaser__shorttext', 0);
            $summaryText = $summary ? trim($summary->plaintext) : '';

            // Image Extraction
            $img = '';
            $picture = $element->find('picture', 0);
            if ($picture) {
                $source = $picture->find('source', 0);
                if ($source && $source->srcset) {
                    // Get the high-res image from srcset (usually the last entry before the comma)
                    $parts = explode(',', $source->srcset);
                    $lastPart = trim(end($parts));
                    $img = explode(' ', $lastPart)[0];
                } else {
                    $imgObj = $picture->find('img', 0);
                    $img = $imgObj ? $imgObj->src : '';
                }
            }

            // Type Detection
            $typePrefix = '';
            if ($element->find('.icon--play', 0)) {
                $typePrefix = ' [Video]';
            } elseif ($element->find('.icon--audio', 0)) {
                $typePrefix = ' [Audio]';
            }

            // Populate Item
            $item['uri'] = $fullUrl;
            $item['title'] = trim($title->plaintext) . $typePrefix;

            $content = '';
            if ($img) {
                $content .= '<img src="' . $img . '" style="max-width:100%; height:auto;" /><br>';
            }
            $content .= '<p>' . $summaryText . '</p>';

            $item['content'] = $content;
            $item['uid'] = md5($fullUrl);

            $this->items[] = $item;
        }

        // Remove duplicates (Tagesschau often repeats items in different sections)
        $this->items = array_map("unserialize", array_unique(array_map("serialize", $this->items)));
    }
}
