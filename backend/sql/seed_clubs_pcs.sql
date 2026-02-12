BEGIN;

INSERT INTO public.clubs (name, city, address)
SELECT v.name, v.city, v.address
FROM (
  VALUES
    ('Top Game', 'Astana', 'Astana, Dinmukhamed Kunayev Street 23'),
    ('BRO', 'Astana', 'Astana, Kuishi Dina Street 31'),
    ('Xan.exe', 'Astana', 'Astana, Heydar Aliyev Street 3'),
    ('Prime Game Hub', 'Astana', 'Astana, Kerei and Zhanibek Khandar Street 14/2'),
    ('Yamato Cyber Club', 'Astana', 'Astana, Syganak Street 21/1')
) AS v(name, city, address)
WHERE NOT EXISTS (
  SELECT 1
  FROM public.clubs c
  WHERE c.name = v.name
);

INSERT INTO public.pcs (club_id, pc_number, status)
SELECT c.id, p.pc_number, 'ACTIVE'
FROM public.clubs c
JOIN (
  VALUES
    ('Top Game', 1), ('Top Game', 2), ('Top Game', 3), ('Top Game', 4), ('Top Game', 5),
    ('Top Game', 6), ('Top Game', 7), ('Top Game', 8), ('Top Game', 9), ('Top Game', 10),
    ('BRO', 1), ('BRO', 2), ('BRO', 3), ('BRO', 4), ('BRO', 5),
    ('BRO', 6), ('BRO', 7), ('BRO', 8), ('BRO', 9), ('BRO', 10),
    ('Xan.exe', 1), ('Xan.exe', 2), ('Xan.exe', 3), ('Xan.exe', 4), ('Xan.exe', 5),
    ('Xan.exe', 6), ('Xan.exe', 7), ('Xan.exe', 8), ('Xan.exe', 9), ('Xan.exe', 10),
    ('Prime Game Hub', 1), ('Prime Game Hub', 2), ('Prime Game Hub', 3), ('Prime Game Hub', 4), ('Prime Game Hub', 5),
    ('Prime Game Hub', 6), ('Prime Game Hub', 7), ('Prime Game Hub', 8), ('Prime Game Hub', 9), ('Prime Game Hub', 10),
    ('Yamato Cyber Club', 1), ('Yamato Cyber Club', 2), ('Yamato Cyber Club', 3), ('Yamato Cyber Club', 4), ('Yamato Cyber Club', 5),
    ('Yamato Cyber Club', 6), ('Yamato Cyber Club', 7), ('Yamato Cyber Club', 8), ('Yamato Cyber Club', 9), ('Yamato Cyber Club', 10)
) AS p(club_name, pc_number)
  ON p.club_name = c.name
WHERE NOT EXISTS (
  SELECT 1
  FROM public.pcs x
  WHERE x.club_id = c.id AND x.pc_number = p.pc_number
);

COMMIT;

