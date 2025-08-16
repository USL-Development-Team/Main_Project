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
	('00000000-0000-0000-0000-000000000000', 'a4d8e62d-3c58-4a84-b3e8-3780549d69b0', '{"action":"user_signedup","actor_id":"32481211-bd29-4581-9b9a-50ce02924dff","actor_name":"mogtron","actor_username":"reilly.kyle101@gmail.com","actor_via_sso":false,"log_type":"team","traits":{"provider":"discord"}}', '2025-08-16 18:27:46.590086+00', '');


--
-- Data for Name: flow_state; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: users; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--

INSERT INTO "auth"."users" ("instance_id", "id", "aud", "role", "email", "encrypted_password", "email_confirmed_at", "invited_at", "confirmation_token", "confirmation_sent_at", "recovery_token", "recovery_sent_at", "email_change_token_new", "email_change", "email_change_sent_at", "last_sign_in_at", "raw_app_meta_data", "raw_user_meta_data", "is_super_admin", "created_at", "updated_at", "phone", "phone_confirmed_at", "phone_change", "phone_change_token", "phone_change_sent_at", "email_change_token_current", "email_change_confirm_status", "banned_until", "reauthentication_token", "reauthentication_sent_at", "is_sso_user", "deleted_at", "is_anonymous") VALUES
	('00000000-0000-0000-0000-000000000000', '32481211-bd29-4581-9b9a-50ce02924dff', 'authenticated', 'authenticated', 'reilly.kyle101@gmail.com', NULL, '2025-08-16 18:27:46.590661+00', NULL, '', NULL, '', NULL, '', '', NULL, '2025-08-16 18:27:46.591902+00', '{"provider": "discord", "providers": ["discord"]}', '{"iss": "https://discord.com/api", "sub": "354474826192388127", "name": "mogtron#0", "email": "reilly.kyle101@gmail.com", "picture": "https://cdn.discordapp.com/avatars/354474826192388127/c50ea819dfdc72747e00c6f67d4ade97.png", "full_name": "mogtron", "avatar_url": "https://cdn.discordapp.com/avatars/354474826192388127/c50ea819dfdc72747e00c6f67d4ade97.png", "provider_id": "354474826192388127", "custom_claims": {"global_name": "mogtron"}, "email_verified": true, "phone_verified": false}', NULL, '2025-08-16 18:27:46.584892+00', '2025-08-16 18:27:46.594037+00', NULL, NULL, '', '', NULL, '', 0, NULL, '', NULL, false, NULL, false);


--
-- Data for Name: identities; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--

INSERT INTO "auth"."identities" ("provider_id", "user_id", "identity_data", "provider", "last_sign_in_at", "created_at", "updated_at", "id") VALUES
	('354474826192388127', '32481211-bd29-4581-9b9a-50ce02924dff', '{"iss": "https://discord.com/api", "sub": "354474826192388127", "name": "mogtron#0", "email": "reilly.kyle101@gmail.com", "picture": "https://cdn.discordapp.com/avatars/354474826192388127/c50ea819dfdc72747e00c6f67d4ade97.png", "full_name": "mogtron", "avatar_url": "https://cdn.discordapp.com/avatars/354474826192388127/c50ea819dfdc72747e00c6f67d4ade97.png", "provider_id": "354474826192388127", "custom_claims": {"global_name": "mogtron"}, "email_verified": true, "phone_verified": false}', 'discord', '2025-08-16 18:27:46.587667+00', '2025-08-16 18:27:46.58769+00', '2025-08-16 18:27:46.58769+00', 'da2b13a8-b077-4ba7-bffd-3c26bd3f74bd');


--
-- Data for Name: instances; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--



--
-- Data for Name: sessions; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--

INSERT INTO "auth"."sessions" ("id", "user_id", "created_at", "updated_at", "factor_id", "aal", "not_after", "refreshed_at", "user_agent", "ip", "tag") VALUES
	('1287365f-bf71-4d8f-abe1-87b7dcbf5649', '32481211-bd29-4581-9b9a-50ce02924dff', '2025-08-16 18:27:46.591934+00', '2025-08-16 18:27:46.591934+00', NULL, 'aal1', NULL, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36', '172.20.0.1', NULL);


--
-- Data for Name: mfa_amr_claims; Type: TABLE DATA; Schema: auth; Owner: supabase_auth_admin
--

INSERT INTO "auth"."mfa_amr_claims" ("session_id", "created_at", "updated_at", "authentication_method", "id") VALUES
	('1287365f-bf71-4d8f-abe1-87b7dcbf5649', '2025-08-16 18:27:46.594377+00', '2025-08-16 18:27:46.594377+00', 'oauth', 'aa71d4eb-81cc-4a26-8753-b51d957a8d91');


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
	('00000000-0000-0000-0000-000000000000', 1, 'en547ij3x7li', '32481211-bd29-4581-9b9a-50ce02924dff', false, '2025-08-16 18:27:46.592977+00', '2025-08-16 18:27:46.592977+00', NULL, '1287365f-bf71-4d8f-abe1-87b7dcbf5649');


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
	(1, '1390537743385231451', 'Underrated Soccer League (USL)', true, '{"discord": {"bot_command_prefix": "!usl", "leaderboard_channel_id": null, "announcement_channel_id": null}, "permissions": {"admin_role_ids": [], "moderator_role_ids": []}}', '2025-08-16 17:34:28.088855+00', '2025-08-16 17:34:28.088855+00');


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
	(2, 'ayejoshy', '837466622670667776', true, false, 0, 1998.540000, 5.082000, '2025-08-13', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:25:56.836511+00'),
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
	(48, 'i8rawnuggies', '996987225050992710', true, false, 0, 1000.000000, 8.330000, '2025-08-16', '2025-08-16 18:25:56.836511+00', '2025-08-16 18:32:31.317642+00');


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
-- Data for Name: prefixes; Type: TABLE DATA; Schema: storage; Owner: supabase_storage_admin
--



--
-- Data for Name: s3_multipart_uploads; Type: TABLE DATA; Schema: storage; Owner: supabase_storage_admin
--



--
-- Data for Name: s3_multipart_uploads_parts; Type: TABLE DATA; Schema: storage; Owner: supabase_storage_admin
--



--
-- Data for Name: hooks; Type: TABLE DATA; Schema: supabase_functions; Owner: supabase_functions_admin
--



--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE SET; Schema: auth; Owner: supabase_auth_admin
--

SELECT pg_catalog.setval('"auth"."refresh_tokens_id_seq"', 1, true);


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

SELECT pg_catalog.setval('"public"."usl_users_id_seq"', 50, true);


--
-- Name: hooks_id_seq; Type: SEQUENCE SET; Schema: supabase_functions; Owner: supabase_functions_admin
--

SELECT pg_catalog.setval('"supabase_functions"."hooks_id_seq"', 1, false);


--
-- PostgreSQL database dump complete
--

RESET ALL;
