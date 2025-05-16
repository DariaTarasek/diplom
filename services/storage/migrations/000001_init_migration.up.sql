--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

-- Started on 2025-05-16 17:25:08

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 231 (class 1259 OID 16497)
-- Name: administrators; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS administrators (
                                                     user_id integer NOT NULL,
                                                     first_name character varying(256) NOT NULL,
    second_name character varying(256) NOT NULL,
    surname character varying(256),
    phone_number character varying(11),
    email character varying(256) NOT NULL,
    gender character(1) NOT NULL
    );


ALTER TABLE administrators OWNER TO postgres;

--
-- TOC entry 236 (class 1259 OID 24706)
-- Name: appointments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS appointments (
                                                   id integer NOT NULL,
                                                   doctor_id integer NOT NULL,
                                                   date date NOT NULL,
                                                   "time" time without time zone NOT NULL,
                                                   patient_id integer,
                                                   second_name character varying(256) NOT NULL,
    first_name character varying(256) NOT NULL,
    surname character varying(256),
    birth_date date NOT NULL,
    gender character(1) NOT NULL,
    phone_number character varying(11) NOT NULL,
    status character varying(20) DEFAULT 'unconfirmed'::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT status CHECK (((status)::text = ANY ((ARRAY['unconfirmed'::character varying, 'confirmed'::character varying, 'completed'::character varying, 'cancelled'::character varying])::text[])))
    );


ALTER TABLE appointments OWNER TO postgres;

--
-- TOC entry 235 (class 1259 OID 24705)
-- Name: appointments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS appointments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE appointments_id_seq OWNER TO postgres;

--
-- TOC entry 4902 (class 0 OID 0)
-- Dependencies: 235
-- Name: appointments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE appointments_id_seq OWNED BY appointments.id;


--
-- TOC entry 240 (class 1259 OID 24741)
-- Name: daily_overrides; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS daily_overrides (
                                                      id integer NOT NULL,
                                                      doctor_id integer NOT NULL,
                                                      date date NOT NULL,
                                                      start_time time without time zone,
                                                      end_time time without time zone,
                                                      slot_duration_minutes integer,
                                                      is_day_off boolean DEFAULT false
);


ALTER TABLE daily_overrides OWNER TO postgres;

--
-- TOC entry 239 (class 1259 OID 24740)
-- Name: daily_overrides_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS daily_overrides_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE daily_overrides_id_seq OWNER TO postgres;

--
-- TOC entry 4903 (class 0 OID 0)
-- Dependencies: 239
-- Name: daily_overrides_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE daily_overrides_id_seq OWNED BY daily_overrides.id;


--
-- TOC entry 221 (class 1259 OID 16412)
-- Name: doctor_specializations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS doctor_specializations (
                                                             doctor_id integer NOT NULL,
                                                             specialization_id integer NOT NULL
);


ALTER TABLE doctor_specializations OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 16390)
-- Name: doctors; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS doctors (
                                              user_id integer NOT NULL,
                                              first_name character varying(256) NOT NULL,
    second_name character varying(256) NOT NULL,
    surname character varying(256),
    phone_number character varying(11),
    email character varying(256) NOT NULL,
    education text,
    experience integer,
    gender character(1) NOT NULL
    );


ALTER TABLE doctors OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16389)
-- Name: doctors_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS doctors_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE doctors_id_seq OWNER TO postgres;

--
-- TOC entry 4904 (class 0 OID 0)
-- Dependencies: 217
-- Name: doctors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE doctors_id_seq OWNED BY doctors.user_id;


--
-- TOC entry 242 (class 1259 OID 24754)
-- Name: materials; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS materials (
                                                id integer NOT NULL,
                                                name character varying(256) NOT NULL,
    price integer NOT NULL
    );


ALTER TABLE materials OWNER TO postgres;

--
-- TOC entry 241 (class 1259 OID 24753)
-- Name: materials_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS materials_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE materials_id_seq OWNER TO postgres;

