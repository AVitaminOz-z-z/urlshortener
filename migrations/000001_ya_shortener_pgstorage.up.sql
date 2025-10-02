-- storage
drop table if exists storage;
--$$--
create table if not exists storage (
    id          bigserial,
    url         text,
    short_url   text,
    rand        text
);
--$$--
alter table storage add constraint pk__storage__id primary key (id);
--$$--
create unique index if not exists ux__storage__short_url on storage(short_url);
--$$--
create unique index if not exists ux__storage__url on storage(url);
--$$--

-- fn__load_from_file_storage
drop function if exists fn__load_from_file_storage;
--$$--
create or replace function fn__load_from_file_storage(fs_data jsonb) returns bigint as $$
declare
    fs_key  text    = 'POSTStorage';
    d_rc    bigint  = 0;
begin
    with data as (
        select fs_data -> fs_key as json
    ), keys as (
        select jsonb_object_keys(data.json) f_url, data.json from data)
    insert into storage (url, short_url, rand)
    select f_url, (json -> f_url ->> 0), (json -> f_url ->> 1) from keys
    on conflict (short_url) do nothing;
    get diagnostics d_rc = ROW_COUNT;
    return d_rc;
end
$$ language plpgsql;
--$$--
grant execute on function fn__load_from_file_storage(jsonb) to public;
--$$--

-- fn__gen_random_string
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

-- fn__insert_short_url
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

-- fn__return_short_url
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

-- fn__return_full_url
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
