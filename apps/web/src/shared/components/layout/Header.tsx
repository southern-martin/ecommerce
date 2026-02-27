import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
  Search,
  ShoppingCart,
  Menu,
  User,
  LogOut,
  Package,
  LayoutDashboard,
  Shield,
  Zap,
  Heart,
  Store,
} from "lucide-react";
import { Button } from "@/shared/components/ui/button";
import { Input } from "@/shared/components/ui/input";
import { Badge } from "@/shared/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/components/ui/dropdown-menu";
import { MobileNav } from "./MobileNav";
import { useAuth } from "@/shared/hooks/useAuth";
import { useCartStore } from "@/shared/stores/cart.store";

export function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const navigate = useNavigate();

  const cartItemCount = useCartStore((s) => s.itemCount());
  const { isAuthenticated, user, logout, isSeller, isAdmin } = useAuth();

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchQuery)}`);
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto flex h-16 items-center gap-4 px-4">
        {/* Mobile menu button */}
        <Button
          variant="ghost"
          size="icon"
          className="shrink-0 md:hidden"
          onClick={() => setMobileMenuOpen(true)}
        >
          <Menu className="h-5 w-5" />
          <span className="sr-only">Toggle menu</span>
        </Button>

        {/* Logo */}
        <Link to="/" className="flex shrink-0 items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <Store className="h-4 w-4" />
          </div>
          <span className="text-xl font-bold tracking-tight">
            <span className="text-primary">Market</span>
            <span className="text-foreground">Hub</span>
          </span>
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden items-center gap-1 md:flex">
          <Link
            to="/products"
            className="rounded-lg px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          >
            Shop
          </Link>
          <Link
            to="/categories"
            className="rounded-lg px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          >
            Categories
          </Link>
          <Link
            to="/promotions"
            className="flex items-center gap-1.5 rounded-lg px-3 py-2 text-sm font-medium text-orange-600 transition-colors hover:bg-orange-50"
          >
            <Zap className="h-4 w-4" />
            Deals
          </Link>
        </nav>

        {/* Search Bar */}
        <form
          onSubmit={handleSearch}
          className="hidden flex-1 md:flex"
        >
          <div className="relative w-full max-w-md mx-auto">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Search products, brands, categories..."
              className="w-full rounded-full border-muted-foreground/20 bg-muted/50 pl-10 pr-4 focus-visible:ring-primary"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
        </form>

        <div className="flex items-center gap-1 ml-auto">
          {/* Wishlist */}
          <Button variant="ghost" size="icon" asChild className="hidden sm:flex">
            <Link to="/account/wishlist">
              <Heart className="h-5 w-5" />
              <span className="sr-only">Wishlist</span>
            </Link>
          </Button>

          {/* Cart */}
          <Button variant="ghost" size="icon" asChild className="relative">
            <Link to="/cart">
              <ShoppingCart className="h-5 w-5" />
              {cartItemCount > 0 && (
                <Badge className="absolute -right-1 -top-1 h-5 min-w-5 rounded-full p-0 flex items-center justify-center text-[10px] font-bold">
                  {cartItemCount}
                </Badge>
              )}
              <span className="sr-only">Cart</span>
            </Link>
          </Button>

          {/* User Menu */}
          {isAuthenticated ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon">
                  <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 text-primary">
                    <User className="h-4 w-4" />
                  </div>
                  <span className="sr-only">User menu</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56">
                <div className="px-3 py-2">
                  <p className="text-sm font-medium">{user?.email}</p>
                  <p className="text-xs text-muted-foreground capitalize">{user?.role || 'buyer'}</p>
                </div>
                <DropdownMenuSeparator />
                <DropdownMenuItem asChild>
                  <Link to="/account/profile" className="flex items-center">
                    <User className="mr-2 h-4 w-4" />
                    Profile
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link to="/account/orders" className="flex items-center">
                    <Package className="mr-2 h-4 w-4" />
                    Orders
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link to="/account/wishlist" className="flex items-center">
                    <Heart className="mr-2 h-4 w-4" />
                    Wishlist
                  </Link>
                </DropdownMenuItem>
                {(isSeller || isAdmin) && (
                  <>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem asChild>
                      <Link to="/seller" className="flex items-center">
                        <LayoutDashboard className="mr-2 h-4 w-4" />
                        Seller Dashboard
                      </Link>
                    </DropdownMenuItem>
                  </>
                )}
                {isAdmin && (
                  <DropdownMenuItem asChild>
                    <Link to="/admin" className="flex items-center">
                      <Shield className="mr-2 h-4 w-4" />
                      Admin Panel
                    </Link>
                  </DropdownMenuItem>
                )}
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout} className="flex items-center cursor-pointer text-red-600 focus:text-red-600">
                  <LogOut className="mr-2 h-4 w-4" />
                  Sign Out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <div className="flex items-center gap-2 ml-2">
              <Button variant="ghost" size="sm" asChild>
                <Link to="/login">Sign In</Link>
              </Button>
              <Button size="sm" asChild className="hidden sm:flex">
                <Link to="/register">Sign Up</Link>
              </Button>
            </div>
          )}
        </div>
      </div>

      {/* Mobile Navigation Drawer */}
      <MobileNav open={mobileMenuOpen} onClose={() => setMobileMenuOpen(false)} />
    </header>
  );
}
