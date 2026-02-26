import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import apiClient from '@/shared/lib/api-client';

export default function StaticPage() {
  const { slug } = useParams<{ slug: string }>();

  const { data: page, isLoading } = useQuery({
    queryKey: ['cms-page', slug],
    queryFn: async () => {
      const res = await apiClient.get(`/cms/pages/${slug}`);
      return res.data;
    },
    enabled: !!slug,
  });

  if (isLoading) return <div className="p-6">Loading...</div>;

  if (!page) return <div className="p-6 text-center text-muted-foreground">Page not found.</div>;

  return (
    <div className="prose mx-auto max-w-3xl py-8">
      <h1>{(page as any).title}</h1>
      <div dangerouslySetInnerHTML={{ __html: (page as any).content || '' }} />
    </div>
  );
}
