import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

/// Represents an item being returned.
class ReturnItem {
  final String id;
  final String productId;
  final String productName;
  final String? imageUrl;
  final int quantity;
  final double price;

  const ReturnItem({
    required this.id,
    required this.productId,
    required this.productName,
    this.imageUrl,
    required this.quantity,
    required this.price,
  });

  factory ReturnItem.fromJson(Map<String, dynamic> json) {
    return ReturnItem(
      id: json['id'] as String,
      productId: json['productId'] as String,
      productName: json['productName'] as String,
      imageUrl: json['imageUrl'] as String?,
      quantity: json['quantity'] as int,
      price: (json['price'] as num).toDouble(),
    );
  }
}

/// Represents a return request from a customer.
class SellerReturn {
  final String id;
  final String returnNumber;
  final String orderId;
  final String orderNumber;
  final String customerName;
  final String reason;
  final String status; // pending, approved, rejected
  final List<ReturnItem> items;
  final String? note;
  final DateTime createdAt;
  final DateTime updatedAt;

  const SellerReturn({
    required this.id,
    required this.returnNumber,
    required this.orderId,
    required this.orderNumber,
    required this.customerName,
    required this.reason,
    required this.status,
    required this.items,
    this.note,
    required this.createdAt,
    required this.updatedAt,
  });

  factory SellerReturn.fromJson(Map<String, dynamic> json) {
    return SellerReturn(
      id: json['id'] as String,
      returnNumber: json['returnNumber'] as String,
      orderId: json['orderId'] as String,
      orderNumber: json['orderNumber'] as String,
      customerName: json['customerName'] as String,
      reason: json['reason'] as String,
      status: json['status'] as String,
      items: (json['items'] as List<dynamic>?)
              ?.map((e) => ReturnItem.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      note: json['note'] as String?,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }
}

/// Paginated result wrapper for returns.
class PaginatedReturns {
  final List<SellerReturn> returns;
  final int totalCount;
  final int currentPage;
  final int totalPages;

  const PaginatedReturns({
    required this.returns,
    required this.totalCount,
    required this.currentPage,
    required this.totalPages,
  });

  factory PaginatedReturns.fromJson(Map<String, dynamic> json) {
    return PaginatedReturns(
      returns: (json['data'] as List<dynamic>)
          .map((e) => SellerReturn.fromJson(e as Map<String, dynamic>))
          .toList(),
      totalCount: json['totalCount'] as int,
      currentPage: json['currentPage'] as int,
      totalPages: json['totalPages'] as int,
    );
  }
}

/// Repository for managing seller returns.
class SellerReturnRepository {
  final ApiClient _apiClient;

  SellerReturnRepository({required ApiClient apiClient})
      : _apiClient = apiClient;

  /// Fetches paginated list of returns with optional status filter.
  Future<PaginatedReturns> getReturns({
    int page = 1,
    String? status,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'scope': 'seller',
    };
    if (status != null) queryParams['status'] = status;

    final response = await _apiClient.get(
      ApiEndpoints.returns,
      queryParameters: queryParams,
    );
    return PaginatedReturns.fromJson(response.data as Map<String, dynamic>);
  }

  /// Fetches a single return by ID.
  Future<SellerReturn> getReturnById(String id) async {
    final response = await _apiClient.get(
      '${ApiEndpoints.returns}/$id',
    );
    return SellerReturn.fromJson(response.data as Map<String, dynamic>);
  }

  /// Approves a return request with an optional note.
  Future<SellerReturn> approveReturn(String id, {String? note}) async {
    final body = <String, dynamic>{};
    if (note != null) body['note'] = note;

    final response = await _apiClient.post(
      '${ApiEndpoints.returns}/$id/approve',
      data: body,
    );
    return SellerReturn.fromJson(response.data as Map<String, dynamic>);
  }

  /// Rejects a return request with a required note.
  Future<SellerReturn> rejectReturn(String id, String note) async {
    final response = await _apiClient.post(
      '${ApiEndpoints.returns}/$id/reject',
      data: {'note': note},
    );
    return SellerReturn.fromJson(response.data as Map<String, dynamic>);
  }
}
