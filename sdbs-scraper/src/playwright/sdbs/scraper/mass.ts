import type { Page } from "playwright";
import {
	assertOnSdbsPage,
	getSideMenuSpectralLinks,
	tryAcceptDisclaimerWithMinDelay,
} from "../index.js";
import { getImageBytes } from "./utils/image.js";

export const scrapeMassSpecData = async (page: Page) => {
	await assertOnSdbsPage(page);

	await tryAcceptDisclaimerWithMinDelay(page);
	if (!(await clickMassSpecNavAnchor(page))) {
		return null;
	}

	await tryAcceptDisclaimerWithMinDelay(page);
	const massSpecImage = await getMassSpecImage(page);

	await tryAcceptDisclaimerWithMinDelay(page);
	const massSpecPeakData = await getMassSpecPeakData(page);

	return {
		image: massSpecImage,
		data: massSpecPeakData,
	};
};

const clickMassSpecNavAnchor = async (page: Page): Promise<boolean> => {
	const links = await getSideMenuSpectralLinks(
		page,
		async (el) => !!(await el.innerText()).match(/Mass\s{0,2}:/i),
	);
	if (links.length === 0) return false;
	if (links.length > 1)
		throw new Error(
			"More than one Mass link found in side menu; selection is ambiguous.",
		);
	await links[0].click();
	await page.waitForLoadState("networkidle");
	// Wait for the mass spec image to be visible
	await page.waitForSelector("#BodyContentPlaceHolder_SpectralInfoMode > img", {
		state: "visible",
		timeout: 10000,
	});
	return true;
};

const getMassSpecImage = (page: Page) => {
	return getImageBytes(page, "#BodyContentPlaceHolder_SpectralInfoMode > img");
};

type MassSpecPeakData = {
	"m/z": number;
	relativeIntensity: number;
}[];

const getMassSpecPeakData = async (page: Page) => {
	const button = await page.$(
		"#BodyContentPlaceHolder_SpectralInfoMode > input[type=submit]",
	);
	if (!button) {
		throw new Error("Mass spec peak data download button not found.");
	}
	await button.click();
	await page.waitForLoadState("networkidle");

	const thresholdInput = await page.$("#BodyContentPlaceHolder_ThresholdText");
	if (!thresholdInput) {
		throw new Error(
			"Threshold input (#BodyContentPlaceHolder_ThresholdText) not found.",
		);
	}
	await thresholdInput.fill("0");
	const downloadButton = await page.$(
		"#BodyContentPlaceHolder_PeakinfoMode > input[type=submit]:nth-child(3)",
	);
	if (!downloadButton) {
		throw new Error(
			"Download button for peak data not found (#BodyContentPlaceHolder_PeakinfoMode > input[type=submit]:nth-child(3)).",
		);
	}
	await downloadButton.click();
	await page.waitForLoadState("networkidle");

	const thresholdTable = await page.$("#BodyContentPlaceHolder_ThresholdTable");
	if (!thresholdTable) {
		throw new Error(
			"Threshold table (#BodyContentPlaceHolder_ThresholdTable) not found.",
		);
	}

	const tableText = (
		await thresholdTable.evaluate((node) => node.textContent)
	).trim();
	const splitText = tableText.split(/\s+/);

	const result: MassSpecPeakData = [];
	for (let i = 0; i < splitText.length; i += 2) {
		result.push({
			"m/z": parseFloat(splitText[i]),
			relativeIntensity: parseFloat(splitText[i + 1]),
		});
	}

	return result;
};
