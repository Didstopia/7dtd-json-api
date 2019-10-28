# Compiler image
FROM didstopia/base:go-alpine-3.5 AS go-builder

# Copy the project 
COPY . /tmp/7dtd-json-api/
WORKDIR /tmp/7dtd-json-api/

# Install dependencies
RUN make deps

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/7dtd-json-api



# Runtime image
FROM scratch

# Copy certificates
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the built binary
COPY --from=go-builder /go/bin/7dtd-json-api /go/bin/7dtd-json-api

# Expose environment variables
ENV SERVER  "localhost"
ENV PORT    "8081"
ENV PASSWORD ""

# Run the binary
ENTRYPOINT ["/go/bin/7dtd-json-api"]
