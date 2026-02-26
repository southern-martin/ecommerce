export const API_URL = import.meta.env.VITE_API_URL || '/api/v1';
export const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8000/ws';

export const ORDER_STATUSES = [
  'pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded',
] as const;

export const RETURN_STATUSES = [
  'requested', 'approved', 'rejected', 'shipped_back', 'received', 'refunded',
] as const;

export const ROLES = ['buyer', 'seller', 'admin'] as const;

export const LOYALTY_TIERS = ['bronze', 'silver', 'gold', 'platinum'] as const;

export const COUPON_TYPES = ['percentage', 'fixed_amount', 'free_shipping'] as const;

export const PAGE_SIZES = [10, 20, 50, 100] as const;
