import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { WalletBalanceCard } from '../components/WalletBalanceCard';
import { WalletTransactionTable } from '../components/WalletTransactionTable';
import { PayoutTable } from '../components/PayoutTable';
import { PayoutRequestForm } from '../components/PayoutRequestForm';
import {
  useSellerWalletBalance,
  useSellerWalletTransactions,
  useSellerPayouts,
  useRequestPayout,
} from '../hooks/useSellerWallet';

export default function SellerWalletPage() {
  const [txPage, setTxPage] = useState(1);
  const [payoutPage, setPayoutPage] = useState(1);

  const { data: wallet, isLoading: walletLoading } = useSellerWalletBalance();
  const { data: txData, isLoading: txLoading } = useSellerWalletTransactions(txPage);
  const { data: payoutData, isLoading: payoutsLoading } = useSellerPayouts(payoutPage);
  const requestPayout = useRequestPayout();

  const txTotalPages = txData ? Math.ceil(txData.total / txData.page_size) : 0;
  const payoutTotalPages = payoutData ? Math.ceil(payoutData.total / payoutData.page_size) : 0;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Wallet & Payouts</h1>

      {walletLoading ? (
        <Skeleton className="h-24" />
      ) : wallet ? (
        <WalletBalanceCard
          availableBalance={wallet.available_balance}
          pendingBalance={wallet.pending_balance}
          currency={wallet.currency}
        />
      ) : (
        <Card>
          <CardContent className="py-8 text-center text-muted-foreground">
            Wallet not available. Complete a sale to activate your wallet.
          </CardContent>
        </Card>
      )}

      <Tabs defaultValue="transactions">
        <TabsList>
          <TabsTrigger value="transactions">Transactions</TabsTrigger>
          <TabsTrigger value="payouts">Payouts</TabsTrigger>
          <TabsTrigger value="request">Request Payout</TabsTrigger>
        </TabsList>

        <TabsContent value="transactions">
          {txLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : txData && txData.transactions.length > 0 ? (
            <>
              <WalletTransactionTable transactions={txData.transactions} />
              {txTotalPages > 1 && (
                <div className="mt-6 flex items-center justify-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={txPage === 1}
                    onClick={() => setTxPage((p) => p - 1)}
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <span className="text-sm text-muted-foreground">
                    Page {txPage} of {txTotalPages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={txPage === txTotalPages}
                    onClick={() => setTxPage((p) => p + 1)}
                  >
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </div>
              )}
            </>
          ) : (
            <p className="py-8 text-center text-muted-foreground">No transactions yet.</p>
          )}
        </TabsContent>

        <TabsContent value="payouts">
          {payoutsLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : payoutData && payoutData.payouts.length > 0 ? (
            <>
              <PayoutTable payouts={payoutData.payouts} />
              {payoutTotalPages > 1 && (
                <div className="mt-6 flex items-center justify-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={payoutPage === 1}
                    onClick={() => setPayoutPage((p) => p - 1)}
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <span className="text-sm text-muted-foreground">
                    Page {payoutPage} of {payoutTotalPages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={payoutPage === payoutTotalPages}
                    onClick={() => setPayoutPage((p) => p + 1)}
                  >
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </div>
              )}
            </>
          ) : (
            <p className="py-8 text-center text-muted-foreground">No payouts yet.</p>
          )}
        </TabsContent>

        <TabsContent value="request">
          <Card>
            <CardHeader>
              <CardTitle>Request a Payout</CardTitle>
            </CardHeader>
            <CardContent>
              <PayoutRequestForm
                availableBalance={wallet?.available_balance ?? 0}
                onSubmit={(data) => requestPayout.mutate(data)}
                isPending={requestPayout.isPending}
              />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
