Act as an expert in Web Scraping and RSS-Bridge development. I will provide the **HTML source code** of a website below.

Please analyze the DOM structure and generate the following two deliverables:

**1. JavaScript Debugging Snippet:**
Write a script I can paste into the Chrome/Firefox browser console to verify the CSS selectors.
*   **Loop:** Identify the main container for news items/articles.
*   ** robust Extraction:** Extract the Title, URL, Summary/Teaser, and Category. Look for specific attributes (like `data-id`, `aria-label`, or `data-teaser`) which are often cleaner than standard `innerText`.
*   **Images:** Handle lazy loading. Check for `data-src`, `srcset`, or `<picture>` tags before falling back to `src` to ensure high-quality images are found.
*   **Type Detection:** Attempt to detect if an item is a Video, Paid content (Plus), or standard article based on classes, icons (SVGs), or attributes.
*   **Output:** Print the result using `console.table()`.

**2. PHP Bridge Class:**
Write the complete PHP class for `rss-bridge` based on the logic above.
*   Inherit from `BridgeAbstract`.
*   Use `$html = getSimpleHTMLDOM($url);` and standard find methods.
*   Translate the JS logic into PHP (e.g., handling relative URLs, extracting attributes).
*   **No Docker or Installation instructions**—provide only the code.

**Here is the HTML source:**

[PASTE YOUR HTML HERE]
