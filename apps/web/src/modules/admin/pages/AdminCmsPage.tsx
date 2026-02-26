import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/shared/components/ui/dialog';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { Plus, Pencil, Trash2, Globe, Loader2 } from 'lucide-react';
import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { ConfirmDialog } from '@/shared/components/data/ConfirmDialog';
import {
  useAdminBanners,
  useCreateBanner,
  useUpdateBanner,
  useDeleteBanner,
  useAdminPages,
  useCreatePage,
  useUpdatePage,
  useDeletePage,
  usePublishPage,
} from '../hooks/useAdminCms';

// Banner form schema
const bannerSchema = z.object({
  title: z.string().min(1, 'Title is required'),
  subtitle: z.string().optional(),
  image_url: z.string().url('Must be a valid URL'),
  link_url: z.string().optional(),
  is_active: z.boolean().default(true),
  sort_order: z.coerce.number().default(0),
});

type BannerFormValues = z.infer<typeof bannerSchema>;

// Page form schema
const pageSchema = z.object({
  title: z.string().min(1, 'Title is required'),
  slug: z.string().min(1, 'Slug is required'),
  content: z.string().min(1, 'Content is required'),
  meta_title: z.string().optional(),
  meta_description: z.string().optional(),
});

type PageFormValues = z.infer<typeof pageSchema>;

