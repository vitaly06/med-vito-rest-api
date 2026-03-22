--
-- PostgreSQL database dump
--

\restrict i4fGPb0eOhTnTjn0EFtAKuaz3Fw7sGV1kUneg1OivaqplrSVolBkovyySBXKgMD

-- Dumped from database version 17.6
-- Dumped by pg_dump version 17.6

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
-- Name: BannerModerate; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."BannerModerate" AS ENUM (
    'MODERATE',
    'APPROVED',
    'DENIDED'
);


ALTER TYPE public."BannerModerate" OWNER TO postgres;

--
-- Name: BannerPlace; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."BannerPlace" AS ENUM (
    'PRODUCT_FEED',
    'PROFILE',
    'FAVORITES',
    'CHATS'
);


ALTER TYPE public."BannerPlace" OWNER TO postgres;

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
-- Name: ReviewModerate; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."ReviewModerate" AS ENUM (
    'MODERATE',
    'APPROVED',
    'DENIDED'
);


ALTER TYPE public."ReviewModerate" OWNER TO postgres;

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
    "updatedAt" timestamp(3) without time zone NOT NULL,
    place public."BannerPlace" NOT NULL,
    "navigateToUrl" text NOT NULL,
    name text NOT NULL,
    "userId" integer NOT NULL,
    "moderateState" public."BannerModerate" DEFAULT 'MODERATE'::public."BannerModerate" NOT NULL
);


ALTER TABLE public."Banner" OWNER TO postgres;

--
-- Name: BannerView; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."BannerView" (
    id integer NOT NULL,
    "bannerId" integer NOT NULL,
    "userId" integer,
    "ipAddress" text,
    "viewedAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."BannerView" OWNER TO postgres;

--
-- Name: BannerView_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."BannerView_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."BannerView_id_seq" OWNER TO postgres;

--
-- Name: BannerView_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."BannerView_id_seq" OWNED BY public."BannerView".id;


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
-- Name: Log; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Log" (
    id integer NOT NULL,
    "userId" integer NOT NULL,
    action text NOT NULL
);


ALTER TABLE public."Log" OWNER TO postgres;

--
-- Name: Log_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Log_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Log_id_seq" OWNER TO postgres;

--
-- Name: Log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Log_id_seq" OWNED BY public."Log".id;


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
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "moderateState" public."ReviewModerate" DEFAULT 'MODERATE'::public."ReviewModerate" NOT NULL
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
    "freeAdsLimit" integer DEFAULT 6 NOT NULL,
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
-- Name: BannerView id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."BannerView" ALTER COLUMN id SET DEFAULT nextval('public."BannerView_id_seq"'::regclass);


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
-- Name: Log id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Log" ALTER COLUMN id SET DEFAULT nextval('public."Log_id_seq"'::regclass);


--
-- Name: Message id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message" ALTER COLUMN id SET DEFAULT nextval('public."Message_id_seq"'::regclass);


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
-- Name: BannerView BannerView_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."BannerView"
    ADD CONSTRAINT "BannerView_pkey" PRIMARY KEY (id);


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
-- Name: Log Log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Log"
    ADD CONSTRAINT "Log_pkey" PRIMARY KEY (id);


--
-- Name: Message Message_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_pkey" PRIMARY KEY (id);


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
-- Name: BannerView_bannerId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "BannerView_bannerId_idx" ON public."BannerView" USING btree ("bannerId");


--
-- Name: BannerView_userId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "BannerView_userId_idx" ON public."BannerView" USING btree ("userId");


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
-- Name: BannerView BannerView_bannerId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."BannerView"
    ADD CONSTRAINT "BannerView_bannerId_fkey" FOREIGN KEY ("bannerId") REFERENCES public."Banner"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: BannerView BannerView_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."BannerView"
    ADD CONSTRAINT "BannerView_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: Banner Banner_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Banner"
    ADD CONSTRAINT "Banner_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


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
-- Name: Log Log_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Log"
    ADD CONSTRAINT "Log_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


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

\unrestrict i4fGPb0eOhTnTjn0EFtAKuaz3Fw7sGV1kUneg1OivaqplrSVolBkovyySBXKgMD

