Here is the complete package.

### Part 1: JavaScript Debugger (Browser Console)
Run this in the console on `spiegel.de` to verify the data extraction logic before you run the bridge.

```javascript
// Spiegel.de Debugger
let items = [];

// Find all article containers
let articles = document.querySelectorAll('article');

articles.forEach((article) => {
    let item = {};

    // 1. Get the unique Article ID (Spiegel specific)
    let articleId = article.getAttribute('data-sara-article-id');

    // 2. Find the Headline
    // Prefer aria-label for the clean full text
    item['title'] = article.getAttribute('aria-label');
    if (!item['title']) {
        let h2 = article.querySelector('h2');
        item['title'] = h2 ? h2.innerText.trim() : 'No Title';
    }

    // 3. Find the URL
    // Strategy: Look for a link containing the article ID to ensure we get the right one
    let linkEl = null;
    if (articleId) {
        linkEl = article.querySelector(`a[href*="${articleId}"]`);
    }
    // Fallback: finding the first link inside the h2
    if (!linkEl) {
        linkEl = article.querySelector('h2 a');
    }

    if (linkEl) {
        item['uri'] = linkEl.href;
    } else {
        item['uri'] = 'No URL found';
    }

    // 4. Extract Category (TopMark)
    // Spiegel uses red text above headlines for categories
    let catEl = article.querySelector('[data-target-teaser-el="topMark"]');
    item['category'] = catEl ? catEl.innerText.trim() : 'News';

    // 5. Extract Teaser Text
    let textEl = article.querySelector('[data-target-teaser-el="text"]');
    item['content'] = textEl ? textEl.innerText.trim() : '';

    // 6. Extract Image
    // Spiegel uses lazy loading, so we look for 'data-src' or <source> tags
    let imgEl = article.querySelector('img');
    if (imgEl) {
        // Try data-src first, then srcset, then src
        let src = imgEl.getAttribute('data-src') || imgEl.currentSrc || imgEl.src;
        // Filter out tiny placeholder SVGs
        if (src && !src.startsWith('data:image/svg')) {
            item['image'] = src;
        }
    }

    // 7. Detect Type (Video/Paid)
    let teaserAttr = article.getAttribute('data-target-teaser') || '';
    if (teaserAttr.includes('video') || article.innerHTML.includes('spon-video')) {
        item['type'] = 'Video';
        item['title'] = '[Video] ' + item['title'];
    } else if (article.innerHTML.includes('paid-flag') || article.innerHTML.includes('spon-splus')) {
        item['type'] = 'Paid';
        item['title'] = '[S+] ' + item['title'];
    } else {
        item['type'] = 'Article';
    }

    items.push(item);
});

console.table(items, ['type', 'category', 'title', 'uri']);
```

---

### Part 2: The PHP Bridge Code
Save this file as `SpiegelBridge.php`. Note that the filename **must** match the class name exactly.

