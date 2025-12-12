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
-- Name: ProductModerate; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."ProductModerate" AS ENUM (
    'MODERATE',
    'APPROVED',
    'DENIDED'
);


ALTER TYPE public."ProductModerate" OWNER TO postgres;

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
-- Name: OkseiProduct; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."OkseiProduct" (
    id integer NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    price integer NOT NULL,
    image text,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."OkseiProduct" OWNER TO postgres;

--
-- Name: OkseiProduct_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."OkseiProduct_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."OkseiProduct_id_seq" OWNER TO postgres;

--
-- Name: OkseiProduct_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."OkseiProduct_id_seq" OWNED BY public."OkseiProduct".id;


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
    address text NOT NULL,
    images text[] DEFAULT ARRAY[]::text[],
    "categoryId" integer NOT NULL,
    "subCategoryId" integer NOT NULL,
    "userId" integer NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL,
    "typeId" integer,
    "videoUrl" text,
    "isHide" boolean DEFAULT false NOT NULL,
    "moderateState" public."ProductModerate" DEFAULT 'MODERATE'::public."ProductModerate" NOT NULL,
    "moderationRejectionReason" text
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
-- Name: ProductPromotion; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."ProductPromotion" (
    id integer NOT NULL,
    "productId" integer NOT NULL,
    "promotionId" integer NOT NULL,
    "userId" integer NOT NULL,
    days integer NOT NULL,
    "totalPrice" integer NOT NULL,
    "startDate" timestamp(3) without time zone NOT NULL,
    "endDate" timestamp(3) without time zone NOT NULL,
    "isActive" boolean DEFAULT true NOT NULL,
    "isPaid" boolean DEFAULT false NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."ProductPromotion" OWNER TO postgres;

--
-- Name: ProductPromotion_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."ProductPromotion_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."ProductPromotion_id_seq" OWNER TO postgres;

--
-- Name: ProductPromotion_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."ProductPromotion_id_seq" OWNED BY public."ProductPromotion".id;


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
-- Name: Promotion; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Promotion" (
    id integer NOT NULL,
    name text NOT NULL,
    "pricePerDay" integer NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."Promotion" OWNER TO postgres;

--
-- Name: Promotion_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Promotion_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Promotion_id_seq" OWNER TO postgres;

--
-- Name: Promotion_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Promotion_id_seq" OWNED BY public."Promotion".id;


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
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL,
    rating integer,
    "isResetVerified" boolean DEFAULT false NOT NULL,
    "roleId" integer,
    "isAnswersCall" boolean DEFAULT false,
    photo text,
    "isEmailVerified" boolean DEFAULT false NOT NULL,
    balance double precision DEFAULT 0 NOT NULL
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
-- Name: OkseiProduct id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."OkseiProduct" ALTER COLUMN id SET DEFAULT nextval('public."OkseiProduct_id_seq"'::regclass);


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
-- Name: ProductPromotion id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductPromotion" ALTER COLUMN id SET DEFAULT nextval('public."ProductPromotion_id_seq"'::regclass);


--
-- Name: ProductView id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductView" ALTER COLUMN id SET DEFAULT nextval('public."ProductView_id_seq"'::regclass);


--
-- Name: Promotion id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Promotion" ALTER COLUMN id SET DEFAULT nextval('public."Promotion_id_seq"'::regclass);


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
1	13	7	16	0	0	\N	2025-11-28 09:15:51.647	2025-11-28 09:15:51.647	2025-11-28 09:15:51.647
2	16	11	9	0	0	\N	2025-11-28 09:16:07.289	2025-11-28 09:16:07.289	2025-11-28 09:16:07.289
3	16	13	9	0	0	\N	2025-11-28 09:16:15.369	2025-11-28 09:16:15.369	2025-11-28 09:16:15.369
4	17	14	13	0	0	\N	2025-11-28 09:16:50.74	2025-11-28 09:16:50.74	2025-11-28 09:16:50.74
5	18	7	15	0	0	\N	2025-11-28 09:17:36.46	2025-11-28 09:17:36.46	2025-11-28 09:17:36.46
6	22	11	9	0	0	\N	2025-11-28 09:18:57.647	2025-11-28 09:18:57.647	2025-11-28 09:18:57.647
7	20	19	16	0	0	\N	2025-11-28 09:20:00.074	2025-11-28 09:20:00.074	2025-11-28 09:20:00.074
9	23	5	18	0	0	\N	2025-11-28 09:43:43.309	2025-11-28 09:43:43.309	2025-11-28 09:43:43.309
42	11	53	10	0	0	\N	2025-12-01 08:12:52.805	2025-12-01 08:12:52.805	2025-12-01 08:12:52.805
75	28	86	20	0	0	\N	2025-12-01 08:28:55.846	2025-12-01 08:28:55.846	2025-12-01 08:28:55.846
108	61	5	53	0	0	\N	2025-12-01 08:36:19.933	2025-12-01 08:36:19.933	2025-12-01 08:36:19.933
141	20	5	16	0	0	\N	2025-12-01 12:30:18.663	2025-12-01 12:30:18.663	2025-12-01 12:30:18.663
142	28	119	20	0	0	\N	2025-12-01 14:29:48.006	2025-12-01 14:29:48.006	2025-12-01 14:29:48.006
143	23	119	18	0	0	\N	2025-12-01 14:46:54.623	2025-12-01 14:46:54.623	2025-12-01 14:46:54.623
145	252	121	120	0	1	2	2025-12-02 11:33:33.784	2025-12-02 11:33:25.608	2025-12-02 11:47:19.434
144	28	5	20	0	1	1	2025-12-02 06:39:40.026	2025-12-02 06:32:54.893	2025-12-04 06:34:08.592
\.


--
-- Data for Name: FavoriteAction; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."FavoriteAction" (id, "userId", "productId", "addedAt") FROM stdin;
3	14	20	2025-11-28 09:18:22.746
4	9	22	2025-11-28 09:18:35.259
5	14	6	2025-11-28 09:21:05.722
6	14	10	2025-11-28 09:21:07.97
7	7	6	2025-11-28 09:21:12.339
8	5	23	2025-11-29 09:00:09.809
41	86	28	2025-12-01 09:22:11.203
74	119	28	2025-12-01 14:29:45.996
75	5	127	2025-12-02 07:34:17.402
76	5	94	2025-12-02 07:34:18.185
77	5	21	2025-12-02 07:46:33.084
78	20	282	2025-12-03 00:00:00.214
\.


--
-- Data for Name: Message; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Message" (id, content, "senderId", "chatId", "isRead", "readAt", "createdAt", "updatedAt") FROM stdin;
1	тест	5	144	f	\N	2025-12-02 06:39:40.018	2025-12-02 06:39:40.018
2	Куда цену задрал? 200 край	121	145	f	\N	2025-12-02 11:33:33.781	2025-12-02 11:33:33.781
\.


--
-- Data for Name: OkseiProduct; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."OkseiProduct" (id, name, description, price, image, "createdAt") FROM stdin;
1	iPhone 15 Pro	Новый iPhone 15 Pro в отличном состоянии	120000	https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1c68c479-ade3-43ff-91eb-b8428b46ed74.jpg	2025-12-12 08:53:25.175
\.


--
-- Data for Name: PhoneNumberView; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."PhoneNumberView" (id, "viewedById", "viewedUserId", "viewedAt") FROM stdin;
\.


--
-- Data for Name: Product; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") FROM stdin;
7	Бусы б/у	1000	USED	Красные, из жемчуга	г Екатеринбург, ул Чкалова	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png}	1	3	16	2025-11-28 09:11:49.34	2025-12-03 19:00:32.343	\N	\N	f	APPROVED	\N
282	Кофемашина Thomson CF20A02	11399	NEW	Рабoчaя бытoвая тexника, намного дeшевлe, чем в мaгaзинe;\r\n\r\n\r\n\r\nНeт тapы для мoлoкa\r\n\r\n- Любыe пpoверки при caмoвывозе;\r\n\r\n- Пpи пpиемке товара вся теxника пpoвepяется на рaбoтоспocoбнoсть;\r\n\r\n- Oтправляeм Авитo доcтaвкой;\r\n\r\n- Пpи дoставке тoвap упaковываeтся по высшему урoвню.\r\n\r\nBитpинный образец:\r\n\r\n• товар новый, стоял на витрине в магазине;\r\n\r\n• может быть повреждена заводская упаковка;\r\n\r\n• возможны незначительные потёртости или повреждения корпуса, которые никак не влияют на работоспособность.\r\n\r\nЗа фотографиями дефектов обращайтесь в лс\r\n\r\nСамовывоз возможен из 2-х точек: Метро Текстильщики, Метро Шипиловская.\r\n\r\nВ нашем профиле большой ассортимент разнообразной бытовой техники. Советуем заглянуть!\r\n\r\nБольше техники в нашем телеграмм-канале\r\n\r\nПереходите там большие скидки!\r\n\r\n	44, улица Кирова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/673d92ac-2d1c-415f-8beb-02bd64a3b69d.png}	1	10	120	2025-12-02 11:34:25.201	2025-12-03 19:00:32.44	\N	\N	f	MODERATE	\N
23	ВАЗ 2107	435000	NEW	Продаётся готовый проект под RDS. Соответствует всем стандартам турниров и сходок. Гарантия на проект год.	Степной, Дзержинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/af9ca37d-87d9-44bc-b0aa-b1fc99737315.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/7a511cc5-d998-49c3-8f53-5dd88abd875b.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1d044f53-c12f-4841-9f4b-b486e551411a.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5c8cfe09-cd90-489a-8aad-0f8c2e80f6f4.jpg}	1	10	18	2025-11-28 09:18:17.344	2025-12-03 19:00:32.467	\N	\N	f	APPROVED	\N
94	Очередной товар дня!	35000	NEW	1) пусть будет текст\r\n2) здесь еще что-то\r\n**\r\n💥\r\n🟩\r\nККЕКЕКЕЕУУЦКУ""\r\n                                              ЦЕНТР\r\n          ТАБУЛЯЦИЯ СмещЕНИЕ\r\n\r\n	18, улица Расковой, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/6638b79f-2357-46ff-9010-ba9175ce50db.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c02fa4fd-6284-45a5-8cd7-61583db872fe.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fe8ad9f4-5664-4832-b5be-dc1f4df2adcf.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5488c8b5-af91-4294-85c6-0bb7d48145b6.jpg}	1	12	86	2025-12-01 08:35:56.623	2025-12-03 19:00:32.477	\N	\N	f	APPROVED	\N
233	Зипка	5000	NEW	Кофта теплая на замке	35, улица 9 Января, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f64e149a-e711-4173-85ea-98db13c3ca1e.png}	1	1	7	2025-12-02 10:59:24.476	2025-12-03 19:00:32.504	\N	\N	f	APPROVED	\N
235	Украшения ручной работы	1000	NEW	Украшения ручной работы на заказ по Вашим эскизам/фото. Стоимость украшений на фото 1000р. Срок изготовления: 4-7 дней.\r\n\r\n	46, улица 9 Января, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/4164bbf3-aa1d-4dac-b361-dd22fc5c2001.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/471637e1-9898-41d2-9dc1-c112c642c296.png}	1	7	120	2025-12-02 11:01:00.313	2025-12-03 19:00:32.509	\N	\N	f	APPROVED	\N
250	Капельница	200	NEW	Просто капельница	г Оренбург, ул Харьковская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b0a71fb2-6719-4085-b554-d16b5cf9b2a2.webp}	1	15	121	2025-12-02 11:11:52.698	2025-12-03 19:00:32.533	\N	\N	f	APPROVED	\N
251	Люлька	2000	USED	Люлька детская	68, улица Кичигина, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/51e36a8d-6904-4a08-96bc-d1d449241608.png}	1	2	7	2025-12-02 11:11:54.214	2025-12-03 19:00:32.536	\N	\N	f	APPROVED	\N
13	Протеин 1000гр	1500	NEW	Вкус шоколад, 1000 грамм	2, улица 13-я Линия, Линии, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460040, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a745c362-a754-4a92-abb9-b8969bebead7.png}	1	8	16	2025-11-28 09:14:04.157	2025-12-03 19:00:32.364	\N	\N	f	APPROVED	\N
21	Набор украшений для пирсинга	4000	NEW	\N	12А, Больничный проезд, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3c250f31-591b-4346-b1fb-1b3bf70f2c73.webp}	1	3	10	2025-11-28 09:17:54.801	2025-12-03 19:00:32.368	\N	\N	f	APPROVED	\N
27	Детские книжки по математике	1000	USED	Превосходный источник знаний для вашего ребенка	Hawthorne Street, Кламат-Фолс, Klamath County, Орегон, 97601, Соединённые Штаты Америки	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5e8875c2-aec8-4b2f-b618-2e220defa9cf.webp,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/2248f3b6-63b3-44e7-84a4-ef35b2d7bcdc.jpg}	1	2	15	2025-11-28 09:21:17.846	2025-12-03 19:00:32.37	\N	\N	f	APPROVED	\N
252	Минский Бургер с курицей	330	NEW	По-белорусски вкусный! Бургер с сочной куриной котлетой в хрустящей панировке, румяным картофельным оладушком, свежим салатом, двумя ломтиками нежного сыра, хрустящим ароматным беконом, маринованными огурчиками, нежным соусом «Сметана-укроп», и всё это — на воздушной горячей булочке с хрустящей крошкой.	54, улица Кирова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3458ebaf-7b3f-4f39-b1f4-5a53322d9e64.png}	1	8	120	2025-12-02 11:12:20.477	2025-12-03 19:00:32.541	\N	\N	f	APPROVED	\N
268	Алоэ вера лечебный 3 года, есть 1 год	200	USED	Алое Вера, лечебное 3х детки, есть однолетки	В. И. Ленину, Ленинская улица, Кузнечный, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/17ad0794-e12d-433a-9958-528bba02bf87.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/09b01358-7870-4abd-a66b-7cfcab7ecec9.png}	1	9	120	2025-12-02 11:26:21.919	2025-12-03 19:00:32.563	\N	\N	f	APPROVED	\N
271	Фундук культурный	280	USED	Прoдaю фундук 2024г cбopa, собственный небoльшой cад в предгopьяx Кавказa, бeз xимии тoлькo органикa.\r\n\r\nBcе сopта выращиваемые мной имеют лучшиe вкусoвыe xaрактериcтики и oтносятcя к cтoлoвым сoртам, oбладaют плотным ядpoм и приятным выpaженным маcляниcтым вкуcoм, который не сравним с дешёвыми cетевыми безвкусными орешками.\r\nПредлагаю микс сортов Трапезунд, Анаклиури, Президент.\r\n\r\nВозможна доставка авитодоставкой до 20кг или транспортной компанией от 30.	19/2, улица Бурзянцева, Кузнечный, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a0505026-7df1-43c2-a19c-16e30e07a690.png}	1	9	120	2025-12-02 11:27:45.043	2025-12-03 19:00:32.568	\N	\N	f	APPROVED	\N
160	Графин в виде рыбы	500	NEW	Замечательный графин в виде рыбы	г Оренбург, ул Киевская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/67533ddf-2495-4aed-b405-a68922a398bf.jpg}	1	11	121	2025-12-02 10:50:32.345	2025-12-03 19:00:32.373	\N	\N	f	APPROVED	\N
195	Сковорода антипригарная	1000	NEW	Сковорода. Можно пожарить все что угодно	г Оренбург, ул Днепропетровская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/299e425b-f6c0-49bb-a7e8-3c7c591ce39d.jpg}	1	11	121	2025-12-02 10:53:18.109	2025-12-03 19:00:32.377	\N	\N	f	APPROVED	\N
228	Стакан	200	NEW	Просто стакан.	г Оренбург, ул Житомирская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/4cfed187-e55c-4014-8c59-cd2450aca91e.jpg}	1	11	121	2025-12-02 10:55:45.209	2025-12-03 19:00:32.38	\N	\N	f	APPROVED	\N
348	папавпа	55454	NEW	павпвапа	«Урал», Ленинский район, Пригородный, Пригородный сельсовет, Оренбургский район, Оренбургская область, Приволжский федеральный округ, 460041, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/396aecf1-b980-45ec-bd5a-ba238f1fdefb.jpg}	1	5	5	2025-12-03 19:36:04.58	2025-12-03 19:36:04.58	\N	\N	f	MODERATE	\N
243	Ингалятор	2000	NEW	Ингалятор для ингаляций	г Оренбург, ул Луганская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/08d81171-c46d-4857-b7d9-bb7a983d5ab4.jpg}	1	15	121	2025-12-02 11:05:14.587	2025-12-03 19:00:32.386	\N	\N	f	APPROVED	\N
246	Массажный стол	4000	NEW	Просто массажный стол	г Оренбург, ул Запорожская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/002d563c-310a-4287-ac15-5826d88e5d37.jpg}	1	15	121	2025-12-02 11:07:04.422	2025-12-03 19:00:32.392	\N	\N	f	APPROVED	\N
244	Посуда детская	1500	NEW	Детская посуда для кормления	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f252c198-af82-42a9-a80b-e42a052caae3.png}	1	2	7	2025-12-02 11:06:21.305	2025-12-09 06:13:48.017	\N	\N	f	APPROVED	\N
245	Украшения	1000	USED	Продам укрошенияБраслет -500\r\nСерьги - 300\r\nКольцо 10 - 250\r\nВсе вместе 1000	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a28160dd-9d06-4750-b6fe-6045f6a3df8b.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e49061f1-c633-4254-bc95-b9e06ae322ae.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/15b89ff9-deb0-4067-80a9-77151e9ad946.png}	1	7	120	2025-12-02 11:06:25.377	2025-12-09 06:13:48.017	\N	\N	f	APPROVED	\N
248	Нутриен energy питание	2700	NEW	Cмecь Nutrien enеrgy, диетичeское лечeбноe  питание,\r\n\r\nПитаниe для oнкoбольных , питаниe для ocлaблeнных, питание пocлe опеpaции, питание, обогащённое витаминaми и микрoэлeмeнтaми.\r\n\r\nПродукт готовый к упoтpеблению 200 мл, 300 ккaл.\r\nПoдxoдит для онкoбольных, пoслeопepациoнныx взрослыx и дeтeй с 3 лет для вoсстановления сил и энергии.	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/6f42b7d6-7a7d-46ca-85c3-611b159a8a0a.png}	1	8	120	2025-12-02 11:08:58.23	2025-12-09 06:13:48.017	\N	\N	f	APPROVED	\N
259	Крем для рук	200	NEW	Просто крем для рук	г Оренбург, поселок Нижнесакмарский, ул Николаевская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fb99e19c-4e45-43b7-90fc-1c48c50106d7.webp}	1	3	121	2025-12-02 11:17:30.55	2025-12-03 19:00:32.394	\N	\N	f	APPROVED	\N
260	Пельмени домашние	380	NEW	Пpeдcтaвляeм вaшему вниманию пельмени, манты, хинкaли, ваpеники pучнoй лeпки.	18, Матросский переулок, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/19b43a4f-6ae7-4539-bee3-78a337e8e3c8.png}	1	8	120	2025-12-02 11:17:35.121	2025-12-03 19:00:32.4	\N	\N	f	APPROVED	\N
264	Матрас	3000	USED	Матрас для восстанволения	61А, улица Орлова, Новостройка, Промышленный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/7d6cda50-98e5-42bc-9a31-08b62188d9fa.png}	1	5	7	2025-12-02 11:22:24.486	2025-12-03 19:00:32.405	\N	\N	f	APPROVED	\N
261	Компрессорный ингалятор	2000	NEW	Компрессорный ингалятор	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f178392b-1965-4569-b91c-c6efd48b56da.png}	1	5	7	2025-12-02 11:18:32.771	2025-12-09 06:13:48.017	\N	\N	f	APPROVED	\N
262	Кресло-коляска	5000	USED	Кресло-коляска для инвалидов Ortonica Olvia 30	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/695dac97-c44a-42c0-9307-2cc8ff3bcaab.png}	1	5	7	2025-12-02 11:21:05.818	2025-12-09 06:13:48.017	\N	\N	f	APPROVED	\N
267	Спрей	600	USED	Защитная пленка для кожи	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/002307c3-8394-46f2-b18a-952a389efc6d.png}	1	5	7	2025-12-02 11:25:57	2025-12-09 06:13:48.017	\N	\N	f	APPROVED	\N
315	Пароочиститель для дома мощный, новые	1650	NEW	Унивеpсaльный паровой очиcтитель – этo эффективная бытовaя тeхникa для убоpки дoмa, coздaнная для удобствa и экoнoмии времeни. Этoт мощный парогенератoр cтaнeт вaшим надежным помощникoм в бoрьбе c зaгpязнeниями на куxне, мебeли и другиx повepхнocтях.\r\n	Мегаполис, 22, улица Володарского, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/965b0669-84ea-48c0-80e3-3e5ce0fe5022.png}	1	10	120	2025-12-02 11:35:21.921	2025-12-03 19:00:32.411	\N	\N	f	MODERATE	\N
15	SWEETPEEPS золотые украшения	7000	NEW	Золотые украшения с фианитами	Уральская улица, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f3a550d2-fc17-4548-943e-b62d66c014eb.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/afa0dc23-c3de-4ab1-8844-a4506bb49309.png}	1	7	7	2025-11-28 09:14:30.517	2025-12-03 19:00:32.414	\N	\N	f	APPROVED	\N
17	Chevrolet Corvette C7	8500000	USED	Корвет был угнан у курседа	Оренбургский Колледж Экономики и Информатики, 11, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/d0c85309-e830-459b-9f87-a672313a465e.jpg}	1	7	13	2025-11-28 09:15:31.784	2025-12-03 19:00:32.416	\N	\N	f	APPROVED	\N
19	ДМРВ на ваз 2107	7000	NEW	Датчик массового расхода воздуха	48, улица Коминтерна, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/05f3e4d2-118c-4262-a12a-8804cddb31d7.webp}	1	10	19	2025-11-28 09:16:40.023	2025-12-03 19:00:32.42	\N	https://yandex.ru/video/preview/13520813755431483017	f	APPROVED	\N
274	Эублефар	4000	USED	Продаются малыши эублефары различных морф. Едят разморозку, линяют хорошо, все процессы в норме.\r\n	2, Госпитальный переулок, Аренда, Промышленный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/99f1ee8b-a4e7-400c-9591-74ad68a831b6.png}	1	9	120	2025-12-02 11:29:46.562	2025-12-03 19:00:32.422	\N	\N	f	APPROVED	\N
5	Собака	100	USED	Собака овчарка	3/5, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b48c5093-8d65-47e0-84be-e1736ceffbe9.png}	1	9	11	2025-11-28 09:10:38.847	2025-12-03 19:00:32.426	\N	\N	f	APPROVED	\N
9	Игуана	285000	USED	Xoроший спoкoйный пaрень в самoм рaсцвeтe игуаниx cил.\r\n\r\nЗовут Яша, 19 лет, любит тепло и голубику.	26Б, улица Шевченко, Дзержинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/77280a5a-a493-49e1-aeb4-5bb7cbe97653.webp}	1	9	18	2025-11-28 09:11:59.098	2025-12-03 19:00:32.428	\N	\N	f	APPROVED	\N
18	Monster Energy Pipeline Punch	250	NEW	Тонизирующий напиток с изысканным вкусом!	Hawthorne Street, Кламат-Фолс, Klamath County, Орегон, 97601, Соединённые Штаты Америки	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/bfa67311-f947-41c7-b275-e7d28e1db313.jpg}	1	8	15	2025-11-28 09:16:21.207	2025-12-03 19:00:32.431	\N	https://vk.com/video-129440544_456249335	f	APPROVED	\N
20	Кресло-горилла	170000	NEW	Кресло-горилла удобное, выполнено из лучших материалов.	37А, Илекская улица, село имени 9 Января, Красноуральский сельсовет, Оренбургский район, Оренбургская область, Приволжский федеральный округ, 460501, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/228256c7-91a9-4c50-9a2e-a2804917075b.png}	1	12	16	2025-11-28 09:16:46.544	2025-12-03 19:00:32.434	\N	\N	f	APPROVED	\N
22	Майка	20000	USED	очень крутые маечки с аниме принтами, у2к вайб имеется🪽размер S, полиэстер\r\nцена 500 рублей за штуку\r\nпо любым вопросам пишите!!\r\n\r\n	3, улица Аксакова, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/33229854-0ccb-4ecc-b7a8-3ab2965d8fdc.png}	1	1	9	2025-11-28 09:18:16.763	2025-12-03 19:00:32.438	\N	\N	f	APPROVED	\N
316	Соковыжималка caso CP 300 Pro	4500	USED	CASO – нeмeцкaя торговая маркa бытовoй техники, принадлежащaя кoмпaнии Braukmann GmbH. Cоковыжималкa CASО СP 330 Prо предназначена для цитруcовыx cpeднего и крупногo pазмеpoв. Koрпуc прибоpа и cито для жмыxa выполнeны из нeржавeющeй cтaли. Автoмaтичeский старт плавнo зaпускает двигатель мощностью 160 Вт, функция «капля – стоп» обеспечивает чистоту рабочего места. В идеальном состоянии.\r\n\r\n	23/2, Пролетарская улица, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f2216887-329e-4d48-a605-e9b8bad18686.png}	1	10	120	2025-12-02 11:36:09.481	2025-12-03 19:00:32.444	\N	\N	f	MODERATE	\N
6	Посуда для сервировки Estetic	3500	NEW	Вся посуда выполнена в минималистичных стилях, из качественных материалов, подойдет на каждый день	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3b0ddf2a-19a7-4ecf-99c9-9dfc705d35c7.png}	1	11	7	2025-11-28 09:11:35.533	2025-12-03 19:00:32.447	\N	\N	f	APPROVED	\N
8	Кошка	10	USED	Кошка домашняя	5, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/8f6c52ee-023b-41d0-befc-39612d968abf.webp}	1	9	11	2025-11-28 09:11:58.971	2025-12-03 19:00:32.449	\N	\N	f	APPROVED	\N
10	Ford Mustang	2500000	NEW	Самый лучший автомобиль в мире	Оренбургский Колледж Экономики и Информатики, 11, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c91c9b99-afe9-4d10-8a81-b3cee1c12296.jpg}	1	10	13	2025-11-28 09:12:03.576	2025-12-03 19:00:32.451	\N	\N	f	APPROVED	\N
11	Набор золотых украшений	2000	NEW	\N	Лицей №2, Красная улица, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e7670256-5300-4433-82a5-edf31f999776.webp}	1	7	10	2025-11-28 09:12:06.874	2025-12-03 19:00:32.454	\N	\N	f	APPROVED	\N
12	Салонный фильтр на ваз 2110	1000	NEW	салонный фильтр подходит на автомобили ваз2110,2112	20, улица Кобозева, Кузнечный, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/22cb6191-445a-4ee0-80f2-92cc96055093.webp}	1	10	19	2025-11-28 09:12:25.319	2025-12-03 19:00:32.456	\N	https://yandex.ru/video/preview/9506785745966413491	f	APPROVED	\N
14	Козы камерунские	3000	NEW	Продаются козочки камерунские,разного возраста, есть два козлика для покрытия, покрытие 3 тыс	"Воздух" конный клуб, 9, Бассейный переулок, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/bfba4d6b-5f64-40b1-a0e2-8f232b9140ea.webp}	1	9	18	2025-11-28 09:14:23.82	2025-12-03 19:00:32.459	\N	\N	f	APPROVED	\N
16	Платье горничной	1200	USED	платье горничной в хорошем состоянии , нету только ободка осталась только от него ткань, если нужно доп фото пишите, к платью идет бантик и фартук\r\n\r\n	г Оренбург	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/dd3f08a3-bdcc-4593-b35f-1e186ce5262a.png}	1	1	9	2025-11-28 09:14:31.733	2025-12-03 19:00:32.462	\N	\N	f	APPROVED	\N
232	Этно украшения	300	USED	Украшения в этническом стиле! серьги, браслеты, ожерелья, броши и т, д	26, улица Кирова, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fe07794a-1f08-4d05-8569-a90fa9c75a56.png}	1	7	120	2025-12-02 10:58:49.981	2025-12-03 19:00:32.464	\N	\N	f	APPROVED	\N
24	Разобранный кубик рубика	10	USED	не смог собрать	Оренбургский Колледж Экономики и Информатики, 11, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c55064c4-df45-485c-80a3-7253e48ff798.jfif,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/d4b88f23-5430-4948-96f8-b521befb052b.jpg}	1	2	13	2025-11-28 09:19:52.442	2025-12-03 19:00:32.469	\N	\N	f	APPROVED	\N
28	Тест	20000	NEW	Описание	Вита Экспресс, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f76a26f7-4801-4c7f-9166-6b2869b5a765.jpg}	1	8	20	2025-12-01 05:50:33.37	2025-12-03 19:00:32.471	\N	\N	f	APPROVED	\N
61	Кресло-коляска	45000	USED	новая	г Оренбург	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/0a43c4d3-222d-421b-8528-7e3e59cc909a.jpg}	1	15	53	2025-12-01 08:10:41.21	2025-12-03 19:00:32.475	\N	\N	f	APPROVED	\N
127	Медицинское кресло	15798	NEW	Инвалидное кресло для комфортной и активной жизни.\r\n*  Мягкое сиденье и удобная спинка обеспечат комфорт даже при длительном использовании. Легко складывается для транспортировки.\r\n*  Регулируется под индивидуальные потребности. [Указать преимущества, например, наличие подголовника, антиопрокидыватели. \r\n\r\n✈✈✈✈✈ Можно отправить!\r\n\r\nЦена реальная. Звоните или пишите" 	г Оренбург, пр-кт Победы, д 10	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5fe6a6e0-d9a6-418d-bca1-dda8509a758f.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/659a2e57-a129-468f-9dea-c500bce1dcaa.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b4410b84-3d31-43b7-a9b0-0d10516b503b.jpg}	1	5	86	2025-12-01 09:07:28.717	2025-12-03 19:00:32.482	\N	\N	f	APPROVED	\N
229	Кресло офисное	5000	NEW	Удобное кресло	г Оренбург, ул Богдана Хмельницкого	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a0763806-d3e2-4699-9674-487603f386a3.jpg}	1	12	121	2025-12-02 10:56:54.12	2025-12-03 19:00:32.484	\N	\N	f	APPROVED	\N
227	Платье	2000	NEW	Платье летнее разных расцветок	2, улица Богдана Хмельницкого, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460024, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/497ae48d-8b31-443d-8a6f-9f00be0ac793.png}	1	1	7	2025-12-02 10:54:44.305	2025-12-03 19:00:32.489	\N	\N	f	APPROVED	\N
226	Ложка	100	NEW	Просто ложка	г Оренбург, ул Львовская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1582c805-78ef-49e4-a5bc-c875f429af60.webp}	1	11	121	2025-12-02 10:54:33.568	2025-12-03 19:00:32.491	\N	\N	f	APPROVED	\N
194	Вилка	100	NEW	Просто вилка	г Оренбург, ул Одесская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/d66fc1ce-a773-4ff1-8d0c-eaaf5515e495.webp}	1	11	121	2025-12-02 10:51:56.847	2025-12-03 19:00:32.495	\N	\N	f	APPROVED	\N
193	Джинсы	2500	NEW	Джинсы в новом состояние	48, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/935c3f3a-3998-4c35-97aa-d83c3b4c3beb.png}	1	1	7	2025-12-02 10:51:54.306	2025-12-03 19:00:32.497	\N	\N	f	APPROVED	\N
230	Свитер	3000	NEW	Свитер теплый из мягкой ткани	3/1, Телевизионный переулок, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460024, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/8ee766c3-6a41-4953-881d-c48ff14a1add.png}	1	1	7	2025-12-02 10:57:39.627	2025-12-03 19:00:32.498	\N	\N	f	APPROVED	\N
231	Садовые качели	10000	NEW	Просто качели. Качаться весело	г Оренбург, ул Шевченко	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5ab65b78-ba32-48e4-aa0e-a45203f815ab.jpg}	1	12	121	2025-12-02 10:58:20.061	2025-12-03 19:00:32.501	\N	\N	f	APPROVED	\N
234	Табурет	500	NEW	Просто табурет.	г Оренбург, ул Полтавская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/2a965649-9cf7-4e34-a427-8926bc88b2c9.jpg}	1	12	121	2025-12-02 11:00:14.233	2025-12-03 19:00:32.505	\N	\N	f	APPROVED	\N
236	Диван	6000	NEW	Просто диван	г Оренбург, ул Гоголя	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/28a67960-b058-4732-8bf1-501b4d4cca5a.webp}	1	12	121	2025-12-02 11:01:04.406	2025-12-03 19:00:32.511	\N	\N	f	APPROVED	\N
237	Дубленка	8000	NEW	Дубленка зимняя	улица Рокоссовского, Горка, Дзержинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b6118a2c-fce6-4c04-826b-b71e8953afe4.png}	1	1	7	2025-12-02 11:01:27.168	2025-12-03 19:00:32.512	\N	\N	f	APPROVED	\N
238	Кровать	10000	NEW	Удобная кровать. Евродвушка	г Оренбург, Крымский пер	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a83a8772-60ab-46e0-ac5e-3130ef9deb81.webp}	1	12	121	2025-12-02 11:02:27.961	2025-12-03 19:00:32.515	\N	\N	f	APPROVED	\N
240	Тонометр	2000	NEW	Тонометр. Давление меряет еще что-то там	г Оренбург, ул Донецкая	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/50400626-7e5c-41a0-8b2a-d7d753a627cd.jpg}	1	15	121	2025-12-02 11:03:58.743	2025-12-03 19:00:32.52	\N	\N	f	APPROVED	\N
241	Серебряные украшения	1500	USED	Пoд номeром 1: сеpежки с розoвым камнeм 1000 рублей. Под номeрoм 2: нaбop cepeжки и кольцо с жeлтым кaмнeм 2000 рублей зa нaбор. Под номером 3: набoр cepeжки, кольцо и подвеcкa с зелeным кaмнeм 2000 pублей зa набоp. Под нoмеpoм 5: сepежки с рoзoвым кaмнем 1000 рублeй. Серебрянaя цeпoчка 2000 рублей. Кольцо с белым камнем и две подвески с белыми камнями- по 500 рублей каждая. Серебро все в хорошем состоянии	Фармленд, 52, Советская улица, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/053660d3-2e5f-4f54-a2f8-c11767cd53fc.png}	1	7	120	2025-12-02 11:04:37.805	2025-12-03 19:00:32.524	\N	\N	f	APPROVED	\N
242	Детские игрушки	1000	USED	Набор детский игрушек	24, Луговая улица, Восточный, Сотки, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e8e9b6ba-a396-477d-b8c3-000cc9e85c0f.png}	1	2	7	2025-12-02 11:04:58.661	2025-12-03 19:00:32.526	\N	\N	f	APPROVED	\N
247	Ванночка	3000	USED	Ванна для купания новорожденного	6Б, Телевизионный переулок, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460024, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/ef7cf831-d301-4470-8e08-45a9662ccc25.png}	1	2	7	2025-12-02 11:08:32.807	2025-12-03 19:00:32.527	\N	\N	f	APPROVED	\N
249	Кроватка	3000	USED	Кроватка для новорожденных	199, Комсомольская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/42b2d16e-20eb-456c-8b73-3237d225a549.png}	1	2	7	2025-12-02 11:10:27.655	2025-12-03 19:00:32.532	\N	\N	f	APPROVED	\N
239	Украшения в русском стиле	2800	NEW	Украшения в русском стиле из натуральных камней и керамических бусин с подвесками ручной работы: неваляшки, Петушки, лошадки.	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1c859471-e841-471b-a425-cd312047cd68.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/9114495f-1926-48c5-bcbd-336ef851b323.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3dd7eee0-ee94-440c-84e0-16d5dfb7740e.png}	1	7	120	2025-12-02 11:03:22.662	2025-12-09 06:13:48.017	\N	\N	f	APPROVED	\N
253	Кровать	10000	USED	Кровать детская	139, Ташкентская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5b84facc-aec9-43ff-aa6e-e6a5922f31f2.png}	1	2	7	2025-12-02 11:13:18.658	2025-12-03 19:00:32.538	\N	\N	f	APPROVED	\N
254	Пара Флэт Уайт	363	NEW	Пара Флэт Уайт по выгодной цене. Доступно только в доставке!	30, улица 8 Марта, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/61625d76-10a4-423b-9040-f9871c898a6b.png}	1	8	120	2025-12-02 11:14:00.643	2025-12-03 19:00:32.544	\N	\N	f	APPROVED	\N
255	Катетер	150	NEW	Просто катетер	г Оренбург, ул Севастопольская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5a2c2b79-558b-46e7-b267-6ae194c526b9.jpg}	1	15	121	2025-12-02 11:15:13.932	2025-12-03 19:00:32.546	\N	\N	f	APPROVED	\N
256	Кофеин в таблетках	160	NEW	Просто кофеин	г Оренбург, мкр Ростошинские пруды, Керченский пер	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/845037dc-1299-407d-a608-b0bfccd6de8d.webp}	1	3	121	2025-12-02 11:16:08.622	2025-12-03 19:00:32.549	\N	\N	f	APPROVED	\N
257	Часы	500	NEW	Часы громкоговорители	113, Невельская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b7eb5ecd-4fb4-4d07-ac74-b5146c8080ba.png}	1	5	7	2025-12-02 11:16:24.995	2025-12-03 19:00:32.552	\N	\N	f	APPROVED	\N
258	Домашние полуфабрикаты, пельмени и тд	650	NEW	Прoдаём cвoю дoмaшнюю пpодукцию из магазина и пpинимаeм закaзы.Продукция oчeнь вкусная, из домaшниx яиц. Xaляль. Фaрш делаeм caми, ни одной жилки плёнки тaм нeт.\r\n	улица Цвиллинга, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c14573ef-45ac-488a-a7c5-663ffee7150e.png}	1	8	120	2025-12-02 11:16:29.78	2025-12-03 19:00:32.554	\N	\N	f	APPROVED	\N
263	Померанский шпиц, щенок	1	USED	Продаетcя очapовательная мини дeвочкa помepанcкoгo шпицa.28.09.2025 гoдa poждeния.\r\nДoкументы: Вет пacпoрт прививки oбpаботки по возрасту.\r\nОчeнь лаcкoвaя игpивая контактная .\r\nПpиучeна к пелeнки.\r\nKушaeт суxой коpм\r\nОтличнo ладит c дeтьми и другими живoтными .\r\nИщeм добрыe зaботливыe руки.\r\nРoдитeли:\r\nМама - померaнский шпиц, белый окрас (3,5 кг)\r\nПапа - померанский шпиц, пати колор (3 кг)\r\nБудет не больше 2,5 кг.	77/2, улица Терешковой, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/26aadd6d-3f95-4315-9a51-c59257705c32.png}	1	9	120	2025-12-02 11:21:39.796	2025-12-03 19:00:32.557	\N	\N	f	APPROVED	\N
265	Средство для удаления тейпов	500	NEW	Средство для удаления тейпов	41, улица Терешковой, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/9092c6b3-18bf-4960-b6c0-ae294784dd18.png}	1	5	7	2025-12-02 11:24:04.588	2025-12-03 19:00:32.559	\N	\N	f	APPROVED	\N
266	Котёнок в добрые руки	1	USED	котёнок около 4 месяцев, стерелизован, мальчикрыжий, очень активный, игривый, с другими животными и детьми ладит. очень ласковый, постоянно мурчит	14, улица Терешковой, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/333a3c16-155f-4e82-9fbe-a878937a6f9f.png}	1	9	120	2025-12-02 11:24:58.868	2025-12-03 19:00:32.56	\N	\N	f	APPROVED	\N
273	Шампунь Гарньер	500	NEW	Просто шампунь	г Оренбург, ул Сумская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/72f30e85-f179-47a0-a982-11bf90113e0e.jpg}	1	3	121	2025-12-02 11:29:29.551	2025-12-03 19:00:32.565	\N	\N	f	APPROVED	\N
272	Концелярия	700	NEW	Канцелярия для школы набор линеек y2k эстетика бант кролик	128, Орская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5e603622-4d00-4003-a685-cca9d4c78cf7.png}	1	6	7	2025-12-02 11:29:06.905	2025-12-03 19:00:32.566	\N	\N	f	APPROVED	\N
270	Духи 	3500	NEW	Духи Dior Sauvage	г Оренбург, ул Черниговская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fa8257b2-6b96-4b79-b10e-df3d7c723129.jpg}	1	3	121	2025-12-02 11:27:20.06	2025-12-03 19:00:32.572	\N	\N	f	APPROVED	\N
276	Блокнот	300	NEW	Блокнот Осенняя эстетика	92, улица Орджоникидзе, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/28d3f9c3-0ef8-4de7-8460-8d402294aa14.png}	1	6	7	2025-12-02 11:31:05.266	2025-12-03 19:00:32.578	\N	\N	f	APPROVED	\N
277	Тени для век	2000	NEW	Просто тени	г Оренбург, ул Житомирская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3a20e814-cdaf-40a4-8252-fd825161c268.webp}	1	3	121	2025-12-02 11:31:50.394	2025-12-03 19:00:32.58	\N	\N	f	APPROVED	\N
279	Морозилки ларь Бирюса, Pozis, Kraft и другие	15990	NEW	Бoльшой выбoр мopoзильныx камер (вepтикaльныe, лapи) разных oбъёмoв в нaличии в Орeнбуpге!\r\n\r\nА так же в наличии огромный выбoр бытoвoй тexники по оптовым ценaм!\r\n\r\n	Вишнёвая улица, СНТ "ЮЖНЫЙ УРАЛ ОФИЦЕРОВ ЗАПАСА И ОТСТАВКИ", Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/598ec012-3f41-4d02-9721-750964a49125.png}	1	10	120	2025-12-02 11:32:20.271	2025-12-03 19:00:32.583	\N	\N	f	APPROVED	\N
280	Стиральная машина бу	7000	USED	Стиральныe машины б.у. 🚛 Бecплaтная доставкa по гoроду ✅Гарaнтия до 12 меcяцeв пo чeку + пocлeгарантийнoe oбслуживаниe.\r\n\r\n	25, Краснознамённая улица, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/6f7d4e74-5748-4fa3-b1d3-a22f8aa6a061.png}	1	10	120	2025-12-02 11:33:16.603	2025-12-03 19:00:32.585	\N	\N	f	APPROVED	\N
281	Закладки для учебников 	300	NEW	Закладки для книг, «Книжная эстетика»	5, улица Макаровой-Мутновой, Новостройка, Промышленный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a2609c6c-a3c3-4d4a-bd93-60020e455210.png}	1	6	7	2025-12-02 11:33:30.034	2025-12-03 19:00:32.587	\N	\N	f	APPROVED	\N
269	Концтовары	700	NEW	Набор канцтоваров для школы и офиса Лапки котика 5 предметов	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e97b5f41-387d-42c2-bc90-a342ab3403bb.png}	1	6	7	2025-12-02 11:27:09.759	2025-12-09 06:13:48.017	\N	\N	f	APPROVED	\N
275	Наклейки	700	NEW	Наклейки для ежедневника Школьная эстетика	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/ab80ea7d-50f6-4bcf-beff-ec2a68d97299.png}	1	6	7	2025-12-02 11:30:16.141	2025-12-09 06:13:48.017	\N	\N	f	APPROVED	\N
278	Пенал	590	NEW	Милый эстетичный большой пенал школьный	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5090f83f-e7b2-42af-bc8e-2937850f8952.png}	1	6	7	2025-12-02 11:31:59.433	2025-12-09 06:13:48.017	\N	\N	f	APPROVED	\N
\.


