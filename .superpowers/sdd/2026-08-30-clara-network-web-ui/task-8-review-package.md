BASE 157c32a

0f711eb feat(web): overview + transactions dashboards

=== STAT ===
 web/package-lock.json                   | 519 ++++++++++++++++++++++++++++++++
 web/package.json                        |   1 +
 web/src/app/(app)/ops/page.tsx          |  21 +-
 web/src/app/(app)/overview/page.tsx     |  20 +-
 web/src/app/(app)/transactions/page.tsx |  46 +++
 web/src/components/cards/stat-card.tsx  |  10 +
 web/src/lib/adminapi.ts                 |  13 +
 web/src/lib/money/minor.test.ts         |  10 +
 web/src/lib/money/minor.ts              |   2 +-
 9 files changed, 638 insertions(+), 4 deletions(-)
diff --git a/web/package-lock.json b/web/package-lock.json
index f902622..553941b 100644
--- a/web/package-lock.json
+++ b/web/package-lock.json
@@ -24,20 +24,21 @@
         "tw-animate-css": "^1.4.0"
       },
       "devDependencies": {
         "@tailwindcss/postcss": "^4",
         "@types/node": "^20",
         "@types/react": "^19",
         "@types/react-dom": "^19",
         "eslint": "^9",
         "eslint-config-next": "16.3.3",
         "tailwindcss": "^4",
+        "tsx": "^4.23.13",
         "typescript": "^5"
       }
     },
     "node_modules/@alloc/quick-lru": {
       "version": "5.2.0",
       "dev": true,
       "license": "MIT",
       "engines": {
         "node": ">=10"
       },
@@ -747,20 +748,462 @@
       "version": "1.2.1",
       "resolved": "https://registry.npmjs.org/@emnapi/wasi-threads/-/wasi-threads-1.2.1.tgz",
       "integrity": "sha512-uTII7OYF+/Mes/MrcIOYp5yOtSMLBWSIoLPpcgwipoiKbli6k322tcoFsxoIIxPDqW01SQGAgko4EzZi2BNv2w==",
       "dev": true,
       "license": "MIT",
       "optional": true,
       "dependencies": {
         "tslib": "^2.4.0"
       }
     },
