# rdiam(beta version)

[![Test](https://github.com/snowhork/rdiam/actions/workflows/test.yml/badge.svg)](https://github.com/snowhork/rdiam/actions/workflows/test.yml)
[![Apache-2.0](https://img.shields.io/github/license/snowhork/rdiam)](LICENSE)

`rdiam` is CLI for redash group members, inspired by [bqiam](https://github.com/hirosassa/bqiam).
This is beta version.

## Install
```
# requirement: go version >= 1.16.0
go install github.com/snowhork/rdiam
```

## Requirement
* go >= 1.16.0
* Redash version >= 7.0.0

`rdiam` uses Redash internal API. With another version, `rdiam` may not work.

## Usage
### Setting
In first, you can interactively set your redash endpoint and your redash user API key as below.

```
Enter you Redash endpoint (e.g. https://redash.yourdomain.com): 
Enter your Redash user API Key: 
```
You can get user API KEY at `<your redash domain>/users/me`.

Your settings are written at `~/.rdiam.yaml`.

### Example
```
rdiam add -u user1@email.com,user2@emali.com -g group1,group2
```

TODO: write explanation. (Sorry)

```
rdiam inspect query 1234
```
