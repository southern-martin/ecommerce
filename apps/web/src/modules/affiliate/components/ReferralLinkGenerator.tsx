import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Copy, Check, Link as LinkIcon } from 'lucide-react';
import { useGenerateLink } from '../hooks/useAffiliateStats';

export function ReferralLinkGenerator() {
  const [copied, setCopied] = useState(false);
  const generateLink = useGenerateLink();

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-base">
          <LinkIcon className="h-5 w-5" />
          Referral Link Generator
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <Button
          onClick={() => generateLink.mutate(undefined)}
          disabled={generateLink.isPending}
        >
          Generate New Link
        </Button>

        {generateLink.data && (
          <div className="space-y-2">
            <div className="flex gap-2">
              <Input value={generateLink.data.link} readOnly className="bg-muted" />
              <Button
                variant="outline"
                size="icon"
                onClick={() => handleCopy(generateLink.data!.link)}
              >
                {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
              </Button>
            </div>
            <p className="text-xs text-muted-foreground">
              Referral code: <span className="font-mono font-medium">{generateLink.data.code}</span>
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
