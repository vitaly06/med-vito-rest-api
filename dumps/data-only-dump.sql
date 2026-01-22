--
-- PostgreSQL database dump
--

\restrict 4NiEYDyKlxvlgGqkeb5HvaMLHy97mCg56PzipZ49km854s9NGzFsWS5JEI3INTO

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
-- Data for Name: Role; Type: TABLE DATA; Schema: public; Owner: -
--

SET SESSION AUTHORIZATION DEFAULT;

ALTER TABLE public."Role" DISABLE TRIGGER ALL;

INSERT INTO public."Role" (id, name) VALUES (1, 'default');
INSERT INTO public."Role" (id, name) VALUES (2, 'moderator');
INSERT INTO public."Role" (id, name) VALUES (3, 'admin');


ALTER TABLE public."Role" ENABLE TRIGGER ALL;

--
-- Data for Name: User; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."User" DISABLE TRIGGER ALL;

INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (7106521, 'Попов Матвей Иванович', 'vitaly.sadikov1@yandex.ru', '+79510341677', '$2b$10$05FMyE494pfJScN9OF98COs6yLacnIIE2gueMbTS8s1/PNzaYrA6C', 'INDIVIDUAL', '2025-11-06 19:33:46.625', '2026-01-22 06:14:35.664', NULL, false, 3, false, 'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/users/eac42b51-e66a-4d76-bad2-c6db0efd947b.jpg', true, 0, 1800, false, 12, '2026-01-11 17:11:38.151', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (3432589, 'Исаев Максим Андреевич', 'sima.isaev2305@mail.ru', '+79501859919', '$2b$10$VI6Gb9KuiHWEnbndcyi1WemTTQgKWwVhpcOfnEEj7W18T8Gw.TPou', 'INDIVIDUAL', '2025-11-28 09:06:55.938', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (3235109, 'Арзамасцев Даниил', 'arzamastsevdaniil@gmail.com', '+79068346355', '$2b$10$NvJVMH9Kn16C7hSuCtRAf./yj8/jgaeUg2ZI0IAkxt2Tc/Cf5DR8G', 'INDIVIDUAL', '2025-12-01 05:48:10.726', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (9371169, 'Захаров АР ВЛ', 'Zahar83s@mail.ru', '+79878600551', '$2b$10$TfLU49EmrMYrTPd46fQv6.QNkD3tEE2WnHVmy8qIdYzHVOX4PLe4q', 'INDIVIDUAL', '2025-11-28 09:07:21.428', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (4761896, 'Гатин Ян Талгатович', 'ggg2107@gmail.com', '+79228386030', '$2b$10$aUbIJdrSn4qPvErIPV8E6uo162lESkmE7orVVIrS/2v8/k8qUQjvm', 'INDIVIDUAL', '2025-11-28 09:08:47.126', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (8633592, 'Махар Святой Рог', 'vmahauri029@gmail.com', '+79123557497', '$2b$10$UbWFDK5KoI92FFzmWZw.s.jslpRNGreNJFQi30q4ZWI9lB02sqegS', 'INDIVIDUAL', '2025-11-28 09:07:05.955', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (6251884, 'Попов Матвей Иванович', 'trrina04@mail.ru', '+79878993845', '$2b$10$cfHgsH42YXRqYPpoZbbhAuFK4bg.81DSzN4JNMGmkLffNma7mLmB.', 'INDIVIDUAL', '2025-12-03 19:26:12.827', '2025-12-08 12:30:51.217', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (8964288, 'Макаров Николай', 'bapenick445@gmail.com', '+79225387481', '$2b$10$DHSa1l.0cj7MK.b7ATupL.f7yXnjfGBUEr7Wezf1wul9x2z2eOIkO', 'INDIVIDUAL', '2025-11-28 09:07:33.445', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (6053931, 'Голосняк Юлия Викторовна', 'juliagolosnyak@mail.ru', '+79328538922', '$2b$10$9VP3OmZRjdumTgAJWCBGGe5ozGVZG0Z/okvuWwUdx1wxmJG7brTES', 'INDIVIDUAL', '2025-11-28 09:07:19.394', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (8261539, 'Подрядов Екатерина Сергеевна', 'podradovakata91@gmail.com', '+79083234725', '$2b$10$sdWaXECQtpyEqc61gS4MrOlsoz4nsjYb1gGC1xD2VVFgr/pUqwB3m', 'INDIVIDUAL', '2025-11-28 09:07:29.962', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (5966833, 'Каверина Мария', 'kunafina_ruslana7@mail.ru', '+79228362555', '$2b$10$AY/2V0DgPQ1.ZorhEmTMfOb4o8hq1EkOR9qkHx4/RgG7Cq6OFAOo2', 'INDIVIDUAL', '2025-11-28 09:07:42.429', '2025-12-08 12:30:51.217', NULL, false, 1, true, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (2287442, 'Абвгдеивич Егор Константинович', 'barabulkabarabulka@gmail.com', '+72280303111', '$2b$10$PPEwZxCaLahLuE4XtqI2k.UxgqrcfBgCoXBHT1EUoq86kYraokwz2', 'INDIVIDUAL', '2025-11-28 09:08:14.573', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (6038643, 'дар', 'bdi-2006@mail.ru', '+79123400130', '$2b$10$TROWXU059pwS6Q98JIfGDOL1kzA0oohdraWoB3ZxpEgGqEU//.qQ6', 'INDIVIDUAL', '2025-11-28 09:06:52.861', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (7391202, 'kostyukov', 'geronimoprofitop@gmail.com', '+79228744883', '$2b$10$ulXOXoQl7aAYjf7uJ2opGOApWYjLTVFSWBrWyYAjJp80HAeDl97OS', 'INDIVIDUAL', '2025-11-28 09:07:57.477', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (6669460, 'Афонасьев Афиларет Михайлович', 'pr.actual@mail.ru', '+79082734009', '$2b$10$R0pbgCnq1AVwe9phmKu1GOT0emg48XzDbtYRBEn/xEyCFd8aNYX7y', 'INDIVIDUAL', '2025-12-01 08:28:35.989', '2025-12-08 12:30:51.217', NULL, false, 1, true, 'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/users/71116356-ea56-4dd5-ac1a-86c5a6e2e11b.jpg', false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (1208299, 'Кокеев Фирилл Батькович', 'test@test.com', '+79953501391', '$2b$10$0GEA/Uvq4NrHTLuOetQTXuoviQG19DrdEX4NIFUwD.54aF7ePJveO', 'INDIVIDUAL', '2025-11-28 09:07:44.576', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (2321239, 'Прокофьева Валерия Денисовна', 'lin.ferr@mail.ru', '+79225406669', '$2b$10$7mnxrJ2LJ0S5RoBoo8gVteXYR.o2kM/nnm07SpxHT37YZqEghfVAC', 'INDIVIDUAL', '2025-11-28 09:08:42.207', '2025-12-19 12:04:30.85', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (2681599, 'Корякина Ирина', 'ikoryakina47@gmail.com', '+79228579009', '$2b$10$48dtDNK6DIH0yBgup4eqeeG8k5NPkHuhqBNvQ2yCJqayB3sNthYOS', 'INDIVIDUAL', '2025-12-01 08:08:29.883', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (7132269, 'йцукенгшщзх', 'qwertyui123@gmail.com', '+75678903456', '$2b$10$hhmWdTv8RdWeJ1ofHOjaTuKBgOo2JUky9za7NTJ.uCcfrH3W2CK/S', 'INDIVIDUAL', '2025-12-01 14:29:11.538', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (4146092, 'Фокеев Кирилл', 'test1@test.com', '+71234567890', '$2b$10$FELoBjJj0J8IeMy2YhKlIeniLkjz86fijJS2HOFJ3XvJ3fnIulg2i', 'INDIVIDUAL', '2025-12-02 10:48:41.186', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);
INSERT INTO public."User" (id, "fullName", email, "phoneNumber", password, "profileType", "createdAt", "updatedAt", rating, "isResetVerified", "roleId", "isAnswersCall", photo, "isEmailVerified", balance, "bonusBalance", "isBanned", "freeAdsLimit", "lastAdLimitReset", "usedFreeAds") VALUES (9851099, 'Черешков Данила Алексеевич', 'chereshkov.da2006@gmail.com', '+79123431910', '$2b$10$hvt0jXBTO6PcqEzKYDKYUO7hivY2kCsC/7Bzwix242L8YDeP6UgnW', 'INDIVIDUAL', '2025-12-02 10:47:25.87', '2025-12-08 12:30:43.354', NULL, false, 1, false, NULL, false, 0, 0, false, 12, '2025-12-24 18:33:10.973', 0);


ALTER TABLE public."User" ENABLE TRIGGER ALL;

--
-- Data for Name: Banner; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."Banner" DISABLE TRIGGER ALL;

INSERT INTO public."Banner" (id, "photoUrl", "createdAt", "updatedAt", place, "navigateToUrl", name, "userId") VALUES (1, 'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/3aaca434-943c-423f-b2b7-2a0e27354f30.png', '2026-01-11 19:18:55.768', '2026-01-11 19:18:55.768', 'PRODUCT_FEED', 'https://yandex.ru', 'Yandex Browser', 7106521);
INSERT INTO public."Banner" (id, "photoUrl", "createdAt", "updatedAt", place, "navigateToUrl", name, "userId") VALUES (2, 'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/88f01b62-b8f4-4fcf-9de7-e7160a2cf286.png', '2026-01-11 19:19:22.773', '2026-01-11 19:19:22.773', 'PROFILE', 'https://google.com', 'Google Browser', 7106521);
INSERT INTO public."Banner" (id, "photoUrl", "createdAt", "updatedAt", place, "navigateToUrl", name, "userId") VALUES (3, 'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/ed332202-ca96-47fb-8b49-425cacd3e739.png', '2026-01-11 19:19:39.78', '2026-01-11 19:19:39.78', 'FAVORITES', 'https://mail.ru', 'Mail.ru', 7106521);
INSERT INTO public."Banner" (id, "photoUrl", "createdAt", "updatedAt", place, "navigateToUrl", name, "userId") VALUES (4, 'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/db7395f3-dd21-49fc-9278-393642b85f19.png', '2026-01-11 19:19:52.414', '2026-01-11 19:19:52.414', 'CHATS', 'https://github.com', 'Github', 7106521);


ALTER TABLE public."Banner" ENABLE TRIGGER ALL;

--
-- Data for Name: BannerView; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."BannerView" DISABLE TRIGGER ALL;



ALTER TABLE public."BannerView" ENABLE TRIGGER ALL;

--
-- Data for Name: Category; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."Category" DISABLE TRIGGER ALL;

INSERT INTO public."Category" (id, name, "createdAt", slug, "updatedAt") VALUES (1, 'Личные вещи', '2025-12-15 19:18:08.497', 'lichnye-veschi', '2025-12-15 17:21:30.479');


ALTER TABLE public."Category" ENABLE TRIGGER ALL;

--
-- Data for Name: SubCategory; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."SubCategory" DISABLE TRIGGER ALL;

INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (1, 'Одежда', 1, '2025-12-15 19:18:08.513', 'odezhda', '2025-12-15 17:21:30.498');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (2, 'Детские товары', 1, '2025-12-15 19:18:08.513', 'detskie-tovary', '2025-12-15 17:21:30.503');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (5, 'Средства реабилитации', 1, '2025-12-15 19:18:08.513', 'sredstva-reabilitatsii', '2025-12-15 17:21:30.508');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (6, 'Школьные товары', 1, '2025-12-15 19:18:08.513', 'shkol-nye-tovary', '2025-12-15 17:21:30.513');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (7, 'Украшения', 1, '2025-12-15 19:18:08.513', 'ukrasheniya', '2025-12-15 17:21:30.518');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (8, 'Продукты питания', 1, '2025-12-15 19:18:08.513', 'produkty-pitaniya', '2025-12-15 17:21:30.522');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (9, 'Животные, растения', 1, '2025-12-15 19:18:08.513', 'zhivotnye-rasteniya', '2025-12-15 17:21:30.527');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (10, 'Бытовая техника', 1, '2025-12-15 19:18:08.513', 'bytovaya-tehnika', '2025-12-15 17:21:30.531');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (11, 'Посуда', 1, '2025-12-15 19:18:08.513', 'posuda', '2025-12-15 17:21:30.536');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (12, 'Мебель', 1, '2025-12-15 19:18:08.513', 'mebel', '2025-12-15 17:21:30.54');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (15, 'Медицинские товары', 1, '2025-12-15 19:18:08.513', 'meditsinskie-tovary', '2025-12-15 17:21:30.544');
INSERT INTO public."SubCategory" (id, name, "categoryId", "createdAt", slug, "updatedAt") VALUES (3, 'Красота и здоровье', 1, '2025-12-15 19:18:08.513', 'krasota-i-zdorov-e', '2025-12-15 17:21:30.548');


ALTER TABLE public."SubCategory" ENABLE TRIGGER ALL;

--
-- Data for Name: SubcategotyType; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."SubcategotyType" DISABLE TRIGGER ALL;

INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (1, 'Мужская', 1, '2025-12-15 19:18:08.514', 'muzhskaya', '2025-12-15 17:21:30.561');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (2, 'Женская', 1, '2025-12-15 19:18:08.514', 'zhenskaya', '2025-12-15 17:21:30.565');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (3, 'Детская девочки', 1, '2025-12-15 19:18:08.514', 'detskaya-devochki', '2025-12-15 17:21:30.57');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (4, 'Детская мальчики', 1, '2025-12-15 19:18:08.514', 'detskaya-mal-chiki', '2025-12-15 17:21:30.574');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (5, 'Ткани, текстиль и фурнитура', 1, '2025-12-15 19:18:08.514', 'tkani-tekstil-i-furnitura', '2025-12-15 17:21:30.578');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (6, 'Сумки, рюкзаки', 1, '2025-12-15 19:18:08.514', 'sumki-ryukzaki', '2025-12-15 17:21:30.582');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (7, 'Аксессуары', 1, '2025-12-15 19:18:08.514', 'aksessuary', '2025-12-15 17:21:30.585');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (8, 'Обувь', 1, '2025-12-15 19:18:08.514', 'obuv', '2025-12-15 17:21:30.59');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (9, 'Игрушки', 2, '2025-12-15 19:18:08.514', 'igrushki', '2025-12-15 17:21:30.595');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (10, 'Детская мебель', 2, '2025-12-15 19:18:08.514', 'detskaya-mebel', '2025-12-15 17:21:30.599');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (11, 'Коляски детские', 2, '2025-12-15 19:18:08.514', 'kolyaski-detskie', '2025-12-15 17:21:30.603');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (12, 'Велосипеды и самокаты', 2, '2025-12-15 19:18:08.514', 'velosipedy-i-samokaty', '2025-12-15 17:21:30.608');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (13, 'Детское питание и посуда', 2, '2025-12-15 19:18:08.514', 'detskoe-pitanie-i-posuda', '2025-12-15 17:21:30.612');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (14, 'Образовательные товары', 2, '2025-12-15 19:18:08.514', 'obrazovatel-nye-tovary', '2025-12-15 17:21:30.616');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (15, 'Уход и гигиена', 2, '2025-12-15 19:18:08.514', 'uhod-i-gigiena', '2025-12-15 17:21:30.62');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (16, 'Косметика для ухода за кожей', 3, '2025-12-15 19:18:08.514', 'kosmetika-dlya-uhoda-za-kozhey', '2025-12-15 17:21:30.624');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (17, 'Средства для ухода за волосами', 3, '2025-12-15 19:18:08.514', 'sredstva-dlya-uhoda-za-volosami', '2025-12-15 17:21:30.629');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (18, 'Уход и гигиена', 3, '2025-12-15 19:18:08.514', 'uhod-i-gigiena', '2025-12-15 17:21:30.633');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (19, 'Приборы и аксессуары', 3, '2025-12-15 19:18:08.514', 'pribory-i-aksessuary', '2025-12-15 17:21:30.636');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (20, 'Парфюмерия', 3, '2025-12-15 19:18:08.514', 'parfyumeriya', '2025-12-15 17:21:30.64');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (21, 'Макияж', 3, '2025-12-15 19:18:08.514', 'makiyazh', '2025-12-15 17:21:30.645');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (22, 'Бады', 3, '2025-12-15 19:18:08.514', 'bady', '2025-12-15 17:21:30.648');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (26, 'Измерительные приборы', 5, '2025-12-15 19:18:08.514', 'izmeritel-nye-pribory', '2025-12-15 17:21:30.653');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (27, 'Ортопедия (бандажи, корсеты)', 5, '2025-12-15 19:18:08.514', 'ortopediya-bandazhi-korsety', '2025-12-15 17:21:30.655');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (28, 'Уходовая косметика', 5, '2025-12-15 19:18:08.514', 'uhodovaya-kosmetika', '2025-12-15 17:21:30.659');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (29, 'Кресла-коляски', 5, '2025-12-15 19:18:08.514', 'kresla-kolyaski', '2025-12-15 17:21:30.662');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (30, 'Спецодежда, трикотаж, компрессионное белье', 5, '2025-12-15 19:18:08.514', 'spetsodezhda-trikotazh-kompressionnoe-bel-e', '2025-12-15 17:21:30.666');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (31, 'Подгузники, пеленки, прокладки', 5, '2025-12-15 19:18:08.514', 'podguzniki-pelenki-prokladki', '2025-12-15 17:21:30.669');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (32, 'Катетеры', 5, '2025-12-15 19:18:08.514', 'katetery', '2025-12-15 17:21:30.673');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (33, 'Средства ухода за стомой', 5, '2025-12-15 19:18:08.514', 'sredstva-uhoda-za-stomoy', '2025-12-15 17:21:30.675');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (34, 'Кресла-стулья санитарные', 5, '2025-12-15 19:18:08.514', 'kresla-stul-ya-sanitarnye', '2025-12-15 17:21:30.678');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (35, 'Специальные устройства', 5, '2025-12-15 19:18:08.514', 'spetsial-nye-ustroystva', '2025-12-15 17:21:30.682');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (36, 'Калоприемники, уроприемники', 5, '2025-12-15 19:18:08.514', 'kalopriemniki-uropriemniki', '2025-12-15 17:21:30.686');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (37, 'Трости, костыли', 5, '2025-12-15 19:18:08.514', 'trosti-kostyli', '2025-12-15 17:21:30.69');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (38, 'Вертикализаторы, опоры', 5, '2025-12-15 19:18:08.514', 'vertikalizatory-opory', '2025-12-15 17:21:30.694');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (39, 'Матрасы', 5, '2025-12-15 19:18:08.514', 'matrasy', '2025-12-15 17:21:30.696');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (40, 'Кровати медицинские', 5, '2025-12-15 19:18:08.514', 'krovati-meditsinskie', '2025-12-15 17:21:30.701');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (41, 'Письменные принадлежности', 6, '2025-12-15 19:18:08.514', 'pis-mennye-prinadlezhnosti', '2025-12-15 17:21:30.704');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (42, 'Бумажная продукция', 6, '2025-12-15 19:18:08.514', 'bumazhnaya-produktsiya', '2025-12-15 17:21:30.708');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (43, 'Принадлежности для рисования и творчества', 6, '2025-12-15 19:18:08.514', 'prinadlezhnosti-dlya-risovaniya-i-tvorchestva', '2025-12-15 17:21:30.71');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (44, 'Органайзеры и хранение', 6, '2025-12-15 19:18:08.514', 'organayzery-i-hranenie', '2025-12-15 17:21:30.714');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (45, 'Учебные пособия и инструменты', 6, '2025-12-15 19:18:08.514', 'uchebnye-posobiya-i-instrumenty', '2025-12-15 17:21:30.717');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (46, 'Рюкзаки и сумки', 6, '2025-12-15 19:18:08.514', 'ryukzaki-i-sumki', '2025-12-15 17:21:30.721');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (47, 'Прочее', 6, '2025-12-15 19:18:08.514', 'prochee', '2025-12-15 17:21:30.723');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (48, 'Ювелирные изделия', 7, '2025-12-15 19:18:08.514', 'yuvelirnye-izdeliya', '2025-12-15 17:21:30.727');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (49, 'Бижутерия', 7, '2025-12-15 19:18:08.514', 'bizhuteriya', '2025-12-15 17:21:30.73');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (50, 'Часы', 7, '2025-12-15 19:18:08.514', 'chasy', '2025-12-15 17:21:30.733');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (51, 'Готовые продукты', 8, '2025-12-15 19:18:08.514', 'gotovye-produkty', '2025-12-15 17:21:30.737');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (52, 'Напитки', 8, '2025-12-15 19:18:08.514', 'napitki', '2025-12-15 17:21:30.74');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (53, 'Заморозки, полуфабрикаты', 8, '2025-12-15 19:18:08.514', 'zamorozki-polufabrikaty', '2025-12-15 17:21:30.744');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (54, 'Домашние животные', 9, '2025-12-15 19:18:08.514', 'domashnie-zhivotnye', '2025-12-15 17:21:30.748');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (55, 'С/х животные', 9, '2025-12-15 19:18:08.514', 's-h-zhivotnye', '2025-12-15 17:21:30.751');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (56, 'Рептилии', 9, '2025-12-15 19:18:08.514', 'reptilii', '2025-12-15 17:21:30.756');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (57, 'Растения комнатные', 9, '2025-12-15 19:18:08.514', 'rasteniya-komnatnye', '2025-12-15 17:21:30.759');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (58, 'Культурные растения', 9, '2025-12-15 19:18:08.514', 'kul-turnye-rasteniya', '2025-12-15 17:21:30.764');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (59, 'Декоративные уличные растения', 9, '2025-12-15 19:18:08.514', 'dekorativnye-ulichnye-rasteniya', '2025-12-15 17:21:30.768');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (61, 'Кухонная', 10, '2025-12-15 19:18:08.514', 'kuhonnaya', '2025-12-15 17:21:30.777');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (62, 'Бытовая', 10, '2025-12-15 19:18:08.514', 'bytovaya', '2025-12-15 17:21:30.779');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (63, 'Для приготовления пищи', 11, '2025-12-15 19:18:08.514', 'dlya-prigotovleniya-pischi', '2025-12-15 17:21:30.783');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (64, 'Для хранения', 11, '2025-12-15 19:18:08.514', 'dlya-hraneniya', '2025-12-15 17:21:30.786');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (65, 'Для сервировки', 11, '2025-12-15 19:18:08.514', 'dlya-servirovki', '2025-12-15 17:21:30.79');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (66, 'Для приёма пищи', 11, '2025-12-15 19:18:08.514', 'dlya-priema-pischi', '2025-12-15 17:21:30.793');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (67, 'Мягкая мебель', 12, '2025-12-15 19:18:08.514', 'myagkaya-mebel', '2025-12-15 17:21:30.798');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (68, 'Корпусная мебель', 12, '2025-12-15 19:18:08.514', 'korpusnaya-mebel', '2025-12-15 17:21:30.801');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (69, 'Мебель для кухни', 12, '2025-12-15 19:18:08.514', 'mebel-dlya-kuhni', '2025-12-15 17:21:30.805');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (70, 'Мебель для спальни', 12, '2025-12-15 19:18:08.514', 'mebel-dlya-spal-ni', '2025-12-15 17:21:30.808');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (71, 'Садовая мебель', 12, '2025-12-15 19:18:08.514', 'sadovaya-mebel', '2025-12-15 17:21:30.813');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (72, 'Офисная мебель', 12, '2025-12-15 19:18:08.514', 'ofisnaya-mebel', '2025-12-15 17:21:30.815');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (74, 'Оборудование для клиник', 15, '2025-12-15 19:18:08.514', 'oborudovanie-dlya-klinik', '2025-12-15 17:21:30.821');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (75, 'Медицинская мебель', 15, '2025-12-15 19:18:08.514', 'meditsinskaya-mebel', '2025-12-15 17:21:30.825');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (60, 'Доп товары (горшки, грунт, кормилки, поилки, средства по уходу за растениями, инструменты, корма, игрушки, клетки, аксессуары)', 9, '2025-12-15 19:18:08.514', 'dop-tovary-gorshki-grunt-kormilki-poilki-sredstva-po-uhodu-za-rasteniyami-instrumenty-korma-igrushki-kletki-aksessuary', '2025-12-15 17:21:30.772');
INSERT INTO public."SubcategotyType" (id, name, "subcategoryId", "createdAt", slug, "updatedAt") VALUES (73, 'Диагностическое оборудование', 15, '2025-12-15 19:18:08.514', 'diagnosticheskoe-oborudovanie', '2025-12-15 17:21:30.819');


ALTER TABLE public."SubcategotyType" ENABLE TRIGGER ALL;

--
-- Data for Name: Product; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."Product" DISABLE TRIGGER ALL;

INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (5231119, 'Бусы б/у', 1000, 'USED', 'Красные, из жемчуга', 'г Екатеринбург, ул Чкалова', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/cfb8be90-e717-49b0-a1ef-c0f0ed43b623.png}', 1, 3, 7391202, '2025-11-28 09:11:49.34', '2025-12-16 09:53:24.161', 16, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (8640334, 'Нутриен energy питание', 2700, 'NEW', 'Cмecь Nutrien enеrgy, диетичeское лечeбноe  питание,

Питаниe для oнкoбольных , питаниe для ocлaблeнных, питание пocлe опеpaции, питание, обогащённое витаминaми и микрoэлeмeнтaми.

Продукт готовый к упoтpеблению 200 мл, 300 ккaл.
Пoдxoдит для онкoбольных, пoслeопepациoнныx взрослыx и дeтeй с 3 лет для вoсстановления сил и энергии.', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/6f42b7d6-7a7d-46ca-85c3-611b159a8a0a.png}', 1, 8, 9851099, '2025-12-02 11:08:58.23', '2025-12-16 09:53:24.246', 51, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9305563, 'Капельница', 200, 'NEW', 'Просто капельница', 'г Оренбург, ул Харьковская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b0a71fb2-6719-4085-b554-d16b5cf9b2a2.webp}', 1, 15, 4146092, '2025-12-02 11:11:52.698', '2025-12-16 09:53:24.191', 74, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (7384341, 'ВАЗ 2107', 435000, 'NEW', 'Продаётся готовый проект под RDS. Соответствует всем стандартам турниров и сходок. Гарантия на проект год.', 'Степной, Дзержинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/af9ca37d-87d9-44bc-b0aa-b1fc99737315.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/7a511cc5-d998-49c3-8f53-5dd88abd875b.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1d044f53-c12f-4841-9f4b-b486e551411a.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5c8cfe09-cd90-489a-8aad-0f8c2e80f6f4.jpg}', 1, 10, 2321239, '2025-11-28 09:18:17.344', '2025-12-16 09:53:24.18', 61, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2161612, 'Очередной товар дня!', 35000, 'NEW', '1) пусть будет текст
2) здесь еще что-то
**
💥
🟩
ККЕКЕКЕЕУУЦКУ""
                                              ЦЕНТР
          ТАБУЛЯЦИЯ СмещЕНИЕ

', '18, улица Расковой, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/6638b79f-2357-46ff-9010-ba9175ce50db.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c02fa4fd-6284-45a5-8cd7-61583db872fe.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fe8ad9f4-5664-4832-b5be-dc1f4df2adcf.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5488c8b5-af91-4294-85c6-0bb7d48145b6.jpg}', 1, 12, 6669460, '2025-12-01 08:35:56.623', '2025-12-16 09:53:24.183', 67, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9122333, 'Зипка', 5000, 'NEW', 'Кофта теплая на замке', '35, улица 9 Января, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f64e149a-e711-4173-85ea-98db13c3ca1e.png}', 1, 1, 6038643, '2025-12-02 10:59:24.476', '2025-12-16 09:53:24.186', 2, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (7659684, 'Протеин 1000гр', 1500, 'NEW', 'Вкус шоколад, 1000 грамм', '2, улица 13-я Линия, Линии, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460040, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a745c362-a754-4a92-abb9-b8969bebead7.png}', 1, 8, 7391202, '2025-11-28 09:14:04.157', '2025-12-16 09:53:24.199', 51, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9262881, 'Набор украшений для пирсинга', 4000, 'NEW', NULL, '12А, Больничный проезд, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3c250f31-591b-4346-b1fb-1b3bf70f2c73.webp}', 1, 3, 6053931, '2025-11-28 09:17:54.801', '2025-12-16 09:53:24.203', 16, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (4215912, 'Детские книжки по математике', 1000, 'USED', 'Превосходный источник знаний для вашего ребенка', 'Hawthorne Street, Кламат-Фолс, Klamath County, Орегон, 97601, Соединённые Штаты Америки', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5e8875c2-aec8-4b2f-b618-2e220defa9cf.webp,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/2248f3b6-63b3-44e7-84a4-ef35b2d7bcdc.jpg}', 1, 2, 1208299, '2025-11-28 09:21:17.846', '2025-12-16 09:53:24.208', 14, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (1512888, 'Померанский шпиц, щенок', 1, 'USED', 'Продаетcя очapовательная мини дeвочкa помepанcкoгo шпицa.28.09.2025 гoдa poждeния.
Дoкументы: Вет пacпoрт прививки oбpаботки по возрасту.
Очeнь лаcкoвaя игpивая контактная .
Пpиучeна к пелeнки.
Kушaeт суxой коpм
Отличнo ладит c дeтьми и другими живoтными .
Ищeм добрыe зaботливыe руки.
Рoдитeли:
Мама - померaнский шпиц, белый окрас (3,5 кг)
Папа - померанский шпиц, пати колор (3 кг)
Будет не больше 2,5 кг.', '77/2, улица Терешковой, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/26aadd6d-3f95-4315-9a51-c59257705c32.png}', 1, 9, 9851099, '2025-12-02 11:21:39.796', '2025-12-16 09:53:24.422', 54, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (5142108, 'Котёнок в добрые руки', 1, 'USED', 'котёнок около 4 месяцев, стерелизован, мальчикрыжий, очень активный, игривый, с другими животными и детьми ладит. очень ласковый, постоянно мурчит', '14, улица Терешковой, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/333a3c16-155f-4e82-9fbe-a878937a6f9f.png}', 1, 9, 9851099, '2025-12-02 11:24:58.868', '2025-12-16 09:53:24.429', 54, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (5868178, 'папавпа', 55454, 'NEW', 'павпвапаfdggdfgf212121', '«Урал», Ленинский район, Пригородный, Пригородный сельсовет, Оренбургский район, Оренбургская область, Приволжский федеральный округ, 460041, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/396aecf1-b980-45ec-bd5a-ba238f1fdefb.jpg}', 1, 5, 7106521, '2025-12-03 19:36:04.58', '2026-01-14 19:08:58.604', 27, NULL, false, 'MODERATE', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (4332941, 'Графин в виде рыбы', 500, 'NEW', 'Замечательный графин в виде рыбы', 'г Оренбург, ул Киевская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/67533ddf-2495-4aed-b405-a68922a398bf.jpg}', 1, 11, 4146092, '2025-12-02 10:50:32.345', '2025-12-16 09:53:24.22', 63, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (3982248, 'Сковорода антипригарная', 1000, 'NEW', 'Сковорода. Можно пожарить все что угодно', 'г Оренбург, ул Днепропетровская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/299e425b-f6c0-49bb-a7e8-3c7c591ce39d.jpg}', 1, 11, 4146092, '2025-12-02 10:53:18.109', '2025-12-16 09:53:24.223', 63, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (4758351, 'Медицинское кресло', 15798, 'NEW', 'Инвалидное кресло для комфортной и активной жизни.
*  Мягкое сиденье и удобная спинка обеспечат комфорт даже при длительном использовании. Легко складывается для транспортировки.
*  Регулируется под индивидуальные потребности. [Указать преимущества, например, наличие подголовника, антиопрокидыватели. 

✈✈✈✈✈ Можно отправить!

Цена реальная. Звоните или пишите" ', 'г Оренбург, пр-кт Победы, д 10', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5fe6a6e0-d9a6-418d-bca1-dda8509a758f.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/659a2e57-a129-468f-9dea-c500bce1dcaa.jpg,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b4410b84-3d31-43b7-a9b0-0d10516b503b.jpg}', 1, 5, 6669460, '2025-12-01 09:07:28.717', '2025-12-16 09:53:24.346', 29, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2693271, 'Стакан', 200, 'NEW', 'Просто стакан.', 'г Оренбург, ул Житомирская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/4cfed187-e55c-4014-8c59-cd2450aca91e.jpg}', 1, 11, 4146092, '2025-12-02 10:55:45.209', '2025-12-16 09:53:24.226', 63, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (3563632, 'Ингалятор', 2000, 'NEW', 'Ингалятор для ингаляций', 'г Оренбург, ул Луганская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/08d81171-c46d-4857-b7d9-bb7a983d5ab4.jpg}', 1, 15, 4146092, '2025-12-02 11:05:14.587', '2025-12-16 09:53:24.234', 73, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2865910, 'Посуда детская', 1500, 'NEW', 'Детская посуда для кормления', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f252c198-af82-42a9-a80b-e42a052caae3.png}', 1, 2, 6038643, '2025-12-02 11:06:21.305', '2025-12-16 09:53:24.24', 13, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (8882052, 'Алоэ вера лечебный 3 года, есть 1 год', 200, 'USED', 'Алое Вера, лечебное 3х детки, есть однолетки', 'В. И. Ленину, Ленинская улица, Кузнечный, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/17ad0794-e12d-433a-9958-528bba02bf87.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/09b01358-7870-4abd-a66b-7cfcab7ecec9.png}', 1, 9, 9851099, '2025-12-02 11:26:21.919', '2025-12-16 09:53:24.215', 54, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2016352, 'Украшения', 1000, 'USED', 'Продам укрошенияБраслет -500
Серьги - 300
Кольцо 10 - 250
Все вместе 1000', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a28160dd-9d06-4750-b6fe-6045f6a3df8b.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e49061f1-c633-4254-bc95-b9e06ae322ae.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/15b89ff9-deb0-4067-80a9-77151e9ad946.png}', 1, 7, 9851099, '2025-12-02 11:06:25.377', '2025-12-16 09:53:24.242', 48, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (8948419, 'Массажный стол', 4000, 'NEW', 'Просто массажный стол', 'г Оренбург, ул Запорожская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/002d563c-310a-4287-ac15-5826d88e5d37.jpg}', 1, 15, 4146092, '2025-12-02 11:07:04.422', '2025-12-16 09:53:24.237', 74, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6091694, 'Домашние полуфабрикаты, пельмени и тд', 650, 'NEW', 'Прoдаём cвoю дoмaшнюю пpодукцию из магазина и пpинимаeм закaзы.Продукция oчeнь вкусная, из домaшниx яиц. Xaляль. Фaрш делаeм caми, ни одной жилки плёнки тaм нeт.
', 'улица Цвиллинга, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c14573ef-45ac-488a-a7c5-663ffee7150e.png}', 1, 8, 9851099, '2025-12-02 11:16:29.78', '2025-12-16 09:53:24.418', 51, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2273041, 'Спрей', 600, 'USED', 'Защитная пленка для кожи', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/002307c3-8394-46f2-b18a-952a389efc6d.png}', 1, 5, 6038643, '2025-12-02 11:25:57', '2025-12-16 09:53:24.276', 28, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6617171, 'Кресло-коляска', 5000, 'USED', 'Кресло-коляска для инвалидов Ortonica Olvia 30', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/695dac97-c44a-42c0-9307-2cc8ff3bcaab.png}', 1, 5, 6038643, '2025-12-02 11:21:05.818', '2025-12-16 09:53:24.271', 29, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6628130, 'Матрас', 3000, 'USED', 'Матрас для восстанволения', '61А, улица Орлова, Новостройка, Промышленный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/7d6cda50-98e5-42bc-9a31-08b62188d9fa.png}', 1, 5, 6038643, '2025-12-02 11:22:24.486', '2025-12-16 09:53:24.259', 26, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (1063797, 'Блокнот', 300, 'NEW', 'Блокнот Осенняя эстетика', '92, улица Орджоникидзе, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/28d3f9c3-0ef8-4de7-8460-8d402294aa14.png}', 1, 6, 6038643, '2025-12-02 11:31:05.266', '2025-12-16 09:53:24.44', 41, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6193207, 'Вилка', 100, 'NEW', 'Просто вилка', 'г Оренбург, ул Одесская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/d66fc1ce-a773-4ff1-8d0c-eaaf5515e495.webp}', 1, 11, 4146092, '2025-12-02 10:51:56.847', '2025-12-16 09:53:24.358', 63, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (4372887, 'Компрессорный ингалятор', 2000, 'NEW', 'Компрессорный ингалятор', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f178392b-1965-4569-b91c-c6efd48b56da.png}', 1, 5, 6038643, '2025-12-02 11:18:32.771', '2025-12-16 09:53:24.266', 35, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2854985, 'Садовые качели', 10000, 'NEW', 'Просто качели. Качаться весело', 'г Оренбург, ул Шевченко', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5ab65b78-ba32-48e4-aa0e-a45203f815ab.jpg}', 1, 12, 4146092, '2025-12-02 10:58:20.061', '2025-12-16 09:53:24.369', 67, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6734788, 'Кроватка', 3000, 'USED', 'Кроватка для новорожденных', '199, Комсомольская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/42b2d16e-20eb-456c-8b73-3237d225a549.png}', 1, 2, 6038643, '2025-12-02 11:10:27.655', '2025-12-16 09:53:24.395', 10, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9265239, 'Люлька', 2000, 'USED', 'Люлька детская', '68, улица Кичигина, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/51e36a8d-6904-4a08-96bc-d1d449241608.png}', 1, 2, 6038643, '2025-12-02 11:11:54.214', '2025-12-16 09:53:24.196', 10, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (8507601, 'Крем для рук', 200, 'NEW', 'Просто крем для рук', 'г Оренбург, поселок Нижнесакмарский, ул Николаевская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fb99e19c-4e45-43b7-90fc-1c48c50106d7.webp}', 1, 3, 4146092, '2025-12-02 11:17:30.55', '2025-12-16 09:53:24.251', 16, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (7270506, 'Пароочиститель для дома мощный, новые', 1650, 'NEW', 'Унивеpсaльный паровой очиcтитель – этo эффективная бытовaя тeхникa для убоpки дoмa, coздaнная для удобствa и экoнoмии времeни. Этoт мощный парогенератoр cтaнeт вaшим надежным помощникoм в бoрьбе c зaгpязнeниями на куxне, мебeли и другиx повepхнocтях.
', 'Мегаполис, 22, улица Володарского, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/965b0669-84ea-48c0-80e3-3e5ce0fe5022.png}', 1, 10, 9851099, '2025-12-02 11:35:21.921', '2025-12-16 09:53:24.282', 61, NULL, false, 'MODERATE', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (3207807, 'SWEETPEEPS золотые украшения', 7000, 'NEW', 'Золотые украшения с фианитами', 'Уральская улица, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f3a550d2-fc17-4548-943e-b62d66c014eb.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/afa0dc23-c3de-4ab1-8844-a4506bb49309.png}', 1, 7, 6038643, '2025-11-28 09:14:30.517', '2025-12-16 09:53:24.284', 48, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9256863, 'Chevrolet Corvette C7', 8500000, 'USED', 'Корвет был угнан у курседа', 'Оренбургский Колледж Экономики и Информатики, 11, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/d0c85309-e830-459b-9f87-a672313a465e.jpg}', 1, 7, 8964288, '2025-11-28 09:15:31.784', '2025-12-16 09:53:24.287', 48, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (5609249, 'ДМРВ на ваз 2107', 7000, 'NEW', 'Датчик массового расхода воздуха', '48, улица Коминтерна, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/05f3e4d2-118c-4262-a12a-8804cddb31d7.webp}', 1, 10, 4761896, '2025-11-28 09:16:40.023', '2025-12-16 09:53:24.29', 61, 'https://yandex.ru/video/preview/13520813755431483017', false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2388612, 'Фундук культурный', 280, 'USED', 'Прoдaю фундук 2024г cбopa, собственный небoльшой cад в предгopьяx Кавказa, бeз xимии тoлькo органикa.

