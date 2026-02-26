import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../constants/route_names.dart';
import '../../features/auth/presentation/pages/login_page.dart';
import '../../features/auth/presentation/pages/register_page.dart';
import '../../features/home/presentation/pages/home_page.dart';
import '../../features/shop/presentation/pages/product_list_page.dart';
import '../../features/shop/presentation/pages/product_detail_page.dart';
import '../../features/search/presentation/pages/search_page.dart';
import '../../features/cart/presentation/pages/cart_page.dart';
import '../../features/checkout/presentation/pages/checkout_page.dart';
import '../../features/checkout/presentation/pages/order_confirmation_page.dart';
import '../../features/orders/presentation/pages/orders_page.dart';
import '../../features/orders/presentation/pages/order_detail_page.dart';
import '../../features/returns/presentation/pages/return_list_page.dart';
import '../../features/returns/presentation/pages/return_request_page.dart';
import '../../features/tracking/presentation/pages/tracking_page.dart';
import '../../features/profile/presentation/pages/profile_page.dart';
import '../../features/profile/presentation/pages/addresses_page.dart';
import '../../features/wishlist/presentation/pages/wishlist_page.dart';
import '../../features/loyalty/presentation/pages/loyalty_page.dart';
import '../../features/affiliate/presentation/pages/affiliate_page.dart';
import '../../features/notifications/presentation/pages/notifications_page.dart';
import '../../features/chat/presentation/pages/conversation_list_page.dart';
import '../../features/chat/presentation/pages/chat_page.dart';
import '../../features/ai/presentation/pages/ai_assistant_page.dart';
import '../../features/reviews/presentation/pages/write_review_page.dart';

final _isAuthenticatedProvider = StateProvider<bool>((ref) => false);

final appRouterProvider = Provider<GoRouter>((ref) {
  final isAuthenticated = ref.watch(_isAuthenticatedProvider);

  return GoRouter(
    initialLocation: RouteNames.home,
    debugLogDiagnostics: true,
    redirect: (BuildContext context, GoRouterState state) {
      final loggingIn = state.matchedLocation == RouteNames.login;
      final registering = state.matchedLocation == RouteNames.register;

      final publicRoutes = [
        RouteNames.home,
        RouteNames.login,
        RouteNames.register,
        '/products',
        '/search',
      ];

      final isPublic = publicRoutes.any(
        (route) => state.matchedLocation.startsWith(route),
      );

      if (!isAuthenticated && !isPublic) {
        return RouteNames.login;
      }

      if (isAuthenticated && (loggingIn || registering)) {
        return RouteNames.home;
      }

      return null;
    },
    routes: [
      GoRoute(
        path: RouteNames.home,
        name: 'home',
        builder: (context, state) => const HomePage(),
      ),
      GoRoute(
        path: RouteNames.login,
        name: 'login',
        builder: (context, state) => const LoginPage(),
      ),
      GoRoute(
        path: RouteNames.register,
        name: 'register',
        builder: (context, state) => const RegisterPage(),
      ),
      GoRoute(
        path: RouteNames.products,
        name: 'products',
        builder: (context, state) => const ProductListPage(),
      ),
      GoRoute(
        path: RouteNames.productDetail,
        name: 'productDetail',
        builder: (context, state) => ProductDetailPage(
          slug: state.pathParameters['slug']!,
        ),
      ),
      GoRoute(
        path: RouteNames.search,
        name: 'search',
        builder: (context, state) => const SearchPage(),
      ),
      GoRoute(
        path: RouteNames.cart,
        name: 'cart',
        builder: (context, state) => const CartPage(),
      ),
      GoRoute(
        path: RouteNames.checkout,
        name: 'checkout',
        builder: (context, state) => const CheckoutPage(),
      ),
      GoRoute(
        path: RouteNames.orderConfirmation,
        name: 'orderConfirmation',
        builder: (context, state) => OrderConfirmationPage(
          orderId: state.pathParameters['id']!,
        ),
      ),
      GoRoute(
        path: RouteNames.account,
        name: 'account',
        builder: (context, state) => const ProfilePage(),
      ),
      GoRoute(
        path: RouteNames.accountOrders,
        name: 'accountOrders',
        builder: (context, state) => const OrdersPage(),
      ),
      GoRoute(
        path: RouteNames.orderDetail,
        name: 'orderDetail',
        builder: (context, state) => OrderDetailPage(
          orderId: state.pathParameters['id']!,
        ),
      ),
      GoRoute(
        path: RouteNames.accountAddresses,
        name: 'accountAddresses',
        builder: (context, state) => const AddressesPage(),
      ),
      GoRoute(
        path: RouteNames.accountReturns,
        name: 'accountReturns',
        builder: (context, state) => const ReturnListPage(),
      ),
      GoRoute(
        path: RouteNames.returnRequest,
        name: 'returnRequest',
        builder: (context, state) => const ReturnRequestPage(),
      ),
      GoRoute(
        path: RouteNames.accountWishlist,
        name: 'accountWishlist',
        builder: (context, state) => const WishlistPage(),
      ),
      GoRoute(
        path: RouteNames.accountLoyalty,
        name: 'accountLoyalty',
        builder: (context, state) => const LoyaltyPage(),
      ),
      GoRoute(
        path: RouteNames.accountAffiliate,
        name: 'accountAffiliate',
        builder: (context, state) => const AffiliatePage(),
      ),
      GoRoute(
        path: RouteNames.notifications,
        name: 'notifications',
        builder: (context, state) => const NotificationsPage(),
      ),
      GoRoute(
        path: RouteNames.chat,
        name: 'chat',
        builder: (context, state) => const ConversationListPage(),
      ),
      GoRoute(
        path: RouteNames.chatDetail,
        name: 'chatDetail',
        builder: (context, state) => ChatPage(
          conversationId: state.pathParameters['id']!,
        ),
      ),
      GoRoute(
        path: RouteNames.aiAssistant,
        name: 'aiAssistant',
        builder: (context, state) => const AIAssistantPage(),
      ),
      GoRoute(
        path: RouteNames.tracking,
        name: 'tracking',
        builder: (context, state) => TrackingPage(
          orderId: state.pathParameters['orderId']!,
        ),
      ),
      GoRoute(
        path: RouteNames.writeReview,
        name: 'writeReview',
        builder: (context, state) => WriteReviewPage(
          productId: state.uri.queryParameters['productId'] ?? '',
        ),
      ),
    ],
    errorBuilder: (context, state) => Scaffold(
      appBar: AppBar(title: const Text('Page Not Found')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error_outline, size: 64, color: Colors.grey),
            const SizedBox(height: 16),
            Text(
              'Page not found',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 8),
            Text(
              state.matchedLocation,
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: Colors.grey,
                  ),
            ),
            const SizedBox(height: 24),
            FilledButton(
              onPressed: () => context.go(RouteNames.home),
              child: const Text('Go Home'),
            ),
          ],
        ),
      ),
    ),
  );
});
