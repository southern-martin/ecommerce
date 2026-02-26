import { lazy, Suspense } from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import {
  LayoutDashboard, Package, ShoppingCart, Users, Ticket,
  Truck, RotateCcw, BarChart3, Tag, FileText, Star,
  AlertTriangle, Link, Calculator, Settings,
} from 'lucide-react';
import { RootLayout } from '@/shared/components/layout/RootLayout';
import { DashboardLayout } from '@/shared/components/layout/DashboardLayout';
import { LoadingSpinner } from '@/shared/components/feedback/LoadingSpinner';
import { PageNotFound } from '@/shared/components/feedback/PageNotFound';
import { AuthGuard } from '@/shared/guards/AuthGuard';
import { GuestGuard } from '@/shared/guards/GuestGuard';
import { RoleGuard } from '@/shared/guards/RoleGuard';
import { SellerGuard } from '@/shared/guards/SellerGuard';

// ── Public pages ──
const LoginPage = lazy(() => import('@/modules/auth/pages/LoginPage'));
const RegisterPage = lazy(() => import('@/modules/auth/pages/RegisterPage'));
const HomePage = lazy(() => import('@/modules/shop/pages/HomePage'));
const ProductListPage = lazy(() => import('@/modules/shop/pages/ProductListPage'));
const ProductDetailPage = lazy(() => import('@/modules/shop/pages/ProductDetailPage'));
const SearchResultsPage = lazy(() => import('@/modules/search/pages/SearchResultsPage'));
const PromotionsPage = lazy(() => import('@/modules/promotions/pages/PromotionsPage'));
const StaticPage = lazy(() => import('@/modules/cms/pages/StaticPage'));

// ── Authenticated pages ──
const CartPage = lazy(() => import('@/modules/cart/pages/CartPage'));
const CheckoutPage = lazy(() => import('@/modules/checkout/pages/CheckoutPage'));
const OrderConfirmationPage = lazy(() => import('@/modules/checkout/pages/OrderConfirmationPage'));

// ── Account pages ──
const ProfilePage = lazy(() => import('@/modules/account/pages/ProfilePage'));
const OrdersPage = lazy(() => import('@/modules/account/pages/OrdersPage'));
const AddressesPage = lazy(() => import('@/modules/account/pages/AddressesPage'));
const ReturnListPage = lazy(() => import('@/modules/returns/pages/ReturnListPage'));
const ReturnRequestPage = lazy(() => import('@/modules/returns/pages/ReturnRequestPage'));
const NotificationsPage = lazy(() => import('@/modules/notifications/pages/NotificationsPage'));
const ChatPage = lazy(() => import('@/modules/chat/pages/ChatPage'));
const ShippingTrackingPage = lazy(() => import('@/modules/shipping/pages/ShippingTrackingPage'));
const LoyaltyDashboardPage = lazy(() => import('@/modules/loyalty/pages/LoyaltyDashboardPage'));
const AffiliateDashboardPage = lazy(() => import('@/modules/affiliate/pages/AffiliateDashboardPage'));

// ── Seller pages ──
const SellerDashboardPage = lazy(() => import('@/modules/seller/pages/SellerDashboardPage'));
const SellerProductsPage = lazy(() => import('@/modules/seller/pages/SellerProductsPage'));
const SellerProductNewPage = lazy(() => import('@/modules/seller/pages/SellerProductNewPage'));
const SellerProductEditPage = lazy(() => import('@/modules/seller/pages/SellerProductEditPage'));
const SellerOrdersPage = lazy(() => import('@/modules/seller/pages/SellerOrdersPage'));
const SellerOrderDetailPage = lazy(() => import('@/modules/seller/pages/SellerOrderDetailPage'));
const SellerCouponsPage = lazy(() => import('@/modules/seller/pages/SellerCouponsPage'));
const SellerShipmentsPage = lazy(() => import('@/modules/seller/pages/SellerShipmentsPage'));
const SellerReturnsPage = lazy(() => import('@/modules/seller/pages/SellerReturnsPage'));
const SellerAnalyticsPage = lazy(() => import('@/modules/seller/pages/SellerAnalyticsPage'));

