import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/di/injection.dart';
import '../../data/shipment_repository.dart';

/// A form page for creating a new shipment.
///
/// Includes fields for order ID, carrier selection, tracking number,
/// optional weight, and optional dimensions (L x W x H).
class CreateShipmentPage extends StatefulWidget {
  const CreateShipmentPage({super.key});

  @override
  State<CreateShipmentPage> createState() => _CreateShipmentPageState();
}

class _CreateShipmentPageState extends State<CreateShipmentPage> {
  final _formKey = GlobalKey<FormState>();
  final _orderIdController = TextEditingController();
  final _trackingController = TextEditingController();
  final _weightController = TextEditingController();
  final _lengthController = TextEditingController();
  final _widthController = TextEditingController();
  final _heightController = TextEditingController();

  String? _selectedCarrier;
  bool _isSaving = false;

  static const _carriers = ['FedEx', 'UPS', 'USPS', 'DHL'];

  @override
  void dispose() {
    _orderIdController.dispose();
    _trackingController.dispose();
    _weightController.dispose();
    _lengthController.dispose();
    _widthController.dispose();
    _heightController.dispose();
    super.dispose();
  }

  Future<void> _createShipment() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isSaving = true);

    try {
      final weight = _weightController.text.isNotEmpty
          ? double.tryParse(_weightController.text)
          : null;

      ShipmentDimensions? dimensions;
      if (_lengthController.text.isNotEmpty &&
          _widthController.text.isNotEmpty &&
          _heightController.text.isNotEmpty) {
        dimensions = ShipmentDimensions(
          length: double.parse(_lengthController.text),
          width: double.parse(_widthController.text),
          height: double.parse(_heightController.text),
        );
      }

      await getIt<ShipmentRepository>().createShipment(
        orderId: _orderIdController.text.trim(),
        carrier: _selectedCarrier!,
        trackingNumber: _trackingController.text.trim(),
        weight: weight,
        dimensions: dimensions,
      );

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Shipment created successfully')),
        );
        context.pop();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to create shipment: $e'),
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
        title: const Text('Create Shipment'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // Order ID
              TextFormField(
                controller: _orderIdController,
                decoration: const InputDecoration(
                  labelText: 'Order ID',
                  hintText: 'Enter the order ID',
                  prefixIcon: Icon(Icons.receipt_long_outlined),
                ),
                validator: (val) {
                  if (val == null || val.trim().isEmpty) {
                    return 'Order ID is required';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Carrier dropdown
              DropdownButtonFormField<String>(
                value: _selectedCarrier,
                decoration: const InputDecoration(
                  labelText: 'Carrier',
                  prefixIcon: Icon(Icons.local_shipping_outlined),
                ),
                items: _carriers.map((carrier) {
                  return DropdownMenuItem(value: carrier, child: Text(carrier));
                }).toList(),
                onChanged: (val) {
                  setState(() => _selectedCarrier = val);
                },
                validator: (val) {
                  if (val == null || val.isEmpty) {
                    return 'Please select a carrier';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Tracking number
              TextFormField(
                controller: _trackingController,
                decoration: const InputDecoration(
                  labelText: 'Tracking Number',
                  hintText: 'Enter tracking number',
                  prefixIcon: Icon(Icons.qr_code),
                ),
                validator: (val) {
                  if (val == null || val.trim().isEmpty) {
                    return 'Tracking number is required';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 24),

              // Optional section header
              Text(
                'Optional Details',
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 12),

              // Weight
              TextFormField(
                controller: _weightController,
                decoration: const InputDecoration(
                  labelText: 'Weight (lbs)',
                  hintText: 'Optional',
                  prefixIcon: Icon(Icons.scale_outlined),
                ),
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                validator: (val) {
                  if (val != null && val.isNotEmpty) {
                    if (double.tryParse(val) == null || double.parse(val) <= 0) {
                      return 'Enter a valid weight';
                    }
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Dimensions (L x W x H)
              Text(
                'Dimensions (inches)',
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: theme.colorScheme.onSurfaceVariant,
                ),
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  Expanded(
                    child: TextFormField(
                      controller: _lengthController,
                      decoration: const InputDecoration(
                        labelText: 'L',
                        hintText: 'Length',
                      ),
                      keyboardType: const TextInputType.numberWithOptions(decimal: true),
                      validator: (val) {
                        if (val != null && val.isNotEmpty) {
                          if (double.tryParse(val) == null || double.parse(val) <= 0) {
                            return 'Invalid';
                          }
                        }
                        return null;
                      },
                    ),
                  ),
                  const Padding(
                    padding: EdgeInsets.symmetric(horizontal: 8),
                    child: Text('x'),
                  ),
                  Expanded(
                    child: TextFormField(
                      controller: _widthController,
                      decoration: const InputDecoration(
                        labelText: 'W',
                        hintText: 'Width',
                      ),
                      keyboardType: const TextInputType.numberWithOptions(decimal: true),
                      validator: (val) {
                        if (val != null && val.isNotEmpty) {
                          if (double.tryParse(val) == null || double.parse(val) <= 0) {
                            return 'Invalid';
                          }
                        }
                        return null;
                      },
                    ),
                  ),
                  const Padding(
                    padding: EdgeInsets.symmetric(horizontal: 8),
                    child: Text('x'),
                  ),
                  Expanded(
                    child: TextFormField(
                      controller: _heightController,
                      decoration: const InputDecoration(
                        labelText: 'H',
                        hintText: 'Height',
                      ),
                      keyboardType: const TextInputType.numberWithOptions(decimal: true),
                      validator: (val) {
                        if (val != null && val.isNotEmpty) {
                          if (double.tryParse(val) == null || double.parse(val) <= 0) {
                            return 'Invalid';
                          }
                        }
                        return null;
                      },
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 32),

              // Create button
              FilledButton(
                onPressed: _isSaving ? null : _createShipment,
                child: _isSaving
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          color: Colors.white,
                        ),
                      )
                    : const Text('Create Shipment'),
              ),
              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }
}
