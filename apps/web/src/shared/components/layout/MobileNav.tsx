import { Link } from "react-router-dom";
import { X, Search, Zap, ShoppingCart, User, LogOut } from "lucide-react";
import { Button } from "@/shared/components/ui/button";
import { Input } from "@/shared/components/ui/input";
import { Separator } from "@/shared/components/ui/separator";
import { cn } from "@/shared/lib/utils";

interface MobileNavProps {
  open: boolean;
  onClose: () => void;
}

export function MobileNav({ open, onClose }: MobileNavProps) {
  return (
    <>
      {/* Overlay */}
      {open && (
        <div
          className="fixed inset-0 z-50 bg-black/50"
          onClick={onClose}
        />
      )}

      {/* Drawer */}
      <div
        className={cn(
          "fixed inset-y-0 left-0 z-50 w-72 bg-background shadow-lg transform transition-transform duration-300 ease-in-out",
          open ? "translate-x-0" : "-translate-x-full"
        )}
      >
        <div className="flex items-center justify-between p-4">
          <span className="text-lg font-bold">Store</span>
          <Button variant="ghost" size="icon" onClick={onClose}>
            <X className="h-5 w-5" />
            <span className="sr-only">Close menu</span>
          </Button>
        </div>

        <div className="px-4 pb-4">
          <div className="relative">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Search products..."
              className="pl-8"
            />
          </div>
        </div>

        <Separator />

        <nav className="flex flex-col p-4 space-y-1">
          <Link
            to="/shop"
            onClick={onClose}
            className="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium hover:bg-accent"
          >
            <ShoppingCart className="h-4 w-4" />
            Shop
          </Link>
          <Link
            to="/flash-sales"
            onClick={onClose}
            className="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium hover:bg-accent"
          >
            <Zap className="h-4 w-4" />
            Flash Sales
          </Link>

          <Separator className="my-2" />

          <Link
            to="/cart"
            onClick={onClose}
            className="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium hover:bg-accent"
          >
            <ShoppingCart className="h-4 w-4" />
            Cart
          </Link>
          <Link
            to="/account/profile"
            onClick={onClose}
            className="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium hover:bg-accent"
          >
            <User className="h-4 w-4" />
            Profile
          </Link>

          <Separator className="my-2" />

          <button className="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium hover:bg-accent text-left">
            <LogOut className="h-4 w-4" />
            Logout
          </button>
        </nav>
      </div>
    </>
  );
}
