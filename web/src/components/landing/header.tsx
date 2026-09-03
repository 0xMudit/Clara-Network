// src/components/landing/header.tsx
import Link from "next/link";
import { Landmark, Star } from "lucide-react";
import { ThemeToggle } from "@/components/theme-toggle";
import { GithubIcon } from "./github-icon";

const GITHUB_URL = "https://github.com/0xMudit/Clara-Network";

export function LandingHeader() {
  return (
    <header className="absolute inset-x-0 top-0 z-30">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-2.5 transition-opacity hover:opacity-80">
          <div className="flex size-9 items-center justify-center rounded-xl bg-white/10 ring-1 ring-white/20 backdrop-blur-sm">
            <Landmark className="size-4.5" />
          </div>
          <span className="text-base font-semibold tracking-tight">
            Clara Network
          </span>
        </Link>

        {/* Right: theme toggle + GitHub */}
        <div className="flex items-center gap-2">
          <ThemeToggle />
          <a
            href={GITHUB_URL}
            target="_blank"
            rel="noopener noreferrer"
            className="group inline-flex items-center gap-2 rounded-full bg-white/10 px-3.5 py-1.5 text-xs font-medium ring-1 ring-white/15 backdrop-blur-sm transition-all hover:bg-white/15 hover:ring-white/25"
          >
            <GithubIcon className="size-3.5" />
            <span className="hidden sm:inline">Star on GitHub</span>
            <span className="inline sm:hidden">GitHub</span>
            <span className="inline-flex items-center gap-1 text-white/60 transition-colors group-hover:text-white/90">
              <Star className="size-3.5 fill-current" />
            </span>
          </a>
        </div>
      </div>
    </header>
  );
}