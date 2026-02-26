import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../constants/route_names.dart';
import '../di/injection.dart';
import '../widgets/seller_scaffold.dart';
import '../../features/auth/data/auth_repository.dart';
import '../../features/auth/presentation/pages/login_page.dart';
import '../../features/dashboard/presentation/pages/dashboard_page.dart';
import '../../features/products/presentation/pages/product_list_page.dart';
import '../../features/products/presentation/pages/product_form_page.dart';
import '../../features/orders/presentation/pages/order_list_page.dart';
import '../../features/orders/presentation/pages/order_detail_page.dart';
import '../../features/returns/presentation/pages/return_list_page.dart';
import '../../features/returns/presentation/pages/return_detail_page.dart';
import '../../features/shipments/presentation/pages/shipment_list_page.dart';
import '../../features/shipments/presentation/pages/create_shipment_page.dart';
import '../../features/coupons/presentation/pages/coupon_list_page.dart';
import '../../features/coupons/presentation/pages/coupon_form_page.dart';
import '../../features/analytics/presentation/pages/analytics_page.dart';
import '../../features/payouts/presentation/pages/payouts_page.dart';

final _rootNavigatorKey = GlobalKey<NavigatorState>();
final _shellNavigatorKey = GlobalKey<NavigatorState>();

final GoRouter appRouter = GoRouter(
  navigatorKey: _rootNavigatorKey,
  initialLocation: RouteNames.dashboard,
  redirect: (BuildContext context, GoRouterState state) async {
    final authRepo = getIt<SellerAuthRepository>();
    final isLoggedIn = await authRepo.isAuthenticated();
    final isLoginRoute = state.matchedLocation == RouteNames.login;

    if (!isLoggedIn && !isLoginRoute) {
      return RouteNames.login;
    }

    if (isLoggedIn && isLoginRoute) {
      return RouteNames.dashboard;
    }

    return null;
  },
  routes: [
    GoRoute(
      path: RouteNames.login,
      builder: (context, state) => const SellerLoginPage(),
    ),
    ShellRoute(
      navigatorKey: _shellNavigatorKey,
      builder: (context, state, child) => SellerScaffold(child: child),
      routes: [
        GoRoute(
          path: RouteNames.dashboard,
          builder: (context, state) => const DashboardPage(),
        ),
        GoRoute(
          path: RouteNames.products,
          builder: (context, state) => const SellerProductListPage(),
        ),
        GoRoute(
          path: RouteNames.orders,
          builder: (context, state) => const SellerOrderListPage(),
        ),
      ],
    ),
    GoRoute(
      path: RouteNames.productNew,
      builder: (context, state) => const ProductFormPage(),
    ),
    GoRoute(
      path: RouteNames.productEdit,
      builder: (context, state) {
        final productId = state.pathParameters['id']!;
        return ProductFormPage(productId: productId);
      },
    ),
    GoRoute(
      path: RouteNames.orderDetail,
      builder: (context, state) {
        final orderId = state.pathParameters['id']!;
        return SellerOrderDetailPage(orderId: orderId);
      },
    ),
    GoRoute(
      path: RouteNames.returns,
      builder: (context, state) => const SellerReturnListPage(),
    ),
    GoRoute(
      path: RouteNames.returnDetail,
      builder: (context, state) {
        final returnId = state.pathParameters['id']!;
        return SellerReturnDetailPage(returnId: returnId);
      },
    ),
    GoRoute(
      path: RouteNames.shipments,
      builder: (context, state) => const ShipmentListPage(),
    ),
    GoRoute(
      path: '/shipments/new',
      builder: (context, state) => const CreateShipmentPage(),
    ),
    GoRoute(
      path: RouteNames.coupons,
      builder: (context, state) => const CouponListPage(),
    ),
    GoRoute(
      path: RouteNames.couponNew,
      builder: (context, state) => const CouponFormPage(),
    ),
    GoRoute(
      path: RouteNames.analytics,
      builder: (context, state) => const AnalyticsPage(),
    ),
    GoRoute(
      path: RouteNames.payouts,
      builder: (context, state) => const PayoutsPage(),
    ),
  ],
);