Bcе сopта выращиваемые мной имеют лучшиe вкусoвыe xaрактериcтики и oтносятcя к cтoлoвым сoртам, oбладaют плотным ядpoм и приятным выpaженным маcляниcтым вкуcoм, который не сравним с дешёвыми cетевыми безвкусными орешками.
Предлагаю микс сортов Трапезунд, Анаклиури, Президент.

Возможна доставка авитодоставкой до 20кг или транспортной компанией от 30.', '19/2, улица Бурзянцева, Кузнечный, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a0505026-7df1-43c2-a19c-16e30e07a690.png}', 1, 9, 9851099, '2025-12-02 11:27:45.043', '2025-12-16 09:53:24.218', 54, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2139014, 'Пельмени домашние', 380, 'NEW', 'Пpeдcтaвляeм вaшему вниманию пельмени, манты, хинкaли, ваpеники pучнoй лeпки.', '18, Матросский переулок, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/19b43a4f-6ae7-4539-bee3-78a337e8e3c8.png}', 1, 8, 9851099, '2025-12-02 11:17:35.121', '2025-12-16 09:53:24.254', 51, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (4966297, 'Эублефар', 4000, 'USED', 'Продаются малыши эублефары различных морф. Едят разморозку, линяют хорошо, все процессы в норме.
', '2, Госпитальный переулок, Аренда, Промышленный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/99f1ee8b-a4e7-400c-9591-74ad68a831b6.png}', 1, 9, 9851099, '2025-12-02 11:29:46.562', '2025-12-16 09:53:24.292', 54, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (4081087, 'Собака', 100, 'USED', 'Собака овчарка', '3/5, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b48c5093-8d65-47e0-84be-e1736ceffbe9.png}', 1, 9, 9371169, '2025-11-28 09:10:38.847', '2025-12-16 09:53:24.295', 54, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (1979749, 'Игуана', 285000, 'USED', 'Xoроший спoкoйный пaрень в самoм рaсцвeтe игуаниx cил.

