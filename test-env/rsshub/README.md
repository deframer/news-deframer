# Custom RSSHub Routes

## Mapping Logic

**URL Pattern:** `http://localhost:8004/<NAMESPACE>/<ROUTE_PATH>`

1.  **NAMESPACE**: The folder name inside `custom_routes/` (e.g., `example`).
    *   *Metadata defined in `namespace.ts`.*
2.  **ROUTE_PATH**: The `path` property exported in your route file (e.g., `feed.ts`).

## Example

*   **File:** `custom_routes/example/feed.ts`
    ```typescript
    export const route = { path: '/feed', ... }
    ```
*   **Access:** <http://localhost:8004/example/feed>

Build in examples:

- <http://localhost:8004/github/issue/DIYgod/RSSHub>
- <http://localhost:8004/youtube/user/PewDiePie>
- more <https://github.com/DIYgod/RSSHub/tree/master/lib/routes>


## Apply Changes

```bash
docker compose up -d --build
```
