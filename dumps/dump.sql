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

--
-- Name: TicketPriority; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."TicketPriority" AS ENUM (
    'LOW',
    'MEDIUM',
    'HIGH',
    'URGENT'
);


ALTER TYPE public."TicketPriority" OWNER TO postgres;

--
-- Name: TicketStatus; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."TicketStatus" AS ENUM (
    'OPEN',
    'IN_PROGRESS',
    'RESOLVED',
    'CLOSED'
);


ALTER TYPE public."TicketStatus" OWNER TO postgres;

--
-- Name: TicketTheme; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."TicketTheme" AS ENUM (
    'TECHNICAL_ISSUE',
    'ACCOUNT_PROBLEM',
    'PAYMENT_ISSUE',
    'PRODUCT_QUESTION',
    'COMPLAINT',
    'SUGGESTION',
    'OTHER'
);


ALTER TYPE public."TicketTheme" OWNER TO postgres;

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
    description text,
    address text,
    images text[] DEFAULT ARRAY[]::text[],
    "categoryId" integer NOT NULL,
    "subCategoryId" integer NOT NULL,
    "userId" integer NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL,
    "typeId" integer,
    "videoUrl" text
);


ALTER TABLE public."Product" OWNER TO postgres;

--
-- Name: ProductFieldValue; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."ProductFieldValue" (
    id integer NOT NULL,
    value text NOT NULL,
    "fieldId" integer NOT NULL,
    "productId" integer NOT NULL
);


ALTER TABLE public."ProductFieldValue" OWNER TO postgres;

--
-- Name: ProductFieldValue_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."ProductFieldValue_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."ProductFieldValue_id_seq" OWNER TO postgres;

--
-- Name: ProductFieldValue_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."ProductFieldValue_id_seq" OWNED BY public."ProductFieldValue".id;


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
-- Name: Role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Role" (
    id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public."Role" OWNER TO postgres;

--
-- Name: Role_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Role_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Role_id_seq" OWNER TO postgres;

