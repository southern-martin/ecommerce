import 'package:ecommerce_api_client/ecommerce_api_client.dart';

import '../../home/data/home_repository.dart';

class PaginatedResponse<T> {
  final List<T> data;
  final int total;
  final int page;
  final int pageSize;

  const PaginatedResponse({
    required this.data,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  bool get hasMore => page * pageSize < total;
}

class SearchRepository {
  final ApiClient _apiClient;

  SearchRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<PaginatedResponse<ProductSummary>> search({
    required String query,
    Map<String, dynamic>? filters,
    int page = 1,
    int pageSize = 20,
  }) async {
    final queryParams = <String, dynamic>{
      'q': query,
      'page': page,
      'pageSize': pageSize,
      ...?filters,
    };

    final response = await _apiClient.get('/search', queryParameters: queryParams);
    final data = response.data as Map<String, dynamic>;

    final products = (data['data'] as List<dynamic>)
        .map((e) => ProductSummary.fromJson(e as Map<String, dynamic>))
        .toList();

    return PaginatedResponse<ProductSummary>(
      data: products,
      total: data['total'] as int,
      page: data['page'] as int,
      pageSize: data['pageSize'] as int,
    );
  }

  Future<List<String>> getSuggestions(String query) async {
    final response = await _apiClient.get('/search/suggestions', queryParameters: {
      'q': query,
    });
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => e as String).toList();
  }
}
