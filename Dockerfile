FROM alpine:edge AS dist
ARG TARGETPLATFORM

# dist is only there if goreleaser has created it before this build
COPY dist /dist/
RUN set -eux; \
  case "${TARGETPLATFORM:-}" in \
    linux/amd64)  f="dist/cli_linux_amd64_v1/restic-runner"; ;; \
    linux/arm)    f="dist/cli_linux_arm_7/restic-runner"; ;; \
    linux/arm64)  f="dist/cli_linux_arm64_v8.0/restic-runner"; ;; \
    *) echo "unknown TARGETPLATFORM: '${TARGETPLATFORM:-}'"; exit 1; ;; \
  esac; \
  cp "${f}" /; \
  chmod +x /restic-runner;

FROM alpine:edge
RUN apk add --no-cache \
    'ca-certificates>=20251003-r0' \
    'restic>=0.18.1' \
  ;
COPY --from=dist /restic-runner /

ENTRYPOINT [ "/restic-runner" ]
