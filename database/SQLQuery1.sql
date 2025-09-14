

-- This script creates the database schema for the University Management System.
-- It is based on the latest ERD, using Faculty/Department terminology.

-- 1. ROLES Table
-- Defines user roles within the system (e.g., 'Student', 'Instructor', 'Admin').
CREATE TABLE ROLES (
    role_id INT PRIMARY KEY IDENTITY(1,1),
    role_name VARCHAR(50) NOT NULL UNIQUE
);

-- 2. USERS Table
-- Central table for authentication. All students and instructors will have an entry here.
CREATE TABLE USERS (
    user_id INT PRIMARY KEY IDENTITY(1,1),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id INT,
    FOREIGN KEY (role_id) REFERENCES ROLES(role_id)
);

-- 3. FACULTIES Table (Formerly DEPARTMENTS)
-- Stores large academic divisions (e.g., "Faculty of Engineering", "Faculty of Arts").
CREATE TABLE FACULTIES (
    faculty_id INT PRIMARY KEY IDENTITY(1,1),
    name VARCHAR(255) NOT NULL
);

-- 4. DEPARTMENTS Table (Formerly PROGRAMS)
-- Stores academic departments within a faculty (e.g., "Computer Science", "History").
CREATE TABLE DEPARTMENTS (
    department_id INT PRIMARY KEY IDENTITY(1,1),
    department_name VARCHAR(255) NOT NULL,
    faculty_id INT,
    FOREIGN KEY (faculty_id) REFERENCES FACULTIES(faculty_id)
);

-- 5. INSTRUCTORS Table
-- Stores profile information for instructors. Linked to a specific department.
CREATE TABLE INSTRUCTORS (
    instructor_id INT PRIMARY KEY IDENTITY(1,1),
    user_id INT NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    department_id INT,
    FOREIGN KEY (user_id) REFERENCES USERS(user_id),
    FOREIGN KEY (department_id) REFERENCES DEPARTMENTS(department_id)
);

-- 6. STUDENTS Table
-- Contains profile information for students, linked to their major/department.
CREATE TABLE STUDENTS (
    student_id INT PRIMARY KEY IDENTITY(1,1),
    user_id INT NOT NULL UNIQUE,
    department_id INT,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE,
    FOREIGN KEY (user_id) REFERENCES USERS(user_id),
    FOREIGN KEY (department_id) REFERENCES DEPARTMENTS(department_id)
);

-- 7. COURSES Table
-- The catalog of courses the university offers. Each course belongs to a DEPARTMENT.
CREATE TABLE COURSES (
    course_id INT PRIMARY KEY IDENTITY(1,1),
    course_code VARCHAR(20) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    credits INT NOT NULL,
    department_id INT,
    FOREIGN KEY (department_id) REFERENCES DEPARTMENTS(department_id)
);

-- 8. COURSE_PREREQUISITES Table
-- Defines prerequisite relationships between courses.
CREATE TABLE COURSE_PREREQUISITES (
    course_id INT NOT NULL,
    prerequisite_id INT NOT NULL,
    PRIMARY KEY (course_id, prerequisite_id),
    FOREIGN KEY (course_id) REFERENCES COURSES(course_id),
    FOREIGN KEY (prerequisite_id) REFERENCES COURSES(course_id)
);

-- 9. BOOKS Table
-- A catalog of all textbooks.
CREATE TABLE BOOKS (
    book_id INT PRIMARY KEY IDENTITY(1,1),
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255),
    isbn VARCHAR(20) UNIQUE,
    publisher VARCHAR(100),
    publication_year INT
);

-- 10. COURSE_BOOKS Table
-- Links courses to their required or recommended books.
CREATE TABLE COURSE_BOOKS (
    course_id INT NOT NULL,
    book_id INT NOT NULL,
    usage_type VARCHAR(50) NOT NULL, -- e.g., 'Required', 'Recommended'
    PRIMARY KEY (course_id, book_id),
    FOREIGN KEY (course_id) REFERENCES COURSES(course_id),
    FOREIGN KEY (book_id) REFERENCES BOOKS(book_id)
);

-- 11. COURSE_OFFERINGS Table
-- Represents a specific instance of a course being taught in a given semester and year.
CREATE TABLE COURSE_OFFERINGS (
    offering_id INT PRIMARY KEY IDENTITY(1,1),
    course_id INT NOT NULL,
    instructor_id INT NOT NULL,
    semester VARCHAR(50) NOT NULL,
    year INT NOT NULL,
    capacity INT,
    FOREIGN KEY (course_id) REFERENCES COURSES(course_id),
    FOREIGN KEY (instructor_id) REFERENCES INSTRUCTORS(instructor_id)
);

-- 12. ENROLLMENTS Table
-- Connects STUDENTS and COURSE_OFFERINGS, representing a student's enrollment in a course.
CREATE TABLE ENROLLMENTS (
    student_id INT NOT NULL,
    offering_id INT NOT NULL,
    grade VARCHAR(5),
    PRIMARY KEY (student_id, offering_id), -- A student can only enroll in the same offering once
    FOREIGN KEY (student_id) REFERENCES STUDENTS(student_id),
    FOREIGN KEY (offering_id) REFERENCES COURSE_OFFERINGS(offering_id)
);

-- 13. CLASS_SESSIONS Table
-- Defines the actual timetable (day, time, location) for a course offering.
CREATE TABLE CLASS_SESSIONS (
    session_id INT PRIMARY KEY IDENTITY(1,1),
    offering_id INT NOT NULL,
    day_of_week VARCHAR(20) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    location VARCHAR(100),
    FOREIGN KEY (offering_id) REFERENCES COURSE_OFFERINGS(offering_id)
);

