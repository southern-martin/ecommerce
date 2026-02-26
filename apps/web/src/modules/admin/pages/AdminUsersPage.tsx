import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { UserManagementTable } from '../components/UserManagementTable';
import { useAdminUsers, useDeleteUser } from '../hooks/useAdminUsers';

export default function AdminUsersPage() {
  const [page, setPage] = useState(1);
  const { data, isLoading } = useAdminUsers(page);
  const deleteUser = useDeleteUser();
  const totalPages = data ? Math.ceil(data.total / data.page_size) : 0;

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold">User Management</h1>

      {data && data.data.length > 0 ? (
        <>
          <UserManagementTable
            users={data.data}
            onDelete={(id) => deleteUser.mutate(id)}
          />
          {totalPages > 1 && (
            <div className="mt-6 flex items-center justify-center gap-2">
              <Button variant="outline" size="sm" disabled={page === 1} onClick={() => setPage((p) => p - 1)}>
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <span className="text-sm text-muted-foreground">Page {page} of {totalPages}</span>
              <Button variant="outline" size="sm" disabled={page === totalPages} onClick={() => setPage((p) => p + 1)}>
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          )}
        </>
      ) : (
        <p className="py-8 text-center text-muted-foreground">No users found.</p>
      )}
    </div>
  );
}
