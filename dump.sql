--
-- PostgreSQL database dump
--

-- Dumped from database version 10.5 (Debian 10.5-1.pgdg90+1)
-- Dumped by pg_dump version 10.5 (Debian 10.5-1.pgdg90+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE ONLY public.products DROP CONSTRAINT products_store_id_fkey;
ALTER TABLE ONLY public.orders DROP CONSTRAINT orders_store_id_fkey;
ALTER TABLE ONLY public.items DROP CONSTRAINT items_store_id_fkey;
ALTER TABLE ONLY public.items DROP CONSTRAINT items_product_id_fkey;
ALTER TABLE ONLY public.items DROP CONSTRAINT items_order_id_fkey;
ALTER TABLE ONLY public.stores DROP CONSTRAINT stores_pkey;
ALTER TABLE ONLY public.products DROP CONSTRAINT products_pkey;
ALTER TABLE ONLY public.orders DROP CONSTRAINT orders_pkey;
ALTER TABLE ONLY public.items DROP CONSTRAINT items_pkey;
ALTER TABLE public.stores ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.products ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.orders ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.items ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE public.stores_id_seq;
DROP TABLE public.stores;
DROP SEQUENCE public.products_id_seq;
DROP TABLE public.products;
DROP SEQUENCE public.orders_id_seq;
DROP TABLE public.orders;
DROP SEQUENCE public.items_id_seq;
DROP TABLE public.items;
DROP EXTENSION plpgsql;
DROP SCHEMA public;
--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres-dev
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO "postgres-dev";

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres-dev
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: items; Type: TABLE; Schema: public; Owner: postgres-dev
--

CREATE TABLE public.items (
    id integer NOT NULL,
    product_id integer,
    order_id integer,
    store_id integer NOT NULL
);


ALTER TABLE public.items OWNER TO "postgres-dev";

--
-- Name: items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres-dev
--

CREATE SEQUENCE public.items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.items_id_seq OWNER TO "postgres-dev";

--
-- Name: items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres-dev
--

ALTER SEQUENCE public.items_id_seq OWNED BY public.items.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres-dev
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    total real NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    store_id integer NOT NULL
);


ALTER TABLE public.orders OWNER TO "postgres-dev";

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres-dev
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.orders_id_seq OWNER TO "postgres-dev";

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres-dev
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres-dev
--

CREATE TABLE public.products (
    id integer NOT NULL,
    name text NOT NULL,
    price real NOT NULL,
    store_id integer NOT NULL
);


ALTER TABLE public.products OWNER TO "postgres-dev";

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres-dev
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.products_id_seq OWNER TO "postgres-dev";

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres-dev
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: stores; Type: TABLE; Schema: public; Owner: postgres-dev
--

CREATE TABLE public.stores (
    id integer NOT NULL,
    name text NOT NULL,
    description text
);


ALTER TABLE public.stores OWNER TO "postgres-dev";

--
-- Name: stores_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres-dev
--

CREATE SEQUENCE public.stores_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.stores_id_seq OWNER TO "postgres-dev";

--
-- Name: stores_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres-dev
--

ALTER SEQUENCE public.stores_id_seq OWNED BY public.stores.id;


--
-- Name: items id; Type: DEFAULT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.items ALTER COLUMN id SET DEFAULT nextval('public.items_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: stores id; Type: DEFAULT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.stores ALTER COLUMN id SET DEFAULT nextval('public.stores_id_seq'::regclass);


--
-- Data for Name: items; Type: TABLE DATA; Schema: public; Owner: postgres-dev
--

COPY public.items (id, product_id, order_id, store_id) FROM stdin;
109	69	37	1
110	68	\N	1
111	68	\N	1
112	68	37	1
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres-dev
--

COPY public.orders (id, total, created, store_id) FROM stdin;
37	2500	2018-09-22 20:02:25.568247	1
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres-dev
--

COPY public.products (id, name, price, store_id) FROM stdin;
68	iPhone	1000	1
69	iPad	1500	1
\.


--
-- Data for Name: stores; Type: TABLE DATA; Schema: public; Owner: postgres-dev
--

COPY public.stores (id, name, description) FROM stdin;
1	My very cool store	Descriptive description
\.


--
-- Name: items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres-dev
--

SELECT pg_catalog.setval('public.items_id_seq', 141, true);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres-dev
--

SELECT pg_catalog.setval('public.orders_id_seq', 69, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres-dev
--

SELECT pg_catalog.setval('public.products_id_seq', 101, true);


--
-- Name: stores_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres-dev
--

SELECT pg_catalog.setval('public.stores_id_seq', 34, true);


--
-- Name: items items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: stores stores_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT stores_pkey PRIMARY KEY (id);


--
-- Name: items items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: items items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: items items_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.stores(id);


--
-- Name: orders orders_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.stores(id);


--
-- Name: products products_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.stores(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres-dev
--

GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

