import type { Page } from "playwright";

export const getImageBytes = async (page: Page, selector: string) => {
	await page.waitForLoadState("networkidle");
	await page.waitForSelector(selector, {
		state: "visible",
		timeout: 10000,
	});
	const imgHandle = await page.$(selector);
	if (!imgHandle) {
		throw new Error(`Image not found for selector: ${selector}`);
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
