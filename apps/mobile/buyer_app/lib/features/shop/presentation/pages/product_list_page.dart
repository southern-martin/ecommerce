import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';

import '../../../../core/di/injection.dart';
import '../../data/product_repository.dart';
import '../../domain/filter_state.dart';
import '../../../home/data/home_repository.dart';
import '../widgets/filter_bottom_sheet.dart';

class ProductListPage extends StatefulWidget {
  const ProductListPage({super.key});

  @override
  State<ProductListPage> createState() => _ProductListPageState();
}

class _ProductListPageState extends State<ProductListPage> {
  final ProductRepository _productRepo = getIt<ProductRepository>();
  final RefreshController _refreshController = RefreshController();
  final ScrollController _scrollController = ScrollController();

  List<ProductSummary> _products = [];
  List<Category> _categories = [];
  FilterState _filter = const FilterState();
  bool _isLoading = true;
  bool _isLoadingMore = false;
  int _totalProducts = 0;

  @override
  void initState() {
    super.initState();
    _loadCategories();
    _loadProducts();
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _refreshController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent - 200) {
      _loadMore();
    }
  }

  Future<void> _loadCategories() async {
    try {
      final categories = await _productRepo.getCategories();
      if (mounted) {
        setState(() => _categories = categories);
      }
    } catch (_) {}
  }

  Future<void> _loadProducts() async {
    setState(() => _isLoading = true);
    try {
      final result = await _productRepo.getProducts(
        category: _filter.category,
        minPrice: _filter.minPrice,
        maxPrice: _filter.maxPrice,
        rating: _filter.rating,
        inStock: _filter.inStock,
        page: 1,
        pageSize: _filter.pageSize,
        sort: _filter.sort.value,
        search: _filter.search,
      );

      if (mounted) {
        setState(() {
          _products = result.products;
          _totalProducts = result.total;
          _filter = _filter.copyWith(page: 1);
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  Future<void> _loadMore() async {
    if (_isLoadingMore || _products.length >= _totalProducts) return;

    setState(() => _isLoadingMore = true);
    try {
      final nextPage = _filter.page + 1;
      final result = await _productRepo.getProducts(
        category: _filter.category,
        minPrice: _filter.minPrice,
        maxPrice: _filter.maxPrice,
        rating: _filter.rating,
        inStock: _filter.inStock,
        page: nextPage,
        pageSize: _filter.pageSize,
        sort: _filter.sort.value,
        search: _filter.search,
      );

      if (mounted) {
        setState(() {
          _products.addAll(result.products);
          _filter = _filter.copyWith(page: nextPage);
          _isLoadingMore = false;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() => _isLoadingMore = false);
      }
    }
  }

  Future<void> _onRefresh() async {
    await _loadProducts();
    _refreshController.refreshCompleted();
  }

  void _showSortMenu() {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Padding(
            padding: EdgeInsets.all(16),
            child: Text(
              'Sort By',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
          ),
          ...SortOption.values.map((option) => ListTile(
                title: Text(option.label),
                trailing: _filter.sort == option
                    ? Icon(Icons.check, color: Theme.of(context).colorScheme.primary)
                    : null,
                onTap: () {
                  setState(() {
                    _filter = _filter.copyWith(sort: option, page: 1);
                  });
                  Navigator.pop(context);
                  _loadProducts();
                },
              )),
          const SizedBox(height: 16),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Products'),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () => context.push('/search'),
          ),
        ],
      ),
      body: Column(
        children: [
          // Filter & Sort bar
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            decoration: BoxDecoration(
              color: theme.scaffoldBackgroundColor,
              border: Border(
                bottom: BorderSide(color: Colors.grey.shade200),
              ),
            ),
            child: Row(
              children: [
                Expanded(
                  child: Text(
                    '$_totalProducts products',
                    style: theme.textTheme.bodyMedium?.copyWith(color: Colors.grey),
                  ),
                ),
                TextButton.icon(
                  onPressed: _showSortMenu,
                  icon: const Icon(Icons.sort, size: 18),
                  label: Text(_filter.sort.label),
                ),
                const SizedBox(width: 8),
                Badge(
                  isLabelVisible: _filter.hasActiveFilters,
                  label: Text('${_filter.activeFilterCount}'),
                  child: IconButton(
                    icon: const Icon(Icons.filter_list),
                    onPressed: () {
                      FilterBottomSheet.show(
                        context: context,
                        currentFilter: _filter,
                        categories: _categories,
                        onApply: (newFilter) {
                          setState(() => _filter = newFilter);
                          _loadProducts();
                        },
                      );
                    },
                  ),
                ),
              ],
            ),
          ),

          // Product grid
          Expanded(
            child: _isLoading
                ? const Center(child: CircularProgressIndicator())
                : _products.isEmpty
                    ? Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(Icons.inventory_2_outlined,
                                size: 64, color: Colors.grey.shade400),
                            const SizedBox(height: 16),
                            Text(
                              'No products found',
                              style: theme.textTheme.titleMedium?.copyWith(
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextButton(
                              onPressed: () {
                                setState(() => _filter = const FilterState());
                                _loadProducts();
                              },
                              child: const Text('Clear Filters'),
                            ),
                          ],
                        ),
                      )
                    : SmartRefresher(
                        controller: _refreshController,
                        onRefresh: _onRefresh,
                        child: GridView.builder(
                          controller: _scrollController,
                          padding: const EdgeInsets.all(16),
                          gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                            crossAxisCount: 2,
                            childAspectRatio: 0.65,
                            crossAxisSpacing: 12,
                            mainAxisSpacing: 12,
                          ),
                          itemCount: _products.length + (_isLoadingMore ? 2 : 0),
                          itemBuilder: (context, index) {
                            if (index >= _products.length) {
                              return const Center(child: CircularProgressIndicator());
                            }
                            final product = _products[index];
                            return _ProductGridCard(
                              product: product,
                              onTap: () => context.push('/products/${product.slug}'),
                            );
                          },
                        ),
                      ),
          ),
        ],
      ),
    );
  }
}

