import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@/test/test-utils';
import { UserManagementTable } from '../UserManagementTable';
import type { User } from '@/shared/types/user.types';

vi.mock('lucide-react', () => ({
  Trash2: (props: Record<string, unknown>) => <svg data-testid="trash-icon" {...props} />,
}));

describe('UserManagementTable', () => {
  const mockUsers: User[] = [
    {
      id: 'u1',
      email: 'alice@example.com',
      first_name: 'Alice',
      last_name: 'Smith',
      role: 'admin',
      is_verified: true,
      created_at: '2024-06-15T00:00:00Z',
    },
    {
      id: 'u2',
      email: 'bob@example.com',
      first_name: 'Bob',
      last_name: 'Jones',
      role: 'buyer',
      is_verified: false,
      created_at: '2024-08-20T00:00:00Z',
    },
  ];

  it('should render table headers', () => {
    render(<UserManagementTable users={mockUsers} />);

    expect(screen.getByText('User')).toBeDefined();
    expect(screen.getByText('Role')).toBeDefined();
    expect(screen.getByText('Status')).toBeDefined();
    expect(screen.getByText('Joined')).toBeDefined();
    expect(screen.getByText('Actions')).toBeDefined();
  });

  it('should render user names and emails', () => {
    render(<UserManagementTable users={mockUsers} />);

    expect(screen.getByText('Alice Smith')).toBeDefined();
    expect(screen.getByText('alice@example.com')).toBeDefined();
    expect(screen.getByText('Bob Jones')).toBeDefined();
    expect(screen.getByText('bob@example.com')).toBeDefined();
  });

  it('should display verification status badges', () => {
    render(<UserManagementTable users={mockUsers} />);

    expect(screen.getByText('Verified')).toBeDefined();
    expect(screen.getByText('Unverified')).toBeDefined();
  });

  it('should call onDelete when delete button is clicked', () => {
    const onDelete = vi.fn();
    render(<UserManagementTable users={mockUsers} onDelete={onDelete} />);

    const deleteButtons = screen.getAllByRole('button');
    fireEvent.click(deleteButtons[0]);

    expect(onDelete).toHaveBeenCalledWith('u1');
  });

  it('should render empty table when no users provided', () => {
    render(<UserManagementTable users={[]} />);

    expect(screen.getByText('User')).toBeDefined();
    expect(screen.queryAllByRole('button')).toHaveLength(0);
  });
});
