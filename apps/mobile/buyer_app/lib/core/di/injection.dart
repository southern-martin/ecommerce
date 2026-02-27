import 'package:get_it/get_it.dart';
import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

import '../../features/auth/data/auth_repository.dart';
import '../../features/home/data/home_repository.dart';
import '../../features/shop/data/product_repository.dart';
import '../../features/search/data/search_repository.dart';
import '../../features/cart/data/cart_repository.dart';
import '../../features/checkout/data/checkout_repository.dart';
import '../../features/orders/data/order_repository.dart';
import '../../features/returns/data/return_repository.dart';
import '../../features/tracking/data/tracking_repository.dart';
import '../../features/profile/data/profile_repository.dart';
import '../../features/reviews/data/review_repository.dart';
import '../../features/loyalty/data/loyalty_repository.dart';
import '../../features/affiliate/data/affiliate_repository.dart';
import '../../features/notifications/data/notification_repository.dart';
import '../../features/chat/data/chat_repository.dart';
import '../../features/ai/data/ai_repository.dart';
import '../../features/wishlist/data/wishlist_repository.dart';

final getIt = GetIt.instance;

Future<void> configureDependencies() async {
  // Core services
  getIt.registerLazySingleton<SecureStorage>(() => SecureStorage());

  // API Client
  getIt.registerLazySingleton<ApiClient>(
    () => ApiClient(baseUrl: const String.fromEnvironment('API_BASE_URL', defaultValue: 'https://api.example.com')),
  );

  // Repositories
  getIt.registerLazySingleton<AuthRepository>(
    () => AuthRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<HomeRepository>(
    () => HomeRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<ProductRepository>(
    () => ProductRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<SearchRepository>(
    () => SearchRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<CartRepository>(
    () => CartRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<CheckoutRepository>(
    () => CheckoutRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<OrderRepository>(
    () => OrderRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<ReturnRepository>(
    () => ReturnRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<TrackingRepository>(
    () => TrackingRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<ProfileRepository>(
    () => ProfileRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<ReviewRepository>(
    () => ReviewRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<LoyaltyRepository>(
    () => LoyaltyRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<AffiliateRepository>(
    () => AffiliateRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<NotificationRepository>(
    () => NotificationRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<ChatRepository>(
    () => ChatRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<AIRepository>(
    () => AIRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<WishlistRepository>(
    () => WishlistRepository(apiClient: getIt<ApiClient>()),
  );
}
