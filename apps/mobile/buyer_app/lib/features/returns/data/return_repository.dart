import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class ReturnItem {
  final String id;
  final String productName;
  final String imageUrl;
  final int quantity;

  const ReturnItem({
    required this.id,
    required this.productName,
    required this.imageUrl,
    required this.quantity,
  });

  factory ReturnItem.fromJson(Map<String, dynamic> json) {
    return ReturnItem(
      id: json['id'] as String,
      productName: json['productName'] as String,
      imageUrl: json['imageUrl'] as String,
      quantity: json['quantity'] as int,
    );
  }
}

class ReturnRequest {
  final String id;
  final String returnNumber;
  final String orderId;
  final String orderNumber;
  final String status;
  final String reason;
  final String? description;
  final DateTime createdAt;
  final List<ReturnItem> items;

  const ReturnRequest({
    required this.id,
    required this.returnNumber,
    required this.orderId,
    required this.orderNumber,
    required this.status,
    required this.reason,
    this.description,
    required this.createdAt,
    required this.items,
  });

  factory ReturnRequest.fromJson(Map<String, dynamic> json) {
    return ReturnRequest(
      id: json['id'] as String,
      returnNumber: json['returnNumber'] as String,
      orderId: json['orderId'] as String,
      orderNumber: json['orderNumber'] as String,
      status: json['status'] as String,
      reason: json['reason'] as String,
      description: json['description'] as String?,
      createdAt: DateTime.parse(json['createdAt'] as String),
      items: (json['items'] as List<dynamic>?)
              ?.map((e) => ReturnItem.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
    );
  }
}

class ReturnRepository {
  final ApiClient _apiClient;

  ReturnRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<List<ReturnRequest>> getReturns() async {
    final response = await _apiClient.get('/returns');
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => ReturnRequest.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<ReturnRequest> createReturn({
    required String orderId,
    required List<Map<String, dynamic>> items,
    required String reason,
    String? description,
  }) async {
    final response = await _apiClient.post('/returns', data: {
      'orderId': orderId,
      'items': items,
      'reason': reason,
      if (description != null) 'description': description,
    });
    return ReturnRequest.fromJson(response.data as Map<String, dynamic>);
  }

  Future<ReturnRequest> getReturnById(String id) async {
    final response = await _apiClient.get('/returns/$id');
    return ReturnRequest.fromJson(response.data as Map<String, dynamic>);
  }
}
