FROM alpine

COPY emojify-cache /

ENTRYPOINT ["/emojify-cache"]
