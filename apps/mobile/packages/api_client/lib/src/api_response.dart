/// Generic wrapper for API responses.
///
/// Provides consistent deserialization for both single-item and
/// paginated responses from the server.
///
/// Example:
/// ```dart
/// final response = await apiClient.get('/products/123');
/// final apiResponse = ApiResponse<Product>.fromJson(
///   response.data,
///   (json) => Product.fromJson(json as Map<String, dynamic>),
/// );
/// print(apiResponse.data?.name);
/// ```
class ApiResponse<T> {
  /// The deserialized data payload. May be `null` for empty responses.
  final T? data;

  /// Human-readable message from the server (e.g., "Product created").
  final String? message;

  /// Whether the request was successful according to the server.
  final bool success;

  /// Optional metadata from the server (e.g., timing, version info).
  final Map<String, dynamic>? meta;

  const ApiResponse({
    this.data,
    this.message,
    this.success = true,
    this.meta,
  });

  /// Deserializes a JSON map into an [ApiResponse].
  ///
  /// [fromJsonT] is a factory function that converts the `data` field of the
  /// JSON into an instance of [T].
  factory ApiResponse.fromJson(
    Map<String, dynamic> json,
    T Function(Object? json) fromJsonT,
  ) {
    return ApiResponse<T>(
      data: json.containsKey('data') ? fromJsonT(json['data']) : null,
      message: json['message'] as String?,
      success: json['success'] as bool? ?? true,
      meta: json['meta'] as Map<String, dynamic>?,
    );
  }

  /// Converts this response to a JSON map.
  ///
  /// [toJsonT] converts the [data] field back to JSON. If not provided,
  /// the data field is omitted.
  Map<String, dynamic> toJson([Object? Function(T value)? toJsonT]) {
    return {
      if (data != null && toJsonT != null) 'data': toJsonT(data as T),
      if (message != null) 'message': message,
      'success': success,
      if (meta != null) 'meta': meta,
    };
  }

  @override
  String toString() =>
      'ApiResponse(success: $success, message: $message, data: $data)';
}

/// Wrapper for paginated API responses.
///
/// Example:
/// ```dart
/// final response = await apiClient.get('/products', queryParameters: {'page': 1});
/// final paginated = PaginatedResponse<Product>.fromJson(
///   response.data,
///   (json) => Product.fromJson(json as Map<String, dynamic>),
/// );
/// print('${paginated.items.length} of ${paginated.total} products');
/// ```
class PaginatedResponse<T> {
  /// The list of items on the current page.
  final List<T> items;

  /// The total number of items across all pages.
  final int total;

  /// The current page number (1-based).
  final int page;

  /// The number of items per page.
  final int pageSize;

  const PaginatedResponse({
    required this.items,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  /// The total number of pages.
  int get totalPages => (total / pageSize).ceil();

  /// Whether there is a next page.
  bool get hasNextPage => page < totalPages;

  /// Whether there is a previous page.
  bool get hasPreviousPage => page > 1;

  /// Whether the current page is the first page.
  bool get isFirstPage => page == 1;

  /// Whether the current page is the last page.
  bool get isLastPage => page >= totalPages;

  /// Deserializes a JSON map into a [PaginatedResponse].
  ///
  /// Expects the JSON to have:
  /// - `items` or `data`: a list of serialized items
  /// - `total` or `totalCount`: total item count
  /// - `page` or `currentPage`: current page number
  /// - `pageSize` or `perPage` or `limit`: items per page
  factory PaginatedResponse.fromJson(
    Map<String, dynamic> json,
    T Function(Object? json) fromJsonT,
  ) {
    final rawItems = (json['items'] ?? json['data']) as List<dynamic>? ?? [];
    final items = rawItems.map((item) => fromJsonT(item)).toList();

    return PaginatedResponse<T>(
      items: items,
      total: (json['total'] ?? json['totalCount'] ?? 0) as int,
      page: (json['page'] ?? json['currentPage'] ?? 1) as int,
      pageSize: (json['pageSize'] ?? json['perPage'] ?? json['limit'] ?? 20) as int,
    );
  }

  /// Converts this paginated response to a JSON map.
  Map<String, dynamic> toJson([Object? Function(T value)? toJsonT]) {
    return {
      'items': toJsonT != null
          ? items.map((item) => toJsonT(item)).toList()
          : items,
      'total': total,
      'page': page,
      'pageSize': pageSize,
    };
  }

  @override
  String toString() =>
      'PaginatedResponse(page: $page/$totalPages, '
      'items: ${items.length}, total: $total)';
}