--
-- TOC entry 4905 (class 0 OID 0)
-- Dependencies: 241
-- Name: materials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE materials_id_seq OWNED BY materials.id;


--
-- TOC entry 230 (class 1259 OID 16480)
-- Name: patients; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS patients (
                                               user_id integer NOT NULL,
                                               first_name character varying NOT NULL,
                                               second_name character varying(256) NOT NULL,
    surname character varying(256),
    email character varying(256),
    birth_date date NOT NULL,
    phone_number character varying(11) NOT NULL,
    gender character(1) NOT NULL
    );


ALTER TABLE patients OWNER TO postgres;

--
-- TOC entry 233 (class 1259 OID 16522)
-- Name: permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS permissions (
                                                  id integer NOT NULL,
                                                  name text NOT NULL
);


ALTER TABLE permissions OWNER TO postgres;

--
-- TOC entry 232 (class 1259 OID 16521)
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS permissions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE permissions_id_seq OWNER TO postgres;

--
-- TOC entry 4906 (class 0 OID 0)
-- Dependencies: 232
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE permissions_id_seq OWNED BY permissions.id;


--
-- TOC entry 229 (class 1259 OID 16462)
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS roles (
                                            id smallint NOT NULL,
                                            name text NOT NULL
);


ALTER TABLE roles OWNER TO postgres;

--
-- TOC entry 228 (class 1259 OID 16461)
-- Name: role_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS role_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE role_id_seq OWNER TO postgres;

--
-- TOC entry 4907 (class 0 OID 0)
-- Dependencies: 228
-- Name: role_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE role_id_seq OWNED BY roles.id;


--
-- TOC entry 234 (class 1259 OID 16531)
-- Name: role_permission; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS role_permission (
                                                      role_id integer NOT NULL,
                                                      permission_id integer NOT NULL
);


ALTER TABLE role_permission OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 16428)
-- Name: service_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS service_types (
                                                    id integer NOT NULL,
                                                    name text NOT NULL
);


ALTER TABLE service_types OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 16427)
-- Name: service_types_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS service_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE service_types_id_seq OWNER TO postgres;

--
-- TOC entry 4908 (class 0 OID 0)
-- Dependencies: 222
-- Name: service_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE service_types_id_seq OWNED BY service_types.id;


--
-- TOC entry 225 (class 1259 OID 16437)
-- Name: services; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS services (
                                               id integer NOT NULL,
                                               name text NOT NULL,
                                               price integer,
                                               type integer NOT NULL
);


ALTER TABLE services OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 16436)
-- Name: services_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS services_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE services_id_seq OWNER TO postgres;

--
-- TOC entry 4909 (class 0 OID 0)
-- Dependencies: 224
-- Name: services_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE services_id_seq OWNED BY services.id;


--
-- TOC entry 220 (class 1259 OID 16399)
-- Name: specializations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS specializations (
                                                      id integer NOT NULL,
                                                      name text NOT NULL
);


ALTER TABLE specializations OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 16398)
-- Name: specializations_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS specializations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE specializations_id_seq OWNER TO postgres;

--
-- TOC entry 4910 (class 0 OID 0)
-- Dependencies: 219
-- Name: specializations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE specializations_id_seq OWNED BY specializations.id;


--
-- TOC entry 227 (class 1259 OID 16451)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS users (
                                            id integer NOT NULL,
                                            login character varying(256) NOT NULL,
    password text NOT NULL,
    role smallint NOT NULL
    );


ALTER TABLE users OWNER TO postgres;

--
-- TOC entry 226 (class 1259 OID 16450)
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE user_id_seq OWNER TO postgres;

--
-- TOC entry 4911 (class 0 OID 0)
-- Dependencies: 226
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE user_id_seq OWNED BY users.id;


