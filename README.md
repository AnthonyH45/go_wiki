# go_wiki
Improved from this https://golang.org/doc/articles/wiki/

Steps to run:
- `git clone https://github.com/AnthonyH45/go_wiki.git`
- `cd go_wiki`
- `go build wiki.go`
- `./wiki`
- Open up a web browser, go to `localhost:8080/home/` and you should see some text.

The most difficult part of this was the regex validation for URL paths.

Adding styling, the home HTML template was not hard, except for the {{ }} stuff that I have never seen before, so with some googling, I was able to find a solution to list out the txt files and add links.


Pretty fun, the regex was the hardest part

Overall, I would say this was a fun experience, but did take me a few days to figure out the redirect and regex URL stuff.
9/10, I think if someone is bored, this is a nice project to spend some time on
