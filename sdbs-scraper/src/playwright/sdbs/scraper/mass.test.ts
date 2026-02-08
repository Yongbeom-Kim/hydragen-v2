import type { Browser, Page } from "playwright";
import { afterAll, beforeAll, describe, expect, it } from "vitest";
import { launchBrowser, launchPage } from "../../launch.js";
import { navigateTo, tryAcceptDisclaimer } from "../index.js";
import { scrapeMassSpecData } from "./mass.js";

describe("scrapeMassSpecData", () => {
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
		{ sdbsNumber: 3022, description: "SDBS 3022" },
		{ sdbsNumber: 3151, description: "SDBS 3151" },
		{ sdbsNumber: 3308, description: "SDBS 3308" },
		{ sdbsNumber: 8120, description: "SDBS 8120" },
	];

	for (const { sdbsNumber, description } of testCases) {
		it(`should scrape mass spec data for ${description}`, async () => {
			await navigateTo(page, sdbsNumber);
			await tryAcceptDisclaimer(page);

			const result = await scrapeMassSpecData(page);
			expect(result).toMatchSnapshot(`${description}-peak-data`);
		});
	}
});
