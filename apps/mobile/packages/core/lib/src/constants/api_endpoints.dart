/// Centralized API endpoint constants.
///
/// All paths are relative to the base URL configured in the API client.
class ApiEndpoints {
  ApiEndpoints._();

  // ---------------------------------------------------------------------------
  // Auth
  // ---------------------------------------------------------------------------
  static const String login = '/api/v1/auth/login';
  static const String register = '/api/v1/auth/register';
  static const String refreshToken = '/api/v1/auth/refresh';
  static const String logout = '/api/v1/auth/logout';
  static const String forgotPassword = '/api/v1/auth/forgot-password';
  static const String resetPassword = '/api/v1/auth/reset-password';
  static const String verifyEmail = '/api/v1/auth/verify-email';
  static const String changePassword = '/api/v1/auth/change-password';

  // ---------------------------------------------------------------------------
  // Users
  // ---------------------------------------------------------------------------
  static const String userProfile = '/api/v1/users/me';
  static const String updateProfile = '/api/v1/users/me';
  static const String uploadAvatar = '/api/v1/users/me/avatar';
  static const String deleteAccount = '/api/v1/users/me';

  /// Returns the endpoint for a specific user by [id].
  static String userById(String id) => '/api/v1/users/$id';

  // ---------------------------------------------------------------------------
  // Products
  // ---------------------------------------------------------------------------
  static const String products = '/api/v1/products';
  static const String featuredProducts = '/api/v1/products/featured';
  static const String popularProducts = '/api/v1/products/popular';
  static const String newArrivals = '/api/v1/products/new-arrivals';

  /// Returns the endpoint for a specific product by [id].
  static String productById(String id) => '/api/v1/products/$id';

  /// Returns the endpoint for reviews of a specific product by [productId].
  static String productReviews(String productId) =>
      '/api/v1/products/$productId/reviews';

  // ---------------------------------------------------------------------------
  // Categories
  // ---------------------------------------------------------------------------
  static const String categories = '/api/v1/categories';

  /// Returns the endpoint for a specific category by [id].
  static String categoryById(String id) => '/api/v1/categories/$id';

  /// Returns the products under a specific category.
  static String categoryProducts(String categoryId) =>
      '/api/v1/categories/$categoryId/products';

  // ---------------------------------------------------------------------------
  // Cart
  // ---------------------------------------------------------------------------
  static const String cart = '/api/v1/cart';
  static const String addToCart = '/api/v1/cart/items';

  /// Returns the endpoint to update or remove a specific cart item.
  static String cartItem(String itemId) => '/api/v1/cart/items/$itemId';

  static const String clearCart = '/api/v1/cart/clear';

  // ---------------------------------------------------------------------------
  // Orders
  // ---------------------------------------------------------------------------
  static const String orders = '/api/v1/orders';

  /// Returns the endpoint for a specific order by [id].
  static String orderById(String id) => '/api/v1/orders/$id';

  /// Returns the endpoint to cancel an order.
  static String cancelOrder(String id) => '/api/v1/orders/$id/cancel';

  /// Returns the endpoint for tracking info of an order.
  static String orderTracking(String id) => '/api/v1/orders/$id/tracking';

  // ---------------------------------------------------------------------------
  // Checkout
  // ---------------------------------------------------------------------------
  static const String checkout = '/api/v1/checkout';
  static const String applyCoupon = '/api/v1/checkout/coupon';
  static const String shippingMethods = '/api/v1/checkout/shipping-methods';

  // ---------------------------------------------------------------------------
  // Search
  // ---------------------------------------------------------------------------
  static const String search = '/api/v1/search';
  static const String searchSuggestions = '/api/v1/search/suggestions';

  // ---------------------------------------------------------------------------
  // Reviews
  // ---------------------------------------------------------------------------
  static const String reviews = '/api/v1/reviews';

  /// Returns the endpoint for a specific review by [id].
  static String reviewById(String id) => '/api/v1/reviews/$id';

  // ---------------------------------------------------------------------------
  // Wishlist
  // ---------------------------------------------------------------------------
  static const String wishlist = '/api/v1/wishlist';

  /// Returns the endpoint to add/remove a specific product from wishlist.
  static String wishlistItem(String productId) =>
      '/api/v1/wishlist/$productId';

  // ---------------------------------------------------------------------------
  // Notifications
  // ---------------------------------------------------------------------------
  static const String notifications = '/api/v1/notifications';
  static const String markAllNotificationsRead =
      '/api/v1/notifications/mark-all-read';

  /// Returns the endpoint for a specific notification by [id].
  static String notificationById(String id) => '/api/v1/notifications/$id';

  // ---------------------------------------------------------------------------
  // Chat
  // ---------------------------------------------------------------------------
  static const String conversations = '/api/v1/chat/conversations';

  /// Returns the endpoint for a specific conversation.
  static String conversationById(String id) =>
      '/api/v1/chat/conversations/$id';

  /// Returns the messages endpoint for a specific conversation.
  static String conversationMessages(String conversationId) =>
      '/api/v1/chat/conversations/$conversationId/messages';

  // ---------------------------------------------------------------------------
  // Addresses
  // ---------------------------------------------------------------------------
  static const String addresses = '/api/v1/addresses';

  /// Returns the endpoint for a specific address by [id].
  static String addressById(String id) => '/api/v1/addresses/$id';

  // ---------------------------------------------------------------------------
  // Returns
  // ---------------------------------------------------------------------------
  static const String returns = '/api/v1/returns';

  /// Returns the endpoint for a specific return request by [id].
  static String returnById(String id) => '/api/v1/returns/$id';

  // ---------------------------------------------------------------------------
  // Seller
  // ---------------------------------------------------------------------------
  static const String sellerDashboard = '/api/v1/seller/dashboard';
  static const String sellerProducts = '/api/v1/seller/products';
  static const String sellerOrders = '/api/v1/seller/orders';
  static const String sellerAnalytics = '/api/v1/seller/analytics';
  static const String sellerPayouts = '/api/v1/seller/payouts';
  static const String sellerCoupons = '/api/v1/seller/coupons';

  // ---------------------------------------------------------------------------
  // Loyalty & Affiliate
  // ---------------------------------------------------------------------------
  static const String loyaltyPoints = '/api/v1/loyalty/points';
  static const String loyaltyRewards = '/api/v1/loyalty/rewards';
  static const String affiliateInfo = '/api/v1/affiliate';
  static const String affiliateEarnings = '/api/v1/affiliate/earnings';

  // ---------------------------------------------------------------------------
  // Shipping
  // ---------------------------------------------------------------------------
  static const String shippingRates = '/api/v1/shipping/rates';
  static const String trackShipment = '/api/v1/shipping/track';

  // ---------------------------------------------------------------------------
  // Promotions
  // ---------------------------------------------------------------------------
  static const String promotions = '/api/v1/promotions';
  static const String validateCoupon = '/api/v1/promotions/validate';
}