--
-- Data for Name: ProductFieldValue; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."ProductFieldValue" (id, value, "fieldId", "productId") FROM stdin;
28	52	1	16
29	Черный	2	16
30	Шелк	3	16
31	Gucci	4	16
32	Платье	5	16
33	Горничная	6	16
34	50	1	22
35	Радужный	2	22
36	Хлопок	3	22
37	Demix	4	22
38	Аниме	5	22
39	Да	6	22
40	15	54	24
41	разный	58	24
42	3х3	62	24
43	не знаю	67	24
48	1-5	54	27
49	Книжный	58	27
50	Книжный	62	27
51	Книжные	67	27
52	черный	86	61
85	Активная модель 1000	136	127
86	, прогулочная	137	127
87	Железяка	138	127
88	16 кг 500 грамм	139	127
89	120 кгили 0,12 тонны, или 1,2 центнера	140	127
90	низкопрофильные	141	127
91	нет	142	127
92	ручное	143	127
93	мягкая сидушка	144	127
94	да	145	127
95	черный	146	127
118	s	7	193
119	голубой	8	193
120	джинса	9	193
121	Gloria	10	193
122	джинсы	11	193
123	багги	12	193
151	XS-L	7	227
152	разные	8	227
153	хлопок	9	227
154	Dasha	10	227
155	DV	11	227
156	летние	12	227
157	S-L	7	230
158	белый	8	230
159	норка	9	230
160	red	10	230
161	Sweet	11	230
162	свитер	12	230
163	S-L	7	233
164	бежевый	8	233
165	хлопок	9	233
166	Bant	10	233
167	BD	11	233
168	зипка	12	233
169	s	7	237
170	голубо-белый	8	237
171	джинса	9	237
172	VK	10	237
173	vk	11	237
174	дубленка	12	237
175	Белый	84	240
176	Да	87	240
177	Omron	90	240
178	Да	95	240
179	2-3	51	242
180	5	59	242
181	разные	63	242
182	15-30 см	64	242
183	Белый	84	243
184	Да	87	243
185	Omron	90	243
186	Да	95	243
187	50см	50	244
188	грязно-синий	57	244
189	1-3	60	244
190	50см	70	244
191	Да	85	246
192	Zenet	88	246
193	Нет	92	246
194	Синий	94	246
195	белый	49	247
196	100см	55	247
197	0-1	61	247
198	200см	73	247
199	100см	53	249
200	0-1	65	249
201	100см	72	249
202	белый	75	249
203	Да	85	250
204	KMED	88	250
205	Да	92	250
206	Белый	94	250
207	200см	53	251
208	0-1	65	251
209	200см	72	251
210	бело-коричневый	75	251
211	1,5 метра	53	253
212	0-7	65	253
213	1,5 метра	72	253
214	белый	75	253
215	Да	85	255
216	KMD	88	255
217	Да	92	255
218	Белый	94	255
219	Таблетки	83	256
220	часы	187	257
221	америка	188	257
222	часы	189	257
223	60см	190	257
224	Крем	82	259
225	Компрессорный ингалятор	187	261
226	Omron Comp Air NE-C300 Complete	188	261
227	Небулайзер OMRON C300 Complete — прибор, работающий в 3 режимах ингаляции. 	189	261
228	70см	190	261
229	кресло-коляска	136	262
230	сидячий	137	262
231	метал	138	262
232	7кг	139	262
233	100кг	140	262
234	 Колесная база, не выступающая за габариты коляски	141	262
235	нет	142	262
236	автоматическое	143	262
237	нет	144	262
238	есть	145	262
239	черный	146	262
240	clinar	132	265
241	балончик	133	265
242	2 года	134	265
243	америка	135	265
244	уход	132	267
245	спрей	133	267
246	5 лет	134	267
247	dinax	135	267
248	Духи	77	270
249	Шампунь	78	273
250	Палетка с тенями	80	277
251	папва	117	348
252	папва	118	348
253	апавпва	119	348
254	ир	120	348
255	рпрпр	121	348
256	пнпнп	122	348
\.


