import { Link } from "react-router-dom";
import { Separator } from "@/shared/components/ui/separator";

export function Footer() {
  return (
    <footer className="border-t bg-background">
      <div className="container mx-auto px-4 py-10">
        <div className="grid grid-cols-2 gap-8 md:grid-cols-4">
          <div>
            <h3 className="text-sm font-semibold">Shop</h3>
            <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
              <li>
                <Link to="/shop" className="hover:text-foreground">
                  All Products
                </Link>
              </li>
              <li>
                <Link to="/flash-sales" className="hover:text-foreground">
                  Flash Sales
                </Link>
              </li>
              <li>
                <Link to="/shop?category=new" className="hover:text-foreground">
                  New Arrivals
                </Link>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold">Account</h3>
            <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
              <li>
                <Link to="/account/profile" className="hover:text-foreground">
                  Profile
                </Link>
              </li>
              <li>
                <Link to="/account/orders" className="hover:text-foreground">
                  Orders
                </Link>
              </li>
              <li>
                <Link to="/account/wishlist" className="hover:text-foreground">
                  Wishlist
                </Link>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold">Seller</h3>
            <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
              <li>
                <Link to="/seller/dashboard" className="hover:text-foreground">
                  Seller Dashboard
                </Link>
              </li>
              <li>
                <Link to="/seller/products" className="hover:text-foreground">
                  Manage Products
                </Link>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold">Support</h3>
            <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
              <li>
                <Link to="/help" className="hover:text-foreground">
                  Help Center
                </Link>
              </li>
              <li>
                <Link to="/returns" className="hover:text-foreground">
                  Returns
                </Link>
              </li>
              <li>
                <Link to="/shipping" className="hover:text-foreground">
                  Shipping Info
                </Link>
              </li>
            </ul>
          </div>
        </div>
        <Separator className="my-8" />
        <div className="flex flex-col items-center justify-between gap-4 md:flex-row">
          <p className="text-sm text-muted-foreground">
            &copy; {new Date().getFullYear()} Store. All rights reserved.
          </p>
          <div className="flex space-x-4 text-sm text-muted-foreground">
            <Link to="/privacy" className="hover:text-foreground">
              Privacy Policy
            </Link>
            <Link to="/terms" className="hover:text-foreground">
              Terms of Service
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