function BannerForm({
  defaultValues,
  onSubmit,
  isPending,
  submitLabel = 'Create Banner',
}: {
  defaultValues?: Partial<BannerFormValues>;
  onSubmit: (data: BannerFormValues) => void;
  isPending?: boolean;
  submitLabel?: string;
}) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<BannerFormValues>({
    resolver: zodResolver(bannerSchema),
    defaultValues: { title: '', subtitle: '', image_url: '', link_url: '', is_active: true, sort_order: 0, ...defaultValues },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="banner-title">Title</Label>
        <Input id="banner-title" {...register('title')} />
        {errors.title && <p className="text-sm text-destructive">{errors.title.message}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="banner-subtitle">Subtitle</Label>
        <Input id="banner-subtitle" {...register('subtitle')} />
      </div>
      <div className="space-y-2">
        <Label htmlFor="banner-image">Image URL</Label>
        <Input id="banner-image" {...register('image_url')} placeholder="https://..." />
        {errors.image_url && <p className="text-sm text-destructive">{errors.image_url.message}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="banner-link">Link URL</Label>
        <Input id="banner-link" {...register('link_url')} />
      </div>
      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="banner-sort">Sort Order</Label>
          <Input id="banner-sort" type="number" {...register('sort_order')} />
        </div>
        <div className="flex items-center gap-2 pt-6">
          <input type="checkbox" id="banner-active" {...register('is_active')} className="rounded" />
          <Label htmlFor="banner-active">Active</Label>
        </div>
      </div>
      <Button type="submit" disabled={isPending}>
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {submitLabel}
      </Button>
    </form>
  );
}

function PageForm({
  defaultValues,
  onSubmit,
  isPending,
  submitLabel = 'Create Page',
}: {
  defaultValues?: Partial<PageFormValues>;
  onSubmit: (data: PageFormValues) => void;
  isPending?: boolean;
  submitLabel?: string;
}) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<PageFormValues>({
    resolver: zodResolver(pageSchema),
    defaultValues: { title: '', slug: '', content: '', meta_title: '', meta_description: '', ...defaultValues },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="page-title">Title</Label>
        <Input id="page-title" {...register('title')} />
        {errors.title && <p className="text-sm text-destructive">{errors.title.message}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="page-slug">Slug</Label>
        <Input id="page-slug" {...register('slug')} />
        {errors.slug && <p className="text-sm text-destructive">{errors.slug.message}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="page-content">Content</Label>
        <textarea
          id="page-content"
          {...register('content')}
          rows={6}
          className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
        {errors.content && <p className="text-sm text-destructive">{errors.content.message}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="page-meta-title">Meta Title</Label>
        <Input id="page-meta-title" {...register('meta_title')} />
      </div>
      <div className="space-y-2">
        <Label htmlFor="page-meta-desc">Meta Description</Label>
        <Input id="page-meta-desc" {...register('meta_description')} />
      </div>
      <Button type="submit" disabled={isPending}>
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {submitLabel}
      </Button>
    </form>
  );
}

export default function AdminCmsPage() {
  const [bannerDialogOpen, setBannerDialogOpen] = useState(false);
  const [editingBanner, setEditingBanner] = useState<any>(null);
  const [pageDialogOpen, setPageDialogOpen] = useState(false);
  const [editingPage, setEditingPage] = useState<any>(null);

  const { data: banners, isLoading: bannersLoading } = useAdminBanners();
  const createBanner = useCreateBanner();
  const updateBanner = useUpdateBanner();
  const deleteBanner = useDeleteBanner();

  const { data: pages, isLoading: pagesLoading } = useAdminPages();
  const createPage = useCreatePage();
  const updatePage = useUpdatePage();
  const deletePage = useDeletePage();
  const publishPage = usePublishPage();

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Content Management</h1>

      <Tabs defaultValue="banners">
        <TabsList>
          <TabsTrigger value="banners">Banners</TabsTrigger>
          <TabsTrigger value="pages">Pages</TabsTrigger>
        </TabsList>

        {/* Banners Tab */}
        <TabsContent value="banners" className="space-y-4">
          <div className="flex justify-end">
            <Dialog open={bannerDialogOpen} onOpenChange={setBannerDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Banner
                </Button>
              </DialogTrigger>
              <DialogContent className="max-w-lg">
                <DialogHeader>
                  <DialogTitle>Create Banner</DialogTitle>
                </DialogHeader>
                <BannerForm
                  onSubmit={(data) =>
                    createBanner.mutate(data, {
                      onSuccess: () => setBannerDialogOpen(false),
                    })
                  }
                  isPending={createBanner.isPending}
                />
              </DialogContent>
            </Dialog>
          </div>

          {bannersLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Image</TableHead>
                  <TableHead>Title</TableHead>
                  <TableHead>Link</TableHead>
                  <TableHead>Order</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[100px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(!banners || banners.length === 0) ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-muted-foreground">
                      No banners found.
                    </TableCell>
                  </TableRow>
                ) : (
                  banners.map((banner) => (
                    <TableRow key={banner.id}>
                      <TableCell>
                        <img
                          src={banner.image_url}
                          alt={banner.title}
                          className="h-10 w-16 rounded object-cover"
                        />
                      </TableCell>
                      <TableCell className="font-medium">{banner.title}</TableCell>
                      <TableCell className="max-w-[150px] truncate text-muted-foreground">
                        {banner.link_url || '-'}
                      </TableCell>
                      <TableCell>{banner.sort_order}</TableCell>
                      <TableCell>
                        <StatusBadge status={banner.is_active ? 'active' : 'inactive'} />
                      </TableCell>
                      <TableCell>
                        <div className="flex gap-1">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => setEditingBanner(banner)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <ConfirmDialog
                            title="Delete Banner"
                            description={`Delete banner "${banner.title}"? This cannot be undone.`}
                            onConfirm={() => deleteBanner.mutate(banner.id)}
                            isPending={deleteBanner.isPending}
                            trigger={
                              <Button variant="ghost" size="sm">
                                <Trash2 className="h-4 w-4 text-destructive" />
                              </Button>
                            }
                          />
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          )}

          <Dialog open={!!editingBanner} onOpenChange={(open) => !open && setEditingBanner(null)}>
            <DialogContent className="max-w-lg">
              <DialogHeader>
                <DialogTitle>Edit Banner</DialogTitle>
              </DialogHeader>
              {editingBanner && (
                <BannerForm
                  defaultValues={editingBanner}
                  onSubmit={(data) =>
                    updateBanner.mutate(
                      { id: editingBanner.id, data },
                      { onSuccess: () => setEditingBanner(null) }
                    )
                  }
                  isPending={updateBanner.isPending}
                  submitLabel="Update Banner"
                />
              )}
            </DialogContent>
          </Dialog>
        </TabsContent>

        {/* Pages Tab */}
        <TabsContent value="pages" className="space-y-4">
          <div className="flex justify-end">
            <Dialog open={pageDialogOpen} onOpenChange={setPageDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Page
                </Button>
              </DialogTrigger>
              <DialogContent className="max-w-lg">
                <DialogHeader>
                  <DialogTitle>Create Page</DialogTitle>
                </DialogHeader>
                <PageForm
                  onSubmit={(data) =>
                    createPage.mutate(data, {
                      onSuccess: () => setPageDialogOpen(false),
                    })
                  }
                  isPending={createPage.isPending}
                />
              </DialogContent>
            </Dialog>
          </div>

          {pagesLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Title</TableHead>
                  <TableHead>Slug</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[140px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(!pages || pages.length === 0) ? (
                  <TableRow>
                    <TableCell colSpan={4} className="text-center text-muted-foreground">
                      No pages found.
                    </TableCell>
                  </TableRow>
                ) : (
                  pages.map((page) => (
                    <TableRow key={page.id}>
                      <TableCell className="font-medium">{page.title}</TableCell>
                      <TableCell className="text-muted-foreground">/{page.slug}</TableCell>
                      <TableCell>
                        <StatusBadge status={page.published ? 'published' : 'draft'} />
                      </TableCell>
                      <TableCell>
                        <div className="flex gap-1">
                          {!page.published && (
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => publishPage.mutate(page.id)}
                              disabled={publishPage.isPending}
                            >
                              <Globe className="h-4 w-4" />
                            </Button>
                          )}
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => setEditingPage(page)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <ConfirmDialog
                            title="Delete Page"
                            description={`Delete page "${page.title}"? This cannot be undone.`}
                            onConfirm={() => deletePage.mutate(page.id)}
                            isPending={deletePage.isPending}
                            trigger={
                              <Button variant="ghost" size="sm">
                                <Trash2 className="h-4 w-4 text-destructive" />
                              </Button>
                            }
                          />
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          )}

          <Dialog open={!!editingPage} onOpenChange={(open) => !open && setEditingPage(null)}>
            <DialogContent className="max-w-lg">
              <DialogHeader>
                <DialogTitle>Edit Page</DialogTitle>
              </DialogHeader>
              {editingPage && (
                <PageForm
                  defaultValues={editingPage}
                  onSubmit={(data) =>
                    updatePage.mutate(
                      { id: editingPage.id, data },
                      { onSuccess: () => setEditingPage(null) }
                    )
                  }
                  isPending={updatePage.isPending}
                  submitLabel="Update Page"
                />
              )}
            </DialogContent>
          </Dialog>
        </TabsContent>
      </Tabs>
    </div>
  );
}
