BASE 872db5a

5bd97ac feat(web): add shadcn ui with stripe theme

=== STAT ===
 web/components.json                             |   25 +
 web/package-lock.json                           | 3316 +++++++++++++++++++++--
 web/package.json                                |   10 +-
 web/src/app/globals.css                         |  131 +-
 web/src/app/layout.tsx                          |    6 +-
 web/src/components/providers/theme-provider.tsx |   19 +
 web/src/components/theme-provider.tsx           |    1 +
 web/src/components/theme-switcher.tsx           |  158 ++
 web/src/components/theme-toggle.tsx             |   31 +
 web/src/components/ui/button.tsx                |   58 +
 web/src/components/ui/dropdown-menu.tsx         |  268 ++
 web/src/components/ui/scroll-area.tsx           |   55 +
 web/src/lib/themes-config.ts                    |  331 +++
 web/src/lib/utils.ts                            |    6 +
 web/src/styles/themes/amber-minimal.css         |  124 +
 web/src/styles/themes/bold-tech.css             |  134 +
 web/src/styles/themes/bubblegum.css             |  134 +
 web/src/styles/themes/caffeine.css              |  132 +
 web/src/styles/themes/candyland.css             |  122 +
 web/src/styles/themes/catppuccin.css            |  134 +
 web/src/styles/themes/claude.css                |  132 +
 web/src/styles/themes/claymorphism.css          |  129 +
 web/src/styles/themes/clean-slate.css           |  124 +
 web/src/styles/themes/cosmic-night.css          |  134 +
 web/src/styles/themes/cyberpunk.css             |  124 +
 web/src/styles/themes/doom-64.css               |  124 +
 web/src/styles/themes/elegant-luxury.css        |  134 +
 web/src/styles/themes/github.css                |  114 +
 web/src/styles/themes/graphite.css              |  124 +
 web/src/styles/themes/index.css                 |   43 +
 web/src/styles/themes/kick.css                  |  126 +
 web/src/styles/themes/kodama-grove.css          |  134 +
 web/src/styles/themes/midnight-bloom.css        |  124 +
 web/src/styles/themes/mocha-mousse.css          |  134 +
 web/src/styles/themes/modern-minimal.css        |  122 +
 web/src/styles/themes/mono.css                  |  122 +
 web/src/styles/themes/nature.css                |  122 +
 web/src/styles/themes/neo-brutalism.css         |  122 +
 web/src/styles/themes/northern-lights.css       |  122 +
 web/src/styles/themes/ocean-breeze.css          |  124 +
 web/src/styles/themes/pastel-dreams.css         |  124 +
 web/src/styles/themes/perpetuity.css            |  134 +
 web/src/styles/themes/quantum-rose.css          |  134 +
 web/src/styles/themes/retro-arcade.css          |  134 +
 web/src/styles/themes/solar-dusk.css            |  129 +
 web/src/styles/themes/spotify.css               |  114 +
 web/src/styles/themes/starry-night.css          |  126 +
 web/src/styles/themes/stripe.css                |  128 +
 web/src/styles/themes/sunset-horizon.css        |  124 +
 web/src/styles/themes/supabase.css              |  124 +
 web/src/styles/themes/t3-chat.css               |  132 +
 web/src/styles/themes/tangerine.css             |  124 +
 web/src/styles/themes/twitch.css                |  128 +
 web/src/styles/themes/twitter.css               |  134 +
 web/src/styles/themes/vercel.css                |  124 +
 web/src/styles/themes/vintage-paper.css         |  134 +
 web/src/styles/themes/windows98.css             |  123 +
 57 files changed, 9608 insertions(+), 191 deletions(-)
diff --git a/web/components.json b/web/components.json
new file mode 100644
index 0000000..8d886db
--- /dev/null
+++ b/web/components.json
@@ -0,0 +1,25 @@
+{
+  "$schema": "https://ui.shadcn.com/schema.json",
+  "style": "base-nova",
+  "rsc": true,
+  "tsx": true,
+  "tailwind": {
+    "config": "",
+    "css": "src/app/globals.css",
+    "baseColor": "neutral",
+    "cssVariables": true,
+    "prefix": ""
+  },
+  "iconLibrary": "lucide",
+  "rtl": false,
+  "aliases": {
+    "components": "@/components",
+    "utils": "@/lib/utils",
+    "ui": "@/components/ui",
+    "lib": "@/lib",
+    "hooks": "@/hooks"
+  },
+  "menuColor": "default",
+  "menuAccent": "subtle",
+  "registries": {}
+}
diff --git a/web/package-lock.json b/web/package-lock.json
index 0ad08cf..0604189 100644
--- a/web/package-lock.json
+++ b/web/package-lock.json
@@ -1,23 +1,31 @@
 {
   "name": "web",
   "version": "0.1.0",
   "lockfileVersion": 3,
   "requires": true,
   "packages": {
     "": {
       "name": "web",
       "version": "0.1.0",
       "dependencies": {
+        "@base-ui/react": "^1.7.0",
+        "class-variance-authority": "^0.7.1",
+        "clsx": "^2.1.1",
+        "lucide-react": "^1.37.0",
         "next": "16.3.3",
+        "next-themes": "^0.4.6",
         "react": "19.2.8",
-        "react-dom": "19.2.8"
+        "react-dom": "19.2.8",
+        "shadcn": "^4.19.0",
+        "tailwind-merge": "^3.6.0",
+        "tw-animate-css": "^1.4.0"
       },
       "devDependencies": {
         "@tailwindcss/postcss": "^4",
         "@types/node": "^20",
         "@types/react": "^19",
         "@types/react-dom": "^19",
         "eslint": "^9",
         "eslint-config-next": "16.3.3",
         "tailwindcss": "^4",
         "typescript": "^5"
@@ -29,42 +37,39 @@
       "license": "MIT",
       "engines": {
         "node": ">=10"
       },
       "funding": {
         "url": "https://github.com/sponsors/sindresorhus"
       }
     },
     "node_modules/@babel/code-frame": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/helper-validator-identifier": "^7.29.7",
         "js-tokens": "^4.0.0",
         "picocolors": "^1.1.1"
       },
       "engines": {
         "node": ">=6.9.0"
       }
     },
     "node_modules/@babel/compat-data": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=6.9.0"
       }
     },
     "node_modules/@babel/core": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/code-frame": "^7.29.7",
         "@babel/generator": "^7.29.7",
         "@babel/helper-compilation-targets": "^7.29.7",
         "@babel/helper-module-transforms": "^7.29.7",
         "@babel/helpers": "^7.29.7",
         "@babel/parser": "^7.29.7",
         "@babel/template": "^7.29.7",
         "@babel/traverse": "^7.29.7",
@@ -79,176 +84,647 @@
       "engines": {
         "node": ">=6.9.0"
       },
       "funding": {
         "type": "opencollective",
         "url": "https://opencollective.com/babel"
       }
     },
     "node_modules/@babel/generator": {
       "version": "7.29.8",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/parser": "^7.29.8",
         "@babel/types": "^7.29.8",
         "@jridgewell/gen-mapping": "^0.3.12",
         "@jridgewell/trace-mapping": "^0.3.28",
         "jsesc": "^3.0.2"
       },
       "engines": {
         "node": ">=6.9.0"
       }
     },
+    "node_modules/@babel/helper-annotate-as-pure": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/helper-annotate-as-pure/-/helper-annotate-as-pure-7.29.7.tgz",
+      "integrity": "sha512-OoK6239jHPuSQOoS0kfTVKn0b/rVTk0seKq4Gd2UMLtmOVLjDC0ki3e+c90Trqv2gMfvJFqkiljrr568+qddiw==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/types": "^7.29.7"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      }
+    },
     "node_modules/@babel/helper-compilation-targets": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/compat-data": "^7.29.7",
         "@babel/helper-validator-option": "^7.29.7",
         "browserslist": "^4.24.0",
         "lru-cache": "^5.1.1",
         "semver": "^6.3.1"
       },
       "engines": {
         "node": ">=6.9.0"
       }
     },
+    "node_modules/@babel/helper-create-class-features-plugin": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/helper-create-class-features-plugin/-/helper-create-class-features-plugin-7.29.7.tgz",
+      "integrity": "sha512-IY3ZD9Tmooqr3TUhc3DUWxiuo8xx1DWLhd5M7hQ+ZWJamqM2BbalrBJb2MisSLoYorOj75U03qULCxQTY9r3hg==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/helper-annotate-as-pure": "^7.29.7",
+        "@babel/helper-member-expression-to-functions": "^7.29.7",
+        "@babel/helper-optimise-call-expression": "^7.29.7",
+        "@babel/helper-replace-supers": "^7.29.7",
+        "@babel/helper-skip-transparent-expression-wrappers": "^7.29.7",
+        "@babel/traverse": "^7.29.7",
+        "semver": "^6.3.1"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      },
+      "peerDependencies": {
+        "@babel/core": "^7.0.0"
+      }
+    },
     "node_modules/@babel/helper-globals": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=6.9.0"
       }
     },
+    "node_modules/@babel/helper-member-expression-to-functions": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/helper-member-expression-to-functions/-/helper-member-expression-to-functions-7.29.7.tgz",
+      "integrity": "sha512-j+7JYmk1JYDtACIGj0QJqqWZjoUpMoEikQGADMaHgCMCSDqd2+P32rfcibUNrGOMWrlzK1WJBdxrB3JJQZwWtg==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/traverse": "^7.29.7",
+        "@babel/types": "^7.29.7"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      }
+    },
     "node_modules/@babel/helper-module-imports": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/traverse": "^7.29.7",
         "@babel/types": "^7.29.7"
       },
       "engines": {
         "node": ">=6.9.0"
       }
     },
     "node_modules/@babel/helper-module-transforms": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/helper-module-imports": "^7.29.7",
         "@babel/helper-validator-identifier": "^7.29.7",
         "@babel/traverse": "^7.29.7"
       },
       "engines": {
         "node": ">=6.9.0"
       },
       "peerDependencies": {
         "@babel/core": "^7.0.0"
       }
     },
+    "node_modules/@babel/helper-optimise-call-expression": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/helper-optimise-call-expression/-/helper-optimise-call-expression-7.29.7.tgz",
+      "integrity": "sha512-+kmGVjcT9RGYzoDwdwEqEvGgKe3BYq+O1iGzjFubaNgZHwYHP6lsF2Yghf4kEuv9BV7tYDZ913aBW9am6YKong==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/types": "^7.29.7"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      }
+    },
+    "node_modules/@babel/helper-plugin-utils": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/helper-plugin-utils/-/helper-plugin-utils-7.29.7.tgz",
+      "integrity": "sha512-G7sHYigPY17oO5SYWnfD/0MTBwVR781S/JI643e/JhUYgVgWE/61SoW3NH9KWUKyKq5LVh3npif99Wkt6j86Jw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6.9.0"
+      }
+    },
+    "node_modules/@babel/helper-replace-supers": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/helper-replace-supers/-/helper-replace-supers-7.29.7.tgz",
+      "integrity": "sha512-atfGXWSeCiF4DnKZIfmJfQRkSw9b9gNNXR1kqKjbhG4pGYCOnkp8OcTB8E3NXjBu8NpheSnOeNKz8KT7UNFTmQ==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/helper-member-expression-to-functions": "^7.29.7",
+        "@babel/helper-optimise-call-expression": "^7.29.7",
+        "@babel/traverse": "^7.29.7"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      },
+      "peerDependencies": {
+        "@babel/core": "^7.0.0"
+      }
+    },
+    "node_modules/@babel/helper-skip-transparent-expression-wrappers": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/helper-skip-transparent-expression-wrappers/-/helper-skip-transparent-expression-wrappers-7.29.7.tgz",
+      "integrity": "sha512-brcMGQaVzIeUb+6/bs1Av0f8YuNNjKY2JyvfRCsFuFsdKccEQ5Ges2y74D74NZ1Rz8lKJ9ksJkfqwQFJ/iNEyQ==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/traverse": "^7.29.7",
+        "@babel/types": "^7.29.7"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      }
+    },
     "node_modules/@babel/helper-string-parser": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=6.9.0"
       }
     },
     "node_modules/@babel/helper-validator-identifier": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=6.9.0"
       }
     },
     "node_modules/@babel/helper-validator-option": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=6.9.0"
       }
     },
     "node_modules/@babel/helpers": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/template": "^7.29.7",
         "@babel/types": "^7.29.7"
       },
       "engines": {
         "node": ">=6.9.0"
       }
     },
     "node_modules/@babel/parser": {
       "version": "7.29.8",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/types": "^7.29.8"
       },
       "bin": {
         "parser": "bin/babel-parser.js"
       },
       "engines": {
         "node": ">=6.0.0"
       }
     },
+    "node_modules/@babel/plugin-syntax-jsx": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/plugin-syntax-jsx/-/plugin-syntax-jsx-7.29.7.tgz",
+      "integrity": "sha512-TSu8+mHCoEaaCDEZ0I3+6mvTBYR4PCxQwf2z9/r5Tbztv6NaLR3B9thGTTxX2WGuGHJqRiAbKPeGTJ5XWXVg6A==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/helper-plugin-utils": "^7.29.7"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      },
+      "peerDependencies": {
+        "@babel/core": "^7.0.0-0"
+      }
+    },
+    "node_modules/@babel/plugin-syntax-typescript": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/plugin-syntax-typescript/-/plugin-syntax-typescript-7.29.7.tgz",
+      "integrity": "sha512-ngr+82Sh0xMz25TPCZi+nC2iTzjfCdWS2ONXTp/PtSCHCgaCNBpdMqgvJ2ccdLlClVZ7sisIgB914j/JFe+RZA==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/helper-plugin-utils": "^7.29.7"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      },
+      "peerDependencies": {
+        "@babel/core": "^7.0.0-0"
+      }
+    },
+    "node_modules/@babel/plugin-transform-modules-commonjs": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/plugin-transform-modules-commonjs/-/plugin-transform-modules-commonjs-7.29.7.tgz",
+      "integrity": "sha512-j0vCldybPC5b5dwCQOJ21uKtHzt7hxLygJTg9eF1ScfaikEDNfzn94XoW5Fi+seBR0nCyL23xaBFFkq7dTM8XQ==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/helper-module-transforms": "^7.29.7",
+        "@babel/helper-plugin-utils": "^7.29.7"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      },
+      "peerDependencies": {
+        "@babel/core": "^7.0.0-0"
+      }
+    },
+    "node_modules/@babel/plugin-transform-typescript": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/plugin-transform-typescript/-/plugin-transform-typescript-7.29.7.tgz",
+      "integrity": "sha512-jK52h8LaLc7JarhQV2ofeFMts4H7vnOXnqZNA6fYglBTZewRBE51KWt3BUltW1P+KoPsYkHoJeXePuz4zo2LMw==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/helper-annotate-as-pure": "^7.29.7",
+        "@babel/helper-create-class-features-plugin": "^7.29.7",
+        "@babel/helper-plugin-utils": "^7.29.7",
+        "@babel/helper-skip-transparent-expression-wrappers": "^7.29.7",
+        "@babel/plugin-syntax-typescript": "^7.29.7"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      },
+      "peerDependencies": {
+        "@babel/core": "^7.0.0-0"
+      }
+    },
+    "node_modules/@babel/preset-typescript": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/preset-typescript/-/preset-typescript-7.29.7.tgz",
+      "integrity": "sha512-/Foi8vKY2EVbed/1eZx0gJEEwHAIxogrySI7rULcRIvhZzbvoE/b5qG5Ghc0WKAFKOHA9SD1x7RsFlOYdutIiQ==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/helper-plugin-utils": "^7.29.7",
+        "@babel/helper-validator-option": "^7.29.7",
+        "@babel/plugin-syntax-jsx": "^7.29.7",
+        "@babel/plugin-transform-modules-commonjs": "^7.29.7",
+        "@babel/plugin-transform-typescript": "^7.29.7"
+      },
+      "engines": {
+        "node": ">=6.9.0"
+      },
+      "peerDependencies": {
+        "@babel/core": "^7.0.0-0"
+      }
+    },
+    "node_modules/@babel/runtime": {
+      "version": "7.29.7",
+      "resolved": "https://registry.npmjs.org/@babel/runtime/-/runtime-7.29.7.tgz",
+      "integrity": "sha512-Nq8OhGWiZIZGV6hLHoyAKLLcJihP/xFeBMGJoUrxTX2psI8dCifzLhZISFb+VWS3wFMRDmCGw5R+dOySCqPLhw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6.9.0"
+      }
+    },
     "node_modules/@babel/template": {
       "version": "7.29.7",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/code-frame": "^7.29.7",
         "@babel/parser": "^7.29.7",
         "@babel/types": "^7.29.7"
       },
       "engines": {
         "node": ">=6.9.0"
       }
     },
     "node_modules/@babel/traverse": {
       "version": "7.29.8",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/code-frame": "^7.29.7",
         "@babel/generator": "^7.29.8",
         "@babel/helper-globals": "^7.29.7",
         "@babel/parser": "^7.29.8",
         "@babel/template": "^7.29.7",
         "@babel/types": "^7.29.8",
         "debug": "^4.3.1"
       },
       "engines": {
         "node": ">=6.9.0"
       }
     },
     "node_modules/@babel/types": {
       "version": "7.29.8",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@babel/helper-string-parser": "^7.29.7",
         "@babel/helper-validator-identifier": "^7.29.7"
       },
       "engines": {
         "node": ">=6.9.0"
       }
     },
+    "node_modules/@base-ui/react": {
+      "version": "1.7.0",
+      "resolved": "https://registry.npmjs.org/@base-ui/react/-/react-1.7.0.tgz",
+      "integrity": "sha512-j+8QjX44C32jrXD/qyEAGpFr70FRpGL2CY61mQd9nBPWN737CK0xxD1ceJ055rW4RtdvFDT1e7otzdlfxvsYug==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/runtime": "^7.29.2",
+        "@base-ui/utils": "0.3.2",
+        "@floating-ui/react-dom": "^2.1.9",
+        "@floating-ui/utils": "^0.2.12",
+        "use-sync-external-store": "^1.6.0"
+      },
+      "engines": {
+        "node": ">=14.0.0"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/mui-org"
+      },
+      "peerDependencies": {
+        "@date-fns/tz": "^1.2.0",
+        "@types/react": "^17 || ^18 || ^19",
+        "date-fns": "^4.0.0",
+        "react": "^17 || ^18 || ^19",
+        "react-dom": "^17 || ^18 || ^19"
+      },
+      "peerDependenciesMeta": {
+        "@date-fns/tz": {
+          "optional": true
+        },
+        "@types/react": {
+          "optional": true
+        },
+        "date-fns": {
+          "optional": true
+        }
+      }
+    },
+    "node_modules/@base-ui/utils": {
+      "version": "0.3.2",
+      "resolved": "https://registry.npmjs.org/@base-ui/utils/-/utils-0.3.2.tgz",
+      "integrity": "sha512-oWy1aq/I2GmYjpl4PhEAhzflF8VPGKgZeq0xAWTbfD5KBWyxcN0ZP2+WHSUm/5Z6lVMBDLReLcoXwSYoRc/zNQ==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/runtime": "^7.29.2",
+        "@floating-ui/utils": "^0.2.12",
+        "reselect": "^5.2.0",
+        "use-sync-external-store": "^1.6.0"
+      },
+      "peerDependencies": {
+        "@types/react": "^17 || ^18 || ^19",
+        "react": "^17 || ^18 || ^19",
+        "react-dom": "^17 || ^18 || ^19"
+      },
+      "peerDependenciesMeta": {
+        "@types/react": {
+          "optional": true
+        }
+      }
+    },
+    "node_modules/@dotenvx/dotenvx": {
+      "version": "1.75.1",
+      "resolved": "https://registry.npmjs.org/@dotenvx/dotenvx/-/dotenvx-1.75.1.tgz",
+      "integrity": "sha512-/BITOC9dmS/edY2zQwZNicQ059O6RKabtQfyEafV0nGtfYRNHYy1DIPiYVcov40+tob9hfmBnbR963dS+EQ1DQ==",
+      "license": "BSD-3-Clause",
+      "dependencies": {
+        "@dotenvx/primitives": "^0.8.0",
+        "commander": "^11.1.0",
+        "conf": "^10.2.0",
+        "dotenv": "^17.2.1",
+        "enquirer": "^2.4.1",
+        "env-paths": "^2.2.1",
+        "execa": "^5.1.1",
+        "fdir": "^6.2.0",
+        "ignore": "^5.3.0",
+        "object-treeify": "1.1.33",
+        "open": "^8.4.2",
+        "picomatch": "^4.0.4",
+        "systeminformation": "^5.22.11",
+        "undici": "^7.11.0",
+        "which": "^4.0.0",
+        "yocto-spinner": "^1.1.0"
+      },
+      "bin": {
+        "dotenvx": "src/cli/dotenvx.js"
+      },
+      "funding": {
+        "url": "https://dotenvx.com"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/commander": {
+      "version": "11.1.0",
+      "resolved": "https://registry.npmjs.org/commander/-/commander-11.1.0.tgz",
+      "integrity": "sha512-yPVavfyCcRhmorC7rWlkHn15b4wDVgVmBA7kV4QVBsF7kv/9TKJAbAXVTxvTnwP8HHKjRCJDClKbciiYS7p0DQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=16"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/define-lazy-prop": {
+      "version": "2.0.0",
+      "resolved": "https://registry.npmjs.org/define-lazy-prop/-/define-lazy-prop-2.0.0.tgz",
+      "integrity": "sha512-Ds09qNh8yw3khSjiJjiUInaGX9xlqZDY7JVryGxdxV7NPeuqQfplOpQ66yJFZut3jLa5zOwkXw1g9EI2uKh4Og==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=8"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/execa": {
+      "version": "5.1.1",
+      "resolved": "https://registry.npmjs.org/execa/-/execa-5.1.1.tgz",
+      "integrity": "sha512-8uSpZZocAZRBAPIEINJj3Lo9HyGitllczc27Eh5YYojjMFMn8yHMDMaUHE2Jqfq05D/wucwI4JGURyXt1vchyg==",
+      "license": "MIT",
+      "dependencies": {
+        "cross-spawn": "^7.0.3",
+        "get-stream": "^6.0.0",
+        "human-signals": "^2.1.0",
+        "is-stream": "^2.0.0",
+        "merge-stream": "^2.0.0",
+        "npm-run-path": "^4.0.1",
+        "onetime": "^5.1.2",
+        "signal-exit": "^3.0.3",
+        "strip-final-newline": "^2.0.0"
+      },
+      "engines": {
+        "node": ">=10"
+      },
+      "funding": {
+        "url": "https://github.com/sindresorhus/execa?sponsor=1"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/fdir": {
+      "version": "6.5.0",
+      "resolved": "https://registry.npmjs.org/fdir/-/fdir-6.5.0.tgz",
+      "integrity": "sha512-tIbYtZbucOs0BRGqPJkshJUYdL+SDH7dVM8gjy+ERp3WAUjLEFJE+02kanyHtwjWOnwrKYBiwAmM0p4kLJAnXg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12.0.0"
+      },
+      "peerDependencies": {
+        "picomatch": "^3 || ^4"
+      },
+      "peerDependenciesMeta": {
+        "picomatch": {
+          "optional": true
+        }
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/get-stream": {
+      "version": "6.0.1",
+      "resolved": "https://registry.npmjs.org/get-stream/-/get-stream-6.0.1.tgz",
+      "integrity": "sha512-ts6Wi+2j3jQjqi70w5AlN8DFnkSwC+MqmxEzdEALB2qXZYV3X/b1CTfgPLGJNMeAWxdPfU8FO1ms3NUfaHCPYg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=10"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/human-signals": {
+      "version": "2.1.0",
+      "resolved": "https://registry.npmjs.org/human-signals/-/human-signals-2.1.0.tgz",
+      "integrity": "sha512-B4FFZ6q/T2jhhksgkbEW3HBvWIfDW85snkQgawt07S7J5QXTk6BkNV+0yAeZrM5QpMAdYlocGoljn0sJ/WQkFw==",
+      "license": "Apache-2.0",
+      "engines": {
+        "node": ">=10.17.0"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/is-docker": {
+      "version": "2.2.1",
+      "resolved": "https://registry.npmjs.org/is-docker/-/is-docker-2.2.1.tgz",
+      "integrity": "sha512-F+i2BKsFrH66iaUFc0woD8sLy8getkwTwtOBjvs56Cx4CgJDeKQeqfz8wAYiSb8JOprWhHH5p77PbmYCvvUuXQ==",
+      "license": "MIT",
+      "bin": {
+        "is-docker": "cli.js"
+      },
+      "engines": {
+        "node": ">=8"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/is-stream": {
+      "version": "2.0.1",
+      "resolved": "https://registry.npmjs.org/is-stream/-/is-stream-2.0.1.tgz",
+      "integrity": "sha512-hFoiJiTl63nn+kstHGBtewWSKnQLpyb155KHheA1l39uvtO9nWIop1p3udqPcUd/xbF1VLMO4n7OI6p7RbngDg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=8"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/is-wsl": {
+      "version": "2.2.0",
+      "resolved": "https://registry.npmjs.org/is-wsl/-/is-wsl-2.2.0.tgz",
+      "integrity": "sha512-fKzAra0rGJUUBwGBgNkHZuToZcn+TtXHpeCgmkMJMMYx1sQDYaCSyjJBSCa2nH1DGm7s3n1oBnohoVTBaN7Lww==",
+      "license": "MIT",
+      "dependencies": {
+        "is-docker": "^2.0.0"
+      },
+      "engines": {
+        "node": ">=8"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/isexe": {
+      "version": "3.1.5",
+      "resolved": "https://registry.npmjs.org/isexe/-/isexe-3.1.5.tgz",
+      "integrity": "sha512-6B3tLtFqtQS4ekarvLVMZ+X+VlvQekbe4taUkf/rhVO3d/h0M2rfARm/pXLcPEsjjMsFgrFgSrhQIxcSVrBz8w==",
+      "license": "BlueOak-1.0.0",
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/npm-run-path": {
+      "version": "4.0.1",
+      "resolved": "https://registry.npmjs.org/npm-run-path/-/npm-run-path-4.0.1.tgz",
+      "integrity": "sha512-S48WzZW777zhNIrn7gxOlISNAqi9ZC/uQFnRdbeIHhZhCA6UqpkOT8T1G7BvfdgP4Er8gF4sUbaS0i7QvIfCWw==",
+      "license": "MIT",
+      "dependencies": {
+        "path-key": "^3.0.0"
+      },
+      "engines": {
+        "node": ">=8"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/open": {
+      "version": "8.4.2",
+      "resolved": "https://registry.npmjs.org/open/-/open-8.4.2.tgz",
+      "integrity": "sha512-7x81NCL719oNbsq/3mh+hVrAWmFuEYUqrq/Iw3kUzH8ReypT9QQ0BLoJS7/G9k6N81XjW4qHWtjWwe/9eLy1EQ==",
+      "license": "MIT",
+      "dependencies": {
+        "define-lazy-prop": "^2.0.0",
+        "is-docker": "^2.1.1",
+        "is-wsl": "^2.2.0"
+      },
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/picomatch": {
+      "version": "4.0.7",
+      "resolved": "https://registry.npmjs.org/picomatch/-/picomatch-4.0.7.tgz",
+      "integrity": "sha512-qcJu88Q2IWqJsDD529JKMdwGm/dvInW4HvQnRwiH9JtihJvzGOscDtHE3x1pBKeUOTysQ8kVmLnJ2kJu7yhcGA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/jonschlinkert"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/signal-exit": {
+      "version": "3.0.7",
+      "resolved": "https://registry.npmjs.org/signal-exit/-/signal-exit-3.0.7.tgz",
+      "integrity": "sha512-wnD2ZE+l+SPC/uoS0vXeE9L1+0wuaMqKlfz9AMUo38JsyLSBWSFcHR1Rri62LZc12vLr1gb3jl7iwQhgwpAbGQ==",
+      "license": "ISC"
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/strip-final-newline": {
+      "version": "2.0.0",
+      "resolved": "https://registry.npmjs.org/strip-final-newline/-/strip-final-newline-2.0.0.tgz",
+      "integrity": "sha512-BrpvfNAE3dcvq7ll3xVumzjKjZQ5tI1sEUIKr3Uoks0XUl45St3FlatVqef9prk4jRDzhW6WZg+3bk93y6pLjA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6"
+      }
+    },
+    "node_modules/@dotenvx/dotenvx/node_modules/which": {
+      "version": "4.0.0",
+      "resolved": "https://registry.npmjs.org/which/-/which-4.0.0.tgz",
+      "integrity": "sha512-GlaYyEb07DPxYCKhKzplCWBJtvxZcZMrL+4UkrTSJHHPyZU4mYYTv3qaOe77H7EODLSSopAUFAc6W8U4yqvscg==",
+      "license": "ISC",
+      "dependencies": {
+        "isexe": "^3.1.1"
+      },
+      "bin": {
+        "node-which": "bin/which.js"
+      },
+      "engines": {
+        "node": "^16.13.0 || >=18.0.0"
+      }
+    },
+    "node_modules/@dotenvx/primitives": {
+      "version": "0.8.0",
+      "resolved": "https://registry.npmjs.org/@dotenvx/primitives/-/primitives-0.8.0.tgz",
+      "integrity": "sha512-VYJy0uhFm9zTJ1TxBaW/pA8bjbOM/OttaNMwZ1RHG4JKyRG7DhSdiqD1ipQoAyoD22olUtxbP78W9xY3Wd11bg==",
+      "license": "BSD-3-Clause"
+    },
     "node_modules/@emnapi/core": {
       "version": "1.10.0",
       "resolved": "https://registry.npmjs.org/@emnapi/core/-/core-1.10.0.tgz",
       "integrity": "sha512-yq6OkJ4p82CAfPl0u9mQebQHKPJkY7WrIuk205cTYnYe+k2Z8YBh11FrbRG/H6ihirqcacOgl2BIO8oyMQLeXw==",
       "dev": true,
       "license": "MIT",
       "optional": true,
       "dependencies": {
         "@emnapi/wasi-threads": "1.2.1",
         "tslib": "^2.4.0"
@@ -392,20 +868,70 @@
       "dev": true,
       "license": "Apache-2.0",
       "dependencies": {
         "@eslint/core": "^0.17.0",
         "levn": "^0.4.1"
       },
       "engines": {
         "node": "^18.18.0 || ^20.9.0 || >=21.1.0"
       }
     },
+    "node_modules/@floating-ui/core": {
+      "version": "1.8.0",
+      "resolved": "https://registry.npmjs.org/@floating-ui/core/-/core-1.8.0.tgz",
+      "integrity": "sha512-0CIZ5itps/8x7BG8dEIhs53BvCUH2PCoogtakwRTut+Arm58sJooJ0AuZhLw2HJYIR5cMLNPBSS728sPho2khQ==",
+      "license": "MIT",
+      "dependencies": {
+        "@floating-ui/utils": "^0.2.12"
+      }
+    },
+    "node_modules/@floating-ui/dom": {
+      "version": "1.8.0",
+      "resolved": "https://registry.npmjs.org/@floating-ui/dom/-/dom-1.8.0.tgz",
+      "integrity": "sha512-yXSrzeHZBTZadLOlfyhCkJHNeLJnHRnRInwdZ40L7ZiaAtrBwoYlsDrX3v5zB1Utk7CLfzcOVnVVWoXEky7Ceg==",
+      "license": "MIT",
+      "dependencies": {
+        "@floating-ui/core": "^1.8.0",
+        "@floating-ui/utils": "^0.2.12"
+      }
+    },
+    "node_modules/@floating-ui/react-dom": {
+      "version": "2.1.9",
+      "resolved": "https://registry.npmjs.org/@floating-ui/react-dom/-/react-dom-2.1.9.tgz",
+      "integrity": "sha512-JDjEFGCpImxDCA7JJKviA0M9+RtmJdj0m/NVU5IMgBK+AmZouAQQ7/+2GLH0GXXY0YMw9oXPB8hKdbPYg5QLYg==",
+      "license": "MIT",
+      "dependencies": {
+        "@floating-ui/dom": "^1.8.0"
+      },
+      "peerDependencies": {
+        "react": ">=16.8.0",
+        "react-dom": ">=16.8.0"
+      }
+    },
+    "node_modules/@floating-ui/utils": {
+      "version": "0.2.12",
+      "resolved": "https://registry.npmjs.org/@floating-ui/utils/-/utils-0.2.12.tgz",
+      "integrity": "sha512-HpCo8tmWzLVad5s2d19EhAz5zqrrQ6s69qd6moPMQvkOuSwDT1YgRfWSVuc4ennqrgv3OHppiOGMQ7oC13yIww==",
+      "license": "MIT"
+    },
+    "node_modules/@hono/node-server": {
+      "version": "2.1.1",
+      "resolved": "https://registry.npmjs.org/@hono/node-server/-/node-server-2.1.1.tgz",
+      "integrity": "sha512-ELuehkj5VCBdgEw9zs+ivkKwyzzUCSQuE96YmiPvn1ECBoZCczbFXJLeEGMTYjphP6gydh4pHMqEYPVMYUVgQg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=20"
+      },
+      "peerDependencies": {
+        "hono": "^4"
+      }
+    },
     "node_modules/@humanfs/core": {
       "version": "0.19.2",
       "dev": true,
       "license": "Apache-2.0",
       "dependencies": {
         "@humanfs/types": "^0.15.0"
       },
       "engines": {
         "node": ">=18.18.0"
       }
@@ -995,58 +1521,115 @@
       ],
       "engines": {
         "node": ">=20.9.0"
       },
       "funding": {
         "url": "https://opencollective.com/libvips"
       }
     },
     "node_modules/@jridgewell/gen-mapping": {
       "version": "0.3.13",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@jridgewell/sourcemap-codec": "^1.5.0",
         "@jridgewell/trace-mapping": "^0.3.24"
       }
     },
     "node_modules/@jridgewell/remapping": {
       "version": "2.3.5",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@jridgewell/gen-mapping": "^0.3.5",
         "@jridgewell/trace-mapping": "^0.3.24"
       }
     },
     "node_modules/@jridgewell/resolve-uri": {
       "version": "3.1.2",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=6.0.0"
       }
     },
     "node_modules/@jridgewell/sourcemap-codec": {
       "version": "1.6.0",
-      "dev": true,
       "license": "MIT"
     },
     "node_modules/@jridgewell/trace-mapping": {
       "version": "0.3.31",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@jridgewell/resolve-uri": "^3.1.0",
         "@jridgewell/sourcemap-codec": "^1.4.14"
       }
     },
+    "node_modules/@modelcontextprotocol/sdk": {
+      "version": "1.30.0",
+      "resolved": "https://registry.npmjs.org/@modelcontextprotocol/sdk/-/sdk-1.30.0.tgz",
+      "integrity": "sha512-xKd8OIzlqNzcqcNumGAa6g+PW2kjD5vrpcKOnfldAUPP3j7lnqMPwlTXQm8gF+UwH72z0lqaRbjr9hqGz0eITA==",
+      "license": "MIT",
+      "dependencies": {
+        "@hono/node-server": "^1.19.9 || ^2.0.5",
+        "ajv": "^8.17.1",
+        "ajv-formats": "^3.0.1",
+        "content-type": "^1.0.5",
+        "cors": "^2.8.5",
+        "cross-spawn": "^7.0.5",
+        "eventsource": "^3.0.2",
+        "eventsource-parser": "^3.0.0",
+        "express": "^5.2.1",
+        "express-rate-limit": "^8.2.1",
+        "hono": "^4.11.4",
+        "jose": "^6.1.3",
+        "json-schema-typed": "^8.0.2",
+        "pkce-challenge": "^5.0.0",
+        "raw-body": "^3.0.0",
+        "zod": "^3.25 || ^4.0",
+        "zod-to-json-schema": "^3.25.1"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "peerDependencies": {
+        "@cfworker/json-schema": "^4.1.1",
+        "zod": "^3.25 || ^4.0"
+      },
+      "peerDependenciesMeta": {
+        "@cfworker/json-schema": {
+          "optional": true
+        },
+        "zod": {
+          "optional": false
+        }
+      }
+    },
+    "node_modules/@modelcontextprotocol/sdk/node_modules/ajv": {
+      "version": "8.20.0",
+      "resolved": "https://registry.npmjs.org/ajv/-/ajv-8.20.0.tgz",
+      "integrity": "sha512-Thbli+OlOj+iMPYFBVBfJ3OmCAnaSyNn4M1vz9T6Gka5Jt9ba/HIR56joy65tY6kx/FCF5VXNB819Y7/GUrBGA==",
+      "license": "MIT",
+      "dependencies": {
+        "fast-deep-equal": "^3.1.3",
+        "fast-uri": "^3.0.1",
+        "json-schema-traverse": "^1.0.0",
+        "require-from-string": "^2.0.2"
+      },
+      "funding": {
+        "type": "github",
+        "url": "https://github.com/sponsors/epoberezkin"
+      }
+    },
+    "node_modules/@modelcontextprotocol/sdk/node_modules/json-schema-traverse": {
+      "version": "1.0.0",
+      "resolved": "https://registry.npmjs.org/json-schema-traverse/-/json-schema-traverse-1.0.0.tgz",
+      "integrity": "sha512-NM8/P9n3XjXhIZn1lLhkFaACTOURQXjWhV4BA/RnOv8xvgqtqpAX9IO4mRQxSx1Rlo4tqzeqb0sOlruaOy3dug==",
+      "license": "MIT"
+    },
     "node_modules/@napi-rs/wasm-runtime": {
       "version": "1.2.3",
       "resolved": "https://registry.npmjs.org/@napi-rs/wasm-runtime/-/wasm-runtime-1.2.3.tgz",
       "integrity": "sha512-UMduMbqO5s5zF2NkNacMT/yK5Y5QiKvWr2+50bzIIxFDwVJ2h49b+oyjaCGPhJxd2/gC2x39EHv/gHVuu36x2Q==",
       "dev": true,
       "license": "MIT",
       "optional": true,
       "dependencies": {
         "@tybys/wasm-util": "^0.10.3"
       },
@@ -1236,41 +1819,38 @@
       "optional": true,
       "os": [
         "win32"
       ],
       "engines": {
         "node": ">= 10"
       }
     },
     "node_modules/@nodelib/fs.scandir": {
       "version": "2.1.5",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@nodelib/fs.stat": "2.0.5",
         "run-parallel": "^1.1.9"
       },
       "engines": {
         "node": ">= 8"
       }
     },
     "node_modules/@nodelib/fs.stat": {
       "version": "2.0.5",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 8"
       }
     },
     "node_modules/@nodelib/fs.walk": {
       "version": "1.2.8",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "@nodelib/fs.scandir": "2.1.5",
         "fastq": "^1.6.0"
       },
       "engines": {
         "node": ">= 8"
       }
     },
     "node_modules/@nolyfill/is-core-module": {
@@ -1279,20 +1859,38 @@
       "license": "MIT",
       "engines": {
         "node": ">=12.4.0"
       }
     },
     "node_modules/@rtsao/scc": {
       "version": "1.1.0",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/@sec-ant/readable-stream": {
+      "version": "0.4.1",
+      "resolved": "https://registry.npmjs.org/@sec-ant/readable-stream/-/readable-stream-0.4.1.tgz",
+      "integrity": "sha512-831qok9r2t8AlxLko40y2ebgSDhenenCatLVeW/uBtnHPyhHOvG0C7TvfgecV+wHzIm5KUICgzmVpWS+IMEAeg==",
+      "license": "MIT"
+    },
+    "node_modules/@sindresorhus/merge-streams": {
+      "version": "4.0.0",
+      "resolved": "https://registry.npmjs.org/@sindresorhus/merge-streams/-/merge-streams-4.0.0.tgz",
+      "integrity": "sha512-tlqY9xq5ukxTUZBmoOp+m61cqwQD5pHJtFY3Mn8CA8ps6yghLH/Hw8UPdqg4OLmFW3IFlcXnQNmo/dh8HzXYIQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/@swc/helpers": {
       "version": "0.5.23",
       "license": "Apache-2.0",
       "dependencies": {
         "tslib": "^2.8.0"
       }
     },
     "node_modules/@tailwindcss/node": {
       "version": "4.3.3",
       "dev": true,
@@ -1561,20 +2159,95 @@
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "@alloc/quick-lru": "^5.2.0",
         "@tailwindcss/node": "4.3.3",
         "@tailwindcss/oxide": "4.3.3",
         "postcss": "^8.5.16",
         "tailwindcss": "4.3.3"
       }
     },
+    "node_modules/@ts-morph/common": {
+      "version": "0.27.0",
+      "resolved": "https://registry.npmjs.org/@ts-morph/common/-/common-0.27.0.tgz",
+      "integrity": "sha512-Wf29UqxWDpc+i61k3oIOzcUfQt79PIT9y/MWfAGlrkjg6lBC1hwDECLXPVJAhWjiGbfBCxZd65F/LIZF3+jeJQ==",
+      "license": "MIT",
+      "dependencies": {
+        "fast-glob": "^3.3.3",
+        "minimatch": "^10.0.1",
+        "path-browserify": "^1.0.1"
+      }
+    },
+    "node_modules/@ts-morph/common/node_modules/balanced-match": {
+      "version": "4.0.4",
+      "resolved": "https://registry.npmjs.org/balanced-match/-/balanced-match-4.0.4.tgz",
+      "integrity": "sha512-BLrgEcRTwX2o6gGxGOCNyMvGSp35YofuYzw9h1IMTRmKqttAZZVU67bdb9Pr2vUHA8+j3i2tJfjO6C6+4myGTA==",
+      "license": "MIT",
+      "engines": {
+        "node": "18 || 20 || >=22"
+      }
+    },
+    "node_modules/@ts-morph/common/node_modules/brace-expansion": {
+      "version": "5.0.9",
+      "resolved": "https://registry.npmjs.org/brace-expansion/-/brace-expansion-5.0.9.tgz",
+      "integrity": "sha512-ScQ4IuvIEF1TMlP7Zt+vjJ//9zlPb2SDcxWxM3bk8s6t6GGdJ7KO1dCcTidOPJKePW30LE/2cT7wCyPho9/Wxg==",
+      "license": "MIT",
+      "dependencies": {
+        "balanced-match": "^4.0.2"
+      },
+      "engines": {
+        "node": "20 || >=22"
+      }
+    },
+    "node_modules/@ts-morph/common/node_modules/fast-glob": {
+      "version": "3.3.3",
+      "resolved": "https://registry.npmjs.org/fast-glob/-/fast-glob-3.3.3.tgz",
+      "integrity": "sha512-7MptL8U0cqcFdzIzwOTHoilX9x5BrNqye7Z/LuC7kCMRio1EMSyqRK3BEAUD7sXRq4iT4AzTVuZdhgQ2TCvYLg==",
+      "license": "MIT",
+      "dependencies": {
+        "@nodelib/fs.stat": "^2.0.2",
+        "@nodelib/fs.walk": "^1.2.3",
+        "glob-parent": "^5.1.2",
+        "merge2": "^1.3.0",
+        "micromatch": "^4.0.8"
+      },
+      "engines": {
+        "node": ">=8.6.0"
+      }
+    },
+    "node_modules/@ts-morph/common/node_modules/glob-parent": {
+      "version": "5.1.2",
+      "resolved": "https://registry.npmjs.org/glob-parent/-/glob-parent-5.1.2.tgz",
+      "integrity": "sha512-AOIgSQCepiJYwP3ARnGx+5VnTu2HBYdzbGP45eLw1vr3zB3vZLeyed1sC9hnbcOc9/SrMyM5RPQrkGz4aS9Zow==",
+      "license": "ISC",
+      "dependencies": {
+        "is-glob": "^4.0.1"
+      },
+      "engines": {
+        "node": ">= 6"
+      }
+    },
+    "node_modules/@ts-morph/common/node_modules/minimatch": {
+      "version": "10.2.6",
+      "resolved": "https://registry.npmjs.org/minimatch/-/minimatch-10.2.6.tgz",
+      "integrity": "sha512-vpLQEs+VLCr1nU0BXS07maYoFwlDAH0gngQuuttxIwutDFEMHq2blX+8vpgxDdK3J1PwjCJiep77OitTZ4Ll1A==",
+      "license": "BlueOak-1.0.0",
+      "dependencies": {
+        "brace-expansion": "^5.0.8"
+      },
+      "engines": {
+        "node": "18 || 20 || >=22"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/isaacs"
+      }
+    },
     "node_modules/@tybys/wasm-util": {
       "version": "0.10.3",
       "resolved": "https://registry.npmjs.org/@tybys/wasm-util/-/wasm-util-0.10.3.tgz",
       "integrity": "sha512-F3fo1MYrRJYL3zER0OUOmkutjr1Vp23m7OsSgp7nq4SP6OqX6C/56XFIPAl5bt3zaBRjmW7SGz3u/6LwFpYcOg==",
       "dev": true,
       "license": "MIT",
       "optional": true,
       "dependencies": {
         "tslib": "^2.4.0"
       }
@@ -1597,34 +2270,40 @@
     "node_modules/@types/node": {
       "version": "20.19.43",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "undici-types": "~6.21.0"
       }
     },
     "node_modules/@types/react": {
       "version": "19.2.18",
-      "dev": true,
+      "devOptional": true,
       "license": "MIT",
       "dependencies": {
         "csstype": "^3.2.2"
       }
     },
     "node_modules/@types/react-dom": {
       "version": "19.2.5",
       "dev": true,
       "license": "MIT",
       "peerDependencies": {
         "@types/react": "^19.2.0"
       }
     },
+    "node_modules/@types/validate-npm-package-name": {
+      "version": "4.0.2",
+      "resolved": "https://registry.npmjs.org/@types/validate-npm-package-name/-/validate-npm-package-name-4.0.2.tgz",
+      "integrity": "sha512-lrpDziQipxCEeK5kWxvljWYhUvOiB2A9izZd9B2AFarYAkqZshb4lPbRs7zKEic6eGtH8V/2qJW+dPp9OtF6bw==",
+      "license": "MIT"
+    },
     "node_modules/@typescript-eslint/eslint-plugin": {
       "version": "8.68.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "@eslint-community/regexpp": "^4.12.2",
         "@typescript-eslint/scope-manager": "8.68.0",
         "@typescript-eslint/type-utils": "8.68.0",
         "@typescript-eslint/utils": "8.68.0",
         "@typescript-eslint/visitor-keys": "8.68.0",
@@ -2226,71 +2905,140 @@
       "cpu": [
         "x64"
       ],
       "dev": true,
       "license": "MIT",
       "optional": true,
       "os": [
         "win32"
       ]
     },
+    "node_modules/accepts": {
+      "version": "2.0.0",
+      "resolved": "https://registry.npmjs.org/accepts/-/accepts-2.0.0.tgz",
+      "integrity": "sha512-5cvg6CtKwfgdmVqY1WIiXKc3Q1bkRqGLi+2W/6ao+6Y7gu/RCwRuAhGEzh5B4KlszSuTLgZYuqFqo5bImjNKng==",
+      "license": "MIT",
+      "dependencies": {
+        "mime-types": "^3.0.0",
+        "negotiator": "^1.0.0"
+      },
+      "engines": {
+        "node": ">= 0.6"
+      }
+    },
     "node_modules/acorn": {
       "version": "8.18.0",
       "dev": true,
       "license": "MIT",
       "bin": {
         "acorn": "bin/acorn"
       },
       "engines": {
         "node": ">=0.4.0"
       }
     },
     "node_modules/acorn-jsx": {
       "version": "5.3.2",
       "dev": true,
       "license": "MIT",
       "peerDependencies": {
         "acorn": "^6.0.0 || ^7.0.0 || ^8.0.0"
       }
     },
-    "node_modules/ajv": {
-      "version": "6.15.0",
-      "dev": true,
+    "node_modules/ajv": {
+      "version": "6.15.0",
+      "dev": true,
+      "license": "MIT",
+      "dependencies": {
+        "fast-deep-equal": "^3.1.1",
+        "fast-json-stable-stringify": "^2.0.0",
+        "json-schema-traverse": "^0.4.1",
+        "uri-js": "^4.2.2"
+      },
+      "funding": {
+        "type": "github",
+        "url": "https://github.com/sponsors/epoberezkin"
+      }
+    },
+    "node_modules/ajv-formats": {
+      "version": "3.0.1",
+      "resolved": "https://registry.npmjs.org/ajv-formats/-/ajv-formats-3.0.1.tgz",
+      "integrity": "sha512-8iUql50EUR+uUcdRQ3HDqa6EVyo3docL8g5WJ3FNcWmu62IbkGUue/pEyLBW8VGKKucTPgqeks4fIU1DA4yowQ==",
+      "license": "MIT",
+      "dependencies": {
+        "ajv": "^8.0.0"
+      },
+      "peerDependencies": {
+        "ajv": "^8.0.0"
+      },
+      "peerDependenciesMeta": {
+        "ajv": {
+          "optional": true
+        }
+      }
+    },
+    "node_modules/ajv-formats/node_modules/ajv": {
+      "version": "8.20.0",
+      "resolved": "https://registry.npmjs.org/ajv/-/ajv-8.20.0.tgz",
+      "integrity": "sha512-Thbli+OlOj+iMPYFBVBfJ3OmCAnaSyNn4M1vz9T6Gka5Jt9ba/HIR56joy65tY6kx/FCF5VXNB819Y7/GUrBGA==",
       "license": "MIT",
       "dependencies": {
-        "fast-deep-equal": "^3.1.1",
-        "fast-json-stable-stringify": "^2.0.0",
-        "json-schema-traverse": "^0.4.1",
-        "uri-js": "^4.2.2"
+        "fast-deep-equal": "^3.1.3",
+        "fast-uri": "^3.0.1",
+        "json-schema-traverse": "^1.0.0",
+        "require-from-string": "^2.0.2"
       },
       "funding": {
         "type": "github",
         "url": "https://github.com/sponsors/epoberezkin"
       }
     },
+    "node_modules/ajv-formats/node_modules/json-schema-traverse": {
+      "version": "1.0.0",
+      "resolved": "https://registry.npmjs.org/json-schema-traverse/-/json-schema-traverse-1.0.0.tgz",
+      "integrity": "sha512-NM8/P9n3XjXhIZn1lLhkFaACTOURQXjWhV4BA/RnOv8xvgqtqpAX9IO4mRQxSx1Rlo4tqzeqb0sOlruaOy3dug==",
+      "license": "MIT"
+    },
+    "node_modules/ansi-colors": {
+      "version": "4.1.3",
+      "resolved": "https://registry.npmjs.org/ansi-colors/-/ansi-colors-4.1.3.tgz",
+      "integrity": "sha512-/6w/C21Pm1A7aZitlI5Ni/2J6FFQN8i1Cvz3kHABAAbw93v/NlvKdVOqz7CCWz/3iv/JplRSEEZ83XION15ovw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6"
+      }
+    },
+    "node_modules/ansi-regex": {
+      "version": "5.0.1",
+      "resolved": "https://registry.npmjs.org/ansi-regex/-/ansi-regex-5.0.1.tgz",
+      "integrity": "sha512-quJQXlTSUGL2LH9SUXo8VwsY4soanhgo6LNSm84E1LBcE8s3O0wpdiRzyR9z/ZZJMlMWv37qOOb9pdJlMUEKFQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=8"
+      }
+    },
     "node_modules/ansi-styles": {
       "version": "4.3.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "color-convert": "^2.0.1"
       },
       "engines": {
         "node": ">=8"
       },
       "funding": {
         "url": "https://github.com/chalk/ansi-styles?sponsor=1"
       }
     },
     "node_modules/argparse": {
       "version": "2.0.1",
-      "dev": true,
       "license": "Python-2.0"
     },
     "node_modules/aria-query": {
       "version": "5.3.2",
       "dev": true,
       "license": "Apache-2.0",
       "engines": {
         "node": ">= 0.4"
       }
     },
@@ -2431,33 +3179,54 @@
         "get-intrinsic": "^1.2.6",
         "is-array-buffer": "^3.0.4"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/ast-types": {
+      "version": "0.16.1",
+      "resolved": "https://registry.npmjs.org/ast-types/-/ast-types-0.16.1.tgz",
+      "integrity": "sha512-6t10qk83GOG8p0vKmaCr8eiilZwO171AvbROMtvvNiwrTly62t+7XkA8RdIIVbpMhCASAsxgAzdRSwh6nw/5Dg==",
+      "license": "MIT",
+      "dependencies": {
+        "tslib": "^2.0.1"
+      },
+      "engines": {
+        "node": ">=4"
+      }
+    },
     "node_modules/ast-types-flow": {
       "version": "0.0.8",
       "dev": true,
       "license": "MIT"
     },
     "node_modules/async-function": {
       "version": "1.0.0",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       }
     },
+    "node_modules/atomically": {
+      "version": "1.7.0",
+      "resolved": "https://registry.npmjs.org/atomically/-/atomically-1.7.0.tgz",
+      "integrity": "sha512-Xcz9l0z7y9yQ9rdDaxlmaI4uJHf/T8g9hOEzJcsEqX2SjCj4J20uK7+ldkDHMbpJDK76wF7xEIgxc/vSlsfw5w==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=10.12.0"
+      }
+    },
     "node_modules/available-typed-arrays": {
       "version": "1.0.7",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "possible-typed-array-names": "^1.0.0"
       },
       "engines": {
         "node": ">= 0.4"
       },
@@ -2489,43 +3258,78 @@
     "node_modules/baseline-browser-mapping": {
       "version": "2.11.20",
       "license": "Apache-2.0",
       "bin": {
         "baseline-browser-mapping": "dist/cli.cjs"
       },
       "engines": {
         "node": ">=6.0.0"
       }
     },
+    "node_modules/body-parser": {
+      "version": "2.3.0",
+      "resolved": "https://registry.npmjs.org/body-parser/-/body-parser-2.3.0.tgz",
+      "integrity": "sha512-2cGmJupaNgg+QUwVLAucDuWuoMZ6EX9iHDRswZ5lsNYEmwPaRknMPCLZz07yTzVq/83p4o/wzbDZbBrTvGGTIw==",
+      "license": "MIT",
+      "dependencies": {
+        "bytes": "^3.1.2",
+        "content-type": "^2.0.0",
+        "debug": "^4.4.3",
+        "http-errors": "^2.0.1",
+        "iconv-lite": "^0.7.2",
+        "on-finished": "^2.4.1",
+        "qs": "^6.15.2",
+        "raw-body": "^3.0.2",
+        "type-is": "^2.1.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
+    "node_modules/body-parser/node_modules/content-type": {
+      "version": "2.1.0",
+      "resolved": "https://registry.npmjs.org/content-type/-/content-type-2.1.0.tgz",
+      "integrity": "sha512-mj7UPXE0jaqaOsukNZRUEfEi2AcL7C/vwmwcHV0O97eO1E1pxBZuyjlZrx5seTaNBg1U6+o35wpa35Qfcc+7ag==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
     "node_modules/brace-expansion": {
       "version": "1.1.18",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "balanced-match": "^1.0.0",
         "concat-map": "0.0.1"
       }
     },
     "node_modules/braces": {
       "version": "3.0.3",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "fill-range": "^7.1.1"
       },
       "engines": {
         "node": ">=8"
       }
     },
     "node_modules/browserslist": {
       "version": "4.28.8",
-      "dev": true,
       "funding": [
         {
           "type": "opencollective",
           "url": "https://opencollective.com/browserslist"
         },
         {
           "type": "tidelift",
           "url": "https://tidelift.com/funding/github/npm/browserslist"
         },
         {
@@ -2541,151 +3345,409 @@
         "node-releases": "^2.0.53",
         "update-browserslist-db": "^1.3.0"
       },
       "bin": {
         "browserslist": "cli.js"
       },
       "engines": {
         "node": "^6 || ^7 || ^8 || ^9 || ^10 || ^11 || ^12 || >=13.7"
       }
     },
+    "node_modules/bundle-name": {
+      "version": "4.1.0",
+      "resolved": "https://registry.npmjs.org/bundle-name/-/bundle-name-4.1.0.tgz",
+      "integrity": "sha512-tjwM5exMg6BGRI+kNmTntNsvdZS1X8BFYS6tnJ2hdH0kVxM6/eVZ2xy+FqStSWvYmtfFMDLIxurorHwDKfDz5Q==",
+      "license": "MIT",
+      "dependencies": {
+        "run-applescript": "^7.0.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/bytes": {
+      "version": "3.1.2",
+      "resolved": "https://registry.npmjs.org/bytes/-/bytes-3.1.2.tgz",
+      "integrity": "sha512-/Nf7TyzTx6S3yRJObOAV7956r8cr2+Oj8AC5dt8wSP3BQAoeX58NoHyCU8P8zGkNXStjTSi6fzO6F0pBdcYbEg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.8"
+      }
+    },
     "node_modules/call-bind": {
       "version": "1.0.9",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bind-apply-helpers": "^1.0.2",
         "es-define-property": "^1.0.1",
         "get-intrinsic": "^1.3.0",
         "set-function-length": "^1.2.2"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
-        "url": "https://github.com/sponsors/ljharb"
+        "url": "https://github.com/sponsors/ljharb"
+      }
+    },
+    "node_modules/call-bind-apply-helpers": {
+      "version": "1.0.2",
+      "license": "MIT",
+      "dependencies": {
+        "es-errors": "^1.3.0",
+        "function-bind": "^1.1.2"
+      },
+      "engines": {
+        "node": ">= 0.4"
+      }
+    },
+    "node_modules/call-bound": {
+      "version": "1.0.4",
+      "license": "MIT",
+      "dependencies": {
+        "call-bind-apply-helpers": "^1.0.2",
+        "get-intrinsic": "^1.3.0"
+      },
+      "engines": {
+        "node": ">= 0.4"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/ljharb"
+      }
+    },
+    "node_modules/callsites": {
+      "version": "3.1.0",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6"
+      }
+    },
+    "node_modules/caniuse-lite": {
+      "version": "1.0.30001810",
+      "funding": [
+        {
+          "type": "opencollective",
+          "url": "https://opencollective.com/browserslist"
+        },
+        {
+          "type": "tidelift",
+          "url": "https://tidelift.com/funding/github/npm/caniuse-lite"
+        },
+        {
+          "type": "github",
+          "url": "https://github.com/sponsors/ai"
+        }
+      ],
+      "license": "CC-BY-4.0"
+    },
+    "node_modules/chalk": {
+      "version": "4.1.2",
+      "dev": true,
+      "license": "MIT",
+      "dependencies": {
+        "ansi-styles": "^4.1.0",
+        "supports-color": "^7.1.0"
+      },
+      "engines": {
+        "node": ">=10"
+      },
+      "funding": {
+        "url": "https://github.com/chalk/chalk?sponsor=1"
+      }
+    },
+    "node_modules/class-variance-authority": {
+      "version": "0.7.1",
+      "resolved": "https://registry.npmjs.org/class-variance-authority/-/class-variance-authority-0.7.1.tgz",
+      "integrity": "sha512-Ka+9Trutv7G8M6WT6SeiRWz792K5qEqIGEGzXKhAE6xOWAY6pPH8U+9IY3oCMv6kqTmLsv7Xh/2w2RigkePMsg==",
+      "license": "Apache-2.0",
+      "dependencies": {
+        "clsx": "^2.1.1"
+      },
+      "funding": {
+        "url": "https://polar.sh/cva"
+      }
+    },
+    "node_modules/cli-cursor": {
+      "version": "5.0.0",
+      "resolved": "https://registry.npmjs.org/cli-cursor/-/cli-cursor-5.0.0.tgz",
+      "integrity": "sha512-aCj4O5wKyszjMmDT4tZj93kxyydN/K5zPWSCe6/0AV/AA1pqe5ZBIw0a2ZfPQV7lL5/yb5HsUreJ6UFAF1tEQw==",
+      "license": "MIT",
+      "dependencies": {
+        "restore-cursor": "^5.0.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/cli-spinners": {
+      "version": "2.9.2",
+      "resolved": "https://registry.npmjs.org/cli-spinners/-/cli-spinners-2.9.2.tgz",
+      "integrity": "sha512-ywqV+5MmyL4E7ybXgKys4DugZbX0FC6LnwrhjuykIjnK9k8OQacQ7axGKnjDXWNhns0xot3bZI5h55H8yo9cJg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/client-only": {
+      "version": "0.0.1",
+      "license": "MIT"
+    },
+    "node_modules/clsx": {
+      "version": "2.1.1",
+      "resolved": "https://registry.npmjs.org/clsx/-/clsx-2.1.1.tgz",
+      "integrity": "sha512-eYm0QWBtUrBWZWG0d386OGAw16Z995PiOVo2B7bjWSbHedGl5e0ZWaq65kOGgUSNesEIDkB9ISbTg/JK9dhCZA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6"
+      }
+    },
+    "node_modules/code-block-writer": {
+      "version": "13.0.3",
+      "resolved": "https://registry.npmjs.org/code-block-writer/-/code-block-writer-13.0.3.tgz",
+      "integrity": "sha512-Oofo0pq3IKnsFtuHqSF7TqBfr71aeyZDVJ0HpmqB7FBM2qEigL0iPONSCZSO9pE9dZTAxANe5XHG9Uy0YMv8cg==",
+      "license": "MIT"
+    },
+    "node_modules/color-convert": {
+      "version": "2.0.1",
+      "dev": true,
+      "license": "MIT",
+      "dependencies": {
+        "color-name": "~1.1.4"
+      },
+      "engines": {
+        "node": ">=7.0.0"
+      }
+    },
+    "node_modules/color-name": {
+      "version": "1.1.4",
+      "dev": true,
+      "license": "MIT"
+    },
+    "node_modules/commander": {
+      "version": "14.0.3",
+      "resolved": "https://registry.npmjs.org/commander/-/commander-14.0.3.tgz",
+      "integrity": "sha512-H+y0Jo/T1RZ9qPP4Eh1pkcQcLRglraJaSLoyOtHxu6AapkjWVCy2Sit1QQ4x3Dng8qDlSsZEet7g5Pq06MvTgw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=20"
+      }
+    },
+    "node_modules/concat-map": {
+      "version": "0.0.1",
+      "dev": true,
+      "license": "MIT"
+    },
+    "node_modules/conf": {
+      "version": "10.2.0",
+      "resolved": "https://registry.npmjs.org/conf/-/conf-10.2.0.tgz",
+      "integrity": "sha512-8fLl9F04EJqjSqH+QjITQfJF8BrOVaYr1jewVgSRAEWePfxT0sku4w2hrGQ60BC/TNLGQ2pgxNlTbWQmMPFvXg==",
+      "license": "MIT",
+      "dependencies": {
+        "ajv": "^8.6.3",
+        "ajv-formats": "^2.1.1",
+        "atomically": "^1.7.0",
+        "debounce-fn": "^4.0.0",
+        "dot-prop": "^6.0.1",
+        "env-paths": "^2.2.1",
+        "json-schema-typed": "^7.0.3",
+        "onetime": "^5.1.2",
+        "pkg-up": "^3.1.0",
+        "semver": "^7.3.5"
+      },
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/conf/node_modules/ajv": {
+      "version": "8.20.0",
+      "resolved": "https://registry.npmjs.org/ajv/-/ajv-8.20.0.tgz",
+      "integrity": "sha512-Thbli+OlOj+iMPYFBVBfJ3OmCAnaSyNn4M1vz9T6Gka5Jt9ba/HIR56joy65tY6kx/FCF5VXNB819Y7/GUrBGA==",
+      "license": "MIT",
+      "dependencies": {
+        "fast-deep-equal": "^3.1.3",
+        "fast-uri": "^3.0.1",
+        "json-schema-traverse": "^1.0.0",
+        "require-from-string": "^2.0.2"
+      },
+      "funding": {
+        "type": "github",
+        "url": "https://github.com/sponsors/epoberezkin"
+      }
+    },
+    "node_modules/conf/node_modules/ajv-formats": {
+      "version": "2.1.1",
+      "resolved": "https://registry.npmjs.org/ajv-formats/-/ajv-formats-2.1.1.tgz",
+      "integrity": "sha512-Wx0Kx52hxE7C18hkMEggYlEifqWZtYaRgouJor+WMdPnQyEK13vgEWyVNup7SoeeoLMsr4kf5h6dOW11I15MUA==",
+      "license": "MIT",
+      "dependencies": {
+        "ajv": "^8.0.0"
+      },
+      "peerDependencies": {
+        "ajv": "^8.0.0"
+      },
+      "peerDependenciesMeta": {
+        "ajv": {
+          "optional": true
+        }
+      }
+    },
+    "node_modules/conf/node_modules/json-schema-traverse": {
+      "version": "1.0.0",
+      "resolved": "https://registry.npmjs.org/json-schema-traverse/-/json-schema-traverse-1.0.0.tgz",
+      "integrity": "sha512-NM8/P9n3XjXhIZn1lLhkFaACTOURQXjWhV4BA/RnOv8xvgqtqpAX9IO4mRQxSx1Rlo4tqzeqb0sOlruaOy3dug==",
+      "license": "MIT"
+    },
+    "node_modules/conf/node_modules/json-schema-typed": {
+      "version": "7.0.3",
+      "resolved": "https://registry.npmjs.org/json-schema-typed/-/json-schema-typed-7.0.3.tgz",
+      "integrity": "sha512-7DE8mpG+/fVw+dTpjbxnx47TaMnDfOI1jwft9g1VybltZCduyRQPJPvc+zzKY9WPHxhPWczyFuYa6I8Mw4iU5A==",
+      "license": "BSD-2-Clause"
+    },
+    "node_modules/conf/node_modules/semver": {
+      "version": "7.8.5",
+      "resolved": "https://registry.npmjs.org/semver/-/semver-7.8.5.tgz",
+      "integrity": "sha512-Y7/KDsb8LjooZpwaqGyulO6DQlksgCncchHGk+sZIY4SBvUocMBEFH5Ur1fI4dV+Jvl0w6cjvucaIi40puRioA==",
+      "license": "ISC",
+      "bin": {
+        "semver": "bin/semver.js"
+      },
+      "engines": {
+        "node": ">=10"
+      }
+    },
+    "node_modules/content-disposition": {
+      "version": "1.1.0",
+      "resolved": "https://registry.npmjs.org/content-disposition/-/content-disposition-1.1.0.tgz",
+      "integrity": "sha512-5jRCH9Z/+DRP7rkvY83B+yGIGX96OYdJmzngqnw2SBSxqCFPd0w2km3s5iawpGX8krnwSGmF0FW5Nhr0Hfai3g==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
       }
     },
-    "node_modules/call-bind-apply-helpers": {
-      "version": "1.0.2",
-      "dev": true,
+    "node_modules/content-type": {
+      "version": "1.0.5",
+      "resolved": "https://registry.npmjs.org/content-type/-/content-type-1.0.5.tgz",
+      "integrity": "sha512-nTjqfcBFEipKdXCv4YDQWCfmcLZKm81ldF0pAopTvyrFGVbcR6P/VAAd5G7N+0tTr8QqiU0tFadD6FK4NtJwOA==",
       "license": "MIT",
-      "dependencies": {
-        "es-errors": "^1.3.0",
-        "function-bind": "^1.1.2"
-      },
       "engines": {
-        "node": ">= 0.4"
+        "node": ">= 0.6"
       }
     },
-    "node_modules/call-bound": {
-      "version": "1.0.4",
-      "dev": true,
+    "node_modules/convert-source-map": {
+      "version": "2.0.0",
+      "license": "MIT"
+    },
+    "node_modules/cookie": {
+      "version": "0.7.2",
+      "resolved": "https://registry.npmjs.org/cookie/-/cookie-0.7.2.tgz",
+      "integrity": "sha512-yki5XnKuf750l50uGTllt6kKILY4nQ1eNIQatoXEByZ5dWgnKqbnqmTrBE5B4N7lrMJKQ2ytWMiTO2o0v6Ew/w==",
       "license": "MIT",
-      "dependencies": {
-        "call-bind-apply-helpers": "^1.0.2",
-        "get-intrinsic": "^1.3.0"
-      },
       "engines": {
-        "node": ">= 0.4"
-      },
-      "funding": {
-        "url": "https://github.com/sponsors/ljharb"
+        "node": ">= 0.6"
       }
     },
-    "node_modules/callsites": {
-      "version": "3.1.0",
-      "dev": true,
+    "node_modules/cookie-signature": {
+      "version": "1.2.2",
+      "resolved": "https://registry.npmjs.org/cookie-signature/-/cookie-signature-1.2.2.tgz",
+      "integrity": "sha512-D76uU73ulSXrD1UXF4KE2TMxVVwhsnCgfAyTg9k8P6KGZjlXKrOLe4dJQKI3Bxi5wjesZoFXJWElNWBjPZMbhg==",
       "license": "MIT",
       "engines": {
-        "node": ">=6"
+        "node": ">=6.6.0"
       }
     },
-    "node_modules/caniuse-lite": {
-      "version": "1.0.30001810",
-      "funding": [
-        {
-          "type": "opencollective",
-          "url": "https://opencollective.com/browserslist"
-        },
-        {
-          "type": "tidelift",
-          "url": "https://tidelift.com/funding/github/npm/caniuse-lite"
-        },
-        {
-          "type": "github",
-          "url": "https://github.com/sponsors/ai"
-        }
-      ],
-      "license": "CC-BY-4.0"
-    },
-    "node_modules/chalk": {
-      "version": "4.1.2",
-      "dev": true,
+    "node_modules/cors": {
+      "version": "2.8.6",
+      "resolved": "https://registry.npmjs.org/cors/-/cors-2.8.6.tgz",
+      "integrity": "sha512-tJtZBBHA6vjIAaF6EnIaq6laBBP9aq/Y3ouVJjEfoHbRBcHBAHYcMh/w8LDrk2PvIMMq8gmopa5D4V8RmbrxGw==",
       "license": "MIT",
       "dependencies": {
-        "ansi-styles": "^4.1.0",
-        "supports-color": "^7.1.0"
+        "object-assign": "^4",
+        "vary": "^1"
       },
       "engines": {
-        "node": ">=10"
+        "node": ">= 0.10"
       },
       "funding": {
-        "url": "https://github.com/chalk/chalk?sponsor=1"
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
       }
     },
-    "node_modules/client-only": {
-      "version": "0.0.1",
-      "license": "MIT"
-    },
-    "node_modules/color-convert": {
-      "version": "2.0.1",
-      "dev": true,
+    "node_modules/cosmiconfig": {
+      "version": "9.0.2",
+      "resolved": "https://registry.npmjs.org/cosmiconfig/-/cosmiconfig-9.0.2.tgz",
+      "integrity": "sha512-gtTZxTDau1wL7Y7zifc2dd8jHSK/k6BTx/2Xp/BpdlAdnlYWFVt7qhJqgwi7637yRwRQ3qL4ZidbB4I8tA5VOg==",
       "license": "MIT",
       "dependencies": {
-        "color-name": "~1.1.4"
+        "env-paths": "^2.2.1",
+        "import-fresh": "^3.3.0",
+        "js-yaml": "^4.1.0",
+        "parse-json": "^5.2.0"
       },
       "engines": {
-        "node": ">=7.0.0"
+        "node": ">=14"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/d-fischer"
+      },
+      "peerDependencies": {
+        "typescript": ">=4.9.5"
+      },
+      "peerDependenciesMeta": {
+        "typescript": {
+          "optional": true
+        }
       }
     },
-    "node_modules/color-name": {
-      "version": "1.1.4",
-      "dev": true,
-      "license": "MIT"
-    },
-    "node_modules/concat-map": {
-      "version": "0.0.1",
-      "dev": true,
-      "license": "MIT"
-    },
-    "node_modules/convert-source-map": {
-      "version": "2.0.0",
-      "dev": true,
-      "license": "MIT"
-    },
     "node_modules/cross-spawn": {
       "version": "7.0.6",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "path-key": "^3.1.0",
         "shebang-command": "^2.0.0",
         "which": "^2.0.1"
       },
       "engines": {
         "node": ">= 8"
       }
     },
+    "node_modules/cssesc": {
+      "version": "3.0.0",
+      "resolved": "https://registry.npmjs.org/cssesc/-/cssesc-3.0.0.tgz",
+      "integrity": "sha512-/Tb/JcjK111nNScGob5MNtsntNM1aCNUDipB/TkwZFhyDrrE47SOx/18wF2bbjgc3ZzCSKW1T5nt5EbFoAz/Vg==",
+      "license": "MIT",
+      "bin": {
+        "cssesc": "bin/cssesc"
+      },
+      "engines": {
+        "node": ">=4"
+      }
+    },
     "node_modules/csstype": {
       "version": "3.2.3",
-      "dev": true,
+      "devOptional": true,
       "license": "MIT"
     },
     "node_modules/damerau-levenshtein": {
       "version": "1.0.8",
       "dev": true,
       "license": "BSD-2-Clause"
     },
     "node_modules/data-view-buffer": {
       "version": "1.0.2",
       "dev": true,
@@ -2727,127 +3789,293 @@
         "es-errors": "^1.3.0",
         "is-data-view": "^1.0.1"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/debounce-fn": {
+      "version": "4.0.0",
+      "resolved": "https://registry.npmjs.org/debounce-fn/-/debounce-fn-4.0.0.tgz",
+      "integrity": "sha512-8pYCQiL9Xdcg0UPSD3d+0KMlOjp+KGU5EPwYddgzQ7DATsg4fuUDjQtsYLmWjnk2obnNHgV3vE2Y4jejSOJVBQ==",
+      "license": "MIT",
+      "dependencies": {
+        "mimic-fn": "^3.0.0"
+      },
+      "engines": {
+        "node": ">=10"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/debug": {
       "version": "4.4.3",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "ms": "^2.1.3"
       },
       "engines": {
         "node": ">=6.0"
       },
       "peerDependenciesMeta": {
         "supports-color": {
           "optional": true
         }
       }
     },
+    "node_modules/dedent": {
+      "version": "1.7.2",
+      "resolved": "https://registry.npmjs.org/dedent/-/dedent-1.7.2.tgz",
+      "integrity": "sha512-WzMx3mW98SN+zn3hgemf4OzdmyNhhhKz5Ay0pUfQiMQ3e1g+xmTJWp/pKdwKVXhdSkAEGIIzqeuWrL3mV/AXbA==",
+      "license": "MIT",
+      "peerDependencies": {
+        "babel-plugin-macros": "^3.1.0"
+      },
+      "peerDependenciesMeta": {
+        "babel-plugin-macros": {
+          "optional": true
+        }
+      }
+    },
     "node_modules/deep-is": {
       "version": "0.1.4",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/deepmerge": {
+      "version": "4.3.1",
+      "resolved": "https://registry.npmjs.org/deepmerge/-/deepmerge-4.3.1.tgz",
+      "integrity": "sha512-3sUqbMEc77XqpdNO7FRyRog+eW3ph+GYCbj+rK+uYyRMuwsVy0rMiVtPn+QJlKFvWP/1PYpapqYn0Me2knFn+A==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=0.10.0"
+      }
+    },
+    "node_modules/default-browser": {
+      "version": "5.5.1",
+      "resolved": "https://registry.npmjs.org/default-browser/-/default-browser-5.5.1.tgz",
+      "integrity": "sha512-m1pAzaJgZ/gssEqlOhJkPJp8Xly7QyW6xcrkUa2KKcDeDSEMP7X8xipU3snUcfisTQx0w1AGae+9UtJSfVnXGw==",
+      "license": "MIT",
+      "dependencies": {
+        "bundle-name": "^4.1.0",
+        "default-browser-id": "^5.0.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/default-browser-id": {
+      "version": "5.0.1",
+      "resolved": "https://registry.npmjs.org/default-browser-id/-/default-browser-id-5.0.1.tgz",
+      "integrity": "sha512-x1VCxdX4t+8wVfd1so/9w+vQ4vx7lKd2Qp5tDRutErwmR85OgmfX7RlLRMWafRMY7hbEiXIbudNrjOAPa/hL8Q==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/define-data-property": {
       "version": "1.1.4",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "es-define-property": "^1.0.0",
         "es-errors": "^1.3.0",
         "gopd": "^1.0.1"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/define-lazy-prop": {
+      "version": "3.0.0",
+      "resolved": "https://registry.npmjs.org/define-lazy-prop/-/define-lazy-prop-3.0.0.tgz",
+      "integrity": "sha512-N+MeXYoqr3pOgn8xfyRPREN7gHakLYjhsHhWGT3fWAiL4IkAt0iDw14QiiEm2bE30c5XX5q0FtAA3CK5f9/BUg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/define-properties": {
       "version": "1.2.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "define-data-property": "^1.0.1",
         "has-property-descriptors": "^1.0.0",
         "object-keys": "^1.1.1"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/depd": {
+      "version": "2.0.0",
+      "resolved": "https://registry.npmjs.org/depd/-/depd-2.0.0.tgz",
+      "integrity": "sha512-g7nH6P6dyDioJogAAGprGpCtVImJhpPk/roCzdb3fIh61/s/nPsfR6onyMwkCAR/OlC3yBC0lESvUoQEAssIrw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.8"
+      }
+    },
     "node_modules/detect-libc": {
       "version": "2.1.2",
       "devOptional": true,
       "license": "Apache-2.0",
       "engines": {
         "node": ">=8"
       }
     },
+    "node_modules/diff": {
+      "version": "8.0.4",
+      "resolved": "https://registry.npmjs.org/diff/-/diff-8.0.4.tgz",
+      "integrity": "sha512-DPi0FmjiSU5EvQV0++GFDOJ9ASQUVFh5kD+OzOnYdi7n3Wpm9hWWGfB/O2blfHcMVTL5WkQXSnRiK9makhrcnw==",
+      "license": "BSD-3-Clause",
+      "engines": {
+        "node": ">=0.3.1"
+      }
+    },
     "node_modules/doctrine": {
       "version": "2.1.0",
       "dev": true,
       "license": "Apache-2.0",
       "dependencies": {
         "esutils": "^2.0.2"
       },
       "engines": {
         "node": ">=0.10.0"
       }
     },
+    "node_modules/dot-prop": {
+      "version": "6.0.1",
+      "resolved": "https://registry.npmjs.org/dot-prop/-/dot-prop-6.0.1.tgz",
+      "integrity": "sha512-tE7ztYzXHIeyvc7N+hR3oi7FIbf/NIjVP9hmAt3yMXzrQ072/fpjGLx2GxNxGxUl5V73MEqYzioOMoVhGMJ5cA==",
+      "license": "MIT",
+      "dependencies": {
+        "is-obj": "^2.0.0"
+      },
+      "engines": {
+        "node": ">=10"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/dotenv": {
+      "version": "17.4.2",
+      "resolved": "https://registry.npmjs.org/dotenv/-/dotenv-17.4.2.tgz",
+      "integrity": "sha512-nI4U3TottKAcAD9LLud4Cb7b2QztQMUEfHbvhTH09bqXTxnSie8WnjPALV/WMCrJZ6UV/qHJ6L03OqO3LcdYZw==",
+      "license": "BSD-2-Clause",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://dotenvx.com"
+      }
+    },
     "node_modules/dunder-proto": {
       "version": "1.0.1",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bind-apply-helpers": "^1.0.1",
         "es-errors": "^1.3.0",
         "gopd": "^1.2.0"
       },
       "engines": {
         "node": ">= 0.4"
       }
     },
+    "node_modules/ee-first": {
+      "version": "1.1.1",
+      "resolved": "https://registry.npmjs.org/ee-first/-/ee-first-1.1.1.tgz",
+      "integrity": "sha512-WMwm9LhRUo+WUaRN+vRuETqG89IgZphVSNkdFgeb6sS/E4OrDIN7t48CAewSHXc6C8lefD8KKfr5vY61brQlow==",
+      "license": "MIT"
+    },
     "node_modules/electron-to-chromium": {
       "version": "1.5.416",
-      "dev": true,
       "license": "ISC"
     },
     "node_modules/emoji-regex": {
       "version": "9.2.2",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/encodeurl": {
+      "version": "2.0.0",
+      "resolved": "https://registry.npmjs.org/encodeurl/-/encodeurl-2.0.0.tgz",
+      "integrity": "sha512-Q0n9HRi4m6JuGIV1eFlmvJB7ZEVxu93IrMyiMsGC0lrMJMWzRgx6WGquyfQgZVb31vhGgXnfmPNNXmxnOkRBrg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.8"
+      }
+    },
     "node_modules/enhanced-resolve": {
       "version": "5.24.5",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "graceful-fs": "^4.2.4",
         "tapable": "^2.3.3"
       },
       "engines": {
         "node": ">=10.13.0"
       }
     },
+    "node_modules/enquirer": {
+      "version": "2.4.1",
+      "resolved": "https://registry.npmjs.org/enquirer/-/enquirer-2.4.1.tgz",
+      "integrity": "sha512-rRqJg/6gd538VHvR3PSrdRBb/1Vy2YfzHqzvbhGIQpDRKIa4FgV/54b5Q1xYSxOOwKvjXweS26E0Q+nAMwp2pQ==",
+      "license": "MIT",
+      "dependencies": {
+        "ansi-colors": "^4.1.1",
+        "strip-ansi": "^6.0.1"
+      },
+      "engines": {
+        "node": ">=8.6"
+      }
+    },
+    "node_modules/env-paths": {
+      "version": "2.2.1",
+      "resolved": "https://registry.npmjs.org/env-paths/-/env-paths-2.2.1.tgz",
+      "integrity": "sha512-+h1lkLKhZMTYjog1VEpJNG7NZJWcuc2DDk/qsqSTRRCOXiLjeQ1d1/udrUGhqMxUgAlwKNZ0cf2uqan5GLuS2A==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6"
+      }
+    },
+    "node_modules/error-ex": {
+      "version": "1.3.4",
+      "resolved": "https://registry.npmjs.org/error-ex/-/error-ex-1.3.4.tgz",
+      "integrity": "sha512-sqQamAnR14VgCr1A618A3sGrygcpK+HEbenA/HiEAkkUwcZIIB/tgWqHFxWgOyDh4nB4JCRimh79dR5Ywc9MDQ==",
+      "license": "MIT",
+      "dependencies": {
+        "is-arrayish": "^0.2.1"
+      }
+    },
     "node_modules/es-abstract": {
       "version": "1.24.2",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "array-buffer-byte-length": "^1.0.2",
         "arraybuffer.prototype.slice": "^1.0.4",
         "available-typed-arrays": "^1.0.7",
         "call-bind": "^1.0.8",
         "call-bound": "^1.0.4",
@@ -2920,29 +4148,27 @@
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/es-define-property": {
       "version": "1.0.1",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       }
     },
     "node_modules/es-errors": {
       "version": "1.3.0",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       }
     },
     "node_modules/es-iterator-helpers": {
       "version": "1.4.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
@@ -2962,21 +4188,20 @@
         "internal-slot": "^1.1.0",
         "iterator.prototype": "^1.1.5",
         "math-intrinsics": "^1.1.0"
       },
       "engines": {
         "node": ">= 0.4"
       }
     },
     "node_modules/es-object-atoms": {
       "version": "1.1.2",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "es-errors": "^1.3.0"
       },
       "engines": {
         "node": ">= 0.4"
       }
     },
     "node_modules/es-set-tostringtag": {
       "version": "2.1.0",
@@ -3017,26 +4242,31 @@
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/escalade": {
       "version": "3.2.0",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=6"
       }
     },
+    "node_modules/escape-html": {
+      "version": "1.0.3",
+      "resolved": "https://registry.npmjs.org/escape-html/-/escape-html-1.0.3.tgz",
+      "integrity": "sha512-NiSupZ4OeuGwr68lGIeym/ksIZMJodUGOSCZ/FSnTxcrekbvqrgdUxlJOMpijaKZVjAJrWrGs/6Jy8OMuyj9ow==",
+      "license": "MIT"
+    },
     "node_modules/escape-string-regexp": {
       "version": "4.0.0",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=10"
       },
       "funding": {
         "url": "https://github.com/sponsors/sindresorhus"
       }
@@ -3362,20 +4592,33 @@
         "acorn-jsx": "^5.3.2",
         "eslint-visitor-keys": "^4.2.1"
       },
       "engines": {
         "node": "^18.18.0 || ^20.9.0 || >=21.1.0"
       },
       "funding": {
         "url": "https://opencollective.com/eslint"
       }
     },
+    "node_modules/esprima": {
+      "version": "4.0.1",
+      "resolved": "https://registry.npmjs.org/esprima/-/esprima-4.0.1.tgz",
+      "integrity": "sha512-eGuFFw7Upda+g4p+QHvnW0RyTX/SVeJBDM/gCtMARO0cLuT2HcEKnTPvhjV6aGeqrCB/sbNop0Kszm0jsaWU4A==",
+      "license": "BSD-2-Clause",
+      "bin": {
+        "esparse": "bin/esparse.js",
+        "esvalidate": "bin/esvalidate.js"
+      },
+      "engines": {
+        "node": ">=4"
+      }
+    },
     "node_modules/esquery": {
       "version": "1.7.0",
       "dev": true,
       "license": "BSD-3-Clause",
       "dependencies": {
         "estraverse": "^5.1.0"
       },
       "engines": {
         "node": ">=0.10"
       }
@@ -3392,31 +4635,148 @@
       }
     },
     "node_modules/estraverse": {
       "version": "5.3.0",
       "dev": true,
       "license": "BSD-2-Clause",
       "engines": {
         "node": ">=4.0"
       }
     },
-    "node_modules/esutils": {
-      "version": "2.0.3",
-      "dev": true,
-      "license": "BSD-2-Clause",
+    "node_modules/esutils": {
+      "version": "2.0.3",
+      "dev": true,
+      "license": "BSD-2-Clause",
+      "engines": {
+        "node": ">=0.10.0"
+      }
+    },
+    "node_modules/etag": {
+      "version": "1.8.1",
+      "resolved": "https://registry.npmjs.org/etag/-/etag-1.8.1.tgz",
+      "integrity": "sha512-aIL5Fx7mawVa300al2BnEE4iNvo1qETxLrPI/o05L7z6go7fCw1J6EQmbK4FmJ2AS7kgVF/KEZWufBfdClMcPg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.6"
+      }
+    },
+    "node_modules/eventsource": {
+      "version": "3.0.7",
+      "resolved": "https://registry.npmjs.org/eventsource/-/eventsource-3.0.7.tgz",
+      "integrity": "sha512-CRT1WTyuQoD771GW56XEZFQ/ZoSfWid1alKGDYMmkt2yl8UXrVR4pspqWNEcqKvVIzg6PAltWjxcSSPrboA4iA==",
+      "license": "MIT",
+      "dependencies": {
+        "eventsource-parser": "^3.0.1"
+      },
+      "engines": {
+        "node": ">=18.0.0"
+      }
+    },
+    "node_modules/eventsource-parser": {
+      "version": "3.1.1",
+      "resolved": "https://registry.npmjs.org/eventsource-parser/-/eventsource-parser-3.1.1.tgz",
+      "integrity": "sha512-EKN1vKAMcZ8MlYMpaNuxN6R9yakzH6uajHcHVTqWJzvu5pWw9DyhbP35HH8MVBQ+dZjAfDxk+A8NiR9KWaXiyQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18.0.0"
+      }
+    },
+    "node_modules/execa": {
+      "version": "9.6.1",
+      "resolved": "https://registry.npmjs.org/execa/-/execa-9.6.1.tgz",
+      "integrity": "sha512-9Be3ZoN4LmYR90tUoVu2te2BsbzHfhJyfEiAVfz7N5/zv+jduIfLrV2xdQXOHbaD6KgpGdO9PRPM1Y4Q9QkPkA==",
+      "license": "MIT",
+      "dependencies": {
+        "@sindresorhus/merge-streams": "^4.0.0",
+        "cross-spawn": "^7.0.6",
+        "figures": "^6.1.0",
+        "get-stream": "^9.0.0",
+        "human-signals": "^8.0.1",
+        "is-plain-obj": "^4.1.0",
+        "is-stream": "^4.0.1",
+        "npm-run-path": "^6.0.0",
+        "pretty-ms": "^9.2.0",
+        "signal-exit": "^4.1.0",
+        "strip-final-newline": "^4.0.0",
+        "yoctocolors": "^2.1.1"
+      },
+      "engines": {
+        "node": "^18.19.0 || >=20.5.0"
+      },
+      "funding": {
+        "url": "https://github.com/sindresorhus/execa?sponsor=1"
+      }
+    },
+    "node_modules/express": {
+      "version": "5.2.1",
+      "resolved": "https://registry.npmjs.org/express/-/express-5.2.1.tgz",
+      "integrity": "sha512-hIS4idWWai69NezIdRt2xFVofaF4j+6INOpJlVOLDO8zXGpUVEVzIYk12UUi2JzjEzWL3IOAxcTubgz9Po0yXw==",
+      "license": "MIT",
+      "dependencies": {
+        "accepts": "^2.0.0",
+        "body-parser": "^2.2.1",
+        "content-disposition": "^1.0.0",
+        "content-type": "^1.0.5",
+        "cookie": "^0.7.1",
+        "cookie-signature": "^1.2.1",
+        "debug": "^4.4.0",
+        "depd": "^2.0.0",
+        "encodeurl": "^2.0.0",
+        "escape-html": "^1.0.3",
+        "etag": "^1.8.1",
+        "finalhandler": "^2.1.0",
+        "fresh": "^2.0.0",
+        "http-errors": "^2.0.0",
+        "merge-descriptors": "^2.0.0",
+        "mime-types": "^3.0.0",
+        "on-finished": "^2.4.1",
+        "once": "^1.4.0",
+        "parseurl": "^1.3.3",
+        "proxy-addr": "^2.0.7",
+        "qs": "^6.14.0",
+        "range-parser": "^1.2.1",
+        "router": "^2.2.0",
+        "send": "^1.1.0",
+        "serve-static": "^2.2.0",
+        "statuses": "^2.0.1",
+        "type-is": "^2.0.1",
+        "vary": "^1.1.2"
+      },
+      "engines": {
+        "node": ">= 18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
+    "node_modules/express-rate-limit": {
+      "version": "8.7.0",
+      "resolved": "https://registry.npmjs.org/express-rate-limit/-/express-rate-limit-8.7.0.tgz",
+      "integrity": "sha512-hOwV7WOxXfjRpAM1DSJWZDXx3GhplwD8IfwuwvogD8i1Qnkgosw/H45s4ZnFAUHDAhPjlY9hLBvJhKmGMyY26g==",
+      "license": "MIT",
+      "dependencies": {
+        "debug": "^4.4.3",
+        "ip-address": "^10.2.0"
+      },
       "engines": {
-        "node": ">=0.10.0"
+        "node": ">= 16"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/express-rate-limit"
+      },
+      "peerDependencies": {
+        "express": ">= 4.11"
       }
     },
     "node_modules/fast-deep-equal": {
       "version": "3.1.3",
-      "dev": true,
       "license": "MIT"
     },
     "node_modules/fast-glob": {
       "version": "3.3.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "@nodelib/fs.stat": "^2.0.2",
         "@nodelib/fs.walk": "^1.2.3",
         "glob-parent": "^5.1.2",
@@ -3441,50 +4801,100 @@
     "node_modules/fast-json-stable-stringify": {
       "version": "2.1.0",
       "dev": true,
       "license": "MIT"
     },
     "node_modules/fast-levenshtein": {
       "version": "2.0.6",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/fast-uri": {
+      "version": "3.1.6",
+      "resolved": "https://registry.npmjs.org/fast-uri/-/fast-uri-3.1.6.tgz",
+      "integrity": "sha512-7Ical1vFEMr0onbVzEDIreM22I4khW+fzyQPwvAFWBp1iwdshSZRsL4jjRvPG9JP1uiqMHRto+YU6R2/CzDz5Q==",
+      "funding": [
+        {
+          "type": "github",
+          "url": "https://github.com/sponsors/fastify"
+        },
+        {
+          "type": "opencollective",
+          "url": "https://opencollective.com/fastify"
+        }
+      ],
+      "license": "BSD-3-Clause"
+    },
     "node_modules/fastq": {
       "version": "1.20.3",
-      "dev": true,
       "license": "ISC",
       "dependencies": {
         "reusify": "^1.0.4"
       }
     },
+    "node_modules/figures": {
+      "version": "6.1.0",
+      "resolved": "https://registry.npmjs.org/figures/-/figures-6.1.0.tgz",
+      "integrity": "sha512-d+l3qxjSesT4V7v2fh+QnmFnUWv9lSpjarhShNTgBOfA0ttejbQUAlHLitbjkoRiDulW0OPoQPYIGhIC8ohejg==",
+      "license": "MIT",
+      "dependencies": {
+        "is-unicode-supported": "^2.0.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/file-entry-cache": {
       "version": "8.0.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "flat-cache": "^4.0.0"
       },
       "engines": {
         "node": ">=16.0.0"
       }
     },
     "node_modules/fill-range": {
       "version": "7.1.1",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "to-regex-range": "^5.0.1"
       },
       "engines": {
         "node": ">=8"
       }
     },
+    "node_modules/finalhandler": {
+      "version": "2.1.1",
+      "resolved": "https://registry.npmjs.org/finalhandler/-/finalhandler-2.1.1.tgz",
+      "integrity": "sha512-S8KoZgRZN+a5rNwqTxlZZePjT/4cnm0ROV70LedRHZ0p8u9fRID0hJUZQpkKLzro8LfmC8sx23bY6tVNxv8pQA==",
+      "license": "MIT",
+      "dependencies": {
+        "debug": "^4.4.0",
+        "encodeurl": "^2.0.0",
+        "escape-html": "^1.0.3",
+        "on-finished": "^2.4.1",
+        "parseurl": "^1.3.3",
+        "statuses": "^2.0.1"
+      },
+      "engines": {
+        "node": ">= 18.0.0"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
     "node_modules/find-up": {
       "version": "5.0.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "locate-path": "^6.0.0",
         "path-exists": "^4.0.0"
       },
       "engines": {
         "node": ">=10"
@@ -3517,23 +4927,54 @@
       "dependencies": {
         "is-callable": "^1.2.7"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/forwarded": {
+      "version": "0.2.0",
+      "resolved": "https://registry.npmjs.org/forwarded/-/forwarded-0.2.0.tgz",
+      "integrity": "sha512-buRG0fpBtRHSTCOASe6hD258tEubFoRLb4ZNA6NxMVHNw2gOcwHo9wyablzMzOA5z9xA9L1KNjk/Nt6MT9aYow==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.6"
+      }
+    },
+    "node_modules/fresh": {
+      "version": "2.0.0",
+      "resolved": "https://registry.npmjs.org/fresh/-/fresh-2.0.0.tgz",
+      "integrity": "sha512-Rx/WycZ60HOaqLKAi6cHRKKI7zxWbJ31MhntmtwMoaTeF7XFH9hhBp8vITaMidfljRQ6eYWCKkaTK+ykVJHP2A==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.8"
+      }
+    },
+    "node_modules/fs-extra": {
+      "version": "11.4.0",
+      "resolved": "https://registry.npmjs.org/fs-extra/-/fs-extra-11.4.0.tgz",
+      "integrity": "sha512-EQsFzMUJkCKGr1ePqlYADkIUmHW1s3ZXr5Yqy6wbGrfUCphpl2maM/kyOIRA2HpP3AaFQTZXD4ldjek+nccddA==",
+      "license": "MIT",
+      "dependencies": {
+        "graceful-fs": "^4.2.0",
+        "jsonfile": "^6.0.1",
+        "universalify": "^2.0.0"
+      },
+      "engines": {
+        "node": ">=14.14"
+      }
+    },
     "node_modules/function-bind": {
       "version": "1.1.2",
-      "dev": true,
       "license": "MIT",
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/function.prototype.name": {
       "version": "1.2.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
@@ -3555,71 +4996,114 @@
       }
     },
     "node_modules/functions-have-names": {
       "version": "1.2.3",
       "dev": true,
       "license": "MIT",
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/fuzzysort": {
+      "version": "3.1.0",
+      "resolved": "https://registry.npmjs.org/fuzzysort/-/fuzzysort-3.1.0.tgz",
+      "integrity": "sha512-sR9BNCjBg6LNgwvxlBd0sBABvQitkLzoVY9MYYROQVX/FvfJ4Mai9LsGhDgd8qYdds0bY77VzYd5iuB+v5rwQQ==",
+      "license": "MIT"
+    },
     "node_modules/generator-function": {
       "version": "2.0.1",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       }
     },
     "node_modules/gensync": {
       "version": "1.0.0-beta.2",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=6.9.0"
       }
     },
+    "node_modules/get-east-asian-width": {
+      "version": "1.6.0",
+      "resolved": "https://registry.npmjs.org/get-east-asian-width/-/get-east-asian-width-1.6.0.tgz",
+      "integrity": "sha512-QRbvDIbx6YklUe6RxeTeleMR0yv3cYH6PsPZHcnVn7xv7zO1BHN8r0XETu8n6Ye3Q+ahtSarc3WgtNWmehIBfA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/get-intrinsic": {
       "version": "1.3.0",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bind-apply-helpers": "^1.0.2",
         "es-define-property": "^1.0.1",
         "es-errors": "^1.3.0",
         "es-object-atoms": "^1.1.1",
         "function-bind": "^1.1.2",
         "get-proto": "^1.0.1",
         "gopd": "^1.2.0",
         "has-symbols": "^1.1.0",
         "hasown": "^2.0.2",
         "math-intrinsics": "^1.1.0"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/get-own-enumerable-keys": {
+      "version": "1.0.0",
+      "resolved": "https://registry.npmjs.org/get-own-enumerable-keys/-/get-own-enumerable-keys-1.0.0.tgz",
+      "integrity": "sha512-PKsK2FSrQCyxcGHsGrLDcK0lx+0Ke+6e8KFFozA9/fIQLhQzPaRvJFdcz7+Axg3jUH/Mq+NI4xa5u/UT2tQskA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=14.16"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/get-proto": {
       "version": "1.0.1",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "dunder-proto": "^1.0.1",
         "es-object-atoms": "^1.0.0"
       },
       "engines": {
         "node": ">= 0.4"
       }
     },
+    "node_modules/get-stream": {
+      "version": "9.0.1",
+      "resolved": "https://registry.npmjs.org/get-stream/-/get-stream-9.0.1.tgz",
+      "integrity": "sha512-kVCxPF3vQM/N0B1PmoqVUqgHP+EeVjmZSQn+1oCRPxd2P21P2F19lIgbR3HBosbB1PUhOAoctJnfEn2GbN2eZA==",
+      "license": "MIT",
+      "dependencies": {
+        "@sec-ant/readable-stream": "^0.4.1",
+        "is-stream": "^4.0.1"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/get-symbol-description": {
       "version": "1.1.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bound": "^1.0.3",
         "es-errors": "^1.3.0",
         "get-intrinsic": "^1.2.6"
       },
       "engines": {
@@ -3672,32 +5156,30 @@
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/gopd": {
       "version": "1.2.0",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/graceful-fs": {
       "version": "4.2.11",
-      "dev": true,
       "license": "ISC"
     },
     "node_modules/has-bigints": {
       "version": "1.1.0",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
@@ -3732,21 +5214,20 @@
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/has-symbols": {
       "version": "1.1.0",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/has-tostringtag": {
       "version": "1.0.2",
@@ -3757,21 +5238,20 @@
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/hasown": {
       "version": "2.0.4",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "function-bind": "^1.1.2"
       },
       "engines": {
         "node": ">= 0.4"
       }
     },
     "node_modules/hermes-estree": {
       "version": "0.25.1",
@@ -3779,80 +5259,162 @@
       "license": "MIT"
     },
     "node_modules/hermes-parser": {
       "version": "0.25.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "hermes-estree": "0.25.1"
       }
     },
+    "node_modules/hono": {
+      "version": "4.13.5",
+      "resolved": "https://registry.npmjs.org/hono/-/hono-4.13.5.tgz",
+      "integrity": "sha512-O6+/eCYRkzzzy0rPWwKLiGBR1nFuUPZynnwjxN1MBA62NNqbT0wQEzQyK2gSO5yDIDB336sXQleAhOHrzlYyKw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=16.9.0"
+      }
+    },
+    "node_modules/http-errors": {
+      "version": "2.0.1",
+      "resolved": "https://registry.npmjs.org/http-errors/-/http-errors-2.0.1.tgz",
+      "integrity": "sha512-4FbRdAX+bSdmo4AUFuS0WNiPz8NgFt+r8ThgNWmlrjQjt1Q7ZR9+zTlce2859x4KSXrwIsaeTqDoKQmtP8pLmQ==",
+      "license": "MIT",
+      "dependencies": {
+        "depd": "~2.0.0",
+        "inherits": "~2.0.4",
+        "setprototypeof": "~1.2.0",
+        "statuses": "~2.0.2",
+        "toidentifier": "~1.0.1"
+      },
+      "engines": {
+        "node": ">= 0.8"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
+    "node_modules/human-signals": {
+      "version": "8.0.1",
+      "resolved": "https://registry.npmjs.org/human-signals/-/human-signals-8.0.1.tgz",
+      "integrity": "sha512-eKCa6bwnJhvxj14kZk5NCPc6Hb6BdsU9DZcOnmQKSnO1VKrfV0zCvtttPZUsBvjmNDn8rpcJfpwSYnHBjc95MQ==",
+      "license": "Apache-2.0",
+      "engines": {
+        "node": ">=18.18.0"
+      }
+    },
+    "node_modules/iconv-lite": {
+      "version": "0.7.3",
+      "resolved": "https://registry.npmjs.org/iconv-lite/-/iconv-lite-0.7.3.tgz",
+      "integrity": "sha512-IKXpvIzjnC9XTAUbVBcMfGS0EPaIXtW6v+zr+RRp+hqULEpo0owZax6wyRwPOJbWbzjYspQwusTsfVr0ifh4uQ==",
+      "license": "MIT",
+      "dependencies": {
+        "safer-buffer": ">= 2.1.2 < 3.0.0"
+      },
+      "engines": {
+        "node": ">=0.10.0"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
     "node_modules/ignore": {
       "version": "5.3.2",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 4"
       }
     },
     "node_modules/import-fresh": {
       "version": "3.3.1",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "parent-module": "^1.0.0",
         "resolve-from": "^4.0.0"
       },
       "engines": {
         "node": ">=6"
       },
       "funding": {
         "url": "https://github.com/sponsors/sindresorhus"
       }
     },
     "node_modules/imurmurhash": {
       "version": "0.1.4",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=0.8.19"
       }
     },
+    "node_modules/inherits": {
+      "version": "2.0.4",
+      "resolved": "https://registry.npmjs.org/inherits/-/inherits-2.0.4.tgz",
+      "integrity": "sha512-k/vGaX4/Yla3WzyMCvTQOXYeIHvqOKtnqBduzTHpzpQZzAskKMhZ2K+EnBiSM9zGSoIFeMpXKxa4dYeZIQqewQ==",
+      "license": "ISC"
+    },
     "node_modules/internal-slot": {
       "version": "1.1.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "es-errors": "^1.3.0",
         "hasown": "^2.0.2",
         "side-channel": "^1.1.0"
       },
       "engines": {
         "node": ">= 0.4"
       }
     },
+    "node_modules/ip-address": {
+      "version": "10.7.0",
+      "resolved": "https://registry.npmjs.org/ip-address/-/ip-address-10.7.0.tgz",
+      "integrity": "sha512-BGFsyJd5mpXp3rK6jIdADLNgpJUK1jnjzvYF8lK+VyDab9JAmqN0YOKDdP17HlgKb2+ehPgDc8EtnRLbGCAMhA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 12"
+      }
+    },
+    "node_modules/ipaddr.js": {
+      "version": "1.9.1",
+      "resolved": "https://registry.npmjs.org/ipaddr.js/-/ipaddr.js-1.9.1.tgz",
+      "integrity": "sha512-0KI/607xoxSToH7GjN1FfSbLoU0+btTicjsQSWQlh/hZykN8KpmMf7uYwPW3R+akZ6R/w18ZlXSHBYXiYUPO3g==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.10"
+      }
+    },
     "node_modules/is-array-buffer": {
       "version": "3.0.5",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bind": "^1.0.8",
         "call-bound": "^1.0.3",
         "get-intrinsic": "^1.2.6"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/is-arrayish": {
+      "version": "0.2.1",
+      "resolved": "https://registry.npmjs.org/is-arrayish/-/is-arrayish-0.2.1.tgz",
+      "integrity": "sha512-zz06S8t0ozoDXMG+ube26zeCTNXcKIPJZJi8hBrF4idCLms4CG9QtK7qBl1boi5ODzFpjswb5JPmHCbMpjaYzg==",
+      "license": "MIT"
+    },
     "node_modules/is-async-function": {
       "version": "2.1.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "async-function": "^1.0.0",
         "call-bound": "^1.0.3",
         "get-proto": "^1.0.1",
         "has-tostringtag": "^1.0.2",
         "safe-regex-test": "^1.1.0"
@@ -3961,37 +5523,51 @@
         "call-bound": "^1.0.2",
         "has-tostringtag": "^1.0.2"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/is-docker": {
+      "version": "3.0.0",
+      "resolved": "https://registry.npmjs.org/is-docker/-/is-docker-3.0.0.tgz",
+      "integrity": "sha512-eljcgEDlEns/7AXFosB5K/2nCM4P7FQPkGc/DWLy5rmFEWvZayGrik1d9/QIY5nJ4f9YsVvBkA6kJpHn9rISdQ==",
+      "license": "MIT",
+      "bin": {
+        "is-docker": "cli.js"
+      },
+      "engines": {
+        "node": "^12.20.0 || ^14.13.1 || >=16.0.0"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/is-document.all": {
       "version": "1.0.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bound": "^1.0.4"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/is-extglob": {
       "version": "2.1.1",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=0.10.0"
       }
     },
     "node_modules/is-finalizationregistry": {
       "version": "1.1.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
@@ -4017,29 +5593,70 @@
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/is-glob": {
       "version": "4.0.3",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "is-extglob": "^2.1.1"
       },
       "engines": {
         "node": ">=0.10.0"
       }
     },
+    "node_modules/is-in-ssh": {
+      "version": "1.0.0",
+      "resolved": "https://registry.npmjs.org/is-in-ssh/-/is-in-ssh-1.0.0.tgz",
+      "integrity": "sha512-jYa6Q9rH90kR1vKB6NM7qqd1mge3Fx4Dhw5TVlK1MUBqhEOuCagrEHMevNuCcbECmXZ0ThXkRm+Ymr51HwEPAw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=20"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/is-inside-container": {
+      "version": "1.0.0",
+      "resolved": "https://registry.npmjs.org/is-inside-container/-/is-inside-container-1.0.0.tgz",
+      "integrity": "sha512-KIYLCCJghfHZxqjYBE7rEy0OBuTd5xCHS7tHVgvCLkx7StIoaxwNW3hCALgEUjFfeRk+MG/Qxmp/vtETEF3tRA==",
+      "license": "MIT",
+      "dependencies": {
+        "is-docker": "^3.0.0"
+      },
+      "bin": {
+        "is-inside-container": "cli.js"
+      },
+      "engines": {
+        "node": ">=14.16"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/is-interactive": {
+      "version": "2.0.0",
+      "resolved": "https://registry.npmjs.org/is-interactive/-/is-interactive-2.0.0.tgz",
+      "integrity": "sha512-qP1vozQRI+BMOPcjFzrjXuQvdak2pHNUMZoeG2eRbiSqyvbEf/wQtEOTOX1guk6E3t36RkaqiSt8A/6YElNxLQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/is-map": {
       "version": "2.0.3",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
@@ -4050,58 +5667,96 @@
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/is-number": {
       "version": "7.0.0",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=0.12.0"
       }
     },
     "node_modules/is-number-object": {
       "version": "1.1.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bound": "^1.0.3",
         "has-tostringtag": "^1.0.2"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/is-obj": {
+      "version": "2.0.0",
+      "resolved": "https://registry.npmjs.org/is-obj/-/is-obj-2.0.0.tgz",
+      "integrity": "sha512-drqDG3cbczxxEJRoOXcOjtdp1J/lyp1mNn0xaznRs8+muBhgQcrnbspox5X5fOw0HnMnbfDzvnEMEtqDEJEo8w==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=8"
+      }
+    },
+    "node_modules/is-plain-obj": {
+      "version": "4.1.0",
+      "resolved": "https://registry.npmjs.org/is-plain-obj/-/is-plain-obj-4.1.0.tgz",
+      "integrity": "sha512-+Pgi+vMuUNkJyExiMBt5IlFoMyKnr5zhJ4Uspz58WOhBF5QoIZkFyNHIbBAtHwzVAgk5RtndVNsDRN61/mmDqg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/is-promise": {
+      "version": "4.0.0",
+      "resolved": "https://registry.npmjs.org/is-promise/-/is-promise-4.0.0.tgz",
+      "integrity": "sha512-hvpoI6korhJMnej285dSg6nu1+e6uxs7zG3BYAm5byqDsgJNWwxzM6z6iZiAgQR4TJ30JmBTOwqZUw3WlyH3AQ==",
+      "license": "MIT"
+    },
     "node_modules/is-regex": {
       "version": "1.2.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bound": "^1.0.2",
         "gopd": "^1.2.0",
         "has-tostringtag": "^1.0.2",
         "hasown": "^2.0.2"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/is-regexp": {
+      "version": "3.1.0",
+      "resolved": "https://registry.npmjs.org/is-regexp/-/is-regexp-3.1.0.tgz",
+      "integrity": "sha512-rbku49cWloU5bSMI+zaRaXdQHXnthP6DZ/vLnfdSKyL4zUzuWnomtOEiZZOd+ioQ+avFo/qau3KPTc7Fjy1uPA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/is-set": {
       "version": "2.0.3",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
@@ -4113,20 +5768,32 @@
       "dependencies": {
         "call-bound": "^1.0.3"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/is-stream": {
+      "version": "4.0.1",
+      "resolved": "https://registry.npmjs.org/is-stream/-/is-stream-4.0.1.tgz",
+      "integrity": "sha512-Dnz92NInDqYckGEUJv689RbRiTSEHCQ7wOVeALbkOz999YpqT46yMRIGtSNl2iCL1waAZSx40+h59NV/EwzV/A==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/is-string": {
       "version": "1.1.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bound": "^1.0.3",
         "has-tostringtag": "^1.0.2"
       },
       "engines": {
         "node": ">= 0.4"
@@ -4158,20 +5825,32 @@
       "dependencies": {
         "which-typed-array": "^1.1.16"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/is-unicode-supported": {
+      "version": "2.1.0",
+      "resolved": "https://registry.npmjs.org/is-unicode-supported/-/is-unicode-supported-2.1.0.tgz",
+      "integrity": "sha512-mE00Gnza5EEB3Ds0HfMyllZzbBrmLOX3vfWoj9A9PEnTfratQ/BcaJOuMhnkhjXvb2+FkY3VuHqtAGpTPmglFQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/is-weakmap": {
       "version": "2.0.2",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
@@ -4198,28 +5877,42 @@
         "call-bound": "^1.0.3",
         "get-intrinsic": "^1.2.6"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/is-wsl": {
+      "version": "3.1.1",
+      "resolved": "https://registry.npmjs.org/is-wsl/-/is-wsl-3.1.1.tgz",
+      "integrity": "sha512-e6rvdUCiQCAuumZslxRJWR/Doq4VpPR82kqclvcS0efgt430SlGIk05vdCN58+VrzgtIcfNODjozVielycD4Sw==",
+      "license": "MIT",
+      "dependencies": {
+        "is-inside-container": "^1.0.0"
+      },
+      "engines": {
+        "node": ">=16"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/isarray": {
       "version": "2.0.5",
       "dev": true,
       "license": "MIT"
     },
     "node_modules/isexe": {
       "version": "2.0.0",
-      "dev": true,
       "license": "ISC"
     },
     "node_modules/iterator.prototype": {
       "version": "1.1.5",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "define-data-property": "^1.1.4",
         "es-object-atoms": "^1.0.0",
         "get-intrinsic": "^1.2.6",
@@ -4232,83 +5925,112 @@
       }
     },
     "node_modules/jiti": {
       "version": "2.7.0",
       "dev": true,
       "license": "MIT",
       "bin": {
         "jiti": "lib/jiti-cli.mjs"
       }
     },
+    "node_modules/jose": {
+      "version": "6.2.10",
+      "resolved": "https://registry.npmjs.org/jose/-/jose-6.2.10.tgz",
+      "integrity": "sha512-iiW7J9qRFlGxvCOIBDBDxFePQSn7ZMAnrYGhrrOo6siO/MIqwfyilLR27pkfDgUk+raLuzADS8A3S/KLBisc0g==",
+      "license": "MIT",
+      "funding": {
+        "url": "https://github.com/sponsors/panva"
+      }
+    },
     "node_modules/js-tokens": {
       "version": "4.0.0",
-      "dev": true,
       "license": "MIT"
     },
     "node_modules/js-yaml": {
       "version": "4.3.2",
-      "dev": true,
       "funding": [
         {
           "type": "github",
           "url": "https://github.com/sponsors/puzrin"
         },
         {
           "type": "github",
           "url": "https://github.com/sponsors/nodeca"
         }
       ],
       "license": "MIT",
       "dependencies": {
         "argparse": "^2.0.1"
       },
       "bin": {
         "js-yaml": "bin/js-yaml.js"
       }
     },
     "node_modules/jsesc": {
       "version": "3.1.0",
-      "dev": true,
       "license": "MIT",
       "bin": {
         "jsesc": "bin/jsesc"
       },
       "engines": {
         "node": ">=6"
       }
     },
     "node_modules/json-buffer": {
       "version": "3.0.1",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/json-parse-even-better-errors": {
+      "version": "2.3.1",
+      "resolved": "https://registry.npmjs.org/json-parse-even-better-errors/-/json-parse-even-better-errors-2.3.1.tgz",
+      "integrity": "sha512-xyFwyhro/JEof6Ghe2iz2NcXoj2sloNsWr/XsERDK/oiPCfaNhl5ONfp+jQdAZRQQ0IJWNzH9zIZF7li91kh2w==",
+      "license": "MIT"
+    },
     "node_modules/json-schema-traverse": {
       "version": "0.4.1",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/json-schema-typed": {
+      "version": "8.0.2",
+      "resolved": "https://registry.npmjs.org/json-schema-typed/-/json-schema-typed-8.0.2.tgz",
+      "integrity": "sha512-fQhoXdcvc3V28x7C7BMs4P5+kNlgUURe2jmUT1T//oBRMDrqy1QPelJimwZGo7Hg9VPV3EQV5Bnq4hbFy2vetA==",
+      "license": "BSD-2-Clause"
+    },
     "node_modules/json-stable-stringify-without-jsonify": {
       "version": "1.0.1",
       "dev": true,
       "license": "MIT"
     },
     "node_modules/json5": {
       "version": "2.2.3",
-      "dev": true,
       "license": "MIT",
       "bin": {
         "json5": "lib/cli.js"
       },
       "engines": {
         "node": ">=6"
       }
     },
+    "node_modules/jsonfile": {
+      "version": "6.2.1",
+      "resolved": "https://registry.npmjs.org/jsonfile/-/jsonfile-6.2.1.tgz",
+      "integrity": "sha512-zwOTdL3rFQ/lRdBnntKVOX6k5cKJwEc1HdilT71BWEu7J41gXIB2MRp+vxduPSwZJPWBxEzv4yH1wYLJGUHX4Q==",
+      "license": "MIT",
+      "dependencies": {
+        "universalify": "^2.0.0"
+      },
+      "optionalDependencies": {
+        "graceful-fs": "^4.1.6"
+      }
+    },
     "node_modules/jsx-ast-utils": {
       "version": "3.3.5",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "array-includes": "^3.1.6",
         "array.prototype.flat": "^1.3.1",
         "object.assign": "^4.1.4",
         "object.values": "^1.1.6"
       },
@@ -4317,20 +6039,29 @@
       }
     },
     "node_modules/keyv": {
       "version": "4.5.4",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "json-buffer": "3.0.1"
       }
     },
+    "node_modules/kleur": {
+      "version": "4.1.5",
+      "resolved": "https://registry.npmjs.org/kleur/-/kleur-4.1.5.tgz",
+      "integrity": "sha512-o+NO+8WrRiQEE4/7nwRJhN1HWpVmJm511pBHUxPLtp0BUISzlBplORYSmTclCnJvQq2tKu/sgl3xVpkc7ZWuQQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6"
+      }
+    },
     "node_modules/language-subtag-registry": {
       "version": "0.3.23",
       "dev": true,
       "license": "CC0-1.0"
     },
     "node_modules/language-tags": {
       "version": "1.0.9",
       "dev": true,
       "license": "MIT",
       "dependencies": {
@@ -4614,116 +6345,242 @@
         "win32"
       ],
       "engines": {
         "node": ">= 12.0.0"
       },
       "funding": {
         "type": "opencollective",
         "url": "https://opencollective.com/parcel"
       }
     },
+    "node_modules/lines-and-columns": {
+      "version": "1.2.4",
+      "resolved": "https://registry.npmjs.org/lines-and-columns/-/lines-and-columns-1.2.4.tgz",
+      "integrity": "sha512-7ylylesZQ/PV29jhEDl3Ufjo6ZX7gCqJr5F7PKrqc93v7fzSymt1BpwEU8nAUXs8qzzvqhbjhK5QZg6Mt/HkBg==",
+      "license": "MIT"
+    },
     "node_modules/locate-path": {
       "version": "6.0.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "p-locate": "^5.0.0"
       },
       "engines": {
         "node": ">=10"
       },
       "funding": {
         "url": "https://github.com/sponsors/sindresorhus"
       }
     },
     "node_modules/lodash.merge": {
       "version": "4.6.2",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/log-symbols": {
+      "version": "6.0.0",
+      "resolved": "https://registry.npmjs.org/log-symbols/-/log-symbols-6.0.0.tgz",
+      "integrity": "sha512-i24m8rpwhmPIS4zscNzK6MSEhk0DUWa/8iYQWxhffV8jkI4Phvs3F+quL5xvS0gdQR0FyTCMMH33Y78dDTzzIw==",
+      "license": "MIT",
+      "dependencies": {
+        "chalk": "^5.3.0",
+        "is-unicode-supported": "^1.3.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/log-symbols/node_modules/chalk": {
+      "version": "5.6.2",
+      "resolved": "https://registry.npmjs.org/chalk/-/chalk-5.6.2.tgz",
+      "integrity": "sha512-7NzBL0rN6fMUW+f7A6Io4h40qQlG+xGmtMxfbnH/K7TAtt8JQWVQK+6g0UXKMeVJoyV5EkkNsErQ8pVD3bLHbA==",
+      "license": "MIT",
+      "engines": {
+        "node": "^12.17.0 || ^14.13 || >=16.0.0"
+      },
+      "funding": {
+        "url": "https://github.com/chalk/chalk?sponsor=1"
+      }
+    },
+    "node_modules/log-symbols/node_modules/is-unicode-supported": {
+      "version": "1.3.0",
+      "resolved": "https://registry.npmjs.org/is-unicode-supported/-/is-unicode-supported-1.3.0.tgz",
+      "integrity": "sha512-43r2mRvz+8JRIKnWJ+3j8JtjRKZ6GmjzfaE/qiBJnikNnYv/6bagRJ1kUhNk8R5EX/GkobD+r+sfxCPJsiKBLQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/loose-envify": {
       "version": "1.4.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "js-tokens": "^3.0.0 || ^4.0.0"
       },
       "bin": {
         "loose-envify": "cli.js"
       }
     },
     "node_modules/lru-cache": {
       "version": "5.1.1",
-      "dev": true,
       "license": "ISC",
       "dependencies": {
         "yallist": "^3.0.2"
       }
     },
+    "node_modules/lucide-react": {
+      "version": "1.37.0",
+      "resolved": "https://registry.npmjs.org/lucide-react/-/lucide-react-1.37.0.tgz",
+      "integrity": "sha512-LPsB4rD1TD6wZu1djKOf9vUnS1jTNaHbolXebXDgiTdb6jeA1agIJhJsIybCmjKmQClcOaal1o1OaiYahEftyQ==",
+      "license": "ISC",
+      "peerDependencies": {
+        "react": "^16.5.1 || ^17.0.0 || ^18.0.0 || ^19.0.0"
+      }
+    },
     "node_modules/magic-string": {
       "version": "0.30.21",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "@jridgewell/sourcemap-codec": "^1.5.5"
       }
     },
     "node_modules/math-intrinsics": {
       "version": "1.1.0",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       }
     },
+    "node_modules/media-typer": {
+      "version": "1.1.1",
+      "resolved": "https://registry.npmjs.org/media-typer/-/media-typer-1.1.1.tgz",
+      "integrity": "sha512-yz3xRaG20c6/BOzvYoDaGtPmGscs7YivItZEEqe6GbwNfHuxu9YNmvnEkMzKldAGY4/80pRcQRZSEnhquk9XuQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.8"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
+    "node_modules/merge-descriptors": {
+      "version": "2.0.0",
+      "resolved": "https://registry.npmjs.org/merge-descriptors/-/merge-descriptors-2.0.0.tgz",
+      "integrity": "sha512-Snk314V5ayFLhp3fkUREub6WtjBfPdCPY1Ln8/8munuLuiYhsABgBVWsozAG+MWMbVEvcdcpbi9R7ww22l9Q3g==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/merge-stream": {
+      "version": "2.0.0",
+      "resolved": "https://registry.npmjs.org/merge-stream/-/merge-stream-2.0.0.tgz",
+      "integrity": "sha512-abv/qOcuPfk3URPfDzmZU1LKmuw8kT+0nIHvKrKgFrwifol/doWcdA4ZqsWQ8ENrFKkd67Mfpo/LovbIUsbt3w==",
+      "license": "MIT"
+    },
     "node_modules/merge2": {
       "version": "1.4.1",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 8"
       }
     },
     "node_modules/micromatch": {
       "version": "4.0.8",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "braces": "^3.0.3",
         "picomatch": "^2.3.1"
       },
       "engines": {
         "node": ">=8.6"
       }
     },
+    "node_modules/mime-db": {
+      "version": "1.54.0",
+      "resolved": "https://registry.npmjs.org/mime-db/-/mime-db-1.54.0.tgz",
+      "integrity": "sha512-aU5EJuIN2WDemCcAp2vFBfp/m4EAhWJnUNSSw0ixs7/kXbd6Pg64EmwJkNdFhB8aWt1sH2CTXrLxo/iAGV3oPQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.6"
+      }
+    },
+    "node_modules/mime-types": {
+      "version": "3.0.2",
+      "resolved": "https://registry.npmjs.org/mime-types/-/mime-types-3.0.2.tgz",
+      "integrity": "sha512-Lbgzdk0h4juoQ9fCKXW4by0UJqj+nOOrI9MJ1sSj4nI8aI2eo1qmvQEie4VD1glsS250n15LsWsYtCugiStS5A==",
+      "license": "MIT",
+      "dependencies": {
+        "mime-db": "^1.54.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
+    "node_modules/mimic-fn": {
+      "version": "3.1.0",
+      "resolved": "https://registry.npmjs.org/mimic-fn/-/mimic-fn-3.1.0.tgz",
+      "integrity": "sha512-Ysbi9uYW9hFyfrThdDEQuykN4Ey6BuwPD2kpI5ES/nFTDn/98yxYNLZJcgUAKPT/mcrLLKaGzJR9YVxJrIdASQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=8"
+      }
+    },
+    "node_modules/mimic-function": {
+      "version": "5.0.1",
+      "resolved": "https://registry.npmjs.org/mimic-function/-/mimic-function-5.0.1.tgz",
+      "integrity": "sha512-VP79XUPxV2CigYP3jWwAUFSku2aKqBH7uTAapFWCBqutsbmDo96KY5o8uh6U+/YSIn5OxJnXp73beVkpqMIGhA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/minimatch": {
       "version": "3.1.5",
       "dev": true,
       "license": "ISC",
       "dependencies": {
         "brace-expansion": "^1.1.7"
       },
       "engines": {
         "node": "*"
       }
     },
     "node_modules/minimist": {
       "version": "1.2.8",
-      "dev": true,
       "license": "MIT",
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/ms": {
       "version": "2.1.3",
-      "dev": true,
       "license": "MIT"
     },
     "node_modules/nanoid": {
       "version": "3.3.18",
       "funding": [
         {
           "type": "github",
           "url": "https://github.com/sponsors/ai"
         }
       ],
@@ -4747,20 +6604,49 @@
       },
       "funding": {
         "url": "https://opencollective.com/napi-postinstall"
       }
     },
     "node_modules/natural-compare": {
       "version": "1.4.0",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/negotiator": {
+      "version": "1.1.0",
+      "resolved": "https://registry.npmjs.org/negotiator/-/negotiator-1.1.0.tgz",
+      "integrity": "sha512-NMPBRMJgiQHjbd8phG3Vebdx4kZ1H121rbl5IkMqeOsahptB9BKo/d7oJ3zTXqTgagn2bWlNSXkh0QUGM31RYg==",
+      "license": "MIT",
+      "dependencies": {
+        "content-type": "^2.1.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
+    "node_modules/negotiator/node_modules/content-type": {
+      "version": "2.1.0",
+      "resolved": "https://registry.npmjs.org/content-type/-/content-type-2.1.0.tgz",
+      "integrity": "sha512-mj7UPXE0jaqaOsukNZRUEfEi2AcL7C/vwmwcHV0O97eO1E1pxBZuyjlZrx5seTaNBg1U6+o35wpa35Qfcc+7ag==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
     "node_modules/next": {
       "version": "16.3.3",
       "license": "MIT",
       "dependencies": {
         "@next/env": "16.3.3",
         "@swc/helpers": "0.5.23",
         "baseline-browser-mapping": "^2.9.19",
         "caniuse-lite": "^1.0.30001579",
         "postcss": "8.5.23",
         "styled-jsx": "5.1.6"
@@ -4798,20 +6684,30 @@
           "optional": true
         },
         "babel-plugin-react-compiler": {
           "optional": true
         },
         "sass": {
           "optional": true
         }
       }
     },
+    "node_modules/next-themes": {
+      "version": "0.4.6",
+      "resolved": "https://registry.npmjs.org/next-themes/-/next-themes-0.4.6.tgz",
+      "integrity": "sha512-pZvgD5L0IEvX5/9GWyHMf3m8BKiVQwsCMHfoFosXtXBMnaS0ZnIJ9ST4b4NqLVKDEm8QBxoNNGNaBv2JNF6XNA==",
+      "license": "MIT",
+      "peerDependencies": {
+        "react": "^16.8 || ^17 || ^18 || ^19 || ^19.0.0-rc",
+        "react-dom": "^16.8 || ^17 || ^18 || ^19 || ^19.0.0-rc"
+      }
+    },
     "node_modules/next/node_modules/postcss": {
       "version": "8.5.23",
       "funding": [
         {
           "type": "opencollective",
           "url": "https://opencollective.com/postcss/"
         },
         {
           "type": "tidelift",
           "url": "https://tidelift.com/funding/github/npm/postcss"
@@ -4843,53 +6739,87 @@
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/node-releases": {
       "version": "2.0.54",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=18"
       }
     },
+    "node_modules/npm-run-path": {
+      "version": "6.0.0",
+      "resolved": "https://registry.npmjs.org/npm-run-path/-/npm-run-path-6.0.0.tgz",
+      "integrity": "sha512-9qny7Z9DsQU8Ou39ERsPU4OZQlSTP47ShQzuKZ6PRXpYLtIFgl/DEBYEXKlvcEa+9tHVcK8CF81Y2V72qaZhWA==",
+      "license": "MIT",
+      "dependencies": {
+        "path-key": "^4.0.0",
+        "unicorn-magic": "^0.3.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/npm-run-path/node_modules/path-key": {
+      "version": "4.0.0",
+      "resolved": "https://registry.npmjs.org/path-key/-/path-key-4.0.0.tgz",
+      "integrity": "sha512-haREypq7xkM7ErfgIyA0z+Bj4AGKlMSdlQE2jvJo6huWD1EdkKYV+G/T4nq0YEF2vgTT8kqMFKo1uHn950r4SQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/object-assign": {
       "version": "4.1.1",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=0.10.0"
       }
     },
     "node_modules/object-inspect": {
       "version": "1.13.4",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/object-keys": {
       "version": "1.1.1",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       }
     },
+    "node_modules/object-treeify": {
+      "version": "1.1.33",
+      "resolved": "https://registry.npmjs.org/object-treeify/-/object-treeify-1.1.33.tgz",
+      "integrity": "sha512-EFVjAYfzWqWsBMRHPMAXLCDIJnpMhdWAqR7xG6M6a2cs6PMFpl/+Z20w9zDW4vkxOFfddegBKq9Rehd0bxWE7A==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 10"
+      }
+    },
     "node_modules/object.assign": {
       "version": "4.1.7",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bind": "^1.0.8",
         "call-bound": "^1.0.3",
         "define-properties": "^1.2.1",
         "es-object-atoms": "^1.0.0",
         "has-symbols": "^1.1.0",
@@ -4956,36 +6886,163 @@
         "define-properties": "^1.2.1",
         "es-object-atoms": "^1.0.0"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/on-finished": {
+      "version": "2.4.1",
+      "resolved": "https://registry.npmjs.org/on-finished/-/on-finished-2.4.1.tgz",
+      "integrity": "sha512-oVlzkg3ENAhCk2zdv7IJwd/QUD4z2RxRwpkcGY8psCVcCYZNq4wYnVWALHM+brtuJjePWiYF/ClmuDr8Ch5+kg==",
+      "license": "MIT",
+      "dependencies": {
+        "ee-first": "1.1.1"
+      },
+      "engines": {
+        "node": ">= 0.8"
+      }
+    },
+    "node_modules/once": {
+      "version": "1.4.0",
+      "resolved": "https://registry.npmjs.org/once/-/once-1.4.0.tgz",
+      "integrity": "sha512-lNaJgI+2Q5URQBkccEKHTQOPaXdUxnZZElQTZY0MFUAuaEqe1E+Nyvgdz/aIyNi6Z9MzO5dv1H8n58/GELp3+w==",
+      "license": "ISC",
+      "dependencies": {
+        "wrappy": "1"
+      }
+    },
+    "node_modules/onetime": {
+      "version": "5.1.2",
+      "resolved": "https://registry.npmjs.org/onetime/-/onetime-5.1.2.tgz",
+      "integrity": "sha512-kbpaSSGJTWdAY5KPVeMOKXSrPtr8C8C7wodJbcsd51jRnmD+GZu8Y0VoU6Dm5Z4vWr0Ig/1NKuWRKf7j5aaYSg==",
+      "license": "MIT",
+      "dependencies": {
+        "mimic-fn": "^2.1.0"
+      },
+      "engines": {
+        "node": ">=6"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/onetime/node_modules/mimic-fn": {
+      "version": "2.1.0",
+      "resolved": "https://registry.npmjs.org/mimic-fn/-/mimic-fn-2.1.0.tgz",
+      "integrity": "sha512-OqbOk5oEQeAZ8WXWydlu9HJjz9WVdEIvamMCcXmuqUYjTknH/sqsWvhQ3vgwKFRR1HpjvNBKQ37nbJgYzGqGcg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6"
+      }
+    },
+    "node_modules/open": {
+      "version": "11.0.2",
+      "resolved": "https://registry.npmjs.org/open/-/open-11.0.2.tgz",
+      "integrity": "sha512-RWqF+pBSkqecEvCKOn8QYhaNdRMJDZRIrlS/7rTDdLHaPcfXGCZ/h8zb413NfvdeAV0MR7T1yJcA34/q+CSm1Q==",
+      "license": "MIT",
+      "dependencies": {
+        "default-browser": "^5.5.1",
+        "define-lazy-prop": "^3.0.0",
+        "is-in-ssh": "^1.0.0",
+        "is-inside-container": "^1.0.0",
+        "powershell-utils": "^0.2.1",
+        "wsl-utils": "^1.0.0"
+      },
+      "engines": {
+        "node": ">=20"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/optionator": {
       "version": "0.9.4",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "deep-is": "^0.1.3",
         "fast-levenshtein": "^2.0.6",
         "levn": "^0.4.1",
         "prelude-ls": "^1.2.1",
         "type-check": "^0.4.0",
         "word-wrap": "^1.2.5"
       },
       "engines": {
         "node": ">= 0.8.0"
       }
     },
+    "node_modules/ora": {
+      "version": "8.2.0",
+      "resolved": "https://registry.npmjs.org/ora/-/ora-8.2.0.tgz",
+      "integrity": "sha512-weP+BZ8MVNnlCm8c0Qdc1WSWq4Qn7I+9CJGm7Qali6g44e/PUzbjNqJX5NJ9ljlNMosfJvg1fKEGILklK9cwnw==",
+      "license": "MIT",
+      "dependencies": {
+        "chalk": "^5.3.0",
+        "cli-cursor": "^5.0.0",
+        "cli-spinners": "^2.9.2",
+        "is-interactive": "^2.0.0",
+        "is-unicode-supported": "^2.0.0",
+        "log-symbols": "^6.0.0",
+        "stdin-discarder": "^0.2.2",
+        "string-width": "^7.2.0",
+        "strip-ansi": "^7.1.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/ora/node_modules/ansi-regex": {
+      "version": "6.3.0",
+      "resolved": "https://registry.npmjs.org/ansi-regex/-/ansi-regex-6.3.0.tgz",
+      "integrity": "sha512-WpDfL7NO6j7tH88IDBNVdUJxDh9nmCteAVW9dsep846XdwF4naCBK+/tGLX3KJgcpgMRXCFlTM2hKGoK9FsdrQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/chalk/ansi-regex?sponsor=1"
+      }
+    },
+    "node_modules/ora/node_modules/chalk": {
+      "version": "5.6.2",
+      "resolved": "https://registry.npmjs.org/chalk/-/chalk-5.6.2.tgz",
+      "integrity": "sha512-7NzBL0rN6fMUW+f7A6Io4h40qQlG+xGmtMxfbnH/K7TAtt8JQWVQK+6g0UXKMeVJoyV5EkkNsErQ8pVD3bLHbA==",
+      "license": "MIT",
+      "engines": {
+        "node": "^12.17.0 || ^14.13 || >=16.0.0"
+      },
+      "funding": {
+        "url": "https://github.com/chalk/chalk?sponsor=1"
+      }
+    },
+    "node_modules/ora/node_modules/strip-ansi": {
+      "version": "7.2.0",
+      "resolved": "https://registry.npmjs.org/strip-ansi/-/strip-ansi-7.2.0.tgz",
+      "integrity": "sha512-yDPMNjp4WyfYBkHnjIRLfca1i6KMyGCtsVgoKe/z1+6vukgaENdgGBZt+ZmKPc4gavvEZ5OgHfHdrazhgNyG7w==",
+      "license": "MIT",
+      "dependencies": {
+        "ansi-regex": "^6.2.2"
+      },
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/chalk/strip-ansi?sponsor=1"
+      }
+    },
     "node_modules/own-keys": {
       "version": "1.0.2",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bound": "^1.0.4",
         "get-intrinsic": "^1.3.0",
         "object-keys": "^1.1.1",
         "safe-push-apply": "^1.0.0"
       },
@@ -5017,78 +7074,220 @@
       "dependencies": {
         "p-limit": "^3.0.2"
       },
       "engines": {
         "node": ">=10"
       },
       "funding": {
         "url": "https://github.com/sponsors/sindresorhus"
       }
     },
+    "node_modules/p-try": {
+      "version": "2.2.0",
+      "resolved": "https://registry.npmjs.org/p-try/-/p-try-2.2.0.tgz",
+      "integrity": "sha512-R4nPAVTAU0B9D35/Gk3uJf/7XYbQcyohSKdvAxIRSNghFl4e71hVoGnBNQz9cWaXxO2I10KTC+3jMdvvoKw6dQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6"
+      }
+    },
     "node_modules/parent-module": {
       "version": "1.0.1",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "callsites": "^3.0.0"
       },
       "engines": {
-        "node": ">=6"
+        "node": ">=6"
+      }
+    },
+    "node_modules/parse-json": {
+      "version": "5.2.0",
+      "resolved": "https://registry.npmjs.org/parse-json/-/parse-json-5.2.0.tgz",
+      "integrity": "sha512-ayCKvm/phCGxOkYRSCM82iDwct8/EonSEgCSxWxD7ve6jHggsFl4fZVQBPRNgQoKiuV/odhFrGzQXZwbifC8Rg==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/code-frame": "^7.0.0",
+        "error-ex": "^1.3.1",
+        "json-parse-even-better-errors": "^2.3.0",
+        "lines-and-columns": "^1.1.6"
+      },
+      "engines": {
+        "node": ">=8"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/parse-ms": {
+      "version": "4.0.0",
+      "resolved": "https://registry.npmjs.org/parse-ms/-/parse-ms-4.0.0.tgz",
+      "integrity": "sha512-TXfryirbmq34y8QBwgqCVLi+8oA3oWx2eAnSn62ITyEhEYaWRlVZ2DvMM9eZbMs/RfxPu/PK/aBLyGj4IrqMHw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/parseurl": {
+      "version": "1.3.3",
+      "resolved": "https://registry.npmjs.org/parseurl/-/parseurl-1.3.3.tgz",
+      "integrity": "sha512-CiyeOxFT/JZyN5m0z9PfXw4SCBJ6Sygz1Dpl0wqjlhDEGGBP1GnsUVEL0p63hoG1fcj3fHynXi9NYO4nWOL+qQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.8"
       }
     },
+    "node_modules/path-browserify": {
+      "version": "1.0.1",
+      "resolved": "https://registry.npmjs.org/path-browserify/-/path-browserify-1.0.1.tgz",
+      "integrity": "sha512-b7uo2UCUOYZcnF/3ID0lulOJi/bafxa1xPe7ZPsammBSpjSWQkjNxlt635YGS2MiR9GjvuXCtz2emr3jbsz98g==",
+      "license": "MIT"
+    },
     "node_modules/path-exists": {
       "version": "4.0.0",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=8"
       }
     },
     "node_modules/path-key": {
       "version": "3.1.1",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=8"
       }
     },
     "node_modules/path-parse": {
       "version": "1.0.7",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/path-to-regexp": {
+      "version": "8.4.2",
+      "resolved": "https://registry.npmjs.org/path-to-regexp/-/path-to-regexp-8.4.2.tgz",
+      "integrity": "sha512-qRcuIdP69NPm4qbACK+aDogI5CBDMi1jKe0ry5rSQJz8JVLsC7jV8XpiJjGRLLol3N+R5ihGYcrPLTno6pAdBA==",
+      "license": "MIT",
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
     "node_modules/picocolors": {
       "version": "1.1.1",
       "license": "ISC"
     },
     "node_modules/picomatch": {
       "version": "2.3.2",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=8.6"
       },
       "funding": {
         "url": "https://github.com/sponsors/jonschlinkert"
       }
     },
+    "node_modules/pkce-challenge": {
+      "version": "5.0.1",
+      "resolved": "https://registry.npmjs.org/pkce-challenge/-/pkce-challenge-5.0.1.tgz",
+      "integrity": "sha512-wQ0b/W4Fr01qtpHlqSqspcj3EhBvimsdh0KlHhH8HRZnMsEa0ea2fTULOXOS9ccQr3om+GcGRk4e+isrZWV8qQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=16.20.0"
+      }
+    },
+    "node_modules/pkg-up": {
+      "version": "3.1.0",
+      "resolved": "https://registry.npmjs.org/pkg-up/-/pkg-up-3.1.0.tgz",
+      "integrity": "sha512-nDywThFk1i4BQK4twPQ6TA4RT8bDY96yeuCVBWL3ePARCiEKDRSrNGbFIgUJpLp+XeIR65v8ra7WuJOFUBtkMA==",
+      "license": "MIT",
+      "dependencies": {
+        "find-up": "^3.0.0"
+      },
+      "engines": {
+        "node": ">=8"
+      }
+    },
+    "node_modules/pkg-up/node_modules/find-up": {
+      "version": "3.0.0",
+      "resolved": "https://registry.npmjs.org/find-up/-/find-up-3.0.0.tgz",
+      "integrity": "sha512-1yD6RmLI1XBfxugvORwlck6f75tYL+iR0jqwsOrOxMZyGYqUuDhJ0l4AXdO1iX/FTs9cBAMEk1gWSEx1kSbylg==",
+      "license": "MIT",
+      "dependencies": {
+        "locate-path": "^3.0.0"
+      },
+      "engines": {
+        "node": ">=6"
+      }
+    },
+    "node_modules/pkg-up/node_modules/locate-path": {
+      "version": "3.0.0",
+      "resolved": "https://registry.npmjs.org/locate-path/-/locate-path-3.0.0.tgz",
+      "integrity": "sha512-7AO748wWnIhNqAuaty2ZWHkQHRSNfPVIsPIfwEOWO22AmaoVrWavlOcMR5nzTLNYvp36X220/maaRsrec1G65A==",
+      "license": "MIT",
+      "dependencies": {
+        "p-locate": "^3.0.0",
+        "path-exists": "^3.0.0"
+      },
+      "engines": {
+        "node": ">=6"
+      }
+    },
+    "node_modules/pkg-up/node_modules/p-limit": {
+      "version": "2.3.0",
+      "resolved": "https://registry.npmjs.org/p-limit/-/p-limit-2.3.0.tgz",
+      "integrity": "sha512-//88mFWSJx8lxCzwdAABTJL2MyWB12+eIY7MDL2SqLmAkeKU9qxRvWuSyTjm3FUmpBEMuFfckAIqEaVGUDxb6w==",
+      "license": "MIT",
+      "dependencies": {
+        "p-try": "^2.0.0"
+      },
+      "engines": {
+        "node": ">=6"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/pkg-up/node_modules/p-locate": {
+      "version": "3.0.0",
+      "resolved": "https://registry.npmjs.org/p-locate/-/p-locate-3.0.0.tgz",
+      "integrity": "sha512-x+12w/To+4GFfgJhBEpiDcLozRJGegY+Ei7/z0tSLkMmxGZNybVMSfWj9aJn8Z5Fc7dBUNJOOVgPv2H7IwulSQ==",
+      "license": "MIT",
+      "dependencies": {
+        "p-limit": "^2.0.0"
+      },
+      "engines": {
+        "node": ">=6"
+      }
+    },
+    "node_modules/pkg-up/node_modules/path-exists": {
+      "version": "3.0.0",
+      "resolved": "https://registry.npmjs.org/path-exists/-/path-exists-3.0.0.tgz",
+      "integrity": "sha512-bpC7GYwiDYQ4wYLe+FA8lhRjhQCMcQGuSgGGqDkg/QerRWw9CmGRT0iSOVRSZJ29NMLZgIzqaljJ63oaL4NIJQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=4"
+      }
+    },
     "node_modules/possible-typed-array-names": {
       "version": "1.1.0",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       }
     },
     "node_modules/postcss": {
       "version": "8.5.26",
-      "dev": true,
       "funding": [
         {
           "type": "opencollective",
           "url": "https://opencollective.com/postcss/"
         },
         {
           "type": "tidelift",
           "url": "https://tidelift.com/funding/github/npm/postcss"
         },
         {
@@ -5099,65 +7298,183 @@
       "license": "MIT",
       "dependencies": {
         "nanoid": "^3.3.17",
         "picocolors": "^1.1.1",
         "source-map-js": "^1.2.1"
       },
       "engines": {
         "node": "^10 || ^12 || >=14"
       }
     },
+    "node_modules/postcss-selector-parser": {
+      "version": "7.1.5",
+      "resolved": "https://registry.npmjs.org/postcss-selector-parser/-/postcss-selector-parser-7.1.5.tgz",
+      "integrity": "sha512-KvvtD7SrlBP7dlgkBghEE3r84CABm5SmV2aNcG4oCA+qDnJ/tvKonFVvwWAyyWUEwxuNawdfEAZKP9zM3oZ2Uw==",
+      "license": "MIT",
+      "dependencies": {
+        "cssesc": "^3.0.0",
+        "util-deprecate": "^1.0.2"
+      },
+      "engines": {
+        "node": ">=4"
+      }
+    },
+    "node_modules/powershell-utils": {
+      "version": "0.2.1",
+      "resolved": "https://registry.npmjs.org/powershell-utils/-/powershell-utils-0.2.1.tgz",
+      "integrity": "sha512-C+y9x90UElAddDZmV4qOx9W53B61PO7cIqWz2dQsWlwswuq4mr8NEwytdGKboYbQlGZ3awrkTeNvcZiZNHnQ8A==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=20"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/prelude-ls": {
       "version": "1.2.1",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.8.0"
       }
     },
+    "node_modules/pretty-ms": {
+      "version": "9.3.1",
+      "resolved": "https://registry.npmjs.org/pretty-ms/-/pretty-ms-9.3.1.tgz",
+      "integrity": "sha512-HzMy3Geq23nVALD/M2LliU+F+M+gVNsvkQWWqeBZ8HDiCgzo6YPJ/Omrmtq24EFrIsk0a3EkQGEd7bDOo+IhGA==",
+      "license": "MIT",
+      "dependencies": {
+        "parse-ms": "^4.0.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/prompts": {
+      "version": "2.4.2",
+      "resolved": "https://registry.npmjs.org/prompts/-/prompts-2.4.2.tgz",
+      "integrity": "sha512-NxNv/kLguCA7p3jE8oL2aEBsrJWgAakBpgmgK6lpPWV+WuOmY6r2/zbAVnP+T8bQlA0nzHXSJSJW0Hq7ylaD2Q==",
+      "license": "MIT",
+      "dependencies": {
+        "kleur": "^3.0.3",
+        "sisteransi": "^1.0.5"
+      },
+      "engines": {
+        "node": ">= 6"
+      }
+    },
+    "node_modules/prompts/node_modules/kleur": {
+      "version": "3.0.3",
+      "resolved": "https://registry.npmjs.org/kleur/-/kleur-3.0.3.tgz",
+      "integrity": "sha512-eTIzlVOSUR+JxdDFepEYcBMtZ9Qqdef+rnzWdRZuMbOywu5tO2w2N7rqjoANZ5k9vywhL6Br1VRjUIgTQx4E8w==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=6"
+      }
+    },
     "node_modules/prop-types": {
       "version": "15.8.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "loose-envify": "^1.4.0",
         "object-assign": "^4.1.1",
         "react-is": "^16.13.1"
       }
     },
+    "node_modules/proxy-addr": {
+      "version": "2.0.7",
+      "resolved": "https://registry.npmjs.org/proxy-addr/-/proxy-addr-2.0.7.tgz",
+      "integrity": "sha512-llQsMLSUDUPT44jdrU/O37qlnifitDP+ZwrmmZcoSKyLKvtZxpyV0n2/bD/N4tBAAZ/gJEdZU7KMraoK1+XYAg==",
+      "license": "MIT",
+      "dependencies": {
+        "forwarded": "0.2.0",
+        "ipaddr.js": "1.9.1"
+      },
+      "engines": {
+        "node": ">= 0.10"
+      }
+    },
     "node_modules/punycode": {
       "version": "2.3.1",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=6"
       }
     },
+    "node_modules/qs": {
+      "version": "6.16.0",
+      "resolved": "https://registry.npmjs.org/qs/-/qs-6.16.0.tgz",
+      "integrity": "sha512-h6fhOIaRrID2CbEY2fqs+7t+UXZo+MLAnU5gRIq85uFtdiUPCdsApMlHhXogKVM4HM2DVbIjGNTTYH2OcmP1vA==",
+      "license": "BSD-3-Clause",
+      "dependencies": {
+        "es-define-property": "^1.0.1",
+        "side-channel": "^1.1.1"
+      },
+      "engines": {
+        "node": ">=0.6"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/ljharb"
+      }
+    },
     "node_modules/queue-microtask": {
       "version": "1.2.3",
-      "dev": true,
       "funding": [
         {
           "type": "github",
           "url": "https://github.com/sponsors/feross"
         },
         {
           "type": "patreon",
           "url": "https://www.patreon.com/feross"
         },
         {
           "type": "consulting",
           "url": "https://feross.org/support"
         }
       ],
       "license": "MIT"
     },
+    "node_modules/range-parser": {
+      "version": "1.3.0",
+      "resolved": "https://registry.npmjs.org/range-parser/-/range-parser-1.3.0.tgz",
+      "integrity": "sha512-hek2mFQpPuI4E1BBKrSto+BU3e3x4xuarsbiwr3+lf7p44juvFMV0XFWQAP3xUyqXA4RrXLIoaSUGbSt056ZMw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.6"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
+    "node_modules/raw-body": {
+      "version": "3.0.2",
+      "resolved": "https://registry.npmjs.org/raw-body/-/raw-body-3.0.2.tgz",
+      "integrity": "sha512-K5zQjDllxWkf7Z5xJdV0/B0WTNqx6vxG70zJE4N0kBs4LovmEYWJzQGxC9bS9RAKu3bgM40lrd5zoLJ12MQ5BA==",
+      "license": "MIT",
+      "dependencies": {
+        "bytes": "~3.1.2",
+        "http-errors": "~2.0.1",
+        "iconv-lite": "~0.7.0",
+        "unpipe": "~1.0.0"
+      },
+      "engines": {
+        "node": ">= 0.10"
+      }
+    },
     "node_modules/react": {
       "version": "19.2.8",
       "license": "MIT",
       "engines": {
         "node": ">=0.10.0"
       }
     },
     "node_modules/react-dom": {
       "version": "19.2.8",
       "license": "MIT",
@@ -5166,20 +7483,36 @@
       },
       "peerDependencies": {
         "react": "^19.2.8"
       }
     },
     "node_modules/react-is": {
       "version": "16.13.1",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/recast": {
+      "version": "0.23.21",
+      "resolved": "https://registry.npmjs.org/recast/-/recast-0.23.21.tgz",
+      "integrity": "sha512-mFAyJq9vUbSTARLZUvAEf1z3YxlvAwswbmxMx2mPA/MSm4KmpwvwvhsH/NIrZhyOuwD60Lzyw2qh83uCbgTPYw==",
+      "license": "MIT",
+      "dependencies": {
+        "ast-types": "^0.16.1",
+        "esprima": "~4.0.0",
+        "source-map": "~0.6.1",
+        "tiny-invariant": "^1.3.3",
+        "tslib": "^2.0.1"
+      },
+      "engines": {
+        "node": ">= 4"
+      }
+    },
     "node_modules/reflect.getprototypeof": {
       "version": "1.0.10",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bind": "^1.0.8",
         "define-properties": "^1.2.1",
         "es-abstract": "^1.23.9",
         "es-errors": "^1.3.0",
         "es-object-atoms": "^1.0.0",
@@ -5206,20 +7539,35 @@
         "gopd": "^1.2.0",
         "set-function-name": "^2.0.2"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/require-from-string": {
+      "version": "2.0.2",
+      "resolved": "https://registry.npmjs.org/require-from-string/-/require-from-string-2.0.2.tgz",
+      "integrity": "sha512-Xf0nWe6RseziFMu+Ap9biiUbmplq6S9/p+7w7YXP/JBHhrUDDUhwa+vANyubuqfZWTveU//DYVGsDG7RKL/vEw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=0.10.0"
+      }
+    },
+    "node_modules/reselect": {
+      "version": "5.3.0",
+      "resolved": "https://registry.npmjs.org/reselect/-/reselect-5.3.0.tgz",
+      "integrity": "sha512-XGoLeRAVzUTcJ1qkxPQhDJyIZ5d6zzZD9nT7AEZOaaU9UbWclhycElmhO+VD5bFeLuzhPBaOV2oXC8uG35ZSpg==",
+      "license": "MIT"
+    },
     "node_modules/resolve": {
       "version": "2.0.0-next.7",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "es-errors": "^1.3.0",
         "is-core-module": "^2.16.2",
         "node-exports-info": "^1.6.0",
         "object-keys": "^1.1.1",
         "path-parse": "^1.0.7",
@@ -5230,46 +7578,102 @@
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/resolve-from": {
       "version": "4.0.0",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=4"
       }
     },
     "node_modules/resolve-pkg-maps": {
       "version": "1.0.0",
       "dev": true,
       "license": "MIT",
       "funding": {
         "url": "https://github.com/privatenumber/resolve-pkg-maps?sponsor=1"
       }
     },
+    "node_modules/restore-cursor": {
+      "version": "5.1.0",
+      "resolved": "https://registry.npmjs.org/restore-cursor/-/restore-cursor-5.1.0.tgz",
+      "integrity": "sha512-oMA2dcrw6u0YfxJQXm342bFKX/E4sG9rbTzO9ptUcR/e8A33cHuvStiYOwH7fszkZlZ1z/ta9AAoPk2F4qIOHA==",
+      "license": "MIT",
+      "dependencies": {
+        "onetime": "^7.0.0",
+        "signal-exit": "^4.1.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/restore-cursor/node_modules/onetime": {
+      "version": "7.0.0",
+      "resolved": "https://registry.npmjs.org/onetime/-/onetime-7.0.0.tgz",
+      "integrity": "sha512-VXJjc87FScF88uafS3JllDgvAm+c/Slfz06lorj2uAY34rlUu0Nt+v8wreiImcrgAjjIHp1rXpTDlLOGw29WwQ==",
+      "license": "MIT",
+      "dependencies": {
+        "mimic-function": "^5.0.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/reusify": {
       "version": "1.1.0",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "iojs": ">=1.0.0",
         "node": ">=0.10.0"
       }
     },
+    "node_modules/router": {
+      "version": "2.2.0",
+      "resolved": "https://registry.npmjs.org/router/-/router-2.2.0.tgz",
+      "integrity": "sha512-nLTrUKm2UyiL7rlhapu/Zl45FwNgkZGaCpZbIHajDYgwlJCOzLSk+cIPAnsEqV955GjILJnKbdQC1nVPz+gAYQ==",
+      "license": "MIT",
+      "dependencies": {
+        "debug": "^4.4.0",
+        "depd": "^2.0.0",
+        "is-promise": "^4.0.0",
+        "parseurl": "^1.3.3",
+        "path-to-regexp": "^8.0.0"
+      },
+      "engines": {
+        "node": ">= 18"
+      }
+    },
+    "node_modules/run-applescript": {
+      "version": "7.1.0",
+      "resolved": "https://registry.npmjs.org/run-applescript/-/run-applescript-7.1.0.tgz",
+      "integrity": "sha512-DPe5pVFaAsinSaV6QjQ6gdiedWDcRCbUuiQfQa2wmWV7+xC9bGulGI8+TdRmoFkAPaBXk8CrAbnlY2ISniJ47Q==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/run-parallel": {
       "version": "1.2.0",
-      "dev": true,
       "funding": [
         {
           "type": "github",
           "url": "https://github.com/sponsors/feross"
         },
         {
           "type": "patreon",
           "url": "https://www.patreon.com/feross"
         },
         {
@@ -5324,32 +7728,82 @@
         "es-errors": "^1.3.0",
         "is-regex": "^1.2.1"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/safer-buffer": {
+      "version": "2.1.2",
+      "resolved": "https://registry.npmjs.org/safer-buffer/-/safer-buffer-2.1.2.tgz",
+      "integrity": "sha512-YZo3K82SD7Riyi0E1EQPojLz7kpepnSQI9IyPbHHg1XXXevb5dJI7tpyN2ADxGcQbHG7vcyRHk0cbwqcQriUtg==",
+      "license": "MIT"
+    },
     "node_modules/scheduler": {
       "version": "0.27.0",
       "license": "MIT"
     },
     "node_modules/semver": {
       "version": "6.3.1",
-      "dev": true,
       "license": "ISC",
       "bin": {
         "semver": "bin/semver.js"
       }
     },
+    "node_modules/send": {
+      "version": "1.2.1",
+      "resolved": "https://registry.npmjs.org/send/-/send-1.2.1.tgz",
+      "integrity": "sha512-1gnZf7DFcoIcajTjTwjwuDjzuz4PPcY2StKPlsGAQ1+YH20IRVrBaXSWmdjowTJ6u8Rc01PoYOGHXfP1mYcZNQ==",
+      "license": "MIT",
+      "dependencies": {
+        "debug": "^4.4.3",
+        "encodeurl": "^2.0.0",
+        "escape-html": "^1.0.3",
+        "etag": "^1.8.1",
+        "fresh": "^2.0.0",
+        "http-errors": "^2.0.1",
+        "mime-types": "^3.0.2",
+        "ms": "^2.1.3",
+        "on-finished": "^2.4.1",
+        "range-parser": "^1.2.1",
+        "statuses": "^2.0.2"
+      },
+      "engines": {
+        "node": ">= 18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
+    "node_modules/serve-static": {
+      "version": "2.2.1",
+      "resolved": "https://registry.npmjs.org/serve-static/-/serve-static-2.2.1.tgz",
+      "integrity": "sha512-xRXBn0pPqQTVQiC8wyQrKs2MOlX24zQ0POGaj0kultvoOCstBQM5yvOhAVSUwOMjQtTvsPWoNCHfPGwaaQJhTw==",
+      "license": "MIT",
+      "dependencies": {
+        "encodeurl": "^2.0.0",
+        "escape-html": "^1.0.3",
+        "parseurl": "^1.3.3",
+        "send": "^1.2.0"
+      },
+      "engines": {
+        "node": ">= 18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
     "node_modules/set-function-length": {
       "version": "1.2.2",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "define-data-property": "^1.1.4",
         "es-errors": "^1.3.0",
         "function-bind": "^1.1.2",
         "get-intrinsic": "^1.2.4",
         "gopd": "^1.0.1",
@@ -5379,20 +7833,124 @@
       "license": "MIT",
       "dependencies": {
         "dunder-proto": "^1.0.1",
         "es-errors": "^1.3.0",
         "es-object-atoms": "^1.0.0"
       },
       "engines": {
         "node": ">= 0.4"
       }
     },
+    "node_modules/setprototypeof": {
+      "version": "1.2.0",
+      "resolved": "https://registry.npmjs.org/setprototypeof/-/setprototypeof-1.2.0.tgz",
+      "integrity": "sha512-E5LDX7Wrp85Kil5bhZv46j8jOeboKq5JMmYM3gVGdGH8xFpPWXUMsNrlODCrkoxMEeNi/XZIwuRvY4XNwYMJpw==",
+      "license": "ISC"
+    },
+    "node_modules/shadcn": {
+      "version": "4.19.0",
+      "resolved": "https://registry.npmjs.org/shadcn/-/shadcn-4.19.0.tgz",
+      "integrity": "sha512-EQF6R+CUXTsEP2BpyhxrUEAFesrtFD1POvVOf5jM+wkgtA4kG1EW1+1Wlmi9LqiprSL681JJXBcS0u8WkVVVyQ==",
+      "license": "MIT",
+      "dependencies": {
+        "@babel/core": "^7.28.0",
+        "@babel/parser": "^7.28.0",
+        "@babel/plugin-transform-typescript": "^7.28.0",
+        "@babel/preset-typescript": "^7.27.1",
+        "@dotenvx/dotenvx": "^1.48.4",
+        "@modelcontextprotocol/sdk": "^1.26.0",
+        "@types/validate-npm-package-name": "^4.0.2",
+        "browserslist": "^4.26.2",
+        "commander": "^14.0.0",
+        "cosmiconfig": "^9.0.0",
+        "dedent": "^1.6.0",
+        "deepmerge": "^4.3.1",
+        "diff": "^8.0.2",
+        "execa": "^9.6.0",
+        "fast-glob": "^3.3.3",
+        "fs-extra": "^11.3.1",
+        "fuzzysort": "^3.1.0",
+        "kleur": "^4.1.5",
+        "open": "^11.0.0",
+        "ora": "^8.2.0",
+        "postcss": "^8.5.6",
+        "postcss-selector-parser": "^7.1.0",
+        "prompts": "^2.4.2",
+        "recast": "^0.23.11",
+        "socks": "^2.8.8",
+        "stringify-object": "^5.0.0",
+        "tailwind-merge": "^3.0.1",
+        "ts-morph": "^26.0.0",
+        "tsconfig-paths": "^4.2.0",
+        "undici": "^7.27.2",
+        "validate-npm-package-name": "^7.0.1",
+        "zod": "^3.24.1",
+        "zod-to-json-schema": "^3.24.6"
+      },
+      "bin": {
+        "shadcn": "dist/index.js"
+      },
+      "engines": {
+        "node": ">=20.18.1"
+      }
+    },
+    "node_modules/shadcn/node_modules/fast-glob": {
+      "version": "3.3.3",
+      "resolved": "https://registry.npmjs.org/fast-glob/-/fast-glob-3.3.3.tgz",
+      "integrity": "sha512-7MptL8U0cqcFdzIzwOTHoilX9x5BrNqye7Z/LuC7kCMRio1EMSyqRK3BEAUD7sXRq4iT4AzTVuZdhgQ2TCvYLg==",
+      "license": "MIT",
+      "dependencies": {
+        "@nodelib/fs.stat": "^2.0.2",
+        "@nodelib/fs.walk": "^1.2.3",
+        "glob-parent": "^5.1.2",
+        "merge2": "^1.3.0",
+        "micromatch": "^4.0.8"
+      },
+      "engines": {
+        "node": ">=8.6.0"
+      }
+    },
+    "node_modules/shadcn/node_modules/glob-parent": {
+      "version": "5.1.2",
+      "resolved": "https://registry.npmjs.org/glob-parent/-/glob-parent-5.1.2.tgz",
+      "integrity": "sha512-AOIgSQCepiJYwP3ARnGx+5VnTu2HBYdzbGP45eLw1vr3zB3vZLeyed1sC9hnbcOc9/SrMyM5RPQrkGz4aS9Zow==",
+      "license": "ISC",
+      "dependencies": {
+        "is-glob": "^4.0.1"
+      },
+      "engines": {
+        "node": ">= 6"
+      }
+    },
+    "node_modules/shadcn/node_modules/tsconfig-paths": {
+      "version": "4.2.0",
+      "resolved": "https://registry.npmjs.org/tsconfig-paths/-/tsconfig-paths-4.2.0.tgz",
+      "integrity": "sha512-NoZ4roiN7LnbKn9QqE1amc9DJfzvZXxF4xDavcOWt1BPkdx+m+0gJuPM+S0vCe7zTJMYUP0R8pO2XMr+Y8oLIg==",
+      "license": "MIT",
+      "dependencies": {
+        "json5": "^2.2.2",
+        "minimist": "^1.2.6",
+        "strip-bom": "^3.0.0"
+      },
+      "engines": {
+        "node": ">=6"
+      }
+    },
+    "node_modules/shadcn/node_modules/zod": {
+      "version": "3.25.76",
+      "resolved": "https://registry.npmjs.org/zod/-/zod-3.25.76.tgz",
+      "integrity": "sha512-gzUt/qt81nXsFGKIFcC3YnfEAx5NkunCfnDlvuBSSFS02bcXu4Lmea0AFIUwbLWxWPx3d9p8S5QoaujKcNQxcQ==",
+      "license": "MIT",
+      "funding": {
+        "url": "https://github.com/sponsors/colinhacks"
+      }
+    },
     "node_modules/sharp": {
       "version": "0.35.4",
       "license": "Apache-2.0",
       "optional": true,
       "dependencies": {
         "@img/colour": "^1.1.0",
         "detect-libc": "^2.1.2",
         "semver": "^7.8.5"
       },
       "engines": {
@@ -5440,129 +7998,245 @@
       "optional": true,
       "bin": {
         "semver": "bin/semver.js"
       },
       "engines": {
         "node": ">=10"
       }
     },
     "node_modules/shebang-command": {
       "version": "2.0.0",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "shebang-regex": "^3.0.0"
       },
       "engines": {
         "node": ">=8"
       }
     },
     "node_modules/shebang-regex": {
       "version": "3.0.0",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=8"
       }
     },
     "node_modules/side-channel": {
       "version": "1.1.1",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "es-errors": "^1.3.0",
         "object-inspect": "^1.13.4",
         "side-channel-list": "^1.0.1",
         "side-channel-map": "^1.0.1",
         "side-channel-weakmap": "^1.0.2"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/side-channel-list": {
       "version": "1.0.1",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "es-errors": "^1.3.0",
         "object-inspect": "^1.13.4"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/side-channel-map": {
       "version": "1.0.1",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bound": "^1.0.2",
         "es-errors": "^1.3.0",
         "get-intrinsic": "^1.2.5",
         "object-inspect": "^1.13.3"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/side-channel-weakmap": {
       "version": "1.0.2",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bound": "^1.0.2",
         "es-errors": "^1.3.0",
         "get-intrinsic": "^1.2.5",
         "object-inspect": "^1.13.3",
         "side-channel-map": "^1.0.1"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/signal-exit": {
+      "version": "4.1.0",
+      "resolved": "https://registry.npmjs.org/signal-exit/-/signal-exit-4.1.0.tgz",
+      "integrity": "sha512-bzyZ1e88w9O1iNJbKnOlvYTrWPDl46O1bG0D3XInv+9tkPrxrN8jUUTiFlDkkmKWgn1M6CfIA13SuGqOa9Korw==",
+      "license": "ISC",
+      "engines": {
+        "node": ">=14"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/isaacs"
+      }
+    },
+    "node_modules/sisteransi": {
+      "version": "1.0.5",
+      "resolved": "https://registry.npmjs.org/sisteransi/-/sisteransi-1.0.5.tgz",
+      "integrity": "sha512-bLGGlR1QxBcynn2d5YmDX4MGjlZvy2MRBDRNHLJ8VI6l6+9FUiyTFNJ0IveOSP0bcXgVDPRcfGqA0pjaqUpfVg==",
+      "license": "MIT"
+    },
+    "node_modules/smart-buffer": {
+      "version": "4.2.0",
+      "resolved": "https://registry.npmjs.org/smart-buffer/-/smart-buffer-4.2.0.tgz",
+      "integrity": "sha512-94hK0Hh8rPqQl2xXc3HsaBoOXKV20MToPkcXvwbISWLEs+64sBq5kFgn2kJDHb1Pry9yrP0dxrCI9RRci7RXKg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 6.0.0",
+        "npm": ">= 3.0.0"
+      }
+    },
+    "node_modules/socks": {
+      "version": "2.8.9",
+      "resolved": "https://registry.npmjs.org/socks/-/socks-2.8.9.tgz",
+      "integrity": "sha512-LJhUYUvItdQ0LkJTmPeaEObWXAqFyfmP85x0tch/ez9cahmhlBBLbIqDFnvBnUJGagb0JbIQrkBs1wJ+yRYpEw==",
+      "license": "MIT",
+      "dependencies": {
+        "ip-address": "^10.1.1",
+        "smart-buffer": "^4.2.0"
+      },
+      "engines": {
+        "node": ">= 10.0.0",
+        "npm": ">= 3.0.0"
+      }
+    },
+    "node_modules/source-map": {
+      "version": "0.6.1",
+      "resolved": "https://registry.npmjs.org/source-map/-/source-map-0.6.1.tgz",
+      "integrity": "sha512-UjgapumWlbMhkBgzT7Ykc5YXUT46F0iKu8SGXq0bcwP5dz/h0Plj6enJqjz1Zbq2l5WaqYnrVbwWOWMyF3F47g==",
+      "license": "BSD-3-Clause",
+      "engines": {
+        "node": ">=0.10.0"
+      }
+    },
     "node_modules/source-map-js": {
       "version": "1.2.1",
       "license": "BSD-3-Clause",
       "engines": {
         "node": ">=0.10.0"
       }
     },
     "node_modules/stable-hash": {
       "version": "0.0.5",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/statuses": {
+      "version": "2.0.2",
+      "resolved": "https://registry.npmjs.org/statuses/-/statuses-2.0.2.tgz",
+      "integrity": "sha512-DvEy55V3DB7uknRo+4iOGT5fP1slR8wQohVdknigZPMpMstaKJQWhwiYBACJE3Ul2pTnATihhBYnRhZQHGBiRw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.8"
+      }
+    },
+    "node_modules/stdin-discarder": {
+      "version": "0.2.2",
+      "resolved": "https://registry.npmjs.org/stdin-discarder/-/stdin-discarder-0.2.2.tgz",
+      "integrity": "sha512-UhDfHmA92YAlNnCfhmq0VeNL5bDbiZGg7sZ2IvPsXubGkiNa9EC+tUTsjBRsYUAz87btI6/1wf4XoVvQ3uRnmQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/stop-iteration-iterator": {
       "version": "1.1.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "es-errors": "^1.3.0",
         "internal-slot": "^1.1.0"
       },
       "engines": {
         "node": ">= 0.4"
       }
     },
+    "node_modules/string-width": {
+      "version": "7.2.0",
+      "resolved": "https://registry.npmjs.org/string-width/-/string-width-7.2.0.tgz",
+      "integrity": "sha512-tsaTIkKW9b4N+AEj+SVA+WhJzV7/zMhcSu78mLKWSk7cXMOSHsBKFWUs0fWwq8QyK3MgJBQRX6Gbi4kYbdvGkQ==",
+      "license": "MIT",
+      "dependencies": {
+        "emoji-regex": "^10.3.0",
+        "get-east-asian-width": "^1.0.0",
+        "strip-ansi": "^7.1.0"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/string-width/node_modules/ansi-regex": {
+      "version": "6.3.0",
+      "resolved": "https://registry.npmjs.org/ansi-regex/-/ansi-regex-6.3.0.tgz",
+      "integrity": "sha512-WpDfL7NO6j7tH88IDBNVdUJxDh9nmCteAVW9dsep846XdwF4naCBK+/tGLX3KJgcpgMRXCFlTM2hKGoK9FsdrQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/chalk/ansi-regex?sponsor=1"
+      }
+    },
+    "node_modules/string-width/node_modules/emoji-regex": {
+      "version": "10.6.0",
+      "resolved": "https://registry.npmjs.org/emoji-regex/-/emoji-regex-10.6.0.tgz",
+      "integrity": "sha512-toUI84YS5YmxW219erniWD0CIVOo46xGKColeNQRgOzDorgBi1v4D71/OFzgD9GO2UGKIv1C3Sp8DAn0+j5w7A==",
+      "license": "MIT"
+    },
+    "node_modules/string-width/node_modules/strip-ansi": {
+      "version": "7.2.0",
+      "resolved": "https://registry.npmjs.org/strip-ansi/-/strip-ansi-7.2.0.tgz",
+      "integrity": "sha512-yDPMNjp4WyfYBkHnjIRLfca1i6KMyGCtsVgoKe/z1+6vukgaENdgGBZt+ZmKPc4gavvEZ5OgHfHdrazhgNyG7w==",
+      "license": "MIT",
+      "dependencies": {
+        "ansi-regex": "^6.2.2"
+      },
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/chalk/strip-ansi?sponsor=1"
+      }
+    },
     "node_modules/string.prototype.includes": {
       "version": "2.0.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bind": "^1.0.7",
         "define-properties": "^1.2.1",
         "es-abstract": "^1.23.3"
       },
       "engines": {
@@ -5651,28 +8325,80 @@
         "define-properties": "^1.2.1",
         "es-object-atoms": "^1.0.0"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/stringify-object": {
+      "version": "5.0.0",
+      "resolved": "https://registry.npmjs.org/stringify-object/-/stringify-object-5.0.0.tgz",
+      "integrity": "sha512-zaJYxz2FtcMb4f+g60KsRNFOpVMUyuJgA51Zi5Z1DOTC3S59+OQiVOzE9GZt0x72uBGWKsQIuBKeF9iusmKFsg==",
+      "license": "BSD-2-Clause",
+      "dependencies": {
+        "get-own-enumerable-keys": "^1.0.0",
+        "is-obj": "^3.0.0",
+        "is-regexp": "^3.1.0"
+      },
+      "engines": {
+        "node": ">=14.16"
+      },
+      "funding": {
+        "url": "https://github.com/yeoman/stringify-object?sponsor=1"
+      }
+    },
+    "node_modules/stringify-object/node_modules/is-obj": {
+      "version": "3.0.0",
+      "resolved": "https://registry.npmjs.org/is-obj/-/is-obj-3.0.0.tgz",
+      "integrity": "sha512-IlsXEHOjtKhpN8r/tRFj2nDyTmHvcfNeu/nrRIcXE17ROeatXchkojffa1SpdqW4cr/Fj6QkEf/Gn4zf6KKvEQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=12"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/strip-ansi": {
+      "version": "6.0.1",
+      "resolved": "https://registry.npmjs.org/strip-ansi/-/strip-ansi-6.0.1.tgz",
+      "integrity": "sha512-Y38VPSHcqkFrCpFnQ9vuSXmquuv5oXOKpGeT6aGrr3o3Gc9AlVa6JBfUSOCnbxGGZF+/0ooI7KrPuUSztUdU5A==",
+      "license": "MIT",
+      "dependencies": {
+        "ansi-regex": "^5.0.1"
+      },
+      "engines": {
+        "node": ">=8"
+      }
+    },
     "node_modules/strip-bom": {
       "version": "3.0.0",
-      "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=4"
       }
     },
+    "node_modules/strip-final-newline": {
+      "version": "4.0.0",
+      "resolved": "https://registry.npmjs.org/strip-final-newline/-/strip-final-newline-4.0.0.tgz",
+      "integrity": "sha512-aulFJcD6YK8V1G7iRB5tigAP4TsHBZZrOV8pjV++zdUwmeV8uzbY7yn6h9MswN62adStNZFuCIx4haBnRuMDaw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/strip-json-comments": {
       "version": "3.1.1",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=8"
       },
       "funding": {
         "url": "https://github.com/sponsors/sindresorhus"
       }
@@ -5713,37 +8439,79 @@
       "version": "1.0.0",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/systeminformation": {
+      "version": "5.33.6",
+      "resolved": "https://registry.npmjs.org/systeminformation/-/systeminformation-5.33.6.tgz",
+      "integrity": "sha512-hMOQG/eRUzuopuYGGdl8ntkau0nEC7fOaRoTUg1RSr2GTQIk2VNa76DA0+ApajkGfzmcgAupgIP/vt+jtoe5EA==",
+      "license": "MIT",
+      "os": [
+        "darwin",
+        "linux",
+        "win32",
+        "freebsd",
+        "openbsd",
+        "netbsd",
+        "sunos",
+        "android"
+      ],
+      "bin": {
+        "systeminformation": "lib/cli.js"
+      },
+      "engines": {
+        "node": ">=10.0.0"
+      },
+      "funding": {
+        "type": "Buy me a coffee",
+        "url": "https://www.buymeacoffee.com/systeminfo"
+      }
+    },
+    "node_modules/tailwind-merge": {
+      "version": "3.6.0",
+      "resolved": "https://registry.npmjs.org/tailwind-merge/-/tailwind-merge-3.6.0.tgz",
+      "integrity": "sha512-uxL7qAVQriqRQPAyK3pj66VqskWqoZ37PW94jwOTwNfq/z9oyu1V+eqrZqtR2+fCiXdYOZe/Modt8GtvqNzu+w==",
+      "license": "MIT",
+      "funding": {
+        "type": "github",
+        "url": "https://github.com/sponsors/dcastil"
+      }
+    },
     "node_modules/tailwindcss": {
       "version": "4.3.3",
       "dev": true,
       "license": "MIT"
     },
     "node_modules/tapable": {
       "version": "2.3.3",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=6"
       },
       "funding": {
         "type": "opencollective",
         "url": "https://opencollective.com/webpack"
       }
     },
+    "node_modules/tiny-invariant": {
+      "version": "1.3.3",
+      "resolved": "https://registry.npmjs.org/tiny-invariant/-/tiny-invariant-1.3.3.tgz",
+      "integrity": "sha512-+FbBPE1o9QAYvviau/qC5SE3caw21q3xkvWKBtja5vgqOWIHHJ3ioaq1VPfn/Szqctz2bU/oYeKd9/z5BL+PVg==",
+      "license": "MIT"
+    },
     "node_modules/tinyglobby": {
       "version": "0.2.17",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "fdir": "^6.5.0",
         "picomatch": "^4.0.4"
       },
       "engines": {
         "node": ">=12.0.0"
@@ -5774,40 +8542,58 @@
       "license": "MIT",
       "engines": {
         "node": ">=12"
       },
       "funding": {
         "url": "https://github.com/sponsors/jonschlinkert"
       }
     },
     "node_modules/to-regex-range": {
       "version": "5.0.1",
-      "dev": true,
       "license": "MIT",
       "dependencies": {
         "is-number": "^7.0.0"
       },
       "engines": {
         "node": ">=8.0"
       }
     },
+    "node_modules/toidentifier": {
+      "version": "1.0.1",
+      "resolved": "https://registry.npmjs.org/toidentifier/-/toidentifier-1.0.1.tgz",
+      "integrity": "sha512-o5sSPKEkg/DIQNmH43V0/uerLrpzVedkUh8tGNvaeXpfpuwjKenlSox/2O/BTlZUtEe+JG7s5YhEz608PlAHRA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=0.6"
+      }
+    },
     "node_modules/ts-api-utils": {
       "version": "2.5.0",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=18.12"
       },
       "peerDependencies": {
         "typescript": ">=4.8.4"
       }
     },
+    "node_modules/ts-morph": {
+      "version": "26.0.0",
+      "resolved": "https://registry.npmjs.org/ts-morph/-/ts-morph-26.0.0.tgz",
+      "integrity": "sha512-ztMO++owQnz8c/gIENcM9XfCEzgoGphTv+nKpYNM1bgsdOVC/jRZuEBf6N+mLLDNg68Kl+GgUZfOySaRiG1/Ug==",
+      "license": "MIT",
+      "dependencies": {
+        "@ts-morph/common": "~0.27.0",
+        "code-block-writer": "^13.0.3"
+      }
+    },
     "node_modules/tsconfig-paths": {
       "version": "3.15.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "@types/json5": "^0.0.29",
         "json5": "^1.0.2",
         "minimist": "^1.2.6",
         "strip-bom": "^3.0.0"
       }
@@ -5820,31 +8606,71 @@
         "minimist": "^1.2.0"
       },
       "bin": {
         "json5": "lib/cli.js"
       }
     },
     "node_modules/tslib": {
       "version": "2.8.1",
       "license": "0BSD"
     },
+    "node_modules/tw-animate-css": {
+      "version": "1.4.0",
+      "resolved": "https://registry.npmjs.org/tw-animate-css/-/tw-animate-css-1.4.0.tgz",
+      "integrity": "sha512-7bziOlRqH0hJx80h/3mbicLW7o8qLsH5+RaLR2t+OHM3D0JlWGODQKQ4cxbK7WlvmUxpcj6Kgu6EKqjrGFe3QQ==",
+      "license": "MIT",
+      "funding": {
+        "url": "https://github.com/sponsors/Wombosvideo"
+      }
+    },
     "node_modules/type-check": {
       "version": "0.4.0",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "prelude-ls": "^1.2.1"
       },
       "engines": {
         "node": ">= 0.8.0"
       }
     },
+    "node_modules/type-is": {
+      "version": "2.1.0",
+      "resolved": "https://registry.npmjs.org/type-is/-/type-is-2.1.0.tgz",
+      "integrity": "sha512-faYHw0anBbc/kWF3zFTEnxSFOAGUX9GFbOBthvDdLsIlEoWOFOtS0zgCiQYwIskL9iGXZL3kAXD8OoZ4GmMATA==",
+      "license": "MIT",
+      "dependencies": {
+        "content-type": "^2.0.0",
+        "media-typer": "^1.1.0",
+        "mime-types": "^3.0.0"
+      },
+      "engines": {
+        "node": ">= 18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
+    "node_modules/type-is/node_modules/content-type": {
+      "version": "2.1.0",
+      "resolved": "https://registry.npmjs.org/content-type/-/content-type-2.1.0.tgz",
+      "integrity": "sha512-mj7UPXE0jaqaOsukNZRUEfEi2AcL7C/vwmwcHV0O97eO1E1pxBZuyjlZrx5seTaNBg1U6+o35wpa35Qfcc+7ag==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
     "node_modules/typed-array-buffer": {
       "version": "1.0.3",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "call-bound": "^1.0.3",
         "es-errors": "^1.3.0",
         "is-typed-array": "^1.1.14"
       },
       "engines": {
@@ -5903,21 +8729,21 @@
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/typescript": {
       "version": "5.9.3",
-      "dev": true,
+      "devOptional": true,
       "license": "Apache-2.0",
       "bin": {
         "tsc": "bin/tsc",
         "tsserver": "bin/tsserver"
       },
       "engines": {
         "node": ">=14.17"
       }
     },
     "node_modules/typescript-eslint": {
@@ -5952,25 +8778,64 @@
         "has-symbols": "^1.1.0",
         "which-boxed-primitive": "^1.1.1"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/undici": {
+      "version": "7.29.0",
+      "resolved": "https://registry.npmjs.org/undici/-/undici-7.29.0.tgz",
+      "integrity": "sha512-IDxfleLmmbSskfWSUATiN1nfn2rDuvnMOqb5CWR92iIfojA0Ud+ulOAAEQ57LPr9rWmsreUyf5lwyao+7GNNVw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=20.18.1"
+      }
+    },
     "node_modules/undici-types": {
       "version": "6.21.0",
       "dev": true,
       "license": "MIT"
     },
+    "node_modules/unicorn-magic": {
+      "version": "0.3.0",
+      "resolved": "https://registry.npmjs.org/unicorn-magic/-/unicorn-magic-0.3.0.tgz",
+      "integrity": "sha512-+QBBXBCvifc56fsbuxZQ6Sic3wqqc3WWaqxs58gvJrcOuN83HGTCwz3oS5phzU9LthRNE9VrJCFCLUgHeeFnfA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/universalify": {
+      "version": "2.0.1",
+      "resolved": "https://registry.npmjs.org/universalify/-/universalify-2.0.1.tgz",
+      "integrity": "sha512-gptHNQghINnc/vTGIk0SOFGFNXw7JVrlRUtConJRlvaw6DuX0wO5Jeko9sWrMBhh+PsYAZ7oXAiOnf/UKogyiw==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 10.0.0"
+      }
+    },
+    "node_modules/unpipe": {
+      "version": "1.0.0",
+      "resolved": "https://registry.npmjs.org/unpipe/-/unpipe-1.0.0.tgz",
+      "integrity": "sha512-pjy2bYhSsufwWlKwPc+l3cN7+wuJlK6uz0YdJEOlQDbl6jo/YlPi4mb8agUkVC8BF7V8NuzeyPNqRksA3hztKQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.8"
+      }
+    },
     "node_modules/unrs-resolver": {
       "version": "1.12.2",
       "dev": true,
       "hasInstallScript": true,
       "license": "MIT",
       "dependencies": {
         "napi-postinstall": "^0.3.4"
       },
       "funding": {
         "url": "https://opencollective.com/unrs-resolver"
@@ -5995,21 +8860,20 @@
         "@unrs/resolver-binding-linux-x64-musl": "1.12.2",
         "@unrs/resolver-binding-openharmony-arm64": "1.12.2",
         "@unrs/resolver-binding-wasm32-wasi": "1.12.2",
         "@unrs/resolver-binding-win32-arm64-msvc": "1.12.2",
         "@unrs/resolver-binding-win32-ia32-msvc": "1.12.2",
         "@unrs/resolver-binding-win32-x64-msvc": "1.12.2"
       }
     },
     "node_modules/update-browserslist-db": {
       "version": "1.3.2",
-      "dev": true,
       "funding": [
         {
           "type": "opencollective",
           "url": "https://opencollective.com/browserslist"
         },
         {
           "type": "tidelift",
           "url": "https://tidelift.com/funding/github/npm/browserslist"
         },
         {
@@ -6030,23 +8894,55 @@
       }
     },
     "node_modules/uri-js": {
       "version": "4.4.1",
       "dev": true,
       "license": "BSD-2-Clause",
       "dependencies": {
         "punycode": "^2.1.0"
       }
     },
+    "node_modules/use-sync-external-store": {
+      "version": "1.6.0",
+      "resolved": "https://registry.npmjs.org/use-sync-external-store/-/use-sync-external-store-1.6.0.tgz",
+      "integrity": "sha512-Pp6GSwGP/NrPIrxVFAIkOQeyw8lFenOHijQWkUTrDvrF4ALqylP2C/KCkeS9dpUM3KvYRQhna5vt7IL95+ZQ9w==",
+      "license": "MIT",
+      "peerDependencies": {
+        "react": "^16.8.0 || ^17.0.0 || ^18.0.0 || ^19.0.0"
+      }
+    },
+    "node_modules/util-deprecate": {
+      "version": "1.0.2",
+      "resolved": "https://registry.npmjs.org/util-deprecate/-/util-deprecate-1.0.2.tgz",
+      "integrity": "sha512-EPD5q1uXyFxJpCrLnCc1nHnq3gOa6DZBocAIiI2TaSCA7VCJ1UJDMagCzIkXNsUYfD1daK//LTEQ8xiIbrHtcw==",
+      "license": "MIT"
+    },
+    "node_modules/validate-npm-package-name": {
+      "version": "7.0.2",
+      "resolved": "https://registry.npmjs.org/validate-npm-package-name/-/validate-npm-package-name-7.0.2.tgz",
+      "integrity": "sha512-hVDIBwsRruT73PbK7uP5ebUt+ezEtCmzZz3F59BSr2F6OVFnJ/6h8liuvdLrQ88Xmnk6/+xGGuq+pG9WwTuy3A==",
+      "license": "ISC",
+      "engines": {
+        "node": "^20.17.0 || >=22.9.0"
+      }
+    },
+    "node_modules/vary": {
+      "version": "1.1.2",
+      "resolved": "https://registry.npmjs.org/vary/-/vary-1.1.2.tgz",
+      "integrity": "sha512-BNGbWLfd0eUPabhkXUVm0j8uuvREyTh5ovRa/dyow/BqAbZJyC+5fU+IzQOzmAKzYqYRAISoRhdQr3eIZ/PXqg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">= 0.8"
+      }
+    },
     "node_modules/which": {
       "version": "2.0.2",
-      "dev": true,
       "license": "ISC",
       "dependencies": {
         "isexe": "^2.0.0"
       },
       "bin": {
         "node-which": "bin/node-which"
       },
       "engines": {
         "node": ">= 8"
       }
@@ -6133,44 +9029,112 @@
       }
     },
     "node_modules/word-wrap": {
       "version": "1.2.5",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=0.10.0"
       }
     },
+    "node_modules/wrappy": {
+      "version": "1.0.2",
+      "resolved": "https://registry.npmjs.org/wrappy/-/wrappy-1.0.2.tgz",
+      "integrity": "sha512-l4Sp/DRseor9wL6EvV2+TuQn63dMkPjZ/sp9XkghTEbV9KlPS1xUsZ3u7/IQO4wxtcFB4bgpQPRcR3QCvezPcQ==",
+      "license": "ISC"
+    },
+    "node_modules/wsl-utils": {
+      "version": "1.0.0",
+      "resolved": "https://registry.npmjs.org/wsl-utils/-/wsl-utils-1.0.0.tgz",
+      "integrity": "sha512-Hl0ZOAs672vg+06kfujwRhoS6/jehvULrlFkuF2dRu6pHgA8U06h3xqNIqNNU1LTXPcedxByAR4GS6pwQK0mgA==",
+      "license": "MIT",
+      "dependencies": {
+        "is-wsl": "^3.1.0",
+        "powershell-utils": "^0.1.0"
+      },
+      "engines": {
+        "node": ">=20"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/wsl-utils/node_modules/powershell-utils": {
+      "version": "0.1.0",
+      "resolved": "https://registry.npmjs.org/powershell-utils/-/powershell-utils-0.1.0.tgz",
+      "integrity": "sha512-dM0jVuXJPsDN6DvRpea484tCUaMiXWjuCn++HGTqUWzGDjv5tZkEZldAJ/UMlqRYGFrD/etByo4/xOuC/snX2A==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=20"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/yallist": {
       "version": "3.1.1",
-      "dev": true,
       "license": "ISC"
     },
     "node_modules/yocto-queue": {
       "version": "0.1.0",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=10"
       },
       "funding": {
         "url": "https://github.com/sponsors/sindresorhus"
       }
     },
+    "node_modules/yocto-spinner": {
+      "version": "1.2.2",
+      "resolved": "https://registry.npmjs.org/yocto-spinner/-/yocto-spinner-1.2.2.tgz",
+      "integrity": "sha512-DODGl1wJjA/s5pnJFKau9lIYHT81lnhob1i3e1TjxZRxEhWRKl74nTbWE6H5KlkViQQTo/Z29YFdxzTZAMY3ng==",
+      "license": "MIT",
+      "dependencies": {
+        "yoctocolors": "^2.1.1"
+      },
+      "engines": {
+        "node": ">=18.19"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
+    "node_modules/yoctocolors": {
+      "version": "2.2.0",
+      "resolved": "https://registry.npmjs.org/yoctocolors/-/yoctocolors-2.2.0.tgz",
+      "integrity": "sha512-xYqdZFUK/VYazNl/oCDYN+3WloWQwMfZxBoiNt6qNyk+xfOdi598muWE42rNZFp1kNOiqW936q5RhUdnpqElSg==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "url": "https://github.com/sponsors/sindresorhus"
+      }
+    },
     "node_modules/zod": {
       "version": "4.5.4",
-      "dev": true,
       "license": "MIT",
       "funding": {
         "url": "https://github.com/sponsors/colinhacks"
       }
     },
+    "node_modules/zod-to-json-schema": {
+      "version": "3.25.2",
+      "resolved": "https://registry.npmjs.org/zod-to-json-schema/-/zod-to-json-schema-3.25.2.tgz",
+      "integrity": "sha512-O/PgfnpT1xKSDeQYSCfRI5Gy3hPf91mKVDuYLUHZJMiDFptvP41MSnWofm8dnCm0256ZNfZIM7DSzuSMAFnjHA==",
+      "license": "ISC",
+      "peerDependencies": {
+        "zod": "^3.25.28 || ^4"
+      }
+    },
     "node_modules/zod-validation-error": {
       "version": "4.0.2",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=18.0.0"
       },
       "peerDependencies": {
         "zod": "^3.25.0 || ^4.0.0"
       }
diff --git a/web/package.json b/web/package.json
index fe0eeb4..455e64c 100644
--- a/web/package.json
+++ b/web/package.json
@@ -2,23 +2,31 @@
   "name": "web",
   "version": "0.1.0",
   "private": true,
   "scripts": {
     "dev": "next dev",
     "build": "next build",
     "start": "next start",
     "lint": "eslint"
   },
   "dependencies": {
+    "@base-ui/react": "^1.7.0",
+    "class-variance-authority": "^0.7.1",
+    "clsx": "^2.1.1",
+    "lucide-react": "^1.37.0",
     "next": "16.3.3",
+    "next-themes": "^0.4.6",
     "react": "19.2.8",
-    "react-dom": "19.2.8"
+    "react-dom": "19.2.8",
+    "shadcn": "^4.19.0",
+    "tailwind-merge": "^3.6.0",
+    "tw-animate-css": "^1.4.0"
   },
   "devDependencies": {
     "@tailwindcss/postcss": "^4",
     "@types/node": "^20",
     "@types/react": "^19",
     "@types/react-dom": "^19",
     "eslint": "^9",
     "eslint-config-next": "16.3.3",
     "tailwindcss": "^4",
     "typescript": "^5"
diff --git a/web/src/app/globals.css b/web/src/app/globals.css
index a2dc41e..1d301eb 100644
--- a/web/src/app/globals.css
+++ b/web/src/app/globals.css
@@ -1,26 +1,131 @@
 @import "tailwindcss";
+@import "tw-animate-css";
+@import "shadcn/tailwind.css";
+@import "../styles/themes/index.css";
 
-:root {
-  --background: #ffffff;
-  --foreground: #171717;
-}
+@custom-variant dark (&:is([data-theme$="-dark"] *));
 
 @theme inline {
   --color-background: var(--background);
   --color-foreground: var(--foreground);
   --font-sans: var(--font-geist-sans);
   --font-mono: var(--font-geist-mono);
+  --font-heading: var(--font-sans);
+  --color-sidebar-ring: var(--sidebar-ring);
+  --color-sidebar-border: var(--sidebar-border);
+  --color-sidebar-accent-foreground: var(--sidebar-accent-foreground);
+  --color-sidebar-accent: var(--sidebar-accent);
+  --color-sidebar-primary-foreground: var(--sidebar-primary-foreground);
+  --color-sidebar-primary: var(--sidebar-primary);
+  --color-sidebar-foreground: var(--sidebar-foreground);
+  --color-sidebar: var(--sidebar);
+  --color-chart-5: var(--chart-5);
+  --color-chart-4: var(--chart-4);
+  --color-chart-3: var(--chart-3);
+  --color-chart-2: var(--chart-2);
+  --color-chart-1: var(--chart-1);
+  --color-ring: var(--ring);
+  --color-input: var(--input);
+  --color-border: var(--border);
+  --color-destructive: var(--destructive);
+  --color-accent-foreground: var(--accent-foreground);
+  --color-accent: var(--accent);
+  --color-muted-foreground: var(--muted-foreground);
+  --color-muted: var(--muted);
+  --color-secondary-foreground: var(--secondary-foreground);
+  --color-secondary: var(--secondary);
+  --color-primary-foreground: var(--primary-foreground);
+  --color-primary: var(--primary);
+  --color-popover-foreground: var(--popover-foreground);
+  --color-popover: var(--popover);
+  --color-card-foreground: var(--card-foreground);
+  --color-card: var(--card);
+  --radius-sm: calc(var(--radius) * 0.6);
+  --radius-md: calc(var(--radius) * 0.8);
+  --radius-lg: var(--radius);
+  --radius-xl: calc(var(--radius) * 1.4);
+  --radius-2xl: calc(var(--radius) * 1.8);
+  --radius-3xl: calc(var(--radius) * 2.2);
+  --radius-4xl: calc(var(--radius) * 2.6);
 }
 
-@media (prefers-color-scheme: dark) {
-  :root {
-    --background: #0a0a0a;
-    --foreground: #ededed;
-  }
+:root {
+  --background: oklch(1 0 0);
+  --foreground: oklch(0.145 0 0);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.145 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.145 0 0);
+  --primary: oklch(0.205 0 0);
+  --primary-foreground: oklch(0.985 0 0);
+  --secondary: oklch(0.97 0 0);
+  --secondary-foreground: oklch(0.205 0 0);
+  --muted: oklch(0.97 0 0);
+  --muted-foreground: oklch(0.556 0 0);
+  --accent: oklch(0.97 0 0);
+  --accent-foreground: oklch(0.205 0 0);
+  --destructive: oklch(0.577 0.245 27.325);
+  --border: oklch(0.922 0 0);
+  --input: oklch(0.922 0 0);
+  --ring: oklch(0.708 0 0);
+  --chart-1: oklch(0.87 0 0);
+  --chart-2: oklch(0.556 0 0);
+  --chart-3: oklch(0.439 0 0);
+  --chart-4: oklch(0.371 0 0);
+  --chart-5: oklch(0.269 0 0);
+  --radius: 0.625rem;
+  --sidebar: oklch(0.985 0 0);
+  --sidebar-foreground: oklch(0.145 0 0);
+  --sidebar-primary: oklch(0.205 0 0);
+  --sidebar-primary-foreground: oklch(0.985 0 0);
+  --sidebar-accent: oklch(0.97 0 0);
+  --sidebar-accent-foreground: oklch(0.205 0 0);
+  --sidebar-border: oklch(0.922 0 0);
+  --sidebar-ring: oklch(0.708 0 0);
 }
 
-body {
-  background: var(--background);
-  color: var(--foreground);
-  font-family: Arial, Helvetica, sans-serif;
+.dark {
+  --background: oklch(0.145 0 0);
+  --foreground: oklch(0.985 0 0);
+  --card: oklch(0.205 0 0);
+  --card-foreground: oklch(0.985 0 0);
+  --popover: oklch(0.205 0 0);
+  --popover-foreground: oklch(0.985 0 0);
+  --primary: oklch(0.922 0 0);
+  --primary-foreground: oklch(0.205 0 0);
+  --secondary: oklch(0.269 0 0);
+  --secondary-foreground: oklch(0.985 0 0);
+  --muted: oklch(0.269 0 0);
+  --muted-foreground: oklch(0.708 0 0);
+  --accent: oklch(0.269 0 0);
+  --accent-foreground: oklch(0.985 0 0);
+  --destructive: oklch(0.704 0.191 22.216);
+  --border: oklch(1 0 0 / 10%);
+  --input: oklch(1 0 0 / 15%);
+  --ring: oklch(0.556 0 0);
+  --chart-1: oklch(0.87 0 0);
+  --chart-2: oklch(0.556 0 0);
+  --chart-3: oklch(0.439 0 0);
+  --chart-4: oklch(0.371 0 0);
+  --chart-5: oklch(0.269 0 0);
+  --sidebar: oklch(0.205 0 0);
+  --sidebar-foreground: oklch(0.985 0 0);
+  --sidebar-primary: oklch(0.488 0.243 264.376);
+  --sidebar-primary-foreground: oklch(0.985 0 0);
+  --sidebar-accent: oklch(0.269 0 0);
+  --sidebar-accent-foreground: oklch(0.985 0 0);
+  --sidebar-border: oklch(1 0 0 / 10%);
+  --sidebar-ring: oklch(0.556 0 0);
 }
+
+@layer base {
+  * {
+    @apply border-border outline-ring/50;
+  }
+  body {
+    @apply bg-background text-foreground;
+  }
+  html {
+    @apply font-sans;
+  }
+}
\ No newline at end of file
diff --git a/web/src/app/layout.tsx b/web/src/app/layout.tsx
index 9852c15..b22cb99 100644
--- a/web/src/app/layout.tsx
+++ b/web/src/app/layout.tsx
@@ -1,29 +1,33 @@
 import type { Metadata } from "next";
 import { Geist, Geist_Mono } from "next/font/google";
 import "./globals.css";
+import { ThemeProvider } from "@/components/theme-provider";
 
 const geistSans = Geist({
   variable: "--font-geist-sans",
   subsets: ["latin"],
 });
 
 const geistMono = Geist_Mono({
   variable: "--font-geist-mono",
   subsets: ["latin"],
 });
 
 export const metadata: Metadata = {
   title: "Create Next App",
   description: "Generated by create next app",
 };
 
 export default function RootLayout({ children }: LayoutProps<"/">) {
   return (
     <html
       lang="en"
+      suppressHydrationWarning
       className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
     >
-      <body className="min-h-full flex flex-col">{children}</body>
+      <body className="min-h-full flex flex-col">
+        <ThemeProvider>{children}</ThemeProvider>
+      </body>
     </html>
   );
 }
diff --git a/web/src/components/providers/theme-provider.tsx b/web/src/components/providers/theme-provider.tsx
new file mode 100644
index 0000000..d9eff81
--- /dev/null
+++ b/web/src/components/providers/theme-provider.tsx
@@ -0,0 +1,19 @@
+"use client";
+
+import { ThemeProvider as NextThemesProvider } from "next-themes";
+import { ReactNode } from "react";
+import { allThemeValues, DEFAULT_THEME } from "@/lib/themes-config";
+
+export function ThemeProvider({ children }: { children: ReactNode }) {
+  return (
+    <NextThemesProvider
+      attribute="data-theme"
+      themes={allThemeValues}
+      defaultTheme={DEFAULT_THEME}
+      enableSystem={false}
+      disableTransitionOnChange
+    >
+      {children}
+    </NextThemesProvider>
+  );
+}
diff --git a/web/src/components/theme-provider.tsx b/web/src/components/theme-provider.tsx
new file mode 100644
index 0000000..2c2809e
--- /dev/null
+++ b/web/src/components/theme-provider.tsx
@@ -0,0 +1 @@
+export { ThemeProvider } from "@/components/providers/theme-provider";
diff --git a/web/src/components/theme-switcher.tsx b/web/src/components/theme-switcher.tsx
new file mode 100644
index 0000000..abe64b5
--- /dev/null
+++ b/web/src/components/theme-switcher.tsx
@@ -0,0 +1,158 @@
+"use client";
+
+import { Check, Moon, Palette, Sun } from "lucide-react";
+import {
+  DropdownMenu,
+  DropdownMenuContent,
+  DropdownMenuItem,
+  DropdownMenuLabel,
+  DropdownMenuPortal,
+  DropdownMenuSeparator,
+  DropdownMenuSub,
+  DropdownMenuSubContent,
+  DropdownMenuSubTrigger,
+  DropdownMenuTrigger,
+} from "@/components/ui/dropdown-menu";
+import { sortedThemes, themes } from "@/lib/themes-config";
+import { useSyncExternalStore } from "react";
+
+import { Button } from "@/components/ui/button";
+import { ScrollArea } from "@/components/ui/scroll-area";
+import { cn } from "@/lib/utils";
+import { useTheme } from "next-themes";
+
+// Helper to get current color theme and mode from theme string like "catppuccin-dark"
+function parseTheme(theme: string | undefined): {
+  colorTheme: string;
+  mode: "light" | "dark";
+} {
+  if (!theme) return { colorTheme: "default", mode: "dark" };
+
+  if (theme.endsWith("-dark")) {
+    return { colorTheme: theme.replace("-dark", ""), mode: "dark" };
+  }
+  if (theme.endsWith("-light")) {
+    return { colorTheme: theme.replace("-light", ""), mode: "light" };
+  }
+  return { colorTheme: "default", mode: "dark" };
+}
+
+export function ThemeSwitcher() {
+  const mounted = useSyncExternalStore(
+    () => () => {},
+    () => true,
+    () => false
+  );
+  const { theme, setTheme } = useTheme();
+
+  const { colorTheme, mode } = parseTheme(theme);
+
+  const setColorTheme = (name: string) => {
+    setTheme(`${name}-${mode}`);
+  };
+
+  const setMode = (newMode: "light" | "dark") => {
+    setTheme(`${colorTheme}-${newMode}`);
+  };
+
+  if (!mounted) {
+    return (
+      <Button variant="ghost" size="icon" className="h-9 w-9">
+        <Sun className="h-4 w-4" />
+      </Button>
+    );
+  }
+
+  const currentThemeConfig =
+    themes.find((t) => t.name === colorTheme) || themes[0];
+  const isDark = mode === "dark";
+
+  return (
+    <DropdownMenu>
+      <DropdownMenuTrigger render={<Button variant="ghost" size="icon" className="h-9 w-9" />}>{isDark ? <Moon className="h-4 w-4" /> : <Sun className="h-4 w-4" />}<span className="sr-only">Toggle theme</span></DropdownMenuTrigger>
+      <DropdownMenuContent align="end" className="w-48">
+        <DropdownMenuLabel className="flex items-center gap-2">
+          <div
+            className="h-3 w-3 rounded-full"
+            style={{
+              backgroundColor: isDark
+                ? currentThemeConfig.primaryDark
+                : currentThemeConfig.primaryLight,
+            }}
+          />
+          {currentThemeConfig.title}
+        </DropdownMenuLabel>
+        <DropdownMenuSeparator />
+
+        {/* Mode Toggle */}
+        <DropdownMenuItem
+          onClick={() => setMode("light")}
+          className={cn("flex items-center gap-2 cursor-pointer")}
+        >
+          <Sun className="h-4 w-4" />
+          <span className="flex-1">Light</span>
+          {mode === "light" && <Check className="h-4 w-4 text-primary" />}
+        </DropdownMenuItem>
+        <DropdownMenuItem
+          onClick={() => setMode("dark")}
+          className={cn("flex items-center gap-2 cursor-pointer")}
+        >
+          <Moon className="h-4 w-4" />
+          <span className="flex-1">Dark</span>
+          {mode === "dark" && <Check className="h-4 w-4 text-primary" />}
+        </DropdownMenuItem>
+
+        <DropdownMenuSeparator />
+
+        {/* Color Theme Submenu */}
+        <DropdownMenuSub>
+          <DropdownMenuSubTrigger className="flex items-center gap-2">
+            <Palette className="h-4 w-4" />
+            <span>Color Theme</span>
+          </DropdownMenuSubTrigger>
+          <DropdownMenuPortal>
+            <DropdownMenuSubContent className="w-52">
+              <ScrollArea className="h-80">
+                <div className="p-1">
+                  {sortedThemes.map((t) => (
+                    <DropdownMenuItem
+                      key={t.name}
+                      onClick={() => setColorTheme(t.name)}
+                      className="flex items-center gap-3 cursor-pointer"
+                    >
+                      <div
+                        className="h-4 w-4 rounded-full border border-border shrink-0"
+                        style={{
+                          backgroundColor: isDark
+                            ? t.primaryDark
+                            : t.primaryLight,
+                        }}
+                      />
+                      <span className="flex-1">{t.title}</span>
+                      {colorTheme === t.name && (
+                        <Check className="h-4 w-4 text-primary" />
+                      )}
+                    </DropdownMenuItem>
+                  ))}
+                </div>
+              </ScrollArea>
+              <DropdownMenuSeparator />
+              <div className="px-2 py-1.5 text-xs text-muted-foreground text-center">
+                Themes by{" "}
+                <a
+                  href="https://tweakcn.com"
+                  target="_blank"
+                  rel="noopener noreferrer"
+                  className="font-medium text-foreground hover:underline"
+                  onClick={(e) => e.stopPropagation()}
+                >
+                  tweakcn
+                </a>
+              </div>
+            </DropdownMenuSubContent>
+          </DropdownMenuPortal>
+        </DropdownMenuSub>
+      </DropdownMenuContent>
+    </DropdownMenu>
+  );
+}
diff --git a/web/src/components/theme-toggle.tsx b/web/src/components/theme-toggle.tsx
new file mode 100644
index 0000000..51c36f8
--- /dev/null
+++ b/web/src/components/theme-toggle.tsx
@@ -0,0 +1,31 @@
+"use client";
+
+import { Moon, Sun } from "lucide-react";
+import { useTheme } from "next-themes";
+
+import { Button } from "@/components/ui/button";
+
+// Tweakcn themes use the "{name}-dark" / "{name}-light" scheme, so next-themes
+// resolves to values like "stripe-dark". Detect dark via the "-dark" suffix and
+// toggle between the current color theme's light / dark variants.
+export function ThemeToggle() {
+  const { resolvedTheme, setTheme } = useTheme();
+
+  const isDark = resolvedTheme?.endsWith("-dark");
+  const colorTheme = resolvedTheme?.replace(/-(light|dark)$/, "") || "stripe";
+
+  const toggle = () =>
+    setTheme(`${colorTheme}${isDark ? "-light" : "-dark"}`);
+
+  return (
+    <Button
+      variant="outline"
+      size="icon"
+      aria-label="Toggle theme"
+      onClick={toggle}
+    >
+      <Sun className="h-4 w-4 dark:hidden" />
+      <Moon className="hidden h-4 w-4 dark:block" />
+    </Button>
+  );
+}
diff --git a/web/src/components/ui/button.tsx b/web/src/components/ui/button.tsx
new file mode 100644
index 0000000..b033601
--- /dev/null
+++ b/web/src/components/ui/button.tsx
@@ -0,0 +1,58 @@
+import { Button as ButtonPrimitive } from "@base-ui/react/button"
+import { cva, type VariantProps } from "class-variance-authority"
+
+import { cn } from "@/lib/utils"
+
+const buttonVariants = cva(
+  "group/button inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
+  {
+    variants: {
+      variant: {
+        default: "bg-primary text-primary-foreground hover:bg-primary/80",
+        outline:
+          "border-border bg-background hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50",
+        secondary:
+          "bg-secondary text-secondary-foreground hover:bg-[color-mix(in_oklch,var(--secondary),var(--foreground)_5%)] aria-expanded:bg-secondary aria-expanded:text-secondary-foreground",
+        ghost:
+          "hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:hover:bg-muted/50",
+        destructive:
+          "bg-destructive/10 text-destructive hover:bg-destructive/20 focus-visible:border-destructive/40 focus-visible:ring-destructive/20 dark:bg-destructive/20 dark:hover:bg-destructive/30 dark:focus-visible:ring-destructive/40",
+        link: "text-primary underline-offset-4 hover:underline",
+      },
+      size: {
+        default:
+          "h-8 gap-1.5 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2",
+        xs: "h-6 gap-1 rounded-[min(var(--radius-md),10px)] px-2 text-xs in-data-[slot=button-group]:rounded-lg has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3",
+        sm: "h-7 gap-1 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] in-data-[slot=button-group]:rounded-lg has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3.5",
+        lg: "h-9 gap-1.5 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2",
+        icon: "size-8",
+        "icon-xs":
+          "size-6 rounded-[min(var(--radius-md),10px)] in-data-[slot=button-group]:rounded-lg [&_svg:not([class*='size-'])]:size-3",
+        "icon-sm":
+          "size-7 rounded-[min(var(--radius-md),12px)] in-data-[slot=button-group]:rounded-lg",
+        "icon-lg": "size-9",
+      },
+    },
+    defaultVariants: {
+      variant: "default",
+      size: "default",
+    },
+  }
+)
+
+function Button({
+  className,
+  variant = "default",
+  size = "default",
+  ...props
+}: ButtonPrimitive.Props & VariantProps<typeof buttonVariants>) {
+  return (
+    <ButtonPrimitive
+      data-slot="button"
+      className={cn(buttonVariants({ variant, size, className }))}
+      {...props}
+    />
+  )
+}
+
+export { Button, buttonVariants }
diff --git a/web/src/components/ui/dropdown-menu.tsx b/web/src/components/ui/dropdown-menu.tsx
new file mode 100644
index 0000000..9d5ebbd
--- /dev/null
+++ b/web/src/components/ui/dropdown-menu.tsx
@@ -0,0 +1,268 @@
+"use client"
+
+import * as React from "react"
+import { Menu as MenuPrimitive } from "@base-ui/react/menu"
+
+import { cn } from "@/lib/utils"
+import { ChevronRightIcon, CheckIcon } from "lucide-react"
+
+function DropdownMenu({ ...props }: MenuPrimitive.Root.Props) {
+  return <MenuPrimitive.Root data-slot="dropdown-menu" {...props} />
+}
+
+function DropdownMenuPortal({ ...props }: MenuPrimitive.Portal.Props) {
+  return <MenuPrimitive.Portal data-slot="dropdown-menu-portal" {...props} />
+}
+
+function DropdownMenuTrigger({ ...props }: MenuPrimitive.Trigger.Props) {
+  return <MenuPrimitive.Trigger data-slot="dropdown-menu-trigger" {...props} />
+}
+
+function DropdownMenuContent({
+  align = "start",
+  alignOffset = 0,
+  side = "bottom",
+  sideOffset = 4,
+  className,
+  ...props
+}: MenuPrimitive.Popup.Props &
+  Pick<
+    MenuPrimitive.Positioner.Props,
+    "align" | "alignOffset" | "side" | "sideOffset"
+  >) {
+  return (
+    <MenuPrimitive.Portal>
+      <MenuPrimitive.Positioner
+        className="isolate z-50 outline-none"
+        align={align}
+        alignOffset={alignOffset}
+        side={side}
+        sideOffset={sideOffset}
+      >
+        <MenuPrimitive.Popup
+          data-slot="dropdown-menu-content"
+          className={cn("z-50 max-h-(--available-height) w-(--anchor-width) min-w-32 origin-(--transform-origin) overflow-x-hidden overflow-y-auto rounded-lg bg-popover p-1 text-popover-foreground shadow-md ring-1 ring-foreground/10 duration-100 outline-none data-[side=bottom]:slide-in-from-top-2 data-[side=inline-end]:slide-in-from-left-2 data-[side=inline-start]:slide-in-from-right-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:overflow-hidden data-closed:fade-out-0 data-closed:zoom-out-95", className )}
+          {...props}
+        />
+      </MenuPrimitive.Positioner>
+    </MenuPrimitive.Portal>
+  )
+}
+
+function DropdownMenuGroup({ ...props }: MenuPrimitive.Group.Props) {
+  return <MenuPrimitive.Group data-slot="dropdown-menu-group" {...props} />
+}
+
+function DropdownMenuLabel({
+  className,
+  inset,
+  ...props
+}: MenuPrimitive.GroupLabel.Props & {
+  inset?: boolean
+}) {
+  return (
+    <MenuPrimitive.GroupLabel
+      data-slot="dropdown-menu-label"
+      data-inset={inset}
+      className={cn(
+        "px-1.5 py-1 text-xs font-medium text-muted-foreground data-inset:pl-7",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+function DropdownMenuItem({
+  className,
+  inset,
+  variant = "default",
+  ...props
+}: MenuPrimitive.Item.Props & {
+  inset?: boolean
+  variant?: "default" | "destructive"
+}) {
+  return (
+    <MenuPrimitive.Item
+      data-slot="dropdown-menu-item"
+      data-inset={inset}
+      data-variant={variant}
+      className={cn(
+        "group/dropdown-menu-item relative flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground not-data-[variant=destructive]:focus:**:text-accent-foreground data-inset:pl-7 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 data-[variant=destructive]:*:[svg]:text-destructive",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+function DropdownMenuSub({ ...props }: MenuPrimitive.SubmenuRoot.Props) {
+  return <MenuPrimitive.SubmenuRoot data-slot="dropdown-menu-sub" {...props} />
+}
+
+function DropdownMenuSubTrigger({
+  className,
+  inset,
+  children,
+  ...props
+}: MenuPrimitive.SubmenuTrigger.Props & {
+  inset?: boolean
+}) {
+  return (
+    <MenuPrimitive.SubmenuTrigger
+      data-slot="dropdown-menu-sub-trigger"
+      data-inset={inset}
+      className={cn(
+        "flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground not-data-[variant=destructive]:focus:**:text-accent-foreground data-inset:pl-7 data-popup-open:bg-accent data-popup-open:text-accent-foreground data-open:bg-accent data-open:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
+        className
+      )}
+      {...props}
+    >
+      {children}
+      <ChevronRightIcon className="ml-auto" />
+    </MenuPrimitive.SubmenuTrigger>
+  )
+}
+
+function DropdownMenuSubContent({
+  align = "start",
+  alignOffset = -3,
+  side = "right",
+  sideOffset = 0,
+  className,
+  ...props
+}: React.ComponentProps<typeof DropdownMenuContent>) {
+  return (
+    <DropdownMenuContent
+      data-slot="dropdown-menu-sub-content"
+      className={cn("w-auto min-w-[96px] rounded-lg bg-popover p-1 text-popover-foreground shadow-lg ring-1 ring-foreground/10 duration-100 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95", className )}
+      align={align}
+      alignOffset={alignOffset}
+      side={side}
+      sideOffset={sideOffset}
+      {...props}
+    />
+  )
+}
+
+function DropdownMenuCheckboxItem({
+  className,
+  children,
+  checked,
+  inset,
+  ...props
+}: MenuPrimitive.CheckboxItem.Props & {
+  inset?: boolean
+}) {
+  return (
+    <MenuPrimitive.CheckboxItem
+      data-slot="dropdown-menu-checkbox-item"
+      data-inset={inset}
+      className={cn(
+        "relative flex cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground focus:**:text-accent-foreground data-inset:pl-7 data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
+        className
+      )}
+      checked={checked}
+      {...props}
+    >
+      <span
+        className="pointer-events-none absolute right-2 flex items-center justify-center"
+        data-slot="dropdown-menu-checkbox-item-indicator"
+      >
+        <MenuPrimitive.CheckboxItemIndicator>
+          <CheckIcon
+          />
+        </MenuPrimitive.CheckboxItemIndicator>
+      </span>
+      {children}
+    </MenuPrimitive.CheckboxItem>
+  )
+}
+
+function DropdownMenuRadioGroup({ ...props }: MenuPrimitive.RadioGroup.Props) {
+  return (
+    <MenuPrimitive.RadioGroup
+      data-slot="dropdown-menu-radio-group"
+      {...props}
+    />
+  )
+}
+
+function DropdownMenuRadioItem({
+  className,
+  children,
+  inset,
+  ...props
+}: MenuPrimitive.RadioItem.Props & {
+  inset?: boolean
+}) {
+  return (
+    <MenuPrimitive.RadioItem
+      data-slot="dropdown-menu-radio-item"
+      data-inset={inset}
+      className={cn(
+        "relative flex cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground focus:**:text-accent-foreground data-inset:pl-7 data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
+        className
+      )}
+      {...props}
+    >
+      <span
+        className="pointer-events-none absolute right-2 flex items-center justify-center"
+        data-slot="dropdown-menu-radio-item-indicator"
+      >
+        <MenuPrimitive.RadioItemIndicator>
+          <CheckIcon
+          />
+        </MenuPrimitive.RadioItemIndicator>
+      </span>
+      {children}
+    </MenuPrimitive.RadioItem>
+  )
+}
+
+function DropdownMenuSeparator({
+  className,
+  ...props
+}: MenuPrimitive.Separator.Props) {
+  return (
+    <MenuPrimitive.Separator
+      data-slot="dropdown-menu-separator"
+      className={cn("-mx-1 my-1 h-px bg-border", className)}
+      {...props}
+    />
+  )
+}
+
+function DropdownMenuShortcut({
+  className,
+  ...props
+}: React.ComponentProps<"span">) {
+  return (
+    <span
+      data-slot="dropdown-menu-shortcut"
+      className={cn(
+        "ml-auto text-xs tracking-widest text-muted-foreground group-focus/dropdown-menu-item:text-accent-foreground",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+export {
+  DropdownMenu,
+  DropdownMenuPortal,
+  DropdownMenuTrigger,
+  DropdownMenuContent,
+  DropdownMenuGroup,
+  DropdownMenuLabel,
+  DropdownMenuItem,
+  DropdownMenuCheckboxItem,
+  DropdownMenuRadioGroup,
+  DropdownMenuRadioItem,
+  DropdownMenuSeparator,
+  DropdownMenuShortcut,
+  DropdownMenuSub,
+  DropdownMenuSubTrigger,
+  DropdownMenuSubContent,
+}
diff --git a/web/src/components/ui/scroll-area.tsx b/web/src/components/ui/scroll-area.tsx
new file mode 100644
index 0000000..84c1e9f
--- /dev/null
+++ b/web/src/components/ui/scroll-area.tsx
@@ -0,0 +1,55 @@
+"use client"
+
+import * as React from "react"
+import { ScrollArea as ScrollAreaPrimitive } from "@base-ui/react/scroll-area"
+
+import { cn } from "@/lib/utils"
+
+function ScrollArea({
+  className,
+  children,
+  ...props
+}: ScrollAreaPrimitive.Root.Props) {
+  return (
+    <ScrollAreaPrimitive.Root
+      data-slot="scroll-area"
+      className={cn("relative", className)}
+      {...props}
+    >
+      <ScrollAreaPrimitive.Viewport
+        data-slot="scroll-area-viewport"
+        className="size-full rounded-[inherit] transition-[color,box-shadow] outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1"
+      >
+        {children}
+      </ScrollAreaPrimitive.Viewport>
+      <ScrollBar />
+      <ScrollAreaPrimitive.Corner />
+    </ScrollAreaPrimitive.Root>
+  )
+}
+
+function ScrollBar({
+  className,
+  orientation = "vertical",
+  ...props
+}: ScrollAreaPrimitive.Scrollbar.Props) {
+  return (
+    <ScrollAreaPrimitive.Scrollbar
+      data-slot="scroll-area-scrollbar"
+      data-orientation={orientation}
+      orientation={orientation}
+      className={cn(
+        "flex touch-none p-px transition-colors select-none data-horizontal:h-2.5 data-horizontal:flex-col data-horizontal:border-t data-horizontal:border-t-transparent data-vertical:h-full data-vertical:w-2.5 data-vertical:border-l data-vertical:border-l-transparent",
+        className
+      )}
+      {...props}
+    >
+      <ScrollAreaPrimitive.Thumb
+        data-slot="scroll-area-thumb"
+        className="relative flex-1 rounded-full bg-border"
+      />
+    </ScrollAreaPrimitive.Scrollbar>
+  )
+}
+
+export { ScrollArea, ScrollBar }
diff --git a/web/src/lib/themes-config.ts b/web/src/lib/themes-config.ts
new file mode 100644
index 0000000..4068e8a
--- /dev/null
+++ b/web/src/lib/themes-config.ts
@@ -0,0 +1,331 @@
+export interface ThemeConfig {
+  name: string;
+  title: string;
+  primaryLight: string;
+  primaryDark: string;
+  fontSans: string;
+}
+
+export const themes: ThemeConfig[] = [
+  {
+    name: "default",
+    title: "Default",
+    primaryLight: "oklch(0.2050 0 0)",
+    primaryDark: "oklch(0.9220 0 0)",
+    fontSans: "ui-sans-serif, system-ui, sans-serif",
+  },
+  {
+    name: "amber-minimal",
+    title: "Amber Minimal",
+    primaryLight: "oklch(0.77 0.16 70.08)",
+    primaryDark: "oklch(0.77 0.16 70.08)",
+    fontSans: "Inter, sans-serif",
+  },
+  {
+    name: "bold-tech",
+    title: "Bold Tech",
+    primaryLight: "oklch(0.61 0.22 292.72)",
+    primaryDark: "oklch(0.61 0.22 292.72)",
+    fontSans: "Roboto, sans-serif",
+  },
+  {
+    name: "bubblegum",
+    title: "Bubblegum",
+    primaryLight: "oklch(0.62 0.18 348.14)",
+    primaryDark: "oklch(0.92 0.08 87.67)",
+    fontSans: "Poppins, sans-serif",
+  },
+  {
+    name: "caffeine",
+    title: "Caffeine",
+    primaryLight: "oklch(0.43 0.04 41.99)",
+    primaryDark: "oklch(0.92 0.05 66.17)",
+    fontSans: "system-ui, sans-serif",
+  },
+  {
+    name: "candyland",
+    title: "Candyland",
+    primaryLight: "oklch(0.87 0.07 7.09)",
+    primaryDark: "oklch(0.8 0.14 349.23)",
+    fontSans: "Poppins, sans-serif",
+  },
+  {
+    name: "catppuccin",
+    title: "Catppuccin",
+    primaryLight: "oklch(0.55 0.25 297.02)",
+    primaryDark: "oklch(0.79 0.12 304.77)",
+    fontSans: "Montserrat, sans-serif",
+  },
+  {
+    name: "claude",
+    title: "Claude",
+    primaryLight: "oklch(0.62 0.14 39.04)",
+    primaryDark: "oklch(0.67 0.13 38.76)",
+    fontSans: "system-ui, sans-serif",
+  },
+  {
+    name: "claymorphism",
+    title: "Claymorphism",
+    primaryLight: "oklch(0.59 0.2 277.12)",
+    primaryDark: "oklch(0.68 0.16 276.93)",
+    fontSans: "Plus Jakarta Sans, sans-serif",
+  },
+  {
+    name: "clean-slate",
+    title: "Clean Slate",
+    primaryLight: "oklch(0.59 0.2 277.12)",
+    primaryDark: "oklch(0.68 0.16 276.93)",
+    fontSans: "Inter, sans-serif",
+  },
+  {
+    name: "cosmic-night",
+    title: "Cosmic Night",
+    primaryLight: "oklch(0.54 0.18 288.03)",
+    primaryDark: "oklch(0.72 0.16 290.40)",
+    fontSans: "Inter, sans-serif",
+  },
+  {
+    name: "cyberpunk",
+    title: "Cyberpunk",
+    primaryLight: "oklch(0.67 0.29 341.41)",
+    primaryDark: "oklch(0.67 0.29 341.41)",
+    fontSans: "Outfit, sans-serif",
+  },
+  {
+    name: "doom-64",
+    title: "Doom 64",
+    primaryLight: "oklch(0.5 0.19 27.48)",
+    primaryDark: "oklch(0.61 0.21 27.03)",
+    fontSans: "Oxanium, sans-serif",
+  },
+  {
+    name: "elegant-luxury",
+    title: "Elegant Luxury",
+    primaryLight: "oklch(0.47 0.15 24.94)",
+    primaryDark: "oklch(0.51 0.19 27.52)",
+    fontSans: "Poppins, sans-serif",
+  },
+  {
+    name: "graphite",
+    title: "Graphite",
+    primaryLight: "oklch(0.49 0 0)",
+    primaryDark: "oklch(0.71 0 0)",
+    fontSans: "Inter, sans-serif",
+  },
+  {
+    name: "kodama-grove",
+    title: "Kodama Grove",
+    primaryLight: "oklch(0.67 0.11 118.91)",
+    primaryDark: "oklch(0.68 0.06 132.45)",
+    fontSans: "Merriweather, serif",
+  },
+  {
+    name: "midnight-bloom",
+    title: "Midnight Bloom",
+    primaryLight: "oklch(0.57 0.20 283.08)",
+    primaryDark: "oklch(0.57 0.20 283.08)",
+    fontSans: "Montserrat, sans-serif",
+  },
+  {
+    name: "mocha-mousse",
+    title: "Mocha Mousse",
+    primaryLight: "oklch(0.61 0.06 44.36)",
+    primaryDark: "oklch(0.73 0.05 52.33)",
+    fontSans: "DM Sans, sans-serif",
+  },
+  {
+    name: "modern-minimal",
+    title: "Modern Minimal",
+    primaryLight: "oklch(0.62 0.19 259.81)",
+    primaryDark: "oklch(0.62 0.19 259.81)",
+    fontSans: "Inter, sans-serif",
+  },
+  {
+    name: "mono",
+    title: "Mono",
+    primaryLight: "oklch(0.56 0 0)",
+    primaryDark: "oklch(0.56 0 0)",
+    fontSans: "Geist Mono, monospace",
+  },
+  {
+    name: "nature",
+    title: "Nature",
+    primaryLight: "oklch(0.52 0.13 144.17)",
+    primaryDark: "oklch(0.67 0.16 144.21)",
+    fontSans: "Montserrat, sans-serif",
+  },
+  {
+    name: "neo-brutalism",
+    title: "Neo Brutalism",
+    primaryLight: "oklch(0.65 0.24 26.97)",
+    primaryDark: "oklch(0.70 0.19 23.19)",
+    fontSans: "DM Sans, sans-serif",
+  },
+  {
+    name: "northern-lights",
+    title: "Northern Lights",
+    primaryLight: "oklch(0.65 0.15 150.31)",
+    primaryDark: "oklch(0.65 0.15 150.31)",
+    fontSans: "Plus Jakarta Sans, sans-serif",
+  },
+  {
+    name: "ocean-breeze",
+    title: "Ocean Breeze",
+    primaryLight: "oklch(0.72 0.19 149.58)",
+    primaryDark: "oklch(0.77 0.15 163.22)",
+    fontSans: "DM Sans, sans-serif",
+  },
+  {
+    name: "pastel-dreams",
+    title: "Pastel Dreams",
+    primaryLight: "oklch(0.71 0.16 293.54)",
+    primaryDark: "oklch(0.79 0.12 295.75)",
+    fontSans: "Open Sans, sans-serif",
+  },
+  {
+    name: "perpetuity",
+    title: "Perpetuity",
+    primaryLight: "oklch(0.56 0.09 203.28)",
+    primaryDark: "oklch(0.85 0.13 195.04)",
+    fontSans: "Courier New, monospace",
+  },
+  {
+    name: "quantum-rose",
+    title: "Quantum Rose",
+    primaryLight: "oklch(0.6 0.24 0.13)",
+    primaryDark: "oklch(0.75 0.23 332.02)",
+    fontSans: "Poppins, sans-serif",
+  },
+  {
+    name: "retro-arcade",
+    title: "Retro Arcade",
+    primaryLight: "oklch(0.59 0.2 355.89)",
+    primaryDark: "oklch(0.59 0.2 355.89)",
+    fontSans: "Outfit, sans-serif",
+  },
+  {
+    name: "solar-dusk",
+    title: "Solar Dusk",
+    primaryLight: "oklch(0.56 0.15 49)",
+    primaryDark: "oklch(0.7 0.19 47.6)",
+    fontSans: "Oxanium, sans-serif",
+  },
+  {
+    name: "starry-night",
+    title: "Starry Night",
+    primaryLight: "oklch(0.48 0.12 263.38)",
+    primaryDark: "oklch(0.48 0.12 263.38)",
+    fontSans: "Libre Baskerville, serif",
+  },
+  {
+    name: "supabase",
+    title: "Supabase",
+    primaryLight: "oklch(0.83 0.13 160.91)",
+    primaryDark: "oklch(0.44 0.1 156.76)",
+    fontSans: "Outfit, sans-serif",
+  },
+  {
+    name: "sunset-horizon",
+    title: "Sunset Horizon",
+    primaryLight: "oklch(0.74 0.16 34.71)",
+    primaryDark: "oklch(0.74 0.16 34.71)",
+    fontSans: "Montserrat, sans-serif",
+  },
+  {
+    name: "t3-chat",
+    title: "T3 Chat",
+    primaryLight: "oklch(0.53 0.14 355.2)",
+    primaryDark: "oklch(0.46 0.19 4.1)",
+    fontSans: "system-ui, sans-serif",
+  },
+  {
+    name: "tangerine",
+    title: "Tangerine",
+    primaryLight: "oklch(0.64 0.17 36.44)",
+    primaryDark: "oklch(0.64 0.17 36.44)",
+    fontSans: "Inter, sans-serif",
+  },
+  {
+    name: "twitter",
+    title: "Twitter",
+    primaryLight: "oklch(0.67 0.16 245)",
+    primaryDark: "oklch(0.67 0.16 245.01)",
+    fontSans: "Open Sans, sans-serif",
+  },
+  {
+    name: "vercel",
+    title: "Vercel",
+    primaryLight: "oklch(0 0 0)",
+    primaryDark: "oklch(1 0 0)",
+    fontSans: "Geist, sans-serif",
+  },
+  {
+    name: "vintage-paper",
+    title: "Vintage Paper",
+    primaryLight: "oklch(0.62 0.08 65.54)",
+    primaryDark: "oklch(0.73 0.06 66.7)",
+    fontSans: "Libre Baskerville, serif",
+  },
+  {
+    name: "twitch",
+    title: "Twitch",
+    primaryLight: "oklch(0.54 0.24 292)",
+    primaryDark: "oklch(0.60 0.22 292)",
+    fontSans: "Inter, sans-serif",
+  },
+  {
+    name: "kick",
+    title: "Kick",
+    primaryLight: "oklch(0.75 0.26 135)",
+    primaryDark: "oklch(0.83 0.28 135)",
+    fontSans: "Inter, sans-serif",
+  },
+  {
+    name: "spotify",
+    title: "Spotify",
+    primaryLight: "oklch(0.64 0.20 155)",
+    primaryDark: "oklch(0.68 0.20 155)",
+    fontSans: "Montserrat, sans-serif",
+  },
+  {
+    name: "stripe",
+    title: "Stripe",
+    primaryLight: "oklch(0.55 0.24 280)",
+    primaryDark: "oklch(0.65 0.22 280)",
+    fontSans: "system-ui, sans-serif",
+  },
+  {
+    name: "github",
+    title: "GitHub",
+    primaryLight: "oklch(0.52 0.16 145)",
+    primaryDark: "oklch(0.62 0.16 145)",
+    fontSans: "system-ui, sans-serif",
+  },
+  {
+    name: "windows98",
+    title: "Windows 98",
+    primaryLight: "oklch(0.2711 0.1879 264.0520)",
+    primaryDark: "oklch(0.2711 0.1879 264.0520)",
+    fontSans: "Pixelify Sans, sans-serif",
+  },
+];
+
+// Sort themes alphabetically by title, but keep default first
+export const sortedThemes = [
+  themes[0],
+  ...themes.slice(1).sort((a, b) => a.title.localeCompare(b.title)),
+];
+
+export const themeNames = themes.map((t) => t.name);
+
+// Generate all theme values for next-themes (name-light and name-dark variants)
+export const allThemeValues = themes.flatMap((t) => [
+  `${t.name}-light`,
+  `${t.name}-dark`,
+]);
+
+// Default theme
+export const DEFAULT_THEME = "stripe-dark";
+
+// Convenience alias for the default theme (used by layout)
+export const defaultTheme = "stripe-dark";
diff --git a/web/src/lib/utils.ts b/web/src/lib/utils.ts
new file mode 100644
index 0000000..bd0c391
--- /dev/null
+++ b/web/src/lib/utils.ts
@@ -0,0 +1,6 @@
+import { clsx, type ClassValue } from "clsx"
+import { twMerge } from "tailwind-merge"
+
+export function cn(...inputs: ClassValue[]) {
+  return twMerge(clsx(inputs))
+}
diff --git a/web/src/styles/themes/amber-minimal.css b/web/src/styles/themes/amber-minimal.css
new file mode 100644
index 0000000..6335939
--- /dev/null
+++ b/web/src/styles/themes/amber-minimal.css
@@ -0,0 +1,124 @@
+/* Amber Minimal Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="amber-minimal-light"] {
+  --background: oklch(1 0 0);
+  --foreground: oklch(0.27 0 0);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.27 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.27 0 0);
+  --primary: oklch(0.77 0.16 70.08);
+  --primary-foreground: oklch(0 0 0);
+  --secondary: oklch(0.97 0 264.54);
+  --secondary-foreground: oklch(0.45 0.03 256.8);
+  --muted: oklch(0.98 0 247.84);
+  --muted-foreground: oklch(0.55 0.02 264.36);
+  --accent: oklch(0.99 0.02 95.28);
+  --accent-foreground: oklch(0.47 0.12 46.2);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.93 0.01 264.53);
+  --input: oklch(0.93 0.01 264.53);
+  --ring: oklch(0.77 0.16 70.08);
+  --chart-1: oklch(0.77 0.16 70.08);
+  --chart-2: oklch(0.67 0.16 58.32);
+  --chart-3: oklch(0.56 0.15 49);
+  --chart-4: oklch(0.47 0.12 46.2);
+  --chart-5: oklch(0.41 0.11 45.9);
+  --radius: 0.375rem;
+  --sidebar: oklch(0.98 0 247.84);
+  --sidebar-foreground: oklch(0.27 0 0);
+  --sidebar-primary: oklch(0.77 0.16 70.08);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.99 0.02 95.28);
+  --sidebar-accent-foreground: oklch(0.47 0.12 46.2);
+  --sidebar-border: oklch(0.93 0.01 264.53);
+  --sidebar-ring: oklch(0.77 0.16 70.08);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 8px;
+  --shadow-spread: -1px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 2px 4px -2px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 4px 6px -2px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 8px 10px -2px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 4px 8px -1px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="amber-minimal-dark"] {
+  --background: oklch(0.2 0 0);
+  --foreground: oklch(0.92 0 0);
+  --card: oklch(0.27 0 0);
+  --card-foreground: oklch(0.92 0 0);
+  --popover: oklch(0.27 0 0);
+  --popover-foreground: oklch(0.92 0 0);
+  --primary: oklch(0.77 0.16 70.08);
+  --primary-foreground: oklch(0 0 0);
+  --secondary: oklch(0.27 0 0);
+  --secondary-foreground: oklch(0.92 0 0);
+  --muted: oklch(0.27 0 0);
+  --muted-foreground: oklch(0.72 0 0);
+  --accent: oklch(0.47 0.12 46.2);
+  --accent-foreground: oklch(0.92 0.12 95.75);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.37 0 0);
+  --input: oklch(0.37 0 0);
+  --ring: oklch(0.77 0.16 70.08);
+  --chart-1: oklch(0.84 0.16 84.43);
+  --chart-2: oklch(0.67 0.16 58.32);
+  --chart-3: oklch(0.47 0.12 46.2);
+  --chart-4: oklch(0.56 0.15 49);
+  --chart-5: oklch(0.47 0.12 46.2);
+  --radius: 0.375rem;
+  --sidebar: oklch(0.17 0 0);
+  --sidebar-foreground: oklch(0.92 0 0);
+  --sidebar-primary: oklch(0.77 0.16 70.08);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.47 0.12 46.2);
+  --sidebar-accent-foreground: oklch(0.92 0.12 95.75);
+  --sidebar-border: oklch(0.37 0 0);
+  --sidebar-ring: oklch(0.77 0.16 70.08);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 8px;
+  --shadow-spread: -1px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 2px 4px -2px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 4px 6px -2px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 8px 10px -2px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 4px 8px -1px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/bold-tech.css b/web/src/styles/themes/bold-tech.css
new file mode 100644
index 0000000..d687395
--- /dev/null
+++ b/web/src/styles/themes/bold-tech.css
@@ -0,0 +1,134 @@
+/* Bold Tech Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="bold-tech-light"] {
+  --background: oklch(1 0 0);
+  --foreground: oklch(0.36 0.14 278.7);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.36 0.14 278.7);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.36 0.14 278.7);
+  --primary: oklch(0.61 0.22 292.72);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.96 0.02 295.19);
+  --secondary-foreground: oklch(0.46 0.21 277.02);
+  --muted: oklch(0.97 0.02 293.76);
+  --muted-foreground: oklch(0.54 0.25 293.01);
+  --accent: oklch(0.93 0.03 255.59);
+  --accent-foreground: oklch(0.42 0.18 265.64);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.93 0.03 272.79);
+  --input: oklch(0.93 0.03 272.79);
+  --ring: oklch(0.61 0.22 292.72);
+  --chart-1: oklch(0.61 0.22 292.72);
+  --chart-2: oklch(0.54 0.25 293.01);
+  --chart-3: oklch(0.49 0.24 292.58);
+  --chart-4: oklch(0.43 0.21 292.76);
+  --chart-5: oklch(0.38 0.18 293.74);
+  --radius: 0.625rem;
+  --sidebar: oklch(0.97 0.02 293.76);
+  --sidebar-foreground: oklch(0.36 0.14 278.7);
+  --sidebar-primary: oklch(0.61 0.22 292.72);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.93 0.03 255.59);
+  --sidebar-accent-foreground: oklch(0.42 0.18 265.64);
+  --sidebar-border: oklch(0.93 0.03 272.79);
+  --sidebar-ring: oklch(0.61 0.22 292.72);
+  --font-sans: Roboto, sans-serif;
+  --font-serif: Playfair Display, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(255 86% 66%);
+  --shadow-opacity: 0.2;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 2px 4px 0px hsl(255 86% 66% / 0.1);
+  --shadow-xs: 2px 2px 4px 0px hsl(255 86% 66% / 0.1);
+  --shadow-sm:
+    2px 2px 4px 0px hsl(255 86% 66% / 0.2),
+    2px 1px 2px -1px hsl(255 86% 66% / 0.2);
+  --shadow:
+    2px 2px 4px 0px hsl(255 86% 66% / 0.2),
+    2px 1px 2px -1px hsl(255 86% 66% / 0.2);
+  --shadow-md:
+    2px 2px 4px 0px hsl(255 86% 66% / 0.2),
+    2px 2px 4px -1px hsl(255 86% 66% / 0.2);
+  --shadow-lg:
+    2px 2px 4px 0px hsl(255 86% 66% / 0.2),
+    2px 4px 6px -1px hsl(255 86% 66% / 0.2);
+  --shadow-xl:
+    2px 2px 4px 0px hsl(255 86% 66% / 0.2),
+    2px 8px 10px -1px hsl(255 86% 66% / 0.2);
+  --shadow-2xl: 2px 2px 4px 0px hsl(255 86% 66% / 0.5);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="bold-tech-dark"] {
+  --background: oklch(0.21 0.04 265.75);
+  --foreground: oklch(0.93 0.03 272.79);
+  --card: oklch(0.26 0.09 281.29);
+  --card-foreground: oklch(0.93 0.03 272.79);
+  --popover: oklch(0.26 0.09 281.29);
+  --popover-foreground: oklch(0.93 0.03 272.79);
+  --primary: oklch(0.61 0.22 292.72);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.26 0.09 281.29);
+  --secondary-foreground: oklch(0.93 0.03 272.79);
+  --muted: oklch(0.26 0.09 281.29);
+  --muted-foreground: oklch(0.81 0.1 293.57);
+  --accent: oklch(0.46 0.21 277.02);
+  --accent-foreground: oklch(0.93 0.03 272.79);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.28 0.14 291.09);
+  --input: oklch(0.28 0.14 291.09);
+  --ring: oklch(0.61 0.22 292.72);
+  --chart-1: oklch(0.71 0.16 293.54);
+  --chart-2: oklch(0.61 0.22 292.72);
+  --chart-3: oklch(0.54 0.25 293.01);
+  --chart-4: oklch(0.49 0.24 292.58);
+  --chart-5: oklch(0.43 0.21 292.76);
+  --radius: 0.625rem;
+  --sidebar: oklch(0.21 0.04 265.75);
+  --sidebar-foreground: oklch(0.93 0.03 272.79);
+  --sidebar-primary: oklch(0.61 0.22 292.72);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.46 0.21 277.02);
+  --sidebar-accent-foreground: oklch(0.93 0.03 272.79);
+  --sidebar-border: oklch(0.28 0.14 291.09);
+  --sidebar-ring: oklch(0.61 0.22 292.72);
+  --font-sans: Roboto, sans-serif;
+  --font-serif: Playfair Display, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(255 86% 66%);
+  --shadow-opacity: 0.2;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 2px 4px 0px hsl(255 86% 66% / 0.1);
+  --shadow-xs: 2px 2px 4px 0px hsl(255 86% 66% / 0.1);
+  --shadow-sm:
+    2px 2px 4px 0px hsl(255 86% 66% / 0.2),
+    2px 1px 2px -1px hsl(255 86% 66% / 0.2);
+  --shadow:
+    2px 2px 4px 0px hsl(255 86% 66% / 0.2),
+    2px 1px 2px -1px hsl(255 86% 66% / 0.2);
+  --shadow-md:
+    2px 2px 4px 0px hsl(255 86% 66% / 0.2),
+    2px 2px 4px -1px hsl(255 86% 66% / 0.2);
+  --shadow-lg:
+    2px 2px 4px 0px hsl(255 86% 66% / 0.2),
+    2px 4px 6px -1px hsl(255 86% 66% / 0.2);
+  --shadow-xl:
+    2px 2px 4px 0px hsl(255 86% 66% / 0.2),
+    2px 8px 10px -1px hsl(255 86% 66% / 0.2);
+  --shadow-2xl: 2px 2px 4px 0px hsl(255 86% 66% / 0.5);
+}
diff --git a/web/src/styles/themes/bubblegum.css b/web/src/styles/themes/bubblegum.css
new file mode 100644
index 0000000..740757d
--- /dev/null
+++ b/web/src/styles/themes/bubblegum.css
@@ -0,0 +1,134 @@
+/* Bubblegum Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="bubblegum-light"] {
+  --background: oklch(0.94 0.02 345.7);
+  --foreground: oklch(0.47 0 0);
+  --card: oklch(0.95 0.05 86.89);
+  --card-foreground: oklch(0.47 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.47 0 0);
+  --primary: oklch(0.62 0.18 348.14);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.81 0.07 198.19);
+  --secondary-foreground: oklch(0.32 0 0);
+  --muted: oklch(0.88 0.05 212.1);
+  --muted-foreground: oklch(0.58 0 0);
+  --accent: oklch(0.92 0.08 87.67);
+  --accent-foreground: oklch(0.32 0 0);
+  --destructive: oklch(0.71 0.17 21.96);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.62 0.18 348.14);
+  --input: oklch(0.92 0 0);
+  --ring: oklch(0.7 0.16 350.75);
+  --chart-1: oklch(0.7 0.16 350.75);
+  --chart-2: oklch(0.82 0.08 212.09);
+  --chart-3: oklch(0.92 0.08 87.67);
+  --chart-4: oklch(0.8 0.11 348.18);
+  --chart-5: oklch(0.62 0.19 353.91);
+  --radius: 0.4rem;
+  --sidebar: oklch(0.91 0.04 343.09);
+  --sidebar-foreground: oklch(0.32 0 0);
+  --sidebar-primary: oklch(0.66 0.21 354.31);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.82 0.11 346.02);
+  --sidebar-accent-foreground: oklch(0.32 0 0);
+  --sidebar-border: oklch(0.95 0.03 307.17);
+  --sidebar-ring: oklch(0.66 0.21 354.31);
+  --font-sans: Poppins, sans-serif;
+  --font-serif: Lora, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(325.78 58.18% 56.86% / 0.5);
+  --shadow-opacity: 1;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 3px;
+  --shadow-offset-y: 3px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 3px 3px 0px 0px hsl(325.78 58.18% 56.86% / 0.5);
+  --shadow-xs: 3px 3px 0px 0px hsl(325.78 58.18% 56.86% / 0.5);
+  --shadow-sm:
+    3px 3px 0px 0px hsl(325.78 58.18% 56.86% / 1),
+    3px 1px 2px -1px hsl(325.78 58.18% 56.86% / 1);
+  --shadow:
+    3px 3px 0px 0px hsl(325.78 58.18% 56.86% / 1),
+    3px 1px 2px -1px hsl(325.78 58.18% 56.86% / 1);
+  --shadow-md:
+    3px 3px 0px 0px hsl(325.78 58.18% 56.86% / 1),
+    3px 2px 4px -1px hsl(325.78 58.18% 56.86% / 1);
+  --shadow-lg:
+    3px 3px 0px 0px hsl(325.78 58.18% 56.86% / 1),
+    3px 4px 6px -1px hsl(325.78 58.18% 56.86% / 1);
+  --shadow-xl:
+    3px 3px 0px 0px hsl(325.78 58.18% 56.86% / 1),
+    3px 8px 10px -1px hsl(325.78 58.18% 56.86% / 1);
+  --shadow-2xl: 3px 3px 0px 0px hsl(325.78 58.18% 56.86% / 2.5);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="bubblegum-dark"] {
+  --background: oklch(0.25 0.03 234.16);
+  --foreground: oklch(0.93 0.02 349.08);
+  --card: oklch(0.29 0.03 233.54);
+  --card-foreground: oklch(0.93 0.02 349.08);
+  --popover: oklch(0.29 0.03 233.54);
+  --popover-foreground: oklch(0.93 0.02 349.08);
+  --primary: oklch(0.92 0.08 87.67);
+  --primary-foreground: oklch(0.25 0.03 234.16);
+  --secondary: oklch(0.78 0.08 4.13);
+  --secondary-foreground: oklch(0.25 0.03 234.16);
+  --muted: oklch(0.27 0.01 255.58);
+  --muted-foreground: oklch(0.78 0.08 4.13);
+  --accent: oklch(0.67 0.1 356.98);
+  --accent-foreground: oklch(0.93 0.02 349.08);
+  --destructive: oklch(0.67 0.18 350.36);
+  --destructive-foreground: oklch(0.25 0.03 234.16);
+  --border: oklch(0.39 0.04 242.22);
+  --input: oklch(0.31 0.03 232);
+  --ring: oklch(0.7 0.09 201.87);
+  --chart-1: oklch(0.7 0.09 201.87);
+  --chart-2: oklch(0.78 0.08 4.13);
+  --chart-3: oklch(0.67 0.1 356.98);
+  --chart-4: oklch(0.44 0.07 217.08);
+  --chart-5: oklch(0.27 0.01 255.58);
+  --radius: 0.4rem;
+  --sidebar: oklch(0.23 0.03 235.97);
+  --sidebar-foreground: oklch(0.97 0 264.54);
+  --sidebar-primary: oklch(0.66 0.21 354.31);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.82 0.11 346.02);
+  --sidebar-accent-foreground: oklch(0.28 0.03 256.85);
+  --sidebar-border: oklch(0.37 0.03 259.73);
+  --sidebar-ring: oklch(0.66 0.21 354.31);
+  --font-sans: Poppins, sans-serif;
+  --font-serif: Lora, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: #324859;
+  --shadow-opacity: 1;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 3px;
+  --shadow-offset-y: 3px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 3px 3px 0px 0px hsl(206.15 28.06% 27.25% / 0.5);
+  --shadow-xs: 3px 3px 0px 0px hsl(206.15 28.06% 27.25% / 0.5);
+  --shadow-sm:
+    3px 3px 0px 0px hsl(206.15 28.06% 27.25% / 1),
+    3px 1px 2px -1px hsl(206.15 28.06% 27.25% / 1);
+  --shadow:
+    3px 3px 0px 0px hsl(206.15 28.06% 27.25% / 1),
+    3px 1px 2px -1px hsl(206.15 28.06% 27.25% / 1);
+  --shadow-md:
+    3px 3px 0px 0px hsl(206.15 28.06% 27.25% / 1),
+    3px 2px 4px -1px hsl(206.15 28.06% 27.25% / 1);
+  --shadow-lg:
+    3px 3px 0px 0px hsl(206.15 28.06% 27.25% / 1),
+    3px 4px 6px -1px hsl(206.15 28.06% 27.25% / 1);
+  --shadow-xl:
+    3px 3px 0px 0px hsl(206.15 28.06% 27.25% / 1),
+    3px 8px 10px -1px hsl(206.15 28.06% 27.25% / 1);
+  --shadow-2xl: 3px 3px 0px 0px hsl(206.15 28.06% 27.25% / 2.5);
+}
diff --git a/web/src/styles/themes/caffeine.css b/web/src/styles/themes/caffeine.css
new file mode 100644
index 0000000..ade79d9
--- /dev/null
+++ b/web/src/styles/themes/caffeine.css
@@ -0,0 +1,132 @@
+/* Caffeine Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="caffeine-light"] {
+  --background: oklch(0.98 0 0);
+  --foreground: oklch(0.24 0 0);
+  --card: oklch(0.99 0 0);
+  --card-foreground: oklch(0.24 0 0);
+  --popover: oklch(0.99 0 0);
+  --popover-foreground: oklch(0.24 0 0);
+  --primary: oklch(0.43 0.04 41.99);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.92 0.07 74.37);
+  --secondary-foreground: oklch(0.35 0.07 40.83);
+  --muted: oklch(0.95 0 0);
+  --muted-foreground: oklch(0.5 0 0);
+  --accent: oklch(0.93 0 0);
+  --accent-foreground: oklch(0.24 0 0);
+  --destructive: oklch(0.63 0.19 33.34);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.88 0 0);
+  --input: oklch(0.88 0 0);
+  --ring: oklch(0.43 0.04 41.99);
+  --chart-1: oklch(0.43 0.04 41.99);
+  --chart-2: oklch(0.92 0.07 74.37);
+  --chart-3: oklch(0.93 0 0);
+  --chart-4: oklch(0.94 0.05 75.5);
+  --chart-5: oklch(0.43 0.04 41.67);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.99 0 0);
+  --sidebar-foreground: oklch(0.26 0 0);
+  --sidebar-primary: oklch(0.33 0 0);
+  --sidebar-primary-foreground: oklch(0.99 0 0);
+  --sidebar-accent: oklch(0.98 0 0);
+  --sidebar-accent-foreground: oklch(0.33 0 0);
+  --sidebar-border: oklch(0.94 0 0);
+  --sidebar-ring: oklch(0.77 0 0);
+  --font-sans:
+    ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
+    Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif,
+    "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono:
+    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
+    "Courier New", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="caffeine-dark"] {
+  --background: oklch(0.18 0 0);
+  --foreground: oklch(0.95 0 0);
+  --card: oklch(0.21 0 0);
+  --card-foreground: oklch(0.95 0 0);
+  --popover: oklch(0.21 0 0);
+  --popover-foreground: oklch(0.95 0 0);
+  --primary: oklch(0.92 0.05 66.17);
+  --primary-foreground: oklch(0.2 0.02 200.2);
+  --secondary: oklch(0.32 0.02 63.7);
+  --secondary-foreground: oklch(0.92 0.05 66.17);
+  --muted: oklch(0.25 0 0);
+  --muted-foreground: oklch(0.77 0 0);
+  --accent: oklch(0.29 0 0);
+  --accent-foreground: oklch(0.95 0 0);
+  --destructive: oklch(0.63 0.19 33.34);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.24 0.01 91.75);
+  --input: oklch(0.4 0 0);
+  --ring: oklch(0.92 0.05 66.17);
+  --chart-1: oklch(0.92 0.05 66.17);
+  --chart-2: oklch(0.32 0.02 63.7);
+  --chart-3: oklch(0.29 0 0);
+  --chart-4: oklch(0.35 0.02 67);
+  --chart-5: oklch(0.92 0.05 67.09);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.21 0.01 285.89);
+  --sidebar-foreground: oklch(0.97 0 286.38);
+  --sidebar-primary: oklch(0.49 0.22 264.38);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.27 0.01 286.03);
+  --sidebar-accent-foreground: oklch(0.97 0 286.38);
+  --sidebar-border: oklch(0.27 0.01 286.03);
+  --sidebar-ring: oklch(0.87 0.01 286.29);
+  --font-sans:
+    ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
+    Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif,
+    "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono:
+    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
+    "Courier New", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/candyland.css b/web/src/styles/themes/candyland.css
new file mode 100644
index 0000000..a729fef
--- /dev/null
+++ b/web/src/styles/themes/candyland.css
@@ -0,0 +1,122 @@
+/* Candyland Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="candyland-light"] {
+  --background: oklch(0.98 0 228.78);
+  --foreground: oklch(0.32 0 0);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.32 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.32 0 0);
+  --primary: oklch(0.87 0.07 7.09);
+  --primary-foreground: oklch(0 0 0);
+  --secondary: oklch(0.81 0.08 225.75);
+  --secondary-foreground: oklch(0 0 0);
+  --muted: oklch(0.88 0.03 98.1);
+  --muted-foreground: oklch(0.54 0 0);
+  --accent: oklch(0.97 0.21 109.77);
+  --accent-foreground: oklch(0 0 0);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.87 0 0);
+  --input: oklch(0.87 0 0);
+  --ring: oklch(0.87 0.07 7.09);
+  --chart-1: oklch(0.87 0.07 7.09);
+  --chart-2: oklch(0.81 0.08 225.75);
+  --chart-3: oklch(0.97 0.21 109.77);
+  --chart-4: oklch(0.8 0.14 349.23);
+  --chart-5: oklch(0.74 0.23 142.85);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.98 0 228.78);
+  --sidebar-foreground: oklch(0.32 0 0);
+  --sidebar-primary: oklch(0.87 0.07 7.09);
+  --sidebar-primary-foreground: oklch(0 0 0);
+  --sidebar-accent: oklch(0.97 0.21 109.77);
+  --sidebar-accent-foreground: oklch(0 0 0);
+  --sidebar-border: oklch(0.87 0 0);
+  --sidebar-ring: oklch(0.87 0.07 7.09);
+  --font-sans: Poppins, sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: Roboto Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="candyland-dark"] {
+  --background: oklch(0.23 0.01 264.29);
+  --foreground: oklch(0.92 0 0);
+  --card: oklch(0.32 0.01 223.67);
+  --card-foreground: oklch(0.92 0 0);
+  --popover: oklch(0.32 0.01 223.67);
+  --popover-foreground: oklch(0.92 0 0);
+  --primary: oklch(0.8 0.14 349.23);
+  --primary-foreground: oklch(0 0 0);
+  --secondary: oklch(0.74 0.23 142.85);
+  --secondary-foreground: oklch(0 0 0);
+  --muted: oklch(0.39 0 0);
+  --muted-foreground: oklch(0.72 0 0);
+  --accent: oklch(0.81 0.08 225.75);
+  --accent-foreground: oklch(0 0 0);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.39 0 0);
+  --input: oklch(0.39 0 0);
+  --ring: oklch(0.8 0.14 349.23);
+  --chart-1: oklch(0.8 0.14 349.23);
+  --chart-2: oklch(0.74 0.23 142.85);
+  --chart-3: oklch(0.81 0.08 225.75);
+  --chart-4: oklch(0.97 0.21 109.77);
+  --chart-5: oklch(0.87 0.18 90.38);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.23 0.01 264.29);
+  --sidebar-foreground: oklch(0.92 0 0);
+  --sidebar-primary: oklch(0.8 0.14 349.23);
+  --sidebar-primary-foreground: oklch(0 0 0);
+  --sidebar-accent: oklch(0.81 0.08 225.75);
+  --sidebar-accent-foreground: oklch(0 0 0);
+  --sidebar-border: oklch(0.39 0 0);
+  --sidebar-ring: oklch(0.8 0.14 349.23);
+  --font-sans: Poppins, sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: Roboto Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/catppuccin.css b/web/src/styles/themes/catppuccin.css
new file mode 100644
index 0000000..64dc0db
--- /dev/null
+++ b/web/src/styles/themes/catppuccin.css
@@ -0,0 +1,134 @@
+/* Catppuccin Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="catppuccin-light"] {
+  --background: oklch(0.96 0.01 264.53);
+  --foreground: oklch(0.44 0.04 279.33);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.44 0.04 279.33);
+  --popover: oklch(0.86 0.01 268.48);
+  --popover-foreground: oklch(0.44 0.04 279.33);
+  --primary: oklch(0.55 0.25 297.02);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.86 0.01 268.48);
+  --secondary-foreground: oklch(0.44 0.04 279.33);
+  --muted: oklch(0.91 0.01 264.51);
+  --muted-foreground: oklch(0.55 0.03 279.08);
+  --accent: oklch(0.68 0.14 235.38);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.55 0.22 19.81);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.81 0.02 271.2);
+  --input: oklch(0.86 0.01 268.48);
+  --ring: oklch(0.55 0.25 297.02);
+  --chart-1: oklch(0.55 0.25 297.02);
+  --chart-2: oklch(0.68 0.14 235.38);
+  --chart-3: oklch(0.63 0.18 140.44);
+  --chart-4: oklch(0.69 0.2 42.43);
+  --chart-5: oklch(0.71 0.1 33.1);
+  --radius: 0.35rem;
+  --sidebar: oklch(0.93 0.01 264.52);
+  --sidebar-foreground: oklch(0.44 0.04 279.33);
+  --sidebar-primary: oklch(0.55 0.25 297.02);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.68 0.14 235.38);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0.81 0.02 271.2);
+  --sidebar-ring: oklch(0.55 0.25 297.02);
+  --font-sans: Montserrat, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(240 30% 25%);
+  --shadow-opacity: 0.12;
+  --shadow-blur: 6px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 6px 0px hsl(240 30% 25% / 0.06);
+  --shadow-xs: 0px 4px 6px 0px hsl(240 30% 25% / 0.06);
+  --shadow-sm:
+    0px 4px 6px 0px hsl(240 30% 25% / 0.12),
+    0px 1px 2px -1px hsl(240 30% 25% / 0.12);
+  --shadow:
+    0px 4px 6px 0px hsl(240 30% 25% / 0.12),
+    0px 1px 2px -1px hsl(240 30% 25% / 0.12);
+  --shadow-md:
+    0px 4px 6px 0px hsl(240 30% 25% / 0.12),
+    0px 2px 4px -1px hsl(240 30% 25% / 0.12);
+  --shadow-lg:
+    0px 4px 6px 0px hsl(240 30% 25% / 0.12),
+    0px 4px 6px -1px hsl(240 30% 25% / 0.12);
+  --shadow-xl:
+    0px 4px 6px 0px hsl(240 30% 25% / 0.12),
+    0px 8px 10px -1px hsl(240 30% 25% / 0.12);
+  --shadow-2xl: 0px 4px 6px 0px hsl(240 30% 25% / 0.3);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="catppuccin-dark"] {
+  --background: oklch(0.22 0.03 284.06);
+  --foreground: oklch(0.88 0.04 272.28);
+  --card: oklch(0.24 0.03 283.91);
+  --card-foreground: oklch(0.88 0.04 272.28);
+  --popover: oklch(0.4 0.03 280.15);
+  --popover-foreground: oklch(0.88 0.04 272.28);
+  --primary: oklch(0.79 0.12 304.77);
+  --primary-foreground: oklch(0.24 0.03 283.91);
+  --secondary: oklch(0.48 0.03 278.64);
+  --secondary-foreground: oklch(0.88 0.04 272.28);
+  --muted: oklch(0.3 0.03 276.21);
+  --muted-foreground: oklch(0.75 0.04 273.93);
+  --accent: oklch(0.85 0.08 210.25);
+  --accent-foreground: oklch(0.24 0.03 283.91);
+  --destructive: oklch(0.76 0.13 2.76);
+  --destructive-foreground: oklch(0.24 0.03 283.91);
+  --border: oklch(0.32 0.03 281.98);
+  --input: oklch(0.32 0.03 281.98);
+  --ring: oklch(0.79 0.12 304.77);
+  --chart-1: oklch(0.79 0.12 304.77);
+  --chart-2: oklch(0.85 0.08 210.25);
+  --chart-3: oklch(0.86 0.11 142.72);
+  --chart-4: oklch(0.82 0.1 52.63);
+  --chart-5: oklch(0.92 0.02 30.49);
+  --radius: 0.35rem;
+  --sidebar: oklch(0.18 0.02 284.2);
+  --sidebar-foreground: oklch(0.88 0.04 272.28);
+  --sidebar-primary: oklch(0.79 0.12 304.77);
+  --sidebar-primary-foreground: oklch(0.24 0.03 283.91);
+  --sidebar-accent: oklch(0.85 0.08 210.25);
+  --sidebar-accent-foreground: oklch(0.24 0.03 283.91);
+  --sidebar-border: oklch(0.4 0.03 280.15);
+  --sidebar-ring: oklch(0.79 0.12 304.77);
+  --font-sans: Montserrat, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(240 30% 25%);
+  --shadow-opacity: 0.12;
+  --shadow-blur: 6px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 6px 0px hsl(240 30% 25% / 0.06);
+  --shadow-xs: 0px 4px 6px 0px hsl(240 30% 25% / 0.06);
+  --shadow-sm:
+    0px 4px 6px 0px hsl(240 30% 25% / 0.12),
+    0px 1px 2px -1px hsl(240 30% 25% / 0.12);
+  --shadow:
+    0px 4px 6px 0px hsl(240 30% 25% / 0.12),
+    0px 1px 2px -1px hsl(240 30% 25% / 0.12);
+  --shadow-md:
+    0px 4px 6px 0px hsl(240 30% 25% / 0.12),
+    0px 2px 4px -1px hsl(240 30% 25% / 0.12);
+  --shadow-lg:
+    0px 4px 6px 0px hsl(240 30% 25% / 0.12),
+    0px 4px 6px -1px hsl(240 30% 25% / 0.12);
+  --shadow-xl:
+    0px 4px 6px 0px hsl(240 30% 25% / 0.12),
+    0px 8px 10px -1px hsl(240 30% 25% / 0.12);
+  --shadow-2xl: 0px 4px 6px 0px hsl(240 30% 25% / 0.3);
+}
diff --git a/web/src/styles/themes/claude.css b/web/src/styles/themes/claude.css
new file mode 100644
index 0000000..b681193
--- /dev/null
+++ b/web/src/styles/themes/claude.css
@@ -0,0 +1,132 @@
+/* Claude Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="claude-light"] {
+  --background: oklch(0.98 0.01 95.1);
+  --foreground: oklch(0.34 0.03 95.72);
+  --card: oklch(0.98 0.01 95.1);
+  --card-foreground: oklch(0.19 0 106.59);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.27 0.02 98.94);
+  --primary: oklch(0.62 0.14 39.04);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.92 0.01 92.99);
+  --secondary-foreground: oklch(0.43 0.02 98.6);
+  --muted: oklch(0.93 0.02 90.24);
+  --muted-foreground: oklch(0.61 0.01 97.42);
+  --accent: oklch(0.92 0.01 92.99);
+  --accent-foreground: oklch(0.27 0.02 98.94);
+  --destructive: oklch(0.19 0 106.59);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.88 0.01 97.36);
+  --input: oklch(0.76 0.02 98.35);
+  --ring: oklch(0.59 0.17 253.06);
+  --chart-1: oklch(0.56 0.13 43);
+  --chart-2: oklch(0.69 0.16 290.41);
+  --chart-3: oklch(0.88 0.03 93.13);
+  --chart-4: oklch(0.88 0.04 298.18);
+  --chart-5: oklch(0.56 0.13 42.06);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.97 0.01 98.88);
+  --sidebar-foreground: oklch(0.36 0.01 106.65);
+  --sidebar-primary: oklch(0.62 0.14 39.04);
+  --sidebar-primary-foreground: oklch(0.99 0 0);
+  --sidebar-accent: oklch(0.92 0.01 92.99);
+  --sidebar-accent-foreground: oklch(0.33 0 0);
+  --sidebar-border: oklch(0.94 0 0);
+  --sidebar-ring: oklch(0.77 0 0);
+  --font-sans:
+    ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
+    Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif,
+    "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono:
+    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
+    "Courier New", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="claude-dark"] {
+  --background: oklch(0.27 0 106.64);
+  --foreground: oklch(0.81 0.01 93.01);
+  --card: oklch(0.27 0 106.64);
+  --card-foreground: oklch(0.98 0.01 95.1);
+  --popover: oklch(0.31 0 106.6);
+  --popover-foreground: oklch(0.92 0 106.48);
+  --primary: oklch(0.67 0.13 38.76);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.98 0.01 95.1);
+  --secondary-foreground: oklch(0.31 0 106.6);
+  --muted: oklch(0.22 0 106.71);
+  --muted-foreground: oklch(0.77 0.02 99.07);
+  --accent: oklch(0.21 0.01 95.42);
+  --accent-foreground: oklch(0.97 0.01 98.88);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.36 0.01 106.89);
+  --input: oklch(0.43 0.01 100.22);
+  --ring: oklch(0.59 0.17 253.06);
+  --chart-1: oklch(0.56 0.13 43);
+  --chart-2: oklch(0.69 0.16 290.41);
+  --chart-3: oklch(0.21 0.01 95.42);
+  --chart-4: oklch(0.31 0.05 289.32);
+  --chart-5: oklch(0.56 0.13 42.06);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.24 0 67.71);
+  --sidebar-foreground: oklch(0.81 0.01 93.01);
+  --sidebar-primary: oklch(0.33 0 0);
+  --sidebar-primary-foreground: oklch(0.99 0 0);
+  --sidebar-accent: oklch(0.17 0 106.62);
+  --sidebar-accent-foreground: oklch(0.81 0.01 93.01);
+  --sidebar-border: oklch(0.94 0 0);
+  --sidebar-ring: oklch(0.77 0 0);
+  --font-sans:
+    ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
+    Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif,
+    "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono:
+    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
+    "Courier New", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/claymorphism.css b/web/src/styles/themes/claymorphism.css
new file mode 100644
index 0000000..82fd4f7
--- /dev/null
+++ b/web/src/styles/themes/claymorphism.css
@@ -0,0 +1,129 @@
+/* Claymorphism Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="claymorphism-light"] {
+  --background: oklch(0.92 0 48.72);
+  --foreground: oklch(0.28 0.04 260.03);
+  --card: oklch(0.97 0 106.42);
+  --card-foreground: oklch(0.28 0.04 260.03);
+  --popover: oklch(0.97 0 106.42);
+  --popover-foreground: oklch(0.28 0.04 260.03);
+  --primary: oklch(0.59 0.2 277.12);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.87 0 56.37);
+  --secondary-foreground: oklch(0.45 0.03 256.8);
+  --muted: oklch(0.92 0 48.72);
+  --muted-foreground: oklch(0.55 0.02 264.36);
+  --accent: oklch(0.94 0.03 321.94);
+  --accent-foreground: oklch(0.37 0.03 259.73);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.87 0 56.37);
+  --input: oklch(0.87 0 56.37);
+  --ring: oklch(0.59 0.2 277.12);
+  --chart-1: oklch(0.59 0.2 277.12);
+  --chart-2: oklch(0.51 0.23 276.97);
+  --chart-3: oklch(0.46 0.21 277.02);
+  --chart-4: oklch(0.4 0.18 277.37);
+  --chart-5: oklch(0.36 0.14 278.7);
+  --radius: 1.25rem;
+  --sidebar: oklch(0.87 0 56.37);
+  --sidebar-foreground: oklch(0.28 0.04 260.03);
+  --sidebar-primary: oklch(0.59 0.2 277.12);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.94 0.03 321.94);
+  --sidebar-accent-foreground: oklch(0.37 0.03 259.73);
+  --sidebar-border: oklch(0.87 0 56.37);
+  --sidebar-ring: oklch(0.59 0.2 277.12);
+  --font-sans: Plus Jakarta Sans, sans-serif;
+  --font-serif: Lora, serif;
+  --font-mono: Roboto Mono, monospace;
+  --shadow-color: hsl(240 4% 60%);
+  --shadow-opacity: 0.18;
+  --shadow-blur: 10px;
+  --shadow-spread: 4px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 2px 10px 4px hsl(240 4% 60% / 0.09);
+  --shadow-xs: 2px 2px 10px 4px hsl(240 4% 60% / 0.09);
+  --shadow-sm:
+    2px 2px 10px 4px hsl(240 4% 60% / 0.18),
+    2px 1px 2px 3px hsl(240 4% 60% / 0.18);
+  --shadow:
+    2px 2px 10px 4px hsl(240 4% 60% / 0.18),
+    2px 1px 2px 3px hsl(240 4% 60% / 0.18);
+  --shadow-md:
+    2px 2px 10px 4px hsl(240 4% 60% / 0.18),
+    2px 2px 4px 3px hsl(240 4% 60% / 0.18);
+  --shadow-lg:
+    2px 2px 10px 4px hsl(240 4% 60% / 0.18),
+    2px 4px 6px 3px hsl(240 4% 60% / 0.18);
+  --shadow-xl:
+    2px 2px 10px 4px hsl(240 4% 60% / 0.18),
+    2px 8px 10px 3px hsl(240 4% 60% / 0.18);
+  --shadow-2xl: 2px 2px 10px 4px hsl(240 4% 60% / 0.45);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="claymorphism-dark"] {
+  --background: oklch(0.22 0.01 67.44);
+  --foreground: oklch(0.93 0.01 255.51);
+  --card: oklch(0.28 0.01 59.34);
+  --card-foreground: oklch(0.93 0.01 255.51);
+  --popover: oklch(0.28 0.01 59.34);
+  --popover-foreground: oklch(0.93 0.01 255.51);
+  --primary: oklch(0.68 0.16 276.93);
+  --primary-foreground: oklch(0.22 0.01 67.44);
+  --secondary: oklch(0.34 0.01 59.42);
+  --secondary-foreground: oklch(0.87 0.01 258.34);
+  --muted: oklch(0.28 0.01 59.34);
+  --muted-foreground: oklch(0.71 0.02 261.32);
+  --accent: oklch(0.39 0.01 59.47);
+  --accent-foreground: oklch(0.87 0.01 258.34);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(0.22 0.01 67.44);
+  --border: oklch(0.34 0.01 59.42);
+  --input: oklch(0.34 0.01 59.42);
+  --ring: oklch(0.68 0.16 276.93);
+  --chart-1: oklch(0.68 0.16 276.93);
+  --chart-2: oklch(0.59 0.2 277.12);
+  --chart-3: oklch(0.51 0.23 276.97);
+  --chart-4: oklch(0.46 0.21 277.02);
+  --chart-5: oklch(0.4 0.18 277.37);
+  --radius: 1.25rem;
+  --sidebar: oklch(0.34 0.01 59.42);
+  --sidebar-foreground: oklch(0.93 0.01 255.51);
+  --sidebar-primary: oklch(0.68 0.16 276.93);
+  --sidebar-primary-foreground: oklch(0.22 0.01 67.44);
+  --sidebar-accent: oklch(0.39 0.01 59.47);
+  --sidebar-accent-foreground: oklch(0.87 0.01 258.34);
+  --sidebar-border: oklch(0.34 0.01 59.42);
+  --sidebar-ring: oklch(0.68 0.16 276.93);
+  --font-sans: Plus Jakarta Sans, sans-serif;
+  --font-serif: Lora, serif;
+  --font-mono: Roboto Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.18;
+  --shadow-blur: 10px;
+  --shadow-spread: 4px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 2px 10px 4px hsl(0 0% 0% / 0.09);
+  --shadow-xs: 2px 2px 10px 4px hsl(0 0% 0% / 0.09);
+  --shadow-sm:
+    2px 2px 10px 4px hsl(0 0% 0% / 0.18), 2px 1px 2px 3px hsl(0 0% 0% / 0.18);
+  --shadow:
+    2px 2px 10px 4px hsl(0 0% 0% / 0.18), 2px 1px 2px 3px hsl(0 0% 0% / 0.18);
+  --shadow-md:
+    2px 2px 10px 4px hsl(0 0% 0% / 0.18), 2px 2px 4px 3px hsl(0 0% 0% / 0.18);
+  --shadow-lg:
+    2px 2px 10px 4px hsl(0 0% 0% / 0.18), 2px 4px 6px 3px hsl(0 0% 0% / 0.18);
+  --shadow-xl:
+    2px 2px 10px 4px hsl(0 0% 0% / 0.18), 2px 8px 10px 3px hsl(0 0% 0% / 0.18);
+  --shadow-2xl: 2px 2px 10px 4px hsl(0 0% 0% / 0.45);
+}
diff --git a/web/src/styles/themes/clean-slate.css b/web/src/styles/themes/clean-slate.css
new file mode 100644
index 0000000..606da96
--- /dev/null
+++ b/web/src/styles/themes/clean-slate.css
@@ -0,0 +1,124 @@
+/* Clean Slate Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="clean-slate-light"] {
+  --background: oklch(0.98 0 247.86);
+  --foreground: oklch(0.28 0.04 260.03);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.28 0.04 260.03);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.28 0.04 260.03);
+  --primary: oklch(0.59 0.2 277.12);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.93 0.01 264.53);
+  --secondary-foreground: oklch(0.37 0.03 259.73);
+  --muted: oklch(0.97 0 264.54);
+  --muted-foreground: oklch(0.55 0.02 264.36);
+  --accent: oklch(0.93 0.03 272.79);
+  --accent-foreground: oklch(0.37 0.03 259.73);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.87 0.01 258.34);
+  --input: oklch(0.87 0.01 258.34);
+  --ring: oklch(0.59 0.2 277.12);
+  --chart-1: oklch(0.59 0.2 277.12);
+  --chart-2: oklch(0.51 0.23 276.97);
+  --chart-3: oklch(0.46 0.21 277.02);
+  --chart-4: oklch(0.4 0.18 277.37);
+  --chart-5: oklch(0.36 0.14 278.7);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.97 0 264.54);
+  --sidebar-foreground: oklch(0.28 0.04 260.03);
+  --sidebar-primary: oklch(0.59 0.2 277.12);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.93 0.03 272.79);
+  --sidebar-accent-foreground: oklch(0.37 0.03 259.73);
+  --sidebar-border: oklch(0.87 0.01 258.34);
+  --sidebar-ring: oklch(0.59 0.2 277.12);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Merriweather, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 8px;
+  --shadow-spread: -1px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 2px 4px -2px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 4px 6px -2px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 8px 10px -2px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 4px 8px -1px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="clean-slate-dark"] {
+  --background: oklch(0.21 0.04 265.75);
+  --foreground: oklch(0.93 0.01 255.51);
+  --card: oklch(0.28 0.04 260.03);
+  --card-foreground: oklch(0.93 0.01 255.51);
+  --popover: oklch(0.28 0.04 260.03);
+  --popover-foreground: oklch(0.93 0.01 255.51);
+  --primary: oklch(0.68 0.16 276.93);
+  --primary-foreground: oklch(0.21 0.04 265.75);
+  --secondary: oklch(0.34 0.03 260.91);
+  --secondary-foreground: oklch(0.87 0.01 258.34);
+  --muted: oklch(0.28 0.04 260.03);
+  --muted-foreground: oklch(0.71 0.02 261.32);
+  --accent: oklch(0.37 0.03 259.73);
+  --accent-foreground: oklch(0.87 0.01 258.34);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(0.21 0.04 265.75);
+  --border: oklch(0.45 0.03 256.8);
+  --input: oklch(0.45 0.03 256.8);
+  --ring: oklch(0.68 0.16 276.93);
+  --chart-1: oklch(0.68 0.16 276.93);
+  --chart-2: oklch(0.59 0.2 277.12);
+  --chart-3: oklch(0.51 0.23 276.97);
+  --chart-4: oklch(0.46 0.21 277.02);
+  --chart-5: oklch(0.4 0.18 277.37);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.28 0.04 260.03);
+  --sidebar-foreground: oklch(0.93 0.01 255.51);
+  --sidebar-primary: oklch(0.68 0.16 276.93);
+  --sidebar-primary-foreground: oklch(0.21 0.04 265.75);
+  --sidebar-accent: oklch(0.37 0.03 259.73);
+  --sidebar-accent-foreground: oklch(0.87 0.01 258.34);
+  --sidebar-border: oklch(0.45 0.03 256.8);
+  --sidebar-ring: oklch(0.68 0.16 276.93);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Merriweather, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 8px;
+  --shadow-spread: -1px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 2px 4px -2px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 4px 6px -2px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 8px 10px -2px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 4px 8px -1px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/cosmic-night.css b/web/src/styles/themes/cosmic-night.css
new file mode 100644
index 0000000..ee9c21f
--- /dev/null
+++ b/web/src/styles/themes/cosmic-night.css
@@ -0,0 +1,134 @@
+/* Cosmic Night Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="cosmic-night-light"] {
+  --background: oklch(0.97 0.01 286.15);
+  --foreground: oklch(0.3 0.06 282.42);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.3 0.06 282.42);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.3 0.06 282.42);
+  --primary: oklch(0.54 0.18 288.03);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.92 0.04 292.69);
+  --secondary-foreground: oklch(0.41 0.1 288.17);
+  --muted: oklch(0.96 0.01 286.15);
+  --muted-foreground: oklch(0.54 0.05 284.74);
+  --accent: oklch(0.92 0.04 262.14);
+  --accent-foreground: oklch(0.3 0.06 282.42);
+  --destructive: oklch(0.69 0.21 14.99);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.91 0.02 285.96);
+  --input: oklch(0.91 0.02 285.96);
+  --ring: oklch(0.54 0.18 288.03);
+  --chart-1: oklch(0.54 0.18 288.03);
+  --chart-2: oklch(0.7 0.16 288.99);
+  --chart-3: oklch(0.57 0.21 276.71);
+  --chart-4: oklch(0.64 0.19 281.81);
+  --chart-5: oklch(0.45 0.18 279.38);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.96 0.01 286.15);
+  --sidebar-foreground: oklch(0.3 0.06 282.42);
+  --sidebar-primary: oklch(0.54 0.18 288.03);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.92 0.04 262.14);
+  --sidebar-accent-foreground: oklch(0.3 0.06 282.42);
+  --sidebar-border: oklch(0.91 0.02 285.96);
+  --sidebar-ring: oklch(0.54 0.18 288.03);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(240 30% 25%);
+  --shadow-opacity: 0.12;
+  --shadow-blur: 10px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 10px 0px hsl(240 30% 25% / 0.06);
+  --shadow-xs: 0px 4px 10px 0px hsl(240 30% 25% / 0.06);
+  --shadow-sm:
+    0px 4px 10px 0px hsl(240 30% 25% / 0.12),
+    0px 1px 2px -1px hsl(240 30% 25% / 0.12);
+  --shadow:
+    0px 4px 10px 0px hsl(240 30% 25% / 0.12),
+    0px 1px 2px -1px hsl(240 30% 25% / 0.12);
+  --shadow-md:
+    0px 4px 10px 0px hsl(240 30% 25% / 0.12),
+    0px 2px 4px -1px hsl(240 30% 25% / 0.12);
+  --shadow-lg:
+    0px 4px 10px 0px hsl(240 30% 25% / 0.12),
+    0px 4px 6px -1px hsl(240 30% 25% / 0.12);
+  --shadow-xl:
+    0px 4px 10px 0px hsl(240 30% 25% / 0.12),
+    0px 8px 10px -1px hsl(240 30% 25% / 0.12);
+  --shadow-2xl: 0px 4px 10px 0px hsl(240 30% 25% / 0.3);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="cosmic-night-dark"] {
+  --background: oklch(0.17 0.02 283.8);
+  --foreground: oklch(0.92 0.03 285.88);
+  --card: oklch(0.23 0.04 282.93);
+  --card-foreground: oklch(0.92 0.03 285.88);
+  --popover: oklch(0.23 0.04 282.93);
+  --popover-foreground: oklch(0.92 0.03 285.88);
+  --primary: oklch(0.72 0.16 290.4);
+  --primary-foreground: oklch(0.17 0.02 283.8);
+  --secondary: oklch(0.31 0.07 283.46);
+  --secondary-foreground: oklch(0.84 0.08 285.91);
+  --muted: oklch(0.27 0.06 281.44);
+  --muted-foreground: oklch(0.72 0.05 285.17);
+  --accent: oklch(0.34 0.08 280.97);
+  --accent-foreground: oklch(0.92 0.03 285.88);
+  --destructive: oklch(0.69 0.21 14.99);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.33 0.06 282.58);
+  --input: oklch(0.33 0.06 282.58);
+  --ring: oklch(0.72 0.16 290.4);
+  --chart-1: oklch(0.72 0.16 290.4);
+  --chart-2: oklch(0.64 0.1 274.91);
+  --chart-3: oklch(0.75 0.12 244.75);
+  --chart-4: oklch(0.71 0.1 186.68);
+  --chart-5: oklch(0.75 0.18 346.81);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.23 0.04 282.93);
+  --sidebar-foreground: oklch(0.92 0.03 285.88);
+  --sidebar-primary: oklch(0.72 0.16 290.4);
+  --sidebar-primary-foreground: oklch(0.17 0.02 283.8);
+  --sidebar-accent: oklch(0.34 0.08 280.97);
+  --sidebar-accent-foreground: oklch(0.92 0.03 285.88);
+  --sidebar-border: oklch(0.33 0.06 282.58);
+  --sidebar-ring: oklch(0.72 0.16 290.4);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(240 30% 25%);
+  --shadow-opacity: 0.12;
+  --shadow-blur: 10px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 10px 0px hsl(240 30% 25% / 0.06);
+  --shadow-xs: 0px 4px 10px 0px hsl(240 30% 25% / 0.06);
+  --shadow-sm:
+    0px 4px 10px 0px hsl(240 30% 25% / 0.12),
+    0px 1px 2px -1px hsl(240 30% 25% / 0.12);
+  --shadow:
+    0px 4px 10px 0px hsl(240 30% 25% / 0.12),
+    0px 1px 2px -1px hsl(240 30% 25% / 0.12);
+  --shadow-md:
+    0px 4px 10px 0px hsl(240 30% 25% / 0.12),
+    0px 2px 4px -1px hsl(240 30% 25% / 0.12);
+  --shadow-lg:
+    0px 4px 10px 0px hsl(240 30% 25% / 0.12),
+    0px 4px 6px -1px hsl(240 30% 25% / 0.12);
+  --shadow-xl:
+    0px 4px 10px 0px hsl(240 30% 25% / 0.12),
+    0px 8px 10px -1px hsl(240 30% 25% / 0.12);
+  --shadow-2xl: 0px 4px 10px 0px hsl(240 30% 25% / 0.3);
+}
diff --git a/web/src/styles/themes/cyberpunk.css b/web/src/styles/themes/cyberpunk.css
new file mode 100644
index 0000000..a6b1560
--- /dev/null
+++ b/web/src/styles/themes/cyberpunk.css
@@ -0,0 +1,124 @@
+/* Cyberpunk Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="cyberpunk-light"] {
+  --background: oklch(0.98 0 247.84);
+  --foreground: oklch(0.16 0.04 281.83);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.16 0.04 281.83);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.16 0.04 281.83);
+  --primary: oklch(0.67 0.29 341.41);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.96 0.02 286.02);
+  --secondary-foreground: oklch(0.16 0.04 281.83);
+  --muted: oklch(0.96 0.02 286.02);
+  --muted-foreground: oklch(0.16 0.04 281.83);
+  --accent: oklch(0.89 0.17 171.27);
+  --accent-foreground: oklch(0.16 0.04 281.83);
+  --destructive: oklch(0.65 0.23 34.04);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.92 0.01 225.09);
+  --input: oklch(0.92 0.01 225.09);
+  --ring: oklch(0.67 0.29 341.41);
+  --chart-1: oklch(0.67 0.29 341.41);
+  --chart-2: oklch(0.55 0.29 299.1);
+  --chart-3: oklch(0.84 0.15 209.29);
+  --chart-4: oklch(0.89 0.17 171.27);
+  --chart-5: oklch(0.92 0.19 101.41);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.96 0.02 286.02);
+  --sidebar-foreground: oklch(0.16 0.04 281.83);
+  --sidebar-primary: oklch(0.67 0.29 341.41);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.89 0.17 171.27);
+  --sidebar-accent-foreground: oklch(0.16 0.04 281.83);
+  --sidebar-border: oklch(0.92 0.01 225.09);
+  --sidebar-ring: oklch(0.67 0.29 341.41);
+  --font-sans: Outfit, sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 8px;
+  --shadow-spread: -2px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 8px -2px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 4px 8px -2px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 4px 8px -2px hsl(0 0% 0% / 0.1), 0px 1px 2px -3px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 4px 8px -2px hsl(0 0% 0% / 0.1), 0px 1px 2px -3px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 4px 8px -2px hsl(0 0% 0% / 0.1), 0px 2px 4px -3px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 4px 8px -2px hsl(0 0% 0% / 0.1), 0px 4px 6px -3px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 4px 8px -2px hsl(0 0% 0% / 0.1), 0px 8px 10px -3px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 4px 8px -2px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="cyberpunk-dark"] {
+  --background: oklch(0.16 0.04 281.83);
+  --foreground: oklch(0.95 0.01 260.73);
+  --card: oklch(0.25 0.06 281.14);
+  --card-foreground: oklch(0.95 0.01 260.73);
+  --popover: oklch(0.25 0.06 281.14);
+  --popover-foreground: oklch(0.95 0.01 260.73);
+  --primary: oklch(0.67 0.29 341.41);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.25 0.06 281.14);
+  --secondary-foreground: oklch(0.95 0.01 260.73);
+  --muted: oklch(0.25 0.06 281.14);
+  --muted-foreground: oklch(0.62 0.05 278.1);
+  --accent: oklch(0.89 0.17 171.27);
+  --accent-foreground: oklch(0.16 0.04 281.83);
+  --destructive: oklch(0.65 0.23 34.04);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.33 0.08 280.79);
+  --input: oklch(0.33 0.08 280.79);
+  --ring: oklch(0.67 0.29 341.41);
+  --chart-1: oklch(0.67 0.29 341.41);
+  --chart-2: oklch(0.55 0.29 299.1);
+  --chart-3: oklch(0.84 0.15 209.29);
+  --chart-4: oklch(0.89 0.17 171.27);
+  --chart-5: oklch(0.92 0.19 101.41);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.16 0.04 281.83);
+  --sidebar-foreground: oklch(0.95 0.01 260.73);
+  --sidebar-primary: oklch(0.67 0.29 341.41);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.89 0.17 171.27);
+  --sidebar-accent-foreground: oklch(0.16 0.04 281.83);
+  --sidebar-border: oklch(0.33 0.08 280.79);
+  --sidebar-ring: oklch(0.67 0.29 341.41);
+  --font-sans: Outfit, sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 8px;
+  --shadow-spread: -2px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 8px -2px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 4px 8px -2px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 4px 8px -2px hsl(0 0% 0% / 0.1), 0px 1px 2px -3px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 4px 8px -2px hsl(0 0% 0% / 0.1), 0px 1px 2px -3px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 4px 8px -2px hsl(0 0% 0% / 0.1), 0px 2px 4px -3px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 4px 8px -2px hsl(0 0% 0% / 0.1), 0px 4px 6px -3px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 4px 8px -2px hsl(0 0% 0% / 0.1), 0px 8px 10px -3px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 4px 8px -2px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/doom-64.css b/web/src/styles/themes/doom-64.css
new file mode 100644
index 0000000..f7503f2
--- /dev/null
+++ b/web/src/styles/themes/doom-64.css
@@ -0,0 +1,124 @@
+/* Doom 64 Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="doom-64-light"] {
+  --background: oklch(0.85 0 0);
+  --foreground: oklch(0.24 0 0);
+  --card: oklch(0.76 0 0);
+  --card-foreground: oklch(0.24 0 0);
+  --popover: oklch(0.76 0 0);
+  --popover-foreground: oklch(0.24 0 0);
+  --primary: oklch(0.5 0.19 27.48);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.5 0.09 126.19);
+  --secondary-foreground: oklch(1 0 0);
+  --muted: oklch(0.78 0 0);
+  --muted-foreground: oklch(0.41 0 0);
+  --accent: oklch(0.59 0.1 245.74);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.71 0.2 46.46);
+  --destructive-foreground: oklch(0 0 0);
+  --border: oklch(0.43 0 0);
+  --input: oklch(0.43 0 0);
+  --ring: oklch(0.5 0.19 27.48);
+  --chart-1: oklch(0.5 0.19 27.48);
+  --chart-2: oklch(0.5 0.09 126.19);
+  --chart-3: oklch(0.59 0.1 245.74);
+  --chart-4: oklch(0.71 0.2 46.46);
+  --chart-5: oklch(0.57 0.04 40.43);
+  --radius: 0px;
+  --sidebar: oklch(0.76 0 0);
+  --sidebar-foreground: oklch(0.24 0 0);
+  --sidebar-primary: oklch(0.5 0.19 27.48);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.59 0.1 245.74);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0.43 0 0);
+  --sidebar-ring: oklch(0.5 0.19 27.48);
+  --font-sans: "Oxanium", sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: "Source Code Pro", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.4;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 2px 4px 0px hsl(0 0% 0% / 0.2);
+  --shadow-xs: 0px 2px 4px 0px hsl(0 0% 0% / 0.2);
+  --shadow-sm:
+    0px 2px 4px 0px hsl(0 0% 0% / 0.4), 0px 1px 2px -1px hsl(0 0% 0% / 0.4);
+  --shadow:
+    0px 2px 4px 0px hsl(0 0% 0% / 0.4), 0px 1px 2px -1px hsl(0 0% 0% / 0.4);
+  --shadow-md:
+    0px 2px 4px 0px hsl(0 0% 0% / 0.4), 0px 2px 4px -1px hsl(0 0% 0% / 0.4);
+  --shadow-lg:
+    0px 2px 4px 0px hsl(0 0% 0% / 0.4), 0px 4px 6px -1px hsl(0 0% 0% / 0.4);
+  --shadow-xl:
+    0px 2px 4px 0px hsl(0 0% 0% / 0.4), 0px 8px 10px -1px hsl(0 0% 0% / 0.4);
+  --shadow-2xl: 0px 2px 4px 0px hsl(0 0% 0% / 1);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="doom-64-dark"] {
+  --background: oklch(0.22 0 0);
+  --foreground: oklch(0.91 0 0);
+  --card: oklch(0.29 0 0);
+  --card-foreground: oklch(0.91 0 0);
+  --popover: oklch(0.29 0 0);
+  --popover-foreground: oklch(0.91 0 0);
+  --primary: oklch(0.61 0.21 27.03);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.64 0.15 133.01);
+  --secondary-foreground: oklch(0 0 0);
+  --muted: oklch(0.26 0 0);
+  --muted-foreground: oklch(0.71 0 0);
+  --accent: oklch(0.75 0.12 244.75);
+  --accent-foreground: oklch(0 0 0);
+  --destructive: oklch(0.78 0.17 68.09);
+  --destructive-foreground: oklch(0 0 0);
+  --border: oklch(0.41 0 0);
+  --input: oklch(0.41 0 0);
+  --ring: oklch(0.61 0.21 27.03);
+  --chart-1: oklch(0.61 0.21 27.03);
+  --chart-2: oklch(0.64 0.15 133.01);
+  --chart-3: oklch(0.75 0.12 244.75);
+  --chart-4: oklch(0.78 0.17 68.09);
+  --chart-5: oklch(0.65 0.03 40.8);
+  --radius: 0px;
+  --sidebar: oklch(0.19 0 0);
+  --sidebar-foreground: oklch(0.91 0 0);
+  --sidebar-primary: oklch(0.61 0.21 27.03);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.75 0.12 244.75);
+  --sidebar-accent-foreground: oklch(0 0 0);
+  --sidebar-border: oklch(0.41 0 0);
+  --sidebar-ring: oklch(0.61 0.21 27.03);
+  --font-sans: "Oxanium", sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: "Source Code Pro", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.6;
+  --shadow-blur: 5px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 2px 5px 0px hsl(0 0% 0% / 0.3);
+  --shadow-xs: 0px 2px 5px 0px hsl(0 0% 0% / 0.3);
+  --shadow-sm:
+    0px 2px 5px 0px hsl(0 0% 0% / 0.6), 0px 1px 2px -1px hsl(0 0% 0% / 0.6);
+  --shadow:
+    0px 2px 5px 0px hsl(0 0% 0% / 0.6), 0px 1px 2px -1px hsl(0 0% 0% / 0.6);
+  --shadow-md:
+    0px 2px 5px 0px hsl(0 0% 0% / 0.6), 0px 2px 4px -1px hsl(0 0% 0% / 0.6);
+  --shadow-lg:
+    0px 2px 5px 0px hsl(0 0% 0% / 0.6), 0px 4px 6px -1px hsl(0 0% 0% / 0.6);
+  --shadow-xl:
+    0px 2px 5px 0px hsl(0 0% 0% / 0.6), 0px 8px 10px -1px hsl(0 0% 0% / 0.6);
+  --shadow-2xl: 0px 2px 5px 0px hsl(0 0% 0% / 1.5);
+}
diff --git a/web/src/styles/themes/elegant-luxury.css b/web/src/styles/themes/elegant-luxury.css
new file mode 100644
index 0000000..b5843de
--- /dev/null
+++ b/web/src/styles/themes/elegant-luxury.css
@@ -0,0 +1,134 @@
+/* Elegant Luxury Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="elegant-luxury-light"] {
+  --background: oklch(0.98 0 56.38);
+  --foreground: oklch(0.22 0 0);
+  --card: oklch(0.98 0 56.38);
+  --card-foreground: oklch(0.22 0 0);
+  --popover: oklch(0.98 0 56.38);
+  --popover-foreground: oklch(0.22 0 0);
+  --primary: oklch(0.47 0.15 24.94);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.96 0.04 89.09);
+  --secondary-foreground: oklch(0.48 0.1 75.12);
+  --muted: oklch(0.94 0.01 53.44);
+  --muted-foreground: oklch(0.44 0.01 73.64);
+  --accent: oklch(0.96 0.06 95.62);
+  --accent-foreground: oklch(0.4 0.13 25.72);
+  --destructive: oklch(0.44 0.16 26.9);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.94 0.03 80.99);
+  --input: oklch(0.94 0.03 80.99);
+  --ring: oklch(0.47 0.15 24.94);
+  --chart-1: oklch(0.51 0.19 27.52);
+  --chart-2: oklch(0.47 0.15 24.94);
+  --chart-3: oklch(0.4 0.13 25.72);
+  --chart-4: oklch(0.56 0.15 49);
+  --chart-5: oklch(0.47 0.12 46.2);
+  --radius: 0.375rem;
+  --sidebar: oklch(0.94 0.01 53.44);
+  --sidebar-foreground: oklch(0.22 0 0);
+  --sidebar-primary: oklch(0.47 0.15 24.94);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.96 0.06 95.62);
+  --sidebar-accent-foreground: oklch(0.4 0.13 25.72);
+  --sidebar-border: oklch(0.94 0.03 80.99);
+  --sidebar-ring: oklch(0.47 0.15 24.94);
+  --font-sans: Poppins, sans-serif;
+  --font-serif: Libre Baskerville, serif;
+  --font-mono: IBM Plex Mono, monospace;
+  --shadow-color: hsl(0 63% 18%);
+  --shadow-opacity: 0.12;
+  --shadow-blur: 16px;
+  --shadow-spread: -2px;
+  --shadow-offset-x: 1px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 1px 1px 16px -2px hsl(0 63% 18% / 0.06);
+  --shadow-xs: 1px 1px 16px -2px hsl(0 63% 18% / 0.06);
+  --shadow-sm:
+    1px 1px 16px -2px hsl(0 63% 18% / 0.12),
+    1px 1px 2px -3px hsl(0 63% 18% / 0.12);
+  --shadow:
+    1px 1px 16px -2px hsl(0 63% 18% / 0.12),
+    1px 1px 2px -3px hsl(0 63% 18% / 0.12);
+  --shadow-md:
+    1px 1px 16px -2px hsl(0 63% 18% / 0.12),
+    1px 2px 4px -3px hsl(0 63% 18% / 0.12);
+  --shadow-lg:
+    1px 1px 16px -2px hsl(0 63% 18% / 0.12),
+    1px 4px 6px -3px hsl(0 63% 18% / 0.12);
+  --shadow-xl:
+    1px 1px 16px -2px hsl(0 63% 18% / 0.12),
+    1px 8px 10px -3px hsl(0 63% 18% / 0.12);
+  --shadow-2xl: 1px 1px 16px -2px hsl(0 63% 18% / 0.3);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="elegant-luxury-dark"] {
+  --background: oklch(0.22 0.01 56.04);
+  --foreground: oklch(0.97 0 106.42);
+  --card: oklch(0.27 0.01 34.3);
+  --card-foreground: oklch(0.97 0 106.42);
+  --popover: oklch(0.27 0.01 34.3);
+  --popover-foreground: oklch(0.97 0 106.42);
+  --primary: oklch(0.51 0.19 27.52);
+  --primary-foreground: oklch(0.98 0 56.38);
+  --secondary: oklch(0.47 0.12 46.2);
+  --secondary-foreground: oklch(0.96 0.06 95.62);
+  --muted: oklch(0.27 0.01 34.3);
+  --muted-foreground: oklch(0.87 0 56.37);
+  --accent: oklch(0.56 0.15 49);
+  --accent-foreground: oklch(0.96 0.06 95.62);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.37 0.01 67.56);
+  --input: oklch(0.37 0.01 67.56);
+  --ring: oklch(0.51 0.19 27.52);
+  --chart-1: oklch(0.71 0.17 22.22);
+  --chart-2: oklch(0.64 0.21 25.33);
+  --chart-3: oklch(0.58 0.22 27.33);
+  --chart-4: oklch(0.84 0.16 84.43);
+  --chart-5: oklch(0.77 0.16 70.08);
+  --radius: 0.375rem;
+  --sidebar: oklch(0.22 0.01 56.04);
+  --sidebar-foreground: oklch(0.97 0 106.42);
+  --sidebar-primary: oklch(0.51 0.19 27.52);
+  --sidebar-primary-foreground: oklch(0.98 0 56.38);
+  --sidebar-accent: oklch(0.56 0.15 49);
+  --sidebar-accent-foreground: oklch(0.96 0.06 95.62);
+  --sidebar-border: oklch(0.37 0.01 67.56);
+  --sidebar-ring: oklch(0.51 0.19 27.52);
+  --font-sans: Poppins, sans-serif;
+  --font-serif: Libre Baskerville, serif;
+  --font-mono: IBM Plex Mono, monospace;
+  --shadow-color: hsl(0 63% 18%);
+  --shadow-opacity: 0.12;
+  --shadow-blur: 16px;
+  --shadow-spread: -2px;
+  --shadow-offset-x: 1px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 1px 1px 16px -2px hsl(0 63% 18% / 0.06);
+  --shadow-xs: 1px 1px 16px -2px hsl(0 63% 18% / 0.06);
+  --shadow-sm:
+    1px 1px 16px -2px hsl(0 63% 18% / 0.12),
+    1px 1px 2px -3px hsl(0 63% 18% / 0.12);
+  --shadow:
+    1px 1px 16px -2px hsl(0 63% 18% / 0.12),
+    1px 1px 2px -3px hsl(0 63% 18% / 0.12);
+  --shadow-md:
+    1px 1px 16px -2px hsl(0 63% 18% / 0.12),
+    1px 2px 4px -3px hsl(0 63% 18% / 0.12);
+  --shadow-lg:
+    1px 1px 16px -2px hsl(0 63% 18% / 0.12),
+    1px 4px 6px -3px hsl(0 63% 18% / 0.12);
+  --shadow-xl:
+    1px 1px 16px -2px hsl(0 63% 18% / 0.12),
+    1px 8px 10px -3px hsl(0 63% 18% / 0.12);
+  --shadow-2xl: 1px 1px 16px -2px hsl(0 63% 18% / 0.3);
+}
diff --git a/web/src/styles/themes/github.css b/web/src/styles/themes/github.css
new file mode 100644
index 0000000..057f8b4
--- /dev/null
+++ b/web/src/styles/themes/github.css
@@ -0,0 +1,114 @@
+/* GitHub Theme - Official GitHub brand colors (Primer Design System)  */
+
+/* Light mode */
+[data-theme="github-light"] {
+  --background: oklch(1 0 0);
+  --foreground: oklch(0.25 0.015 250);
+  --card: oklch(0.98 0.003 250);
+  --card-foreground: oklch(0.25 0.015 250);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.25 0.015 250);
+  --primary: oklch(0.50 0.15 145);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.96 0.005 250);
+  --secondary-foreground: oklch(0.30 0.015 250);
+  --muted: oklch(0.97 0.003 250);
+  --muted-foreground: oklch(0.50 0.015 250);
+  --accent: oklch(0.57 0.18 260);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.55 0.20 25);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.87 0.008 250);
+  --input: oklch(0.97 0.003 250);
+  --ring: oklch(0.57 0.18 260);
+  --chart-1: oklch(0.50 0.15 145);
+  --chart-2: oklch(0.57 0.18 260);
+  --chart-3: oklch(0.65 0.15 320);
+  --chart-4: oklch(0.70 0.12 40);
+  --chart-5: oklch(0.60 0.14 180);
+  --radius: 0.375rem;
+  --sidebar: oklch(0.97 0.003 250);
+  --sidebar-foreground: oklch(0.25 0.015 250);
+  --sidebar-primary: oklch(0.50 0.15 145);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.96 0.005 250);
+  --sidebar-accent-foreground: oklch(0.30 0.015 250);
+  --sidebar-border: oklch(0.87 0.008 250);
+  --sidebar-ring: oklch(0.57 0.18 260);
+  --font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans", Helvetica, Arial, sans-serif;
+  --font-serif: ui-serif, Georgia, serif;
+  --font-mono: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace;
+  --shadow-color: hsl(215 20% 20%);
+  --shadow-opacity: 0.06;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: -0.01em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 0 0px hsl(215 20% 20% / 0.04);
+  --shadow-xs: 0 1px 0 0px hsl(215 20% 20% / 0.04), 0 1px 3px 0px hsl(215 20% 20% / 0.06);
+  --shadow-sm: 0 1px 0 0px hsl(215 20% 20% / 0.04), 0 2px 4px 0px hsl(215 20% 20% / 0.08);
+  --shadow: 0 1px 0 0px hsl(215 20% 20% / 0.04), 0 2px 4px 0px hsl(215 20% 20% / 0.08);
+  --shadow-md: 0 1px 0 0px hsl(215 20% 20% / 0.04), 0 4px 8px 0px hsl(215 20% 20% / 0.1);
+  --shadow-lg: 0 1px 0 0px hsl(215 20% 20% / 0.04), 0 8px 16px 0px hsl(215 20% 20% / 0.1);
+  --shadow-xl: 0 1px 0 0px hsl(215 20% 20% / 0.04), 0 12px 28px 0px hsl(215 20% 20% / 0.12);
+  --shadow-2xl: 0 24px 48px 0px hsl(215 20% 20% / 0.2);
+  --tracking-normal: -0.01em;
+}
+
+/* Dark mode - #0D1117 is GitHub's signature dark background */
+[data-theme="github-dark"] {
+  --background: oklch(0.15 0.012 250);
+  --foreground: oklch(0.85 0.01 250);
+  --card: oklch(0.18 0.012 250);
+  --card-foreground: oklch(0.85 0.01 250);
+  --popover: oklch(0.20 0.012 250);
+  --popover-foreground: oklch(0.85 0.01 250);
+  --primary: oklch(0.55 0.15 145);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.22 0.012 250);
+  --secondary-foreground: oklch(0.80 0.01 250);
+  --muted: oklch(0.20 0.01 250);
+  --muted-foreground: oklch(0.58 0.012 250);
+  --accent: oklch(0.57 0.18 260);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.60 0.18 25);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.30 0.01 250);
+  --input: oklch(0.18 0.012 250);
+  --ring: oklch(0.57 0.18 260);
+  --chart-1: oklch(0.55 0.15 145);
+  --chart-2: oklch(0.57 0.18 260);
+  --chart-3: oklch(0.70 0.14 320);
+  --chart-4: oklch(0.75 0.12 40);
+  --chart-5: oklch(0.65 0.12 180);
+  --radius: 0.375rem;
+  --sidebar: oklch(0.13 0.01 250);
+  --sidebar-foreground: oklch(0.85 0.01 250);
+  --sidebar-primary: oklch(0.55 0.15 145);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.22 0.012 250);
+  --sidebar-accent-foreground: oklch(0.80 0.01 250);
+  --sidebar-border: oklch(0.30 0.01 250);
+  --sidebar-ring: oklch(0.57 0.18 260);
+  --font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans", Helvetica, Arial, sans-serif;
+  --font-serif: ui-serif, Georgia, serif;
+  --font-mono: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.4;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: -0.01em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 0 0 1px hsl(0 0% 0% / 0.3);
+  --shadow-xs: 0 0 0 1px hsl(0 0% 0% / 0.3), 0 1px 3px 0px hsl(0 0% 0% / 0.3);
+  --shadow-sm: 0 0 0 1px hsl(0 0% 0% / 0.3), 0 2px 4px 0px hsl(0 0% 0% / 0.4);
+  --shadow: 0 0 0 1px hsl(0 0% 0% / 0.3), 0 2px 4px 0px hsl(0 0% 0% / 0.4);
+  --shadow-md: 0 0 0 1px hsl(0 0% 0% / 0.3), 0 4px 8px 0px hsl(0 0% 0% / 0.4);
+  --shadow-lg: 0 0 0 1px hsl(0 0% 0% / 0.3), 0 8px 16px 0px hsl(0 0% 0% / 0.4);
+  --shadow-xl: 0 0 0 1px hsl(0 0% 0% / 0.3), 0 12px 28px 0px hsl(0 0% 0% / 0.5);
+  --shadow-2xl: 0 24px 48px 0px hsl(0 0% 0% / 0.6);
+}
diff --git a/web/src/styles/themes/graphite.css b/web/src/styles/themes/graphite.css
new file mode 100644
index 0000000..7598f4a
--- /dev/null
+++ b/web/src/styles/themes/graphite.css
@@ -0,0 +1,124 @@
+/* Graphite Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="graphite-light"] {
+  --background: oklch(0.96 0 0);
+  --foreground: oklch(0.32 0 0);
+  --card: oklch(0.97 0 0);
+  --card-foreground: oklch(0.32 0 0);
+  --popover: oklch(0.97 0 0);
+  --popover-foreground: oklch(0.32 0 0);
+  --primary: oklch(0.49 0 0);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.91 0 0);
+  --secondary-foreground: oklch(0.32 0 0);
+  --muted: oklch(0.89 0 0);
+  --muted-foreground: oklch(0.51 0 0);
+  --accent: oklch(0.81 0 0);
+  --accent-foreground: oklch(0.32 0 0);
+  --destructive: oklch(0.56 0.19 25.86);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.86 0 0);
+  --input: oklch(0.91 0 0);
+  --ring: oklch(0.49 0 0);
+  --chart-1: oklch(0.49 0 0);
+  --chart-2: oklch(0.49 0.04 196.03);
+  --chart-3: oklch(0.65 0 0);
+  --chart-4: oklch(0.73 0 0);
+  --chart-5: oklch(0.81 0 0);
+  --radius: 0.35rem;
+  --sidebar: oklch(0.94 0 0);
+  --sidebar-foreground: oklch(0.32 0 0);
+  --sidebar-primary: oklch(0.49 0 0);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.81 0 0);
+  --sidebar-accent-foreground: oklch(0.32 0 0);
+  --sidebar-border: oklch(0.86 0 0);
+  --sidebar-ring: oklch(0.49 0 0);
+  --font-sans: Montserrat, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(0 0% 20% / 0.1);
+  --shadow-opacity: 0.15;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 2px 0px 0px hsl(0 0% 20% / 0.07);
+  --shadow-xs: 0px 2px 0px 0px hsl(0 0% 20% / 0.07);
+  --shadow-sm:
+    0px 2px 0px 0px hsl(0 0% 20% / 0.15), 0px 1px 2px -1px hsl(0 0% 20% / 0.15);
+  --shadow:
+    0px 2px 0px 0px hsl(0 0% 20% / 0.15), 0px 1px 2px -1px hsl(0 0% 20% / 0.15);
+  --shadow-md:
+    0px 2px 0px 0px hsl(0 0% 20% / 0.15), 0px 2px 4px -1px hsl(0 0% 20% / 0.15);
+  --shadow-lg:
+    0px 2px 0px 0px hsl(0 0% 20% / 0.15), 0px 4px 6px -1px hsl(0 0% 20% / 0.15);
+  --shadow-xl:
+    0px 2px 0px 0px hsl(0 0% 20% / 0.15), 0px 8px 10px -1px hsl(0 0% 20% / 0.15);
+  --shadow-2xl: 0px 2px 0px 0px hsl(0 0% 20% / 0.38);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="graphite-dark"] {
+  --background: oklch(0.22 0 0);
+  --foreground: oklch(0.89 0 0);
+  --card: oklch(0.24 0 0);
+  --card-foreground: oklch(0.89 0 0);
+  --popover: oklch(0.24 0 0);
+  --popover-foreground: oklch(0.89 0 0);
+  --primary: oklch(0.71 0 0);
+  --primary-foreground: oklch(0.22 0 0);
+  --secondary: oklch(0.31 0 0);
+  --secondary-foreground: oklch(0.89 0 0);
+  --muted: oklch(0.29 0 0);
+  --muted-foreground: oklch(0.6 0 0);
+  --accent: oklch(0.37 0 0);
+  --accent-foreground: oklch(0.89 0 0);
+  --destructive: oklch(0.66 0.15 22.17);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.33 0 0);
+  --input: oklch(0.31 0 0);
+  --ring: oklch(0.71 0 0);
+  --chart-1: oklch(0.71 0 0);
+  --chart-2: oklch(0.67 0.03 206.35);
+  --chart-3: oklch(0.55 0 0);
+  --chart-4: oklch(0.46 0 0);
+  --chart-5: oklch(0.37 0 0);
+  --radius: 0.35rem;
+  --sidebar: oklch(0.24 0 0);
+  --sidebar-foreground: oklch(0.89 0 0);
+  --sidebar-primary: oklch(0.71 0 0);
+  --sidebar-primary-foreground: oklch(0.22 0 0);
+  --sidebar-accent: oklch(0.37 0 0);
+  --sidebar-accent-foreground: oklch(0.89 0 0);
+  --sidebar-border: oklch(0.33 0 0);
+  --sidebar-ring: oklch(0.71 0 0);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(0 0% 20% / 0.1);
+  --shadow-opacity: 0.15;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 2px 0px 0px hsl(0 0% 20% / 0.07);
+  --shadow-xs: 0px 2px 0px 0px hsl(0 0% 20% / 0.07);
+  --shadow-sm:
+    0px 2px 0px 0px hsl(0 0% 20% / 0.15), 0px 1px 2px -1px hsl(0 0% 20% / 0.15);
+  --shadow:
+    0px 2px 0px 0px hsl(0 0% 20% / 0.15), 0px 1px 2px -1px hsl(0 0% 20% / 0.15);
+  --shadow-md:
+    0px 2px 0px 0px hsl(0 0% 20% / 0.15), 0px 2px 4px -1px hsl(0 0% 20% / 0.15);
+  --shadow-lg:
+    0px 2px 0px 0px hsl(0 0% 20% / 0.15), 0px 4px 6px -1px hsl(0 0% 20% / 0.15);
+  --shadow-xl:
+    0px 2px 0px 0px hsl(0 0% 20% / 0.15), 0px 8px 10px -1px hsl(0 0% 20% / 0.15);
+  --shadow-2xl: 0px 2px 0px 0px hsl(0 0% 20% / 0.38);
+}
diff --git a/web/src/styles/themes/index.css b/web/src/styles/themes/index.css
new file mode 100644
index 0000000..d9984ab
--- /dev/null
+++ b/web/src/styles/themes/index.css
@@ -0,0 +1,43 @@
+/* Theme imports */
+@import "./modern-minimal.css";
+@import "./t3-chat.css";
+@import "./twitter.css";
+@import "./mocha-mousse.css";
+@import "./bubblegum.css";
+@import "./doom-64.css";
+@import "./catppuccin.css";
+@import "./graphite.css";
+@import "./perpetuity.css";
+@import "./kodama-grove.css";
+@import "./cosmic-night.css";
+@import "./tangerine.css";
+@import "./quantum-rose.css";
+@import "./nature.css";
+@import "./bold-tech.css";
+@import "./elegant-luxury.css";
+@import "./amber-minimal.css";
+@import "./supabase.css";
+@import "./neo-brutalism.css";
+@import "./solar-dusk.css";
+@import "./claymorphism.css";
+@import "./cyberpunk.css";
+@import "./pastel-dreams.css";
+@import "./clean-slate.css";
+@import "./caffeine.css";
+@import "./ocean-breeze.css";
+@import "./retro-arcade.css";
+@import "./midnight-bloom.css";
+@import "./candyland.css";
+@import "./northern-lights.css";
+@import "./vintage-paper.css";
+@import "./sunset-horizon.css";
+@import "./starry-night.css";
+@import "./claude.css";
+@import "./vercel.css";
+@import "./mono.css";
+@import "./twitch.css";
+@import "./kick.css";
+@import "./spotify.css";
+@import "./stripe.css";
+@import "./github.css";
+@import "./windows98.css";
diff --git a/web/src/styles/themes/kick.css b/web/src/styles/themes/kick.css
new file mode 100644
index 0000000..be7d25e
--- /dev/null
+++ b/web/src/styles/themes/kick.css
@@ -0,0 +1,126 @@
+/* Kick Theme - Official Kick brand colors */
+
+/* Light mode */
+[data-theme="kick-light"] {
+  --background: oklch(1 0 0);
+  --foreground: oklch(0.12 0 0);
+  --card: oklch(0.98 0 0);
+  --card-foreground: oklch(0.12 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.12 0 0);
+  --primary: oklch(0.75 0.28 128);
+  --primary-foreground: oklch(0.12 0 0);
+  --secondary: oklch(0.95 0.003 0);
+  --secondary-foreground: oklch(0.2 0 0);
+  --muted: oklch(0.96 0 0);
+  --muted-foreground: oklch(0.5 0 0);
+  --accent: oklch(0.8 0.25 128);
+  --accent-foreground: oklch(0.12 0 0);
+  --destructive: oklch(0.55 0.22 25);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.9 0 0);
+  --input: oklch(0.92 0 0);
+  --ring: oklch(0.75 0.28 128);
+  --chart-1: oklch(0.75 0.28 128);
+  --chart-2: oklch(0.65 0.2 180);
+  --chart-3: oklch(0.7 0.18 220);
+  --chart-4: oklch(0.6 0.15 280);
+  --chart-5: oklch(0.55 0.2 340);
+  --radius: 0.625rem;
+  --sidebar: oklch(0.98 0 0);
+  --sidebar-foreground: oklch(0.12 0 0);
+  --sidebar-primary: oklch(0.75 0.28 128);
+  --sidebar-primary-foreground: oklch(0.12 0 0);
+  --sidebar-accent: oklch(0.95 0.003 0);
+  --sidebar-accent-foreground: oklch(0.2 0 0);
+  --sidebar-border: oklch(0.9 0 0);
+  --sidebar-ring: oklch(0.75 0.28 128);
+  --font-sans: Inter, ui-sans-serif, system-ui, sans-serif;
+  --font-serif: ui-serif, Georgia, serif;
+  --font-mono: ui-monospace, SFMono-Regular, "SF Mono", monospace;
+  --shadow-color: hsl(130 80% 15%);
+  --shadow-opacity: 0.08;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 2px;
+  --letter-spacing: -0.01em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 2px 0px hsl(130 80% 15% / 0.04);
+  --shadow-xs: 0 1px 3px 0px hsl(130 80% 15% / 0.06);
+  --shadow-sm:
+    0 2px 4px 0px hsl(130 80% 15% / 0.08),
+    0 1px 2px -1px hsl(130 80% 15% / 0.08);
+  --shadow:
+    0 2px 4px 0px hsl(130 80% 15% / 0.08),
+    0 1px 2px -1px hsl(130 80% 15% / 0.08);
+  --shadow-md:
+    0 4px 6px 0px hsl(130 80% 15% / 0.1), 0 2px 4px -1px hsl(130 80% 15% / 0.1);
+  --shadow-lg:
+    0 8px 10px 0px hsl(130 80% 15% / 0.1), 0 4px 6px -1px hsl(130 80% 15% / 0.1);
+  --shadow-xl:
+    0 16px 20px 0px hsl(130 80% 15% / 0.12),
+    0 8px 10px -1px hsl(130 80% 15% / 0.12);
+  --shadow-2xl: 0 24px 48px 0px hsl(130 80% 15% / 0.2);
+  --tracking-normal: -0.01em;
+}
+
+/* Dark mode - #0D0D0D is Kick's signature black background */
+[data-theme="kick-dark"] {
+  --background: oklch(0.12 0 0);
+  --foreground: oklch(1 0 0);
+  --card: oklch(0.15 0.003 250);
+  --card-foreground: oklch(1 0 0);
+  --popover: oklch(0.17 0.005 250);
+  --popover-foreground: oklch(1 0 0);
+  --primary: oklch(0.85 0.3 128);
+  --primary-foreground: oklch(0.12 0 0);
+  --secondary: oklch(0.2 0.003 250);
+  --secondary-foreground: oklch(0.9 0 0);
+  --muted: oklch(0.18 0.003 250);
+  --muted-foreground: oklch(0.6 0 0);
+  --accent: oklch(0.85 0.3 128);
+  --accent-foreground: oklch(0.12 0 0);
+  --destructive: oklch(0.6 0.2 25);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.26 0.005 250);
+  --input: oklch(0.2 0.003 250);
+  --ring: oklch(0.85 0.3 128);
+  --chart-1: oklch(0.85 0.3 128);
+  --chart-2: oklch(0.7 0.18 180);
+  --chart-3: oklch(0.75 0.16 220);
+  --chart-4: oklch(0.65 0.14 280);
+  --chart-5: oklch(0.6 0.18 340);
+  --radius: 0.625rem;
+  --sidebar: oklch(0.1 0 0);
+  --sidebar-foreground: oklch(1 0 0);
+  --sidebar-primary: oklch(0.85 0.3 128);
+  --sidebar-primary-foreground: oklch(0.12 0 0);
+  --sidebar-accent: oklch(0.2 0.003 250);
+  --sidebar-accent-foreground: oklch(0.9 0 0);
+  --sidebar-border: oklch(0.26 0.005 250);
+  --sidebar-ring: oklch(0.85 0.3 128);
+  --font-sans: Inter, ui-sans-serif, system-ui, sans-serif;
+  --font-serif: ui-serif, Georgia, serif;
+  --font-mono: ui-monospace, SFMono-Regular, "SF Mono", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.5;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 2px;
+  --letter-spacing: -0.01em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 2px 0px hsl(0 0% 0% / 0.25);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.35);
+  --shadow-sm:
+    0 2px 4px 0px hsl(0 0% 0% / 0.5), 0 1px 2px -1px hsl(0 0% 0% / 0.5);
+  --shadow: 0 2px 4px 0px hsl(0 0% 0% / 0.5), 0 1px 2px -1px hsl(0 0% 0% / 0.5);
+  --shadow-md:
+    0 4px 6px 0px hsl(0 0% 0% / 0.5), 0 2px 4px -1px hsl(0 0% 0% / 0.5);
+  --shadow-lg:
+    0 8px 10px 0px hsl(0 0% 0% / 0.5), 0 4px 6px -1px hsl(0 0% 0% / 0.5);
+  --shadow-xl:
+    0 16px 20px 0px hsl(0 0% 0% / 0.5), 0 8px 10px -1px hsl(0 0% 0% / 0.5);
+  --shadow-2xl: 0 24px 48px 0px hsl(0 0% 0% / 0.6);
+}
diff --git a/web/src/styles/themes/kodama-grove.css b/web/src/styles/themes/kodama-grove.css
new file mode 100644
index 0000000..de867ab
--- /dev/null
+++ b/web/src/styles/themes/kodama-grove.css
@@ -0,0 +1,134 @@
+/* Kodama Grove Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="kodama-grove-light"] {
+  --background: oklch(0.88 0.05 91.79);
+  --foreground: oklch(0.43 0.03 59.22);
+  --card: oklch(0.89 0.04 87.57);
+  --card-foreground: oklch(0.43 0.03 59.22);
+  --popover: oklch(0.94 0.03 89.85);
+  --popover-foreground: oklch(0.43 0.03 59.22);
+  --primary: oklch(0.67 0.11 118.91);
+  --primary-foreground: oklch(0.99 0.01 88.64);
+  --secondary: oklch(0.85 0.06 91.15);
+  --secondary-foreground: oklch(0.43 0.03 59.22);
+  --muted: oklch(0.85 0.06 91.15);
+  --muted-foreground: oklch(0.58 0.03 60.93);
+  --accent: oklch(0.84 0.07 90.33);
+  --accent-foreground: oklch(0.43 0.03 59.22);
+  --destructive: oklch(0.71 0.1 29.98);
+  --destructive-foreground: oklch(0.98 0.01 91.48);
+  --border: oklch(0.69 0.04 59.84);
+  --input: oklch(0.84 0.07 90.33);
+  --ring: oklch(0.73 0.06 130.85);
+  --chart-1: oklch(0.73 0.06 130.85);
+  --chart-2: oklch(0.68 0.06 132.45);
+  --chart-3: oklch(0.82 0.03 136.65);
+  --chart-4: oklch(0.59 0.05 137.62);
+  --chart-5: oklch(0.52 0.04 137.19);
+  --radius: 0.425rem;
+  --sidebar: oklch(0.86 0.06 90.52);
+  --sidebar-foreground: oklch(0.43 0.03 59.22);
+  --sidebar-primary: oklch(0.73 0.06 130.85);
+  --sidebar-primary-foreground: oklch(0.99 0.01 88.64);
+  --sidebar-accent: oklch(0.92 0.02 88);
+  --sidebar-accent-foreground: oklch(0.43 0.03 59.22);
+  --sidebar-border: oklch(0.91 0.02 88);
+  --sidebar-ring: oklch(0.73 0.06 130.85);
+  --font-sans: Merriweather, serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(88 22% 35% / 0.15);
+  --shadow-opacity: 0.15;
+  --shadow-blur: 2px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 3px;
+  --shadow-offset-y: 3px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 3px 3px 2px 0px hsl(88 22% 35% / 0.07);
+  --shadow-xs: 3px 3px 2px 0px hsl(88 22% 35% / 0.07);
+  --shadow-sm:
+    3px 3px 2px 0px hsl(88 22% 35% / 0.15),
+    3px 1px 2px -1px hsl(88 22% 35% / 0.15);
+  --shadow:
+    3px 3px 2px 0px hsl(88 22% 35% / 0.15),
+    3px 1px 2px -1px hsl(88 22% 35% / 0.15);
+  --shadow-md:
+    3px 3px 2px 0px hsl(88 22% 35% / 0.15),
+    3px 2px 4px -1px hsl(88 22% 35% / 0.15);
+  --shadow-lg:
+    3px 3px 2px 0px hsl(88 22% 35% / 0.15),
+    3px 4px 6px -1px hsl(88 22% 35% / 0.15);
+  --shadow-xl:
+    3px 3px 2px 0px hsl(88 22% 35% / 0.15),
+    3px 8px 10px -1px hsl(88 22% 35% / 0.15);
+  --shadow-2xl: 3px 3px 2px 0px hsl(88 22% 35% / 0.38);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="kodama-grove-dark"] {
+  --background: oklch(0.33 0.02 88.07);
+  --foreground: oklch(0.92 0.02 82.12);
+  --card: oklch(0.36 0.02 82.33);
+  --card-foreground: oklch(0.92 0.02 82.12);
+  --popover: oklch(0.36 0.02 82.33);
+  --popover-foreground: oklch(0.92 0.02 82.12);
+  --primary: oklch(0.68 0.06 132.45);
+  --primary-foreground: oklch(0.27 0.01 61.02);
+  --secondary: oklch(0.44 0.02 84.55);
+  --secondary-foreground: oklch(0.92 0.02 82.12);
+  --muted: oklch(0.39 0.02 82.71);
+  --muted-foreground: oklch(0.71 0.02 73.62);
+  --accent: oklch(0.65 0.07 90.76);
+  --accent-foreground: oklch(0.27 0.01 61.02);
+  --destructive: oklch(0.63 0.08 31.3);
+  --destructive-foreground: oklch(0.94 0.02 84.59);
+  --border: oklch(0.44 0.02 84.55);
+  --input: oklch(0.44 0.02 84.55);
+  --ring: oklch(0.68 0.06 132.45);
+  --chart-1: oklch(0.68 0.06 132.45);
+  --chart-2: oklch(0.73 0.06 130.85);
+  --chart-3: oklch(0.59 0.05 137.62);
+  --chart-4: oklch(0.65 0.07 90.76);
+  --chart-5: oklch(0.52 0.04 137.19);
+  --radius: 0.425rem;
+  --sidebar: oklch(0.33 0.02 88.07);
+  --sidebar-foreground: oklch(0.92 0.02 82.12);
+  --sidebar-primary: oklch(0.68 0.06 132.45);
+  --sidebar-primary-foreground: oklch(0.27 0.01 61.02);
+  --sidebar-accent: oklch(0.65 0.07 90.76);
+  --sidebar-accent-foreground: oklch(0.27 0.01 61.02);
+  --sidebar-border: oklch(0.44 0.02 84.55);
+  --sidebar-ring: oklch(0.68 0.06 132.45);
+  --font-sans: Merriweather, serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(88 22% 35% / 0.15);
+  --shadow-opacity: 0.15;
+  --shadow-blur: 2px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 3px;
+  --shadow-offset-y: 3px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 3px 3px 2px 0px hsl(88 22% 35% / 0.07);
+  --shadow-xs: 3px 3px 2px 0px hsl(88 22% 35% / 0.07);
+  --shadow-sm:
+    3px 3px 2px 0px hsl(88 22% 35% / 0.15),
+    3px 1px 2px -1px hsl(88 22% 35% / 0.15);
+  --shadow:
+    3px 3px 2px 0px hsl(88 22% 35% / 0.15),
+    3px 1px 2px -1px hsl(88 22% 35% / 0.15);
+  --shadow-md:
+    3px 3px 2px 0px hsl(88 22% 35% / 0.15),
+    3px 2px 4px -1px hsl(88 22% 35% / 0.15);
+  --shadow-lg:
+    3px 3px 2px 0px hsl(88 22% 35% / 0.15),
+    3px 4px 6px -1px hsl(88 22% 35% / 0.15);
+  --shadow-xl:
+    3px 3px 2px 0px hsl(88 22% 35% / 0.15),
+    3px 8px 10px -1px hsl(88 22% 35% / 0.15);
+  --shadow-2xl: 3px 3px 2px 0px hsl(88 22% 35% / 0.38);
+}
diff --git a/web/src/styles/themes/midnight-bloom.css b/web/src/styles/themes/midnight-bloom.css
new file mode 100644
index 0000000..5b19b89
--- /dev/null
+++ b/web/src/styles/themes/midnight-bloom.css
@@ -0,0 +1,124 @@
+/* Midnight Bloom Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="midnight-bloom-light"] {
+  --background: oklch(0.98 0 0);
+  --foreground: oklch(0.32 0 0);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.32 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.32 0 0);
+  --primary: oklch(0.57 0.2 283.08);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.82 0.07 249.35);
+  --secondary-foreground: oklch(0.32 0 0);
+  --muted: oklch(0.82 0.02 91.62);
+  --muted-foreground: oklch(0.54 0 0);
+  --accent: oklch(0.65 0.06 117.43);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.87 0 0);
+  --input: oklch(0.87 0 0);
+  --ring: oklch(0.57 0.2 283.08);
+  --chart-1: oklch(0.57 0.2 283.08);
+  --chart-2: oklch(0.53 0.17 314.65);
+  --chart-3: oklch(0.34 0.18 301.68);
+  --chart-4: oklch(0.67 0.14 261.34);
+  --chart-5: oklch(0.59 0.1 245.74);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.98 0 0);
+  --sidebar-foreground: oklch(0.32 0 0);
+  --sidebar-primary: oklch(0.57 0.2 283.08);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.65 0.06 117.43);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0.87 0 0);
+  --sidebar-ring: oklch(0.57 0.2 283.08);
+  --font-sans: Montserrat, sans-serif;
+  --font-serif: Playfair Display, serif;
+  --font-mono: Source Code Pro, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 10px;
+  --shadow-spread: -2px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 5px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 5px 10px -2px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 5px 10px -2px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 5px 10px -2px hsl(0 0% 0% / 0.1), 0px 1px 2px -3px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 5px 10px -2px hsl(0 0% 0% / 0.1), 0px 1px 2px -3px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 5px 10px -2px hsl(0 0% 0% / 0.1), 0px 2px 4px -3px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 5px 10px -2px hsl(0 0% 0% / 0.1), 0px 4px 6px -3px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 5px 10px -2px hsl(0 0% 0% / 0.1), 0px 8px 10px -3px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 5px 10px -2px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="midnight-bloom-dark"] {
+  --background: oklch(0.23 0.01 264.29);
+  --foreground: oklch(0.92 0 0);
+  --card: oklch(0.32 0.01 223.67);
+  --card-foreground: oklch(0.92 0 0);
+  --popover: oklch(0.32 0.01 223.67);
+  --popover-foreground: oklch(0.92 0 0);
+  --primary: oklch(0.57 0.2 283.08);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.34 0.18 301.68);
+  --secondary-foreground: oklch(0.92 0 0);
+  --muted: oklch(0.39 0 0);
+  --muted-foreground: oklch(0.72 0 0);
+  --accent: oklch(0.67 0.14 261.34);
+  --accent-foreground: oklch(0.92 0 0);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.39 0 0);
+  --input: oklch(0.39 0 0);
+  --ring: oklch(0.57 0.2 283.08);
+  --chart-1: oklch(0.57 0.2 283.08);
+  --chart-2: oklch(0.53 0.17 314.65);
+  --chart-3: oklch(0.34 0.18 301.68);
+  --chart-4: oklch(0.67 0.14 261.34);
+  --chart-5: oklch(0.59 0.1 245.74);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.23 0.01 264.29);
+  --sidebar-foreground: oklch(0.92 0 0);
+  --sidebar-primary: oklch(0.57 0.2 283.08);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.67 0.14 261.34);
+  --sidebar-accent-foreground: oklch(0.92 0 0);
+  --sidebar-border: oklch(0.39 0 0);
+  --sidebar-ring: oklch(0.57 0.2 283.08);
+  --font-sans: Montserrat, sans-serif;
+  --font-serif: Playfair Display, serif;
+  --font-mono: Source Code Pro, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 10px;
+  --shadow-spread: -2px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 5px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 5px 10px -2px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 5px 10px -2px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 5px 10px -2px hsl(0 0% 0% / 0.1), 0px 1px 2px -3px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 5px 10px -2px hsl(0 0% 0% / 0.1), 0px 1px 2px -3px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 5px 10px -2px hsl(0 0% 0% / 0.1), 0px 2px 4px -3px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 5px 10px -2px hsl(0 0% 0% / 0.1), 0px 4px 6px -3px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 5px 10px -2px hsl(0 0% 0% / 0.1), 0px 8px 10px -3px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 5px 10px -2px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/mocha-mousse.css b/web/src/styles/themes/mocha-mousse.css
new file mode 100644
index 0000000..9a8e359
--- /dev/null
+++ b/web/src/styles/themes/mocha-mousse.css
@@ -0,0 +1,134 @@
+/* Mocha Mousse Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="mocha-mousse-light"] {
+  --background: oklch(0.95 0.01 102.46);
+  --foreground: oklch(0.41 0.03 40.36);
+  --card: oklch(0.95 0.01 102.46);
+  --card-foreground: oklch(0.41 0.03 40.36);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.41 0.03 40.36);
+  --primary: oklch(0.61 0.06 44.36);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.75 0.04 80.55);
+  --secondary-foreground: oklch(1 0 0);
+  --muted: oklch(0.85 0.04 49.09);
+  --muted-foreground: oklch(0.54 0.05 37.21);
+  --accent: oklch(0.85 0.04 49.09);
+  --accent-foreground: oklch(0.41 0.03 40.36);
+  --destructive: oklch(0.22 0.01 52.96);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.75 0.04 80.55);
+  --input: oklch(0.75 0.04 80.55);
+  --ring: oklch(0.61 0.06 44.36);
+  --chart-1: oklch(0.61 0.06 44.36);
+  --chart-2: oklch(0.54 0.05 37.21);
+  --chart-3: oklch(0.73 0.05 52.33);
+  --chart-4: oklch(0.75 0.04 80.55);
+  --chart-5: oklch(0.64 0.04 52.39);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.89 0.03 49.57);
+  --sidebar-foreground: oklch(0.41 0.03 40.36);
+  --sidebar-primary: oklch(0.61 0.06 44.36);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.73 0.05 52.33);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0.64 0.04 52.39);
+  --sidebar-ring: oklch(0.61 0.06 44.36);
+  --font-sans: DM Sans, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: Menlo, monospace;
+  --shadow-color: hsl(20 18% 51% / 0.4);
+  --shadow-opacity: 0.11;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 2px 0px 0px hsl(20 18% 51% / 0.06);
+  --shadow-xs: 2px 2px 0px 0px hsl(20 18% 51% / 0.06);
+  --shadow-sm:
+    2px 2px 0px 0px hsl(20 18% 51% / 0.11),
+    2px 1px 2px -1px hsl(20 18% 51% / 0.11);
+  --shadow:
+    2px 2px 0px 0px hsl(20 18% 51% / 0.11),
+    2px 1px 2px -1px hsl(20 18% 51% / 0.11);
+  --shadow-md:
+    2px 2px 0px 0px hsl(20 18% 51% / 0.11),
+    2px 2px 4px -1px hsl(20 18% 51% / 0.11);
+  --shadow-lg:
+    2px 2px 0px 0px hsl(20 18% 51% / 0.11),
+    2px 4px 6px -1px hsl(20 18% 51% / 0.11);
+  --shadow-xl:
+    2px 2px 0px 0px hsl(20 18% 51% / 0.11),
+    2px 8px 10px -1px hsl(20 18% 51% / 0.11);
+  --shadow-2xl: 2px 2px 0px 0px hsl(20 18% 51% / 0.28);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="mocha-mousse-dark"] {
+  --background: oklch(0.27 0.01 48.18);
+  --foreground: oklch(0.95 0.01 102.46);
+  --card: oklch(0.33 0.02 50.89);
+  --card-foreground: oklch(0.95 0.01 102.46);
+  --popover: oklch(0.33 0.02 50.89);
+  --popover-foreground: oklch(0.95 0.01 102.46);
+  --primary: oklch(0.73 0.05 52.33);
+  --primary-foreground: oklch(0.27 0.01 48.18);
+  --secondary: oklch(0.54 0.05 37.21);
+  --secondary-foreground: oklch(0.95 0.01 102.46);
+  --muted: oklch(0.41 0.03 40.36);
+  --muted-foreground: oklch(0.76 0.04 50.86);
+  --accent: oklch(0.75 0.04 80.55);
+  --accent-foreground: oklch(0.27 0.01 48.18);
+  --destructive: oklch(0.69 0.14 21.46);
+  --destructive-foreground: oklch(0.27 0.01 48.18);
+  --border: oklch(0.41 0.03 40.36);
+  --input: oklch(0.41 0.03 40.36);
+  --ring: oklch(0.73 0.05 52.33);
+  --chart-1: oklch(0.73 0.05 52.33);
+  --chart-2: oklch(0.75 0.04 80.55);
+  --chart-3: oklch(0.61 0.06 44.36);
+  --chart-4: oklch(0.54 0.05 37.21);
+  --chart-5: oklch(0.64 0.04 52.39);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.22 0.01 52.96);
+  --sidebar-foreground: oklch(0.95 0.01 102.46);
+  --sidebar-primary: oklch(0.73 0.05 52.33);
+  --sidebar-primary-foreground: oklch(0.22 0.01 52.96);
+  --sidebar-accent: oklch(0.75 0.04 80.55);
+  --sidebar-accent-foreground: oklch(0.22 0.01 52.96);
+  --sidebar-border: oklch(0.41 0.03 40.36);
+  --sidebar-ring: oklch(0.73 0.05 52.33);
+  --font-sans: DM Sans, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: Menlo, monospace;
+  --shadow-color: hsl(20 18% 30% / 0.5);
+  --shadow-opacity: 0.11;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 2px 0px 0px hsl(20 18% 30% / 0.06);
+  --shadow-xs: 2px 2px 0px 0px hsl(20 18% 30% / 0.06);
+  --shadow-sm:
+    2px 2px 0px 0px hsl(20 18% 30% / 0.11),
+    2px 1px 2px -1px hsl(20 18% 30% / 0.11);
+  --shadow:
+    2px 2px 0px 0px hsl(20 18% 30% / 0.11),
+    2px 1px 2px -1px hsl(20 18% 30% / 0.11);
+  --shadow-md:
+    2px 2px 0px 0px hsl(20 18% 30% / 0.11),
+    2px 2px 4px -1px hsl(20 18% 30% / 0.11);
+  --shadow-lg:
+    2px 2px 0px 0px hsl(20 18% 30% / 0.11),
+    2px 4px 6px -1px hsl(20 18% 30% / 0.11);
+  --shadow-xl:
+    2px 2px 0px 0px hsl(20 18% 30% / 0.11),
+    2px 8px 10px -1px hsl(20 18% 30% / 0.11);
+  --shadow-2xl: 2px 2px 0px 0px hsl(20 18% 30% / 0.28);
+}
diff --git a/web/src/styles/themes/modern-minimal.css b/web/src/styles/themes/modern-minimal.css
new file mode 100644
index 0000000..5dca3a8
--- /dev/null
+++ b/web/src/styles/themes/modern-minimal.css
@@ -0,0 +1,122 @@
+/* Modern Minimal Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="modern-minimal-light"] {
+  --background: oklch(1 0 0);
+  --foreground: oklch(0.32 0 0);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.32 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.32 0 0);
+  --primary: oklch(0.62 0.19 259.81);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.97 0 264.54);
+  --secondary-foreground: oklch(0.45 0.03 256.8);
+  --muted: oklch(0.98 0 247.84);
+  --muted-foreground: oklch(0.55 0.02 264.36);
+  --accent: oklch(0.95 0.03 236.82);
+  --accent-foreground: oklch(0.38 0.14 265.52);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.93 0.01 264.53);
+  --input: oklch(0.93 0.01 264.53);
+  --ring: oklch(0.62 0.19 259.81);
+  --chart-1: oklch(0.62 0.19 259.81);
+  --chart-2: oklch(0.55 0.22 262.88);
+  --chart-3: oklch(0.49 0.22 264.38);
+  --chart-4: oklch(0.42 0.18 265.64);
+  --chart-5: oklch(0.38 0.14 265.52);
+  --radius: 0.375rem;
+  --sidebar: oklch(0.98 0 247.84);
+  --sidebar-foreground: oklch(0.32 0 0);
+  --sidebar-primary: oklch(0.62 0.19 259.81);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.95 0.03 236.82);
+  --sidebar-accent-foreground: oklch(0.38 0.14 265.52);
+  --sidebar-border: oklch(0.93 0.01 264.53);
+  --sidebar-ring: oklch(0.62 0.19 259.81);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="modern-minimal-dark"] {
+  --background: oklch(0.2 0 0);
+  --foreground: oklch(0.92 0 0);
+  --card: oklch(0.27 0 0);
+  --card-foreground: oklch(0.92 0 0);
+  --popover: oklch(0.27 0 0);
+  --popover-foreground: oklch(0.92 0 0);
+  --primary: oklch(0.62 0.19 259.81);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.27 0 0);
+  --secondary-foreground: oklch(0.92 0 0);
+  --muted: oklch(0.27 0 0);
+  --muted-foreground: oklch(0.72 0 0);
+  --accent: oklch(0.38 0.14 265.52);
+  --accent-foreground: oklch(0.88 0.06 254.13);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.37 0 0);
+  --input: oklch(0.37 0 0);
+  --ring: oklch(0.62 0.19 259.81);
+  --chart-1: oklch(0.71 0.14 254.62);
+  --chart-2: oklch(0.62 0.19 259.81);
+  --chart-3: oklch(0.55 0.22 262.88);
+  --chart-4: oklch(0.49 0.22 264.38);
+  --chart-5: oklch(0.42 0.18 265.64);
+  --radius: 0.375rem;
+  --sidebar: oklch(0.2 0 0);
+  --sidebar-foreground: oklch(0.92 0 0);
+  --sidebar-primary: oklch(0.62 0.19 259.81);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.38 0.14 265.52);
+  --sidebar-accent-foreground: oklch(0.88 0.06 254.13);
+  --sidebar-border: oklch(0.37 0 0);
+  --sidebar-ring: oklch(0.62 0.19 259.81);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/mono.css b/web/src/styles/themes/mono.css
new file mode 100644
index 0000000..01f0428
--- /dev/null
+++ b/web/src/styles/themes/mono.css
@@ -0,0 +1,122 @@
+/* Mono Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="mono-light"] {
+  --background: oklch(1 0 0);
+  --foreground: oklch(0.14 0 0);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.14 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.14 0 0);
+  --primary: oklch(0.56 0 0);
+  --primary-foreground: oklch(0.99 0 0);
+  --secondary: oklch(0.97 0 0);
+  --secondary-foreground: oklch(0.2 0 0);
+  --muted: oklch(0.97 0 0);
+  --muted-foreground: oklch(0.55 0 0);
+  --accent: oklch(0.97 0 0);
+  --accent-foreground: oklch(0.2 0 0);
+  --destructive: oklch(0.58 0.24 28.48);
+  --destructive-foreground: oklch(0.97 0 0);
+  --border: oklch(0.92 0 0);
+  --input: oklch(0.92 0 0);
+  --ring: oklch(0.71 0 0);
+  --chart-1: oklch(0.56 0 0);
+  --chart-2: oklch(0.56 0 0);
+  --chart-3: oklch(0.56 0 0);
+  --chart-4: oklch(0.56 0 0);
+  --chart-5: oklch(0.56 0 0);
+  --radius: 0rem;
+  --sidebar: oklch(0.99 0 0);
+  --sidebar-foreground: oklch(0.14 0 0);
+  --sidebar-primary: oklch(0.2 0 0);
+  --sidebar-primary-foreground: oklch(0.99 0 0);
+  --sidebar-accent: oklch(0.97 0 0);
+  --sidebar-accent-foreground: oklch(0.2 0 0);
+  --sidebar-border: oklch(0.92 0 0);
+  --sidebar-ring: oklch(0.71 0 0);
+  --font-sans: Geist Mono, monospace;
+  --font-serif: Geist Mono, monospace;
+  --font-mono: Geist Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 1px 0px 0px hsl(0 0% 0% / 0);
+  --shadow-xs: 0px 1px 0px 0px hsl(0 0% 0% / 0);
+  --shadow-sm:
+    0px 1px 0px 0px hsl(0 0% 0% / 0), 0px 1px 2px -1px hsl(0 0% 0% / 0);
+  --shadow: 0px 1px 0px 0px hsl(0 0% 0% / 0), 0px 1px 2px -1px hsl(0 0% 0% / 0);
+  --shadow-md:
+    0px 1px 0px 0px hsl(0 0% 0% / 0), 0px 2px 4px -1px hsl(0 0% 0% / 0);
+  --shadow-lg:
+    0px 1px 0px 0px hsl(0 0% 0% / 0), 0px 4px 6px -1px hsl(0 0% 0% / 0);
+  --shadow-xl:
+    0px 1px 0px 0px hsl(0 0% 0% / 0), 0px 8px 10px -1px hsl(0 0% 0% / 0);
+  --shadow-2xl: 0px 1px 0px 0px hsl(0 0% 0% / 0);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="mono-dark"] {
+  --background: oklch(0.14 0 0);
+  --foreground: oklch(0.99 0 0);
+  --card: oklch(0.21 0 0);
+  --card-foreground: oklch(0.99 0 0);
+  --popover: oklch(0.27 0 0);
+  --popover-foreground: oklch(0.99 0 0);
+  --primary: oklch(0.56 0 0);
+  --primary-foreground: oklch(0.99 0 0);
+  --secondary: oklch(0.27 0 0);
+  --secondary-foreground: oklch(0.99 0 0);
+  --muted: oklch(0.27 0 0);
+  --muted-foreground: oklch(0.71 0 0);
+  --accent: oklch(0.37 0 0);
+  --accent-foreground: oklch(0.99 0 0);
+  --destructive: oklch(0.7 0.19 22.23);
+  --destructive-foreground: oklch(0.27 0 0);
+  --border: oklch(0.34 0 0);
+  --input: oklch(0.44 0 0);
+  --ring: oklch(0.56 0 0);
+  --chart-1: oklch(0.56 0 0);
+  --chart-2: oklch(0.56 0 0);
+  --chart-3: oklch(0.56 0 0);
+  --chart-4: oklch(0.56 0 0);
+  --chart-5: oklch(0.56 0 0);
+  --radius: 0rem;
+  --sidebar: oklch(0.2 0 0);
+  --sidebar-foreground: oklch(0.99 0 0);
+  --sidebar-primary: oklch(0.99 0 0);
+  --sidebar-primary-foreground: oklch(0.2 0 0);
+  --sidebar-accent: oklch(0.27 0 0);
+  --sidebar-accent-foreground: oklch(0.99 0 0);
+  --sidebar-border: oklch(1 0 0);
+  --sidebar-ring: oklch(0.44 0 0);
+  --font-sans: Geist Mono, monospace;
+  --font-serif: Geist Mono, monospace;
+  --font-mono: Geist Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 1px 0px 0px hsl(0 0% 0% / 0);
+  --shadow-xs: 0px 1px 0px 0px hsl(0 0% 0% / 0);
+  --shadow-sm:
+    0px 1px 0px 0px hsl(0 0% 0% / 0), 0px 1px 2px -1px hsl(0 0% 0% / 0);
+  --shadow: 0px 1px 0px 0px hsl(0 0% 0% / 0), 0px 1px 2px -1px hsl(0 0% 0% / 0);
+  --shadow-md:
+    0px 1px 0px 0px hsl(0 0% 0% / 0), 0px 2px 4px -1px hsl(0 0% 0% / 0);
+  --shadow-lg:
+    0px 1px 0px 0px hsl(0 0% 0% / 0), 0px 4px 6px -1px hsl(0 0% 0% / 0);
+  --shadow-xl:
+    0px 1px 0px 0px hsl(0 0% 0% / 0), 0px 8px 10px -1px hsl(0 0% 0% / 0);
+  --shadow-2xl: 0px 1px 0px 0px hsl(0 0% 0% / 0);
+}
diff --git a/web/src/styles/themes/nature.css b/web/src/styles/themes/nature.css
new file mode 100644
index 0000000..1a099f7
--- /dev/null
+++ b/web/src/styles/themes/nature.css
@@ -0,0 +1,122 @@
+/* Nature Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="nature-light"] {
+  --background: oklch(0.97 0.01 80.72);
+  --foreground: oklch(0.3 0.04 30.2);
+  --card: oklch(0.97 0.01 80.72);
+  --card-foreground: oklch(0.3 0.04 30.2);
+  --popover: oklch(0.97 0.01 80.72);
+  --popover-foreground: oklch(0.3 0.04 30.2);
+  --primary: oklch(0.52 0.13 144.17);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.96 0.02 147.64);
+  --secondary-foreground: oklch(0.43 0.12 144.31);
+  --muted: oklch(0.94 0.01 74.42);
+  --muted-foreground: oklch(0.45 0.05 39.21);
+  --accent: oklch(0.9 0.05 146.04);
+  --accent-foreground: oklch(0.43 0.12 144.31);
+  --destructive: oklch(0.54 0.19 26.72);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.88 0.02 74.64);
+  --input: oklch(0.88 0.02 74.64);
+  --ring: oklch(0.52 0.13 144.17);
+  --chart-1: oklch(0.67 0.16 144.21);
+  --chart-2: oklch(0.58 0.14 144.18);
+  --chart-3: oklch(0.52 0.13 144.17);
+  --chart-4: oklch(0.43 0.12 144.31);
+  --chart-5: oklch(0.22 0.05 145.73);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.94 0.01 74.42);
+  --sidebar-foreground: oklch(0.3 0.04 30.2);
+  --sidebar-primary: oklch(0.52 0.13 144.17);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.9 0.05 146.04);
+  --sidebar-accent-foreground: oklch(0.43 0.12 144.31);
+  --sidebar-border: oklch(0.88 0.02 74.64);
+  --sidebar-ring: oklch(0.52 0.13 144.17);
+  --font-sans: Montserrat, sans-serif;
+  --font-serif: Merriweather, serif;
+  --font-mono: Source Code Pro, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="nature-dark"] {
+  --background: oklch(0.27 0.03 150.77);
+  --foreground: oklch(0.94 0.01 72.66);
+  --card: oklch(0.33 0.03 146.99);
+  --card-foreground: oklch(0.94 0.01 72.66);
+  --popover: oklch(0.33 0.03 146.99);
+  --popover-foreground: oklch(0.94 0.01 72.66);
+  --primary: oklch(0.67 0.16 144.21);
+  --primary-foreground: oklch(0.22 0.05 145.73);
+  --secondary: oklch(0.39 0.03 142.99);
+  --secondary-foreground: oklch(0.9 0.02 142.55);
+  --muted: oklch(0.33 0.03 146.99);
+  --muted-foreground: oklch(0.86 0.02 76.1);
+  --accent: oklch(0.58 0.14 144.18);
+  --accent-foreground: oklch(0.94 0.01 72.66);
+  --destructive: oklch(0.54 0.19 26.72);
+  --destructive-foreground: oklch(0.94 0.01 72.66);
+  --border: oklch(0.39 0.03 142.99);
+  --input: oklch(0.39 0.03 142.99);
+  --ring: oklch(0.67 0.16 144.21);
+  --chart-1: oklch(0.77 0.12 145.3);
+  --chart-2: oklch(0.72 0.14 144.89);
+  --chart-3: oklch(0.67 0.16 144.21);
+  --chart-4: oklch(0.63 0.15 144.2);
+  --chart-5: oklch(0.58 0.14 144.18);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.27 0.03 150.77);
+  --sidebar-foreground: oklch(0.94 0.01 72.66);
+  --sidebar-primary: oklch(0.67 0.16 144.21);
+  --sidebar-primary-foreground: oklch(0.22 0.05 145.73);
+  --sidebar-accent: oklch(0.58 0.14 144.18);
+  --sidebar-accent-foreground: oklch(0.94 0.01 72.66);
+  --sidebar-border: oklch(0.39 0.03 142.99);
+  --sidebar-ring: oklch(0.67 0.16 144.21);
+  --font-sans: Montserrat, sans-serif;
+  --font-serif: Merriweather, serif;
+  --font-mono: Source Code Pro, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/neo-brutalism.css b/web/src/styles/themes/neo-brutalism.css
new file mode 100644
index 0000000..4f82061
--- /dev/null
+++ b/web/src/styles/themes/neo-brutalism.css
@@ -0,0 +1,122 @@
+/* Neo Brutalism Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="neo-brutalism-light"] {
+  --background: oklch(1 0 0);
+  --foreground: oklch(0 0 0);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0 0 0);
+  --primary: oklch(0.65 0.24 26.97);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.97 0.21 109.77);
+  --secondary-foreground: oklch(0 0 0);
+  --muted: oklch(0.96 0 0);
+  --muted-foreground: oklch(0.32 0 0);
+  --accent: oklch(0.56 0.24 260.82);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0 0 0);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0 0 0);
+  --input: oklch(0 0 0);
+  --ring: oklch(0.65 0.24 26.97);
+  --chart-1: oklch(0.65 0.24 26.97);
+  --chart-2: oklch(0.97 0.21 109.77);
+  --chart-3: oklch(0.56 0.24 260.82);
+  --chart-4: oklch(0.73 0.25 142.5);
+  --chart-5: oklch(0.59 0.27 328.36);
+  --radius: 0px;
+  --sidebar: oklch(0.96 0 0);
+  --sidebar-foreground: oklch(0 0 0);
+  --sidebar-primary: oklch(0.65 0.24 26.97);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.56 0.24 260.82);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0 0 0);
+  --sidebar-ring: oklch(0.65 0.24 26.97);
+  --font-sans: DM Sans, sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: Space Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 1;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 4px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 4px 4px 0px 0px hsl(0 0% 0% / 0.5);
+  --shadow-xs: 4px 4px 0px 0px hsl(0 0% 0% / 0.5);
+  --shadow-sm:
+    4px 4px 0px 0px hsl(0 0% 0% / 1), 4px 1px 2px -1px hsl(0 0% 0% / 1);
+  --shadow: 4px 4px 0px 0px hsl(0 0% 0% / 1), 4px 1px 2px -1px hsl(0 0% 0% / 1);
+  --shadow-md:
+    4px 4px 0px 0px hsl(0 0% 0% / 1), 4px 2px 4px -1px hsl(0 0% 0% / 1);
+  --shadow-lg:
+    4px 4px 0px 0px hsl(0 0% 0% / 1), 4px 4px 6px -1px hsl(0 0% 0% / 1);
+  --shadow-xl:
+    4px 4px 0px 0px hsl(0 0% 0% / 1), 4px 8px 10px -1px hsl(0 0% 0% / 1);
+  --shadow-2xl: 4px 4px 0px 0px hsl(0 0% 0% / 2.5);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="neo-brutalism-dark"] {
+  --background: oklch(0 0 0);
+  --foreground: oklch(1 0 0);
+  --card: oklch(0.32 0 0);
+  --card-foreground: oklch(1 0 0);
+  --popover: oklch(0.32 0 0);
+  --popover-foreground: oklch(1 0 0);
+  --primary: oklch(0.7 0.19 23.19);
+  --primary-foreground: oklch(0 0 0);
+  --secondary: oklch(0.97 0.2 109.62);
+  --secondary-foreground: oklch(0 0 0);
+  --muted: oklch(0.32 0 0);
+  --muted-foreground: oklch(0.85 0 0);
+  --accent: oklch(0.68 0.18 252.26);
+  --accent-foreground: oklch(0 0 0);
+  --destructive: oklch(1 0 0);
+  --destructive-foreground: oklch(0 0 0);
+  --border: oklch(1 0 0);
+  --input: oklch(1 0 0);
+  --ring: oklch(0.7 0.19 23.19);
+  --chart-1: oklch(0.7 0.19 23.19);
+  --chart-2: oklch(0.97 0.2 109.62);
+  --chart-3: oklch(0.68 0.18 252.26);
+  --chart-4: oklch(0.74 0.23 142.85);
+  --chart-5: oklch(0.61 0.25 328.07);
+  --radius: 0px;
+  --sidebar: oklch(0 0 0);
+  --sidebar-foreground: oklch(1 0 0);
+  --sidebar-primary: oklch(0.7 0.19 23.19);
+  --sidebar-primary-foreground: oklch(0 0 0);
+  --sidebar-accent: oklch(0.68 0.18 252.26);
+  --sidebar-accent-foreground: oklch(0 0 0);
+  --sidebar-border: oklch(1 0 0);
+  --sidebar-ring: oklch(0.7 0.19 23.19);
+  --font-sans: DM Sans, sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: Space Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 1;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 4px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 4px 4px 0px 0px hsl(0 0% 0% / 0.5);
+  --shadow-xs: 4px 4px 0px 0px hsl(0 0% 0% / 0.5);
+  --shadow-sm:
+    4px 4px 0px 0px hsl(0 0% 0% / 1), 4px 1px 2px -1px hsl(0 0% 0% / 1);
+  --shadow: 4px 4px 0px 0px hsl(0 0% 0% / 1), 4px 1px 2px -1px hsl(0 0% 0% / 1);
+  --shadow-md:
+    4px 4px 0px 0px hsl(0 0% 0% / 1), 4px 2px 4px -1px hsl(0 0% 0% / 1);
+  --shadow-lg:
+    4px 4px 0px 0px hsl(0 0% 0% / 1), 4px 4px 6px -1px hsl(0 0% 0% / 1);
+  --shadow-xl:
+    4px 4px 0px 0px hsl(0 0% 0% / 1), 4px 8px 10px -1px hsl(0 0% 0% / 1);
+  --shadow-2xl: 4px 4px 0px 0px hsl(0 0% 0% / 2.5);
+}
diff --git a/web/src/styles/themes/northern-lights.css b/web/src/styles/themes/northern-lights.css
new file mode 100644
index 0000000..83ecd69
--- /dev/null
+++ b/web/src/styles/themes/northern-lights.css
@@ -0,0 +1,122 @@
+/* Northern Lights Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="northern-lights-light"] {
+  --background: oklch(0.98 0 286.38);
+  --foreground: oklch(0.32 0 0);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.32 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.32 0 0);
+  --primary: oklch(0.65 0.15 150.31);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.67 0.14 261.34);
+  --secondary-foreground: oklch(1 0 0);
+  --muted: oklch(0.88 0.03 98.1);
+  --muted-foreground: oklch(0.54 0 0);
+  --accent: oklch(0.83 0.11 211.96);
+  --accent-foreground: oklch(0.32 0 0);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.87 0 0);
+  --input: oklch(0.87 0 0);
+  --ring: oklch(0.65 0.15 150.31);
+  --chart-1: oklch(0.65 0.15 150.31);
+  --chart-2: oklch(0.67 0.14 261.34);
+  --chart-3: oklch(0.83 0.11 211.96);
+  --chart-4: oklch(0.59 0.1 245.74);
+  --chart-5: oklch(0.59 0.16 148.24);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.98 0 286.38);
+  --sidebar-foreground: oklch(0.32 0 0);
+  --sidebar-primary: oklch(0.65 0.15 150.31);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.83 0.11 211.96);
+  --sidebar-accent-foreground: oklch(0.32 0 0);
+  --sidebar-border: oklch(0.87 0 0);
+  --sidebar-ring: oklch(0.65 0.15 150.31);
+  --font-sans: Plus Jakarta Sans, sans-serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="northern-lights-dark"] {
+  --background: oklch(0.23 0.01 264.29);
+  --foreground: oklch(0.92 0 0);
+  --card: oklch(0.32 0.01 223.67);
+  --card-foreground: oklch(0.92 0 0);
+  --popover: oklch(0.32 0.01 223.67);
+  --popover-foreground: oklch(0.92 0 0);
+  --primary: oklch(0.65 0.15 150.31);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.59 0.1 245.74);
+  --secondary-foreground: oklch(0.92 0 0);
+  --muted: oklch(0.39 0 0);
+  --muted-foreground: oklch(0.72 0 0);
+  --accent: oklch(0.67 0.14 261.34);
+  --accent-foreground: oklch(0.92 0 0);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.39 0 0);
+  --input: oklch(0.39 0 0);
+  --ring: oklch(0.65 0.15 150.31);
+  --chart-1: oklch(0.65 0.15 150.31);
+  --chart-2: oklch(0.59 0.1 245.74);
+  --chart-3: oklch(0.67 0.14 261.34);
+  --chart-4: oklch(0.83 0.11 211.96);
+  --chart-5: oklch(0.59 0.16 148.24);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.23 0.01 264.29);
+  --sidebar-foreground: oklch(0.92 0 0);
+  --sidebar-primary: oklch(0.65 0.15 150.31);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.67 0.14 261.34);
+  --sidebar-accent-foreground: oklch(0.92 0 0);
+  --sidebar-border: oklch(0.39 0 0);
+  --sidebar-ring: oklch(0.65 0.15 150.31);
+  --font-sans: Plus Jakarta Sans, sans-serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/ocean-breeze.css b/web/src/styles/themes/ocean-breeze.css
new file mode 100644
index 0000000..23cef25
--- /dev/null
+++ b/web/src/styles/themes/ocean-breeze.css
@@ -0,0 +1,124 @@
+/* Ocean Breeze Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="ocean-breeze-light"] {
+  --background: oklch(0.98 0.01 244.25);
+  --foreground: oklch(0.37 0.03 259.73);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.37 0.03 259.73);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.37 0.03 259.73);
+  --primary: oklch(0.72 0.19 149.58);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.95 0.03 236.82);
+  --secondary-foreground: oklch(0.45 0.03 256.8);
+  --muted: oklch(0.97 0 264.54);
+  --muted-foreground: oklch(0.55 0.02 264.36);
+  --accent: oklch(0.95 0.05 163.05);
+  --accent-foreground: oklch(0.37 0.03 259.73);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.93 0.01 264.53);
+  --input: oklch(0.93 0.01 264.53);
+  --ring: oklch(0.72 0.19 149.58);
+  --chart-1: oklch(0.72 0.19 149.58);
+  --chart-2: oklch(0.7 0.15 162.48);
+  --chart-3: oklch(0.6 0.13 163.23);
+  --chart-4: oklch(0.51 0.1 165.61);
+  --chart-5: oklch(0.43 0.09 166.91);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.95 0.03 236.82);
+  --sidebar-foreground: oklch(0.37 0.03 259.73);
+  --sidebar-primary: oklch(0.72 0.19 149.58);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.95 0.05 163.05);
+  --sidebar-accent-foreground: oklch(0.37 0.03 259.73);
+  --sidebar-border: oklch(0.93 0.01 264.53);
+  --sidebar-ring: oklch(0.72 0.19 149.58);
+  --font-sans: DM Sans, sans-serif;
+  --font-serif: Lora, serif;
+  --font-mono: IBM Plex Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 8px;
+  --shadow-spread: -1px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 2px 4px -2px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 4px 6px -2px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 8px 10px -2px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 4px 8px -1px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="ocean-breeze-dark"] {
+  --background: oklch(0.21 0.04 265.75);
+  --foreground: oklch(0.87 0.01 258.34);
+  --card: oklch(0.28 0.04 260.03);
+  --card-foreground: oklch(0.87 0.01 258.34);
+  --popover: oklch(0.28 0.04 260.03);
+  --popover-foreground: oklch(0.87 0.01 258.34);
+  --primary: oklch(0.77 0.15 163.22);
+  --primary-foreground: oklch(0.21 0.04 265.75);
+  --secondary: oklch(0.34 0.03 260.91);
+  --secondary-foreground: oklch(0.71 0.01 286.07);
+  --muted: oklch(0.28 0.04 260.03);
+  --muted-foreground: oklch(0.55 0.02 264.36);
+  --accent: oklch(0.37 0.03 259.73);
+  --accent-foreground: oklch(0.71 0.01 286.07);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(0.21 0.04 265.75);
+  --border: oklch(0.45 0.03 256.8);
+  --input: oklch(0.45 0.03 256.8);
+  --ring: oklch(0.77 0.15 163.22);
+  --chart-1: oklch(0.77 0.15 163.22);
+  --chart-2: oklch(0.78 0.13 181.91);
+  --chart-3: oklch(0.72 0.19 149.58);
+  --chart-4: oklch(0.7 0.15 162.48);
+  --chart-5: oklch(0.6 0.13 163.23);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.28 0.04 260.03);
+  --sidebar-foreground: oklch(0.87 0.01 258.34);
+  --sidebar-primary: oklch(0.77 0.15 163.22);
+  --sidebar-primary-foreground: oklch(0.21 0.04 265.75);
+  --sidebar-accent: oklch(0.37 0.03 259.73);
+  --sidebar-accent-foreground: oklch(0.71 0.01 286.07);
+  --sidebar-border: oklch(0.45 0.03 256.8);
+  --sidebar-ring: oklch(0.77 0.15 163.22);
+  --font-sans: DM Sans, sans-serif;
+  --font-serif: Lora, serif;
+  --font-mono: IBM Plex Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 8px;
+  --shadow-spread: -1px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 4px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 4px 8px -1px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 1px 2px -2px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 2px 4px -2px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 4px 6px -2px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 4px 8px -1px hsl(0 0% 0% / 0.1), 0px 8px 10px -2px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 4px 8px -1px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/pastel-dreams.css b/web/src/styles/themes/pastel-dreams.css
new file mode 100644
index 0000000..97b34df
--- /dev/null
+++ b/web/src/styles/themes/pastel-dreams.css
@@ -0,0 +1,124 @@
+/* Pastel Dreams Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="pastel-dreams-light"] {
+  --background: oklch(0.97 0.01 314.78);
+  --foreground: oklch(0.37 0.03 259.73);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.37 0.03 259.73);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.37 0.03 259.73);
+  --primary: oklch(0.71 0.16 293.54);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.91 0.05 306.09);
+  --secondary-foreground: oklch(0.45 0.03 256.8);
+  --muted: oklch(0.95 0.03 307.17);
+  --muted-foreground: oklch(0.55 0.02 264.36);
+  --accent: oklch(0.94 0.03 321.94);
+  --accent-foreground: oklch(0.37 0.03 259.73);
+  --destructive: oklch(0.81 0.1 19.57);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.91 0.05 306.09);
+  --input: oklch(0.91 0.05 306.09);
+  --ring: oklch(0.71 0.16 293.54);
+  --chart-1: oklch(0.71 0.16 293.54);
+  --chart-2: oklch(0.61 0.22 292.72);
+  --chart-3: oklch(0.54 0.25 293.01);
+  --chart-4: oklch(0.49 0.24 292.58);
+  --chart-5: oklch(0.43 0.21 292.76);
+  --radius: 1.5rem;
+  --sidebar: oklch(0.91 0.05 306.09);
+  --sidebar-foreground: oklch(0.37 0.03 259.73);
+  --sidebar-primary: oklch(0.71 0.16 293.54);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.94 0.03 321.94);
+  --sidebar-accent-foreground: oklch(0.37 0.03 259.73);
+  --sidebar-border: oklch(0.91 0.05 306.09);
+  --sidebar-ring: oklch(0.71 0.16 293.54);
+  --font-sans: Open Sans, sans-serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: IBM Plex Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.08;
+  --shadow-blur: 16px;
+  --shadow-spread: -4px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 8px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 8px 16px -4px hsl(0 0% 0% / 0.04);
+  --shadow-xs: 0px 8px 16px -4px hsl(0 0% 0% / 0.04);
+  --shadow-sm:
+    0px 8px 16px -4px hsl(0 0% 0% / 0.08), 0px 1px 2px -5px hsl(0 0% 0% / 0.08);
+  --shadow:
+    0px 8px 16px -4px hsl(0 0% 0% / 0.08), 0px 1px 2px -5px hsl(0 0% 0% / 0.08);
+  --shadow-md:
+    0px 8px 16px -4px hsl(0 0% 0% / 0.08), 0px 2px 4px -5px hsl(0 0% 0% / 0.08);
+  --shadow-lg:
+    0px 8px 16px -4px hsl(0 0% 0% / 0.08), 0px 4px 6px -5px hsl(0 0% 0% / 0.08);
+  --shadow-xl:
+    0px 8px 16px -4px hsl(0 0% 0% / 0.08), 0px 8px 10px -5px hsl(0 0% 0% / 0.08);
+  --shadow-2xl: 0px 8px 16px -4px hsl(0 0% 0% / 0.2);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="pastel-dreams-dark"] {
+  --background: oklch(0.22 0.01 56.04);
+  --foreground: oklch(0.93 0.03 272.79);
+  --card: oklch(0.28 0.03 307.23);
+  --card-foreground: oklch(0.93 0.03 272.79);
+  --popover: oklch(0.28 0.03 307.23);
+  --popover-foreground: oklch(0.93 0.03 272.79);
+  --primary: oklch(0.79 0.12 295.75);
+  --primary-foreground: oklch(0.22 0.01 56.04);
+  --secondary: oklch(0.34 0.04 308.85);
+  --secondary-foreground: oklch(0.87 0.01 258.34);
+  --muted: oklch(0.28 0.03 307.23);
+  --muted-foreground: oklch(0.71 0.02 261.32);
+  --accent: oklch(0.39 0.05 304.64);
+  --accent-foreground: oklch(0.87 0.01 258.34);
+  --destructive: oklch(0.81 0.1 19.57);
+  --destructive-foreground: oklch(0.22 0.01 56.04);
+  --border: oklch(0.34 0.04 308.85);
+  --input: oklch(0.34 0.04 308.85);
+  --ring: oklch(0.79 0.12 295.75);
+  --chart-1: oklch(0.79 0.12 295.75);
+  --chart-2: oklch(0.71 0.16 293.54);
+  --chart-3: oklch(0.61 0.22 292.72);
+  --chart-4: oklch(0.54 0.25 293.01);
+  --chart-5: oklch(0.49 0.24 292.58);
+  --radius: 1.5rem;
+  --sidebar: oklch(0.34 0.04 308.85);
+  --sidebar-foreground: oklch(0.93 0.03 272.79);
+  --sidebar-primary: oklch(0.79 0.12 295.75);
+  --sidebar-primary-foreground: oklch(0.22 0.01 56.04);
+  --sidebar-accent: oklch(0.39 0.05 304.64);
+  --sidebar-accent-foreground: oklch(0.87 0.01 258.34);
+  --sidebar-border: oklch(0.34 0.04 308.85);
+  --sidebar-ring: oklch(0.79 0.12 295.75);
+  --font-sans: Open Sans, sans-serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: IBM Plex Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.08;
+  --shadow-blur: 16px;
+  --shadow-spread: -4px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 8px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 8px 16px -4px hsl(0 0% 0% / 0.04);
+  --shadow-xs: 0px 8px 16px -4px hsl(0 0% 0% / 0.04);
+  --shadow-sm:
+    0px 8px 16px -4px hsl(0 0% 0% / 0.08), 0px 1px 2px -5px hsl(0 0% 0% / 0.08);
+  --shadow:
+    0px 8px 16px -4px hsl(0 0% 0% / 0.08), 0px 1px 2px -5px hsl(0 0% 0% / 0.08);
+  --shadow-md:
+    0px 8px 16px -4px hsl(0 0% 0% / 0.08), 0px 2px 4px -5px hsl(0 0% 0% / 0.08);
+  --shadow-lg:
+    0px 8px 16px -4px hsl(0 0% 0% / 0.08), 0px 4px 6px -5px hsl(0 0% 0% / 0.08);
+  --shadow-xl:
+    0px 8px 16px -4px hsl(0 0% 0% / 0.08), 0px 8px 10px -5px hsl(0 0% 0% / 0.08);
+  --shadow-2xl: 0px 8px 16px -4px hsl(0 0% 0% / 0.2);
+}
diff --git a/web/src/styles/themes/perpetuity.css b/web/src/styles/themes/perpetuity.css
new file mode 100644
index 0000000..b4a7cce
--- /dev/null
+++ b/web/src/styles/themes/perpetuity.css
@@ -0,0 +1,134 @@
+/* Perpetuity Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="perpetuity-light"] {
+  --background: oklch(0.95 0.01 197.01);
+  --foreground: oklch(0.38 0.06 212.66);
+  --card: oklch(0.97 0.01 197.07);
+  --card-foreground: oklch(0.38 0.06 212.66);
+  --popover: oklch(0.97 0.01 197.07);
+  --popover-foreground: oklch(0.38 0.06 212.66);
+  --primary: oklch(0.56 0.09 203.28);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.92 0.02 196.84);
+  --secondary-foreground: oklch(0.38 0.06 212.66);
+  --muted: oklch(0.93 0.01 196.97);
+  --muted-foreground: oklch(0.54 0.06 201.57);
+  --accent: oklch(0.9 0.03 201.89);
+  --accent-foreground: oklch(0.38 0.06 212.66);
+  --destructive: oklch(0.57 0.19 25.54);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.89 0.02 204.41);
+  --input: oklch(0.92 0.02 196.84);
+  --ring: oklch(0.56 0.09 203.28);
+  --chart-1: oklch(0.56 0.09 203.28);
+  --chart-2: oklch(0.64 0.1 201.59);
+  --chart-3: oklch(0.71 0.11 201.25);
+  --chart-4: oklch(0.77 0.1 201.18);
+  --chart-5: oklch(0.83 0.08 200.97);
+  --radius: 0.125rem;
+  --sidebar: oklch(0.93 0.02 205.32);
+  --sidebar-foreground: oklch(0.38 0.06 212.66);
+  --sidebar-primary: oklch(0.56 0.09 203.28);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.9 0.03 201.89);
+  --sidebar-accent-foreground: oklch(0.38 0.06 212.66);
+  --sidebar-border: oklch(0.89 0.02 204.41);
+  --sidebar-ring: oklch(0.56 0.09 203.28);
+  --font-sans: Courier New, monospace;
+  --font-serif: Courier New, monospace;
+  --font-mono: Courier New, monospace;
+  --shadow-color: hsl(185 70% 30% / 0.15);
+  --shadow-opacity: 0.15;
+  --shadow-blur: 2px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 1px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 1px 1px 2px 0px hsl(185 70% 30% / 0.07);
+  --shadow-xs: 1px 1px 2px 0px hsl(185 70% 30% / 0.07);
+  --shadow-sm:
+    1px 1px 2px 0px hsl(185 70% 30% / 0.15),
+    1px 1px 2px -1px hsl(185 70% 30% / 0.15);
+  --shadow:
+    1px 1px 2px 0px hsl(185 70% 30% / 0.15),
+    1px 1px 2px -1px hsl(185 70% 30% / 0.15);
+  --shadow-md:
+    1px 1px 2px 0px hsl(185 70% 30% / 0.15),
+    1px 2px 4px -1px hsl(185 70% 30% / 0.15);
+  --shadow-lg:
+    1px 1px 2px 0px hsl(185 70% 30% / 0.15),
+    1px 4px 6px -1px hsl(185 70% 30% / 0.15);
+  --shadow-xl:
+    1px 1px 2px 0px hsl(185 70% 30% / 0.15),
+    1px 8px 10px -1px hsl(185 70% 30% / 0.15);
+  --shadow-2xl: 1px 1px 2px 0px hsl(185 70% 30% / 0.38);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="perpetuity-dark"] {
+  --background: oklch(0.21 0.02 224.45);
+  --foreground: oklch(0.85 0.13 195.04);
+  --card: oklch(0.23 0.03 216.07);
+  --card-foreground: oklch(0.85 0.13 195.04);
+  --popover: oklch(0.23 0.03 216.07);
+  --popover-foreground: oklch(0.85 0.13 195.04);
+  --primary: oklch(0.85 0.13 195.04);
+  --primary-foreground: oklch(0.21 0.02 224.45);
+  --secondary: oklch(0.38 0.06 216.5);
+  --secondary-foreground: oklch(0.85 0.13 195.04);
+  --muted: oklch(0.29 0.04 218.82);
+  --muted-foreground: oklch(0.66 0.1 195.05);
+  --accent: oklch(0.38 0.06 216.5);
+  --accent-foreground: oklch(0.85 0.13 195.04);
+  --destructive: oklch(0.62 0.21 25.81);
+  --destructive-foreground: oklch(0.96 0 0);
+  --border: oklch(0.38 0.06 216.5);
+  --input: oklch(0.38 0.06 216.5);
+  --ring: oklch(0.85 0.13 195.04);
+  --chart-1: oklch(0.85 0.13 195.04);
+  --chart-2: oklch(0.66 0.1 195.05);
+  --chart-3: oklch(0.58 0.08 195.07);
+  --chart-4: oklch(0.43 0.06 202.62);
+  --chart-5: oklch(0.31 0.05 204.16);
+  --radius: 0.125rem;
+  --sidebar: oklch(0.21 0.02 224.45);
+  --sidebar-foreground: oklch(0.85 0.13 195.04);
+  --sidebar-primary: oklch(0.85 0.13 195.04);
+  --sidebar-primary-foreground: oklch(0.21 0.02 224.45);
+  --sidebar-accent: oklch(0.38 0.06 216.5);
+  --sidebar-accent-foreground: oklch(0.85 0.13 195.04);
+  --sidebar-border: oklch(0.38 0.06 216.5);
+  --sidebar-ring: oklch(0.85 0.13 195.04);
+  --font-sans: Source Code Pro, monospace;
+  --font-serif: Source Code Pro, monospace;
+  --font-mono: Source Code Pro, monospace;
+  --shadow-color: hsl(180 70% 60% / 0.2);
+  --shadow-opacity: 0.2;
+  --shadow-blur: 2px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 1px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 1px 1px 2px 0px hsl(180 70% 60% / 0.1);
+  --shadow-xs: 1px 1px 2px 0px hsl(180 70% 60% / 0.1);
+  --shadow-sm:
+    1px 1px 2px 0px hsl(180 70% 60% / 0.2),
+    1px 1px 2px -1px hsl(180 70% 60% / 0.2);
+  --shadow:
+    1px 1px 2px 0px hsl(180 70% 60% / 0.2),
+    1px 1px 2px -1px hsl(180 70% 60% / 0.2);
+  --shadow-md:
+    1px 1px 2px 0px hsl(180 70% 60% / 0.2),
+    1px 2px 4px -1px hsl(180 70% 60% / 0.2);
+  --shadow-lg:
+    1px 1px 2px 0px hsl(180 70% 60% / 0.2),
+    1px 4px 6px -1px hsl(180 70% 60% / 0.2);
+  --shadow-xl:
+    1px 1px 2px 0px hsl(180 70% 60% / 0.2),
+    1px 8px 10px -1px hsl(180 70% 60% / 0.2);
+  --shadow-2xl: 1px 1px 2px 0px hsl(180 70% 60% / 0.5);
+}
diff --git a/web/src/styles/themes/quantum-rose.css b/web/src/styles/themes/quantum-rose.css
new file mode 100644
index 0000000..520cced
--- /dev/null
+++ b/web/src/styles/themes/quantum-rose.css
@@ -0,0 +1,134 @@
+/* Quantum Rose Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="quantum-rose-light"] {
+  --background: oklch(0.97 0.02 343.93);
+  --foreground: oklch(0.44 0.17 352.38);
+  --card: oklch(0.98 0.01 339.33);
+  --card-foreground: oklch(0.44 0.17 352.38);
+  --popover: oklch(0.98 0.01 339.33);
+  --popover-foreground: oklch(0.44 0.17 352.38);
+  --primary: oklch(0.6 0.24 0.13);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.92 0.07 326.13);
+  --secondary-foreground: oklch(0.44 0.17 352.38);
+  --muted: oklch(0.94 0.04 344.26);
+  --muted-foreground: oklch(0.57 0.17 352.05);
+  --accent: oklch(0.88 0.08 344.88);
+  --accent-foreground: oklch(0.44 0.17 352.38);
+  --destructive: oklch(0.58 0.19 6.34);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.89 0.07 344.39);
+  --input: oklch(0.92 0.07 326.13);
+  --ring: oklch(0.6 0.24 0.13);
+  --chart-1: oklch(0.6 0.24 0.13);
+  --chart-2: oklch(0.6 0.17 345.04);
+  --chart-3: oklch(0.6 0.12 311.8);
+  --chart-4: oklch(0.58 0.12 283.29);
+  --chart-5: oklch(0.65 0.19 267.97);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.96 0.02 345.75);
+  --sidebar-foreground: oklch(0.44 0.17 352.38);
+  --sidebar-primary: oklch(0.6 0.24 0.13);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.88 0.08 344.88);
+  --sidebar-accent-foreground: oklch(0.44 0.17 352.38);
+  --sidebar-border: oklch(0.93 0.04 343.31);
+  --sidebar-ring: oklch(0.6 0.24 0.13);
+  --font-sans: Poppins, sans-serif;
+  --font-serif: Playfair Display, serif;
+  --font-mono: Space Mono, monospace;
+  --shadow-color: hsl(330 70% 30% / 0.12);
+  --shadow-opacity: 0.18;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 3px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 3px 0px 0px hsl(330 70% 30% / 0.09);
+  --shadow-xs: 0px 3px 0px 0px hsl(330 70% 30% / 0.09);
+  --shadow-sm:
+    0px 3px 0px 0px hsl(330 70% 30% / 0.18),
+    0px 1px 2px -1px hsl(330 70% 30% / 0.18);
+  --shadow:
+    0px 3px 0px 0px hsl(330 70% 30% / 0.18),
+    0px 1px 2px -1px hsl(330 70% 30% / 0.18);
+  --shadow-md:
+    0px 3px 0px 0px hsl(330 70% 30% / 0.18),
+    0px 2px 4px -1px hsl(330 70% 30% / 0.18);
+  --shadow-lg:
+    0px 3px 0px 0px hsl(330 70% 30% / 0.18),
+    0px 4px 6px -1px hsl(330 70% 30% / 0.18);
+  --shadow-xl:
+    0px 3px 0px 0px hsl(330 70% 30% / 0.18),
+    0px 8px 10px -1px hsl(330 70% 30% / 0.18);
+  --shadow-2xl: 0px 3px 0px 0px hsl(330 70% 30% / 0.45);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="quantum-rose-dark"] {
+  --background: oklch(0.18 0.05 313.72);
+  --foreground: oklch(0.86 0.13 326.64);
+  --card: oklch(0.24 0.07 313.23);
+  --card-foreground: oklch(0.86 0.13 326.64);
+  --popover: oklch(0.24 0.07 313.23);
+  --popover-foreground: oklch(0.86 0.13 326.64);
+  --primary: oklch(0.75 0.23 332.02);
+  --primary-foreground: oklch(0.16 0.05 327.57);
+  --secondary: oklch(0.32 0.09 319.65);
+  --secondary-foreground: oklch(0.86 0.13 326.64);
+  --muted: oklch(0.27 0.08 312.35);
+  --muted-foreground: oklch(0.71 0.16 327.11);
+  --accent: oklch(0.36 0.12 325.77);
+  --accent-foreground: oklch(0.86 0.13 326.64);
+  --destructive: oklch(0.65 0.24 7.17);
+  --destructive-foreground: oklch(0.98 0 0);
+  --border: oklch(0.33 0.12 313.54);
+  --input: oklch(0.32 0.09 319.65);
+  --ring: oklch(0.75 0.23 332.02);
+  --chart-1: oklch(0.75 0.23 332.02);
+  --chart-2: oklch(0.65 0.22 317.63);
+  --chart-3: oklch(0.62 0.22 292.77);
+  --chart-4: oklch(0.61 0.16 278.72);
+  --chart-5: oklch(0.62 0.2 268.05);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.19 0.05 311.4);
+  --sidebar-foreground: oklch(0.86 0.13 326.64);
+  --sidebar-primary: oklch(0.75 0.23 332.02);
+  --sidebar-primary-foreground: oklch(0.16 0.05 327.57);
+  --sidebar-accent: oklch(0.36 0.12 325.77);
+  --sidebar-accent-foreground: oklch(0.86 0.13 326.64);
+  --sidebar-border: oklch(0.33 0.12 313.54);
+  --sidebar-ring: oklch(0.75 0.23 332.02);
+  --font-sans: Quicksand, sans-serif;
+  --font-serif: Playfair Display, serif;
+  --font-mono: Space Mono, monospace;
+  --shadow-color: hsl(300 80% 50% / 0.25);
+  --shadow-opacity: 0.18;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 3px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 3px 0px 0px hsl(300 80% 50% / 0.09);
+  --shadow-xs: 0px 3px 0px 0px hsl(300 80% 50% / 0.09);
+  --shadow-sm:
+    0px 3px 0px 0px hsl(300 80% 50% / 0.18),
+    0px 1px 2px -1px hsl(300 80% 50% / 0.18);
+  --shadow:
+    0px 3px 0px 0px hsl(300 80% 50% / 0.18),
+    0px 1px 2px -1px hsl(300 80% 50% / 0.18);
+  --shadow-md:
+    0px 3px 0px 0px hsl(300 80% 50% / 0.18),
+    0px 2px 4px -1px hsl(300 80% 50% / 0.18);
+  --shadow-lg:
+    0px 3px 0px 0px hsl(300 80% 50% / 0.18),
+    0px 4px 6px -1px hsl(300 80% 50% / 0.18);
+  --shadow-xl:
+    0px 3px 0px 0px hsl(300 80% 50% / 0.18),
+    0px 8px 10px -1px hsl(300 80% 50% / 0.18);
+  --shadow-2xl: 0px 3px 0px 0px hsl(300 80% 50% / 0.45);
+}
diff --git a/web/src/styles/themes/retro-arcade.css b/web/src/styles/themes/retro-arcade.css
new file mode 100644
index 0000000..e6c16f9
--- /dev/null
+++ b/web/src/styles/themes/retro-arcade.css
@@ -0,0 +1,134 @@
+/* Retro Arcade Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="retro-arcade-light"] {
+  --background: oklch(0.97 0.03 90.1);
+  --foreground: oklch(0.31 0.05 219.65);
+  --card: oklch(0.93 0.03 92.4);
+  --card-foreground: oklch(0.31 0.05 219.65);
+  --popover: oklch(0.93 0.03 92.4);
+  --popover-foreground: oklch(0.31 0.05 219.65);
+  --primary: oklch(0.59 0.2 355.89);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.64 0.1 187.38);
+  --secondary-foreground: oklch(1 0 0);
+  --muted: oklch(0.7 0.02 196.79);
+  --muted-foreground: oklch(0.31 0.05 219.65);
+  --accent: oklch(0.58 0.17 39.5);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.59 0.21 27.12);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.65 0.02 205.26);
+  --input: oklch(0.65 0.02 205.26);
+  --ring: oklch(0.59 0.2 355.89);
+  --chart-1: oklch(0.61 0.14 244.93);
+  --chart-2: oklch(0.64 0.1 187.38);
+  --chart-3: oklch(0.59 0.2 355.89);
+  --chart-4: oklch(0.58 0.17 39.5);
+  --chart-5: oklch(0.59 0.21 27.12);
+  --radius: 0.25rem;
+  --sidebar: oklch(0.97 0.03 90.1);
+  --sidebar-foreground: oklch(0.31 0.05 219.65);
+  --sidebar-primary: oklch(0.59 0.2 355.89);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.64 0.1 187.38);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0.65 0.02 205.26);
+  --sidebar-ring: oklch(0.59 0.2 355.89);
+  --font-sans: Outfit, sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: Space Mono, monospace;
+  --shadow-color: hsl(196 83% 10%);
+  --shadow-opacity: 0.15;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 2px 4px 0px hsl(196 83% 10% / 0.07);
+  --shadow-xs: 2px 2px 4px 0px hsl(196 83% 10% / 0.07);
+  --shadow-sm:
+    2px 2px 4px 0px hsl(196 83% 10% / 0.15),
+    2px 1px 2px -1px hsl(196 83% 10% / 0.15);
+  --shadow:
+    2px 2px 4px 0px hsl(196 83% 10% / 0.15),
+    2px 1px 2px -1px hsl(196 83% 10% / 0.15);
+  --shadow-md:
+    2px 2px 4px 0px hsl(196 83% 10% / 0.15),
+    2px 2px 4px -1px hsl(196 83% 10% / 0.15);
+  --shadow-lg:
+    2px 2px 4px 0px hsl(196 83% 10% / 0.15),
+    2px 4px 6px -1px hsl(196 83% 10% / 0.15);
+  --shadow-xl:
+    2px 2px 4px 0px hsl(196 83% 10% / 0.15),
+    2px 8px 10px -1px hsl(196 83% 10% / 0.15);
+  --shadow-2xl: 2px 2px 4px 0px hsl(196 83% 10% / 0.38);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="retro-arcade-dark"] {
+  --background: oklch(0.27 0.05 219.82);
+  --foreground: oklch(0.7 0.02 196.79);
+  --card: oklch(0.31 0.05 219.65);
+  --card-foreground: oklch(0.7 0.02 196.79);
+  --popover: oklch(0.31 0.05 219.65);
+  --popover-foreground: oklch(0.7 0.02 196.79);
+  --primary: oklch(0.59 0.2 355.89);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.64 0.1 187.38);
+  --secondary-foreground: oklch(1 0 0);
+  --muted: oklch(0.52 0.03 219.14);
+  --muted-foreground: oklch(0.7 0.02 196.79);
+  --accent: oklch(0.58 0.17 39.5);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.59 0.21 27.12);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.52 0.03 219.14);
+  --input: oklch(0.52 0.03 219.14);
+  --ring: oklch(0.59 0.2 355.89);
+  --chart-1: oklch(0.61 0.14 244.93);
+  --chart-2: oklch(0.64 0.1 187.38);
+  --chart-3: oklch(0.59 0.2 355.89);
+  --chart-4: oklch(0.58 0.17 39.5);
+  --chart-5: oklch(0.59 0.21 27.12);
+  --radius: 0.25rem;
+  --sidebar: oklch(0.27 0.05 219.82);
+  --sidebar-foreground: oklch(0.7 0.02 196.79);
+  --sidebar-primary: oklch(0.59 0.2 355.89);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.64 0.1 187.38);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0.52 0.03 219.14);
+  --sidebar-ring: oklch(0.59 0.2 355.89);
+  --font-sans: Outfit, sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: Space Mono, monospace;
+  --shadow-color: hsl(196 83% 10%);
+  --shadow-opacity: 0.15;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 2px 4px 0px hsl(196 83% 10% / 0.07);
+  --shadow-xs: 2px 2px 4px 0px hsl(196 83% 10% / 0.07);
+  --shadow-sm:
+    2px 2px 4px 0px hsl(196 83% 10% / 0.15),
+    2px 1px 2px -1px hsl(196 83% 10% / 0.15);
+  --shadow:
+    2px 2px 4px 0px hsl(196 83% 10% / 0.15),
+    2px 1px 2px -1px hsl(196 83% 10% / 0.15);
+  --shadow-md:
+    2px 2px 4px 0px hsl(196 83% 10% / 0.15),
+    2px 2px 4px -1px hsl(196 83% 10% / 0.15);
+  --shadow-lg:
+    2px 2px 4px 0px hsl(196 83% 10% / 0.15),
+    2px 4px 6px -1px hsl(196 83% 10% / 0.15);
+  --shadow-xl:
+    2px 2px 4px 0px hsl(196 83% 10% / 0.15),
+    2px 8px 10px -1px hsl(196 83% 10% / 0.15);
+  --shadow-2xl: 2px 2px 4px 0px hsl(196 83% 10% / 0.38);
+}
diff --git a/web/src/styles/themes/solar-dusk.css b/web/src/styles/themes/solar-dusk.css
new file mode 100644
index 0000000..7e8d91c
--- /dev/null
+++ b/web/src/styles/themes/solar-dusk.css
@@ -0,0 +1,129 @@
+/* Solar Dusk Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="solar-dusk-light"] {
+  --background: oklch(0.99 0.01 84.57);
+  --foreground: oklch(0.37 0.03 49.61);
+  --card: oklch(0.97 0.01 78.28);
+  --card-foreground: oklch(0.37 0.03 49.61);
+  --popover: oklch(0.97 0.01 78.28);
+  --popover-foreground: oklch(0.37 0.03 49.61);
+  --primary: oklch(0.56 0.15 49);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.83 0.08 74.44);
+  --secondary-foreground: oklch(0.44 0.01 73.64);
+  --muted: oklch(0.94 0.02 83.26);
+  --muted-foreground: oklch(0.55 0.01 58.07);
+  --accent: oklch(0.9 0.05 74.99);
+  --accent-foreground: oklch(0.44 0.01 73.64);
+  --destructive: oklch(0.44 0.16 26.9);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.89 0.04 89.7);
+  --input: oklch(0.89 0.04 89.7);
+  --ring: oklch(0.56 0.15 49);
+  --chart-1: oklch(0.56 0.15 49);
+  --chart-2: oklch(0.55 0.01 58.07);
+  --chart-3: oklch(0.55 0.12 66.44);
+  --chart-4: oklch(0.55 0.01 58.07);
+  --chart-5: oklch(0.68 0.14 75.83);
+  --radius: 0.3rem;
+  --sidebar: oklch(0.94 0.02 83.26);
+  --sidebar-foreground: oklch(0.37 0.03 49.61);
+  --sidebar-primary: oklch(0.56 0.15 49);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.55 0.12 66.44);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0.89 0.04 89.7);
+  --sidebar-ring: oklch(0.56 0.15 49);
+  --font-sans: Oxanium, sans-serif;
+  --font-serif: Merriweather, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(28 18% 25%);
+  --shadow-opacity: 0.18;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 2px 3px 0px hsl(28 18% 25% / 0.09);
+  --shadow-xs: 0px 2px 3px 0px hsl(28 18% 25% / 0.09);
+  --shadow-sm:
+    0px 2px 3px 0px hsl(28 18% 25% / 0.18),
+    0px 1px 2px -1px hsl(28 18% 25% / 0.18);
+  --shadow:
+    0px 2px 3px 0px hsl(28 18% 25% / 0.18),
+    0px 1px 2px -1px hsl(28 18% 25% / 0.18);
+  --shadow-md:
+    0px 2px 3px 0px hsl(28 18% 25% / 0.18),
+    0px 2px 4px -1px hsl(28 18% 25% / 0.18);
+  --shadow-lg:
+    0px 2px 3px 0px hsl(28 18% 25% / 0.18),
+    0px 4px 6px -1px hsl(28 18% 25% / 0.18);
+  --shadow-xl:
+    0px 2px 3px 0px hsl(28 18% 25% / 0.18),
+    0px 8px 10px -1px hsl(28 18% 25% / 0.18);
+  --shadow-2xl: 0px 2px 3px 0px hsl(28 18% 25% / 0.45);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="solar-dusk-dark"] {
+  --background: oklch(0.22 0.01 56.04);
+  --foreground: oklch(0.97 0 106.42);
+  --card: oklch(0.27 0.01 34.3);
+  --card-foreground: oklch(0.97 0 106.42);
+  --popover: oklch(0.27 0.01 34.3);
+  --popover-foreground: oklch(0.97 0 106.42);
+  --primary: oklch(0.7 0.19 47.6);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.44 0.01 73.64);
+  --secondary-foreground: oklch(0.92 0 48.72);
+  --muted: oklch(0.27 0.01 34.3);
+  --muted-foreground: oklch(0.72 0.01 56.26);
+  --accent: oklch(0.36 0.05 229.32);
+  --accent-foreground: oklch(0.92 0 48.72);
+  --destructive: oklch(0.58 0.22 27.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.37 0.01 67.56);
+  --input: oklch(0.37 0.01 67.56);
+  --ring: oklch(0.7 0.19 47.6);
+  --chart-1: oklch(0.7 0.19 47.6);
+  --chart-2: oklch(0.68 0.15 237.32);
+  --chart-3: oklch(0.8 0.16 86.05);
+  --chart-4: oklch(0.72 0.01 56.26);
+  --chart-5: oklch(0.55 0.01 58.07);
+  --radius: 0.3rem;
+  --sidebar: oklch(0.27 0.01 34.3);
+  --sidebar-foreground: oklch(0.97 0 106.42);
+  --sidebar-primary: oklch(0.7 0.19 47.6);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.68 0.15 237.32);
+  --sidebar-accent-foreground: oklch(0.28 0.07 254.54);
+  --sidebar-border: oklch(0.37 0.01 67.56);
+  --sidebar-ring: oklch(0.7 0.19 47.6);
+  --font-sans: Oxanium, sans-serif;
+  --font-serif: Merriweather, serif;
+  --font-mono: Fira Code, monospace;
+  --shadow-color: hsl(0 0% 5%);
+  --shadow-opacity: 0.18;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 2px 3px 0px hsl(0 0% 5% / 0.09);
+  --shadow-xs: 0px 2px 3px 0px hsl(0 0% 5% / 0.09);
+  --shadow-sm:
+    0px 2px 3px 0px hsl(0 0% 5% / 0.18), 0px 1px 2px -1px hsl(0 0% 5% / 0.18);
+  --shadow:
+    0px 2px 3px 0px hsl(0 0% 5% / 0.18), 0px 1px 2px -1px hsl(0 0% 5% / 0.18);
+  --shadow-md:
+    0px 2px 3px 0px hsl(0 0% 5% / 0.18), 0px 2px 4px -1px hsl(0 0% 5% / 0.18);
+  --shadow-lg:
+    0px 2px 3px 0px hsl(0 0% 5% / 0.18), 0px 4px 6px -1px hsl(0 0% 5% / 0.18);
+  --shadow-xl:
+    0px 2px 3px 0px hsl(0 0% 5% / 0.18), 0px 8px 10px -1px hsl(0 0% 5% / 0.18);
+  --shadow-2xl: 0px 2px 3px 0px hsl(0 0% 5% / 0.45);
+}
diff --git a/web/src/styles/themes/spotify.css b/web/src/styles/themes/spotify.css
new file mode 100644
index 0000000..e933ef0
--- /dev/null
+++ b/web/src/styles/themes/spotify.css
@@ -0,0 +1,114 @@
+/* Spotify Theme - Official Spotify brand colors */
+
+/* Light mode */
+[data-theme="spotify-light"] {
+  --background: oklch(1 0 0);
+  --foreground: oklch(0.16 0.008 25);
+  --card: oklch(0.98 0.003 0);
+  --card-foreground: oklch(0.16 0.008 25);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.16 0.008 25);
+  --primary: oklch(0.64 0.18 152);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.95 0.008 152);
+  --secondary-foreground: oklch(0.20 0.01 25);
+  --muted: oklch(0.96 0.003 0);
+  --muted-foreground: oklch(0.50 0.01 0);
+  --accent: oklch(0.72 0.20 150);
+  --accent-foreground: oklch(0.16 0.008 25);
+  --destructive: oklch(0.55 0.22 25);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.90 0.003 0);
+  --input: oklch(0.92 0.003 0);
+  --ring: oklch(0.64 0.18 152);
+  --chart-1: oklch(0.64 0.18 152);
+  --chart-2: oklch(0.72 0.20 150);
+  --chart-3: oklch(0.60 0.15 200);
+  --chart-4: oklch(0.70 0.12 280);
+  --chart-5: oklch(0.65 0.18 340);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.16 0.008 25);
+  --sidebar-foreground: oklch(0.95 0.003 0);
+  --sidebar-primary: oklch(0.64 0.18 152);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.25 0.01 25);
+  --sidebar-accent-foreground: oklch(0.95 0.003 0);
+  --sidebar-border: oklch(0.28 0.01 25);
+  --sidebar-ring: oklch(0.64 0.18 152);
+  --font-sans: Montserrat, ui-sans-serif, system-ui, sans-serif;
+  --font-serif: ui-serif, Georgia, serif;
+  --font-mono: ui-monospace, SFMono-Regular, monospace;
+  --shadow-color: hsl(150 60% 10%);
+  --shadow-opacity: 0.08;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 2px;
+  --letter-spacing: -0.01em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 2px 0px hsl(150 60% 10% / 0.04);
+  --shadow-xs: 0 1px 3px 0px hsl(150 60% 10% / 0.06);
+  --shadow-sm: 0 2px 4px 0px hsl(150 60% 10% / 0.08), 0 1px 2px -1px hsl(150 60% 10% / 0.08);
+  --shadow: 0 2px 4px 0px hsl(150 60% 10% / 0.08), 0 1px 2px -1px hsl(150 60% 10% / 0.08);
+  --shadow-md: 0 4px 6px 0px hsl(150 60% 10% / 0.1), 0 2px 4px -1px hsl(150 60% 10% / 0.1);
+  --shadow-lg: 0 8px 10px 0px hsl(150 60% 10% / 0.1), 0 4px 6px -1px hsl(150 60% 10% / 0.1);
+  --shadow-xl: 0 16px 20px 0px hsl(150 60% 10% / 0.12), 0 8px 10px -1px hsl(150 60% 10% / 0.12);
+  --shadow-2xl: 0 24px 48px 0px hsl(150 60% 10% / 0.2);
+  --tracking-normal: -0.01em;
+}
+
+/* Dark mode - #121212 is Spotify's signature dark background */
+[data-theme="spotify-dark"] {
+  --background: oklch(0.145 0 0);
+  --foreground: oklch(1 0 0);
+  --card: oklch(0.18 0 0);
+  --card-foreground: oklch(1 0 0);
+  --popover: oklch(0.20 0 0);
+  --popover-foreground: oklch(1 0 0);
+  --primary: oklch(0.64 0.18 152);
+  --primary-foreground: oklch(0.145 0 0);
+  --secondary: oklch(0.22 0 0);
+  --secondary-foreground: oklch(0.90 0 0);
+  --muted: oklch(0.20 0 0);
+  --muted-foreground: oklch(0.65 0 0);
+  --accent: oklch(0.72 0.20 150);
+  --accent-foreground: oklch(0.145 0 0);
+  --destructive: oklch(0.60 0.20 25);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.28 0 0);
+  --input: oklch(0.22 0 0);
+  --ring: oklch(0.64 0.18 152);
+  --chart-1: oklch(0.64 0.18 152);
+  --chart-2: oklch(0.72 0.20 150);
+  --chart-3: oklch(0.65 0.14 200);
+  --chart-4: oklch(0.75 0.12 280);
+  --chart-5: oklch(0.70 0.16 340);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.12 0 0);
+  --sidebar-foreground: oklch(1 0 0);
+  --sidebar-primary: oklch(0.64 0.18 152);
+  --sidebar-primary-foreground: oklch(0.145 0 0);
+  --sidebar-accent: oklch(0.22 0 0);
+  --sidebar-accent-foreground: oklch(0.90 0 0);
+  --sidebar-border: oklch(0.28 0 0);
+  --sidebar-ring: oklch(0.64 0.18 152);
+  --font-sans: Montserrat, ui-sans-serif, system-ui, sans-serif;
+  --font-serif: ui-serif, Georgia, serif;
+  --font-mono: ui-monospace, SFMono-Regular, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.5;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 2px;
+  --letter-spacing: -0.01em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 2px 0px hsl(0 0% 0% / 0.25);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.35);
+  --shadow-sm: 0 2px 4px 0px hsl(0 0% 0% / 0.5), 0 1px 2px -1px hsl(0 0% 0% / 0.5);
+  --shadow: 0 2px 4px 0px hsl(0 0% 0% / 0.5), 0 1px 2px -1px hsl(0 0% 0% / 0.5);
+  --shadow-md: 0 4px 6px 0px hsl(0 0% 0% / 0.5), 0 2px 4px -1px hsl(0 0% 0% / 0.5);
+  --shadow-lg: 0 8px 10px 0px hsl(0 0% 0% / 0.5), 0 4px 6px -1px hsl(0 0% 0% / 0.5);
+  --shadow-xl: 0 16px 20px 0px hsl(0 0% 0% / 0.5), 0 8px 10px -1px hsl(0 0% 0% / 0.5);
+  --shadow-2xl: 0 24px 48px 0px hsl(0 0% 0% / 0.6);
+}
diff --git a/web/src/styles/themes/starry-night.css b/web/src/styles/themes/starry-night.css
new file mode 100644
index 0000000..bcbc72f
--- /dev/null
+++ b/web/src/styles/themes/starry-night.css
@@ -0,0 +1,126 @@
+/* Starry Night Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="starry-night-light"] {
+  --background: oklch(0.98 0 258.32);
+  --foreground: oklch(0.26 0.04 268.07);
+  --card: oklch(0.93 0.01 251.56);
+  --card-foreground: oklch(0.26 0.04 268.07);
+  --popover: oklch(0.99 0.03 98.05);
+  --popover-foreground: oklch(0.26 0.04 268.07);
+  --primary: oklch(0.48 0.12 263.38);
+  --primary-foreground: oklch(0.99 0.03 98.05);
+  --secondary: oklch(0.86 0.12 81.01);
+  --secondary-foreground: oklch(0.26 0.04 268.07);
+  --muted: oklch(0.92 0.01 106.56);
+  --muted-foreground: oklch(0.48 0.12 263.38);
+  --accent: oklch(0.69 0.07 234.04);
+  --accent-foreground: oklch(0.99 0.03 98.05);
+  --destructive: oklch(0.26 0.04 322.53);
+  --destructive-foreground: oklch(0.99 0.03 98.05);
+  --border: oklch(0.78 0.02 251.19);
+  --input: oklch(0.69 0.07 234.04);
+  --ring: oklch(0.86 0.12 81.01);
+  --chart-1: oklch(0.48 0.12 263.38);
+  --chart-2: oklch(0.86 0.12 81.01);
+  --chart-3: oklch(0.69 0.07 234.04);
+  --chart-4: oklch(0.78 0.02 251.19);
+  --chart-5: oklch(0.26 0.04 322.53);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.93 0.01 251.56);
+  --sidebar-foreground: oklch(0.26 0.04 268.07);
+  --sidebar-primary: oklch(0.48 0.12 263.38);
+  --sidebar-primary-foreground: oklch(0.99 0.03 98.05);
+  --sidebar-accent: oklch(0.86 0.12 81.01);
+  --sidebar-accent-foreground: oklch(0.26 0.04 268.07);
+  --sidebar-border: oklch(0.78 0.02 251.19);
+  --sidebar-ring: oklch(0.86 0.12 81.01);
+  --font-sans: Libre Baskerville, serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono:
+    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
+    "Courier New", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="starry-night-dark"] {
+  --background: oklch(0.22 0.02 275.84);
+  --foreground: oklch(0.94 0.01 266.7);
+  --card: oklch(0.27 0.04 281.3);
+  --card-foreground: oklch(0.94 0.01 266.7);
+  --popover: oklch(0.27 0.04 281.3);
+  --popover-foreground: oklch(0.91 0.14 95.11);
+  --primary: oklch(0.48 0.12 263.38);
+  --primary-foreground: oklch(0.91 0.14 95.11);
+  --secondary: oklch(0.91 0.14 95.11);
+  --secondary-foreground: oklch(0.27 0.04 281.3);
+  --muted: oklch(0.27 0.04 281.3);
+  --muted-foreground: oklch(0.62 0.04 262.04);
+  --accent: oklch(0.85 0.05 264.78);
+  --accent-foreground: oklch(0.22 0.02 275.84);
+  --destructive: oklch(0.53 0.12 357.11);
+  --destructive-foreground: oklch(0.91 0.14 95.11);
+  --border: oklch(0.31 0.03 281.77);
+  --input: oklch(0.48 0.12 263.38);
+  --ring: oklch(0.91 0.14 95.11);
+  --chart-1: oklch(0.48 0.12 263.38);
+  --chart-2: oklch(0.91 0.14 95.11);
+  --chart-3: oklch(0.69 0.07 234.04);
+  --chart-4: oklch(0.62 0.04 262.04);
+  --chart-5: oklch(0.53 0.12 357.11);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.27 0.04 281.3);
+  --sidebar-foreground: oklch(0.94 0.01 266.7);
+  --sidebar-primary: oklch(0.48 0.12 263.38);
+  --sidebar-primary-foreground: oklch(0.91 0.14 95.11);
+  --sidebar-accent: oklch(0.91 0.14 95.11);
+  --sidebar-accent-foreground: oklch(0.27 0.04 281.3);
+  --sidebar-border: oklch(0.31 0.03 281.77);
+  --sidebar-ring: oklch(0.91 0.14 95.11);
+  --font-sans: Libre Baskerville, serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono:
+    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
+    "Courier New", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/stripe.css b/web/src/styles/themes/stripe.css
new file mode 100644
index 0000000..4ab63d1
--- /dev/null
+++ b/web/src/styles/themes/stripe.css
@@ -0,0 +1,128 @@
+/* Stripe Theme - Official Stripe brand colors */
+
+/* Light mode */
+[data-theme="stripe-light"] {
+  --background: oklch(0.98 0.006 240);
+  --foreground: oklch(0.22 0.05 250);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.22 0.05 250);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.22 0.05 250);
+  --primary: oklch(0.55 0.25 285);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.95 0.008 240);
+  --secondary-foreground: oklch(0.3 0.04 250);
+  --muted: oklch(0.96 0.006 240);
+  --muted-foreground: oklch(0.45 0.03 250);
+  --accent: oklch(0.82 0.14 200);
+  --accent-foreground: oklch(0.22 0.05 200);
+  --destructive: oklch(0.55 0.22 25);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.9 0.008 240);
+  --input: oklch(0.92 0.008 240);
+  --ring: oklch(0.55 0.25 285);
+  --chart-1: oklch(0.55 0.25 285);
+  --chart-2: oklch(0.82 0.14 200);
+  --chart-3: oklch(0.65 0.18 320);
+  --chart-4: oklch(0.7 0.15 150);
+  --chart-5: oklch(0.6 0.2 30);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.22 0.05 250);
+  --sidebar-foreground: oklch(0.95 0.006 240);
+  --sidebar-primary: oklch(0.6 0.22 285);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.3 0.05 250);
+  --sidebar-accent-foreground: oklch(0.95 0.006 240);
+  --sidebar-border: oklch(0.35 0.05 250);
+  --sidebar-ring: oklch(0.6 0.22 285);
+  --font-sans:
+    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
+  --font-serif: ui-serif, Georgia, serif;
+  --font-mono: ui-monospace, SFMono-Regular, "SF Mono", monospace;
+  --shadow-color: hsl(220 60% 20%);
+  --shadow-opacity: 0.08;
+  --shadow-blur: 6px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 2px;
+  --letter-spacing: -0.01em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 2px 0px hsl(220 60% 20% / 0.04);
+  --shadow-xs: 0 1px 3px 0px hsl(220 60% 20% / 0.06);
+  --shadow-sm:
+    0 2px 4px 0px hsl(220 60% 20% / 0.08),
+    0 1px 2px -1px hsl(220 60% 20% / 0.08);
+  --shadow:
+    0 2px 4px 0px hsl(220 60% 20% / 0.08),
+    0 1px 2px -1px hsl(220 60% 20% / 0.08);
+  --shadow-md:
+    0 4px 8px 0px hsl(220 60% 20% / 0.1), 0 2px 4px -1px hsl(220 60% 20% / 0.1);
+  --shadow-lg:
+    0 8px 16px 0px hsl(220 60% 20% / 0.1), 0 4px 8px -1px hsl(220 60% 20% / 0.1);
+  --shadow-xl:
+    0 16px 24px 0px hsl(220 60% 20% / 0.12),
+    0 8px 12px -1px hsl(220 60% 20% / 0.12);
+  --shadow-2xl: 0 24px 48px 0px hsl(220 60% 20% / 0.2);
+  --tracking-normal: -0.01em;
+}
+
+/* Dark mode - #0A2540 is Stripe's signature navy background */
+[data-theme="stripe-dark"] {
+  --background: oklch(0.22 0.05 250);
+  --foreground: oklch(0.95 0.006 240);
+  --card: oklch(0.26 0.05 250);
+  --card-foreground: oklch(0.95 0.006 240);
+  --popover: oklch(0.28 0.05 250);
+  --popover-foreground: oklch(0.95 0.006 240);
+  --primary: oklch(0.6 0.23 285);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.3 0.05 250);
+  --secondary-foreground: oklch(0.9 0.006 240);
+  --muted: oklch(0.28 0.04 250);
+  --muted-foreground: oklch(0.65 0.02 240);
+  --accent: oklch(0.82 0.14 200);
+  --accent-foreground: oklch(0.22 0.05 200);
+  --destructive: oklch(0.6 0.2 25);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.35 0.05 250);
+  --input: oklch(0.3 0.05 250);
+  --ring: oklch(0.6 0.23 285);
+  --chart-1: oklch(0.6 0.23 285);
+  --chart-2: oklch(0.82 0.14 200);
+  --chart-3: oklch(0.7 0.16 320);
+  --chart-4: oklch(0.75 0.14 150);
+  --chart-5: oklch(0.65 0.18 30);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.18 0.05 250);
+  --sidebar-foreground: oklch(0.95 0.006 240);
+  --sidebar-primary: oklch(0.6 0.23 285);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.3 0.05 250);
+  --sidebar-accent-foreground: oklch(0.9 0.006 240);
+  --sidebar-border: oklch(0.35 0.05 250);
+  --sidebar-ring: oklch(0.6 0.23 285);
+  --font-sans:
+    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
+  --font-serif: ui-serif, Georgia, serif;
+  --font-mono: ui-monospace, SFMono-Regular, "SF Mono", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.4;
+  --shadow-blur: 6px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 2px;
+  --letter-spacing: -0.01em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 2px 0px hsl(0 0% 0% / 0.2);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.3);
+  --shadow-sm:
+    0 2px 4px 0px hsl(0 0% 0% / 0.4), 0 1px 2px -1px hsl(0 0% 0% / 0.4);
+  --shadow: 0 2px 4px 0px hsl(0 0% 0% / 0.4), 0 1px 2px -1px hsl(0 0% 0% / 0.4);
+  --shadow-md:
+    0 4px 8px 0px hsl(0 0% 0% / 0.4), 0 2px 4px -1px hsl(0 0% 0% / 0.4);
+  --shadow-lg:
+    0 8px 16px 0px hsl(0 0% 0% / 0.4), 0 4px 8px -1px hsl(0 0% 0% / 0.4);
+  --shadow-xl:
+    0 16px 24px 0px hsl(0 0% 0% / 0.4), 0 8px 12px -1px hsl(0 0% 0% / 0.4);
+  --shadow-2xl: 0 24px 48px 0px hsl(0 0% 0% / 0.5);
+}
diff --git a/web/src/styles/themes/sunset-horizon.css b/web/src/styles/themes/sunset-horizon.css
new file mode 100644
index 0000000..0e4463d
--- /dev/null
+++ b/web/src/styles/themes/sunset-horizon.css
@@ -0,0 +1,124 @@
+/* Sunset Horizon Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="sunset-horizon-light"] {
+  --background: oklch(0.99 0.01 56.32);
+  --foreground: oklch(0.34 0.01 2.77);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.34 0.01 2.77);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.34 0.01 2.77);
+  --primary: oklch(0.74 0.16 34.71);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.96 0.02 28.9);
+  --secondary-foreground: oklch(0.56 0.13 32.74);
+  --muted: oklch(0.97 0.02 39.4);
+  --muted-foreground: oklch(0.55 0.01 58.07);
+  --accent: oklch(0.83 0.11 58);
+  --accent-foreground: oklch(0.34 0.01 2.77);
+  --destructive: oklch(0.61 0.21 22.24);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.93 0.04 38.69);
+  --input: oklch(0.93 0.04 38.69);
+  --ring: oklch(0.74 0.16 34.71);
+  --chart-1: oklch(0.74 0.16 34.71);
+  --chart-2: oklch(0.83 0.11 58);
+  --chart-3: oklch(0.88 0.08 54.93);
+  --chart-4: oklch(0.82 0.11 40.89);
+  --chart-5: oklch(0.64 0.13 32.07);
+  --radius: 0.625rem;
+  --sidebar: oklch(0.97 0.02 39.4);
+  --sidebar-foreground: oklch(0.34 0.01 2.77);
+  --sidebar-primary: oklch(0.74 0.16 34.71);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.83 0.11 58);
+  --sidebar-accent-foreground: oklch(0.34 0.01 2.77);
+  --sidebar-border: oklch(0.93 0.04 38.69);
+  --sidebar-ring: oklch(0.74 0.16 34.71);
+  --font-sans: Montserrat, sans-serif;
+  --font-serif: Merriweather, serif;
+  --font-mono: Ubuntu Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.09;
+  --shadow-blur: 12px;
+  --shadow-spread: -3px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 6px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 6px 12px -3px hsl(0 0% 0% / 0.04);
+  --shadow-xs: 0px 6px 12px -3px hsl(0 0% 0% / 0.04);
+  --shadow-sm:
+    0px 6px 12px -3px hsl(0 0% 0% / 0.09), 0px 1px 2px -4px hsl(0 0% 0% / 0.09);
+  --shadow:
+    0px 6px 12px -3px hsl(0 0% 0% / 0.09), 0px 1px 2px -4px hsl(0 0% 0% / 0.09);
+  --shadow-md:
+    0px 6px 12px -3px hsl(0 0% 0% / 0.09), 0px 2px 4px -4px hsl(0 0% 0% / 0.09);
+  --shadow-lg:
+    0px 6px 12px -3px hsl(0 0% 0% / 0.09), 0px 4px 6px -4px hsl(0 0% 0% / 0.09);
+  --shadow-xl:
+    0px 6px 12px -3px hsl(0 0% 0% / 0.09), 0px 8px 10px -4px hsl(0 0% 0% / 0.09);
+  --shadow-2xl: 0px 6px 12px -3px hsl(0 0% 0% / 0.22);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="sunset-horizon-dark"] {
+  --background: oklch(0.26 0.02 352.4);
+  --foreground: oklch(0.94 0.01 51.32);
+  --card: oklch(0.32 0.02 341.45);
+  --card-foreground: oklch(0.94 0.01 51.32);
+  --popover: oklch(0.32 0.02 341.45);
+  --popover-foreground: oklch(0.94 0.01 51.32);
+  --primary: oklch(0.74 0.16 34.71);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.36 0.02 342.27);
+  --secondary-foreground: oklch(0.94 0.01 51.32);
+  --muted: oklch(0.32 0.02 341.45);
+  --muted-foreground: oklch(0.84 0.02 52.63);
+  --accent: oklch(0.83 0.11 58);
+  --accent-foreground: oklch(0.26 0.02 352.4);
+  --destructive: oklch(0.61 0.21 22.24);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.36 0.02 342.27);
+  --input: oklch(0.36 0.02 342.27);
+  --ring: oklch(0.74 0.16 34.71);
+  --chart-1: oklch(0.74 0.16 34.71);
+  --chart-2: oklch(0.83 0.11 58);
+  --chart-3: oklch(0.88 0.08 54.93);
+  --chart-4: oklch(0.82 0.11 40.89);
+  --chart-5: oklch(0.64 0.13 32.07);
+  --radius: 0.625rem;
+  --sidebar: oklch(0.26 0.02 352.4);
+  --sidebar-foreground: oklch(0.94 0.01 51.32);
+  --sidebar-primary: oklch(0.74 0.16 34.71);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.83 0.11 58);
+  --sidebar-accent-foreground: oklch(0.26 0.02 352.4);
+  --sidebar-border: oklch(0.36 0.02 342.27);
+  --sidebar-ring: oklch(0.74 0.16 34.71);
+  --font-sans: Montserrat, sans-serif;
+  --font-serif: Merriweather, serif;
+  --font-mono: Ubuntu Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.09;
+  --shadow-blur: 12px;
+  --shadow-spread: -3px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 6px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 6px 12px -3px hsl(0 0% 0% / 0.04);
+  --shadow-xs: 0px 6px 12px -3px hsl(0 0% 0% / 0.04);
+  --shadow-sm:
+    0px 6px 12px -3px hsl(0 0% 0% / 0.09), 0px 1px 2px -4px hsl(0 0% 0% / 0.09);
+  --shadow:
+    0px 6px 12px -3px hsl(0 0% 0% / 0.09), 0px 1px 2px -4px hsl(0 0% 0% / 0.09);
+  --shadow-md:
+    0px 6px 12px -3px hsl(0 0% 0% / 0.09), 0px 2px 4px -4px hsl(0 0% 0% / 0.09);
+  --shadow-lg:
+    0px 6px 12px -3px hsl(0 0% 0% / 0.09), 0px 4px 6px -4px hsl(0 0% 0% / 0.09);
+  --shadow-xl:
+    0px 6px 12px -3px hsl(0 0% 0% / 0.09), 0px 8px 10px -4px hsl(0 0% 0% / 0.09);
+  --shadow-2xl: 0px 6px 12px -3px hsl(0 0% 0% / 0.22);
+}
diff --git a/web/src/styles/themes/supabase.css b/web/src/styles/themes/supabase.css
new file mode 100644
index 0000000..28d6d00
--- /dev/null
+++ b/web/src/styles/themes/supabase.css
@@ -0,0 +1,124 @@
+/* Supabase Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="supabase-light"] {
+  --background: oklch(0.99 0 0);
+  --foreground: oklch(0.2 0 0);
+  --card: oklch(0.99 0 0);
+  --card-foreground: oklch(0.2 0 0);
+  --popover: oklch(0.99 0 0);
+  --popover-foreground: oklch(0.44 0 0);
+  --primary: oklch(0.83 0.13 160.91);
+  --primary-foreground: oklch(0.26 0.01 166.46);
+  --secondary: oklch(0.99 0 0);
+  --secondary-foreground: oklch(0.2 0 0);
+  --muted: oklch(0.95 0 0);
+  --muted-foreground: oklch(0.24 0 0);
+  --accent: oklch(0.95 0 0);
+  --accent-foreground: oklch(0.24 0 0);
+  --destructive: oklch(0.55 0.19 32.73);
+  --destructive-foreground: oklch(0.99 0 17.21);
+  --border: oklch(0.9 0 0);
+  --input: oklch(0.97 0 0);
+  --ring: oklch(0.83 0.13 160.91);
+  --chart-1: oklch(0.83 0.13 160.91);
+  --chart-2: oklch(0.62 0.19 259.81);
+  --chart-3: oklch(0.61 0.22 292.72);
+  --chart-4: oklch(0.77 0.16 70.08);
+  --chart-5: oklch(0.7 0.15 162.48);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.99 0 0);
+  --sidebar-foreground: oklch(0.55 0 0);
+  --sidebar-primary: oklch(0.83 0.13 160.91);
+  --sidebar-primary-foreground: oklch(0.26 0.01 166.46);
+  --sidebar-accent: oklch(0.95 0 0);
+  --sidebar-accent-foreground: oklch(0.24 0 0);
+  --sidebar-border: oklch(0.9 0 0);
+  --sidebar-ring: oklch(0.83 0.13 160.91);
+  --font-sans: Outfit, sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: monospace;
+  --shadow-color: #000000;
+  --shadow-opacity: 0.17;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0.025em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 1px 3px 0px hsl(0 0% 0% / 0.09);
+  --shadow-xs: 0px 1px 3px 0px hsl(0 0% 0% / 0.09);
+  --shadow-sm:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.17), 0px 1px 2px -1px hsl(0 0% 0% / 0.17);
+  --shadow:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.17), 0px 1px 2px -1px hsl(0 0% 0% / 0.17);
+  --shadow-md:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.17), 0px 2px 4px -1px hsl(0 0% 0% / 0.17);
+  --shadow-lg:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.17), 0px 4px 6px -1px hsl(0 0% 0% / 0.17);
+  --shadow-xl:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.17), 0px 8px 10px -1px hsl(0 0% 0% / 0.17);
+  --shadow-2xl: 0px 1px 3px 0px hsl(0 0% 0% / 0.43);
+  --tracking-normal: 0.025em;
+}
+
+/* Dark mode */
+[data-theme="supabase-dark"] {
+  --background: oklch(0.18 0 0);
+  --foreground: oklch(0.93 0.01 255.51);
+  --card: oklch(0.2 0 0);
+  --card-foreground: oklch(0.93 0.01 255.51);
+  --popover: oklch(0.26 0 0);
+  --popover-foreground: oklch(0.73 0 0);
+  --primary: oklch(0.44 0.1 156.76);
+  --primary-foreground: oklch(0.92 0.01 167.16);
+  --secondary: oklch(0.26 0 0);
+  --secondary-foreground: oklch(0.99 0 0);
+  --muted: oklch(0.24 0 0);
+  --muted-foreground: oklch(0.71 0 0);
+  --accent: oklch(0.31 0 0);
+  --accent-foreground: oklch(0.99 0 0);
+  --destructive: oklch(0.31 0.09 29.79);
+  --destructive-foreground: oklch(0.94 0 34.31);
+  --border: oklch(0.28 0 0);
+  --input: oklch(0.26 0 0);
+  --ring: oklch(0.8 0.18 151.71);
+  --chart-1: oklch(0.8 0.18 151.71);
+  --chart-2: oklch(0.71 0.14 254.62);
+  --chart-3: oklch(0.71 0.16 293.54);
+  --chart-4: oklch(0.84 0.16 84.43);
+  --chart-5: oklch(0.78 0.13 181.91);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.18 0 0);
+  --sidebar-foreground: oklch(0.63 0 0);
+  --sidebar-primary: oklch(0.44 0.1 156.76);
+  --sidebar-primary-foreground: oklch(0.92 0.01 167.16);
+  --sidebar-accent: oklch(0.31 0 0);
+  --sidebar-accent-foreground: oklch(0.99 0 0);
+  --sidebar-border: oklch(0.28 0 0);
+  --sidebar-ring: oklch(0.8 0.18 151.71);
+  --font-sans: Outfit, sans-serif;
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono: monospace;
+  --shadow-color: #000000;
+  --shadow-opacity: 0.17;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0.025em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 1px 3px 0px hsl(0 0% 0% / 0.09);
+  --shadow-xs: 0px 1px 3px 0px hsl(0 0% 0% / 0.09);
+  --shadow-sm:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.17), 0px 1px 2px -1px hsl(0 0% 0% / 0.17);
+  --shadow:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.17), 0px 1px 2px -1px hsl(0 0% 0% / 0.17);
+  --shadow-md:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.17), 0px 2px 4px -1px hsl(0 0% 0% / 0.17);
+  --shadow-lg:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.17), 0px 4px 6px -1px hsl(0 0% 0% / 0.17);
+  --shadow-xl:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.17), 0px 8px 10px -1px hsl(0 0% 0% / 0.17);
+  --shadow-2xl: 0px 1px 3px 0px hsl(0 0% 0% / 0.43);
+}
diff --git a/web/src/styles/themes/t3-chat.css b/web/src/styles/themes/t3-chat.css
new file mode 100644
index 0000000..3581686
--- /dev/null
+++ b/web/src/styles/themes/t3-chat.css
@@ -0,0 +1,132 @@
+/* T3 Chat Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="t3-chat-light"] {
+  --background: oklch(0.98 0.01 325.64);
+  --foreground: oklch(0.33 0.12 325.04);
+  --card: oklch(0.98 0.01 325.64);
+  --card-foreground: oklch(0.33 0.12 325.04);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.33 0.12 325.04);
+  --primary: oklch(0.53 0.14 355.2);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.87 0.07 334.9);
+  --secondary-foreground: oklch(0.44 0.13 324.8);
+  --muted: oklch(0.94 0.03 331.55);
+  --muted-foreground: oklch(0.49 0.12 324.45);
+  --accent: oklch(0.87 0.07 334.9);
+  --accent-foreground: oklch(0.44 0.13 324.8);
+  --destructive: oklch(0.52 0.14 20.83);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.86 0.08 328.91);
+  --input: oklch(0.85 0.06 336.6);
+  --ring: oklch(0.59 0.22 0.58);
+  --chart-1: oklch(0.6 0.24 344.47);
+  --chart-2: oklch(0.44 0.23 300.62);
+  --chart-3: oklch(0.38 0.04 226.15);
+  --chart-4: oklch(0.83 0.12 88.35);
+  --chart-5: oklch(0.78 0.13 59);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.94 0.03 320.58);
+  --sidebar-foreground: oklch(0.49 0.19 354.54);
+  --sidebar-primary: oklch(0.4 0.03 285.2);
+  --sidebar-primary-foreground: oklch(0.97 0.01 337.52);
+  --sidebar-accent: oklch(0.98 0 106.42);
+  --sidebar-accent-foreground: oklch(0.4 0.03 285.2);
+  --sidebar-border: oklch(0.94 0 48.72);
+  --sidebar-ring: oklch(0.59 0.22 0.58);
+  --font-sans:
+    ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
+    Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif,
+    "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono:
+    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
+    "Courier New", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="t3-chat-dark"] {
+  --background: oklch(0.24 0.02 307.53);
+  --foreground: oklch(0.84 0.04 309.54);
+  --card: oklch(0.28 0.02 307.54);
+  --card-foreground: oklch(0.85 0.03 341.46);
+  --popover: oklch(0.15 0.01 338.9);
+  --popover-foreground: oklch(0.96 0.01 341.8);
+  --primary: oklch(0.46 0.19 4.1);
+  --primary-foreground: oklch(0.86 0.06 346.37);
+  --secondary: oklch(0.31 0.03 310.06);
+  --secondary-foreground: oklch(0.85 0.04 307.96);
+  --muted: oklch(0.26 0.02 309.47);
+  --muted-foreground: oklch(0.79 0.04 307.1);
+  --accent: oklch(0.36 0.05 308.49);
+  --accent-foreground: oklch(0.96 0.01 341.8);
+  --destructive: oklch(0.23 0.05 12.61);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.33 0.02 343.45);
+  --input: oklch(0.34 0.02 332.83);
+  --ring: oklch(0.59 0.22 0.58);
+  --chart-1: oklch(0.53 0.14 355.2);
+  --chart-2: oklch(0.56 0.19 306.86);
+  --chart-3: oklch(0.72 0.15 60.58);
+  --chart-4: oklch(0.62 0.2 312.74);
+  --chart-5: oklch(0.61 0.21 6.14);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.19 0.02 331.05);
+  --sidebar-foreground: oklch(0.86 0.03 343.66);
+  --sidebar-primary: oklch(0.49 0.22 264.38);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.23 0.03 338.2);
+  --sidebar-accent-foreground: oklch(0.97 0 286.38);
+  --sidebar-border: oklch(0 0 0);
+  --sidebar-ring: oklch(0.59 0.22 0.58);
+  --font-sans:
+    ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
+    Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif,
+    "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
+  --font-serif: ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
+  --font-mono:
+    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
+    "Courier New", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow: 0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0 1px 3px 0px hsl(0 0% 0% / 0.1), 0 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0 1px 3px 0px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/tangerine.css b/web/src/styles/themes/tangerine.css
new file mode 100644
index 0000000..96de47a
--- /dev/null
+++ b/web/src/styles/themes/tangerine.css
@@ -0,0 +1,124 @@
+/* Tangerine Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="tangerine-light"] {
+  --background: oklch(0.94 0 236.5);
+  --foreground: oklch(0.32 0 0);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.32 0 0);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.32 0 0);
+  --primary: oklch(0.64 0.17 36.44);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.97 0 264.54);
+  --secondary-foreground: oklch(0.45 0.03 256.8);
+  --muted: oklch(0.98 0 247.84);
+  --muted-foreground: oklch(0.55 0.02 264.36);
+  --accent: oklch(0.91 0.02 243.82);
+  --accent-foreground: oklch(0.38 0.14 265.52);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.9 0.01 247.88);
+  --input: oklch(0.97 0 264.54);
+  --ring: oklch(0.64 0.17 36.44);
+  --chart-1: oklch(0.72 0.06 248.68);
+  --chart-2: oklch(0.79 0.09 35.96);
+  --chart-3: oklch(0.58 0.08 254.16);
+  --chart-4: oklch(0.5 0.08 259.49);
+  --chart-5: oklch(0.42 0.1 264.03);
+  --radius: 0.75rem;
+  --sidebar: oklch(0.9 0 258.33);
+  --sidebar-foreground: oklch(0.32 0 0);
+  --sidebar-primary: oklch(0.64 0.17 36.44);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.91 0.02 243.82);
+  --sidebar-accent-foreground: oklch(0.38 0.14 265.52);
+  --sidebar-border: oklch(0.93 0.01 264.53);
+  --sidebar-ring: oklch(0.64 0.17 36.44);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.1), 0px 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.1), 0px 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.1), 0px 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.1), 0px 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.1), 0px 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 1px 3px 0px hsl(0 0% 0% / 0.25);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="tangerine-dark"] {
+  --background: oklch(0.26 0.03 262.67);
+  --foreground: oklch(0.92 0 0);
+  --card: oklch(0.31 0.03 268.64);
+  --card-foreground: oklch(0.92 0 0);
+  --popover: oklch(0.29 0.02 268.4);
+  --popover-foreground: oklch(0.92 0 0);
+  --primary: oklch(0.64 0.17 36.44);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.31 0.03 266.71);
+  --secondary-foreground: oklch(0.92 0 0);
+  --muted: oklch(0.31 0.03 266.71);
+  --muted-foreground: oklch(0.72 0 0);
+  --accent: oklch(0.34 0.06 267.59);
+  --accent-foreground: oklch(0.88 0.06 254.13);
+  --destructive: oklch(0.64 0.21 25.33);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.38 0.03 269.73);
+  --input: oklch(0.38 0.03 269.73);
+  --ring: oklch(0.64 0.17 36.44);
+  --chart-1: oklch(0.72 0.06 248.68);
+  --chart-2: oklch(0.77 0.09 34.19);
+  --chart-3: oklch(0.58 0.08 254.16);
+  --chart-4: oklch(0.5 0.08 259.49);
+  --chart-5: oklch(0.42 0.1 264.03);
+  --radius: 0.75rem;
+  --sidebar: oklch(0.31 0.03 267.74);
+  --sidebar-foreground: oklch(0.92 0 0);
+  --sidebar-primary: oklch(0.64 0.17 36.44);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.34 0.06 267.59);
+  --sidebar-accent-foreground: oklch(0.88 0.06 254.13);
+  --sidebar-border: oklch(0.38 0.03 269.73);
+  --sidebar-ring: oklch(0.64 0.17 36.44);
+  --font-sans: Inter, sans-serif;
+  --font-serif: Source Serif 4, serif;
+  --font-mono: JetBrains Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.1;
+  --shadow-blur: 3px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-xs: 0px 1px 3px 0px hsl(0 0% 0% / 0.05);
+  --shadow-sm:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.1), 0px 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.1), 0px 1px 2px -1px hsl(0 0% 0% / 0.1);
+  --shadow-md:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.1), 0px 2px 4px -1px hsl(0 0% 0% / 0.1);
+  --shadow-lg:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.1), 0px 4px 6px -1px hsl(0 0% 0% / 0.1);
+  --shadow-xl:
+    0px 1px 3px 0px hsl(0 0% 0% / 0.1), 0px 8px 10px -1px hsl(0 0% 0% / 0.1);
+  --shadow-2xl: 0px 1px 3px 0px hsl(0 0% 0% / 0.25);
+}
diff --git a/web/src/styles/themes/twitch.css b/web/src/styles/themes/twitch.css
new file mode 100644
index 0000000..48d4e4a
--- /dev/null
+++ b/web/src/styles/themes/twitch.css
@@ -0,0 +1,128 @@
+/* Twitch Theme - Official Twitch brand colors */
+
+/* Light mode */
+[data-theme="twitch-light"] {
+  --background: oklch(0.97 0.002 0);
+  --foreground: oklch(0.14 0.005 285);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0.14 0.005 285);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.14 0.005 285);
+  --primary: oklch(0.54 0.27 290);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.94 0.015 290);
+  --secondary-foreground: oklch(0.2 0.01 285);
+  --muted: oklch(0.95 0.005 0);
+  --muted-foreground: oklch(0.5 0.01 285);
+  --accent: oklch(0.48 0.26 290);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.55 0.22 25);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.9 0.005 0);
+  --input: oklch(0.92 0.005 0);
+  --ring: oklch(0.54 0.27 290);
+  --chart-1: oklch(0.54 0.27 290);
+  --chart-2: oklch(0.48 0.26 290);
+  --chart-3: oklch(0.65 0.18 200);
+  --chart-4: oklch(0.7 0.15 150);
+  --chart-5: oklch(0.6 0.2 340);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.98 0.002 0);
+  --sidebar-foreground: oklch(0.14 0.005 285);
+  --sidebar-primary: oklch(0.54 0.27 290);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.94 0.015 290);
+  --sidebar-accent-foreground: oklch(0.2 0.01 285);
+  --sidebar-border: oklch(0.9 0.005 0);
+  --sidebar-ring: oklch(0.54 0.27 290);
+  --font-sans: Inter, ui-sans-serif, system-ui, sans-serif;
+  --font-serif: ui-serif, Georgia, serif;
+  --font-mono: ui-monospace, SFMono-Regular, "SF Mono", monospace;
+  --shadow-color: hsl(263 100% 10%);
+  --shadow-opacity: 0.08;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 2px;
+  --letter-spacing: -0.01em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 2px 0px hsl(263 100% 10% / 0.04);
+  --shadow-xs: 0 1px 3px 0px hsl(263 100% 10% / 0.06);
+  --shadow-sm:
+    0 2px 4px 0px hsl(263 100% 10% / 0.08),
+    0 1px 2px -1px hsl(263 100% 10% / 0.08);
+  --shadow:
+    0 2px 4px 0px hsl(263 100% 10% / 0.08),
+    0 1px 2px -1px hsl(263 100% 10% / 0.08);
+  --shadow-md:
+    0 4px 6px 0px hsl(263 100% 10% / 0.1),
+    0 2px 4px -1px hsl(263 100% 10% / 0.1);
+  --shadow-lg:
+    0 8px 10px 0px hsl(263 100% 10% / 0.1),
+    0 4px 6px -1px hsl(263 100% 10% / 0.1);
+  --shadow-xl:
+    0 16px 20px 0px hsl(263 100% 10% / 0.12),
+    0 8px 10px -1px hsl(263 100% 10% / 0.12);
+  --shadow-2xl: 0 24px 48px 0px hsl(263 100% 10% / 0.2);
+  --tracking-normal: -0.01em;
+}
+
+/* Dark mode - #0E0E10 is Twitch's signature dark background */
+[data-theme="twitch-dark"] {
+  --background: oklch(0.14 0.005 285);
+  --foreground: oklch(0.95 0.003 0);
+  --card: oklch(0.18 0.008 285);
+  --card-foreground: oklch(0.95 0.003 0);
+  --popover: oklch(0.2 0.008 285);
+  --popover-foreground: oklch(0.95 0.003 0);
+  --primary: oklch(0.54 0.27 290);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.22 0.008 285);
+  --secondary-foreground: oklch(0.9 0.003 0);
+  --muted: oklch(0.2 0.006 285);
+  --muted-foreground: oklch(0.65 0.01 285);
+  --accent: oklch(0.72 0.18 290);
+  --accent-foreground: oklch(0.14 0.005 285);
+  --destructive: oklch(0.6 0.2 25);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.28 0.01 285);
+  --input: oklch(0.22 0.008 285);
+  --ring: oklch(0.54 0.27 290);
+  --chart-1: oklch(0.54 0.27 290);
+  --chart-2: oklch(0.72 0.18 290);
+  --chart-3: oklch(0.7 0.16 200);
+  --chart-4: oklch(0.75 0.14 150);
+  --chart-5: oklch(0.65 0.18 340);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.12 0.005 285);
+  --sidebar-foreground: oklch(0.95 0.003 0);
+  --sidebar-primary: oklch(0.54 0.27 290);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.22 0.008 285);
+  --sidebar-accent-foreground: oklch(0.9 0.003 0);
+  --sidebar-border: oklch(0.28 0.01 285);
+  --sidebar-ring: oklch(0.54 0.27 290);
+  --font-sans: Inter, ui-sans-serif, system-ui, sans-serif;
+  --font-serif: ui-serif, Georgia, serif;
+  --font-mono: ui-monospace, SFMono-Regular, "SF Mono", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.4;
+  --shadow-blur: 4px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0;
+  --shadow-offset-y: 2px;
+  --letter-spacing: -0.01em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0 1px 2px 0px hsl(0 0% 0% / 0.2);
+  --shadow-xs: 0 1px 3px 0px hsl(0 0% 0% / 0.3);
+  --shadow-sm:
+    0 2px 4px 0px hsl(0 0% 0% / 0.4), 0 1px 2px -1px hsl(0 0% 0% / 0.4);
+  --shadow: 0 2px 4px 0px hsl(0 0% 0% / 0.4), 0 1px 2px -1px hsl(0 0% 0% / 0.4);
+  --shadow-md:
+    0 4px 6px 0px hsl(0 0% 0% / 0.4), 0 2px 4px -1px hsl(0 0% 0% / 0.4);
+  --shadow-lg:
+    0 8px 10px 0px hsl(0 0% 0% / 0.4), 0 4px 6px -1px hsl(0 0% 0% / 0.4);
+  --shadow-xl:
+    0 16px 20px 0px hsl(0 0% 0% / 0.4), 0 8px 10px -1px hsl(0 0% 0% / 0.4);
+  --shadow-2xl: 0 24px 48px 0px hsl(0 0% 0% / 0.5);
+}
diff --git a/web/src/styles/themes/twitter.css b/web/src/styles/themes/twitter.css
new file mode 100644
index 0000000..42b7fa8
--- /dev/null
+++ b/web/src/styles/themes/twitter.css
@@ -0,0 +1,134 @@
+/* Twitter Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="twitter-light"] {
+  --background: oklch(1 0 0);
+  --foreground: oklch(0.19 0.01 248.51);
+  --card: oklch(0.98 0 197.14);
+  --card-foreground: oklch(0.19 0.01 248.51);
+  --popover: oklch(1 0 0);
+  --popover-foreground: oklch(0.19 0.01 248.51);
+  --primary: oklch(0.67 0.16 245);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.19 0.01 248.51);
+  --secondary-foreground: oklch(1 0 0);
+  --muted: oklch(0.92 0 286.37);
+  --muted-foreground: oklch(0.19 0.01 248.51);
+  --accent: oklch(0.94 0.02 250.85);
+  --accent-foreground: oklch(0.67 0.16 245);
+  --destructive: oklch(0.62 0.24 25.77);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.93 0.01 231.66);
+  --input: oklch(0.98 0 228.78);
+  --ring: oklch(0.68 0.16 243.35);
+  --chart-1: oklch(0.67 0.16 245);
+  --chart-2: oklch(0.69 0.16 160.35);
+  --chart-3: oklch(0.82 0.16 82.53);
+  --chart-4: oklch(0.71 0.18 151.71);
+  --chart-5: oklch(0.59 0.22 10.58);
+  --radius: 1.3rem;
+  --sidebar: oklch(0.98 0 197.14);
+  --sidebar-foreground: oklch(0.19 0.01 248.51);
+  --sidebar-primary: oklch(0.67 0.16 245);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.94 0.02 250.85);
+  --sidebar-accent-foreground: oklch(0.67 0.16 245);
+  --sidebar-border: oklch(0.93 0.01 238.52);
+  --sidebar-ring: oklch(0.68 0.16 243.35);
+  --font-sans: Open Sans, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: Menlo, monospace;
+  --shadow-color: rgba(29, 161, 242, 0.15);
+  --shadow-opacity: 0;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-xs: 0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-sm:
+    0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0),
+    0px 1px 2px -1px hsl(202.82 89.12% 53.14% / 0);
+  --shadow:
+    0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0),
+    0px 1px 2px -1px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-md:
+    0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0),
+    0px 2px 4px -1px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-lg:
+    0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0),
+    0px 4px 6px -1px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-xl:
+    0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0),
+    0px 8px 10px -1px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-2xl: 0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="twitter-dark"] {
+  --background: oklch(0 0 0);
+  --foreground: oklch(0.93 0 228.79);
+  --card: oklch(0.21 0.01 274.53);
+  --card-foreground: oklch(0.89 0 0);
+  --popover: oklch(0 0 0);
+  --popover-foreground: oklch(0.93 0 228.79);
+  --primary: oklch(0.67 0.16 245.01);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.96 0 219.53);
+  --secondary-foreground: oklch(0.19 0.01 248.51);
+  --muted: oklch(0.21 0 0);
+  --muted-foreground: oklch(0.56 0.01 247.97);
+  --accent: oklch(0.19 0.03 242.55);
+  --accent-foreground: oklch(0.67 0.16 245.01);
+  --destructive: oklch(0.62 0.24 25.77);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.27 0 248);
+  --input: oklch(0.3 0.03 244.82);
+  --ring: oklch(0.68 0.16 243.35);
+  --chart-1: oklch(0.67 0.16 245);
+  --chart-2: oklch(0.69 0.16 160.35);
+  --chart-3: oklch(0.82 0.16 82.53);
+  --chart-4: oklch(0.71 0.18 151.71);
+  --chart-5: oklch(0.59 0.22 10.58);
+  --radius: 1.3rem;
+  --sidebar: oklch(0.21 0.01 274.53);
+  --sidebar-foreground: oklch(0.89 0 0);
+  --sidebar-primary: oklch(0.68 0.16 243.35);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.19 0.03 242.55);
+  --sidebar-accent-foreground: oklch(0.67 0.16 245.01);
+  --sidebar-border: oklch(0.38 0.02 240.59);
+  --sidebar-ring: oklch(0.68 0.16 243.35);
+  --font-sans: Open Sans, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: Menlo, monospace;
+  --shadow-color: rgba(29, 161, 242, 0.25);
+  --shadow-opacity: 0;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-xs: 0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-sm:
+    0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0),
+    0px 1px 2px -1px hsl(202.82 89.12% 53.14% / 0);
+  --shadow:
+    0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0),
+    0px 1px 2px -1px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-md:
+    0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0),
+    0px 2px 4px -1px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-lg:
+    0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0),
+    0px 4px 6px -1px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-xl:
+    0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0),
+    0px 8px 10px -1px hsl(202.82 89.12% 53.14% / 0);
+  --shadow-2xl: 0px 2px 0px 0px hsl(202.82 89.12% 53.14% / 0);
+}
diff --git a/web/src/styles/themes/vercel.css b/web/src/styles/themes/vercel.css
new file mode 100644
index 0000000..8d18b3c
--- /dev/null
+++ b/web/src/styles/themes/vercel.css
@@ -0,0 +1,124 @@
+/* Vercel Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="vercel-light"] {
+  --background: oklch(0.99 0 0);
+  --foreground: oklch(0 0 0);
+  --card: oklch(1 0 0);
+  --card-foreground: oklch(0 0 0);
+  --popover: oklch(0.99 0 0);
+  --popover-foreground: oklch(0 0 0);
+  --primary: oklch(0 0 0);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.94 0 0);
+  --secondary-foreground: oklch(0 0 0);
+  --muted: oklch(0.97 0 0);
+  --muted-foreground: oklch(0.44 0 0);
+  --accent: oklch(0.94 0 0);
+  --accent-foreground: oklch(0 0 0);
+  --destructive: oklch(0.63 0.19 23.03);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.92 0 0);
+  --input: oklch(0.94 0 0);
+  --ring: oklch(0 0 0);
+  --chart-1: oklch(0.81 0.17 75.35);
+  --chart-2: oklch(0.55 0.22 264.53);
+  --chart-3: oklch(0.72 0 0);
+  --chart-4: oklch(0.92 0 0);
+  --chart-5: oklch(0.56 0 0);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.99 0 0);
+  --sidebar-foreground: oklch(0 0 0);
+  --sidebar-primary: oklch(0 0 0);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.94 0 0);
+  --sidebar-accent-foreground: oklch(0 0 0);
+  --sidebar-border: oklch(0.94 0 0);
+  --sidebar-ring: oklch(0 0 0);
+  --font-sans: Geist, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: Geist Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.18;
+  --shadow-blur: 2px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 1px 2px 0px hsl(0 0% 0% / 0.09);
+  --shadow-xs: 0px 1px 2px 0px hsl(0 0% 0% / 0.09);
+  --shadow-sm:
+    0px 1px 2px 0px hsl(0 0% 0% / 0.18), 0px 1px 2px -1px hsl(0 0% 0% / 0.18);
+  --shadow:
+    0px 1px 2px 0px hsl(0 0% 0% / 0.18), 0px 1px 2px -1px hsl(0 0% 0% / 0.18);
+  --shadow-md:
+    0px 1px 2px 0px hsl(0 0% 0% / 0.18), 0px 2px 4px -1px hsl(0 0% 0% / 0.18);
+  --shadow-lg:
+    0px 1px 2px 0px hsl(0 0% 0% / 0.18), 0px 4px 6px -1px hsl(0 0% 0% / 0.18);
+  --shadow-xl:
+    0px 1px 2px 0px hsl(0 0% 0% / 0.18), 0px 8px 10px -1px hsl(0 0% 0% / 0.18);
+  --shadow-2xl: 0px 1px 2px 0px hsl(0 0% 0% / 0.45);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="vercel-dark"] {
+  --background: oklch(0 0 0);
+  --foreground: oklch(1 0 0);
+  --card: oklch(0.14 0 0);
+  --card-foreground: oklch(1 0 0);
+  --popover: oklch(0.18 0 0);
+  --popover-foreground: oklch(1 0 0);
+  --primary: oklch(1 0 0);
+  --primary-foreground: oklch(0 0 0);
+  --secondary: oklch(0.25 0 0);
+  --secondary-foreground: oklch(1 0 0);
+  --muted: oklch(0.23 0 0);
+  --muted-foreground: oklch(0.72 0 0);
+  --accent: oklch(0.32 0 0);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.69 0.2 23.91);
+  --destructive-foreground: oklch(0 0 0);
+  --border: oklch(0.26 0 0);
+  --input: oklch(0.32 0 0);
+  --ring: oklch(0.72 0 0);
+  --chart-1: oklch(0.81 0.17 75.35);
+  --chart-2: oklch(0.58 0.21 260.84);
+  --chart-3: oklch(0.56 0 0);
+  --chart-4: oklch(0.44 0 0);
+  --chart-5: oklch(0.92 0 0);
+  --radius: 0.5rem;
+  --sidebar: oklch(0.18 0 0);
+  --sidebar-foreground: oklch(1 0 0);
+  --sidebar-primary: oklch(1 0 0);
+  --sidebar-primary-foreground: oklch(0 0 0);
+  --sidebar-accent: oklch(0.32 0 0);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0.32 0 0);
+  --sidebar-ring: oklch(0.72 0 0);
+  --font-sans: Geist, sans-serif;
+  --font-serif: Georgia, serif;
+  --font-mono: Geist Mono, monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 0.18;
+  --shadow-blur: 2px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 0px;
+  --shadow-offset-y: 1px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 0px 1px 2px 0px hsl(0 0% 0% / 0.09);
+  --shadow-xs: 0px 1px 2px 0px hsl(0 0% 0% / 0.09);
+  --shadow-sm:
+    0px 1px 2px 0px hsl(0 0% 0% / 0.18), 0px 1px 2px -1px hsl(0 0% 0% / 0.18);
+  --shadow:
+    0px 1px 2px 0px hsl(0 0% 0% / 0.18), 0px 1px 2px -1px hsl(0 0% 0% / 0.18);
+  --shadow-md:
+    0px 1px 2px 0px hsl(0 0% 0% / 0.18), 0px 2px 4px -1px hsl(0 0% 0% / 0.18);
+  --shadow-lg:
+    0px 1px 2px 0px hsl(0 0% 0% / 0.18), 0px 4px 6px -1px hsl(0 0% 0% / 0.18);
+  --shadow-xl:
+    0px 1px 2px 0px hsl(0 0% 0% / 0.18), 0px 8px 10px -1px hsl(0 0% 0% / 0.18);
+  --shadow-2xl: 0px 1px 2px 0px hsl(0 0% 0% / 0.45);
+}
diff --git a/web/src/styles/themes/vintage-paper.css b/web/src/styles/themes/vintage-paper.css
new file mode 100644
index 0000000..b15d976
--- /dev/null
+++ b/web/src/styles/themes/vintage-paper.css
@@ -0,0 +1,134 @@
+/* Vintage Paper Theme - Auto-generated from tweakcn */
+
+/* Light mode */
+[data-theme="vintage-paper-light"] {
+  --background: oklch(0.96 0.02 90.24);
+  --foreground: oklch(0.38 0.02 64.34);
+  --card: oklch(0.99 0.01 87.47);
+  --card-foreground: oklch(0.38 0.02 64.34);
+  --popover: oklch(0.99 0.01 87.47);
+  --popover-foreground: oklch(0.38 0.02 64.34);
+  --primary: oklch(0.62 0.08 65.54);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.88 0.03 85.57);
+  --secondary-foreground: oklch(0.43 0.03 64.93);
+  --muted: oklch(0.92 0.02 83.06);
+  --muted-foreground: oklch(0.54 0.04 71.17);
+  --accent: oklch(0.83 0.04 88.81);
+  --accent-foreground: oklch(0.38 0.02 64.34);
+  --destructive: oklch(0.55 0.14 32.91);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.86 0.03 84.59);
+  --input: oklch(0.86 0.03 84.59);
+  --ring: oklch(0.62 0.08 65.54);
+  --chart-1: oklch(0.62 0.08 65.54);
+  --chart-2: oklch(0.56 0.06 68.58);
+  --chart-3: oklch(0.49 0.06 72.68);
+  --chart-4: oklch(0.68 0.06 64.78);
+  --chart-5: oklch(0.73 0.06 66.7);
+  --radius: 0.25rem;
+  --sidebar: oklch(0.92 0.02 83.06);
+  --sidebar-foreground: oklch(0.38 0.02 64.34);
+  --sidebar-primary: oklch(0.62 0.08 65.54);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.83 0.04 88.81);
+  --sidebar-accent-foreground: oklch(0.38 0.02 64.34);
+  --sidebar-border: oklch(0.86 0.03 84.59);
+  --sidebar-ring: oklch(0.62 0.08 65.54);
+  --font-sans: Libre Baskerville, serif;
+  --font-serif: Lora, serif;
+  --font-mono: IBM Plex Mono, monospace;
+  --shadow-color: hsl(28 13% 20%);
+  --shadow-opacity: 0.12;
+  --shadow-blur: 5px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 3px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 3px 5px 0px hsl(28 13% 20% / 0.06);
+  --shadow-xs: 2px 3px 5px 0px hsl(28 13% 20% / 0.06);
+  --shadow-sm:
+    2px 3px 5px 0px hsl(28 13% 20% / 0.12),
+    2px 1px 2px -1px hsl(28 13% 20% / 0.12);
+  --shadow:
+    2px 3px 5px 0px hsl(28 13% 20% / 0.12),
+    2px 1px 2px -1px hsl(28 13% 20% / 0.12);
+  --shadow-md:
+    2px 3px 5px 0px hsl(28 13% 20% / 0.12),
+    2px 2px 4px -1px hsl(28 13% 20% / 0.12);
+  --shadow-lg:
+    2px 3px 5px 0px hsl(28 13% 20% / 0.12),
+    2px 4px 6px -1px hsl(28 13% 20% / 0.12);
+  --shadow-xl:
+    2px 3px 5px 0px hsl(28 13% 20% / 0.12),
+    2px 8px 10px -1px hsl(28 13% 20% / 0.12);
+  --shadow-2xl: 2px 3px 5px 0px hsl(28 13% 20% / 0.3);
+  --tracking-normal: 0em;
+}
+
+/* Dark mode */
+[data-theme="vintage-paper-dark"] {
+  --background: oklch(0.27 0.01 57.65);
+  --foreground: oklch(0.92 0.02 83.06);
+  --card: oklch(0.32 0.02 59.06);
+  --card-foreground: oklch(0.92 0.02 83.06);
+  --popover: oklch(0.32 0.02 59.06);
+  --popover-foreground: oklch(0.92 0.02 83.06);
+  --primary: oklch(0.73 0.06 66.7);
+  --primary-foreground: oklch(0.27 0.01 57.65);
+  --secondary: oklch(0.38 0.02 57.13);
+  --secondary-foreground: oklch(0.92 0.02 83.06);
+  --muted: oklch(0.32 0.02 59.06);
+  --muted-foreground: oklch(0.8 0.02 82.11);
+  --accent: oklch(0.42 0.03 56.34);
+  --accent-foreground: oklch(0.92 0.02 83.06);
+  --destructive: oklch(0.55 0.14 32.91);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.38 0.02 57.13);
+  --input: oklch(0.38 0.02 57.13);
+  --ring: oklch(0.73 0.06 66.7);
+  --chart-1: oklch(0.73 0.06 66.7);
+  --chart-2: oklch(0.68 0.06 64.78);
+  --chart-3: oklch(0.62 0.08 65.54);
+  --chart-4: oklch(0.56 0.06 68.58);
+  --chart-5: oklch(0.49 0.06 72.68);
+  --radius: 0.25rem;
+  --sidebar: oklch(0.27 0.01 57.65);
+  --sidebar-foreground: oklch(0.92 0.02 83.06);
+  --sidebar-primary: oklch(0.73 0.06 66.7);
+  --sidebar-primary-foreground: oklch(0.27 0.01 57.65);
+  --sidebar-accent: oklch(0.42 0.03 56.34);
+  --sidebar-accent-foreground: oklch(0.92 0.02 83.06);
+  --sidebar-border: oklch(0.38 0.02 57.13);
+  --sidebar-ring: oklch(0.73 0.06 66.7);
+  --font-sans: Libre Baskerville, serif;
+  --font-serif: Lora, serif;
+  --font-mono: IBM Plex Mono, monospace;
+  --shadow-color: hsl(28 13% 20%);
+  --shadow-opacity: 0.12;
+  --shadow-blur: 5px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 3px;
+  --letter-spacing: 0em;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 3px 5px 0px hsl(28 13% 20% / 0.06);
+  --shadow-xs: 2px 3px 5px 0px hsl(28 13% 20% / 0.06);
+  --shadow-sm:
+    2px 3px 5px 0px hsl(28 13% 20% / 0.12),
+    2px 1px 2px -1px hsl(28 13% 20% / 0.12);
+  --shadow:
+    2px 3px 5px 0px hsl(28 13% 20% / 0.12),
+    2px 1px 2px -1px hsl(28 13% 20% / 0.12);
+  --shadow-md:
+    2px 3px 5px 0px hsl(28 13% 20% / 0.12),
+    2px 2px 4px -1px hsl(28 13% 20% / 0.12);
+  --shadow-lg:
+    2px 3px 5px 0px hsl(28 13% 20% / 0.12),
+    2px 4px 6px -1px hsl(28 13% 20% / 0.12);
+  --shadow-xl:
+    2px 3px 5px 0px hsl(28 13% 20% / 0.12),
+    2px 8px 10px -1px hsl(28 13% 20% / 0.12);
+  --shadow-2xl: 2px 3px 5px 0px hsl(28 13% 20% / 0.3);
+}
diff --git a/web/src/styles/themes/windows98.css b/web/src/styles/themes/windows98.css
new file mode 100644
index 0000000..39337a0
--- /dev/null
+++ b/web/src/styles/themes/windows98.css
@@ -0,0 +1,123 @@
+/* Windows 98 Theme */
+
+/* Light mode */
+[data-theme="windows98-light"] {
+  --background: oklch(0.5431 0.0927 194.7689);
+  --foreground: oklch(0 0 0);
+  --card: oklch(0.8078 0 0);
+  --card-foreground: oklch(0 0 0);
+  --popover: oklch(0.8078 0 0);
+  --popover-foreground: oklch(0 0 0);
+  --primary: oklch(0.2711 0.1879 264.052);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.5999 0 0);
+  --secondary-foreground: oklch(0 0 0);
+  --muted: oklch(0.5999 0 0);
+  --muted-foreground: oklch(0 0 0);
+  --accent: oklch(0.2711 0.1879 264.052);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.628 0.2577 29.2339);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.5999 0 0);
+  --input: oklch(1 0 0);
+  --ring: oklch(0.2711 0.1879 264.052);
+  --chart-1: oklch(0.2711 0.1879 264.052);
+  --chart-2: oklch(0.5198 0.1769 142.4953);
+  --chart-3: oklch(0.4209 0.1935 328.3634);
+  --chart-4: oklch(0.628 0.2577 29.2339);
+  --chart-5: oklch(0.5807 0.1266 109.7692);
+  --radius: 0px;
+  --sidebar: oklch(0.8078 0 0);
+  --sidebar-foreground: oklch(0 0 0);
+  --sidebar-primary: oklch(0.2711 0.1879 264.052);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.2711 0.1879 264.052);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0.5999 0 0);
+  --sidebar-ring: oklch(0.2711 0.1879 264.052);
+  --font-sans: "Pixelify Sans", sans-serif;
+  --font-serif: "Pixelify Sans", serif;
+  --font-mono: "Pixelify Sans", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 1;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0.025rem;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 2px 0px 0px hsl(0 0% 0% / 0.5);
+  --shadow-xs: 2px 2px 0px 0px hsl(0 0% 0% / 0.5);
+  --shadow-sm:
+    2px 2px 0px 0px hsl(0 0% 0% / 1), 2px 1px 2px -1px hsl(0 0% 0% / 1);
+  --shadow: 2px 2px 0px 0px hsl(0 0% 0% / 1), 2px 1px 2px -1px hsl(0 0% 0% / 1);
+  --shadow-md:
+    2px 2px 0px 0px hsl(0 0% 0% / 1), 2px 2px 4px -1px hsl(0 0% 0% / 1);
+  --shadow-lg:
+    2px 2px 0px 0px hsl(0 0% 0% / 1), 2px 4px 6px -1px hsl(0 0% 0% / 1);
+  --shadow-xl:
+    2px 2px 0px 0px hsl(0 0% 0% / 1), 2px 8px 10px -1px hsl(0 0% 0% / 1);
+  --shadow-2xl: 2px 2px 0px 0px hsl(0 0% 0% / 2.5);
+  --tracking-normal: 0.025rem;
+}
+
+/* Dark mode */
+[data-theme="windows98-dark"] {
+  --background: oklch(0.5431 0.0927 194.7689);
+  --foreground: oklch(0 0 0);
+  --card: oklch(0.8078 0 0);
+  --card-foreground: oklch(0 0 0);
+  --popover: oklch(0.8078 0 0);
+  --popover-foreground: oklch(0 0 0);
+  --primary: oklch(0.2711 0.1879 264.052);
+  --primary-foreground: oklch(1 0 0);
+  --secondary: oklch(0.5999 0 0);
+  --secondary-foreground: oklch(0 0 0);
+  --muted: oklch(0.5999 0 0);
+  --muted-foreground: oklch(1 0 0);
+  --accent: oklch(0.2711 0.1879 264.052);
+  --accent-foreground: oklch(1 0 0);
+  --destructive: oklch(0.628 0.2577 29.2339);
+  --destructive-foreground: oklch(1 0 0);
+  --border: oklch(0.5999 0 0);
+  --input: oklch(1 0 0);
+  --ring: oklch(0.2711 0.1879 264.052);
+  --chart-1: oklch(0.2711 0.1879 264.052);
+  --chart-2: oklch(0.5431 0.0927 194.7689);
+  --chart-3: oklch(0.4209 0.1935 328.3634);
+  --chart-4: oklch(0.628 0.2577 29.2339);
+  --chart-5: oklch(0.8664 0.2948 142.4953);
+  --radius: 0px;
+  --sidebar: oklch(0.8078 0 0);
+  --sidebar-foreground: oklch(0 0 0);
+  --sidebar-primary: oklch(0.2711 0.1879 264.052);
+  --sidebar-primary-foreground: oklch(1 0 0);
+  --sidebar-accent: oklch(0.2711 0.1879 264.052);
+  --sidebar-accent-foreground: oklch(1 0 0);
+  --sidebar-border: oklch(0.5999 0 0);
+  --sidebar-ring: oklch(0.2711 0.1879 264.052);
+  --font-sans: "Pixelify Sans", sans-serif;
+  --font-serif: "Pixelify Sans", serif;
+  --font-mono: "Pixelify Sans", monospace;
+  --shadow-color: hsl(0 0% 0%);
+  --shadow-opacity: 1;
+  --shadow-blur: 0px;
+  --shadow-spread: 0px;
+  --shadow-offset-x: 2px;
+  --shadow-offset-y: 2px;
+  --letter-spacing: 0.025rem;
+  --spacing: 0.25rem;
+  --shadow-2xs: 2px 2px 0px 0px hsl(0 0% 0% / 0.5);
+  --shadow-xs: 2px 2px 0px 0px hsl(0 0% 0% / 0.5);
+  --shadow-sm:
+    2px 2px 0px 0px hsl(0 0% 0% / 1), 2px 1px 2px -1px hsl(0 0% 0% / 1);
+  --shadow: 2px 2px 0px 0px hsl(0 0% 0% / 1), 2px 1px 2px -1px hsl(0 0% 0% / 1);
+  --shadow-md:
+    2px 2px 0px 0px hsl(0 0% 0% / 1), 2px 2px 4px -1px hsl(0 0% 0% / 1);
+  --shadow-lg:
+    2px 2px 0px 0px hsl(0 0% 0% / 1), 2px 4px 6px -1px hsl(0 0% 0% / 1);
+  --shadow-xl:
+    2px 2px 0px 0px hsl(0 0% 0% / 1), 2px 8px 10px -1px hsl(0 0% 0% / 1);
+  --shadow-2xl: 2px 2px 0px 0px hsl(0 0% 0% / 2.5);
+  --tracking-normal: 0.025rem;
+}
