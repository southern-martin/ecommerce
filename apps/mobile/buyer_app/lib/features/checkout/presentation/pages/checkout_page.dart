import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/di/injection.dart';
import '../../data/checkout_repository.dart';

class CheckoutPage extends StatefulWidget {
  const CheckoutPage({super.key});

  @override
  State<CheckoutPage> createState() => _CheckoutPageState();
}

class _CheckoutPageState extends State<CheckoutPage> {
  final CheckoutRepository _checkoutRepo = getIt<CheckoutRepository>();

  int _currentStep = 0;
  bool _isLoading = false;
  bool _isPlacingOrder = false;

  // Address
  final _nameController = TextEditingController();
  final _streetController = TextEditingController();
  final _cityController = TextEditingController();
  final _stateController = TextEditingController();
  final _zipController = TextEditingController();
  final _phoneController = TextEditingController();
  String? _selectedAddressId;

  // Shipping
  List<ShippingRate> _shippingRates = [];
  String? _selectedShippingId;

  // Payment
  List<PaymentMethod> _paymentMethods = [];
  String? _selectedPaymentId;

  @override
  void dispose() {
    _nameController.dispose();
    _streetController.dispose();
    _cityController.dispose();
    _stateController.dispose();
    _zipController.dispose();
    _phoneController.dispose();
    super.dispose();
  }

  Future<void> _loadShippingRates() async {
    if (_selectedAddressId == null) return;
    setState(() => _isLoading = true);
    try {
      final rates = await _checkoutRepo.getShippingRates(_selectedAddressId!);
      if (mounted) {
        setState(() {
          _shippingRates = rates;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isLoading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load shipping rates: $e')),
        );
      }
    }
  }

  Future<void> _loadPaymentMethods() async {
    setState(() => _isLoading = true);
    try {
      final methods = await _checkoutRepo.getPaymentMethods();
      if (mounted) {
        setState(() {
          _paymentMethods = methods;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isLoading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load payment methods: $e')),
        );
      }
    }
  }

  Future<void> _placeOrder() async {
    if (_selectedShippingId == null ||
        _selectedPaymentId == null ||
        _selectedAddressId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please complete all steps')),
      );
      return;
    }

    setState(() => _isPlacingOrder = true);
    try {
      final result = await _checkoutRepo.placeOrder(
        shippingMethodId: _selectedShippingId!,
        paymentMethodId: _selectedPaymentId!,
        addressId: _selectedAddressId!,
      );
      if (mounted) {
        context.go('/order-confirmation/${result.orderId}');
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isPlacingOrder = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to place order: $e')),
        );
      }
    }
  }