// ── Admin pages ──
const AdminDashboardPage = lazy(() => import('@/modules/admin/pages/AdminDashboardPage'));
const AdminUsersPage = lazy(() => import('@/modules/admin/pages/AdminUsersPage'));
const AdminOrdersPage = lazy(() => import('@/modules/admin/pages/AdminOrdersPage'));
const AdminProductsPage = lazy(() => import('@/modules/admin/pages/AdminProductsPage'));
const AdminPromotionsPage = lazy(() => import('@/modules/admin/pages/AdminPromotionsPage'));
const AdminCmsPage = lazy(() => import('@/modules/admin/pages/AdminCmsPage'));
const AdminReviewsPage = lazy(() => import('@/modules/admin/pages/AdminReviewsPage'));
const AdminDisputesPage = lazy(() => import('@/modules/admin/pages/AdminDisputesPage'));
const AdminCarriersPage = lazy(() => import('@/modules/admin/pages/AdminCarriersPage'));
const AdminAffiliatePage = lazy(() => import('@/modules/admin/pages/AdminAffiliatePage'));
const AdminTaxPage = lazy(() => import('@/modules/admin/pages/AdminTaxPage'));
const AdminSettingsPage = lazy(() => import('@/modules/admin/pages/AdminSettingsPage'));

function SuspenseWrapper({ children }: { children: React.ReactNode }) {
  return <Suspense fallback={<LoadingSpinner message="Loading..." />}>{children}</Suspense>;
}

