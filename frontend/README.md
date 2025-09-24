# Frontend — PhD Student Portal (React Router v7 + TanStack Query + Tailwind)

## Pages (starter)
- **Login**: email + password, link to "Forgot password"
- **Forgot password**: submits email, expects a reset link (Mailpit in dev)
- **Reset password**: accepts token from URL, sets a new password
- **Dashboard**: placeholder for checklist, uploads, comments

## Routing
- Centralized in `main.tsx` using TanStack Router v7 (root route with `AppLayout`).
- Protect routes by checking presence of `token` in `localStorage` (add proper guards later).

## API Client
- `src/api/client.ts` provides a lightweight fetch wrapper including `Authorization` header.

## Styling
- Tailwind is configured. You can add shadcn/ui components later with `npx shadcn-ui init`.

## Run
```bash
npm i
VITE_API_URL=http://localhost:8080/api npm run dev
```


## UI & Interactions (v4)
- **shadcn-style components**: local stubs in `src/components/ui/*` (Button, Card, Input, Textarea). You can replace them by running:
  ```bash
  npx shadcn-ui@latest init
  # then generate components you want and swap imports
  ```
- **Micro-interactions**: Framer Motion used in Checklist & Advisor Inbox.
- **Document detail page** (`/doc/:id` wired by you): auto-detects S3 support; uses pre-signed upload if available, else local multipart.
- **Threaded comments** with @mentions search (uses `/api/admin/users`).


## Role-Based Route Guards
We implement guards using TanStack Router's `beforeLoad`:
- `requireAuth()` redirects to `/login`.
- `requireRole('admin',...)` blocks unauthorized users (renders 403).

See `src/auth/auth.ts` and guards attached to routes in `src/main.tsx`.

## Design System (shadcn/ui — vendored)
A minimal set of shadcn-style components is pre-included under `src/components/ui/`:
- Button, Card, Input, Textarea, Badge, Accordion, DropdownMenu, plus `theme.css` tokens.
You may replace them with official `shadcn/ui` via:
```bash
npx shadcn-ui@latest init
# generate components (button, card, input, textarea, badge, accordion, dropdown-menu, toast, etc.)
```
Then switch imports to the generated aliases if desired.


## Auth via /me
- The app now relies on `GET /api/me` for user identity/role. JWT decoding in the client is no longer needed.
- Role-aware **TopNav** only shows links the user can access.

## Forms & Toasts
- `react-hook-form` + `zod` for validation (example: Login form).
- Simple toast system in `src/lib/toast.tsx` (replace with shadcn/ui `toast` later if desired).

## Structure
- Common utilities grouped under:
  - `src/config/` (e.g., `app.ts`),
  - `src/hooks/` (e.g., `useMe`),
  - `src/lib/` (e.g., `toast`).

## Ports
- Frontend talks to `VITE_API_URL` (default `http://localhost:8280/api`).


## Navigation & Roles
- Top navigation adapts to the user returned from `/me` (role-aware).

## Forms & Validation
- `react-hook-form` + `zod` for login and admin create user.

## Toasts
- Global toast provider in `src/components/toast.tsx`. Use `const { push } = useToast()`.

## Vertical Progress
- `src/components/VerticalProgress.tsx` shows a mobile-friendly vertical stepper for a module.


## Breadcrumbs & Active Nav
- The top navigation highlights the active route.
- Breadcrumbs appear under the header and follow TanStack Router matches.

## Document Detail Linking & PDF Preview
- Checklist "Open" buttons create/open a document and navigate to `/documents/:docId`.
- On the Document page, "Open latest version" tries S3 pre-signed GET, falls back to a local download route.


## Documents
- `/students/:id/documents` API lists documents; UI page `DocumentsList` renders cards with links.
- Checklist's **Open** action creates a document (if needed) and routes to `/documents/:docId`.
- Document Detail shows an inline **PDF preview** (first page) using `react-pdf`, and a link to open the full file.

## Mentions UI
- Chips + avatars (initials) with a pick list; remove chips inline. Drives mentions array for comment submissions.