--
-- Data for Name: ProductPromotion; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."ProductPromotion" (id, "productId", "promotionId", "userId", days, "totalPrice", "startDate", "endDate", "isActive", "isPaid", "createdAt", "updatedAt") FROM stdin;
\.


--
-- Data for Name: ProductView; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."ProductView" (id, "viewedById", "productId", "viewedAt") FROM stdin;
8	7	17	2025-11-28 09:15:34.351
5	7	16	2025-11-28 09:15:40.734
6	7	14	2025-11-28 09:15:46.542
11	7	13	2025-11-28 09:15:50.575
12	13	16	2025-11-28 09:16:05.662
13	11	16	2025-11-28 09:16:05.942
14	14	18	2025-11-28 09:16:42.097
15	14	17	2025-11-28 09:16:48.332
16	7	20	2025-11-28 09:17:09.157
17	7	19	2025-11-28 09:17:15.39
18	7	18	2025-11-28 09:17:21.48
19	7	7	2025-11-28 09:17:47.718
20	7	9	2025-11-28 09:18:00.208
21	7	8	2025-11-28 09:18:07.082
7	7	10	2025-11-28 09:18:13.977
23	7	5	2025-11-28 09:18:17.402
24	7	22	2025-11-28 09:18:26.145
25	7	21	2025-11-28 09:18:35.85
26	11	22	2025-11-28 09:18:54.527
29	14	5	2025-11-28 09:19:16.228
28	9	5	2025-11-28 09:19:28.543
31	19	20	2025-11-28 09:19:39.18
32	7	24	2025-11-28 09:20:15.947
33	5	24	2025-11-28 09:20:19.239
35	18	16	2025-11-28 09:20:26.159
27	18	22	2025-11-28 09:20:38.981
39	18	19	2025-11-28 09:20:43.139
40	18	17	2025-11-28 09:20:48.02
42	18	12	2025-11-28 09:20:52.446
43	13	18	2025-11-28 09:20:59.678
44	9	21	2025-11-28 09:21:04.291
45	9	24	2025-11-28 09:21:06.787
46	13	20	2025-11-28 09:21:09.011
47	9	23	2025-11-28 09:21:09.232
48	9	20	2025-11-28 09:21:12.953
49	9	19	2025-11-28 09:21:16.127
51	13	15	2025-11-28 09:21:19.648
52	9	27	2025-11-28 09:21:21.713
54	13	8	2025-11-28 09:21:24.529
70	13	14	2025-11-28 09:22:21.23
71	13	27	2025-11-28 09:22:26.679
105	13	9	2025-11-28 09:49:35.399
106	13	7	2025-11-28 09:50:04.663
239	5	18	2025-12-01 07:48:19.685
270	53	11	2025-12-01 08:12:43.181
336	86	28	2025-12-01 08:37:06.268
369	5	61	2025-12-01 08:48:56.032
468	5	94	2025-12-01 09:22:41.051
237	5	27	2025-12-01 09:24:12.715
536	5	19	2025-12-01 09:24:23.12
206	5	22	2025-12-01 09:26:15.811
569	5	17	2025-12-01 09:26:25.118
173	5	16	2025-12-01 09:26:28.453
204	5	14	2025-12-01 09:26:35.713
238	5	15	2025-12-01 09:26:48.769
305	5	13	2025-12-01 12:03:10.898
602	119	10	2025-12-01 14:29:41.095
603	119	28	2025-12-01 14:29:43.789
604	119	24	2025-12-01 14:37:03.378
605	119	23	2025-12-01 14:46:32.887
606	119	27	2025-12-01 14:54:23.873
205	5	28	2025-12-02 06:32:52.922
567	5	127	2025-12-02 07:32:36.871
504	5	12	2025-12-02 07:39:03.095
50	5	21	2025-12-02 07:46:38.367
617	121	252	2025-12-02 11:33:23.238
618	120	94	2025-12-02 11:42:27.288
651	5	277	2025-12-03 09:51:19.021
652	5	247	2025-12-03 09:51:25.113
654	5	252	2025-12-03 09:52:30.64
656	5	234	2025-12-03 09:52:39.77
653	5	264	2025-12-03 16:38:43.192
658	5	316	2025-12-03 16:38:47.965
660	5	242	2025-12-03 16:41:06.056
662	5	263	2025-12-03 16:54:39.706
72	5	23	2025-12-03 16:55:51.09
172	5	20	2025-12-03 16:57:03.686
659	5	271	2025-12-03 17:01:43.056
\.


