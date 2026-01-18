# Build stage
FROM ruby:3.2-slim AS builder

WORKDIR /app

# Reintroducem build-essential pentru gem-urile mici (bigdecimal, pg, json)
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

COPY Gemfile Gemfile.lock ./

# IMPORTANT: force_ruby_platform trebuie să fie false (sau șters, căci e default false)
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