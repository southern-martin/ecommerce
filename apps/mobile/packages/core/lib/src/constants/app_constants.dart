/// Application-wide constants.
class AppConstants {
  AppConstants._();

  /// Display name of the application.
  static const String appName = 'Ecommerce';

  /// Current application version string.
  static const String appVersion = '1.0.0';

  /// Build number.
  static const int buildNumber = 1;

  // ---------------------------------------------------------------------------
  // Networking
  // ---------------------------------------------------------------------------

  /// Default connection timeout in milliseconds.
  static const int connectionTimeout = 30000;

  /// Default receive timeout in milliseconds.
  static const int receiveTimeout = 30000;

  /// Default send timeout in milliseconds.
  static const int sendTimeout = 30000;

  /// Maximum number of automatic retry attempts for failed requests.
  static const int maxRetryAttempts = 3;

  /// Delay between retry attempts in milliseconds.
  static const int retryDelay = 1000;

  // ---------------------------------------------------------------------------
  // Pagination
  // ---------------------------------------------------------------------------

  /// Default number of items per page for paginated lists.
  static const int defaultPageSize = 20;

  /// Maximum number of items per page.
  static const int maxPageSize = 100;

  // ---------------------------------------------------------------------------
  // Cache
  // ---------------------------------------------------------------------------

  /// Default cache duration in minutes.
  static const int cacheDurationMinutes = 30;

  /// Maximum number of items to keep in memory cache.
  static const int maxCacheItems = 200;

  // ---------------------------------------------------------------------------
  // UI
  // ---------------------------------------------------------------------------

  /// Default animation duration in milliseconds.
  static const int animationDuration = 300;

  /// Short animation duration for subtle transitions.
  static const int animationDurationShort = 150;

  /// Long animation duration for page transitions.
  static const int animationDurationLong = 500;

  /// Debounce delay for search input in milliseconds.
  static const int searchDebounceMs = 400;

  /// Maximum number of recent searches to store.
  static const int maxRecentSearches = 10;

  // ---------------------------------------------------------------------------
  // Validation
  // ---------------------------------------------------------------------------

  /// Minimum password length.
  static const int minPasswordLength = 8;

  /// Maximum password length.
  static const int maxPasswordLength = 128;

  /// Maximum file upload size in bytes (10 MB).
  static const int maxFileUploadSize = 10 * 1024 * 1024;

  /// Allowed image extensions for upload.
  static const List<String> allowedImageExtensions = [
    'jpg',
    'jpeg',
    'png',
    'webp',
    'gif',
  ];

  // ---------------------------------------------------------------------------
  // Storage keys
  // ---------------------------------------------------------------------------

  /// Key for persisting the selected theme mode.
  static const String themeKey = 'theme_mode';

  /// Key for persisting the onboarding completion flag.
  static const String onboardingCompleteKey = 'onboarding_complete';

  /// Key for persisting recent search queries.
  static const String recentSearchesKey = 'recent_searches';

  /// Key for persisting the selected locale.
  static const String localeKey = 'locale';

  /// Key for persisting notification preferences.
  static const String notificationPrefsKey = 'notification_preferences';
}
