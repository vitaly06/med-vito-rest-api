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
-- Name: Banner; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Banner" (
    id integer NOT NULL,
    "photoUrl" text NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."Banner" OWNER TO postgres;

--
-- Name: Banner_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Banner_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Banner_id_seq" OWNER TO postgres;

--
-- Name: Banner_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Banner_id_seq" OWNED BY public."Banner".id;


--
-- Name: Category; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Category" (
    id integer NOT NULL,
    name text NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    slug text DEFAULT ''::text NOT NULL,
    "updatedAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
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
    "productId" integer,
    "buyerId" integer NOT NULL,
    "sellerId" integer NOT NULL,
    "unreadCountBuyer" integer DEFAULT 0 NOT NULL,
    "unreadCountSeller" integer DEFAULT 0 NOT NULL,
    "lastMessageId" integer,
    "lastMessageAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL,
    "isModerationChat" boolean DEFAULT false NOT NULL
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
    "updatedAt" timestamp(3) without time zone NOT NULL,
    "relatedProductId" integer
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
-- Name: Payment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Payment" (
    id integer NOT NULL,
    "orderId" text NOT NULL,
    "paymentId" text NOT NULL,
    "userId" integer NOT NULL,
    amount double precision NOT NULL,
    status text NOT NULL,
    "paymentUrl" text,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."Payment" OWNER TO postgres;

--
-- Name: Payment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Payment_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Payment_id_seq" OWNER TO postgres;

--
-- Name: Payment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Payment_id_seq" OWNED BY public."Payment".id;


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
    "categoryId" integer NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    slug text DEFAULT ''::text NOT NULL,
    "updatedAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
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
    "subcategoryId" integer NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    slug text DEFAULT ''::text NOT NULL,
    "updatedAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
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
    balance double precision DEFAULT 0 NOT NULL,
    "bonusBalance" double precision DEFAULT 0 NOT NULL,
    "isBanned" boolean DEFAULT false NOT NULL,
    "freeAdsLimit" integer DEFAULT 12 NOT NULL,
    "lastAdLimitReset" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "usedFreeAds" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."User" OWNER TO postgres;

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
-- Name: Banner id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Banner" ALTER COLUMN id SET DEFAULT nextval('public."Banner_id_seq"'::regclass);


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
-- Name: Payment id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Payment" ALTER COLUMN id SET DEFAULT nextval('public."Payment_id_seq"'::regclass);


--
-- Name: PhoneNumberView id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."PhoneNumberView" ALTER COLUMN id SET DEFAULT nextval('public."PhoneNumberView_id_seq"'::regclass);


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
-- Data for Name: Banner; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Banner" (id, "photoUrl", "createdAt", "updatedAt") FROM stdin;
\.


--
-- Data for Name: Category; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Category" (id, name, "createdAt", slug, "updatedAt") FROM stdin;
1	Личные вещи	2025-12-15 19:18:08.497	lichnye-veschi	2025-12-15 17:21:30.479
\.


--
-- Data for Name: Chat; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") FROM stdin;
42	9863001	2681599	6053931	0	0	\N	2025-12-01 08:12:52.805	2025-12-01 08:12:52.805	2025-12-01 08:12:52.805	f
1	7659684	6038643	7391202	0	0	\N	2025-11-28 09:15:51.647	2025-11-28 09:15:51.647	2025-11-28 09:15:51.647	f
2	9915250	9371169	8633592	0	0	\N	2025-11-28 09:16:07.289	2025-11-28 09:16:07.289	2025-11-28 09:16:07.289	f
3	9915250	8964288	8633592	0	0	\N	2025-11-28 09:16:15.369	2025-11-28 09:16:15.369	2025-11-28 09:16:15.369	f
4	9256863	5966833	8964288	0	0	\N	2025-11-28 09:16:50.74	2025-11-28 09:16:50.74	2025-11-28 09:16:50.74	f
5	5611056	6038643	1208299	0	0	\N	2025-11-28 09:17:36.46	2025-11-28 09:17:36.46	2025-11-28 09:17:36.46	f
7	1250840	4761896	7391202	0	0	\N	2025-11-28 09:20:00.074	2025-11-28 09:20:00.074	2025-11-28 09:20:00.074	f
141	1250840	7106521	7391202	0	0	\N	2025-12-01 12:30:18.663	2025-12-01 12:30:18.663	2025-12-01 12:30:18.663	f
6	9368305	9371169	8633592	0	0	\N	2025-11-28 09:18:57.647	2025-11-28 09:18:57.647	2025-11-28 09:18:57.647	f
9	7384341	7106521	2321239	0	0	\N	2025-11-28 09:43:43.309	2025-11-28 09:43:43.309	2025-11-28 09:43:43.309	f
143	7384341	7132269	2321239	0	0	\N	2025-12-01 14:46:54.623	2025-12-01 14:46:54.623	2025-12-01 14:46:54.623	f
144	2105765	7106521	3235109	0	1	1	2025-12-02 06:39:40.026	2025-12-02 06:32:54.893	2025-12-04 06:34:08.592	f
142	2105765	7132269	3235109	0	0	\N	2025-12-01 14:29:48.006	2025-12-01 14:29:48.006	2025-12-01 14:29:48.006	f
75	2105765	6669460	3235109	0	0	\N	2025-12-01 08:28:55.846	2025-12-01 08:28:55.846	2025-12-01 08:28:55.846	f
108	1122280	7106521	2681599	0	0	\N	2025-12-01 08:36:19.933	2025-12-01 08:36:19.933	2025-12-01 08:36:19.933	f
145	3437684	4146092	9851099	0	1	2	2025-12-02 11:33:33.784	2025-12-02 11:33:25.608	2025-12-02 11:47:19.434	f
\.


--
-- Data for Name: FavoriteAction; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."FavoriteAction" (id, "userId", "productId", "addedAt") FROM stdin;
5	5966833	3334788	2025-11-28 09:21:05.722
7	6038643	3334788	2025-11-28 09:21:12.339
6	5966833	1300264	2025-11-28 09:21:07.97
3	5966833	1250840	2025-11-28 09:18:22.746
77	7106521	9262881	2025-12-02 07:46:33.084
4	8633592	9368305	2025-11-28 09:18:35.259
8	7106521	7384341	2025-11-29 09:00:09.809
74	7132269	2105765	2025-12-01 14:29:45.996
41	6669460	2105765	2025-12-01 09:22:11.203
76	7106521	2161612	2025-12-02 07:34:18.185
75	7106521	4758351	2025-12-02 07:34:17.402
78	3235109	8776759	2025-12-03 00:00:00.214
\.


--
-- Data for Name: Message; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Message" (id, content, "senderId", "chatId", "isRead", "readAt", "createdAt", "updatedAt", "relatedProductId") FROM stdin;
1	тест	7106521	144	f	\N	2025-12-02 06:39:40.018	2025-12-02 06:39:40.018	\N
2	Куда цену задрал? 200 край	4146092	145	f	\N	2025-12-02 11:33:33.781	2025-12-02 11:33:33.781	\N
\.


--
-- Data for Name: OkseiProduct; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."OkseiProduct" (id, name, description, price, image, "createdAt") FROM stdin;
1	iPhone 15 Pro	Новый iPhone 15 Pro в отличном состоянии	120000	https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1c68c479-ade3-43ff-91eb-b8428b46ed74.jpg	2025-12-12 08:53:25.175
\.


