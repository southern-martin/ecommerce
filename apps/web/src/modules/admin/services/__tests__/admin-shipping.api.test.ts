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

import { adminShippingApi } from '../admin-shipping.api';

describe('adminShippingApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getCarriers', () => {
    it('should fetch carriers and unwrap nested data', async () => {
      const carriers = [{ id: 'c1', code: 'ups', name: 'UPS', tracking_url_template: 'https://ups.com/{tracking}', is_active: true }];
      mockApiClient.get.mockResolvedValue({ data: { data: carriers } });

      const result = await adminShippingApi.getCarriers();

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/carriers');
      expect(result).toEqual(carriers);
    });

    it('should fall back to data when nested data is undefined', async () => {
      const carriers = [{ code: 'fedex', name: 'FedEx', tracking_url_template: '', is_active: true }];
      mockApiClient.get.mockResolvedValue({ data: carriers });

      const result = await adminShippingApi.getCarriers();

      expect(result).toEqual(carriers);
    });
  });

  describe('createCarrier', () => {
    it('should post carrier data and return created carrier', async () => {
      const data = { code: 'dhl', name: 'DHL', tracking_url_template: 'https://dhl.com/{id}', is_active: true };
      const created = { id: 'c2', ...data };
      mockApiClient.post.mockResolvedValue({ data: { data: created } });

      const result = await adminShippingApi.createCarrier(data);

      expect(mockApiClient.post).toHaveBeenCalledWith('/admin/carriers', data);
      expect(result).toEqual(created);
    });
  });

  describe('updateCarrier', () => {
    it('should patch carrier by code and return updated carrier', async () => {
      const updated = { id: 'c1', code: 'ups', name: 'UPS Ground', tracking_url_template: '', is_active: false };
      mockApiClient.patch.mockResolvedValue({ data: { data: updated } });

      const result = await adminShippingApi.updateCarrier('ups', { name: 'UPS Ground', is_active: false });

      expect(mockApiClient.patch).toHaveBeenCalledWith('/admin/carriers/ups', { name: 'UPS Ground', is_active: false });
      expect(result).toEqual(updated);
    });
  });
});
