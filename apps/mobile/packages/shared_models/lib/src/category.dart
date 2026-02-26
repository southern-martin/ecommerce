import 'package:json_annotation/json_annotation.dart';

part 'category.g.dart';

/// Represents a product category in the catalog hierarchy.
@JsonSerializable()
class Category {
  final String id;
  final String name;
  final String slug;
  final String? parentId;
  final String? imageUrl;
  final List<Category> children;

  const Category({
    required this.id,
    required this.name,
    required this.slug,
    this.parentId,
    this.imageUrl,
    this.children = const [],
  });

  /// Whether this is a root-level category.
  bool get isRoot => parentId == null;

  /// Whether this category has subcategories.
  bool get hasChildren => children.isNotEmpty;

  factory Category.fromJson(Map<String, dynamic> json) =>
      _$CategoryFromJson(json);

  Map<String, dynamic> toJson() => _$CategoryToJson(this);

  Category copyWith({
    String? id,
    String? name,
    String? slug,
    String? parentId,
    String? imageUrl,
    List<Category>? children,
  }) {
    return Category(
      id: id ?? this.id,
      name: name ?? this.name,
      slug: slug ?? this.slug,
      parentId: parentId ?? this.parentId,
      imageUrl: imageUrl ?? this.imageUrl,
      children: children ?? this.children,
    );
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Category && runtimeType == other.runtimeType && id == other.id;

  @override
  int get hashCode => id.hashCode;

  @override
  String toString() => 'Category(id: $id, name: $name, slug: $slug)';
}
