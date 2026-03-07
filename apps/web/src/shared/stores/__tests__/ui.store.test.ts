import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useUIStore } from '../ui.store';

describe('useUIStore', () => {
  beforeEach(() => {
    useUIStore.setState({
      sidebarOpen: true,
      theme: 'light',
      cartDrawerOpen: false,
    });
    // Reset DOM classList
    document.documentElement.classList.remove('dark');
  });

  describe('initial state', () => {
    it('has sidebar open by default', () => {
      expect(useUIStore.getState().sidebarOpen).toBe(true);
    });

    it('has light theme by default', () => {
      expect(useUIStore.getState().theme).toBe('light');
    });

    it('has cart drawer closed by default', () => {
      expect(useUIStore.getState().cartDrawerOpen).toBe(false);
    });
  });

  describe('toggleSidebar', () => {
    it('closes sidebar when open', () => {
      useUIStore.getState().toggleSidebar();
      expect(useUIStore.getState().sidebarOpen).toBe(false);
    });

    it('opens sidebar when closed', () => {
      useUIStore.setState({ sidebarOpen: false });
      useUIStore.getState().toggleSidebar();
      expect(useUIStore.getState().sidebarOpen).toBe(true);
    });

    it('toggles back and forth', () => {
      useUIStore.getState().toggleSidebar();
      expect(useUIStore.getState().sidebarOpen).toBe(false);
      useUIStore.getState().toggleSidebar();
      expect(useUIStore.getState().sidebarOpen).toBe(true);
    });
  });

  describe('toggleTheme', () => {
    it('switches from light to dark', () => {
      useUIStore.getState().toggleTheme();
      expect(useUIStore.getState().theme).toBe('dark');
    });

    it('switches from dark to light', () => {
      useUIStore.setState({ theme: 'dark' });
      useUIStore.getState().toggleTheme();
      expect(useUIStore.getState().theme).toBe('light');
    });

    it('adds dark class to document when switching to dark', () => {
      useUIStore.getState().toggleTheme();
      expect(document.documentElement.classList.contains('dark')).toBe(true);
    });

    it('removes dark class from document when switching to light', () => {
      document.documentElement.classList.add('dark');
      useUIStore.setState({ theme: 'dark' });
      useUIStore.getState().toggleTheme();
      expect(document.documentElement.classList.contains('dark')).toBe(false);
    });
  });

  describe('setCartDrawerOpen', () => {
    it('opens cart drawer', () => {
      useUIStore.getState().setCartDrawerOpen(true);
      expect(useUIStore.getState().cartDrawerOpen).toBe(true);
    });

    it('closes cart drawer', () => {
      useUIStore.setState({ cartDrawerOpen: true });
      useUIStore.getState().setCartDrawerOpen(false);
      expect(useUIStore.getState().cartDrawerOpen).toBe(false);
    });
  });
});
