# Makefile

.PHONY: load_get
load_get:
	@echo "[LOAD][GET] fetching posts by user..."
	ab -n 1000000 -c 100 -H 'Content-Type: application/json; charset=utf-8' 'http://0.0.0.0:8600/api/makala/v1/feed?page=1&count=27&author=t2_author12&started_fetching_at_unix_nano_utc=1668279680'

.PHONY: load_post
load_post:
	@echo "[LOAD][GET] fetching posts by user..."
	ab -n 1000000 -c 100 -p new_post.json -H 'Content-Type: application/json; charset=utf-8' -T 'application/json; charset=utf-8' 'http://0.0.0.0:8600/api/makala/v1/post'