--
-- Data for Name: Payment; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Payment" (id, "orderId", "paymentId", "userId", amount, status, "paymentUrl", "createdAt", "updatedAt") FROM stdin;
1	7106521-1766690209318	7629971919	7106521	1000	PENDING	https://pay.tbank.ru/ahpkMYdA	2025-12-25 19:16:49.644	2025-12-25 19:16:49.644
2	7106521-1766690350760	7629983326	7106521	1000	PENDING	https://pay.tbank.ru/mYjJpEkG	2025-12-25 19:19:11.103	2025-12-25 19:19:11.103
3	7106521-1766690451846	7629991661	7106521	1000	PENDING	https://pay.tbank.ru/AqmD5LpC	2025-12-25 19:20:52.116	2025-12-25 19:20:52.116
4	7106521-1766690912537	7630030365	7106521	1000	PENDING	https://pay.tbank.ru/4KnshkYJ	2025-12-25 19:28:32.789	2025-12-25 19:28:32.789
5	7106521-1766690970273	7630035082	7106521	10	PENDING	https://pay.tbank.ru/w0hLiyV6	2025-12-25 19:29:30.546	2025-12-25 19:29:30.546
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
5231119	Бусы б/у	1000	USED	Красные, из жемчуга	г Екатеринбург, ул Чкалова	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png}	1	3	7391202	2025-11-28 09:11:49.34	2025-12-16 09:53:24.161	16	\N	f	APPROVED	\N
8640334	Нутриен energy питание	2700	NEW	Cмecь Nutrien enеrgy, диетичeское лечeбноe  питание,\r\n\r\nПитаниe для oнкoбольных , питаниe для ocлaблeнных, питание пocлe опеpaции, питание, обогащённое витаминaми и микрoэлeмeнтaми.\r\n\r\nПродукт готовый к упoтpеблению 200 мл, 300 ккaл.\r\nПoдxoдит для онкoбольных, пoслeопepациoнныx взрослыx и дeтeй с 3 лет для вoсстановления сил и энергии.	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/6f42b7d6-7a7d-46ca-85c3-611b159a8a0a.png}	1	8	9851099	2025-12-02 11:08:58.23	2025-12-16 09:53:24.246	51	\N	f	APPROVED	\N
9305563	Капельница	200	NEW	Просто капельница	г Оренбург, ул Харьковская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b0a71fb2-6719-4085-b554-d16b5cf9b2a2.webp}	1	15	4146092	2025-12-02 11:11:52.698	2025-12-16 09:53:24.191	74	\N	f	APPROVED	\N
7384341	ВАЗ 2107	435000	NEW	Продаётся готовый проект под RDS. Соответствует всем стандартам турниров и сходок. Гарантия на проект год.	Степной, Дзержинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/af9ca37d-87d9-44bc-b0aa-b1fc99737315.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/7a511cc5-d998-49c3-8f53-5dd88abd875b.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1d044f53-c12f-4841-9f4b-b486e551411a.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5c8cfe09-cd90-489a-8aad-0f8c2e80f6f4.jpg}	1	10	2321239	2025-11-28 09:18:17.344	2025-12-16 09:53:24.18	61	\N	f	APPROVED	\N
2161612	Очередной товар дня!	35000	NEW	1) пусть будет текст\r\n2) здесь еще что-то\r\n**\r\n💥\r\n🟩\r\nККЕКЕКЕЕУУЦКУ""\r\n                                              ЦЕНТР\r\n          ТАБУЛЯЦИЯ СмещЕНИЕ\r\n\r\n	18, улица Расковой, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/6638b79f-2357-46ff-9010-ba9175ce50db.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c02fa4fd-6284-45a5-8cd7-61583db872fe.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fe8ad9f4-5664-4832-b5be-dc1f4df2adcf.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5488c8b5-af91-4294-85c6-0bb7d48145b6.jpg}	1	12	6669460	2025-12-01 08:35:56.623	2025-12-16 09:53:24.183	67	\N	f	APPROVED	\N
9122333	Зипка	5000	NEW	Кофта теплая на замке	35, улица 9 Января, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f64e149a-e711-4173-85ea-98db13c3ca1e.png}	1	1	6038643	2025-12-02 10:59:24.476	2025-12-16 09:53:24.186	2	\N	f	APPROVED	\N
7659684	Протеин 1000гр	1500	NEW	Вкус шоколад, 1000 грамм	2, улица 13-я Линия, Линии, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460040, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a745c362-a754-4a92-abb9-b8969bebead7.png}	1	8	7391202	2025-11-28 09:14:04.157	2025-12-16 09:53:24.199	51	\N	f	APPROVED	\N
9262881	Набор украшений для пирсинга	4000	NEW	\N	12А, Больничный проезд, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3c250f31-591b-4346-b1fb-1b3bf70f2c73.webp}	1	3	6053931	2025-11-28 09:17:54.801	2025-12-16 09:53:24.203	16	\N	f	APPROVED	\N
4215912	Детские книжки по математике	1000	USED	Превосходный источник знаний для вашего ребенка	Hawthorne Street, Кламат-Фолс, Klamath County, Орегон, 97601, Соединённые Штаты Америки	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5e8875c2-aec8-4b2f-b618-2e220defa9cf.webp,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/2248f3b6-63b3-44e7-84a4-ef35b2d7bcdc.jpg}	1	2	1208299	2025-11-28 09:21:17.846	2025-12-16 09:53:24.208	14	\N	f	APPROVED	\N
1512888	Померанский шпиц, щенок	1	USED	Продаетcя очapовательная мини дeвочкa помepанcкoгo шпицa.28.09.2025 гoдa poждeния.\r\nДoкументы: Вет пacпoрт прививки oбpаботки по возрасту.\r\nОчeнь лаcкoвaя игpивая контактная .\r\nПpиучeна к пелeнки.\r\nKушaeт суxой коpм\r\nОтличнo ладит c дeтьми и другими живoтными .\r\nИщeм добрыe зaботливыe руки.\r\nРoдитeли:\r\nМама - померaнский шпиц, белый окрас (3,5 кг)\r\nПапа - померанский шпиц, пати колор (3 кг)\r\nБудет не больше 2,5 кг.	77/2, улица Терешковой, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/26aadd6d-3f95-4315-9a51-c59257705c32.png}	1	9	9851099	2025-12-02 11:21:39.796	2025-12-16 09:53:24.422	54	\N	f	APPROVED	\N
5142108	Котёнок в добрые руки	1	USED	котёнок около 4 месяцев, стерелизован, мальчикрыжий, очень активный, игривый, с другими животными и детьми ладит. очень ласковый, постоянно мурчит	14, улица Терешковой, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/333a3c16-155f-4e82-9fbe-a878937a6f9f.png}	1	9	9851099	2025-12-02 11:24:58.868	2025-12-16 09:53:24.429	54	\N	f	APPROVED	\N
4332941	Графин в виде рыбы	500	NEW	Замечательный графин в виде рыбы	г Оренбург, ул Киевская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/67533ddf-2495-4aed-b405-a68922a398bf.jpg}	1	11	4146092	2025-12-02 10:50:32.345	2025-12-16 09:53:24.22	63	\N	f	APPROVED	\N
3982248	Сковорода антипригарная	1000	NEW	Сковорода. Можно пожарить все что угодно	г Оренбург, ул Днепропетровская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/299e425b-f6c0-49bb-a7e8-3c7c591ce39d.jpg}	1	11	4146092	2025-12-02 10:53:18.109	2025-12-16 09:53:24.223	63	\N	f	APPROVED	\N
5868178	папавпа	55454	NEW	павпвапа	«Урал», Ленинский район, Пригородный, Пригородный сельсовет, Оренбургский район, Оренбургская область, Приволжский федеральный округ, 460041, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/396aecf1-b980-45ec-bd5a-ba238f1fdefb.jpg}	1	5	7106521	2025-12-03 19:36:04.58	2025-12-16 09:53:24.232	27	\N	f	MODERATE	\N
4758351	Медицинское кресло	15798	NEW	Инвалидное кресло для комфортной и активной жизни.\r\n*  Мягкое сиденье и удобная спинка обеспечат комфорт даже при длительном использовании. Легко складывается для транспортировки.\r\n*  Регулируется под индивидуальные потребности. [Указать преимущества, например, наличие подголовника, антиопрокидыватели. \r\n\r\n✈✈✈✈✈ Можно отправить!\r\n\r\nЦена реальная. Звоните или пишите" 	г Оренбург, пр-кт Победы, д 10	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5fe6a6e0-d9a6-418d-bca1-dda8509a758f.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/659a2e57-a129-468f-9dea-c500bce1dcaa.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b4410b84-3d31-43b7-a9b0-0d10516b503b.jpg}	1	5	6669460	2025-12-01 09:07:28.717	2025-12-16 09:53:24.346	29	\N	f	APPROVED	\N
2693271	Стакан	200	NEW	Просто стакан.	г Оренбург, ул Житомирская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/4cfed187-e55c-4014-8c59-cd2450aca91e.jpg}	1	11	4146092	2025-12-02 10:55:45.209	2025-12-16 09:53:24.226	63	\N	f	APPROVED	\N
3563632	Ингалятор	2000	NEW	Ингалятор для ингаляций	г Оренбург, ул Луганская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/08d81171-c46d-4857-b7d9-bb7a983d5ab4.jpg}	1	15	4146092	2025-12-02 11:05:14.587	2025-12-16 09:53:24.234	73	\N	f	APPROVED	\N
2865910	Посуда детская	1500	NEW	Детская посуда для кормления	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f252c198-af82-42a9-a80b-e42a052caae3.png}	1	2	6038643	2025-12-02 11:06:21.305	2025-12-16 09:53:24.24	13	\N	f	APPROVED	\N
8882052	Алоэ вера лечебный 3 года, есть 1 год	200	USED	Алое Вера, лечебное 3х детки, есть однолетки	В. И. Ленину, Ленинская улица, Кузнечный, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/17ad0794-e12d-433a-9958-528bba02bf87.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/09b01358-7870-4abd-a66b-7cfcab7ecec9.png}	1	9	9851099	2025-12-02 11:26:21.919	2025-12-16 09:53:24.215	54	\N	f	APPROVED	\N
2016352	Украшения	1000	USED	Продам укрошенияБраслет -500\r\nСерьги - 300\r\nКольцо 10 - 250\r\nВсе вместе 1000	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a28160dd-9d06-4750-b6fe-6045f6a3df8b.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e49061f1-c633-4254-bc95-b9e06ae322ae.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/15b89ff9-deb0-4067-80a9-77151e9ad946.png}	1	7	9851099	2025-12-02 11:06:25.377	2025-12-16 09:53:24.242	48	\N	f	APPROVED	\N
8948419	Массажный стол	4000	NEW	Просто массажный стол	г Оренбург, ул Запорожская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/002d563c-310a-4287-ac15-5826d88e5d37.jpg}	1	15	4146092	2025-12-02 11:07:04.422	2025-12-16 09:53:24.237	74	\N	f	APPROVED	\N
6091694	Домашние полуфабрикаты, пельмени и тд	650	NEW	Прoдаём cвoю дoмaшнюю пpодукцию из магазина и пpинимаeм закaзы.Продукция oчeнь вкусная, из домaшниx яиц. Xaляль. Фaрш делаeм caми, ни одной жилки плёнки тaм нeт.\r\n	улица Цвиллинга, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c14573ef-45ac-488a-a7c5-663ffee7150e.png}	1	8	9851099	2025-12-02 11:16:29.78	2025-12-16 09:53:24.418	51	\N	f	APPROVED	\N
2273041	Спрей	600	USED	Защитная пленка для кожи	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/002307c3-8394-46f2-b18a-952a389efc6d.png}	1	5	6038643	2025-12-02 11:25:57	2025-12-16 09:53:24.276	28	\N	f	APPROVED	\N
6617171	Кресло-коляска	5000	USED	Кресло-коляска для инвалидов Ortonica Olvia 30	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/695dac97-c44a-42c0-9307-2cc8ff3bcaab.png}	1	5	6038643	2025-12-02 11:21:05.818	2025-12-16 09:53:24.271	29	\N	f	APPROVED	\N
6628130	Матрас	3000	USED	Матрас для восстанволения	61А, улица Орлова, Новостройка, Промышленный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/7d6cda50-98e5-42bc-9a31-08b62188d9fa.png}	1	5	6038643	2025-12-02 11:22:24.486	2025-12-16 09:53:24.259	26	\N	f	APPROVED	\N
1063797	Блокнот	300	NEW	Блокнот Осенняя эстетика	92, улица Орджоникидзе, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/28d3f9c3-0ef8-4de7-8460-8d402294aa14.png}	1	6	6038643	2025-12-02 11:31:05.266	2025-12-16 09:53:24.44	41	\N	f	APPROVED	\N
6193207	Вилка	100	NEW	Просто вилка	г Оренбург, ул Одесская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/d66fc1ce-a773-4ff1-8d0c-eaaf5515e495.webp}	1	11	4146092	2025-12-02 10:51:56.847	2025-12-16 09:53:24.358	63	\N	f	APPROVED	\N
4372887	Компрессорный ингалятор	2000	NEW	Компрессорный ингалятор	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f178392b-1965-4569-b91c-c6efd48b56da.png}	1	5	6038643	2025-12-02 11:18:32.771	2025-12-16 09:53:24.266	35	\N	f	APPROVED	\N
2854985	Садовые качели	10000	NEW	Просто качели. Качаться весело	г Оренбург, ул Шевченко	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5ab65b78-ba32-48e4-aa0e-a45203f815ab.jpg}	1	12	4146092	2025-12-02 10:58:20.061	2025-12-16 09:53:24.369	67	\N	f	APPROVED	\N
6734788	Кроватка	3000	USED	Кроватка для новорожденных	199, Комсомольская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/42b2d16e-20eb-456c-8b73-3237d225a549.png}	1	2	6038643	2025-12-02 11:10:27.655	2025-12-16 09:53:24.395	10	\N	f	APPROVED	\N
9265239	Люлька	2000	USED	Люлька детская	68, улица Кичигина, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/51e36a8d-6904-4a08-96bc-d1d449241608.png}	1	2	6038643	2025-12-02 11:11:54.214	2025-12-16 09:53:24.196	10	\N	f	APPROVED	\N
8507601	Крем для рук	200	NEW	Просто крем для рук	г Оренбург, поселок Нижнесакмарский, ул Николаевская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fb99e19c-4e45-43b7-90fc-1c48c50106d7.webp}	1	3	4146092	2025-12-02 11:17:30.55	2025-12-16 09:53:24.251	16	\N	f	APPROVED	\N
7270506	Пароочиститель для дома мощный, новые	1650	NEW	Унивеpсaльный паровой очиcтитель – этo эффективная бытовaя тeхникa для убоpки дoмa, coздaнная для удобствa и экoнoмии времeни. Этoт мощный парогенератoр cтaнeт вaшим надежным помощникoм в бoрьбе c зaгpязнeниями на куxне, мебeли и другиx повepхнocтях.\r\n	Мегаполис, 22, улица Володарского, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/965b0669-84ea-48c0-80e3-3e5ce0fe5022.png}	1	10	9851099	2025-12-02 11:35:21.921	2025-12-16 09:53:24.282	61	\N	f	MODERATE	\N
3207807	SWEETPEEPS золотые украшения	7000	NEW	Золотые украшения с фианитами	Уральская улица, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f3a550d2-fc17-4548-943e-b62d66c014eb.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/afa0dc23-c3de-4ab1-8844-a4506bb49309.png}	1	7	6038643	2025-11-28 09:14:30.517	2025-12-16 09:53:24.284	48	\N	f	APPROVED	\N
9256863	Chevrolet Corvette C7	8500000	USED	Корвет был угнан у курседа	Оренбургский Колледж Экономики и Информатики, 11, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/d0c85309-e830-459b-9f87-a672313a465e.jpg}	1	7	8964288	2025-11-28 09:15:31.784	2025-12-16 09:53:24.287	48	\N	f	APPROVED	\N
5609249	ДМРВ на ваз 2107	7000	NEW	Датчик массового расхода воздуха	48, улица Коминтерна, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/05f3e4d2-118c-4262-a12a-8804cddb31d7.webp}	1	10	4761896	2025-11-28 09:16:40.023	2025-12-16 09:53:24.29	61	https://yandex.ru/video/preview/13520813755431483017	f	APPROVED	\N
2388612	Фундук культурный	280	USED	Прoдaю фундук 2024г cбopa, собственный небoльшой cад в предгopьяx Кавказa, бeз xимии тoлькo органикa.\r\n\r\nBcе сopта выращиваемые мной имеют лучшиe вкусoвыe xaрактериcтики и oтносятcя к cтoлoвым сoртам, oбладaют плотным ядpoм и приятным выpaженным маcляниcтым вкуcoм, который не сравним с дешёвыми cетевыми безвкусными орешками.\r\nПредлагаю микс сортов Трапезунд, Анаклиури, Президент.\r\n\r\nВозможна доставка авитодоставкой до 20кг или транспортной компанией от 30.	19/2, улица Бурзянцева, Кузнечный, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a0505026-7df1-43c2-a19c-16e30e07a690.png}	1	9	9851099	2025-12-02 11:27:45.043	2025-12-16 09:53:24.218	54	\N	f	APPROVED	\N
2139014	Пельмени домашние	380	NEW	Пpeдcтaвляeм вaшему вниманию пельмени, манты, хинкaли, ваpеники pучнoй лeпки.	18, Матросский переулок, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/19b43a4f-6ae7-4539-bee3-78a337e8e3c8.png}	1	8	9851099	2025-12-02 11:17:35.121	2025-12-16 09:53:24.254	51	\N	f	APPROVED	\N
4966297	Эублефар	4000	USED	Продаются малыши эублефары различных морф. Едят разморозку, линяют хорошо, все процессы в норме.\r\n	2, Госпитальный переулок, Аренда, Промышленный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/99f1ee8b-a4e7-400c-9591-74ad68a831b6.png}	1	9	9851099	2025-12-02 11:29:46.562	2025-12-16 09:53:24.292	54	\N	f	APPROVED	\N
4081087	Собака	100	USED	Собака овчарка	3/5, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b48c5093-8d65-47e0-84be-e1736ceffbe9.png}	1	9	9371169	2025-11-28 09:10:38.847	2025-12-16 09:53:24.295	54	\N	f	APPROVED	\N
1979749	Игуана	285000	USED	Xoроший спoкoйный пaрень в самoм рaсцвeтe игуаниx cил.\r\n\r\nЗовут Яша, 19 лет, любит тепло и голубику.	26Б, улица Шевченко, Дзержинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/77280a5a-a493-49e1-aeb4-5bb7cbe97653.webp}	1	9	2321239	2025-11-28 09:11:59.098	2025-12-16 09:53:24.299	54	\N	f	APPROVED	\N
5611056	Monster Energy Pipeline Punch	250	NEW	Тонизирующий напиток с изысканным вкусом!	Hawthorne Street, Кламат-Фолс, Klamath County, Орегон, 97601, Соединённые Штаты Америки	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/bfa67311-f947-41c7-b275-e7d28e1db313.jpg}	1	8	1208299	2025-11-28 09:16:21.207	2025-12-16 09:53:24.302	51	https://vk.com/video-129440544_456249335	f	APPROVED	\N
8083712	Концтовары	700	NEW	Набор канцтоваров для школы и офиса Лапки котика 5 предметов	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e97b5f41-387d-42c2-bc90-a342ab3403bb.png}	1	6	6038643	2025-12-02 11:27:09.759	2025-12-16 09:53:24.455	41	\N	f	APPROVED	\N
1961051	Наклейки	700	NEW	Наклейки для ежедневника Школьная эстетика	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/ab80ea7d-50f6-4bcf-beff-ec2a68d97299.png}	1	6	6038643	2025-12-02 11:30:16.141	2025-12-16 09:53:24.457	41	\N	f	APPROVED	\N
8257036	Пенал	590	NEW	Милый эстетичный большой пенал школьный	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5090f83f-e7b2-42af-bc8e-2937850f8952.png}	1	6	6038643	2025-12-02 11:31:59.433	2025-12-16 09:53:24.459	41	\N	f	APPROVED	\N
3824376	Украшения ручной работы	1000	NEW	Украшения ручной работы на заказ по Вашим эскизам/фото. Стоимость украшений на фото 1000р. Срок изготовления: 4-7 дней.\r\n\r\n	46, улица 9 Января, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/4164bbf3-aa1d-4dac-b361-dd22fc5c2001.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/471637e1-9898-41d2-9dc1-c112c642c296.png}	1	7	9851099	2025-12-02 11:01:00.313	2025-12-16 09:53:24.189	48	\N	f	APPROVED	\N
3244052	Дубленка	8000	NEW	Дубленка зимняя	улица Рокоссовского, Горка, Дзержинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b6118a2c-fce6-4c04-826b-b71e8953afe4.png}	1	1	6038643	2025-12-02 11:01:27.168	2025-12-16 09:53:24.376	2	\N	f	APPROVED	\N
8776759	Кофемашина Thomson CF20A02	11399	NEW	Рабoчaя бытoвая тexника, намного дeшевлe, чем в мaгaзинe;\r\n\r\n\r\n\r\nНeт тapы для мoлoкa\r\n\r\n- Любыe пpoверки при caмoвывозе;\r\n\r\n- Пpи пpиемке товара вся теxника пpoвepяется на рaбoтоспocoбнoсть;\r\n\r\n- Oтправляeм Авитo доcтaвкой;\r\n\r\n- Пpи дoставке тoвap упaковываeтся по высшему урoвню.\r\n\r\nBитpинный образец:\r\n\r\n• товар новый, стоял на витрине в магазине;\r\n\r\n• может быть повреждена заводская упаковка;\r\n\r\n• возможны незначительные потёртости или повреждения корпуса, которые никак не влияют на работоспособность.\r\n\r\nЗа фотографиями дефектов обращайтесь в лс\r\n\r\nСамовывоз возможен из 2-х точек: Метро Текстильщики, Метро Шипиловская.\r\n\r\nВ нашем профиле большой ассортимент разнообразной бытовой техники. Советуем заглянуть!\r\n\r\nБольше техники в нашем телеграмм-канале\r\n\r\nПереходите там большие скидки!\r\n\r\n	44, улица Кирова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/673d92ac-2d1c-415f-8beb-02bd64a3b69d.png}	1	10	9851099	2025-12-02 11:34:25.201	2025-12-16 09:53:24.174	61	\N	f	MODERATE	\N
1250840	Кресло-горилла	170000	NEW	Кресло-горилла удобное, выполнено из лучших материалов.	37А, Илекская улица, село имени 9 Января, Красноуральский сельсовет, Оренбургский район, Оренбургская область, Приволжский федеральный округ, 460501, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/228256c7-91a9-4c50-9a2e-a2804917075b.png}	1	12	7391202	2025-11-28 09:16:46.544	2025-12-16 09:53:24.305	67	\N	f	APPROVED	\N
9368305	Майка	20000	USED	очень крутые маечки с аниме принтами, у2к вайб имеется🪽размер S, полиэстер\r\nцена 500 рублей за штуку\r\nпо любым вопросам пишите!!\r\n\r\n	3, улица Аксакова, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/33229854-0ccb-4ecc-b7a8-3ab2965d8fdc.png}	1	1	8633592	2025-11-28 09:18:16.763	2025-12-16 09:53:24.307	1	\N	f	APPROVED	\N
3334788	Посуда для сервировки Estetic	3500	NEW	Вся посуда выполнена в минималистичных стилях, из качественных материалов, подойдет на каждый день	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3b0ddf2a-19a7-4ecf-99c9-9dfc705d35c7.png}	1	11	6038643	2025-11-28 09:11:35.533	2025-12-16 09:53:24.316	63	\N	f	APPROVED	\N
6901799	Кошка	10	USED	Кошка домашняя	5, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/8f6c52ee-023b-41d0-befc-39612d968abf.webp}	1	9	9371169	2025-11-28 09:11:58.971	2025-12-16 09:53:24.319	54	\N	f	APPROVED	\N
4224343	Салонный фильтр на ваз 2110	1000	NEW	салонный фильтр подходит на автомобили ваз2110,2112	20, улица Кобозева, Кузнечный, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/22cb6191-445a-4ee0-80f2-92cc96055093.webp}	1	10	4761896	2025-11-28 09:12:25.319	2025-12-16 09:53:24.326	61	https://yandex.ru/video/preview/9506785745966413491	f	APPROVED	\N
1300264	Ford Mustang	2500000	NEW	Самый лучший автомобиль в мире	Оренбургский Колледж Экономики и Информатики, 11, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c91c9b99-afe9-4d10-8a81-b3cee1c12296.jpg}	1	10	8964288	2025-11-28 09:12:03.576	2025-12-16 09:53:24.322	61	\N	f	APPROVED	\N
1970246	Козы камерунские	3000	NEW	Продаются козочки камерунские,разного возраста, есть два козлика для покрытия, покрытие 3 тыс	"Воздух" конный клуб, 9, Бассейный переулок, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/bfba4d6b-5f64-40b1-a0e2-8f232b9140ea.webp}	1	9	2321239	2025-11-28 09:14:23.82	2025-12-16 09:53:24.33	54	\N	f	APPROVED	\N
9915250	Платье горничной	1200	USED	платье горничной в хорошем состоянии , нету только ободка осталась только от него ткань, если нужно доп фото пишите, к платью идет бантик и фартук\r\n\r\n	г Оренбург	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/dd3f08a3-bdcc-4593-b35f-1e186ce5262a.png}	1	1	8633592	2025-11-28 09:14:31.733	2025-12-16 09:53:24.334	1	\N	f	APPROVED	\N
9863001	Набор золотых украшений	2000	NEW	\N	Лицей №2, Красная улица, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e7670256-5300-4433-82a5-edf31f999776.webp}	1	7	6053931	2025-11-28 09:12:06.874	2025-12-16 09:53:24.324	48	\N	f	APPROVED	\N
7378626	Этно украшения	300	USED	Украшения в этническом стиле! серьги, браслеты, ожерелья, броши и т, д	26, улица Кирова, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fe07794a-1f08-4d05-8569-a90fa9c75a56.png}	1	7	9851099	2025-12-02 10:58:49.981	2025-12-16 09:53:24.336	48	\N	f	APPROVED	\N
5510664	Джинсы	2500	NEW	Джинсы в новом состояние	48, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/935c3f3a-3998-4c35-97aa-d83c3b4c3beb.png}	1	1	6038643	2025-12-02 10:51:54.306	2025-12-16 09:53:24.361	2	\N	f	APPROVED	\N
2207276	Ложка	100	NEW	Просто ложка	г Оренбург, ул Львовская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1582c805-78ef-49e4-a5bc-c875f429af60.webp}	1	11	4146092	2025-12-02 10:54:33.568	2025-12-16 09:53:24.356	63	\N	f	APPROVED	\N
4523969	Платье	2000	NEW	Платье летнее разных расцветок	2, улица Богдана Хмельницкого, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460024, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/497ae48d-8b31-443d-8a6f-9f00be0ac793.png}	1	1	6038643	2025-12-02 10:54:44.305	2025-12-16 09:53:24.353	2	\N	f	APPROVED	\N
3506516	Разобранный кубик рубика	10	USED	не смог собрать	Оренбургский Колледж Экономики и Информатики, 11, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c55064c4-df45-485c-80a3-7253e48ff798.jfif,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/d4b88f23-5430-4948-96f8-b521befb052b.jpg}	1	2	8964288	2025-11-28 09:19:52.442	2025-12-16 09:53:24.339	14	\N	f	APPROVED	\N
6003323	Кресло офисное	5000	NEW	Удобное кресло	г Оренбург, ул Богдана Хмельницкого	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a0763806-d3e2-4699-9674-487603f386a3.jpg}	1	12	4146092	2025-12-02 10:56:54.12	2025-12-16 09:53:24.35	67	\N	f	APPROVED	\N
2105765	Тест	20000	NEW	Описание	Вита Экспресс, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f76a26f7-4801-4c7f-9166-6b2869b5a765.jpg}	1	8	3235109	2025-12-01 05:50:33.37	2025-12-16 09:53:24.341	51	\N	f	APPROVED	\N
1122280	Кресло-коляска	45000	USED	новая	г Оренбург	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/0a43c4d3-222d-421b-8528-7e3e59cc909a.jpg}	1	15	2681599	2025-12-01 08:10:41.21	2025-12-16 09:53:24.342	75	\N	f	APPROVED	\N
5902819	Свитер	3000	NEW	Свитер теплый из мягкой ткани	3/1, Телевизионный переулок, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460024, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/8ee766c3-6a41-4953-881d-c48ff14a1add.png}	1	1	6038643	2025-12-02 10:57:39.627	2025-12-16 09:53:24.366	2	\N	f	APPROVED	\N
4267180	Табурет	500	NEW	Просто табурет.	г Оренбург, ул Полтавская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/2a965649-9cf7-4e34-a427-8926bc88b2c9.jpg}	1	12	4146092	2025-12-02 11:00:14.233	2025-12-16 09:53:24.371	67	\N	f	APPROVED	\N
6752957	Диван	6000	NEW	Просто диван	г Оренбург, ул Гоголя	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/28a67960-b058-4732-8bf1-501b4d4cca5a.webp}	1	12	4146092	2025-12-02 11:01:04.406	2025-12-16 09:53:24.373	67	\N	f	APPROVED	\N
7213485	Кровать	10000	NEW	Удобная кровать. Евродвушка	г Оренбург, Крымский пер	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a83a8772-60ab-46e0-ac5e-3130ef9deb81.webp}	1	12	4146092	2025-12-02 11:02:27.961	2025-12-16 09:53:24.381	67	\N	f	APPROVED	\N
3184247	Украшения в русском стиле	2800	NEW	Украшения в русском стиле из натуральных камней и керамических бусин с подвесками ручной работы: неваляшки, Петушки, лошадки.	Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1c859471-e841-471b-a425-cd312047cd68.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/9114495f-1926-48c5-bcbd-336ef851b323.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3dd7eee0-ee94-440c-84e0-16d5dfb7740e.png}	1	7	9851099	2025-12-02 11:03:22.662	2025-12-16 09:53:24.4	48	\N	f	APPROVED	\N
9380113	Детские игрушки	1000	USED	Набор детский игрушек	24, Луговая улица, Восточный, Сотки, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e8e9b6ba-a396-477d-b8c3-000cc9e85c0f.png}	1	2	6038643	2025-12-02 11:04:58.661	2025-12-16 09:53:24.389	9	\N	f	APPROVED	\N
6497808	Тонометр	2000	NEW	Тонометр. Давление меряет еще что-то там	г Оренбург, ул Донецкая	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/50400626-7e5c-41a0-8b2a-d7d753a627cd.jpg}	1	15	4146092	2025-12-02 11:03:58.743	2025-12-16 09:53:24.384	73	\N	f	APPROVED	\N
9783545	Ванночка	3000	USED	Ванна для купания новорожденного	6Б, Телевизионный переулок, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460024, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/ef7cf831-d301-4470-8e08-45a9662ccc25.png}	1	2	6038643	2025-12-02 11:08:32.807	2025-12-16 09:53:24.392	15	\N	f	APPROVED	\N
8436378	Серебряные украшения	1500	USED	Пoд номeром 1: сеpежки с розoвым камнeм 1000 рублей. Под номeрoм 2: нaбop cepeжки и кольцо с жeлтым кaмнeм 2000 рублей зa нaбор. Под номером 3: набoр cepeжки, кольцо и подвеcкa с зелeным кaмнeм 2000 pублей зa набоp. Под нoмеpoм 5: сepежки с рoзoвым кaмнем 1000 рублeй. Серебрянaя цeпoчка 2000 рублей. Кольцо с белым камнем и две подвески с белыми камнями- по 500 рублей каждая. Серебро все в хорошем состоянии	Фармленд, 52, Советская улица, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/053660d3-2e5f-4f54-a2f8-c11767cd53fc.png}	1	7	9851099	2025-12-02 11:04:37.805	2025-12-16 09:53:24.387	48	\N	f	APPROVED	\N
6300121	Кровать	10000	USED	Кровать детская	139, Ташкентская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5b84facc-aec9-43ff-aa6e-e6a5922f31f2.png}	1	2	6038643	2025-12-02 11:13:18.658	2025-12-16 09:53:24.403	10	\N	f	APPROVED	\N
1885272	Пара Флэт Уайт	363	NEW	Пара Флэт Уайт по выгодной цене. Доступно только в доставке!	30, улица 8 Марта, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/61625d76-10a4-423b-9040-f9871c898a6b.png}	1	8	9851099	2025-12-02 11:14:00.643	2025-12-16 09:53:24.405	51	\N	f	APPROVED	\N
7718497	Катетер	150	NEW	Просто катетер	г Оренбург, ул Севастопольская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5a2c2b79-558b-46e7-b267-6ae194c526b9.jpg}	1	15	4146092	2025-12-02 11:15:13.932	2025-12-16 09:53:24.408	74	\N	f	APPROVED	\N
3217337	Кофеин в таблетках	160	NEW	Просто кофеин	г Оренбург, мкр Ростошинские пруды, Керченский пер	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/845037dc-1299-407d-a608-b0bfccd6de8d.webp}	1	3	4146092	2025-12-02 11:16:08.622	2025-12-16 09:53:24.41	22	\N	f	APPROVED	\N
5492285	Часы	500	NEW	Часы громкоговорители	113, Невельская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b7eb5ecd-4fb4-4d07-ac74-b5146c8080ba.png}	1	5	6038643	2025-12-02 11:16:24.995	2025-12-16 09:53:24.415	35	\N	f	APPROVED	\N
1314227	Средство для удаления тейпов	500	NEW	Средство для удаления тейпов	41, улица Терешковой, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/9092c6b3-18bf-4960-b6c0-ae294784dd18.png}	1	5	6038643	2025-12-02 11:24:04.588	2025-12-16 09:53:24.425	28	\N	f	APPROVED	\N
6883587	Шампунь Гарньер	500	NEW	Просто шампунь	г Оренбург, ул Сумская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/72f30e85-f179-47a0-a982-11bf90113e0e.jpg}	1	3	4146092	2025-12-02 11:29:29.551	2025-12-16 09:53:24.432	17	\N	f	APPROVED	\N
9956819	Тени для век	2000	NEW	Просто тени	г Оренбург, ул Житомирская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3a20e814-cdaf-40a4-8252-fd825161c268.webp}	1	3	4146092	2025-12-02 11:31:50.394	2025-12-16 09:53:24.442	21	\N	f	APPROVED	\N
9500725	Морозилки ларь Бирюса, Pozis, Kraft и другие	15990	NEW	Бoльшой выбoр мopoзильныx камер (вepтикaльныe, лapи) разных oбъёмoв в нaличии в Орeнбуpге!\r\n\r\nА так же в наличии огромный выбoр бытoвoй тexники по оптовым ценaм!\r\n\r\n	Вишнёвая улица, СНТ "ЮЖНЫЙ УРАЛ ОФИЦЕРОВ ЗАПАСА И ОТСТАВКИ", Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/598ec012-3f41-4d02-9721-750964a49125.png}	1	10	9851099	2025-12-02 11:32:20.271	2025-12-16 09:53:24.444	61	\N	f	APPROVED	\N
7162519	Стиральная машина бу	7000	USED	Стиральныe машины б.у. 🚛 Бecплaтная доставкa по гoроду ✅Гарaнтия до 12 меcяцeв пo чeку + пocлeгарантийнoe oбслуживаниe.\r\n\r\n	25, Краснознамённая улица, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/6f7d4e74-5748-4fa3-b1d3-a22f8aa6a061.png}	1	10	9851099	2025-12-02 11:33:16.603	2025-12-16 09:53:24.449	61	\N	f	APPROVED	\N
9042977	Закладки для учебников 	300	NEW	Закладки для книг, «Книжная эстетика»	5, улица Макаровой-Мутновой, Новостройка, Промышленный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a2609c6c-a3c3-4d4a-bd93-60020e455210.png}	1	6	6038643	2025-12-02 11:33:30.034	2025-12-16 09:53:24.452	41	\N	f	APPROVED	\N
7566163	Духи 	3500	NEW	Духи Dior Sauvage	г Оренбург, ул Черниговская	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fa8257b2-6b96-4b79-b10e-df3d7c723129.jpg}	1	3	4146092	2025-12-02 11:27:20.06	2025-12-16 09:53:24.438	20	\N	f	APPROVED	\N
2568373	Концелярия	700	NEW	Канцелярия для школы набор линеек y2k эстетика бант кролик	128, Орская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5e603622-4d00-4003-a685-cca9d4c78cf7.png}	1	6	6038643	2025-12-02 11:29:06.905	2025-12-16 09:53:24.434	41	\N	f	APPROVED	\N
3437684	Минский Бургер с курицей	330	NEW	По-белорусски вкусный! Бургер с сочной куриной котлетой в хрустящей панировке, румяным картофельным оладушком, свежим салатом, двумя ломтиками нежного сыра, хрустящим ароматным беконом, маринованными огурчиками, нежным соусом «Сметана-укроп», и всё это — на воздушной горячей булочке с хрустящей крошкой.	54, улица Кирова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3458ebaf-7b3f-4f39-b1f4-5a53322d9e64.png}	1	8	9851099	2025-12-02 11:12:20.477	2025-12-16 09:53:24.21	51	\N	f	APPROVED	\N
6218446	Соковыжималка caso CP 300 Pro	4500	USED	CASO – нeмeцкaя торговая маркa бытовoй техники, принадлежащaя кoмпaнии Braukmann GmbH. Cоковыжималкa CASО СP 330 Prо предназначена для цитруcовыx cpeднего и крупногo pазмеpoв. Koрпуc прибоpа и cито для жмыxa выполнeны из нeржавeющeй cтaли. Автoмaтичeский старт плавнo зaпускает двигатель мощностью 160 Вт, функция «капля – стоп» обеспечивает чистоту рабочего места. В идеальном состоянии.\r\n\r\n	23/2, Пролетарская улица, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия	{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f2216887-329e-4d48-a605-e9b8bad18686.png}	1	10	9851099	2025-12-02 11:36:09.481	2025-12-16 09:53:24.309	61	\N	f	MODERATE	\N
\.


