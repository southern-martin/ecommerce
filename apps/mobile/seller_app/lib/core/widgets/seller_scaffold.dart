import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../constants/route_names.dart';
import '../di/injection.dart';
import '../../features/auth/data/auth_repository.dart';

class SellerScaffold extends StatefulWidget {
  final Widget child;

  const SellerScaffold({super.key, required this.child});

  @override
  State<SellerScaffold> createState() => _SellerScaffoldState();
}

class _SellerScaffoldState extends State<SellerScaffold> {
  int _selectedIndex = 0;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    _updateSelectedIndex();
  }

  void _updateSelectedIndex() {
    final location = GoRouterState.of(context).matchedLocation;
    if (location == RouteNames.dashboard) {
      _selectedIndex = 0;
    } else if (location.startsWith('/products')) {
      _selectedIndex = 1;
    } else if (location.startsWith('/orders')) {
      _selectedIndex = 2;
    } else {
      _selectedIndex = 3;
    }
  }

  void _onDestinationSelected(int index) {
    switch (index) {
      case 0:
        context.go(RouteNames.dashboard);
        break;
      case 1:
        context.go(RouteNames.products);
        break;
      case 2:
        context.go(RouteNames.orders);
        break;
      case 3:
        _showMoreMenu();
        break;
    }
  }

  void _showMoreMenu() {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const SizedBox(height: 8),
            Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: Colors.grey.shade300,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            const SizedBox(height: 16),
            Text(
              'More',
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 8),
            _MoreMenuItem(
              icon: Icons.assignment_return_outlined,
              label: 'Returns',
              onTap: () {
                Navigator.pop(context);
                this.context.push(RouteNames.returns);
              },
            ),
            _MoreMenuItem(
              icon: Icons.local_shipping_outlined,
              label: 'Shipments',
              onTap: () {
                Navigator.pop(context);
                this.context.push(RouteNames.shipments);
              },
            ),
            _MoreMenuItem(
              icon: Icons.confirmation_number_outlined,
              label: 'Coupons',
              onTap: () {
                Navigator.pop(context);
                this.context.push(RouteNames.coupons);
              },
            ),
            _MoreMenuItem(
              icon: Icons.analytics_outlined,
              label: 'Analytics',
              onTap: () {
                Navigator.pop(context);
                this.context.push(RouteNames.analytics);
              },
            ),
            _MoreMenuItem(
              icon: Icons.account_balance_wallet_outlined,
              label: 'Payouts',
              onTap: () {
                Navigator.pop(context);
                this.context.push(RouteNames.payouts);
              },
            ),
            const Divider(),
            _MoreMenuItem(
              icon: Icons.settings_outlined,
              label: 'Settings',
              onTap: () {
                Navigator.pop(context);
              },
            ),
            _MoreMenuItem(
              icon: Icons.logout,
              label: 'Logout',
              isDestructive: true,
              onTap: () async {
                Navigator.pop(context);
                final authRepo = getIt<SellerAuthRepository>();
                await authRepo.logout();
                if (mounted) {
                  context.go(RouteNames.login);
                }
              },
            ),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: widget.child,
      bottomNavigationBar: NavigationBar(
        selectedIndex: _selectedIndex,
        onDestinationSelected: _onDestinationSelected,
        destinations: const [
          NavigationDestination(
            icon: Icon(Icons.dashboard_outlined),
            selectedIcon: Icon(Icons.dashboard),
            label: 'Dashboard',
          ),
          NavigationDestination(
            icon: Icon(Icons.inventory_2_outlined),
            selectedIcon: Icon(Icons.inventory_2),
            label: 'Products',
          ),
          NavigationDestination(
            icon: Icon(Icons.receipt_long_outlined),
            selectedIcon: Icon(Icons.receipt_long),
            label: 'Orders',
          ),
          NavigationDestination(
            icon: Icon(Icons.more_horiz),
            selectedIcon: Icon(Icons.more_horiz),
            label: 'More',
          ),
        ],
      ),
    );
  }
}

class _MoreMenuItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final bool isDestructive;

  const _MoreMenuItem({
    required this.icon,
    required this.label,
    required this.onTap,
    this.isDestructive = false,
  });

  @override
  Widget build(BuildContext context) {
    final color = isDestructive
        ? Theme.of(context).colorScheme.error
        : Theme.of(context).colorScheme.onSurface;

    return ListTile(
      leading: Icon(icon, color: color),
      title: Text(label, style: TextStyle(color: color)),
      onTap: onTap,
    );
  }
}
