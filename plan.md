# Plan

## Goal

Keep RSS as an input format, but remove the RSS proxy and stop enhancing rendered items in the sync path.

## Steps

- [x] Remove the RSS proxy contract and service implementation.
- [ ] Split the current worker into two explicit modes: an ingester and a thinker.
- [ ] Stop injecting deframer-specific `<item>` extensions during sync.
- [ ] Keep raw `title` and `description` values intact in stored items.
- [ ] Move thinker processing into its own loop so sync and analysis run separately.
- [ ] Rebuild cached feed output so existing stored feeds no longer carry original-title/description deframer tags.
- [ ] Update tests and docs to match the new text-first pipeline.

## Expected Result

- RSS remains an upstream source.
- The service no longer proxies RSS responses.
- Sync stores raw item text and metadata only.
- Thinking happens asynchronously in a separate worker loop.


## Anything else

- [x] Update documentation
