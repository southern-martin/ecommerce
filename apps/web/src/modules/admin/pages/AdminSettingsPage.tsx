import { useState } from 'react';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/shared/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';

export default function AdminSettingsPage() {
  const [platformName, setPlatformName] = useState('My eCommerce');
  const [currency, setCurrency] = useState('USD');
  const [maintenanceMode, setMaintenanceMode] = useState(false);

  const [orderNotifications, setOrderNotifications] = useState(true);
  const [userNotifications, setUserNotifications] = useState(true);
  const [sellerNotifications, setSellerNotifications] = useState(false);
  const [reviewNotifications, setReviewNotifications] = useState(true);

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold">Platform Settings</h1>

      <Tabs defaultValue="general">
        <TabsList>
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="notifications">Notifications</TabsTrigger>
        </TabsList>

        <TabsContent value="general">
          <Card>
            <CardHeader>
              <CardTitle>General Settings</CardTitle>
              <CardDescription>Configure basic platform settings.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6 max-w-md">
              <div className="space-y-2">
                <Label htmlFor="platform-name">Platform Name</Label>
                <Input
                  id="platform-name"
                  value={platformName}
                  onChange={(e) => setPlatformName(e.target.value)}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="currency">Default Currency</Label>
                <Select value={currency} onValueChange={setCurrency}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="USD">USD - US Dollar</SelectItem>
                    <SelectItem value="EUR">EUR - Euro</SelectItem>
                    <SelectItem value="GBP">GBP - British Pound</SelectItem>
                    <SelectItem value="CAD">CAD - Canadian Dollar</SelectItem>
                    <SelectItem value="AUD">AUD - Australian Dollar</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="flex items-center justify-between rounded-lg border p-4">
                <div>
                  <p className="text-sm font-medium">Maintenance Mode</p>
                  <p className="text-sm text-muted-foreground">
                    Temporarily disable public access to the storefront.
                  </p>
                </div>
                <button
                  type="button"
                  role="switch"
                  aria-checked={maintenanceMode}
                  onClick={() => setMaintenanceMode(!maintenanceMode)}
                  className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring ${
                    maintenanceMode ? 'bg-primary' : 'bg-input'
                  }`}
                >
                  <span
                    className={`pointer-events-none block h-5 w-5 rounded-full bg-background shadow-lg ring-0 transition-transform ${
                      maintenanceMode ? 'translate-x-5' : 'translate-x-0'
                    }`}
                  />
                </button>
              </div>

              <Button>Save Settings</Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="notifications">
          <Card>
            <CardHeader>
              <CardTitle>Email Notifications</CardTitle>
              <CardDescription>
                Configure which email notifications are sent to administrators.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4 max-w-md">
              {[
                {
                  id: 'order',
                  label: 'New Order Notifications',
                  description: 'Receive email when a new order is placed.',
                  checked: orderNotifications,
                  onChange: setOrderNotifications,
                },
                {
                  id: 'user',
                  label: 'New User Registrations',
                  description: 'Receive email when a new user registers.',
                  checked: userNotifications,
                  onChange: setUserNotifications,
                },
                {
                  id: 'seller',
                  label: 'Seller Applications',
                  description: 'Receive email when a seller applies for approval.',
                  checked: sellerNotifications,
                  onChange: setSellerNotifications,
                },
                {
                  id: 'review',
                  label: 'New Reviews',
                  description: 'Receive email when a new product review is submitted.',
                  checked: reviewNotifications,
                  onChange: setReviewNotifications,
                },
              ].map((item) => (
                <div
                  key={item.id}
                  className="flex items-center justify-between rounded-lg border p-4"
                >
                  <div>
                    <p className="text-sm font-medium">{item.label}</p>
                    <p className="text-sm text-muted-foreground">{item.description}</p>
                  </div>
                  <button
                    type="button"
                    role="switch"
                    aria-checked={item.checked}
                    onClick={() => item.onChange(!item.checked)}
                    className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring ${
                      item.checked ? 'bg-primary' : 'bg-input'
                    }`}
                  >
                    <span
                      className={`pointer-events-none block h-5 w-5 rounded-full bg-background shadow-lg ring-0 transition-transform ${
                        item.checked ? 'translate-x-5' : 'translate-x-0'
                      }`}
                    />
                  </button>
                </div>
              ))}

              <Button>Save Notification Preferences</Button>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
