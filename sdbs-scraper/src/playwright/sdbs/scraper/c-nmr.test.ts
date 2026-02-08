import type { Browser, Page } from "playwright";
import { afterAll, beforeAll, describe, expect, it } from "vitest";
import { launchBrowser, launchPage } from "../../launch.js";
import { navigateTo, tryAcceptDisclaimer } from "../index.js";
import { scrapeCNmrSpecData } from "./c-nmr.js";

describe("scrapeCNmrSpecData", () => {
	let browser: Browser;
	let page: Page;

	beforeAll(async () => {
		browser = await launchBrowser({ headless: false });
		page = await launchPage(browser);
	});

	afterAll(async () => {
		await browser.close();
	});

	const testCases = [
		{ sdbsNumber: 3302, description: "methanol", expectedLength: 5 },
	];

	for (const { sdbsNumber, description, expectedLength } of testCases) {
		it(`should scrape C-NMR spec data for ${description}`, async () => {
			await navigateTo(page, sdbsNumber);
			await tryAcceptDisclaimer(page);

			const result = await scrapeCNmrSpecData(page);
			expect(result).toHaveLength(expectedLength);
			expect(result).toMatchSnapshot(`${description}-cnmr-data`);
		});
	}
});
