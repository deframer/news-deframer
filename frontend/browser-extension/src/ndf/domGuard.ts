const OVERLAY_SELECTORS = ['.rtb-box', '#rteadblocker', '.adb-root', '[data-adblock-message]'];

/**
 * Removes every attribute from the provided element. We mainly call this on
 * <html> and <body> so any anti-adblock classes/styles applied by the page are
 * wiped before they take effect.
 */
export const resetElementAttributes = (el: Element | null) => {
  if (!el) return;
  while (el.attributes.length > 0) {
    el.removeAttribute(el.attributes[0].name);
  }
};

/**
 * Deletes the common overlay containers that sites inject when they want to
 * display "please disable your ad blocker" dialogs.
 */
export const purgeOverlayElements = () => {
  OVERLAY_SELECTORS.forEach((selector) => {
    document.querySelectorAll(selector).forEach((el) => el.remove());
  });
};

type ShouldKeepOptions = {
  rootEl?: HTMLElement | null;
  allowedIds?: Set<string>;
};

/**
 * Returns true if the node belongs to our “safe” DOM (global styles, the root
 * wrapper, etc.). Everything else is considered hostile and removed.
 */
export const shouldKeepNode = (node: Node, options: ShouldKeepOptions = {}) => {
  const { rootEl = null, allowedIds } = options;
  if (node === document.documentElement || node === document.head || node === document.body) {
    return true;
  }
  if (node === rootEl) {
    return true;
  }
  if (node.nodeName === 'TITLE') {
    return true;
  }
  if (node instanceof Element && allowedIds?.has(node.id)) {
    return true;
  }
  return false;
};

export type PreemptGuard = {
  release: () => void;
};

/**
 * Installs a synchronous DOM shield that hides <html>, strips attributes, and
 * deletes hostile nodes (scripts, iframes, overlays) before they can execute.
 * call release() once you have replaced the page with our own UI.
 */
export const installPreemptiveDomGuard = (opts?: { allowedIds?: string[] }) => {
  const html = document.documentElement;
  if (!html) return null;

  let body = document.body;
  if (!body) {
    body = document.createElement('body');
    html.appendChild(body);
  }

  resetElementAttributes(html);
  resetElementAttributes(body);

  const hideStyle = document.createElement('style');
  hideStyle.id = 'ndf-preempt-hide-style';
  hideStyle.textContent = 'html, body { display: none !important; }';
  (document.head || html).appendChild(hideStyle);

  const allowedIds = new Set(['ndf-preempt-hide-style', ...(opts?.allowedIds ?? [])]);

  // Suppress intrusive alert/confirm/prompt popups some sites use in anti-adblock flows.
  const originalAlert = window.alert;
  const originalConfirm = window.confirm;
  const originalPrompt = window.prompt;
  const originalPrint = window.print;
  window.alert = () => undefined;
  window.confirm = () => true;
  window.prompt = () => '';
  window.print = () => undefined;

  const observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
      if (
        mutation.type === 'attributes' &&
        (mutation.target === html || mutation.target === document.body)
      ) {
        resetElementAttributes(mutation.target as Element);
        return;
      }

      mutation.addedNodes.forEach((node) => {
        if (shouldKeepNode(node, { allowedIds })) {
          if (node instanceof Element) {
            // If a kept element (like body) is re-added, strip its attributes immediately
            resetElementAttributes(node);
          }
          // Recursively clean children of kept containers (body/head)
          // to ensure no hostile nodes sneak in via a replaced container.
          if (node === document.body || node === document.head) {
            Array.from(node.childNodes).forEach((child) => {
              if (!shouldKeepNode(child, { allowedIds })) {
                child.parentNode?.removeChild(child);
              }
            });
          }
          return;
        }

        // Remove anything else
        node.parentNode?.removeChild(node);
      });
    });

    purgeOverlayElements();
  });

  // Observe the entire document to catch root element replacements
  observer.observe(document, {
    childList: true,
    subtree: true,
    attributes: true,
  });

  return {
    release: () => {
      observer.disconnect();
      hideStyle.remove();
      html.style.removeProperty('display');
      document.body?.style.removeProperty('display');
      window.alert = originalAlert;
      window.confirm = originalConfirm;
      window.prompt = originalPrompt;
      window.print = originalPrint;
    },
  } satisfies PreemptGuard;
};
