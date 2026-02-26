import 'package:flutter/material.dart';

import '../../domain/filter_state.dart';
import '../../../home/data/home_repository.dart';

class FilterBottomSheet extends StatefulWidget {
  final FilterState currentFilter;
  final List<Category> categories;
  final ValueChanged<FilterState> onApply;

  const FilterBottomSheet({
    super.key,
    required this.currentFilter,
    required this.categories,
    required this.onApply,
  });

  static Future<FilterState?> show({
    required BuildContext context,
    required FilterState currentFilter,
    required List<Category> categories,
    required ValueChanged<FilterState> onApply,
  }) {
    return showModalBottomSheet<FilterState>(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (_) => FilterBottomSheet(
        currentFilter: currentFilter,
        categories: categories,
        onApply: onApply,
      ),
    );
  }

  @override
  State<FilterBottomSheet> createState() => _FilterBottomSheetState();
}

class _FilterBottomSheetState extends State<FilterBottomSheet> {
  late String? _selectedCategory;
  late RangeValues _priceRange;
  late double? _minRating;
  late bool _inStockOnly;

  @override
  void initState() {
    super.initState();
    _selectedCategory = widget.currentFilter.category;
    _priceRange = RangeValues(
      widget.currentFilter.minPrice ?? 0,
      widget.currentFilter.maxPrice ?? 1000,
    );
    _minRating = widget.currentFilter.rating;
    _inStockOnly = widget.currentFilter.inStock ?? false;
  }

  void _applyFilters() {
    final filter = widget.currentFilter.copyWith(
      category: _selectedCategory,
      minPrice: _priceRange.start > 0 ? _priceRange.start : null,
      maxPrice: _priceRange.end < 1000 ? _priceRange.end : null,
      rating: _minRating,
      inStock: _inStockOnly ? true : null,
      page: 1,
      clearCategory: _selectedCategory == null,
      clearMinPrice: _priceRange.start <= 0,
      clearMaxPrice: _priceRange.end >= 1000,
      clearRating: _minRating == null,
      clearInStock: !_inStockOnly,
    );
    widget.onApply(filter);
    Navigator.pop(context, filter);
  }

  void _resetFilters() {
    setState(() {
      _selectedCategory = null;
      _priceRange = const RangeValues(0, 1000);
      _minRating = null;
      _inStockOnly = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return DraggableScrollableSheet(
      initialChildSize: 0.75,
      minChildSize: 0.5,
      maxChildSize: 0.9,
      expand: false,
      builder: (_, scrollController) => Column(
        children: [
          // Handle bar
          Container(
            margin: const EdgeInsets.only(top: 12),
            width: 40,
            height: 4,
            decoration: BoxDecoration(
              color: Colors.grey.shade300,
              borderRadius: BorderRadius.circular(2),
            ),
          ),

          // Header
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Filters',
                  style: theme.textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                TextButton(
                  onPressed: _resetFilters,
                  child: const Text('Reset'),
                ),
              ],
            ),
          ),
          const Divider(height: 1),

          // Filter content
          Expanded(
            child: ListView(
              controller: scrollController,
              padding: const EdgeInsets.all(16),
              children: [
                // Category
                Text(
                  'Category',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 12),
                Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: widget.categories.map((category) {
                    final isSelected = _selectedCategory == category.slug;
                    return FilterChip(
                      label: Text(category.name),
                      selected: isSelected,
                      onSelected: (selected) {
                        setState(() {
                          _selectedCategory = selected ? category.slug : null;
                        });
                      },
                    );
                  }).toList(),
                ),
                const SizedBox(height: 24),

                // Price range
                Text(
                  'Price Range',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 12),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text('\$${_priceRange.start.toInt()}'),
                    Text('\$${_priceRange.end.toInt()}'),
                  ],
                ),
                RangeSlider(
                  values: _priceRange,
                  min: 0,
                  max: 1000,
                  divisions: 100,
                  labels: RangeLabels(
                    '\$${_priceRange.start.toInt()}',
                    '\$${_priceRange.end.toInt()}',
                  ),
                  onChanged: (values) {
                    setState(() => _priceRange = values);
                  },
                ),
                const SizedBox(height: 24),

                // Rating
                Text(
                  'Minimum Rating',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 12),
                Row(
                  children: List.generate(5, (index) {
                    final ratingValue = (index + 1).toDouble();
                    final isActive = _minRating != null && _minRating! >= ratingValue;
                    return GestureDetector(
                      onTap: () {
                        setState(() {
                          _minRating = _minRating == ratingValue ? null : ratingValue;
                        });
                      },
                      child: Padding(
                        padding: const EdgeInsets.only(right: 8),
                        child: Chip(
                          label: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Icon(
                                Icons.star,
                                size: 16,
                                color: isActive ? Colors.amber : Colors.grey,
                              ),
                              const SizedBox(width: 4),
                              Text('${index + 1}+'),
                            ],
                          ),
                          backgroundColor: isActive
                              ? theme.colorScheme.primaryContainer
                              : null,
                        ),
                      ),
                    );
                  }),
                ),
                const SizedBox(height: 24),

                // In stock
                SwitchListTile(
                  title: const Text('In Stock Only'),
                  value: _inStockOnly,
                  onChanged: (value) {
                    setState(() => _inStockOnly = value);
                  },
                  contentPadding: EdgeInsets.zero,
                ),
              ],
            ),
          ),

          // Apply button
          Padding(
            padding: const EdgeInsets.all(16),
            child: FilledButton(
              onPressed: _applyFilters,
              child: const Text('Apply Filters'),
            ),
          ),
        ],
      ),
    );
  }
}
