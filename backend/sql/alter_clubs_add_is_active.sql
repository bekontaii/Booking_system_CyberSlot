BEGIN;

ALTER TABLE public.clubs
ADD COLUMN IF NOT EXISTS is_active boolean NOT NULL DEFAULT true;

UPDATE public.clubs
SET is_active = true
WHERE is_active IS NULL;

COMMIT;

