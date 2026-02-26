import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class Review {
  final String id;
  final String productId;
  final String userId;
  final String userName;
  final String? userAvatarUrl;
  final int rating;
  final String title;
  final String comment;
  final List<String> imageUrls;
  final DateTime createdAt;

  const Review({
    required this.id,
    required this.productId,
    required this.userId,
    required this.userName,
    this.userAvatarUrl,
    required this.rating,
    required this.title,
    required this.comment,
    this.imageUrls = const [],
    required this.createdAt,
  });

  factory Review.fromJson(Map<String, dynamic> json) {
    return Review(
      id: json['id'] as String,
      productId: json['productId'] as String,
      userId: json['userId'] as String,
      userName: json['userName'] as String,
      userAvatarUrl: json['userAvatarUrl'] as String?,
      rating: json['rating'] as int,
      title: json['title'] as String,
      comment: json['comment'] as String,
      imageUrls: (json['imageUrls'] as List<dynamic>?)
              ?.map((e) => e as String)
              .toList() ??
          [],
      createdAt: DateTime.parse(json['createdAt'] as String),
    );
  }
}

class PaginatedReviews {
  final List<Review> reviews;
  final int total;
  final double averageRating;

  const PaginatedReviews({
    required this.reviews,
    required this.total,
    required this.averageRating,
  });

  factory PaginatedReviews.fromJson(Map<String, dynamic> json) {
    return PaginatedReviews(
      reviews: (json['data'] as List<dynamic>)
          .map((e) => Review.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: json['total'] as int,
      averageRating: (json['averageRating'] as num? ?? 0).toDouble(),
    );
  }
}

class ReviewRepository {
  final ApiClient _apiClient;

  ReviewRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<PaginatedReviews> getProductReviews(
    String productId, {
    int page = 1,
  }) async {
    final response = await _apiClient.get(
      '/products/$productId/reviews',
      queryParameters: {'page': page},
    );
    return PaginatedReviews.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Review> createReview({
    required String productId,
    required int rating,
    required String title,
    required String comment,
    List<String> imageUrls = const [],
  }) async {
    final response = await _apiClient.post(
      '/products/$productId/reviews',
      data: {
        'rating': rating,
        'title': title,
        'comment': comment,
        if (imageUrls.isNotEmpty) 'imageUrls': imageUrls,
      },
    );
    return Review.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Review> updateReview({
    required String reviewId,
    int? rating,
    String? title,
    String? comment,
    List<String>? imageUrls,
  }) async {
    final response = await _apiClient.put(
      '/reviews/$reviewId',
      data: {
        if (rating != null) 'rating': rating,
        if (title != null) 'title': title,
        if (comment != null) 'comment': comment,
        if (imageUrls != null) 'imageUrls': imageUrls,
      },
    );
    return Review.fromJson(response.data as Map<String, dynamic>);
  }

  Future<void> deleteReview(String reviewId) async {
    await _apiClient.delete('/reviews/$reviewId');
  }
}
