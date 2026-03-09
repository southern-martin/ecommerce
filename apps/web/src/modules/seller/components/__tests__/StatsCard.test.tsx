import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { StatsCard } from '../StatsCard';
import { DollarSign } from 'lucide-react';

describe('StatsCard', () => {
  it('should render label and value', () => {
    render(<StatsCard icon={DollarSign} label="Revenue" value="$1,000" />);

    expect(screen.getByText('Revenue')).toBeInTheDocument();
    expect(screen.getByText('$1,000')).toBeInTheDocument();
  });

  it('should render positive trend with percentage', () => {
    render(<StatsCard icon={DollarSign} label="Revenue" value="$1,000" trend={12} />);

    expect(screen.getByText('12%')).toBeInTheDocument();
  });

  it('should render negative trend', () => {
    render(<StatsCard icon={DollarSign} label="Revenue" value="$1,000" trend={-5} />);

    expect(screen.getByText('5%')).toBeInTheDocument();
  });

  it('should not render trend when undefined', () => {
    const { container } = render(<StatsCard icon={DollarSign} label="Revenue" value="$1,000" />);

    // No trend element should be present — check no percentage text
    expect(container.querySelector('.text-green-600')).toBeNull();
    expect(container.querySelector('.text-red-600')).toBeNull();
  });
});