Зовут Яша, 19 лет, любит тепло и голубику.', '26Б, улица Шевченко, Дзержинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/77280a5a-a493-49e1-aeb4-5bb7cbe97653.webp}', 1, 9, 2321239, '2025-11-28 09:11:59.098', '2025-12-16 09:53:24.299', 54, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (5611056, 'Monster Energy Pipeline Punch', 250, 'NEW', 'Тонизирующий напиток с изысканным вкусом!', 'Hawthorne Street, Кламат-Фолс, Klamath County, Орегон, 97601, Соединённые Штаты Америки', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/bfa67311-f947-41c7-b275-e7d28e1db313.jpg}', 1, 8, 1208299, '2025-11-28 09:16:21.207', '2025-12-16 09:53:24.302', 51, 'https://vk.com/video-129440544_456249335', false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (8083712, 'Концтовары', 700, 'NEW', 'Набор канцтоваров для школы и офиса Лапки котика 5 предметов', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e97b5f41-387d-42c2-bc90-a342ab3403bb.png}', 1, 6, 6038643, '2025-12-02 11:27:09.759', '2025-12-16 09:53:24.455', 41, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (1961051, 'Наклейки', 700, 'NEW', 'Наклейки для ежедневника Школьная эстетика', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/ab80ea7d-50f6-4bcf-beff-ec2a68d97299.png}', 1, 6, 6038643, '2025-12-02 11:30:16.141', '2025-12-16 09:53:24.457', 41, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (8257036, 'Пенал', 590, 'NEW', 'Милый эстетичный большой пенал школьный', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5090f83f-e7b2-42af-bc8e-2937850f8952.png}', 1, 6, 6038643, '2025-12-02 11:31:59.433', '2025-12-16 09:53:24.459', 41, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (3824376, 'Украшения ручной работы', 1000, 'NEW', 'Украшения ручной работы на заказ по Вашим эскизам/фото. Стоимость украшений на фото 1000р. Срок изготовления: 4-7 дней.