class _ProductGridCard extends StatelessWidget {
  final ProductSummary product;
  final VoidCallback? onTap;

  const _ProductGridCard({required this.product, this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return GestureDetector(
      onTap: onTap,
      child: Card(
        clipBehavior: Clip.antiAlias,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Expanded(
              flex: 3,
              child: Stack(
                fit: StackFit.expand,
                children: [
                  CachedNetworkImage(
                    imageUrl: product.imageUrl,
                    fit: BoxFit.cover,
                    placeholder: (_, __) => Container(
                      color: Colors.grey.shade100,
                    ),
                    errorWidget: (_, __, ___) => Container(
                      color: Colors.grey.shade200,
                      child: const Icon(Icons.image_not_supported),
                    ),
                  ),
                  if (product.hasDiscount)
                    Positioned(
                      top: 8,
                      left: 8,
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                        decoration: BoxDecoration(
                          color: Colors.red,
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          '-${product.discountPercentage.toInt()}%',
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 11,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ),
                  Positioned(
                    top: 8,
                    right: 8,
                    child: CircleAvatar(
                      radius: 16,
                      backgroundColor: Colors.white,
                      child: IconButton(
                        icon: const Icon(Icons.favorite_border, size: 16),
                        padding: EdgeInsets.zero,
                        onPressed: () {
                          // TODO: Add to wishlist
                        },
                      ),
                    ),
                  ),
                ],
              ),
            ),
            Expanded(
              flex: 2,
              child: Padding(
                padding: const EdgeInsets.all(8),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      product.name,
                      style: theme.textTheme.bodySmall?.copyWith(
                        fontWeight: FontWeight.w500,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const Spacer(),
                    Row(
                      children: [
                        const Icon(Icons.star, size: 14, color: Colors.amber),
                        const SizedBox(width: 2),
                        Text(
                          '${product.rating.toStringAsFixed(1)} (${product.reviewCount})',
                          style: theme.textTheme.labelSmall?.copyWith(
                            color: Colors.grey,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        Text(
                          '\$${product.price.toStringAsFixed(2)}',
                          style: theme.textTheme.titleSmall?.copyWith(
                            fontWeight: FontWeight.bold,
                            color: theme.colorScheme.primary,
                          ),
                        ),
                        if (product.hasDiscount) ...[
                          const SizedBox(width: 4),
                          Text(
                            '\$${product.compareAtPrice!.toStringAsFixed(2)}',
                            style: theme.textTheme.labelSmall?.copyWith(
                              decoration: TextDecoration.lineThrough,
                              color: Colors.grey,
                            ),
                          ),
                        ],
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
