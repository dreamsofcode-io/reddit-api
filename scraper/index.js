const playwright = require("playwright");
const logger = require("./logger");
const queue = require("./sqs");

async function parseComment(e) {
  const things = await e.$$("> .sitetable > .thing");
  let comments = [];
  for (const thing of things) {
    let thingClass = await things[0].getAttribute("class");
    let children = await parseComment(await thing.$(".child"));
    let isDeleted = thingClass.includes("deleted");
    let author = isDeleted
      ? ""
      : await thing.$eval(".author", (el) => el.innerText);
    let time = await thing.$eval("time", (el) => el.getAttribute("datetime"));
    let comment = isDeleted
      ? ""
      : await thing.$eval("div.md", (el) => el.innerText.trim());
    let points = isDeleted
      ? ""
      : await thing.$eval("span.score", (el) => el.innerText.trim());

    comments.push({ author, time, comment, points, children, isDeleted });
  }

  return comments;
}

async function getPostData({ page, post }) {
  logger.info("getting details for post", { post: post });

  await page.goto(post.url);

  const sitetable = await page.$("div.sitetable");
  const thing = await sitetable.$(".thing");

  let id = post.id;
  let subreddit = post.subreddit;
  let dataType = await thing.getAttribute("data-type");
  let dataURL = await thing.getAttribute("data-url");
  let isPromoted = (await thing.getAttribute("data-promoted")) === "true";
  let isGallery = (await thing.getAttribute("data-gallery")) === "true";
  let title = await page.$eval("a.title", (el) => el.innerText);
  let points = parseInt(await sitetable.$(".score.unvoted").innerText);
  let text = await sitetable.$("div.usertext-body").innerText;
  let comments = await parseComment(await page.$("div.commentarea"));

  return {
    id,
    subreddit,
    dataType,
    dataURL,
    isPromoted,
    isGallery,
    title,
    timestamp: post.timestamp,
    author: post.author,
    url: post.url,
    points: isNaN(points) ? 0 : points,
    text,
    comments,
  };
}

async function getPostsOnPage(page) {
  logger.info("getting posts for page");
  const elements = await page.$$(".thing");

  let posts = [];

  for (const element of elements) {
    const id = await element.getAttribute("data-fullname");
    const subreddit = await element.getAttribute("data-subreddit-prefixed");

    const time = await element.$("time");
    if (time == null) {
      continue;
    }

    const timestamp = Date.parse(await time.getAttribute("datetime"));
    const author = await element.$eval(".author", (el) => el.innerText);
    const url = await element.$eval("a.comments", (el) =>
      el.getAttribute("href")
    );

    posts.push({ id, subreddit, timestamp, author, url });
  }

  return posts;
}

async function main() {
  const browser = await playwright.chromium.launch({
    headless: false,
  });

  const context = await browser.newContext();
  const page = await context.newPage();

  await page.goto("https://old.reddit.com/r/programming/new/");
  logger.info("connected!");

  let hour = 1000 * 60 * 60;

  let now = Date.now();
  let cutoff = Date.now() - 24 * hour;
  let earliest = new Date();

  let posts = [];
  while (cutoff < earliest) {
    let pagePosts = await getPostsOnPage(page);
    if (pagePosts.length == 0) {
      break;
    }

    posts = posts.concat(pagePosts);
    let earliestPost = posts[posts.length - 1];
    earliest = earliestPost.timestamp;

    if (earliestPost.timestamp < cutoff) {
      break;
    }

    let nextPageURL = await page.$eval(".next-button a", (el) => el.href);
    await page.goto(nextPageURL);
  }

  posts = posts.filter((post) => post.timestamp > cutoff);

  let data = [];

  for (const post of posts) {
    let postData = await getPostData({ post, page });
    data.push(postData);
  }

  const nowStr = new Date().toISOString();

  await queue.publish(data.map((post) => ({ ...post, scrapedAt: nowStr })));

  logger.info(`got ${data.length} posts`);

  await browser.close();
}

if (require.main === module) {
  main();
}