+    "node_modules/@esbuild/aix-ppc64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/aix-ppc64/-/aix-ppc64-0.28.2.tgz",
+      "integrity": "sha512-XExcO+dvLKvVtNTibSTBej1NCAbaGhWn9Ww1ZPx80qsahhPFe/8jgWP0IchNe0F3HwkU7n8ejhH8bjonqht8mQ==",
+      "cpu": [
+        "ppc64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "aix"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/android-arm": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/android-arm/-/android-arm-0.28.2.tgz",
+      "integrity": "sha512-kXXoiPVVGQcnIYGOeaovwOURpniDBpSq4A03qkQ+BMQqtGG6HYap3xne9C1O1yo4TR3qxlCX5IqqmX6fFo2Lqg==",
+      "cpu": [
+        "arm"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "android"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/android-arm64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/android-arm64/-/android-arm64-0.28.2.tgz",
+      "integrity": "sha512-5YfKeeI8qWfBZIX+u2xZC3Zlb3Os/gLS2sbEKM+I4ZOcsWmHS2WLysCcQZDAFRslDUU5Oiq44gf6PYN1vGwG5A==",
+      "cpu": [
+        "arm64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "android"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/android-x64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/android-x64/-/android-x64-0.28.2.tgz",
+      "integrity": "sha512-O387ite7SzUyCcy3JQX4P4bLtEA7bLLkx+esve5JHnyYfNTxcVpXZo9jhdB0lTKN44gztELTdU7nS8Nr16Fs1Q==",
+      "cpu": [
+        "x64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "android"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/darwin-arm64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/darwin-arm64/-/darwin-arm64-0.28.2.tgz",
+      "integrity": "sha512-n4KqkOQrraxHJcgjM1RvwbigfQKIKJVpM7xp+KsxiyUSrRdIXnt73VhrPAx0fV44hgfmIVKjxMN9J1t5jySVkw==",
+      "cpu": [
+        "arm64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "darwin"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/darwin-x64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/darwin-x64/-/darwin-x64-0.28.2.tgz",
+      "integrity": "sha512-uq6suIWYP37qzGddBKPw5QEQPi6HiLGsO7UmkpfyaYNQ3D+rN6w6WfwH+nuqcGXWvawGwxOEroO4YGnFh95azw==",
+      "cpu": [
+        "x64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "darwin"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/freebsd-arm64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/freebsd-arm64/-/freebsd-arm64-0.28.2.tgz",
+      "integrity": "sha512-n+I0BTSRIoy+d6RPKnEVwql5UwBJolytvY4mAOIEJorKlqgPII8ix6slVVrfZ5Tnj7glIZvloylbB/EJPMWEXw==",
+      "cpu": [
+        "arm64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "freebsd"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/freebsd-x64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/freebsd-x64/-/freebsd-x64-0.28.2.tgz",
+      "integrity": "sha512-78XJTJkvPs0kz2w61301PJjXl4g7q3JqiYMZ/M/yVI73EHBrCRTgkhu9oqG7vPqq+a/yadEW8aD+agKlk5xrmg==",
+      "cpu": [
+        "x64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "freebsd"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/linux-arm": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/linux-arm/-/linux-arm-0.28.2.tgz",
+      "integrity": "sha512-XlDnu2q5yoqems+xay6wSAcg9DDD7K9RLKZEBOMZm3ckNpJBvOX20tSfby8KfrrhINDyv9V2YVZKY/SpoGJI8w==",
+      "cpu": [
+        "arm"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "linux"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/linux-arm64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/linux-arm64/-/linux-arm64-0.28.2.tgz",
+      "integrity": "sha512-pW4AC0P3it8c7do9MVM4p51FzHzdM/TZrerurgRcHJ2WTa1VQ1CIq18xncfpBJw4ojkiZZrKW2yIBWBP92j6Ug==",
+      "cpu": [
+        "arm64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "linux"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/linux-ia32": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/linux-ia32/-/linux-ia32-0.28.2.tgz",
+      "integrity": "sha512-CYbnj78HsIeA+DhgUKgFCfvNsTHFhMMrinUrMZpDXJXKN8T3XViTZ/+wtHeVxEWY8ewSzTFN+nRmSwO2tZaLUQ==",
+      "cpu": [
+        "ia32"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "linux"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/linux-loong64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/linux-loong64/-/linux-loong64-0.28.2.tgz",
+      "integrity": "sha512-buwkd8nsph4R+ajRvw0qM5Hja/TXQow3ptzWO2EbG/cqcIkHloRrdlBtQlshyYGTNFvfkfJ5tpPLVkY4DtsPfQ==",
+      "cpu": [
+        "loong64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "linux"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/linux-mips64el": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/linux-mips64el/-/linux-mips64el-0.28.2.tgz",
+      "integrity": "sha512-ZVykbDyk7519VwiNb9Lcj9m8XM6v5V9uKPvrEMkkEedVewf+0itkhahp4HDpgERXhwLRpWFypsGbG/J8s0QjJA==",
+      "cpu": [
+        "mips64el"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "linux"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/linux-ppc64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/linux-ppc64/-/linux-ppc64-0.28.2.tgz",
+      "integrity": "sha512-CAXl+Dtd9UUuJd8pKKdwh6MLm3MUMiqMPmhZ3tTSXPqfyQ3vDl6R5hZdZ/kYojK4ofXtdfSv1tFq8XzWx3heNQ==",
+      "cpu": [
+        "ppc64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "linux"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/linux-riscv64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/linux-riscv64/-/linux-riscv64-0.28.2.tgz",
+      "integrity": "sha512-GeXCej4IQtU1B+QlDV8W/RRvbzI3O/Stss+/bCXv4lZls5WGRtu2a+3JkA3i4qIUlMXpcHebWpF8AkJhATowuA==",
+      "cpu": [
+        "riscv64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "linux"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/linux-s390x": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/linux-s390x/-/linux-s390x-0.28.2.tgz",
+      "integrity": "sha512-3H1weTYZPxt/WOhByszQZybS9w5lKzUn1FDMsgEChbHWQwHYQQRfBxgCcZvPhjHfKyJjIievvMmEUawJrdY9Dg==",
+      "cpu": [
+        "s390x"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "linux"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/linux-x64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/linux-x64/-/linux-x64-0.28.2.tgz",
+      "integrity": "sha512-4xTZr1FUmSoQW4XIWmit3tzQrUTZM+N3P0XV8xROKYF50XfI7xeO90+1bZvNwxIufQ9hDQVRJH5YhgPVF8A/HQ==",
+      "cpu": [
+        "x64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "linux"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/netbsd-arm64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/netbsd-arm64/-/netbsd-arm64-0.28.2.tgz",
+      "integrity": "sha512-sSATRjPeDBg3pdgHoQfoYBob11Kk1FGa9lui5RIHZCoCkJa9QKlvl3/vKz2usCmYYjs7ymJR/2Nnsqe+Hjt5nw==",
+      "cpu": [
+        "arm64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "netbsd"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/netbsd-x64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/netbsd-x64/-/netbsd-x64-0.28.2.tgz",
+      "integrity": "sha512-lqnzCV+mM0gIADaKihiCg6ifgfU2L3h5E33rNQBN1Y4MaVGnzryzmvvf7UHxprpQdE8hpqLolJ9Rl+SkIRDpyw==",
+      "cpu": [
+        "x64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "netbsd"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/openbsd-arm64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/openbsd-arm64/-/openbsd-arm64-0.28.2.tgz",
+      "integrity": "sha512-AL2qJILH7lNjrDmCQDvdxMfAUIv8KMNZOvrwAQ8i8//ntL9FflhOyMJ8OZSMBb8/AWXe3/5v5S20y3zCoZWKoQ==",
+      "cpu": [
+        "arm64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "openbsd"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/openbsd-x64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/openbsd-x64/-/openbsd-x64-0.28.2.tgz",
+      "integrity": "sha512-QtiuPytchRyC4rwUKhexJdQKvDuZ6hWloi3igqPQNUJCS1/v9EiO3UTOXR6A3FoMo4fnAKbWJdqaIwhOzh8qEw==",
+      "cpu": [
+        "x64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "openbsd"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/openharmony-arm64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/openharmony-arm64/-/openharmony-arm64-0.28.2.tgz",
+      "integrity": "sha512-WkhYDmpTjLvGlScA1rwjRUmhl4k8oXR3cIbtqWmELgU/dFeHHlEllxDvdWcNJV9rbzCexB5vz8gtNewWLgCT7Q==",
+      "cpu": [
+        "arm64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "openharmony"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/sunos-x64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/sunos-x64/-/sunos-x64-0.28.2.tgz",
+      "integrity": "sha512-GPMSkTOtMnv2U2F8gxe4Io6qmVs+YKyp832Etqqxr0hFngmXQ3rzwytelm3GIn7T4VviRUlf3sOgBOiTdvaf7g==",
+      "cpu": [
+        "x64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "sunos"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/win32-arm64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/win32-arm64/-/win32-arm64-0.28.2.tgz",
+      "integrity": "sha512-PIhhEkE9uPBleRBrQEJpUn7MBnibZzbGzYWPmY3x+YoVg/95zbjB4CxPPOQ8l5tYYM4mMaCthF8/1DIfBQQyWQ==",
+      "cpu": [
+        "arm64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "win32"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/win32-ia32": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/win32-ia32/-/win32-ia32-0.28.2.tgz",
+      "integrity": "sha512-YmJbfTlvU7Sdn9BB+4PRES4oB6pxgS37MAONj+hBr/cpXS1aBPKXxNnDbu+QCWPj0o9dgyxeq79g6c5P8KeuYA==",
+      "cpu": [
+        "ia32"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "win32"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
+    "node_modules/@esbuild/win32-x64": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/@esbuild/win32-x64/-/win32-x64-0.28.2.tgz",
+      "integrity": "sha512-5ebpxr3nWMzrL/rnUI755Jkuee0bHL/Gq0WTF9lvcpv73wAp5eu8MfBUgWK9bhWvZjj7yX8etf/8tI8Ney695g==",
+      "cpu": [
+        "x64"
+      ],
+      "dev": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "win32"
+      ],
+      "engines": {
+        "node": ">=18"
+      }
+    },
     "node_modules/@eslint-community/eslint-utils": {
       "version": "4.10.1",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "eslint-visitor-keys": "^3.4.3"
       },
       "engines": {
         "node": "^12.22.0 || ^14.17.0 || >=16.0.0"
       },
@@ -4360,20 +4803,62 @@
         "is-date-object": "^1.1.0",
         "is-symbol": "^1.1.1"
       },
       "engines": {
         "node": ">= 0.4"
       },
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
+    "node_modules/esbuild": {
+      "version": "0.28.2",
+      "resolved": "https://registry.npmjs.org/esbuild/-/esbuild-0.28.2.tgz",
+      "integrity": "sha512-HKVLS8dvII+xoKW9kmqxbRKrnWEXfJJr/FZhhJmiqIB0e053QNYFqOBouTMO/k5sID4MvCiUCvv8b9M4h32wIA==",
+      "dev": true,
+      "hasInstallScript": true,
+      "license": "MIT",
+      "bin": {
+        "esbuild": "bin/esbuild"
+      },
+      "engines": {
+        "node": ">=18"
+      },
+      "optionalDependencies": {
+        "@esbuild/aix-ppc64": "0.28.2",
+        "@esbuild/android-arm": "0.28.2",
+        "@esbuild/android-arm64": "0.28.2",
+        "@esbuild/android-x64": "0.28.2",
+        "@esbuild/darwin-arm64": "0.28.2",
+        "@esbuild/darwin-x64": "0.28.2",
+        "@esbuild/freebsd-arm64": "0.28.2",
+        "@esbuild/freebsd-x64": "0.28.2",
+        "@esbuild/linux-arm": "0.28.2",
+        "@esbuild/linux-arm64": "0.28.2",
+        "@esbuild/linux-ia32": "0.28.2",
+        "@esbuild/linux-loong64": "0.28.2",
+        "@esbuild/linux-mips64el": "0.28.2",
+        "@esbuild/linux-ppc64": "0.28.2",
+        "@esbuild/linux-riscv64": "0.28.2",
+        "@esbuild/linux-s390x": "0.28.2",
+        "@esbuild/linux-x64": "0.28.2",
+        "@esbuild/netbsd-arm64": "0.28.2",
+        "@esbuild/netbsd-x64": "0.28.2",
+        "@esbuild/openbsd-arm64": "0.28.2",
+        "@esbuild/openbsd-x64": "0.28.2",
+        "@esbuild/openharmony-arm64": "0.28.2",
+        "@esbuild/sunos-x64": "0.28.2",
+        "@esbuild/win32-arm64": "0.28.2",
+        "@esbuild/win32-ia32": "0.28.2",
+        "@esbuild/win32-x64": "0.28.2"
+      }
+    },
     "node_modules/escalade": {
       "version": "3.2.0",
       "license": "MIT",
       "engines": {
         "node": ">=6"
       }
     },
     "node_modules/escape-html": {
       "version": "1.0.3",
       "resolved": "https://registry.npmjs.org/escape-html/-/escape-html-1.0.3.tgz",
@@ -5079,20 +5564,35 @@
       "license": "MIT",
       "dependencies": {
         "graceful-fs": "^4.2.0",
         "jsonfile": "^6.0.1",
         "universalify": "^2.0.0"
       },
       "engines": {
         "node": ">=14.14"
       }
     },
+    "node_modules/fsevents": {
+      "version": "2.3.3",
+      "resolved": "https://registry.npmjs.org/fsevents/-/fsevents-2.3.3.tgz",
+      "integrity": "sha512-5xoDfX+fL7faATnagmWPpbFtwh/R77WmMMqqHGS65C3vvB0YHrgF+B1YmZ3441tMj5n63k0212XNoJwzlhffQw==",
+      "dev": true,
+      "hasInstallScript": true,
+      "license": "MIT",
+      "optional": true,
+      "os": [
+        "darwin"
+      ],
+      "engines": {
+        "node": "^8.16.0 || ^10.6.0 || >=11.0.0"
+      }
+    },
     "node_modules/function-bind": {
       "version": "1.1.2",
       "license": "MIT",
       "funding": {
         "url": "https://github.com/sponsors/ljharb"
       }
     },
     "node_modules/function.prototype.name": {
       "version": "1.2.0",
       "dev": true,
@@ -8741,20 +9241,39 @@
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
+    "node_modules/tsx": {
+      "version": "4.23.13",
+      "resolved": "https://registry.npmjs.org/tsx/-/tsx-4.23.13.tgz",
+      "integrity": "sha512-BL5MGkRln6aDYhb0xbQlEAGw743BaZYWdbWtdJOBriYJboKgUUYCadFp2/FpBBZquBC/ezNBn7wMMPx7FDZUDw==",
+      "dev": true,
+      "license": "MIT",
+      "dependencies": {
+        "esbuild": "~0.28.0"
+      },
+      "bin": {
+        "tsx": "dist/cli.mjs"
+      },
+      "engines": {
+        "node": ">=18.0.0"
+      },
+      "optionalDependencies": {
+        "fsevents": "~2.3.3"
+      }
+    },
     "node_modules/tw-animate-css": {
       "version": "1.4.0",
       "resolved": "https://registry.npmjs.org/tw-animate-css/-/tw-animate-css-1.4.0.tgz",
       "integrity": "sha512-7bziOlRqH0hJx80h/3mbicLW7o8qLsH5+RaLR2t+OHM3D0JlWGODQKQ4cxbK7WlvmUxpcj6Kgu6EKqjrGFe3QQ==",
       "license": "MIT",
       "funding": {
         "url": "https://github.com/sponsors/Wombosvideo"
       }
     },
     "node_modules/type-check": {
diff --git a/web/package.json b/web/package.json
index fff09c8..86e1609 100644
--- a/web/package.json
+++ b/web/package.json
@@ -25,13 +25,14 @@
     "tw-animate-css": "^1.4.0"
   },
   "devDependencies": {
     "@tailwindcss/postcss": "^4",
     "@types/node": "^20",
     "@types/react": "^19",
     "@types/react-dom": "^19",
     "eslint": "^9",
     "eslint-config-next": "16.3.3",
     "tailwindcss": "^4",
+    "tsx": "^4.23.13",
     "typescript": "^5"
   }
 }
diff --git a/web/src/app/(app)/ops/page.tsx b/web/src/app/(app)/ops/page.tsx
index 9f9c057..5496e4a 100644
--- a/web/src/app/(app)/ops/page.tsx
+++ b/web/src/app/(app)/ops/page.tsx
@@ -1,10 +1,27 @@
 import { notFound } from "next/navigation";
+import Link from "next/link";
 import { createServerClient } from "@/lib/supabase/server";
 import { roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { StatCard } from "@/components/cards/stat-card";
+import type { DashboardSummary } from "@/types/admin";
 
 export default async function OpsPage() {
   const supabase = await createServerClient();
   const { data } = await supabase.auth.getUser();
   if (roleFromAppMetadata(data.user?.app_metadata) !== "scheme_operator") notFound();
-  return <h1 className="text-2xl font-semibold">Operations</h1>;
-}
\ No newline at end of file
+  const d = await fetchAdmin<DashboardSummary>("/dashboard");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Operations</h1>
+      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
+        <StatCard title="Transactions today" value={d.transactions.toLocaleString()} hint="Authorizations via the switch" />
+        <StatCard title="Clearing records" value={d.clearingRecords.toLocaleString()} hint="Captured this settlement window" />
+        <StatCard title="Merchants onboarded" value={d.merchants.toLocaleString()} />
+      </div>
+      <p className="text-sm text-muted-foreground">
+        <Link href="/transactions" className="underline underline-offset-4 hover:text-foreground">View transaction log ÔåÆ</Link>
+      </p>
+    </div>
+  );
+}
diff --git a/web/src/app/(app)/overview/page.tsx b/web/src/app/(app)/overview/page.tsx
index 8d0ec73..69bbbe6 100644
--- a/web/src/app/(app)/overview/page.tsx
+++ b/web/src/app/(app)/overview/page.tsx
@@ -1,3 +1,21 @@
+// src/app/(app)/overview/page.tsx
+import { fetchAdmin } from "@/lib/adminapi";
+import { StatCard } from "@/components/cards/stat-card";
+import type { DashboardSummary } from "@/types/admin";
+
 export default async function OverviewPage() {
-  return <h1 className="text-2xl font-semibold">Overview</h1>;
+  const d = await fetchAdmin<DashboardSummary>("/dashboard");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Network overview</h1>
+      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
+        <StatCard title="Transactions" value={d.transactions.toLocaleString()} />
+        <StatCard title="Clearing records" value={d.clearingRecords.toLocaleString()} />
+        <StatCard title="Merchants" value={d.merchants.toLocaleString()} />
+        <StatCard title="Disputes" value={d.disputes.toLocaleString()} />
+        <StatCard title="Cards" value={d.cards.toLocaleString()} />
+        <StatCard title="Tokens" value={d.tokens.toLocaleString()} />
+      </div>
+    </div>
+  );
 }
\ No newline at end of file
diff --git a/web/src/app/(app)/transactions/page.tsx b/web/src/app/(app)/transactions/page.tsx
new file mode 100644
index 0000000..1b08d3d
--- /dev/null
+++ b/web/src/app/(app)/transactions/page.tsx
@@ -0,0 +1,46 @@
+// src/app/(app)/transactions/page.tsx
+import { fetchAdmin } from "@/lib/adminapi";
+import { fmtTs } from "@/lib/date";
+import type { Page } from "@/types/admin";
+
+interface AuditEvent {
+  stan: string;
+  mti: string;
+  pan: string;
+  amount: string;
+  responseCode: string;
+  destination: string;
+  createdAt: string;
+}
+
+export default async function TransactionsPage() {
+  const page = await fetchAdmin<Page<AuditEvent>>("/transactions?limit=50");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Transactions</h1>
+      <div className="rounded-lg border">
+        <table className="w-full text-sm">
+          <thead><tr className="border-b text-left text-muted-foreground">
+            <th className="px-3 py-2">STAN</th><th className="px-3 py-2">MTI</th><th className="px-3 py-2">PAN</th>
+            <th className="px-3 py-2">Amount</th><th className="px-3 py-2">Resp</th><th className="px-3 py-2">Dest</th>
+            <th className="px-3 py-2">Time</th>
+          </tr></thead>
+          <tbody>
+            {page.items.map(t => (
+              <tr key={t.stan} className="border-b last:border-0">
+                <td className="px-3 py-2 font-mono">{t.stan}</td>
+                <td className="px-3 py-2 font-mono">{t.mti}</td>
+                <td className="px-3 py-2 font-mono">{t.pan}</td>
+                <td className="px-3 py-2">{t.amount}</td>
+                <td className="px-3 py-2 font-mono">{t.responseCode}</td>
+                <td className="px-3 py-2 font-mono">{t.destination}</td>
+                <td className="px-3 py-2 text-muted-foreground">{fmtTs(t.createdAt)}</td>
+              </tr>
+            ))}
+          </tbody>
+        </table>
+        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No transactions yet ÔÇö run the seed (Task 10).</p>}
+      </div>
+    </div>
+  );
+}
\ No newline at end of file
diff --git a/web/src/components/cards/stat-card.tsx b/web/src/components/cards/stat-card.tsx
new file mode 100644
index 0000000..83076d8
--- /dev/null
+++ b/web/src/components/cards/stat-card.tsx
@@ -0,0 +1,10 @@
+// src/components/cards/stat-card.tsx
+import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
+
+export function StatCard({ title, value, hint }: { title: string; value: string; hint?: string }) {
+  return (
+    <Card><CardHeader className="pb-2"><CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle></CardHeader>
+      <CardContent><div className="text-2xl font-semibold">{value}</div>
+        {hint && <p className="mt-1 text-xs text-muted-foreground">{hint}</p>}</CardContent></Card>
+  );
+}
\ No newline at end of file
diff --git a/web/src/lib/adminapi.ts b/web/src/lib/adminapi.ts
new file mode 100644
index 0000000..6829915
--- /dev/null
+++ b/web/src/lib/adminapi.ts
@@ -0,0 +1,13 @@
+// src/lib/adminapi.ts
+import "server-only";
+
+export async function fetchAdmin<T>(path: string): Promise<T> {
+  const url = `/api/data${path}`;
+  const res = await fetch(url, { cache: "no-store" });
+  if (!res.ok) {
+    const body = await res.json().catch(() => null);
+    const msg = body && (body as { error?: string }).error ? (body as { error: string }).error : String(res.status);
+    throw new Error(`adminapi ${path}: ${msg}`);
+  }
+  return res.json() as Promise<T>;
+}
\ No newline at end of file
diff --git a/web/src/lib/money/minor.test.ts b/web/src/lib/money/minor.test.ts
new file mode 100644
index 0000000..aa47d06
--- /dev/null
+++ b/web/src/lib/money/minor.test.ts
@@ -0,0 +1,10 @@
+// src/lib/money/minor.test.ts
+import assert from "node:assert/strict";
+import { test } from "node:test";
+import { fmtMinor } from "./minor";
+
+test("fmtMinor formats minor units", () => {
+  assert.equal(fmtMinor(123456), "1,234.56 EUR");
+  assert.equal(fmtMinor(0), "0.00 EUR");
+  assert.equal(fmtMinor(-5000), "-50.00 EUR");
+});
\ No newline at end of file
diff --git a/web/src/lib/money/minor.ts b/web/src/lib/money/minor.ts
index 8dbcbd5..b9e18ee 100644
--- a/web/src/lib/money/minor.ts
+++ b/web/src/lib/money/minor.ts
@@ -1,8 +1,8 @@
 export function fmtMinor(minor: number, currency = "EUR"): string {
   const sign = minor < 0 ? "-" : "";
   const abs = Math.abs(minor);
   const major = Math.floor(abs / 100);
   const frac = String(abs % 100).padStart(2, "0");
   const intl = new Intl.NumberFormat("en-US").format(major);
-  return `${sign} ${intl}.${frac} ${currency}`.trim();
+  return `${sign}${intl}.${frac} ${currency}`;
 }
