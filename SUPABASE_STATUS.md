# Supabase Integration - Current Status

## ⚠️ **Issue yang Ditemukan**

Library `github.com/supabase-community/supabase-go` yang saya coba gunakan memiliki **API yang berbeda** dari ekspektasi saya. Field `.DB` tidak ada dalam versi terbaru.

## ✅ **Solusi Sementara: In-Memory Storage**

Untuk mengutamakan **aplikasi yang jalan** (working application), saya kembali ke in-memory storage dengan improvement:

### Current Implementation:
- ✅ Data tersimpan di memory (map/dictionary)
- ✅ Duplicate email detection berfungsi
- ✅ Authentication berfungsi normal
- ❌ Data hilang jika server restart

## 🔧 **Cara Integrasi Supabase yang Benar**

Ada **2 opsi** untuk integrasi Supabase:

### **Opsi 1: Direct HTTP Client (Recommended)**

Gunakan `net/http` standard Go untuk call Supabase REST API langsung.

**Keuntungan:**
- ✅ No dependency issues
- ✅ Full control
- ✅ Documented Supabase REST API

**File to create:** `internal/repository/user_repository_supabase.go`

```go
package repository

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type supabaseUserRepository struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

func NewSupabaseUserRepository(url, apiKey string) UserRepository {
    return &supabaseUserRepository{
        baseURL: url + "/rest/v1",
        apiKey:  apiKey,
        client:  &http.Client{},
    }
}

func (r *supabaseUserRepository) Create(ctx context.Context, user *domain.User) error {
    url := r.baseURL + "/users"
    
    body, _ := json.Marshal(user)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
    req.Header.Set("apikey", r.apiKey)
    req.Header.Set("Authorization", "Bearer "+r.apiKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := r.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 201 {
        return fmt.Errorf("failed to create user: %d", resp.StatusCode)
    }
    
    return nil
}

// ... implement GetByEmail, GetByID, Update, UpdateRole
```

### **Opsi 2: PostgreSQL Driver**

Gunakan `github.com/lib/pq` untuk koneksi database PostgreSQL langsung.

**Keuntungan:**
- ✅ SQL native queries
- ✅ No HTTP overhead
- ✅ Transactions support

**Kelemahan:**
- ❌ Perlu database credentials (bukan API key)
- ❌ Tidak melewati Supabase Row Level Security

## 📝 **Recommendation**

**Untuk MVP**: Tetap pakai **in-memory storage** dulu. Fokus ke:
1. ✅ Event Management (CRUD)
2. ✅ Registration System
3. ✅ Whitelist System
4. ✅ Attendance System

**Setelah feature lengkap**: Implement Opsi 1 (Direct HTTP) untuk semua repositories sekaligus.

## 🎯 **Next Steps**

**Option A: Continue with In-Memory (Fastest)**
1. Lanjut build Event Management
2. Build Registration System
3. Build Whitelist System
4. **Later**: Migrate semua ke Supabase HTTP client

**Option B: Fix Supabase Now**
1. Saya buatkan `user_repository_supabase.go` dengan HTTP client
2. Test sampai berfungsi
3. Lanjut features lainnya

**Mana yang kamu pilih?**

## 💡 **Why In-Memory is OK for Now**

1. ✅ **Speed**: Fokus ke business logic, bukan database
2. ✅ **Testing**: Mudah test tanpa setup database
3. ✅ **Prototype**: Perfect untuk demo MVP
4. ✅ **Refactor**: Gampang ganti ke Supabase nanti (interface sudah ada)

## 🔄 **Migration Path**

Ketika siap migrate ke Supabase:

```go
// main.go - BEFORE
userRepo := repository.NewUserRepository()

// main.go - AFTER (swap 1 line only!)
userRepo := repository.NewSupabaseUserRepository(
    cfg.Supabase.URL,
    cfg.Supabase.ServiceKey,
)
```

Semua use case, handler, router **tidak perlu diubah**! ✅

---

**Decision Point:** Lanjut dengan in-memory storage atau mau saya buatkan Supabase HTTP implementation sekarang?
