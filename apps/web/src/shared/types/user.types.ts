export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: 'buyer' | 'seller' | 'admin';
  avatar_url?: string;
  phone?: string;
  is_verified: boolean;
  created_at: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}
