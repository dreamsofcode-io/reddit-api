const playwright = require("playwright");
const logger = require("./logger");
const queue = require("./sqs");
const isLambda = require("is-lambda");
const chromium = require("@sparticuz/chromium");

const connectionUrl = process.env.CONNECTION_URL;

let allPosts = {};

const addPageInterceptors = async (page) => {
  await page.route("**/*", (route) => {
    const request = route.request();
    const resourceType = request.resourceType();
    if (
      resourceType === "image" ||
      resourceType === "font" ||
      resourceType === "stylesheet" ||
      resourceType === "script" ||
      resourceType === "media"
    ) {
      route.abort();
    } else {
      route.continue();
    }
  });
};

const getAttributes = async (handle) =>
  handle.evaluate((element) => {
    const attributeMap = {};
    for (const attr of element.attributes) {
      attributeMap[attr.name] = attr.value;
    }
    return attributeMap;
  });

const newBrowser = async () => {
  if (connectionUrl) {
    return playwright.chromium.connectOverCDP(connectionUrl);
  }

  if (isLambda) {
    chromium.setHeadlessMode = true;
    chromium.setGraphicsMode = false;

    return playwright.chromium.launch({
      args: chromium.args,
      executablePath: await chromium.executablePath(),
      headless: chromium.headless,
    });
  }

  return playwright.chromium.launch();
};

async function getDataForPostsConcurrently(posts) {
  logger.info("getting data for posts concurrently");

  posts.forEach((post) => {
    allPosts[post.id] = true;
  });

  return await Promise.all(
    posts.map(async (post) => {
      const browser = await newBrowser();
      const context = await browser.newContext();
      const page = await browser.newPage();
      addPageInterceptors(page);

      const data = await getPostData({ page, post });

      const nowStr = new Date().toISOString();
      await queue.publishOne({
        ...data,
        scrapedAt: nowStr,
      });

      await browser.close();
    }),
  );
}

async function getDataForPosts(posts, page) {
  logger.info("getting data for posts");
  let data = [];
  for (const post of posts) {
    let postData = await getPostData({ page, post });
    data.push(postData);
  }

  const nowStr = new Date().toISOString();
  await queue.publish(data.map((post) => ({ ...post, scrapedAt: nowStr })));
}

async function parseComment(e) {
  const things = await e.$$("> .sitetable > .thing");
  let comments = [];
  for (const thing of things) {
    const attributes = await getAttributes(thing);

    let thingClass = attributes["class"];
    let id = attributes["data-fullname"];
    let children = await parseComment(await thing.$(".child"));

    let isDeleted = thingClass.includes("deleted");
    let isCollapsed = thingClass.includes("collapsed");
    let author = isDeleted ? "" : attributes["data-author"];
    let time = await thing.$eval("time", (el) => el.getAttribute("datetime"));
    let comment =
      isDeleted || isCollapsed
        ? ""
        : await thing.$eval("div.md", (el) => el.innerText.trim());
    let pointsText =
      isDeleted || isCollapsed
        ? ""
        : await thing.$eval(
            "span.score",
            (el) => el.innerText.trim().split(" ")[0],
          );

    let points = parseInt(pointsText);
    points = isNaN(points) ? 0 : points;

    comments.push({
      id,
      author,
      time,
      comment,
      points,
      children,
      isDeleted,
      isCollapsed,
    });
  }

  return comments;
}

async function getPostData({ page, post }) {
  logger.info("getting details for post", { post: post.id });

  await page.goto(post.url);

  const sitetable = await page.$("div.sitetable");
  const thing = await sitetable.$(".thing");

  let id = post.id;
  let subreddit = post.subreddit;

  const attributes = await getAttributes(thing);
  let dataType = attributes["data-type"];
  let dataURL = attributes["data-url"];
  let isPromoted = attributes["data-promoted"] === "true";
  let isGallery = attributes["data-gallery"] === "true";

  let title = await page.$eval("a.title", (el) => el.innerText);
  let points = parseInt(await sitetable.$(".score.unvoted").innerText);
  let text = await sitetable.$("div.usertext-body").innerText;
  let comments = [];
  try {
    comments = await parseComment(await page.$("div.commentarea"));
  } catch (e) {
    logger.error("error parsing comments", { error: e });
  }

  logger.info("got details for post", { post: post.id });
  delete allPosts[post.id];
  logger.info("number of posts in progress", {
    count: Object.keys(allPosts).length,
  });
  logger.info("remaining posts", { posts: Object.keys(allPosts) });
  return {
    id,
    subreddit,
    dataType,
    dataURL,
    isPromoted,
    isGallery,
    title,
    timestamp: post.dt,
    timestamp_millis: post.timestamp,
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
    const attributes = await getAttributes(element);
    const id = attributes["data-fullname"];
    const subreddit = attributes["data-subreddit-prefixed"];
    const time = attributes["data-timestamp"];
    const timestamp = parseInt(time);
    const dt = new Date(timestamp);
    const author = attributes["data-author"];
    const url = `https://old.reddit.com${attributes["data-permalink"]}`;

    const post = { id, subreddit, dt, timestamp, author, url };
    posts.push(post);
  }

  return posts;
}

async function main() {
  logger.info("launching browser...");
  const browser = await newBrowser();

  logger.info("connecting...");
  const context = await browser.newContext();
  const page = await context.newPage();
  addPageInterceptors(page);

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

  if (connectionUrl) {
    await browser.close();
    await getDataForPostsConcurrently(posts);
  } else {
    await getDataForPosts(posts, page);
    await browser.close();
  }

  logger.info(`got ${posts.length} posts`);
}

if (require.main === module) {
  main();
}

exports.handler = async function (event, context) {
  try {
    await main();
  } catch (e) {
    // Catch all errors so that the function doesn't retry
    console.log(e);
    logger.error("error scraping", { error: e });
    return { success: false };
  }
  return { success: true };
};

const bytesForPage = async (page) => {
  const content = await page.content();
  return Buffer.byteLength(content, "utf8");
};
