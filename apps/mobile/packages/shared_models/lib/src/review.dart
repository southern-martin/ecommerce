import 'package:json_annotation/json_annotation.dart';

part 'review.g.dart';

/// Represents a product review submitted by a user.
@JsonSerializable()
class Review {
  final String id;
  final String userId;
  final String userName;
  final String productId;
  final double rating;
  final String title;
  final String comment;
  final List<String> images;
  final DateTime createdAt;

  const Review({
    required this.id,
    required this.userId,
    required this.userName,
    required this.productId,
    required this.rating,
    required this.title,
    required this.comment,
    this.images = const [],
    required this.createdAt,
  });

  /// Whether the review includes images.
  bool get hasImages => images.isNotEmpty;

  factory Review.fromJson(Map<String, dynamic> json) =>
      _$ReviewFromJson(json);

  Map<String, dynamic> toJson() => _$ReviewToJson(this);

  Review copyWith({
    String? id,
    String? userId,
    String? userName,
    String? productId,
    double? rating,
    String? title,
    String? comment,
    List<String>? images,
    DateTime? createdAt,
  }) {
    return Review(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      userName: userName ?? this.userName,
      productId: productId ?? this.productId,
      rating: rating ?? this.rating,
      title: title ?? this.title,
      comment: comment ?? this.comment,
      images: images ?? this.images,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Review && runtimeType == other.runtimeType && id == other.id;

  @override
  int get hashCode => id.hashCode;

  @override
  String toString() =>
      'Review(id: $id, productId: $productId, rating: $rating)';
}
