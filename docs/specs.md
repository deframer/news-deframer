# Specification

## Test and Helper Systems

### FreshRSS

- **Purpose**:
  - Creates RSS Feeds of Websites that don't have them.
  - Creates better RSS Feeds of Websites that have a limited RSS feed support
- **Technology*:
  - PHP Service
  - Downloads the `index.html`
  - Directly parses XPath (no rendering of Javascript) - workaround - directly access Json
  - Requires per upstream domain PHP Classes / simple to create
- **Public**: No
- **Service URL**:
    - Internal: http://freshrss
    - External: http://localhost:8001

### Dummy News Site

- **Purpose**:
  - A Dummy Website that looks like a News Site
  - The Website offers RSS feeds
  - Articles and Images can be added via API
- **Technology*:
  - Word Press Stack
  - Free `ColorMag` Theme: <https://wordpress.org/themes/colormag/>
- **Public**: No
- **Service URL**:
    - Internal: http://wordpress
    - External: http://wordpress:8002

## News Deframer

### RSS Deframer Proxy

- **Purpose**:
  - Deframing Proxy for any RSS Feeds.
  - Any Input RSS feed can be used and should be deframed.
```bash
ORIGINAL_URL="https://rss.nytimes.com/services/xml/rss/nyt/World.xml"
URL_ENCODED=$(echo -n "${ORIGINAL_URL}" | jq -sRr @uri)
REPLACEMENT_URL="https://deframer.example.com/?url=${URL_ENCODED&embedded=true&max_score=0.5"
```
- **Result**
  - The result is a valid RSS 2.0 document.
  - It has new fields as described in the <ALGORITHM.md> document.
- **Constraints**
  - Cache the results (full XML Feed for a domain) is served from memory. Every 15 min we poll the upstream RSS to check if we have new items.
  - Cache the AI calculated results. An `<item>` of the RSS document might be already deframed.
- **Scaling**
  - We want to scale this service to server 10.000 Users at a single point of time
  - We want to support 5 New Networks - an RSS feed has 30 Items. We assume we poll it every 15 min.
  - We assume of the 30 items a max of 50% are new (or needs updates).
- **Additional Features**
  - The service will have an additional endpoint.
  - Except direclty getting XML, we want to know if we have a deframing result of an URL that is known by a Link of an Item.
  - We can early break if we don't know the base domain e.g.
  - The result is a Json Document / status code if the domain or upstream feed is unknown.
  - We get a Json representation of the enhanced `<item>` as described in the algorithm document.
```bash
REPLACEMENT_URL="https://api.example.com/?url=${URL_ENCODED&embedded=true&max_score=0.5"
```
- **Open**:
  - I want to write this in golang (all possible parts)
  - I don't know how many services we should use e.g. ingester, poller, http server, ...
  - I have no idea where to put API keys but also the list of upstream RSS Feeds we should support. It needs to be configured by the admin - the user can send us any requests
  - The system needs to run with single docker instances on a development computer
  - The Docker containers need to scale on kubernetes
