
SET session_replication_role = replica;

--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: audit_log_entries; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--

INSERT INTO "auth"."audit_log_entries" ("instance_id", "id", "payload", "created_at", "ip_address") VALUES
	('00000000-0000-0000-0000-000000000000', '485ea5fb-d69a-40d0-a10c-2cfda99cc0b0', '{"action":"user_signedup","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"team","traits":{"provider":"discord"}}', '2025-08-16 21:11:51.573142+00', ''),
	('00000000-0000-0000-0000-000000000000', 'fb7697e8-9ae7-4380-b8e4-d1be2f7a73ef', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-16 21:12:57.367823+00', ''),
	('00000000-0000-0000-0000-000000000000', '69fba3d6-18c3-4adc-82ab-c49d03671313', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-16 22:27:40.471553+00', ''),
	('00000000-0000-0000-0000-000000000000', '3e2a3817-6da8-44a5-b746-2e1f3c2cc191', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 00:38:46.703967+00', ''),
	('00000000-0000-0000-0000-000000000000', '0add2658-86d5-44be-b51f-de85dd59d4f4', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 02:11:47.433172+00', ''),
	('00000000-0000-0000-0000-000000000000', '90cf0994-52bd-487d-a847-8d8170e7cb2c', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 02:53:55.791488+00', ''),
	('00000000-0000-0000-0000-000000000000', 'f21b8335-f778-45b4-9f60-f17c2f0f88e7', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 02:58:05.572128+00', ''),
	('00000000-0000-0000-0000-000000000000', '1c8244c9-5366-449c-b426-ef55c3c858e8', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 08:21:12.117517+00', ''),
	('00000000-0000-0000-0000-000000000000', '76b5c5e2-3889-4a55-a87a-2a8d55ba3ef6', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 08:37:37.148197+00', ''),
	('00000000-0000-0000-0000-000000000000', '0ce5632c-11d6-4446-b388-875dfed105f4', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 09:21:57.236562+00', ''),
	('00000000-0000-0000-0000-000000000000', '7328ee7e-7063-48c3-8ce5-9f94876b18e8', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 09:24:33.836438+00', ''),
	('00000000-0000-0000-0000-000000000000', 'f6958e32-8ccb-4314-bb9c-223dd9140648', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 09:32:07.693775+00', ''),
	('00000000-0000-0000-0000-000000000000', '20714513-4793-48c8-af63-347cdf6c1b7a', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 16:29:09.74315+00', ''),
	('00000000-0000-0000-0000-000000000000', '76dd649a-29e2-4a85-a0d6-77bfaca96ba0', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 16:29:45.599156+00', ''),
	('00000000-0000-0000-0000-000000000000', 'f12844f4-af54-4909-bae2-1b34b6009491', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 16:49:16.10101+00', ''),
	('00000000-0000-0000-0000-000000000000', 'f9d0e150-7e7b-4a1e-9852-83830a3a8d0a', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 18:08:00.16107+00', ''),
	('00000000-0000-0000-0000-000000000000', '66a763bc-b91c-4036-a55b-819a784f0d16', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 18:15:59.794476+00', ''),
	('00000000-0000-0000-0000-000000000000', '8e065c8c-58fb-4343-8408-dcd17b88af11', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 18:41:00.982772+00', ''),
	('00000000-0000-0000-0000-000000000000', 'c7990a6d-9c27-4d80-87d5-4e2beb7afd2e', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 21:08:23.820217+00', ''),
	('00000000-0000-0000-0000-000000000000', '12754867-bd4b-4b41-a0c3-ee64fd11d234', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 21:09:45.149233+00', ''),
	('00000000-0000-0000-0000-000000000000', '018a91ac-ab4e-49f6-9097-c5d7f1f156b7', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 21:10:49.965824+00', ''),
	('00000000-0000-0000-0000-000000000000', '84b30081-a030-46c1-85a5-d4bcae3fdaa6', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 21:38:30.202993+00', ''),
	('00000000-0000-0000-0000-000000000000', '0716fbcd-c747-44ca-a1a7-5cea286a313b', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 21:46:01.419024+00', ''),
	('00000000-0000-0000-0000-000000000000', '0c125264-884e-4195-a906-4cf7b92b45ba', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 21:48:17.558224+00', ''),
	('00000000-0000-0000-0000-000000000000', 'aa748310-7d4f-48af-a596-90d216b9a53b', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-17 21:56:57.446128+00', ''),
	('00000000-0000-0000-0000-000000000000', '0d7cab45-8ce0-46c9-b8b4-9dfc4bc662ef', '{"action":"user_signedup","actor_id":"5c0f44a8-7beb-45fd-81e2-023cdf7cde73","actor_name":"rockytheii","actor_username":"rockohamilton09@gmail.com","actor_via_sso":false,"log_type":"team","traits":{"provider":"discord"}}', '2025-08-18 02:03:09.280604+00', ''),
	('00000000-0000-0000-0000-000000000000', '1a713c08-43f3-48ce-913a-d99af3d4c6e3', '{"action":"login","actor_id":"38977817-1066-40d9-ab5e-a8a8ab8e667d","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"account","traits":{"provider":"discord"}}', '2025-08-18 02:37:50.014404+00', '');


--
-- Data for Name: flow_state; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: users; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--

INSERT INTO "auth"."users" ("instance_id", "id", "aud", "role", "email", "encrypted_password", "email_confirmed_at", "invited_at", "confirmation_token", "confirmation_sent_at", "recovery_token", "recovery_sent_at", "email_change_token_new", "email_change", "email_change_sent_at", "last_sign_in_at", "raw_app_meta_data", "raw_user_meta_data", "is_super_admin", "created_at", "updated_at", "phone", "phone_confirmed_at", "phone_change", "phone_change_token", "phone_change_sent_at", "email_change_token_current", "email_change_confirm_status", "banned_until", "reauthentication_token", "reauthentication_sent_at", "is_sso_user", "deleted_at", "is_anonymous") VALUES
	('00000000-0000-0000-0000-000000000000', '38977817-1066-40d9-ab5e-a8a8ab8e667d', 'authenticated', 'authenticated', 'reilly.kyle101@gmail.com', NULL, '2025-08-16 21:11:51.577506+00', NULL, '', NULL, '', NULL, '', '', NULL, '2025-08-18 02:37:50.019355+00', '{"provider": "discord", "providers": ["discord"]}', '{"iss": "https://discord.com/api", "sub": "354474826192388127", "name": "mogtron#0", "email": "reilly.kyle101@gmail.com", "picture": "https://cdn.discordapp.com/avatars/354474826192388127/c50ea819dfdc72747e00c6f67d4ade97.png", "full_name": "mogtron", "avatar_url": "https://cdn.discordapp.com/avatars/354474826192388127/c50ea819dfdc72747e00c6f67d4ade97.png", "provider_id": "354474826192388127", "custom_claims": {"global_name": "mogtron"}, "email_verified": true, "phone_verified": false}', NULL, '2025-08-16 21:11:51.543303+00', '2025-08-18 02:37:50.039424+00', NULL, NULL, '', '', NULL, '', 0, NULL, '', NULL, false, NULL, false),
	('00000000-0000-0000-0000-000000000000', '5c0f44a8-7beb-45fd-81e2-023cdf7cde73', 'authenticated', 'authenticated', 'rockohamilton09@gmail.com', NULL, '2025-08-18 02:03:09.292703+00', NULL, '', NULL, '', NULL, '', '', NULL, '2025-08-18 02:03:09.3026+00', '{"provider": "discord", "providers": ["discord"]}', '{"iss": "https://discord.com/api", "sub": "679038415576104971", "name": "rockytheii#0", "email": "rockohamilton09@gmail.com", "picture": "https://cdn.discordapp.com/avatars/679038415576104971/2da887db42c84e1712cbe1b33a1526f8.png", "full_name": "rockytheii", "avatar_url": "https://cdn.discordapp.com/avatars/679038415576104971/2da887db42c84e1712cbe1b33a1526f8.png", "provider_id": "679038415576104971", "custom_claims": {"global_name": "Edge"}, "email_verified": true, "phone_verified": false}', NULL, '2025-08-18 02:03:09.197113+00', '2025-08-18 02:03:09.351705+00', NULL, NULL, '', '', NULL, '', 0, NULL, '', NULL, false, NULL, false);


--
-- Data for Name: identities; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--

INSERT INTO "auth"."identities" ("provider_id", "user_id", "identity_data", "provider", "last_sign_in_at", "created_at", "updated_at", "id") VALUES
	('679038415576104971', '5c0f44a8-7beb-45fd-81e2-023cdf7cde73', '{"iss": "https://discord.com/api", "sub": "679038415576104971", "name": "rockytheii#0", "email": "rockohamilton09@gmail.com", "picture": "https://cdn.discordapp.com/avatars/679038415576104971/2da887db42c84e1712cbe1b33a1526f8.png", "full_name": "rockytheii", "avatar_url": "https://cdn.discordapp.com/avatars/679038415576104971/2da887db42c84e1712cbe1b33a1526f8.png", "provider_id": "679038415576104971", "custom_claims": {"global_name": "Edge"}, "email_verified": true, "phone_verified": false}', 'discord', '2025-08-18 02:03:09.262032+00', '2025-08-18 02:03:09.262102+00', '2025-08-18 02:03:09.262102+00', '503962c2-ce10-41e8-91ca-8548e4e6ee5c'),
	('354474826192388127', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '{"iss": "https://discord.com/api", "sub": "354474826192388127", "name": "mogtron#0", "email": "reilly.kyle101@gmail.com", "picture": "https://cdn.discordapp.com/avatars/354474826192388127/c50ea819dfdc72747e00c6f67d4ade97.png", "full_name": "mogtron", "avatar_url": "https://cdn.discordapp.com/avatars/354474826192388127/c50ea819dfdc72747e00c6f67d4ade97.png", "provider_id": "354474826192388127", "custom_claims": {"global_name": "mogtron"}, "email_verified": true, "phone_verified": false}', 'discord', '2025-08-16 21:11:51.567016+00', '2025-08-16 21:11:51.567073+00', '2025-08-18 02:37:49.995101+00', '41d2270d-7ef8-4b8c-a08f-b25583273e26');


--
-- Data for Name: instances; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: sessions; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--

INSERT INTO "auth"."sessions" ("id", "user_id", "created_at", "updated_at", "factor_id", "aal", "not_after", "refreshed_at", "user_agent", "ip", "tag") VALUES
	('cdf3361f-10b3-4c3a-a8e8-b936e380ddf3', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-16 21:11:51.586415+00', '2025-08-16 21:11:51.586415+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('68e2e1af-49e7-4b72-8a66-063be3b747cd', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-16 21:12:57.368696+00', '2025-08-16 21:12:57.368696+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('0b3bf938-5cb0-4d61-bef6-d3c638829b4a', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-16 22:27:40.480343+00', '2025-08-16 22:27:40.480343+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('0f9c65c2-b02e-468b-8767-c579a39ab012', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 00:38:46.72083+00', '2025-08-17 00:38:46.72083+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('d2c10324-150a-4a3c-9eb3-4be9d48ab10b', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 02:11:47.449678+00', '2025-08-17 02:11:47.449678+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('c196da87-c477-4420-bf4e-0a2f98847813', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 02:53:55.797573+00', '2025-08-17 02:53:55.797573+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('1af0f694-dde0-4f6a-bdc3-020e01cb6471', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 02:58:05.579867+00', '2025-08-17 02:58:05.579867+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('b207d6e4-3052-4a9a-94dc-22f4c4adc773', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 08:21:12.127149+00', '2025-08-17 08:21:12.127149+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('07bc4793-461e-4271-b8e0-d90144004f07', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 08:37:37.167128+00', '2025-08-17 08:37:37.167128+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('01afda98-30e4-4392-8d42-f128905c0ae9', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 09:21:57.248318+00', '2025-08-17 09:21:57.248318+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('3892c9b3-c615-4b72-9a01-82cfcbb498f3', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 09:24:33.837313+00', '2025-08-17 09:24:33.837313+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('9b9ad240-2f3a-424a-bd4e-ba51cb344241', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 09:32:07.695954+00', '2025-08-17 09:32:07.695954+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('fe4b839c-2e93-421d-ac8b-05724aaee6fa', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 16:29:09.754+00', '2025-08-17 16:29:09.754+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('62aeb8df-43ec-415b-a939-3620695c4cd5', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 16:29:45.599848+00', '2025-08-17 16:29:45.599848+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('8af87453-0756-4a01-8016-3d41db98f30d', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 16:49:16.110527+00', '2025-08-17 16:49:16.110527+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('3285adf3-807b-4ad6-897b-bb3813627252', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 18:08:00.18013+00', '2025-08-17 18:08:00.18013+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('b8fea0e7-d98f-4d03-a1b4-7ad3e35dc6a4', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 18:15:59.798701+00', '2025-08-17 18:15:59.798701+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('82866f49-a1a9-4fc6-b062-c5202034a3b2', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 18:41:00.995652+00', '2025-08-17 18:41:00.995652+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('7793598e-fd04-411e-bb91-26b84fe2d2f4', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 21:08:23.838736+00', '2025-08-17 21:08:23.838736+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('e04ed28b-de59-44a8-b1c2-6da78d6b0a3f', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 21:09:45.150716+00', '2025-08-17 21:09:45.150716+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('b35db76c-528d-4aff-b70e-8384a3a246f9', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 21:10:49.967288+00', '2025-08-17 21:10:49.967288+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('17b134d2-9c01-4cb3-b06e-a9c5c0740aa0', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 21:38:30.210116+00', '2025-08-17 21:38:30.210116+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('df7cd37e-0fb7-4233-bac3-44cb8455ef40', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 21:46:01.423933+00', '2025-08-17 21:46:01.423933+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('b31998bf-d128-4ef0-a0de-dc97866c6147', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 21:48:17.55908+00', '2025-08-17 21:48:17.55908+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('559e7859-07f8-4713-b7d7-562ad666f0fa', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-17 21:56:57.453065+00', '2025-08-17 21:56:57.453065+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL),
	('836f28e5-039e-4058-91d0-67d7654e8cfe', '5c0f44a8-7beb-45fd-81e2-023cdf7cde73', '2025-08-18 02:03:09.303736+00', '2025-08-18 02:03:09.303736+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '192.230.137.53', NULL),
	('58f51ad0-bdce-40fe-83dc-5a86d5a023d4', '38977817-1066-40d9-ab5e-a8a8ab8e667d', '2025-08-18 02:37:50.021182+00', '2025-08-18 02:37:50.021182+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '24.16.165.99', NULL);


--
-- Data for Name: mfa_amr_claims; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--

INSERT INTO "auth"."mfa_amr_claims" ("session_id", "created_at", "updated_at", "authentication_method", "id") VALUES
	('cdf3361f-10b3-4c3a-a8e8-b936e380ddf3', '2025-08-16 21:11:51.627139+00', '2025-08-16 21:11:51.627139+00', 'oauth', '6f4cc2ed-5d91-43f0-a932-aab894e7e82c'),
	('68e2e1af-49e7-4b72-8a66-063be3b747cd', '2025-08-16 21:12:57.372019+00', '2025-08-16 21:12:57.372019+00', 'oauth', 'c5cdb168-c9bd-49d4-92a5-d414848849e8'),
	('0b3bf938-5cb0-4d61-bef6-d3c638829b4a', '2025-08-16 22:27:40.499164+00', '2025-08-16 22:27:40.499164+00', 'oauth', 'd969781c-3379-415d-8c83-558d9ceea6ca'),
	('0f9c65c2-b02e-468b-8767-c579a39ab012', '2025-08-17 00:38:46.757147+00', '2025-08-17 00:38:46.757147+00', 'oauth', '2f4d57fd-775e-4c33-9737-6f4bc9b880db'),
	('d2c10324-150a-4a3c-9eb3-4be9d48ab10b', '2025-08-17 02:11:47.488777+00', '2025-08-17 02:11:47.488777+00', 'oauth', 'dce5a873-8cec-40f5-a83e-20d09678bf97'),
	('c196da87-c477-4420-bf4e-0a2f98847813', '2025-08-17 02:53:55.814549+00', '2025-08-17 02:53:55.814549+00', 'oauth', '2c93a353-aa22-4b11-adb3-f8bb57ed4db5'),
	('1af0f694-dde0-4f6a-bdc3-020e01cb6471', '2025-08-17 02:58:05.591417+00', '2025-08-17 02:58:05.591417+00', 'oauth', 'c7b59d73-cf6e-48b0-97e6-c3d03d49f917'),
	('b207d6e4-3052-4a9a-94dc-22f4c4adc773', '2025-08-17 08:21:12.150251+00', '2025-08-17 08:21:12.150251+00', 'oauth', '75b8e4c5-322a-4cca-8959-a6693d6c2aa5'),
	('07bc4793-461e-4271-b8e0-d90144004f07', '2025-08-17 08:37:37.201237+00', '2025-08-17 08:37:37.201237+00', 'oauth', '4ab897e2-d1f2-404d-9f8f-9ab2e44d59d0'),
	('01afda98-30e4-4392-8d42-f128905c0ae9', '2025-08-17 09:21:57.272652+00', '2025-08-17 09:21:57.272652+00', 'oauth', 'e2d52455-d604-42b3-8392-c2c2654a2d21'),
	('3892c9b3-c615-4b72-9a01-82cfcbb498f3', '2025-08-17 09:24:33.841113+00', '2025-08-17 09:24:33.841113+00', 'oauth', '9739404d-3f6e-4e3f-853f-13e9606f18b1'),
	('9b9ad240-2f3a-424a-bd4e-ba51cb344241', '2025-08-17 09:32:07.701781+00', '2025-08-17 09:32:07.701781+00', 'oauth', 'f1e4f486-b5c0-4ae5-a710-a8e0c766f265'),
	('fe4b839c-2e93-421d-ac8b-05724aaee6fa', '2025-08-17 16:29:09.791064+00', '2025-08-17 16:29:09.791064+00', 'oauth', '158e1d3d-5d64-457e-a212-331ac3907ade'),
	('62aeb8df-43ec-415b-a939-3620695c4cd5', '2025-08-17 16:29:45.603633+00', '2025-08-17 16:29:45.603633+00', 'oauth', 'c25a0f74-10ed-4c2d-9b2a-28bc53e55057'),
	('8af87453-0756-4a01-8016-3d41db98f30d', '2025-08-17 16:49:16.124717+00', '2025-08-17 16:49:16.124717+00', 'oauth', '3bf38ff6-0ab8-4fc5-83dd-543df4c5136d'),
	('3285adf3-807b-4ad6-897b-bb3813627252', '2025-08-17 18:08:00.220842+00', '2025-08-17 18:08:00.220842+00', 'oauth', '010320ac-b745-4a15-8657-c03cf582f6b1'),
	('b8fea0e7-d98f-4d03-a1b4-7ad3e35dc6a4', '2025-08-17 18:15:59.815205+00', '2025-08-17 18:15:59.815205+00', 'oauth', '6bab4494-6790-40d9-8067-ad6ed1774b21'),
	('82866f49-a1a9-4fc6-b062-c5202034a3b2', '2025-08-17 18:41:01.019376+00', '2025-08-17 18:41:01.019376+00', 'oauth', 'f6744854-7c86-4253-bccc-f636572cf57c'),
	('7793598e-fd04-411e-bb91-26b84fe2d2f4', '2025-08-17 21:08:23.878474+00', '2025-08-17 21:08:23.878474+00', 'oauth', '3f65b1f5-ab58-42a0-a332-6a95d4e7ddc5'),
	('e04ed28b-de59-44a8-b1c2-6da78d6b0a3f', '2025-08-17 21:09:45.156832+00', '2025-08-17 21:09:45.156832+00', 'oauth', '8c3af26f-e4e0-4911-af3c-e9c6e679eea6'),
	('b35db76c-528d-4aff-b70e-8384a3a246f9', '2025-08-17 21:10:49.970068+00', '2025-08-17 21:10:49.970068+00', 'oauth', 'af7bb37a-755e-4ae6-8678-515d724b2d57'),
	('17b134d2-9c01-4cb3-b06e-a9c5c0740aa0', '2025-08-17 21:38:30.220322+00', '2025-08-17 21:38:30.220322+00', 'oauth', 'f04cdb83-6e2f-48e5-b01b-d2e59cb13d30'),
	('df7cd37e-0fb7-4233-bac3-44cb8455ef40', '2025-08-17 21:46:01.432187+00', '2025-08-17 21:46:01.432187+00', 'oauth', 'e05be537-340a-4702-9d6a-80f0d58f38a8'),
	('b31998bf-d128-4ef0-a0de-dc97866c6147', '2025-08-17 21:48:17.562691+00', '2025-08-17 21:48:17.562691+00', 'oauth', 'efc6308d-16ce-41a4-8086-c0640aaee911'),
	('559e7859-07f8-4713-b7d7-562ad666f0fa', '2025-08-17 21:56:57.459914+00', '2025-08-17 21:56:57.459914+00', 'oauth', 'c9758ff7-e767-40ee-a710-60cea887aa4e'),
	('836f28e5-039e-4058-91d0-67d7654e8cfe', '2025-08-18 02:03:09.35221+00', '2025-08-18 02:03:09.35221+00', 'oauth', '5f0a2dc4-319b-40e2-b1f9-5883df0c19f4'),
	('58f51ad0-bdce-40fe-83dc-5a86d5a023d4', '2025-08-18 02:37:50.039974+00', '2025-08-18 02:37:50.039974+00', 'oauth', 'febe8405-2337-4292-be9e-c6950adebae4');


--
-- Data for Name: mfa_factors; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: mfa_challenges; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: one_time_tokens; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--

INSERT INTO "auth"."refresh_tokens" ("instance_id", "id", "token", "user_id", "revoked", "created_at", "updated_at", "parent", "session_id") VALUES
	('00000000-0000-0000-0000-000000000000', 2, 'ielwj3agr4pr', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-16 21:11:51.602781+00', '2025-08-16 21:11:51.602781+00', NULL, 'cdf3361f-10b3-4c3a-a8e8-b936e380ddf3'),
	('00000000-0000-0000-0000-000000000000', 3, 'jj46okkbifat', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-16 21:12:57.369726+00', '2025-08-16 21:12:57.369726+00', NULL, '68e2e1af-49e7-4b72-8a66-063be3b747cd'),
	('00000000-0000-0000-0000-000000000000', 4, 'tr27df73snzi', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-16 22:27:40.488847+00', '2025-08-16 22:27:40.488847+00', NULL, '0b3bf938-5cb0-4d61-bef6-d3c638829b4a'),
	('00000000-0000-0000-0000-000000000000', 5, '5awqcjx6wz35', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 00:38:46.737177+00', '2025-08-17 00:38:46.737177+00', NULL, '0f9c65c2-b02e-468b-8767-c579a39ab012'),
	('00000000-0000-0000-0000-000000000000', 6, 'wawiwb62g4rk', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 02:11:47.465359+00', '2025-08-17 02:11:47.465359+00', NULL, 'd2c10324-150a-4a3c-9eb3-4be9d48ab10b'),
	('00000000-0000-0000-0000-000000000000', 7, 'oceu3nkjgohj', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 02:53:55.802049+00', '2025-08-17 02:53:55.802049+00', NULL, 'c196da87-c477-4420-bf4e-0a2f98847813'),
	('00000000-0000-0000-0000-000000000000', 8, '75trlvt5l67o', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 02:58:05.582691+00', '2025-08-17 02:58:05.582691+00', NULL, '1af0f694-dde0-4f6a-bdc3-020e01cb6471'),
	('00000000-0000-0000-0000-000000000000', 9, 'cjz4qvymlnru', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 08:21:12.137517+00', '2025-08-17 08:21:12.137517+00', NULL, 'b207d6e4-3052-4a9a-94dc-22f4c4adc773'),
	('00000000-0000-0000-0000-000000000000', 10, 'aftlps56ex5w', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 08:37:37.174824+00', '2025-08-17 08:37:37.174824+00', NULL, '07bc4793-461e-4271-b8e0-d90144004f07'),
	('00000000-0000-0000-0000-000000000000', 11, 'bip7wbt4hrqg', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 09:21:57.255611+00', '2025-08-17 09:21:57.255611+00', NULL, '01afda98-30e4-4392-8d42-f128905c0ae9'),
	('00000000-0000-0000-0000-000000000000', 12, 'htcpalosf2gc', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 09:24:33.838806+00', '2025-08-17 09:24:33.838806+00', NULL, '3892c9b3-c615-4b72-9a01-82cfcbb498f3'),
	('00000000-0000-0000-0000-000000000000', 13, 'jurbhkjhxosw', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 09:32:07.697044+00', '2025-08-17 09:32:07.697044+00', NULL, '9b9ad240-2f3a-424a-bd4e-ba51cb344241'),
	('00000000-0000-0000-0000-000000000000', 14, 'tz4735dcv76d', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 16:29:09.767595+00', '2025-08-17 16:29:09.767595+00', NULL, 'fe4b839c-2e93-421d-ac8b-05724aaee6fa'),
	('00000000-0000-0000-0000-000000000000', 15, 'ldd3g3qeeaxu', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 16:29:45.601864+00', '2025-08-17 16:29:45.601864+00', NULL, '62aeb8df-43ec-415b-a939-3620695c4cd5'),
	('00000000-0000-0000-0000-000000000000', 16, 'bdebi7tuh2uj', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 16:49:16.116804+00', '2025-08-17 16:49:16.116804+00', NULL, '8af87453-0756-4a01-8016-3d41db98f30d'),
	('00000000-0000-0000-0000-000000000000', 17, 'omi3ldhbptma', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 18:08:00.19467+00', '2025-08-17 18:08:00.19467+00', NULL, '3285adf3-807b-4ad6-897b-bb3813627252'),
	('00000000-0000-0000-0000-000000000000', 18, 'pgun5djeeugz', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 18:15:59.801511+00', '2025-08-17 18:15:59.801511+00', NULL, 'b8fea0e7-d98f-4d03-a1b4-7ad3e35dc6a4'),
	('00000000-0000-0000-0000-000000000000', 19, 'xoozvclprfqe', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 18:41:01.003454+00', '2025-08-17 18:41:01.003454+00', NULL, '82866f49-a1a9-4fc6-b062-c5202034a3b2'),
	('00000000-0000-0000-0000-000000000000', 20, 'fhduv67x3kl5', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 21:08:23.850795+00', '2025-08-17 21:08:23.850795+00', NULL, '7793598e-fd04-411e-bb91-26b84fe2d2f4'),
	('00000000-0000-0000-0000-000000000000', 21, 'muqslrosr52s', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 21:09:45.153702+00', '2025-08-17 21:09:45.153702+00', NULL, 'e04ed28b-de59-44a8-b1c2-6da78d6b0a3f'),
	('00000000-0000-0000-0000-000000000000', 22, 'vgrjhk5dv4r7', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 21:10:49.968325+00', '2025-08-17 21:10:49.968325+00', NULL, 'b35db76c-528d-4aff-b70e-8384a3a246f9'),
	('00000000-0000-0000-0000-000000000000', 23, 'bqjnede6gfnc', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 21:38:30.214256+00', '2025-08-17 21:38:30.214256+00', NULL, '17b134d2-9c01-4cb3-b06e-a9c5c0740aa0'),
	('00000000-0000-0000-0000-000000000000', 24, 'yktx7mx3esay', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 21:46:01.427597+00', '2025-08-17 21:46:01.427597+00', NULL, 'df7cd37e-0fb7-4233-bac3-44cb8455ef40'),
	('00000000-0000-0000-0000-000000000000', 25, 'ifbrjfwosgfh', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 21:48:17.560251+00', '2025-08-17 21:48:17.560251+00', NULL, 'b31998bf-d128-4ef0-a0de-dc97866c6147'),
	('00000000-0000-0000-0000-000000000000', 26, 'bnkocmn3odly', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-17 21:56:57.455559+00', '2025-08-17 21:56:57.455559+00', NULL, '559e7859-07f8-4713-b7d7-562ad666f0fa'),
	('00000000-0000-0000-0000-000000000000', 27, 'iz5n4heselu4', '5c0f44a8-7beb-45fd-81e2-023cdf7cde73', false, '2025-08-18 02:03:09.323931+00', '2025-08-18 02:03:09.323931+00', NULL, '836f28e5-039e-4058-91d0-67d7654e8cfe'),
	('00000000-0000-0000-0000-000000000000', 28, 's5yoeegiqdhi', '38977817-1066-40d9-ab5e-a8a8ab8e667d', false, '2025-08-18 02:37:50.027734+00', '2025-08-18 02:37:50.027734+00', NULL, '58f51ad0-bdce-40fe-83dc-5a86d5a023d4');


--
-- Data for Name: sso_providers; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: saml_providers; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: saml_relay_states; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: sso_domains; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: guilds; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO "public"."guilds" ("id", "discord_guild_id", "name", "active", "config", "created_at", "updated_at") VALUES
	(1, '1390537743385231451', 'Underrated Soccer League (USL)', true, '{"discord": {"bot_command_prefix": "!usl", "leaderboard_channel_id": null, "announcement_channel_id": null}, "permissions": {"admin_role_ids": [], "moderator_role_ids": []}}', '2025-08-14 21:03:23.287436+00', '2025-08-14 21:03:23.287436+00');


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: player_effective_mmr; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: player_historical_mmr; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: user_guild_memberships; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: user_trackers; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: usl_users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO "public"."usl_users" ("id", "name", "discord_id", "active", "banned", "mmr", "trueskill_mu", "trueskill_sigma", "trueskill_last_updated", "created_at", "updated_at") VALUES
	(1, 'oay', '544209988931944479', true, false, 0, 1998.650000, 7.632000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(3, 'yousif', '1159876823023353937', true, false, 0, 1997.490000, 6.617000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(4, 'Haiku', '753368582053953549', true, false, 0, 1995.240000, 6.767000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(5, 'Colt', '253702175530811394', true, false, 0, 1989.720000, 7.989000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(6, 'vareuew', '1008814187524407316', true, false, 0, 1988.000000, 5.826000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(7, 'kckiller', '542428031600427038', true, false, 0, 1983.900000, 7.138000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(8, 'Draconis', '679409909963554857', true, false, 0, 1982.860000, 7.449000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(9, 'q', '1017945853257850890', true, false, 0, 1975.930000, 6.133000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(10, 'sunkami', '692259430045319168', true, false, 0, 1974.600000, 5.333000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(11, 'Glassy', '1227487314221989942', true, false, 0, 1968.670000, 5.696000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(12, 'Winner P0V', '702327959213572157', true, false, 0, 1968.400000, 7.440000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(13, 'Dobby', '1377191739361857608', true, false, 0, 1954.900000, 7.725000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(14, 'Bola', '1138584884584120410', true, false, 0, 1942.580000, 6.639000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(15, 'dan', '527708595257606159', true, false, 0, 1932.160000, 7.027000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(16, 'Vièws.', '1303761806996803716', true, false, 0, 1923.800000, 7.254000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(17, 'Daniel', '807039329405763675', true, false, 0, 1918.760000, 7.213000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(18, 'Papi Acyx', '1118319407085666355', true, false, 0, 1913.230000, 7.874000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(19, 'frogahontas', '782482525939040296', true, false, 0, 1905.980000, 6.207000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(20, 'JoshuaBarnette7', '539902602310320155', true, false, 0, 1896.800000, 5.951000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(21, 'BOPPIN', '452653811115491349', true, false, 0, 1888.710000, 7.208000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(22, 'Castellan GT', '646100513942798340', true, false, 0, 1887.530000, 7.654000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(23, 'PELTvs', '525107026775375897', true, false, 0, 1883.760000, 6.860000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(24, 'dylan', '546836567176511528', true, false, 0, 1871.080000, 5.683000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(25, 'boogiegx', '593181293873856542', true, false, 0, 1867.880000, 7.395000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(26, 'Jesus502', '348128035502948372', true, false, 0, 1867.540000, 7.236000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(27, 'CrazyKid', '774739667567378442', true, false, 0, 1866.360000, 8.303000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(28, 'korbev', '778442295912824862', true, false, 0, 1850.950000, 7.693000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(29, 'Pxlse', '1201316903419904053', true, false, 0, 1845.340000, 6.361000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(30, 'BigChad', '721175279954821193', true, false, 0, 1824.070000, 6.184000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(31, 'mogtron', '354474826192388127', true, false, 0, 1805.780000, 7.874000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(32, 'Tauros', '1249438021833592903', true, false, 0, 1800.650000, 5.820000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(33, 'Slotha_ 🕊🕊🤍', '612537818811727873', true, false, 0, 1798.380000, 7.384000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(34, 'Req', '1121093280944377886', true, false, 0, 1776.360000, 5.876000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(35, 'Ryno', '1061036540454768700', true, false, 0, 1718.260000, 7.705000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(36, 'Qzr', '1377742398104404081', true, false, 0, 1665.570000, 7.311000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(37, 'lankzy', '1016033972272234496', true, false, 0, 1664.930000, 6.945000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(38, 'batman', '1198665474720923682', true, false, 0, 1652.560000, 7.307000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(39, 'Jeremy', '1091569905716973598', true, false, 0, 1650.520000, 7.902000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(40, 'streakyarc', '572607984333750292', true, false, 0, 1626.290000, 6.990000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(41, 'Nova', '1041658474595110934', true, false, 0, 1584.420000, 5.581000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(42, 'Sabotore', '662499562702766128', true, false, 0, 1574.570000, 6.051000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(43, 'Yams', '846154908310306826', true, false, 0, 1408.040000, 8.252000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(44, 'Oliver', '775481808347201586', true, false, 0, 1158.690000, 8.259000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(45, 'Kayio', '1200544340750106654', true, false, 0, 775.190000, 8.015000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(46, 'Rocky', '679038415576104971', true, false, 0, 709.510000, 6.720000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(47, 'Zapdos', '1029063532945350777', true, false, 0, 182.510000, 8.291000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(49, 'IodineOdin', '1082059903155327036', true, false, 0, 1854.400000, 5.563000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
	(50, 'eppflow', '1152631991284547625', true, false, 0, 1000.000000, 8.333333, NULL, '2025-08-16 18:27:05.360662+00', '2025-08-16 18:27:05.360662+00'),
	(48, 'i8rawnuggies', '996987225050992710', true, false, 0, 1000.000000, 8.330000, '2025-08-16', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:32:31.317642+00'),
	(51, 'feed', '318171715744169986', true, false, 0, 1000.000000, 8.333000, NULL, '2025-08-17 18:46:04.162513+00', '2025-08-17 18:46:04.162513+00'),
	(52, 'Jee', '647194725090066438', true, false, 0, 1000.000000, 8.333000, NULL, '2025-08-17 21:17:42.492023+00', '2025-08-17 21:17:42.492023+00'),
	(2, 'ayejoshy', '837466622670667776', true, false, 0, 1998.660000, 6.778000, '2025-08-18', '2025-08-16 18:25:56.836511+00', '2025-08-18 02:05:05.822681+00');


--
-- Data for Name: usl_user_trackers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO "public"."usl_user_trackers" ("id", "discord_id", "url", "ones_current_season_peak", "ones_previous_season_peak", "ones_all_time_peak", "ones_current_season_games_played", "ones_previous_season_games_played", "twos_current_season_peak", "twos_previous_season_peak", "twos_all_time_peak", "twos_current_season_games_played", "twos_previous_season_games_played", "threes_current_season_peak", "threes_previous_season_peak", "threes_all_time_peak", "threes_current_season_games_played", "threes_previous_season_games_played", "last_updated", "valid", "mmr", "created_at", "updated_at") VALUES
	(51, '544209988931944479', 'https://rocketleague.tracker.network/rocket-league/profile/steam/76561199021028007/overview', 1140, 1184, 1302, 0, 23, 1660, 2002, 2002, 0, 170, 1559, 1753, 1753, 0, 110, '2025-07-31', true, 1902, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(52, '1303761806996803716', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Vi%C3%A8ws./overview', 1047, 698, 1047, 55, 2, 1636, 1954, 1954, 49, 108, 1094, 1047, 1094, 79, 30, '2025-07-30', true, 1887, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(53, '837466622670667776', 'https://rocketleague.tracker.network/rocket-league/profile/epic/ayejoshy.TTV/overview', 1365, 1311, 1365, 351, 931, 1801, 1885, 1968, 728, 772, 1398, 1444, 1735, 5, 1, '2025-07-30', true, 1902, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(54, '679409909963554857', 'https://rocketleague.tracker.network/rocket-league/profile/steam/76561198021624143/overview', 917, 1178, 1214, 1, 58, 1034, 1883, 1883, 7, 189, 1071, 1265, 1305, 0, 127, '2025-07-24', true, 1900, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(55, '1159876823023353937', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Tempur./overview', 1150, 1155, 1155, 19, 89, 1704, 1762, 1762, 93, 383, 1690, 1571, 1572, 0, 90, '2025-07-23', true, 1902, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(56, '753368582053953549', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Haiku.189%20-%20TT/overview', 1046, 1064, 1064, 18, 22, 1620, 1725, 1725, 99, 63, 1751, 1577, 1751, 99, 207, '2025-07-28', true, 1902, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(57, '1008814187524407316', 'https://rocketleague.tracker.network/rocket-league/profile/epic/SungVare%20Woo/overview', 1025, 994, 1050, 16, 80, 1622, 1704, 1704, 99, 653, 1437, 1404, 1437, 77, 107, '2025-07-30', true, 1901, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(58, '253702175530811394', 'https://rocketleague.tracker.network/rocket-league/profile/epic/iColtRL/overview', 1048, 1156, 1251, 0, 7, 1535, 1696, 1826, 2, 100, 1224, 1370, 1783, 5, 30, '2025-07-25', true, 1901, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(59, '1152631991284547625', 'https://rocketleague.tracker.network/rocket-league/profile/epic/eppflo/overview', 1080, 1054, 1095, 17, 13, 1626, 1599, 1626, 220, 118, 1501, 1416, 1501, 16, 133, '2025-07-28', true, 1902, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(60, '593181293873856542', 'https://rocketleague.tracker.network/rocket-league/profile/epic/boogie5m/overview', 921, 962, 1049, 0, 7, 1456, 1574, 1576, 80, 404, 975, 1005, 1016, 0, 21, '2025-07-29', true, 1871, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(61, '1201316903419904053', 'https://rocketleague.tracker.network/rocket-league/profile/epic/TrulzyX/overview', 1005, 866, 1286, 11, 114, 1482, 1502, 1502, 109, 497, 1069, 852, 1069, 60, 39, '2025-07-25', true, 1864, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(62, '782482525939040296', 'https://rocketleague.tracker.network/rocket-league/profile/epic/frogahontas/overview', 952, 1009, 1009, 9, 107, 1461, 1500, 1500, 298, 384, 1007, 1058, 1058, 3, 48, '2025-07-28', true, 1883, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(63, '692259430045319168', 'https://rocketleague.tracker.network/rocket-league/profile/epic/sunkami_/overview', 964, 1010, 1051, 88, 181, 1437, 1489, 1491, 274, 452, 1202, 1246, 1246, 79, 210, '2025-07-22', true, 1899, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(64, '1017945853257850890', 'https://rocketleague.tracker.network/rocket-league/profile/xbl/Call%20XI/overview', 1070, 1041, 1176, 104, 67, 1581, 1480, 1606, 268, 212, 1263, 1187, 1452, 31, 7, '2025-07-28', true, 1899, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(65, '996987225050992710', 'https://rocketleague.tracker.network/rocket-league/profile/epic/60o0o0o7/overview', 931, 922, 959, 2, 0, 1443, 1461, 1648, 3, 1, 957, 986, 1017, 0, 0, '2025-07-30', false, NULL, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(66, '702327959213572157', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Winner%20P0V/overview', 951, 930, 1072, 6, 5, 1463, 1457, 1603, 69, 40, 1191, 1415, 1429, 90, 65, '2025-07-28', true, 1898, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(67, '1227487314221989942', 'https://rocketleague.tracker.network/rocket-league/profile/epic/%C4%9Elassy/overview', 914, 860, 914, 21, 27, 1583, 1452, 1583, 151, 422, 1417, 1536, 1536, 134, 367, '2025-07-25', true, 1897, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(68, '1118319407085666355', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Acyx_/overview', 894, 894, 894, 24, 0, 1437, 1437, 1437, 43, 0, 1118, 1118, 1118, 41, 0, '2025-07-28', true, 1886, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(69, '546836567176511528', 'https://rocketleague.tracker.network/rocket-league/profile/steam/76561199001930515/overview', 991, 1023, 1051, 135, 132, 1110, 1437, 1481, 267, 175, 1110, 922, 1110, 74, 22, '2025-07-27', true, 1877, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(70, '542428031600427038', 'https://rocketleague.tracker.network/rocket-league/profile/epic/kckillerxx/overview', 797, 813, 997, 0, 0, 1343, 1421, 1556, 36, 220, 1317, 1399, 1522, 10, 195, '2025-07-26', true, 1901, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(71, '1377191739361857608', 'https://rocketleague.tracker.network/rocket-league/profile/epic/MysticAria/overview', 831, 849, 900, 0, 0, 1246, 1350, 1652, 22, 49, 1217, 1300, 1321, 44, 88, '2025-07-25', true, 1896, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(72, '807039329405763675', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Daniel%20rikow/overview', 1007, 979, 1007, 74, 16, 1522, 1337, 1522, 155, 67, 1105, 1064, 1105, 10, 12, '2025-07-28', true, 1887, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(73, '527708595257606159', 'https://rocketleague.tracker.network/rocket-league/profile/epic/xDracO0/overview', 820, 836, 964, 0, 3, 1353, 1303, 1454, 103, 194, 1153, 1177, 1217, 50, 118, '2025-07-29', true, 1891, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(74, '721175279954821193', 'https://rocketleague.tracker.network/rocket-league/profile/epic/PastT6/overview', 986, 936, 986, 84, 211, 1416, 1243, 1406, 301, 151, 987, 937, 1071, 7, 22, '2025-07-28', true, 1854, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(75, '539902602310320155', 'https://rocketleague.tracker.network/rocket-league/profile/epic/JoshuaBarnette7/overview', 816, 860, 860, 0, 15, 1151, 1243, 1274, 144, 193, 1255, 1318, 1318, 86, 575, '2025-07-26', true, 1883, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(76, '1377742398104404081', 'https://rocketleague.tracker.network/rocket-league/profile/xbl/I%20Dont%20Shave599/overview', 799, 869, 869, 1, 22, 1175, 1234, 1246, 103, 431, 808, 704, 865, 17, 19, '2025-07-26', true, 1663, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(77, '1138584884584120410', 'https://rocketleague.tracker.network/rocket-league/profile/epic/66Bola/overview', 914, 914, 914, 13, 0, 1383, 1209, 1383, 250, 49, 1242, 1011, 1242, 145, 7, '2025-07-28', true, 1893, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(78, '774739667567378442', 'https://rocketleague.tracker.network/rocket-league/profile/epic/CrazyKid2579/overview', 736, 749, 790, 1, 0, 1131, 1202, 1257, 6, 9, 971, 983, 1032, 1, 0, '2025-07-30', true, 1878, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(79, '452653811115491349', 'https://rocketleague.tracker.network/rocket-league/profile/psn/FCBOPPIN/overview', 800, 818, 825, 0, 4, 1455, 1197, 1455, 184, 58, 1069, 974, 974, 70, 13, '2025-07-28', true, 1879, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(80, '646100513942798340', 'https://rocketleague.tracker.network/rocket-league/profile/xbl/CASTELLAN%20Gt/overview', 777, 791, 818, 0, 0, 1213, 1195, 1213, 215, 121, 875, 910, 920, 0, 5, '2025-07-28', true, 1882, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(81, '348128035502948372', 'https://rocketleague.tracker.network/rocket-league/profile/epic/jesusr502/overview', 735, 746, 909, 0, 0, 1225, 1169, 1291, 268, 111, 1115, 1000, 1115, 23, 0, '2025-07-26', true, 1878, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(82, '354474826192388127', 'https://rocketleague.tracker.network/rocket-league/profile/steam/76561198051701160/overview', 800, 799, 849, 31, 37, 1145, 1133, 1155, 32, 60, 1123, 1167, 1167, 123, 42, '2025-07-14', true, 1072, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:31:43.426809+00'),
	(84, '1249438021833592903', 'https://rocketleague.tracker.network/rocket-league/profile/psn/Tauros_301/overview', 816, 824, 824, 53, 41, 1165, 1125, 1165, 169, 648, 1048, 1150, 1150, 80, 168, '2025-07-29', true, 1860, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(85, '525107026775375897', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Peltdiditagain1/overview', 740, 739, 798, 0, 7, 993, 1110, 1260, 0, 7, 1158, 1154, 1224, 79, 1111, '2025-08-10', true, 1882, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(86, '572607984333750292', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Streakyboi1905/overview', 728, 746, 746, 1, 15, 1080, 1103, 1103, 52, 594, 832, 920, 920, 18, 74, '2025-07-26', true, 1766, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(87, '662499562702766128', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Sabotore/overview', 688, 669, 688, 35, 63, 1117, 1083, 1117, 213, 500, 899, 886, 899, 104, 267, '2025-07-27', true, 1655, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(88, '1061036540454768700', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Ryno%E3%83%AC/overview', 819, 836, 836, 2, 3, 1123, 1075, 1123, 40, 169, 1008, 982, 1046, 10, 34, '2025-07-31', true, 1836, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(89, '1091569905716973598', 'https://rocketleague.tracker.network/rocket-league/profile/psn/Jemgent1/overview', 679, 679, 728, 3, 0, 1034, 1034, 1034, 185, 0, 917, 917, 924, 3, 0, '2025-07-27', true, 1812, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(90, '1121093280944377886', 'https://rocketleague.tracker.network/rocket-league/profile/epic/ReqRL/overview', 877, 744, 877, 87, 11, 1107, 1033, 1107, 438, 275, 1090, 700, 1090, 19, 0, '2025-07-27', true, 1855, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(91, '1016033972272234496', 'https://rocketleague.tracker.network/rocket-league/profile/epic/jnko./overview', 759, 747, 759, 25, 58, 1144, 1025, 1144, 185, 341, 760, 678, 760, 6, 2, '2025-07-28', true, 1792, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(92, '778442295912824862', 'https://rocketleague.tracker.network/rocket-league/profile/epic/korbev%20%E3%83%83/overview', 822, 756, 823, 13, 0, 1238, 1006, 1212, 187, 0, 971, 859, 971, 6, 0, '2025-07-30', true, 1870, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(93, '1041658474595110934', 'https://rocketleague.tracker.network/rocket-league/profile/psn/meek_s_/overview', 775, 763, 775, 206, 206, 1137, 1004, 1137, 561, 366, 829, 818, 829, 7, 32, '2025-07-26', true, 1705, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(94, '846154908310306826', 'https://rocketleague.tracker.network/rocket-league/profile/epic/TheTypicalOwl/overview', 714, 724, 823, 0, 0, 903, 904, 1084, 25, 27, 842, 862, 1101, 0, 4, '2025-07-30', true, 1564, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(95, '1198665474720923682', 'https://rocketleague.tracker.network/rocket-league/profile/xbl/Aranara2boy/overview', 1006, 713, 1006, 2, 22, 1432, 883, 1432, 163, 147, 839, 901, 901, 24, 65, '2025-07-28', true, 1755, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(96, '775481808347201586', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Firz.RL/overview', 657, 677, 724, 3, 0, 801, 769, 851, 20, 14, 589, 589, 598, 0, 0, '2025-07-25', true, 907, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(97, '679038415576104971', 'https://rocketleague.tracker.network/rocket-league/profile/epic/Rocky%20The%20II/overview', 582, 582, 582, 111, 111, 661, 689, 689, 190, 190, 505, 505, 533, 68, 68, '2025-07-24', true, 164, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(98, '1029063532945350777', 'https://rocketleague.tracker.network/rocket-league/profile/epic/m%C3%B8ltr%C3%ABs/overview', 474, 476, 512, 1, 14, 0, 495, 555, 9, 28, 0, 568, 570, 0, 1, '2025-07-26', true, 6, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(99, '1200544340750106654', 'https://rocketleague.tracker.network/rocket-league/profile/epic/KxiyoRL/overview', 435, 256, 612, 15, 0, 824, 319, 824, 267, 15, 497, 489, 497, 59, 2, '2025-07-26', true, 492, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(100, '1082059903155327036', 'https://rocketleague.tracker.network/rocket-league/profile/epic/IodineOdin/overview', 798, 782, 825, 8, 0, 1231, 1169, 1234, 323, 120, 940, 1176, 1234, 131, 200, '2025-08-13', true, 1854, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:27:12.98033+00'),
	(83, '612537818811727873', 'https://rocketleague.tracker.network/rocket-league/profile/epic/slotha_%20-washed-/overview', 1002, 1049, 1117, 0, 11, 1673, 1761, 1761, 67, 252, 1491, 1560, 1560, 9, 100, '2025-07-15', true, 1651, '2025-08-16 18:27:12.98033+00', '2025-08-16 18:30:03.720271+00');


--
-- Data for Name: buckets; Type: TABLE DATA; Schema: storage; Owner: supabase_storage_admin
--



--
-- Data for Name: objects; Type: TABLE DATA; Schema: storage; Owner: supabase_storage_admin
--



--
-- Data for Name: s3_multipart_uploads; Type: TABLE DATA; Schema: storage; Owner: supabase_storage_admin
--



--
-- Data for Name: s3_multipart_uploads_parts; Type: TABLE DATA; Schema: storage; Owner: supabase_storage_admin
--



--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE SET; Schema: auth; Owner: supabase_auth_admin
--

SELECT pg_catalog.setval('"auth"."refresh_tokens_id_seq"', 28, true);


--
-- Name: guilds_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('"public"."guilds_id_seq"', 1, true);


--
-- Name: player_effective_mmr_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('"public"."player_effective_mmr_id_seq"', 1, false);


--
-- Name: player_historical_mmr_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('"public"."player_historical_mmr_id_seq"', 1, false);


--
-- Name: user_guild_memberships_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('"public"."user_guild_memberships_id_seq"', 1, false);


--
-- Name: user_trackers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('"public"."user_trackers_id_seq"', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('"public"."users_id_seq"', 1, false);


--
-- Name: usl_user_trackers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('"public"."usl_user_trackers_id_seq"', 100, true);


--
-- Name: usl_users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('"public"."usl_users_id_seq"', 52, true);


--
-- PostgreSQL database dump complete
--

RESET ALL;
