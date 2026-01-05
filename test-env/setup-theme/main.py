#!/usr/bin/env python3
import argparse
import sys
from playwright.sync_api import sync_playwright

def wait_with_log(page, delay):
    print(f"sleeping {delay}s...")
    page.wait_for_timeout(delay * 1000)

def step_login(page, base_url, delay, user, password):
    print("Logging in...")

    url = base_url + "/wp-admin/"
    print(f"Navigating to {url}")
    page.goto(url)

    page.fill("#user_login", user)
    page.fill("#user_pass", password)
    page.click("#wp-submit")

    page.wait_for_url(base_url + "/wp-admin/")
    wait_with_log(page, delay)

def step_theme_colormag_starter(page, base_url, theme_import_time, delay):
    print("Setting up theme...")

    url = base_url + "/wp-admin/themes.php?page=tg-starter-templates#/import/colormag/colormag"
    print(f"Navigating to {url}")

    page.goto(url)
    wait_with_log(page, delay)

    print(f"Clicking Continue")
    page.click("text=Continue")
    wait_with_log(page, delay)

    print(f"Clicking Start Import")
    page.click("text=Start Import")
    wait_with_log(page, delay)

    print("Clicking Start Import in dialog")
    page.locator('div[role="dialog"]').get_by_text("Start Import").click()
    wait_with_log(page, theme_import_time)

    wait_with_log(page, delay)

def is_installed(page, base_url, delay):
    print(f"Checking installation status at {base_url}...")
    page.goto(base_url)
    res = False
    if page.locator('img[src*="CM-Logo-Main.png"]').count() > 0:
        print("Theme detected (Logo found).")
        res = True
    wait_with_log(page, delay)
    return res

def main():
    parser = argparse.ArgumentParser(description="Playwright automation script")
    parser.add_argument("--url", default="http://localhost:8000", help="Target URL")
    parser.add_argument("--delay", type=int, default=5, help="Delay in seconds")
    parser.add_argument("--theme_import_time", type=int, default=600, help="Theme import time in seconds (default: 10 minutes)")
    parser.add_argument("--user", default="admin", help="Username")
    parser.add_argument("--password", default="password", help="Password")
    parser.add_argument("--headless", action="store_true", help="Run in headless mode")
    args = parser.parse_args()

    print(f"Running setup-theme...")
    print(f"Target URL: {args.url}")

    with sync_playwright() as p:
        print(f"Launching chromium...")
        browser = p.chromium.launch(headless=args.headless)
        page = browser.new_page()
        wait_with_log(page, args.delay)

        step_login(page, args.url, args.delay, args.user, args.password)
        if not is_installed(page, args.url, args.delay):
            print("Theme not installed installed, using the colormag starter.")
            step_theme_colormag_starter(page, args.url, args.theme_import_time, args.delay)
            if not is_installed(page, args.url, args.delay):
                print("Error: Theme installation failed verification.")
                browser.close()
                sys.exit(1)
        else:
            print("Theme already installed, skipping setup.")

        browser.close()

if __name__ == "__main__":
    main()
