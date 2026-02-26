import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { ChevronLeft, ChevronRight, Check, X } from 'lucide-react';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';
import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { formatPrice, formatDate } from '@/shared/lib/utils';
import {
  useAffiliateProgram,
  useUpdateAffiliateProgram,
  useAffiliatePayouts,
  useUpdatePayoutStatus,
} from '../hooks/useAdminAffiliate';

const programSchema = z.object({
  commission_rate: z.coerce.number().min(0).max(100),
  cookie_duration_days: z.coerce.number().min(1),
  min_payout_amount: z.coerce.number().min(0),
  payout_schedule: z.string().min(1, 'Payout schedule is required'),
});

type ProgramFormValues = z.infer<typeof programSchema>;

export default function AdminAffiliatePage() {
  const [page, setPage] = useState(1);

  const { data: program, isLoading: programLoading } = useAffiliateProgram();
  const updateProgram = useUpdateAffiliateProgram();
  const { data: payoutsData, isLoading: payoutsLoading } = useAffiliatePayouts(page);
  const updatePayoutStatus = useUpdatePayoutStatus();

  const payouts = payoutsData?.data ?? [];
  const totalPages = payoutsData ? Math.ceil(payoutsData.total / payoutsData.page_size) : 0;

  const form = useForm<ProgramFormValues>({
    resolver: zodResolver(programSchema),
    values: program
      ? {
          commission_rate: program.commission_rate,
          cookie_duration_days: program.cookie_duration_days,
          min_payout_amount: program.min_payout_amount,
          payout_schedule: program.payout_schedule,
        }
      : undefined,
  });

  const onSaveProgram = (values: ProgramFormValues) => {
    updateProgram.mutate(values);
  };

  if (programLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold">Affiliate Management</h1>

      <Tabs defaultValue="settings">
        <TabsList>
          <TabsTrigger value="settings">Program Settings</TabsTrigger>
          <TabsTrigger value="payouts">Payouts</TabsTrigger>
        </TabsList>

        <TabsContent value="settings">
          <Card>
            <CardHeader>
              <CardTitle>Program Settings</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={form.handleSubmit(onSaveProgram)} className="space-y-4 max-w-md">
                <div className="space-y-2">
                  <Label htmlFor="commission_rate">Commission Rate (%)</Label>
                  <Input
                    id="commission_rate"
                    type="number"
                    step="0.1"
                    {...form.register('commission_rate')}
                  />
                  {form.formState.errors.commission_rate && (
                    <p className="text-sm text-destructive">
                      {form.formState.errors.commission_rate.message}
                    </p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="cookie_duration_days">Cookie Duration (days)</Label>
                  <Input
                    id="cookie_duration_days"
                    type="number"
                    {...form.register('cookie_duration_days')}
                  />
                  {form.formState.errors.cookie_duration_days && (
                    <p className="text-sm text-destructive">
                      {form.formState.errors.cookie_duration_days.message}
                    </p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="min_payout_amount">Minimum Payout Amount</Label>
                  <Input
                    id="min_payout_amount"
                    type="number"
                    step="0.01"
                    {...form.register('min_payout_amount')}
                  />
                  {form.formState.errors.min_payout_amount && (
                    <p className="text-sm text-destructive">
                      {form.formState.errors.min_payout_amount.message}
                    </p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="payout_schedule">Payout Schedule</Label>
                  <Select
                    value={form.watch('payout_schedule')}
                    onValueChange={(val) => form.setValue('payout_schedule', val)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select schedule" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="weekly">Weekly</SelectItem>
                      <SelectItem value="biweekly">Biweekly</SelectItem>
                      <SelectItem value="monthly">Monthly</SelectItem>
                    </SelectContent>
                  </Select>
                  {form.formState.errors.payout_schedule && (
                    <p className="text-sm text-destructive">
                      {form.formState.errors.payout_schedule.message}
                    </p>
                  )}
                </div>
                <Button type="submit" disabled={updateProgram.isPending}>
                  Save Settings
                </Button>
              </form>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="payouts">
          {payoutsLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : payouts.length > 0 ? (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>User Email</TableHead>
                    <TableHead>Amount</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Requested At</TableHead>
                    <TableHead className="w-32">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {payouts.map((payout: { id: string; user_email: string; amount: number; status: string; requested_at: string }) => (
                    <TableRow key={payout.id}>
                      <TableCell>{payout.user_email}</TableCell>
                      <TableCell>{formatPrice(payout.amount)}</TableCell>
                      <TableCell>
                        <StatusBadge status={payout.status} />
                      </TableCell>
                      <TableCell>{formatDate(payout.requested_at)}</TableCell>
                      <TableCell>
                        {payout.status === 'pending' && (
                          <div className="flex gap-1">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() =>
                                updatePayoutStatus.mutate({ id: payout.id, status: 'approved' })
                              }
                              disabled={updatePayoutStatus.isPending}
                            >
                              <Check className="h-4 w-4 text-green-600" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() =>
                                updatePayoutStatus.mutate({ id: payout.id, status: 'rejected' })
                              }
                              disabled={updatePayoutStatus.isPending}
                            >
                              <X className="h-4 w-4 text-red-600" />
                            </Button>
                          </div>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
              {totalPages > 1 && (
                <div className="mt-6 flex items-center justify-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={page === 1}
                    onClick={() => setPage((p) => p - 1)}
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <span className="text-sm text-muted-foreground">
                    Page {page} of {totalPages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={page === totalPages}
                    onClick={() => setPage((p) => p + 1)}
                  >
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </div>
              )}
            </>
          ) : (
            <p className="py-8 text-center text-muted-foreground">No payouts found.</p>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}
