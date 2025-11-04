--
-- PostgreSQL database dump
--


-- Dumped from database version 17.6
-- Dumped by pg_dump version 17.6

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

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO postgres;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS '';


--
-- Name: ProductState; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."ProductState" AS ENUM (
    'NEW',
    'USED'
);


ALTER TYPE public."ProductState" OWNER TO postgres;

--
-- Name: ProfileType; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."ProfileType" AS ENUM (
    'INDIVIDUAL',
    'OOO',
    'IP'
);


ALTER TYPE public."ProfileType" OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: Category; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Category" (
    id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public."Category" OWNER TO postgres;

--
-- Name: Category_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Category_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Category_id_seq" OWNER TO postgres;

--
-- Name: Category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Category_id_seq" OWNED BY public."Category".id;


--
-- Name: Chat; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Chat" (
    id integer NOT NULL,
    "productId" integer NOT NULL,
    "buyerId" integer NOT NULL,
    "sellerId" integer NOT NULL,
    "unreadCountBuyer" integer DEFAULT 0 NOT NULL,
    "unreadCountSeller" integer DEFAULT 0 NOT NULL,
    "lastMessageId" integer,
    "lastMessageAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."Chat" OWNER TO postgres;

--
-- Name: Chat_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Chat_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Chat_id_seq" OWNER TO postgres;

--
-- Name: Chat_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Chat_id_seq" OWNED BY public."Chat".id;


