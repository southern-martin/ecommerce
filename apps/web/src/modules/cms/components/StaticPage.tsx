import { Skeleton } from '@/shared/components/ui/skeleton';
import type { StaticPage as StaticPageType } from '../services/cms.api';

interface StaticPageProps {
  page: StaticPageType;
  isLoading?: boolean;
}

export function StaticPage({ page, isLoading }: StaticPageProps) {
  if (isLoading) {
    return (
      <div className="mx-auto max-w-3xl space-y-4">
        <Skeleton className="h-10 w-3/4" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
      </div>
    );
  }

  return (
    <article className="mx-auto max-w-3xl">
      <h1 className="mb-6 text-3xl font-bold">{page.title}</h1>
      <div
        className="prose prose-neutral max-w-none dark:prose-invert"
        dangerouslySetInnerHTML={{ __html: page.content }}
      />
    </article>
  );
}
