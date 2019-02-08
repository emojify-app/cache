FROM scratch

COPY emojify-cache /

ENTRYPOINT ["/emojify-cache"]
