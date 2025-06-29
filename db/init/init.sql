--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9 (Ubuntu 16.9-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.9 (Ubuntu 16.9-0ubuntu0.24.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: portfolioassets; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.portfolioassets (
    asset_id integer NOT NULL,
    portfolio_id integer NOT NULL,
    symbol character varying(20) NOT NULL,
    quantity numeric(18,8) NOT NULL,
    average_buy_price numeric(18,8),
    last_updated timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.portfolioassets OWNER TO postgres;

--
-- Name: portfolioassets_asset_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.portfolioassets_asset_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.portfolioassets_asset_id_seq OWNER TO postgres;

--
-- Name: portfolioassets_asset_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.portfolioassets_asset_id_seq OWNED BY public.portfolioassets.asset_id;


--
-- Name: portfolios; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.portfolios (
    portfolio_id integer NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.portfolios OWNER TO postgres;

--
-- Name: portfolios_portfolio_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.portfolios_portfolio_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.portfolios_portfolio_id_seq OWNER TO postgres;

--
-- Name: portfolios_portfolio_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.portfolios_portfolio_id_seq OWNED BY public.portfolios.portfolio_id;


--
-- Name: resettokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resettokens (
    id integer NOT NULL,
    token character varying(50) NOT NULL,
    email character varying(50) NOT NULL,
    expdate timestamp without time zone NOT NULL
);


ALTER TABLE public.resettokens OWNER TO postgres;

--
-- Name: resetTokens_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."resetTokens_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."resetTokens_id_seq" OWNER TO postgres;

--
-- Name: resetTokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."resetTokens_id_seq" OWNED BY public.resettokens.id;


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    session_id text NOT NULL,
    user_id uuid NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions (
    transaction_id integer NOT NULL,
    portfolio_id integer NOT NULL,
    symbol character varying(20) NOT NULL,
    trade_type character varying(10) NOT NULL,
    quantity numeric(18,8) NOT NULL,
    price_per_unit numeric(18,8) NOT NULL,
    total_amount numeric(18,8) NOT NULL,
    transaction_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT transactions_trade_type_check CHECK (((trade_type)::text = ANY ((ARRAY['BUY'::character varying, 'SELL'::character varying])::text[])))
);


ALTER TABLE public.transactions OWNER TO postgres;

--
-- Name: transactions_transaction_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.transactions_transaction_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.transactions_transaction_id_seq OWNER TO postgres;

--
-- Name: transactions_transaction_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.transactions_transaction_id_seq OWNED BY public.transactions.transaction_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    username character varying(50) NOT NULL,
    password character varying(100) NOT NULL,
    name character varying(50) NOT NULL,
    surname character varying(50) NOT NULL,
    email character varying(75) NOT NULL,
    role character varying(25) NOT NULL,
    pp character varying(150),
    created_at timestamp without time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: portfolioassets asset_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolioassets ALTER COLUMN asset_id SET DEFAULT nextval('public.portfolioassets_asset_id_seq'::regclass);


--
-- Name: portfolios portfolio_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolios ALTER COLUMN portfolio_id SET DEFAULT nextval('public.portfolios_portfolio_id_seq'::regclass);


--
-- Name: resettokens id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resettokens ALTER COLUMN id SET DEFAULT nextval('public."resetTokens_id_seq"'::regclass);


--
-- Name: transactions transaction_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions ALTER COLUMN transaction_id SET DEFAULT nextval('public.transactions_transaction_id_seq'::regclass);


--
-- Data for Name: portfolioassets; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.portfolioassets (asset_id, portfolio_id, symbol, quantity, average_buy_price, last_updated) FROM stdin;
3	2	BTCUSD	1.00000000	107262.75000000	2025-06-28 18:00:40.354699
2	1	BTCUSD	2.00000000	107569.16000000	2025-06-27 19:07:42.29942
1	1	AAPL	120.00000000	201.08000000	2025-06-27 03:03:44.343895
\.


--
-- Data for Name: portfolios; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.portfolios (portfolio_id, user_id, created_at, updated_at) FROM stdin;
1	bf3befab-388a-46b8-9944-204e326e8697	2025-06-27 02:57:29.122191	2025-06-27 02:57:29.122191
2	60b234f5-e3de-4942-ab96-4b748dc217da	2025-06-28 18:00:40.344316	2025-06-28 18:00:40.344316
\.


--
-- Data for Name: resettokens; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.resettokens (id, token, email, expdate) FROM stdin;
1	c--4qv7tICLfX7YI6BCrADxxHMhHKmrXNvfwPejL1WE=	batuhanalun1999@hotmail.com	2025-05-27 00:00:00
2	S_gsLvypmMOM3R16L0qGRVI0Lv7B9l22Voi9XgdSBcw=	batuhanalun1999@hotmail.com	2025-05-27 00:00:00
3	Y2jPFww9IMmnEH44IcZEX8aospL0BYLiYEV95HqMxVI=	batuhanalun1999@hotmail.com	2025-05-27 00:00:00
4	zJC-Dct9sgKVRcwznknREGlgLwHd8_U4ulUvGjAO3n4=	batuhanalun1999@hotmail.com	2025-05-27 00:00:00
5	8INiXLQn8BiXPq0BYLR6SlffFGBZuMwb1Jult4MRkk0=	batuhanalun1999@hotmail.com	2025-05-27 00:00:00
6	KZZ3nZZNNnUWZI1mXw6FrNfNJGsTZuUUyHIEDwsrjCA=	batuhanalun1999@hotmail.com	2025-05-27 00:00:00
7	dr3N-CV5wXuSWcYtNOhL8YrAweCohWXyu8SOwRruYEw=	batuhanalun1999@hotmail.com	2025-05-27 00:00:00
8	5gb0tY882wPCH_8NAiCRkX-p1ZxnD7FczpnXfbJz0vo=	batuhanalun1999@hotmail.com	2025-05-27 00:00:00
9	qTKWdriYHvN8nuDKz17Q4wucJ8IupLYfUEPVbO8DBmU=	batuhanalun1999@hotmail.com	2025-05-27 00:00:00
10	6q75RO006mqJ7kjU7kTfZ2FrHVQJiG_BS0bvbopYJ-M=	batuhanalun199@hotmail.com	2025-06-28 00:00:00
12	d29nS3hGdgwq51pRPjnBCAgVn2KjE8XTATP8YJot1Rg=	batuhanalun199@hotmail.com	2025-06-28 17:40:41.687482
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sessions (session_id, user_id, expires_at, created_at) FROM stdin;
6yupC_Vl2PxaFAEdLEeW1p0AYwNudYCN9Xpyw0cYLXQ=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 22:03:56.262775	2025-06-22 22:03:56.263012
5ALoQSZPL9y2ELM5wKe_jwTrVVgXBcsGtLjX_VDJzCA=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 22:54:49.802404	2025-06-22 22:54:49.802673
KdS9eqo6Zp8S2aBD9N_Ax8sJchUeRvfndI2jMv19QdY=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 22:55:06.126799	2025-06-22 22:55:06.127004
EeIEEl5izIdGzGUrPLfARUa1nsxWXcr2JWwO0UNqOkQ=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 22:56:03.196081	2025-06-22 22:56:03.19623
LiF6cHiazTCpqcEibB1-Tz8G0cYC2kFrxeTEL9LgiMo=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 22:59:29.295304	2025-06-22 22:59:29.295609
mAe9wI0X1psmwdLhTZi6yDob8UFBkwrUBYUhKmm0PiU=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 22:59:42.98006	2025-06-22 22:59:42.980219
6I0H1zR2CnuecMkDEmCMDLRWKdzx3jZjPheN8W9I0GM=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:05:00.797116	2025-06-22 23:05:00.797419
VqAGNa_MHffMNnK0VjgOm8vb1FViSzBkKr7nQhiVZuA=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:05:17.084108	2025-06-22 23:05:17.084249
3tSXo76zrwB7boEjvZcEVpo5CBk2mWkx0E1In6bxYHk=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:12:02.684966	2025-06-22 23:12:02.685206
rZwYxdzF7aFSrBvdqqqII-VE8O5c4uUu27o_KCvmuiA=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:17:02.845092	2025-06-22 23:17:02.845429
pH2f-vPbk_RyxD3Bu1j_TMmrgkwzS_-DoiW6jvXwgdw=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:17:22.684665	2025-06-22 23:17:22.6848
07CBx87SikCmG9R-O7Uvi8ClQHAeVigl7vRoeG72tmg=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:17:42.557544	2025-06-22 23:17:42.557698
3e6vpURTEc-Wzn72zOlKqvnnB2Xh1McY_q7SQRc4Guw=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:18:25.579724	2025-06-22 23:18:25.580036
QB00LBZPWqrxYYxph2_8vpjrHAihVoqI-IP4wDaDvto=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:22:13.386139	2025-06-22 23:22:13.38628
c53lRElYiCvUW3jvc-jIFMMGwECgRln9jiebg6ivCs8=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:22:19.115477	2025-06-22 23:22:19.11564
hWh--aRWi4vNGIJRL7k1DftxdIVjjOFdiI1OMQjpGlk=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:23:44.198269	2025-06-22 23:23:44.198515
6FugQLZi0ondd0z2RNwDAttSsdU1U040s5J4RNLrZGQ=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-23 23:24:07.600233	2025-06-22 23:24:07.600365
8jktj4oxkRD8HWI9CHQFwptbAT00f_J0IftE3KlCOQo=	da71cee6-3520-4f5b-87f2-8edf17470998	2025-06-24 14:58:50.826172	2025-06-23 14:58:50.826797
guXM36BHpEgMY29QwUbFWHLDOSmz0biF6asPIBDsILw=	da71cee6-3520-4f5b-87f2-8edf17470998	2025-06-24 14:59:22.395385	2025-06-23 14:59:22.395646
jur_IC34nvtUEgHvISWKrLGwrmEWm_fMMkv7iC7oIIQ=	da71cee6-3520-4f5b-87f2-8edf17470998	2025-06-24 15:00:41.090254	2025-06-23 15:00:41.090413
bTMZNK7h0rnGkGzrXZrFCqxFIdjGZLp_QI0Gk5Gp7EA=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-24 15:01:56.126496	2025-06-23 15:01:56.126737
2pgx2Tw_w59hkXjglaVtHq22_XfGVoztQjWRTCwmq68=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:41:45.488755	2025-06-24 13:41:45.488995
IZG0ZuGtkhkCOP5FQ1ZwiS3dU7pi2TjTdkYNko2gGSI=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:44:18.154003	2025-06-24 13:44:18.15426
ANC70_KpHEGIj0avGJt6yCWS8q6Sct0V-qG-hKXyGlc=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:44:34.772032	2025-06-24 13:44:34.772194
HjaVcoEBPatGR5xzYwB_Dh-szNDN41eAcUML7Gc8MqI=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:45:58.590501	2025-06-24 13:45:58.590732
DFM45u9JqSw5QgBxVjK4YYJCzQLNorardJMxdaGlASk=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:47:49.048022	2025-06-24 13:47:49.048158
ag6Xd-okD1HhuKxUzF5683jaRA3qVczW1Q5knl2-0XM=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:49:42.424164	2025-06-24 13:49:42.42445
eta4WCMZkoRQ5pHTekJvUgY3XdULrKyWmM3EYW3vkCg=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:53:10.182355	2025-06-24 13:53:10.182606
M69of31K5_csrxwq6tkp9AhtymrLcKVr13_GFQVOJ-0=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:53:56.699244	2025-06-24 13:53:56.69948
Vd0ne5EFYwJTgvMkyrfPd34DbWOrk3WbDlH0uc9HGQc=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:54:37.435928	2025-06-24 13:54:37.436166
JF0u4SUCizti3OSzTbAa-z7V97CTwQ7h6v9Hhf2n3vc=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:56:01.057086	2025-06-24 13:56:01.057378
k8n8f62W05acZvKwGv72JzN1UB48F5gamSy29n-ZPxI=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 13:59:13.014271	2025-06-24 13:59:13.014764
o3myVb43Ttj7Xkamm5mEPqkXfyAtNgETjRTf-pRlfQc=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 14:02:08.643583	2025-06-24 14:02:08.643908
7Nsc8mgYkyMXb54ab7MO8wcLS4NBiHXNntcs3GDDsts=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 20:22:58.380102	2025-06-24 20:22:58.380773
swKMataXMCBeSWj1wwScaac-tmC_G_-bQIpmbjgYB2M=	798f251e-06ff-40f5-b50c-dbc3305ac1b9	2025-06-25 20:58:47.27445	2025-06-24 20:58:47.274634
lPxKFRO_iVJshLDSXSSbmKAjaA_vbNuy057kluJzh10=	da71cee6-3520-4f5b-87f2-8edf17470998	2025-06-25 21:19:09.335211	2025-06-24 21:19:09.335424
CZkcxBZR89cQXmjXDxTH1eZfa4v2dt8fDEjBX0zoXrw=	da71cee6-3520-4f5b-87f2-8edf17470998	2025-06-25 21:35:44.132157	2025-06-24 21:35:44.132414
RfS6Aync4ixqiaHt6FR95zhfjmOnu8o-1MCVpR5EmWE=	da71cee6-3520-4f5b-87f2-8edf17470998	2025-06-25 21:46:05.186486	2025-06-24 21:46:05.186719
aTMQbrGuafJ9OjjAmTtMRd0rodNvseVdXPmgaccBz2M=	da71cee6-3520-4f5b-87f2-8edf17470998	2025-06-25 21:54:56.388276	2025-06-24 21:54:56.388544
FXI2IrOGAQtOt6OpeTY5BMSDuPOPUiT2QtpZZlrZRkw=	778a6f9b-6c8c-4154-92c6-9e8aae5e0109	2025-06-25 22:55:40.702529	2025-06-24 22:55:40.702767
I9OdW3LxSnBk400AS-25WhJ9Oz_-rT6ECeDABvGN4mA=	778a6f9b-6c8c-4154-92c6-9e8aae5e0109	2025-06-26 13:31:56.817616	2025-06-25 13:31:56.818863
eXiZqzTey33b3Bc7JOgdPsrNTvhXbiXGJ1nFQ_kLTrc=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-26 14:39:18.169852	2025-06-25 14:39:18.17017
NADkWi5GtHBkC4GgHIfKel1ktpD_XAL-fYyU4pVgZIE=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-26 14:48:21.513029	2025-06-25 14:48:21.513259
XI-urRHB17txVL7GVym67MYIsmUWyjbvtHAKyNJDFnE=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-26 14:49:36.847092	2025-06-25 14:49:36.847279
nr8jsO6lMdNj7MR7oco2OhyQiq7gK02YB5Cn9gkRbMs=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-26 15:08:32.642467	2025-06-25 15:08:32.642762
QoPrVlT3Gt5PnyOk8HIMggznkm2nLad68dOK3wEAY70=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-27 15:03:23.652415	2025-06-26 15:03:23.653453
ROLVVUgujtGG9QWCPuqOK89E6mUfEckqtvMKbmZZMTg=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-28 19:05:20.980776	2025-06-27 19:05:20.981517
d_xjjfDeyjGzGNAr1x1zvqeZHcM5EgLQep_UEytVp_Q=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-28 20:29:34.093132	2025-06-27 20:29:34.093312
Auvehcacd1YVU1XtAOjcjsn03ognfBWhTjs_9Z2FynA=	60b234f5-e3de-4942-ab96-4b748dc217da	2025-06-29 16:35:34.468609	2025-06-28 16:35:34.468905
y9cgorm6kZ7uwLjMsh5xp0Va8_jSh2wu9BIhzCqPg3E=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-29 17:17:58.630976	2025-06-28 17:17:58.631284
yeb9tSeG-rMz1y3sSLqyUyR4pyQvW-qcCM11wla-1Y8=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-30 11:49:43.082979	2025-06-29 11:49:43.084082
Wa2vA-ljyl2219bPgEJdnB9dmyqR8blnBrFnT1yh4W0=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-30 12:18:04.319235	2025-06-29 12:18:04.319601
wL1im528AOb6OAZ-78qUYpgbVzJ2Y9Rbovj_xjaBVGE=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-30 12:37:07.625715	2025-06-29 12:37:07.626036
-OuRySwnBX-DQVUy4hhL2AHUNfQMXs6e4C_kaUE7PWc=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-30 12:37:40.595357	2025-06-29 12:37:40.595679
IIk1HGykdRhlinJXXQk0RE4CQ9ta5wHcQ9ExBEgA4dU=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-30 12:38:14.394358	2025-06-29 12:38:14.394645
DLeVN2hDgoB10mfdxTxwJrj1xzZQYXzSapi-sAPa_8E=	bf3befab-388a-46b8-9944-204e326e8697	2025-06-30 12:40:26.253579	2025-06-29 12:40:26.253916
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (transaction_id, portfolio_id, symbol, trade_type, quantity, price_per_unit, total_amount, transaction_date) FROM stdin;
3	1	AAPL	BUY	10.00000000	201.00000000	2010.00000000	2025-06-27 03:03:44.342203
4	1	BTCUSD	BUY	1.00000000	107421.32000000	107421.32000000	2025-06-27 19:07:42.292446
6	1	AAPL	BUY	100.00000000	202.24000000	20224.00000000	2025-06-27 20:23:08.545748
7	2	BTCUSD	BUY	1.00000000	107262.75000000	107262.75000000	2025-06-28 18:00:40.348999
11	1	BTCUSD	BUY	1.00000000	107717.00000000	107717.00000000	2025-06-29 12:04:05.937379
13	1	AAPL	SELL	10.00000000	201.08000000	2010.80000000	2025-06-29 12:09:27.553947
14	1	AAPL	SELL	10.00000000	201.08000000	2010.80000000	2025-06-29 12:19:14.220693
15	1	AAPL	SELL	10.00000000	201.08000000	2010.80000000	2025-06-29 12:31:45.566387
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, password, name, surname, email, role, pp, created_at) FROM stdin;
7c73493e-1296-4346-8bb7-d5db5039c621	sa	$2a$10$XFWZuPYORHEwcZfnieOnB.TSGwPWfdDWyG1ipgWQ/bCUPiBSob0YO	sa	sa	sa@sa.com	user	\N	\N
4b5deeaa-65cf-4033-8b6b-c56ae078cf05	batualun	$2a$10$7AFgtwNG1EJmbLEaVhkvwe11FYfk6oQIRcHbsA6wtuNAQ0OiDHNy6	Batu	Alun	batuhanalun1999@hotmail.com	user	\N	\N
a80c423e-44f0-48e3-9fd2-974a5d9611d1	asd	$2a$10$mkuJnBbSteWOq5XRCthZpuYgIqTTPacKIvjf072mEXCI6HgCsPfLi	sda	asd	sa@hotmail.com	user	\N	\N
fca9e2a6-2519-49c9-933e-6fb3cb5227c9	asa	$2a$10$8C4/s51SM7F6I80hUvrvDuCPVx4NOHjYe2jASVP3P/wZzZJncIW7u	sa	s	asdasd@asdasd	user	\N	\N
da71cee6-3520-4f5b-87f2-8edf17470998	deneme1	$2a$10$Z186ebiq0klaCyPZxaK0x.fGNgHllILw.Rc8Ax.FQZr/UXw/m8Hd2	deneme1	deneme1	deneme1@example.com	user	\N	\N
8672b683-51b8-45b7-85da-d409d02d3a64	asdasd	$2a$10$LgAHW6hwE/jqD4lwnD9WauRHlks9e7Umj.pvHumg9ELSZIAH8jnMm	asda	asda	asdads@example.com	user	../../../static/avatars/fbc0bd59-8ce0-4d7a-bde9-52c6fa0aa64f.png	2025-06-24 22:54:40.882705
efb4fc2c-6dda-496b-a5cc-a0c32308f446	deneme	$2a$10$eY3kGe0vbvWFDMcpHZxt3usx/n9y0u3dGUy39oHtWJ46dYrU3gyFy	deneme	deneme	deneme@example.com	user	../../../static/avatars/c9a793dc-6087-4a42-ab44-69c5c9d9d29e.png	2025-06-25 14:32:23.299861
a5dc780f-9a10-42a2-a8b9-1aefbed06cf0	sdasda	$2a$10$qDEIBAkbZdA3VaFEl7ryjOQx5Tyxp78I8dItzNqyJaN14kXMmX2NC	adsa	sda	asdasd@example.com	user	../../../static/default-avatar.png	2025-06-25 14:33:50.662195
60b234f5-e3de-4942-ab96-4b748dc217da	batuhanalun	$2a$10$WjYolSwXvh9kRn1BjdpIYO896q7SuFesemas2cfgJbHfAouszAmLK	batuhan	alun	batuhanalun199@hotmail.com	user	../../../static/default-avatar.png	2025-06-28 16:35:25.766483
bf3befab-388a-46b8-9944-204e326e8697	deneme12	$2a$10$SqF6UY7gRECsu/dHLoSFbuXDN4POl2EBPBynRX7Z0d/tzMPEgtWqu	deneme12	deneme	deneme12@example.com	user	./static/avatars/e1070f48-ec5d-4cc0-b555-0623b171b2dd.png	2025-06-25 14:39:10.529927
\.


--
-- Name: portfolioassets_asset_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.portfolioassets_asset_id_seq', 3, true);


--
-- Name: portfolios_portfolio_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.portfolios_portfolio_id_seq', 2, true);


--
-- Name: resetTokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."resetTokens_id_seq"', 12, true);


--
-- Name: transactions_transaction_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.transactions_transaction_id_seq', 15, true);


--
-- Name: portfolioassets portfolioassets_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolioassets
    ADD CONSTRAINT portfolioassets_pkey PRIMARY KEY (asset_id);


--
-- Name: portfolioassets portfolioassets_portfolio_id_symbol_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolioassets
    ADD CONSTRAINT portfolioassets_portfolio_id_symbol_key UNIQUE (portfolio_id, symbol);


--
-- Name: portfolios portfolios_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolios
    ADD CONSTRAINT portfolios_pkey PRIMARY KEY (portfolio_id);


--
-- Name: portfolios portfolios_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolios
    ADD CONSTRAINT portfolios_user_id_key UNIQUE (user_id);


--
-- Name: resettokens resetTokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resettokens
    ADD CONSTRAINT "resetTokens_pkey" PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (session_id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (transaction_id);


--
-- Name: users usernameUnique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT "usernameUnique" UNIQUE (username, email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: portfolioassets portfolioassets_portfolio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolioassets
    ADD CONSTRAINT portfolioassets_portfolio_id_fkey FOREIGN KEY (portfolio_id) REFERENCES public.portfolios(portfolio_id);


--
-- Name: portfolios portfolios_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.portfolios
    ADD CONSTRAINT portfolios_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: transactions transactions_portfolio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_portfolio_id_fkey FOREIGN KEY (portfolio_id) REFERENCES public.portfolios(portfolio_id);


--
-- PostgreSQL database dump complete
--

