create table if not exists articles_content (
    id int auto_increment,
    text longblob not null,

    primary key (id)
);

create table if not exists articles_info (
    id int auto_increment,
    name varchar(255) not null,
    decription text not null,
    price int not null,
    article_content int not null,

    primary key (id),
    foreign key (article_content) references articles_content(id) 
);

