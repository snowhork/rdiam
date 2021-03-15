# rdiam(alpha version)
`rdiam` is CLI for redash group members, inspired by [bqiam](https://github.com/hirosassa/bqiam).
This is alpha version.

## usage
You can get user API KEY at `<your redash domain>/users/me`.

```
REDASH_ENDPOINT="your redash domain" REDASH_API_KEY="your user api key" \
go run ./ add -u user1@email.com,user2@emali.com -g group1,group2
```

In the future version, the endpoint and api_key may be managed by setting file.