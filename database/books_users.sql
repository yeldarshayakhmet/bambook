create table if not exists books_users (
    book_id bigint references books(id) on update cascade on delete restrict,
    user_id bigint references users(id) on update cascade on delete restrict,
    constraint book_user_pkey primary key (book_id, user_id)
)