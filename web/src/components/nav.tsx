import { HOME_BY_ROLE, DASHBOARD_ACCESS, ROLE_LABEL, type Role } from "@/lib/roles";
import { LogoutButton } from "@/components/logout-button";
import { ThemeToggle } from "@/components/theme-toggle";
import Link from "next/link";

const TITLES: Record<string, string> = {
  "/overview": "Overview", "/ops": "Operations", "/issuer": "Issuer", "/acquirer": "Acquirer",
  "/merchant": "Merchant", "/transactions": "Transactions", "/clearing": "Clearing",
  "/settlement": "Settlement", "/ledger": "Ledger", "/cards": "Cards", "/tokens": "Tokens",
  "/merchants": "Merchants", "/funding": "Funding", "/disputes": "Disputes",
};

export default function Nav({ role }: { role: Role }) {
  const items = DASHBOARD_ACCESS[role] ?? [];
  return (
    <header className="sticky top-0 z-40 border-b bg-background/95 backdrop-blur">
      <div className="flex h-14 items-center justify-between px-4">
        <div className="flex items-center gap-6">
          <Link href={HOME_BY_ROLE[role]} className="font-semibold tracking-tight">Clara Network</Link>
          <nav className="hidden items-center gap-1 md:flex">
            {items.map(p => (
              <Link key={p} href={p} className="rounded-md px-3 py-1.5 text-sm text-muted-foreground hover:bg-muted hover:text-foreground">
                {TITLES[p] ?? p}
              </Link>
            ))}
          </nav>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-sm text-muted-foreground">{ROLE_LABEL[role]}</span>
          <ThemeToggle />
          <LogoutButton />
        </div>
      </div>
    </header>
  );
}