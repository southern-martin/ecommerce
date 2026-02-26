import { useParams } from 'react-router-dom';
import { StaticPage } from '../components/StaticPage';
import { usePage } from '../hooks/usePages';

export default function StaticPageView() {
  const { slug } = useParams<{ slug: string }>();
  const { data: page, isLoading } = usePage(slug!);

  if (!page && !isLoading) {
    return (
      <div className="py-16 text-center">
        <p className="text-lg text-muted-foreground">Page not found.</p>
      </div>
    );
  }

  return page ? <StaticPage page={page} isLoading={isLoading} /> : null;
}
