import type { ElementHandle, Page } from "playwright";
import {
	assertOnSdbsPage,
	getSideMenuSpectralLinks,
	tryAcceptDisclaimerWithMinDelay,
} from "../index.js";
import { getImageBytes } from "./utils/image.js";

export const scrapeCNmrSpecData = async (page: Page) => {
	await assertOnSdbsPage(page);

	await tryAcceptDisclaimerWithMinDelay(page);
	const anchors = await getCNmrNavAnchors(page);
	const additionalInfos = anchors.map(
		async (anchor) => (await anchor.innerText()).trim().split(/\s*:\s*/)[1],
	);

	const result = [];

	for (let i = 0; i < anchors.length; ++i) {
		const anchors = await getCNmrNavAnchors(page);
		const anchor = anchors[i];
		const additionalInfo = await additionalInfos[i];
		await anchor.click();
		await tryAcceptDisclaimerWithMinDelay(page, 500);

		await tryAcceptDisclaimerWithMinDelay(page);

		const image = await getCNmrImage(page);
		await tryAcceptDisclaimerWithMinDelay(page);
		const data = await getCNmrPeakData(page);

		result.push({
			additionalInfo,
			image,
			data,
		});
	}

	return result;
};

const getCNmrNavAnchors = (
	page: Page,
): Promise<ElementHandle<HTMLAnchorElement>[]> => {
	return getSideMenuSpectralLinks(
		page,
		async (el) =>
			!!(await el.innerText()).match(/13\s{0,2}C\s{0,2}NMR\s{0,2}:/i),
	);
};

const getCNmrImage = (page: Page) => {
	return getImageBytes(page, "#pnlMain > div > img");
};

type CNmrPeakData = {
	ppm: number;
	intensity: number;
	assignedCarbon: number;
}[];

const getCNmrPeakData = async (page: Page) => {
	const selectors = [
		"#BodyContentPlaceHolder_OldShiftTable",
		"#BodyContentPlaceHolder_NewShiftTable",
	];

	// Promise.race over waitForSelector for both selectors:
	const waitPromises = selectors.map((selector) =>
		page
			.waitForSelector(selector, {
				state: "visible",
				timeout: 5000,
			})
			.then(() => selector),
	);

	let foundSelector: string | null = null;
	try {
		foundSelector = await Promise.race(waitPromises);
	} catch (_) {
		// All selectors timed out
	}

	if (!foundSelector) {
		throw new Error(
			"Threshold table (#BodyContentPlaceHolder_OldShiftTable or #BodyContentPlaceHolder_NewShiftTable) not found.",
		);
	}

	const specTable = await page.$(foundSelector);

	if (!specTable) {
		throw new Error(
			"Threshold table found by selector but could not retrieve element handle.",
		);
	}

	const tableText = (await specTable.innerText()).trim();
	const rows = tableText.split(/\n+/);

	const result: CNmrPeakData = [];
	for (const row of rows) {
		const cols = row.trim().split(/\s+/);
		if (cols.length < 2) continue;
		const ppm = parseFloat(cols[0]);
		const intensity = parseFloat(cols[1]);
		const assignedCarbon = parseFloat(cols[2]);

		if (Number.isNaN(ppm) && Number.isNaN(intensity)) {
			continue;
		}

		result.push({ ppm, intensity, assignedCarbon });
	}

	return result;
};
