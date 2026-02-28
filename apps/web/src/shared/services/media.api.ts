import apiClient from '@/shared/lib/api-client';

export interface MediaFile {
  id: string;
  owner_id: string;
  owner_type: string;
  file_name: string;
  original_name: string;
  content_type: string;
  size_bytes: number;
  url: string;
  thumbnail_url: string;
  status: string;
  created_at: string;
}

export async function uploadImage(
  file: File,
  ownerType = 'product',
): Promise<MediaFile> {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('owner_type', ownerType);
  const { data } = await apiClient.post('/media/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
  return data.media;
}