```php
<?php
class SpiegelBridge extends BridgeAbstract {
    const NAME = 'Der Spiegel Bridge';
    const URI = 'https://www.spiegel.de';
    const DESCRIPTION = 'Returns the latest news from Spiegel.de';
    const MAINTAINER = 'YourName';
    const CACHE_TIMEOUT = 900; // 15 minutes

    public function collectData() {
        $html = getSimpleHTMLDOM(self::URI);

        // Limit the number of items to prevent timeouts
        $limit = 20;
        $count = 0;

        foreach($html->find('article') as $element) {
            if ($count >= $limit) break;

            $item = array();

            // 1. Title
            $item['title'] = $element->getAttribute('aria-label');
            if (empty($item['title'])) {
                $h2 = $element->find('h2', 0);
                if ($h2) $item['title'] = strip_tags($h2->plaintext);
            }

            // 2. URI
            // Use the data-sara-article-id to find the specific link
            $articleId = $element->getAttribute('data-sara-article-id');
            $urlElement = null;

            if ($articleId) {
                // Find anchor tag containing the ID in href
                $urlElement = $element->find('a[href*="' . $articleId . '"]', 0);
            }

            if (!$urlElement) {
                // Fallback to H2 link
                $urlElement = $element->find('h2 a', 0);
            }

            if ($urlElement) {
                $item['uri'] = $urlElement->href;
                // Handle relative URLs
                if (substr($item['uri'], 0, 4) !== 'http') {
                    $item['uri'] = self::URI . $item['uri'];
                }
            }

            // 3. Content / Teaser
            $summaryNode = $element->find('[data-target-teaser-el="text"]', 0);
            $item['content'] = '';

            if ($summaryNode) {
                $item['content'] .= '<p>' . $summaryNode->plaintext . '</p>';
            }

            // 4. Categories (TopMark)
            $topMark = $element->find('[data-target-teaser-el="topMark"]', 0);
            if ($topMark) {
                $item['categories'] = array($topMark->plaintext);
            }

            // 5. Images
            // Spiegel uses lazy loading (data-src)
            $img = $element->find('img', 0);
            if ($img) {
                $src = $img->getAttribute('data-src');
                if (!$src) $src = $img->src;

                // Filter out base64 SVG placeholders
                if ($src && strpos($src, 'data:image/svg') === false) {
                    $item['enclosures'] = array($src);
                    $item['content'] = '<img src="' . $src . '" /><br>' . $item['content'];
                }
            }

            // 6. Detect Type (Video/Plus) and modify title
            $teaserAttr = $element->getAttribute('data-target-teaser') ?? '';
            $innerHtml = $element->innertext;

            if (strpos($teaserAttr, 'video') !== false || strpos($innerHtml, 'spon-video') !== false) {
                $item['title'] = '[Video] ' . $item['title'];
            } elseif (strpos($innerHtml, 'paid-flag') !== false || strpos($innerHtml, 'plus-paid') !== false) {
                $item['title'] = '[S+] ' . $item['title'];
            }

            // Only add if we found a link
            if (!empty($item['uri'])) {
                $this->items[] = $item;
                $count++;
            }
        }
    }
}
```

---

### Part 3: Installation with Official Docker Container

The official Docker container works by pulling code from the repository. To add your custom bridge, you need to **mount** your local file into the container's `bridges/` folder.

#### Option A: Using Docker Compose (Recommended)

1.  Create a folder on your computer (e.g., `rssbridge`).
2.  Inside that folder, create a subfolder named `custom-bridges`.
3.  Place your `SpiegelBridge.php` file inside `custom-bridges`.
4.  Create a `docker-compose.yml` file in the main folder:

```yaml
version: '3'
services:
  rss-bridge:
    image: rssbridge/rss-bridge:latest
    container_name: rss-bridge
    ports:
      - "3000:80"
    volumes:
      # Mount your custom bridge file into the container
      - ./custom-bridges/SpiegelBridge.php:/app/bridges/SpiegelBridge.php
    restart: unless-stopped
```

5.  Run: `docker-compose up -d`
6.  Open `http://localhost:3000` in your browser. You should see **Der Spiegel Bridge** in the list.

#### Option B: Using Docker Run (Command Line)

If you just want to run it once without a compose file:

1.  Navigate to the folder where `SpiegelBridge.php` is located.
2.  Run the following command (this maps the single file into the container):

```bash
docker run -d \
  --name rss-bridge \
  -p 3000:80 \
  -v $(pwd)/SpiegelBridge.php:/app/bridges/SpiegelBridge.php \
  rssbridge/rss-bridge:latest
```

*(Note: If you are on Windows PowerShell, replace `$(pwd)` with `${PWD}`).*

#### Troubleshooting "Allow List"
If the bridge appears but gives a "Bridge not whitelisted" error when you click it:

1.  You need to enable it in the configuration.
2.  The easiest way in Docker is to allow **all** bridges by setting an environment variable in your `docker-compose.yml`:

```yaml
    environment:
      - RSSBRIDGE_whitelist=*
```

Or just copy `config.default.ini.php` to `config.ini.php` and mount that as a volume as well if you want specific configuration.