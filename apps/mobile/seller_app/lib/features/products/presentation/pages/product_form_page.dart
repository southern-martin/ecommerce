import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/di/injection.dart';
import '../../data/product_repository.dart';
import '../../domain/product_form_state.dart';
import '../widgets/variant_editor.dart';

/// A comprehensive product form page used for both creating and editing products.
///
/// When [productId] is provided, the form loads existing product data for editing.
/// Otherwise, it behaves as a create form.
class ProductFormPage extends StatefulWidget {
  final String? productId;

  const ProductFormPage({super.key, this.productId});

  @override
  State<ProductFormPage> createState() => _ProductFormPageState();
}

class _ProductFormPageState extends State<ProductFormPage> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _priceController = TextEditingController();
  final _compareAtPriceController = TextEditingController();
  final _stockController = TextEditingController();

  bool _isLoading = false;
  bool _isSaving = false;
  String _selectedCategoryId = '';
  List<String> _images = [];
  List<ProductVariantFormEntry> _variants = [];

  bool get _isEditing => widget.productId != null;

  static const _categories = [
    {'id': 'cat_1', 'name': 'Electronics'},
    {'id': 'cat_2', 'name': 'Clothing'},
    {'id': 'cat_3', 'name': 'Home & Garden'},
    {'id': 'cat_4', 'name': 'Sports & Outdoors'},
    {'id': 'cat_5', 'name': 'Books'},
  ];

  @override
  void initState() {
    super.initState();
    if (_isEditing) {
      _loadProduct();
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    _descriptionController.dispose();
    _priceController.dispose();
    _compareAtPriceController.dispose();
    _stockController.dispose();
    super.dispose();
  }

  Future<void> _loadProduct() async {
    setState(() => _isLoading = true);

    // In a real app, fetch the product by ID
    final repo = getIt<SellerProductRepository>();
    final result = await repo.getMyProducts();
    final product = result.products.firstWhere(
      (p) => p.id == widget.productId,
      orElse: () => result.products.first,
    );

    _nameController.text = product.name;
    _descriptionController.text = product.description;
    _priceController.text = product.price.toStringAsFixed(2);
    if (product.compareAtPrice != null) {
      _compareAtPriceController.text =
          product.compareAtPrice!.toStringAsFixed(2);
    }
    _stockController.text = product.stockQuantity.toString();
    _selectedCategoryId = product.categoryId;
    _images = List.from(product.imageUrls);
    _variants = product.variants
        .map((v) => ProductVariantFormEntry(
              name: v.name,
              value: v.value,
              priceModifier: v.priceModifier,
              stock: v.stock,
            ))
        .toList();

    setState(() => _isLoading = false);
  }

  Future<void> _pickImages() async {
    final picker = ImagePicker();
    final pickedFiles = await picker.pickMultiImage();

    if (pickedFiles.isNotEmpty) {
      setState(() {
        _images.addAll(pickedFiles.map((f) => f.path));
      });
    }
  }

  void _removeImage(int index) {
    setState(() {
      _images.removeAt(index);
    });
  }

  Future<void> _saveProduct({bool publish = false}) async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isSaving = true);

    final repo = getIt<SellerProductRepository>();
    final variants = _variants
        .map((v) => ProductVariant(
              id: '',
              name: v.name,
              value: v.value,
              priceModifier: v.priceModifier,
              stock: v.stock,
            ))
        .toList();

    try {
      if (_isEditing) {
        await repo.updateProduct(
          id: widget.productId!,
          name: _nameController.text.trim(),
          description: _descriptionController.text.trim(),
          price: double.parse(_priceController.text),
          compareAtPrice: _compareAtPriceController.text.isNotEmpty
              ? double.parse(_compareAtPriceController.text)
              : null,
          categoryId: _selectedCategoryId,
          stockQuantity: int.parse(_stockController.text),
          imageUrls: _images,
          variants: variants,
        );
      } else {
        await repo.createProduct(
          name: _nameController.text.trim(),
          description: _descriptionController.text.trim(),
          price: double.parse(_priceController.text),
          compareAtPrice: _compareAtPriceController.text.isNotEmpty
              ? double.parse(_compareAtPriceController.text)
              : null,
          categoryId: _selectedCategoryId,
          stockQuantity: int.parse(_stockController.text),
          imageUrls: _images,
          variants: variants,
        );
      }

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(_isEditing
                ? 'Product updated successfully'
                : 'Product created successfully'),
          ),
        );
        context.pop();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to save product: $e'),
            backgroundColor: Theme.of(context).colorScheme.error,
          ),
        );
      }
    } finally {
      if (mounted) setState(() => _isSaving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    if (_isLoading) {
      return Scaffold(
        appBar: AppBar(
          title: Text(_isEditing ? 'Edit Product' : 'New Product'),
        ),
        body: const Center(child: CircularProgressIndicator()),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: Text(_isEditing ? 'Edit Product' : 'New Product'),
        actions: [
          TextButton(
            onPressed: _isSaving ? null : () => _saveProduct(),
            child: const Text('Save Draft'),
          ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // Name
              TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(
                  labelText: 'Product Name',
                  hintText: 'Enter product name',
                ),
                validator: (val) {
                  if (val == null || val.trim().isEmpty) {
                    return 'Product name is required';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Description
              TextFormField(
                controller: _descriptionController,
                decoration: const InputDecoration(
                  labelText: 'Description',
                  hintText: 'Describe your product...',
                  alignLabelWithHint: true,
                ),
                maxLines: 5,
                validator: (val) {
                  if (val == null || val.trim().isEmpty) {
                    return 'Description is required';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Price and Compare at Price
              Row(
                children: [
                  Expanded(
                    child: TextFormField(
                      controller: _priceController,
                      decoration: const InputDecoration(
                        labelText: 'Price',
                        prefixText: '\$ ',
                      ),
                      keyboardType: const TextInputType.numberWithOptions(
                          decimal: true),
                      validator: (val) {
                        if (val == null || val.isEmpty) {
                          return 'Price is required';
                        }
                        if (double.tryParse(val) == null ||
                            double.parse(val) <= 0) {
                          return 'Enter a valid price';
                        }
                        return null;
                      },
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: TextFormField(
                      controller: _compareAtPriceController,
                      decoration: const InputDecoration(
                        labelText: 'Compare at Price',
                        prefixText: '\$ ',
                        hintText: 'Optional',
                      ),
                      keyboardType: const TextInputType.numberWithOptions(
                          decimal: true),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),

              // Category dropdown
              DropdownButtonFormField<String>(
                value: _selectedCategoryId.isNotEmpty
                    ? _selectedCategoryId
                    : null,
                decoration: const InputDecoration(
                  labelText: 'Category',
                ),
                items: _categories
                    .map((cat) => DropdownMenuItem(
                          value: cat['id'],
                          child: Text(cat['name']!),
                        ))
                    .toList(),
                onChanged: (val) {
                  if (val != null) {
                    setState(() => _selectedCategoryId = val);
                  }
                },
                validator: (val) {
                  if (val == null || val.isEmpty) {
                    return 'Please select a category';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Stock quantity
              TextFormField(
                controller: _stockController,
                decoration: const InputDecoration(
                  labelText: 'Stock Quantity',
                ),
                keyboardType: TextInputType.number,
                validator: (val) {
                  if (val == null || val.isEmpty) {
                    return 'Stock quantity is required';
                  }
                  if (int.tryParse(val) == null) {
                    return 'Enter a valid number';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 24),

              // Image uploader
              Text(
                'Images',
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 8),
              SizedBox(
                height: 100,
                child: ListView(
                  scrollDirection: Axis.horizontal,
                  children: [
                    // Add image button
                    InkWell(
                      onTap: _pickImages,
                      borderRadius: BorderRadius.circular(8),
                      child: Container(
                        width: 100,
                        height: 100,
                        decoration: BoxDecoration(
                          border: Border.all(color: Colors.grey.shade300),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(Icons.add_photo_alternate_outlined,
                                color: theme.colorScheme.primary),
                            const SizedBox(height: 4),
                            Text('Add',
                                style: theme.textTheme.bodySmall?.copyWith(
                                  color: theme.colorScheme.primary,
                                )),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    // Existing images
                    ..._images.asMap().entries.map((entry) {
                      return Padding(
                        padding: const EdgeInsets.only(right: 8),
                        child: Stack(
                          children: [
                            ClipRRect(
                              borderRadius: BorderRadius.circular(8),
                              child: entry.value.startsWith('http')
                                  ? Image.network(
                                      entry.value,
                                      width: 100,
                                      height: 100,
                                      fit: BoxFit.cover,
                                      errorBuilder: (_, __, ___) =>
                                          Container(
                                        width: 100,
                                        height: 100,
                                        color: Colors.grey.shade200,
                                        child: const Icon(Icons.broken_image),
                                      ),
                                    )
                                  : Container(
                                      width: 100,
                                      height: 100,
                                      color: Colors.grey.shade200,
                                      child: const Icon(Icons.image),
                                    ),
                            ),
                            Positioned(
                              top: 4,
                              right: 4,
                              child: GestureDetector(
                                onTap: () => _removeImage(entry.key),
                                child: Container(
                                  padding: const EdgeInsets.all(2),
                                  decoration: const BoxDecoration(
                                    color: Colors.red,
                                    shape: BoxShape.circle,
                                  ),
                                  child: const Icon(
                                    Icons.close,
                                    size: 14,
                                    color: Colors.white,
                                  ),
                                ),
                              ),
                            ),
                          ],
                        ),
                      );
                    }),
                  ],
                ),
              ),
              const SizedBox(height: 24),

              // Variant editor
              VariantEditor(
                variants: _variants,
                onChanged: (updated) {
                  setState(() => _variants = updated);
                },
              ),
              const SizedBox(height: 32),

              // Save / Publish buttons
              FilledButton(
                onPressed:
                    _isSaving ? null : () => _saveProduct(publish: true),
                child: _isSaving
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          color: Colors.white,
                        ),
                      )
                    : Text(_isEditing ? 'Update & Publish' : 'Publish'),
              ),
              const SizedBox(height: 12),
              OutlinedButton(
                onPressed: _isSaving ? null : () => _saveProduct(),
                child: const Text('Save as Draft'),
              ),
              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }
}
