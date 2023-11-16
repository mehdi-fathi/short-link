CREATE TYPE links_status AS ENUM ('pending','approve','reject','deleted') ;

alter table links
    add status links_status default 'pending';



