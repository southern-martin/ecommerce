import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/route_names.dart';
import '../../../../core/di/injection.dart';
import '../../data/auth_repository.dart';

/// A wrapper widget that checks authentication state and redirects
/// to the login page if the seller is not authenticated.
///
/// Wrap any page that requires authentication with this widget.
class SellerAuthWrapper extends StatefulWidget {
  final Widget child;

  const SellerAuthWrapper({super.key, required this.child});

  @override
  State<SellerAuthWrapper> createState() => _SellerAuthWrapperState();
}

class _SellerAuthWrapperState extends State<SellerAuthWrapper> {
  bool _isChecking = true;
  bool _isAuthenticated = false;

  @override
  void initState() {
    super.initState();
    _checkAuth();
  }

  Future<void> _checkAuth() async {
    final authRepo = getIt<SellerAuthRepository>();
    final authenticated = await authRepo.isAuthenticated();

    if (!mounted) return;

    if (!authenticated) {
      context.go(RouteNames.login);
      return;
    }

    setState(() {
      _isChecking = false;
      _isAuthenticated = true;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_isChecking) {
      return const Scaffold(
        body: Center(
          child: CircularProgressIndicator(),
        ),
      );
    }

    if (!_isAuthenticated) {
      return const SizedBox.shrink();
    }

    return widget.child;
  }
}
