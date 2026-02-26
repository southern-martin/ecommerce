import { RoleGuard } from './RoleGuard';

export function SellerGuard({ children }: { children?: React.ReactNode }) {
  return <RoleGuard roles={['seller', 'admin']}>{children}</RoleGuard>;
}