--
-- Data for Name: ProductFieldValue; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."ProductFieldValue" (id, value, "fieldId", "productId") FROM stdin;
34	50	1	9368305
35	Радужный	2	9368305
36	Хлопок	3	9368305
37	Demix	4	9368305
38	Аниме	5	9368305
39	Да	6	9368305
40	15	54	3506516
41	разный	58	3506516
42	3х3	62	3506516
43	не знаю	67	3506516
48	1-5	54	4215912
49	Книжный	58	4215912
50	Книжный	62	4215912
51	Книжные	67	4215912
52	черный	86	1122280
85	Активная модель 1000	136	4758351
86	, прогулочная	137	4758351
87	Железяка	138	4758351
91	нет	142	4758351
118	s	7	5510664
119	голубой	8	5510664
120	джинса	9	5510664
121	Gloria	10	5510664
122	джинсы	11	5510664
123	багги	12	5510664
151	XS-L	7	4523969
152	разные	8	4523969
153	хлопок	9	4523969
154	Dasha	10	4523969
155	DV	11	4523969
156	летние	12	4523969
157	S-L	7	5902819
158	белый	8	5902819
159	норка	9	5902819
160	red	10	5902819
161	Sweet	11	5902819
162	свитер	12	5902819
163	S-L	7	9122333
164	бежевый	8	9122333
165	хлопок	9	9122333
166	Bant	10	9122333
167	BD	11	9122333
168	зипка	12	9122333
169	s	7	3244052
170	голубо-белый	8	3244052
171	джинса	9	3244052
172	VK	10	3244052
173	vk	11	3244052
174	дубленка	12	3244052
175	Белый	84	6497808
176	Да	87	6497808
177	Omron	90	6497808
178	Да	95	6497808
179	2-3	51	9380113
180	5	59	9380113
181	разные	63	9380113
182	15-30 см	64	9380113
183	Белый	84	3563632
184	Да	87	3563632
185	Omron	90	3563632
186	Да	95	3563632
187	50см	50	2865910
188	грязно-синий	57	2865910
189	1-3	60	2865910
190	50см	70	2865910
191	Да	85	8948419
192	Zenet	88	8948419
193	Нет	92	8948419
194	Синий	94	8948419
195	белый	49	9783545
196	100см	55	9783545
197	0-1	61	9783545
198	200см	73	9783545
199	100см	53	6734788
200	0-1	65	6734788
201	100см	72	6734788
202	белый	75	6734788
203	Да	85	9305563
204	KMED	88	9305563
205	Да	92	9305563
206	Белый	94	9305563
207	200см	53	9265239
208	0-1	65	9265239
209	200см	72	9265239
210	бело-коричневый	75	9265239
211	1,5 метра	53	6300121
212	0-7	65	6300121
213	1,5 метра	72	6300121
214	белый	75	6300121
215	Да	85	7718497
216	KMD	88	7718497
217	Да	92	7718497
218	Белый	94	7718497
219	Таблетки	83	3217337
220	часы	187	5492285
221	америка	188	5492285
222	часы	189	5492285
223	60см	190	5492285
224	Крем	82	8507601
225	Компрессорный ингалятор	187	4372887
226	Omron Comp Air NE-C300 Complete	188	4372887
227	Небулайзер OMRON C300 Complete — прибор, работающий в 3 режимах ингаляции. 	189	4372887
228	70см	190	4372887
229	кресло-коляска	136	6617171
230	сидячий	137	6617171
231	метал	138	6617171
232	7кг	139	6617171
233	100кг	140	6617171
235	нет	142	6617171
240	clinar	132	1314227
241	балончик	133	1314227
242	2 года	134	1314227
243	америка	135	1314227
244	уход	132	2273041
245	спрей	133	2273041
246	5 лет	134	2273041
247	dinax	135	2273041
248	Духи	77	7566163
249	Шампунь	78	6883587
250	Палетка с тенями	80	9956819
251	папва	117	5868178
252	папва	118	5868178
253	апавпва	119	5868178
254	ир	120	5868178
255	рпрпр	121	5868178
256	пнпнп	122	5868178
28	52	1	9915250
29	Черный	2	9915250
30	Шелк	3	9915250
31	Gucci	4	9915250
32	Платье	5	9915250
33	Горничная	6	9915250
88	16 кг 500 грамм	139	4758351
89	120 кгили 0,12 тонны, или 1,2 центнера	140	4758351
90	низкопрофильные	141	4758351
92	ручное	143	4758351
93	мягкая сидушка	144	4758351
94	да	145	4758351
95	черный	146	4758351
234	 Колесная база, не выступающая за габариты коляски	141	6617171
236	автоматическое	143	6617171
237	нет	144	6617171
238	есть	145	6617171
239	черный	146	6617171
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
29	5966833	4081087	2025-11-28 09:19:16.228
23	6038643	4081087	2025-11-28 09:18:17.402
28	8633592	4081087	2025-11-28 09:19:28.543
19	6038643	5231119	2025-11-28 09:17:47.718
106	8964288	5231119	2025-11-28 09:50:04.663
21	6038643	6901799	2025-11-28 09:18:07.082
54	8964288	6901799	2025-11-28 09:21:24.529
20	6038643	1979749	2025-11-28 09:18:00.208
105	8964288	1979749	2025-11-28 09:49:35.399
7	6038643	1300264	2025-11-28 09:18:13.977
602	7132269	1300264	2025-12-01 14:29:41.095
270	2681599	9863001	2025-12-01 08:12:43.181
42	2321239	4224343	2025-11-28 09:20:52.446
504	7106521	4224343	2025-12-02 07:39:03.095
11	6038643	7659684	2025-11-28 09:15:50.575
305	7106521	7659684	2025-12-01 12:03:10.898
6	6038643	1970246	2025-11-28 09:15:46.542
204	7106521	1970246	2025-12-01 09:26:35.713
70	8964288	1970246	2025-11-28 09:22:21.23
238	7106521	3207807	2025-12-01 09:26:48.769
51	8964288	3207807	2025-11-28 09:21:19.648
5	6038643	9915250	2025-11-28 09:15:40.734
35	2321239	9915250	2025-11-28 09:20:26.159
13	9371169	9915250	2025-11-28 09:16:05.942
173	7106521	9915250	2025-12-01 09:26:28.453
12	8964288	9915250	2025-11-28 09:16:05.662
15	5966833	9256863	2025-11-28 09:16:48.332
8	6038643	9256863	2025-11-28 09:15:34.351
40	2321239	9256863	2025-11-28 09:20:48.02
569	7106521	9256863	2025-12-01 09:26:25.118
14	5966833	5611056	2025-11-28 09:16:42.097
18	6038643	5611056	2025-11-28 09:17:21.48
239	7106521	5611056	2025-12-01 07:48:19.685
43	8964288	5611056	2025-11-28 09:20:59.678
17	6038643	5609249	2025-11-28 09:17:15.39
39	2321239	5609249	2025-11-28 09:20:43.139
49	8633592	5609249	2025-11-28 09:21:16.127
536	7106521	5609249	2025-12-01 09:24:23.12
16	6038643	1250840	2025-11-28 09:17:09.157
31	4761896	1250840	2025-11-28 09:19:39.18
48	8633592	1250840	2025-11-28 09:21:12.953
172	7106521	1250840	2025-12-03 16:57:03.686
46	8964288	1250840	2025-11-28 09:21:09.011
25	6038643	9262881	2025-11-28 09:18:35.85
44	8633592	9262881	2025-11-28 09:21:04.291
50	7106521	9262881	2025-12-02 07:46:38.367
24	6038643	9368305	2025-11-28 09:18:26.145
27	2321239	9368305	2025-11-28 09:20:38.981
26	9371169	9368305	2025-11-28 09:18:54.527
206	7106521	9368305	2025-12-01 09:26:15.811
47	8633592	7384341	2025-11-28 09:21:09.232
72	7106521	7384341	2025-12-03 16:55:51.09
605	7132269	7384341	2025-12-01 14:46:32.887
32	6038643	3506516	2025-11-28 09:20:15.947
45	8633592	3506516	2025-11-28 09:21:06.787
33	7106521	3506516	2025-11-28 09:20:19.239
604	7132269	3506516	2025-12-01 14:37:03.378
52	8633592	4215912	2025-11-28 09:21:21.713
237	7106521	4215912	2025-12-01 09:24:12.715
71	8964288	4215912	2025-11-28 09:22:26.679
606	7132269	4215912	2025-12-01 14:54:23.873
205	7106521	2105765	2025-12-02 06:32:52.922
603	7132269	2105765	2025-12-01 14:29:43.789
336	6669460	2105765	2025-12-01 08:37:06.268
369	7106521	1122280	2025-12-01 08:48:56.032
468	7106521	2161612	2025-12-01 09:22:41.051
618	9851099	2161612	2025-12-02 11:42:27.288
567	7106521	4758351	2025-12-02 07:32:36.871
656	7106521	4267180	2025-12-03 09:52:39.77
660	7106521	9380113	2025-12-03 16:41:06.056
652	7106521	9783545	2025-12-03 09:51:25.113
654	7106521	3437684	2025-12-03 09:52:30.64
617	4146092	3437684	2025-12-02 11:33:23.238
662	7106521	1512888	2025-12-03 16:54:39.706
653	7106521	6628130	2025-12-03 16:38:43.192
659	7106521	2388612	2025-12-03 17:01:43.056
651	7106521	9956819	2025-12-03 09:51:19.021
658	7106521	6218446	2025-12-03 16:38:47.965
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

