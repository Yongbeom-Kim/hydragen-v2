import { type Browser, chromium } from "playwright";

export type LaunchBrowserOptions = {
	headless?: boolean;
};

export const launchBrowser = (options: LaunchBrowserOptions = {}) => {
	const { headless = false } = options;
	return chromium.launch({
		headless,
		args: headless ? [] : ["--ozone-platform=x11"],
	});
};

export const launchPage = (browser: Browser) => {
	return browser.newPage();
};
