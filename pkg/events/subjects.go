package events

// NATS JetStream subject constants for all ecommerce domain events.
const (
	// Auth events
	SubjectUserRegistered = "auth.user.registered"
	SubjectUserLoggedIn   = "auth.user.logged_in"
	SubjectUserLoggedOut  = "auth.user.logged_out"
	SubjectPasswordReset  = "auth.password.reset"

	// User events
	SubjectUserCreated        = "user.created"
	SubjectUserUpdated        = "user.updated"
	SubjectUserDeleted        = "user.deleted"
	SubjectUserEmailVerified  = "user.email.verified"
	SubjectUserProfileUpdated = "user.profile.updated"

	// Product events
	SubjectProductCreated     = "product.created"
	SubjectProductUpdated     = "product.updated"
	SubjectProductDeleted     = "product.deleted"
	SubjectProductStockUpdate = "product.stock.updated"
	SubjectProductPriceUpdate = "product.price.updated"

	// Cart events
	SubjectCartItemAdded   = "cart.item.added"
	SubjectCartItemRemoved = "cart.item.removed"
	SubjectCartItemUpdated = "cart.item.updated"
	SubjectCartCleared     = "cart.cleared"

	// Order events
	SubjectOrderCreated    = "order.created"
	SubjectOrderConfirmed  = "order.confirmed"
	SubjectOrderCancelled  = "order.cancelled"
	SubjectOrderShipped    = "order.shipped"
	SubjectOrderDelivered  = "order.delivered"
	SubjectOrderRefunded   = "order.refunded"
	SubjectOrderCompleted  = "order.completed"

	// Payment events
	SubjectPaymentInitiated = "payment.initiated"
	SubjectPaymentCompleted = "payment.completed"
	SubjectPaymentFailed    = "payment.failed"
	SubjectPaymentRefunded  = "payment.refunded"

	// Notification events
	SubjectNotificationEmail = "notification.email"
	SubjectNotificationSMS   = "notification.sms"
	SubjectNotificationPush  = "notification.push"

	// Review events
	SubjectReviewCreated  = "review.created"
	SubjectReviewUpdated  = "review.updated"
	SubjectReviewDeleted  = "review.deleted"
	SubjectReviewApproved = "review.approved"

	// Search events
	SubjectSearchIndexProduct = "search.index.product"
	SubjectSearchRemoveProduct = "search.remove.product"

	// Shipping events
	SubjectShipmentCreated  = "shipping.shipment.created"
	SubjectShipmentUpdated  = "shipping.shipment.updated"
	SubjectShipmentDelivered = "shipping.shipment.delivered"

	// Promotion events
	SubjectPromotionCreated  = "promotion.created"
	SubjectPromotionExpired  = "promotion.expired"
	SubjectPromotionApplied  = "promotion.applied"

	// Return events
	SubjectReturnRequested = "return.requested"
	SubjectReturnApproved  = "return.approved"
	SubjectReturnRejected  = "return.rejected"
	SubjectReturnCompleted = "return.completed"

	// Dispute events
	SubjectDisputeOpened   = "dispute.opened"
	SubjectDisputeResolved = "dispute.resolved"

	// Loyalty events
	SubjectLoyaltyPointsEarned   = "loyalty.points.earned"
	SubjectLoyaltyPointsRedeemed = "loyalty.points.redeemed"
	SubjectLoyaltyTierUpgraded   = "loyalty.tier.upgraded"

	// Affiliate events
	SubjectAffiliateClickTracked      = "affiliate.click.tracked"
	SubjectAffiliateConversionTracked = "affiliate.conversion.tracked"
	SubjectAffiliatePayoutRequested   = "affiliate.payout.requested"

	// Coupon events
	SubjectCouponRedeemed   = "coupon.redeemed"
	SubjectFlashSaleStarted = "flash_sale.started"
	SubjectFlashSaleEnded   = "flash_sale.ended"

	// Chat events
	SubjectChatMessageSent     = "chat.message.sent"
	SubjectChatConversationNew = "chat.conversation.new"

	// AI events
	SubjectAIEmbeddingReady       = "ai.embedding.ready"
	SubjectAIRecommendationReady  = "ai.recommendation.ready"
	SubjectAIDescriptionGenerated = "ai.description.generated"
)

// Stream names for NATS JetStream.
const (
	StreamAuth         = "AUTH"
	StreamUser         = "USER"
	StreamProduct      = "PRODUCT"
	StreamCart         = "CART"
	StreamOrder        = "ORDER"
	StreamPayment      = "PAYMENT"
	StreamNotification = "NOTIFICATION"
	StreamReview       = "REVIEW"
	StreamSearch       = "SEARCH"
	StreamShipping     = "SHIPPING"
	StreamPromotion    = "PROMOTION"
	StreamReturn       = "RETURN"
	StreamDispute      = "DISPUTE"
	StreamLoyalty      = "LOYALTY"
	StreamAffiliate    = "AFFILIATE"
	StreamChat         = "CHAT"
	StreamAI           = "AI"
)
