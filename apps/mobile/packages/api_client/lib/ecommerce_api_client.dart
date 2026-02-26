/// HTTP API client package for the ecommerce mobile application.
///
/// Provides a pre-configured Dio client with JWT authentication,
/// automatic token refresh, error mapping, and retry logic.
library ecommerce_api_client;

export 'src/api_client.dart';
export 'src/api_response.dart';
export 'src/interceptors/auth_interceptor.dart';
export 'src/interceptors/error_interceptor.dart';
