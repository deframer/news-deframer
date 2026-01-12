# Test and Helper Systems

## 1.1 RSS-Bridge
- **Purpose**:
  - Serves as the "Connector" logic.
  - Generates standard RSS Feeds for websites that lack them.
  - Sanitizes/Standardizes broken or limited upstream feeds before the Deframer touches them.
- **Technology**:
  - PHP Service (Community Standard).
  - Fetches `index.html`, parses via XPath/CSS Selectors.
  - **Constraint**: Does not render JavaScript (server-side scraping only).
- **Network**:
  - **Internal URL**: `http://rssbridge` (Accessible by Deframer).
  - **External URL**: `http://localhost:8002` (Accessible by Developer).

## 1.2 Dummy News Site
- **Purpose**:
  - Acts as the "Upstream Source."
  - Provides a controllable environment to test edge cases (broken images, paywalls, weird encodings).
  - Allows injection of specific test articles via API.
- **Technology**:
  - WordPress Stack.
  - Theme: `ColorMag` (Simulates a standard magazine layout).
- **Network**:
  - **Internal URL**: `http://wordpress` (Accessible by Deframer).
  - **External URL**: `http://localhost:8003` (Accessible by Developer).

## 1.3 FreshRSS
- **Purpose**:
  - Simple Web based RSS Reader
  - Can read internal and external RSS Feeds
- **Technology**:
  - Docker
- **Network**:
  - **External URL**: `http://localhost:8001` (Accessible by Developer).

