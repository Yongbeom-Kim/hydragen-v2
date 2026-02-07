import { chromium } from "playwright";
import { launchBrowser, launchPage } from "./playwright/launch.js";
import { navigateTo, tryAcceptDisclaimer } from "./playwright/sdbs/index.js";
import { scrapeMassSpecData } from "./playwright/sdbs/scraper/mass-spec.js";

const browser = await launchBrowser();
const page = await launchPage(browser);
await navigateTo(page, 3);
await tryAcceptDisclaimer(page);
await scrapeMassSpecData(page);

await page.waitForTimeout(600000); // keep it open