--
-- Data for Name: Promotion; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Promotion" (id, name, "pricePerDay", "createdAt", "updatedAt") FROM stdin;
1	Стандарт	50	2025-12-08 12:37:51.475	2025-12-08 12:37:32.223
2	Люкс	100	2025-12-08 12:37:51.475	2025-12-08 12:37:44.761
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
5	Средства реабилитации	1
6	Школьные товары	1
7	Украшения	1
8	Продукты питания	1
9	Животные, растения	1
10	Бытовая техника	1
11	Посуда	1
12	Мебель	1
15	Медицинские товары	1
3	Красота и здоровье	1
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

COPY public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance) FROM stdin;
7	дар	bdi-2006@mail.ru	+79123400130	$2b$10$TROWXU059pwS6Q98JIfGDOL1kzA0oohdraWoB3ZxpEgGqEU//.qQ6	INDIVIDUAL	2025-11-28 09:06:52.861	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
8	Исаев Максим Андреевич	sima.isaev2305@mail.ru	+79501859919	$2b$10$VI6Gb9KuiHWEnbndcyi1WemTTQgKWwVhpcOfnEEj7W18T8Gw.TPou	INDIVIDUAL	2025-11-28 09:06:55.938	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
9	Махар Святой Рог	vmahauri029@gmail.com	+79123557497	$2b$10$UbWFDK5KoI92FFzmWZw.s.jslpRNGreNJFQi30q4ZWI9lB02sqegS	INDIVIDUAL	2025-11-28 09:07:05.955	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
10	Голосняк Юлия Викторовна	juliagolosnyak@mail.ru	+79328538922	$2b$10$9VP3OmZRjdumTgAJWCBGGe5ozGVZG0Z/okvuWwUdx1wxmJG7brTES	INDIVIDUAL	2025-11-28 09:07:19.394	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
11	Захаров АР ВЛ	Zahar83s@mail.ru	+79878600551	$2b$10$TfLU49EmrMYrTPd46fQv6.QNkD3tEE2WnHVmy8qIdYzHVOX4PLe4q	INDIVIDUAL	2025-11-28 09:07:21.428	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
12	Подрядов Екатерина Сергеевна	podradovakata91@gmail.com	+79083234725	$2b$10$sdWaXECQtpyEqc61gS4MrOlsoz4nsjYb1gGC1xD2VVFgr/pUqwB3m	INDIVIDUAL	2025-11-28 09:07:29.962	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
13	Макаров Николай	bapenick445@gmail.com	+79225387481	$2b$10$DHSa1l.0cj7MK.b7ATupL.f7yXnjfGBUEr7Wezf1wul9x2z2eOIkO	INDIVIDUAL	2025-11-28 09:07:33.445	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
15	Кокеев Фирилл Батькович	test@test.com	+79953501391	$2b$10$0GEA/Uvq4NrHTLuOetQTXuoviQG19DrdEX4NIFUwD.54aF7ePJveO	INDIVIDUAL	2025-11-28 09:07:44.576	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
17	Абвгдеивич Егор Константинович	barabulkabarabulka@gmail.com	+72280303111	$2b$10$PPEwZxCaLahLuE4XtqI2k.UxgqrcfBgCoXBHT1EUoq86kYraokwz2	INDIVIDUAL	2025-11-28 09:08:14.573	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
16	kostyukov	geronimoprofitop@gmail.com	+79228744883	$2b$10$ulXOXoQl7aAYjf7uJ2opGOApWYjLTVFSWBrWyYAjJp80HAeDl97OS	INDIVIDUAL	2025-11-28 09:07:57.477	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
18	Прокофьева Валерия Денисовна	lin.ferr@mail.ru	+79225406669	$2b$10$7mnxrJ2LJ0S5RoBoo8gVteXYR.o2kM/nnm07SpxHT37YZqEghfVAC	INDIVIDUAL	2025-11-28 09:08:42.207	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
19	Гатин Ян Талгатович	ggg2107@gmail.com	+79228386030	$2b$10$aUbIJdrSn4qPvErIPV8E6uo162lESkmE7orVVIrS/2v8/k8qUQjvm	INDIVIDUAL	2025-11-28 09:08:47.126	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
20	Арзамасцев Даниил	arzamastsevdaniil@gmail.com	+79068346355	$2b$10$NvJVMH9Kn16C7hSuCtRAf./yj8/jgaeUg2ZI0IAkxt2Tc/Cf5DR8G	INDIVIDUAL	2025-12-01 05:48:10.726	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
86	Афонасьев Афиларет Михайлович	pr.actual@mail.ru	+79082734009	$2b$10$R0pbgCnq1AVwe9phmKu1GOT0emg48XzDbtYRBEn/xEyCFd8aNYX7y	INDIVIDUAL	2025-12-01 08:28:35.989	2025-12-08 12:30:51.217	\N	f	1	t	https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/users/71116356-ea56-4dd5-ac1a-86c5a6e2e11b.jpg	f	0
14	Каверина Мария	kunafina_ruslana7@mail.ru	+79228362555	$2b$10$AY/2V0DgPQ1.ZorhEmTMfOb4o8hq1EkOR9qkHx4/RgG7Cq6OFAOo2	INDIVIDUAL	2025-11-28 09:07:42.429	2025-12-08 12:30:51.217	\N	f	1	t	\N	f	0
122	Попов Матвей Иванович	trrina04@mail.ru	+79878993845	$2b$10$cfHgsH42YXRqYPpoZbbhAuFK4bg.81DSzN4JNMGmkLffNma7mLmB.	INDIVIDUAL	2025-12-03 19:26:12.827	2025-12-08 12:30:51.217	\N	f	1	f	\N	f	0
5	Попов Матвей Иванович	vitaly.sadikov1@yandex.ru	+79510341677	$2b$10$05FMyE494pfJScN9OF98COs6yLacnIIE2gueMbTS8s1/PNzaYrA6C	INDIVIDUAL	2025-11-06 19:33:46.625	2025-12-08 12:30:51.217	\N	f	3	f	https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/users/eac42b51-e66a-4d76-bad2-c6db0efd947b.jpg	t	0
6	Садиков Виталий Дмитриевич	vitaly.sadikov2@yandex.ru	+79510341676	$2b$10$Tsi0whXkdERT2AvjSe6Jn.v6ba.K3sTDPXT6AzWMlkpahIY.LxDSS	INDIVIDUAL	2025-11-06 19:33:55.742	2025-12-08 14:16:57.863	\N	f	1	\N	\N	f	500000
53	Корякина Ирина	ikoryakina47@gmail.com	+79228579009	$2b$10$48dtDNK6DIH0yBgup4eqeeG8k5NPkHuhqBNvQ2yCJqayB3sNthYOS	INDIVIDUAL	2025-12-01 08:08:29.883	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
119	йцукенгшщзх	qwertyui123@gmail.com	+75678903456	$2b$10$hhmWdTv8RdWeJ1ofHOjaTuKBgOo2JUky9za7NTJ.uCcfrH3W2CK/S	INDIVIDUAL	2025-12-01 14:29:11.538	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
121	Фокеев Кирилл	test1@test.com	+71234567890	$2b$10$FELoBjJj0J8IeMy2YhKlIeniLkjz86fijJS2HOFJ3XvJ3fnIulg2i	INDIVIDUAL	2025-12-02 10:48:41.186	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
120	Черешков Данила Алексеевич	chereshkov.da2006@gmail.com	+79123431910	$2b$10$hvt0jXBTO6PcqEzKYDKYUO7hivY2kCsC/7Bzwix242L8YDeP6UgnW	INDIVIDUAL	2025-12-02 10:47:25.87	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0
\.


