# whats-that-card
Search terms to help you find that hearthstone card

Run
	go run main.go
to start the server.

To set up key config for mashape, run:
	heroku login
	heroku create # if you need to create
	heroku config:add MASHAPE_KEY="$MASHAPE_KEY"
which assumes you have a local env variable called `MASHAPE_KEY`, then
	git push heroku master
to push to heroku.

When active you can run `heroku logs` to check the output.

