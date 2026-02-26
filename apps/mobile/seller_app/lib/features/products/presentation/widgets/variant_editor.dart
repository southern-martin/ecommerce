import 'package:flutter/material.dart';

import '../../domain/product_form_state.dart';

/// An expandable section for managing product variants.
///
/// Displays a list of variant rows where each row contains fields for name,
/// value, price modifier, and stock. Supports adding and removing variants.
class VariantEditor extends StatelessWidget {
  final List<ProductVariantFormEntry> variants;
  final ValueChanged<List<ProductVariantFormEntry>> onChanged;

  const VariantEditor({
    super.key,
    required this.variants,
    required this.onChanged,
  });

  void _addVariant() {
    final updated = List<ProductVariantFormEntry>.from(variants);
    updated.add(const ProductVariantFormEntry());
    onChanged(updated);
  }

  void _removeVariant(int index) {
    final updated = List<ProductVariantFormEntry>.from(variants);
    updated.removeAt(index);
    onChanged(updated);
  }

  void _updateVariant(int index, ProductVariantFormEntry entry) {
    final updated = List<ProductVariantFormEntry>.from(variants);
    updated[index] = entry;
    onChanged(updated);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return ExpansionTile(
      title: Text(
        'Variants (${variants.length})',
        style: theme.textTheme.titleMedium?.copyWith(
          fontWeight: FontWeight.w600,
        ),
      ),
      subtitle: Text(
        'Add size, color, or other options',
        style: theme.textTheme.bodySmall,
      ),
      initiallyExpanded: variants.isNotEmpty,
      children: [
        if (variants.isEmpty)
          Padding(
            padding: const EdgeInsets.all(16),
            child: Text(
              'No variants added yet. Click the button below to add one.',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
              textAlign: TextAlign.center,
            ),
          ),
        ...variants.asMap().entries.map((entry) {
          final index = entry.key;
          final variant = entry.value;
          return _VariantRow(
            index: index,
            variant: variant,
            onChanged: (updated) => _updateVariant(index, updated),
            onRemove: () => _removeVariant(index),
          );
        }),
        Padding(
          padding: const EdgeInsets.all(16),
          child: OutlinedButton.icon(
            onPressed: _addVariant,
            icon: const Icon(Icons.add),
            label: const Text('Add Variant'),
          ),
        ),
      ],
    );
  }
}

class _VariantRow extends StatelessWidget {
  final int index;
  final ProductVariantFormEntry variant;
  final ValueChanged<ProductVariantFormEntry> onChanged;
  final VoidCallback onRemove;

  const _VariantRow({
    required this.index,
    required this.variant,
    required this.onChanged,
    required this.onRemove,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(
                  'Variant ${index + 1}',
                  style: theme.textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const Spacer(),
                IconButton(
                  icon: const Icon(Icons.delete_outline, color: Colors.red),
                  onPressed: onRemove,
                  iconSize: 20,
                  constraints: const BoxConstraints(),
                  padding: EdgeInsets.zero,
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Expanded(
                  child: TextFormField(
                    initialValue: variant.name,
                    decoration: const InputDecoration(
                      labelText: 'Name',
                      hintText: 'e.g. Size',
                      isDense: true,
                    ),
                    onChanged: (val) =>
                        onChanged(variant.copyWith(name: val)),
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: TextFormField(
                    initialValue: variant.value,
                    decoration: const InputDecoration(
                      labelText: 'Value',
                      hintText: 'e.g. Large',
                      isDense: true,
                    ),
                    onChanged: (val) =>
                        onChanged(variant.copyWith(value: val)),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Expanded(
                  child: TextFormField(
                    initialValue: variant.priceModifier != 0
                        ? variant.priceModifier.toString()
                        : '',
                    decoration: const InputDecoration(
                      labelText: 'Price +/-',
                      hintText: '0.00',
                      prefixText: '\$ ',
                      isDense: true,
                    ),
                    keyboardType:
                        const TextInputType.numberWithOptions(decimal: true),
                    onChanged: (val) => onChanged(variant.copyWith(
                      priceModifier: double.tryParse(val) ?? 0,
                    )),
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: TextFormField(
                    initialValue: variant.stock != 0
                        ? variant.stock.toString()
                        : '',
                    decoration: const InputDecoration(
                      labelText: 'Stock',
                      hintText: '0',
                      isDense: true,
                    ),
                    keyboardType: TextInputType.number,
                    onChanged: (val) => onChanged(variant.copyWith(
                      stock: int.tryParse(val) ?? 0,
                    )),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
