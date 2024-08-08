--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2
-- Dumped by pg_dump version 16.2

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
-- Name: entries; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.entries (
    entry_id bigint NOT NULL,
    product_name character varying NOT NULL,
    product_price integer NOT NULL,
    quantity_added integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.entries OWNER TO root;

--
-- Name: entries_entry_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.entries_entry_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.entries_entry_id_seq OWNER TO root;

--
-- Name: entries_entry_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.entries_entry_id_seq OWNED BY public.entries.entry_id;


--
-- Name: invoices; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.invoices (
    invoice_id bigint NOT NULL,
    invoice_number character varying NOT NULL,
    user_invoice_id integer NOT NULL,
    user_invoice_username character varying NOT NULL,
    invoice_data jsonb NOT NULL,
    invoice_pdf bytea NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    invoice_date timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.invoices OWNER TO root;

--
-- Name: invoices_invoice_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.invoices_invoice_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.invoices_invoice_id_seq OWNER TO root;

--
-- Name: invoices_invoice_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.invoices_invoice_id_seq OWNED BY public.invoices.invoice_id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.products (
    product_id bigint NOT NULL,
    product_name character varying NOT NULL,
    unit_price integer NOT NULL,
    packsize character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.products OWNER TO root;

--
-- Name: products_product_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.products_product_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_product_id_seq OWNER TO root;

--
-- Name: products_product_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.products_product_id_seq OWNED BY public.products.product_id;


--
-- Name: receipts; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.receipts (
    receipt_id bigint NOT NULL,
    receipt_number character varying NOT NULL,
    user_receipt_id integer NOT NULL,
    user_receipt_username character varying NOT NULL,
    receipt_data json,
    receipt_pdf bytea NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.receipts OWNER TO root;

--
-- Name: receipts_receipt_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.receipts_receipt_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.receipts_receipt_id_seq OWNER TO root;

--
-- Name: receipts_receipt_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.receipts_receipt_id_seq OWNED BY public.receipts.receipt_id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO root;

--
-- Name: stock_value; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.stock_value (
    user_id integer NOT NULL,
    value bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public.stock_value OWNER TO root;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.transactions (
    transaction_id character varying NOT NULL,
    transaction_user_id integer NOT NULL,
    amount integer NOT NULL,
    status boolean DEFAULT false NOT NULL,
    data_sold json,
    phone_number character varying DEFAULT '00'::character varying NOT NULL,
    mpesa_receipt_number character varying DEFAULT 'No Receipt Number'::character varying NOT NULL,
    result_description character varying DEFAULT 'Mpesa Not Called Back Description'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.transactions OWNER TO root;

--
-- Name: users; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.users (
    user_id bigint NOT NULL,
    username character varying NOT NULL,
    password character varying NOT NULL,
    email character varying NOT NULL,
    phone_number character varying NOT NULL,
    address character varying NOT NULL,
    stock json,
    role character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.users OWNER TO root;

--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.users_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_user_id_seq OWNER TO root;

--
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.users_user_id_seq OWNED BY public.users.user_id;


--
-- Name: entries entry_id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.entries ALTER COLUMN entry_id SET DEFAULT nextval('public.entries_entry_id_seq'::regclass);


--
-- Name: invoices invoice_id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.invoices ALTER COLUMN invoice_id SET DEFAULT nextval('public.invoices_invoice_id_seq'::regclass);


--
-- Name: products product_id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.products ALTER COLUMN product_id SET DEFAULT nextval('public.products_product_id_seq'::regclass);


--
-- Name: receipts receipt_id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.receipts ALTER COLUMN receipt_id SET DEFAULT nextval('public.receipts_receipt_id_seq'::regclass);


--
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_user_id_seq'::regclass);


--
-- Name: entries entries_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.entries
    ADD CONSTRAINT entries_pkey PRIMARY KEY (entry_id);


--
-- Name: invoices invoices_invoice_number_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_invoice_number_key UNIQUE (invoice_number);


--
-- Name: invoices invoices_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_pkey PRIMARY KEY (invoice_id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (product_id);


--
-- Name: products products_product_name_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_product_name_key UNIQUE (product_name);


--
-- Name: receipts receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_pkey PRIMARY KEY (receipt_id);


--
-- Name: receipts receipts_receipt_number_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_receipt_number_key UNIQUE (receipt_number);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (transaction_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_phone_number_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_number_key UNIQUE (phone_number);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: invoices_invoice_number_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX invoices_invoice_number_idx ON public.invoices USING btree (invoice_number);


--
-- Name: invoices_user_invoice_id_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX invoices_user_invoice_id_idx ON public.invoices USING btree (user_invoice_id);


--
-- Name: products_product_name_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX products_product_name_idx ON public.products USING btree (product_name);


--
-- Name: receipts_receipt_number_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX receipts_receipt_number_idx ON public.receipts USING btree (receipt_number);


--
-- Name: receipts_user_receipt_id_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX receipts_user_receipt_id_idx ON public.receipts USING btree (user_receipt_id);


--
-- Name: users_username_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX users_username_idx ON public.users USING btree (username);


--
-- Name: stock_value fk_stock_value_user_id; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.stock_value
    ADD CONSTRAINT fk_stock_value_user_id FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: transactions fk_transaction_user_id; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT fk_transaction_user_id FOREIGN KEY (transaction_user_id) REFERENCES public.users(user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: invoices fk_user_invoice_id; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT fk_user_invoice_id FOREIGN KEY (user_invoice_id) REFERENCES public.users(user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: invoices fk_user_invoice_username; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT fk_user_invoice_username FOREIGN KEY (user_invoice_username) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: receipts fk_user_receipt_id; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT fk_user_receipt_id FOREIGN KEY (user_receipt_id) REFERENCES public.users(user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: receipts fk_user_receipt_username; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT fk_user_receipt_username FOREIGN KEY (user_receipt_username) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: receipts receipts_receipt_number_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_receipt_number_fkey FOREIGN KEY (receipt_number) REFERENCES public.transactions(transaction_id);


--
-- PostgreSQL database dump complete
--