export const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    errorElement: <PageNotFound />,
    children: [
      // Public routes
      {
        index: true,
        element: (
          <SuspenseWrapper>
            <HomePage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'products',
        element: (
          <SuspenseWrapper>
            <ProductListPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'products/:slug',
        element: (
          <SuspenseWrapper>
            <ProductDetailPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'search',
        element: (
          <SuspenseWrapper>
            <SearchResultsPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'promotions',
        element: (
          <SuspenseWrapper>
            <PromotionsPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'pages/:slug',
        element: (
          <SuspenseWrapper>
            <StaticPage />
          </SuspenseWrapper>
        ),
      },

      // Guest-only routes (redirect to home if logged in)
      {
        path: 'login',
        element: (
          <GuestGuard>
            <SuspenseWrapper>
              <LoginPage />
            </SuspenseWrapper>
          </GuestGuard>
        ),
      },
      {
        path: 'register',
        element: (
          <GuestGuard>
            <SuspenseWrapper>
              <RegisterPage />
            </SuspenseWrapper>
          </GuestGuard>
        ),
      },

      // Authenticated routes
      {
        path: 'cart',
        element: (
          <AuthGuard>
            <SuspenseWrapper>
              <CartPage />
            </SuspenseWrapper>
          </AuthGuard>
        ),
      },
      {
        path: 'checkout',
        element: (
          <AuthGuard>
            <SuspenseWrapper>
              <CheckoutPage />
            </SuspenseWrapper>
          </AuthGuard>
        ),
      },
      {
        path: 'order-confirmation/:orderId',
        element: (
          <AuthGuard>
            <SuspenseWrapper>
              <OrderConfirmationPage />
            </SuspenseWrapper>
          </AuthGuard>
        ),
      },

      // Account section
      {
        path: 'account',
        element: <AuthGuard />,
        children: [
          { index: true, element: <Navigate to="profile" replace /> },
          {
            path: 'profile',
            element: (
              <SuspenseWrapper>
                <ProfilePage />
              </SuspenseWrapper>
            ),
          },
          {
            path: 'orders',
            element: (
              <SuspenseWrapper>
                <OrdersPage />
              </SuspenseWrapper>
            ),
          },
          {
            path: 'addresses',
            element: (
              <SuspenseWrapper>
                <AddressesPage />
              </SuspenseWrapper>
            ),
          },
          {
            path: 'returns',
            element: (
              <SuspenseWrapper>
                <ReturnListPage />
              </SuspenseWrapper>
            ),
          },
          {
            path: 'returns/new',
            element: (
              <SuspenseWrapper>
                <ReturnRequestPage />
              </SuspenseWrapper>
            ),
          },
          {
            path: 'notifications',
            element: (
              <SuspenseWrapper>
                <NotificationsPage />
              </SuspenseWrapper>
            ),
          },
          {
            path: 'chat',
            element: (
              <SuspenseWrapper>
                <ChatPage />
              </SuspenseWrapper>
            ),
          },
          {
            path: 'shipping/:orderId',
            element: (
              <SuspenseWrapper>
                <ShippingTrackingPage />
              </SuspenseWrapper>
            ),
          },
          {
            path: 'loyalty',
            element: (
              <SuspenseWrapper>
                <LoyaltyDashboardPage />
              </SuspenseWrapper>
            ),
          },
          {
            path: 'affiliate',
            element: (
              <SuspenseWrapper>
                <AffiliateDashboardPage />
              </SuspenseWrapper>
            ),
          },
        ],
      },
    ],
  },

  // ── Seller dashboard ──
  {
    path: '/seller',
    element: (
      <SellerGuard>
        <DashboardLayout sidebarItems={[
          { title: 'Dashboard', href: '/seller', icon: LayoutDashboard },
          { title: 'Products', href: '/seller/products', icon: Package },
          { title: 'Orders', href: '/seller/orders', icon: ShoppingCart },
          { title: 'Coupons', href: '/seller/coupons', icon: Ticket },
          { title: 'Shipments', href: '/seller/shipments', icon: Truck },
          { title: 'Returns', href: '/seller/returns', icon: RotateCcw },
          { title: 'Analytics', href: '/seller/analytics', icon: BarChart3 },
        ]} />
      </SellerGuard>
    ),
    errorElement: <PageNotFound />,
    children: [
      {
        index: true,
        element: (
          <SuspenseWrapper>
            <SellerDashboardPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'products',
        element: (
          <SuspenseWrapper>
            <SellerProductsPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'products/new',
        element: (
          <SuspenseWrapper>
            <SellerProductNewPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'products/:id/edit',
        element: (
          <SuspenseWrapper>
            <SellerProductEditPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'orders',
        element: (
          <SuspenseWrapper>
            <SellerOrdersPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'orders/:id',
        element: (
          <SuspenseWrapper>
            <SellerOrderDetailPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'coupons',
        element: (
          <SuspenseWrapper>
            <SellerCouponsPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'shipments',
        element: (
          <SuspenseWrapper>
            <SellerShipmentsPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'returns',
        element: (
          <SuspenseWrapper>
            <SellerReturnsPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'analytics',
        element: (
          <SuspenseWrapper>
            <SellerAnalyticsPage />
          </SuspenseWrapper>
        ),
      },
    ],
  },

  // ── Admin dashboard ──
  {
    path: '/admin',
    element: (
      <RoleGuard allowedRoles={['admin']}>
        <DashboardLayout sidebarItems={[
          { title: 'Dashboard', href: '/admin', icon: LayoutDashboard },
          { title: 'Users', href: '/admin/users', icon: Users },
          { title: 'Orders', href: '/admin/orders', icon: ShoppingCart },
          { title: 'Products', href: '/admin/products', icon: Package },
          { title: 'Promotions', href: '/admin/promotions', icon: Tag },
          { title: 'CMS', href: '/admin/cms', icon: FileText },
          { title: 'Reviews', href: '/admin/reviews', icon: Star },
          { title: 'Disputes', href: '/admin/disputes', icon: AlertTriangle },
          { title: 'Carriers', href: '/admin/carriers', icon: Truck },
          { title: 'Affiliates', href: '/admin/affiliates', icon: Link },
          { title: 'Tax', href: '/admin/tax', icon: Calculator },
          { title: 'Settings', href: '/admin/settings', icon: Settings },
        ]} />
      </RoleGuard>
    ),
    errorElement: <PageNotFound />,
    children: [
      {
        index: true,
        element: (
          <SuspenseWrapper>
            <AdminDashboardPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'users',
        element: (
          <SuspenseWrapper>
            <AdminUsersPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'orders',
        element: (
          <SuspenseWrapper>
            <AdminOrdersPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'products',
        element: (
          <SuspenseWrapper>
            <AdminProductsPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'promotions',
        element: (
          <SuspenseWrapper>
            <AdminPromotionsPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'cms',
        element: (
          <SuspenseWrapper>
            <AdminCmsPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'reviews',
        element: (
          <SuspenseWrapper>
            <AdminReviewsPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'disputes',
        element: (
          <SuspenseWrapper>
            <AdminDisputesPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'carriers',
        element: (
          <SuspenseWrapper>
            <AdminCarriersPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'affiliates',
        element: (
          <SuspenseWrapper>
            <AdminAffiliatePage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'tax',
        element: (
          <SuspenseWrapper>
            <AdminTaxPage />
          </SuspenseWrapper>
        ),
      },
      {
        path: 'settings',
        element: (
          <SuspenseWrapper>
            <AdminSettingsPage />
          </SuspenseWrapper>
        ),
      },
    ],
  },

  // Catch-all
  { path: '*', element: <PageNotFound /> },
]);
