FROM golang:1.19.4 as development
ENV GO111MODULE=on

# Add a work directory
WORKDIR /app
COPY sl-exporter ./
COPY sl-exporter-config.yaml ./

# Expose ports
EXPOSE 4317
EXPOSE 4318
EXPOSE 9204

# Start app
ENTRYPOINT ["./sl-exporter"]
CMD ["--config", "sl-exporter-config.yaml"]