--
-- Data for Name: _UserFavorites; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."_UserFavorites" ("A", "B") FROM stdin;
20	14
6	14
10	14
6	7
23	5
28	86
28	119
127	5
94	5
21	5
282	20
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

SELECT pg_catalog.setval('public."Chat_id_seq"', 145, true);


--
-- Name: FavoriteAction_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."FavoriteAction_id_seq"', 79, true);


--
-- Name: Message_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Message_id_seq"', 2, true);


--
-- Name: OkseiProduct_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."OkseiProduct_id_seq"', 1, true);


--
-- Name: PhoneNumberView_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."PhoneNumberView_id_seq"', 1, true);


--
-- Name: ProductFieldValue_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."ProductFieldValue_id_seq"', 256, true);


--
-- Name: ProductPromotion_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."ProductPromotion_id_seq"', 1, false);


--
-- Name: ProductView_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."ProductView_id_seq"', 667, true);


--
-- Name: Product_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Product_id_seq"', 348, true);


--
-- Name: Promotion_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Promotion_id_seq"', 2, true);


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

SELECT pg_catalog.setval('public."User_id_seq"', 122, true);


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
-- Name: OkseiProduct OkseiProduct_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."OkseiProduct"
    ADD CONSTRAINT "OkseiProduct_pkey" PRIMARY KEY (id);


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
-- Name: ProductPromotion ProductPromotion_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductPromotion"
    ADD CONSTRAINT "ProductPromotion_pkey" PRIMARY KEY (id);


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
-- Name: Promotion Promotion_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Promotion"
    ADD CONSTRAINT "Promotion_pkey" PRIMARY KEY (id);


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
-- Name: Promotion_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Promotion_name_key" ON public."Promotion" USING btree (name);


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
-- Name: ProductPromotion ProductPromotion_productId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductPromotion"
    ADD CONSTRAINT "ProductPromotion_productId_fkey" FOREIGN KEY ("productId") REFERENCES public."Product"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: ProductPromotion ProductPromotion_promotionId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductPromotion"
    ADD CONSTRAINT "ProductPromotion_promotionId_fkey" FOREIGN KEY ("promotionId") REFERENCES public."Promotion"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: ProductPromotion ProductPromotion_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ProductPromotion"
    ADD CONSTRAINT "ProductPromotion_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


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
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT CREATE ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--


