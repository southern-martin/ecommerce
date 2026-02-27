import { User } from 'lucide-react';
import { PageLayout } from '@/shared/components/layout/PageLayout';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { ProfileForm } from '../components/ProfileForm';
import { useProfile } from '../hooks/useProfile';

export default function ProfilePage() {
  const { data: user, isLoading, updateProfile } = useProfile();

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!user) return null;

  return (
    <PageLayout
      title="My Profile"
      icon={User}
      breadcrumbs={[
        { label: 'Account', href: '/account/profile' },
        { label: 'Profile' },
      ]}
    >
      <ProfileForm
        user={user}
        onSubmit={(data) => updateProfile.mutate(data)}
        isPending={updateProfile.isPending}
      />
    </PageLayout>
  );
}
