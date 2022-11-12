### Description
Ths microservice is the simplest, but production-ready version of http://old.makala.com :) 

Enjoy it!

### Developers
In order to start the microservice, you need to run external services with the following command:
```shell
make services_up
```
Start the microservice itself with the command:
```shell
make api_up
```
API will be accessible at the port :8600 by default

### About Tests
The microservice is divided into three layers:
* handler or request parser (./server/rest/handler/*.go)
* service or business logic (./service/postfeed/*)
* datastore & cache (./provider/*)

<b>Note:</b> tests in the microservice are to show how I cover my code with tests. 
I did not intend to cover all the code with tests, 
because I think this is not necessary at the moment

### Unit Tests
<b>Run all</b>
```shell
make unit_test
```

<b>Test for handlers</b>
```shell
go test -v ./server/rest/handler/*.go
```
<b>Test for service</b>
```shell
go test -v ./service/postfeed/*.go
```
<b>Test for providers</b>

> Redis
> 
> ```shell
> go test -v ./provider/adstore/*.go
> ```

> PostgreSQL
> 
> There is no need to write for PostgreSQL as sqlc generates a type-safe, tested go code.

### Integration Tests
<b>Run all</b>
```shell
make integration_test
```

> <b>Test for Redis (cache)</b>
> ```shell
> go test -v tests/integration/provider/adstore/ads_test.go	
> go test -v tests/integration/provider/feedstore/feed_test.go
> ```

> <b>Test for PostgreSQL (cache)</b>
> ```shell
> go test -v tests/integration/provider/poststore/poststore_test.go
> ```


### Task 
Task is to build a microservice that:
1. accepts posts and stores them
2. generates a feed based on the score of each post
3. generates a feed, which contains promoted posts (ads)

### Solution
<b>Posts & Ads:</b>
* All posts are stored in PostgreSQL datastore, including promoted ones (ads)
* Feed and ads are cached in Redis (ids of posts are stored in cache)

<b>Feed:</b>
* There is a background worker, which periodically creates a feed in Redis and
  updates versions
* Feed is stored in Redis in a sorted set. The set is sorted based on the score of each post.
* The number of posts in a feed is limited. 
You can set a limit (at this moment it is 1000)
* Redis holds two versions of feed 'new' and 'old'. 
* 'New' version for users, who started scrolling feed (started with page = 0 or without started_fetching_at_unix_nano_utc in query)
* 'Old' version for users, who were scrolling, while feed gets updated

<b>Feed Update & Background worker:</b>
* Background worker goes through all posts in PostgreSQL & generates a feed with limited number\
* The newly created sorted set, will replace 'new' version of feed, while it replaces 'old' version
* Microservice remembers when user started to scroll feed. This is done based on the feed page (page = 0 is the start of feed scroll)
* Microservice also remembers feed update time
* Example of how feed version is picked:

Let's say a user started scrolling feed at 10:00:00 (hh:mm:ss)
If feed was updated at 09:00:00, then a user accesses the 'new' version of feed.
Then let's say feed is updated at 10:05:00 then a user will get the 'old' version of a feed.

As at 10:05:00
newly-generated-feed (by background worker) -- copied --> key_feed_version_new -- copied --> key_feed_version_old

<b>Note:</b> it is still not perfect (there should be several versions of feed for those, who scroll for a long time)

<b>User Feed & Ads</b>
* Ads (promoted posts) itself are stored in PostgreSQL, while their ids will be stored in Redis in a list
* Microservice remembers index of ad, which a user has seen. This means users will not miss any of the ads.
* If the index has reached the end, it starts from the beginning. This means that users will receive all ads in a circular fashion


### API Reference
You can find documentation of the microservice [here](./docs)

POST /api/makala/v1/post
> * there are restrictions that are described in the task (.pdf file)

GET /api/makala/v1/feed?count=10&page=0&author=t2_authorme&started_fetching_at_unix_nano_utc=1746560444234496670
> * count - number of posts (default 27)
> * page - page number, if it is 0 then started_fetching_at_unix_nano_utc will be set to current time & feed version will be 'new'
> * author - at this moment, there is no restriction on the name of an author. This must be taken from token/header/session,
>  but to make you life easier I put it here
> * started_fetching_at_unix_nano_utc - time, when a user started to scroll feed. This is used to switch between versions of feed


#### Requests

<b>Fetch feed 'new' version of feed:</b>
```shell
curl --location --request GET 'http://0.0.0.0:8600/api/makala/v1/feed?page=0&count=27&author=t2_author12' \
--header 'Content-Type: application/json; charset=utf-8'
```
```shell
curl --location --request GET 'http://0.0.0.0:8600/api/makala/v1/feed?page=1&count=27&author=t2_author12&started_fetching_at_unix_nano_utc=<current_time>' \
--header 'Content-Type: application/json; charset=utf-8'
```

<b>Fetch feed 'old' version of feed:</b>
```shell
curl --location --request GET 'http://0.0.0.0:8600/api/makala/v1/feed?page=1&count=27&author=t2_author12&started_fetching_at_unix_nano_utc=1146674939827247437' \
--header 'Content-Type: application/json; charset=utf-8'
```

Create post:
```shell
curl --location --request POST 'http://0.0.0.0:8600/api/makala/v1/post' \
--header 'Content-Type: application/json' \
--data-raw '{
    "author": "t2_author12",
    "content": "some content",
    "nsfw": false,
    "promoted": true,
    "score": 15.05,
    "submakala": "submakala",
    "title": "title 1"
}'
```
```shell
curl --location --request POST 'http://0.0.0.0:8600/api/makala/v1/post' \
--header 'Content-Type: application/json' \
--data-raw '{
    "author": "t2_author12",
    "link": "https://makala.com",
    "nsfw": false,
    "promoted": true,
    "score": 15.05,
    "submakala": "submakala",
    "title": "title 1"
}'
```

### Points of improvement 
* The microservice is not fully covered with unit tests, 
because this is just a task, not a microservice shifting to production :)
* The url validation might be improved. For example, we should check a protocol (https, http).
* The integration test should be commented properly and needs refactoring
* Ids of promoted posts are stored in Redis, while they are fetched from PostgreSQL. 
The storage in Redis could be eliminated

