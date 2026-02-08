import type { ElementHandle, Page } from "playwright";

export const navigateTo = async (page: Page, sdbsNumber: number) => {
	await page.goto(
		`https://sdbs.db.aist.go.jp/CompoundView.aspx?sdbsno=${sdbsNumber}`,
	);
	await page.waitForLoadState("networkidle");
};

export const tryAcceptDisclaimer = async (page: Page) => {
	const disclaimer = await page.$(".DisclaimeraAcceptClass");
	if (disclaimer) {
		const input = await disclaimer.$("input");
		if (input) {
			await input.click();
		}
	}
	await page.waitForLoadState("networkidle");
};

export const tryAcceptDisclaimerWithMinDelay = async (
	page: Page,
	minDelayMs: number = 500,
) => {
	try {
		await page.waitForLoadState("networkidle");
		const start = Date.now();
		const disclaimer = await page.$(".DisclaimeraAcceptClass");
		if (disclaimer) {
			const input = await disclaimer.$("input");
			if (input) {
				await input.click();
			}
		}
		await page.waitForLoadState("networkidle");
		const elapsed = Date.now() - start;
		if (elapsed < minDelayMs) {
			await page.waitForTimeout(minDelayMs - elapsed);
		}
	} catch {
		// do nothing
	}
};

export const assertOnSdbsPage = (page: Page) => {
	const url = page.url();
	if (!url.startsWith("https://sdbs.db.aist.go.jp/CompoundView.aspx")) {
		throw new Error(`Not on a valid SDBS CompoundView page, got URL: ${url}`);
	}
};

export const getSideMenuSpectralLinks = async (
	page: Page,
	predicate?: (el: ElementHandle<HTMLAnchorElement>) => Promise<boolean>,
): Promise<ElementHandle<HTMLAnchorElement>[]> => {
	const links = (await page.$$(
		"#CtlMasterSideMenu_SpectralLink a",
	)) as ElementHandle<HTMLAnchorElement>[];
	if (!predicate) return links;

	const filtered = [];
	for (const link of links) {
		if (await predicate(link)) {
			filtered.push(link);
		}
	}
	return filtered;
};
