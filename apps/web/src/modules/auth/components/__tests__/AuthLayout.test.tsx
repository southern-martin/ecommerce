import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { AuthLayout } from '../AuthLayout';

describe('AuthLayout', () => {
  it('renders the title', () => {
    render(<AuthLayout title="Sign In">content</AuthLayout>);
    expect(screen.getByText('Sign In')).toBeInTheDocument();
  });

  it('renders children content', () => {
    render(
      <AuthLayout title="Login">
        <form data-testid="login-form">
          <input placeholder="Email" />
        </form>
      </AuthLayout>
    );
    expect(screen.getByTestId('login-form')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Email')).toBeInTheDocument();
  });

  it('renders footer when provided', () => {
    render(
      <AuthLayout title="Login" footer={<span>Don't have an account?</span>}>
        <div>form</div>
      </AuthLayout>
    );
    expect(screen.getByText("Don't have an account?")).toBeInTheDocument();
  });

  it('does not render footer section when not provided', () => {
    const { container } = render(
      <AuthLayout title="Login">
        <div>form</div>
      </AuthLayout>
    );
    // Only CardContent children should be the children div, no extra footer div
    const cardContent = container.querySelector('.text-center.text-sm');
    expect(cardContent).toBeNull();
  });

  it('applies centered layout classes', () => {
    const { container } = render(<AuthLayout title="Test">content</AuthLayout>);
    const wrapper = container.firstChild as HTMLElement;
    expect(wrapper.className).toContain('flex');
    expect(wrapper.className).toContain('min-h-screen');
    expect(wrapper.className).toContain('items-center');
    expect(wrapper.className).toContain('justify-center');
  });
});
