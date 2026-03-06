import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { StatusBadge } from '../StatusBadge';

describe('StatusBadge', () => {
  it('renders the status text', () => {
    render(<StatusBadge status="pending" />);
    expect(screen.getByText('pending')).toBeInTheDocument();
  });

  it('replaces underscores with spaces', () => {
    render(<StatusBadge status="shipped_back" />);
    expect(screen.getByText('shipped back')).toBeInTheDocument();
  });

  it('applies the correct color class for known statuses', () => {
    const { container } = render(<StatusBadge status="delivered" />);
    const badge = container.firstChild as HTMLElement;
    expect(badge.className).toContain('bg-green-100');
    expect(badge.className).toContain('text-green-800');
  });

  it('applies yellow classes for pending status', () => {
    const { container } = render(<StatusBadge status="pending" />);
    const badge = container.firstChild as HTMLElement;
    expect(badge.className).toContain('bg-yellow-100');
    expect(badge.className).toContain('text-yellow-800');
  });

  it('applies red classes for cancelled status', () => {
    const { container } = render(<StatusBadge status="cancelled" />);
    const badge = container.firstChild as HTMLElement;
    expect(badge.className).toContain('bg-red-100');
    expect(badge.className).toContain('text-red-800');
  });

  it('applies blue classes for confirmed status', () => {
    const { container } = render(<StatusBadge status="confirmed" />);
    const badge = container.firstChild as HTMLElement;
    expect(badge.className).toContain('bg-blue-100');
    expect(badge.className).toContain('text-blue-800');
  });

  it('falls back to gray for unknown statuses', () => {
    const { container } = render(<StatusBadge status="some_unknown_status" />);
    const badge = container.firstChild as HTMLElement;
    expect(badge.className).toContain('bg-gray-100');
    expect(badge.className).toContain('text-gray-800');
    expect(screen.getByText('some unknown status')).toBeInTheDocument();
  });

  it('applies capitalize class', () => {
    const { container } = render(<StatusBadge status="active" />);
    const badge = container.firstChild as HTMLElement;
    expect(badge.className).toContain('capitalize');
  });

  it('merges additional className', () => {
    const { container } = render(<StatusBadge status="active" className="ml-2" />);
    const badge = container.firstChild as HTMLElement;
    expect(badge.className).toContain('ml-2');
  });
});
