import { Link } from "react-router-dom";
import { Store, Mail, MapPin, Phone } from "lucide-react";

export function Footer() {
  return (
    <footer className="border-t bg-muted/30">
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-2 gap-8 md:grid-cols-5">
          {/* Brand */}
          <div className="col-span-2 md:col-span-1">
            <Link to="/" className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                <Store className="h-4 w-4" />
              </div>
              <span className="text-lg font-bold">
                <span className="text-primary">Market</span>Hub
              </span>
            </Link>
            <p className="mt-4 text-sm text-muted-foreground leading-relaxed">
              Your one-stop marketplace for quality products at the best prices.
            </p>
            <div className="mt-4 space-y-2 text-sm text-muted-foreground">
              <div className="flex items-center gap-2">
                <Mail className="h-4 w-4" />
                support@markethub.com
              </div>
              <div className="flex items-center gap-2">
                <Phone className="h-4 w-4" />
                1-800-MARKET
              </div>
              <div className="flex items-center gap-2">
                <MapPin className="h-4 w-4" />
                San Francisco, CA
              </div>
            </div>
          </div>

          {/* Shop */}
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider">Shop</h3>
            <ul className="mt-4 space-y-3 text-sm text-muted-foreground">
              <li>
                <Link to="/products" className="transition-colors hover:text-foreground">
                  All Products
                </Link>
              </li>
              <li>
                <Link to="/categories" className="transition-colors hover:text-foreground">
                  Categories
                </Link>
              </li>
              <li>
                <Link to="/promotions" className="transition-colors hover:text-foreground">
                  Flash Sales
                </Link>
              </li>
              <li>
                <Link to="/products?sort=newest" className="transition-colors hover:text-foreground">
                  New Arrivals
                </Link>
              </li>
            </ul>
          </div>

          {/* Account */}
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider">Account</h3>
            <ul className="mt-4 space-y-3 text-sm text-muted-foreground">
              <li>
                <Link to="/account/profile" className="transition-colors hover:text-foreground">
                  My Profile
                </Link>
              </li>
              <li>
                <Link to="/account/orders" className="transition-colors hover:text-foreground">
                  Order History
                </Link>
              </li>
              <li>
                <Link to="/account/wishlist" className="transition-colors hover:text-foreground">
                  Wishlist
                </Link>
              </li>
              <li>
                <Link to="/cart" className="transition-colors hover:text-foreground">
                  Shopping Cart
                </Link>
              </li>
            </ul>
          </div>

          {/* Sell */}
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider">Sell</h3>
            <ul className="mt-4 space-y-3 text-sm text-muted-foreground">
              <li>
                <Link to="/seller/dashboard" className="transition-colors hover:text-foreground">
                  Seller Center
                </Link>
              </li>
              <li>
                <Link to="/seller/products" className="transition-colors hover:text-foreground">
                  Manage Products
                </Link>
              </li>
              <li>
                <Link to="/help" className="transition-colors hover:text-foreground">
                  Seller Help
                </Link>
              </li>
            </ul>
          </div>

          {/* Support */}
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider">Support</h3>
            <ul className="mt-4 space-y-3 text-sm text-muted-foreground">
              <li>
                <Link to="/help" className="transition-colors hover:text-foreground">
                  Help Center
                </Link>
              </li>
              <li>
                <Link to="/returns" className="transition-colors hover:text-foreground">
                  Returns & Refunds
                </Link>
              </li>
              <li>
                <Link to="/shipping" className="transition-colors hover:text-foreground">
                  Shipping Info
                </Link>
              </li>
              <li>
                <Link to="/privacy" className="transition-colors hover:text-foreground">
                  Privacy Policy
                </Link>
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom bar */}
        <div className="mt-12 flex flex-col items-center justify-between gap-4 border-t pt-8 md:flex-row">
          <p className="text-sm text-muted-foreground">
            &copy; {new Date().getFullYear()} MarketHub. All rights reserved.
          </p>
          <div className="flex items-center gap-6 text-sm text-muted-foreground">
            <Link to="/privacy" className="transition-colors hover:text-foreground">
              Privacy
            </Link>
            <Link to="/terms" className="transition-colors hover:text-foreground">
              Terms
            </Link>
            <Link to="/help" className="transition-colors hover:text-foreground">
              Support
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
