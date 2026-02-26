import { Outlet } from "react-router-dom";
import { Sidebar, type SidebarNavItem } from "./Sidebar";

interface DashboardLayoutProps {
  sidebarItems: SidebarNavItem[];
  sidebarTitle?: string;
}

export function DashboardLayout({
  sidebarItems,
  sidebarTitle,
}: DashboardLayoutProps) {
  return (
    <div className="flex h-[calc(100vh-4rem)]">
      <Sidebar items={sidebarItems} title={sidebarTitle} />
      <main className="flex-1 overflow-auto p-6">
        <Outlet />
      </main>
    </div>
  );
}
