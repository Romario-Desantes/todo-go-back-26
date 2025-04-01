ALTER TABLE
    public.tasks
ADD
    COLUMN user_id REFERENCES public.users(id);
