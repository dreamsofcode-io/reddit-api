const {
  SQSClient,
  SendMessageBatchCommand,
  SendMessageCommand,
} = require("@aws-sdk/client-sqs");

const client = new SQSClient({ region: "us-east-1" });
const queueURL = process.env.QUEUE_URL;

async function publishChunk(chunk) {
  const command = new SendMessageBatchCommand({
    QueueUrl: queueURL,
    Entries: chunk,
  });

  await client.send(command);
}

exports.publishOne = async function (post) {
  const command = new SendMessageCommand({
    QueueUrl: queueURL,
    MessageBody: JSON.stringify(post),
    Id: post.id,
  });

  await client.send(command);
};

exports.publish = async function (posts) {
  const msgs = posts.map((post) => ({
    Id: post.id,
    MessageBody: JSON.stringify(post),
  }));

  const chunkSize = 10;
  for (let i = 0; i < msgs.length; i += chunkSize) {
    let chunk = msgs.slice(i, i + chunkSize);
    await publishChunk(chunk);
  }
};