COPY public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") FROM stdin;
1	Одежда	1	2025-12-15 19:18:08.513	odezhda	2025-12-15 17:21:30.498
2	Детские товары	1	2025-12-15 19:18:08.513	detskie-tovary	2025-12-15 17:21:30.503
5	Средства реабилитации	1	2025-12-15 19:18:08.513	sredstva-reabilitatsii	2025-12-15 17:21:30.508
6	Школьные товары	1	2025-12-15 19:18:08.513	shkol-nye-tovary	2025-12-15 17:21:30.513
7	Украшения	1	2025-12-15 19:18:08.513	ukrasheniya	2025-12-15 17:21:30.518
8	Продукты питания	1	2025-12-15 19:18:08.513	produkty-pitaniya	2025-12-15 17:21:30.522
9	Животные, растения	1	2025-12-15 19:18:08.513	zhivotnye-rasteniya	2025-12-15 17:21:30.527
10	Бытовая техника	1	2025-12-15 19:18:08.513	bytovaya-tehnika	2025-12-15 17:21:30.531
11	Посуда	1	2025-12-15 19:18:08.513	posuda	2025-12-15 17:21:30.536
12	Мебель	1	2025-12-15 19:18:08.513	mebel	2025-12-15 17:21:30.54
15	Медицинские товары	1	2025-12-15 19:18:08.513	meditsinskie-tovary	2025-12-15 17:21:30.544
3	Красота и здоровье	1	2025-12-15 19:18:08.513	krasota-i-zdorov-e	2025-12-15 17:21:30.548
\.


