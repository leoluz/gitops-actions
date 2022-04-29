# Gitops Actions (GoA)

# Motivation

As part of the Argo team, we have to maintain our Twitter account for the open
source project. This process consists of someone writing a text for a tweet and
asking the other member to review the text over Slack. Once the text is approved
it is then copied from slack by someone with access in the project's Twitter
account and manually creating a tweet. As the project grows and more member want
to collaborate, this process quickly becomes clunky and tedious.

This project aims to address this issue allowing configuring a Github repository
to trigger actions, like creating a tweet, once a pull request is approved and
merged.

# Goal

Gitops-Actions (or GoA) enables the execution of actions using a gitops
approach. Actions are implemented in the code base and can execute any task
triggered by a git event. A few examples of actions are: 
- Posting a tweet
- Create a new google doc
- Sending message over Slack
- Creating a ticket in Jira

Basically anything that exposes an API can be implemented as an action.

Check the existing actions and how to proper configure them in the
[Actions](#actions) section bellow.

# How

This project publishes a docker container that can be used in a CI tool to
execute pre-defined actions based on git events. It will inspect for modified
files in a pre-determined directory (`go-actions`) to decide which actions to
trigger. For example:

```
$ git status -s
A  go-actions/tweet/11-howto.txt
```

Once commited and merged the contents of file `11-howto.txt` will be
sent to the `tweet` action which will create a new tweet with that text.

An easy way to configure GoA is invoking it from a Github Action.
This is an example of a Github Action configuration that will invoke GoA with
the required environment variables:

```
$ cat .github/workflows/push-main.yml
name: gitops-actions
on:
  push:
    branches:
      - main
jobs:
  some-job:
    runs-on: ubuntu-latest
    container:
      image: leoluz/goa:latest
      env:
        GOA_BASE_SHA: ${{ github.event.before }}
        GOA_EVENT_SHA: ${{ github.event.after }}
        GOA_EVENT_REF_NAME: ${{ github.ref_name }}
        GOA_REPO_URL: ${{ github.server_url }}/${{ github.repository }}
    steps:
    - run: goa
```

# Actions

## Tweet

The Tweet action allows creating tweets based on git events like merging a file
in the `main` branch for example. For this action to work, it requires a twitter
account configured with an application that will be allowed to tweet on behalf
of that user. 

### Twitter Configuration

In order to register a Twitter application enable a developer access in your
Twitter account by accessing https://developer.twitter.com/ and signing up.

Once the developer account is active, follow the steps bellow:

- Create a new project. Example: `Gitops-Actions`
- Under the created project, create a new application. Example: `GitTweet`
- Take note of the Consumer key and Consumer Secret generated
- Inside the application configuration click the `edit` button under the "User
  authentication settings" section and do:
    - Enable OAuth 1.0a authentication method
    - Type of App: Automated App or Bot
    - Oauth 1.0a Settings: Read and Write
    - Callback URI: http://localhost
    - click Save
- Click the "Keys and tokens" tab at the top of the page
- Under the "Authentication Tokens" section click "Generate" button for "Access
  Token and Secret".
- Take note of the access token and access token secret generated.

### Github Repository Configuration

Create a new repository or choose an existing one to be used to trigger events to
create tweets. Access the repository settings page and follow the steps bellow:
- In the left menu click on the item Secrets > Actions
- In the top right corner click the button "New repository secret" button
- Create one repository secret for each of the 4 generated tokens created during
  the Twitter configuration. Example: TWITTER_CONSUMER_KEY,
  TWITTER_CONSUMER_SECRET, TWITTER_ACCESS_TOKEN and TWITTER_ACCESS_TOKEN_SECRET

Assuming you created all 4 secrets with the suggested names above, just create a
new github action with the following content:

```
$ cat .github/workflows/push-main.yml
name: gitops-actions
on:
  push:
    branches:
      - main
jobs:
  publish_twitter:
    runs-on: ubuntu-latest
    container:
      image: leoluz/goa:latest
      env:
        GOA_BASE_SHA: ${{ github.event.before }}
        GOA_EVENT_SHA: ${{ github.event.after }}
        GOA_EVENT_REF_NAME: ${{ github.ref_name }}
        GOA_REPO_URL: ${{ github.server_url }}/${{ github.repository }}
        GOA_TWITTER_CONSUMER_KEY: ${{ secrets.TWITTER_CONSUMER_KEY }}
        GOA_TWITTER_CONSUMER_SECRET: ${{ secrets.TWITTER_CONSUMER_SECRET }}
        GOA_TWITTER_ACCESS_TOKEN: ${{ secrets.TWITTER_ACCESS_TOKEN }}
        GOA_TWITTER_ACCESS_TOKEN_SECRET: ${{ secrets.TWITTER_ACCESS_TOKEN_SECRET }}
    steps:
    - run: goa
```

Once this is done, every new file created under `go-actions/tweet` folder will
have its content publised in the configured Twitter account once merged in the
`main` branch. Check [this repo](https://github.com/leoluz/kube-tests) for a
working example.
