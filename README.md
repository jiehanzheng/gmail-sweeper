# Gmail Sweeper

I subscribe to a lot of newsletters, but I often can't go through all of them.  Most newsletters kind of "expire" after a few days, creating a backlog of newsletters that I will likely never read.

This program aims to fix that by looking for newsletters in a Gmail inbox, categorizing them, then retaining only N newsletters per category.  When scheduled to run as a cron job, this program can give you a clean inbox, which leads to a much better life for you!

## Setting this up

### Clone this repository

### Getting `credentials.json`

Complete [Step 1 in this guide](https://developers.google.com/gmail/api/quickstart/go) and download `credentials.json` into this directory.

### Edit `sweep_policy/policy.go`

* `Fetch` lets you define the scope of the messages you look at during each clean up.
* `Group` gives you an email message and asks you for a `string` as group ID.
* `Sweep` gives you a group ID and asks how you define "old" for this group, and what you want to do with the old messages.

### Run once

```bash
go build && ./gmail-sweeper
```

Note that you will need either Go 1.11 (or later) or [vgo](https://github.com/golang/go/wiki/vgo) for `go build` to work.

You will be prompted to set up your `token.json`.  Follow the on-screen instructions.

### Set up `cron` to perform automatic cleanups

Here is an example crontab:
```
37 * * * * /home/jiehan/gmail-sweeper/gmail-sweeper 2>&1 | tee /home/jiehan/gmail-sweeper/last-run.log
```

## New to Go; no shaming deal

This is my first project in Go, and I've decided to open source this to make it easier to share with a few friends.  At the same time, I hope some random strangers can benefit from using this tool as well.  But do not shame me for writing silly Go code!  I welcome any suggestions though.