--
-- Data for Name: SubcategotyType; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") FROM stdin;
1	Мужская	1	2025-12-15 19:18:08.514	muzhskaya	2025-12-15 17:21:30.561
2	Женская	1	2025-12-15 19:18:08.514	zhenskaya	2025-12-15 17:21:30.565
3	Детская девочки	1	2025-12-15 19:18:08.514	detskaya-devochki	2025-12-15 17:21:30.57
4	Детская мальчики	1	2025-12-15 19:18:08.514	detskaya-mal-chiki	2025-12-15 17:21:30.574
5	Ткани, текстиль и фурнитура	1	2025-12-15 19:18:08.514	tkani-tekstil-i-furnitura	2025-12-15 17:21:30.578
6	Сумки, рюкзаки	1	2025-12-15 19:18:08.514	sumki-ryukzaki	2025-12-15 17:21:30.582
7	Аксессуары	1	2025-12-15 19:18:08.514	aksessuary	2025-12-15 17:21:30.585
8	Обувь	1	2025-12-15 19:18:08.514	obuv	2025-12-15 17:21:30.59
9	Игрушки	2	2025-12-15 19:18:08.514	igrushki	2025-12-15 17:21:30.595
10	Детская мебель	2	2025-12-15 19:18:08.514	detskaya-mebel	2025-12-15 17:21:30.599
11	Коляски детские	2	2025-12-15 19:18:08.514	kolyaski-detskie	2025-12-15 17:21:30.603
12	Велосипеды и самокаты	2	2025-12-15 19:18:08.514	velosipedy-i-samokaty	2025-12-15 17:21:30.608
13	Детское питание и посуда	2	2025-12-15 19:18:08.514	detskoe-pitanie-i-posuda	2025-12-15 17:21:30.612
14	Образовательные товары	2	2025-12-15 19:18:08.514	obrazovatel-nye-tovary	2025-12-15 17:21:30.616
15	Уход и гигиена	2	2025-12-15 19:18:08.514	uhod-i-gigiena	2025-12-15 17:21:30.62
16	Косметика для ухода за кожей	3	2025-12-15 19:18:08.514	kosmetika-dlya-uhoda-za-kozhey	2025-12-15 17:21:30.624
17	Средства для ухода за волосами	3	2025-12-15 19:18:08.514	sredstva-dlya-uhoda-za-volosami	2025-12-15 17:21:30.629
18	Уход и гигиена	3	2025-12-15 19:18:08.514	uhod-i-gigiena	2025-12-15 17:21:30.633
19	Приборы и аксессуары	3	2025-12-15 19:18:08.514	pribory-i-aksessuary	2025-12-15 17:21:30.636
20	Парфюмерия	3	2025-12-15 19:18:08.514	parfyumeriya	2025-12-15 17:21:30.64
21	Макияж	3	2025-12-15 19:18:08.514	makiyazh	2025-12-15 17:21:30.645
22	Бады	3	2025-12-15 19:18:08.514	bady	2025-12-15 17:21:30.648
26	Измерительные приборы	5	2025-12-15 19:18:08.514	izmeritel-nye-pribory	2025-12-15 17:21:30.653
27	Ортопедия (бандажи, корсеты)	5	2025-12-15 19:18:08.514	ortopediya-bandazhi-korsety	2025-12-15 17:21:30.655
28	Уходовая косметика	5	2025-12-15 19:18:08.514	uhodovaya-kosmetika	2025-12-15 17:21:30.659
29	Кресла-коляски	5	2025-12-15 19:18:08.514	kresla-kolyaski	2025-12-15 17:21:30.662
30	Спецодежда, трикотаж, компрессионное белье	5	2025-12-15 19:18:08.514	spetsodezhda-trikotazh-kompressionnoe-bel-e	2025-12-15 17:21:30.666
31	Подгузники, пеленки, прокладки	5	2025-12-15 19:18:08.514	podguzniki-pelenki-prokladki	2025-12-15 17:21:30.669
32	Катетеры	5	2025-12-15 19:18:08.514	katetery	2025-12-15 17:21:30.673
33	Средства ухода за стомой	5	2025-12-15 19:18:08.514	sredstva-uhoda-za-stomoy	2025-12-15 17:21:30.675
34	Кресла-стулья санитарные	5	2025-12-15 19:18:08.514	kresla-stul-ya-sanitarnye	2025-12-15 17:21:30.678
35	Специальные устройства	5	2025-12-15 19:18:08.514	spetsial-nye-ustroystva	2025-12-15 17:21:30.682
36	Калоприемники, уроприемники	5	2025-12-15 19:18:08.514	kalopriemniki-uropriemniki	2025-12-15 17:21:30.686
37	Трости, костыли	5	2025-12-15 19:18:08.514	trosti-kostyli	2025-12-15 17:21:30.69
38	Вертикализаторы, опоры	5	2025-12-15 19:18:08.514	vertikalizatory-opory	2025-12-15 17:21:30.694
39	Матрасы	5	2025-12-15 19:18:08.514	matrasy	2025-12-15 17:21:30.696
40	Кровати медицинские	5	2025-12-15 19:18:08.514	krovati-meditsinskie	2025-12-15 17:21:30.701
41	Письменные принадлежности	6	2025-12-15 19:18:08.514	pis-mennye-prinadlezhnosti	2025-12-15 17:21:30.704
42	Бумажная продукция	6	2025-12-15 19:18:08.514	bumazhnaya-produktsiya	2025-12-15 17:21:30.708
43	Принадлежности для рисования и творчества	6	2025-12-15 19:18:08.514	prinadlezhnosti-dlya-risovaniya-i-tvorchestva	2025-12-15 17:21:30.71
44	Органайзеры и хранение	6	2025-12-15 19:18:08.514	organayzery-i-hranenie	2025-12-15 17:21:30.714
45	Учебные пособия и инструменты	6	2025-12-15 19:18:08.514	uchebnye-posobiya-i-instrumenty	2025-12-15 17:21:30.717
46	Рюкзаки и сумки	6	2025-12-15 19:18:08.514	ryukzaki-i-sumki	2025-12-15 17:21:30.721
47	Прочее	6	2025-12-15 19:18:08.514	prochee	2025-12-15 17:21:30.723
48	Ювелирные изделия	7	2025-12-15 19:18:08.514	yuvelirnye-izdeliya	2025-12-15 17:21:30.727
49	Бижутерия	7	2025-12-15 19:18:08.514	bizhuteriya	2025-12-15 17:21:30.73
50	Часы	7	2025-12-15 19:18:08.514	chasy	2025-12-15 17:21:30.733
51	Готовые продукты	8	2025-12-15 19:18:08.514	gotovye-produkty	2025-12-15 17:21:30.737
52	Напитки	8	2025-12-15 19:18:08.514	napitki	2025-12-15 17:21:30.74
53	Заморозки, полуфабрикаты	8	2025-12-15 19:18:08.514	zamorozki-polufabrikaty	2025-12-15 17:21:30.744
54	Домашние животные	9	2025-12-15 19:18:08.514	domashnie-zhivotnye	2025-12-15 17:21:30.748
55	С/х животные	9	2025-12-15 19:18:08.514	s-h-zhivotnye	2025-12-15 17:21:30.751
56	Рептилии	9	2025-12-15 19:18:08.514	reptilii	2025-12-15 17:21:30.756
57	Растения комнатные	9	2025-12-15 19:18:08.514	rasteniya-komnatnye	2025-12-15 17:21:30.759
58	Культурные растения	9	2025-12-15 19:18:08.514	kul-turnye-rasteniya	2025-12-15 17:21:30.764
59	Декоративные уличные растения	9	2025-12-15 19:18:08.514	dekorativnye-ulichnye-rasteniya	2025-12-15 17:21:30.768
61	Кухонная	10	2025-12-15 19:18:08.514	kuhonnaya	2025-12-15 17:21:30.777
62	Бытовая	10	2025-12-15 19:18:08.514	bytovaya	2025-12-15 17:21:30.779
63	Для приготовления пищи	11	2025-12-15 19:18:08.514	dlya-prigotovleniya-pischi	2025-12-15 17:21:30.783
64	Для хранения	11	2025-12-15 19:18:08.514	dlya-hraneniya	2025-12-15 17:21:30.786
65	Для сервировки	11	2025-12-15 19:18:08.514	dlya-servirovki	2025-12-15 17:21:30.79
66	Для приёма пищи	11	2025-12-15 19:18:08.514	dlya-priema-pischi	2025-12-15 17:21:30.793
67	Мягкая мебель	12	2025-12-15 19:18:08.514	myagkaya-mebel	2025-12-15 17:21:30.798
68	Корпусная мебель	12	2025-12-15 19:18:08.514	korpusnaya-mebel	2025-12-15 17:21:30.801
69	Мебель для кухни	12	2025-12-15 19:18:08.514	mebel-dlya-kuhni	2025-12-15 17:21:30.805
70	Мебель для спальни	12	2025-12-15 19:18:08.514	mebel-dlya-spal-ni	2025-12-15 17:21:30.808
71	Садовая мебель	12	2025-12-15 19:18:08.514	sadovaya-mebel	2025-12-15 17:21:30.813
72	Офисная мебель	12	2025-12-15 19:18:08.514	ofisnaya-mebel	2025-12-15 17:21:30.815
74	Оборудование для клиник	15	2025-12-15 19:18:08.514	oborudovanie-dlya-klinik	2025-12-15 17:21:30.821
75	Медицинская мебель	15	2025-12-15 19:18:08.514	meditsinskaya-mebel	2025-12-15 17:21:30.825
60	Доп товары (горшки, грунт, кормилки, поилки, средства по уходу за растениями, инструменты, корма, игрушки, клетки, аксессуары)	9	2025-12-15 19:18:08.514	dop-tovary-gorshki-grunt-kormilki-poilki-sredstva-po-uhodu-za-rasteniyami-instrumenty-korma-igrushki-kletki-aksessuary	2025-12-15 17:21:30.772
73	Диагностическое оборудование	15	2025-12-15 19:18:08.514	diagnosticheskoe-oborudovanie	2025-12-15 17:21:30.819
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

