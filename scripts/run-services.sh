#!/usr/bin/env zsh
# Run all Go services locally with correct env vars
# Usage: ./scripts/run-services.sh [service_name]
# Without args: starts all services
# With arg: starts only that service

export POSTGRES_PORT=15432 POSTGRES_USER=ecommerce POSTGRES_PASSWORD=ecommerce_secret
export REDIS_URL="localhost:16379" NATS_URL="nats://localhost:14222"
export OTEL_EXPORTER_ENDPOINT="http://localhost:4318"

GO=/opt/homebrew/bin/go
SERVICES=(auth user product cart order payment tax shipping return search review media notification cms promotion loyalty affiliate chat ai)
PIDS=()

# trap to kill all on exit
cleanup() { for pid in $PIDS; do kill $pid 2>/dev/null; done; exit; }
trap cleanup INT TERM

if [[ -n "$1" ]]; then
  # Run single service
  echo "Starting $1 service..."
  cd services/$1 && $GO run cmd/main.go
else
  # Run all services
  for svc in $SERVICES; do
    # Create DB if not exists
    docker exec ecommerce-postgres psql -U ecommerce -d postgres -c "CREATE DATABASE ecommerce_$svc;" 2>/dev/null

    echo "Starting $svc..."
    (cd services/$svc && $GO run cmd/main.go) &
    PIDS+=($!)
    sleep 0.5
  done

  echo ""
  echo "All services started! PIDs: $PIDS"
  echo "Press Ctrl+C to stop all services"
  wait
fi
