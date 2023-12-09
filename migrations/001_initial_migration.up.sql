CREATE TABLE public.chat_history (
	chat_id bigserial NOT null,
	user_id_a int8 NOT NULL,
	user_id_b int8 NOT NULL,
	sender_user_id int8 NOT NULL,
	message text NOT NULL,
	reply_time timestamp NOT null,
	CONSTRAINT chat_history_pkey PRIMARY KEY (chat_id)
);

CREATE TABLE public.group_chat_history (
	group_chat_id bigserial NOT null,
	group_id int8 NOT NULL,
	sender_user_id int8 NOT NULL,
	message text NOT NULL,
	reply_time timestamp NOT null,
	CONSTRAINT room_pkey PRIMARY KEY (group_chat_id)
);

CREATE TABLE IF NOT EXISTS public.user (
	id bigserial NOT NULL PRIMARY KEY,
	name VARCHAR(255) NULL,
	username VARCHAR(255) NOT  NULL,
	password VARCHAR(255) NOT NULL
);