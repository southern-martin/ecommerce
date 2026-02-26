import 'package:get_it/get_it.dart';
import 'package:injectable/injectable.dart';

import '../logger/app_logger.dart';
import '../storage/secure_storage.dart';
import '../storage/local_storage.dart';

/// Global GetIt service locator instance.
final GetIt getIt = GetIt.instance;

/// Configures all dependencies using GetIt service locator.
///
/// Call this method early in the app startup (e.g., in main()) to register
/// all singleton and factory services before they are used.
///
/// [environment] can be used to differentiate between 'dev', 'staging',
/// and 'prod' configurations.
@InjectableInit()
Future<void> configureDependencies({String environment = 'prod'}) async {
  // Register core services

  // Logger — singleton so all parts of the app share the same logger instance
  getIt.registerLazySingleton<AppLogger>(() => AppLogger());

  // Secure storage — singleton for token and credential persistence
  getIt.registerLazySingleton<SecureStorage>(() => SecureStorage());

  // Local storage — async singleton because SharedPreferences.getInstance()
  // returns a Future
  getIt.registerSingletonAsync<LocalStorage>(() async {
    final localStorage = LocalStorage();
    await localStorage.init();
    return localStorage;
  });

  // Wait for all async singletons to complete initialization
  await getIt.allReady();

  getIt<AppLogger>().info(
    'Dependencies configured for environment: $environment',
  );
}

/// Resets all registered dependencies. Useful for testing.
Future<void> resetDependencies() async {
  await getIt.reset();
}
