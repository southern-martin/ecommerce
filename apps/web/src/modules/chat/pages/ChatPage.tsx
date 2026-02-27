import { useState } from 'react';
import { MessageCircle, Send } from 'lucide-react';
import { PageLayout } from '@/shared/components/layout/PageLayout';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Card } from '@/shared/components/ui/card';

export default function ChatPage() {
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState<{ text: string; sender: 'user' | 'system' }[]>([
    { text: 'Welcome to customer support! How can we help you?', sender: 'system' },
  ]);

  const handleSend = () => {
    if (!message.trim()) return;
    setMessages((prev) => [...prev, { text: message, sender: 'user' }]);
    setMessage('');
    // Simulate response
    setTimeout(() => {
      setMessages((prev) => [...prev, { text: 'Thank you for your message. A support agent will respond shortly.', sender: 'system' }]);
    }, 1000);
  };

  return (
    <PageLayout
      title="Chat Support"
      icon={MessageCircle}
      breadcrumbs={[
        { label: 'Account', href: '/account/profile' },
        { label: 'Chat' },
      ]}
    >
      <Card className="flex flex-col overflow-hidden h-[calc(100vh-16rem)]">
        <div className="flex-1 overflow-y-auto p-4 space-y-3">
          {messages.map((msg, i) => (
            <div key={i} className={`flex ${msg.sender === 'user' ? 'justify-end' : 'justify-start'}`}>
              <div className={`max-w-[70%] rounded-lg px-4 py-2 text-sm ${msg.sender === 'user' ? 'bg-primary text-primary-foreground' : 'bg-muted'}`}>
                {msg.text}
              </div>
            </div>
          ))}
        </div>
        <div className="border-t p-4 flex gap-2">
          <Input
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSend()}
            placeholder="Type a message..."
          />
          <Button size="icon" onClick={handleSend}><Send className="h-4 w-4" /></Button>
        </div>
      </Card>
    </PageLayout>
  );
}
