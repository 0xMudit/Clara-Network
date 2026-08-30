// syntax check via: node --check scripts/create-users.mjs
//
// Creates (or updates) the five Clara Network demo users in Supabase Auth using
// the service-role key and the admin API. Idempotent: existing users are
// matched by exact email and only have app_metadata.role (re)set; missing users
// are created with email_confirm: true.
//
// Run via `npm run db:users` from web/ (passes --env-file=.env.local, which must
// contain SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY).
import { createClient } from "@supabase/supabase-js";

const DEMO_PASSWORD = "ClaraDemo!2026";

const MANIFEST = [
  { role: "scheme_operator", email: "scheme-operator@clara.demo", password: DEMO_PASSWORD },
  { role: "issuer", email: "issuer@clara.demo", password: DEMO_PASSWORD },
  { role: "acquirer", email: "acquirer@clara.demo", password: DEMO_PASSWORD },
  { role: "merchant", email: "merchant@clara.demo", password: DEMO_PASSWORD },
  { role: "viewer", email: "viewer@clara.demo", password: DEMO_PASSWORD },
];

const supabaseUrl = process.env.SUPABASE_URL;
const serviceRoleKey = process.env.SUPABASE_SERVICE_ROLE_KEY;
if (!supabaseUrl || !serviceRoleKey) {
  console.error(
    "SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set. " +
      'Run via `npm run db:users` (passes --env-file=.env.local) with both keys present in web/.env.local.'
  );
  process.exit(1);
}

const supabase = createClient(supabaseUrl, serviceRoleKey, {
  auth: { autoRefreshToken: false, persistSession: false },
});

async function findByEmail(email) {
  let page = 1;
  for (;;) {
    const { data, error } = await supabase.auth.admin.listUsers({ page, perPage: 200 });
    if (error) throw error;
    const match = data.users.find((u) => u.email === email);
    if (match) return match;
    if (data.users.length === 0 || page >= data.lastPage) return null;
    page += 1;
  }
}

let failed = false;
for (const demo of MANIFEST) {
  try {
    const existing = await findByEmail(demo.email);
    if (existing) {
      const { error } = await supabase.auth.admin.updateUserById(existing.id, {
        app_metadata: { role: demo.role },
      });
      if (error) throw error;
      console.log(`updated ${demo.email} -> role ${demo.role} (id ${existing.id})`);
    } else {
      const { data, error } = await supabase.auth.admin.createUser({
        email: demo.email,
        password: demo.password,
        email_confirm: true,
        app_metadata: { role: demo.role },
      });
      if (error) throw error;
      console.log(`created ${demo.email} -> role ${demo.role} (id ${data.user.id})`);
    }
  } catch (err) {
    failed = true;
    console.error(`FAILED ${demo.email}: ${err?.message ?? err}`);
  }
}

if (failed) {
  console.error("one or more demo users failed (see above)");
  process.exit(1);
}
console.log("all demo users ready");