', '46, улица 9 Января, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/4164bbf3-aa1d-4dac-b361-dd22fc5c2001.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/471637e1-9898-41d2-9dc1-c112c642c296.png}', 1, 7, 9851099, '2025-12-02 11:01:00.313', '2025-12-16 09:53:24.189', 48, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (3244052, 'Дубленка', 8000, 'NEW', 'Дубленка зимняя', 'улица Рокоссовского, Горка, Дзержинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b6118a2c-fce6-4c04-826b-b71e8953afe4.png}', 1, 1, 6038643, '2025-12-02 11:01:27.168', '2025-12-16 09:53:24.376', 2, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (8776759, 'Кофемашина Thomson CF20A02', 11399, 'NEW', 'Рабoчaя бытoвая тexника, намного дeшевлe, чем в мaгaзинe;



Нeт тapы для мoлoкa

- Любыe пpoверки при caмoвывозе;

- Пpи пpиемке товара вся теxника пpoвepяется на рaбoтоспocoбнoсть;

- Oтправляeм Авитo доcтaвкой;

- Пpи дoставке тoвap упaковываeтся по высшему урoвню.

Bитpинный образец:

• товар новый, стоял на витрине в магазине;

• может быть повреждена заводская упаковка;

• возможны незначительные потёртости или повреждения корпуса, которые никак не влияют на работоспособность.

За фотографиями дефектов обращайтесь в лс

Самовывоз возможен из 2-х точек: Метро Текстильщики, Метро Шипиловская.

В нашем профиле большой ассортимент разнообразной бытовой техники. Советуем заглянуть!

Больше техники в нашем телеграмм-канале

Переходите там большие скидки!

', '44, улица Кирова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/673d92ac-2d1c-415f-8beb-02bd64a3b69d.png}', 1, 10, 9851099, '2025-12-02 11:34:25.201', '2025-12-16 09:53:24.174', 61, NULL, false, 'MODERATE', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (1250840, 'Кресло-горилла', 170000, 'NEW', 'Кресло-горилла удобное, выполнено из лучших материалов.', '37А, Илекская улица, село имени 9 Января, Красноуральский сельсовет, Оренбургский район, Оренбургская область, Приволжский федеральный округ, 460501, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/228256c7-91a9-4c50-9a2e-a2804917075b.png}', 1, 12, 7391202, '2025-11-28 09:16:46.544', '2025-12-16 09:53:24.305', 67, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9368305, 'Майка', 20000, 'USED', 'очень крутые маечки с аниме принтами, у2к вайб имеется🪽размер S, полиэстер
цена 500 рублей за штуку
по любым вопросам пишите!!