--
-- Name: Role_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Role_id_seq" OWNED BY public."Role".id;


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
-- Name: SupportMessage; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."SupportMessage" (
    id integer NOT NULL,
    "ticketId" integer NOT NULL,
    "authorId" integer NOT NULL,
    text text NOT NULL,
    "sentAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."SupportMessage" OWNER TO postgres;

--
-- Name: SupportMessage_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."SupportMessage_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."SupportMessage_id_seq" OWNER TO postgres;

--
-- Name: SupportMessage_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."SupportMessage_id_seq" OWNED BY public."SupportMessage".id;


--
-- Name: SupportTicket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."SupportTicket" (
    id integer NOT NULL,
    theme public."TicketTheme" NOT NULL,
    subject text NOT NULL,
    status public."TicketStatus" DEFAULT 'OPEN'::public."TicketStatus" NOT NULL,
    priority public."TicketPriority" DEFAULT 'MEDIUM'::public."TicketPriority" NOT NULL,
    "userId" integer NOT NULL,
    "moderatorId" integer,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."SupportTicket" OWNER TO postgres;

--
-- Name: SupportTicket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."SupportTicket_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."SupportTicket_id_seq" OWNER TO postgres;

--
-- Name: SupportTicket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."SupportTicket_id_seq" OWNED BY public."SupportTicket".id;


--
-- Name: TypeField; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."TypeField" (
    id integer NOT NULL,
    name text NOT NULL,
    "isRequired" boolean DEFAULT false NOT NULL,
    "typeId" integer NOT NULL
);


ALTER TABLE public."TypeField" OWNER TO postgres;

--
-- Name: TypeField_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."TypeField_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."TypeField_id_seq" OWNER TO postgres;

--
-- Name: TypeField_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."TypeField_id_seq" OWNED BY public."TypeField".id;


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
    "isResetVerified" boolean DEFAULT false NOT NULL,
    "roleId" integer,
    "isAnswersCall" boolean DEFAULT false,
    photo text,
    "isEmailVerified" boolean DEFAULT false NOT NULL
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
-- Name: ProductFieldValue id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductFieldValue" ALTER COLUMN id SET DEFAULT nextval('public."ProductFieldValue_id_seq"'::regclass);


--
-- Name: ProductView id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductView" ALTER COLUMN id SET DEFAULT nextval('public."ProductView_id_seq"'::regclass);


--
-- Name: Review id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Review" ALTER COLUMN id SET DEFAULT nextval('public."Review_id_seq"'::regclass);


--
-- Name: Role id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Role" ALTER COLUMN id SET DEFAULT nextval('public."Role_id_seq"'::regclass);


--
-- Name: SubCategory id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SubCategory" ALTER COLUMN id SET DEFAULT nextval('public."SubCategory_id_seq"'::regclass);


--
-- Name: SubcategotyType id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SubcategotyType" ALTER COLUMN id SET DEFAULT nextval('public."SubcategotyType_id_seq"'::regclass);


--
-- Name: SupportMessage id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SupportMessage" ALTER COLUMN id SET DEFAULT nextval('public."SupportMessage_id_seq"'::regclass);


--
-- Name: SupportTicket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SupportTicket" ALTER COLUMN id SET DEFAULT nextval('public."SupportTicket_id_seq"'::regclass);


--
-- Name: TypeField id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."TypeField" ALTER COLUMN id SET DEFAULT nextval('public."TypeField_id_seq"'::regclass);


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
2	6	2	2025-11-14 11:42:45.616
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
\.


--
-- Data for Name: Product; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl") FROM stdin;
2	iPhone 15 Pro	120000	NEW	Новый iPhone 15 Pro в отличном состоянии	г. Москва, ул. Тверская, д. 1	{/uploads/product/images-1762934382600-525832298.jpg}	1	1	5	2025-11-12 07:59:42.638	2025-11-12 07:59:42.638	\N	\N
4	Футболка Ronaldo	359999	NEW	Футболка Ronaldo с автографом месси	г. Москва, ул. Тверская, д. 1	{/uploads/product/images-1763556943231-999932093.jpg,/uploads/product/images-1763556943232-639717408.jpg}	1	1	5	2025-11-19 12:55:43.255	2025-11-19 12:55:43.255	1	\N
\.


--
-- Data for Name: ProductFieldValue; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."ProductFieldValue" (id, value, "fieldId", "productId") FROM stdin;
26	XXL	1	4
27	Чёрный	2	4
\.


--
-- Data for Name: ProductView; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."ProductView" (id, "viewedById", "productId", "viewedAt") FROM stdin;
4	6	2	2025-11-14 11:42:49.53
\.


--
-- Data for Name: Review; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Review" (id, "reviewedById", text, rating, "reviewedUserId", "reviewedAt") FROM stdin;
\.


--
-- Data for Name: Role; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Role" (id, name) FROM stdin;
1	default
2	moderator
3	admin
\.


--
-- Data for Name: SubCategory; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."SubCategory" (id, name, "categoryId") FROM stdin;
1	Одежда	1
2	Детские товары	1
3	Красота и Здоровье (уход и гигиена, средства, приборы и аксессуары, парфюмерия)	1
5	Средства реабилитации	1
6	Школьные товары	1
7	Украшения	1
8	Продукты питания	1
9	Животные, растения	1
10	Бытовая техника	1
11	Посуда	1
12	Мебель	1
15	Медицинские товары	1
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
73	Диагностическое оборудование	15
74	Оборудование для клиник	15
75	Медицинская мебель	15
\.


--
-- Data for Name: SupportMessage; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."SupportMessage" (id, "ticketId", "authorId", text, "sentAt") FROM stdin;
\.


--
-- Data for Name: SupportTicket; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."SupportTicket" (id, theme, subject, status, priority, "userId", "moderatorId", "createdAt", "updatedAt") FROM stdin;
\.


--
-- Data for Name: TypeField; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."TypeField" (id, name, "isRequired", "typeId") FROM stdin;
1	Размер	f	1
2	Цвет	f	1
3	Материал	f	1
4	Бренд	f	1
5	Название	f	1
6	Вид	f	1
7	Размер	f	2
8	Цвет	f	2
9	Материал	f	2
10	Бренд	f	2
11	Название	f	2
12	Вид	f	2
13	Размер	f	3
14	Цвет	f	3
15	Материал	f	3
16	Бренд	f	3
17	Название	f	3
18	Вид	f	3
19	Размер	f	4
20	Цвет	f	4
21	Материал	f	4
22	Бренд	f	4
23	Название	f	4
24	Вид	f	4
25	Размер	f	5
26	Цвет	f	5
27	Материал	f	5
28	Бренд	f	5
29	Название	f	5
30	Вид	f	5
31	Размер	f	6
32	Цвет	f	6
33	Материал	f	6
34	Бренд	f	6
35	Название	f	6
36	Вид	f	6
37	Размер	f	7
38	Цвет	f	7
39	Материал	f	7
40	Бренд	f	7
41	Название	f	7
42	Вид	f	7
43	Размер	f	8
44	Цвет	f	8
45	Материал	f	8
46	Бренд	f	8
47	Название	f	8
48	Вид	f	8
49	Цвет	f	15
50	Размер	f	13
51	Возраст	f	9
52	Габариты	f	12
53	Габариты	f	10
54	Возраст	f	14
55	Размер	f	15
56	Габариты	f	11
57	Цвет	f	13
58	Цвет	f	14
59	Размер	f	9
60	Возраст	f	13
61	Возраст	f	15
62	Размер	f	14
63	Цвет	f	9
64	Габариты	f	9
65	Возраст	f	10
66	Возраст	f	12
67	Габариты	f	14
68	Возраст	f	11
69	Цвет	f	11
70	Габариты	f	13
71	Размер	f	12
72	Размер	f	10
73	Габариты	f	15
74	Цвет	f	12
75	Цвет	f	10
76	Размер	f	11
77	Вид	f	20
78	Вид	f	17
79	Вид	f	18
80	Вид	f	21
81	Вид	f	19
82	Вид	f	16
83	Вид	f	22
84	Цвет	f	73
85	Наличие сертификата	f	74
86	Цвет	f	75
87	Портативность	f	73
88	Бренд	f	74
89	Портативность	f	75
90	Бренд	f	73
91	Бренд	f	75
92	Портативность	f	74
93	Наличие сертификата	f	75
94	Цвет	f	74
95	Наличие сертификата	f	73
96	Тип питания	f	26
97	Диапазон измерений	f	26
98	Бренд	f	26
99	Вид	f	26
100	Комплектация	f	26
101	Замеры аритмии	f	26
102	Индикаторы	f	26
103	Точность измерений	f	26
104	Производитель	f	26
105	Метод измерения	f	26
106	Память	f	26
107	Тип	f	26
108	Калибровка	f	26
109	Объем капли	f	26
110	Погрешность	f	26
111	Гибкость	f	26
112	Размер	f	26
113	Время измерения	f	26
114	Функции маркировки	f	26
115	Подсветка	f	26
116	Звуковой сигнал	f	26
117	Ребра жесткости	f	27
118	Вид	f	27
119	Конструктивные особенности	f	27
120	Область применения	f	27
121	Производитель	f	27
122	Степень фиксации	f	27
123	Гипоаллергенность	f	27
124	Назначение	f	27
125	Затяжки	f	27
126	Цвет	f	27
127	Размер	f	27
128	Шнурки	f	27
129	Возрастная группа	f	27
130	Материал	f	27
131	Пол	f	27
132	Тип	f	28
133	Вид	f	28
134	Срок годности	f	28
135	Производитель	f	28
136	Тип	f	29
137	Вид	f	29
138	Материал рамы	f	29
139	Вес	f	29
140	Грузоподъёмность	f	29
141	Колёса	f	29
142	Аккумулятор	f	29
143	Управление	f	29
144	Доп функции	f	29
145	Складная конструкция	f	29
146	Цвет	f	29
147	Материалы	f	30
148	Гипоаллергенность	f	30
149	Степень компрессии	f	30
150	Размер	f	30
151	Цвет	f	30
152	Защитные свойства	f	30
153	Доп функции	f	30
154	Производитель	f	30
155	Тип	f	31
156	Размер	f	31
157	Впитываемость	f	31
158	Материал впитывающего слоя	f	31
159	Материал внешнего слоя	f	31
160	Материал внутреннего слоя	f	31
161	Вид	f	31
162	Возраст	f	31
163	Доп свойства	f	31
164	Цвет	f	31
165	Производитель	f	31
166	Вид	f	32
167	Материал	f	32
168	Тип	f	32
169	Размер	f	32
170	Доп функции	f	32
171	Срок годности	f	32
172	Производитель	f	32
173	Тип	f	33
174	Вид стомы	f	33
175	Размер	f	33
176	Производитель	f	33
177	Тип	f	34
178	Материал рамы	f	34
179	Материал сиденья и спинки	f	34
180	Регулировка высоты сидений	f	34
181	Регулировка высоты и положения подлокотников	f	34
182	Размер	f	34
183	Доп опции	f	34
184	Максимальная нагрузка	f	34
185	Цвет	f	34
186	Производитель	f	34
187	Вид	f	35
188	Производитель	f	35
189	Характеристики устройства	f	35
190	Габариты	f	35
191	Вид	f	36
192	Тип	f	36
193	Материалы	f	36
194	Объём мешков	f	36
195	Диаметр пластин	f	36
196	Производитель	f	36
197	Наличие фильтров	f	36
198	Наличие клапанов	f	36
199	Наличие градуировки для измерения	f	36
200	Тип	f	37
201	Вид	f	37
202	По поддержке	f	37
203	Регулировка высоты	f	37
204	Материал опор	f	37
205	Вид наконечника	f	37
206	Допустимая нагрузка	f	37
207	Производитель	f	37
208	Тип рукоятки	f	37
209	Противоскользящий наконечник	f	37
\.


--
-- Data for Name: User; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "refreshToken", "refreshTokenExpiresAt", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified") FROM stdin;
5	Попов Матвей Иванович	vitaly.sadikov1@yandex.ru	+79510341677	$2b$10$05FMyE494pfJScN9OF98COs6yLacnIIE2gueMbTS8s1/PNzaYrA6C	INDIVIDUAL	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjUsImlhdCI6MTc2MzY0NDU0MCwiZXhwIjoxNzY0MjQ5MzQwfQ.olhsHykkbhIVVnvO8jC6YIZw7XfGZmHIGJTUCezkUkA	2025-11-27 13:15:40.03	2025-11-06 19:33:46.625	2025-11-20 13:31:39.524	\N	f	3	\N	/uploads/user/photo-1763627456668-926795266.png	t
6	Садиков Виталий Дмитриевич	vitaly.sadikov2@yandex.ru	+79510341676	$2b$10$Tsi0whXkdERT2AvjSe6Jn.v6ba.K3sTDPXT6AzWMlkpahIY.LxDSS	INDIVIDUAL	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjYsImlhdCI6MTc2MzYyMzUxNiwiZXhwIjoxNzY0MjI4MzE2fQ.Zaq_RCqkQqG3afjGnecJfuhc7HlVjcbHJzgjsh9FS8U	2025-11-27 07:25:16.021	2025-11-06 19:33:55.742	2025-11-20 07:25:16.031	\N	f	1	\N	\N	f
\.


--
-- Data for Name: _UserFavorites; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."_UserFavorites" ("A", "B") FROM stdin;
2	6
\.


--
-- Data for Name: _prisma_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public._prisma_migrations (id, checksum, finished_at, migration_name, logs, rolled_back_at, started_at, applied_steps_count) FROM stdin;
6c3c162d-a2a2-429c-a947-535a7d1ce19e	9a97e81893bf2ec4fc3cbb193511b77b89e82561206c399240ad61a1ef3f5411	2025-11-06 21:13:06.102356+02	20251008175218_add_refresh_token_to_user	\N	\N	2025-11-06 21:13:06.059517+02	1
e3b26ca6-1680-4576-bcfe-22ede4738205	cb65e46beba4ea36f2090bc28cb45f2d037cd9ceda9b5f126c4cae06f4c46b68	2025-11-06 21:13:06.165458+02	20251012105519_add_product_categories_images	\N	\N	2025-11-06 21:13:06.105344+02	1
f261c8c3-3a9f-4cd1-9edb-77561e67b37f	f1c9c744537ed418b626a499157343dd8afb663cdb8c6b3252b271c2c8c9f603	2025-11-06 21:13:06.180541+02	20251014052438_add_user_favorites	\N	\N	2025-11-06 21:13:06.167153+02	1
bfb3baed-143f-4caf-8d65-c07bdd63a80f	03d3193b4270d1c37aa9f566c2eef2d9ee62ab10ef42ef75f3af96f130928504	2025-11-06 21:13:06.237422+02	20251104121014_add_phone_number_view_stats	\N	\N	2025-11-06 21:13:06.181534+02	1
df00ed7f-80fe-4458-8f8b-83fe88304cc8	8d8cfe1eacb1a375fc8254a31aa50217e946af0d14f48f9a25000d7ccec5bfcc	2025-11-06 21:13:16.694481+02	20251106191316_add_support_system	\N	\N	2025-11-06 21:13:16.57959+02	1
677d62ff-b994-413b-95cc-1e84acebc01a	9a97e81893bf2ec4fc3cbb193511b77b89e82561206c399240ad61a1ef3f5411	2025-11-04 14:10:04.704194+02	20251008175218_add_refresh_token_to_user	\N	\N	2025-11-04 14:10:04.69034+02	1
19a6379e-4e35-46d4-9f7c-e26039da52f6	cb65e46beba4ea36f2090bc28cb45f2d037cd9ceda9b5f126c4cae06f4c46b68	2025-11-04 14:10:04.730779+02	20251012105519_add_product_categories_images	\N	\N	2025-11-04 14:10:04.70509+02	1
2fa6215e-913c-49c2-b810-2a3e54e4c771	f1c9c744537ed418b626a499157343dd8afb663cdb8c6b3252b271c2c8c9f603	2025-11-04 14:10:04.742188+02	20251014052438_add_user_favorites	\N	\N	2025-11-04 14:10:04.731824+02	1
414c4d82-d2d0-4784-bf25-ac6de655e60c	03d3193b4270d1c37aa9f566c2eef2d9ee62ab10ef42ef75f3af96f130928504	2025-11-04 14:10:14.808667+02	20251104121014_add_phone_number_view_stats	\N	\N	2025-11-04 14:10:14.749453+02	1
\.


--
-- Name: Category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Category_id_seq"', 2, true);


--
-- Name: Chat_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Chat_id_seq"', 1, false);


--
-- Name: FavoriteAction_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."FavoriteAction_id_seq"', 2, true);


--
-- Name: Message_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Message_id_seq"', 1, false);


--
-- Name: PhoneNumberView_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."PhoneNumberView_id_seq"', 1, true);


--
-- Name: ProductFieldValue_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."ProductFieldValue_id_seq"', 27, true);


--
-- Name: ProductView_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."ProductView_id_seq"', 4, true);


--
-- Name: Product_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Product_id_seq"', 4, true);


--
-- Name: Review_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Review_id_seq"', 2, true);


--
-- Name: Role_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Role_id_seq"', 3, true);


--
-- Name: SubCategory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."SubCategory_id_seq"', 15, true);


--
-- Name: SubcategotyType_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."SubcategotyType_id_seq"', 75, true);


--
-- Name: SupportMessage_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."SupportMessage_id_seq"', 1, false);


--
-- Name: SupportTicket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."SupportTicket_id_seq"', 1, false);


--
-- Name: TypeField_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."TypeField_id_seq"', 209, true);


--
-- Name: User_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."User_id_seq"', 6, true);


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
-- Name: ProductFieldValue ProductFieldValue_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductFieldValue"
    ADD CONSTRAINT "ProductFieldValue_pkey" PRIMARY KEY (id);


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
-- Name: Role Role_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Role"
    ADD CONSTRAINT "Role_pkey" PRIMARY KEY (id);


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
-- Name: SupportMessage SupportMessage_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SupportMessage"
    ADD CONSTRAINT "SupportMessage_pkey" PRIMARY KEY (id);


--
-- Name: SupportTicket SupportTicket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SupportTicket"
    ADD CONSTRAINT "SupportTicket_pkey" PRIMARY KEY (id);


--
-- Name: TypeField TypeField_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."TypeField"
    ADD CONSTRAINT "TypeField_pkey" PRIMARY KEY (id);


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
-- Name: ProductFieldValue_fieldId_productId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "ProductFieldValue_fieldId_productId_key" ON public."ProductFieldValue" USING btree ("fieldId", "productId");


--
-- Name: ProductView_viewedById_productId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "ProductView_viewedById_productId_key" ON public."ProductView" USING btree ("viewedById", "productId");


--
-- Name: Review_reviewedById_reviewedUserId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Review_reviewedById_reviewedUserId_key" ON public."Review" USING btree ("reviewedById", "reviewedUserId");


--
-- Name: Role_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Role_name_key" ON public."Role" USING btree (name);


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
-- Name: ProductFieldValue ProductFieldValue_fieldId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductFieldValue"
    ADD CONSTRAINT "ProductFieldValue_fieldId_fkey" FOREIGN KEY ("fieldId") REFERENCES public."TypeField"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: ProductFieldValue ProductFieldValue_productId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductFieldValue"
    ADD CONSTRAINT "ProductFieldValue_productId_fkey" FOREIGN KEY ("productId") REFERENCES public."Product"(id) ON UPDATE CASCADE ON DELETE CASCADE;


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
-- Name: Product Product_typeId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Product"
    ADD CONSTRAINT "Product_typeId_fkey" FOREIGN KEY ("typeId") REFERENCES public."SubcategotyType"(id) ON UPDATE CASCADE ON DELETE SET NULL;


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
-- Name: SupportMessage SupportMessage_authorId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SupportMessage"
    ADD CONSTRAINT "SupportMessage_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: SupportMessage SupportMessage_ticketId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SupportMessage"
    ADD CONSTRAINT "SupportMessage_ticketId_fkey" FOREIGN KEY ("ticketId") REFERENCES public."SupportTicket"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: SupportTicket SupportTicket_moderatorId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SupportTicket"
    ADD CONSTRAINT "SupportTicket_moderatorId_fkey" FOREIGN KEY ("moderatorId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: SupportTicket SupportTicket_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."SupportTicket"
    ADD CONSTRAINT "SupportTicket_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TypeField TypeField_typeId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."TypeField"
    ADD CONSTRAINT "TypeField_typeId_fkey" FOREIGN KEY ("typeId") REFERENCES public."SubcategotyType"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: User User_roleId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_roleId_fkey" FOREIGN KEY ("roleId") REFERENCES public."Role"(id) ON UPDATE CASCADE ON DELETE SET NULL;


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


