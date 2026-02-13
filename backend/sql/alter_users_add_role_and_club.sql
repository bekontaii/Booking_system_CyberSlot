BEGIN;

ALTER TABLE public.users
ADD COLUMN IF NOT EXISTS role text NOT NULL DEFAULT 'USER';

ALTER TABLE public.users
ADD COLUMN IF NOT EXISTS club_id bigint NULL REFERENCES public.clubs(id);

ALTER TABLE public.users
DROP CONSTRAINT IF EXISTS users_role_check;

ALTER TABLE public.users
ADD CONSTRAINT users_role_check
CHECK (role IN ('USER', 'CLUB_ADMIN', 'SITE_ADMIN'));

ALTER TABLE public.users
DROP CONSTRAINT IF EXISTS users_club_admin_scope_check;

ALTER TABLE public.users
ADD CONSTRAINT users_club_admin_scope_check
CHECK (
  (role = 'CLUB_ADMIN' AND club_id IS NOT NULL)
  OR (role IN ('USER', 'SITE_ADMIN') AND club_id IS NULL)
);

COMMIT;
