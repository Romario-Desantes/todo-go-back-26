ALTER TABLE
    public.tasks
ADD
    COLUMN user_id integer REFERENCES public.users(id);
