import { lazy, Suspense } from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import { LayoutDashboard, Package, ShoppingCart, Users } from 'lucide-react';
import { RootLayout } from '@/shared/components/layout/RootLayout';
import { DashboardLayout } from '@/shared/components/layout/DashboardLayout';
import { LoadingSpinner } from '@/shared/components/feedback/LoadingSpinner';
import { PageNotFound } from '@/shared/components/feedback/PageNotFound';
import { AuthGuard } from '@/shared/guards/AuthGuard';
import { GuestGuard } from '@/shared/guards/GuestGuard';
import { RoleGuard } from '@/shared/guards/RoleGuard';
import { SellerGuard } from '@/shared/guards/SellerGuard';

// Lazy-loaded pages
const LoginPage = lazy(() => import('@/modules/auth/pages/LoginPage'));
const RegisterPage = lazy(() => import('@/modules/auth/pages/RegisterPage'));

const HomePage = lazy(() => import('@/modules/shop/pages/HomePage'));
const ProductListPage = lazy(() => import('@/modules/shop/pages/ProductListPage'));
const ProductDetailPage = lazy(() => import('@/modules/shop/pages/ProductDetailPage'));

const SearchResultsPage = lazy(() => import('@/modules/search/pages/SearchResultsPage'));

const CartPage = lazy(() => import('@/modules/cart/pages/CartPage'));

const CheckoutPage = lazy(() => import('@/modules/checkout/pages/CheckoutPage'));
const OrderConfirmationPage = lazy(() => import('@/modules/checkout/pages/OrderConfirmationPage'));

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

const PromotionsPage = lazy(() => import('@/modules/promotions/pages/PromotionsPage'));

const StaticPage = lazy(() => import('@/modules/cms/pages/StaticPage'));

// Seller
const SellerDashboardPage = lazy(() => import('@/modules/seller/pages/SellerDashboardPage'));
const SellerProductsPage = lazy(() => import('@/modules/seller/pages/SellerProductsPage'));
const SellerOrdersPage = lazy(() => import('@/modules/seller/pages/SellerOrdersPage'));

// Admin
const AdminDashboardPage = lazy(() => import('@/modules/admin/pages/AdminDashboardPage'));
const AdminUsersPage = lazy(() => import('@/modules/admin/pages/AdminUsersPage'));
const AdminOrdersPage = lazy(() => import('@/modules/admin/pages/AdminOrdersPage'));

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

  // Seller dashboard (separate layout)
  {
    path: '/seller',
    element: (
      <SellerGuard>
        <DashboardLayout sidebarItems={[
          { title: 'Dashboard', href: '/seller', icon: LayoutDashboard },
          { title: 'Products', href: '/seller/products', icon: Package },
          { title: 'Orders', href: '/seller/orders', icon: ShoppingCart },
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
        path: 'orders',
        element: (
          <SuspenseWrapper>
            <SellerOrdersPage />
          </SuspenseWrapper>
        ),
      },
    ],
  },

  // Admin dashboard (separate layout)
  {
    path: '/admin',
    element: (
      <RoleGuard allowedRoles={['admin']}>
        <DashboardLayout sidebarItems={[
          { title: 'Dashboard', href: '/admin', icon: LayoutDashboard },
          { title: 'Users', href: '/admin/users', icon: Users },
          { title: 'Orders', href: '/admin/orders', icon: ShoppingCart },
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
    ],
  },

  // Catch-all
  { path: '*', element: <PageNotFound /> },
]);
