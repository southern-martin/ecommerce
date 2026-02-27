import 'package:get_it/get_it.dart';
import 'package:ecommerce_core/ecommerce_core.dart';
import 'package:ecommerce_api_client/ecommerce_api_client.dart';

import '../../features/auth/data/auth_repository.dart';
import '../../features/dashboard/data/dashboard_repository.dart';
import '../../features/products/data/product_repository.dart';
import '../../features/orders/data/order_repository.dart';
import '../../features/returns/data/return_repository.dart';
import '../../features/shipments/data/shipment_repository.dart';
import '../../features/coupons/data/coupon_repository.dart';
import '../../features/analytics/data/analytics_repository.dart';
import '../../features/payouts/data/payout_repository.dart';

final GetIt getIt = GetIt.instance;

/// Configures all dependencies for the Seller App using GetIt service locator.
///
/// Call this in main() before runApp to register all singleton and factory
/// services.
Future<void> configureDependencies({String environment = 'prod'}) async {
  // Core services
  getIt.registerLazySingleton<SecureStorage>(() => SecureStorage());

  // API Client
  getIt.registerLazySingleton<ApiClient>(
    () => ApiClient(
      baseUrl: const String.fromEnvironment(
        'API_BASE_URL',
        defaultValue: 'https://api.example.com',
      ),
      secureStorage: getIt<SecureStorage>(),
    ),
  );

  // Auth
  getIt.registerLazySingleton<SellerAuthRepository>(
    () => SellerAuthRepository(
      apiClient: getIt<ApiClient>(),
      secureStorage: getIt<SecureStorage>(),
    ),
  );

  // Feature repositories
  getIt.registerLazySingleton<DashboardRepository>(
    () => DashboardRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<SellerProductRepository>(
    () => SellerProductRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<SellerOrderRepository>(
    () => SellerOrderRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<SellerReturnRepository>(
    () => SellerReturnRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<ShipmentRepository>(
    () => ShipmentRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<CouponRepository>(
    () => CouponRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<AnalyticsRepository>(
    () => AnalyticsRepository(apiClient: getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<PayoutRepository>(
    () => PayoutRepository(apiClient: getIt<ApiClient>()),
  );
}

/// Resets all registered dependencies. Useful for testing.
Future<void> resetDependencies() async {
  await getIt.reset();
}