COPY public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") FROM stdin;
3432589	Исаев Максим Андреевич	sima.isaev2305@mail.ru	+79501859919	$2b$10$VI6Gb9KuiHWEnbndcyi1WemTTQgKWwVhpcOfnEEj7W18T8Gw.TPou	INDIVIDUAL	2025-11-28 09:06:55.938	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
3235109	Арзамасцев Даниил	arzamastsevdaniil@gmail.com	+79068346355	$2b$10$NvJVMH9Kn16C7hSuCtRAf./yj8/jgaeUg2ZI0IAkxt2Tc/Cf5DR8G	INDIVIDUAL	2025-12-01 05:48:10.726	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
9371169	Захаров АР ВЛ	Zahar83s@mail.ru	+79878600551	$2b$10$TfLU49EmrMYrTPd46fQv6.QNkD3tEE2WnHVmy8qIdYzHVOX4PLe4q	INDIVIDUAL	2025-11-28 09:07:21.428	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
4761896	Гатин Ян Талгатович	ggg2107@gmail.com	+79228386030	$2b$10$aUbIJdrSn4qPvErIPV8E6uo162lESkmE7orVVIrS/2v8/k8qUQjvm	INDIVIDUAL	2025-11-28 09:08:47.126	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
8633592	Махар Святой Рог	vmahauri029@gmail.com	+79123557497	$2b$10$UbWFDK5KoI92FFzmWZw.s.jslpRNGreNJFQi30q4ZWI9lB02sqegS	INDIVIDUAL	2025-11-28 09:07:05.955	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
6251884	Попов Матвей Иванович	trrina04@mail.ru	+79878993845	$2b$10$cfHgsH42YXRqYPpoZbbhAuFK4bg.81DSzN4JNMGmkLffNma7mLmB.	INDIVIDUAL	2025-12-03 19:26:12.827	2025-12-08 12:30:51.217	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
8964288	Макаров Николай	bapenick445@gmail.com	+79225387481	$2b$10$DHSa1l.0cj7MK.b7ATupL.f7yXnjfGBUEr7Wezf1wul9x2z2eOIkO	INDIVIDUAL	2025-11-28 09:07:33.445	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
6053931	Голосняк Юлия Викторовна	juliagolosnyak@mail.ru	+79328538922	$2b$10$9VP3OmZRjdumTgAJWCBGGe5ozGVZG0Z/okvuWwUdx1wxmJG7brTES	INDIVIDUAL	2025-11-28 09:07:19.394	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
8261539	Подрядов Екатерина Сергеевна	podradovakata91@gmail.com	+79083234725	$2b$10$sdWaXECQtpyEqc61gS4MrOlsoz4nsjYb1gGC1xD2VVFgr/pUqwB3m	INDIVIDUAL	2025-11-28 09:07:29.962	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
5966833	Каверина Мария	kunafina_ruslana7@mail.ru	+79228362555	$2b$10$AY/2V0DgPQ1.ZorhEmTMfOb4o8hq1EkOR9qkHx4/RgG7Cq6OFAOo2	INDIVIDUAL	2025-11-28 09:07:42.429	2025-12-08 12:30:51.217	\N	f	1	t	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
2287442	Абвгдеивич Егор Константинович	barabulkabarabulka@gmail.com	+72280303111	$2b$10$PPEwZxCaLahLuE4XtqI2k.UxgqrcfBgCoXBHT1EUoq86kYraokwz2	INDIVIDUAL	2025-11-28 09:08:14.573	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
6038643	дар	bdi-2006@mail.ru	+79123400130	$2b$10$TROWXU059pwS6Q98JIfGDOL1kzA0oohdraWoB3ZxpEgGqEU//.qQ6	INDIVIDUAL	2025-11-28 09:06:52.861	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
7391202	kostyukov	geronimoprofitop@gmail.com	+79228744883	$2b$10$ulXOXoQl7aAYjf7uJ2opGOApWYjLTVFSWBrWyYAjJp80HAeDl97OS	INDIVIDUAL	2025-11-28 09:07:57.477	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
6669460	Афонасьев Афиларет Михайлович	pr.actual@mail.ru	+79082734009	$2b$10$R0pbgCnq1AVwe9phmKu1GOT0emg48XzDbtYRBEn/xEyCFd8aNYX7y	INDIVIDUAL	2025-12-01 08:28:35.989	2025-12-08 12:30:51.217	\N	f	1	t	https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/users/71116356-ea56-4dd5-ac1a-86c5a6e2e11b.jpg	f	0	0	f	12	2025-12-24 18:33:10.973	0
1208299	Кокеев Фирилл Батькович	test@test.com	+79953501391	$2b$10$0GEA/Uvq4NrHTLuOetQTXuoviQG19DrdEX4NIFUwD.54aF7ePJveO	INDIVIDUAL	2025-11-28 09:07:44.576	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
7106521	Попов Матвей Иванович	vitaly.sadikov1@yandex.ru	+79510341677	$2b$10$05FMyE494pfJScN9OF98COs6yLacnIIE2gueMbTS8s1/PNzaYrA6C	INDIVIDUAL	2025-11-06 19:33:46.625	2025-12-19 11:36:27.898	\N	f	3	f	https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/users/eac42b51-e66a-4d76-bad2-c6db0efd947b.jpg	t	0	2500	f	12	2025-12-24 18:33:10.973	0
2321239	Прокофьева Валерия Денисовна	lin.ferr@mail.ru	+79225406669	$2b$10$7mnxrJ2LJ0S5RoBoo8gVteXYR.o2kM/nnm07SpxHT37YZqEghfVAC	INDIVIDUAL	2025-11-28 09:08:42.207	2025-12-19 12:04:30.85	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
2681599	Корякина Ирина	ikoryakina47@gmail.com	+79228579009	$2b$10$48dtDNK6DIH0yBgup4eqeeG8k5NPkHuhqBNvQ2yCJqayB3sNthYOS	INDIVIDUAL	2025-12-01 08:08:29.883	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
7132269	йцукенгшщзх	qwertyui123@gmail.com	+75678903456	$2b$10$hhmWdTv8RdWeJ1ofHOjaTuKBgOo2JUky9za7NTJ.uCcfrH3W2CK/S	INDIVIDUAL	2025-12-01 14:29:11.538	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
4146092	Фокеев Кирилл	test1@test.com	+71234567890	$2b$10$FELoBjJj0J8IeMy2YhKlIeniLkjz86fijJS2HOFJ3XvJ3fnIulg2i	INDIVIDUAL	2025-12-02 10:48:41.186	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
9851099	Черешков Данила Алексеевич	chereshkov.da2006@gmail.com	+79123431910	$2b$10$hvt0jXBTO6PcqEzKYDKYUO7hivY2kCsC/7Bzwix242L8YDeP6UgnW	INDIVIDUAL	2025-12-02 10:47:25.87	2025-12-08 12:30:43.354	\N	f	1	f	\N	f	0	0	f	12	2025-12-24 18:33:10.973	0
\.


