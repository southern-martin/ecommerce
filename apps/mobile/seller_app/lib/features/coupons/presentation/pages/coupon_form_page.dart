import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../../../../core/di/injection.dart';
import '../../data/coupon_repository.dart';

/// Form page for creating or editing a coupon.
///
/// When [coupon] is provided, the form pre-fills fields for editing.
/// Otherwise it behaves as a creation form.
class CouponFormPage extends StatefulWidget {
  final Coupon? coupon;

  const CouponFormPage({super.key, this.coupon});

  @override
  State<CouponFormPage> createState() => _CouponFormPageState();
}

class _CouponFormPageState extends State<CouponFormPage> {
  final _formKey = GlobalKey<FormState>();
  final _codeController = TextEditingController();
  final _valueController = TextEditingController();
  final _minOrderController = TextEditingController();
  final _maxUsesController = TextEditingController();

  String _discountType = 'percentage';
  DateTime? _expiryDate;
  bool _isSaving = false;

  bool get _isEditing => widget.coupon != null;

  @override
  void initState() {
    super.initState();
    if (_isEditing) {
      final c = widget.coupon!;
      _codeController.text = c.code;
      _discountType = c.discountType;
      _valueController.text = c.discountValue.toStringAsFixed(
          c.discountType == 'percentage' ? 0 : 2);
      if (c.minimumOrderAmount != null) {
        _minOrderController.text = c.minimumOrderAmount!.toStringAsFixed(2);
      }
      if (c.maxUses != null) {
        _maxUsesController.text = c.maxUses.toString();
      }
      _expiryDate = c.expiryDate;
    }
  }

  @override
  void dispose() {
    _codeController.dispose();
    _valueController.dispose();
    _minOrderController.dispose();
    _maxUsesController.dispose();
    super.dispose();
  }

  Future<void> _pickExpiryDate() async {
    final now = DateTime.now();
    final picked = await showDatePicker(
      context: context,
      initialDate: _expiryDate ?? now.add(const Duration(days: 30)),
      firstDate: now,
      lastDate: now.add(const Duration(days: 365 * 2)),
    );
    if (picked != null) {
      setState(() => _expiryDate = picked);
    }
  }

  Future<void> _saveCoupon() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isSaving = true);

    final repo = getIt<CouponRepository>();
    final data = <String, dynamic>{
      'code': _codeController.text.trim().toUpperCase(),
      'discountType': _discountType,
      'discountValue': double.parse(_valueController.text),
      'minimumOrderAmount': _minOrderController.text.isNotEmpty
          ? double.parse(_minOrderController.text)
          : null,
      'maxUses': _maxUsesController.text.isNotEmpty
          ? int.parse(_maxUsesController.text)
          : null,
      'expiryDate': _expiryDate,
    };

    try {
      if (_isEditing) {
        await repo.updateCoupon(widget.coupon!.id, data);
      } else {
        await repo.createCoupon(data);
      }

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(_isEditing
                ? 'Coupon updated successfully'
                : 'Coupon created successfully'),
          ),
        );
        Navigator.pop(context, true);
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to save coupon: $e'),
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

    return Scaffold(
      appBar: AppBar(
        title: Text(_isEditing ? 'Edit Coupon' : 'Create Coupon'),
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // Coupon code
              TextFormField(
                controller: _codeController,
                decoration: const InputDecoration(
                  labelText: 'Coupon Code',
                  hintText: 'e.g. SUMMER25',
                ),
                textCapitalization: TextCapitalization.characters,
                validator: (val) {
                  if (val == null || val.trim().isEmpty) {
                    return 'Coupon code is required';
                  }
                  if (val.trim().length < 3) {
                    return 'Code must be at least 3 characters';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),

              // Discount type segmented button
              Text('Discount Type',
                  style: theme.textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.w600,
                  )),
              const SizedBox(height: 8),
              SegmentedButton<String>(
                segments: const [
                  ButtonSegment(
                    value: 'percentage',
                    label: Text('Percentage'),
                    icon: Icon(Icons.percent),
                  ),
                  ButtonSegment(
                    value: 'fixed',
                    label: Text('Fixed'),
                    icon: Icon(Icons.attach_money),
                  ),
                ],
                selected: {_discountType},
                onSelectionChanged: (selected) {
                  setState(() => _discountType = selected.first);
                },
              ),
              const SizedBox(height: 20),

              // Discount value
              TextFormField(
                controller: _valueController,
                decoration: InputDecoration(
                  labelText: 'Discount Value',
                  prefixText: _discountType == 'fixed' ? '\$ ' : null,
                  suffixText: _discountType == 'percentage' ? '%' : null,
                ),
                keyboardType:
                    const TextInputType.numberWithOptions(decimal: true),
                validator: (val) {
                  if (val == null || val.isEmpty) {
                    return 'Discount value is required';
                  }
                  final parsed = double.tryParse(val);
                  if (parsed == null || parsed <= 0) {
                    return 'Enter a valid positive number';
                  }
                  if (_discountType == 'percentage' && parsed > 100) {
                    return 'Percentage cannot exceed 100';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Minimum order amount
              TextFormField(
                controller: _minOrderController,
                decoration: const InputDecoration(
                  labelText: 'Minimum Order Amount',
                  prefixText: '\$ ',
                  hintText: 'Optional',
                ),
                keyboardType:
                    const TextInputType.numberWithOptions(decimal: true),
                validator: (val) {
                  if (val != null && val.isNotEmpty) {
                    final parsed = double.tryParse(val);
                    if (parsed == null || parsed < 0) {
                      return 'Enter a valid amount';
                    }
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Max uses
              TextFormField(
                controller: _maxUsesController,
                decoration: const InputDecoration(
                  labelText: 'Max Uses',
                  hintText: 'Optional (unlimited if empty)',
                ),
                keyboardType: TextInputType.number,
                validator: (val) {
                  if (val != null && val.isNotEmpty) {
                    final parsed = int.tryParse(val);
                    if (parsed == null || parsed <= 0) {
                      return 'Enter a valid positive number';
                    }
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),

              // Expiry date picker
              Text('Expiry Date',
                  style: theme.textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.w600,
                  )),
              const SizedBox(height: 8),
              InkWell(
                onTap: _pickExpiryDate,
                borderRadius: BorderRadius.circular(8),
                child: InputDecorator(
                  decoration: InputDecoration(
                    suffixIcon: const Icon(Icons.calendar_today),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: Text(
                    _expiryDate != null
                        ? DateFormat('MMM d, yyyy').format(_expiryDate!)
                        : 'No expiry date (optional)',
                    style: TextStyle(
                      color: _expiryDate != null
                          ? theme.colorScheme.onSurface
                          : theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                ),
              ),
              if (_expiryDate != null) ...[
                const SizedBox(height: 4),
                Align(
                  alignment: Alignment.centerRight,
                  child: TextButton(
                    onPressed: () => setState(() => _expiryDate = null),
                    child: const Text('Clear date'),
                  ),
                ),
              ],
              const SizedBox(height: 32),

              // Save button
              SizedBox(
                height: 48,
                child: ElevatedButton(
                  onPressed: _isSaving ? null : _saveCoupon,
                  child: _isSaving
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Text('Save Coupon'),
                ),
              ),
              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }
}
