import { describe, it, expect, vi, beforeEach } from 'vitest';

const mockApiClient = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  patch: vi.fn(),
  delete: vi.fn(),
}));

vi.mock('@/shared/lib/api-client', () => ({
  default: mockApiClient,
}));

import { sellerShippingApi } from '../seller-shipping.api';

describe('sellerShippingApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getShipments', () => {
    it('should fetch shipments with pagination', async () => {
      const response = { data: [{ id: 'sh1', tracking_number: 'TRK001' }], total: 1 };
      mockApiClient.get.mockResolvedValue({ data: response });

      const result = await sellerShippingApi.getShipments({ page: 1, page_size: 10 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/shipments', { params: { page: 1, page_size: 10 } });
      expect(result).toEqual(response);
    });

    it('should work without params', async () => {
      mockApiClient.get.mockResolvedValue({ data: { data: [] } });

      await sellerShippingApi.getShipments();

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/shipments', { params: undefined });
    });
  });

  describe('getCarriers', () => {
    it('should fetch available carriers', async () => {
      const carriers = [{ id: 'car1', carrier_code: 'fedex', carrier_name: 'FedEx', is_active: true }];
      mockApiClient.get.mockResolvedValue({ data: { data: carriers } });

      const result = await sellerShippingApi.getCarriers();

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/carriers');
      expect(result).toEqual(carriers);
    });
  });

  describe('setupCarrier', () => {
    it('should set up a carrier for the seller', async () => {
      const input = { carrier_code: 'ups', account_number: 'ACC123' };
      const created = { id: 'car2', ...input, carrier_name: 'UPS', is_active: true };
      mockApiClient.post.mockResolvedValue({ data: { data: created } });

      const result = await sellerShippingApi.setupCarrier(input);

      expect(mockApiClient.post).toHaveBeenCalledWith('/seller/carriers', input);
      expect(result).toEqual(created);
    });
  });

  describe('createShipment', () => {
    it('should create a shipment and return it', async () => {
      const input = {
        order_id: 'o1',
        carrier_code: 'fedex',
        origin: { street: '123 A St', city: 'LA', state: 'CA', postal_code: '90001', country: 'US' },
        destination: { street: '456 B St', city: 'NY', state: 'NY', postal_code: '10001', country: 'US' },
        weight_grams: 500,
        items: [{ product_id: 'p1', product_name: 'Widget', quantity: 1 }],
      };
      const shipment = { id: 'sh2', order_id: 'o1', tracking_number: 'TRK002', carrier: 'fedex', status: 'created' };
      mockApiClient.post.mockResolvedValue({ data: { shipment } });

      const result = await sellerShippingApi.createShipment(input);

      expect(mockApiClient.post).toHaveBeenCalledWith('/shipments', input);
      expect(result).toEqual(shipment);
    });
  });

  describe('getShipmentsByOrderId', () => {
    it('should return the first shipment for an order', async () => {
      const shipment = { id: 'sh1', order_id: 'o1', tracking_number: 'TRK001' };
      mockApiClient.get.mockResolvedValue({ data: { data: [shipment] } });

      const result = await sellerShippingApi.getShipmentsByOrderId('o1');

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/shipments', {
        params: { order_id: 'o1', page_size: 1 },
      });
      expect(result).toEqual(shipment);
    });

    it('should return null when no shipments exist for the order', async () => {
      mockApiClient.get.mockResolvedValue({ data: { data: [] } });

      const result = await sellerShippingApi.getShipmentsByOrderId('o999');

      expect(result).toBeNull();
    });
  });
});
