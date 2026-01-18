# Build stage
FROM ruby:3.2-slim AS builder

WORKDIR /app

# 1. Instalează doar ce este strict necesar pentru alte gems (ex: pg are nevoie de libpq-dev)
# Am eliminat build-essential pentru a preveni compilarea lentă a gRPC
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

COPY Gemfile Gemfile.lock ./

# 2. Forțăm Bundler să NU compileze, ci să folosească platforma binară
RUN bundle config set --local force_ruby_platform false && \
    bundle config set --local without 'development test' && \
    bundle install

# Final stage
FROM ruby:3.2-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq-dev \
    wget \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local/bundle /usr/local/bundle
COPY . .

EXPOSE 8083 50053

HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8083/health || exit 1

CMD ["ruby", "server.rb"]