--
-- Name: FavoriteAction; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."FavoriteAction" (
    id integer NOT NULL,
    "userId" integer NOT NULL,
    "productId" integer NOT NULL,
    "addedAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."FavoriteAction" OWNER TO postgres;

--
-- Name: FavoriteAction_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."FavoriteAction_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."FavoriteAction_id_seq" OWNER TO postgres;

--
-- Name: FavoriteAction_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."FavoriteAction_id_seq" OWNED BY public."FavoriteAction".id;


--
-- Name: Message; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Message" (
    id integer NOT NULL,
    content text NOT NULL,
    "senderId" integer NOT NULL,
    "chatId" integer NOT NULL,
    "isRead" boolean DEFAULT false NOT NULL,
    "readAt" timestamp(3) without time zone,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."Message" OWNER TO postgres;

--
-- Name: Message_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Message_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Message_id_seq" OWNER TO postgres;

--
-- Name: Message_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Message_id_seq" OWNED BY public."Message".id;


--
-- Name: PhoneNumberView; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."PhoneNumberView" (
    id integer NOT NULL,
    "viewedById" integer NOT NULL,
    "viewedUserId" integer NOT NULL,
    "viewedAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."PhoneNumberView" OWNER TO postgres;

--
-- Name: PhoneNumberView_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."PhoneNumberView_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."PhoneNumberView_id_seq" OWNER TO postgres;

--
-- Name: PhoneNumberView_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."PhoneNumberView_id_seq" OWNED BY public."PhoneNumberView".id;


--
-- Name: Product; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Product" (
    id integer NOT NULL,
    name text NOT NULL,
    price integer NOT NULL,
    state public."ProductState" NOT NULL,
    brand text,
    model text,
    description text,
    address text,
    images text[] DEFAULT ARRAY[]::text[],
    "categoryId" integer NOT NULL,
    "subCategoryId" integer NOT NULL,
    "userId" integer NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."Product" OWNER TO postgres;

--
-- Name: ProductView; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."ProductView" (
    id integer NOT NULL,
    "viewedById" integer NOT NULL,
    "productId" integer NOT NULL,
    "viewedAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."ProductView" OWNER TO postgres;

--
-- Name: ProductView_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."ProductView_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."ProductView_id_seq" OWNER TO postgres;

--
-- Name: ProductView_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."ProductView_id_seq" OWNED BY public."ProductView".id;


--
-- Name: Product_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Product_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Product_id_seq" OWNER TO postgres;

--
-- Name: Product_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Product_id_seq" OWNED BY public."Product".id;


--
-- Name: Review; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Review" (
    id integer NOT NULL,
    "reviewedById" integer NOT NULL,
    text text,
    rating double precision NOT NULL,
    "reviewedUserId" integer NOT NULL,
    "reviewedAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."Review" OWNER TO postgres;

--
-- Name: Review_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Review_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Review_id_seq" OWNER TO postgres;

--
-- Name: Review_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Review_id_seq" OWNED BY public."Review".id;


--
-- Name: SubCategory; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."SubCategory" (
    id integer NOT NULL,
    name text NOT NULL,
    "categoryId" integer NOT NULL
);


ALTER TABLE public."SubCategory" OWNER TO postgres;

--
-- Name: SubCategory_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."SubCategory_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."SubCategory_id_seq" OWNER TO postgres;

--
-- Name: SubCategory_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."SubCategory_id_seq" OWNED BY public."SubCategory".id;


--
-- Name: SubcategotyType; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."SubcategotyType" (
    id integer NOT NULL,
    name text NOT NULL,
    "subcategoryId" integer NOT NULL
);


ALTER TABLE public."SubcategotyType" OWNER TO postgres;

--
-- Name: SubcategotyType_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."SubcategotyType_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."SubcategotyType_id_seq" OWNER TO postgres;

--
-- Name: SubcategotyType_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."SubcategotyType_id_seq" OWNED BY public."SubcategotyType".id;


--
-- Name: User; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."User" (
    id integer NOT NULL,
    "fullName" text NOT NULL,
    email text NOT NULL,
    "phoneNumber" text NOT NULL,
    password text NOT NULL,
    "profileType" public."ProfileType" DEFAULT 'INDIVIDUAL'::public."ProfileType" NOT NULL,
    "refreshToken" text,
    "refreshTokenExpiresAt" timestamp(3) without time zone,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL,
    rating integer,
    "isResetVerified" boolean DEFAULT false NOT NULL
);


ALTER TABLE public."User" OWNER TO postgres;

--
-- Name: User_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."User_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."User_id_seq" OWNER TO postgres;

--
-- Name: User_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."User_id_seq" OWNED BY public."User".id;


--
-- Name: _UserFavorites; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."_UserFavorites" (
    "A" integer NOT NULL,
    "B" integer NOT NULL
);


ALTER TABLE public."_UserFavorites" OWNER TO postgres;

--
-- Name: _prisma_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public._prisma_migrations (
    id character varying(36) NOT NULL,
    checksum character varying(64) NOT NULL,
    finished_at timestamp with time zone,
    migration_name character varying(255) NOT NULL,
    logs text,
    rolled_back_at timestamp with time zone,
    started_at timestamp with time zone DEFAULT now() NOT NULL,
    applied_steps_count integer DEFAULT 0 NOT NULL
);


ALTER TABLE public._prisma_migrations OWNER TO postgres;

--
-- Name: Category id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Category" ALTER COLUMN id SET DEFAULT nextval('public."Category_id_seq"'::regclass);


--
-- Name: Chat id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Chat" ALTER COLUMN id SET DEFAULT nextval('public."Chat_id_seq"'::regclass);


--
-- Name: FavoriteAction id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."FavoriteAction" ALTER COLUMN id SET DEFAULT nextval('public."FavoriteAction_id_seq"'::regclass);


--
-- Name: Message id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message" ALTER COLUMN id SET DEFAULT nextval('public."Message_id_seq"'::regclass);


--
-- Name: PhoneNumberView id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."PhoneNumberView" ALTER COLUMN id SET DEFAULT nextval('public."PhoneNumberView_id_seq"'::regclass);


--
-- Name: Product id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Product" ALTER COLUMN id SET DEFAULT nextval('public."Product_id_seq"'::regclass);


--
-- Name: ProductView id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductView" ALTER COLUMN id SET DEFAULT nextval('public."ProductView_id_seq"'::regclass);


--
-- Name: Review id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Review" ALTER COLUMN id SET DEFAULT nextval('public."Review_id_seq"'::regclass);


--
-- Name: SubCategory id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SubCategory" ALTER COLUMN id SET DEFAULT nextval('public."SubCategory_id_seq"'::regclass);


--
-- Name: SubcategotyType id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SubcategotyType" ALTER COLUMN id SET DEFAULT nextval('public."SubcategotyType_id_seq"'::regclass);


--
-- Name: User id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User" ALTER COLUMN id SET DEFAULT nextval('public."User_id_seq"'::regclass);


--
-- Data for Name: Category; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Category" (id, name) FROM stdin;
1	Личные вещи
\.


--
-- Data for Name: Chat; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt") FROM stdin;
\.


--
-- Data for Name: FavoriteAction; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."FavoriteAction" (id, "userId", "productId", "addedAt") FROM stdin;
\.


--
-- Data for Name: Message; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Message" (id, content, "senderId", "chatId", "isRead", "readAt", "createdAt", "updatedAt") FROM stdin;
\.


--
-- Data for Name: PhoneNumberView; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."PhoneNumberView" (id, "viewedById", "viewedUserId", "viewedAt") FROM stdin;
1	2	1	2025-11-04 12:17:15.834
\.


--
-- Data for Name: Product; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Product" (id, name, price, state, brand, model, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt") FROM stdin;
1	Ручка шариковая	40	NEW	RYCHKA	Rychka 15 pro max	Новый Rychka 15 Pro max в отличном состоянии	г. Москва, ул. Тверская, д. 1	{/uploads/product/images-1762275324649-128468681.jpg}	1	6	1	2025-11-04 16:55:24.669	2025-11-04 16:55:24.669
\.


--
-- Data for Name: ProductView; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."ProductView" (id, "viewedById", "productId", "viewedAt") FROM stdin;
\.


--
-- Data for Name: Review; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Review" (id, "reviewedById", text, rating, "reviewedUserId", "reviewedAt") FROM stdin;
1	1	Коляска очень вместительная! 5 звёзд!	5	2	2025-11-04 21:12:15.924
2	3	Коляска пришла в плохом состоянии!	3.5	2	2025-11-04 21:15:55.997
\.


--
-- Data for Name: SubCategory; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."SubCategory" (id, name, "categoryId") FROM stdin;
1	Одежда	1
2	Детские товары	1
3	Красота и Здоровье (уход и гигиена, средства, приборы и аксессуары, парфюмерия)	1
4	Медицинские товары	1
5	Средства реабилитации	1
6	Школьные товары	1
7	Украшения	1
8	Продукты питания	1
9	Животные, растения	1
10	Бытовая техника	1
11	Посуда	1
12	Мебель	1
\.


--
-- Data for Name: SubcategotyType; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."SubcategotyType" (id, name, "subcategoryId") FROM stdin;
1	Мужская	1
2	Женская	1
3	Детская девочки	1
4	Детская мальчики	1
5	Ткани, текстиль и фурнитура	1
6	Сумки, рюкзаки	1
7	Аксессуары	1
8	Обувь	1
9	Игрушки	2
10	Детская мебель	2
11	Коляски детские	2
12	Велосипеды и самокаты	2
13	Детское питание и посуда	2
14	Образовательные товары	2
15	Уход и гигиена	2
16	Косметика для ухода за кожей	3
17	Средства для ухода за волосами	3
18	Уход и гигиена	3
19	Приборы и аксессуары	3
20	Парфюмерия	3
21	Макияж	3
22	Бады	3
23	Диагностическое оборудование	4
24	Оборудование для клиник	4
25	Медицинская мебель	4
26	Измерительные приборы	5
27	Ортопедия (бандажи, корсеты)	5
28	Уходовая косметика	5
29	Кресла-коляски	5
30	Спецодежда, трикотаж, компрессионное белье	5
31	Подгузники, пеленки, прокладки	5
32	Катетеры	5
33	Средства ухода за стомой	5
34	Кресла-стулья санитарные	5
35	Специальные устройства	5
36	Калоприемники, уроприемники	5
37	Трости, костыли	5
38	Вертикализаторы, опоры	5
39	Матрасы	5
40	Кровати медицинские	5
41	Письменные принадлежности	6
42	Бумажная продукция	6
43	Принадлежности для рисования и творчества	6
44	Органайзеры и хранение	6
45	Учебные пособия и инструменты	6
46	Рюкзаки и сумки	6
47	Прочее	6
48	Ювелирные изделия	7
49	Бижутерия	7
50	Часы	7
51	Готовые продукты	8
52	Напитки	8
53	Заморозки, полуфабрикаты	8
54	Домашние животные	9
55	С/х животные	9
56	Рептилии	9
57	Растения комнатные	9
58	Культурные растения	9
59	Декоративные уличные растения	9
60	Доп товары (горшки, грунт, кормилки, поилки, средства по уходу за растениями, инструменты, корма, игрушки, клетки, аксессуары)	9
61	Кухонная	10
62	Бытовая	10
63	Для приготовления пищи	11
64	Для хранения	11
65	Для сервировки	11
66	Для приёма пищи	11
67	Мягкая мебель	12
68	Корпусная мебель	12
69	Мебель для кухни	12
70	Мебель для спальни	12
71	Садовая мебель	12
72	Офисная мебель	12
\.


--
-- Data for Name: User; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "refreshToken", "refreshTokenExpiresAt", "createdAt", "updatedAt", rating, "isResetVerified") FROM stdin;
1	Садиков Виталий Дмитриевич	vitaly.sadikov1@yandex.ru	+79510341677	$2b$10$X2estm9N44tB2hHPYAqbnecrXM/NF9iiSlJBIgL2opsHeO8t8J9n6	INDIVIDUAL	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEsImlhdCI6MTc2MjI5MDY0MCwiZXhwIjoxNzYyODk1NDQwfQ.0vwdYFk-GzxmCR2B5TjROMu8fqdRXaG_K0fvJ1nmcuY	2025-11-11 21:10:40.674	2025-11-04 12:16:36.672	2025-11-04 21:10:40.679	\N	f
3	Садиков Виталий Дмитриевич	vitaly.sadikov3@yandex.ru	+79510341675	$2b$10$M2j5JoP.BteJNzPie9k/X.NyAEJHzGbOzwekVf2jXCSRAQmdddQ2q	INDIVIDUAL	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjMsImlhdCI6MTc2MjI5MDkzNCwiZXhwIjoxNzYyODk1NzM0fQ.y2ZlwL4Vfr3z6tDOeVeHnaXqB8YDIoWTmZgaWT0YPAA	2025-11-11 21:15:34.856	2025-11-04 21:15:12.141	2025-11-04 21:15:34.861	\N	f
2	Садиков Виталий Дмитриевич	vitaly.sadikov2@yandex.ru	+79510341676	$2b$10$5C4dgckQJEBbOtu7ElWeHuMA438DcqEzfJKsuFWwITAsLLOd3HsZC	INDIVIDUAL	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjIsImlhdCI6MTc2MjI5MDk3NCwiZXhwIjoxNzYyODk1Nzc0fQ.6c_Hzo16e2E3idzrSCfN8VLQBZeP2Mt1X3pGOUntSJg	2025-11-11 21:16:14.133	2025-11-04 12:16:46.707	2025-11-04 21:16:14.137	\N	f
\.


--
-- Data for Name: _UserFavorites; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."_UserFavorites" ("A", "B") FROM stdin;
\.


--
-- Data for Name: _prisma_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public._prisma_migrations (id, checksum, finished_at, migration_name, logs, rolled_back_at, started_at, applied_steps_count) FROM stdin;
677d62ff-b994-413b-95cc-1e84acebc01a	9a97e81893bf2ec4fc3cbb193511b77b89e82561206c399240ad61a1ef3f5411	2025-11-04 14:10:04.704194+02	20251008175218_add_refresh_token_to_user	\N	\N	2025-11-04 14:10:04.69034+02	1
19a6379e-4e35-46d4-9f7c-e26039da52f6	cb65e46beba4ea36f2090bc28cb45f2d037cd9ceda9b5f126c4cae06f4c46b68	2025-11-04 14:10:04.730779+02	20251012105519_add_product_categories_images	\N	\N	2025-11-04 14:10:04.70509+02	1
2fa6215e-913c-49c2-b810-2a3e54e4c771	f1c9c744537ed418b626a499157343dd8afb663cdb8c6b3252b271c2c8c9f603	2025-11-04 14:10:04.742188+02	20251014052438_add_user_favorites	\N	\N	2025-11-04 14:10:04.731824+02	1
414c4d82-d2d0-4784-bf25-ac6de655e60c	03d3193b4270d1c37aa9f566c2eef2d9ee62ab10ef42ef75f3af96f130928504	2025-11-04 14:10:14.808667+02	20251104121014_add_phone_number_view_stats	\N	\N	2025-11-04 14:10:14.749453+02	1
\.


--
-- Name: Category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Category_id_seq"', 1, true);


--
-- Name: Chat_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Chat_id_seq"', 1, false);


--
-- Name: FavoriteAction_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."FavoriteAction_id_seq"', 1, false);


--
-- Name: Message_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Message_id_seq"', 1, false);


--
-- Name: PhoneNumberView_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."PhoneNumberView_id_seq"', 1, true);


--
-- Name: ProductView_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."ProductView_id_seq"', 1, false);


--
-- Name: Product_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Product_id_seq"', 1, true);


--
-- Name: Review_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Review_id_seq"', 2, true);


--
-- Name: SubCategory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."SubCategory_id_seq"', 12, true);


--
-- Name: SubcategotyType_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."SubcategotyType_id_seq"', 72, true);


--
-- Name: User_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."User_id_seq"', 3, true);


--
-- Name: Category Category_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Category"
    ADD CONSTRAINT "Category_pkey" PRIMARY KEY (id);


--
-- Name: Chat Chat_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Chat"
    ADD CONSTRAINT "Chat_pkey" PRIMARY KEY (id);


--
-- Name: FavoriteAction FavoriteAction_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."FavoriteAction"
    ADD CONSTRAINT "FavoriteAction_pkey" PRIMARY KEY (id);


--
-- Name: Message Message_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_pkey" PRIMARY KEY (id);


--
-- Name: PhoneNumberView PhoneNumberView_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."PhoneNumberView"
    ADD CONSTRAINT "PhoneNumberView_pkey" PRIMARY KEY (id);


--
-- Name: ProductView ProductView_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductView"
    ADD CONSTRAINT "ProductView_pkey" PRIMARY KEY (id);


--
-- Name: Product Product_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Product"
    ADD CONSTRAINT "Product_pkey" PRIMARY KEY (id);


--
-- Name: Review Review_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Review"
    ADD CONSTRAINT "Review_pkey" PRIMARY KEY (id);


--
-- Name: SubCategory SubCategory_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SubCategory"
    ADD CONSTRAINT "SubCategory_pkey" PRIMARY KEY (id);


--
-- Name: SubcategotyType SubcategotyType_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SubcategotyType"
    ADD CONSTRAINT "SubcategotyType_pkey" PRIMARY KEY (id);


--
-- Name: User User_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_pkey" PRIMARY KEY (id);


--
-- Name: _UserFavorites _UserFavorites_AB_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."_UserFavorites"
    ADD CONSTRAINT "_UserFavorites_AB_pkey" PRIMARY KEY ("A", "B");


--
-- Name: _prisma_migrations _prisma_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public._prisma_migrations
    ADD CONSTRAINT _prisma_migrations_pkey PRIMARY KEY (id);


--
-- Name: Chat_buyerId_sellerId_productId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Chat_buyerId_sellerId_productId_key" ON public."Chat" USING btree ("buyerId", "sellerId", "productId");


--
-- Name: FavoriteAction_userId_productId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "FavoriteAction_userId_productId_key" ON public."FavoriteAction" USING btree ("userId", "productId");


--
-- Name: PhoneNumberView_viewedById_viewedUserId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "PhoneNumberView_viewedById_viewedUserId_key" ON public."PhoneNumberView" USING btree ("viewedById", "viewedUserId");


--
-- Name: ProductView_viewedById_productId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "ProductView_viewedById_productId_key" ON public."ProductView" USING btree ("viewedById", "productId");


--
-- Name: Review_reviewedById_reviewedUserId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Review_reviewedById_reviewedUserId_key" ON public."Review" USING btree ("reviewedById", "reviewedUserId");


--
-- Name: User_email_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "User_email_key" ON public."User" USING btree (email);


--
-- Name: User_phoneNumber_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "User_phoneNumber_key" ON public."User" USING btree ("phoneNumber");


--
-- Name: _UserFavorites_B_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "_UserFavorites_B_index" ON public."_UserFavorites" USING btree ("B");


--
-- Name: Chat Chat_buyerId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Chat"
    ADD CONSTRAINT "Chat_buyerId_fkey" FOREIGN KEY ("buyerId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Chat Chat_lastMessageId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Chat"
    ADD CONSTRAINT "Chat_lastMessageId_fkey" FOREIGN KEY ("lastMessageId") REFERENCES public."Message"(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: Chat Chat_productId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Chat"
    ADD CONSTRAINT "Chat_productId_fkey" FOREIGN KEY ("productId") REFERENCES public."Product"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Chat Chat_sellerId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Chat"
    ADD CONSTRAINT "Chat_sellerId_fkey" FOREIGN KEY ("sellerId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: FavoriteAction FavoriteAction_productId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."FavoriteAction"
    ADD CONSTRAINT "FavoriteAction_productId_fkey" FOREIGN KEY ("productId") REFERENCES public."Product"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: FavoriteAction FavoriteAction_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."FavoriteAction"
    ADD CONSTRAINT "FavoriteAction_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Message Message_chatId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_chatId_fkey" FOREIGN KEY ("chatId") REFERENCES public."Chat"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Message Message_senderId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_senderId_fkey" FOREIGN KEY ("senderId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: PhoneNumberView PhoneNumberView_viewedById_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."PhoneNumberView"
    ADD CONSTRAINT "PhoneNumberView_viewedById_fkey" FOREIGN KEY ("viewedById") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: PhoneNumberView PhoneNumberView_viewedUserId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."PhoneNumberView"
    ADD CONSTRAINT "PhoneNumberView_viewedUserId_fkey" FOREIGN KEY ("viewedUserId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: ProductView ProductView_productId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductView"
    ADD CONSTRAINT "ProductView_productId_fkey" FOREIGN KEY ("productId") REFERENCES public."Product"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: ProductView ProductView_viewedById_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductView"
    ADD CONSTRAINT "ProductView_viewedById_fkey" FOREIGN KEY ("viewedById") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Product Product_categoryId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Product"
    ADD CONSTRAINT "Product_categoryId_fkey" FOREIGN KEY ("categoryId") REFERENCES public."Category"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Product Product_subCategoryId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Product"
    ADD CONSTRAINT "Product_subCategoryId_fkey" FOREIGN KEY ("subCategoryId") REFERENCES public."SubCategory"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Product Product_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Product"
    ADD CONSTRAINT "Product_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Review Review_reviewedById_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Review"
    ADD CONSTRAINT "Review_reviewedById_fkey" FOREIGN KEY ("reviewedById") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Review Review_reviewedUserId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Review"
    ADD CONSTRAINT "Review_reviewedUserId_fkey" FOREIGN KEY ("reviewedUserId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: SubCategory SubCategory_categoryId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SubCategory"
    ADD CONSTRAINT "SubCategory_categoryId_fkey" FOREIGN KEY ("categoryId") REFERENCES public."Category"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: SubcategotyType SubcategotyType_subcategoryId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SubcategotyType"
    ADD CONSTRAINT "SubcategotyType_subcategoryId_fkey" FOREIGN KEY ("subcategoryId") REFERENCES public."SubCategory"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: _UserFavorites _UserFavorites_A_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."_UserFavorites"
    ADD CONSTRAINT "_UserFavorites_A_fkey" FOREIGN KEY ("A") REFERENCES public."Product"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: _UserFavorites _UserFavorites_B_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."_UserFavorites"
    ADD CONSTRAINT "_UserFavorites_B_fkey" FOREIGN KEY ("B") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;


--
-- PostgreSQL database dump complete
--


