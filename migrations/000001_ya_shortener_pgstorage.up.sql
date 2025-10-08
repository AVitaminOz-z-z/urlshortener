--$$--drop table if exists storage;
create table if not exists storage (
   id          bigserial,
   url         text,
   short_url   text,
   rand        text
);
--$$--
alter table storage drop constraint if exists pk__storage__id;
--$$--
alter table storage add constraint pk__storage__id primary key (id);
--$$--
create unique index if not exists ux__storage__short_url on storage(short_url);
--$$--
create unique index if not exists ux__storage__url on storage(url);
--$$--
drop function if exists fn__gen_random_string;
--$$--
create or replace function fn__gen_random_string (max_len int) returns text as $$
begin
    return(
        select string_agg(substr(characters, (random() * length(characters) + 1)::int, 1), '')
        from (values('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789')) as symbols(characters)
                 join generate_series(1, max_len) on true)::text;
end
$$ language plpgsql;
--$$--
grant execute on function fn__gen_random_string(int) to public;
--$$--
drop function if exists fn__insert_short_url;
--$$--
create or replace function fn__insert_short_url(full_url text) returns text as $$
declare
    max_len  int = 16;
    rnd_seed text = fn__gen_random_string(max_len);
begin
    if not exists(select short_url from storage where url = full_url) then
        insert into storage (url, short_url, rand)
        select full_url, encode(sha256((full_url || rnd_seed)::bytea), 'hex'), rnd_seed
        on conflict (url) do nothing;
    end if;
    return (select short_url from storage where url = full_url);
end
$$ language plpgsql;
--$$--
grant execute on function fn__insert_short_url(text) to public;
--$$--
drop function if exists fn__return_short_url;
--$$--
create or replace function fn__return_short_url(full_url text) returns text as $$
begin
    return fn__insert_short_url(full_url);
end
$$ language plpgsql;
--$$--
grant execute on function fn__return_short_url(text) to public;
--$$--
drop function if exists fn__return_full_url;
--$$--
create or replace function fn__return_full_url(short text) returns text as $$
declare
    url text = (select url from storage where short_url = short);
begin
    return coalesce(url, '');
end
$$ language plpgsql;
--$$--
grant execute on function fn__return_full_url(text) to public;
--$$--
drop function if exists fn__insert_short_url_on_conflict;
--$$--
create or replace function fn__insert_short_url_on_conflict(full_url text, prefix text default '') returns table (
    data        text,
    conflict    bool
) as $$
declare
    max_len     int  = 16;
    rnd_seed    text = fn__gen_random_string(max_len);
    is_conflict bool;
begin
    with cte as
        (insert into storage (url, short_url, rand)
        select full_url, encode(sha256((full_url || rnd_seed)::bytea), 'hex'), rnd_seed
        on conflict (url) do nothing
        returning *)
    select not exists(select cte.* from cte) into is_conflict;
    return query (select prefix || short_url, is_conflict from storage where url = full_url limit 1);
end
$$ language plpgsql;
--$$--
grant execute on function fn__insert_short_url_on_conflict(text, text) to public;
--$$--
drop function if exists fn__return_short_url_on_conflict;
--$$--
create or replace function fn__return_short_url_on_conflict(full_url text, prefix text default '') returns table (
    data        text,
    conflict    bool
) as $$
begin
    return query (select * from fn__insert_short_url_on_conflict(full_url, prefix));
end
$$ language plpgsql;
--$$--
grant execute on function fn__return_short_url_on_conflict(text, text) to public;
--$$--
drop function if exists fn__return_batch_short_urls;
--$$--
create or replace function fn__return_batch_short_urls(batch jsonb, prefix text default '') returns jsonb as $$
begin
    return (
        select to_jsonb(array_agg(t1.*))
        from (select t.correlation_id, prefix || fn__return_short_url(t.original_url) as short_url
              from  jsonb_to_recordset(batch)
                        as t(
                             "correlation_id" text,
                             "original_url" text
                      ))t1
    );
end
$$ language plpgsql;
--$$--
grant execute on function fn__return_batch_short_urls(jsonb, text) to public;