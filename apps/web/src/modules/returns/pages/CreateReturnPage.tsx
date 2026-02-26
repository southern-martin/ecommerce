import { useNavigate, useSearchParams } from 'react-router-dom';
import { ReturnRequestForm } from '../components/ReturnRequestForm';
import { useCreateReturn } from '../hooks/useReturns';

export default function CreateReturnPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const orderId = searchParams.get('order_id') ?? undefined;
  const createReturn = useCreateReturn();

  return (
    <div className="mx-auto max-w-2xl">
      <h1 className="mb-6 text-2xl font-bold">Request a Return</h1>
      <ReturnRequestForm
        orderId={orderId}
        onSubmit={(data) =>
          createReturn.mutate(
            { ...data, items: [] },
            { onSuccess: () => navigate('/account/returns') }
          )
        }
        isPending={createReturn.isPending}
      />
    </div>
  );
}
