run:

windows:
	GOOS=windows GOARCH=amd64 go install


linux:
	GOOS=linux GOARCH=amd64 go build -o ~/release/gte