  void _onStepContinue() {
    if (_currentStep == 0) {
      // Validate address
      if (_selectedAddressId == null &&
          (_streetController.text.isEmpty || _cityController.text.isEmpty)) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Please enter a shipping address')),
        );
        return;
      }
      _selectedAddressId ??= 'new_address';
      _loadShippingRates();
      setState(() => _currentStep = 1);
    } else if (_currentStep == 1) {
      if (_selectedShippingId == null) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Please select a shipping method')),
        );
        return;
      }
      _loadPaymentMethods();
      setState(() => _currentStep = 2);
    } else if (_currentStep == 2) {
      _placeOrder();
    }
  }

  void _onStepCancel() {
    if (_currentStep > 0) {
      setState(() => _currentStep -= 1);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Checkout'),
      ),
      body: Stepper(
        currentStep: _currentStep,
        onStepContinue: _onStepContinue,
        onStepCancel: _onStepCancel,
        controlsBuilder: (context, details) {
          return Padding(
            padding: const EdgeInsets.only(top: 16),
            child: Row(
              children: [
                if (_currentStep == 2)
                  Expanded(
                    child: ElevatedButton(
                      onPressed: _isPlacingOrder ? null : details.onStepContinue,
                      child: _isPlacingOrder
                          ? const SizedBox(
                              height: 20,
                              width: 20,
                              child: CircularProgressIndicator(strokeWidth: 2),
                            )
                          : const Text('Place Order'),
                    ),
                  )
                else
                  Expanded(
                    child: FilledButton(
                      onPressed: details.onStepContinue,
                      child: const Text('Continue'),
                    ),
                  ),
                if (_currentStep > 0) ...[
                  const SizedBox(width: 12),
                  TextButton(
                    onPressed: details.onStepCancel,
                    child: const Text('Back'),
                  ),
                ],
              ],
            ),
          );
        },
        steps: [
          // Step 1: Shipping Address
          Step(
            title: const Text('Shipping Address'),
            subtitle: _selectedAddressId != null
                ? const Text('Address selected')
                : null,
            isActive: _currentStep >= 0,
            state: _currentStep > 0 ? StepState.complete : StepState.indexed,
            content: _buildAddressStep(theme),
          ),

          // Step 2: Shipping Method
          Step(
            title: const Text('Shipping Method'),
            subtitle: _selectedShippingId != null
                ? Text(_shippingRates
                    .where((r) => r.id == _selectedShippingId)
                    .firstOrNull
                    ?.name ?? '')
                : null,
            isActive: _currentStep >= 1,
            state: _currentStep > 1 ? StepState.complete : StepState.indexed,
            content: _buildShippingStep(theme),
          ),

          // Step 3: Review & Pay
          Step(
            title: const Text('Review & Pay'),
            isActive: _currentStep >= 2,
            state: _currentStep > 2 ? StepState.complete : StepState.indexed,
            content: _buildPaymentStep(theme),
          ),
        ],
      ),
    );
  }

  Widget _buildAddressStep(ThemeData theme) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Enter your shipping address',
          style: theme.textTheme.bodyMedium?.copyWith(color: Colors.grey),
        ),
        const SizedBox(height: 16),
        TextField(
          controller: _nameController,
          decoration: const InputDecoration(
            labelText: 'Full Name',
            prefixIcon: Icon(Icons.person_outline),
          ),
        ),
        const SizedBox(height: 12),
        TextField(
          controller: _streetController,
          decoration: const InputDecoration(
            labelText: 'Street Address',
            prefixIcon: Icon(Icons.location_on_outlined),
          ),
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(
              child: TextField(
                controller: _cityController,
                decoration: const InputDecoration(labelText: 'City'),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: TextField(
                controller: _stateController,
                decoration: const InputDecoration(labelText: 'State'),
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(
              child: TextField(
                controller: _zipController,
                decoration: const InputDecoration(labelText: 'ZIP Code'),
                keyboardType: TextInputType.number,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: TextField(
                controller: _phoneController,
                decoration: const InputDecoration(labelText: 'Phone'),
                keyboardType: TextInputType.phone,
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildShippingStep(ThemeData theme) {
    if (_isLoading) {
      return const Center(
        child: Padding(
          padding: EdgeInsets.all(24),
          child: CircularProgressIndicator(),
        ),
      );
    }

    if (_shippingRates.isEmpty) {
      return const Padding(
        padding: EdgeInsets.all(16),
        child: Text('No shipping rates available for this address.'),
      );
    }

    return Column(
      children: _shippingRates.map((rate) {
        return RadioListTile<String>(
          title: Text(rate.name),
          subtitle: Text('${rate.carrier} - ${rate.estimatedDays}'),
          secondary: Text(
            rate.price > 0 ? '\$${rate.price.toStringAsFixed(2)}' : 'Free',
            style: theme.textTheme.titleSmall?.copyWith(
              fontWeight: FontWeight.bold,
              color: theme.colorScheme.primary,
            ),
          ),
          value: rate.id,
          groupValue: _selectedShippingId,
          onChanged: (value) => setState(() => _selectedShippingId = value),
          contentPadding: EdgeInsets.zero,
        );
      }).toList(),
    );
  }

  Widget _buildPaymentStep(ThemeData theme) {
    if (_isLoading) {
      return const Center(
        child: Padding(
          padding: EdgeInsets.all(24),
          child: CircularProgressIndicator(),
        ),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Select Payment Method',
          style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 8),
        ..._paymentMethods.map((method) {
          return RadioListTile<String>(
            title: Text(method.label),
            subtitle: method.last4 != null
                ? Text('${method.brand ?? ''} ending in ${method.last4}')
                : null,
            leading: Icon(
              method.type == 'card'
                  ? Icons.credit_card
                  : method.type == 'paypal'
                      ? Icons.account_balance_wallet
                      : Icons.payment,
            ),
            value: method.id,
            groupValue: _selectedPaymentId,
            onChanged: (value) => setState(() => _selectedPaymentId = value),
            contentPadding: EdgeInsets.zero,
          );
        }),
        const Divider(height: 32),
        Text(
          'Order Summary',
          style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                _buildSummaryRow(theme, 'Shipping Address', _streetController.text.isNotEmpty
                    ? '${_streetController.text}, ${_cityController.text}'
                    : 'Selected address'),
                const Divider(height: 16),
                _buildSummaryRow(theme, 'Shipping Method', _shippingRates
                    .where((r) => r.id == _selectedShippingId)
                    .firstOrNull
                    ?.name ?? 'Not selected'),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildSummaryRow(ThemeData theme, String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: theme.textTheme.bodySmall?.copyWith(color: Colors.grey)),
          Flexible(
            child: Text(
              value,
              style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w500),
              textAlign: TextAlign.end,
            ),
          ),
        ],
      ),
    );
  }
}
