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

import { adminProductApi } from '../admin-product.api';

describe('adminProductApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getCategories', () => {
    it('should fetch categories and return the array', async () => {
      const categories = [{ id: 'c1', name: 'Electronics', slug: 'electronics', created_at: '2024-01-01' }];
      mockApiClient.get.mockResolvedValue({ data: { categories } });

      const result = await adminProductApi.getCategories();

      expect(mockApiClient.get).toHaveBeenCalledWith('/categories');
      expect(result).toEqual(categories);
    });

    it('should return empty array when categories is undefined', async () => {
      mockApiClient.get.mockResolvedValue({ data: {} });

      const result = await adminProductApi.getCategories();

      expect(result).toEqual([]);
    });
  });

  describe('createCategory', () => {
    it('should post category data and return created category', async () => {
      const newCat = { name: 'Books', description: 'All books' };
      const created = { id: 'c2', ...newCat, slug: 'books', created_at: '2024-01-01' };
      mockApiClient.post.mockResolvedValue({ data: created });

      const result = await adminProductApi.createCategory(newCat);

      expect(mockApiClient.post).toHaveBeenCalledWith('/admin/categories', newCat);
      expect(result).toEqual(created);
    });
  });

  describe('updateCategory', () => {
    it('should patch category by id', async () => {
      const updated = { id: 'c1', name: 'Updated', slug: 'updated', created_at: '2024-01-01' };
      mockApiClient.patch.mockResolvedValue({ data: updated });

      const result = await adminProductApi.updateCategory('c1', { name: 'Updated' });

      expect(mockApiClient.patch).toHaveBeenCalledWith('/admin/categories/c1', { name: 'Updated' });
      expect(result).toEqual(updated);
    });
  });

  describe('deleteCategory', () => {
    it('should delete category by id', async () => {
      mockApiClient.delete.mockResolvedValue({});

      await adminProductApi.deleteCategory('c1');

      expect(mockApiClient.delete).toHaveBeenCalledWith('/admin/categories/c1');
    });
  });

  describe('getAttributes', () => {
    it('should fetch attributes and return the array', async () => {
      const attributes = [{ id: 'a1', name: 'Color', type: 'select', required: true, filterable: true, created_at: '2024-01-01' }];
      mockApiClient.get.mockResolvedValue({ data: { attributes } });

      const result = await adminProductApi.getAttributes();

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/attributes');
      expect(result).toEqual(attributes);
    });

    it('should return empty array when attributes is undefined', async () => {
      mockApiClient.get.mockResolvedValue({ data: {} });

      const result = await adminProductApi.getAttributes();

      expect(result).toEqual([]);
    });
  });

  describe('createAttribute', () => {
    it('should post attribute data', async () => {
      const data = { name: 'Size', type: 'select' as const };
      const created = { id: 'a2', ...data, required: false, filterable: false, created_at: '2024-01-01' };
      mockApiClient.post.mockResolvedValue({ data: created });

      const result = await adminProductApi.createAttribute(data);

      expect(mockApiClient.post).toHaveBeenCalledWith('/admin/attributes', data);
      expect(result).toEqual(created);
    });
  });

  describe('deleteAttribute', () => {
    it('should delete attribute by id', async () => {
      mockApiClient.delete.mockResolvedValue({});

      await adminProductApi.deleteAttribute('a1');

      expect(mockApiClient.delete).toHaveBeenCalledWith('/admin/attributes/a1');
    });
  });

  describe('getAttributeGroups', () => {
    it('should fetch attribute groups', async () => {
      const groups = [{ id: 'g1', name: 'Physical', slug: 'physical', sort_order: 1, created_at: '2024-01-01', updated_at: '2024-01-01' }];
      mockApiClient.get.mockResolvedValue({ data: { attribute_groups: groups } });

      const result = await adminProductApi.getAttributeGroups();

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/attribute-groups');
      expect(result).toEqual(groups);
    });
  });

  describe('addAttributeToGroup', () => {
    it('should post attribute assignment to group', async () => {
      mockApiClient.post.mockResolvedValue({});

      await adminProductApi.addAttributeToGroup('g1', { attribute_id: 'a1', sort_order: 1 });

      expect(mockApiClient.post).toHaveBeenCalledWith('/admin/attribute-groups/g1/attributes', {
        attribute_id: 'a1',
        sort_order: 1,
      });
    });
  });

  describe('removeAttributeFromGroup', () => {
    it('should delete attribute from group', async () => {
      mockApiClient.delete.mockResolvedValue({});

      await adminProductApi.removeAttributeFromGroup('g1', 'a1');

      expect(mockApiClient.delete).toHaveBeenCalledWith('/admin/attribute-groups/g1/attributes/a1');
    });
  });
});
