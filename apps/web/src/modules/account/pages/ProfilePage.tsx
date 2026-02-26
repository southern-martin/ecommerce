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
    <div>
      <h1 className="mb-6 text-2xl font-bold">My Profile</h1>
      <ProfileForm
        user={user}
        onSubmit={(data) => updateProfile.mutate(data)}
        isPending={updateProfile.isPending}
      />
    </div>
  );
}