', '3, улица Аксакова, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/33229854-0ccb-4ecc-b7a8-3ab2965d8fdc.png}', 1, 1, 8633592, '2025-11-28 09:18:16.763', '2025-12-16 09:53:24.307', 1, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (3334788, 'Посуда для сервировки Estetic', 3500, 'NEW', 'Вся посуда выполнена в минималистичных стилях, из качественных материалов, подойдет на каждый день', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3b0ddf2a-19a7-4ecf-99c9-9dfc705d35c7.png}', 1, 11, 6038643, '2025-11-28 09:11:35.533', '2025-12-16 09:53:24.316', 63, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6901799, 'Кошка', 10, 'USED', 'Кошка домашняя', '5, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/8f6c52ee-023b-41d0-befc-39612d968abf.webp}', 1, 9, 9371169, '2025-11-28 09:11:58.971', '2025-12-16 09:53:24.319', 54, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (4224343, 'Салонный фильтр на ваз 2110', 1000, 'NEW', 'салонный фильтр подходит на автомобили ваз2110,2112', '20, улица Кобозева, Кузнечный, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/22cb6191-445a-4ee0-80f2-92cc96055093.webp}', 1, 10, 4761896, '2025-11-28 09:12:25.319', '2025-12-16 09:53:24.326', 61, 'https://yandex.ru/video/preview/9506785745966413491', false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (1300264, 'Ford Mustang', 2500000, 'NEW', 'Самый лучший автомобиль в мире', 'Оренбургский Колледж Экономики и Информатики, 11, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c91c9b99-afe9-4d10-8a81-b3cee1c12296.jpg}', 1, 10, 8964288, '2025-11-28 09:12:03.576', '2025-12-16 09:53:24.322', 61, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (1970246, 'Козы камерунские', 3000, 'NEW', 'Продаются козочки камерунские,разного возраста, есть два козлика для покрытия, покрытие 3 тыс', '"Воздух" конный клуб, 9, Бассейный переулок, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/bfba4d6b-5f64-40b1-a0e2-8f232b9140ea.webp}', 1, 9, 2321239, '2025-11-28 09:14:23.82', '2025-12-16 09:53:24.33', 54, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9915250, 'Платье горничной', 1200, 'USED', 'платье горничной в хорошем состоянии , нету только ободка осталась только от него ткань, если нужно доп фото пишите, к платью идет бантик и фартук

', 'г Оренбург', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/dd3f08a3-bdcc-4593-b35f-1e186ce5262a.png}', 1, 1, 8633592, '2025-11-28 09:14:31.733', '2025-12-16 09:53:24.334', 1, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9863001, 'Набор золотых украшений', 2000, 'NEW', NULL, 'Лицей №2, Красная улица, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e7670256-5300-4433-82a5-edf31f999776.webp}', 1, 7, 6053931, '2025-11-28 09:12:06.874', '2025-12-16 09:53:24.324', 48, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (7378626, 'Этно украшения', 300, 'USED', 'Украшения в этническом стиле! серьги, браслеты, ожерелья, броши и т, д', '26, улица Кирова, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fe07794a-1f08-4d05-8569-a90fa9c75a56.png}', 1, 7, 9851099, '2025-12-02 10:58:49.981', '2025-12-16 09:53:24.336', 48, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (5510664, 'Джинсы', 2500, 'NEW', 'Джинсы в новом состояние', '48, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/935c3f3a-3998-4c35-97aa-d83c3b4c3beb.png}', 1, 1, 6038643, '2025-12-02 10:51:54.306', '2025-12-16 09:53:24.361', 2, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2207276, 'Ложка', 100, 'NEW', 'Просто ложка', 'г Оренбург, ул Львовская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1582c805-78ef-49e4-a5bc-c875f429af60.webp}', 1, 11, 4146092, '2025-12-02 10:54:33.568', '2025-12-16 09:53:24.356', 63, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (4523969, 'Платье', 2000, 'NEW', 'Платье летнее разных расцветок', '2, улица Богдана Хмельницкого, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460024, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/497ae48d-8b31-443d-8a6f-9f00be0ac793.png}', 1, 1, 6038643, '2025-12-02 10:54:44.305', '2025-12-16 09:53:24.353', 2, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (3506516, 'Разобранный кубик рубика', 10, 'USED', 'не смог собрать', 'Оренбургский Колледж Экономики и Информатики, 11, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460001, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/c55064c4-df45-485c-80a3-7253e48ff798.jfif,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/d4b88f23-5430-4948-96f8-b521befb052b.jpg}', 1, 2, 8964288, '2025-11-28 09:19:52.442', '2025-12-16 09:53:24.339', 14, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6003323, 'Кресло офисное', 5000, 'NEW', 'Удобное кресло', 'г Оренбург, ул Богдана Хмельницкого', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a0763806-d3e2-4699-9674-487603f386a3.jpg}', 1, 12, 4146092, '2025-12-02 10:56:54.12', '2025-12-16 09:53:24.35', 67, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2105765, 'Тест', 20000, 'NEW', 'Описание', 'Вита Экспресс, улица Чкалова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f76a26f7-4801-4c7f-9166-6b2869b5a765.jpg}', 1, 8, 3235109, '2025-12-01 05:50:33.37', '2025-12-16 09:53:24.341', 51, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (1122280, 'Кресло-коляска', 45000, 'USED', 'новая', 'г Оренбург', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/0a43c4d3-222d-421b-8528-7e3e59cc909a.jpg}', 1, 15, 2681599, '2025-12-01 08:10:41.21', '2025-12-16 09:53:24.342', 75, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (5902819, 'Свитер', 3000, 'NEW', 'Свитер теплый из мягкой ткани', '3/1, Телевизионный переулок, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460024, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/8ee766c3-6a41-4953-881d-c48ff14a1add.png}', 1, 1, 6038643, '2025-12-02 10:57:39.627', '2025-12-16 09:53:24.366', 2, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (4267180, 'Табурет', 500, 'NEW', 'Просто табурет.', 'г Оренбург, ул Полтавская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/2a965649-9cf7-4e34-a427-8926bc88b2c9.jpg}', 1, 12, 4146092, '2025-12-02 11:00:14.233', '2025-12-16 09:53:24.371', 67, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6752957, 'Диван', 6000, 'NEW', 'Просто диван', 'г Оренбург, ул Гоголя', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/28a67960-b058-4732-8bf1-501b4d4cca5a.webp}', 1, 12, 4146092, '2025-12-02 11:01:04.406', '2025-12-16 09:53:24.373', 67, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (7213485, 'Кровать', 10000, 'NEW', 'Удобная кровать. Евродвушка', 'г Оренбург, Крымский пер', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a83a8772-60ab-46e0-ac5e-3130ef9deb81.webp}', 1, 12, 4146092, '2025-12-02 11:02:27.961', '2025-12-16 09:53:24.381', 67, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (3184247, 'Украшения в русском стиле', 2800, 'NEW', 'Украшения в русском стиле из натуральных камней и керамических бусин с подвесками ручной работы: неваляшки, Петушки, лошадки.', 'Оренбургская обл, Оренбургский р-н, тер. СНТ Клуб имени Чкалова, д 11', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/1c859471-e841-471b-a425-cd312047cd68.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/9114495f-1926-48c5-bcbd-336ef851b323.png,https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3dd7eee0-ee94-440c-84e0-16d5dfb7740e.png}', 1, 7, 9851099, '2025-12-02 11:03:22.662', '2025-12-16 09:53:24.4', 48, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9380113, 'Детские игрушки', 1000, 'USED', 'Набор детский игрушек', '24, Луговая улица, Восточный, Сотки, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/e8e9b6ba-a396-477d-b8c3-000cc9e85c0f.png}', 1, 2, 6038643, '2025-12-02 11:04:58.661', '2025-12-16 09:53:24.389', 9, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6497808, 'Тонометр', 2000, 'NEW', 'Тонометр. Давление меряет еще что-то там', 'г Оренбург, ул Донецкая', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/50400626-7e5c-41a0-8b2a-d7d753a627cd.jpg}', 1, 15, 4146092, '2025-12-02 11:03:58.743', '2025-12-16 09:53:24.384', 73, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9783545, 'Ванночка', 3000, 'USED', 'Ванна для купания новорожденного', '6Б, Телевизионный переулок, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460024, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/ef7cf831-d301-4470-8e08-45a9662ccc25.png}', 1, 2, 6038643, '2025-12-02 11:08:32.807', '2025-12-16 09:53:24.392', 15, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (8436378, 'Серебряные украшения', 1500, 'USED', 'Пoд номeром 1: сеpежки с розoвым камнeм 1000 рублей. Под номeрoм 2: нaбop cepeжки и кольцо с жeлтым кaмнeм 2000 рублей зa нaбор. Под номером 3: набoр cepeжки, кольцо и подвеcкa с зелeным кaмнeм 2000 pублей зa набоp. Под нoмеpoм 5: сepежки с рoзoвым кaмнем 1000 рублeй. Серебрянaя цeпoчка 2000 рублей. Кольцо с белым камнем и две подвески с белыми камнями- по 500 рублей каждая. Серебро все в хорошем состоянии', 'Фармленд, 52, Советская улица, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/053660d3-2e5f-4f54-a2f8-c11767cd53fc.png}', 1, 7, 9851099, '2025-12-02 11:04:37.805', '2025-12-16 09:53:24.387', 48, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6300121, 'Кровать', 10000, 'USED', 'Кровать детская', '139, Ташкентская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5b84facc-aec9-43ff-aa6e-e6a5922f31f2.png}', 1, 2, 6038643, '2025-12-02 11:13:18.658', '2025-12-16 09:53:24.403', 10, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (1885272, 'Пара Флэт Уайт', 363, 'NEW', 'Пара Флэт Уайт по выгодной цене. Доступно только в доставке!', '30, улица 8 Марта, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/61625d76-10a4-423b-9040-f9871c898a6b.png}', 1, 8, 9851099, '2025-12-02 11:14:00.643', '2025-12-16 09:53:24.405', 51, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (7718497, 'Катетер', 150, 'NEW', 'Просто катетер', 'г Оренбург, ул Севастопольская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5a2c2b79-558b-46e7-b267-6ae194c526b9.jpg}', 1, 15, 4146092, '2025-12-02 11:15:13.932', '2025-12-16 09:53:24.408', 74, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (3217337, 'Кофеин в таблетках', 160, 'NEW', 'Просто кофеин', 'г Оренбург, мкр Ростошинские пруды, Керченский пер', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/845037dc-1299-407d-a608-b0bfccd6de8d.webp}', 1, 3, 4146092, '2025-12-02 11:16:08.622', '2025-12-16 09:53:24.41', 22, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (5492285, 'Часы', 500, 'NEW', 'Часы громкоговорители', '113, Невельская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/b7eb5ecd-4fb4-4d07-ac74-b5146c8080ba.png}', 1, 5, 6038643, '2025-12-02 11:16:24.995', '2025-12-16 09:53:24.415', 35, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (1314227, 'Средство для удаления тейпов', 500, 'NEW', 'Средство для удаления тейпов', '41, улица Терешковой, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/9092c6b3-18bf-4960-b6c0-ae294784dd18.png}', 1, 5, 6038643, '2025-12-02 11:24:04.588', '2025-12-16 09:53:24.425', 28, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6883587, 'Шампунь Гарньер', 500, 'NEW', 'Просто шампунь', 'г Оренбург, ул Сумская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/72f30e85-f179-47a0-a982-11bf90113e0e.jpg}', 1, 3, 4146092, '2025-12-02 11:29:29.551', '2025-12-16 09:53:24.432', 17, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9956819, 'Тени для век', 2000, 'NEW', 'Просто тени', 'г Оренбург, ул Житомирская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3a20e814-cdaf-40a4-8252-fd825161c268.webp}', 1, 3, 4146092, '2025-12-02 11:31:50.394', '2025-12-16 09:53:24.442', 21, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9500725, 'Морозилки ларь Бирюса, Pozis, Kraft и другие', 15990, 'NEW', 'Бoльшой выбoр мopoзильныx камер (вepтикaльныe, лapи) разных oбъёмoв в нaличии в Орeнбуpге!

А так же в наличии огромный выбoр бытoвoй тexники по оптовым ценaм!

', 'Вишнёвая улица, СНТ "ЮЖНЫЙ УРАЛ ОФИЦЕРОВ ЗАПАСА И ОТСТАВКИ", Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/598ec012-3f41-4d02-9721-750964a49125.png}', 1, 10, 9851099, '2025-12-02 11:32:20.271', '2025-12-16 09:53:24.444', 61, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (7162519, 'Стиральная машина бу', 7000, 'USED', 'Стиральныe машины б.у. 🚛 Бecплaтная доставкa по гoроду ✅Гарaнтия до 12 меcяцeв пo чeку + пocлeгарантийнoe oбслуживаниe.

', '25, Краснознамённая улица, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/6f7d4e74-5748-4fa3-b1d3-a22f8aa6a061.png}', 1, 10, 9851099, '2025-12-02 11:33:16.603', '2025-12-16 09:53:24.449', 61, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (9042977, 'Закладки для учебников ', 300, 'NEW', 'Закладки для книг, «Книжная эстетика»', '5, улица Макаровой-Мутновой, Новостройка, Промышленный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/a2609c6c-a3c3-4d4a-bd93-60020e455210.png}', 1, 6, 6038643, '2025-12-02 11:33:30.034', '2025-12-16 09:53:24.452', 41, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (7566163, 'Духи ', 3500, 'NEW', 'Духи Dior Sauvage', 'г Оренбург, ул Черниговская', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/fa8257b2-6b96-4b79-b10e-df3d7c723129.jpg}', 1, 3, 4146092, '2025-12-02 11:27:20.06', '2025-12-16 09:53:24.438', 20, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (2568373, 'Концелярия', 700, 'NEW', 'Канцелярия для школы набор линеек y2k эстетика бант кролик', '128, Орская улица, Новостройка, Центральный район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/5e603622-4d00-4003-a685-cca9d4c78cf7.png}', 1, 6, 6038643, '2025-12-02 11:29:06.905', '2025-12-16 09:53:24.434', 41, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (3437684, 'Минский Бургер с курицей', 330, 'NEW', 'По-белорусски вкусный! Бургер с сочной куриной котлетой в хрустящей панировке, румяным картофельным оладушком, свежим салатом, двумя ломтиками нежного сыра, хрустящим ароматным беконом, маринованными огурчиками, нежным соусом «Сметана-укроп», и всё это — на воздушной горячей булочке с хрустящей крошкой.', '54, улица Кирова, Форштадт, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/3458ebaf-7b3f-4f39-b1f4-5a53322d9e64.png}', 1, 8, 9851099, '2025-12-02 11:12:20.477', '2025-12-16 09:53:24.21', 51, NULL, false, 'APPROVED', NULL);
INSERT INTO public."Product" (id, name, price, state, description, address, images, "categoryId", "subCategoryId", "userId", "createdAt", "updatedAt", "typeId", "videoUrl", "isHide", "moderateState", "moderationRejectionReason") VALUES (6218446, 'Соковыжималка caso CP 300 Pro', 4500, 'USED', 'CASO – нeмeцкaя торговая маркa бытовoй техники, принадлежащaя кoмпaнии Braukmann GmbH. Cоковыжималкa CASО СP 330 Prо предназначена для цитруcовыx cpeднего и крупногo pазмеpoв. Koрпуc прибоpа и cито для жмыxa выполнeны из нeржавeющeй cтaли. Автoмaтичeский старт плавнo зaпускает двигатель мощностью 160 Вт, функция «капля – стоп» обеспечивает чистоту рабочего места. В идеальном состоянии.

', '23/2, Пролетарская улица, Аренда, Ленинский район, Оренбург, городской округ Оренбург, Оренбургская область, Приволжский федеральный округ, 460000, Россия', '{https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/f2216887-329e-4d48-a605-e9b8bad18686.png}', 1, 10, 9851099, '2025-12-02 11:36:09.481', '2025-12-16 09:53:24.309', 61, NULL, false, 'MODERATE', NULL);


ALTER TABLE public."Product" ENABLE TRIGGER ALL;

--
-- Data for Name: Chat; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."Chat" DISABLE TRIGGER ALL;

INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (42, 9863001, 2681599, 6053931, 0, 0, NULL, '2025-12-01 08:12:52.805', '2025-12-01 08:12:52.805', '2025-12-01 08:12:52.805', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (1, 7659684, 6038643, 7391202, 0, 0, NULL, '2025-11-28 09:15:51.647', '2025-11-28 09:15:51.647', '2025-11-28 09:15:51.647', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (2, 9915250, 9371169, 8633592, 0, 0, NULL, '2025-11-28 09:16:07.289', '2025-11-28 09:16:07.289', '2025-11-28 09:16:07.289', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (3, 9915250, 8964288, 8633592, 0, 0, NULL, '2025-11-28 09:16:15.369', '2025-11-28 09:16:15.369', '2025-11-28 09:16:15.369', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (4, 9256863, 5966833, 8964288, 0, 0, NULL, '2025-11-28 09:16:50.74', '2025-11-28 09:16:50.74', '2025-11-28 09:16:50.74', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (5, 5611056, 6038643, 1208299, 0, 0, NULL, '2025-11-28 09:17:36.46', '2025-11-28 09:17:36.46', '2025-11-28 09:17:36.46', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (7, 1250840, 4761896, 7391202, 0, 0, NULL, '2025-11-28 09:20:00.074', '2025-11-28 09:20:00.074', '2025-11-28 09:20:00.074', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (141, 1250840, 7106521, 7391202, 0, 0, NULL, '2025-12-01 12:30:18.663', '2025-12-01 12:30:18.663', '2025-12-01 12:30:18.663', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (6, 9368305, 9371169, 8633592, 0, 0, NULL, '2025-11-28 09:18:57.647', '2025-11-28 09:18:57.647', '2025-11-28 09:18:57.647', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (9, 7384341, 7106521, 2321239, 0, 0, NULL, '2025-11-28 09:43:43.309', '2025-11-28 09:43:43.309', '2025-11-28 09:43:43.309', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (143, 7384341, 7132269, 2321239, 0, 0, NULL, '2025-12-01 14:46:54.623', '2025-12-01 14:46:54.623', '2025-12-01 14:46:54.623', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (144, 2105765, 7106521, 3235109, 0, 1, 1, '2025-12-02 06:39:40.026', '2025-12-02 06:32:54.893', '2025-12-04 06:34:08.592', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (142, 2105765, 7132269, 3235109, 0, 0, NULL, '2025-12-01 14:29:48.006', '2025-12-01 14:29:48.006', '2025-12-01 14:29:48.006', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (75, 2105765, 6669460, 3235109, 0, 0, NULL, '2025-12-01 08:28:55.846', '2025-12-01 08:28:55.846', '2025-12-01 08:28:55.846', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (108, 1122280, 7106521, 2681599, 0, 0, NULL, '2025-12-01 08:36:19.933', '2025-12-01 08:36:19.933', '2025-12-01 08:36:19.933', false);
INSERT INTO public."Chat" (id, "productId", "buyerId", "sellerId", "unreadCountBuyer", "unreadCountSeller", "lastMessageId", "lastMessageAt", "createdAt", "updatedAt", "isModerationChat") VALUES (145, 3437684, 4146092, 9851099, 0, 1, 2, '2025-12-02 11:33:33.784', '2025-12-02 11:33:25.608', '2025-12-02 11:47:19.434', false);


ALTER TABLE public."Chat" ENABLE TRIGGER ALL;

--
-- Data for Name: FavoriteAction; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."FavoriteAction" DISABLE TRIGGER ALL;

INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (5, 5966833, 3334788, '2025-11-28 09:21:05.722');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (7, 6038643, 3334788, '2025-11-28 09:21:12.339');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (6, 5966833, 1300264, '2025-11-28 09:21:07.97');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (3, 5966833, 1250840, '2025-11-28 09:18:22.746');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (77, 7106521, 9262881, '2025-12-02 07:46:33.084');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (4, 8633592, 9368305, '2025-11-28 09:18:35.259');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (8, 7106521, 7384341, '2025-11-29 09:00:09.809');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (74, 7132269, 2105765, '2025-12-01 14:29:45.996');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (41, 6669460, 2105765, '2025-12-01 09:22:11.203');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (76, 7106521, 2161612, '2025-12-02 07:34:18.185');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (75, 7106521, 4758351, '2025-12-02 07:34:17.402');
INSERT INTO public."FavoriteAction" (id, "userId", "productId", "addedAt") VALUES (78, 3235109, 8776759, '2025-12-03 00:00:00.214');


ALTER TABLE public."FavoriteAction" ENABLE TRIGGER ALL;

--
-- Data for Name: Log; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."Log" DISABLE TRIGGER ALL;

INSERT INTO public."Log" (id, "userId", action) VALUES (1, 1208299, 'Пополнение баланса: id: 1208299; email: kokeev.fil@mail.ru;\nсумма пополнения: 300; баланс: 500; бонусный баланс: 200');
INSERT INTO public."Log" (id, "userId", action) VALUES (2, 2287442, 'Пополнение баланса: id: 1208299; email: kokeev.fil@mail.ru;\nсумма пополнения: 300; баланс: 500; бонусный баланс: 300');


ALTER TABLE public."Log" ENABLE TRIGGER ALL;

--
-- Data for Name: Message; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."Message" DISABLE TRIGGER ALL;

INSERT INTO public."Message" (id, content, "senderId", "chatId", "isRead", "readAt", "createdAt", "updatedAt", "relatedProductId") VALUES (1, 'тест', 7106521, 144, false, NULL, '2025-12-02 06:39:40.018', '2025-12-02 06:39:40.018', NULL);
INSERT INTO public."Message" (id, content, "senderId", "chatId", "isRead", "readAt", "createdAt", "updatedAt", "relatedProductId") VALUES (2, 'Куда цену задрал? 200 край', 4146092, 145, false, NULL, '2025-12-02 11:33:33.781', '2025-12-02 11:33:33.781', NULL);


ALTER TABLE public."Message" ENABLE TRIGGER ALL;

--
-- Data for Name: Payment; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."Payment" DISABLE TRIGGER ALL;

INSERT INTO public."Payment" (id, "orderId", "paymentId", "userId", amount, status, "paymentUrl", "createdAt", "updatedAt") VALUES (1, '7106521-1766690209318', '7629971919', 7106521, 1000, 'PENDING', 'https://pay.tbank.ru/ahpkMYdA', '2025-12-25 19:16:49.644', '2025-12-25 19:16:49.644');
INSERT INTO public."Payment" (id, "orderId", "paymentId", "userId", amount, status, "paymentUrl", "createdAt", "updatedAt") VALUES (2, '7106521-1766690350760', '7629983326', 7106521, 1000, 'PENDING', 'https://pay.tbank.ru/mYjJpEkG', '2025-12-25 19:19:11.103', '2025-12-25 19:19:11.103');
INSERT INTO public."Payment" (id, "orderId", "paymentId", "userId", amount, status, "paymentUrl", "createdAt", "updatedAt") VALUES (3, '7106521-1766690451846', '7629991661', 7106521, 1000, 'PENDING', 'https://pay.tbank.ru/AqmD5LpC', '2025-12-25 19:20:52.116', '2025-12-25 19:20:52.116');
INSERT INTO public."Payment" (id, "orderId", "paymentId", "userId", amount, status, "paymentUrl", "createdAt", "updatedAt") VALUES (4, '7106521-1766690912537', '7630030365', 7106521, 1000, 'PENDING', 'https://pay.tbank.ru/4KnshkYJ', '2025-12-25 19:28:32.789', '2025-12-25 19:28:32.789');
INSERT INTO public."Payment" (id, "orderId", "paymentId", "userId", amount, status, "paymentUrl", "createdAt", "updatedAt") VALUES (5, '7106521-1766690970273', '7630035082', 7106521, 10, 'PENDING', 'https://pay.tbank.ru/w0hLiyV6', '2025-12-25 19:29:30.546', '2025-12-25 19:29:30.546');
INSERT INTO public."Payment" (id, "orderId", "paymentId", "userId", amount, status, "paymentUrl", "createdAt", "updatedAt") VALUES (6, '7106521-1766737026277', '7633516030', 7106521, 1000, 'PENDING', 'https://pay.tbank.ru/Y42XjEyF', '2025-12-26 08:17:06.881', '2025-12-26 08:17:06.881');
INSERT INTO public."Payment" (id, "orderId", "paymentId", "userId", amount, status, "paymentUrl", "createdAt", "updatedAt") VALUES (7, '7106521-1766992423382', '7655497017', 7106521, 1, 'PENDING', 'https://pay.tbank.ru/VhEBQwOm', '2025-12-29 07:13:44.223', '2025-12-29 07:13:44.223');


ALTER TABLE public."Payment" ENABLE TRIGGER ALL;

--
-- Data for Name: PhoneNumberView; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."PhoneNumberView" DISABLE TRIGGER ALL;



ALTER TABLE public."PhoneNumberView" ENABLE TRIGGER ALL;

--
-- Data for Name: TypeField; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."TypeField" DISABLE TRIGGER ALL;

INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (1, 'Размер', false, 1);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (2, 'Цвет', false, 1);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (3, 'Материал', false, 1);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (4, 'Бренд', false, 1);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (5, 'Название', false, 1);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (6, 'Вид', false, 1);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (7, 'Размер', false, 2);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (8, 'Цвет', false, 2);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (9, 'Материал', false, 2);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (10, 'Бренд', false, 2);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (11, 'Название', false, 2);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (12, 'Вид', false, 2);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (13, 'Размер', false, 3);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (14, 'Цвет', false, 3);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (15, 'Материал', false, 3);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (16, 'Бренд', false, 3);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (17, 'Название', false, 3);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (18, 'Вид', false, 3);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (19, 'Размер', false, 4);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (20, 'Цвет', false, 4);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (21, 'Материал', false, 4);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (22, 'Бренд', false, 4);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (23, 'Название', false, 4);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (24, 'Вид', false, 4);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (25, 'Размер', false, 5);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (26, 'Цвет', false, 5);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (27, 'Материал', false, 5);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (28, 'Бренд', false, 5);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (29, 'Название', false, 5);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (30, 'Вид', false, 5);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (31, 'Размер', false, 6);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (32, 'Цвет', false, 6);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (33, 'Материал', false, 6);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (34, 'Бренд', false, 6);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (35, 'Название', false, 6);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (36, 'Вид', false, 6);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (37, 'Размер', false, 7);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (38, 'Цвет', false, 7);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (39, 'Материал', false, 7);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (40, 'Бренд', false, 7);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (41, 'Название', false, 7);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (42, 'Вид', false, 7);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (43, 'Размер', false, 8);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (44, 'Цвет', false, 8);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (45, 'Материал', false, 8);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (46, 'Бренд', false, 8);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (47, 'Название', false, 8);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (48, 'Вид', false, 8);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (49, 'Цвет', false, 15);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (50, 'Размер', false, 13);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (51, 'Возраст', false, 9);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (52, 'Габариты', false, 12);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (53, 'Габариты', false, 10);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (54, 'Возраст', false, 14);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (55, 'Размер', false, 15);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (56, 'Габариты', false, 11);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (57, 'Цвет', false, 13);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (58, 'Цвет', false, 14);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (59, 'Размер', false, 9);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (60, 'Возраст', false, 13);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (61, 'Возраст', false, 15);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (62, 'Размер', false, 14);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (63, 'Цвет', false, 9);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (64, 'Габариты', false, 9);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (65, 'Возраст', false, 10);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (66, 'Возраст', false, 12);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (67, 'Габариты', false, 14);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (68, 'Возраст', false, 11);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (69, 'Цвет', false, 11);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (70, 'Габариты', false, 13);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (71, 'Размер', false, 12);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (72, 'Размер', false, 10);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (73, 'Габариты', false, 15);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (74, 'Цвет', false, 12);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (75, 'Цвет', false, 10);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (76, 'Размер', false, 11);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (77, 'Вид', false, 20);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (78, 'Вид', false, 17);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (79, 'Вид', false, 18);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (80, 'Вид', false, 21);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (81, 'Вид', false, 19);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (82, 'Вид', false, 16);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (83, 'Вид', false, 22);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (84, 'Цвет', false, 73);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (85, 'Наличие сертификата', false, 74);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (86, 'Цвет', false, 75);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (87, 'Портативность', false, 73);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (88, 'Бренд', false, 74);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (89, 'Портативность', false, 75);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (90, 'Бренд', false, 73);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (91, 'Бренд', false, 75);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (92, 'Портативность', false, 74);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (93, 'Наличие сертификата', false, 75);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (94, 'Цвет', false, 74);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (95, 'Наличие сертификата', false, 73);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (96, 'Тип питания', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (97, 'Диапазон измерений', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (98, 'Бренд', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (99, 'Вид', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (100, 'Комплектация', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (101, 'Замеры аритмии', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (102, 'Индикаторы', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (103, 'Точность измерений', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (104, 'Производитель', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (105, 'Метод измерения', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (106, 'Память', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (107, 'Тип', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (108, 'Калибровка', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (109, 'Объем капли', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (110, 'Погрешность', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (111, 'Гибкость', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (112, 'Размер', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (113, 'Время измерения', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (114, 'Функции маркировки', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (115, 'Подсветка', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (116, 'Звуковой сигнал', false, 26);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (117, 'Ребра жесткости', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (118, 'Вид', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (119, 'Конструктивные особенности', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (120, 'Область применения', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (121, 'Производитель', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (122, 'Степень фиксации', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (123, 'Гипоаллергенность', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (124, 'Назначение', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (125, 'Затяжки', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (126, 'Цвет', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (127, 'Размер', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (128, 'Шнурки', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (129, 'Возрастная группа', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (130, 'Материал', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (131, 'Пол', false, 27);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (132, 'Тип', false, 28);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (133, 'Вид', false, 28);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (134, 'Срок годности', false, 28);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (135, 'Производитель', false, 28);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (136, 'Тип', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (137, 'Вид', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (138, 'Материал рамы', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (139, 'Вес', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (140, 'Грузоподъёмность', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (141, 'Колёса', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (142, 'Аккумулятор', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (143, 'Управление', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (144, 'Доп функции', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (145, 'Складная конструкция', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (146, 'Цвет', false, 29);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (147, 'Материалы', false, 30);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (148, 'Гипоаллергенность', false, 30);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (149, 'Степень компрессии', false, 30);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (150, 'Размер', false, 30);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (151, 'Цвет', false, 30);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (152, 'Защитные свойства', false, 30);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (153, 'Доп функции', false, 30);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (154, 'Производитель', false, 30);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (155, 'Тип', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (156, 'Размер', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (157, 'Впитываемость', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (158, 'Материал впитывающего слоя', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (159, 'Материал внешнего слоя', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (160, 'Материал внутреннего слоя', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (161, 'Вид', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (162, 'Возраст', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (163, 'Доп свойства', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (164, 'Цвет', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (165, 'Производитель', false, 31);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (166, 'Вид', false, 32);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (167, 'Материал', false, 32);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (168, 'Тип', false, 32);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (169, 'Размер', false, 32);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (170, 'Доп функции', false, 32);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (171, 'Срок годности', false, 32);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (172, 'Производитель', false, 32);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (173, 'Тип', false, 33);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (174, 'Вид стомы', false, 33);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (175, 'Размер', false, 33);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (176, 'Производитель', false, 33);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (177, 'Тип', false, 34);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (178, 'Материал рамы', false, 34);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (179, 'Материал сиденья и спинки', false, 34);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (180, 'Регулировка высоты сидений', false, 34);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (181, 'Регулировка высоты и положения подлокотников', false, 34);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (182, 'Размер', false, 34);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (183, 'Доп опции', false, 34);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (184, 'Максимальная нагрузка', false, 34);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (185, 'Цвет', false, 34);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (186, 'Производитель', false, 34);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (187, 'Вид', false, 35);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (188, 'Производитель', false, 35);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (189, 'Характеристики устройства', false, 35);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (190, 'Габариты', false, 35);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (191, 'Вид', false, 36);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (192, 'Тип', false, 36);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (193, 'Материалы', false, 36);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (194, 'Объём мешков', false, 36);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (195, 'Диаметр пластин', false, 36);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (196, 'Производитель', false, 36);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (197, 'Наличие фильтров', false, 36);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (198, 'Наличие клапанов', false, 36);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (199, 'Наличие градуировки для измерения', false, 36);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (200, 'Тип', false, 37);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (201, 'Вид', false, 37);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (202, 'По поддержке', false, 37);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (203, 'Регулировка высоты', false, 37);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (204, 'Материал опор', false, 37);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (205, 'Вид наконечника', false, 37);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (206, 'Допустимая нагрузка', false, 37);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (207, 'Производитель', false, 37);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (208, 'Тип рукоятки', false, 37);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (209, 'Противоскользящий наконечник', false, 37);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (210, 'Тип', false, 38);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (211, 'Вид', false, 38);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (212, 'По поддержке', false, 38);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (213, 'Регулировка высоты', false, 38);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (214, 'Материал опор', false, 38);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (215, 'Производитель', false, 38);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (216, 'Назначение', false, 39);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (217, 'Тип', false, 39);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (218, 'Размер', false, 39);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (219, 'Материал', false, 39);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (220, 'Конструкция', false, 39);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (221, 'Тип', false, 40);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (222, 'Количество секций', false, 40);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (223, 'Регулировка', false, 40);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (224, 'Материал рамы', false, 40);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (225, 'Наличие/тип боковых ограждений', false, 40);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (226, 'Регулировка высоты', false, 40);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (227, 'Колёса', false, 40);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (228, 'Максимальная нагрузка', false, 40);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (229, 'Тип', false, 41);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (230, 'Тип', false, 42);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (231, 'Тип', false, 43);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (232, 'Тип', false, 44);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (233, 'Тип', false, 45);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (234, 'Тип', false, 46);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (235, 'Тип', false, 47);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (236, 'Материал', false, 48);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (237, 'Вид', false, 48);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (238, 'Производитель', false, 48);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (239, 'Материал', false, 49);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (240, 'Вид', false, 49);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (241, 'Производитель', false, 49);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (242, 'Материал', false, 50);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (243, 'Вид', false, 50);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (244, 'Производитель', false, 50);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (245, 'Срок годности', false, 51);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (246, 'Состав', false, 51);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (247, 'Способ обработки', false, 51);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (248, 'По способу употребления', false, 51);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (249, 'Объём', false, 51);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (250, 'Вес', false, 51);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (251, 'Срок годности', false, 52);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (252, 'Состав', false, 52);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (253, 'Способ обработки', false, 52);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (254, 'По способу употребления', false, 52);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (255, 'Объём', false, 52);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (256, 'Вес', false, 52);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (257, 'Срок годности', false, 53);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (258, 'Состав', false, 53);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (259, 'Способ обработки', false, 53);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (260, 'По способу употребления', false, 53);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (261, 'Объём', false, 53);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (262, 'Вес', false, 53);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (263, 'Возраст', false, 54);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (264, 'Потребности', false, 54);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (265, 'Вес', false, 54);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (266, 'Возраст', false, 55);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (267, 'Потребности', false, 55);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (268, 'Вес', false, 55);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (269, 'Возраст', false, 56);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (270, 'Потребности', false, 56);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (271, 'Вес', false, 56);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (272, 'Особенности ухода', false, 57);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (273, 'Тип корма/грунта', false, 57);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (274, 'Условия содержания', false, 57);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (275, 'Особенности ухода', false, 58);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (276, 'Тип корма/грунта', false, 58);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (277, 'Условия содержания', false, 58);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (278, 'Особенности ухода', false, 59);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (279, 'Тип корма/грунта', false, 59);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (280, 'Условия содержания', false, 59);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (281, 'Особенности ухода', false, 60);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (282, 'Тип корма/грунта', false, 60);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (283, 'Условия содержания', false, 60);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (284, 'Тип задачи', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (285, 'Классификация', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (286, 'Материал', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (287, 'Уровень шума', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (288, 'Уровень энергопотребления', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (289, 'Безопасность', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (290, 'Способ размещения', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (291, 'Габариты', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (292, 'Вес', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (293, 'Мощность', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (294, 'Производитель', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (295, 'Дополнительные функции', false, 61);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (296, 'Тип задачи', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (297, 'Классификация', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (298, 'Материал', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (299, 'Уровень шума', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (300, 'Уровень энергопотребления', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (301, 'Безопасность', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (302, 'Способ размещения', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (303, 'Габариты', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (304, 'Вес', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (305, 'Мощность', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (306, 'Производитель', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (307, 'Дополнительные функции', false, 62);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (308, 'Материал', false, 63);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (309, 'Форма', false, 63);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (310, 'Размер', false, 63);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (311, 'Комплектация', false, 63);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (312, 'Производитель', false, 63);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (313, 'Цвет', false, 63);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (314, 'Стиль', false, 63);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (315, 'Особенности', false, 63);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (316, 'Материал', false, 64);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (317, 'Форма', false, 64);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (318, 'Размер', false, 64);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (319, 'Комплектация', false, 64);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (320, 'Производитель', false, 64);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (321, 'Цвет', false, 64);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (322, 'Стиль', false, 64);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (323, 'Особенности', false, 64);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (324, 'Материал', false, 65);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (325, 'Форма', false, 65);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (326, 'Размер', false, 65);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (327, 'Комплектация', false, 65);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (328, 'Производитель', false, 65);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (329, 'Цвет', false, 65);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (330, 'Стиль', false, 65);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (331, 'Особенности', false, 65);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (332, 'Материал', false, 66);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (333, 'Форма', false, 66);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (334, 'Размер', false, 66);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (335, 'Комплектация', false, 66);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (336, 'Производитель', false, 66);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (337, 'Цвет', false, 66);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (338, 'Стиль', false, 66);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (339, 'Особенности', false, 66);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (340, 'Цвет', false, 67);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (341, 'Стиль', false, 67);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (342, 'Производитель', false, 67);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (343, 'Материал', false, 67);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (344, 'Конструкция', false, 67);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (345, 'Цвет', false, 68);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (346, 'Стиль', false, 68);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (347, 'Производитель', false, 68);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (348, 'Материал', false, 68);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (349, 'Конструкция', false, 68);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (350, 'Цвет', false, 69);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (351, 'Стиль', false, 69);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (352, 'Производитель', false, 69);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (353, 'Материал', false, 69);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (354, 'Конструкция', false, 69);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (355, 'Цвет', false, 70);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (356, 'Стиль', false, 70);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (357, 'Производитель', false, 70);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (358, 'Материал', false, 70);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (359, 'Конструкция', false, 70);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (360, 'Цвет', false, 71);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (361, 'Стиль', false, 71);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (362, 'Производитель', false, 71);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (363, 'Материал', false, 71);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (364, 'Конструкция', false, 71);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (365, 'Цвет', false, 72);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (366, 'Стиль', false, 72);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (367, 'Производитель', false, 72);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (368, 'Материал', false, 72);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (369, 'Конструкция', false, 72);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (370, 'Цвет', false, 73);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (371, 'Стиль', false, 73);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (372, 'Производитель', false, 73);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (373, 'Материал', false, 73);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (374, 'Конструкция', false, 73);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (375, 'Цвет', false, 74);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (376, 'Стиль', false, 74);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (377, 'Производитель', false, 74);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (378, 'Материал', false, 74);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (379, 'Конструкция', false, 74);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (380, 'Цвет', false, 75);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (381, 'Стиль', false, 75);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (382, 'Производитель', false, 75);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (383, 'Материал', false, 75);
INSERT INTO public."TypeField" (id, name, "isRequired", "typeId") VALUES (384, 'Конструкция', false, 75);


ALTER TABLE public."TypeField" ENABLE TRIGGER ALL;

--
-- Data for Name: ProductFieldValue; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."ProductFieldValue" DISABLE TRIGGER ALL;

INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (34, '50', 1, 9368305);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (35, 'Радужный', 2, 9368305);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (36, 'Хлопок', 3, 9368305);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (37, 'Demix', 4, 9368305);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (38, 'Аниме', 5, 9368305);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (39, 'Да', 6, 9368305);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (40, '15', 54, 3506516);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (41, 'разный', 58, 3506516);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (42, '3х3', 62, 3506516);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (43, 'не знаю', 67, 3506516);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (48, '1-5', 54, 4215912);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (49, 'Книжный', 58, 4215912);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (50, 'Книжный', 62, 4215912);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (51, 'Книжные', 67, 4215912);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (52, 'черный', 86, 1122280);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (85, 'Активная модель 1000', 136, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (86, ', прогулочная', 137, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (87, 'Железяка', 138, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (91, 'нет', 142, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (118, 's', 7, 5510664);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (119, 'голубой', 8, 5510664);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (120, 'джинса', 9, 5510664);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (121, 'Gloria', 10, 5510664);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (122, 'джинсы', 11, 5510664);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (123, 'багги', 12, 5510664);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (151, 'XS-L', 7, 4523969);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (152, 'разные', 8, 4523969);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (153, 'хлопок', 9, 4523969);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (154, 'Dasha', 10, 4523969);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (155, 'DV', 11, 4523969);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (156, 'летние', 12, 4523969);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (157, 'S-L', 7, 5902819);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (158, 'белый', 8, 5902819);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (159, 'норка', 9, 5902819);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (160, 'red', 10, 5902819);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (161, 'Sweet', 11, 5902819);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (162, 'свитер', 12, 5902819);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (163, 'S-L', 7, 9122333);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (164, 'бежевый', 8, 9122333);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (165, 'хлопок', 9, 9122333);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (166, 'Bant', 10, 9122333);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (167, 'BD', 11, 9122333);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (168, 'зипка', 12, 9122333);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (169, 's', 7, 3244052);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (170, 'голубо-белый', 8, 3244052);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (171, 'джинса', 9, 3244052);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (172, 'VK', 10, 3244052);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (173, 'vk', 11, 3244052);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (174, 'дубленка', 12, 3244052);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (175, 'Белый', 84, 6497808);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (176, 'Да', 87, 6497808);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (177, 'Omron', 90, 6497808);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (178, 'Да', 95, 6497808);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (179, '2-3', 51, 9380113);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (180, '5', 59, 9380113);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (181, 'разные', 63, 9380113);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (182, '15-30 см', 64, 9380113);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (183, 'Белый', 84, 3563632);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (184, 'Да', 87, 3563632);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (185, 'Omron', 90, 3563632);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (186, 'Да', 95, 3563632);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (187, '50см', 50, 2865910);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (188, 'грязно-синий', 57, 2865910);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (189, '1-3', 60, 2865910);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (190, '50см', 70, 2865910);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (191, 'Да', 85, 8948419);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (192, 'Zenet', 88, 8948419);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (193, 'Нет', 92, 8948419);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (194, 'Синий', 94, 8948419);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (195, 'белый', 49, 9783545);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (196, '100см', 55, 9783545);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (197, '0-1', 61, 9783545);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (198, '200см', 73, 9783545);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (199, '100см', 53, 6734788);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (200, '0-1', 65, 6734788);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (201, '100см', 72, 6734788);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (202, 'белый', 75, 6734788);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (203, 'Да', 85, 9305563);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (204, 'KMED', 88, 9305563);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (205, 'Да', 92, 9305563);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (206, 'Белый', 94, 9305563);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (207, '200см', 53, 9265239);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (208, '0-1', 65, 9265239);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (209, '200см', 72, 9265239);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (210, 'бело-коричневый', 75, 9265239);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (211, '1,5 метра', 53, 6300121);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (212, '0-7', 65, 6300121);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (213, '1,5 метра', 72, 6300121);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (214, 'белый', 75, 6300121);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (215, 'Да', 85, 7718497);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (216, 'KMD', 88, 7718497);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (217, 'Да', 92, 7718497);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (218, 'Белый', 94, 7718497);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (219, 'Таблетки', 83, 3217337);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (220, 'часы', 187, 5492285);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (221, 'америка', 188, 5492285);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (222, 'часы', 189, 5492285);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (223, '60см', 190, 5492285);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (224, 'Крем', 82, 8507601);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (225, 'Компрессорный ингалятор', 187, 4372887);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (226, 'Omron Comp Air NE-C300 Complete', 188, 4372887);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (227, 'Небулайзер OMRON C300 Complete — прибор, работающий в 3 режимах ингаляции. ', 189, 4372887);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (228, '70см', 190, 4372887);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (229, 'кресло-коляска', 136, 6617171);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (230, 'сидячий', 137, 6617171);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (231, 'метал', 138, 6617171);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (232, '7кг', 139, 6617171);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (233, '100кг', 140, 6617171);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (235, 'нет', 142, 6617171);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (240, 'clinar', 132, 1314227);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (241, 'балончик', 133, 1314227);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (242, '2 года', 134, 1314227);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (243, 'америка', 135, 1314227);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (244, 'уход', 132, 2273041);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (245, 'спрей', 133, 2273041);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (246, '5 лет', 134, 2273041);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (247, 'dinax', 135, 2273041);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (248, 'Духи', 77, 7566163);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (249, 'Шампунь', 78, 6883587);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (250, 'Палетка с тенями', 80, 9956819);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (251, 'папва', 117, 5868178);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (252, 'папва', 118, 5868178);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (253, 'апавпва', 119, 5868178);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (254, 'ир', 120, 5868178);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (255, 'рпрпр', 121, 5868178);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (256, 'пнпнп', 122, 5868178);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (28, '52', 1, 9915250);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (29, 'Черный', 2, 9915250);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (30, 'Шелк', 3, 9915250);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (31, 'Gucci', 4, 9915250);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (32, 'Платье', 5, 9915250);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (33, 'Горничная', 6, 9915250);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (88, '16 кг 500 грамм', 139, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (89, '120 кгили 0,12 тонны, или 1,2 центнера', 140, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (90, 'низкопрофильные', 141, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (92, 'ручное', 143, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (93, 'мягкая сидушка', 144, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (94, 'да', 145, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (95, 'черный', 146, 4758351);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (234, ' Колесная база, не выступающая за габариты коляски', 141, 6617171);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (236, 'автоматическое', 143, 6617171);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (237, 'нет', 144, 6617171);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (238, 'есть', 145, 6617171);
INSERT INTO public."ProductFieldValue" (id, value, "fieldId", "productId") VALUES (239, 'черный', 146, 6617171);


ALTER TABLE public."ProductFieldValue" ENABLE TRIGGER ALL;

--
-- Data for Name: Promotion; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."Promotion" DISABLE TRIGGER ALL;

INSERT INTO public."Promotion" (id, name, "pricePerDay", "createdAt", "updatedAt") VALUES (2, 'Люкс', 100, '2025-12-08 12:37:51.475', '2025-12-08 12:37:44.761');
INSERT INTO public."Promotion" (id, name, "pricePerDay", "createdAt", "updatedAt") VALUES (1, 'Стандарт', 50, '2025-12-08 12:37:51.475', '2025-12-08 12:37:32.223');


ALTER TABLE public."Promotion" ENABLE TRIGGER ALL;

--
-- Data for Name: ProductPromotion; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."ProductPromotion" DISABLE TRIGGER ALL;

INSERT INTO public."ProductPromotion" (id, "productId", "promotionId", "userId", days, "totalPrice", "startDate", "endDate", "isActive", "isPaid", "createdAt", "updatedAt") VALUES (1, 5868178, 2, 7106521, 7, 700, '2026-01-22 06:14:35.62', '2026-01-29 06:14:35.62', true, true, '2026-01-22 06:14:35.628', '2026-01-22 06:14:35.672');


ALTER TABLE public."ProductPromotion" ENABLE TRIGGER ALL;

--
-- Data for Name: ProductView; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."ProductView" DISABLE TRIGGER ALL;

INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (29, 5966833, 4081087, '2025-11-28 09:19:16.228');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (23, 6038643, 4081087, '2025-11-28 09:18:17.402');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (28, 8633592, 4081087, '2025-11-28 09:19:28.543');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (19, 6038643, 5231119, '2025-11-28 09:17:47.718');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (106, 8964288, 5231119, '2025-11-28 09:50:04.663');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (21, 6038643, 6901799, '2025-11-28 09:18:07.082');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (54, 8964288, 6901799, '2025-11-28 09:21:24.529');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (20, 6038643, 1979749, '2025-11-28 09:18:00.208');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (105, 8964288, 1979749, '2025-11-28 09:49:35.399');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (7, 6038643, 1300264, '2025-11-28 09:18:13.977');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (602, 7132269, 1300264, '2025-12-01 14:29:41.095');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (270, 2681599, 9863001, '2025-12-01 08:12:43.181');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (42, 2321239, 4224343, '2025-11-28 09:20:52.446');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (504, 7106521, 4224343, '2025-12-02 07:39:03.095');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (11, 6038643, 7659684, '2025-11-28 09:15:50.575');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (305, 7106521, 7659684, '2025-12-01 12:03:10.898');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (6, 6038643, 1970246, '2025-11-28 09:15:46.542');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (204, 7106521, 1970246, '2025-12-01 09:26:35.713');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (70, 8964288, 1970246, '2025-11-28 09:22:21.23');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (238, 7106521, 3207807, '2025-12-01 09:26:48.769');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (51, 8964288, 3207807, '2025-11-28 09:21:19.648');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (5, 6038643, 9915250, '2025-11-28 09:15:40.734');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (35, 2321239, 9915250, '2025-11-28 09:20:26.159');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (13, 9371169, 9915250, '2025-11-28 09:16:05.942');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (173, 7106521, 9915250, '2025-12-01 09:26:28.453');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (12, 8964288, 9915250, '2025-11-28 09:16:05.662');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (15, 5966833, 9256863, '2025-11-28 09:16:48.332');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (8, 6038643, 9256863, '2025-11-28 09:15:34.351');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (40, 2321239, 9256863, '2025-11-28 09:20:48.02');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (569, 7106521, 9256863, '2025-12-01 09:26:25.118');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (14, 5966833, 5611056, '2025-11-28 09:16:42.097');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (18, 6038643, 5611056, '2025-11-28 09:17:21.48');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (239, 7106521, 5611056, '2025-12-01 07:48:19.685');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (43, 8964288, 5611056, '2025-11-28 09:20:59.678');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (17, 6038643, 5609249, '2025-11-28 09:17:15.39');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (39, 2321239, 5609249, '2025-11-28 09:20:43.139');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (49, 8633592, 5609249, '2025-11-28 09:21:16.127');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (536, 7106521, 5609249, '2025-12-01 09:24:23.12');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (16, 6038643, 1250840, '2025-11-28 09:17:09.157');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (31, 4761896, 1250840, '2025-11-28 09:19:39.18');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (48, 8633592, 1250840, '2025-11-28 09:21:12.953');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (172, 7106521, 1250840, '2025-12-03 16:57:03.686');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (46, 8964288, 1250840, '2025-11-28 09:21:09.011');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (25, 6038643, 9262881, '2025-11-28 09:18:35.85');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (44, 8633592, 9262881, '2025-11-28 09:21:04.291');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (50, 7106521, 9262881, '2025-12-02 07:46:38.367');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (24, 6038643, 9368305, '2025-11-28 09:18:26.145');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (27, 2321239, 9368305, '2025-11-28 09:20:38.981');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (26, 9371169, 9368305, '2025-11-28 09:18:54.527');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (206, 7106521, 9368305, '2025-12-01 09:26:15.811');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (47, 8633592, 7384341, '2025-11-28 09:21:09.232');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (72, 7106521, 7384341, '2025-12-03 16:55:51.09');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (605, 7132269, 7384341, '2025-12-01 14:46:32.887');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (32, 6038643, 3506516, '2025-11-28 09:20:15.947');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (45, 8633592, 3506516, '2025-11-28 09:21:06.787');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (33, 7106521, 3506516, '2025-11-28 09:20:19.239');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (604, 7132269, 3506516, '2025-12-01 14:37:03.378');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (52, 8633592, 4215912, '2025-11-28 09:21:21.713');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (237, 7106521, 4215912, '2025-12-01 09:24:12.715');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (71, 8964288, 4215912, '2025-11-28 09:22:26.679');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (606, 7132269, 4215912, '2025-12-01 14:54:23.873');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (205, 7106521, 2105765, '2025-12-02 06:32:52.922');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (603, 7132269, 2105765, '2025-12-01 14:29:43.789');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (336, 6669460, 2105765, '2025-12-01 08:37:06.268');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (369, 7106521, 1122280, '2025-12-01 08:48:56.032');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (468, 7106521, 2161612, '2025-12-01 09:22:41.051');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (618, 9851099, 2161612, '2025-12-02 11:42:27.288');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (567, 7106521, 4758351, '2025-12-02 07:32:36.871');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (656, 7106521, 4267180, '2025-12-03 09:52:39.77');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (660, 7106521, 9380113, '2025-12-03 16:41:06.056');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (652, 7106521, 9783545, '2025-12-03 09:51:25.113');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (654, 7106521, 3437684, '2025-12-03 09:52:30.64');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (617, 4146092, 3437684, '2025-12-02 11:33:23.238');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (662, 7106521, 1512888, '2025-12-03 16:54:39.706');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (653, 7106521, 6628130, '2025-12-03 16:38:43.192');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (659, 7106521, 2388612, '2025-12-03 17:01:43.056');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (651, 7106521, 9956819, '2025-12-03 09:51:19.021');
INSERT INTO public."ProductView" (id, "viewedById", "productId", "viewedAt") VALUES (658, 7106521, 6218446, '2025-12-03 16:38:47.965');


ALTER TABLE public."ProductView" ENABLE TRIGGER ALL;

--
-- Data for Name: Review; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."Review" DISABLE TRIGGER ALL;



ALTER TABLE public."Review" ENABLE TRIGGER ALL;

--
-- Data for Name: SupportTicket; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."SupportTicket" DISABLE TRIGGER ALL;



ALTER TABLE public."SupportTicket" ENABLE TRIGGER ALL;

--
-- Data for Name: SupportMessage; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."SupportMessage" DISABLE TRIGGER ALL;



ALTER TABLE public."SupportMessage" ENABLE TRIGGER ALL;

--
-- Data for Name: _UserFavorites; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public."_UserFavorites" DISABLE TRIGGER ALL;



ALTER TABLE public."_UserFavorites" ENABLE TRIGGER ALL;

--
-- Data for Name: _prisma_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

ALTER TABLE public._prisma_migrations DISABLE TRIGGER ALL;

INSERT INTO public._prisma_migrations (id, checksum, finished_at, migration_name, logs, rolled_back_at, started_at, applied_steps_count) VALUES ('734e28f4-7c7a-4cc7-a17c-83f7a5d4017d', '7ff5b7017fce2e76bc78544000c57780ed8c77b4d7428fc010d39f5b75cb7b8a', '2026-01-16 15:03:21.231382+02', '20260114212748_init', '', NULL, '2026-01-16 15:03:21.231382+02', 0);
INSERT INTO public._prisma_migrations (id, checksum, finished_at, migration_name, logs, rolled_back_at, started_at, applied_steps_count) VALUES ('8f4df8a4-6dd0-45d1-a23d-ecf62da5e81c', 'd73539483e50e92bc3a4f78cdf878ec785ea549ec07fc3f4d9795a0577f66d0a', '2026-01-16 15:04:20.445351+02', '20260116094500_update_banner', '', NULL, '2026-01-16 15:04:20.445351+02', 0);


ALTER TABLE public._prisma_migrations ENABLE TRIGGER ALL;

--
-- Name: BannerView_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."BannerView_id_seq"', 1, false);


--
-- Name: Banner_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Banner_id_seq"', 4, true);


--
-- Name: Category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Category_id_seq"', 2, true);


--
-- Name: Chat_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Chat_id_seq"', 145, true);


--
-- Name: FavoriteAction_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."FavoriteAction_id_seq"', 79, true);


--
-- Name: Log_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Log_id_seq"', 2, true);


--
-- Name: Message_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Message_id_seq"', 2, true);


--
-- Name: Payment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Payment_id_seq"', 7, true);


--
-- Name: PhoneNumberView_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."PhoneNumberView_id_seq"', 1, true);


--
-- Name: ProductFieldValue_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."ProductFieldValue_id_seq"', 256, true);


--
-- Name: ProductPromotion_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."ProductPromotion_id_seq"', 1, true);


--
-- Name: ProductView_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."ProductView_id_seq"', 667, true);


--
-- Name: Promotion_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Promotion_id_seq"', 2, true);


--
-- Name: Review_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Review_id_seq"', 2, true);


--
-- Name: Role_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Role_id_seq"', 3, true);


--
-- Name: SubCategory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."SubCategory_id_seq"', 15, true);


--
-- Name: SubcategotyType_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."SubcategotyType_id_seq"', 75, true);


--
-- Name: SupportMessage_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."SupportMessage_id_seq"', 1, false);


--
-- Name: SupportTicket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."SupportTicket_id_seq"', 1, false);


--
-- Name: TypeField_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."TypeField_id_seq"', 209, true);


--
-- PostgreSQL database dump complete
--

\unrestrict 4NiEYDyKlxvlgGqkeb5HvaMLHy97mCg56PzipZ49km854s9NGzFsWS5JEI3INTO

