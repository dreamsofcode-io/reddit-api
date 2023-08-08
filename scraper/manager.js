const playwright = require("playwright");

async function newBrowser() {
  return await playwright.chromium.launch({
    headless: false,
  });
}

class BrowserManager {
  #browser = null;

  constructor() {}

  async newPage() {
    if (this.#browser == null) {
      this.#browser = await newBrowser();
    }

    return await this.#browser.newPage();
  }

  async handlePostsData({ posts, parser }) {
    const page = await this.newPage();

    let data = [];

    for (const post of posts) {
      let postData = await parser({ page, post });
      data.push(postData);
    }

    return data;
  }

  async close() {
    if (this.#browser != null) {
      await this.#browser.close();
    }
  }
}

exports.BrowserManager = BrowserManager;
