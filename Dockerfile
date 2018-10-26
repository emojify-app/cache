FROM scratch

COPY cache /

ENTRYPOINT ["/cache"]
