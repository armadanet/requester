FROM codyperakslis/dev-go:1.14.4-v1
COPY requester/go.mod requester/go.sum requester/
ADD spinner/ /home/dev/spinner/
RUN cd requester && go mod download 
ADD requester/ /home/dev/requester 
WORKDIR /home/dev/requester
RUN go build -o main cmd/requester.go
ENTRYPOINT [ "./main" ]