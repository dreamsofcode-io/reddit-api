const playwright = require("playwright");
const connectionURL = process.env.CONNECTION_URL;

async function newBrowser() {
  if (connectionURL) {
    return await playwright.chromium.connectOverCDP({
      wsEndpoint: connectionURL,
    });
  } else {
    return await playwright.chromium.launch({
      headless: false,
    });
  });
}

async function handleLocalPostsData({ page, posts, parser }) {
    let data = [];

    for (const post of posts) {
      let postData = await parser({ page, post });
      data.push(postData);
    }

    return data;
}

async function handleRemotePostsData({ posts, parser }) {
  await Promise.all(posts.map(async (post) => {
    const browser = await newBrowser();
    const page = await newBrowser().newPage();

    let postData = await parser({ page, post });

    await page.close();
    return postData;
  }));

class BrowserManager {
  #browser = null;
  #page = null;

  constructor() {}

  async getPage() {
    if (this.#browser == null) {
      this.#browser = await newBrowser();
      this.#page = await this.#browser.newPage();
    }

    return this.#page;
  }

  async handlePostsData({ posts, parser }) {
    if (connectionURL) {
      return await handleRemotePostsData({ posts, parser });
    } else {
      const page = await this.newPage();
      return await handleLocalPostsData({ page, posts, parser });
    }
  }

  async close() {
    if (this.#browser != null) {
      await this.#browser.close();
    }
  }
}

exports.BrowserManager = BrowserManager;
