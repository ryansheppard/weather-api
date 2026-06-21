FROM debian:13-slim AS builder

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update \
    && apt-get -y --no-install-recommends install \
    curl ca-certificates build-essential pkg-config libssl-dev

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ENV MISE_DATA_DIR="/mise"
ENV MISE_CONFIG_DIR="/mise"
ENV MISE_CACHE_DIR="/mise/cache"
ENV MISE_INSTALL_PATH="/usr/local/bin/mise"
ENV PATH="/mise/shims:$PATH"

RUN curl -fsSl https://mise.run | sh

WORKDIR /app
COPY mise.toml .
RUN --mount=type=cache,target=/mise/cache mise install

RUN --mount=type=cache,target=/usr/local/cargo/git/db,sharing=locked \
    --mount=type=cache,target=/usr/local/cargo/registry,sharing=locked \
    --mount=type=cache,target=/cargo_target,sharing=locked \
    CARGO_TARGET_DIR=/cargo_target \
    cargo install cargo-chef --locked

ENV RUSTC_WRAPPER=sccache \
    SCCACHE_DIR=/sccache

COPY Cargo.toml Cargo.lock /app/
RUN cargo chef prepare --recipe-path recipe.json

RUN --mount=type=cache,target=/usr/local/cargo/registry,sharing=locked \
    --mount=type=cache,target=/usr/local/cargo/git/db,sharing=locked \
    --mount=type=cache,target=$SCCACHE_DIR \
    cargo chef cook --release --recipe-path recipe.json

COPY . .

RUN --mount=type=cache,target=/usr/local/cargo/registry,sharing=locked \
    --mount=type=cache,target=/usr/local/cargo/git/db,sharing=locked \
    --mount=type=cache,target=$SCCACHE_DIR \
    cargo build --release --bin weathe_rs

FROM gcr.io/distroless/cc-debian13:nonroot

COPY --from=builder /app/target/release/weathe_rs /usr/local/bin/weathe_rs

ENTRYPOINT ["/usr/local/bin/weathe_rs"]
