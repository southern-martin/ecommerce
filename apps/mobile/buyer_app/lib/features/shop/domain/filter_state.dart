import 'package:equatable/equatable.dart';

enum SortOption {
  newest('newest', 'Newest'),
  priceAsc('price_asc', 'Price: Low to High'),
  priceDesc('price_desc', 'Price: High to Low'),
  rating('rating', 'Highest Rated'),
  popular('popular', 'Most Popular');

  final String value;
  final String label;
  const SortOption(this.value, this.label);
}

class FilterState extends Equatable {
  final String? category;
  final double? minPrice;
  final double? maxPrice;
  final double? rating;
  final bool? inStock;
  final int page;
  final int pageSize;
  final SortOption sort;
  final String? search;

  const FilterState({
    this.category,
    this.minPrice,
    this.maxPrice,
    this.rating,
    this.inStock,
    this.page = 1,
    this.pageSize = 20,
    this.sort = SortOption.newest,
    this.search,
  });

  bool get hasActiveFilters =>
      category != null ||
      minPrice != null ||
      maxPrice != null ||
      rating != null ||
      inStock != null;

  int get activeFilterCount {
    int count = 0;
    if (category != null) count++;
    if (minPrice != null || maxPrice != null) count++;
    if (rating != null) count++;
    if (inStock != null) count++;
    return count;
  }

  FilterState copyWith({
    String? category,
    double? minPrice,
    double? maxPrice,
    double? rating,
    bool? inStock,
    int? page,
    int? pageSize,
    SortOption? sort,
    String? search,
    bool clearCategory = false,
    bool clearMinPrice = false,
    bool clearMaxPrice = false,
    bool clearRating = false,
    bool clearInStock = false,
    bool clearSearch = false,
  }) {
    return FilterState(
      category: clearCategory ? null : (category ?? this.category),
      minPrice: clearMinPrice ? null : (minPrice ?? this.minPrice),
      maxPrice: clearMaxPrice ? null : (maxPrice ?? this.maxPrice),
      rating: clearRating ? null : (rating ?? this.rating),
      inStock: clearInStock ? null : (inStock ?? this.inStock),
      page: page ?? this.page,
      pageSize: pageSize ?? this.pageSize,
      sort: sort ?? this.sort,
      search: clearSearch ? null : (search ?? this.search),
    );
  }

  FilterState reset() {
    return FilterState(
      page: page,
      pageSize: pageSize,
      sort: sort,
      search: search,
    );
  }

  @override
  List<Object?> get props => [
        category,
        minPrice,
        maxPrice,
        rating,
        inStock,
        page,
        pageSize,
        sort,
        search,
      ];
}
