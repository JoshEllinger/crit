import { test, expect } from '@playwright/test';
import { clearAllComments, loadPage } from './helpers';

// ============================================================
// Share Feature — File Mode (fork: no default share_url)
// ============================================================
test.describe('Share — File Mode', () => {
  test.beforeEach(async ({ request }) => {
    await clearAllComments(request);
  });

  test('share button is hidden when share_url is not configured', async ({ page }) => {
    await loadPage(page);

    await expect(page.locator('#shareBtn')).toBeHidden();
    // No share toast or popover on initial load.
    await expect(page.locator('#toast-share')).toHaveCount(0);
    await expect(page.locator('.share-overlay')).toHaveCount(0);
  });
});
