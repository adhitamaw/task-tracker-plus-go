# 📋 Task Tracker Plus - Portfolio Documentation

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=JSON%20web%20tokens&logoColor=white)
![BBolt](https://img.shields.io/badge/BBolt-DB-blue?style=for-the-badge)

## 🎯 Project Overview

**Task Tracker Plus** adalah aplikasi web manajemen tugas (task management) yang dibangun dengan Go (Golang). Aplikasi ini dirancang untuk membantu pengguna dalam mengorganisir tugas-tugas mereka dengan sistem kategori, prioritas, dan deadline. Project ini mengimplementasikan konsep **Monolithic Architecture** dengan REST API dan pola desain MVC (Model-View-Controller).

### ✨ Key Highlights
- 🔐 **Secure Authentication** - JWT-based authentication dengan session management
- 🎨 **Clean Architecture** - Implementasi Repository Pattern dan Service Layer
- 📊 **User Isolation** - Setiap user memiliki data yang terisolasi (categories & tasks)
- 🗄️ **File-based Database** - Menggunakan BBolt untuk embedded database
- 🚀 **RESTful API** - API endpoints yang lengkap dan terstruktur
- 🌐 **Full Stack** - Backend API dan Frontend Web Client dalam satu aplikasi

---

## 🏗️ Architecture & Design Patterns

### Monolithic Architecture
```
┌─────────────────────────────────────────────────┐
│          Task Tracker Plus Application          │
│                                                  │
│  ┌────────────┐           ┌─────────────┐      │
│  │  Frontend  │◄──────────┤   Backend   │      │
│  │  (Views)   │           │  (REST API) │      │
│  └────────────┘           └─────────────┘      │
│         │                        │              │
│         └────────┬───────────────┘              │
│                  │                              │
│            ┌─────▼──────┐                       │
│            │  Database  │                       │
│            │   (BBolt)  │                       │
│            └────────────┘                       │
└─────────────────────────────────────────────────┘
```

### Layer Architecture
```
┌─────────────────────────────────────────┐
│         Handler Layer (API/Web)          │
│  • REST API Endpoints                    │
│  • Web Page Handlers                     │
│  • Request/Response Processing           │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│          Middleware Layer                │
│  • Authentication (JWT)                  │
│  • Session Validation                    │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│          Service Layer                   │
│  • Business Logic                        │
│  • Data Validation                       │
│  • User Isolation Logic                  │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│        Repository Layer                  │
│  • Data Access Logic                     │
│  • Database Operations                   │
│  • Query Execution                       │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│        Database Layer (BBolt)            │
│  • File-based Storage                    │
│  • ACID Transactions                     │
└─────────────────────────────────────────┘
```

---

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.24
- **Framework**: Gin Web Framework v1.9.0
- **Database**: BBolt (embedded key-value database)
- **Authentication**: JWT (JSON Web Token) - golang-jwt/jwt v3.2.2
- **Password Hashing**: bcrypt (golang.org/x/crypto)

### Frontend
- **Template Engine**: HTML Templates (embedded)
- **CSS Framework**: Tailwind CSS (utility-first CSS)
- **Static Assets**: Served via Gin StaticFS
- **Icons & Images**: SVG assets

### Testing
- **Testing Framework**: Ginkgo v2.1.4
- **Assertion Library**: Gomega v1.19.0

### Additional Libraries
- **HTML Parsing**: goquery v1.8.1 (PuerkitoBio/goquery)
- **Database Driver**: lib/pq v1.10.7 (PostgreSQL driver for potential migration)

---

## 📁 Project Structure

```
fcp-web-application-v3/
│
├── 📂 handler/              # HTTP Request Handlers
│   ├── api/                # REST API Handlers
│   │   ├── user.go        # User API (register, login)
│   │   ├── task.go        # Task CRUD API
│   │   └── category.go    # Category CRUD API
│   │
│   └── web/                # Web Page Handlers
│       ├── auth.go        # Login, Register pages
│       ├── dashboard.go   # Dashboard page
│       ├── task.go        # Task management page
│       ├── category.go    # Category management page
│       ├── home.go        # Landing page
│       └── modals.go      # Modal components
│
├── 📂 service/             # Business Logic Layer
│   ├── user.go           # User business logic
│   ├── task.go           # Task business logic
│   ├── category.go       # Category business logic
│   └── session.go        # Session management
│
├── 📂 repository/          # Data Access Layer
│   ├── user.go           # User data operations
│   ├── task.go           # Task data operations
│   ├── category.go       # Category data operations
│   └── session.go        # Session data operations
│
├── 📂 middleware/          # HTTP Middleware
│   └── auth.go           # JWT Authentication middleware
│
├── 📂 model/              # Data Models
│   ├── model.go          # Core models (User, Task, Category)
│   ├── jwt.go            # JWT claims & config
│   └── response.go       # API response models
│
├── 📂 client/             # HTTP Client (for web handlers)
│   ├── user.go           # User API client
│   ├── task.go           # Task API client
│   └── category.go       # Category API client
│
├── 📂 db/filebased/       # Database Implementation
│   ├── filebased.go      # BBolt database operations
│   └── README.md         # Database documentation
│
├── 📂 views/              # HTML Templates
│   ├── auth/             # Login & Register templates
│   ├── main/             # Dashboard & main pages
│   ├── general/          # Shared templates
│   └── modals/           # Modal templates
│
├── 📂 assets/             # Static Assets
│   ├── avatars/
│   └── icons/
│
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
└── README.md             # Project documentation
```

---

## 🔑 Core Features

### 1. 🔐 User Authentication & Authorization

#### Registration
- Validasi input (fullname, email, password)
- Hashing password menggunakan bcrypt
- Pengecekan email duplikat
- Auto-create default categories untuk user baru

#### Login
- Validasi credentials
- Generate JWT token dengan expiry 1 jam
- Store session di database
- Set HTTP-only cookie untuk security

#### Session Management
- JWT-based authentication
- Token expiry management
- Session cleanup untuk expired tokens
- Middleware untuk protected routes

**Security Features:**
```go
// Password hashing dengan bcrypt
bcrypt.GenerateFromPassword([]byte(password), 8)

// JWT token dengan claims
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// HTTP-only cookie (prevent XSS)
http.Cookie{
    Name:     "session_token",
    Value:    token,
    HttpOnly: true,
}
```

### 2. 📝 Task Management

#### CRUD Operations
- **Create**: Tambah task dengan title, deadline, priority, status, category
- **Read**: List semua tasks dan get task by ID
- **Update**: Edit informasi task
- **Delete**: Hapus task

#### Task Features
- **Priority Levels**: 
  - 1 (🟢 Low) - Prioritas rendah
  - 2 (🟡 Medium) - Prioritas sedang
  - 3 (🔴 High) - Prioritas tinggi
  - Dipilih via dropdown selector dengan visual indicator
- **Status**: 
  - Not Started - Belum dimulai
  - In Progress - Sedang dikerjakan
  - Completed - Selesai
  - Dipilih via dropdown selector
- **Deadline**: Date picker untuk tanggal target penyelesaian
- **Category Association**: Task terkait dengan category via dropdown
- **User Isolation**: User hanya bisa melihat task miliknya sendiri

#### User Interface
Form input task menggunakan komponen modern:
- **Text Input** untuk Title
- **Date Picker** untuk Deadline
- **Dropdown Selector** untuk Priority (dengan emoji visual indicator)
- **Dropdown Selector** untuk Status
- **Dropdown Selector** untuk Category

Ini memberikan user experience yang lebih baik dibanding free text input, mengurangi error input, dan memberikan visual guidance yang jelas.

#### API Endpoints
```
POST   /api/v1/task/add              - Create new task
GET    /api/v1/task/get/:id          - Get task by ID
PUT    /api/v1/task/update/:id       - Update task
DELETE /api/v1/task/delete/:id       - Delete task
GET    /api/v1/task/list             - Get all tasks (by user)
GET    /api/v1/task/category/:id     - Get tasks by category
```

### 3. 🏷️ Category Management

#### CRUD Operations
- **Create**: Tambah category baru
- **Read**: List categories dan get by ID
- **Update**: Edit nama category
- **Delete**: Hapus category (cascade delete tasks)

#### Category Features
- User-specific categories
- Cascade delete untuk terkait tasks
- Auto-create default categories (Work, Personal, Study)

#### API Endpoints
```
POST   /api/v1/category/add          - Create new category
GET    /api/v1/category/get/:id      - Get category by ID
PUT    /api/v1/category/update/:id   - Update category
DELETE /api/v1/category/delete/:id   - Delete category
GET    /api/v1/category/list         - Get all categories (by user)
```

### 4. 🌐 Web Interface

#### Pages
- **Landing Page** (`/`) - Clean homepage dengan branding
- **Login** (`/client/login`) - Secure login form dengan validation
- **Register** (`/client/register`) - User registration dengan password hashing
- **Dashboard** (`/client/dashboard`) - Comprehensive overview:
  - User account information
  - Total task count statistics
  - Task list dengan Category, Deadline, Priority badges, dan Status
  - Quick action button untuk add new task
- **Tasks** (`/client/task`) - Full task management interface
- **Categories** (`/client/category`) - Category organization interface

#### UI/UX Features
- **Responsive design** - Mobile-friendly layout
- **Dynamic content rendering** - Real-time data display
- **Form validation** - Client & server-side validation
- **Modal dialogs** - Smooth interactions
- **Session-based navigation** - Protected routes
- **Dropdown selectors** - User-friendly input untuk Priority, Status, dan Category
- **Visual indicators** - Emoji & color coding untuk priority levels
- **Date picker** - Native date selection untuk deadlines
- **Consistent styling** - Tailwind CSS untuk modern UI

#### Form Design Philosophy
Aplikasi menggunakan **constrained input** pattern untuk mengurangi user error:
- ✅ Dropdown untuk pilihan terbatas (Priority, Status, Category)
- ✅ Date picker untuk format tanggal konsisten
- ✅ Clear labeling dan placeholder text
- ✅ Required field validation
- ❌ Menghindari free text input untuk data terstruktur

---

## 🎨 Database Design

### BBolt Buckets (Tables)

#### 1. Users Bucket
```json
{
  "id": 1,
  "fullname": "John Doe",
  "email": "john@example.com",
  "password": "hashed_password",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

#### 2. Tasks Bucket
```json
{
  "id": 1,
  "title": "Complete Project",
  "deadline": "2026-12-31",
  "priority": 2,
  "status": "In Progress",
  "category_id": 1,
  "user_id": 1
}
```

#### 3. Categories Bucket
```json
{
  "id": 1,
  "name": "Work",
  "user_id": 1
}
```

#### 4. Sessions Bucket
```json
{
  "id": 1,
  "token": "jwt_token_string",
  "email": "john@example.com",
  "expiry": "2026-01-01T01:00:00Z"
}
```

### Data Relationships
```
User (1) ──┬── (N) Categories
           │
           └── (N) Tasks
                    │
Category (1) ────── (N) Tasks
```

---

## 🔄 API Reference

### Authentication Required: 🔒

All endpoints except `/user/register`, `/user/login`, and `/user/list` require authentication via JWT token in cookie.

### User API

#### POST `/api/v1/user/register`
Register new user
```json
// Request
{
  "fullname": "John Doe",
  "email": "john@example.com",
  "password": "Pass123!"
}

// Response (201)
{
  "message": "register success"
}
```

#### POST `/api/v1/user/login`
Login user
```json
// Request
{
  "email": "john@example.com",
  "password": "Pass123!"
}

// Response (200)
{
  "user_id": 1,
  "message": "login success"
}
// + Set cookie: session_token
```

#### GET `/api/v1/user/list` 
Debug endpoint - List all users (no auth required)

#### GET `/api/v1/user/tasks` 🔒
Get user's tasks with categories
```json
// Response (200)
[
  {
    "id": 1,
    "fullname": "John Doe",
    "email": "john@example.com",
    "task": "Complete Project",
    "deadline": "2026-12-31",
    "priority": 2,
    "status": "In Progress",
    "category": "Work"
  }
]
```

### Task API

#### POST `/api/v1/task/add` 🔒
Create new task
```json
// Request
{
  "title": "Complete Project",
  "deadline": "2026-12-31",
  "priority": 2,              // 1=Low, 2=Medium, 3=High
  "status": "Not Started",    // "Not Started" | "In Progress" | "Completed"
  "category_id": 1
}

// Response (201)
{
  "user_id": 1,
  "message": "add task success"
}
```

**Priority Values:**
- `1` - 🟢 Low Priority
- `2` - 🟡 Medium Priority  
- `3` - 🔴 High Priority

**Status Options:**
- `"Not Started"` - Belum dimulai
- `"In Progress"` - Sedang dikerjakan
- `"Completed"` - Selesai

#### GET `/api/v1/task/get/:id` 🔒
Get task by ID

#### PUT `/api/v1/task/update/:id` 🔒
Update task

#### DELETE `/api/v1/task/delete/:id` 🔒
Delete task

#### GET `/api/v1/task/list` 🔒
Get all user's tasks

#### GET `/api/v1/task/category/:id` 🔒
Get tasks by category ID

### Category API

#### POST `/api/v1/category/add` 🔒
Create new category
```json
// Request
{
  "name": "Personal"
}

// Response (201)
{
  "user_id": 1,
  "message": "add category success"
}
```

#### GET `/api/v1/category/get/:id` 🔒
Get category by ID

#### PUT `/api/v1/category/update/:id` 🔒
Update category

#### DELETE `/api/v1/category/delete/:id` 🔒
Delete category

#### GET `/api/v1/category/list` 🔒
Get all user's categories

---

## 🚀 Getting Started

### Prerequisites
- Go 1.24 or higher
- Git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/adhitamaw/task-tracker-plus-go.git
cd task-tracker-plus-go
```

2. **Install dependencies**
```bash
go mod download
```

3. **Run the application**
```bash
go run main.go
```

4. **Access the application**
- Web Interface: `http://localhost:8080`
- API Base URL: `http://localhost:8080/api/v1`

### Environment Variables (Optional)
```bash
# Custom database path (default: file.db)
export APP_DB_PATH="custom_path/file.db"
```

---

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test file
go test -v main_test.go
```

### Test Coverage
The project includes comprehensive tests using Ginkgo and Gomega:
- Unit tests for services
- Integration tests for repositories
- API endpoint tests
- Authentication middleware tests

---

## 🔐 Security Implementations

### 1. Password Security
- **Bcrypt hashing** dengan cost factor 8
- Plain passwords tidak pernah disimpan
- Password hashing otomatis saat register

### 2. Authentication
- **JWT tokens** dengan expiry 1 jam
- **HTTP-only cookies** untuk mencegah XSS attacks
- Token validation di setiap protected endpoint

### 3. Authorization
- **User data isolation** - User hanya bisa akses data miliknya
- **Middleware authentication** untuk protected routes
- **Session validation** setiap request

### 4. Input Validation
- Validasi input di handler layer
- Binding validation menggunakan Gin
- Sanitasi data sebelum database operations

### 5. Database Security
- **Transaction support** via BBolt
- **ACID compliance** untuk data consistency
- **File permissions** (0600) untuk database file

---

## 📊 Key Technical Decisions

### 1. Monolithic Architecture
**Why?** Simplicity untuk project scale ini, easier deployment, dan semua components dalam satu codebase.

### 2. BBolt Database
**Why?** 
- Embedded database (no separate server needed)
- ACID compliant
- File-based storage
- Perfect untuk development dan mid-scale applications
- Mudah backup (single file)

### 3. Repository Pattern
**Why?**
- Separation of concerns
- Easier testing dengan mock repositories
- Flexible untuk database migration
- Clean code architecture

### 4. JWT Authentication
**Why?**
- Stateless authentication
- Scalable untuk distributed systems
- Industry standard
- Secure dengan proper implementation

### 5. Service Layer
**Why?**
- Business logic isolation
- Reusability
- Easier to test
- Clear separation between data access dan business rules

---

## 🎯 Project Achievements

### ✅ Technical Implementation
- [x] Clean Architecture dengan separation of concerns
- [x] RESTful API design principles
- [x] JWT-based authentication system
- [x] User data isolation dan multi-tenancy
- [x] Transaction support untuk data consistency
- [x] Comprehensive error handling
- [x] Middleware untuk cross-cutting concerns
- [x] Responsive web interface
- [x] Unit dan integration testing

### ✅ Features Implemented
- [x] User registration dan authentication
- [x] Task CRUD operations dengan categories
- [x] Category management
- [x] Session management
- [x] Dashboard dengan task overview dan priority display
- [x] Priority management dengan visual indicators (emoji & color-coded badges)
- [x] Status management dengan color-coded badges
- [x] Deadline tracking dengan date picker
- [x] User-specific data isolation
- [x] Responsive mobile-friendly interface

---

## 🔮 Future Enhancements

### Potential Features
- [ ] Task search dan filtering
- [ ] Task sorting (by priority, deadline, status)
- [ ] Task attachments
- [ ] Task comments dan collaboration
- [ ] Email notifications untuk deadlines
- [ ] Task recurring/scheduling
- [ ] Data export (CSV, PDF)
- [ ] Mobile responsive improvements
- [ ] Real-time updates dengan WebSocket
- [ ] Task statistics dan analytics

### Technical Improvements
- [ ] Migration ke PostgreSQL untuk scalability
- [ ] Redis untuk session caching
- [ ] Docker containerization
- [ ] CI/CD pipeline
- [ ] API rate limiting
- [ ] Request logging dan monitoring
- [ ] API documentation dengan Swagger
- [ ] Frontend framework (React/Vue)
- [ ] Microservices architecture

---

## 📝 Lessons Learned

### 1. Clean Architecture Benefits
Implementasi layer architecture membuat code lebih maintainable dan testable. Separation of concerns memudahkan untuk modify satu layer tanpa affecting yang lain.

### 2. Repository Pattern Value
Repository pattern sangat helpful untuk abstraction database operations. Ini membuat code lebih flexible untuk future database migrations.

### 3. JWT vs Session-based Auth
JWT memberikan stateless authentication yang scalable, tapi perlu careful implementation untuk security (expiry, refresh tokens, etc).

### 4. User Data Isolation
Implementing proper user data isolation dari awal sangat penting untuk security dan privacy. Setiap query harus include user ID filter.

### 5. Testing Importance
Comprehensive testing dengan Ginkgo/Gomega membantu catch bugs early dan gives confidence saat refactoring.

---

## 🤝 Contributing

This is a personal portfolio project created as part of Ruangguru Camp Final Project.

Feel free to explore, learn from, or fork this project for educational purposes!

### Get in Touch
If you have questions or suggestions about this project, feel free to reach out through:
- Open an issue on GitHub repository
- Connect via LinkedIn
- Email for professional inquiries

---

## 📄 License

This project is for educational dan portfolio purposes.

---

## 🙏 Acknowledgments

- **Ruangguru Camp** - Untuk program pembelajaran dan guidance
- **Go Community** - Untuk excellent documentation dan libraries
- **Gin Framework** - Untuk powerful web framework
- **BBolt** - Untuk reliable embedded database

---

## 📸 Screenshots

To see the application in action:

1. **Landing Page** (`/`) - Clean homepage with navigation to login/register
2. **Authentication Pages** - Modern login and registration forms with validation
3. **Dashboard** (`/client/dashboard`) - Overview showing:
   - User account info
   - Total task count
   - Task list with Title, Category, Deadline, Priority (🟢🟡🔴), and Status
   - Color-coded badges untuk quick visual recognition
4. **Task Management** (`/client/task`) - Full CRUD interface dengan:
   - Form untuk add task dengan dropdown selectors
   - List semua tasks dengan edit/delete actions
   - Visual priority indicators
5. **Category Management** (`/client/category`) - Manage categories dengan add/delete functionality

*Run the application locally to see the full UI with Tailwind CSS styling!*

---

## 🎓 Technical Skills Demonstrated

This project showcases proficiency in:

### Backend Development
- ✅ Go programming language
- ✅ RESTful API design
- ✅ Web framework (Gin)
- ✅ Authentication & Authorization (JWT)
- ✅ Database design dan operations
- ✅ Clean Architecture patterns

### Software Engineering
- ✅ Design patterns (Repository, Service Layer)
- ✅ Dependency injection
- ✅ Error handling
- ✅ Testing (unit & integration)
- ✅ Code organization dan structure

### Security
- ✅ Password hashing (bcrypt)
- ✅ JWT implementation
- ✅ Session management
- ✅ HTTP security headers
- ✅ Data isolation

### DevOps
- ✅ Go modules management
- ✅ Environment configuration
- ✅ Database migrations
- ✅ Build dan deployment

---

<div align="center">

### ⭐ Thank you for reviewing my project! ⭐

**Built with ❤️ using Go**

</div>
