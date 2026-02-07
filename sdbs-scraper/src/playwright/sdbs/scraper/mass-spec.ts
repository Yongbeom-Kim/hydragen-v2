import type { Page } from "playwright";
import {
	assertOnSdbsPage,
	getSideMenuSpectralLinks,
	tryAcceptDisclaimerWithMinDelay,
} from "../index.js";

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

const getMassSpecImage = async (page: Page) => {
	const imgHandle = await page.$(
		"#BodyContentPlaceHolder_SpectralInfoMode > img",
	);
	if (!imgHandle) {
		throw new Error("Mass spec image not found on the page.");
	}

	const imgBytes = await imgHandle.evaluate(async (img: HTMLImageElement) => {
		return new Promise<Uint8Array>((resolve, reject) => {
			const canvas = document.createElement("canvas");
			canvas.width = img.naturalWidth;
			canvas.height = img.naturalHeight;
			const ctx = canvas.getContext("2d");
			if (!ctx) {
				reject(new Error("Failed to get canvas 2D context."));
				return;
			}
			ctx.drawImage(img, 0, 0);
			canvas.toBlob((blob) => {
				if (!blob) {
					reject(new Error("Failed to convert canvas to blob."));
					return;
				}
				const reader = new FileReader();
				reader.onloadend = () => {
					const arrayBuffer = reader.result as ArrayBuffer;
					resolve(new Uint8Array(arrayBuffer));
				};
				reader.onerror = () => reject(reader.error);
				reader.readAsArrayBuffer(blob);
			}, "image/png");
		});
	});
	return imgBytes;
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
