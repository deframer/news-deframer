# Browser UI Refactor Plan

## Goal

Refactor the browser extension UI so the popup becomes a tiny control surface and the current long-form options UI becomes a fullscreen browser page.

This is a browser-first refactor only. It prepares the codebase for later shared frontend work, but does not introduce mobile separation yet.

## Target UX

### Popup

The popup should stay very small and only provide quick access:

- enable or disable toggle
- connected status
- open settings page
- project link

The popup should not contain the full settings form.

### Fullscreen settings page

The fullscreen settings page becomes the main browser configuration surface.

It should contain the existing settings content:

- backend URL
- username and password
- language
- theme
- search engine URL
- connection test

It should behave like a normal configuration page.

## Important rule

Settings changes made in the fullscreen settings page should only persist configuration.

They should not:

- reload the active tab
- execute scripts in the active tab
- force-refresh the current page
- trigger bypass logic
- attempt to immediately re-apply the extension to the current page

Changes from the fullscreen settings page apply to future page loads and future extension activity.

## Why this direction

The current `src/host/ui/options.tsx` mixes several concerns:

- configuration form
- connection checking
- enable or disable control
- tab update logic
- page reload logic
- content-script coordination

That makes the UI brittle and hard to reorganize.

The refactor should separate:

- configuration UI
- quick popup controls
- browser runtime behavior

## Proposed browser-only structure

```text
src/
  host/
    index.html
    manifest.json
    background.ts
    content.ts

    popup/
      index.tsx
      Popup.tsx

    pages/
      SettingsPage.tsx

    components/
      ToggleSwitch.tsx
      SettingsForm.tsx
      SettingsConnectionSection.tsx
      SettingsAppearanceSection.tsx
      SettingsLanguageSection.tsx
      QuickActions.tsx

  ndf/
    pages/
    components/
    utils/

  shared/
    i18n.ts
    theme.ts
    logger.ts
    settings.ts
    domain-cache.ts
```

## Responsibilities

### `src/host/popup/`

Tiny extension popup only.

Owns:

- enable or disable toggle
- connected status badge
- open settings page action
- project link action

Should not own:

- full settings form
- long-form configuration editing
- browser page update orchestration beyond the quick toggle behavior

### `src/host/pages/SettingsPage.tsx`

Fullscreen browser page for extension settings.

Owns:

- layout for the settings page
- rendering the settings form sections
- loading and saving configuration
- running connection tests

Should not own:

- active tab reload logic
- page bypass logic
- content script injection logic
- current-tab scripting logic

### `src/host/components/`

Host-only reusable UI for popup and settings page.

This includes:

- `ToggleSwitch`
- settings form sections
- quick action components

### `src/host/index.html`

Keep `index.html` in `src/host` as the entry surface for the fullscreen browser page.

## What should be removed from the fullscreen settings flow

The fullscreen settings page should not keep the current update code patterns such as:

- `saveAndRefresh(...)`
- `chrome.tabs.query(...)` for config changes
- `chrome.scripting.executeScript(...)`
- `chrome.tabs.reload(...)`
- domain-aware refresh decisions after settings edits

Those are the parts that currently make configuration changes behave like runtime control logic.

## What should stay in the popup

The popup should keep only the small host-control behavior:

- enable or disable toggle
- current connection status
- open settings page
- project link

If the popup still needs immediate host-side state changes for the toggle, keep that behavior isolated there rather than in the fullscreen settings page.

## Risk assessment

This refactor has relatively low risk if the fullscreen settings page becomes save-only.

### Low-risk part

- moving long-form settings to a fullscreen page
- removing reload and refresh behavior from the settings page
- keeping popup minimal

### Main user-facing tradeoff

- changing a setting in the fullscreen page will not immediately update the current browser tab
- settings apply on future page loads or future extension activity

This is acceptable and preferable because it makes the configuration UI much simpler and less fragile.

## Refactoring steps

### 1. Create browser host UI split

- create `src/host/popup/`
- create `src/host/pages/`
- create `src/host/components/`

### 2. Move the current options UI into a fullscreen page

- create `SettingsPage.tsx`
- move the long-form settings UI out of the current popup entry
- keep `index.html` in `src/host` as the fullscreen page entry

### 3. Shrink the popup

- keep enable or disable toggle
- keep connected status
- add open settings action
- keep project link

### 4. Remove update-oriented settings behavior

- remove `saveAndRefresh(...)` from the settings page flow
- remove active-tab reload logic from settings edits
- remove scripting-based update logic from settings edits

### 5. Keep runtime behavior separate

- popup handles quick control
- settings page handles configuration
- background and content keep browser runtime logic

## Result

After this refactor:

- the popup is small and focused
- the settings UI lives in a proper fullscreen browser page
- settings edits are normal save-only configuration changes
- browser runtime behavior is no longer mixed into the settings screen
- the codebase is in a much better state for later shared frontend extraction
