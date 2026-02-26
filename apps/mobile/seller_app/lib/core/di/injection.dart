import 'package:get_it/get_it.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

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
  // Core services from ecommerce_core
  getIt.registerLazySingleton<SecureStorage>(() => SecureStorage());

  // Auth
  getIt.registerLazySingleton<SellerAuthRepository>(
    () => SellerAuthRepository(secureStorage: getIt<SecureStorage>()),
  );

  // Feature repositories
  getIt.registerLazySingleton<DashboardRepository>(
    () => DashboardRepository(),
  );

  getIt.registerLazySingleton<SellerProductRepository>(
    () => SellerProductRepository(),
  );

  getIt.registerLazySingleton<SellerOrderRepository>(
    () => SellerOrderRepository(),
  );

  getIt.registerLazySingleton<SellerReturnRepository>(
    () => SellerReturnRepository(),
  );

  getIt.registerLazySingleton<ShipmentRepository>(
    () => ShipmentRepository(),
  );

  getIt.registerLazySingleton<CouponRepository>(
    () => CouponRepository(),
  );

  getIt.registerLazySingleton<AnalyticsRepository>(
    () => AnalyticsRepository(),
  );

  getIt.registerLazySingleton<PayoutRepository>(
    () => PayoutRepository(),
  );
}

/// Resets all registered dependencies. Useful for testing.
Future<void> resetDependencies() async {
  await getIt.reset();
}