--
-- TOC entry 238 (class 1259 OID 24727)
-- Name: weekly_schedule; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS weekly_schedule (
                                                      id integer NOT NULL,
                                                      doctor_id integer NOT NULL,
                                                      weekday integer NOT NULL,
                                                      start_time time without time zone NOT NULL,
                                                      end_time time without time zone NOT NULL,
                                                      slot_duration_minutes integer DEFAULT 60 NOT NULL,
                                                      CONSTRAINT weekly_schedule_weekday_check CHECK (((weekday >= 0) AND (weekday <= 6)))
    );


ALTER TABLE weekly_schedule OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 24726)
-- Name: weekly_schedule_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE IF NOT EXISTS weekly_schedule_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE weekly_schedule_id_seq OWNER TO postgres;

--
-- TOC entry 4912 (class 0 OID 0)
-- Dependencies: 237
-- Name: weekly_schedule_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE weekly_schedule_id_seq OWNED BY weekly_schedule.id;


--
-- TOC entry 4698 (class 2604 OID 24709)
-- Name: appointments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY appointments ALTER COLUMN id SET DEFAULT nextval('appointments_id_seq'::regclass);


--
-- TOC entry 4702 (class 2604 OID 24744)
-- Name: daily_overrides id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY daily_overrides ALTER COLUMN id SET DEFAULT nextval('daily_overrides_id_seq'::regclass);


--
-- TOC entry 4704 (class 2604 OID 24757)
-- Name: materials id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY materials ALTER COLUMN id SET DEFAULT nextval('materials_id_seq'::regclass);


--
-- TOC entry 4697 (class 2604 OID 16525)
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY permissions ALTER COLUMN id SET DEFAULT nextval('permissions_id_seq'::regclass);


--
-- TOC entry 4696 (class 2604 OID 16465)
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY roles ALTER COLUMN id SET DEFAULT nextval('role_id_seq'::regclass);


--
-- TOC entry 4693 (class 2604 OID 16431)
-- Name: service_types id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY service_types ALTER COLUMN id SET DEFAULT nextval('service_types_id_seq'::regclass);


--
-- TOC entry 4694 (class 2604 OID 16440)
-- Name: services id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY services ALTER COLUMN id SET DEFAULT nextval('services_id_seq'::regclass);


--
-- TOC entry 4692 (class 2604 OID 16402)
-- Name: specializations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY specializations ALTER COLUMN id SET DEFAULT nextval('specializations_id_seq'::regclass);


--
-- TOC entry 4695 (class 2604 OID 16454)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('user_id_seq'::regclass);


--
-- TOC entry 4700 (class 2604 OID 24730)
-- Name: weekly_schedule id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY weekly_schedule ALTER COLUMN id SET DEFAULT nextval('weekly_schedule_id_seq'::regclass);


--
-- TOC entry 4726 (class 2606 OID 16503)
-- Name: administrators administrators_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY administrators
    ADD CONSTRAINT administrators_pkey PRIMARY KEY (user_id);


--
-- TOC entry 4732 (class 2606 OID 24715)
-- Name: appointments appointments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY appointments
    ADD CONSTRAINT appointments_pkey PRIMARY KEY (id);


--
-- TOC entry 4736 (class 2606 OID 24747)
-- Name: daily_overrides daily_overrides_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY daily_overrides
    ADD CONSTRAINT daily_overrides_pkey PRIMARY KEY (id);


--
-- TOC entry 4712 (class 2606 OID 16416)
-- Name: doctor_specializations doctor_specializations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY doctor_specializations
    ADD CONSTRAINT doctor_specializations_pkey PRIMARY KEY (doctor_id, specialization_id);


--
-- TOC entry 4708 (class 2606 OID 16397)
-- Name: doctors doctors_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY doctors
    ADD CONSTRAINT doctors_pkey PRIMARY KEY (user_id);


--
-- TOC entry 4718 (class 2606 OID 24692)
-- Name: users login; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT login UNIQUE (login);


--
-- TOC entry 4738 (class 2606 OID 24759)
-- Name: materials materials_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY materials
    ADD CONSTRAINT materials_pkey PRIMARY KEY (id);


