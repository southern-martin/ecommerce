import type { LucideIcon } from 'lucide-react';
import { Breadcrumbs, type BreadcrumbItem } from '@/shared/components/ui/breadcrumbs';

interface PageLayoutProps {
  children: React.ReactNode;
  title?: string;
  subtitle?: string;
  icon?: LucideIcon;
  breadcrumbs?: BreadcrumbItem[];
  actions?: React.ReactNode;
  fullWidth?: boolean;
}

export function PageLayout({
  children,
  title,
  subtitle,
  icon: Icon,
  breadcrumbs,
  actions,
  fullWidth = false,
}: PageLayoutProps) {
  const hasHeader = title || breadcrumbs;

  return (
    <div className={fullWidth ? '' : 'container mx-auto px-4 py-8'}>
      {hasHeader && (
        <div className="mb-8">
          {breadcrumbs && breadcrumbs.length > 0 && (
            <div className="mb-4">
              <Breadcrumbs items={breadcrumbs} />
            </div>
          )}
          {title && (
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                {Icon && (
                  <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
                    <Icon className="h-5 w-5" />
                  </div>
                )}
                <div>
                  <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">{title}</h1>
                  {subtitle && (
                    <p className="mt-1 text-muted-foreground">{subtitle}</p>
                  )}
                </div>
              </div>
              {actions && <div className="flex items-center gap-2">{actions}</div>}
            </div>
          )}
        </div>
      )}
      {children}
    </div>
  );
}
