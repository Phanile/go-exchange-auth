-- +goose Up
-- +goose StatementBegin
create table Users (
    id serial primary key,
    email character(50) not null unique,
    pass_hash character(100) not null
);

create table Admins (
    id serial primary key,
    user_id int not null,
    constraint fk_Admins_Users
        foreign key (user_id)
        references Users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exist Users;
drop table if exist Admins;
-- +goose StatementEnd