--
-- TOC entry 4724 (class 2606 OID 16486)
-- Name: patients patients_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY patients
    ADD CONSTRAINT patients_pkey PRIMARY KEY (user_id);


--
-- TOC entry 4728 (class 2606 OID 16529)
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- TOC entry 4730 (class 2606 OID 16535)
-- Name: role_permission role_permission_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY role_permission
    ADD CONSTRAINT role_permission_pkey PRIMARY KEY (role_id, permission_id);


--
-- TOC entry 4722 (class 2606 OID 16469)
-- Name: roles role_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY roles
    ADD CONSTRAINT role_pkey PRIMARY KEY (id);


--
-- TOC entry 4714 (class 2606 OID 16435)
-- Name: service_types service_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY service_types
    ADD CONSTRAINT service_types_pkey PRIMARY KEY (id);


--
-- TOC entry 4716 (class 2606 OID 16444)
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY services
    ADD CONSTRAINT services_pkey PRIMARY KEY (id);


--
-- TOC entry 4710 (class 2606 OID 16406)
-- Name: specializations specializations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY specializations
    ADD CONSTRAINT specializations_pkey PRIMARY KEY (id);


--
-- TOC entry 4720 (class 2606 OID 16458)
-- Name: users user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- TOC entry 4734 (class 2606 OID 24734)
-- Name: weekly_schedule weekly_schedule_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY weekly_schedule
    ADD CONSTRAINT weekly_schedule_pkey PRIMARY KEY (id);


--
-- TOC entry 4751 (class 2606 OID 24748)
-- Name: daily_overrides daily_overrides_doctor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY daily_overrides
    ADD CONSTRAINT daily_overrides_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES doctors(user_id);


--
-- TOC entry 4740 (class 2606 OID 16417)
-- Name: doctor_specializations doc_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY doctor_specializations
    ADD CONSTRAINT doc_id FOREIGN KEY (doctor_id) REFERENCES doctors(user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 4748 (class 2606 OID 24721)
-- Name: appointments doctor_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY appointments
    ADD CONSTRAINT doctor_id FOREIGN KEY (doctor_id) REFERENCES doctors(user_id);


--
-- TOC entry 4744 (class 2606 OID 16487)
-- Name: patients id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY patients
    ADD CONSTRAINT id FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 4739 (class 2606 OID 16492)
-- Name: doctors id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY doctors
    ADD CONSTRAINT id FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- TOC entry 4745 (class 2606 OID 16504)
-- Name: administrators id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY administrators
    ADD CONSTRAINT id FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 4749 (class 2606 OID 24716)
-- Name: appointments patient_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY appointments
    ADD CONSTRAINT patient_id FOREIGN KEY (patient_id) REFERENCES patients(user_id);


--
-- TOC entry 4746 (class 2606 OID 16541)
-- Name: role_permission permission; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY role_permission
    ADD CONSTRAINT permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- TOC entry 4743 (class 2606 OID 16475)
-- Name: users role; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT role FOREIGN KEY (role) REFERENCES roles(id) ON UPDATE CASCADE NOT VALID;


--
-- TOC entry 4747 (class 2606 OID 16536)
-- Name: role_permission role; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY role_permission
    ADD CONSTRAINT role FOREIGN KEY (role_id) REFERENCES roles(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 4742 (class 2606 OID 16445)
-- Name: services service_type; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY services
    ADD CONSTRAINT service_type FOREIGN KEY (type) REFERENCES service_types(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 4741 (class 2606 OID 16422)
-- Name: doctor_specializations spec_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY doctor_specializations
    ADD CONSTRAINT spec_id FOREIGN KEY (specialization_id) REFERENCES specializations(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 4750 (class 2606 OID 24735)
-- Name: weekly_schedule weekly_schedule_doctor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY weekly_schedule
    ADD CONSTRAINT weekly_schedule_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES doctors(user_id);


-- Completed on 2025-05-16 17:25:08

--
-- PostgreSQL database dump complete
--

