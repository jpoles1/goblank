build:
	env GOOS=windows GOARCH=386 go build -o dist/goblank.exe 
run:
	go build && ./goblank