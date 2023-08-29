# Reddit API

This project provides a self serve reddit api in order to combat inflation.

## Requirements

### AWS

This code runs on AWS, so you'll need an AWS account and credentials in order
to deploy.

### Terraform

In order to deploy the code, terraform is used. Version 1.4.6 has been tested
for this code.

To deploy the code, you can run `terraform apply`, and to tear everything down
when completed, run `terraform destroy`

## Deployment

There are a number of services to deploy in this deployment model.

### Scraper

The scraper is the data gatherer. It uses playwright to pull data from
reddit. This data is all of the posts on a subreddit and their associated
comments.

The scraper has a schedule that runs once an hour. This is configurable in the
infrastructure as code.

#### Chrome

The chrome deployment will work with AWS Lambda. To install it, you just press enter
when prompted for the brightdata key which will configure terraform to use the
chrome deployment.

#### Brightdata

To use brightdata, you'll first need to sign up for an account at https://brightdata.com
Then, you'll need to create a new scraping browser solution.

Once you have the connection url, you can then enter it when running
`terraform apply`