--
-- Data for Name: _UserFavorites; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."_UserFavorites" ("A", "B") FROM stdin;
28	86
28	119
23	4081087
127	4081087
94	4081087
21	4081087
4163503	5231119
3235109	1970246
6053931	1970246
4163503	1970246
282	1250840
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
-- Name: Banner_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Banner_id_seq"', 1, false);


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
-- Name: Payment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Payment_id_seq"', 5, true);


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
-- Name: Banner Banner_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Banner"
    ADD CONSTRAINT "Banner_pkey" PRIMARY KEY (id);


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
-- Name: Payment Payment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Payment"
    ADD CONSTRAINT "Payment_pkey" PRIMARY KEY (id);


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
-- Name: Category_slug_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Category_slug_idx" ON public."Category" USING btree (slug);


--
-- Name: Category_slug_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Category_slug_key" ON public."Category" USING btree (slug);


--
-- Name: Chat_buyerId_sellerId_productId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Chat_buyerId_sellerId_productId_key" ON public."Chat" USING btree ("buyerId", "sellerId", "productId");


--
-- Name: FavoriteAction_userId_productId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "FavoriteAction_userId_productId_key" ON public."FavoriteAction" USING btree ("userId", "productId");


--
-- Name: Payment_orderId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Payment_orderId_key" ON public."Payment" USING btree ("orderId");


--
-- Name: Payment_paymentId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Payment_paymentId_key" ON public."Payment" USING btree ("paymentId");


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
-- Name: SubCategory_slug_categoryId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "SubCategory_slug_categoryId_key" ON public."SubCategory" USING btree (slug, "categoryId");


--
-- Name: SubCategory_slug_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "SubCategory_slug_idx" ON public."SubCategory" USING btree (slug);


--
-- Name: SubcategotyType_slug_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "SubcategotyType_slug_idx" ON public."SubcategotyType" USING btree (slug);


--
-- Name: SubcategotyType_slug_subcategoryId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "SubcategotyType_slug_subcategoryId_key" ON public."SubcategotyType" USING btree (slug, "subcategoryId");


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
-- Name: Message Message_relatedProductId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_relatedProductId_fkey" FOREIGN KEY ("relatedProductId") REFERENCES public."Product"(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: Message Message_senderId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_senderId_fkey" FOREIGN KEY ("senderId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Payment Payment_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Payment"
    ADD CONSTRAINT "Payment_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


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


