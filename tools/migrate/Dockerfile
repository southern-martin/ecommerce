FROM postgres:16-alpine

# Install bash for the migration script
RUN apk add --no-cache bash findutils

# Copy migration SQL files from all services
# Each service's migrations are expected at services/<name>/migrations/*.up.sql
# They are copied into /migrations/<service>/ in the image

COPY ../../services/auth/migrations/ /migrations/auth/
COPY ../../services/user/migrations/ /migrations/user/
COPY ../../services/product/migrations/ /migrations/product/
COPY ../../services/cart/migrations/ /migrations/cart/
COPY ../../services/order/migrations/ /migrations/order/
COPY ../../services/payment/migrations/ /migrations/payment/
COPY ../../services/search/migrations/ /migrations/search/
COPY ../../services/review/migrations/ /migrations/review/
COPY ../../services/notification/migrations/ /migrations/notification/
COPY ../../services/chat/migrations/ /migrations/chat/
COPY ../../services/media/migrations/ /migrations/media/
COPY ../../services/cms/migrations/ /migrations/cms/
COPY ../../services/promotion/migrations/ /migrations/promotion/
COPY ../../services/loyalty/migrations/ /migrations/loyalty/
COPY ../../services/affiliate/migrations/ /migrations/affiliate/
COPY ../../services/shipping/migrations/ /migrations/shipping/
COPY ../../services/return/migrations/ /migrations/return/
COPY ../../services/tax/migrations/ /migrations/tax/
COPY ../../services/ai/migrations/ /migrations/ai/

# Copy the migration runner script
COPY run-migrations.sh /run-migrations.sh
RUN chmod +x /run-migrations.sh

ENTRYPOINT ["/run-migrations.sh"]
