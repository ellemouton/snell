create table if not exists articles_content (
    id int auto_increment,
    text longblob not null,

    primary key (id)
);

create table if not exists articles_info (
    id int auto_increment,
    name varchar(255) not null,
    description text not null,
    price int not null,
    content_id int not null,

    primary key (id)
);

